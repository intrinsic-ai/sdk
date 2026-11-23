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

import {clamp} from '@intrinsic-ai/core/math';
import * as three from 'three';

const UP = new three.Vector3(0, 0, 1);
const RIGHT = new three.Vector3(1, 0, 0);
const DELTA_SCALE = Math.PI * 2;
// min angle to Z in radians
const MIN_ANGLE_Z = 0.001;
// Zoom per pixel of wheel scroll. Mouse wheels typically report 100px per
// notch, so one notch zooms by about 100 * 0.0008 = 8%.
const WHEEL_ZOOM_SENSITIVITY = 0.0008;
// Zoom per pixel of trackpad pinch (wheel events with ctrlKey), which report
// much smaller deltas than scrolling.
const WHEEL_PINCH_ZOOM_SENSITIVITY = 0.01;
// Pixels per line/page for wheel events not reported in pixels (e.g. Firefox).
const WHEEL_LINE_HEIGHT_PX = 16;
const WHEEL_PAGE_HEIGHT_PX = 800;
// Caps a single wheel event to avoid huge jumps from fast or free-spinning
// wheels.
const WHEEL_MAX_DELTA_PX = 300;
// Tighter cap for pinch events, so that Ctrl + mouse wheel (which is reported
// the same way) doesn't zoom excessively.
const WHEEL_MAX_PINCH_DELTA_PX = 25;
// Smallest factor by which a single zoom step may scale the camera distance.
const MIN_ZOOM_FACTOR = 0.1;

/**
 * The type of motion the user is performing with the mouse.
 */
export enum MotionType {
  NONE = -1,
  PAN = 0,
  ORBIT = 1,
  ZOOM = 2,
}

/**
 * A mapping of JS pointer button values
 */
export enum PointerButton {
  LEFT = 0,
  MIDDLE = 1,
  RIGHT = 2,
}

/**
 * Event map for IntrOrbitControls.
 */
export interface IntrOrbitControlsEventMap {
  start: unknown;
  end: unknown;
  change: {
    motionType: MotionType;
    oldPos: three.Vector3;
    oldTarget: three.Vector3;
  };
}

/**
 * Custom OrbitControls that removes useless threejs logic and makes the API sane.
 */
export class IntrOrbitControls extends three.EventDispatcher<IntrOrbitControlsEventMap> {
  private enabled = true;
  readonly target = new three.Vector3();
  private moving = false;
  private motionType = MotionType.NONE;
  private readonly startingMouseCoords = new three.Vector2();
  enableZoom = true;
  enablePan = true;

  private readonly buttonToMotionType = new Map<PointerButton, MotionType>([
    [PointerButton.LEFT, MotionType.ORBIT],
    [PointerButton.MIDDLE, MotionType.PAN],
    [PointerButton.RIGHT, MotionType.ZOOM],
  ]);

  constructor(
    readonly camera: three.Camera,
    readonly canvas: HTMLCanvasElement,
    buttonMapping?: {
      pan: PointerButton;
      orbit: PointerButton;
      zoom: PointerButton;
    },
  ) {
    super();
    this.canvas.addEventListener('pointerdown', this.onPointerDown);
    this.canvas.addEventListener('pointermove', this.onPointerMove);
    this.canvas.addEventListener('pointerup', this.onPointerUp);
    // Passive must be false for preventDefault to take effect.
    this.canvas.addEventListener('wheel', this.onWheel, {passive: false});
    this.canvas.addEventListener(
      'lostpointercapture',
      this.onLostPointerCapture,
    );
    this.initTarget();
    if (buttonMapping) {
      this.buttonToMotionType.set(buttonMapping.pan, MotionType.PAN);
      this.buttonToMotionType.set(buttonMapping.orbit, MotionType.ORBIT);
      this.buttonToMotionType.set(buttonMapping.zoom, MotionType.ZOOM);
    }
  }

  private initTarget() {
    // compute the point that intersects a plane that's orthogonal to the camera that also passes through the origin.
    const quat = this.camera.quaternion;
    const pos = this.camera.position;
    const cameraToTarget = new three.Vector3(0, 0, 1).applyQuaternion(quat);
    const cameraToOrigin = pos.clone().sub(new three.Vector3());
    // Don't let the initial target be closer than 1m to the camera.
    const distance = Math.max(1, cameraToTarget.dot(cameraToOrigin));
    this.target.copy(
      new three.Vector3(0, 0, -distance).applyQuaternion(quat).add(pos),
    );
  }

  private mapButtonToMotionType(e: PointerEvent): MotionType {
    let motionType = this.buttonToMotionType.get(e.button) ?? MotionType.NONE;
    if (
      (e.getModifierState('Control') || e.getModifierState('Meta')) &&
      motionType === MotionType.PAN
    ) {
      motionType = MotionType.ORBIT;
    }
    return motionType;
  }

  /** Removes the controls' listeners from the canvas. */
  dispose() {
    this.canvas.removeEventListener('pointerdown', this.onPointerDown);
    this.canvas.removeEventListener('pointermove', this.onPointerMove);
    this.canvas.removeEventListener('pointerup', this.onPointerUp);
    this.canvas.removeEventListener('wheel', this.onWheel);
    this.canvas.removeEventListener(
      'lostpointercapture',
      this.onLostPointerCapture,
    );
  }

  private readonly onPointerDown = (e: PointerEvent) => {
    if (e.shiftKey) return;
    this.handlePointerDown(e);
  };

  private readonly onPointerMove = (e: PointerEvent) => {
    this.handlePointerMove(e);
  };

  private readonly onPointerUp = (e: PointerEvent) => {
    this.handlePointerUp(e);
  };

  private readonly onWheel = (e: WheelEvent) => {
    if (!this.enabled || !this.enableZoom) return;
    // Prevent the page from scrolling, and the browser from zooming the
    // page on trackpad pinch (reported as wheel events with ctrlKey).
    e.preventDefault();
    const oldPos = this.camera.position.clone();
    const oldTarget = this.target.clone();
    this.zoomCamera(this.wheelEventToZoomDelta(e));
    this.dispatchEvent({
      type: 'change',
      'motionType': MotionType.ZOOM,
      oldPos,
      oldTarget,
    });
  };

  // Pointer capture (see `handlePointerDown`) keeps delivering events while
  // dragging outside of the canvas. Only stop if the capture is lost for
  // another reason (e.g. the pointer is cancelled or the window loses focus).
  private readonly onLostPointerCapture = () => {
    if (this.motionType !== MotionType.NONE) {
      this.stopMoving();
    }
  };

  private stopMoving() {
    this.motionType = MotionType.NONE;
    this.moving = false;
    this.dispatchEvent({type: 'end'});
  }

  handlePointerDown(e: PointerEvent) {
    if (!this.enabled) return;
    if (this.motionType !== MotionType.NONE) return;
    this.motionType = this.mapButtonToMotionType(e);
    this.startingMouseCoords.set(e.offsetX, e.offsetY);
    if (this.motionType !== MotionType.NONE) {
      // Keep receiving pointer events while the drag leaves the canvas.
      this.canvas.setPointerCapture(e.pointerId);
    }
  }

  handlePointerMove(e: PointerEvent) {
    if (!this.enabled) return;
    if (this.motionType === MotionType.NONE) return;
    const oldPos = this.camera.position.clone();
    const oldTarget = this.target.clone();
    const offset = this.startingMouseCoords
      .clone()
      .sub(new three.Vector2(e.offsetX, e.offsetY));
    if (!this.moving) {
      this.moving = true;
      this.dispatchEvent({type: 'start'});
    }
    // Panning works in pixels so that the content follows the cursor exactly.
    const cursorMovement = offset.clone().negate();
    const delta = offset
      .divide(new three.Vector2(this.canvas.width, this.canvas.height))
      .multiplyScalar(DELTA_SCALE);
    switch (this.motionType) {
      case MotionType.PAN:
        this.panCamera(cursorMovement);
        break;
      case MotionType.ORBIT:
        this.orbitCamera(delta);
        break;
      case MotionType.ZOOM:
        this.zoomCamera(-delta.y);
        break;
      default:
        break;
    }
    this.startingMouseCoords.set(e.offsetX, e.offsetY);
    this.dispatchEvent({
      type: 'change',
      'motionType': this.motionType,
      oldPos,
      oldTarget,
    });
  }

  handlePointerUp(e: PointerEvent) {
    this.stopMoving();
    if (this.canvas.hasPointerCapture(e.pointerId)) {
      this.canvas.releasePointerCapture(e.pointerId);
    }
  }

  setEnabled(enabled: boolean) {
    this.enabled = enabled;
  }

  /**
   * Pans the camera so that the content at the depth of the orbit target
   * follows the cursor exactly ("grab" panning). This makes panning work
   * naturally regardless of if the user is zoomed very far in or out.
   *
   * @param cursorMovement The movement of the cursor in CSS pixels, with x
   *     pointing right and y pointing down.
   */
  panCamera(cursorMovement: three.Vector2) {
    if (!this.enablePan) return;
    const width = this.canvas.clientWidth;
    const height = this.canvas.clientHeight;
    if (width === 0 || height === 0) return;

    // We take the target's position in NDC then figure out the position `to`
    // in NDC that we want the camera to point at after the pan. This way, the
    // panning appears natural regardless of the zoom. For very zoomed in scenes,
    // a big pan in pixels will result in a small pan in meters. For very zoomed
    // out scenes, a big pan in pixels will result in a big pan in meters.

    // Make sure the view matrices reflect the latest camera pose.
    this.camera.updateMatrixWorld();
    const targetNdc = this.target.clone().project(this.camera);
    // The world point at the target's depth that is under the moved cursor.
    const to = new three.Vector3(
      targetNdc.x + (2 * cursorMovement.x) / width,
      targetNdc.y - (2 * cursorMovement.y) / height,
      targetNdc.z,
    ).unproject(this.camera);
    // Move the camera opposite to the cursor so the content follows it.
    const move = this.target.clone().sub(to);
    if (!isFinite(move.x) || !isFinite(move.y) || !isFinite(move.z)) return;
    this.camera.position.add(move);
    this.target.add(move);
  }

  orbitCamera(delta: three.Vector2) {
    const cameraQuat = this.camera.quaternion.clone();
    const cameraPos = this.camera.position.clone();
    const posToCenter = this.camera.position.clone().sub(this.target);
    const quatX = new three.Quaternion().setFromAxisAngle(UP, delta.x);
    cameraQuat.premultiply(quatX);
    const right = RIGHT.clone().applyQuaternion(cameraQuat);
    // compute the actual rotation while preventing rotation past +Z/-Z
    const angleToZ = UP.clone()
      .applyQuaternion(this.camera.quaternion)
      .angleTo(UP);
    const nextAngle = angleToZ + delta.y;
    const newAngle = clamp(nextAngle, MIN_ANGLE_Z, Math.PI - MIN_ANGLE_Z);
    const quatY = new three.Quaternion().setFromAxisAngle(
      right,
      delta.y + (newAngle - nextAngle),
    );
    cameraQuat.premultiply(quatY);
    posToCenter.applyQuaternion(quatX).applyQuaternion(quatY);
    cameraPos.copy(posToCenter).add(this.target);
    this.camera.position.copy(cameraPos);
    this.camera.quaternion.copy(cameraQuat);
  }

  setPose(position: three.Vector3, quaternion: three.Quaternion) {
    this.camera.position.copy(position);
    this.camera.quaternion.copy(quaternion);
    this.camera.updateMatrixWorld();
  }

  /**
   * Converts a wheel event into a zoom delta for `zoomCamera`.
   *
   * The scroll amount is normalized to pixels and mapped exponentially, so
   * that zooming in and out by the same amount cancels out and the camera can
   * never be moved through the target, no matter how large the delta.
   */
  private wheelEventToZoomDelta(e: WheelEvent): number {
    let deltaPx = e.deltaY;
    if (e.deltaMode === WheelEvent.DOM_DELTA_LINE) {
      deltaPx *= WHEEL_LINE_HEIGHT_PX;
    } else if (e.deltaMode === WheelEvent.DOM_DELTA_PAGE) {
      deltaPx *= this.canvas.clientHeight || WHEEL_PAGE_HEIGHT_PX;
    }
    // Trackpad pinch gestures report much smaller deltas than wheel scrolls.
    const [maxDeltaPx, sensitivity] = e.ctrlKey
      ? [WHEEL_MAX_PINCH_DELTA_PX, WHEEL_PINCH_ZOOM_SENSITIVITY]
      : [WHEEL_MAX_DELTA_PX, WHEEL_ZOOM_SENSITIVITY];
    deltaPx = clamp(deltaPx, -maxDeltaPx, maxDeltaPx);
    return Math.exp(deltaPx * sensitivity) - 1;
  }

  /**
   * Moves the camera towards or away from the target along the line between
   * them.
   *
   * @param zoomDelta The relative change of the distance to the target: the
   *     new distance is `distance * (1 + zoomDelta)`. Negative values zoom in
   *     (e.g. -0.1 moves 10% closer), positive values zoom out (e.g. 0.1 moves
   *     10% further away) and 0 does nothing. The resulting scale factor is
   *     clamped to at least `MIN_ZOOM_FACTOR`, so the camera never reaches or
   *     passes the target.
   */
  zoomCamera(zoomDelta: number) {
    if (!this.enableZoom) return;
    const posToCenter = this.camera.position.clone().sub(this.target);
    // Never scale by a non-positive factor, which would move the camera onto
    // or through the target.
    posToCenter.multiplyScalar(Math.max(1 + zoomDelta, MIN_ZOOM_FACTOR));
    this.camera.position.copy(posToCenter).add(this.target);
  }

  getDistance() {
    return this.camera.position.distanceTo(this.target);
  }

  getMoving() {
    return this.moving;
  }
}
