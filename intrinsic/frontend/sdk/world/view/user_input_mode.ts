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

import {Observable, Subject} from 'rxjs';

/**
 * Abstract class for custom user input modes. Clients can extend this class and
 * implement the required functionality.
 */
export abstract class UserInputMode {
  /**
   * @param modeName A unique name identifying this input mode. Each mode is
   *     expected to export its own name constant (for example
   *     `FRAME_EDITOR_MODE_NAME`) and pass it here. Names must not overlap with
   *     those of any other input mode.
   */
  constructor(readonly modeName: string) {}
  protected readonly exitingModeSubject = new Subject<void>();
  readonly exitingMode: Observable<void> = this.exitingModeSubject;
  mouseMoveHandler?(e: MouseEvent): void;
  mouseClickHandler?(e: MouseEvent): void;
  keyDownHandler?(e: KeyboardEvent): void;
  keyUpHandler?(e: KeyboardEvent): void;
  /**
   * Canvas used to register and unregister mouse events. Some modes need to
   * restrict functionality to the 3d world canvas. If not specified, 'document'
   * is used instead.
   */
  getCanvas?(): HTMLCanvasElement | undefined;

  /**
   * When this mode is active, prevent any other mode from becoming active.
   * Useful for tool modes that need to 'own' user input until their operations
   * are completed
   */
  preventOverride?: boolean;

  /**
   * When this mode is active, the WorldViewer disables click and hover picking
   * of scene nodes. Useful for tool modes that handle pointer input themselves
   * and would otherwise conflict with the object picker.
   */
  disablesScenePicking?: boolean;

  /**
   * Called when the user explicitly enters a mode. Not called if this mode
   * becomes active due to another mode being popped off the stack.
   */
  onEnter?(userData?: unknown): void;

  /**
   * Called when the user explicitly exits a mode. Not called if this mode
   * becomes inactive due to another mode being pushed on the stack.
   */
  onExit?(): void;

  /**
   * Called whenever the mode moves to the top of the stack.
   */
  onActivate?(userData?: unknown): void;

  /**
   * Called whenever the mode leaves the top of the stack
   */
  onDeactivate?(): void;

  /**
   * Derived modes call this method when they intent to exit. This lets the
   * manager know that this mode is exiting, and that it should be disabled and
   * popped from the stack.
   */
  protected requestExit() {
    this.exitingModeSubject.next();
  }
}
