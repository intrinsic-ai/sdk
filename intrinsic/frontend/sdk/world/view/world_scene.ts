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

import type {Injector, Signal} from '@angular/core';
import {signal} from '@angular/core';
import type {Observable} from 'rxjs';
import {Subject} from 'rxjs';
import * as three from 'three';

import type {AssetDepot} from '../loaders/asset_depot';
import type {WorldModel} from '../model/world_model';
import {DebugObjects} from './debug';
import type {ObjectStateManager} from './object_state_manager';
import type {PickerManager} from './picker_manager';
import * as shapes from './shapes';
import {WorldRootNode} from './world_root_node';
import type {
  IntrEntityNode,
  IntrKinematicRootEntityNode,
  IntrRootEntityNode,
} from './world_scene_nodes';

const GROUND_PLANE_SIZE = 10;

/** Dependencies and inputs of a WorldScene. */
export interface WorldSceneDeps {
  // The world to display.
  worldId: Signal<string>;
  // The WorldModel to display. Expected to be the model for `worldId`.
  worldModel: Signal<WorldModel | undefined>;
  // When true, polls the backend for updates to the WorldModel.
  polling: Signal<boolean>;

  assetDepot: AssetDepot;
  pickerManager: PickerManager;
  objectStateManager: ObjectStateManager;
  // Injector used to create effects.
  injector: Injector;
}

/**
 * The WorldScene class is an auxiliary data struct for worldModel to hold alls
 * threeJs objects. Most importantly, it owns a ThreeJs scene that can be set to
 * the WorldView component for rendering. After a WorldModel is loaded, extra
 * geometries needed to be loaded into AssetDepot before a WorldScene can be
 * completely constructed.
 *
 * The owner is responsible for calling `teardown` once the scene is no longer
 * needed.
 */
export class WorldScene extends three.Scene {
  // Observable that notifies subscribers when changes are made that only affect
  // the 3d presentation of the existing scene graph structure e.g.
  // material/visibility changes.
  private readonly updatedSubject = new Subject<void>();
  readonly updated: Observable<void> = this.updatedSubject;
  // Observable that notifies subscribers when nodes are add/removed from the
  // scene graph.
  private readonly mutatedSubject = new Subject<void>();
  readonly mutated: Observable<void> = this.mutatedSubject;

  // The root node of the scene graph for the displayed world. It is owned by
  // this scene: created with it and torn down in `teardown`.
  private readonly rootNode: WorldRootNode;

  get root(): WorldRootNode {
    return this.rootNode;
  }

  // For convenience, we make `root` available as a signal here. The signal
  // contains the current WorldRootNode (or undefined if not yet loaded).
  // It is updated whenever changes are made to the 3d node tree, either structurally
  // (like nodes being added/removed) or pose-only (like nodes being moved).
  readonly rootNodeSignal = signal<WorldRootNode | undefined>(undefined, {
    equal: () => false,
  });

  private static readonly grid = shapes.createGrid(GROUND_PLANE_SIZE, 100);
  readonly grid = WorldScene.grid.clone();

  // Create an invisible ground plane to be used for object placement
  private static readonly groundPlane = new three.Mesh(
    new three.PlaneGeometry(GROUND_PLANE_SIZE, GROUND_PLANE_SIZE),
  );
  readonly groundPlane = WorldScene.groundPlane.clone();

  // Allows adding debug objects to this scene
  private readonly debugObjects = new DebugObjects(this);

  constructor(deps: WorldSceneDeps) {
    super();
    this.rootNode = new WorldRootNode(deps);
    // Set up scene.
    this.add(this.grid);
    this.add(this.groundPlane);
    this.groundPlane.visible = false;

    const root = this.rootNode;
    this.add(root);
    this.applyFramesVisible();
    // The root's subscriptions end when it is torn down in `teardown`.
    root.mutated.subscribe(() => {
      this.applyFramesVisible();
      this.mutatedSubject.next();
      this.rootNodeSignal.set(root);
    });
    root.updated.subscribe(() => {
      this.updatedSubject.next();
      this.rootNodeSignal.set(root);
    });
  }

  /**
   * Tears down the root node and detaches the scene from its parent. The scene
   * should not be used afterwards.
   */
  teardown() {
    // Clear the margin shapes of all roots (including ones added by others,
    // like WorldComparisonOverlay) before tearing down our own root.
    this.clearMarginShapes();
    this.rootNode.teardown();
    this.removeFromParent();
  }

  getVisualShapeNodes(): three.Object3D[] {
    const nodes: three.Object3D[] = [];
    for (const entity of this.rootNode.traverseEntities()) {
      nodes.push(entity.visualNode);
    }
    return nodes;
  }

  getObjectNames(): string[] {
    return this.rootNode.getObjectNames();
  }

  getTFNames(): string[] {
    return this.rootNode.getTFNames();
  }

  // Return the scene node given an object/frame id.
  getNodeById(rid: number): IntrEntityNode | undefined {
    return this.rootNode.getNodeById(rid);
  }

  getRobotNodeById(oid: number): IntrKinematicRootEntityNode | undefined {
    const kinematicNode = this.rootNode.getKinematicNodeById(oid);
    return kinematicNode?.canBeControlledByIcon ? kinematicNode : undefined;
  }

  // Return the node given an entity's id.
  getNodeByEntityId(eid: number): IntrEntityNode | undefined {
    return this.rootNode.getNodeByEntityId(eid);
  }

  getNodeByTFName(name: string): IntrEntityNode | undefined {
    return this.rootNode.getNodeByTFName(name);
  }

  getNodeByObjectId(oid: number): IntrRootEntityNode | undefined {
    return this.rootNode.getNodeByObjectId(oid);
  }

  getNodeByObjectName(name: string): IntrRootEntityNode | undefined {
    return this.rootNode.getNodeByObjectName(name);
  }

  setCollisionGeometryVisibleByEntityId(id: number, visible: boolean) {
    const node = this.getNodeByEntityId(id);
    for (const colliderShape of node?.collisionNode?.children ?? []) {
      colliderShape.visible = visible;
    }
    this.updatedSubject.next();
  }

  setCollisionGeometryVisibleByObjectId(id: number, visible: boolean) {
    const node = this.getNodeByObjectId(id);
    for (const entity of node?.traverseOwnEntities() ?? []) {
      this.setCollisionGeometryVisibleByEntityId(entity.entityId, visible);
    }
  }

  private framesVisible = true;

  setGridVisible(visible: boolean) {
    this.grid.visible = visible;
    this.updatedSubject.next();
  }

  setFramesVisible(visible: boolean) {
    this.framesVisible = visible;
    this.applyFramesVisible();
    this.updatedSubject.next();
  }

  private applyFramesVisible() {
    for (const frame of this.rootNode.traverseFrames()) {
      frame.setGeometryVisible(this.framesVisible);
    }
  }

  setFrameScale(scale: number) {
    for (const frame of this.rootNode.traverseFrames()) {
      frame.helper.scale.set(scale, scale, scale);
      frame.helper.updateMatrixWorld();
    }

    this.updatedSubject.next();
  }

  setHelpersVisible(visible: boolean) {
    for (const helper of this.rootNode.getHelpers()) {
      helper.visible = visible;
    }
  }

  clearMarginShapes() {
    for (const child of this.children) {
      if (child instanceof WorldRootNode) {
        child.clearMarginShapes();
      }
    }
    this.updatedSubject.next();
  }

  updateMarginShapes(marginValue: number) {
    for (const child of this.children) {
      if (child instanceof WorldRootNode) {
        child.updateMarginShapes(marginValue);
      }
    }
    this.updatedSubject.next();
  }

  getDebugObjects(): DebugObjects {
    return this.debugObjects;
  }
}
