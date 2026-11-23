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

import {
  EnvironmentInjector,
  Signal,
  WritableSignal,
  effect,
  signal,
  untracked,
} from '@angular/core';

import {NoopWorldEditor} from '../editor/noop_world_editor';
import {WorldEditor} from '../editor/world_editor';
import {waitForWorldModelToLoad} from '../model/wait_for_world_model';
import {WorldModel} from '../model/world_model';
import {ViewController} from './view_controller';

/**
 * Represents an active editing session for a specific world.
 *
 * Includes the WorldModel, WorldEditor, and any other essential
 * classes for editing a world.
 */
export class WorldWorkspace {
  constructor(
    readonly model: WorldModel,
    readonly editor: WorldEditor,
  ) {}
}

/**
 * Service to manage the lifecycle of WorldWorkspaces, which are like
 * editing sessions for worlds.
 *
 * Workspaces are created up front by the application (see
 * `createWorkspace`) and then looked up by world id. Lookups come
 * in two flavors:
 *  - Synchronous (`getWorkspace`, `getWorldModel`, `getEditor`): return
 *    immediately, but the WorldModel may not have finished loading yet.
 *  - Asynchronous (`getWorkspaceAsync`, `getWorldModelAsync`,
 *    `getEditorAsync`): resolve once the WorldModel has loaded.
 * Both throw (or reject) if no workspace exists for the given world id.
 *
 * Code that runs before a workspace exists (e.g. root services) can use the
 * reactive `workspace` lookup instead, which returns a signal of the workspace
 * that's undefined until the workspace is created.
 *
 * This service also tracks the ViewControllers of the app, keyed by view
 * namespace.
 */
export class WorldWorkspaceService {
  // Maps world id to a signal of that world's workspace, which is undefined
  // while there's none.
  private readonly workspaces = new Map<
    string,
    WritableSignal<WorldWorkspace | undefined>
  >();
  // Maps a WorldModel to its load, shared so that concurrent async lookups
  // wait on (and retry) a single load.
  private readonly modelLoads = new WeakMap<WorldModel, Promise<void>>();
  // Maps view namespace to a signal of the ViewController registered for it,
  // which is undefined while there's none.
  private readonly viewControllers = new Map<
    string,
    WritableSignal<ViewController | undefined>
  >();

  /**
   * @param environmentInjector Injector used to create the effects that wait
   *     for view controllers to be registered.
   */
  constructor(private readonly environmentInjector: EnvironmentInjector) {}

  /**
   * Creates and caches a WorldWorkspace for an already-made WorldModel and the
   * WorldEditor that mutates it. Without an editor, the workspace gets a
   * NoopWorldEditor, which ignores all mutations (i.e., a read-only world).
   * The workspace is available to lookups immediately, even if the model
   * hasn't finished loading yet.
   *
   * Returns the existing workspace if one was already created for this model
   * (ignoring `editor`). Otherwise replaces any existing workspace for the
   * model's world id.
   */
  createWorkspace(
    worldModel: WorldModel,
    editor: WorldEditor = new NoopWorldEditor(worldModel),
  ): WorldWorkspace {
    const existing = this.lookupWorkspace(worldModel.id);
    if (existing?.model === worldModel) {
      return existing;
    }
    const workspace = new WorldWorkspace(worldModel, editor);
    this.workspaceSignal(worldModel.id).set(workspace);
    return workspace;
  }

  /**
   * Clears the cached WorldWorkspace for the given WorldModel. Does nothing if
   * the workspace for the model's world id has since been replaced with one
   * for a different model.
   */
  clearWorkspace(worldModel: WorldModel) {
    if (this.lookupWorkspace(worldModel.id)?.model !== worldModel) {
      return;
    }
    this.workspaceSignal(worldModel.id).set(undefined);
  }

  /** Returns whether a workspace exists for the given world id. */
  hasWorkspace(worldId: string): boolean {
    return this.lookupWorkspace(worldId) !== undefined;
  }

  /**
   * Returns a signal of the workspace for the given world id. It's undefined
   * until the workspace is created, updates when it's replaced, and goes back
   * to undefined when it's cleared.
   *
   * The other lookups on this service aren't reactive, so prefer this one when
   * the workspace might not exist yet.
   */
  workspace(worldId: string): Signal<WorldWorkspace | undefined> {
    return this.workspaceSignal(worldId).asReadonly();
  }

  /**
   * Returns the workspace for the given world id. The workspace's model may
   * not have finished loading yet; use `getWorkspaceAsync` if that matters.
   *
   * @throws if no workspace exists for the given world id.
   */
  getWorkspace(worldId: string): WorldWorkspace {
    const workspace = this.lookupWorkspace(worldId);
    if (!workspace) {
      throw new Error(`No WorldWorkspace exists for world '${worldId}'.`);
    }
    return workspace;
  }

  /** Synchronous shortcut for `getWorkspace(worldId).model`. */
  getWorldModel(worldId: string): WorldModel {
    return this.getWorkspace(worldId).model;
  }

  /** Synchronous shortcut for `getWorkspace(worldId).editor`. */
  getEditor(worldId: string): WorldEditor {
    return this.getWorkspace(worldId).editor;
  }

  /**
   * Returns the workspace for the given world id once its model has loaded.
   *
   * Rejects if no workspace exists for the given world id.
   */
  async getWorkspaceAsync(worldId: string): Promise<WorldWorkspace> {
    const workspace = this.getWorkspace(worldId);
    await this.waitForModelToLoad(workspace.model);
    return workspace;
  }

  /** Shortcut for `getWorkspaceAsync(worldId)` returning just the model. */
  async getWorldModelAsync(worldId: string): Promise<WorldModel> {
    return (await this.getWorkspaceAsync(worldId)).model;
  }

  /** Shortcut for `getWorkspaceAsync(worldId)` returning just the editor. */
  async getEditorAsync(worldId: string): Promise<WorldEditor> {
    return (await this.getWorkspaceAsync(worldId)).editor;
  }

  /**
   * Resets the cached workspaces.
   */
  reset() {
    for (const workspace of this.workspaces.values()) {
      workspace.set(undefined);
    }
  }

  registerViewController(viewController: ViewController) {
    this.viewControllerSignal(viewController.getNamespace()).set(
      viewController,
    );
  }

  /**
   * Unregisters the view controller. Does nothing if another view controller
   * has since been registered for its view namespace.
   */
  unregisterViewController(viewController: ViewController) {
    const registered = this.viewControllerSignal(viewController.getNamespace());
    if (untracked(registered) === viewController) {
      registered.set(undefined);
    }
  }

  createViewControllerSignal(
    viewNamespace: string,
  ): Signal<ViewController | undefined> {
    return this.viewControllerSignal(viewNamespace).asReadonly();
  }

  /** Returns a promise that resolves once the view controller is registered. */
  async getViewController(viewNamespace: string): Promise<ViewController> {
    const registered = this.viewControllerSignal(viewNamespace);
    const viewController = untracked(registered);
    if (viewController) return viewController;
    return new Promise((resolve) => {
      const effectRef = effect(
        () => {
          const viewController = registered();
          if (!viewController) return;
          effectRef.destroy();
          resolve(viewController);
        },
        {injector: this.environmentInjector},
      );
    });
  }

  /** Gets or creates the signal for the given world id's workspace. */
  private workspaceSignal(
    worldId: string,
  ): WritableSignal<WorldWorkspace | undefined> {
    let workspace = this.workspaces.get(worldId);
    if (!workspace) {
      workspace = signal<WorldWorkspace | undefined>(undefined);
      this.workspaces.set(worldId, workspace);
    }
    return workspace;
  }

  /** Non-reactive lookup of the workspace for the given world id. */
  private lookupWorkspace(worldId: string): WorldWorkspace | undefined {
    const workspace = this.workspaces.get(worldId);
    if (!workspace) return undefined;
    return untracked(workspace);
  }

  private waitForModelToLoad(model: WorldModel): Promise<void> {
    let load = this.modelLoads.get(model);
    if (!load) {
      load = waitForWorldModelToLoad(model);
      this.modelLoads.set(model, load);
    }
    return load;
  }

  /** Gets or creates the signal for the given view namespace's controller. */
  private viewControllerSignal(
    viewNamespace: string,
  ): WritableSignal<ViewController | undefined> {
    let viewController = this.viewControllers.get(viewNamespace);
    if (!viewController) {
      viewController = signal<ViewController | undefined>(undefined);
      this.viewControllers.set(viewNamespace, viewController);
    }
    return viewController;
  }
}
