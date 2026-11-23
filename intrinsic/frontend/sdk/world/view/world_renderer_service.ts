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
import {RoomEnvironment} from 'three/examples/jsm/environments/RoomEnvironment.js';
import {HDRLoader} from 'three/examples/jsm/loaders/HDRLoader.js';

import * as colors from './colors';
import {OutlineRenderer} from './outline_renderer';
import {createPMREMGenerator} from './renderer';
import type {WorldViewerOptions} from './world_viewer_options';

// The room environment is the default environment map, used only when the
// `environmentMap` option isn't set.

/**
 * The scale of the room environment's lights. Unscaled, its average radiance
 * is about 1.05; scaled, it is about 0.13, darker than the
 * peppermint_powerplant_2 HDR map (about 0.23) that the renderer's exposure is
 * tuned for.
 */
const ROOM_ENVIRONMENT_INTENSITY = 0.12;

/** The blur of the room environment map, as recommended by three.js. */
const ROOM_ENVIRONMENT_BLUR = 0.04;

/**
 * RendererClient is the public, non-constructible interface for
 * RendererClientImpl. RendererClientImpl is responsible for requesting image
 * buffers from the main rendering service.
 */
export interface RendererClient {
  canvas: HTMLCanvasElement;

  render(scene: three.Scene, camera: three.Camera): void;
  setClearColor(color: three.Color, alpha: number): void;
  setSize(width: number, height: number): void;
  computePlaceholderEnvMap(): three.Texture;
  loadEnvironmentMap(): Promise<three.Texture>;
}

class RendererClientImpl implements RendererClient {
  context: ImageBitmapRenderingContext;
  private clearColor = colors.VIEWER_BACKGROUND_COLOR_LIGHT;
  private clearAlpha = 1;

  constructor(
    private readonly rendererService: WorldRendererService,
    readonly canvas: HTMLCanvasElement,
  ) {
    const context = canvas.getContext('bitmaprenderer');
    if (!context) {
      throw new Error('bitmaprenderer unsupported');
    }
    this.context = context;
  }

  render(scene: three.Scene, camera: three.Camera) {
    // We always have to set the properties before rendering as other clients
    // may have overwritten them.
    this.rendererService.setClearColor(this.clearColor, this.clearAlpha);
    this.rendererService.renderToClientContext(this.context, scene, camera);
  }

  setClearColor(color: three.Color, alpha: number) {
    this.clearColor = color;
    this.clearAlpha = alpha;
  }

  setSize(width: number, height: number) {
    this.canvas.width = width;
    this.canvas.height = height;
  }

  computePlaceholderEnvMap(): three.Texture {
    return this.rendererService.computePlaceholderEnvMap();
  }

  async loadEnvironmentMap(): Promise<three.Texture> {
    return this.rendererService.loadEnvironmentMap();
  }
}

/**
 * A renderer that renders threejs scenes to an offscreen canvas. The offscreen
 * canvas permits us to de-duplicate GPU resources as all vertex/texture data
 * will be instantiated once for a given three.js resource UUID as there's only
 * one actual webgl context being created.
 */
export class WorldRendererService {
  private readonly offscreenCanvas = new OffscreenCanvas(1, 1);
  private readonly renderer?: OutlineRenderer;
  private envMapPromise?: Promise<three.Texture>;
  private pmremGenerator?: three.PMREMGenerator;
  private placeholderEnvMap?: three.Texture;

  constructor(private readonly options: WorldViewerOptions) {
    if (this.options.suppressRendering) {
      this.renderer = undefined;
      return;
    }

    try {
      const failIfMajorPerformanceCaveat = this.options.noSoftwareRendering;
      this.renderer = new OutlineRenderer(
        this.offscreenCanvas,
        failIfMajorPerformanceCaveat,
      );
    } catch (e) {
      // Ignore
    }
  }

  createClient(canvas: HTMLCanvasElement): RendererClient {
    return new RendererClientImpl(this, canvas);
  }

  renderToClientContext(
    context: ImageBitmapRenderingContext,
    scene: three.Scene,
    camera: three.Camera,
  ) {
    const originalLayersMask = camera.layers.mask;
    this.maybeResize(context);

    this.renderer?.render(scene, camera);
    // Do not call transferFromImageBitmap if there is no renderer
    // or some frontend components crash and do not render correctly.
    if (this.renderer) {
      context.transferFromImageBitmap(
        this.offscreenCanvas.transferToImageBitmap(),
      );
    }
    camera.layers.mask = originalLayersMask;
  }

  private maybeResize(context: ImageBitmapRenderingContext) {
    // Client canvas sizes are in CSS pixels.
    const width = context.canvas.width;
    const height = context.canvas.height;
    this.renderer?.ensureCapacity(width, height);
    this.renderer?.setScissor(0, 0, width, height);
  }

  setClearColor(color: three.Color, alpha: number) {
    this.renderer?.setClearColor(color, alpha);
  }

  /**
   * Returns the PMREM generator shared by all environment maps. Its shaders
   * and intermediate render targets are compiled and allocated only once.
   */
  private getEnvGenerator(): three.PMREMGenerator | undefined {
    if (!this.renderer) {
      return undefined;
    }
    this.pmremGenerator ??= createPMREMGenerator(
      this.renderer.getThreeRenderer(),
    );
    return this.pmremGenerator;
  }

  computePlaceholderEnvMap(): three.Texture {
    // The placeholder env-map is simply a uniform 50% grey texture in
    // every direction. It is used as a placeholder while the default
    // environment map loads. It is computed once and shared by all clients.
    if (this.placeholderEnvMap) {
      return this.placeholderEnvMap;
    }
    const gen = this.getEnvGenerator();
    if (!gen) {
      return new three.Texture();
    }
    const b = new Float32Array(new Array(256 * 64).fill(0.5));
    const tex = new three.DataTexture(
      b,
      64,
      64,
      three.RGBAFormat,
      three.FloatType,
    );
    tex.needsUpdate = true;
    this.placeholderEnvMap = gen.fromEquirectangular(tex).texture;
    // The equirectangular source is no longer needed once prefiltered.
    tex.dispose();
    return this.placeholderEnvMap;
  }

  /**
   * Loads the environment map from the `environmentMap` option, prefiltered
   * for lighting. If the option is unset, uses a generic room environment. If
   * the map fails to load, uses the placeholder environment map.
   */
  async loadEnvironmentMap(): Promise<three.Texture> {
    if (this.envMapPromise) {
      return this.envMapPromise;
    }
    this.envMapPromise = new Promise((resolve) => {
      const gen = this.getEnvGenerator();
      if (!gen) {
        resolve(new three.Texture());
        return;
      }
      const source = this.options.environmentMap;
      if (source === undefined) {
        resolve(createRoomEnvironmentMap(gen));
        return;
      }
      if (source instanceof three.Texture) {
        resolve(gen.fromEquirectangular(source).texture);
        return;
      }
      const hdrLoader = new HDRLoader();
      hdrLoader.load(
        source,
        (tex) => {
          const hdrTarget = gen.fromEquirectangular(tex);
          // The equirectangular source is no longer needed once prefiltered.
          tex.dispose();
          resolve(hdrTarget.texture);
        },
        undefined,
        (error) => {
          console.error('Failed to load the environment map', error);
          resolve(this.computePlaceholderEnvMap());
        },
      );
    });
    return this.envMapPromise;
  }

  logRenderingStatsPeriodically() {
    this.renderer?.logRenderingStatsPeriodically();
  }
}

/**
 * Creates the scene of a generic, softly lit room (three.js's
 * `RoomEnvironment`), with its lights scaled for the renderer's exposure.
 */
function createRoomEnvironment(): RoomEnvironment {
  const room = new RoomEnvironment();
  scaleLightIntensity(room, ROOM_ENVIRONMENT_INTENSITY);
  return room;
}

/** Creates an environment map of the room environment, prefiltered for lighting. */
function createRoomEnvironmentMap(gen: three.PMREMGenerator): three.Texture {
  const room = createRoomEnvironment();
  const texture = gen.fromScene(room, ROOM_ENVIRONMENT_BLUR).texture;
  // The room's geometries and materials are no longer needed once prefiltered.
  room.dispose();
  return texture;
}

/**
 * Scales the intensity of all lights and emissive materials in `scene`.
 *
 * TODO(b/588515838) - The reason this scaling is required is that our renderer
 * is tuned to look good with our peppermint env-map. Without this scaling,
 * things look way too bright with the default env-map. Once we cleanup the
 * env-map API this should be removed.
 */
function scaleLightIntensity(scene: three.Object3D, scale: number) {
  const materials = new Set<three.Material>();
  scene.traverse((object) => {
    if (object instanceof three.Light) {
      object.intensity *= scale;
    } else if (object instanceof three.Mesh) {
      const material: three.Material | three.Material[] = object.material;
      for (const m of Array.isArray(material) ? material : [material]) {
        materials.add(m);
      }
    }
  });
  for (const material of materials) {
    if (
      material instanceof three.MeshLambertMaterial ||
      material instanceof three.MeshPhongMaterial ||
      material instanceof three.MeshStandardMaterial
    ) {
      material.emissiveIntensity *= scale;
    }
  }
}
