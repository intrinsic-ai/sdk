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

import * as mathUtils from '../utils/math';
import * as colors from './colors';
import {BVHMesh, createRingGeometry} from './shapes';
import {IntrEntityNode} from './world_scene_nodes';

/**
 * The space in which the transform controls are operating.
 */
export enum TransformSpace {
  WORLD = 'world',
  PARENT = 'parent',
  LOCAL = 'local',
}

// Minimum distance in meters required to snap to a collision surface when
// shift is held.
const SNAP_DISTANCE_THRESHOLD = 0.1;
const ARC_SCALE = 0.75;
const SIZE = 0.15;
const WIDTH = 0.003;
const OCTOHEDRON_SIZE = 0.025;
const ARC_OFFSET = 0.15;

const CYLINDER_GEOMETRY = new three.CylinderGeometry(WIDTH, WIDTH, SIZE, 8);
const CYLINDER_COLLISION_GEOMETRY = new three.CylinderGeometry(
  WIDTH + 0.005,
  WIDTH + 0.005,
  SIZE,
  8,
);
const CONE_GEOMETRY = new three.ConeGeometry(0.01, 0.04, 16);
const CONE_COLLISION_GEOMETRY = new three.ConeGeometry(0.015, 0.05, 16);
const OCTOHEDRON_GEOMETRY = new three.OctahedronGeometry(OCTOHEDRON_SIZE);
const ARC_GEOMETRY = new three.TorusGeometry(
  SIZE * ARC_SCALE,
  WIDTH,
  8,
  16,
  Math.PI / 2 - 2 * ARC_OFFSET,
);
const ARC_COLLISION_GEOMETRY = new three.TorusGeometry(
  SIZE * ARC_SCALE,
  WIDTH + 0.005,
  8,
  16,
  Math.PI / 2 - 2 * ARC_OFFSET,
);
const RING_GUIDE_GEOMETRY = createRingGeometry(SIZE * ARC_SCALE, 64);
const ARROW_GUIDE_GEOMETRY = new three.BufferGeometry().setFromPoints([
  new three.Vector3(0, -100, 0),
  new three.Vector3(0, 100, 0),
]);

const MAT_DEFAULTS = {
  transparent: true,
  opacity: 1,
  toneMapped: false,
  depthTest: false,
};

class Handle extends three.Group {
  private guide?: three.LineSegments;

  constructor(readonly axis: three.Vector3) {
    super();
  }

  fade(enabled: boolean) {
    this.visible = !enabled;
  }

  setGuide(guide: three.LineSegments) {
    this.guide = guide;
    this.add(guide);
  }

  showGuide(visible: boolean) {
    if (this.guide) {
      this.guide.visible = visible;
    }
  }
}

class HandleMaterial extends three.MeshBasicMaterial {
  protected readonly originalColor = this.color.clone();
  protected readonly dimColor = this.color.clone().multiplyScalar(0.3);

  override copy(mat: HandleMaterial): this {
    super.copy(mat);
    this.originalColor.copy(mat.originalColor);
    this.dimColor.copy(mat.dimColor);
    return this;
  }

  setActive(active: boolean) {
    if (active) {
      this.color.copy(this.originalColor);
    } else {
      this.color.copy(this.dimColor);
    }
  }
}

class HandleMesh extends BVHMesh {
  constructor(
    geometry: three.BufferGeometry,
    override material: HandleMaterial,
  ) {
    super(geometry, material);
  }
}

const X_MAT = new HandleMaterial({
  color: colors.X_AXIS.clone(),
  ...MAT_DEFAULTS,
});
const Y_MAT = new HandleMaterial({
  color: colors.Y_AXIS.clone(),
  ...MAT_DEFAULTS,
});
const Z_MAT = new HandleMaterial({
  color: colors.Z_AXIS.clone(),
  ...MAT_DEFAULTS,
});
const GREY_MAT = new HandleMaterial({
  color: new three.Color(0x333333),
  ...MAT_DEFAULTS,
});

function makeArrows(material: HandleMaterial, axis: three.Vector3): Handle {
  const matClone = material.clone();
  const visible = new three.Group().add(
    //The arrow tips
    new HandleMesh(CONE_GEOMETRY, matClone).translateY(
      SIZE + OCTOHEDRON_SIZE + 0.01,
    ),
    new HandleMesh(CONE_GEOMETRY, matClone)
      .rotateX(Math.PI)
      .translateY(SIZE + OCTOHEDRON_SIZE + 0.01),
    // the arrow shaft for just the positive axis
    new HandleMesh(CYLINDER_GEOMETRY, matClone).translateY(
      SIZE / 2 + OCTOHEDRON_SIZE + 0.01,
    ),
  );
  const collision = new three.Group().add(
    //The arrow tips
    new HandleMesh(CONE_COLLISION_GEOMETRY, matClone).translateY(
      SIZE + OCTOHEDRON_SIZE + 0.01,
    ),
    new HandleMesh(CONE_COLLISION_GEOMETRY, matClone)
      .rotateX(Math.PI)
      .translateY(SIZE + OCTOHEDRON_SIZE + 0.01),
    // the arrow shaft for just the positive axis
    new HandleMesh(CYLINDER_COLLISION_GEOMETRY, matClone).translateY(
      SIZE / 2 + OCTOHEDRON_SIZE + 0.01,
    ),
  );
  collision.visible = false;
  const arrow = new Handle(axis.clone()).add(visible, collision);
  arrow.setGuide(new three.LineSegments(ARROW_GUIDE_GEOMETRY, matClone));
  arrow.applyQuaternion(
    new three.Quaternion().setFromUnitVectors(new three.Vector3(0, 1, 0), axis),
  );
  return arrow;
}

function makeRing(material: HandleMaterial, axis: three.Vector3): Handle {
  const matClone = material.clone();
  const toAxis = new three.Quaternion().setFromUnitVectors(
    new three.Vector3(0, 0, 1),
    axis,
  );
  const handle = new Handle(axis.clone());
  for (let i = 0; i < 4; i++) {
    const ring = new HandleMesh(ARC_GEOMETRY, matClone);
    ring.rotateZ(ARC_OFFSET + (i * Math.PI) / 2);
    handle.add(ring);
    // Add an invisible, larger arc that's easier to grab.
    const collision = new HandleMesh(ARC_COLLISION_GEOMETRY, matClone);
    collision.visible = false;
    collision.rotateZ(ARC_OFFSET + (i * Math.PI) / 2);
    handle.add(collision);
  }
  const guideMat = material.clone();
  guideMat.clippingPlanes = [];
  handle.setGuide(new three.LineSegments(RING_GUIDE_GEOMETRY, guideMat));
  handle.applyQuaternion(toAxis);
  return handle;
}
/**
 * Event map for IntrTransformControls.
 */
export interface IntrTransformControlsEventMap extends three.Object3DEventMap {
  start: unknown;
  end: unknown;
  change: {worldMatrix: three.Matrix4};
  'hover-changed': {hovering: boolean};
}

/**
 * A custom TransformControls that places objects to be flush with the hover-dragged surface.
 */
export class IntrTransformControls extends three.Object3D<IntrTransformControlsEventMap> {
  private objectBox = new three.Box3(new three.Vector3(), new three.Vector3());
  private draggingHandle?: Handle;
  private hoveringHandle?: Handle;
  private enabled = true;
  private readonly plane = new three.Plane();
  private readonly prevMouseCoords = new three.Vector2();
  private space = TransformSpace.WORLD;
  private readonly rayCaster = new three.Raycaster();
  private surfacePlacementEnabled = false;
  private readonly initialQuaternion = new three.Quaternion();
  private surfaceGenerator?: () => Iterable<three.Object3D>;
  // Aborted on `dispose()` to remove all DOM event listeners.
  private readonly listenerAbortController = new AbortController();

  private readonly centerHandle = new Handle(new three.Vector3()).add(
    new HandleMesh(OCTOHEDRON_GEOMETRY, GREY_MAT),
  );
  private readonly axisHandles = new Set<Handle>([
    makeArrows(X_MAT, new three.Vector3(1, 0, 0)),
    makeArrows(Y_MAT, new three.Vector3(0, 1, 0)),
    makeArrows(Z_MAT, new three.Vector3(0, 0, 1)),
  ]);
  private readonly rotationHandles = new Set<Handle>([
    makeRing(X_MAT, new three.Vector3(1, 0, 0)),
    makeRing(Y_MAT, new three.Vector3(0, 1, 0)),
    makeRing(Z_MAT, new three.Vector3(0, 0, 1)),
  ]);

  constructor(
    private camera: three.Camera,
    private readonly canvas: HTMLCanvasElement,
  ) {
    super();
    this.add(...this.axisHandles);
    this.add(...this.rotationHandles);
    for (const handle of this.rotationHandles) {
      handle.traverse((n) => {
        if (n instanceof HandleMesh) {
          n.material.clippingPlanes = [this.plane];
        }
      });
    }
    this.add(this.centerHandle);

    const listenerOptions = {signal: this.listenerAbortController.signal};
    this.canvas.addEventListener(
      'pointerdown',
      (e: PointerEvent) => {
        this.prevMouseCoords.copy(this.getNdcFromInput(e));

        if (!this.enabled) return;
        const handle = this.getHandleFromRaycast(e);
        if (handle) {
          this.draggingHandle = handle;
          for (const child of this.children) {
            if (child instanceof Handle && child !== handle) {
              child.fade(true);
            }
          }
          this.dispatchEvent({type: 'start'});
        }
      },
      listenerOptions,
    );

    this.canvas.addEventListener(
      'pointerup',
      () => {
        if (!this.enabled) return;
        if (this.draggingHandle) {
          this.stopDrag();
        }
      },
      listenerOptions,
    );

    this.canvas.addEventListener(
      'pointermove',
      (e: PointerEvent) => {
        if (!this.enabled) return;
        // If the mouse button is not pressed, we should stop dragging. This is
        // necessary for the case where the user clicks on a handle, moves the
        // mouse out of the canvas, and then releases the mouse button.
        if (e.buttons === 0 && this.draggingHandle) {
          this.stopDrag();
        }
        if (this.draggingHandle) {
          this.doHandleDrag(e);
        } else {
          this.doHandleHover(e);
        }
      },
      listenerOptions,
    );
  }

  /** Makes the controls work with the given camera from now on. */
  setCamera(camera: three.Camera) {
    this.camera = camera;
  }

  private stopDrag() {
    this.draggingHandle = undefined;
    for (const child of this.children) {
      if (child instanceof Handle) {
        child.fade(false);
      }
    }
    this.updateBox();
    this.dispatchEvent({type: 'end'});
  }

  private doHandleHover(e: PointerEvent) {
    this.clearFade();
    const handle = this.getHandleFromRaycast(e);
    if (handle) {
      this.fadeAllExcept(handle);
      handle.showGuide(true);
    }
    if (this.hoveringHandle !== handle) {
      this.dispatchEvent({type: 'hover-changed', 'hovering': !!handle});
    }
    this.hoveringHandle = handle;
  }

  private clearFade() {
    for (const child of this.children) {
      if (child instanceof Handle) {
        child.fade(false);
        child.showGuide(false);
      }
    }
  }

  private fadeAllExcept(handle: Handle) {
    for (const child of this.children) {
      if (child instanceof Handle && child !== handle) {
        child.fade(true);
        child.showGuide(false);
      }
    }
  }

  private getIntersection(
    ndc: three.Vector2,
    normal: three.Vector3,
  ): three.Vector3 | undefined {
    const pos = this.getWorldPosition(new three.Vector3());
    this.rayCaster.setFromCamera(ndc, this.camera);
    const ray = this.rayCaster.ray;
    const denom = normal.dot(ray.direction);
    if (Math.abs(denom) > 0.0001) {
      const t = pos.clone().sub(ray.origin).dot(normal) / denom;
      if (t < 0) {
        return undefined;
      }
      return ray.origin.clone().addScaledVector(ray.direction, t);
    }
    return undefined;
  }

  private doHandleDrag(e: PointerEvent) {
    if (!this.draggingHandle) return;
    const mouseCoords = this.getNdcFromInput(e);
    if (this.axisHandles.has(this.draggingHandle)) {
      this.doAxisDrag(mouseCoords, this.draggingHandle);
    } else if (this.rotationHandles.has(this.draggingHandle)) {
      this.doRingDrag(mouseCoords, this.draggingHandle);
    } else if (this.draggingHandle === this.centerHandle) {
      this.doCenterDrag(mouseCoords, this.draggingHandle);
    }
    if (this.shouldSnap(e)) {
      this.doSurfacePlacement(e);
    }
    this.prevMouseCoords.copy(mouseCoords);
  }

  private doAxisDrag(mouseCoords: three.Vector2, handle: Handle) {
    if (!this.parent) return;
    const parentWorldQuat = new three.Quaternion();
    const worldQuat = this.getWorldQuaternion(new three.Quaternion());
    const worldPos = new three.Vector3();
    this.parent.matrixWorld.decompose(
      worldPos,
      parentWorldQuat,
      new three.Vector3(),
    );
    const axis = handle.axis.clone().applyQuaternion(worldQuat);
    const negDir = this.camera.getWorldDirection(new three.Vector3()).negate();
    const normal = axis.clone().cross(negDir).cross(axis);
    const oldIntersection = this.getIntersection(this.prevMouseCoords, normal);
    const newIntersection = this.getIntersection(mouseCoords, normal);
    if (!oldIntersection || !newIntersection) return;
    const delta = newIntersection.clone().sub(oldIntersection);
    worldPos.add(axis.clone().multiplyScalar(delta.dot(axis)));
    const mat = new three.Matrix4().compose(
      worldPos,
      parentWorldQuat,
      new three.Vector3(1, 1, 1),
    );
    this.emitChange(mat);
  }

  private doRingDrag(mouseCoords: three.Vector2, handle: Handle) {
    if (!this.parent) return;
    const parentWorldQuat = new three.Quaternion();
    const worldQuat = this.getWorldQuaternion(new three.Quaternion());
    const worldPos = new three.Vector3();
    this.parent.matrixWorld.decompose(
      worldPos,
      parentWorldQuat,
      new three.Vector3(),
    );
    const worldAxis = handle.axis.clone().applyQuaternion(worldQuat);

    const cameraDir = this.camera
      .getWorldDirection(new three.Vector3())
      .negate();
    if (worldAxis.dot(cameraDir) < 0) {
      worldAxis.negate();
    }
    const oldIntersection = this.getIntersection(
      this.prevMouseCoords,
      worldAxis,
    );
    const newIntersection = this.getIntersection(mouseCoords, worldAxis);
    if (!oldIntersection || !newIntersection) return;
    const v1 = newIntersection.clone().sub(worldPos).normalize();
    const v2 = oldIntersection.clone().sub(worldPos).normalize();
    const quat = new three.Quaternion().setFromUnitVectors(v2, v1);
    parentWorldQuat.premultiply(quat);
    const mat = new three.Matrix4().compose(
      worldPos,
      parentWorldQuat,
      new three.Vector3(1, 1, 1),
    );
    this.emitChange(mat);
  }

  private doCenterDrag(mouseCoords: three.Vector2, handle: Handle) {
    if (!this.parent) return;
    const worldPos = this.getWorldPosition(new three.Vector3());
    const worldQuat = this.getWorldQuaternion(new three.Quaternion());
    const normal = this.camera.getWorldDirection(new three.Vector3()).negate();
    const oldIntersection = this.getIntersection(this.prevMouseCoords, normal);
    const newIntersection = this.getIntersection(mouseCoords, normal);
    if (!oldIntersection || !newIntersection) return;
    const delta = newIntersection.clone().sub(oldIntersection);
    worldPos.add(delta);
    this.emitChange(
      new three.Matrix4().compose(
        worldPos,
        worldQuat,
        new three.Vector3(1, 1, 1),
      ),
    );
  }

  private emitChange(matrix: three.Matrix4) {
    this.dispatchEvent({type: 'change', 'worldMatrix': matrix});
    if (this.parent) {
      mathUtils.updateObjectWorldTransform(this.parent, matrix.clone());
    }
  }

  private doSurfacePlacement(e: PointerEvent) {
    if (!this.parent || !this.surfaceGenerator) return;
    const surfaces = [...this.surfaceGenerator()];
    this.setupRaycaster(e);
    const distanceThreshold =
      this.draggingHandle === this.centerHandle
        ? Infinity
        : SNAP_DISTANCE_THRESHOLD;

    const intersection = this.rayCaster.intersectObjects(surfaces, true)[0];
    if (!intersection) return;
    const normal = intersection.face?.normal;
    if (normal && intersection.distance < distanceThreshold) {
      // The normal returned from the raycast is in object space, so we must first transform it to world space.
      normal.applyQuaternion(
        intersection.object.getWorldQuaternion(new three.Quaternion()),
      );
      const up = new three.Vector3(0, 0, 1);
      const quat = new three.Quaternion()
        .setFromUnitVectors(up, normal)
        .multiply(this.initialQuaternion);
      const pos = intersection.point
        .clone()
        .sub(normal.clone().multiplyScalar(this.objectBox.min.z));
      this.emitChange(
        new three.Matrix4().compose(pos, quat, new three.Vector3(1, 1, 1)),
      );
    }
    this.parent.updateMatrix();
  }

  private getNdcFromInput(e: PointerEvent): three.Vector2 {
    const rect = this.canvas.getBoundingClientRect();
    const ndcX = (e.offsetX / rect.width) * 2 - 1;
    const ndcY = -(e.offsetY / rect.height) * 2 + 1;
    return new three.Vector2(ndcX, ndcY);
  }

  private updateLocalQuaternion() {
    if (this.space === TransformSpace.PARENT) {
      // the grandparent is the parent of the object we're controlling
      if (this.parent) {
        this.quaternion.copy(this.parent.quaternion.clone().invert());
      } else {
        this.quaternion.identity();
      }
    } else if (this.space === TransformSpace.WORLD) {
      if (this.parent) {
        this.quaternion.copy(
          this.parent.getWorldQuaternion(new three.Quaternion()).invert(),
        );
      } else {
        this.quaternion.identity();
      }
    } else {
      // Local space is the identity quaternion.
      this.quaternion.identity();
    }
  }

  attachTo(node: three.Object3D) {
    // preserve the world space orientation when computing the bounding box.
    node.getWorldQuaternion(this.initialQuaternion);
    if (node !== this.parent) {
      node.add(this);
      this.updateBox();
    }
  }

  private updateBox() {
    if (this.parent instanceof IntrEntityNode) {
      const mat = this.parent.matrixWorld.clone();
      mathUtils.updateObjectWorldTransform(
        this.parent,
        new three.Matrix4().makeRotationFromQuaternion(this.initialQuaternion),
      );
      this.objectBox =
        this.parent.collisionNode.children.length > 0
          ? this.parent.getCollisionBounds()
          : new three.Box3(new three.Vector3(), new three.Vector3());
      mathUtils.updateObjectWorldTransform(this.parent, mat);
    }
  }

  private orientArrowsToCamera(
    worldPos: three.Vector3,
    worldQuat: three.Quaternion,
  ) {
    const cameraWorldPos = this.camera.getWorldPosition(new three.Vector3());
    const toCamera = cameraWorldPos.sub(worldPos);
    for (const handle of this.axisHandles) {
      const worldAxis = handle.axis.clone().applyQuaternion(worldQuat);
      handle.scale.setScalar(Math.sign(worldAxis.dot(toCamera)));
    }
  }

  override updateMatrix() {
    this.updateLocalQuaternion();
    // Compute world pose from the parent's matrix instead of reading stale
    // `this.matrixWorld`. In Three.js, `updateMatrix()` is called during
    // traversal *before* the object's `matrixWorld` is resolved with the new parent.
    // Reading `this.matrixWorld` here would return stale data from the previous frame.
    const worldPos = new three.Vector3().copy(this.position);
    const worldQuat = new three.Quaternion().copy(this.quaternion);
    if (this.parent) {
      worldPos.applyMatrix4(this.parent.matrixWorld);
      worldQuat.premultiply(
        this.parent.getWorldQuaternion(new three.Quaternion()),
      );
    }
    const cameraDir = this.camera
      .getWorldDirection(new three.Vector3())
      .negate();
    // Scale by the length of a screen space unit-vector.
    const up = new three.Vector3(0, 1, 0)
      .applyQuaternion(this.camera.quaternion)
      .add(worldPos)
      .project(this.camera);
    this.orientArrowsToCamera(worldPos, worldQuat);
    const scale = worldPos.clone().project(this.camera).distanceTo(up);
    this.scale.setScalar(1 / scale);
    this.plane.setFromNormalAndCoplanarPoint(cameraDir, worldPos);
    super.updateMatrix();
  }

  private getHandleFromRaycast(e: PointerEvent): Handle | undefined {
    this.rayCaster.setFromCamera(this.getNdcFromInput(e), this.camera);
    const surfaces: HandleMesh[] = [];
    this.traverse((n) => {
      if (n instanceof HandleMesh) {
        surfaces.push(n);
      }
    });
    const intersection = this.rayCaster.intersectObjects(surfaces, true)[0];
    if (intersection) {
      let handle = intersection.object;
      while (handle.parent && handle.parent !== this) {
        handle = handle.parent;
      }
      if (handle instanceof Handle) {
        // Reject rotation handles behind the clipping plane
        if (
          this.rotationHandles.has(handle) &&
          this.plane.distanceToPoint(intersection.point) < 0
        ) {
          return undefined;
        }
        return handle;
      }
    }
    return undefined;
  }

  enableSurfacePlacement(enabled: boolean) {
    this.surfacePlacementEnabled = enabled;
  }

  setSpace(space: TransformSpace) {
    this.space = space;
    this.updateMatrix();
  }

  getSpace() {
    return this.space;
  }

  setEnabled(enabled: boolean) {
    this.enabled = enabled;
  }

  detach(): this {
    this.removeFromParent();
    this.objectBox = new three.Box3(new three.Vector3(), new three.Vector3());
    return this;
  }

  /**
   * Detaches the controls and removes all event listeners registered on the
   * canvas. The controls must not be used after calling this.
   */
  override dispose() {
    this.listenerAbortController.abort();
    this.draggingHandle = undefined;
    this.hoveringHandle = undefined;
    this.detach();
    super.dispose();
  }

  private shouldSnap(e: PointerEvent): boolean {
    return (
      !!this.draggingHandle &&
      this.surfacePlacementEnabled &&
      (this.draggingHandle === this.centerHandle ||
        (this.axisHandles.has(this.draggingHandle) &&
          e.getModifierState('Shift')))
    );
  }

  private setupRaycaster(e: PointerEvent) {
    const origin = this.getWorldPosition(new three.Vector3());
    if (!origin || !this.draggingHandle) return;
    if (this.draggingHandle === this.centerHandle) {
      this.rayCaster.setFromCamera(this.getNdcFromInput(e), this.camera);
      return;
    }

    const mouseCoords = this.getNdcFromInput(e);
    const axis = this.draggingHandle.axis
      .clone()
      .applyQuaternion(this.getWorldQuaternion(new three.Quaternion()));
    const projectedOrigin = origin.clone().project(this.camera);
    const projectedAxis = axis
      .clone()
      .add(origin)
      .project(this.camera)
      .sub(projectedOrigin);
    const axisInCamera = new three.Vector2(projectedAxis.x, projectedAxis.y);
    const mouseInCamera = mouseCoords.sub(this.prevMouseCoords);
    const sign = Math.sign(axisInCamera.dot(mouseInCamera));
    this.rayCaster.set(origin, axis.multiplyScalar(sign));
  }

  setSurfaceGenerator(generator: () => Iterable<three.Object3D>) {
    this.surfaceGenerator = generator;
  }
}
