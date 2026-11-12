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
 * @file This file provies utility functions for computing bounding
 * primitives for threejs objects.
 */
import * as three from 'three';

import {isMesh, isPoints} from './three_typechecks';

const BOUNDS_EPSILON = 1e-4;

/**
 * Compute the bounding box of the provided shapes.
 */
export function computeBoundingBox(...shapes: three.Object3D[]): three.Box3 {
  // NOTE(pushkarj): bbox is initialized with {min,max} as {Infinity/-Infinity}.
  const bbox = new three.Box3();
  for (const s of shapes) {
    bbox.expandByObject(s);
  }
  return bbox;
}
function extractVerticesFromObject(object: three.Object3D): three.Vector3[] {
  const points = new Array<three.Vector3>();
  object.updateMatrixWorld();
  object.traverse((child) => {
    if (isMesh(child) || isPoints(child)) {
      const posAttr = child.geometry.getAttribute('position');
      for (let i = 0; i < posAttr.count; i++) {
        // Handle the unlikely case of less than 3D vectors
        const p = new three.Vector3(
          posAttr.getX(i),
          posAttr.getY(i) ?? 0,
          posAttr.getZ(i) ?? 0,
        );
        points.push(p.applyMatrix4(child.matrixWorld));
      }
    }
  });
  return points;
}

/**
 * Compute an oriented bounding box for a mesh.
 */
export function computeOrientedBoundingBox(object: three.Object3D): {
  dimensions: three.Vector3;
  quaternion: three.Quaternion;
  position: three.Vector3;
} {
  const pos = object.position.clone();
  const quat = object.quaternion.clone();
  const scale = object.scale.clone();
  // Un-orient the object before computing the box.
  object.position.copy(new three.Vector3());
  object.quaternion.copy(new three.Quaternion());
  object.updateMatrixWorld();
  const box = new three.Box3().expandByObject(object);
  // Once we have the un-oriented points, then restore the transform.
  object.position.copy(pos);
  object.quaternion.copy(quat);
  object.scale.copy(scale);
  object.updateMatrixWorld();

  const dimensions = box.getSize(new three.Vector3()).addScalar(BOUNDS_EPSILON);
  const center = box.getCenter(new three.Vector3());
  const mat = new three.Matrix4().setPosition(center.divide(scale));
  mat.premultiply(object.matrixWorld);
  const quaternion = new three.Quaternion();
  const position = new three.Vector3();
  mat.decompose(position, quaternion, new three.Vector3());
  return {dimensions, quaternion, position};
}

/**
 * Compute a bounding sphere for a mesh.
 * TODO: Consider a more sophisticated strategy like
 * https://en.wikipedia.org/wiki/Smallest-circle_problem#Welzl's_algorithm or a
 * WASM CGAL module.
 */
export function computeBoundingSphere(object: three.Object3D): three.Sphere {
  return new three.Sphere().setFromPoints(extractVerticesFromObject(object));
}

/**
 * Compute an oriented bounding cylinder for a mesh. Projects the scaled, but
 * un-rotated vertices onto the plane orthogal to the major axis, computes the
 * bounding circle, then constructs the appropriate oriented bounding cylinder.
 */
export function computeOrientedBoundingCylinder(object: three.Object3D): {
  length: number;
  radius: number;
  quaternion: three.Quaternion;
  position: three.Vector3;
} {
  // Lazily find the major axis by taking the longest axis of the bounding box.
  const pos = object.position.clone();
  const quat = object.quaternion.clone();
  const scale = object.scale.clone();
  // Un-orient the object before computing the box.
  object.position.copy(new three.Vector3());
  object.quaternion.copy(new three.Quaternion());
  object.updateMatrixWorld();
  const box = new three.Box3();
  box.expandByObject(object);
  const dims = box.getSize(new three.Vector3());
  let longestAxis = 0;
  let longestAxisLength = 0;
  for (let i = 0; i < 3; i++) {
    if (dims.getComponent(i) > longestAxisLength) {
      longestAxisLength = dims.getComponent(i);
      longestAxis = i;
    }
  }
  // Project the object's vertices onto the plane orthogonal to the longest axis
  const projectionScale = scale.clone().setComponent(longestAxis, 0);
  object.scale.copy(projectionScale);
  object.updateMatrixWorld();
  // Get the center and radius of the disk made from the above projection.
  const sphere = computeBoundingSphere(object);
  // Once we have the disk, we restore the transform.
  object.position.copy(pos);
  object.quaternion.copy(quat);
  object.scale.copy(scale);
  object.updateMatrixWorld();
  // threejs cylinders use up == y, but we use up == z.
  const toAxis = new three.Quaternion().setFromUnitVectors(
    new three.Vector3(0, 0, 1),
    new three.Vector3().setComponent(longestAxis, 1),
  );
  const trans = sphere.center
    .add(
      new three.Vector3().setComponent(
        longestAxis,
        box.min.getComponent(longestAxis) + longestAxisLength / 2,
      ),
    )
    .divide(scale);
  const mat = new three.Matrix4().compose(
    trans,
    toAxis,
    new three.Vector3(1, 1, 1),
  );
  mat.premultiply(object.matrixWorld);
  const quaternion = new three.Quaternion();
  const position = new three.Vector3();
  mat.decompose(position, quaternion, new three.Vector3());
  return {
    radius: sphere.radius + BOUNDS_EPSILON,
    length: longestAxisLength + BOUNDS_EPSILON,
    quaternion,
    position,
  };
}
