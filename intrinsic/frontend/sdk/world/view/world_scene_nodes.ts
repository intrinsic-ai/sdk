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
  Entity,
  FrameEntity,
  JointEntity,
  LinkEntity,
  MotionType,
  SensorEntity,
} from '../model/entity';
import {isMesh, isPoints} from '../utils/three_typechecks';
import {RenderLayers} from './outline_renderer';
import {CollisionShape, FrameAxes, IntrVisualMesh, SceneHelper} from './shapes';

/** All possible entity node types. */
export const ALL_ENTITY_NODE_TYPES = [
  'link',
  'root_link',
  'kinematic_root_link',
  'frame',
  'joint',
  'sensor',
] as const;

/**
 * A union of the possible entity node type strings. Useful
 * for switching over with `node.getNodeType()` to handle all cases.
 */
export type EntityNodeType = (typeof ALL_ENTITY_NODE_TYPES)[number];

/**
 * A union of the concrete subclasses of IntrEntityNode.
 * The idea is that for each concrete type of Entity (see
 * the type `ConcreteEntity` in entity.ts), we should have
 * a corresponding threejs node class to visually represent it.
 *
 * Given an abstract IntrEntityNode, you can get the concrete subclass
 * with `asConcreteNode()`.
 */
export type ConcreteEntityNode =
  | IntrLinkNode
  | IntrRootEntityNode
  | IntrKinematicRootEntityNode
  | IntrFrameNode
  | IntrJointNode
  | IntrSensorNode;

/**
 * Simple wrapper around Entity geometry nodes
 */
export class IntrGeometryNode extends three.Object3D {
  override parent: IntrGeometryGroup | null = null;

  constructor(
    readonly componentIndex: number,
    public geometryName: string,
  ) {
    super();
  }

  override clone(recursive = true): this {
    return new IntrGeometryNode(this.componentIndex, this.geometryName).copy(
      this,
      recursive,
    ) as this;
  }
}

/**
 * Helper class that facilitates asynchronous loading of geometry in a way that preserves ordering with respect to the world geometry component
 */
export class IntrGeometryGroup extends three.Group {
  override children: IntrGeometryNode[] = [];
  override parent: IntrEntityNode | null = null;
  setShapeAt(index: number, node: three.Object3D) {
    if (index < this.children.length && index >= 0) {
      this.children[index].clear();
      this.children[index].add(node);
    }
  }

  getShapeAt(index: number): IntrGeometryNode {
    return this.children[index];
  }

  allocate(size: number) {
    this.clear();
    for (let i = 0; i < size; i++) {
      const dummy = new three.Object3D();
      dummy.matrixAutoUpdate = false;
      this.add(dummy);
    }
  }
}

// Classes for memoization via type narrowing during addition/copies
class VisualGeometryGroup extends IntrGeometryGroup {
  override clone(): this {
    return new VisualGeometryGroup().copy(this) as this;
  }
}

class CollisionGeometryGroup extends IntrGeometryGroup {
  override clone(): this {
    return new CollisionGeometryGroup().copy(this) as this;
  }
}

/**
 * A custom node type that encapsulates Entity concepts.
 * e.g. collision/visual geometry as well as child entities and world object associations.
 */
export abstract class IntrEntityNode extends three.Group {
  visualNode = new VisualGeometryGroup();
  collisionNode = new CollisionGeometryGroup();
  outlined = false;
  isPhantom = false;

  constructor(
    public entityId: number,
    public parentObjectId: number,
    override name: string,
  ) {
    super();
    this.add(this.visualNode);
    this.add(this.collisionNode);
  }

  abstract asConcreteNode(): ConcreteEntityNode;
  abstract getNodeType(): EntityNodeType;

  /**
   * Update this node in-place from an Entity
   *
   * NOTE: this does not reset the node's geometries. That should be
   * handled separately by loading shapes from the AssetDepot.
   */
  protected updateFromEntity(entity: Entity) {
    this.entityId = entity.id;
    this.parentObjectId = entity.parentObject.id;
    this.name = entity.name;
  }

  override add(node: three.Object3D): this {
    if (node instanceof VisualGeometryGroup) {
      this.visualNode.removeFromParent();
      this.visualNode = node;
    } else if (node instanceof CollisionGeometryGroup) {
      this.collisionNode.removeFromParent();
      this.collisionNode = node;
    }
    return super.add(node);
  }

  protected *traverseSelf(): Generator<IntrEntityNode> {
    yield this;
  }

  showHelpers(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      for (const child of entity.children) {
        if (child instanceof SceneHelper) {
          child.visible = enabled;
        }
      }
    }
  }

  getRenderables(): Array<three.Mesh | three.Points> {
    const renderables: Array<three.Mesh | three.Points> = [];
    for (const node of [this.visualNode, this.collisionNode]) {
      node.traverse((node) => {
        if (isMesh(node) || isPoints(node)) {
          renderables.push(node);
        }
      });
    }
    return renderables;
  }

  override copy(from: IntrEntityNode, recursive = true): this {
    super.copy(from as this, recursive);
    for (const child of this.children) {
      if (child instanceof VisualGeometryGroup) {
        this.visualNode = child;
      }
      if (child instanceof CollisionGeometryGroup) {
        this.collisionNode = child;
      }
    }
    return this;
  }

  override clone(recursive = true): this {
    // A variant of the threejs version that accepts our custom arguments.
    // Threejs calls `new this.constructor()`, but TypeScript compiler apprears
    // to wrongfully claim this.constructor is non-constructable. So we instead
    // grab the constructor through other means.
    // See: (https://github.com/microsoft/TypeScript/issues/3841)
    const ctor = Object.getPrototypeOf(this).constructor;
    const clone = new ctor(this.entityId, this.parentObjectId).copy(
      this,
      recursive,
    );
    return clone as this;
  }

  getRobotRoot(): IntrKinematicRootEntityNode | undefined {
    return [this, ...this.traverseParentEntities()].find(
      (entity): entity is IntrKinematicRootEntityNode =>
        entity instanceof IntrKinematicRootEntityNode &&
        entity.canBeControlledByIcon,
    );
  }

  getRoot(): IntrRootEntityNode | undefined {
    if (this instanceof IntrRootEntityNode) {
      return this;
    }
    for (const entity of this.traverseParentEntities()) {
      if (entity instanceof IntrRootEntityNode) {
        return entity;
      }
    }
    return undefined;
  }

  *traverseParentEntities(): Generator<IntrEntityNode> {
    let parent = this.parent;
    while (parent) {
      if (parent instanceof IntrEntityNode) {
        yield parent;
      }
      parent = parent.parent;
    }
  }

  // depth-first traversal of self and the entire subtree including other world objects
  *traverseEntities(): Generator<IntrEntityNode> {
    yield this;
    for (const node of this.children) {
      if (node instanceof IntrEntityNode) {
        yield* node.traverseEntities();
      }
    }
  }

  // depth-first traversal of self and any nodes belonging to the same world object
  protected *traverseObjectEntities(): Generator<IntrEntityNode> {
    const id = this.parentObjectId;
    yield this;
    for (const node of this.children) {
      if (node instanceof IntrEntityNode && node.parentObjectId === id) {
        yield* node.traverseObjectEntities();
      }
    }
  }

  // Our threejs typeings wrongly insist that the type of `clone` is this, so we
  // conform for forward compatibility.
  createPhantom(): this {
    const clone = this.clone();
    for (const node of clone.traverseEntities()) {
      node.isPhantom = true;
      node.enableHighlight(true);
      node.showHelpers(false);
    }
    return clone;
  }

  getVisualBounds(recursive = true) {
    const box = new three.Box3();
    if (recursive) {
      for (const entity of this.traverseEntities()) {
        box.expandByObject(entity.visualNode);
      }
    } else {
      box.expandByObject(this.visualNode);
    }
    return box;
  }

  getCollisionBounds(recursive = true) {
    const box = new three.Box3();
    if (recursive) {
      for (const entity of this.traverseEntities()) {
        box.expandByObject(entity.collisionNode);
      }
    } else {
      box.expandByObject(this.collisionNode);
    }
    return box;
  }

  enableHighlight(enabled: boolean, recursive = true) {
    const entities: Iterable<IntrEntityNode> = recursive
      ? this.traverseSelf()
      : [this];
    for (const entity of entities) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.enableHighlight(enabled);
        }
      });
    }
  }

  enableColorHighlight(enabled: three.Color | false) {
    for (const entity of this.traverseEntities()) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.enableColorHighlight(enabled);
        }
      });
    }
  }

  enableColorHighlightOfEntities(
    enabled: three.Color | false,
    entityIds: number[] | number,
  ) {
    const ids = Array.isArray(entityIds) ? entityIds : [entityIds];
    for (const entity of this.traverseEntities()) {
      if (ids.includes(entity.entityId)) {
        entity.visualNode.traverse((child) => {
          if (child instanceof IntrVisualMesh) {
            child.enableColorHighlight(enabled);
          }
        });
      }
    }
  }

  enableOutline(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      for (const renderable of entity.getRenderables()) {
        renderable.layers.set(
          enabled ? RenderLayers.OUTLINE : RenderLayers.NO_OUTLINE,
        );
      }
    }
  }

  enableDimming(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.material.enableDimming(enabled);
        }
      });
    }
  }

  enableDiffHighlight(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.enableDiffHighlight(enabled);
        }
      });
    }
  }

  enableHatching(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.material.enableHatching(enabled);
        }
      });
    }
  }

  // Shader-level visibility control. Unlike the threejs 'visible' flag, this
  // property allows is to control the visibility of parent and children nodes
  // independently.
  makeInvisible(enabled: boolean) {
    for (const entity of this.traverseSelf()) {
      entity.visualNode.traverse((child) => {
        if (child instanceof IntrVisualMesh) {
          child.makeInvisible(enabled);
        }
      });
      entity.collisionNode.traverse((child) => {
        if (child instanceof CollisionShape) {
          child.makeInvisible(enabled);
        }
      });
    }
  }

  isInvisible(): boolean {
    let invisible = false;
    this.visualNode.traverse((child) => {
      if (child instanceof IntrVisualMesh) {
        invisible ||= child.isInvisible();
      }
    });
    this.collisionNode.traverse((child) => {
      if (child instanceof CollisionShape) {
        invisible ||= child.isInvisible();
      }
    });
    return invisible;
  }

  setGeometryVisible(visible: boolean) {
    for (const entity of this.traverseSelf()) {
      entity.visualNode.visible = visible;
      entity.collisionNode.visible = visible;
    }
  }

  getGeometryVisible(): boolean {
    return this.visualNode.visible || this.collisionNode.visible;
  }

  toggleGeometryVisible() {
    const visible = this.getGeometryVisible();
    this.setGeometryVisible(!visible);
  }

  // function for optional chaining
  setVisible(visible: boolean) {
    this.visible = visible;
  }

  /**
   * Finds the closest ancestor IntrEntityNode of the given Object3D.
   * Returns the node itself if it is an IntrEntityNode.
   */
  static findParentEntity(node: three.Object3D): IntrEntityNode | undefined {
    let current: three.Object3D | null = node;
    while (current) {
      if (current instanceof IntrEntityNode) {
        return current;
      }
      current = current.parent;
    }
    return undefined;
  }
}

/**
 * Child class representing a link Entity.
 */
export class IntrLinkNode extends IntrEntityNode {
  override asConcreteNode(): ConcreteEntityNode {
    return this;
  }

  override getNodeType(): EntityNodeType {
    return 'link';
  }

  /**
   * Update this node in-place from a LinkEntity
   *
   * NOTE: this does not reset the node's geometries. That should be
   * handled separately by loading shapes from the AssetDepot.
   */
  updateFromLinkEntity(link: LinkEntity) {
    this.updateFromEntity(link);
  }
}

/** Child class representing a frame entity */
export class IntrFrameNode extends IntrEntityNode {
  helper = FrameAxes.create();

  constructor(entityId: number, parentObjectId: number, name: string) {
    super(entityId, parentObjectId, name);
    this.add(this.helper);
  }

  override asConcreteNode(): IntrFrameNode {
    return this;
  }

  override getNodeType(): 'frame' {
    return 'frame';
  }

  /**
   * Update this node in-place from a FrameEntity
   */
  updateFromFrameEntity(frame: FrameEntity) {
    this.updateFromEntity(frame);
  }

  override enableHighlight(enabled: boolean, recursive = true) {
    this.helper.enableHighlight(enabled);
  }

  override enableDimming(enabled: boolean) {
    this.helper.enableDimming(enabled);
  }

  override makeInvisible(enabled: boolean) {
    this.helper.makeInvisible(enabled);
  }

  override setGeometryVisible(visible: boolean) {
    this.helper.visible = visible;
  }

  override getGeometryVisible(): boolean {
    return this.helper.visible;
  }

  override toggleGeometryVisible() {
    this.helper.visible = !this.helper.visible;
  }

  override getVisualBounds(recursive = true) {
    return new three.Box3(
      new three.Vector3(),
      new three.Vector3(),
    ).applyMatrix4(this.matrixWorld);
  }

  override getRenderables(): Array<three.Mesh | three.Points> {
    const renderables: Array<three.Mesh | three.Points> = [];
    this.helper.traverse((node) => {
      if (isMesh(node)) {
        renderables.push(node);
      }
    });
    return renderables;
  }
}

/** Class representing a joint of a kinematic entity. */
export class IntrJointNode extends IntrEntityNode {
  private jointValue = 0;
  private readonly jointAxis: three.Vector3;
  private readonly parentTInboardQuat = new three.Quaternion();
  private readonly parentTInboardPos = new three.Vector3();
  private readonly outboardTChildPose = new three.Matrix4();

  constructor(
    entityId: number,
    parentObjectId: number,
    override name: string,
    private joint: JointEntity,
  ) {
    super(entityId, parentObjectId, name);
    this.jointAxis = this.joint.getAxis();
    this.joint
      .getParentTInboard()
      .decompose(
        this.parentTInboardPos,
        this.parentTInboardQuat,
        new three.Vector3(),
      );
    this.outboardTChildPose = this.joint.getOutboardTChild();
    this.updateJointValue(this.joint.value);
  }

  override asConcreteNode(): IntrJointNode {
    return this;
  }

  override getNodeType(): 'joint' {
    return 'joint';
  }

  /**
   * Update this node in-place from a JointEntity
   */
  updateFromJointEntity(joint: JointEntity) {
    this.updateFromEntity(joint);
    this.joint = joint;
    this.jointAxis.copy(joint.getAxis());
    this.joint
      .getParentTInboard()
      .decompose(
        this.parentTInboardPos,
        this.parentTInboardQuat,
        new three.Vector3(),
      );
    this.outboardTChildPose.copy(joint.getOutboardTChild());
    this.updateJointValue(this.joint.value);
  }

  override clone(recursive = true): this {
    return new IntrJointNode(
      this.entityId,
      this.parentObjectId,
      this.name,
      this.joint,
    ).copy(this, recursive) as this;
  }

  updateJointValue(value: number) {
    this.jointValue = value;
    this.quaternion.copy(this.parentTInboardQuat);
    this.position.copy(this.parentTInboardPos);
    const motionType = this.joint.motionType;
    if (motionType === MotionType.REVOLUTE) {
      this.quaternion.multiply(
        new three.Quaternion().setFromAxisAngle(
          this.jointAxis,
          this.jointValue,
        ),
      );
    } else if (motionType === MotionType.PRISMATIC) {
      // In the kinematics component, the joint axis is specified in the joint's space (aka its
      // inboard space). We convert it from the joint space to the joint parent's space here, since
      // we are updating the joint's position with respect to its parent.
      const jointAxisInParentSpace = this.jointAxis
        .clone()
        .applyQuaternion(this.quaternion);
      this.position.add(jointAxisInParentSpace.multiplyScalar(this.jointValue));
    }

    // Apply the outboardTChildPose to the joint pose similar to how it's done
    // in the backend.
    // http://intrinsic/world/world.cc;l=3367-3370;rcl=715922795
    // Ref: go/intrinsic-choreographer-kinematic-tdd
    const updatedPose = new three.Matrix4()
      .compose(
        this.position,
        this.quaternion,
        new three.Vector3(1, 1, 1) /* unit scale */,
      )
      .multiply(this.outboardTChildPose);
    updatedPose.decompose(this.position, this.quaternion, new three.Vector3());

    this.updateMatrix();
  }

  getJointValue(): number {
    return this.jointValue;
  }

  isFixed(): boolean {
    return this.joint.motionType === MotionType.FIXED;
  }
}

/**
 * Child class representing a sensor Entity. Currently, this is just a
 * direct subclass of IntrEntityNode
 */
export class IntrSensorNode extends IntrEntityNode {
  override asConcreteNode(): IntrSensorNode {
    return this;
  }

  override getNodeType(): 'sensor' {
    return 'sensor';
  }

  /**
   * Update this node in-place from a SensorEntity
   */
  updateFromSensorEntity(sensor: SensorEntity) {
    this.updateFromEntity(sensor);
  }
}

/**
 * Child class representing an Entity that is the root of a WorldObject.
 * Note that there isn't a separate threejs Node for the WorldObject itself.
 * Instead, a WorldObject is represented by its IntrRootEntityNode and manipulated
 * using functions on that node.
 *
 * When working with the node hierarchy, an IntrEntityNode represents the root Entity
 * of its WorldObject if-and-only-if it is an instance of IntrRootEntityNode.
 *
 * This is a subclass of IntrLinkNode because only link nodes may be the root Entity
 * of a WorldObject.
 */
export class IntrRootEntityNode extends IntrLinkNode {
  constructor(
    entityId: number,
    parentObjectId: number,
    name: string,
    public objectName: string,
  ) {
    super(entityId, parentObjectId, name);
  }

  override asConcreteNode(): ConcreteEntityNode {
    return this;
  }

  override getNodeType(): EntityNodeType {
    return 'root_link';
  }

  /**
   * Update this node in-place from a LinkEntity
   *
   * NOTE: this does not reset the node's geometries. That should be
   * handled separately by loading shapes from the AssetDepot.
   */
  override updateFromLinkEntity(link: LinkEntity) {
    super.updateFromLinkEntity(link);
    this.objectName = link.parentObject.name;
  }

  // Traverse only the entities belonging to this object
  *traverseOwnEntities(): Generator<IntrEntityNode> {
    for (const node of this.traverseObjectEntities()) {
      if (this.isPhantom === node.isPhantom) {
        yield node;
      }
    }
  }

  *traverseObjectRoots(): Generator<IntrRootEntityNode> {
    for (const node of this.traverseEntities()) {
      if (
        node instanceof IntrRootEntityNode &&
        this.isPhantom === node.isPhantom
      ) {
        yield node;
      }
    }
  }

  *traverseFrames(): Generator<IntrFrameNode> {
    for (const node of this.traverseEntities()) {
      if (node instanceof IntrFrameNode && this.isPhantom === node.isPhantom) {
        yield node;
      }
    }
  }

  protected override *traverseSelf(): Generator<IntrEntityNode> {
    yield* this.traverseOwnEntities();
  }
}

/** Class representing a kinematic object's root entity. Ths kinematic object will be a robot if canBeControlledByIcon is true */
export class IntrKinematicRootEntityNode extends IntrRootEntityNode {
  constructor(
    entityId: number,
    parentObjectId: number,
    name: string,
    objectName: string,
    public canBeControlledByIcon: boolean,
  ) {
    super(entityId, parentObjectId, name, objectName);
  }
  override asConcreteNode(): ConcreteEntityNode {
    return this;
  }

  override getNodeType(): EntityNodeType {
    return 'kinematic_root_link';
  }

  /**
   * Update this node in-place from a LinkEntity
   *
   * NOTE: this does not reset the node's geometries. That should be
   * handled separately by loading shapes from the AssetDepot.
   */
  override updateFromLinkEntity(link: LinkEntity) {
    super.updateFromLinkEntity(link);
    const obj = link.parentObject;
    if (obj.isKinematic()) {
      this.canBeControlledByIcon = obj.isRobot();
    }
  }

  *traverseJointEntities(): Generator<IntrJointNode> {
    for (const entity of this.traverseOwnEntities()) {
      if (entity instanceof IntrJointNode) {
        yield entity;
      }
    }
  }

  updateJointValues(valuesMap: Map<number, number>) {
    for (const entity of this.traverseOwnEntities()) {
      if (entity instanceof IntrJointNode) {
        const value = valuesMap.get(entity.entityId);
        if (value !== undefined) {
          entity.updateJointValue(value);
        }
      }
    }
    this.updateMatrixWorld(true);
  }

  getJointValuesMap(): Map<number, number> {
    const values = new Map<number, number>();
    for (const entity of this.traverseOwnEntities()) {
      if (entity instanceof IntrJointNode && !entity.isFixed()) {
        values.set(entity.entityId, entity.getJointValue());
      }
    }
    return values;
  }

  jointsEqual(other: IntrKinematicRootEntityNode): boolean {
    const ownValues = this.getJointValuesMap();
    const otherValues = other.getJointValuesMap();
    if (otherValues.size !== ownValues.size) {
      return false;
    }
    for (const [entityId, ownValue] of ownValues.entries()) {
      if (otherValues.get(entityId) !== ownValue) {
        return false;
      }
    }
    return true;
  }
}
