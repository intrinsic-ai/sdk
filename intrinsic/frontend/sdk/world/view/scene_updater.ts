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

import {assumeExhaustive} from '@intrinsic-ai/core/typescript';
import {Observable, Subject} from 'rxjs';
import * as three from 'three';

import {AssetDepot} from '../loaders/asset_depot';
import {
  Entity,
  FrameEntity,
  JointEntity,
  LinkEntity,
  SensorEntity,
} from '../model/entity';
import {Geometry} from '../model/geometry';
import {ROOT_ENTITY_ID} from '../model/world_ids';
import {WorldModel} from '../model/world_model';
import * as mathUtils from '../utils/math';
import * as shapes from './shapes';
import {applyMaterialToNode} from './shapes';
import {
  IntrEntityNode,
  IntrFrameNode,
  IntrGeometryGroup,
  IntrJointNode,
  IntrKinematicRootEntityNode,
  IntrLinkNode,
  IntrRootEntityNode,
  IntrSensorNode,
} from './world_scene_nodes';

/**
 * SceneUpdater synchronizes a three.js entity tree rooted at an IntrRootEntityNode
 * with a WorldModel.
 */
export class SceneUpdater {
  // Map of objects based on the entity id.
  private readonly entityNodeMap = new Map<number, IntrEntityNode>();
  // Map of objects based on the object name.
  private readonly objectNodesByName = new Map<string, IntrRootEntityNode>();
  // Map of entity nodes based on the generated TF name for the entity.
  private readonly entityNodesByTFName = new Map<string, IntrEntityNode>();

  private readonly updatedSubject = new Subject<void>();
  readonly updated: Observable<void> = this.updatedSubject;

  // Observable that notifies subscribers when nodes are add/removed from the
  // scene graph.
  private readonly mutatedSubject = new Subject<void>();
  readonly mutated: Observable<void> = this.mutatedSubject;

  private showCollisionGeometry = false;

  constructor(
    private readonly rootNode: IntrRootEntityNode,
    private model: WorldModel | undefined,
    private readonly assetDepot: AssetDepot,
  ) {
    this.resetMapsToRootNode();
  }

  /**
   * Clears the lookup maps and (re-)seeds them with the root node, which is
   * always part of the scene graph regardless of the model's contents.
   */
  private resetMapsToRootNode(): void {
    this.entityNodeMap.clear();
    this.objectNodesByName.clear();
    this.entityNodesByTFName.clear();

    this.entityNodeMap.set(ROOT_ENTITY_ID, this.rootNode);
    this.objectNodesByName.set('root', this.rootNode);
    this.entityNodesByTFName.set('Root', this.rootNode);
  }

  setShowCollisionGeometry(show: boolean): void {
    this.showCollisionGeometry = show;
  }

  getRootNode(): IntrRootEntityNode {
    return this.rootNode;
  }

  setModel(model?: WorldModel): void {
    this.model = model;
  }

  getModel(): WorldModel | undefined {
    return this.model;
  }

  updateSceneFromModel(posesOnly = false) {
    // Update the skeleton of the world scene.
    if (posesOnly) {
      this.updatePoses();
    } else {
      this.updateSkeleton();
      this.updateShapes();
      this.updatePoses();
      this.mutatedSubject.next();
    }
  }

  private entityIdToTfFrameName(entity: Entity): string {
    const resourceName = entity.getResourceName();
    const localName =
      entity.name.length > 0 ? entity.name : entity.id.toString();
    const objectName = entity.parentObject.name;

    if (resourceName) {
      return resourceName + '/' + objectName + '/' + localName;
    } else {
      return objectName + '/' + localName;
    }
  }

  private updatePoses() {
    if (!this.model?.isLoaded()) return;
    let emitUpdate = false;
    for (const entityNode of this.entityNodeMap.values()) {
      if (entityNode instanceof IntrJointNode) {
        const jointEntity = this.model
          .getObject(entityNode.parentObjectId)
          ?.asKinematic()
          ?.getEntityById(entityNode.entityId)
          ?.asJoint();
        if (!jointEntity) {
          continue;
        }
        if (jointEntity.value !== entityNode.getJointValue()) {
          emitUpdate = true;
          entityNode.updateJointValue(jointEntity.value);
        }
      } else {
        const matrix = this.model.getLocalTrans(entityNode.entityId);
        if (
          matrix &&
          !mathUtils.matricesAlmostEqual(matrix, entityNode.matrix)
        ) {
          emitUpdate = true;
          mathUtils.updateObjectLocalTransform(entityNode, matrix, false);
        }
      }
    }
    // Synchronously recompute world matrices across the scene graph before
    // notifying updatedSubject / mutatedSubject subscribers, ensuring
    // downstream consumers observe up-to-date world transforms.
    this.rootNode.updateMatrixWorld(true);
    if (emitUpdate) {
      this.updatedSubject.next();
    }
  }

  // Construct the skeleton of scene graph, without attaching the shape nodes.
  private updateSkeleton() {
    if (!this.model || !this.model.isLoaded()) {
      return;
    }

    // Sanity-check that the model's root entity is a link. The root entity
    // is represented visually by rootNode.
    const rootEntity = this.model.getEntity(ROOT_ENTITY_ID);
    if (rootEntity && !rootEntity.isLink()) {
      throw new Error(
        `The scene's root entity (${rootEntity.name}) is not a link. ` +
          `It has type ${rootEntity.getEntityType()}. This is not supported. ` +
          `A scene's root entity must be a link.`,
      );
    }

    // We first remove all existing nodes and clear our internal maps. However, we keep
    // a copy of the existing entityNodeMap so that we can reuse existing nodes
    // where possible.
    const previousEntityNodes = new Map(this.entityNodeMap);
    for (const [eid, node] of previousEntityNodes.entries()) {
      if (eid !== ROOT_ENTITY_ID) {
        node.removeFromParent();
      }
    }
    // Clear the maps, then re-seed them with the root node.
    this.resetMapsToRootNode();

    const entityIds = [];
    for (const oid of this.model.getObjectIds()) {
      entityIds.push(...(this.model.getObjectEntityIds(oid) ?? []));
    }

    // For each entity, we either update an existing node in-place
    // or create a fresh node to represent the entity.
    for (const eid of entityIds) {
      const entity = this.model.getEntity(eid);
      if (!entity) continue;
      if (eid === ROOT_ENTITY_ID) {
        // We've already handled the root entity
        continue;
      }
      let node = previousEntityNodes.get(eid);
      if (!node || !this.tryUpdateNodeFromEntity(node, entity)) {
        node = this.createNodeForEntity(entity);
      }
      if (node instanceof IntrRootEntityNode) {
        this.objectNodesByName.set(entity.parentObject.name, node);
      }
      this.entityNodeMap.set(eid, node);

      // Map the generated TF name for the entity to the node.
      const tfName = this.entityIdToTfFrameName(entity);
      this.entityNodesByTFName.set(tfName, node);
    }

    // Create the parent/child graph of obj nodes.
    for (const eid of entityIds) {
      const pid = this.model.getEntity(eid)?.getParent()?.id;
      if (!pid) continue;
      const childNode = this.entityNodeMap.get(eid);
      if (!childNode) continue;
      // Add the node as children to parent node.
      this.entityNodeMap.get(pid)?.add(childNode);
    }
  }

  private createNodeForEntity(entity: Entity): IntrEntityNode {
    const eid = entity.id;
    const obj = entity.parentObject;
    const rootEntityId = obj.getRootEntity().id;
    const isRootEntity = rootEntityId === eid;

    // First, we sanity-check that, if this is the root entity, it
    // must also be a link. Only links can be root entities.
    if (isRootEntity && !entity.isLink()) {
      throw new Error(
        `Entity ${entity.name} is the root entity of object ${entity.parentObject.name} ` +
          `but it is not a link. It has type ${entity.getEntityType()}. ` +
          `This is not supported. An object's root entity must be a link.`,
      );
    }

    // Create a threejs node corresponding to the Entity type
    const concrete = entity.asConcreteEntity();
    if (concrete.isLink()) {
      if (isRootEntity) {
        if (obj.isKinematic()) {
          return new IntrKinematicRootEntityNode(
            eid,
            obj.id,
            concrete.name,
            obj.name,
            obj.isRobot(), // Kinematic objects with at least one flange frame count as robots and can be controlled by ICON
          );
        }
        return new IntrRootEntityNode(eid, obj.id, concrete.name, obj.name);
      }
      return new IntrLinkNode(eid, obj.id, concrete.name);
    } else if (concrete.isFrame()) {
      return new IntrFrameNode(eid, obj.id, concrete.name);
    } else if (concrete.isJoint()) {
      return new IntrJointNode(eid, obj.id, concrete.name, concrete);
    } else if (concrete.isSensor()) {
      return new IntrSensorNode(eid, obj.id, concrete.name);
    } else {
      assumeExhaustive(concrete);
      throw new Error('Unhandled entity');
    }
  }

  /**
   * Attempts to update the node in place from the entity, if the node corresponds to the type
   * of the entity. For example, an IntrLinkNode can be updated in place from a LinkEntity, but
   * not from a JointEntity.
   *
   * If the node cannot be updated in place, returns false. This indicates that the caller should
   * create a fresh node to represent the entity.
   */
  private tryUpdateNodeFromEntity(
    node: IntrEntityNode,
    entity: Entity,
  ): boolean {
    if (node instanceof IntrLinkNode && entity instanceof LinkEntity) {
      const expectedType = entity.isObjectRoot()
        ? entity.parentObject.isKinematic()
          ? 'kinematic_root_link'
          : 'root_link'
        : 'link';

      if (node.getNodeType() === expectedType) {
        node.updateFromLinkEntity(entity);
        return true;
      }
      return false;
    } else if (node instanceof IntrFrameNode && entity instanceof FrameEntity) {
      node.updateFromFrameEntity(entity);
      return true;
    } else if (node instanceof IntrJointNode && entity instanceof JointEntity) {
      node.updateFromJointEntity(entity);
      return true;
    } else if (
      node instanceof IntrSensorNode &&
      entity instanceof SensorEntity
    ) {
      node.updateFromSensorEntity(entity);
      return true;
    }
    return false;
  }

  private async getShapeForGeometry(
    geo: Geometry,
    parent: IntrGeometryGroup,
    i: number,
    geometryName: string,
  ): Promise<three.Object3D> {
    const mutated =
      !geo.hasRemoteShape() || !this.assetDepot.hasShape(geo.getShape());
    const subject = mutated ? this.mutatedSubject : this.updatedSubject;
    const shape = await this.assetDepot.getShape(geo.getShape());
    const materialOverrides = geo.getMaterialOverrides();
    if (materialOverrides) {
      applyMaterialToNode(shape, materialOverrides);
    }
    parent.setShapeAt(i, shape);
    parent.children[i].geometryName = geometryName;
    mathUtils.updateObjectLocalTransform(shape, geo.getRefTShape(), true);
    subject.next();
    return shape;
  }

  // Reset shape nodes states and meta data.
  private updateShapes() {
    if (!this.model?.isLoaded()) return;

    for (const [eid, entityNode] of this.entityNodeMap.entries()) {
      // Instantiate each visual geometry.
      const visualGeometriesMap = this.model
        .getEntity(eid)
        ?.asLink()
        ?.getNamedVisualGeometries();
      const visualGeometries = Array.from(visualGeometriesMap?.entries() ?? []);
      // Pre-allocate the array as ordering matters; shapes will be inserted at
      // the same index as the associated GeometryComponent.
      entityNode.visualNode.allocate(visualGeometries.length);
      for (const [i, [geometryName, geo]] of visualGeometries.entries()) {
        void this.getShapeForGeometry(
          geo,
          entityNode.visualNode,
          i,
          geometryName,
        );
      }

      // Instantiate each collision geometry.
      const collisionGeometriesMap = this.model
        .getEntity(eid)
        ?.asLink()
        ?.getNamedCollisionGeometries();
      const collisionGeometries = Array.from(
        collisionGeometriesMap?.entries() ?? [],
      );
      // Pre-allocate the array as ordering matters; shapes will be inserted at
      // the same index as the associated GeometryComponent.
      entityNode.collisionNode.allocate(collisionGeometries.length);
      for (const child of entityNode.collisionNode.children) {
        child.visible = this.showCollisionGeometry;
      }
      for (const [i, [geometryName, geo]] of collisionGeometries.entries()) {
        void this.getShapeForGeometry(
          geo,
          entityNode.collisionNode,
          i,
          geometryName,
        ).then((shape) => {
          shapes.CollisionShape.applyToSubtree(shape);
        });
      }
    }
  }

  getEntityNodeMap(): ReadonlyMap<number, IntrEntityNode> {
    return this.entityNodeMap;
  }

  getObjectNodesByName(): ReadonlyMap<string, IntrRootEntityNode> {
    return this.objectNodesByName;
  }

  getEntityNodesByTFName(): ReadonlyMap<string, IntrEntityNode> {
    return this.entityNodesByTFName;
  }
}
