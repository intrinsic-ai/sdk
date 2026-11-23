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

import * as mathUtil from '../utils/math';
import {WorldScene} from './world_scene';
import {IntrEntityNode} from './world_scene_nodes';

/**
 * Creates a phantom copy of `node` and adds it next to `node` in the scene, at
 * the same transform.
 */
export function createPhantomForNode<T extends IntrEntityNode>(node: T): T {
  const phantom = node.createPhantom();
  node.parent?.add(phantom);
  mathUtil.updateObjectLocalTransform(phantom, node.matrix);
  return phantom;
}

/**
 * Moves every entity of `phantom` to the current world transform of the
 * corresponding original entity in `scene`.
 */
export function syncPhantom(phantom: IntrEntityNode, scene: WorldScene) {
  for (const entity of phantom.traverseEntities()) {
    const original = scene.getNodeById(entity.entityId);
    if (original) {
      mathUtil.updateObjectWorldTransform(entity, original.matrixWorld);
    }
  }
}
