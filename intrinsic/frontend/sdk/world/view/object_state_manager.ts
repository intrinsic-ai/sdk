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

import {Signal, computed, signal} from '@angular/core';

/**
 * Holds per-object view and edit flags of the world viewer: which objects or
 * frames are hidden or locked, which entities show their collision margin and
 * whether collision geometry is visible (per world).
 *
 * This is the source of truth for these flags. Objects and frames are
 * identified by their reference string (see `WorldModel.getReferenceString()`),
 * entities by their entity ID.
 *
 * Hosts of the world viewer must provide it, typically as one shared instance
 * at the app (environment) level, so that the viewers, the `WorldEditor`
 * (which is created in the environment injector) and the trees and panels
 * share the same flags.
 */
export class ObjectStateManager {
  private readonly hiddenObjectsOrFramesSignal = signal<ReadonlySet<string>>(
    new Set(),
  );
  private readonly lockedObjectsOrFramesSignal = signal<ReadonlySet<string>>(
    new Set(),
  );
  private readonly entitiesShowingMarginSignal = signal<ReadonlySet<number>>(
    new Set(),
  );
  private readonly collisionGeometryVisibilitySignal = signal<
    ReadonlyMap<string, boolean>
  >(new Map());

  /** The reference strings of the objects and frames that are hidden. */
  readonly hiddenObjectsOrFrames: Signal<ReadonlySet<string>> =
    this.hiddenObjectsOrFramesSignal.asReadonly();

  /** The reference strings of the objects and frames that are locked. */
  readonly lockedObjectsOrFrames: Signal<ReadonlySet<string>> =
    this.lockedObjectsOrFramesSignal.asReadonly();

  /** The IDs of the entities whose collision margin is shown. */
  readonly entitiesShowingMargin: Signal<ReadonlySet<number>> =
    this.entitiesShowingMarginSignal.asReadonly();

  /**
   * Whether collision geometry is visible, by world ID. Worlds without an
   * entry do not show collision geometry.
   */
  readonly collisionGeometryVisibility: Signal<ReadonlyMap<string, boolean>> =
    this.collisionGeometryVisibilitySignal.asReadonly();

  /** Hides the object or frame `ref`. No-op if `ref` is empty. */
  hideObjectOrFrame(ref: string) {
    if (!ref) return;
    this.hiddenObjectsOrFramesSignal.update((hidden) =>
      new Set(hidden).add(ref),
    );
  }

  /** Shows the object or frame `ref`. */
  showObjectOrFrame(ref: string) {
    this.hiddenObjectsOrFramesSignal.update((hidden) => {
      const newHidden = new Set(hidden);
      newHidden.delete(ref);
      return newHidden;
    });
  }

  /** Shows or hides the object or frame `ref`. No-op if `ref` is empty. */
  setObjectOrFrameVisibility(ref: string, visible: boolean) {
    if (!ref) return;
    if (visible) {
      this.showObjectOrFrame(ref);
    } else {
      this.hideObjectOrFrame(ref);
    }
  }

  /** Replaces the set of hidden objects and frames. */
  setHiddenObjectsOrFrames(hidden: Iterable<string>) {
    this.hiddenObjectsOrFramesSignal.set(new Set(hidden));
  }

  /** Locks the object or frame `ref`. */
  lockObjectOrFrame(ref: string) {
    this.lockedObjectsOrFramesSignal.update((locked) =>
      new Set(locked).add(ref),
    );
  }

  /** Unlocks the object or frame `ref`. */
  unlockObjectOrFrame(ref: string) {
    this.lockedObjectsOrFramesSignal.update((locked) => {
      const newLocked = new Set(locked);
      newLocked.delete(ref);
      return newLocked;
    });
  }

  /** Replaces the set of locked objects and frames. */
  setLockedObjectsOrFrames(locked: Iterable<string>) {
    this.lockedObjectsOrFramesSignal.set(new Set(locked));
  }

  /** Shows the collision margin of the entity `eid`. */
  showEntityMargin(eid: number) {
    this.entitiesShowingMarginSignal.update((shown) => new Set(shown).add(eid));
  }

  /** Hides the collision margin of the entity `eid`. */
  hideEntityMargin(eid: number) {
    this.entitiesShowingMarginSignal.update((shown) => {
      const newShown = new Set(shown);
      newShown.delete(eid);
      return newShown;
    });
  }

  /** Replaces the set of entities whose collision margin is shown. */
  setEntitiesShowingMargin(shown: Iterable<number>) {
    this.entitiesShowingMarginSignal.set(new Set(shown));
  }

  /** Whether collision geometry is visible in the world `worldId`. */
  isCollisionGeometryVisible(worldId: string): boolean {
    return !!this.collisionGeometryVisibility().get(worldId);
  }

  /**
   * Returns a signal of whether collision geometry is visible in the world
   * `worldId`.
   */
  collisionGeometryVisible(worldId: string): Signal<boolean> {
    return computed(() => this.isCollisionGeometryVisible(worldId));
  }

  /** Shows or hides collision geometry in the world `worldId`. */
  setCollisionGeometryVisibility(worldId: string, visible: boolean) {
    this.collisionGeometryVisibilitySignal.update((visibility) =>
      new Map(visibility).set(worldId, visible),
    );
  }

  /** Toggles the visibility of collision geometry in the world `worldId`. */
  toggleCollisionGeometryVisibility(worldId: string) {
    this.setCollisionGeometryVisibility(
      worldId,
      !this.isCollisionGeometryVisible(worldId),
    );
  }
}
