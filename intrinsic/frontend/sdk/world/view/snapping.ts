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

import {DebugObjects} from './debug';
import {IntrEntityNode} from './world_scene_nodes';

/**
 * Computes a bounding box of the object's visual meshes. For the purpose
 * of snapping, we ignore the collision meshes. The visual meshes are typically
 * tighter and for snapping the user likely cares more about
 * the object visually being in the right place.
 *
 * @param obj The object for which we compute a bounding box.
 *
 * @returns A bounding box that contains all visible geometries of the object.
 */
export function getSnappingBoundingBox(obj: IntrEntityNode): three.Box3 {
  const snapBox = new three.Box3();
  for (const entity of obj.traverseEntities()) {
    for (const child of entity.visualNode.children) {
      if (child.visible) {
        snapBox.expandByObject(child);
      }
    }
  }
  return snapBox;
}

/**
 * Computes the distance to the nearest collision surface below a given object.
 *
 * @param raycaster A raycaster used to detect collisions.
 * @param fromObj An object for which we compute the distance to `toMeshes`.
 * @param toMeshes A list of meshes for which we compute the distance from
 *     `fromMeshes`.
 *
 * @returns The distance to the nearest collision, or infinity if nothing is hit.
 */
export function computeDistanceToSurfaceBelow(
  raycaster: three.Raycaster,
  fromObj: IntrEntityNode,
  toMeshes: three.Object3D[],
): number {
  const snapBox = getSnappingBoundingBox(fromObj);

  // We raycast from the bottom-center of the bounding box to "toMeshes".
  // This isn't perfect but is simple and works well for many basic objects.
  // Use an epsilon so that we detect that we've hit something even if the
  // bottom of the object is coincident with the surface.
  const originEpsilon = 0.001;
  const raycastStart = snapBox.getCenter(new three.Vector3());
  raycastStart.z = snapBox.min.z + originEpsilon;
  raycaster.set(raycastStart, new three.Vector3(0, 0, -1));
  const intersections = raycaster.intersectObjects(toMeshes, true);
  if (intersections.length === 0) {
    return Infinity;
  }
  return intersections[0].distance - originEpsilon;
}

/**
 * This function is the same as computeDistanceToSurfaceBelow, but also
 * adds debug objects to visualize the bounding boxes and meshes.
 * Use only for debugging.
 */
export function computeDistanceToSurfaceBelowWithDebug(
  raycaster: three.Raycaster,
  fromObj: IntrEntityNode,
  toMeshes: three.Object3D[],
  debugObjects: DebugObjects,
): number {
  // Draw debug objects for the snap-box and the meshes we want to snap to.
  // Then perform the distance computation as usual.
  const snapBox = getSnappingBoundingBox(fromObj);
  for (const obj of toMeshes) {
    debugObjects.highlightMeshes(obj);
  }
  debugObjects.addBox(snapBox);

  return computeDistanceToSurfaceBelow(raycaster, fromObj, toMeshes);
}
