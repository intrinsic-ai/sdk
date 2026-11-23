// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as three from 'three';

import * as colors from './colors';
import {DEFAULT_EXPOSURE} from './renderer';

/** Layers for use in OutlineRenderer. Layers control what objects get rendered */
export enum RenderLayers {
  NO_OUTLINE = 0,
  OUTLINE = 1,
}

const EDGE_KERNEL_VERTEX_SHADER = `
varying vec2 vUv;

void main() {
  vUv = uv;
  gl_Position = vec4( position, 1.0 );
}`;

// A simpler modified sobel-ish shader that sums the total neighboring pixels
// that don't match the center of the kernel; rather than computing a gradient.
const EDGE_KERNEL_FRAGMENT_SHADER = `
uniform sampler2D tDepth;
uniform vec3 uOutlineColor;
varying vec2 vUv;
out vec4 pc_FragColor;

void main() {
  vec2 invSize = 1.0 / vec2(textureSize(tDepth, 0));
  vec2 texelCoord = vec2(gl_FragCoord.xy);

  // Use texture2D instead of texelFetch for get sampler wrapping support
  float center = texture2D(tDepth, texelCoord * invSize).r;

  // get the center of the kernel
  float ref = float(center == 1.0);

  // A very naive 3x3 image kernel. Larger kernels should be smarter
  int edge = 0;
  for (int i=-1; i <= 1; i++) {
    for (int j=-1; j <= 1; j++) {
      float texel = texture2D(tDepth, (texelCoord + vec2(i, j)) * invSize).r;
      // sum the texels that don't match
      edge += int(floor(texel) != ref);
    }
  }

  float fill = float(center < 1.0);
  pc_FragColor = vec4(uOutlineColor, float(edge) * 0.25 + fill * 0.125);
}`;

type CompositeMaterialProperties = three.MeshBasicMaterialParameters & {
  depthMap: three.Texture;
};

class CompositeMaterial extends three.MeshBasicMaterial {
  depthMap: three.Texture;

  constructor(args: CompositeMaterialProperties) {
    const {depthMap, ...baseArgs} = args;
    super(baseArgs);
    this.depthMap = depthMap;
  }

  override onBeforeCompile(shader: three.ShaderLibShader) {
    shader.uniforms['tDepth'] = {value: this.depthMap};

    shader.fragmentShader = shader.fragmentShader.replace(
      '#include <common>',
      ['#include <common>', 'uniform sampler2D tDepth;'].join('\n'),
    );

    shader.fragmentShader = shader.fragmentShader.replace(
      '#include <opaque_fragment>',
      [
        '#include <opaque_fragment>',
        'gl_FragDepth = texelFetch(tDepth, ivec2(gl_FragCoord.xy), 0).r;',
      ].join('\n'),
    );
  }
}

/**
 * A custom outline renderer that implements our own faster outline shader.
 */
export class OutlineRenderer {
  private readonly renderer: three.WebGLRenderer;
  private readonly mainTarget: three.WebGLRenderTarget;
  private readonly outlineScene = new three.Scene();
  private readonly copyScene = new three.Scene();
  private readonly outlineMaterial: three.ShaderMaterial;
  private clearColor = colors.VIEWER_BACKGROUND_COLOR_LIGHT;
  private clearAlpha = 1;
  // Size of the drawing buffer in CSS pixels.
  private readonly capacity = new three.Vector2();
  // Device pixel ratio that the drawing buffer was allocated for.
  private capacityPixelRatio = 0;

  private readonly orthoCamera = new three.OrthographicCamera(
    -1,
    1,
    1,
    -1,
    0,
    1,
  );

  constructor(
    private readonly canvas: HTMLCanvasElement | OffscreenCanvas,
    failIfMajorPerformanceCaveat = false,
  ) {
    this.renderer = new three.WebGLRenderer({
      alpha: true,
      antialias: true,
      canvas: this.canvas,
      preserveDrawingBuffer: true,
      failIfMajorPerformanceCaveat,
    });
    this.renderer.localClippingEnabled = true;
    this.renderer.outputColorSpace = three.SRGBColorSpace;
    this.renderer.setClearColor(colors.VIEWER_BACKGROUND_COLOR_LIGHT);
    this.renderer.toneMapping = three.ACESFilmicToneMapping;
    this.renderer.toneMappingExposure = DEFAULT_EXPOSURE;
    // We have to manage our own clears.
    this.renderer.autoClear = false;
    this.renderer.setPixelRatio(window.devicePixelRatio);

    const size = this.renderer.getSize(new three.Vector2());

    this.mainTarget = this.maybeCreateMSAARenderTarget(size.x, size.y);
    const depthTexture = this.mainTarget.depthTexture;
    if (!depthTexture) {
      throw new Error('OutlineRenderer: mainTarget.depthTexture is null');
    }

    this.outlineMaterial = new three.ShaderMaterial({
      vertexShader: EDGE_KERNEL_VERTEX_SHADER,
      fragmentShader: EDGE_KERNEL_FRAGMENT_SHADER,
      uniforms: {
        'tDepth': {value: depthTexture},
        'uOutlineColor': {value: colors.PRIMARY_LIGHT},
      },
      glslVersion: three.GLSL3,
      toneMapped: false,
      transparent: true,
      blending: three.NormalBlending,
      depthWrite: false,
      depthTest: false,
    });

    const plane = new three.PlaneGeometry(2, 2);

    this.outlineScene.add(
      new three.Group().add(new three.Mesh(plane, this.outlineMaterial)),
    );

    this.copyScene.add(
      new three.Mesh(
        plane,
        new CompositeMaterial({
          map: this.mainTarget.texture,
          depthMap: depthTexture,
          transparent: true,
          toneMapped: true,
          depthWrite: true,
          depthTest: true,
        }),
      ),
    );

    // Sizes the offscreen target to the drawing buffer.
    this.setSize(size.x, size.y);

    // Handle WebGL context loss and restore events.
    this.canvas.addEventListener('contextlost', (event: Event) => {
      event.preventDefault();
      setTimeout(() => {
        this.renderer.forceContextRestore();
      }, 3000);
    });
  }

  setAnimationLoop(fn: (() => void) | null) {
    this.renderer.setAnimationLoop(fn);
  }

  logRenderingStatsPeriodically() {
    // Logs stats on overall memory use and renderer-specific stats every couple seconds, for debugging
    setInterval(() => {
      // Overall JS stats (note - perf.memory is not available on all browsers)
      let stats_msg = '';
      const perf = performance as Performance & {
        memory?: {usedJSHeapSize: number; jsHeapSizeLimit: number};
      };
      if (perf.memory) {
        const {usedJSHeapSize, jsHeapSizeLimit} = perf.memory;
        const usedSizeMB = (usedJSHeapSize / 1024 / 1024).toFixed(0);
        const limitSizeMB = (jsHeapSizeLimit / 1024 / 1024).toFixed(0);
        stats_msg += `Memory Used: ${usedSizeMB} MB (limit=${limitSizeMB} MB)\n`;
      }

      // Renderer stats
      const {memory, programs} = this.renderer.info;
      stats_msg += `Render Stats: geos=${memory.geometries}, textures=${memory.textures}, programs=${programs ? programs.length : 'NA'}`;
      console.log(stats_msg);
    }, 3000);
  }

  /**
   * Resizes the drawing buffer to `width` x `height` CSS pixels at the current
   * device pixel ratio.
   */
  setSize(width: number, height: number, updateStyle = false) {
    width = Math.max(1, width);
    height = Math.max(1, height);
    const ratio = window.devicePixelRatio;
    this.renderer.setPixelRatio(ratio);
    this.renderer.setSize(width, height, updateStyle);
    const buffer = this.renderer.getDrawingBufferSize(new three.Vector2());
    this.mainTarget.setSize(buffer.x, buffer.y);
    this.capacity.set(width, height);
    this.capacityPixelRatio = ratio;
  }

  /**
   * Grows the drawing buffer to fit the given CSS pixel size, if needed. The
   * buffer never shrinks, so that it can serve the largest client, and it is
   * reallocated when the device pixel ratio changes.
   */
  ensureCapacity(width: number, height: number) {
    if (
      width <= this.capacity.x &&
      height <= this.capacity.y &&
      this.capacityPixelRatio === window.devicePixelRatio
    ) {
      return;
    }
    this.setSize(
      Math.max(width, this.capacity.x),
      Math.max(height, this.capacity.y),
    );
  }

  // Leak the abstraction for async env map loading.
  getThreeRenderer(): three.WebGLRenderer {
    return this.renderer;
  }

  setClearColor(color: three.Color, alpha: number) {
    this.clearColor = color.clone();
    this.clearAlpha = alpha;
  }

  setScissor(x: number, y: number, width: number, height: number) {
    width = Math.max(1, width);
    height = Math.max(1, height);
    this.mainTarget.scissor
      .set(x, y, width, height)
      .multiplyScalar(window.devicePixelRatio);
    // pixel ratio is already handled elsewhere.
    this.renderer.setScissor(x, y, width, height);
  }

  render(scene: three.Scene, camera: three.Camera) {
    // Draw only the objects to be outlined
    // We traverse the camera in case there are attached lights
    camera.traverse((n) => {
      n.layers.set(RenderLayers.OUTLINE);
    });

    // Render the outlined object to a secondary buffer
    this.renderer.setRenderTarget(this.mainTarget);
    this.renderer.setClearColor(this.clearColor, this.clearAlpha);
    this.renderer.clear(true, true, false);
    this.renderer.render(scene, camera);

    // Do a poor man's blit to the main buffer (three.js doesn't support blits???)
    this.renderer.setRenderTarget(null);
    this.renderer.clear(true, true, false);
    this.renderer.render(this.copyScene, this.orthoCamera);

    // Draw the rest of the objects
    camera.traverse((n) => {
      n.layers.set(RenderLayers.NO_OUTLINE);
    });
    // Render the rest of the scene
    this.renderer.setRenderTarget(null);
    this.renderer.render(scene, camera);
    // Draw the outline
    this.renderer.render(this.outlineScene, this.orthoCamera);
  }

  private maybeCreateMSAARenderTarget(
    width: number,
    height: number,
  ): three.WebGLRenderTarget {
    width = Math.max(1, width);
    height = Math.max(1, height);
    const options = {
      depthTexture: new three.DepthTexture(
        width,
        height,
        three.UnsignedIntType,
      ),
    };

    return new three.WebGLRenderTarget(width, height, options);
  }
}
