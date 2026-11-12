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

/**
 * Returns whether `obj` is a `three.Mesh`.
 *
 * @param obj The object to check.
 * @returns True if `obj` is a mesh.
 */
export function isMesh(obj: three.Object3D): obj is three.Mesh {
  // This is more robust than more TS-idiomatic `instanceof` checks, especially if
  // multiple versions of three.js get loaded (say, by mistake).
  // - https://github.com/mrdoob/three.js/issues/6362
  // - https://discourse.threejs.org/t/the-mesh-of-gltf-model-cannot-be-recongnized-as-type-mesh-or-object3d/11952/3
  return (obj as three.Mesh).isMesh === true;
}

/**
 * Returns whether `obj` is a `three.Points`. See `isMesh` for why this does not
 * use `instanceof`.
 *
 * @param obj The object to check.
 * @returns True if `obj` is a points object.
 */
export function isPoints(obj: three.Object3D): obj is three.Points {
  return (obj as three.Points).isPoints === true;
}
