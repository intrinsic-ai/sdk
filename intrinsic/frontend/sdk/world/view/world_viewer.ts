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

import type {
  EffectCleanupRegisterFn,
  EffectRef,
  Injector,
  NgZone,
  Signal,
} from '@angular/core';
import {computed, signal} from '@angular/core';
import {explicitEffect} from '@intrinsic-ai/angular/signals';
import {ReplaySubject, Subscription, merge} from 'rxjs';
import type * as three from 'three';

import type {WorldEditor} from '../editor/world_editor';
import {AssetDepot} from '../loaders/asset_depot';
import {WorldModel} from '../model/world_model';
import type {IntrinsicCamera} from './camera';
import {ObjectStateManager} from './object_state_manager';
import {PickerManager} from './picker_manager';
import {RigidObjectTarget} from './rigid_object_target';
import {SceneHelper} from './shapes';
import type {Transformer} from './transformer';
import {UserInputModeManager} from './user_input_mode_manager';
import type {NodePickEvent} from './view_controller';
import {ViewController} from './view_controller';
import {WorldRendererService} from './world_renderer_service';
import {WorldScene} from './world_scene';
import {IntrEntityNode, IntrFrameNode} from './world_scene_nodes';
import type {WorldViewerOptions} from './world_viewer_options';
import type {WorldViewerPlugin} from './world_viewer_plugin';
import {
  CONTAINER_CLASS,
  registerWorldViewerStyles,
} from './world_viewer_styles';
import {WorldWorkspaceService} from './world_workspace_service';

/**
 * Representation of a camera view.
 */
export interface CameraView {
  pos: three.Vector3;
  target: three.Vector3;
}

/** Dependencies of a `WorldViewer`. */
export interface WorldViewerDeps {
  // The element that the viewer fills. It renders the world into it.
  container: HTMLElement;
  // The world to display.
  worldId: Signal<string>;
  // Whether the world can be edited, e.g. by moving objects with the
  // transformer.
  editable: Signal<boolean>;
  // Key under which the view is registered with the WorldWorkspaceService.
  // Empty to not register it.
  viewNamespace: Signal<string>;
  // Whether to show the ground plane.
  showGroundPlane: Signal<boolean>;
  // Whether to show the orientation gimbal.
  showGimbal: Signal<boolean>;
  // Whether to show frames.
  showFrames: Signal<boolean>;
  // Whether objects can be picked (hovered and selected) with the pointer.
  pickingEnabled: Signal<boolean>;
  // Whether the user can move the camera with the pointer.
  cameraControlsEnabled: Signal<boolean>;
  // The camera used to render the world. Pass the same camera to several
  // viewers to share it between them.
  camera: Signal<IntrinsicCamera>;
  // When true, the WorldRootNode polls the backend for updates to the
  // WorldModel and applies them to the scene.
  pollForWorldUpdates: Signal<boolean>;
  // Holds the workspaces (WorldModel, WorldEditor) of the open worlds, and the
  // registered views.
  worldWorkspaceService: WorldWorkspaceService;
  // Holds per-object view and edit flags, e.g. which objects are hidden or
  // locked.
  objectStateManager: ObjectStateManager;
  // Holds the hovered and selected (picked) objects.
  pickerManager: PickerManager;
  // Loads the geometry and other assets of the world's objects.
  assetDepot: AssetDepot;
  // Shared offscreen renderer, so all views share one WebGL context.
  worldRendererService: WorldRendererService;
  // Application-wide world viewer options.
  options: WorldViewerOptions;
  // Zone used to run rendering and event handling outside Angular.
  ngZone: NgZone;
  // Injector used to create the viewer's effects.
  injector: Injector;
}

/**
 * This class is the main world-viewer component. It takes as input a world to
 * visualize and displays a live 3d rendering of it using threejs. Internally,
 * it listens for updates to the world through the WorldModel interface and updates
 * the threejs scene to keep it in sync with the model.
 *
 * Overview of subcomponents:
 * - ViewController: Contains methods to get and manipulate the camera, frame objects, pick objects, etc.
 * - WorldScene: this is our wrapper around a threejs Scene.
 * - WorldRootNode: The root threejs object of the WorldScene. Contains nodes for all links, joints, frames, etc.
 *  It updates dynamically as the world changes. See `world_scene_nodes.ts` for the different node types and
 * for how to traverse the node tree.
 * - PickerManager: Manages node selection state. Used for reacting to click and hover events.
 *
 * Plugin system:
 * The WorldViewer has a simple plugin system for extending its functionality. See
 * the `registerPlugin`/`unregisterPlugin` methods below. This is useful for adding overlays,
 * visual markers, etc.
 *
 * NOTE:
 * The WorldViewer component and related subcomponents should be kept as minimal as possible.
 * This is the main component of the third-party world-viewer-sdk. Any code here will become a dependency
 * for anyone using the SDK. Avoid depending on backend services, and, whenever possible, add new
 * functionality as optional plugins (see WorldViewPlugin) instead of adding it here.
 *
 * NOTE: The WorldViewer is intentionally not an Angular component. It is meant
 * to hook into a given container element in a generic way so that third-party
 * users may use it with whatever frontend framework they prefer. We currently
 * have Angular dependencies for signals and effects but the goal is to
 * gradually remove these and be fully framework-agnostic (b/544815723).
 * Angular apps use it through the `<world-viewer>` component in
 * `world_viewer/angular`.
 *
 * Usage:
 *
 * ```ts
 * const viewer = new WorldViewer({
 *   container,  // The element to render into.
 *   worldId: signal('world'),
 *   editable: signal(false),
 *   ...  // Other inputs and services, see WorldViewerDeps.
 * });
 * // Once the container is attached to the document:
 * viewer.init();
 * ...
 * // When done with the viewer:
 * viewer.dispose();
 * ```
 */
export class WorldViewer {
  readonly worldId: Signal<string>;
  readonly editable: Signal<boolean>;
  readonly viewNamespace: Signal<string>;

  readonly showGroundPlane: Signal<boolean>;
  readonly showGimbal: Signal<boolean>;
  readonly showFrames: Signal<boolean>;

  readonly pickingEnabled: Signal<boolean>;
  readonly cameraControlsEnabled: Signal<boolean>;

  // The camera used to render the world.
  readonly camera: Signal<IntrinsicCamera>;

  // When set to true, the WorldRootNode will poll the backend for
  // updates to the WorldModel and apply them to the scene.
  readonly pollForWorldUpdates: Signal<boolean>;

  // For compatibility with the old non-signals interface, we expose this
  // signal through a public getter.
  protected readonly sceneSignal: Signal<WorldScene> = computed(
    () => this.viewController.scene,
  );
  get scene(): WorldScene {
    return this.sceneSignal();
  }

  getSceneSignal(): Signal<WorldScene> {
    return this.sceneSignal;
  }

  // Dependencies:
  private readonly worldWorkspaceService: WorldWorkspaceService;
  private readonly objectStateManager: ObjectStateManager;
  readonly pickerManager: PickerManager;
  private readonly injector: Injector;

  // Manages this viewer's user input modes, e.g. the measurement tool.
  private readonly userInputModeManager: UserInputModeManager;

  loaded = new ReplaySubject<boolean>(1);
  private viewInitialized = false;

  // NOTE we assert the properties below out of laziness to due to
  // initialization occurring in the view init hook.  It's generally safe as
  // WorldViewers won't be available to most components until the `loaded`
  // subject emits once.

  // Actor that supports undo-redoable actions to mutate a world.
  // This needs to remain private. Only access editors though
  // WorldWorkspaceService.getEditorAsync(..).
  editor!: WorldEditor;
  get model(): WorldModel {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return this.worldModelSignal()!;
  }

  // We set this signal when the `model` loads, and expose it through
  // `getWorldModelSignal`.
  protected readonly worldModelSignal = signal<WorldModel | undefined>(
    undefined,
  );

  // Transformer enables rigid transform of a particular object in the scene.
  private transformer?: Transformer;

  private pickedObjectIds: number[] = [];

  // Picking is off when the host disables it, or while the active input mode
  // handles pointer input itself.
  protected readonly scenePickingEnabled = computed(
    () =>
      this.pickingEnabled() &&
      !this.userInputModeManager.activeMode()?.disablesScenePicking,
  );

  // ViewController contains access to the canvas and camera controls. It is
  // disposed together with the WorldViewer.
  readonly viewController: ViewController;
  private readonly container: HTMLElement;

  // The effects and subscriptions of this viewer. Ended on dispose.
  private readonly effects: EffectRef[] = [];
  private readonly subscriptions = new Subscription();
  private disposed = false;

  // Plugins registered on this WorldViewer, for extending its functionality.
  // The list is in order of registration. Plugins are unregistered in reverse order
  // on cleanup.
  private readonly plugins: WorldViewerPlugin[] = [];
  // Plugins whose registration was requested before view initialization.
  // We defer registering plugins until the view is initialized to ensure that
  // the WorldViewer is in a known-good state before plugins access and/or modify it.
  private readonly deferredPlugins: WorldViewerPlugin[] = [];

  constructor(deps: WorldViewerDeps) {
    this.worldId = deps.worldId;
    this.editable = deps.editable;
    this.viewNamespace = deps.viewNamespace;
    this.showGroundPlane = deps.showGroundPlane;
    this.showGimbal = deps.showGimbal;
    this.showFrames = deps.showFrames;
    this.pickingEnabled = deps.pickingEnabled;
    this.cameraControlsEnabled = deps.cameraControlsEnabled;
    this.camera = deps.camera;
    this.pollForWorldUpdates = deps.pollForWorldUpdates;
    this.worldWorkspaceService = deps.worldWorkspaceService;
    this.objectStateManager = deps.objectStateManager;
    this.pickerManager = deps.pickerManager;
    this.injector = deps.injector;
    this.container = deps.container;
    this.container.classList.add(CONTAINER_CLASS);

    this.userInputModeManager = new UserInputModeManager(this.injector);

    this.viewController = new ViewController({
      container: this.container,
      worldId: this.worldId,
      worldModel: this.worldModelSignal,
      polling: this.pollForWorldUpdates,
      viewNamespace: this.viewNamespace,
      showGimbal: this.showGimbal,
      pickingEnabled: this.scenePickingEnabled,
      cameraControlsEnabled: this.cameraControlsEnabled,
      camera: this.camera,
      assetDepot: deps.assetDepot,
      pickerManager: this.pickerManager,
      objectStateManager: this.objectStateManager,
      worldRendererService: deps.worldRendererService,
      showFps: signal(deps.options.showFps).asReadonly(),
      ngZone: deps.ngZone,
      injector: this.injector,
    });

    this.addEffect([this.worldId], async ([worldId]) => {
      const model =
        await this.worldWorkspaceService.getWorldModelAsync(worldId);
      this.worldModelSignal.set(model);
      this.editor = await this.worldWorkspaceService.getEditorAsync(worldId);
      this.checkAndEmitLoaded();
    });

    this.addEffect([this.objectStateManager.lockedObjectsOrFrames], () => {
      const unlockedPickedObjectIds = this.pickedObjectIds.filter(
        (ofid) => !this.isLocked(ofid),
      );
      // Only re-attach the transformer if a picked object got locked. The
      // locks are re-synced whenever the world changes, and re-attaching
      // would reset the current transform target (e.g. a robot being posed).
      if (unlockedPickedObjectIds.length === this.pickedObjectIds.length) {
        return;
      }
      this.pickedObjectIds = unlockedPickedObjectIds;
      this.attachTransformerToObjects(this.pickedObjectIds);
    });

    this.addEffect([this.editable], () => {
      this.updateEditable();
    });

    this.addEffect(
      [this.showFrames, this.sceneSignal],
      ([showFrames, scene]) => {
        scene.setFramesVisible(showFrames);
      },
    );

    this.addEffect(
      [this.showGroundPlane, this.sceneSignal],
      ([showGroundPlane, scene]) => {
        scene.setGridVisible(showGroundPlane);
      },
    );
  }

  /** Creates an effect that is destroyed when this viewer is disposed. */
  private addEffect<Values>(
    dependencies: Parameters<typeof explicitEffect<Values>>[0],
    fn: (values: Values, onCleanup: EffectCleanupRegisterFn) => void,
  ) {
    this.effects.push(
      explicitEffect(dependencies, fn, {injector: this.injector}),
    );
  }

  /**
   * Finishes setting up the viewer once its container is attached to the
   * document: registers the styles and the view controller, reacts to picking
   * and the transformer, and registers any deferred plugins.
   */
  init() {
    // The container is attached by now: register the styles of the view's
    // canvas and overlays with its document or shadow root.
    registerWorldViewerStyles(this.container);
    // The container is laid out by now too: size the view to it right away, so
    // that the first frame (and any framing of the scene) uses the right size.
    this.viewController.refreshSize();
    if (this.viewController.getNamespace()) {
      this.worldWorkspaceService.registerViewController(this.viewController);
    }

    this.subscriptions.add(
      this.viewController.nodeClicked.subscribe((pick: NodePickEvent) => {
        if (pick.event.shiftKey) return;
        const node = pick.node;
        if (node instanceof IntrEntityNode) {
          this.handleEntityPick(pick, node);
        } else if (node instanceof SceneHelper) {
          // Handled by plugins (e.g. RobotPosingPlugin) or scene helpers.
        } else {
          this.pickerManager.clearAllPicks();
        }
      }),
    );

    this.subscriptions.add(
      this.viewController.nodeHovered.subscribe((pick: NodePickEvent) => {
        const node = pick.node;
        if (node instanceof IntrEntityNode) {
          // We allow hovering objects or frames in the scene. Typically, hovering
          // an entity should trigger a hover of its parent object. Frames are an
          // exception. Hovering a frame should just hover the frame itself, since
          // we allow picking+moving frames without picking+moving their whole parent.
          const oid =
            node instanceof IntrFrameNode ? node.entityId : node.parentObjectId;
          // Do not allow hover-pick of locked objects
          if (!this.isLocked(oid)) {
            // Hover behavior: normally, hovering any node of an object
            // will hover-pick the whole object. But if we have an object
            // deep-picked, hovering an entity hover-picks just that entity.
            const deepPickedObjectId =
              this.pickerManager.getDeepPickedObjectId();
            if (deepPickedObjectId !== WorldModel.NULL_OBJECT_ID) {
              if (deepPickedObjectId === oid) {
                this.pickerManager.setHoverPickedEntityId(node.entityId);
              } else {
                this.pickerManager.clearHoverPickedEntity();
              }
            } else {
              this.pickerManager.setHoverPickedObjectId(oid);
            }
          }
        } else {
          this.pickerManager.clearHoverPicked();
        }
      }),
    );

    this.transformer = this.viewController.getTransformer();

    this.registerTransformerSubscriptions(this.transformer);

    this.subscriptions.add(
      this.pickerManager.deepPick.subscribe((deepPickedObjectId: number) => {
        const node = this.sceneSignal().getNodeByObjectId(deepPickedObjectId);
        if (!node) return;
        this.viewController.saveView('deep_pick');
        // We frame the deep-picked object every time a deep pick is made,
        // such as when opening the object instance editor.
        this.viewController.frameObjects([node]);
        this.viewController.requestNewFrame();
      }),
    );

    this.subscriptions.add(
      this.pickerManager.newPick.subscribe(
        (pickedObjectIds: ReadonlySet<number>) => {
          this.pickedObjectIds = [...pickedObjectIds];
          this.updateTransformer();
          this.sceneSignal().setHelpersVisible(true);
        },
      ),
    );

    this.subscriptions.add(
      merge(this.sceneSignal().updated, this.sceneSignal().mutated).subscribe(
        () => {
          this.viewController.requestNewFrame();
        },
      ),
    );

    // Make sure the frames' size scale as view zoom changes.
    this.subscriptions.add(
      this.viewController.cameraChanged.subscribe(() => {
        // Scale the frames based on distance so it does not become too small when
        // view is zoomed out. The scaling is square root instead of linear
        // to also avoid the frames become too dominant at a far zoom.
        const scale = Math.max(
          1.0,
          Math.sqrt(this.viewController.getDistance()),
        );
        const pageZoom = window.outerWidth / window.innerWidth;
        // Also scale by the page zoom that user set, which is an
        // accessibility feature user can take advantage off.
        this.sceneSignal().setFrameScale(scale * pageZoom);
      }),
    );

    // Call editable/active setters so that controllers are initially disabled.
    this.updateEditable();
    this.viewInitialized = true;
    this.checkAndEmitLoaded();

    // Register any plugins that were deferred prior to view initialization now
    // that the WorldViewer is in a known-good state.
    while (this.deferredPlugins.length > 0) {
      const plugin = this.deferredPlugins.shift();
      if (plugin) {
        this.registerPlugin(plugin);
      }
    }
  }

  /**
   * Tears down the viewer: unregisters its plugins and its view controller, and
   * disposes the view controller. Safe to call more than once.
   */
  dispose() {
    if (this.disposed) return;
    this.disposed = true;
    // Unregister the plugins in reverse order, in case
    // plugins registered later depend on earlier plugins.
    // Also: unregister plugins before other teardown here. Plugins
    // assume in their `onUnregister` that the WorldViewer
    // is still is a good state.
    this.deferredPlugins.length = 0;
    while (this.plugins.length > 0) {
      const plugin = this.plugins.pop();
      plugin?.onUnregister(this);
    }
    this.worldWorkspaceService.unregisterViewController(this.viewController);
    for (const effectRef of this.effects) {
      effectRef.destroy();
    }
    this.effects.length = 0;
    this.subscriptions.unsubscribe();
    this.userInputModeManager.dispose();
    this.viewController.dispose();
    this.container.classList.remove(CONTAINER_CLASS);
  }

  private checkAndEmitLoaded() {
    if (this.viewInitialized && this.worldModelSignal() && this.editor) {
      this.loaded.next(true);
    }
  }

  /**
   * Returns the world ID associated with this world.
   * TODO(b/505779015) - Users of WorldViewer should access the world ID using this
   * method, instead of accessing the worldId input signal. That input may be removed
   * in the future.
   */
  getWorldId(): string {
    return this.worldId();
  }

  getWorldModelSignal(): Signal<WorldModel | undefined> {
    return this.worldModelSignal.asReadonly();
  }

  /**
   * Returns the transformer used for rigid transform of objects in the scene,
   * or undefined if not available yet.
   */
  getTransformer(): Transformer | undefined {
    return this.transformer;
  }

  /** Returns the object state manager used by this viewer. */
  getObjectStateManager(): ObjectStateManager {
    return this.objectStateManager;
  }

  /** Returns the user input mode manager owned by this viewer. */
  getUserInputModeManager(): UserInputModeManager {
    return this.userInputModeManager;
  }

  /**
   * Registers a plugin with this WorldViewer. If called before the view is
   * initialized, registration is deferred until the view is fully initialized
   * to ensure the WorldViewer is in a known-good state before plugins access
   * and/or modify it.
   */
  registerPlugin(plugin: WorldViewerPlugin): void {
    if (
      this.plugins.includes(plugin) ||
      this.deferredPlugins.includes(plugin)
    ) {
      return;
    }
    if (!this.viewInitialized) {
      // Defer registration until the view is fully initialized so the
      // WorldViewer is in a known-good state before the plugin accesses or
      // modifies it.
      this.deferredPlugins.push(plugin);
      return;
    }
    this.plugins.push(plugin);
    plugin.onRegister(this);
  }

  /**
   * Unregisters a plugin from this WorldViewer.
   */
  unregisterPlugin(plugin: WorldViewerPlugin): void {
    const deferredIndex = this.deferredPlugins.indexOf(plugin);
    if (deferredIndex !== -1) {
      this.deferredPlugins.splice(deferredIndex, 1);
      return;
    }
    const index = this.plugins.indexOf(plugin);
    if (index === -1) {
      return;
    }
    this.plugins.splice(index, 1);
    plugin.onUnregister(this);
  }

  /**
   * Returns a list of currently registered plugins.
   */
  getPlugins(): readonly WorldViewerPlugin[] {
    return this.plugins;
  }

  private handleEntityPick(pick: NodePickEvent, node: IntrEntityNode) {
    const deepPickedObjectId = this.pickerManager.getDeepPickedObjectId();
    if (deepPickedObjectId !== WorldModel.NULL_OBJECT_ID) {
      // When we have an object deep-picked, clicking an entity of that object
      // will pick just that entity (instead of the whole object).
      if (deepPickedObjectId === node.parentObjectId) {
        this.pickerManager.setPickedEntityId(node.entityId);
      } else {
        this.pickerManager.clearPickedEntity();
      }
    } else {
      // We allow picking objects or frames in the scene. Typically, clicking/picking
      // an entity of an object should pick its parent object. FrameEntities are an
      // exception. We allow picking a frame independently of its parent object
      // so that it may be moved with the transformer without moving its whole
      // parent object.
      const oid =
        node instanceof IntrFrameNode ? node.entityId : node.parentObjectId;
      // Do not allow pick of locked objects
      if (this.isLocked(oid)) return;
      // Add or remove object from picked set when Ctrl or Meta are pressed.
      if (
        pick.event.getModifierState('Control') ||
        pick.event.getModifierState('Meta')
      ) {
        if (this.pickedObjectIds.includes(oid)) {
          this.pickerManager.removePickedObjectId(oid);
        } else {
          this.pickerManager.addPickedObjectId(oid);
        }
      } else {
        this.pickerManager.setPickedObjectId(oid);
      }
    }
  }

  private registerTransformerSubscriptions(transformer: Transformer) {
    if (!transformer) {
      return;
    }
    this.subscriptions.add(
      transformer.transformerTouched.subscribe((touched: boolean) => {
        // Disable mouse picking when transformer is hovered over. Stopping to
        // drag the transformer by releasing the mouse button always sends a
        // touched=false message and hence will enable the picker.
        this.pickerManager.setEnabled(!touched);
      }),
    );

    this.subscriptions.add(
      transformer.transformStarted.subscribe(() => {
        // Disable picking and camera control when transformer is dragged.
        // Disabling picking on transformer touched above is not enough.
        // In some cases the user might touch the transfomer axis,
        // complete a drag and re-start a new drag without moving the
        // mouse away from the axis. This prevents the 'transformerTouch'
        // subscriber from being called again to disable the picker.
        this.pickerManager.setEnabled(false);
        this.viewController.setFrozen(true);
      }),
    );

    this.subscriptions.add(
      transformer.transformStopped.subscribe(() => {
        // Resume camera control.
        this.viewController.setFrozen(false);
        // Enable mouse-picker at the end of a drag event.
        this.pickerManager.setEnabled(true);
      }),
    );

    this.subscriptions.add(
      transformer.enabledChanged.subscribe((enabled: boolean) => {
        // The transformer detaches when it gets disabled. Re-attach it to the
        // picked objects when it gets enabled again, e.g. after editing is
        // re-enabled or a user input mode (like the measurement tool) ends.
        if (enabled) {
          this.updateTransformer();
        }
      }),
    );
  }

  // Get the current view of the camera
  get cameraView(): CameraView {
    return {
      pos: this.viewController.camera().position.clone(),
      target: this.viewController.getTarget(),
    };
  }

  private updateEditable() {
    if (!this.sceneSignal().root) return;
    this.transformer?.setEnabled(this.editable());
  }

  // attach the transformer to the TCP of the selected robot if we're on the robot.
  private updateTransformer() {
    this.attachTransformerToObjects(this.pickedObjectIds);
  }

  /**
   * Attaches the transformer to a `RigidObjectTarget` that moves the given
   * objects together, with the controls on the last of them. Locked objects,
   * attachment points and the root object are skipped. Detaches the
   * transformer if no object is left.
   */
  attachTransformerToObjects(objectIds: number[]) {
    const model = this.worldModelSignal();
    const transformer = this.transformer;
    if (!model || !transformer) return;
    const filteredObjectIds = objectIds.filter(
      (ofid: number) =>
        !this.isLocked(ofid) &&
        !model.hasAttachmentPoint(ofid) &&
        ofid !== WorldModel.NULL_OBJECT_ID &&
        ofid !== WorldModel.ROOT_OBJECT_ID,
    );
    const target =
      filteredObjectIds.length > 0
        ? RigidObjectTarget.create(this.sceneSignal(), {
            objectIds: filteredObjectIds,
            activeObjectId: filteredObjectIds[filteredObjectIds.length - 1],
            editor: this.editor,
          })
        : undefined;
    if (target) {
      transformer.attach(target);
    } else {
      transformer.detach();
    }
  }

  /**
   * Snaps the object that the transformer is attached to onto the closest
   * surface below it. No-op if the transformer is not attached to objects
   * (e.g. while posing a robot).
   */
  snapSelectedObjectToSurfaceBelow() {
    const target = this.transformer?.getTransformTarget();
    if (target instanceof RigidObjectTarget) {
      target.snapToSurfaceBelow();
    }
  }

  private isLocked(id: number): boolean {
    const ref = this.worldModelSignal()?.getReferenceString(id);
    if (!ref) return false;
    return this.objectStateManager.lockedObjectsOrFrames().has(ref);
  }
}
