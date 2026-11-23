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

import {WorldEditor} from '../editor/world_editor';
import * as mathUtil from '../utils/math';
import {createPhantomForNode, syncPhantom} from './phantom';
import {computeDistanceToSurfaceBelow} from './snapping';
import {TransformTarget} from './transform_target';
import {WorldScene} from './world_scene';
import {IntrEntityNode} from './world_scene_nodes';

/** Options for `RigidObjectTarget`. */
export interface RigidObjectTargetOptions {
  /** The ID of the object that the transform controls attach to. */
  activeObjectId: number;
  /**
   * The IDs of all objects that are transformed together. Must include
   * `activeObjectId`.
   */
  objectIds: readonly number[];
  /** The editor used to commit the new transforms of the objects. */
  editor: WorldEditor;
}

/**
 * A `TransformTarget` that applies a rigid transform to a set of objects.
 *
 * The transform controls attach to a phantom copy of the active object. The
 * real objects are only moved when a drag ends: the same relative transform is
 * then applied to every object in the set and committed via the `WorldEditor`.
 */
export class RigidObjectTarget implements TransformTarget {
  /** The ID of the object that the transform controls attach to. */
  readonly activeObjectId: number;
  /** The node that is dragged: a phantom copy of the active node. */
  readonly node: IntrEntityNode;
  readonly enableSurfacePlacement = true;

  private readonly objectIds: readonly number[];
  private readonly editor: WorldEditor;

  /**
   * Creates a target for the given objects, or returns `undefined` if the
   * active object is not in the scene.
   */
  static create(
    scene: WorldScene,
    options: RigidObjectTargetOptions,
  ): RigidObjectTarget | undefined {
    const node = scene.getNodeById(options.activeObjectId);
    if (!node) return undefined;
    return new RigidObjectTarget(scene, node, options);
  }

  private constructor(
    private readonly scene: WorldScene,
    activeNode: IntrEntityNode,
    options: RigidObjectTargetOptions,
  ) {
    this.activeObjectId = options.activeObjectId;
    this.objectIds = [...options.objectIds];
    this.editor = options.editor;
    this.node = createPhantomForNode(activeNode);
  }

  get handle(): three.Object3D {
    return this.node;
  }

  onDragEnd(changed: boolean) {
    // Only commit transforms if the user actually moved the handle. We don't
    // want to commit if the user just clicks then deselects the object.
    if (!changed) return;
    const worldTActive = this.scene
      .getNodeById(this.activeObjectId)
      ?.matrixWorld.clone();
    if (!worldTActive) return;
    const activeToNode = mathUtil.getRelativeTransform(
      this.node.matrixWorld.clone(),
      worldTActive,
    );
    const transforms = new Map<number, three.Matrix4>();
    for (const oid of this.objectIds) {
      const worldMatrix = this.scene.getNodeById(oid)?.matrixWorld.clone();
      if (!worldMatrix) continue;
      transforms.set(oid, worldMatrix.multiply(activeToNode));
    }
    void this.editor.batchUpdateTransforms(transforms);
  }

  /**
   * Snaps the active object to the nearest surface below it, meaning the
   * closest surface in the -z world direction (see `snappingSurfaces()`), and
   * commits the new pose via the editor. Use for snapping objects to the
   * enclosure floor or other surfaces. No-op if there is no surface below.
   *
   * Only the active object is moved, the other objects in the set are not.
   */
  snapToSurfaceBelow() {
    const activeObject = this.scene.getNodeById(this.activeObjectId);
    if (!activeObject) return;
    const distanceToSurface = computeDistanceToSurfaceBelow(
      new three.Raycaster(),
      activeObject,
      [...this.snappingSurfaces()],
    );
    if (distanceToSurface === Infinity) {
      // No surface below, no-op.
      return;
    }
    const deltaVec = new three.Vector3(0, 0, -distanceToSurface);

    // Compute the new world position.
    const currentWorldPos = new three.Vector3();
    activeObject.getWorldPosition(currentWorldPos);
    const newWorldPos = currentWorldPos.clone().add(deltaVec);
    const newWorldMat = activeObject.matrixWorld
      .clone()
      .setPosition(newWorldPos);

    // Note - for now, we do not update the other selected objects. Perhaps
    // should snap them all to the surface below in the future.
    void this.editor.batchUpdateTransforms(
      new Map([[this.activeObjectId, newWorldMat]]),
    );
  }

  onWorldUpdated() {
    // Move the phantom back onto the objects it represents.
    syncPhantom(this.node, this.scene);
  }

  *snappingSurfaces(): Generator<three.Object3D> {
    yield this.scene.groundPlane;
    for (const entity of this.scene.root.traverseEntities()) {
      // Make sure we aren't testing against entities within the same object.
      // For the purposes of snapping, we only want to snap to the visual
      // meshes, not to the collision meshes. Reason being the user expects
      // the objects to be visually snapped together. Also the collision meshes
      // may not be tight.
      if (entity.parentObjectId !== this.node.parentObjectId) {
        for (const child of entity.visualNode.children) {
          if (child.visible) {
            yield child;
          }
        }
      }
    }
  }

  dispose() {
    this.node.removeFromParent();
  }
}
