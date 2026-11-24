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

import {HttpErrorResponse, HttpStatusCode} from '@angular/common/http';
import {create} from '@bufbuild/protobuf';
import {extractErrorMessage} from '@intrinsic-ai/angular/errors';
import {ReadonlyPose} from '@intrinsic-ai/core/math';
import {createPoll} from '@intrinsic-ai/core/rxjs/polling';
import type {Pose as PoseProto} from '@intrinsic-ai/types/math/proto/pose_pb';
import {SceneObjectInstanceUpdateSchema} from '@intrinsic-ai/types/scene/proto/v1/scene_object_updates_pb';
import {
  type FrameReference,
  FrameReferenceSchema,
  type ObjectReference,
  ObjectReferenceSchema,
} from '@intrinsic-ai/types/world/proto/object_world_refs_pb';
import type {World} from '@intrinsic-ai/types/world/proto/object_world_service_pb';
import {
  Observable,
  Subject,
  Subscription,
  defer,
  exhaustMap,
  repeat,
  retry,
  throttleTime,
} from 'rxjs';

import {WorldCollisionRule} from '../../model/collision_rules';
import {CollisionRules} from '../../model/collision_rules';
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
import * as worldIdUtils from '../../model/world_ids';
import {ObjectId, ObjectOrFrameId} from '../../model/world_ids';
import {WorldLookup} from '../../model/world_lookup';
import {WorldObject} from '../../model/world_object';
import {
  ObjectWorldService,
  DeleteObjectPolicy as ServiceDeleteObjectPolicy,
} from '../../service_interfaces/object_world_service';
import {collisionSettingsToWorldCollisionRules} from '../proto/collision_rules_proto';
import {geometryUpdateToProto} from '../proto/geometry_proto';
import {
  jointLimitsFromProto,
  jointLimitsUpdateToProto,
} from '../proto/joint_limits_proto';
import {poseProtoToMatrix, poseToProto} from '../proto/pose_proto';
import {
  AssetInstanceDeleter,
  CollisionChecker,
  IkSolver,
  SceneObjectUpdater,
  WorldStateUpdate,
  WorldStreamingSource,
} from './object_world_service_backend_interfaces';
import {buildFromProtoList} from './proto_world_object_builder';

/**
 * We throttle world state updates to 50ms (20hz) to avoid slowing the
 * system down when there are many objects in the scene.
 */
export const WORLD_STATE_UPDATES_THROTTLE_TIME_MS = 50;

/**
 * When our data streaming subscription completes or errors, we attempt
 * to resubscribe indefinitely, with some delay.
 */
export const WORLD_STATE_UPDATES_RETRY_DELAY_MS = 1000;

/** Dependencies of an `ObjectWorldServiceBackend`. */
export interface ObjectWorldServiceBackendDeps {
  readonly objectWorldService: ObjectWorldService;
  /**
   * If provided, world state updates are streamed from it (e.g. over pubsub).
   * Otherwise, the world is polled for updates.
   */
  readonly worldStreamingSource?: WorldStreamingSource;
  /**
   * Used to compute inverse kinematics solutions. If not provided,
   * `getIkSolution` is unsupported and throws an error.
   */
  readonly ikSolver?: IkSolver;
  /**
   * Used to check the world for collisions. If not provided, `checkCollisions`
   * is unsupported and throws an error.
   */
  readonly collisionChecker?: CollisionChecker;
  /**
   * Used to apply geometry and simulation property updates to objects. If not
   * provided, `updateEntityGeometry` and `updateObjectPhysics` are unsupported
   * and throw an error.
   */
  readonly sceneObjectUpdater?: SceneObjectUpdater;
  /**
   * Used to delete objects that are asset instances. If not provided,
   * `deleteObject` and `batchDelete` (with objects) are unsupported and throw
   * an error.
   */
  readonly assetInstanceDeleter?: AssetInstanceDeleter;
}

/**
 * The WorldBackend that gets updates from the ObjectWorldService backend (running on the IPC).
 * This is the primary WorldBackend implementation.
 */
export class ObjectWorldServiceBackend implements WorldBackend {
  // Emits WorldUpdates whenever one is ready
  private readonly worldUpdatesSubject = new Subject<WorldUpdate>();
  private streamingSubscription = Subscription.EMPTY;

  // The proto and root WorldObject from the most recent WorldObject update
  private currentProto?: World;
  private worldRoot?: WorldObject;

  // The collision rules for this world
  private readonly collisionRules = new CollisionRules();

  // Map of structure hash to a Promise returned by `updateWorld`. This is used
  // by `updateWorldForHash` to ensure that the world is only updated once for a
  // given structure hash.
  private readonly hashWorldUpdates = new Map<string, Promise<void>>();

  // If there is a pending request to load this model from service.
  private loadRequestPending = false;

  private readonly objectWorldService: ObjectWorldService;
  private readonly worldStreamingSource?: WorldStreamingSource;
  private readonly ikSolver?: IkSolver;
  private readonly collisionChecker?: CollisionChecker;
  private readonly sceneObjectUpdater?: SceneObjectUpdater;
  private readonly assetInstanceDeleter?: AssetInstanceDeleter;

  TEST_ONLY = {
    // Expose this method for testing
    loadFromProto: (proto: World) => this.loadFromProto(proto),
  };

  constructor(
    private readonly id: string,
    deps: ObjectWorldServiceBackendDeps,
  ) {
    this.objectWorldService = deps.objectWorldService;
    this.worldStreamingSource = deps.worldStreamingSource;
    this.ikSolver = deps.ikSolver;
    this.collisionChecker = deps.collisionChecker;
    this.sceneObjectUpdater = deps.sceneObjectUpdater;
    this.assetInstanceDeleter = deps.assetInstanceDeleter;
  }

  /**
   * This stream emits an event any time a world update is ready.
   */
  getWorldUpdatesStream(): Observable<WorldUpdate> {
    return this.worldUpdatesSubject;
  }

  startStreaming(): void {
    this.stopStreaming();

    if (this.worldStreamingSource) {
      // We subscribe to world state events, but we throttle the updates to less than 60Hz. Deserializing
      // and updating the world is expensive, so we want to do it as infrequently as possible.
      // We also use `exhaustMap` which drops events while the current async callback is running. This is
      // important to prevent a growing backlog of state updates in scenarios where we have poor bandwidth
      // and each backend call may take a long time to complete. Without it, we end up with a backlog of
      // calls to get the `World` proto, preventing the app from updating.
      this.streamingSubscription = this.worldStreamingSource
        .subscribeToWorldState(this.id)
        .pipe(
          throttleTime(WORLD_STATE_UPDATES_THROTTLE_TIME_MS, undefined, {
            leading: true,
            trailing: true,
          }),
          exhaustMap(async (state) => {
            await this.updateInternalStateFromWorldStateResponse(state);
          }),
          // If the stream completes or errors, we attempt to retry indefinitely.
          retry({delay: WORLD_STATE_UPDATES_RETRY_DELAY_MS}),
          repeat({delay: WORLD_STATE_UPDATES_RETRY_DELAY_MS}),
        )
        .subscribe();

      // We do a one time sync to ensure that we have the latest state.
      void this.sync();
    } else {
      this.streamingSubscription = createPoll(() => {
        return defer(() => this.sync());
      }, 50).subscribe();
    }
  }

  stopStreaming(): void {
    this.streamingSubscription.unsubscribe();
  }

  /**
   * When `sync` is called, we fetch the full `World` from
   * `ObjectWorldService` via `updateWorld`.
   */
  async sync(force = false): Promise<void> {
    if (this.loadRequestPending && !force) {
      return;
    }

    this.loadRequestPending = true;
    try {
      await this.updateWorld();
    } finally {
      this.loadRequestPending = false;
    }
  }

  /**
   * Reload the WorldModel from the local cached world proto. As world is
   * reloaded, so will scene graph and tree graph be reloaded. Effectively this
   * would reset all essential client data based on curret copy of world proto
   * on client side. An example use case is when user updated something in the
   * scene or tree view and sent request to backend to update the model
   * but got returned some error, calling this would reset everything.
   */
  reloadFromLocalCachedData(): void {
    if (!this.currentProto) {
      console.warn(`World is not loaded yet, can't reload`);
      return;
    }
    this.loadFromProto(this.currentProto);
  }

  isLoaded(): boolean {
    return !!this.worldRoot;
  }

  /**
   * Given a `WorldStateUpdate`, updates our current WorldObject.
   * A `WorldStateUpdate` conceptually just contains pose updates. It does
   * not contain a full list of the objects in the world. Internally, this function
   * will trigger a call to retrieve the full World proto if it determines
   * that a state update is insufficient to update the world.
   *
   * The current update flow is as follows:
   *
   * Step 1 - Entry-points:
   * - When data-streaming / polling is enabled, we receive a stream of `WorldStateUpdate`
   * updates, which we call `updateInternalStateFromWorldStateResponse` with. We should only ever
   * have one call to `updateInternalStateFromWorldStateResponse` in progress at a time from
   * data-streaming (dropping others that may be received while we are processing, to avoid congestion).
   *
   * Step 2 - State updates vs full updates
   * - In `updateInternalStateFromWorldStateResponse`, we first attempt to call `loadStateFromProto` with
   * the given `stateProto`. This attempts to update our local state from the given state update. However,
   * this may not be possible if `stateProto` is empty or if the world structure has changed since our
   * last full update (as indicated by a mismatched structure hash).
   * - If `loadStateFromProto` succeeds, the update is done and we early out. Otherwise, we call `updateWorld`,
   * which fetches a full `WorldWithProto` object from the ObjectWorldService. We update our full
   * WorldObject from this proto and store the structure hash so that we can detect changes.
   *
   * TODO(b/543582279) - This update flow should be cleaned up to make the code easier to follow. I suspect there are some
   * subtle bugs related to updates being called simulataneously through `sync` and through data streaming.
   *
   * @param stateProto
   * @returns
   */
  private async updateInternalStateFromWorldStateResponse(
    stateProto?: WorldStateUpdate,
  ): Promise<void> {
    let stateUpdated = false;
    try {
      if (stateProto && this.isLoaded()) {
        stateUpdated = this.loadStateFromProto(stateProto);
      }
    } catch (e: unknown) {
      console.warn(
        `Failed to load world state proto ${extractErrorMessage(e)}`,
      );
    }

    // The state was successfully updated. There is no need to update the world
    // in its entirety.
    if (stateUpdated) {
      return;
    }

    const structureHash = stateProto?.worldMetadata?.worldStructureHash;
    if (structureHash) {
      await this.updateWorldForHash(structureHash);
    } else {
      await this.updateWorld();
    }
  }

  /**
   * Ensures that we only run one full update of the world for a single
   * structure hash at the same time. This is required because state updates can
   * come in very quickly and will continue to come in while the world update is
   * running. This means that we may receive many states with a new structure
   * hash before we ever finished updating the world for it.
   */
  private async updateWorldForHash(structureHash: string): Promise<void> {
    // Ensure that we only run one full update of
    // the world for it at the same time. This is required because state updates
    // can come in very quickly and will continue to come in while the world
    // update is running. This means that we can reach this code path with the
    // same hash many times over.
    const existingUpdatePromise = this.hashWorldUpdates.get(structureHash);
    if (existingUpdatePromise) {
      return existingUpdatePromise;
    }

    // It's important not to await the update here. It needs to be added to the
    // map BEFORE synchronous code execution continues.
    const updatePromise = this.updateWorld();
    this.hashWorldUpdates.set(structureHash, updatePromise);

    await updatePromise;
    // Remove the reference to the update promise once it's complete. This is
    // purely a memory optimization to prevent keeping references to update
    // promises around forever.
    this.hashWorldUpdates.delete(structureHash);
  }

  /**
   * Updates the world from the backend. This includes any structural changes.
   * Common cases where this needs to be called is when the model for a world ID
   * is first loaded or when its structure hash changes.
   *
   * Should only be called through `updateWorldForHash` if a world update is
   * needed for a structure hash change.
   */
  private async updateWorld(): Promise<void> {
    // Note: Downstream code relies on this method never throwing.
    try {
      const worldProto = await this.objectWorldService.getWorld(this.id);
      this.loadFromProto(worldProto);
      // NOTE: don't trigger collision checking here since it can be expensive.
      // It's okay to trigger it on demand when the user enables visualization
      // of collision meshes (and collisions).
    } catch (e: unknown) {
      console.warn(`Failed to load world proto ${extractErrorMessage(e)}.`);
    }
  }

  private loadFromProto(proto: World) {
    this.currentProto = proto;
    this.worldRoot = buildFromProtoList(proto.objects);
    if (!this.worldRoot) {
      throw new Error('Failed to build root object from proto');
    }

    if (proto.collisionSettings) {
      this.collisionRules.setRules(
        collisionSettingsToWorldCollisionRules(proto.collisionSettings),
      );
    }

    // Update tree and scene-graph.
    this.worldUpdatesSubject.next(
      new WorldUpdate(this.worldRoot, UpdateType.FullUpdate),
    );
  }

  // Returns true if the state update could be applied. Returns false
  // if we detect that a full world update must be performed (such as
  // when the world structure hash has changed).
  private loadStateFromProto(proto: WorldStateUpdate): boolean {
    if (!this.worldRoot) {
      console.warn(
        'Encountered falsy worldRoot in loadStateFromProto, aborting',
      );
      return false;
    }
    const worldLookup = new WorldLookup(this.worldRoot);
    const newHash = proto.worldMetadata?.worldStructureHash;
    const oldHash = this.currentProto?.worldMetadata?.worldStructureHash;
    if (!newHash || newHash !== oldHash) {
      return false;
    }
    // Set the updated entity poses
    for (const entityPose of proto.entityPoses) {
      if (!entityPose.parentTThis) {
        continue;
      }

      const entityId = entityPose.entityId;
      const mat = poseProtoToMatrix(entityPose.parentTThis);
      worldLookup.entityInitialTrans.set(entityId, mat);
      const entity = worldLookup.entities.get(entityId);
      if (entity) {
        entity.parentTThis.copy(mat);
      }
    }

    for (const newKinematicObjectJoints of proto.kinematicObjectJoints) {
      const oid = worldIdUtils.stringToNumberId(
        newKinematicObjectJoints.objectId,
      );
      const kinematicObject = worldLookup.objects.get(oid)?.asKinematic();
      // Skip objects that aren't KinematicWorldObjects
      if (kinematicObject === undefined) {
        continue;
      }
      const oldKinematicObjectJoints = [...kinematicObject.getJoints()];
      const newJointPositions = newKinematicObjectJoints.jointPositions;
      if (oldKinematicObjectJoints.length !== newJointPositions.length) {
        console.warn(
          `Joint list for robot '${kinematicObject.name}' has wrong number of DoFs. Expected ${oldKinematicObjectJoints.length}, got ${newJointPositions.length}`,
        );
        return false;
      }
      for (const [
        jointIndex,
        jointEntity,
      ] of oldKinematicObjectJoints.entries()) {
        jointEntity.value = newJointPositions[jointIndex];
      }
    }

    const rootObject = worldLookup.objects.get(worldIdUtils.ROOT_OBJECT_ID);
    if (!rootObject) {
      console.warn('Root object missing in loadStateFromProto');
      return false;
    }
    this.worldUpdatesSubject.next(
      new WorldUpdate(rootObject, UpdateType.StateUpdate),
    );
    return true;
  }

  async updateEntityName(entityId: number, name: string): Promise<void> {
    if (!this.worldRoot) {
      throw new Error('World not loaded');
    }
    const worldLookup = new WorldLookup(this.worldRoot);
    if (!worldLookup.isEntityFrame(entityId)) {
      throw new Error(
        'ObjectWorldServiceBackend only supports renaming frames.',
      );
    }
    await this.objectWorldService.updateFrameName(this.id, entityId, name);
  }

  /**
   * Updates the parent of an entity.
   *
   * NOTE: ObjectWorldServiceBackend only supports reparenting frames.
   */
  async updateEntityParent(
    entityId: number,
    newParentId: number,
  ): Promise<void> {
    if (!this.worldRoot) {
      throw new Error('World not loaded yet');
    }
    const worldLookup = new WorldLookup(this.worldRoot);

    if (!worldLookup.isEntityFrame(entityId)) {
      throw new Error(
        `Entity ${entityId} is not a frame. Reparenting non-frame entities is not supported in ObjectWorldServiceBackend.`,
      );
    }

    const parentObjectId = worldLookup.getEntityObjectId(newParentId);
    if (parentObjectId === worldIdUtils.NULL_OBJECT_ID) {
      throw new Error(`Parent entity ${newParentId} not found in world`);
    }
    const parentEntityName = worldLookup.entityNames.get(newParentId);
    if (!parentEntityName) {
      throw new Error(`Parent entity ${newParentId} name not found`);
    }
    await this.objectWorldService.updateFrameParent(
      this.id,
      entityId,
      parentObjectId,
      parentEntityName,
    );
  }

  async deleteEntity(
    entityId: number,
    policy: DeleteEntityPolicy,
  ): Promise<void> {
    if (!this.worldRoot) {
      throw new Error('World not loaded yet');
    }
    const worldLookup = new WorldLookup(this.worldRoot);

    // ObjectWorldServiceBackend only supports deleting frames.
    if (!worldLookup.isEntityFrame(entityId)) {
      throw new Error(
        `Entity ${entityId} is not a frame. Deleting non-frame entities is not supported in ObjectWorldServiceBackend.`,
      );
    }

    let force: boolean;
    switch (policy) {
      case DeleteEntityPolicy.DeleteIfNoChildren:
        force = false;
        break;
      case DeleteEntityPolicy.DeleteRecursively:
        force = true;
        break;
      case DeleteEntityPolicy.DeleteAndReparentChildren:
        throw new Error(
          'DeleteAndReparentChildren policy is not supported in ObjectWorldServiceBackend',
        );
    }

    await this.objectWorldService.deleteFrame(this.id, entityId, force);
  }

  // Transforms
  async updateTransform(
    id: number,
    refId: number,
    relativePose: ReadonlyPose,
    updateId?: number,
  ): Promise<void> {
    await this.objectWorldService.updateTransform(
      this.id,
      id,
      refId,
      poseToProto(relativePose),
      updateId,
    );
  }

  async batchUpdateTransform(
    worldPoses: Map<number, ReadonlyPose>,
  ): Promise<void> {
    const protoPoses = new Map<number, PoseProto>();
    for (const [id, pose] of worldPoses) {
      protoPoses.set(id, poseToProto(pose));
    }
    await this.objectWorldService.batchUpdateTransform(this.id, protoPoses);
  }

  async updateEntityGeometry(
    oid: number,
    eid: number,
    geometryUpdate: GeometryUpdate,
  ): Promise<void> {
    const sceneObjectUpdater = this.getSceneObjectUpdater(
      'updateEntityGeometry',
    );
    if (!this.worldRoot) {
      throw new Error('World not loaded');
    }
    const worldLookup = new WorldLookup(this.worldRoot);
    const object = worldLookup.objects.get(oid);
    if (!object) {
      throw new Error(`Object with ID ${oid} not found in world.`);
    }
    const entityName = worldLookup.entityNames.get(eid);
    if (!entityName) {
      throw new Error(`Entity with ID ${eid} not found in world.`);
    }

    const instanceUpdate = create(SceneObjectInstanceUpdateSchema, {
      update: {
        case: 'updateGeometry',
        value: geometryUpdateToProto(geometryUpdate, entityName),
      },
    });
    await sceneObjectUpdater.updateSceneObject(object.name, instanceUpdate);
  }

  // Frames
  async createFrame(
    parentId: ObjectOrFrameId,
    name: string,
    relativePose: ReadonlyPose,
  ): Promise<CreatedFrame> {
    if (!this.worldRoot) {
      throw new Error('World not loaded');
    }
    const worldLookup = new WorldLookup(this.worldRoot);
    const parentIdStr = worldIdUtils.numberToStringIdForObjectOrFrame(parentId);
    let parentReference: FrameReference | ObjectReference;
    if (worldLookup.frames.has(parentId)) {
      parentReference = create(FrameReferenceSchema, {
        frameReference: {case: 'id', value: parentIdStr},
      });
    } else if (worldLookup.objects.has(parentId)) {
      parentReference = create(ObjectReferenceSchema, {
        objectReference: {case: 'id', value: parentIdStr},
      });
    } else {
      throw new Error(`Parent with ID ${parentId} not found in world.`);
    }
    const frame = await this.objectWorldService.createFrame(
      this.id,
      parentReference,
      name,
      poseToProto(relativePose),
    );
    return {id: frame.id};
  }

  async updateFrameProps(
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<void> {
    await this.objectWorldService.updateFrameProps(
      this.id,
      fId,
      isAttachmentFrame,
    );
  }

  // Objects
  async updateObjectName(oId: number, name: string): Promise<void> {
    await this.objectWorldService.updateObjectName(this.id, oId, name);
  }

  async updateObjectParentAndTransform(
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: ReadonlyPose,
  ): Promise<void> {
    await this.objectWorldService.updateObjectParentAndTransform(
      this.id,
      id,
      pId,
      targetEntityName,
      worldPose && poseToProto(worldPose),
    );
  }

  async deleteObject(oId: number, policy: DeleteObjectPolicy): Promise<void> {
    const assetInstanceDeleter = this.getAssetInstanceDeleter('deleteObject');
    if (policy !== DeleteObjectPolicy.DeleteAndReparentChildren) {
      throw new Error(
        `Unsupported delete object policy: ${policy}. Only ` +
          `DeleteAndReparentChildren is supported by the ObjectWorldServiceBackend.`,
      );
    }

    // NOTE - all objects should be SceneObject instances and thus by resources.
    // These must be deleted using `assetInstanceDeleter` (which also handles deleting
    // the object from the ObjectWorldService). However, there are some legacy cases
    // where an object may not be a resource. To handle this, we fallback to calling
    // `objectWorldService.deleteObject` if the object is not a resource.
    const resourceName = this.getResourceName(oId);
    if (resourceName === undefined) {
      // This shouldn't happen: every object should be a resource. Fall back to
      // deleting the object directly through the ObjectWorldService if it is not
      // a resource.
      console.warn(
        `Object ${oId} is not a resource, deleting it through the ` +
          `ObjectWorldService instead.`,
      );
      await this.objectWorldService.deleteObject(
        this.id,
        oId,
        ServiceDeleteObjectPolicy.DeleteAndReparentChildren,
      );
      return;
    }

    await assetInstanceDeleter.deleteAssetInstance(resourceName, this.id);
  }

  async batchDelete(oids: number[], fids: number[]): Promise<void> {
    if (oids.length > 0) {
      // Check up front so that nothing is deleted if objects can't be.
      this.getAssetInstanceDeleter('batchDelete with objects');
    }
    for (const oid of oids) {
      await this.deleteObject(
        oid,
        DeleteObjectPolicy.DeleteAndReparentChildren,
      );
    }
    for (const fid of fids) {
      // Frames have no reparenting policy: a frame is always deleted along with
      // its child frames.
      try {
        await this.objectWorldService.deleteFrame(
          this.id,
          fid,
          /* force= */ true,
        );
      } catch (err: unknown) {
        // The frame may already be gone if it's parent object or frame was already deleted
        // as part of this batchDelete. We skip NotFound errors and keep going.
        if (
          err instanceof HttpErrorResponse &&
          err.status === HttpStatusCode.NotFound
        ) {
          continue;
        }
        throw err;
      }
    }
  }

  /**
   * Returns the `AssetInstanceDeleter`, or throws if none was provided.
   *
   * @param operation The name of the operation that requires the deleter, used
   *     in the error message.
   */
  private getAssetInstanceDeleter(operation: string): AssetInstanceDeleter {
    if (!this.assetInstanceDeleter) {
      throw new Error(
        `${operation} is not supported: no AssetInstanceDeleter was passed to ` +
          `ObjectWorldServiceBackend.`,
      );
    }
    return this.assetInstanceDeleter;
  }

  /**
   * Returns the `SceneObjectUpdater`, or throws if none was provided.
   *
   * @param operation The name of the operation that requires the updater, used
   *     in the error message.
   */
  private getSceneObjectUpdater(operation: string): SceneObjectUpdater {
    if (!this.sceneObjectUpdater) {
      throw new Error(
        `${operation} is not supported: no SceneObjectUpdater was passed to ` +
          `ObjectWorldServiceBackend.`,
      );
    }
    return this.sceneObjectUpdater;
  }

  /**
   * Returns the name of the resource instance that owns the given object, or
   * `undefined` if the object is not part of a resource instance.
   */
  private getResourceName(oId: number): string | undefined {
    if (!this.worldRoot) {
      return undefined;
    }
    return new WorldLookup(this.worldRoot).resourceNames.get(oId);
  }

  // Properties and Simulation
  getCollisionRules(): WorldCollisionRule[] {
    return [...this.collisionRules.getRules()];
  }

  async updateWorldCollisionSettings(
    collisionsRules: WorldCollisionRule[],
  ): Promise<void> {
    await this.objectWorldService.updateWorldCollisionSettings(
      this.id,
      collisionsRules,
    );
  }

  async checkCollisions(): Promise<readonly CollisionPair[]> {
    if (!this.collisionChecker) {
      throw new Error(
        'checkCollisions is not supported: no CollisionChecker was passed to ' +
          'ObjectWorldServiceBackend.',
      );
    }
    return this.collisionChecker.checkCollisions(this.id);
  }

  async updateObjectPhysics(
    oId: number,
    physicsUpdate: PhysicsUpdate,
  ): Promise<void> {
    const sceneObjectUpdater = this.getSceneObjectUpdater(
      'updateObjectPhysics',
    );
    if (!this.worldRoot) {
      throw new Error('World not loaded');
    }
    const worldLookup = new WorldLookup(this.worldRoot);
    const object = worldLookup.objects.get(oId);
    if (!object) {
      throw new Error(`Object with ID ${oId} not found in world.`);
    }

    const instanceUpdate = create(SceneObjectInstanceUpdateSchema, {
      update: {
        case: 'updateSimulationProperties',
        value: {
          isStatic: physicsUpdate.isStatic,
          isDisabled: physicsUpdate.isDisabled,
        },
      },
    });
    await sceneObjectUpdater.updateSceneObject(object.name, instanceUpdate);
  }

  async updateJoints(robotId: number, jvals: number[]): Promise<void> {
    await this.objectWorldService.updateJoints(this.id, robotId, jvals);
  }

  async updateJointApplicationLimits(
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined> {
    const limits = await this.objectWorldService.updateJointApplicationLimits(
      this.id,
      robotId,
      jointLimitsUpdateToProto(newLimits),
    );
    return limits ? jointLimitsFromProto(limits) : undefined;
  }

  async getIkSolution(
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>> {
    if (!this.ikSolver) {
      throw new Error(
        'getIkSolution is not supported: no IkSolver was passed to ' +
          'ObjectWorldServiceBackend.',
      );
    }
    return this.ikSolver.getIkSolution(this.id, tipId, targetPose);
  }
}
