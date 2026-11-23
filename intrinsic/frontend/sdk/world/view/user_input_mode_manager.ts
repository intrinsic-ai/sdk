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

import {Injector, computed, runInInjectionContext, signal} from '@angular/core';
import {Subscription} from 'rxjs';

import {UserInputMode} from './user_input_mode';

// Private interface to track modes and their assigned identifier.
interface ModeInfo {
  modeName: string;
  mode: UserInputMode;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  mouseMoveBind?: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  mouseClickBind?: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  keyDownHandlerBind?: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  keyUpHandlerBind?: any;
}

/**
 * User-input mode manager. Responsible for managing a stack of objects derived
 * from UserInputMode. Exposes API to register new modes as well as pushing new
 * modes on the stack. While the manager provides an 'exitMode' call, generally
 * a mode itself should determine when it is done and request to exit via its
 * requestExit() call. The manager is notified via activeModeExitSubscription.
 *
 * Each WorldViewer owns one manager, provided at the component level. Get it
 * with `WorldViewer.getUserInputModeManager()` rather than injecting it.
 *
 * The manager differentiates between 'entering' a mode and 'activating' a mode.
 * This allows more fine-grained control of initialization and internal
 * switching within the modes. The table below provides a visual explanation of
 * triggering event and resulting operations.
 *
 * Entering a mode means 'enterMode' was called and a mode will be pushed on the
 * stack. onDeactivate() will be called for the previously active mode.
 * onEnter() and onActivate() will be called for the mode we are entering.
 *
 * Conversely exiting a mode means a mode's requestExit() was invoked. The modes
 * onDeactivate() and onExit() are called and the mode is popped from the top of
 * the stack. onActivate() is called for the mode that just bubbled to the top
 * of the stack.
 */
export class UserInputModeManager {
  private readonly registeredModes = new Map<string, UserInputMode>();
  private readonly modeStack: ModeInfo[] = [];
  private activeModeExitSubscription?: Subscription;

  private readonly activeModeSignal = signal<UserInputMode | undefined>(
    undefined,
  );
  /**
   * The mode currently at the top of the stack, or undefined if the stack is
   * empty. Updated whenever a mode is entered or exited.
   */
  readonly activeMode = this.activeModeSignal.asReadonly();
  /**
   * The name of the mode currently at the top of the stack, or undefined if
   * the stack is empty. Updated whenever a mode is entered or exited.
   */
  readonly activeModeName = computed<string | undefined>(
    () => this.activeMode()?.modeName,
  );

  /**
   * @param injector Injector that mode factories passed to `enterMode` run in,
   *     so that they can `inject` their dependencies.
   */
  constructor(private readonly injector: Injector) {}

  // Clear the internal stack, clear subscriptions and exit.
  dispose() {
    // Unwind from the top of the stack downwards: exitMode() requires the mode
    // being exited to be the currently active one, and pops it off. Walking the
    // index down keeps this bounded even as the stack shrinks.
    for (let i = this.modeStack.length - 1; i >= 0; i--) {
      this.exitMode(this.modeStack[i].modeName);
    }

    this.activeModeExitSubscription?.unsubscribe();
    this.activeModeExitSubscription = undefined;

    this.registeredModes.clear();
  }

  /**
   * Register a new mode with the manager. The mode's name is used to identify
   * the mode for subsequent operations, and must be unique: registering a mode
   * whose name is already on the mode stack is an error.
   *
   * Re-registering a name that is not currently on the stack is allowed, so
   * that a tool can be re-entered with a fresh instance after it has been
   * exited.
   */
  private registerMode(mode: UserInputMode) {
    if (
      this.modeStack.some((modeInfo) => modeInfo.modeName === mode.modeName)
    ) {
      throw new Error(
        `Cannot register mode '${mode.modeName}'. A mode with that name is ` +
          `already on the mode stack. User input mode names must be unique.`,
      );
    }

    this.registeredModes.set(mode.modeName, mode);
  }

  /**
   * Returns true if the stack is empty.
   */
  isEmpty(): boolean {
    return this.modeStack.length === 0;
  }

  /**
   * Get size of stack to know how many modes have been pushed.
   */
  getModeStackSize(): number {
    return this.modeStack.length;
  }

  /**
   * Enters and activates a new mode. The specified userData is passed to the
   * mode's onEnter() function, if one is registered.
   *
   * `modeOrFactory` is either a mode instance, or a function that creates one.
   * A factory is run in the manager's injection context, so modes that call
   * inject() in their constructor can be created from outside one. The factory
   * is not called if the current mode prevents being overridden.
   */
  enterMode(
    modeOrFactory: UserInputMode | (() => UserInputMode),
    userData?: unknown,
  ) {
    // Check whether the new mode is allowed onto the stack. This runs before
    // registering it so that a mode we refuse to enter does not displace the
    // currently registered instance.
    const currentModeData = this.getTop();
    // Don't allow to override, unless the current active mode allows it.
    if (currentModeData?.mode.preventOverride) {
      return;
    }
    const mode =
      modeOrFactory instanceof UserInputMode
        ? modeOrFactory
        : runInInjectionContext(this.injector, modeOrFactory);
    // If we are already the active mode, return.
    if (currentModeData?.modeName === mode.modeName) {
      return;
    }

    this.registerMode(mode);

    if (currentModeData) {
      // Deactivate the current mode.
      if (currentModeData.mode.onDeactivate) {
        currentModeData.mode.onDeactivate();
      }
      // Unregister handlers for the current mode.
      this.unregisterModeHandlers(currentModeData);
    }
    // Activate the new mode and push it on the internal stack
    const modeInfo = this.registerModeHandlers(mode.modeName, mode);
    this.modeStack.push(modeInfo);
    this.activeModeSignal.set(mode);
    if (mode.onEnter) {
      mode.onEnter(userData);
    }
    if (mode.onActivate) {
      mode.onActivate(userData);
    }
  }

  /**
   * Called when the active mode needs to exit. This method is only called when
   * the mode is being popped off the top of the stack. It is not called when
   * the active mode is prehempted by another mode being pushed on the stack.
   */
  exitMode(modeName: string) {
    const mode = this.registeredModes.get(modeName);
    if (!mode) {
      throw new Error(
        `Cannot exit mode. Mode ${modeName} has not been registered`,
      );
    }
    if (modeName !== this.activeModeName()) {
      throw new Error(
        `Cannot exit mode. Mode ${modeName} is not active (${this.activeModeName()})`,
      );
    }
    // Unregister handlers for the current mode, call onDeactivate and onExit.
    const currentModeData = this.getTop();
    if (!currentModeData) return;

    this.unregisterModeHandlers(currentModeData);
    if (mode.onDeactivate) {
      mode.onDeactivate();
    }
    if (mode.onExit) {
      mode.onExit();
    }
    this.modeStack.pop();
    this.activeModeSignal.set(this.getTop()?.mode);

    // Re-activate the mode that is now on top of the stack.
    // We do not call enterMode here because we do not want to re-initialize.
    const newMode = this.getTop();
    if (newMode) {
      this.registerModeHandlers(newMode.modeName, newMode.mode);
      if (newMode.mode.onActivate) {
        newMode.mode.onActivate();
      }
    }
  }

  /**
   * Returns the mode currently at the top of the stack. Returns undefined if
   * the stack is empty.
   */
  private getTop(): ModeInfo | undefined {
    return this.modeStack[this.modeStack.length - 1];
  }

  /**
   * Internal call to manager subscribe to the active mode and register the
   * required input handlers. This method is called whenever a mode becomes
   * active.
   */
  private registerModeHandlers(
    modeName: string,
    mode: UserInputMode,
  ): ModeInfo {
    // Subscribe to the mode exiting
    this.activeModeExitSubscription = mode.exitingMode.subscribe(() => {
      this.exitMode(modeName);
    });

    const modeInfo: ModeInfo = {modeName, mode};
    const canvas = modeInfo.mode.getCanvas?.();
    if (mode.mouseMoveHandler) {
      modeInfo.mouseMoveBind = (event: MouseEvent) => {
        if (mode.mouseMoveHandler) {
          mode.mouseMoveHandler(event);
        }
      };
      if (canvas) {
        canvas.addEventListener('pointermove', modeInfo.mouseMoveBind);
      } else {
        document.addEventListener('pointermove', modeInfo.mouseMoveBind);
      }
    }
    if (mode.mouseClickHandler) {
      modeInfo.mouseClickBind = (event: MouseEvent) => {
        if (mode.mouseClickHandler) {
          mode.mouseClickHandler(event);
        }
      };
      if (canvas) {
        canvas.addEventListener('pointerup', modeInfo.mouseClickBind);
      } else {
        document.addEventListener('pointerup', modeInfo.mouseClickBind);
      }
    }
    if (mode.keyDownHandler) {
      modeInfo.keyDownHandlerBind = (event: KeyboardEvent) => {
        if (mode.keyDownHandler) {
          mode.keyDownHandler(event);
        }
      };
      if (canvas) {
        canvas.addEventListener('keydown', modeInfo.keyDownHandlerBind);
      } else {
        document.addEventListener('keydown', modeInfo.keyDownHandlerBind);
      }
    }
    if (mode.keyUpHandler) {
      modeInfo.keyUpHandlerBind = (event: KeyboardEvent) => {
        if (mode.keyUpHandler) {
          mode.keyUpHandler(event);
        }
      };
      if (canvas) {
        canvas.addEventListener('keyup', modeInfo.keyUpHandlerBind);
      } else {
        document.addEventListener('keyup', modeInfo.keyUpHandlerBind);
      }
    }
    return modeInfo;
  }

  /**
   * Internal call to manager unsubscribe from the active mode and unregister
   * the required input handlers. This method is called whenever a mode becomes
   * inactive.
   */
  private unregisterModeHandlers(modeInfo: ModeInfo) {
    const canvas = modeInfo.mode.getCanvas?.();

    // Clear current subscription, if any
    this.activeModeExitSubscription?.unsubscribe();

    if (modeInfo.mouseMoveBind) {
      if (canvas) {
        canvas.removeEventListener('pointermove', modeInfo.mouseMoveBind);
      } else {
        document.removeEventListener('pointermove', modeInfo.mouseMoveBind);
      }
    }
    if (modeInfo.mouseClickBind) {
      if (canvas) {
        canvas.removeEventListener('pointerup', modeInfo.mouseClickBind);
      } else {
        document.removeEventListener('pointerup', modeInfo.mouseClickBind);
      }
    }
    if (modeInfo.keyDownHandlerBind) {
      if (canvas) {
        canvas.removeEventListener('keydown', modeInfo.keyDownHandlerBind);
      } else {
        document.removeEventListener('keydown', modeInfo.keyDownHandlerBind);
      }
    }
    if (modeInfo.keyUpHandlerBind) {
      if (canvas) {
        canvas.removeEventListener('keyup', modeInfo.keyUpHandlerBind);
      } else {
        document.removeEventListener('keyup', modeInfo.keyUpHandlerBind);
      }
    }
  }
}
