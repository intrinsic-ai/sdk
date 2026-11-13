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
import {Observable} from 'rxjs';

import {WorldCollisionRule} from './collision_rules';
import {GeometryUpdate} from './geometry';
import {JointLimits, JointLimitsUpdate} from './joint_limits';
import {PhysicsUpdate} from './physics';
import {EntityId, ObjectId, ObjectOrFrameId} from './world_ids';
import {WorldObject} from './world_object';

/** A pair of entities that collide with each other. */
export interface CollisionPair {
  readonly leftEntityId: EntityId;
  readonly rightEntityId: EntityId;
}

/**
 * Used to indicate the type of update received from the backend. Listeners
 * can use this to only update what is necessary:
 *
 * FullUpdate:
 * An update to the full world structure (world objects may
 * have been added/removed, poses changed, joint values changed, etc)
 *
 * StateUpdate:
 * An update to the poses and joint values of the current world tree.
 * No updates were made to the structure of the world (like objects being added/removed).
 */
export enum UpdateType {
  FullUpdate = 'FullUpdate',
  StateUpdate = 'StateUpdate',
}

/**
 * An update received from the backend. Contains the latest World tree and the
 * update type.
 */
export class WorldUpdate {
  constructor(
    readonly worldRoot: WorldObject,
    readonly updateType: UpdateType,
  ) {}
}

/**
 * Result of creating a frame.
 */
export interface CreatedFrame {
  id: string;
}

/**
 * The different policies to use when deleting an entity
 */
export enum DeleteEntityPolicy {
  // Delete the entity only if it has no children (otherwise fail).
  DeleteIfNoChildren = 'DeleteIfNoChildren',
  // Delete the entity and reparent all of its children, maintaining their
  // world poses.
  DeleteAndReparentChildren = 'DeleteAndReparentChildren',
  // Delete the entity and all of its children, recursively.
  DeleteRecursively = 'DeleteRecursively',
}

/**
 * The different policies to use when deleting an object.
 */
export enum DeleteObjectPolicy {
  // Delete the object only if it has no child objects (otherwise fail).
  DeleteIfNoChildren = 'DeleteIfNoChildren',
  // Delete the object and reparent all of its child objects onto the deleted
  // object's parent, maintaining their world poses.
  DeleteAndReparentChildren = 'DeleteAndReparentChildren',
  // Delete the object and all of its children, recursively.
  DeleteRecursively = 'DeleteRecursively',
}

/**
 * WorldBackend is an interface for interacting with some arbitrary backend
 * that provides world state updates.
 *
 * The user of WorldBackend subscribes to `getWorldUpdatesStream()` to listen
 * for WorldUpdates, thereby staying in sync with the backend.
 *
 * The primary WorldBackend is the ObjectWorldServiceBackend, which fetches the
 * world info from the World HTTP service, then parses the proto to a
 * world.
 *
 * TODO(b/511291018) - Cleanup this interface.
 */
export interface WorldBackend {
  /**
   * This stream emits a `WorldUpdate` event any time a world update
   * is received (whether by polling or manual calls to `sync()`).
   *
   * Note that subscribing to this stream does not start/stop data streaming. Call
   * `startStreaming` if you would like to receive a stream of updates from the backend.
   */
  getWorldUpdatesStream(): Observable<WorldUpdate>;

  /**
   * Tells the backend to start streaming WorldUpdates from its remote source. Streaming will
   * continue until `stopStreaming` is called.
   *
   * Whenever new updates are received, a WorldUpdate will be sent to the
   * world updates stream.
   */
  startStreaming(): void;

  /**
   * Stops streaming.
   */
  stopStreaming(): void;

  /**
   * Triggers a sync with the source. Meaning - if there is not already
   * an update/poll in progress, force one now. A `FullUpdate` WorldUpdate will
   * get emit when the world is done reading.
   *
   * TODO(b/507945794) - The force flag is confusing here. Isn't sync always forced?
   */
  sync(force?: boolean): Promise<void>;

  /**
   * Triggers a reload of the world from the last response from the backend. Use
   * this when you need to trigger a new `FullUpdate` without actually pulling new
   * state from the backend.
   */
  reloadFromLocalCachedData(): void;

  // Entity operations:

  /**
   * Update the name of an entity
   *
   * NOTE: Certain backends may only allow renaming certain types of Entities (such as
   * only frames). The backend should throw an error if attempting an unsupported rename
   * operation.
   */
  updateEntityName(entityId: number, name: string): Promise<void>;

  /**
   * Update the parent of an entity.
   *
   * @param entityId ID of the entity to reparent
   * @param newParentId ID of the new parent entity
   */
  updateEntityParent(entityId: number, newParentId: number): Promise<void>;

  /**
   * Delete an entity, using the given policy.
   */
  deleteEntity(entityId: number, policy: DeleteEntityPolicy): Promise<void>;

  /**
   * Update an Entity's transform so that it achieves pose `relativePose` with-respect-to the
   * the Entity with `refId`.
   *
   * TODO(b/511291018) - Is this interface really needed? Seems more complicated than necessary.
   * Most of the time we just want to revise parentTThis.
   *
   * @param id The ID of the Entity to transform
   * @param refId The ID of the reference Entity
   * @param relativePose - The desired relative pose from Entity `refId` to Entity `id`, in the space of `refId`
   * @param updateId - If set, the ancestor of Entity `id` that we should move to achieve the desired relative pose.
   * If not set, we move Entity `id` within its parent (by changing its parentTThis pose) to achieve the target
   * relative pose.
   */
  updateTransform(
    id: number,
    refId: number,
    relativePose: ReadonlyPose,
    updateId?: number,
  ): Promise<void>;

  /**
   * Change many Entity poses at once
   * @param worldPoses A map of Entity ID to their new target world poses.
   */
  batchUpdateTransform(worldPoses: Map<number, ReadonlyPose>): Promise<void>;

  /**
   * Updates the geometry of the Entity with id `eid` within
   * WorldObject with id `oid`
   */
  updateEntityGeometry(
    oid: number,
    eid: number,
    geometryUpdate: GeometryUpdate,
  ): Promise<void>;

  // Frame operations:

  /**
   * Create a new frame
   * @param parentId ID of the parent object or frame
   * @param name The name of the new frame (must be unique among all Entities in the WorldObject)
   * @param relativePose The parentTThis to set on the frame
   */
  createFrame(
    parentId: ObjectOrFrameId,
    name: string,
    relativePose: ReadonlyPose,
  ): Promise<CreatedFrame>;

  /**
   * Update the properties of a frame
   */
  updateFrameProps(fId: number, isAttachmentFrame: boolean): Promise<void>;

  // Object operations:

  /**
   * Change a WorldObject's name.
   */
  updateObjectName(oId: number, name: string): Promise<void>;

  /**
   * Updates a WorldObject's parent and parentTThis transform.
   *
   * @param id ID of the WorldObject to reparent
   * @param pId ID of the new parent WorldObject
   * @param targetEntityName If given, reparents to the Entity with this
   * name within the parent WorldObject. Otherwise, reparents to the last Entity
   * of the WorldObject.
   * @param worldPose If given, sets the parentTThis transform of WorldObject `id`
   * so that its pose relative to the world root becomes `worldPose`.
   */
  updateObjectParentAndTransform(
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: ReadonlyPose,
  ): Promise<void>;

  /**
   * Deletes the given WorldObject, applying `policy` to its child objects.
   */
  deleteObject(oId: number, policy: DeleteObjectPolicy): Promise<void>;

  /**
   * Batch delete objects and frames.
   *
   * The child objects of each deleted object are reparented onto the deleted
   * object's parent, preserving their world poses. This does not apply to
   * frames: frames have no reparenting policy, so a frame is always deleted
   * along with its child frames. Deleting a frame never deletes an object.
   */
  batchDelete(oids: number[], fids: number[]): Promise<void>;

  // Get/set simulation properties:

  /**
   * Get the collision rules for this world
   */
  getCollisionRules(): WorldCollisionRule[];

  /**
   * Update the collision rules for this world
   */
  updateWorldCollisionSettings(
    collisionsRules: WorldCollisionRule[],
  ): Promise<void>;

  /**
   * Checks the world for collisions. Returns every pair of colliding entities.
   *
   * TODO(b/507945794) - It seems like there is probably a better place for this, but
   * for now putting it here for compatibility with the existing WorldModel interface.
   */
  checkCollisions(): Promise<readonly CollisionPair[]>;

  /**
   * Update the physics properties for the given object
   */
  updateObjectPhysics(oId: number, physicsUpdate: PhysicsUpdate): Promise<void>;

  /**
   * Updates the joint values of the robot with the given id.
   */
  updateJoints(robotId: number, jvals: number[]): Promise<void>;

  /**
   * Updates the joint limits of the robot with the given id.
   */
  updateJointApplicationLimits(
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined>;

  /**
   * Get IK solution for the input tip frame of a robot and supplied target
   * pose. If successful, the joints of the robot will be returned.
   */
  getIkSolution(
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>>;
}
