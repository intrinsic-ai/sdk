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

import {Entity, FrameEntity, JointEntity, LinkEntity} from './entity';
import * as worldIdUtils from './world_ids';
import {EntityId, FrameId, ObjectId} from './world_ids';

interface PprResourceName {
  source: 'resource';
  name: string;
}

/**
 * A way to tag the asset source of a world object.
 */
export type PprSource = PprResourceName;

/**
 * WorldObjectData is the data needed to instantiate an object in the world
 * model. It usually is sourced by translating a backend source such as the
 * world service into this interface and used via a `WorldObjectBuilder`.
 */
export declare interface WorldObjectData {
  id: ObjectId;
  name: string;
  parentEntityId: EntityId;
  rootEntityId: EntityId;
  parentObjectId: ObjectId;
  parentTThis: three.Matrix4;
  isStatic: boolean;
  isDisabled: boolean;
  children: string[];
  ppr?: PprSource;
  /**
   * User data associated with the object. Values are backend-specific (e.g.
   * `google.protobuf.Any` protos), so consumers must narrow them before use.
   */
  userDataMap?: Map<string, unknown>;
}

/**
 * An implementation of intrinsic/world/proto/robot_payload.proto
 */
export declare interface RobotPayload {
  massKg: number;
  tipTCoG: three.Matrix4;
  inertia: three.Matrix3;
}

/**
 * A world object in the world model. World objects are a collection of entities
 * that may represent physical geometry, joints, or sensors. They also possibly
 * define a set of coordinate frames. World objects are implicitly part of a
 * tree of objects defined by the parents of their root entities, which is
 * modeled via accessors in this class.
 *
 * Structure:
 * On the frontend, we can represent an arbitrary 3d scene using a WorldObject tree. WorldObject
 * trees should be created using the `WorldObjectBuilder` class. These are the main requirements
 * to create a well-formed WorldObject tree:
 *
 * 1. The id of the WorldObject at the root of the tree must be ROOT_OBJECT_ID, and its root entity
 * must have id ROOT_ENTITY_ID.
 *
 * 2. The ids of all WorldObjects and Entities must be unique within the tree. This means that a WorldObject cannot have the
 * same ID as some other WorldObject in the tree, or of any Entity within the tree.
 *
 * 3. Every WorldObject contains a tree of Entities, with exactly one root Entity. Other than the root Entity of the
 * root WorldObject (which has no parent), a WorldObject's root Entity is parented to some Entity in another WorldObject.
 * This forms a parent-child relation between the two WorldObjects.
 * Note that only the root entity of the root WorldObject should have id ROOT_ENTITY_ID.
 *
 * 4. The root Entity of every WorldObject (including the root WorldObject) must be a LinkEntity. Building a
 * WorldObject whose root Entity is of another type (e.g. a FrameEntity) fails. The rest of the world model, as
 * well as the rendering code, relies on this. In particular, an object's root Entity is never a frame.
 */
export class WorldObject {
  /** Unique ID of this object in the world. */
  readonly id: ObjectId;

  /** Human-readable name for this object in the world. */
  readonly name: string;

  /**
   * The parent entity of this object.
   *
   * TODO: b/331242474 -- remove this in favor of `this.rootEntity.getParent()?.id`.
   */
  private parentEntityId: EntityId;
  getParentEntityId(): EntityId | undefined {
    return this.parentEntityId === worldIdUtils.NULL_ENTITY_ID
      ? undefined
      : this.parentEntityId;
  }

  /**
   * The parent object of this object (the parent entity's parent object).
   *
   * TODO: b/331242474 -- remove this in favor of `this.getParent()?.id`.
   */
  private parentObjectId: ObjectId;
  getParentObjectId(): ObjectId | undefined {
    return this.parentObjectId === worldIdUtils.NULL_OBJECT_ID
      ? undefined
      : this.parentObjectId;
  }

  /**
   * The transform between the parent entity within the parent object and the
   * root entity for this object.
   */
  readonly parentTThis: three.Matrix4;

  /** The root entity for this object. */
  private rootEntity?: LinkEntity;
  setRootEntityInternal(rootEntity: LinkEntity) {
    this.rootEntity = rootEntity;
  }

  /** Returns the root entity of this object, which is always a link. */
  getRootEntity(): LinkEntity {
    if (!this.rootEntity) {
      throw new Error(`No root entity found for object ${this.name}`);
    }
    return this.rootEntity;
  }

  /** The entities associated with this object. */
  private readonly entities = new Map<EntityId, Entity>();

  /** Use only in world_object_builder.ts to build a WorldObject. */
  getEntitiesInternal(): Map<EntityId, Entity> {
    return this.entities;
  }

  /** Returns the entity object for the given ID. */
  getEntityById(id: EntityId): Entity | undefined {
    return this.entities.get(id);
  }

  /** Returns all of the entities for this object. */
  *getEntities(): Generator<Entity> {
    for (const e of this.entities.values()) {
      yield e;
    }
  }

  /** Returns the entity with the given name, if it exists. */
  getEntityByName(name: string): Entity | undefined {
    for (const entity of this.entities.values()) {
      if (entity.name === name) {
        return entity;
      }
    }
    return undefined;
  }

  /** Returns true if this object has an entity with the given name. */
  hasEntityWithName(name: string): boolean {
    return this.getEntityByName(name) !== undefined;
  }

  /** Returns the frame entity of this object with the given id. */
  getFrameById(id: FrameId): FrameEntity | undefined {
    const entity = this.entities.get(id);
    return entity?.isFrame() ? entity : undefined;
  }

  /** Returns all of the frame entities of this object. */
  *getFrames(): Generator<FrameEntity> {
    for (const entity of this.entities.values()) {
      if (entity.isFrame()) {
        yield entity;
      }
    }
  }

  /** Returns the frame with the given name, if it exists. */
  getFrameByName(name: string): FrameEntity | undefined {
    const entity = this.getEntityByName(name);
    return entity?.isFrame() ? entity : undefined;
  }

  /**
   * Returns the frames that are direct children of this object, i.e. the
   * frames whose parent entity is not itself a frame of this object.
   */
  *getChildFrames(): Generator<FrameEntity> {
    for (const frame of this.getFrames()) {
      if (!frame.getParentFrameInObject()) {
        yield frame;
      }
    }
  }

  /** Returns a all direct child objects of this object. */
  *getChildren(): Generator<WorldObject> {
    for (const c of this.childObjects) {
      yield c;
    }
  }

  /** Returns the child object of the given name. */
  getChildByName(name: string): WorldObject | undefined {
    return this.childObjects.find((c) => c.name === name);
  }

  /**
   * Finds an object by ID in this object's subtree.
   * Returns this object if the ID matches.
   */
  findObjectInSubtree(id: ObjectId): WorldObject | undefined {
    for (const obj of this.traverseChildren()) {
      if (obj.id === id) {
        return obj;
      }
    }
    return undefined;
  }

  /**
   * Traverses the object structure and returns every child that has this
   * object as an ancestor (including this object, which it returns first).
   */
  *traverseChildren(): Generator<WorldObject> {
    yield this;
    for (const c of this.getChildren()) {
      yield* c.traverseChildren();
    }
  }

  /**
   * The parent object for this object. Is undefined if this object is the root
   * object.
   */
  getParent(): WorldObject | undefined {
    return this.parentObject;
  }
  // We assign to this property in children of this object.
  private parentObject: WorldObject | undefined;

  /**
   * Returns true if the passed object is an ancestor of this object. An object
   * is not considered an ancestor of itself.
   */
  hasAncestor(obj: WorldObject): boolean {
    const ancestor = this.getParent();
    if (!ancestor) return false;
    if (ancestor === obj) return true;
    return ancestor.hasAncestor(obj);
  }

  /** Returns the transform node reference string for this object. */
  getReferenceString(): string {
    return this.name;
  }

  /** Stores the resource name that defined this object. */
  private readonly ppr?: PprSource;

  /** All user data associated with this object. */
  private readonly userDataMap?: Map<string, unknown>;

  /**
   * Returns the user data associated with this object or undefined. Values are
   * backend-specific, so callers must narrow them before use.
   */
  getUserdataMap(): Map<string, unknown> | undefined {
    return this.userDataMap;
  }

  /**
   * Returns truthy if this object originated from a resource instance and
   * narrows type to provide resource name.
   */
  isResource(): this is {getResourceName: () => string} & this {
    return this.ppr?.source === 'resource';
  }

  /**
   * Returns the resource name for this object, if this object originated from a
   * resource instance.
   */
  getResourceName(): string | undefined {
    return this.isResource() ? this.ppr?.name : undefined;
  }

  /** Returns truthy if type is kinematic object and narrows type. */
  isKinematic(): this is KinematicWorldObject {
    return this instanceof KinematicWorldObject;
  }

  /** Returns this object as a kinematic object if possible. */
  asKinematic(): KinematicWorldObject | undefined {
    return this.isKinematic() ? this : undefined;
  }

  computeWorldTThis(): three.Matrix4 {
    return this.getRootEntity().computeWorldTThis();
  }

  /**
   * A simple flag for whether or not this object is static in our physics model.
   */
  readonly isStatic: boolean;

  /**
   * A simple flag for whether or not this object is disabled in simulation.
   */
  readonly isDisabled: boolean;

  /**
   * Returns the child objects for this object.
   * @param params Parameters common to all objects.
   * @param childObjects The children of this object.
   */
  constructor(
    params: WorldObjectData,
    protected readonly childObjects: WorldObject[],
  ) {
    this.id = params.id;
    this.name = params.name;
    this.parentEntityId = params.parentEntityId;
    this.parentObjectId = params.parentObjectId;
    this.parentTThis = params.parentTThis;
    this.ppr = params.ppr;
    this.isStatic = params.isStatic;
    this.isDisabled = params.isDisabled;
    this.userDataMap = params.userDataMap;
    for (const child of this.childObjects) {
      if (child.parentObjectId === this.id) {
        child.parentObject = this;
      } else {
        console.warn(
          `Child object added to object with ID ${this.id} with different parentObjectId (${child.parentObjectId})`,
        );
      }
    }
  }

  /**
   * Reparents this object to the target object. Throws an error under the
   * following conditions:
   *   - this object is an ancestor of the target object
   *   - the target entity is not an entity of the target object.
   */
  reparentTo(target: WorldObject, targetEntity?: Entity) {
    // First check to see if this object is an ancestor of the target.
    let targetAncestor: WorldObject | undefined = target;
    while (targetAncestor !== undefined) {
      if (targetAncestor === this) {
        throw new Error(
          `Cannot reparent '${this.name}' to '${target.name}': '${this.name}' is already an ancestor of '${target.name}' -- would create cycle.`,
        );
      }
      targetAncestor = targetAncestor.getParent();
    }

    // Handle entities first -- if there's a target entity, then reparent the
    // root entity of this object to it. Otherwise, use the target object's
    // base entity.
    const entityTarget = targetEntity ?? target.getRootEntity();
    if (!target.entities.has(entityTarget.id)) {
      throw new Error(
        `Cannot reparent root entity '${
          this.getRootEntity().name
        }' of object '${this.name}' to target object '${
          target.name
        }': target object does not contain target entity '${
          entityTarget.name
        }'`,
      );
    }
    this.getRootEntity().reparentTo(entityTarget);

    // First, remove this object from its parent's list of children.
    const oldParent = this.getParent();
    if (oldParent) {
      const childIdx = oldParent.childObjects.findIndex((o) => o === this);
      if (childIdx < 0) {
        throw new Error(
          `INTERNAL -- WorldObject.reparent: '${this.name}' not found as child of old parent '${oldParent.name}'`,
        );
      }

      oldParent.childObjects.splice(childIdx, 1);
    }

    // Next, set the current parent to the new target.
    this.parentObject = target;
    this.parentObjectId = target.id;
    this.parentEntityId = entityTarget.id;

    // Finally, add this entity to the list of children of the new parent.
    target.childObjects.splice(0, 0, this);
  }
}

/**
 * KinematicWorldObjectData is additional data needed to specialize an object as
 * a KinematicWorldObject in the world model. It usually is sourced by
 * adding properties to a WorldObjectData, using a backend source such as the
 * world service into this interface and via a `WorldObjectBuilder`.
 */
export declare interface KinematicWorldObjectData extends WorldObjectData {
  controlFrequency: number;
  jointOrdering: string[];
  isoFlangeFrames: string[];
  namedConfigurations: Map<string, number[]>;
  payload?: RobotPayload;
}

/**
 * A kinematic object in the world model. Kinematic objects are objects that
 * have one or more joints and a defined kinematics configuration.
 */
export class KinematicWorldObject extends WorldObject {
  /** Control frequency used for this kinematic object. */
  readonly controlFrequency: number;

  /** Named joint configurations used for this kinematic object. */
  readonly namedConfigurations: Map<string, number[]>;

  /** List of joints corresponding to this kinematic object. */
  private readonly joints: JointEntity[] = [];
  private readonly links: LinkEntity[] = [];

  /** Use only in world_object_builder.ts to build a KinematicWorldObject. */
  getJointsInternal(): JointEntity[] {
    return this.joints;
  }

  /** Use only in world_object_builder.ts to build a KinematicWorldObject. */
  getLinksInternal(): LinkEntity[] {
    return this.links;
  }

  *getJoints(): Generator<JointEntity, void, undefined> {
    yield* this.joints;
  }

  *getJointValues(): Generator<number> {
    for (const j of this.getJoints()) {
      yield j.value;
    }
  }

  /**
   * Returns a list of link entities that belong to this object, in breadth-first
   * traversal order.
   */
  *getLinks(): Generator<LinkEntity, void, undefined> {
    yield* this.links;
  }

  /**
   * List of frames according to the ISO 9787 standard. Not every kinematic
   * object has flange frames, but one flange frame for every "robot arm" can
   * be expected.
   */
  private readonly flangeFrames: FrameEntity[] = [];
  /** Use only in world_object_builder.ts to build a KinematicWorldObject. */
  getFlangeFramesInternal(): FrameEntity[] {
    return this.flangeFrames;
  }
  *getFlangeFrames(): Generator<FrameEntity, void, undefined> {
    yield* this.flangeFrames;
  }

  private mountedPayload: RobotPayload | undefined;
  getMountedPayload(): RobotPayload | undefined {
    return this.mountedPayload;
  }

  /**
   * This centralizes a specific heuristic used throughout the frontend to
   * determine whether something is a robot or not.
   *
   * This is not necessarily a fully correct heuristic, but works well enough
   * for now.
   */
  isRobot(): boolean {
    return this.flangeFrames.length > 0;
  }

  /**
   * Returns the child objects for this object.
   * @param data The data needed to define this KinematicWorldObject.
   * @param childObjects The children of this object.
   */
  constructor(
    data: KinematicWorldObjectData,
    protected override readonly childObjects: WorldObject[],
  ) {
    super(data, childObjects);
    this.controlFrequency = data.controlFrequency;

    // Deep copy the payload.
    if (data.payload) {
      this.mountedPayload = {
        massKg: data.payload.massKg,
        tipTCoG: data.payload.tipTCoG.clone(),
        inertia: data.payload.inertia.clone(),
      };
    }

    // Deep copy the configurations.
    this.namedConfigurations = new Map();
    for (const [k, v] of data.namedConfigurations.entries()) {
      this.namedConfigurations.set(k, [...v]);
    }
  }
}
