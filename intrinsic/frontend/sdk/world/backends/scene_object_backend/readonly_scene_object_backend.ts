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
import type {SceneObject} from '@intrinsic-ai/types/scene/proto/v1/scene_object_pb';
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
import {sceneObjectToWorldObject} from './scene_object_to_world_object';

/**
 * A read-only WorldBackend for a SceneObject proto.
 *
 * Like SceneObjectBackend, it generates the world from a SceneObject, but it
 * doesn't connect to the SceneObjectEditService and can't edit the SceneObject.
 * Set the SceneObject to show with `setSceneObject`; every edit operation
 * throws.
 *
 * Streaming Behavior:
 * The ReadonlySceneObjectBackend emits a WorldUpdate any time the SceneObject
 * is set. There is no need to call `start/stopStreaming` here, those are
 * no-ops.
 */
export class ReadonlySceneObjectBackend implements WorldBackend {
  private readonly worldUpdatesSubject = new Subject<WorldUpdate>();
  private worldObject?: WorldObject;
  private sceneObject?: SceneObject;

  /**
   * Sets the SceneObject to show, replacing the current one, and emits a
   * WorldUpdate for it.
   */
  setSceneObject(sceneObject: SceneObject) {
    this.sceneObject = sceneObject;
    this.refreshWorldObject();
  }

  getWorldUpdatesStream(): Observable<WorldUpdate> {
    return this.worldUpdatesSubject;
  }

  startStreaming(): void {
    // This is a no-op. The ReadonlySceneObjectBackend always emits
    // WorldUpdates when the SceneObject is set.
  }

  stopStreaming(): void {
    // This is a no-op. The ReadonlySceneObjectBackend always emits
    // WorldUpdates when the SceneObject is set.
  }

  async sync(force?: boolean): Promise<void> {
    // The SceneObject is local, so there's nothing to fetch. We just refresh
    // the WorldObject from it to trigger a WorldUpdate.
    this.refreshWorldObject();
  }

  reloadFromLocalCachedData(): void {
    this.refreshWorldObject();
  }

  getLatestWorldObject(): WorldObject | undefined {
    return this.worldObject;
  }

  getSceneObject(): SceneObject {
    if (!this.sceneObject) {
      throw new Error(
        'No SceneObject is set. You must call setSceneObject before using the backend.',
      );
    }
    return this.sceneObject;
  }

  // Generates the WorldObject from the current SceneObject proto,
  // and emits an update.
  private refreshWorldObject() {
    if (!this.sceneObject) {
      return;
    }
    this.worldObject = sceneObjectToWorldObject(this.sceneObject);
    if (!this.worldObject) {
      throw new Error(`Failed to convert SceneObject to WorldObject`);
    }

    // All updates are FullUpdates, since the whole SceneObject is replaced.
    this.worldUpdatesSubject.next(
      new WorldUpdate(this.worldObject, UpdateType.FullUpdate),
    );
  }

  // Entity operations:

  async updateEntityName(entityId: number, name: string): Promise<void> {
    throw readonlyError('updateEntityName');
  }

  async updateEntityParent(
    entityId: number,
    newParentId: number,
  ): Promise<void> {
    throw readonlyError('updateEntityParent');
  }

  async deleteEntity(
    entityId: number,
    policy: DeleteEntityPolicy,
  ): Promise<void> {
    throw readonlyError('deleteEntity');
  }

  async updateTransform(
    id: number,
    refId: number,
    relativePose: ReadonlyPose,
    updateId?: number,
  ): Promise<void> {
    throw readonlyError('updateTransform');
  }

  async batchUpdateTransform(
    worldPoses: Map<number, ReadonlyPose>,
  ): Promise<void> {
    throw readonlyError('batchUpdateTransform');
  }

  async updateEntityGeometry(
    oid: number,
    eid: number,
    geometryUpdate: GeometryUpdate,
  ): Promise<void> {
    throw readonlyError('updateEntityGeometry');
  }

  // Frames operations:

  async createFrame(
    parentId: ObjectOrFrameId,
    name: string,
    relativePose: ReadonlyPose,
  ): Promise<CreatedFrame> {
    throw readonlyError('createFrame');
  }

  async updateFrameProps(
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<void> {
    throw readonlyError('updateFrameProps');
  }

  // Object operations:

  async updateObjectName(oId: number, name: string): Promise<void> {
    throw readonlyError('updateObjectName');
  }

  async updateObjectParentAndTransform(
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: ReadonlyPose,
  ): Promise<void> {
    throw readonlyError('updateObjectParentAndTransform');
  }

  async deleteObject(oId: number, policy: DeleteObjectPolicy): Promise<void> {
    throw readonlyError('deleteObject');
  }

  async batchDelete(oids: number[], fids: number[]): Promise<void> {
    throw readonlyError('batchDelete');
  }

  getCollisionRules(): WorldCollisionRule[] {
    throw new Error('Method not yet implemented.');
  }

  async updateWorldCollisionSettings(
    collisionsRules: WorldCollisionRule[],
  ): Promise<void> {
    throw readonlyError('updateWorldCollisionSettings');
  }

  checkCollisions(): Promise<readonly CollisionPair[]> {
    throw new Error('Method not yet implemented.');
  }

  async updateObjectPhysics(
    oId: number,
    physicsUpdate: PhysicsUpdate,
  ): Promise<void> {
    throw readonlyError('updateObjectPhysics');
  }

  async updateJoints(robotId: number, jvals: number[]): Promise<void> {
    throw readonlyError('updateJoints');
  }

  async updateJointApplicationLimits(
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined> {
    throw readonlyError('updateJointApplicationLimits');
  }

  async getIkSolution(
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>> {
    throw new Error('IK solving is not supported for Scene Object editing.');
  }
}

function readonlyError(operation: string): Error {
  return new Error(
    `${operation} is not supported: the ReadonlySceneObjectBackend is read-only.`,
  );
}
