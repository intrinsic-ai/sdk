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

import * as three from 'three';

/**
 * Represents something that a `Transformer` can be attached to.
 *
 * The `Transformer` owns the transform controls (the draggable gizmo) and
 * forwards user interaction to its current target. The target decides what the
 * controls attach to and what to do with the resulting transform, e.g. move a
 * set of objects or pose a robot.
 */
export interface TransformTarget {
  /** The object that the transform controls attach to. */
  readonly handle: three.Object3D;

  /** Whether the controls should snap the handle onto surfaces when dragged. */
  readonly enableSurfacePlacement: boolean;

  /**
   * Called whenever the user drags the controls.
   *
   * @param worldTHandle The new world transform of `handle`.
   */
  onDrag?(worldTHandle: three.Matrix4): void;

  /**
   * Called when the user stopped dragging the controls.
   *
   * @param changed Whether the user actually moved the handle during the drag
   *     (as opposed to just clicking on the controls).
   */
  onDragEnd?(changed: boolean): void;

  /**
   * Called when the world shown in the scene changed (e.g. an object was moved
   * by a command, undo/redo, or snapping), unless the controls are being
   * dragged. Targets whose handle should mirror the world can re-sync here.
   */
  onWorldUpdated?(): void;

  /** Surfaces that the handle can be snapped onto. */
  snappingSurfaces?(): Iterable<three.Object3D>;

  /**
   * Called when the transformer detaches from this target. Removes anything the
   * target added to the scene. The target must not be used afterwards.
   */
  dispose(): void;
}
