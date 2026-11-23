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

import {Observable, Subject, interval} from 'rxjs';
import {finalize, take} from 'rxjs/operators';
import * as three from 'three';

import {Clip} from '../model/entity';
import {updateObjectLocalTransform} from '../utils/math';
import * as colors from './colors';

/** Pinhole camera intrinsics, in pixels. */
export interface CameraIntrinsics {
  /** Height of the image. */
  readonly rows: number;
  /** Width of the image. */
  readonly cols: number;
  readonly focalLengthX: number;
  readonly focalLengthY: number;
  readonly principalPointX: number;
  readonly principalPointY: number;
}

/**
 * The camera projection type.
 */
export enum CameraProjection {
  PERSPECTIVE,
  ORTHOGRAPHIC,
}

// constants for the transition between camera projections.
const TRANSITION_TIME_MS = 150;
const TRANSITION_FREQUENCY_HZ = 60;
const NUM_ITERATIONS = (TRANSITION_TIME_MS / 1000) * TRANSITION_FREQUENCY_HZ;

/**
 * Create the camera used in ViewController for rendering the world. Optionally
 * set the intrinsic parameters.
 */
export class IntrinsicCamera extends three.Camera {
  zoom = 1;
  private readonly orthoCam = new three.OrthographicCamera(
    -0.5,
    0.5,
    0.5,
    -0.5,
  );
  private readonly perspectiveCam = new three.PerspectiveCamera(45, 1);
  private activeCamera: three.OrthographicCamera | three.PerspectiveCamera =
    this.perspectiveCam;
  private readonly intrinsicMatrix = new three.Matrix4();
  intrinsics?: CameraIntrinsics;
  private readonly updatedSubject = new Subject<void>();
  updated: Observable<void> = this.updatedSubject;
  // Make z Axis the up vector.
  override up = new three.Vector3(0, 0, 1);
  touched = false;
  private projection = CameraProjection.PERSPECTIVE;

  constructor() {
    super();
    this.add(this.orthoCam);
    this.add(this.perspectiveCam);
  }

  setProjection(projection: CameraProjection) {
    this.projection = projection;
    this.orthoCam.updateProjectionMatrix();
    this.perspectiveCam.updateProjectionMatrix();
    const fromCam = this.activeCamera;
    const [toCam, smoothFunc] =
      projection === CameraProjection.PERSPECTIVE
        ? [this.perspectiveCam, (t: number) => t ** 2]
        : [this.orthoCam, (t: number) => t ** 0.5];
    const startingMatrix = fromCam.projectionMatrix.clone();
    const endingMatrix = toCam.projectionMatrix.clone();
    interval(1000 / TRANSITION_FREQUENCY_HZ)
      .pipe(
        take(NUM_ITERATIONS),
        finalize(() => {
          this.activeCamera = toCam;
          this.updateProjectionMatrix();
        }),
      )
      .subscribe((iter) => {
        const t = smoothFunc(iter / NUM_ITERATIONS);
        for (let i = 0; i < 16; i++) {
          this.projectionMatrix.elements[i] =
            startingMatrix.elements[i] * (1 - t) + endingMatrix.elements[i] * t;
        }
        this.projectionMatrixInverse.copy(this.projectionMatrix).invert();
        this.updatedSubject.next();
      });
  }

  getProjection() {
    return this.projection;
  }

  setIntrinsics(intrinsics: CameraIntrinsics) {
    this.intrinsics = intrinsics;
    const {rows, cols} = intrinsics;
    const fx = intrinsics.focalLengthX;
    const fy = intrinsics.focalLengthY;
    const px = intrinsics.principalPointX;
    const py = intrinsics.principalPointY;

    // Construct WebGL projection matrix from the camera intrinsics. Note
    // that WebGL maps the screen to -1,1 and that the projection is
    // happening through the homogenous component (division by z happens at
    // a later stage).
    // Example for x coordinate when projecting a point [x, y, z]:
    // xHomogenous = x * (fx/cols*2) + z * (-px/cols*2)
    // wHomogenous = -z
    // xProjected = xHomogenous / wHomogenous = x * (fx/cols*2)/z + (-px/cols*2)
    const cam = this.perspectiveCam;
    this.intrinsicMatrix
      .set(
        /*0*/ (fx / cols) * 2.0,
        0,
        0,
        0,
        /*4*/ 0,
        (fy / rows) * 2.0,
        0,
        0,
        /*8*/ (-px / cols) * 2.0 + 1.0,
        (py / rows) * 2.0 - 1.0,
        -(cam.far + cam.near) / (cam.far - cam.near),
        -1,
        /*12*/ 0,
        0,
        (-2 * cam.far * cam.near) / (cam.far - cam.near),
        0,
      )
      .transpose();
    cam.aspect = cols / rows;
    this.updateProjectionMatrix();
  }

  /**
   * Sets the camera pose to the given object's world frame, e.g. to look
   * through a physical camera in the world, or from a robot's end tip.
   *
   * The object's frame is interpreted using the Intrinsic camera convention
   * (+Z forward, +Y down, +X right). This is a one-off pose update; the camera
   * does not follow the object afterwards.
   */
  setPoseFromObject(object: three.Object3D) {
    updateObjectLocalTransform(this, object.matrixWorld, true);
    // The rendering camera looks along its -Z axis, whereas the object's +Z
    // axis is the viewing direction, so flip around Y.
    this.rotateY(Math.PI);
    // Then rotate around Z to match the Intrinsic camera convention where the
    // y axis points downward and x points right.
    this.rotateZ(Math.PI);
    this.updateMatrixWorld(true);
  }

  // A bunch of delegate accessors to conform to the threejs OrthographicCamera and PerspectiveCamera interface.
  set aspect(aspect: number) {
    this.perspectiveCam.aspect = aspect;
    this.orthoCam.left = aspect * this.orthoCam.bottom;
    this.orthoCam.right = aspect * this.orthoCam.top;
    this.updateProjectionMatrix();
  }

  set near(near: number) {
    this.perspectiveCam.near = near;
    this.orthoCam.near = near;
    this.updateProjectionMatrix();
  }

  set far(far: number) {
    this.perspectiveCam.far = far;
    this.orthoCam.far = far;
    this.updateProjectionMatrix();
  }

  get near() {
    return this.activeCamera.near;
  }

  get far() {
    return this.activeCamera.far;
  }

  get top() {
    return this.orthoCam.top;
  }

  get bottom() {
    return this.orthoCam.bottom;
  }

  get left() {
    return this.orthoCam.left;
  }

  get right() {
    return this.orthoCam.right;
  }

  get fov() {
    return this.perspectiveCam.fov;
  }

  get isOrthographicCamera() {
    return this.projection === CameraProjection.ORTHOGRAPHIC;
  }

  get isPerspectiveCamera() {
    return this.projection === CameraProjection.PERSPECTIVE;
  }

  static createDefault(): IntrinsicCamera {
    return new IntrinsicCamera();
  }

  setOrthoDimension(scale: number) {
    this.orthoCam.top *= scale;
    this.orthoCam.bottom *= scale;
    this.orthoCam.left *= scale;
    this.orthoCam.right *= scale;
    this.updateProjectionMatrix();
  }

  updateProjectionMatrix() {
    if (this.projection === CameraProjection.PERSPECTIVE) {
      if (this.intrinsics) {
        this.projectionMatrix = this.intrinsicMatrix;
      } else {
        this.perspectiveCam.updateProjectionMatrix();
        this.projectionMatrix = this.perspectiveCam.projectionMatrix;
      }
    } else {
      this.orthoCam.updateProjectionMatrix();
      this.projectionMatrix = this.orthoCam.projectionMatrix;
    }
    this.projectionMatrixInverse.copy(this.projectionMatrix).invert();
    this.updateMatrixWorld();
    this.updatedSubject.next();
  }

  override updateMatrixWorld(force?: boolean) {
    const oldWorldMatrix = this.matrixWorld.clone();
    super.updateMatrixWorld(force);
    if (!this.matrixWorld.equals(oldWorldMatrix)) {
      this.updatedSubject.next();
    }
  }

  override clone(): this {
    const clone = new IntrinsicCamera().copy(this) as this;
    if (this.intrinsics) {
      clone.setIntrinsics(this.intrinsics);
    }
    return clone;
  }
}

/**
 * Creates a frustum mesh centered along the +Z axis
 */
export class CameraFrustum extends three.Mesh {
  private readonly camera = IntrinsicCamera.createDefault();

  constructor() {
    // Compute a box and transform it to be on +Z
    const geometry = new three.BoxGeometry(2, 2, 1).applyMatrix4(
      new three.Matrix4().setPosition(0, 0, 0.5),
    );
    const lines = new three.LineSegments(
      new three.EdgesGeometry(geometry),
      new three.LineBasicMaterial({
        color: colors.PRIMARY_LIGHT,
        linewidth: 1,
        side: three.DoubleSide,
        toneMapped: false,
      }),
    );
    super(
      geometry,
      new three.MeshBasicMaterial({
        color: colors.PRIMARY_LIGHT,
        transparent: true,
        opacity: 0.2,
        side: three.DoubleSide,
      }),
    );
    this.add(lines);
    // Set an arbitrary, aesthetically pleasing near plane distance.
    this.camera.near = 0.025;
    this.camera.far = 2;
    this.quaternion.setFromUnitVectors(
      new three.Vector3(0, 0, -1),
      new three.Vector3(0, 0, 1),
    );
    this.name = 'cameraFrustum';
  }

  // threejs doesn't really support non-affine transforms on objects, so we
  // override the normal update to also apply the projection every update.
  override updateMatrix() {
    super.updateMatrix();
    this.matrix.multiply(this.camera.projectionMatrixInverse);
  }

  /**
   * Set the FoV from camera intrinsics.
   * @param intrinsics The camera intrinsics
   * @param clip The optional distances of the near and far clipping planes
   */
  setFovFromIntrinsics(intrinsics: CameraIntrinsics, clip?: Clip) {
    if (clip) {
      [this.camera.near, this.camera.far] = [clip.near, clip.far];
    }
    this.camera.setIntrinsics(intrinsics);
  }

  override clone(): this {
    const clone = new CameraFrustum() as this;
    if (this.camera.intrinsics) {
      clone.camera.setIntrinsics(this.camera.intrinsics);
    }
    clone.visible = false;

    return clone;
  }
}

/**
 * Projects a depth map into a point cloud
 */
export class DepthPoints extends three.Points {
  override material = new three.PointsMaterial({
    size: 1.0,
    sizeAttenuation: false,
  });
  private width = 0;
  private height = 0;

  constructor() {
    super(new three.BufferGeometry());
  }

  static create(): DepthPoints {
    return new DepthPoints();
  }

  // Don't raycast, unless you want to kill the browser :O
  override raycast() {}

  updatePoints(
    width: number,
    height: number,
    pointsBuffer: Float32Array<ArrayBuffer>,
    intensityBuffer?: Uint8Array<ArrayBuffer>,
  ) {
    if (width !== this.width || height !== this.height) {
      this.geometry.dispose();
      this.geometry = new three.BufferGeometry();
      this.width = width;
      this.height = height;
    }

    this.geometry.setAttribute(
      'position',
      new three.BufferAttribute(pointsBuffer, 3),
    );

    // Compute bounding sphere and box manually ignoring NaNs
    let minX = Infinity,
      minY = Infinity,
      minZ = Infinity;
    let maxX = -Infinity,
      maxY = -Infinity,
      maxZ = -Infinity;

    for (let i = 0; i < pointsBuffer.length; i += 3) {
      const x = pointsBuffer[i];
      const y = pointsBuffer[i + 1];
      const z = pointsBuffer[i + 2];

      if (!isNaN(x) && !isNaN(y) && !isNaN(z)) {
        if (x < minX) minX = x;
        if (y < minY) minY = y;
        if (z < minZ) minZ = z;
        if (x > maxX) maxX = x;
        if (y > maxY) maxY = y;
        if (z > maxZ) maxZ = z;
      }
    }

    if (minX !== Infinity) {
      const center = new three.Vector3(
        (minX + maxX) / 2,
        (minY + maxY) / 2,
        (minZ + maxZ) / 2,
      );
      const radius = center.distanceTo(new three.Vector3(maxX, maxY, maxZ));
      this.geometry.boundingSphere = new three.Sphere(center, radius);
      this.geometry.boundingBox = new three.Box3(
        new three.Vector3(minX, minY, minZ),
        new three.Vector3(maxX, maxY, maxZ),
      );
    } else {
      this.geometry.boundingSphere = new three.Sphere(new three.Vector3(), 0);
      this.geometry.boundingBox = new three.Box3(
        new three.Vector3(),
        new three.Vector3(),
      );
    }

    if (intensityBuffer) {
      this.geometry.setAttribute(
        'color',
        new three.BufferAttribute(
          intensityBuffer,
          intensityBuffer.length / (width * height),
          true,
        ),
      );
      this.material.color = colors.WHITE;
      this.material.vertexColors = true;
    } else {
      this.material.color = colors.PRIMARY_LIGHT;
      this.material.vertexColors = false;
    }
    this.material.needsUpdate = true;
  }

  override clone(): this {
    // deliberately do nothing so that phantoms aren't doubling the number of points in the scene.
    const clone = new DepthPoints() as this;
    clone.visible = false;
    return clone;
  }
}
