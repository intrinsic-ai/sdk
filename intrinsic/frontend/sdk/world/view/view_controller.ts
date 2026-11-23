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

import {EffectRef, Injector, NgZone, Signal, computed} from '@angular/core';
import {explicitEffect} from '@intrinsic-ai/angular/signals';
import {Observable, Subject, Subscription, merge} from 'rxjs';
import {takeWhile} from 'rxjs/operators';
import * as three from 'three';

import {AssetDepot} from '../loaders/asset_depot';
import type {WorldModel} from '../model/world_model';
import {isMesh, isPoints} from '../utils/three_typechecks';
import {IntrinsicCamera} from './camera';
import {FrameTracker, FrameUpdateCallbackType} from './frame_tracker';
import {
  IntrOrbitControls,
  MotionType,
  PointerButton,
} from './intr_orbit_controls';
import {ObjectStateManager} from './object_state_manager';
import {PickerManager} from './picker_manager';
import {
  CAMERA_FAR_PLANE,
  CAMERA_NEAR_PLANE,
  ENVIRONMENT_ROTATION,
} from './renderer';
import {FrameAxes, SceneHelper} from './shapes';
import {IntrVisualMesh, IntrVisualPoints} from './shapes';
import {Transformer} from './transformer';
import {ViewWidget} from './view_widget';
import {RendererClient, WorldRendererService} from './world_renderer_service';
import {WorldScene} from './world_scene';
import {IntrEntityNode} from './world_scene_nodes';
import {IntrFrameNode, IntrRootEntityNode} from './world_scene_nodes';
import {
  CANVAS_CLASS,
  FPS_DISPLAY_CLASS,
  TOOLTIP_CLASS,
  TOOLTIP_TEXT_CLASS,
  VIEW_WIDGET_CLASS,
} from './world_viewer_styles';

/** Wraps camera view position information to save/restore views. */
declare interface ViewInformation {
  position: three.Vector3;
  target: three.Vector3;
}

// Dynamic scaling of near and far planes as we get closer to objects.
const MAX_DYNAMIC_PROJECTION_PLANES_SCALING = 0.01;
// Camera distance at which we begin reducing near and far planes.
const DYNAMIC_PROJECTION_MAX_DISTANCE = 1.0;
// Camera distance at which we hit the maximum reduction of near and far planes.
const DYNAMIC_PROJECTION_MIN_DISTANCE = 0.05;
const MAX_ROOT_BBOX_EXTENT = 5;
const MAX_ROOT_BBOX_SPAN = MAX_ROOT_BBOX_EXTENT * 2;
//The number of wheel 'ticks' sent for every zoom in/out event.
const ZOOM_FACTOR = 0.1;

// Set to true to help debug memory and performance issues. When true, the renderer
// will log stats to the console every couple of seconds. Should not be set to true on prod.
const LOG_RENDERING_STATS_PERIODICALLY = false;

/** Creates a default perspective camera with an attached directional light. */
export function createDefaultCamera(): IntrinsicCamera {
  const camera = IntrinsicCamera.createDefault();
  camera.position.x = 12;
  camera.position.y = 16;
  camera.position.z = 20;

  return camera;
}

/**
 * Event emitted when a 3D scene node is picked.
 */
export interface NodePickEvent {
  node?: IntrEntityNode | SceneHelper;
  event: MouseEvent;
}

class TargetWidget extends three.Group {
  private static readonly OFFSET = 0.015;
  private static readonly LENGTH = 0.05;

  private static createLinesGeometry() {
    const positions = [];
    for (let i = 0; i < 3; i++) {
      for (let j = 0; j < 2; j++) {
        const sign = j === 0 ? -1 : 1;
        positions.push(
          new three.Vector3().setComponent(i, sign * TargetWidget.OFFSET),
          new three.Vector3().setComponent(i, sign * TargetWidget.LENGTH),
        );
      }
    }

    return new three.BufferGeometry().setFromPoints(positions);
  }

  private static createRingGeometry() {
    const d = (TargetWidget.LENGTH + TargetWidget.OFFSET) / 2;
    return new three.RingGeometry(d / 2 - 0.003, d / 2 + 0.003, 16, 1);
  }

  private static readonly matConfig = {
    toneMapped: false,
    transparent: true,
    opacity: 0.8,
  };

  private static readonly cylinder = new three.Mesh(
    new three.CylinderGeometry(
      0.004,
      0.004,
      TargetWidget.LENGTH - TargetWidget.OFFSET,
      6,
    ),
    new three.MeshBasicMaterial({
      color: 0xffffff,
      side: three.BackSide,
      ...TargetWidget.matConfig,
    }),
  );

  private static makeCylinder(axis: three.Vector3) {
    const cyl = TargetWidget.cylinder.clone();
    cyl.setRotationFromQuaternion(
      new three.Quaternion().setFromUnitVectors(
        new three.Vector3(0, 1, 0),
        axis,
      ),
    );
    cyl.translateY((TargetWidget.LENGTH + TargetWidget.OFFSET) / 2);
    return cyl;
  }

  private static readonly ring = new three.Mesh(
    TargetWidget.createRingGeometry(),
    new three.MeshBasicMaterial({
      color: 0xffffff,
      depthTest: false,
      depthWrite: false,
      ...TargetWidget.matConfig,
    }),
  );

  private static readonly axes = new three.LineSegments(
    TargetWidget.createLinesGeometry(),
    new three.LineBasicMaterial({
      color: 0x000000,
      ...TargetWidget.matConfig,
    }),
  );

  private readonly ring = TargetWidget.ring.clone();
  private readonly axes = TargetWidget.axes.clone();

  constructor() {
    super();
    this.axes.position.setZ(0.0001);
    this.add(this.axes);
    this.add(this.ring);
    // add our axis cylinders
    this.add(TargetWidget.makeCylinder(new three.Vector3(0, 1, 0)));
    this.add(TargetWidget.makeCylinder(new three.Vector3(0, -1, 0)));
    this.add(TargetWidget.makeCylinder(new three.Vector3(0, 0, 1)));
    this.add(TargetWidget.makeCylinder(new three.Vector3(0, 0, -1)));
    this.add(TargetWidget.makeCylinder(new three.Vector3(1, 0, 0)));
    this.add(TargetWidget.makeCylinder(new three.Vector3(-1, 0, 0)));

    this.ring.onBeforeRender = (
      renderer: three.WebGLRenderer,
      scene: three.Scene,
      camera: three.Camera,
    ) => {
      this.ring.quaternion
        .copy(this.quaternion.clone().invert())
        .multiply(camera.quaternion);
      this.ring.updateMatrixWorld();
    };
  }

  override clone(): this {
    return new TargetWidget() as this;
  }
}

/** Dependencies of a `ViewController`. */
export interface ViewControllerDeps {
  // The element that the view fills. The rendered image is sized to match it,
  // and the view adds its overlays (gimbal, tooltip, frame rate) to it.
  container: HTMLElement;
  // The world to display.
  worldId: Signal<string>;
  // The WorldModel to display. Expected to be the model for `worldId`.
  worldModel: Signal<WorldModel | undefined>;
  // When true, polls the backend for updates to the WorldModel.
  polling: Signal<boolean>;
  viewNamespace: Signal<string>;
  showGimbal: Signal<boolean>;
  pickingEnabled: Signal<boolean>;
  cameraControlsEnabled: Signal<boolean>;
  // Shows a frames-per-second counter over the view.
  showFps: Signal<boolean>;
  // The camera used to render this view. Pass the same camera to several views
  // to share it between them.
  camera: Signal<IntrinsicCamera>;
  assetDepot: AssetDepot;
  pickerManager: PickerManager;
  objectStateManager: ObjectStateManager;
  worldRendererService: WorldRendererService;
  ngZone: NgZone;
  // Injector used to create the view controller's effects.
  injector: Injector;
}

/**
 * The ViewController class manages all changes to the virtual camera that is
 * used for rendering in WorldView. It can be considered an auxiliary class for
 * the WorldView which must set a camera for rendering.
 *
 * The view renders into a canvas that it adds to the given container, along
 * with its overlays (gimbal, tooltip, frame rate). These are styled by class,
 * so the owner must call `registerWorldViewerStyles` with the container once it
 * is attached. Call `dispose` when done.
 */
export class ViewController {
  readonly worldId: Signal<string>;
  readonly worldModel: Signal<WorldModel | undefined>;
  readonly polling: Signal<boolean>;
  readonly viewNamespace: Signal<string>;
  readonly showGimbal: Signal<boolean>;
  readonly pickingEnabled: Signal<boolean>;
  readonly cameraControlsEnabled: Signal<boolean>;
  readonly showFps: Signal<boolean>;
  // The camera used to render this view.
  readonly camera: Signal<IntrinsicCamera>;
  // The scene rendered by this view. It is owned by this view controller:
  // created with it and torn down when the view controller is disposed.
  private readonly worldScene: WorldScene;
  // Subscriptions to the PickerManager for the scene.
  private readonly pickerSubscription: Subscription;

  // The transformer of this view. Created with the view, after the canvas, and
  // switched to the current camera by the camera effect.
  private transformer?: Transformer;
  get scene(): WorldScene {
    return this.worldScene;
  }

  // The canvas the view renders into. Added to the container.
  private readonly canvas: HTMLCanvasElement;

  private renderClient!: RendererClient;
  // Controls of the current camera. Created by the camera effect.
  private orbitControls!: IntrOrbitControls;
  // Whether the camera effect has set up the current camera. Nothing is
  // rendered before that, to avoid rendering from a stale camera state.
  private cameraReady = false;
  // Whether the canvas has been sized to the container. Rendering before that
  // would show a stretched frame that is cleared once the canvas is resized.
  private canvasSized = false;
  private readonly targetWidget = new TargetWidget();
  // Notifies upon any camera changes.
  private readonly cameraChangeSubject = new Subject<void>();
  readonly cameraChanged: Observable<void> = this.cameraChangeSubject;

  // Notifies upon any node being clicked.
  private readonly nodeClickedSubject = new Subject<NodePickEvent>();
  readonly nodeClicked: Observable<NodePickEvent> = this.nodeClickedSubject;
  // Notifies upon any node being hovered.
  private readonly nodeHoveredSubject = new Subject<NodePickEvent>();
  readonly nodeHovered: Observable<NodePickEvent> = this.nodeHoveredSubject;

  // Flag indicating if user is dragging camera at the moment.
  private frozen = false;
  // Used to save and restore camera views.
  private readonly savedViews = new Map<string, ViewInformation>();

  private readonly frameTracker = new FrameTracker();
  private frameCallbackHandle = -1;

  // A text bubble displayed over the canvas, for example to show the name of a
  // hover picked object. Added to the DOM while it is visible.
  private readonly tooltip: Tooltip = createTooltipElement();

  readonly frameRate = computed<number>(() =>
    this.frameTracker.frameRateSignal(),
  );

  private readonly raycaster = new three.Raycaster();

  private readonly container: HTMLElement;
  private readonly injector: Injector;
  private readonly ngZone: NgZone;
  private readonly worldRendererService: WorldRendererService;
  private readonly pickerManager: PickerManager;
  // The effects created by this view controller. Destroyed on dispose.
  private readonly effects: EffectRef[] = [];
  private readonly resizeObserver: ResizeObserver;
  private disposed = false;

  readonly TEST_ONLY = {
    nodeClickedSubject: this.nodeClickedSubject,
    nodeHoveredSubject: this.nodeHoveredSubject,
  };

  constructor(deps: ViewControllerDeps) {
    this.canvas = createCanvas();
    this.worldId = deps.worldId;
    this.worldModel = deps.worldModel;
    this.polling = deps.polling;
    this.viewNamespace = deps.viewNamespace;
    this.showGimbal = deps.showGimbal;
    this.pickingEnabled = deps.pickingEnabled;
    this.cameraControlsEnabled = deps.cameraControlsEnabled;
    this.showFps = deps.showFps;
    this.camera = deps.camera;
    this.container = deps.container;
    this.injector = deps.injector;
    this.ngZone = deps.ngZone;
    this.worldRendererService = deps.worldRendererService;
    this.pickerManager = deps.pickerManager;
    this.worldScene = new WorldScene({
      worldId: this.worldId,
      worldModel: this.worldModel,
      polling: this.polling,
      assetDepot: deps.assetDepot,
      pickerManager: deps.pickerManager,
      objectStateManager: deps.objectStateManager,
      injector: deps.injector,
    });

    // Keep the view sized to the container.
    this.resizeObserver = new ResizeObserver(() => {
      this.refreshSize();
    });
    this.resizeObserver.observe(this.container);

    // Enable for all layers so that the OutlineRenderer doesn't mess up picking
    // as it uses multiple layers.
    this.raycaster.layers.enableAll();

    // Setup the scene. The picker subscriptions end in dispose.
    this.worldScene.add(this.targetWidget);
    this.pickerSubscription = this.setupPickerManagerSubscriptions(
      this.worldScene,
    );

    // Redraw the scene whenever the camera updates.
    this.addEffect([this.camera], ([camera], onCleanup) => {
      const sub = camera.updated.subscribe(() => {
        this.requestNewFrame();
      });
      onCleanup(() => {
        sub.unsubscribe();
      });
    });

    // Reframe the scene until the user moves the camera.
    this.addEffect([this.camera], ([camera], onCleanup) => {
      const sub = this.worldScene.mutated
        .pipe(takeWhile(() => !camera.touched))
        .subscribe(() => {
          this.frameRoot();
        });
      onCleanup(() => {
        sub.unsubscribe();
      });
    });

    // Handle enabling/disabling the orbit controls
    this.addEffect([this.cameraControlsEnabled], ([cameraControlsEnabled]) => {
      this.orbitControls?.setEnabled(cameraControlsEnabled);
    });

    this.container.appendChild(this.canvas);
    this.setupCanvas();

    // Set up the camera controls for the current camera, and tear them down
    // when the camera changes or the view is disposed.
    this.addEffect([this.camera], ([camera], onCleanup) => {
      const orbitControls = this.createOrbitControls(camera);
      this.orbitControls = orbitControls;
      this.transformer?.setCamera(camera);
      if (this.canvas.offsetHeight > 0) {
        this.updateAspectRatio(
          this.canvas.offsetWidth / this.canvas.offsetHeight,
        );
      }
      this.handleCameraChange();
      this.cameraReady = true;
      onCleanup(() => {
        orbitControls.dispose();
      });
    });
  }

  /** Creates the orbit controls that move the given camera. */
  private createOrbitControls(camera: IntrinsicCamera): IntrOrbitControls {
    const orbitControls = new IntrOrbitControls(camera, this.canvas, {
      pan: PointerButton.LEFT,
      zoom: PointerButton.RIGHT,
      orbit: PointerButton.MIDDLE,
    });

    orbitControls.setEnabled(!this.frozen && this.cameraControlsEnabled());
    orbitControls.addEventListener('start', () => {
      camera.touched = true;
    });

    orbitControls.addEventListener('change', (e) => {
      if (e.motionType === MotionType.PAN) {
        this.targetWidget.quaternion.identity();
        this.targetWidget.updateMatrix();
      }
      if (e.motionType === MotionType.ZOOM) {
        const oldPos = e.oldPos;
        const oldTarget = e.oldTarget;
        const scaleRatio =
          orbitControls.getDistance() / oldPos.distanceTo(oldTarget);
        camera.setOrthoDimension(scaleRatio);
      }
      this.handleCameraChange();
    });
    return orbitControls;
  }

  /** Creates an effect that is destroyed when this view is disposed. */
  private addEffect<Values>(
    dependencies: Parameters<typeof explicitEffect<Values>>[0],
    fn: Parameters<typeof explicitEffect<Values>>[1],
  ) {
    this.effects.push(
      explicitEffect(dependencies, fn, {injector: this.injector}),
    );
  }

  /**
   * Tears down the view: stops rendering, removes the overlays, and destroys
   * the effects and the scene. Safe to call more than once.
   */
  dispose() {
    if (this.disposed) return;
    this.disposed = true;
    for (const effectRef of this.effects) {
      effectRef.destroy();
    }
    this.effects.length = 0;
    this.resizeObserver.disconnect();
    this.transformer?.dispose();
    this.transformer = undefined;
    cancelAnimationFrame(this.frameCallbackHandle);
    this.pickerSubscription.unsubscribe();
    this.worldScene.teardown();
    this.tooltip.element.remove();
    this.canvas.remove();
  }

  /** Sets up rendering into the canvas and starts the render loop. */
  private setupCanvas() {
    // Set up renderer.
    this.renderClient = this.worldRendererService.createClient(this.canvas);
    this.renderClient?.setClearColor(new three.Color(0, 0, 0), 0);
    this.updateEnvMap();
    // The container may not be laid out yet. Its size is then set once it is.
    if (this.canvas.offsetHeight > 0) {
      this.updateAspectRatio(
        this.canvas.offsetWidth / this.canvas.offsetHeight,
      );
    }
    if (LOG_RENDERING_STATS_PERIODICALLY) {
      this.worldRendererService.logRenderingStatsPeriodically();
    }

    // Prevent the context menu from popping up for right click or cltr-left
    this.canvas.addEventListener('contextmenu', (e) => {
      e.preventDefault();
    });
    this.canvas.addEventListener('pointerdown', (e: PointerEvent) => {
      if (e.button === 0 && e.shiftKey) {
        const coords = this.getNdcFromInputCoord(e.offsetX, e.offsetY);
        this.setTargetFromRaycast(coords);
      }
    });
    this.canvas.addEventListener('pointermove', this.handleMouseMoveEvent);
    this.canvas.addEventListener('pointerup', this.handleMouseClickEvent);

    // Create the transformer associated with this view. Note that this should only
    // be created after the canvas is ready because it registers canvas event listeners.
    this.transformer = new Transformer(this);

    // Kick off the main loop.
    this.mainLoop();

    this.setupOverlays();
  }

  /**
   * Shows the view widget (gimbal) and the frame rate over the main canvas.
   * The overlays are added after the canvas, so they are drawn on top of it.
   */
  private setupOverlays() {
    // Show the view widget (gimbal) in its own canvas.
    this.addEffect(
      [this.showGimbal, this.camera],
      ([showGimbal, camera], onCleanup) => {
        if (!showGimbal) return;
        const viewWidgetCanvas = createViewWidgetCanvas();
        this.container.appendChild(viewWidgetCanvas);
        const viewWidget = this.ngZone.runOutsideAngular(
          () =>
            new ViewWidget(viewWidgetCanvas, camera, this.worldRendererService),
        );
        const sub = viewWidget.quaternionUpdated.subscribe((quat) => {
          this.updateQuaternion(quat);
        });
        onCleanup(() => {
          sub.unsubscribe();
          viewWidget.dispose();
          viewWidgetCanvas.remove();
        });
      },
    );

    // Show the frame rate.
    const fpsDisplay = createFpsDisplayElement();
    this.addEffect([this.showFps], ([showFps], onCleanup) => {
      if (!showFps) return;
      this.container.appendChild(fpsDisplay);
      onCleanup(() => {
        fpsDisplay.remove();
      });
    });
    this.addEffect([this.frameRate], ([frameRate]) => {
      fpsDisplay.textContent = `FPS: ${frameRate}`;
    });
  }

  private setupPickerManagerSubscriptions(scene: WorldScene): Subscription {
    const sub = new Subscription();

    sub.add(
      merge(
        this.pickerManager.newPick,
        this.pickerManager.newHoverPick,
        this.pickerManager.newEntityPick,
        this.pickerManager.newHoverEntityPick,
      ).subscribe(() => {
        this.requestNewFrame();
      }),
    );

    sub.add(
      this.pickerManager.newPick.subscribe(
        (pickedObjectIds: ReadonlySet<number>) => {
          this.setTooltipVisible(false);
          const oid = [...pickedObjectIds.values()][0] ?? 0;
          const obj = scene.getNodeById(oid);
          if (
            obj instanceof IntrRootEntityNode ||
            obj instanceof IntrFrameNode
          ) {
            const pos = new three.Vector3();
            obj.getWorldPosition(pos);
            const name =
              obj instanceof IntrRootEntityNode ? obj.objectName : obj.name;
            this.showToolTipIn3D(name, pos);
          }
        },
      ),
    );

    sub.add(
      this.pickerManager.newHoverPick.subscribe((objectId: number) => {
        this.setTooltipVisible(false);
        const obj = scene.getNodeById(objectId);
        if (obj instanceof IntrRootEntityNode || obj instanceof IntrFrameNode) {
          const pos = new three.Vector3();
          obj.getWorldPosition(pos);
          const name =
            obj instanceof IntrRootEntityNode ? obj.objectName : obj.name;
          this.showToolTipIn3D(name, pos);
        }
      }),
    );

    return sub;
  }

  private updateEnvMap() {
    if (this.renderClient) {
      this.worldScene.environmentRotation.copy(ENVIRONMENT_ROTATION);
      this.worldScene.environment =
        this.renderClient.computePlaceholderEnvMap();
      this.renderClient
        .loadEnvironmentMap()
        .then((envMap) => {
          this.worldScene.environment = envMap;
          this.requestNewFrame();
        })
        .catch((e) => {
          console.error('Failed to download env map', e);
        });
      this.requestNewFrame();
    }
  }

  updateQuaternion(quat: three.Quaternion) {
    this.orbitTo(quat);
  }

  private updateTarget() {
    if (!this.orbitControls?.target) return;
    this.targetWidget.position.copy(this.orbitControls.target);
    // Scale by the length of a screen space unit-vector.
    const up = new three.Vector3(0, 1, 0)
      .applyQuaternion(this.camera().quaternion)
      .add(this.orbitControls.target)
      .project(this.camera());
    const scale = this.orbitControls.target
      .clone()
      .project(this.camera())
      .distanceTo(up);
    this.targetWidget.scale.setScalar(1 / scale);
    this.targetWidget.updateMatrix();
  }

  /**
   * Explicitly request a new frame to be rendered.
   */
  requestNewFrame() {
    this.updateTarget();
    this.frameTracker.requestNewFrame();
  }

  registerFrameUpdateCallback(
    callback: VoidFunction,
    type: FrameUpdateCallbackType,
  ): number {
    return this.frameTracker.registerFrameUpdateCallback(callback, type);
  }

  unregisterFrameUpdateCallback(id: number) {
    this.frameTracker.unregisterFrameUpdateCallback(id);
  }

  newFrameRequested() {
    return this.frameTracker.updateFrame();
  }

  // Directly set a new perspective camera.
  restore(position: three.Vector3, target: three.Vector3) {
    this.camera().position.copy(position);
    this.setCameraPosTarget(position, target);
  }

  setFrozen(frozen: boolean) {
    this.frozen = frozen || !this.cameraControlsEnabled();
    this.orbitControls.setEnabled(!this.frozen);
  }

  getFrozen(): boolean {
    return this.frozen;
  }

  getTarget(): three.Vector3 {
    return this.orbitControls.target.clone();
  }

  getDistance(): number {
    return this.orbitControls.getDistance();
  }

  zoom(newZoom: number) {
    const distance = this.orbitControls.getDistance();
    this.orbitControls.zoomCamera(newZoom);
    const scaleRatio = this.orbitControls.getDistance() / distance;
    this.camera().setOrthoDimension(scaleRatio);
    this.requestNewFrame();
  }

  changeZoom(zoomIn: boolean) {
    this.zoom(zoomIn ? -ZOOM_FACTOR : ZOOM_FACTOR);
  }

  // Update camera aspect ratio when is window is resized.
  updateAspectRatio(aspect: number) {
    if (this.frozen) return;

    this.camera().aspect = aspect;
    this.camera().updateProjectionMatrix();
  }

  /**
   * Get the normalized display coordinate (NDC) from mouse input coord.
   * Input coord: top-left (0, 0) to bottom-right (rect.width, rect.height);
   * NDC: bottom-left (-1, -1) to top-right (1, 1);
   */
  getNdcFromInputCoord(x: number, y: number): three.Vector2 {
    const rect = this.canvas.getBoundingClientRect();
    const ndcX = (x / rect.width) * 2 - 1;
    const ndcY = -(y / rect.height) * 2 + 1;
    return new three.Vector2(ndcX, ndcY);
  }

  // Returns true if camera is moving.
  isMoving() {
    if (this.frozen) return false;
    return this.orbitControls.getMoving();
  }

  setCameraPose(position: three.Vector3, quaternion: three.Quaternion) {
    this.orbitControls.setPose(position, quaternion);
    this.requestNewFrame();
    this.cameraChangeSubject.next();
  }

  orbitTo(quat: three.Quaternion) {
    const pos = this.camera().position.clone().sub(this.orbitControls.target);
    const cameraQuat = this.camera().quaternion.clone();
    const q = quat.clone().multiply(cameraQuat.invert());
    pos.applyQuaternion(q);
    pos.add(this.orbitControls.target);
    this.setCameraPose(pos, quat);
  }

  // Scale the near and far plane distance dynamically as the distance from
  // the orbitControls target shrinks. This allows us to zoom very close to
  // objects without sacrificing view range at normal view distances.
  private updateCameraProjectionPlanes() {
    const targetDistance = this.camera().position.distanceTo(
      this.orbitControls.target,
    );
    // Remap distance [5cm to 1m] to [0.0, 1.0] range.
    const clampedDistance = Math.min(
      Math.max(targetDistance - DYNAMIC_PROJECTION_MIN_DISTANCE, 0.0),
      DYNAMIC_PROJECTION_MAX_DISTANCE,
    );
    // Reduce near and far plane by the computed factor
    const distanceRatio = Math.max(
      clampedDistance,
      MAX_DYNAMIC_PROJECTION_PLANES_SCALING,
    );
    this.camera().near = CAMERA_NEAR_PLANE * distanceRatio;
    this.camera().far = CAMERA_FAR_PLANE * distanceRatio;
    this.camera().updateProjectionMatrix();
  }

  // Call this when the internal camera's values are updated.
  private handleCameraChange() {
    this.updateCameraProjectionPlanes();
    // This makes sure the light that attaches to the camera follows.
    this.camera().updateMatrixWorld(true);
    // Inform subscribers of camera change, provide a clone of current camera.
    this.cameraChangeSubject.next();
    this.frameTracker.requestNewFrame();
  }

  private setTargetFromRaycast(coords: three.Vector2) {
    if (!this.worldScene.root) {
      return;
    }
    const nodes: three.Object3D[] = [];
    this.worldScene.root.traverseVisible((n) => {
      if (isMesh(n) || isPoints(n)) {
        nodes.push(n);
      }
    });
    const quat = this.camera().quaternion.clone();
    this.raycaster.setFromCamera(coords, this.camera());
    const intersection = this.raycaster.intersectObjects(nodes)[0];
    if (intersection) {
      this.orbitControls.target.copy(intersection.point);
      const normal = intersection.face?.normal;
      if (normal) {
        // The normal returned from the raycast is in object space, so we must first transform it to world space.
        normal.applyQuaternion(
          intersection.object.getWorldQuaternion(new three.Quaternion()),
        );
        quat.setFromUnitVectors(new three.Vector3(0, 0, 1), normal);
      }
    }
    this.targetWidget.quaternion.copy(quat);
    this.targetWidget.updateMatrix();

    this.requestNewFrame();
  }

  frameRoot() {
    if (!this.worldScene.root) return;
    this.frameObjects([this.worldScene.root]);
  }

  // Compute a camera view to look over the bounding box of the given scene.
  // This can be used when loading the world scene, or when user picked
  // one or more objects.
  // If the given `objects` is empty, we frame the whole scene.
  frameObjects(
    objects: Iterable<IntrEntityNode>,
    preserveViewDirection = false,
  ) {
    this.targetWidget.quaternion.identity();
    if (this.frozen) return;

    const objectsToFrame = Array.from(objects);
    if (objectsToFrame.length === 0 && this.worldScene.root) {
      objectsToFrame.push(this.worldScene.root);
    }

    let cumulativeBox = null;
    for (const object of objectsToFrame) {
      const bbox = object.getVisualBounds(true);
      const span = new three.Vector3();
      span.subVectors(bbox.max, bbox.min);

      // Do nothing if the bbox is infinite.
      if (!isFinite(span.x) || !isFinite(span.y) || !isFinite(span.z)) {
        continue;
      }

      (cumulativeBox ??= bbox).union(bbox);
    }

    // If the bounding box exceeds the max allowed size, it is likely due to a
    // geometry error in some asset (such as geometry with a point way off in space). In this
    // case, recover by centering the camera on the origin.
    if (cumulativeBox) {
      const span = new three.Vector3();
      span.subVectors(cumulativeBox.max, cumulativeBox.min);
      if (span.length() > MAX_ROOT_BBOX_SPAN) {
        cumulativeBox.min.setScalar(-MAX_ROOT_BBOX_EXTENT);
        cumulativeBox.max.setScalar(MAX_ROOT_BBOX_EXTENT);
      }
    }

    // If we never computed a valid box, then just frame a dumb 2x2x2 box around
    // the origin.
    cumulativeBox ||= new three.Box3(
      new three.Vector3(-1, -1, -1),
      new three.Vector3(1, 1, 1),
    );

    this.frameBoundingBox(cumulativeBox, preserveViewDirection, 0.5);
  }

  frameObjectsById(objectIds: Iterable<number>, preserveViewDirection = false) {
    const objs: IntrEntityNode[] = [...objectIds].flatMap((id) => {
      return this.worldScene.getNodeById(id) ?? [];
    });
    this.frameObjects(objs, preserveViewDirection);
  }

  // Compute a camera view to look over the specified bounding box.
  frameBoundingBox(
    bbox: three.Box3,
    preserveViewDirection = false,
    minDistance = 0.5,
  ) {
    // Get the center of bbox for camera to lookat.
    const center = bbox.getCenter(new three.Vector3());
    // Compute span of the cumulative box
    const span = bbox.getSize(new three.Vector3());
    // Compute a view distance that guarantee to encompass the span of the
    // worldscene with current camera fov. Limit the distance to the specified
    // minimum to avoid extremely jarring changes of scale.
    const distance = Math.max(
      minDistance,
      (0.7 * Math.max(span.x, span.y, span.z)) /
        Math.tan((Math.PI * this.camera().fov) / 360),
    );
    const newPosition = new three.Vector3();
    if (preserveViewDirection) {
      // Compute normalized view direction vector.
      const viewDir = this.camera().position.clone().sub(center).normalize();

      // Set the new position to be 'distance' away from the focus point
      // preserving the current view direction.
      newPosition.copy(viewDir.multiplyScalar(distance).add(center));
    } else {
      // Look direction is parallel to the y-z plane, 26.6 degrees tilt.
      newPosition.set(center.x, center.y - distance, center.z + distance / 2);
    }
    this.setCameraPosTarget(newPosition, center);
    this.handleCameraChange();
  }

  /**
   * Aligns the camera view with the specified threejs object .
   */
  attachViewToObject(obj: three.Object3D) {
    const pos = new three.Vector3();
    obj.getWorldPosition(pos);
    const dir = new three.Vector3();
    obj.getWorldDirection(dir);
    this.setCameraPosTarget(pos.add(dir), dir);
  }

  /** Save the view using a user-provided string as key. */
  saveView(name: string) {
    const viewInfo: ViewInformation = {
      position: this.camera().position.clone(),
      target: this.orbitControls.target.clone(),
    };
    this.savedViews.set(name, viewInfo);
  }

  private setCameraPosTarget(
    pos: three.Vector3,
    target: three.Vector3,
    center = true,
  ) {
    const oldDistance = this.getDistance();
    this.camera().position.copy(pos);
    this.orbitControls.target.copy(target);
    if (center) {
      this.camera().lookAt(target);
    }
    this.camera().setOrthoDimension(this.getDistance() / oldDistance);
    this.handleCameraChange();
  }

  /** Restore the camera view using a user provided string as key. */
  restoreView(name: string) {
    if (!this.savedViews.has(name)) return;
    const viewInfo = this.savedViews.get(name);
    if (!viewInfo) return;
    this.setCameraPosTarget(viewInfo.position, viewInfo.target);
  }

  getScreenCoordFromPos(pos: three.Vector3): three.Vector2 {
    // Get the NDC coordinate which is from [-1, 1], with X axis pointing left
    // to right, Y axis bottom to up.
    pos.project(this.camera());
    const canvasBound = this.canvas.getBoundingClientRect();
    // Map from in cavas NDC to the actual screen coordinate used by JS:
    // - First map from [-1, 1] to [0, 1], note Y direction need to be
    // inverted.
    pos.x = 0.5 * (1.0 + pos.x);
    pos.y = 0.5 * (1.0 - pos.y);
    // - Then compute absolute window location by offsetting canvas's
    // position.
    return new three.Vector2(
      pos.x * canvasBound.width,
      pos.y * canvasBound.height,
    );
  }

  /**
   * @param txt Content of the tooltip.
   * @param coord Screen position of the top-left anchor of the tooltip.
   * @param tooltipWidth Optional width of the tooltip.
   */
  showToolTipIn2D(txt: string, coord: three.Vector2, tooltipWidth = 0) {
    this.tooltip.text.textContent = txt;
    const style = this.tooltip.element.style;
    style.left = `${coord.x}px`;
    style.top = `${coord.y}px`;
    style.width = tooltipWidth === 0 ? '' : `${tooltipWidth}px`;
    this.setTooltipVisible(true);
  }

  setTooltipVisible(visible: boolean) {
    if (!visible) {
      this.tooltip.element.remove();
      return;
    }
    // The tooltip goes on top of everything else in the view.
    if (!this.tooltip.element.isConnected) {
      this.container.appendChild(this.tooltip.element);
    }
  }

  /**
   * Utility to show tooltip at the position of a threejs  object. Calls into
   * showToolTipIn3D and keeps threejs depedendency out of the caller.
   * @param txt Content of the tooltip
   * @param obj threejs Object3D
   */
  showTooltipAtObjPosition(txt: string, obj: three.Object3D | undefined) {
    if (!obj) return;
    const pos = new three.Vector3();
    obj.getWorldPosition(pos);
    this.showToolTipIn3D(txt, pos);
  }

  /**
   * @param txt Content of the tooltip.
   * @param pos World position of the top-left anchor of the tooltip.
   * @param tooltipWidth Optional width of the tooltip.
   */
  showToolTipIn3D(txt: string, pos: three.Vector3, tooltipWidth = 0) {
    const coord = this.getScreenCoordFromPos(pos);
    this.showToolTipIn2D(txt, coord, tooltipWidth);
  }

  /**
   * Returns the transformer of this view. The transformer is owned by this view
   * and disposed together with it.
   */
  getTransformer(): Transformer {
    if (!this.transformer) {
      throw new Error('The view is not initialized yet.');
    }
    return this.transformer;
  }

  // Return the HTML Canvas element, this can be used for simple tasks such as
  // getting the size of canvas, or sophisticated ones such as drawing extra
  // data to the same canvas (e.g.: HUD) or apply post-processing effects.
  getCanvasElement(): HTMLCanvasElement {
    return this.canvas;
  }

  private mainLoop() {
    // Keep frame requests until the camera and canvas are set up, then render
    // them.
    if (this.cameraReady && this.canvasSized) {
      const wantNewFrame = this.newFrameRequested();
      // While the frame rate is shown, render continuously to measure it.
      if (wantNewFrame || this.showFps()) {
        this.render();
      }
    }

    this.ngZone.runOutsideAngular(() => {
      this.frameCallbackHandle = requestAnimationFrame(() => {
        this.mainLoop();
      });
    });
  }

  // Update viewport size and camera aspect ratio when is window is resized.
  refreshSize() {
    const width = this.container.offsetWidth;
    const height = this.container.offsetHeight;
    this.canvasSized = width > 0 && height > 0;
    if (!this.canvasSized) return;
    this.renderClient.setSize(width, height);
    this.updateAspectRatio(width / height);
    this.requestNewFrame();
  }
  getBaseWorldId() {
    return this.worldScene.root.getWorldId();
  }

  // Render current scene from current camera immediately.
  render() {
    // This makes rendering much faster.
    this.renderClient.render(this.worldScene, this.camera());
  }

  private readonly handleMouseMoveEvent = (e: MouseEvent) => {
    // Skip hover picking while camera is moving, this improves frame rate.
    if (this.isMoving() || !this.pickingEnabled()) {
      return;
    }
    const node = this.pickNodeFromMouseCoords(
      this.getNdcFromInputCoord(e.offsetX, e.offsetY),
    );
    if (node instanceof SceneHelper) {
      const pos = node.getWorldPosition(new three.Vector3());
      node.enableHighlight(true);
      this.showToolTipIn3D(
        node.hoverText,
        pos.add(new three.Vector3(0.015, 0, 0.015)),
      );
    }

    this.nodeHoveredSubject.next({node, event: e});
    return;
  };

  /**
   * Mouse click would pick the first viewray-intersected object.
   * @param e A mouse event object that contains picked location.
   * @returns Returns true if actually intersectioned object.
   */
  private readonly handleMouseClickEvent = async (
    e: MouseEvent,
  ): Promise<void> => {
    this.setTooltipVisible(false);
    if (this.isMoving() || !this.pickingEnabled()) {
      return;
    }

    // Ignore middle mouse button events
    if (e.button === 1) {
      return;
    }

    // Only pick for left click. Right click clears previous picked object.
    if (e.button !== 0) {
      this.pickerManager.clearAllPicks();
      return;
    }
    const node = this.pickNodeFromMouseCoords(
      this.getNdcFromInputCoord(e.offsetX, e.offsetY),
    );
    this.nodeClickedSubject.next({node, event: e});
  };

  private rayCastSceneNDC(normalizedPosition: three.Vector2) {
    this.raycaster.setFromCamera(normalizedPosition, this.camera());
    const nodes: three.Object3D[] = [];
    this.worldScene.traverseVisible((n) => {
      if (
        n instanceof IntrVisualMesh ||
        n instanceof IntrVisualPoints ||
        n instanceof FrameAxes ||
        n instanceof SceneHelper
      ) {
        if (!n.isInvisible()) {
          nodes.push(n);
        }
      }
    });
    const intersections = this.raycaster.intersectObjects(nodes, true);
    return intersections;
  }

  private pickNodeFromMouseCoords(
    normalizedPosition: three.Vector2,
  ): IntrEntityNode | SceneHelper | undefined {
    const intersections = this.rayCastSceneNDC(normalizedPosition);
    if (intersections.length === 0) {
      return undefined;
    }
    // Walk up until we find something we care about. Preferably click on scene helpers.
    let helper = undefined;
    const sceneNodes: Array<IntrEntityNode | SceneHelper> = intersections
      .map((i) => {
        let current: three.Object3D | null = i.object;
        while (current) {
          if (current instanceof IntrEntityNode) {
            return current;
          }
          if (current instanceof SceneHelper) {
            helper = current;
            return current;
          }
          current = current.parent;
        }
        return undefined;
      })
      .filter((n) => !!n);
    return helper ?? sceneNodes[0];
  }

  /**
   * Attempts to pick an IntrEntityNode from the given normalized mouse coordinates.
   * Returns undefined if no entity is picked.
   */
  pickEntityNodeFromMouseCoords(
    normalizedPosition: three.Vector2,
  ): IntrEntityNode | undefined {
    const intersections = this.rayCastSceneNDC(normalizedPosition);
    if (intersections.length === 0) {
      return undefined;
    }
    // See if the nearest Object3D picked is part of some IntrEntityNode
    return IntrEntityNode.findParentEntity(intersections[0].object);
  }

  getNamespace() {
    return this.viewNamespace();
  }
}

/** Creates the (detached) canvas that a ViewController renders into. */
function createCanvas(): HTMLCanvasElement {
  const canvas = document.createElement('canvas');
  canvas.className = CANVAS_CLASS;
  // Allows the canvas to grab focus when clicked.
  canvas.tabIndex = 0;
  return canvas;
}

/** The tooltip of a ViewController. */
interface Tooltip {
  // The tooltip itself. In the DOM while the tooltip is visible.
  element: HTMLDivElement;
  // The element containing the tooltip text.
  text: HTMLElement;
}

/** Creates the (detached) tooltip element of a ViewController. */
function createTooltipElement(): Tooltip {
  const tooltip = document.createElement('div');
  tooltip.className = TOOLTIP_CLASS;
  const text = tooltip.appendChild(document.createElement('span'));
  text.className = TOOLTIP_TEXT_CLASS;
  return {element: tooltip, text};
}

/** Creates the (detached) frame rate display of a ViewController. */
function createFpsDisplayElement(): HTMLDivElement {
  const fpsDisplay = document.createElement('div');
  fpsDisplay.className = FPS_DISPLAY_CLASS;
  return fpsDisplay;
}

/** Creates the (detached) canvas of the view widget (gimbal). */
function createViewWidgetCanvas(): HTMLCanvasElement {
  const canvas = document.createElement('canvas');
  canvas.className = VIEW_WIDGET_CLASS;
  canvas.width = 1;
  canvas.height = 1;
  return canvas;
}
