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

import {signal} from '@angular/core';

/**
 * The callback interface allows any class that needs to update its internal
 * state inside the WorldView::mainLoop to register a handler to its update
 * function. For example, an animation class.
 *
 * @param timeMs The current time.
 * @param deltaMs The delta time since last update.
 */
export interface FrameUpdateCallback {
  (timeMs: number, deltaMs: number): void;
}

/**
 * The mainloop controls the update cyle, but not every update cycle draws a new
 * frame. A new frame is only drawn when requested. A registered
 * FrameUpdateCallback can be invoked either every update or only before a new
 * frame is drawn.
 */
export enum FrameUpdateCallbackType {
  // Always invoke the update callback in the main loop, prior to render.
  ALWAYS = 0,
  // Only invoke the update when a new frame is to be rendered.
  NEW_FRAME_ONLY = 1,
}

/**
 * A FrameTracker object shall be owned by the WorldView and passed to any
 * object that updates the view to request new frames to be drawn. It can also
 * keep central record of the number of frames drawn.
 */
export class FrameTracker {
  // We expose the frameRate as a Signal to allow components to track it as
  // it updates.
  private readonly frameRateSignalInternal = signal<number>(0);
  readonly frameRateSignal = this.frameRateSignalInternal.asReadonly();

  // Frame counter for computing FPS on the fly; resets every second.
  private frameCountInternal = 0;
  private wantNewFrame = false;
  private lastFrameUpdateTimeMs = 0;
  // Last time the FPS is computed and frame count resets.
  private lastFrameCountResetTimeMs = 0;
  private nextCallbackId = 1;
  private readonly frameUpdateCallbacks = new Map<
    number,
    {update: FrameUpdateCallback; type: FrameUpdateCallbackType}
  >();
  constructor() {
    this.lastFrameUpdateTimeMs = this.lastFrameCountResetTimeMs =
      performance.now();
  }

  /**
   * Explicitly request a new frame to be rendered.
   */
  requestNewFrame() {
    this.wantNewFrame = true;
  }

  /**
   * Register a frame update callback.
   * @param update Function to be invoked during frame update. Example:
   *     register callback for animation update. Callbacks are invoked in the
   *     order they are registered.
   * @param type The type of callback invocation.
   * @returns the unique callback id which needs to be used to unregister.
   */
  registerFrameUpdateCallback(
    update: FrameUpdateCallback,
    type: FrameUpdateCallbackType,
  ) {
    const callbackId = this.nextCallbackId;
    this.frameUpdateCallbacks.set(callbackId, {update, type});
    this.nextCallbackId++;
    return callbackId;
  }

  unregisterFrameUpdateCallback(id: number) {
    if (this.frameUpdateCallbacks.has(id)) {
      this.frameUpdateCallbacks.delete(id);
    } else {
      console.warn(`Callback ${id} not found!`);
    }
  }

  /**
   * Invoke registered frame update callbacks. If a new frame update is
   * requested and advance the frame count and return true.
   *
   * Note: This function is only to be called within WorldView's mainLoop
   * function.
   *
   * @returns Returns true if a new frame need to be drawn. WorldView will use
   *     this return value to decide if render a new frame or not.
   */
  updateFrame() {
    const timeMs = performance.now();
    const deltaMs = timeMs - this.lastFrameUpdateTimeMs;
    if (timeMs >= this.lastFrameCountResetTimeMs + 1000) {
      this.frameRateSignalInternal.set(
        Math.floor(
          (this.frameCountInternal * 1000) /
            (timeMs - this.lastFrameCountResetTimeMs),
        ),
      );
      this.lastFrameCountResetTimeMs = timeMs;
      this.frameCountInternal = 0;
    }
    // Increment frame count in the tracker.
    this.frameCountInternal++;

    this.lastFrameUpdateTimeMs = timeMs;
    for (const callback of this.frameUpdateCallbacks.values()) {
      if (
        callback.type === FrameUpdateCallbackType.ALWAYS ||
        this.wantNewFrame
      ) {
        callback.update(timeMs, deltaMs);
      }
    }
    // Return false if no new frame is requested.
    if (!this.wantNewFrame) {
      return false;
    }

    // Reset this flag for next frame.
    this.wantNewFrame = false;

    return true;
  }
}
