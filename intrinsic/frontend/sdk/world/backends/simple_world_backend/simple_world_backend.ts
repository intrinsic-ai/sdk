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

import {ReadonlyPose} from '@intrinsic-ai/core/math';
import {Observable, Subject} from 'rxjs';

import {WorldCollisionRule} from '../../model/collision_rules';
import {GeometryUpdate} from '../../model/geometry';
import {JointLimits, JointLimitsUpdate} from '../../model/joint_limits';
import {PhysicsUpdate} from '../../model/physics';
import {
  CollisionPair,
  CreatedFrame,
  DeleteEntityPolicy,
  DeleteObjectPolicy,
  UpdateType,
  WorldBackend,
  WorldUpdate,
} from '../../model/world_backend';
import {ObjectId, ObjectOrFrameId} from '../../model/world_ids';
import {WorldObject} from '../../model/world_object';

/**
 * A lightweight in-memory implementation of WorldBackend that allows
 * directly setting a WorldObject root, with no-op implementations for
 * backend mutation and network methods. Use for situations where you want
 * to drive the WorldBackend from a WorldObject that you create yourself.
 *
 * Usage:
 * Build the `WorldObject` root structure (e.g. using `buildWorldObjectTree`),
 * then pass it to the constructor or call `setWorldObject()` to push updates:
 *
 * ```ts
 * const backend = new SimpleWorldBackend();
 * const worldModel = new WorldModel('test_world', backend);
 *
 * // Build the WorldObject yourself. For example, using `buildWorldObjectTree` from
 * // simple_world_object_builder.ts.
 * const root = buildWorldObjectTree(
 *   [
 *     {
 *       id: ROOT_OBJECT_ID,
 *       name: 'root',
 *       parentEntityId: NULL_ENTITY_ID,
 *       rootEntityId: ROOT_ENTITY_ID,
 *       parentObjectId: NULL_OBJECT_ID,
 *       parentTThis: new three.Matrix4(),
 *       children: [],
 *     },
 *   ],
 *   [
 *     {
 *       id: ROOT_ENTITY_ID,
 *       name: 'Root',
 *       parentObjectId: ROOT_OBJECT_ID,
 *       parentId: NULL_ENTITY_ID,
 *       parentTThis: new three.Matrix4(),
 *       geometries: new Map(),
 *     },
 *   ],
 * );
 *
 * // Feed the WorldObject to the backend. This triggers a full state update.
 * backend.setWorldObject(root);
 * ```
 */
export class SimpleWorldBackend implements WorldBackend {
  private readonly worldUpdatesSubject = new Subject<WorldUpdate>();
  private currentWorldRoot?: WorldObject;
  private collisionRules: WorldCollisionRule[] = [];

  constructor(
    initialWorldRoot?: WorldObject,
    initialCollisionRules?: WorldCollisionRule[],
  ) {
    if (initialWorldRoot) {
      this.setWorldObject(initialWorldRoot);
    }
    if (initialCollisionRules) {
      this.collisionRules = initialCollisionRules;
    }
  }

  /**
   * Sets the current WorldObject root on the backend and emits a WorldUpdate to all listeners.
   */
  setWorldObject(
    worldRoot: WorldObject,
    updateType: UpdateType = UpdateType.FullUpdate,
  ): void {
    this.currentWorldRoot = worldRoot;
    this.worldUpdatesSubject.next(new WorldUpdate(worldRoot, updateType));
  }

  /**
   * Returns the current root WorldObject if one has been set.
   */
  getWorldObject(): WorldObject | undefined {
    return this.currentWorldRoot;
  }

  // WorldBackend interface methods:

  getWorldUpdatesStream(): Observable<WorldUpdate> {
    return this.worldUpdatesSubject.asObservable();
  }

  startStreaming(): void {
    // No-op
  }

  stopStreaming(): void {
    // No-op
  }

  async sync(force?: boolean): Promise<void> {
    if (this.currentWorldRoot) {
      this.worldUpdatesSubject.next(
        new WorldUpdate(this.currentWorldRoot, UpdateType.FullUpdate),
      );
    }
  }

  reloadFromLocalCachedData(): void {
    if (this.currentWorldRoot) {
      this.worldUpdatesSubject.next(
        new WorldUpdate(this.currentWorldRoot, UpdateType.FullUpdate),
      );
    }
  }

  async updateEntityName(entityId: number, name: string): Promise<void> {
    // No-op
  }

  async updateEntityParent(
    entityId: number,
    newParentId: number,
  ): Promise<void> {
    // No-op
  }

  async deleteEntity(
    entityId: number,
    policy: DeleteEntityPolicy,
  ): Promise<void> {
    // No-op
  }

  async updateTransform(
    id: number,
    refId: number,
    relativePose: ReadonlyPose,
    updateId?: number,
  ): Promise<void> {
    // No-op
  }

  async batchUpdateTransform(
    worldPoses: Map<number, ReadonlyPose>,
  ): Promise<void> {
    // No-op
  }

  async updateEntityGeometry(
    oid: number,
    eid: number,
    geometryUpdate: GeometryUpdate,
  ): Promise<void> {
    // No-op
  }

  async createFrame(
    parentId: ObjectOrFrameId,
    name: string,
    relativePose: ReadonlyPose,
  ): Promise<CreatedFrame> {
    return {id: ''};
  }

  async updateFrameProps(
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<void> {
    // No-op
  }

  async updateObjectName(oId: number, name: string): Promise<void> {
    // No-op
  }

  async updateObjectParentAndTransform(
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: ReadonlyPose,
  ): Promise<void> {
    // No-op
  }

  async deleteObject(oId: number, policy: DeleteObjectPolicy): Promise<void> {
    // No-op
  }

  async batchDelete(oids: number[], fids: number[]): Promise<void> {
    // No-op
  }

  getCollisionRules(): WorldCollisionRule[] {
    return this.collisionRules;
  }

  async updateWorldCollisionSettings(
    collisionsRules: WorldCollisionRule[],
  ): Promise<void> {
    this.collisionRules = collisionsRules;
  }

  async checkCollisions(): Promise<readonly CollisionPair[]> {
    // Collisions are not computed for simple worlds.
    return [];
  }

  async updateObjectPhysics(
    oId: number,
    physicsUpdate: PhysicsUpdate,
  ): Promise<void> {
    // No-op
  }

  async updateJoints(robotId: number, jvals: number[]): Promise<void> {
    // No-op
  }

  async updateJointApplicationLimits(
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined> {
    return undefined;
  }

  async getIkSolution(
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>> {
    return new Map();
  }
}
