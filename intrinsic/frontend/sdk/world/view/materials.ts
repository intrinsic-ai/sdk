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

/**
 * @file Utility functions to create some threejs materials
 */
import * as three from 'three';

import * as colors from './colors';

/**
 * Creates a texture with a colored, outlined circle.
 * The texture is light gray with a darker outline.
 * @param size Size of the canvas used for the sprite texture.
 * @returns texture Created threejs Texture
 */
export function createReachabilityPointTexture(size: number) {
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  const context = canvas.getContext('2d');
  if (context !== null) {
    context.beginPath();
    // Fill the center of the circle with very light gray
    context.arc(
      canvas.width * 0.5,
      canvas.width * 0.5,
      canvas.width * 0.5,
      0,
      2 * Math.PI,
    );
    context.closePath();
    context.fillStyle = `rgb(
      230,
      230,
      230)`;

    context.fill();

    // Draw darker, thick outline.
    const outlineWidth = 3;
    context.lineWidth = outlineWidth;
    context.strokeStyle = `rgb(
      64,
      64,
      64)`;
    context.beginPath();
    context.arc(
      canvas.width * 0.5,
      canvas.width * 0.5,
      canvas.width * 0.5 - outlineWidth * 0.5,
      0,
      2 * Math.PI,
    );
    context.stroke();
  }

  return new three.CanvasTexture(canvas);
}

/**
 * Texture to which we can write text using ABCDiatype font
 * Text can be updated after initial creation.
 */
export class TextDisplayTexture extends three.CanvasTexture {
  constructor() {
    const canvas = document.createElement('canvas');
    canvas.width = 64;
    canvas.height = 64;
    super(canvas);
    return this;
  }

  update(newText: string) {
    const context = this.image.getContext('2d');
    if (!context) {
      return;
    }
    context.clearRect(0, 0, this.image.width, this.image.height);
    context.fillStyle = colors.PRIMARY_LIGHT.getStyle(); //new three.Color(0x000000).getStyle();
    context.beginPath();
    context.roundRect(0, 18, 64, 26, 4);
    context.fill();
    context.fillStyle = `rgb(
    255,
    255,
    255)`;
    context.font = '16px ABCDiatype';
    context.fillText(newText, 10, 36);
  }
}

/**
 * Font families providing Material Symbols, in order of preference. These match
 * the families registered in
 * google3/intrinsic/frontend/shared/components/styles/_icons.scss: the CDN
 * subset first, then the locally served full font used as offline fallback.
 */
const MATERIAL_SYMBOLS_FONT_FAMILIES = [
  'Material Symbols Outlined',
  'Material Symbols Outlined Local',
];

/**
 * Returns the first Material Symbols font family that has a loaded font face,
 * or null if none is loaded yet.
 *
 * Note that `document.fonts.check()` cannot be used for this since it returns
 * true for families without any registered font face, which would lead to the
 * icon name being rendered as plain text in a fallback font.
 */
function findLoadedMaterialSymbolsFontFamily(): string | null {
  const loadedFamilies = new Set<string>();
  document.fonts.forEach((face) => {
    if (face.status === 'loaded') {
      loadedFamilies.add(face.family.replace(/^["']|["']$/g, ''));
    }
  });
  return (
    MATERIAL_SYMBOLS_FONT_FAMILIES.find((family) =>
      loadedFamilies.has(family),
    ) ?? null
  );
}

/**
 * Creates a texture with a colored circle and a mat-icon image.
 * @param size Size of the canvas used for the sprite texture.
 * @param matIconName Name of the material icon. e.g. 'mat-icon-close'
 * @returns texture Created threejs Texture, or null if the icon font is not
 *     loaded yet.
 */
export function createMaterialIconTexture(
  size: number,
  matIconName: string,
  options?: {
    backgroundStyle?: string | CanvasGradient | CanvasPattern;
    iconColor?: string | CanvasGradient | CanvasPattern;
    fontSize?: number;
  },
): three.CanvasTexture | null {
  const fontFamily = findLoadedMaterialSymbolsFontFamily();
  if (fontFamily === null) {
    return null;
  }
  const fontSize = options?.fontSize ?? size;
  const fontStyle = `${fontSize}px '${fontFamily}'`;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  const context = canvas.getContext('2d');
  if (context !== null) {
    context.beginPath();
    // Fill the center of the circle with very light gray
    context.arc(
      canvas.width * 0.5,
      canvas.width * 0.5,
      canvas.width * 0.5,
      0,
      2 * Math.PI,
    );
    context.closePath();
    context.fillStyle =
      options?.backgroundStyle ??
      `rgb(
      230,
      230,
      230)`;

    context.fill();

    // Draw darker, thick checkmark.
    context.fillStyle =
      options?.iconColor ??
      `rgb(
      32,
      32,
      32)`;
    context.font = fontStyle;
    context.fillText(matIconName, (size - fontSize) / 2, (fontSize + size) / 2);
  }

  return new three.CanvasTexture(canvas);
}

/**
 * Creates a sprite material with a colored circle.
 * @param color Color of the material to create.
 * @returns material Created threejs SpriteMaterial
 */
export function createCircleSpriteMaterial(
  color: three.Color,
  label?: string,
  toneMapped = false,
  depthTest = true,
  sizeAttenuation = true,
) {
  const size = 64;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  const context = canvas.getContext('2d');
  if (context !== null) {
    context.beginPath();
    // Size and position of the control points are proportional to size
    context.arc(
      canvas.width * 0.5,
      canvas.width * 0.5,
      canvas.width * 0.5,
      0,
      2 * Math.PI,
    );
    context.closePath();
    context.fillStyle = color.getStyle();
    context.fill();
    if (label) {
      context.font = '48px ABCDiatype';
      context.fillStyle = new three.Color(0xffffff).getStyle();
      context.fillText(label, 18, 48);
    }
  }

  const texture = new three.CanvasTexture(canvas);
  const material = new three.SpriteMaterial({
    map: texture,
    toneMapped,
    depthTest,
    sizeAttenuation,
  });
  return material;
}

// Three.js r134's GLTFLoader uses a custom, unexported material type.  We
// detect it so we can coerce it to a standard physical material;
interface CustomSpecGlossMat {
  color: three.Color;
  glossiness: number;
  glossinessMap: three.Texture | null;
  specular: three.Color;
  specularMap: three.Texture | null;
}

// Bespoke properties defined in GLTFMeshStandardSGMaterial
function isCustomSpecGlossMat(mat: object): mat is CustomSpecGlossMat {
  if (
    'glossiness' in mat &&
    'glossinessMap' in mat &&
    'specular' in mat &&
    'specularMap' in mat
  ) {
    return true;
  }
  return false;
}

interface IntrMaterial {
  enableHatching: (enabled: boolean) => void;
  enableDiffHighlight: (enabled: boolean) => void;
  enableHighlight: (enabled: boolean) => void;
  enableDimming: (enabled: boolean) => void;
  makeInvisible: (enabled: boolean) => void;
}

function patchIntrMaterialShaders(shader: three.ShaderLibShader) {
  // Vertex shader to support material invisibility. This is used
  // to control per-object visibility independently from the parent.
  // If a material is marked as invisible, the corresponding object's
  // vertices are collapsed to the same location as a way to suppress
  // rendering.
  shader.vertexShader = shader.vertexShader.replace(
    '#include <begin_vertex>',
    [
      '#ifdef MAKE_INVISIBLE',
      'vec3 transformed = vec3(0);',
      '#else',
      '#include <begin_vertex>',
      '#endif',
    ].join('\n'),
  );

  shader.fragmentShader = shader.fragmentShader.replace(
    '#include <common>',
    ['#include <common>', 'uniform vec3 uHighlightColor;'].join('\n'),
  );

  shader.fragmentShader = shader.fragmentShader.replace(
    '#include <clipping_planes_fragment>',
    [
      '#include <clipping_planes_fragment>',
      '#ifdef USE_HATCH',
      'ivec2 c = ivec2(gl_FragCoord.xy);',
      'if ((c.x + c.y) % 2 == 0) { discard; }',
      '#endif',
    ].join('\n'),
  );

  shader.fragmentShader = shader.fragmentShader.replace(
    '#include <opaque_fragment>',
    [
      '#include <opaque_fragment>',
      // Overlay the highlight color over the desaturated, shaded color
      'float luma = dot(vec3(0.30, 0.59, 0.11), gl_FragColor.rgb);',
      'gl_FragColor.rgb = dot(uHighlightColor.rgb, vec3(1)) > 0.0 ? mix(vec3(luma), gl_FragColor.rgb, 0.5) * 0.9 : gl_FragColor.rgb;',
      'gl_FragColor += vec4(uHighlightColor, 0) * 0.1;',
    ].join('\n'),
  );

  shader.fragmentShader = shader.fragmentShader.replace(
    '#include <tonemapping_fragment>',
    [
      '#include <tonemapping_fragment>',
      '#ifdef USE_DIMMING',
      'gl_FragColor.rgb *= vec3(0.32);',
      '#endif',
    ].join('\n'),
  );
}

/**
 * Creates a custom mesh standard material that also supports a "highlighted" state.
 */
export class IntrVisualMaterial
  extends three.MeshPhysicalMaterial
  implements IntrMaterial
{
  private shader?: three.ShaderLibShader;
  private selectionHighlighted = false;
  private diffHighlighted = false;
  private colorHighlighted: three.Color | false = false;
  private dimmed = false;
  private hatch = false;
  private invisible = false;

  override customProgramCacheKey() {
    return (
      (this.hatch ? '1' : '0') +
      (this.dimmed ? '1' : '0') +
      (this.invisible ? '1' : '0')
    );
  }

  override onBeforeCompile(shader: three.ShaderLibShader) {
    this.shader = shader;
    this.updateShaderInputs();
    // Insert our custom uniform declaration.
    patchIntrMaterialShaders(this.shader);
  }

  // We override copy to handle common materials generated by the GLTF loader.
  override copy(source: three.Material): this {
    if (source instanceof three.MeshPhysicalMaterial) {
      super.copy(source);
    } else if (source instanceof three.MeshStandardMaterial) {
      // gross hack because our version of threejs's MeshStandardMaterial
      // doesn't play well with MeshPhysicalMaterial's copy method, despite
      // being the parent class.  This calls the grandparent class's copy
      // method.
      three.MeshStandardMaterial.prototype.copy.call(this, source);
      if (isCustomSpecGlossMat(source)) {
        this.metalness = 0;
        this.roughness = 1 - source.glossiness;
        // Don't attempt to handle glossinessMap as it's non-trivial
        this.specularColor = source.specular;
        this.specularColorMap = source.specularMap;
      }
    } else if (source instanceof three.MeshBasicMaterial) {
      this.color.copy(source.color);
      this.map = source.map;
    }

    return this;
  }

  private getHighlightColor(): three.Color {
    const color = new three.Color(0, 0, 0);
    if (this.colorHighlighted instanceof three.Color) {
      color.copy(this.colorHighlighted);
    } else if (this.selectionHighlighted) {
      color.copy(colors.PRIMARY_LIGHT);
    } else if (this.diffHighlighted) {
      color.copy(colors.OVERLAY_DIFF.clone().multiplyScalar(0.5));
    }
    return color;
  }

  updateShaderInputs() {
    if (this.shader) {
      this.shader.uniforms['uHighlightColor'] = {
        value: this.getHighlightColor(),
      };
    }
  }

  enableHatching(enabled: boolean) {
    this.hatch = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.hatch) {
      this.defines['USE_HATCH'] = '';
    } else {
      delete this.defines['USE_HATCH'];
    }

    // Recompile the shader (this should be guarded by our
    // customProgramCacheKey)
    this.needsUpdate = true;
  }

  enableDiffHighlight(enabled: boolean) {
    this.diffHighlighted = enabled;
    this.updateShaderInputs();
  }

  enableHighlight(enabled: boolean) {
    this.selectionHighlighted = enabled;
    this.updateShaderInputs();
  }

  enableColorHighlight(enabled: three.Color | false) {
    this.colorHighlighted = enabled;
    this.updateShaderInputs();
  }

  enableDimming(enabled: boolean) {
    this.dimmed = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.dimmed) {
      this.defines['USE_DIMMING'] = '';
    } else {
      delete this.defines['USE_DIMMING'];
    }
    this.needsUpdate = true;
  }

  isInvisible(): boolean {
    return this.invisible;
  }

  makeInvisible(enabled: boolean) {
    this.invisible = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.invisible) {
      this.defines['MAKE_INVISIBLE'] = '';
    } else {
      delete this.defines['MAKE_INVISIBLE'];
    }
    // Recompile the shader (this should be guarded by our
    // customProgramCacheKey)
    this.needsUpdate = true;
  }
}

/**
 * Creates a custom mesh basic material that also supports a "highlighted" state.
 */
export class IntrBasicMaterial
  extends three.MeshBasicMaterial
  implements IntrMaterial
{
  private shader?: three.ShaderLibShader;
  private selectionHighlighted = false;
  private diffHighlighted = false;
  private hatch = false;
  private dimmed = false;
  private invisible = false;
  override customProgramCacheKey() {
    return (
      (this.hatch ? '1' : '0') +
      (this.dimmed ? '1' : '0') +
      (this.invisible ? '1' : '0')
    );
  }

  override onBeforeCompile(shader: three.ShaderLibShader) {
    this.shader = shader;
    this.updateShaderInputs();
    // Insert our custom uniform declaration.
    patchIntrMaterialShaders(this.shader);
  }

  private getHighlightColor(): three.Color {
    const color = new three.Color(0, 0, 0);
    if (this.selectionHighlighted) {
      color.copy(colors.PRIMARY_LIGHT);
    } else if (this.diffHighlighted) {
      color.copy(colors.OVERLAY_DIFF.clone().multiplyScalar(0.2));
    }
    return color;
  }

  private updateShaderInputs() {
    if (this.shader) {
      this.shader.uniforms['uHighlightColor'] = {
        value: this.getHighlightColor(),
      };
    }
  }

  enableHatching(enabled: boolean) {
    this.hatch = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.hatch) {
      this.defines['USE_HATCH'] = '';
    } else {
      delete this.defines['USE_HATCH'];
    }

    // Recompile the shader (this should be guarded by our
    // customProgramCacheKey)
    this.needsUpdate = true;
  }

  makeInvisible(enabled: boolean) {
    this.invisible = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.invisible) {
      this.defines['MAKE_INVISIBLE'] = '';
    } else {
      delete this.defines['MAKE_INVISIBLE'];
    }

    // Recompile the shader (this should be guarded by our
    // customProgramCacheKey)
    this.needsUpdate = true;
  }

  isInvisible(): boolean {
    return this.invisible;
  }

  enableDiffHighlight(enabled: boolean) {
    this.diffHighlighted = enabled;
    this.updateShaderInputs();
  }

  enableHighlight(enabled: boolean) {
    this.selectionHighlighted = enabled;
    this.updateShaderInputs();
  }

  enableDimming(enabled: boolean) {
    this.dimmed = enabled;
    if (!this.defines) {
      this.defines = {};
    }
    if (this.dimmed) {
      this.defines['USE_DIMMING'] = '';
    } else {
      delete this.defines['USE_DIMMING'];
    }
    this.needsUpdate = true;
  }
}

/**
 * Creates a custom points material that also supports a "highlighted" state.
 */
export class IntrVisualPointsMaterial extends three.PointsMaterial {
  private shader?: three.ShaderLibShader;
  private highlighColor = new three.Color(0, 0, 0);
  override size = 1;
  override vertexColors = true;
  override sizeAttenuation = false;

  override onBeforeCompile(shader: three.ShaderLibShader) {
    this.shader = shader;
    this.updateShaderInputs();
    // Insert our custom uniform declaration.
    this.shader.fragmentShader = this.shader.fragmentShader.replace(
      '#include <common>',
      ['#include <common>', 'uniform vec3 uHighlightColor;'].join('\n'),
    );

    this.shader.fragmentShader = this.shader.fragmentShader.replace(
      '#include <opaque_fragment>',
      [
        '#include <opaque_fragment>',
        'gl_FragColor += vec4(uHighlightColor, 0) * 0.3;',
      ].join('\n'),
    );
  }

  updateShaderInputs() {
    if (this.shader) {
      this.shader.uniforms['uHighlightColor'] = {
        value: this.highlighColor,
      };
    }
  }

  enableHighlight(enabled: boolean) {
    if (enabled) {
      this.highlighColor = colors.PRIMARY_LIGHT;
    } else {
      this.highlighColor = new three.Color(0, 0, 0);
    }
    this.updateShaderInputs();
  }
}

/**
 * Material that uses a shader to deform a geometry to spherical coordinates
 * in order to create a spherical section.
 * Call setTheta to set the extend of the deformation, from 0 to 2PI degrees.
 * @param radius Radius of the deformation sphere.
 * @param firstConeVertexId Any vertex IDs >= this value will not be transformed
 * and forced to the local origin.
 */
export class SphericalDeformerMaterial extends three.MeshBasicMaterial {
  shader?: three.ShaderLibShader;
  theta = 0.0;
  linewidth = 3;
  setTheta(theta: number) {
    this.theta = theta;
    if (this.shader) {
      this.shader.uniforms['theta'].value = this.theta;
    }
  }
  constructor(
    args: three.MeshBasicMaterialParameters & {linewidth?: number},
    radius: number,
    firstConeVertexId: number,
  ) {
    const {linewidth, ...meshBasicArgs} = args;
    super(meshBasicArgs);
    if (linewidth !== undefined) {
      this.linewidth = linewidth;
    }
    this.onBeforeCompile = (shader: three.ShaderLibShader) => {
      // Patch the shader to avoid slowly recomputing the entire vertex buffer (host
      // side) every time we want to update the cone angle; See:
      // https://threejs.org/docs/#api/en/materials/Material.onBeforeCompile

      this.shader = shader;
      this.shader.uniforms['radius'] = {value: radius};
      this.shader.uniforms['firstConeVertexId'] = {
        value: firstConeVertexId,
      };
      this.shader.uniforms['theta'] = {
        value: this.theta,
      };

      // Append our margin uniform declaration
      this.shader.vertexShader = this.shader.vertexShader.replace(
        '#include <common>',
        [
          '#include <common>',
          'uniform float theta;',
          'uniform float radius;',
          'uniform int firstConeVertexId;',
          'vec3 toPolar(in vec3 position)',
          '{',
          '   vec3 pos = vec3(position.x + 0.5, position.y + 0.5, position.z);',
          '   float lon=PI * 2.0 * pos.x;',
          '   float lat=(pos.y * theta) - (PI * 0.5);',
          '   float x = radius * cos(lat) * cos(lon);',
          '   float y = radius * cos(lat) * sin(lon);',
          '   float z = radius * sin(lat);',
          '   return vec3(x, y, z);',
          '}',
        ].join('\n'),
      );

      this.shader.vertexShader = this.shader.vertexShader.replace(
        '#include <begin_vertex>',
        [
          'float sphericalPosScale = clamp(float(firstConeVertexId - gl_VertexID), 0.0, 1.0);',
          'vec3 transformed = toPolar(position) * sphericalPosScale;',
        ].join('\n'),
      );
    };
    this.needsUpdate = true;
  }
}
