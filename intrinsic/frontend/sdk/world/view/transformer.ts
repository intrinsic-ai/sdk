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

import {Observable, Subject, Subscription} from 'rxjs';
import * as three from 'three';

import {FrameUpdateCallbackType} from './frame_tracker';
import {IntrTransformControls, TransformSpace} from './intr_transform_controls';
import {TransformTarget} from './transform_target';
import {WorldScene} from './world_scene';

/**
 * The parts of a view (usually a `ViewController`) that a `Transformer` needs
 * to show its transform controls in.
 */
export interface TransformerView {
  readonly scene: WorldScene;
  camera(): three.Camera;
  getCanvasElement(): HTMLCanvasElement;
  requestNewFrame(): void;
  registerFrameUpdateCallback(
    callback: () => void,
    type: FrameUpdateCallbackType,
  ): number;
  unregisterFrameUpdateCallback(id: number): void;
}

/**
 * The Transformer class is responsible for applying rigid transform to a set
 * objects and/or frames. It owns a TransformControls object that can be set to
 * attach to one or more objects/frames to perform a transform action.
 *
 * What the controls attach to, and what happens when they are dragged, is
 * determined by a `TransformTarget` (see `attach()`), e.g. a
 * `RigidObjectTarget` that moves a set of objects, or a robot posing target
 * (see `RobotPoser`).
 *
 * Things user can do with the transformer now:
 *
 *  - Drag the controls to transform the attached target.
 *  - Toggle between "translation" and "rotation" modes.
 *  - Toggle between "local", "parent" and "world" reference frame.
 *
 * The Transformer class exposes a set of Observables that can be subscribed to
 * act upon distinct events.
 */
export class Transformer {
  private readonly transControls: IntrTransformControls;

  private enabled = true;
  // The target that the transform controls are currently attached to.
  private attachedTarget?: TransformTarget;
  // This variable is used to track if the user actually changed any
  // transforms while manipulating the object (versus just selecting then deselecting it)
  private didChangeTransforms = false;
  // Whether the transform controls are currently being dragged.
  private dragging = false;

  // Observables that notify subscribers of distinct events.
  private readonly transformerTouchedSubject = new Subject<boolean>();
  private readonly transformStartSubject = new Subject<void>();
  private readonly transformStopSubject = new Subject<void>();
  private readonly targetChangedSubject = new Subject<
    TransformTarget | undefined
  >();
  private readonly enabledChangedSubject = new Subject<boolean>();

  // Fires event when one of transformer's axies turned active (hovered).
  readonly transformerTouched: Observable<boolean> =
    this.transformerTouchedSubject;
  // Fires event when transformer started getting dragged.
  readonly transformStarted: Observable<void> = this.transformStartSubject;
  // Fires event when transformer stopped getting dragged.
  readonly transformStopped: Observable<void> = this.transformStopSubject;
  // Fires event with the new target (or `undefined` if detached) whenever the
  // transform controls are attached to a target or detached.
  readonly targetChanged: Observable<TransformTarget | undefined> =
    this.targetChangedSubject;
  // Fires event with the new state whenever the transformer is enabled or
  // disabled.
  readonly enabledChanged: Observable<boolean> = this.enabledChangedSubject;
  private readonly frameUpdateCallbackId: number;
  private readonly worldUpdatedSubscription: Subscription;

  TEST_ONLY = {
    getTransformControls: () => {
      return this.transControls;
    },
  };

  constructor(private readonly view: TransformerView) {
    if (!this.view.scene) {
      throw new Error('The view is not initialized');
    }
    this.transControls = new IntrTransformControls(
      this.view.camera(),
      this.view.getCanvasElement(),
    );

    // Disable camera control when the transformer is being dragged.
    this.transControls.addEventListener('start', () => {
      this.dragging = true;
      // Remember the object pose when action started.
      this.transformStartSubject.next();
      this.didChangeTransforms = false;
    });

    this.transControls.addEventListener('end', () => {
      this.dragging = false;
      this.attachedTarget?.onDragEnd?.(this.didChangeTransforms);
      this.transformStopSubject.next();
    });

    // Keep requesting new frames while object is being transformed.
    this.transControls.addEventListener('change', () => {
      const handle = this.transControls.parent;
      if (handle) {
        this.attachedTarget?.onDrag?.(handle.matrixWorld.clone());
      }
      this.didChangeTransforms = true;
      this.view.requestNewFrame();
    });

    this.transControls.addEventListener('hover-changed', (e) => {
      this.transformerTouchedSubject.next(e.hovering);
      this.view.requestNewFrame();
    });

    // Make sure the helper's visual size scale as view zoom changes.
    this.frameUpdateCallbackId = this.view.registerFrameUpdateCallback(() => {
      if (this.transControls.parent) {
        this.transControls.updateMatrixWorld(true);
      }
    }, FrameUpdateCallbackType.NEW_FRAME_ONLY);

    // Let the target react to changes of the world, e.g. to move its phantom
    // to the new pose of the objects it represents.
    this.worldUpdatedSubscription = this.view.scene.root.updated.subscribe(
      () => {
        if (this.dragging) return;
        this.attachedTarget?.onWorldUpdated?.();
        this.view.requestNewFrame();
      },
    );
  }

  /**
   * Enables/disables the transformer from being able to update an object's
   * transform. Setting transform space and mode is always allowed.
   *
   * Disabling the transformer detaches it from its current target. While
   * disabled, `attach()` does not attach to new targets.
   */
  setEnabled(enabled: boolean) {
    this.transControls.setEnabled(enabled);
    this.transControls.visible = enabled;
    if (enabled === this.enabled) return;
    this.enabled = enabled;
    if (!enabled) {
      this.detach();
    }
    this.enabledChangedSubject.next(enabled);
    this.view.requestNewFrame();
  }

  // Set the transform space.
  setTransformSpace(space: TransformSpace) {
    this.transControls.setSpace(space);
    this.view.requestNewFrame();
  }

  getTransformSpace() {
    return this.transControls.getSpace();
  }

  cycleTransformSpace() {
    const currentSpace = this.transControls.getSpace();
    if (currentSpace === TransformSpace.WORLD) {
      this.transControls.setSpace(TransformSpace.PARENT);
    } else if (currentSpace === TransformSpace.PARENT) {
      this.transControls.setSpace(TransformSpace.LOCAL);
    } else {
      this.transControls.setSpace(TransformSpace.WORLD);
    }
    this.view.requestNewFrame();
  }

  /** The view that the transformer is shown in. */
  getView(): TransformerView {
    return this.view;
  }

  /** The target that the transform controls are attached to, if any. */
  getTransformTarget(): TransformTarget | undefined {
    return this.attachedTarget;
  }

  /**
   * Attaches the transform controls to `target`, detaching from the previous
   * target first.
   *
   * If the transformer is disabled, `target` is disposed right away and the
   * transformer is not attached.
   */
  attach(target: TransformTarget) {
    if (!this.enabled) {
      target.dispose();
      return;
    }
    this.detachTarget();
    this.attachedTarget = target;
    this.transControls.enableSurfacePlacement(target.enableSurfacePlacement);
    this.transControls.attachTo(target.handle);
    this.transControls.setSurfaceGenerator(() => this.getSnappingSurfaces());
    this.transControls.updateMatrixWorld(true);
    this.transControls.setEnabled(this.enabled);
    this.view.requestNewFrame();
    this.targetChangedSubject.next(target);
  }

  /**
   * Detaches the transform controls from the current target and disposes the
   * target. No-op if not attached.
   */
  detach() {
    if (!this.attachedTarget) return;
    this.detachTarget();
    this.targetChangedSubject.next(undefined);
  }

  private detachTarget() {
    if (!this.attachedTarget) return;
    this.attachedTarget.dispose();
    this.attachedTarget = undefined;
    this.transControls.setEnabled(false);
    this.transControls.detach();
    this.view.requestNewFrame();
  }

  /**
   * Detach the current object from transformer.
   */
  detachTransformer() {
    this.detach();
  }

  private getSnappingSurfaces(): Iterable<three.Object3D> {
    return this.attachedTarget?.snappingSurfaces?.() ?? [];
  }

  /**
   * Detaches the transformer, removes all of its event listeners and frame
   * callbacks, and completes all of its observables. The transformer must not
   * be used after calling this.
   */
  dispose() {
    this.worldUpdatedSubscription.unsubscribe();
    this.detach();
    this.view.unregisterFrameUpdateCallback(this.frameUpdateCallbackId);
    this.transControls.dispose();
    this.transformerTouchedSubject.complete();
    this.transformStartSubject.complete();
    this.transformStopSubject.complete();
    this.targetChangedSubject.complete();
    this.enabledChangedSubject.complete();
  }

  /** Makes the transformer work with the given camera from now on. */
  setCamera(camera: three.Camera) {
    this.transControls.setCamera(camera);
    this.transControls.updateMatrixWorld(true);
    this.view.requestNewFrame();
  }
}
