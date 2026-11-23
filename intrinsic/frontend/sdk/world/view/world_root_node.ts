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

import type {EffectRef, Injector, Signal} from '@angular/core';
import {explicitEffect} from '@intrinsic-ai/angular/signals';
import {Observable, Subject, merge} from 'rxjs';
import {takeUntil} from 'rxjs/operators';
import * as three from 'three';

import {AssetDepot} from '../loaders/asset_depot';
import * as worldIdUtils from '../model/world_ids';
import {WorldModel} from '../model/world_model';
import * as mathUtils from '../utils/math';
import {isMesh} from '../utils/three_typechecks';
import {ObjectStateManager} from './object_state_manager';
import {RenderLayers} from './outline_renderer';
import {PickerManager} from './picker_manager';
import {SceneUpdater} from './scene_updater';
import * as shapes from './shapes';
import {
  IntrEntityNode,
  IntrKinematicRootEntityNode,
  IntrRootEntityNode,
} from './world_scene_nodes';

/** Dependencies and inputs of a WorldRootNode. */
export interface WorldRootNodeDeps {
  // The world to display.
  worldId: Signal<string>;
  // The WorldModel to display. Expected to be the model for `worldId`.
  worldModel: Signal<WorldModel | undefined>;
  // When true, polls the backend for updates to the WorldModel.
  polling: Signal<boolean>;

  assetDepot: AssetDepot;
  pickerManager: PickerManager;
  // Source of the hidden objects, entity margins and collision geometry
  // visibility to display.
  objectStateManager: ObjectStateManager;
  // Injector used to create the node's effects.
  injector: Injector;
}

/**
 * WorldSceneRoot represents the root node of the scene graph for a specific
 * world. It provides a simple way to get object/frame/entity nodes by id, as
 * well as controls for visibility, and reacts to picker state changes.
 *
 * The owner is responsible for calling `teardown` once the node is no longer
 * needed.
 */
export class WorldRootNode extends IntrRootEntityNode {
  private readonly worldId: Signal<string>;
  private readonly worldModel: Signal<WorldModel | undefined>;
  private readonly polling: Signal<boolean>;
  private readonly sceneUpdater: SceneUpdater;

  // Emits (and completes) on teardown.
  private readonly destroyed = new Subject<void>();
  // The effects created by this node. Destroyed on teardown.
  private readonly effects: EffectRef[] = [];

  // Observable that notifies subscribers when any nodes in the scene are
  // updated in-place (such as changing poses, visibility, etc).
  protected readonly updatedSubject = new Subject<void>();
  readonly updated: Observable<void> = this.updatedSubject;

  // Observable that notifies subscribers when nodes are add/removed from the
  // scene graph.
  readonly mutated: Observable<void>;

  protected readonly assetDepot: AssetDepot;
  protected readonly pickerManager: PickerManager;
  protected readonly objectStateManager: ObjectStateManager;

  // Map of objects based on the object name.
  private deepPickedObjectId = WorldModel.NULL_OBJECT_ID;
  private entityCollisions: number[] = [];
  entitiesShowingMargin: ReadonlySet<number> = new Set<number>();
  collisionGeometryVisible = false;
  private readonly marginShapes = new Map<three.Mesh, shapes.MarginShape>();
  entityCollisionMarginValue = 0;
  hiddenObjects: ReadonlySet<string> = new Set<string>();
  private pickedObjectIds: ReadonlySet<number> = new Set<number>();
  private pickedEntityId = 0;

  constructor(deps: WorldRootNodeDeps) {
    super(
      worldIdUtils.ROOT_ENTITY_ID,
      WorldModel.ROOT_OBJECT_ID,
      'Root',
      'root',
    );
    this.worldId = deps.worldId;
    this.worldModel = deps.worldModel;
    this.polling = deps.polling;
    this.assetDepot = deps.assetDepot;
    this.pickerManager = deps.pickerManager;
    this.objectStateManager = deps.objectStateManager;
    const effectOptions = {injector: deps.injector};

    // Create the SceneUpdater, which handles updating the IntrEntity tree rooted
    // at this node to match the WorldModel.
    this.sceneUpdater = new SceneUpdater(this, this.model, this.assetDepot);
    this.mutated = this.sceneUpdater.mutated;

    // We forward the SceneUpdater events to our own events (which are more general).
    this.sceneUpdater.updated.pipe(takeUntil(this.destroyed)).subscribe(() => {
      this.updatedSubject.next();
    });

    const modelEffect = explicitEffect(
      [this.worldModel],
      ([model], onCleanup) => {
        if (!model) return;
        this.sceneUpdater.setModel(model);
        const updatedSub = model.updated
          .pipe(takeUntil(this.destroyed))
          .subscribe(() => {
            this.updateSceneFromModel();
          });
        const updatedStateSub = model.updatedState
          .pipe(takeUntil(this.destroyed))
          .subscribe(() => {
            this.updateSceneFromModel(true);
          });
        this.updateSceneFromModel();
        onCleanup(() => {
          updatedSub.unsubscribe();
          updatedStateSub.unsubscribe();
        });
      },
      effectOptions,
    );
    this.effects.push(modelEffect);

    const pollingEffect = explicitEffect(
      [this.worldModel, this.polling],
      ([model, polling], onCleanup) => {
        if (!model || !polling) return;
        const pollingSub = model
          .poll()
          .pipe(takeUntil(this.destroyed))
          .subscribe();
        onCleanup(() => {
          pollingSub.unsubscribe();
        });
      },
      effectOptions,
    );
    this.effects.push(pollingEffect);

    // Collision geometry visibility is stored per world, so this needs to
    // react to both the visibility and the current worldId.
    const collisionGeometryVisibilityEffect = explicitEffect(
      [this.objectStateManager.collisionGeometryVisibility, this.worldId],
      ([visibilityMap, worldId]) => {
        const visible = !!visibilityMap.get(worldId);
        if (visible) {
          void this.checkCollisions(worldId);
        }
        this.setShowCollisionGeometry(visible);
        for (const rootEntity of this.traverseObjectRoots()) {
          for (const entity of rootEntity.traverseOwnEntities()) {
            for (const node of entity.collisionNode.children) {
              node.visible = visible;
            }
          }
        }
        this.updatedSubject.next();
      },
      effectOptions,
    );
    this.effects.push(collisionGeometryVisibilityEffect);

    // Nodes are recreated when the scene graph is rebuilt (e.g. when the model
    // changes or another world is shown), which drops their colliding, hidden
    // and margin state. Re-apply the current state afterwards.
    this.mutated.pipe(takeUntil(this.destroyed)).subscribe(() => {
      this.reapplyHiddenObjectsAndMargins();
      this.updateSceneEntityCollisions();
    });

    // If the scene is loaded, set the picked object again and
    // highlight the object node in the scene if it has changed.
    this.mutated.pipe(takeUntil(this.destroyed)).subscribe(() => {
      for (const oid of this.pickedObjectIds) {
        if (this.model?.hasObjectOrFrame(oid)) {
          this.getNodeById(oid)?.enableHighlight(true);
        } else {
          this.pickerManager.removePickedObjectId(oid);
        }
      }
      if (this.pickedEntityId !== 0) {
        if (this.model?.hasObjectOrFrame(this.pickedEntityId)) {
          this.getNodeById(this.pickedEntityId)?.enableHighlight(true);
        }
      }
    });

    // Entity margin visibility (along with matching collision geo)
    const entityMarginEffect = explicitEffect(
      [this.objectStateManager.entitiesShowingMargin],
      ([entityMarginVisible]) => {
        let sceneChanged = false;
        // Go through all entities showing their margin and if they are not in
        // our visible set yet, show them.
        for (const id of entityMarginVisible) {
          const entityNode = this.getNodeByEntityId(id);
          if (entityNode && !this.entitiesShowingMargin.has(id)) {
            this.showMarginShape(entityNode);
            sceneChanged = true;
          }
        }
        // Vice versa, if they are in our local set, but not in the manager's,
        // they need to become invisible.
        for (const id of this.entitiesShowingMargin) {
          const entityNode = this.getNodeByEntityId(id);
          if (entityNode && !entityMarginVisible.has(id)) {
            this.hideMarginShape(entityNode);
            sceneChanged = true;
          }
        }
        this.entitiesShowingMargin = entityMarginVisible;
        if (sceneChanged) {
          this.updatedSubject.next();
        }
      },
      effectOptions,
    );
    this.effects.push(entityMarginEffect);
    // Object and frame visibility
    const hiddenObjectsEffect = explicitEffect(
      [this.objectStateManager.hiddenObjectsOrFrames],
      ([hiddenObjectsOrFrames]) => {
        let sceneChanged = false;
        // Go through all hidden objects and hide them.
        for (const ref of hiddenObjectsOrFrames) {
          const node = this.getNodeByObjectOrFrameReference(ref);
          if (node) {
            node.makeInvisible(true);
            sceneChanged = true;
          }
        }
        // Vice versa, if they are in our local set, but not hidden anymore,
        // they need to become visible.
        for (const ref of this.hiddenObjects) {
          const node = this.getNodeByObjectOrFrameReference(ref);
          if (node && !hiddenObjectsOrFrames.has(ref)) {
            node.makeInvisible(false);
            sceneChanged = true;
          }
        }
        this.hiddenObjects = hiddenObjectsOrFrames;
        if (sceneChanged) {
          this.updatedSubject.next();
        }
      },
      effectOptions,
    );
    this.effects.push(hiddenObjectsEffect);

    // Object picking
    this.pickerManager.newPick
      .pipe(takeUntil(this.destroyed))
      .subscribe((picked: ReadonlySet<number>) => {
        this.pickedObjectIds = picked;
        this.outlineObjects([...picked.values()]);
      });

    this.pickerManager.newEntityPick
      .pipe(takeUntil(this.destroyed))
      .subscribe((pick: number) => {
        this.pickedEntityId = pick;
      });

    this.pickerManager.newHoverPick
      .pipe(takeUntil(this.destroyed))
      .subscribe((oid: number) => {
        for (const node of this.traverseObjectRoots()) {
          node.enableHighlight(false);
        }
        for (const node of this.traverseFrames()) {
          node.enableHighlight(false);
        }
        this.getNodeById(oid)?.enableHighlight(true);
      });

    // Deep object picking
    this.pickerManager.deepPick
      .pipe(takeUntil(this.destroyed))
      .subscribe((objectId) => {
        this.deepPickedObjectId = objectId;
        // Remove the outline effect if it exists.
        this.outlineObjects([]);
      });

    // Entity picking
    merge(
      this.pickerManager.newEntityPick,
      this.pickerManager.newHoverEntityPick,
    )
      .pipe(takeUntil(this.destroyed))
      .subscribe((eid: number) => {
        const objNode = this.getNodeByObjectId(this.deepPickedObjectId);
        objNode?.enableHighlight(false);
        const entityNode = this.getNodeByEntityId(eid);
        // Since other entities are visible, makes sure we only highlight
        // nodes belonging to the object
        if (entityNode?.parentObjectId === this.deepPickedObjectId) {
          // Note - when an object is deep-picked we only want to
          // highlight the specific picked entity (not its whole subtree).
          entityNode?.enableHighlight(true, false);
        }
      });
  }

  // Map of objects based on the entity id.
  protected get entityNodeMap(): ReadonlyMap<number, IntrEntityNode> {
    return this.sceneUpdater.getEntityNodeMap();
  }
  // Map of objects based on the object name.
  protected get objectNodesByName(): ReadonlyMap<string, IntrRootEntityNode> {
    return this.sceneUpdater.getObjectNodesByName();
  }
  // Map of entity nodes based on the generated TF name for the entity.
  protected get entityNodesByTFName(): ReadonlyMap<string, IntrEntityNode> {
    return this.sceneUpdater.getEntityNodesByTFName();
  }

  setShowCollisionGeometry(show: boolean): void {
    this.sceneUpdater.setShowCollisionGeometry(show);
  }

  /** Returns the id of the world currently displayed by this node. */
  getWorldId(): string {
    return this.worldId();
  }

  /**
   * Destroys the node's effects, stops all of its subscriptions and detaches
   * it from its parent. The node should not be used afterwards.
   */
  teardown(): void {
    for (const effectRef of this.effects) {
      effectRef.destroy();
    }
    this.destroyed.next();
    this.destroyed.complete();
    this.removeFromParent();
  }

  // The WorldModel currently displayed by this node.
  protected get model(): WorldModel | undefined {
    return this.worldModel();
  }

  getModel(): WorldModel | undefined {
    return this.worldModel();
  }

  private updateSceneFromModel(posesOnly = false) {
    this.sceneUpdater.updateSceneFromModel(posesOnly);
  }

  *getHelpers(): Generator<shapes.SceneHelper> {
    for (const root of this.traverseObjectRoots()) {
      for (const child of root.children) {
        if (child instanceof shapes.SceneHelper) {
          yield child;
        }
      }
    }
  }

  // Return the scene node given an object/frame id.
  getNodeById(rid: number): IntrEntityNode | undefined {
    return this.getNodeByObjectId(rid) ?? this.getNodeByEntityId(rid);
  }

  getNodeByObjectOrFrameReference(ref: string): IntrEntityNode | undefined {
    const objectOrFrame = this.model?.getObjectOrFrameByReference(ref);
    if (!objectOrFrame) return;
    return this.getNodeById(objectOrFrame.id);
  }

  getKinematicNodeById(rid: number): IntrKinematicRootEntityNode | undefined {
    const node = this.getNodeById(rid);
    if (node instanceof IntrKinematicRootEntityNode) {
      return node;
    }
    return undefined;
  }

  // Return the node given an entity's id.
  getNodeByEntityId(eid: number): IntrEntityNode | undefined {
    return this.entityNodeMap.get(eid);
  }

  getNodeByTFName(name: string): IntrEntityNode | undefined {
    return this.entityNodesByTFName.get(name);
  }

  getNodeByObjectId(oid: number): IntrRootEntityNode | undefined {
    const eid = this.model?.getObject(oid)?.getRootEntity().id;
    if (eid) {
      const node = this.entityNodeMap.get(eid);
      if (node instanceof IntrRootEntityNode) {
        return node;
      }
    }
    return undefined;
  }

  getNodeByObjectName(name: string): IntrRootEntityNode | undefined {
    return this.objectNodesByName.get(name);
  }

  getTFNames(): string[] {
    return Array.from(this.entityNodesByTFName.keys());
  }

  getObjectNames(): string[] {
    return Array.from(this.objectNodesByName.keys());
  }

  async sync() {
    if (!this.model) return;
    await this.model.sync();
  }

  /** Checks the world for collisions and marks the colliding entities. */
  private async checkCollisions(worldId: string) {
    const model = this.model;
    if (!model) return;
    // This is called without awaiting it, so don't let errors escape.
    const collisionPairs = await model
      .checkCollisions()
      .catch((error: unknown) => {
        console.error('Failed to check collisions:', error);
        return undefined;
      });
    if (!collisionPairs) return;
    // Ignore results for a world that is no longer shown.
    if (this.model !== model || this.worldId() !== worldId) return;
    const collidingEntityIds = new Set<number>();
    for (const {leftEntityId, rightEntityId} of collisionPairs) {
      collidingEntityIds.add(leftEntityId);
      collidingEntityIds.add(rightEntityId);
    }
    this.entityCollisions = [...collidingEntityIds];
    this.updateSceneEntityCollisions();
  }

  private updateSceneEntityCollisions() {
    this.traverse((node) => {
      if (node instanceof shapes.CollisionShape) {
        node.markAsColliding(false);
      }
    });

    for (const eid of this.entityCollisions) {
      const group = this.getNodeByEntityId(eid)?.collisionNode;
      group?.traverse((node) => {
        if (node instanceof shapes.CollisionShape) {
          node.markAsColliding(true);
        }
      });
    }
    this.updatedSubject.next();
  }

  private outlineObjects(oids: number[]) {
    this.traverse((e) => {
      e.layers.set(RenderLayers.NO_OUTLINE);
    });
    for (const oid of oids) {
      const node = this.getNodeByObjectId(oid);
      if (!node) {
        continue;
      }
      for (const entity of node.traverseOwnEntities()) {
        for (const renderable of entity.getRenderables()) {
          renderable.layers.set(RenderLayers.OUTLINE);
        }
      }
    }
  }

  /**
   * Reapplies the object visibility and entity margins to the nodes in the
   * scene.
   */
  private reapplyHiddenObjectsAndMargins() {
    for (const ref of this.hiddenObjects) {
      this.getNodeByObjectOrFrameReference(ref)?.makeInvisible(true);
    }
    for (const id of this.entitiesShowingMargin) {
      const entityNode = this.getNodeByEntityId(id);
      if (entityNode) {
        this.showMarginShape(entityNode);
      }
    }
  }

  private showMarginShape(entityNode: IntrEntityNode) {
    entityNode?.collisionNode.traverse((node) => {
      node.visible = true;
      if (!node.parent) return;
      if (isMesh(node)) {
        let marginShape = this.marginShapes.get(node);
        if (!marginShape) {
          marginShape = new shapes.MarginShape(node);
          this.marginShapes.set(node, marginShape);
        }
        marginShape.removeFromParent();
        marginShape.setMargin(this.entityCollisionMarginValue);
        node.parent.add(marginShape);
        mathUtils.updateObjectLocalTransform(marginShape, node.matrix, true);
      }
    });
  }

  private hideMarginShape(entityNode: IntrEntityNode) {
    for (const shape of this.marginShapes.values()) {
      shape.removeFromParent();
    }
    for (const node of entityNode.collisionNode.children) {
      node.visible = this.collisionGeometryVisible;
    }
  }

  clearMarginShapes() {
    for (const shape of this.marginShapes.values()) {
      shape.removeFromParent();
      shape.geometry.dispose();
    }
    this.marginShapes.clear();
    this.entityCollisionMarginValue = 0;
  }

  updateMarginShapes(marginValue: number) {
    this.entityCollisionMarginValue = marginValue;
    for (const shape of this.marginShapes.values()) {
      shape.setMargin(marginValue);
    }
  }
}
