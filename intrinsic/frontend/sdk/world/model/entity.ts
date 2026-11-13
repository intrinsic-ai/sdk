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

import {
  COLLISION_GEOMETRY_KEY,
  Geometry,
  GeometryShape,
  VISUAL_GEOMETRY_KEY,
} from './geometry';
import {MaterialProperties} from './material_properties';
import * as worldIdUtils from './world_ids';
import {EntityId, ObjectId} from './world_ids';
import {PprSource, WorldObject} from './world_object';

/**
 * EntityData is the data needed to instantiate an entity in the world
 * model. It usually is sourced by translating a backend source such as the
 * world service into this interface and used via a `WorldObjectBuilder`.
 */
export declare interface EntityData {
  id: EntityId;
  name: string;
  parentObjectId: ObjectId;
  parentTThis: three.Matrix4;
  parentId: EntityId;
  ppr?: PprSource;
}

/**
 * A union of the concrete subclasses of Entity
 */
export type ConcreteEntity =
  LinkEntity | JointEntity | FrameEntity | SensorEntity;

/** All possible entity types. */
export const ALL_ENTITY_TYPES = ['Link', 'Joint', 'Frame', 'Sensor'] as const;

/**
 * A union of the possible entity type strings. Useful
 * for switching over with `entity.getEntityType()` to handle all cases.
 */
export type EntityType = (typeof ALL_ENTITY_TYPES)[number];

/**
 * A world entity is a compositional atom of a world object. An object may have
 * many entities, each of which defines either a joint, link, or sensor.  All
 * entities are posed relative to their parent, defining a tree of entities.
 *
 * Every WorldObject contains a tree of Entities. One Entity is the root, which has
 * no parent. The ids of all Entities must be unique within the tree (and cannot collide
 * with WorldObject ids, either). The id of the root Entity of the root WorldObject must
 * be set to ROOT_ENTITY_ID.
 */
export abstract class Entity {
  /** The unique identifier for this entity in a world. */
  readonly id: EntityId;

  /** The name of this entity, as defined in the world. */
  readonly name: string;

  abstract getEntityType(): EntityType;

  /**
   * The parent entity id of this entity. May be root or an entity in a
   * different object than this entity's parent object.
   *
   * TODO: b/331242474 -- remove these in favor of `getParent()?.id`.
   */
  private parentId: EntityId;
  getParentId(): EntityId | undefined {
    return this.parentId === worldIdUtils.NULL_ENTITY_ID
      ? undefined
      : this.parentId;
  }

  /** The transformation between the parent entity and this entity. */
  parentTThis: three.Matrix4;

  /**
   * The parent object ID for this entity.
   *
   * TODO: b/331242474 -- remove these in favor of `parentObject.id`.
   */
  private readonly parentObjectId: ObjectId;
  getParentObjectId(): ObjectId {
    return this.parentObjectId;
  }

  /**
   * The parent entity for this entity. Is undefined if this entity is the root
   * entity.
   */
  getParent(): Entity | undefined {
    return this.parentEntity;
  }

  // We assign to this property in children of this entity.
  private parentEntity: Entity | undefined;

  /** Returns the children Entities of this entity. */
  *getChildren(): Generator<Entity> {
    for (const c of this.children) {
      yield c;
    }
  }

  computeWorldTThis(): three.Matrix4 {
    return (
      this.getParent()?.computeWorldTThis()?.multiply(this.parentTThis) ??
      new three.Matrix4()
    );
  }

  /** Stores the resource name that defined this object. */
  private readonly ppr?: PprSource;

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

  /**
   * Constructs an entity with the properties common to all entities.
   *
   * @param params The entity data.
   * @param parentObject The parent object for this entity.
   *     CAUTION: might not be fully initialized -- best to avoid using it in
   *     the constructor
   * @param children The children for this entity.
   */
  protected constructor(
    params: EntityData,
    readonly parentObject: WorldObject,
    private readonly children: Entity[],
  ) {
    this.id = params.id;
    this.name = params.name;
    this.parentObjectId = params.parentObjectId;
    this.parentTThis = params.parentTThis;
    this.parentId = params.parentId;
    this.parentObjectId = params.parentObjectId;
    this.ppr = params.ppr;

    // Connect this entity to its children.
    for (const child of children) {
      if (child.parentId !== this.id) {
        throw new Error(
          `Entity ${child.id} passed to entity ${this.id} as child, but its parent is ${child.parentId}`,
        );
      }
      child.parentEntity = this;
    }
  }

  /** Returns truthy and narrows type to link */
  isLink(): this is LinkEntity {
    return this instanceof LinkEntity;
  }

  /** Returns this entity as a link if it was created as one */
  asLink(): LinkEntity | undefined {
    return this.isLink() ? this : undefined;
  }

  /** Returns truthy and narrows type to joint */
  isJoint(): this is JointEntity {
    return this instanceof JointEntity;
  }

  /** Returns this entity as a joint if it was created as one */
  asJoint(): JointEntity | undefined {
    return this.isJoint() ? this : undefined;
  }

  /** Returns truthy and narrows type to sensor */
  isSensor(): this is SensorEntity {
    return this instanceof SensorEntity;
  }

  /** Returns this entity as a sensor if it was created as one */
  asSensor(): SensorEntity | undefined {
    return this.isSensor() ? this : undefined;
  }

  /** Returns truthy and narrows type to frame */
  isFrame(): this is FrameEntity {
    return this instanceof FrameEntity;
  }

  /** Returns this entity as a frame if it was created as one */
  asFrame(): FrameEntity | undefined {
    return this.isFrame() ? this : undefined;
  }

  /** Returns this entity as a concrete entity subclass */
  abstract asConcreteEntity(): ConcreteEntity;

  /** Returns true if this entity is the root entity of its parent object. */
  isObjectRoot(): boolean {
    return this.parentObject.getRootEntity().id === this.id;
  }

  /**
   * Reparents the entity to the target entity. Throws an error under the
   * following conditions:
   *   - if this entity is an ancestor of the target entity.
   *   - if this entity is neither the parent object's root entity nor frame.
   */
  reparentTo(target: Entity) {
    // Maybe should be an error?
    if (target === this) return;

    // If this entity is a root entity of its parent object, then we should
    // be reparenting the object to preserve coherence within the model.
    if (!this.isObjectRoot()) {
      const name = this.name;
      const parentName = this.parentObject.name;
      if (!this.isFrame()) {
        throw new Error(
          `Cannot reparent entity '${name}' -- entity is neither the root entity nor a frame of object '${parentName}'.`,
        );
      }
    }

    // First check to see if this entity is an ancestor of the target.
    if (target.isDescendantOf(this)) {
      throw new Error(
        `Cannot reparent entity '${this.name}' to entity '${target.name}': '${this.name}' is already an ancestor of '${target.name}' -- would create cycle.`,
      );
    }

    // First, remove the entity from its parent's list of children.
    const oldParent = this.getParent();
    if (oldParent) {
      const childIdx = oldParent.children.findIndex((e) => e === this);
      if (childIdx < 0) {
        throw new Error(
          `INTERNAL -- Entity.reparent: '${this.name}' not found as child of old parent '${oldParent.name}'`,
        );
      }

      oldParent.children.splice(childIdx, 1);
    }

    // Next, set the current parent to the new target.
    this.parentEntity = target;
    this.parentId = target.id;

    // Finally, add this entity to the list of children of the new parent.
    target.children.splice(0, 0, this);
  }

  /**
   * Returns true if this entity is a descendant of `ancestor`.
   * We do not consider an Entity to be its own descendant (that is,
   * returns false if `ancestor` is `this`)
   */
  isDescendantOf(ancestor: Entity): boolean {
    let current = this.getParent();
    while (current) {
      if (current === ancestor) {
        return true;
      }
      current = current.getParent();
    }
    return false;
  }
}

/**
 * A simple struct containing joint limits as described in the kinematics
 * component of the world
 */
export declare interface JointLimits {
  velocity: number;
  acceleration: number;
  jerk: number;
  effort: number;
  fixed?: {
    upper: number;
    lower: number;
  };
}

/**
 * MotionType describes the type of joint.
 */
export enum MotionType {
  UNDEFINED,
  REVOLUTE,
  PRISMATIC,
  FIXED,
}

/**
 * JointEntityData is additional data needed to specialize an entity as
 * a JointEntity in the world model. It usually is sourced by adding properties
 * to an EntityData, using a backend source such as the world service into this
 * interface and via a `WorldObjectBuilder`.
 */
export declare interface JointEntityData extends EntityData {
  motionType: MotionType;
  axis: three.Vector3;
  rawValue: number;
  systemLimits: JointLimits;
  applicationLimits: JointLimits;
  parentTInboard: three.Matrix4;
  outboardTChild: three.Matrix4;
}

/**
 * A joint entity is a specifc type of entity that contains kinematic
 * information.
 */
export class JointEntity extends Entity {
  readonly motionType: MotionType;
  private readonly axis: three.Vector3;
  value: number;
  readonly systemLimits: JointLimits;
  readonly applicationLimits: JointLimits;
  private readonly parentTInboard: three.Matrix4;
  private readonly outboardTChild: three.Matrix4;

  /**
   * Constructs a joint entity.
   *
   * @param params The joint entity data for this joint
   */
  constructor(
    params: JointEntityData,
    // CAUTION: parent object might not be fully initialized -- best to avoid
    // using it in the constructor
    parentObject: WorldObject,
    children: Entity[],
  ) {
    super(params, parentObject, children);
    this.motionType = params.motionType;
    this.axis = params.axis;
    this.value = params.rawValue;
    this.systemLimits = params.systemLimits;
    this.applicationLimits = params.applicationLimits;
    this.parentTInboard = params.parentTInboard;
    this.outboardTChild = params.outboardTChild;
  }

  override getEntityType(): 'Joint' {
    return 'Joint';
  }

  override asConcreteEntity(): JointEntity {
    return this;
  }

  getAxis(): three.Vector3 {
    return this.axis.clone();
  }

  getParentTInboard(): three.Matrix4 {
    return this.parentTInboard.clone();
  }

  getOutboardTChild(): three.Matrix4 {
    return this.outboardTChild.clone();
  }
}

/**
 * GeometryData is the definition used to specify geometry in the world model.
 * Each link has a number of different geometries, and each geometry is
 * determined by a local transform from the link's pose and a shape describing
 * the geometry.
 */
export declare interface GeometryData {
  refTShape: three.Matrix4;
  shape: GeometryShape;
  materialOverrides?: MaterialProperties;
}

/**
 * LinkEntityData is additional data needed to specialize an entity as
 * a LinkEntity in the world model. It usually is sourced by adding properties
 * to an EntityData, using a backend source such as the world service into this
 * interface and via a `WorldObjectBuilder`.
 */
export declare interface LinkEntityData extends EntityData {
  geometries: Map<string, Map<string, GeometryData>>;
}

/**
 * A link entity is a specifc type of entity that contains physical information.
 * That may include geometry definitions, and potentially collision information.
 */
export class LinkEntity extends Entity {
  // The various geometries defined for this link.
  // E.g., geometry type (e.g., visual) -> (geometry_name -> actual_geometry).
  private readonly geometries: Map<string, Map<string, Geometry>>;

  /**
   * Constructs a link entity.
   */
  constructor(
    params: LinkEntityData,
    // CAUTION: parent object might not be fully initialized -- best to avoid
    // using it in the constructor
    parentObject: WorldObject,
    children: Entity[],
  ) {
    super(params, parentObject, children);
    this.geometries = new Map<string, Map<string, Geometry>>();
    for (const [geometryGroupName, geometryMap] of params.geometries) {
      const newGeometryMap = new Map<string, Geometry>();
      for (const [geometryName, geometryData] of geometryMap) {
        newGeometryMap.set(
          geometryName,
          new Geometry(
            geometryData.refTShape,
            geometryData.shape,
            geometryData.materialOverrides,
          ),
        );
      }
      this.geometries.set(geometryGroupName, newGeometryMap);
    }
  }

  override getEntityType(): 'Link' {
    return 'Link';
  }

  override asConcreteEntity(): LinkEntity {
    return this;
  }

  *getGeometryNames(): Generator<string> {
    yield* this.geometries.keys();
  }

  getGeometriesByName(name: string): Map<string, Geometry> | undefined {
    return this.geometries.get(name);
  }

  getCollisionGeometries(): Geometry[] {
    return [...this.getNamedCollisionGeometries().values()];
  }

  getVisualGeometries(): Geometry[] {
    return [...this.getNamedVisualGeometries().values()];
  }

  getNamedCollisionGeometries(): Map<string, Geometry> {
    return this.getGeometriesByName(COLLISION_GEOMETRY_KEY) ?? new Map();
  }

  getNamedVisualGeometries(): Map<string, Geometry> {
    return this.getGeometriesByName(VISUAL_GEOMETRY_KEY) ?? new Map();
  }

  getGeometries(): Map<string, Map<string, Geometry>> {
    return this.geometries;
  }
}

/**
 * Specifies the near and far clipping plane distances of a camera sensor.
 */
export declare interface Clip {
  near: number;
  far: number;
}

/**
 * Specifies a base camera with optional clipping planes.
 */
declare interface BaseCamera {
  type: 'camera' | 'depth';
  clip?: Clip;
}

/**
 * Specifies a standard intensity camera.
 */
export declare interface Camera extends BaseCamera {
  type: 'camera';
}

/**
 * Specifies a camera also providing depth data.
 */
export declare interface DepthCamera extends BaseCamera {
  type: 'depth';
}

/**
 * Specifies the currently supported camera sensor types.
 */
export type CameraSensor = Camera | DepthCamera;

/**
 * Specifies the currently supported sensor types.
 */
export type Sensor = CameraSensor;

/**
 * Returns true if the sensor is a camera sensor.
 */
export function isCameraSensor(sensor: Sensor): sensor is CameraSensor {
  return sensor.type === 'camera' || sensor.type === 'depth';
}

/**
 * SensorEntityData is additional data needed to specialize an entity as
 * a SensorEntity in the world model. It usually is sourced by adding properties
 * to an EntityData, using a backend source such as the world service into this
 * interface and via a `WorldObjectBuilder`.
 */
export declare interface SensorEntityData extends EntityData {
  sensor?: Sensor;
}

/**
 * A sensor entity is a specifc type of entity that contains sensor data.
 */
export class SensorEntity extends Entity {
  readonly sensor?: Sensor;

  /**
   * Constructs a sensor entity.
   *
   * @param params The sensor data for this entity.
   */
  constructor(
    params: SensorEntityData,
    parentObject: WorldObject,
    children: Entity[],
  ) {
    super(params, parentObject, children);
    this.sensor = params.sensor;
  }

  override getEntityType(): 'Sensor' {
    return 'Sensor';
  }

  override asConcreteEntity(): SensorEntity {
    return this;
  }
}

/**
 * FrameEntityData is additional data needed to specialize an entity as
 * a FrameEntity in the world model. It usually is sourced by adding properties
 * to an EntityData, using a backend source such as the world service into this
 * interface and via a `WorldObjectBuilder`.
 */
export declare interface FrameEntityData extends EntityData {
  isAttachment: boolean;
}

/**
 * A frame entity is an entity that contains no other semantic information than
 * the fact that it's part of the entity tree..
 */
export class FrameEntity extends Entity {
  readonly isAttachmentFrame: boolean;

  constructor(
    params: FrameEntityData,
    parentObject: WorldObject,
    children: Entity[],
  ) {
    super(params, parentObject, children);
    this.isAttachmentFrame = params.isAttachment;
  }

  override getEntityType(): 'Frame' {
    return 'Frame';
  }

  override asConcreteEntity(): FrameEntity {
    return this;
  }

  /**
   * Returns the parent entity of this frame if it is a frame of the same
   * object, or undefined otherwise.
   */
  getParentFrameInObject(): FrameEntity | undefined {
    const parent = this.getParent();
    if (
      parent instanceof FrameEntity &&
      parent.parentObject === this.parentObject
    ) {
      return parent;
    }
    return undefined;
  }

  /** Returns the children of this entity that are frames. */
  *getChildFrames(): Generator<FrameEntity> {
    for (const child of this.getChildren()) {
      if (child instanceof FrameEntity) {
        yield child;
      }
    }
  }

  /**
   * Returns the reference string of this frame,
   * `<object_name><separator><frame_name>`, which uniquely identifies the frame
   * by name.
   */
  getReferenceString(): string {
    return `${this.parentObject.name}${worldIdUtils.REFERENCE_STRING_SEPARATOR}${this.name}`;
  }
}
