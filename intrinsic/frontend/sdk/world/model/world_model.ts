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

import {Signal, computed, signal} from '@angular/core';
import {Pose, ReadonlyPose} from '@intrinsic-ai/core/math';
import {Observable, ReplaySubject, Subject, defer} from 'rxjs';
import {finalize, map, share} from 'rxjs/operators';
import * as three from 'three';

import * as mathUtils from '../utils/math';
import {WorldCollisionRule} from './collision_rules';
import {CameraSensor, Entity, FrameEntity} from './entity';
import {GeometryUpdate} from './geometry';
import {Geometry} from './geometry';
import {JointLimits, JointLimitsUpdate} from './joint_limits';
import {PhysicsUpdate} from './physics';
import {
  CollisionPair,
  CreatedFrame,
  DeleteEntityPolicy,
  DeleteObjectPolicy,
  UpdateType,
  WorldBackend,
  WorldUpdate,
} from './world_backend';
import {REFERENCE_STRING_SEPARATOR} from './world_ids';
import * as worldIdUtils from './world_ids';
import {EntityId, FrameId, ObjectId, ObjectOrFrameId} from './world_ids';
import {WorldLookup} from './world_lookup';
import {WorldObject} from './world_object';

/**
 * The WorldModel class is the primary interface for loading, querying,
 * and modifying a live-updating 3d world. The world itself
 * is represented by a WorldObject tree, which hierarchically represents
 * the links, joints, and frames of a scene, along with all of their transforms,
 * geometry, etc.
 *
 * The WorldModel requires a WorldBackend, which abstracts away the raw
 * interface for interacting with a world data-source. The WorldBackend is responsible
 * for generating a WorldObject tree from its data and emitting WorldObject updates as
 * changes are made to the world. You can access the current WorldObject tree using
 * `getRootObject` and/or `getRootObjectSignal`. The WorldObject tree is meant to be immutable,
 * with the WorldBackend fully regenerating it as the world changes.
 *
 * The WorldModel user should call functions on the `worldModel` to interact with the world,
 * instead of calling functions on the WorldBackend directly. Internally, the WorldModel
 * delegates functionality (like `createFrame`, `deleteFrame`, etc) to the WorldBackend as needed.
 *
 * IMPORTANT:
 * The WorldModel is meant to be fully agnostic to the WorldBackend used. Do not
 * add dependencies here on any services that are specific to a particular world
 * backend (like the ObjectWorldService backend, OnPremAssetsService, etc). Third-party users of
 * the world-viewer SDK may implement their own backends that are entirely independent from Intrinsic
 * Core backend services.
 */
export class WorldModel {
  // Static properties are shared among all WorldModel instances.
  // The Object ID of world origin.
  // NOTE: this is an implicit knowledge of intrinsic's world model.
  static ROOT_OBJECT_ID = worldIdUtils.ROOT_OBJECT_ID;
  // The Entity ID of a null object.
  static NULL_OBJECT_ID = worldIdUtils.NULL_OBJECT_ID;

  // The source of world updates for this WorldModel
  private readonly worldBackend: WorldBackend;

  // The most recent world-tree received from the backend. This
  // will be updated and replaced as new world updates are received.
  // NOTE: We use a custom equality check that always returns false to ensure
  // that consumers are notified when the world tree is mutated in place
  // (e.g., during poses-only StateUpdates), even if the root object reference
  // remains the same.
  private readonly worldRoot = signal<WorldObject | undefined>(undefined, {
    equal: () => false,
  });

  // A lookup map into the received world tree
  private worldLookup?: WorldLookup;

  // A multicast stream that fires whenever a world update is received.
  // Gives `true` for state updates (pose only) and false for full updates.
  private pollStream: Observable<boolean>;

  // Note: Since this class does not have a "destroy" hook, the subjects here
  // will never complete, and the subscribers won't be cleaned up.

  // Observable that notifies subscribers when WorldModel is loaded first time.
  private readonly loadedSubject = new ReplaySubject<boolean>(1);
  readonly loaded: Observable<boolean> = this.loadedSubject;
  // Observable that notifies subscribers when WorldModel is updated.
  private readonly updatedSubject = new Subject<void>();
  readonly updated: Observable<void> = this.updatedSubject;

  // Observable that notifies subscribers when WorldModel state is updated.
  private readonly updatedStateSubject = new Subject<void>();
  readonly updatedState: Observable<void> = this.updatedStateSubject;

  /**
   * Construct a WorldModel instance, but note it is empty and has to be loaded
   * from a world proto.
   */
  constructor(
    readonly id: string,
    worldBackend: WorldBackend,
    noAutoLoad?: boolean,
  ) {
    this.worldBackend = worldBackend;

    // Subscribe to all world updates from the backend. These are received
    // from both manual `sync` calls and from streaming (if enabled).
    this.worldBackend
      .getWorldUpdatesStream()
      .pipe(
        map((update: WorldUpdate) => {
          this.handleWorldUpdate(update);
          const isStateUpdate = update.updateType === UpdateType.StateUpdate;
          if (isStateUpdate) {
            this.updatedStateSubject.next();
          } else {
            this.updatedSubject.next();
          }
          return isStateUpdate;
        }),
      )
      .subscribe();

    // Setup the poll stream. This is a multicast stream that, upon subscribe,
    // turns on streaming on the backend. When the last subscription ends, we
    // turn off streaming on the backend.
    this.pollStream = defer(() => {
      this.worldBackend.startStreaming();
      return this.worldBackend.getWorldUpdatesStream().pipe(
        map(
          (update: WorldUpdate) => update.updateType === UpdateType.StateUpdate,
        ),
        finalize(() => {
          this.worldBackend.stopStreaming();
        }),
      );
    }).pipe(share());

    if (noAutoLoad) {
      this.loadedSubject.next(false);
    } else {
      // Automatically attempt to load the world from backend, ignore errors.
      void this.sync();
    }
  }

  // Returns a stream that emits whenever the world is updated.
  // The event is `true` for state updates (aka pose-only updates) and `false` for
  // full world tree updates (which may include objects added/removed)
  // from the world tree.
  poll(): Observable<boolean> {
    return this.pollStream;
  }

  private handleWorldUpdate(update: WorldUpdate): void {
    const firstLoad = !this.isLoaded();
    // Update our internal WorldRoot and the lookup helper associated with it
    this.worldRoot.set(update.worldRoot);
    this.worldLookup = new WorldLookup(update.worldRoot);

    if (firstLoad) {
      this.loadedSubject.next(true);
    }
  }

  isLoaded(): boolean {
    return !!this.worldRoot();
  }

  // Load the world from backend.
  async sync(force = false): Promise<void> {
    await this.worldBackend.sync(force);
  }

  /**
   * Reload the model from backend. Use this to resolve potential discrepancies
   * between the local cached proto and the backend. For example: when some
   * other client updated the world to backend, use this to refresh.
   */
  refresh(force = false): Observable<void> {
    return defer(() => this.sync(force));
  }

  /**
   * Checks the world for collisions. Returns every pair of colliding entities.
   */
  checkCollisions(): Promise<readonly CollisionPair[]> {
    return this.worldBackend.checkCollisions();
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
    this.worldBackend.reloadFromLocalCachedData();
  }

  async getIkSolution(
    tipId: number,
    targetPose: ReadonlyPose,
  ): Promise<Map<number, number>> {
    return await this.worldBackend.getIkSolution(tipId, targetPose);
  }

  // Returns true if the world contains the given object/frame resource.
  hasObjectOrFrame(rid: ObjectOrFrameId): boolean {
    return this.hasObject(rid) || this.hasFrame(rid);
  }

  /**
   * An attachment point (for EOAT attachment workflow) is a frame tagged with
   * the 'is_attachment_frame' attribute. While we do not allow parenting
   * objects to a frame.
   */
  hasAttachmentPoint(rid: FrameId): boolean {
    return !!this.worldLookup?.frames.get(rid)?.isAttachmentFrame;
  }

  /**
   * Return a list of IDs for all attachment points.
   */
  getAttachmentPointIds(): FrameId[] {
    return this.getFrameIds().filter((fid) => this.hasAttachmentPoint(fid));
  }

  /**
   * Returns the root object of the world. This object is not a physical object,
   * but is the ancestor of all other objects in the world. This object has no
   * parent.
   *
   * This function throws if the world is not loaded yet.
   */
  getRootObject(): WorldObject {
    const worldRoot = this.worldRoot();
    if (!worldRoot) {
      throw new Error(
        `Cannot get root for world '${this.id}'. World is not loaded yet.`,
      );
    }
    return worldRoot;
  }

  /**
   * Returns the root world object as a signal. This signal is set/updated any time we
   * receive a WorldObject update from the backend. This includes FullUpdates and StateUpdates
   * (aka pose-only updates).
   */
  getRootObjectSignal(): Signal<WorldObject | undefined> {
    return this.worldRoot;
  }

  // Returns true if the world contains the given object.
  hasObject(oid: ObjectId): boolean {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return false;
    }
    return this.worldLookup.objects.has(oid);
  }

  // Returns the object with the given id.
  getObject(oid: ObjectId): WorldObject | undefined {
    return this.worldLookup?.objects.get(oid);
  }

  // Returns a Signal<WorldObject | undefined> that updates whenever
  // the WorldObject with the given oid changes in the world. This occurs
  // whenever an update is received from the backend and the WorldObject
  // tree is rebuilt.
  createObjectSignal(
    objectId: Signal<ObjectId>,
  ): Signal<WorldObject | undefined> {
    return computed<WorldObject | undefined>(() => {
      const rootObject = this.getRootObjectSignal();
      return rootObject()?.findObjectInSubtree(objectId());
    });
  }

  // Returns a Signal<Entity | undefined> that updates whenever the
  // Entity with id `entityId` in object with id `objectId` changes.
  // This occurs whenever an update is received from the backend and the
  // WorldObject tree is rebuilt.
  createEntitySignal(
    objectId: Signal<ObjectId>,
    entityId: Signal<EntityId>,
  ): Signal<Entity | undefined> {
    return computed<Entity | undefined>(() => {
      const rootObject = this.getRootObjectSignal();
      const object = rootObject()?.findObjectInSubtree(objectId());
      if (!object) {
        return undefined;
      }
      return object.getEntityById(entityId());
    });
  }

  /**
   * Traverses the `WorldModel` depth-first, starting at the root, and finds a
   * WorldObject that matches `predicate`.
   *
   *
   * @param predicate A function that takes a WorldObject and returns `true` if
   *                  that object is the desired one.
   *
   * @returns The first object `obj` for which `predicate(obj)` returns `true`,
   *         or `undefined` if none of the objects in the world match
   *         `predicate`.
   */
  findObjectBy(
    predicate: (obj: WorldObject) => boolean,
  ): WorldObject | undefined {
    // The Generator that `traverseChildren()` returns does not have `find()`,
    // so we roll our own.
    for (const obj of this.getRootObject().traverseChildren()) {
      if (predicate(obj)) {
        return obj;
      }
    }
    return undefined;
  }

  // Returns true if the world contains the given frame.
  hasFrame(fid: FrameId): boolean {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return false;
    }
    return this.worldLookup.frames.has(fid);
  }

  // Returns the frame for the given id.
  getFrame(fid: FrameId): FrameEntity | undefined {
    return this.worldLookup?.frames.get(fid);
  }

  hasEntity(eid: EntityId): boolean {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return false;
    }
    return this.worldLookup.entities.has(eid);
  }

  getEntity(eid: EntityId): Entity | undefined {
    return this.worldLookup?.entities.get(eid);
  }

  // Get a list of all object ids.
  getObjectIds(sorted?: boolean): ObjectId[] {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return [];
    }
    const res = [...this.worldLookup.objects.keys()];
    // Return object ids sorted by their name.
    // Always keep root at the top due to its specialty.
    if (sorted) {
      res.sort((id1: ObjectId, id2: ObjectId) => {
        if (id1 === WorldModel.ROOT_OBJECT_ID) return -1;
        if (id2 === WorldModel.ROOT_OBJECT_ID) return 1;
        const name1 = this.worldLookup?.objFrameNames.get(id1);
        const name2 = this.worldLookup?.objFrameNames.get(id2);
        if (name1 === undefined || name2 === undefined) {
          return 0;
        }
        return name1.localeCompare(name2);
      });
    }
    return res;
  }

  getObjectOrFrame(
    ofid: ObjectOrFrameId,
  ): WorldObject | FrameEntity | undefined {
    return this.getObject(ofid) ?? this.getFrame(ofid);
  }

  /**
   * @returns the reference string of the object/frame, or undefined if no
   * object or frame with that id exists in the world.
   */
  getReferenceString(ofid: ObjectOrFrameId): string | undefined {
    return this.getObjectOrFrame(ofid)?.getReferenceString();
  }

  /**
   * Gets the name of an object or frame.
   *
   * @returns the name, or undefined if no object or frame with that id exists
   * in the world (or if the world is not loaded yet).
   */
  getName(ofid: ObjectOrFrameId): string | undefined {
    return this.worldLookup?.objFrameNames.get(ofid);
  }

  /**
   * Gets a resource name by its ID. To ensure the object/frame is uniquely
   * identifiable, the name is prepended with the parents' name when one is
   * available (unless the parent is the root element).
   *
   * @returns the qualified name, or undefined if no object or frame with that
   * id exists in the world.
   */
  getFullName(ofid: ObjectOrFrameId): string | undefined {
    const name = this.getName(ofid);
    if (name === undefined) {
      return undefined;
    }
    const parentId = this.getParentId(ofid);
    // This relies on a direct comparison instead of some numeric comparison
    // to be resilient against changes to these values.
    if (
      parentId === WorldModel.NULL_OBJECT_ID ||
      parentId === WorldModel.ROOT_OBJECT_ID
    ) {
      return name;
    }
    const parentName = this.getName(parentId);
    return parentName === undefined ? name : `${parentName}.${name}`;
  }

  // Get a list of  all frame ids.
  getFrameIds(sorted?: boolean): FrameId[] {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return [];
    }
    const res = [...this.worldLookup.frames.keys()];
    // Return frame ids sorted by their name.
    if (sorted) {
      res.sort((id1: FrameId, id2: FrameId) => {
        const name1 = this.worldLookup?.objFrameNames.get(id1);
        const name2 = this.worldLookup?.objFrameNames.get(id2);
        if (name1 === undefined || name2 === undefined) {
          return 0;
        }
        return name1.localeCompare(name2);
      });
    }
    return res;
  }

  // Get a list of frame ids under given object id.
  getObjectFrameIds(oid: ObjectId, sorted?: boolean): FrameId[] {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return [];
    }
    const objFrames = this.worldLookup.objFrames.get(oid);
    if (!objFrames) {
      console.warn(`Object ${oid} does not exist!`);
      return [];
    }
    const res = [...objFrames.keys()];
    // Return object ids sorted by their name.
    if (sorted) {
      res.sort((id1: FrameId, id2: FrameId) => {
        const name1 = this.worldLookup?.objFrameNames.get(id1);
        const name2 = this.worldLookup?.objFrameNames.get(id2);
        if (name1 === undefined || name2 === undefined) {
          return 0;
        }
        return name1.localeCompare(name2);
      });
    }
    return res;
  }

  /**
   * @param name The name of an object.
   * @returns the object id of the or undefined if name doesn't exist.
   */
  getObjectIdByName(name: string): ObjectId | undefined {
    const oids = this.getObjectIds().filter(
      (oid) => this.getName(oid) === name,
    );
    if (oids.length === 0) return undefined;
    return oids[0];
  }

  /**
   * Attempts to get an object ID by its name from the local cache.
   * If it is not found, forces a sync with the backend to hydrate the cache
   * and tries again.
   *
   * @param name The name of the object.
   * @returns A promise resolving to the object id, or undefined if it doesn't exist.
   */
  async getHydratedObjectIdByName(name: string): Promise<ObjectId | undefined> {
    let oid = this.getObjectIdByName(name);
    if (oid === undefined) {
      await this.sync();
      oid = this.getObjectIdByName(name);
    }
    return oid;
  }

  /**
   * Finds an object with `objectName` in `worldModel`.
   *
   * @returns the object with `objectName`, or undefined if there is no object
   * with that name.
   */
  getObjectByName(objectName: string): WorldObject | undefined {
    const objectId = this.getObjectIdByName(objectName);
    if (objectId === undefined) {
      return undefined;
    }
    return this.getObject(objectId);
  }

  // Get the parent id of given input id which can be either object or frame
  // id.
  // If the input frame/object's parent is a frame, returns its frame ID.
  // Otherwise, returns the object ID of its parent object.
  getParentId(rid: ObjectOrFrameId): ObjectOrFrameId {
    if (!this.worldLookup) {
      console.warn(`World is not loaded yet!`);
      return WorldModel.NULL_OBJECT_ID;
    }

    if (this.worldLookup.frames.has(rid)) {
      const frame = this.getFrame(rid);
      if (!frame) {
        return WorldModel.NULL_OBJECT_ID;
      }
      return frame.getParentFrameInObject()?.id ?? frame.parentObject.id;
    }

    // If an object was reparented to a frame (e.g., via `reparent_object` with
    // a `parent_frame`), its root entity's parent is that frame rather than the
    // parent object.
    const parentEntity = this.worldLookup.objects
      .get(rid)
      ?.getRootEntity()
      ?.getParent();
    if (parentEntity && this.worldLookup.frames.has(parentEntity.id)) {
      return parentEntity.id;
    }

    const parent = this.worldLookup.objects.get(rid)?.getParent()?.id;
    // Either object does not exist or parent does not exist (means that input
    // is root).
    return parent ?? WorldModel.NULL_OBJECT_ID;
  }

  // Get child object ids of a given input object id.
  getChildIds(oid: ObjectId, sorted?: boolean): ObjectOrFrameId[] {
    const obj = this.worldLookup?.objects.get(oid);
    if (!obj) {
      return [];
    }
    const res = [...obj.getChildren()].map((o) => o.id);
    if (sorted) {
      res.sort((id1: ObjectOrFrameId, id2: ObjectOrFrameId) => {
        const name1 = this.worldLookup?.objFrameNames.get(id1);
        const name2 = this.worldLookup?.objFrameNames.get(id2);
        if (name1 === undefined || name2 === undefined) {
          return 0;
        }
        return name1.localeCompare(name2);
      });
    }
    return res;
  }

  // Get all entity IDs contained in an object
  getObjectEntityIds(oid: ObjectId): EntityId[] | undefined {
    const obj = this.worldLookup?.objects.get(oid);
    if (!obj) return undefined;
    return [...obj.getEntities()].map((e) => e.id);
  }

  getObjectVisualGeometries(oid: ObjectId): Geometry[] {
    const obj = this.worldLookup?.objects.get(oid);
    if (!obj) return [];
    const res = [];
    for (const e of obj.getEntities()) {
      res.push(...(e.asLink()?.getVisualGeometries() ?? []));
    }
    return res;
  }

  /**
   * Returns true if the input object or frame id is considered 6-dof movable.
   */
  isSixDofMovable(oid: ObjectId): boolean {
    // All objects are movable
    if (this.worldLookup?.objects.has(oid)) return true;

    // Frames are movable if all child entities within the same object are
    // frames.
    const frame = this.worldLookup?.entities.get(oid);
    if (frame instanceof FrameEntity) {
      const childEntities = [...frame.getChildren()];
      while (childEntities.length > 0) {
        const child = childEntities.shift();
        if (child && child.parentObject.id === frame.parentObject.id) {
          if (!(child instanceof FrameEntity)) {
            return false;
          }

          childEntities.push(...child.getChildren());
        }
      }

      return true;
    }

    return false;
  }

  isObjectResource(oid: ObjectId): boolean {
    return !!this.worldLookup?.isObjectResource(oid);
  }

  // Returns the first camera sensor of the object.
  getObjectCamera(oid: ObjectId): CameraSensor | undefined {
    return this.worldLookup?.getObjectCamera(oid);
  }

  // Returns true if any entity under the object have camera sensor component.
  // Note entities of child objects does not count, this only checks entities
  // that compose the object.
  isObjectCamera(oid: ObjectId): boolean {
    return !!this.worldLookup?.isObjectCamera(oid);
  }

  isObjectStatic(oid: ObjectId): boolean {
    return !!this.getObject(oid)?.isStatic;
  }

  isObjectDisabled(oid: ObjectId): boolean {
    return !!this.getObject(oid)?.isDisabled;
  }

  getCollisionRules(): WorldCollisionRule[] {
    return this.worldBackend.getCollisionRules();
  }

  /**
   * Returns the relative pose of this entity as it's defined in the attachment
   * graph.
   */
  getLocalTrans(eid: EntityId): three.Matrix4 | undefined {
    const worldLookup = this.worldLookup;
    if (!worldLookup) return undefined;

    return worldLookup.entityInitialTrans.get(worldLookup.toEntityId(eid));
  }

  /**
   * Computes relative transformation between two entities based on data stored
   * in the proto.
   */
  getRelativeTrans(eid1: EntityId, eid2: EntityId): three.Matrix4 | undefined {
    const worldLookup = this.worldLookup;
    if (!worldLookup) return undefined;

    const wTrans1 = worldLookup.getWorldTrans(worldLookup.toEntityId(eid1));
    const wTrans2 = worldLookup.getWorldTrans(worldLookup.toEntityId(eid2));
    if (!wTrans1 || !wTrans2) return undefined;

    return mathUtils.getRelativeTransform(wTrans2, wTrans1);
  }

  /**
   * Returns the relative pose from one entity to another or undefined if
   * the relative transform could not be retrieved or any of the entities supplied
   * are invalid (equal to the NULL_ENTITY_ID=0)
   */
  getRelativePose(eid1: EntityId, eid2: EntityId): Pose | undefined {
    const trans = this.getRelativeTrans(eid1, eid2);
    if (!trans) return undefined;
    return mathUtils.poseFromMatrix(trans);
  }

  // Traverse the model tree from a given node and return the list of
  // object of frame ids encountered during traversal. There is no traverse
  // order for siblings as they are not ordered originally.
  traverseDFS(
    oid: ObjectId | FrameId,
    includeFrames: boolean,
  ): ObjectOrFrameId[] {
    const object = this.getObject(oid);
    const frame = this.getFrame(oid);
    if (!object && !frame) return [];

    let res = [oid];

    if (includeFrames) {
      for (const f of object?.getChildFrames() ?? []) {
        res.push(f.id);
      }
      for (const f of frame?.getChildFrames() ?? []) {
        res.push(f.id);
      }
    }

    for (const o of object?.getChildren() ?? []) {
      res = res.concat(this.traverseDFS(o.id, includeFrames));
    }
    return res;
  }

  getObjectOrFrameByReference(
    ref: string,
  ): WorldObject | FrameEntity | undefined {
    if (!ref) return undefined;
    const includeFrames = ref.includes(REFERENCE_STRING_SEPARATOR);
    const ids = this.traverseDFS(WorldModel.ROOT_OBJECT_ID, includeFrames);
    for (const id of ids) {
      const frame = this.getFrame(id);
      if (frame && frame.getReferenceString() === ref) {
        return frame;
      }
      const obj = this.getObject(id);
      if (obj && obj.getReferenceString() === ref) {
        return obj;
      }
    }
    return undefined;
  }
  /**
   * Converts an object, frame or entity id to its corresponding entity id (the
   * root entity for objects). Returns undefined if the id is unknown or the
   * world is not loaded.
   */
  toEntityId(id: ObjectOrFrameId | EntityId): EntityId | undefined {
    const eid = this.worldLookup?.toEntityId(id);
    return eid === worldIdUtils.NULL_ENTITY_ID ? undefined : eid;
  }

  /**
   * @returns Whether entity eid1 is ancestor of entity eid2 in this model.
   *
   * NOTE: eid1 has to be ancestor of eid2, if eid1 equals eid2 this returns
   * false.
   */
  isAncestorEntity(eid1: EntityId, eid2: EntityId): boolean {
    return this.worldLookup?.isAncestorEntity(eid1, eid2) ?? false;
  }

  /**
   * @returns Whether object of id1 is ancestor of object of id2 in this model.
   *
   * NOTE: id1 has to be ancestor of id2, if id1 equals id2 this returns
   * false.
   */
  isAncestor(id1: ObjectId, id2: ObjectId): boolean {
    if (!this.hasObject(id1) || !this.hasObject(id2)) {
      return false;
    }
    while ((id2 = this.getParentId(id2))) {
      if (id1 === id2) {
        return true;
      }
    }
    return false;
  }

  async updateTransform(
    id: number,
    refId: number,
    relativePose: ReadonlyPose,
    updateId?: number,
  ): Promise<void> {
    await this.worldBackend.updateTransform(id, refId, relativePose, updateId);
  }

  async batchUpdateTransform(
    worldPoses: Map<number, ReadonlyPose>,
  ): Promise<void> {
    await this.worldBackend.batchUpdateTransform(worldPoses);
  }

  async updateEntityGeometry(
    oid: number,
    eid: number,
    geometryUpdate: GeometryUpdate,
  ): Promise<void> {
    await this.worldBackend.updateEntityGeometry(oid, eid, geometryUpdate);
  }

  async createFrame(
    parentId: ObjectOrFrameId,
    name: string,
    relativePose: ReadonlyPose,
  ): Promise<CreatedFrame> {
    return await this.worldBackend.createFrame(parentId, name, relativePose);
  }

  async deleteEntity(
    entityId: number,
    policy: DeleteEntityPolicy,
  ): Promise<void> {
    await this.worldBackend.deleteEntity(entityId, policy);
  }

  async updateFrameProps(
    fId: number,
    isAttachmentFrame: boolean,
  ): Promise<void> {
    await this.worldBackend.updateFrameProps(fId, isAttachmentFrame);
  }

  async deleteObject(oId: number, policy: DeleteObjectPolicy): Promise<void> {
    await this.worldBackend.deleteObject(oId, policy);
  }

  async updateObjectParentAndTransform(
    id: ObjectId,
    pId: ObjectId,
    targetEntityName?: string,
    worldPose?: ReadonlyPose,
  ): Promise<void> {
    await this.worldBackend.updateObjectParentAndTransform(
      id,
      pId,
      targetEntityName,
      worldPose,
    );
  }

  async updateEntityParent(
    entityId: number,
    newParentId: number,
  ): Promise<void> {
    await this.worldBackend.updateEntityParent(entityId, newParentId);
  }

  async updateObjectName(oId: number, name: string): Promise<void> {
    await this.worldBackend.updateObjectName(oId, name);
  }

  async updateEntityName(entityId: number, name: string): Promise<void> {
    await this.worldBackend.updateEntityName(entityId, name);
  }

  async updateWorldCollisionSettings(
    collisionsRules: WorldCollisionRule[],
  ): Promise<void> {
    await this.worldBackend.updateWorldCollisionSettings(collisionsRules);
  }

  async updateObjectPhysics(
    oId: number,
    physicsUpdate: PhysicsUpdate,
  ): Promise<void> {
    await this.worldBackend.updateObjectPhysics(oId, physicsUpdate);
  }

  async updateJoints(robotId: number, jvals: number[]): Promise<void> {
    await this.worldBackend.updateJoints(robotId, jvals);
  }

  async updateJointApplicationLimits(
    robotId: number,
    newLimits: JointLimitsUpdate,
  ): Promise<JointLimits | undefined> {
    return await this.worldBackend.updateJointApplicationLimits(
      robotId,
      newLimits,
    );
  }

  async batchDelete(oids: number[], fids: number[]): Promise<void> {
    await this.worldBackend.batchDelete(oids, fids);
  }
}
