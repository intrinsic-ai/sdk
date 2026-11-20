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

import type {
  JointLimits,
  JointLimitsUpdate,
} from '@intrinsic-ai/types/kinematics/types/joint_limits_pb';
import type {Pose} from '@intrinsic-ai/types/math/proto/pose_pb';
import type {
  FrameReference,
  ObjectReference,
} from '@intrinsic-ai/types/world/proto/object_world_refs_pb';
import type {
  Frame,
  World,
  Object$ as WorldObject,
} from '@intrinsic-ai/types/world/proto/object_world_service_pb';

import type {WorldCollisionRule} from '../model/collision_rules';
import type {ObjectId} from '../model/world_ids';

/** The different policies to use when deleting an object. */
export enum DeleteObjectPolicy {
  /** Delete the object only if it has no child objects (otherwise fail). */
  DeleteIfNoChildren = 'DeleteIfNoChildren',
  /**
   * Delete the object and reparent all of its child objects onto the deleted
   * object's parent, maintaining their world poses.
   */
  DeleteAndReparentChildren = 'DeleteAndReparentChildren',
  /** Delete the object and all of its children, recursively. */
  DeleteRecursively = 'DeleteRecursively',
}

/**
 * The subset of the object world service that world_viewer/core depends on.
 *
 * Object and frame ids are the numeric ids used by world_viewer/core (see
 * `world_ids`).
 */
export interface ObjectWorldService {
  // World
  /** Gets the full world, including all objects and collision settings. */
  getWorld(worldId: string): Promise<World>;

  // Transforms
  /**
   * Updates the transform of `id` (or of `updateId`, if given) such that the
   * pose of `id` relative to `refId` is `relativePose`.
   */
  updateTransform(
    worldId: string,
    id: number,
    refId: number,
    relativePose: Pose,
    updateId?: number,
  ): Promise<void>;

  /** Sets the world poses of several objects or frames, keyed by id. */
  batchUpdateTransform(
    worldId: string,
    worldPoses: Map<number, Pose>,
  ): Promise<void>;

  // Objects
  /** Renames an object. */
  updateObjectName(
    worldId: string,
    oId: number,
    name: string,
  ): Promise<WorldObject>;

  /**
   * Moves an object to a new parent object, optionally attaching it to the
   * parent entity named `targetEntityName` and setting its world pose.
   */
  updateObjectParentAndTransform(
    worldId: string,
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: Pose,
  ): Promise<void>;

  /** Deletes an object, handling its children according to `policy`. */
  deleteObject(
    worldId: string,
    oId: number,
    policy: DeleteObjectPolicy,
  ): Promise<void>;

  // Frames
  /** Creates a new frame under the given parent frame or object. */
  createFrame(
    worldId: string,
    parentReference: FrameReference | ObjectReference,
    name: string,
    relativePose: Pose,
  ): Promise<Frame>;

  /** Renames a frame. */
  updateFrameName(worldId: string, fId: number, name: string): Promise<Frame>;

  /** Updates the properties of a frame. */
  updateFrameProps(
    worldId: string,
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<Frame>;

  /**
   * Moves a frame to the parent object `pId`, optionally attaching it to the
   * parent entity named `entityTarget`.
   */
  updateFrameParent(
    worldId: string,
    fId: number,
    pId: number,
    entityTarget?: string,
  ): Promise<void>;

  /** Deletes a frame. If `force` is set, its child frames are deleted too. */
  deleteFrame(worldId: string, fId: number, force: boolean): Promise<void>;

  // Joints
  /**
   * Sets the joint positions of a robot. Resolves to the joint positions after
   * the update.
   */
  updateJoints(
    worldId: string,
    robotId: number,
    jvals: number[],
  ): Promise<readonly number[]>;

  /**
   * Updates the joint application limits of a robot. Resolves to the limits
   * after the update, if any.
   */
  updateJointApplicationLimits(
    worldId: string,
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined>;

  // Collisions

  /** Replaces the collision rules of the world. */
  updateWorldCollisionSettings(
    worldId: string,
    collisionRules: WorldCollisionRule[],
  ): Promise<void>;
}
