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
 * A very simple three.js renderer usful for creating screenshots with utilities
 * for generating addition rendering/GPU resources
 * TODO(b/277081479) Test this file
 */
import * as three from 'three';
import {GLTFExporter} from 'three/examples/jsm/exporters/GLTFExporter.js';
import {OutlinePass} from 'three/examples/jsm/postprocessing/OutlinePass.js';

import * as colors from './colors';

/** Default camera near plane distance. */
export const CAMERA_NEAR_PLANE = 0.1;
/** Default camera far plane distance. */
export const CAMERA_FAR_PLANE = 1000.0;
/** Exposure for use with the peppermint_powerplant_2 HDRI + ACES tonemapper. */
export const DEFAULT_EXPOSURE = 4;
/** Power to be used in the un-tonemapped env map. */
export const DEFAULT_ENV_POWER = 0.5;
/**
 * Rotation from the y-up convention of equirectangular environment maps to the
 * z-up world. Rotates the map's +Y (sky) to +Z.
 */
export const ENVIRONMENT_ROTATION: Readonly<three.Euler> = new three.Euler(
  Math.PI / 2,
  0,
  0,
);

/**
 * Exports the scene to gltf format.
 */
export async function sceneToGlb(
  scene: three.Object3D,
  filename: string,
): Promise<File> {
  const exporter = new GLTFExporter();
  return new Promise((resolve, reject) => {
    exporter.parse(
      scene,
      (result: object | ArrayBuffer) => {
        // passing 'binary: true' should always generate an ArrayBuffer.
        resolve(new File([result as ArrayBuffer], filename));
      },
      (error: ErrorEvent) => {
        reject(error);
      },
      {binary: true, embedImages: true},
    );
  });
}

/** Creates an outline pass with reusable defaults. */
export function createOutlinePass(
  size: three.Vector2,
  scene: three.Scene,
  camera: three.Camera,
): OutlinePass {
  const outlinePass = new OutlinePass(size, scene, camera);

  // Set the default setting for outline
  outlinePass.edgeStrength = 10;
  outlinePass.edgeGlow = 0.2;
  outlinePass.edgeThickness = 1;
  outlinePass.visibleEdgeColor.copy(colors.PRIMARY_LIGHT);
  outlinePass.hiddenEdgeColor.copy(colors.PRIMARY_LIGHT);

  return outlinePass;
}

const PMREM_LOD_MIN = 4;
const PMREM_EXTRA_LOD_SIGMA = [0.125, 0.215, 0.35, 0.446, 0.526, 0.582];

/**
 * Subclass of PMREMGenerator that uses the pre-0.186.0 spherical Gaussian blur
 * sigmas across CubeUV LOD levels instead of GGX VNDF convolution, preserving
 * the 0.177.0 IBL irradiance and specular mipmap distribution. This prevents
 * the bright equatorial horizon of peppermint_powerplant_2_256x128.hdr envmap
 * from bleeding into the dark +Z irradiance lobe on upward-facing surfaces
 */
class LegacyPMREMGenerator extends three.PMREMGenerator {
  declare private _renderer: three.WebGLRenderer;
  declare private _lodMax: number;
  declare private _sizeLods: number[];
  declare private _lodMeshes: three.Mesh[];
  declare private _blur: (
    cubeUVRenderTarget: three.WebGLRenderTarget,
    lodIn: number,
    lodOut: number,
    sigma: number,
  ) => void;

  protected _applyPMREM(cubeUVRenderTarget: three.WebGLRenderTarget) {
    const autoClear = this._renderer.autoClear;
    this._renderer.autoClear = false;
    const n = this._lodMeshes.length;
    const sigmas: number[] = [];
    for (let i = 0; i < n; i++) {
      let sigma = 1.0 / this._sizeLods[i];
      if (i > this._lodMax - PMREM_LOD_MIN) {
        sigma = PMREM_EXTRA_LOD_SIGMA[i - this._lodMax + PMREM_LOD_MIN - 1];
      } else if (i === 0) {
        sigma = 0;
      }
      sigmas.push(sigma);
    }
    for (let i = 1; i < n; i++) {
      const sigma = Math.sqrt(
        sigmas[i] * sigmas[i] - sigmas[i - 1] * sigmas[i - 1],
      );
      this._blur(cubeUVRenderTarget, i - 1, i, sigma);
    }
    this._renderer.autoClear = autoClear;
  }
}

/**
 * Creates a PMREMGenerator configured with the 0.177.0 Gaussian blur profile.
 */
export function createPMREMGenerator(
  renderer: three.WebGLRenderer,
): three.PMREMGenerator {
  return new LegacyPMREMGenerator(renderer);
}

/**
 * Creates a a flat env map.
 * @param power The intensity of the raw texel values.
 * @param renderer The threejs renderer necessary to create the GPU resource.
 */
export function createEnvironmentMap(
  power: number,
  renderer: three.WebGLRenderer,
): three.Texture {
  const b = new Float32Array(new Array(256 * 64).fill(power));
  const tex = new three.DataTexture(
    b,
    64,
    64,
    three.RGBAFormat,
    three.FloatType,
  );
  tex.needsUpdate = true;
  const gen = createPMREMGenerator(renderer);
  const target = gen.fromEquirectangular(tex);
  // Only the prefiltered target is kept; the generator's shaders and
  // intermediate targets, and the source texture are released.
  gen.dispose();
  tex.dispose();
  return target.texture;
}
