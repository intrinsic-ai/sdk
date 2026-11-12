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
 * @file  This file contains math utility functions for world_viewer.
 */
import {Pose} from '@intrinsic-ai/core/math';
import * as three from 'three';

/**
 * Take two input matrices and a tolerance number, returns false if the two
 * matrices' elements have a total difference exceeding the tolerance, returns
 * true otherwise.
 *
 * @param mat1 The first input matrix.
 * @param mat2 The second input matrix.
 * @param tolerance The tolerance of two matrices' elements' sum difference.
 * @returns true if the two matrices are considered almost equal.
 */
export function matricesAlmostEqual(
  mat1: three.Matrix4,
  mat2: three.Matrix4,
  tolerance = 1e-6,
): boolean {
  // Subtract each element of the matrix, get the sum.
  let dif = 0;
  for (let i = 0; i < 16; i++) {
    dif += Math.abs(mat1.elements[i] - mat2.elements[i]);
  }
  return dif <= tolerance;
}

/**
 * @param mat The input matrix.
 * @returns true if the input matrix is considered identity.
 */
export function isIdentity(mat: three.Matrix4, tolerance = 1e-6): boolean {
  return matricesAlmostEqual(mat, new three.Matrix4(), tolerance);
}

/**
 * Update the input object's local transform.
 * @param obj The input object.
 * @param localTrans The new local transform to set for the input object.
 * @param updateMatrixWorld If true perform a recursive update to the obj and
 *     all its descendants' world matrix. This operation is potentially
 *     expensive.
 */
export function updateObjectLocalTransform(
  obj: three.Object3D,
  localTrans: three.Matrix4,
  updateMatrixWorld = false,
) {
  // Apply the input transform to the object node.
  localTrans.decompose(obj.position, obj.quaternion, obj.scale);
  obj.matrix.copy(localTrans);
  if (updateMatrixWorld) {
    obj.updateMatrixWorld(true);
  }
}

/**
 * Update the input object's local transform.
 * @param obj The input object.
 * @param worldTrans The new world transform to set for the input object.
 */
export function updateObjectWorldTransform(
  obj: three.Object3D,
  worldTrans: three.Matrix4,
) {
  if (obj.parent) {
    const delta = getRelativeTransform(worldTrans, obj.parent.matrixWorld);
    updateObjectLocalTransform(obj, delta, true);
  } else {
    updateObjectLocalTransform(obj, worldTrans, true);
  }
}

/**
 * @param mat1 Matrix representation of the frame to transform from.
 * @param mat2 Matrix representation of the frame to transform to.
 * @returns the relative transform between two frames.
 *
 * For any point P, the coord of P in frame i is Pi, Mi is the world transform
 * of each frame. So we have:
 *     P = M1 * P1 = M2 * P2;
 *     P2 = M2_Inverse * M1 * P1;
 * The transform M2_Inverse * M1 would take coord of a point from frame 1 to
 * frame 2.
 */
export function getRelativeTransform(
  mat1: three.Matrix4,
  mat2: three.Matrix4,
): three.Matrix4 {
  const mat = mat2.clone();
  mat.invert();
  mat.multiply(mat1);
  return mat;
}

/**
 * Creates a Pose from the translation and rotation of the given matrix. Any
 * scale component of the matrix is ignored.
 */
export function poseFromMatrix(mat: three.Matrix4): Pose {
  const position = new three.Vector3();
  const rotation = new three.Quaternion();
  mat.decompose(position, rotation, new three.Vector3());
  return new Pose(position, rotation);
}
