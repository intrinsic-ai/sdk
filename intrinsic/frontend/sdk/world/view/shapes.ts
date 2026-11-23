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

/**
 * @file Utility functions to create some primitive shapes.
 */
import * as three from 'three';
import {
  MeshBVH,
  acceleratedRaycast,
  computeBoundsTree,
  disposeBoundsTree,
} from 'three-mesh-bvh';

import {MaterialProperties} from '../model/material_properties';
import * as mathUtils from '../utils/math';
import {isMesh, isPoints} from '../utils/three_typechecks';
import * as colors from './colors';
import * as materials from './materials';
import {createCircleSpriteMaterial} from './materials';

const ROBOT_WORK_AREA_HELPER_ALPHA = 0.2;
const ROBOT_WORK_AREA_OUTLINE_STEPS = 64;
const ROBOT_WORK_AREA_NUM_LAT_LON_LINES = 8;
const FRAME_AXIS_LENGTH = 0.04;

/** Material properties used to style collision shapes. */
export interface CollisionShapeMaterialProperties {
  color?: three.Color;
  lineColor?: three.Color;
  opacity?: number;
  lineOpacity?: number;
  transparent?: boolean;
  depthWrite?: boolean;
  wireframe?: boolean;
}

/** Creates a grid with a axes helper marking it's origin. */
export function createGrid(size: number, divisions: number): three.Object3D {
  // NOTE: Threejs uses a right handed system, with Y point up and Z pointing
  // out of the screen.

  const gridHelper = new three.GridHelper(size, divisions);
  // Rotates the grid around X-axis to make it align with the ground plane.
  gridHelper.rotateX(Math.PI / 2);
  // Moves the grid up (along Y-axis) to avoid z-fighting with ground surface at
  // z = 0.
  const DELTA_Z = 0.0001;
  gridHelper.translateY(DELTA_Z);

  const axesHelper = new three.AxesHelper(1);
  axesHelper.setColors(colors.X_AXIS, colors.Y_AXIS, colors.Z_AXIS);
  // Rotates the axes clockwise by 90 deg. (i.e. negative) so that Z points up.
  axesHelper.rotateX(-Math.PI / 2);
  // Avoids z-fighting with the grid.
  axesHelper.renderOrder = 1;
  gridHelper.add(axesHelper);
  gridHelper.updateMatrixWorld(true);
  return gridHelper;
}

/**
 * Create an arrow with given spec.
 * @param baseLength The length of the rod of the arrow.
 * @param tipLength The length of the arrow-shaped tip.
 * @param width The diameter of the rod of the arrow.
 * @param material The material to apply to the arrow.
 */
export function createArrow(
  baseLength: number,
  tipLength: number,
  width: number,
  material: three.Material,
): three.Object3D;

/**
 * Create an arrow with given spec.
 * @param baseLength The length of the rod of the arrow.
 * @param tipLength The length of the arrow-shaped tip.
 * @param width The diameter of the rod of the arrow.
 * @param color The color of the arrow as a hex number.
 * @param alpha If not specified, object will be opaque, otherwise it means the
 *     opacity.
 */
export function createArrow(
  baseLength: number,
  tipLength: number,
  width: number,
  color: number | three.Color,
  alpha?: number,
): three.Object3D;

export function createArrow(
  baseLength: number,
  tipLength: number,
  width: number,
  colorOrMaterial: number | three.Color | three.Material,
  alpha?: number,
): three.Object3D {
  const arrow = new three.Object3D();
  const baseGeometry = new three.CylinderGeometry(
    0.5 * width,
    0.5 * width,
    baseLength,
    16,
  ).rotateX(Math.PI / 2);

  let material: three.Material;
  let color: number | three.Color;

  if (colorOrMaterial instanceof three.Material) {
    material = colorOrMaterial;
    color = (material as three.MeshBasicMaterial).color ?? new three.Color();
  } else {
    color = colorOrMaterial;
    material = new three.MeshBasicMaterial({
      color,
      transparent: true,
      opacity: alpha ?? 1,
      toneMapped: false,
      depthTest: false,
    });
  }

  const base = new three.Mesh(baseGeometry, material);
  const tipGeometry = new three.CylinderGeometry(
    0,
    2 * width,
    tipLength,
    12,
  ).rotateX(Math.PI / 2);
  const tip = new three.Mesh(tipGeometry, material);
  const roundStartPoint = createSphere(
    2 * width,
    new three.MeshBasicMaterial({
      color,
    }),
  );
  base.translateZ(baseLength / 2);
  arrow.add(base);
  tip.translateZ(baseLength);
  tip.renderOrder = Infinity;
  arrow.add(tip);
  roundStartPoint.renderOrder = Infinity;
  arrow.add(roundStartPoint);

  arrow.updateMatrixWorld(true);
  return arrow;
}

/**
 * Create a translucent sphere used to display the work area of a robot.
 * Includes latitude and longitude lines along the surface.
 */
export class RobotWorkAreaHelper extends three.Mesh {
  constructor(radius: number) {
    const sphere = createSphere(
      radius,
      new three.MeshBasicMaterial({
        color: colors.PRIMARY_LIGHT,
        opacity: ROBOT_WORK_AREA_HELPER_ALPHA,
        transparent: true,
      }),
    );
    super(sphere.geometry, sphere.material);

    const lineMaterial = new three.LineBasicMaterial({
      color: colors.PRIMARY_LIGHT,
      transparent: true,
      opacity: 0.5,
      linewidth: 1,
      side: three.DoubleSide,
    });

    // Create a line circle to outline the circumference of the sphere.
    const points = [];
    const extent = Math.PI * 2.0;
    const step = extent / ROBOT_WORK_AREA_OUTLINE_STEPS;

    for (let i = 0; i <= ROBOT_WORK_AREA_OUTLINE_STEPS; i += 1) {
      const a = i * step;
      points.push(
        new three.Vector3(Math.sin(a) * radius, 0, Math.cos(a) * radius),
      );
    }

    // Use it to create longitude lines.
    const geometry = new three.BufferGeometry().setFromPoints(points);
    const lineGeo = new three.Line(geometry, lineMaterial);
    for (let i = 0; i <= ROBOT_WORK_AREA_NUM_LAT_LON_LINES; i++) {
      this.add(
        lineGeo
          .clone()
          .rotateZ((Math.PI / ROBOT_WORK_AREA_NUM_LAT_LON_LINES) * i),
      );
    }

    // Create latitude lines.
    for (let z = 0; z <= ROBOT_WORK_AREA_NUM_LAT_LON_LINES; z++) {
      const latPoints = [];
      for (let i = 0; i <= ROBOT_WORK_AREA_OUTLINE_STEPS; i += 1) {
        const a = i * step;
        const phi = (Math.PI * z) / ROBOT_WORK_AREA_NUM_LAT_LON_LINES;
        latPoints.push(
          new three.Vector3().setFromSphericalCoords(radius, phi, a),
        );
      }
      const geometry = new three.BufferGeometry().setFromPoints(latPoints);
      const lineGeo = new three.Line(geometry, lineMaterial);
      this.add(lineGeo.rotateX(Math.PI / 2));
    }
  }
}

/**
 * A mesh that uses a BVH acceleration structure for raycasting.
 */
export class BVHMesh extends three.Mesh {
  override raycast = acceleratedRaycast;
  constructor(
    geometry: three.BufferGeometry,
    material?: three.Material | three.Material[],
  ) {
    const bvhGeometry = BVHBufferGeometry.createFrom(geometry);
    super(bvhGeometry, material);
  }

  override clone(): this {
    return new BVHMesh(this.geometry, this.material).copy(this) as this;
  }
}

/**
 * Creates a wireframed geometry based on a provided geometry.
 * @param geometry The BufferGeometry used to construct the wireframe
 */
export class CollisionShape extends BVHMesh {
  lines: three.LineSegments | three.Mesh;
  linesMaterial = new three.MeshStandardMaterial({
    color: colors.PRIMARY_LIGHT,
    wireframe: true,
    side: three.DoubleSide,
    transparent: true,
    opacity: 0.15,
    depthTest: true,
    depthWrite: false,
    metalness: 0,
    toneMapped: false,
  });
  override material: three.MeshBasicMaterial;

  constructor(
    geometry: three.BufferGeometry,
    meshColor = colors.PRIMARY_LIGHT,
    lineColor = colors.PRIMARY_LIGHT,
  ) {
    super(geometry);
    this.material = new three.MeshBasicMaterial({
      color: meshColor,
      side: three.DoubleSide,
      transparent: true,
      opacity: 0.24,
      depthTest: true,
      depthWrite: false,
    });

    this.linesMaterial.color.set(lineColor);
    this.lines = new BVHMesh(geometry, this.linesMaterial);
    this.add(this.lines);
  }

  applyMaterialProperties(props: CollisionShapeMaterialProperties) {
    // Updates material properties. If `color` is provided but `lineColor` is omitted, `lineColor`
    // inherits `color`.
    const targetColor = props.color ?? this.material.color;
    const targetLineColor =
      props.lineColor ?? props.color ?? this.linesMaterial.color;
    this.material.color.set(targetColor);
    this.linesMaterial.color.set(targetLineColor);
    this.material.opacity = props.opacity ?? this.material.opacity;
    this.material.transparent = props.transparent ?? this.material.transparent;
    this.material.depthWrite = props.depthWrite ?? this.material.depthWrite;
    this.linesMaterial.opacity =
      props.lineOpacity ?? this.linesMaterial.opacity;
    if (props.wireframe !== undefined) {
      this.lines.visible = props.wireframe;
    }
  }

  markAsColliding(colliding: boolean) {
    if (colliding) {
      this.material.color = new three.Color(1, 0, 0);
      this.linesMaterial.color = new three.Color(1, 0, 0);
    } else {
      this.material.color = colors.PRIMARY_LIGHT;
      this.linesMaterial.color = colors.PRIMARY_LIGHT;
    }
  }

  static applyToSubtree(graph: three.Object3D) {
    // Create a dummy parent to make traversal easier
    if (!graph.parent) {
      const root = new three.Object3D();
      root.add(graph);
    }
    const meshes: three.Mesh[] = [];
    graph.traverse((n: three.Object3D) => {
      if (isMesh(n)) {
        // Don't mutate the scene graph mid-traversal.
        meshes.push(n);
      }
    });

    // Since meshes is built in depth-first order, we need to start from the
    // bottom.
    for (const mesh of meshes.reverse()) {
      const collisionShape = new CollisionShape(mesh.geometry);
      // Swap out the original for the modified version
      mesh.parent?.add(collisionShape);
      mathUtils.updateObjectLocalTransform(collisionShape, mesh.matrix, true);
      mesh.removeFromParent();
    }
  }

  makeInvisible(enabled: boolean) {
    this.material.visible = !enabled;
    this.linesMaterial.visible = !enabled;
  }

  isInvisible(): boolean {
    return !this.material.visible || !this.linesMaterial.visible;
  }

  destroy() {
    this.removeFromParent();
    this.geometry.dispose();
    this.lines.geometry.dispose();
  }
}

/**
 * Helper function that will replace all points or meshes with Intrinsic styled,
 * BVH accelerated versions.
 * @param graph The root of the sub tree you wish to modify
 * @returns The new root of the modified tree.
 */
export function createIntrVisuals(graph: three.Object3D): three.Object3D {
  // Create a dummy parent because the root might change
  const root = new three.Object3D();
  root.matrixAutoUpdate = false;
  root.add(graph);

  const renderables: Array<three.Mesh | three.Points> = [];
  graph.traverse((n: three.Object3D) => {
    if (isMesh(n) || isPoints(n)) {
      // Don't mutate the scene graph mid-traversal.
      renderables.push(n);
    }
  });

  // Since meshes is built in depth-first order, we need to start from the
  // bottom.
  for (const renderable of renderables.reverse()) {
    // Make sure renderable has an updated matrix before we copy it.
    renderable.updateMatrix();
    const patched = isMesh(renderable)
      ? new IntrVisualMesh(renderable.geometry, renderable.material)
      : new IntrVisualPoints(renderable.geometry, renderable.material);
    renderable.parent?.add(patched);
    patched.applyMatrix4(renderable.matrix);
    renderable.removeFromParent();
    renderable.geometry.dispose();
  }
  // Return the root we may have just swapped out.
  return root.children[0];
}

/**
 * Custom Mesh class that encapsulates custom rendering logic and BVH
 * acceleration.
 */
export class IntrVisualMesh extends three.Mesh {
  override raycast = acceleratedRaycast;
  override material: materials.IntrVisualMaterial;
  constructor(
    geometry: three.BufferGeometry,
    material?: three.Material | three.Material[],
  ) {
    const bvh = BVHBufferGeometry.createFrom(geometry);
    super(bvh);
    const m = Array.isArray(material) ? material[0] : material;
    this.material = new materials.IntrVisualMaterial();
    if (m) {
      this.material.copy(m);
    }
    // Compute normals in the shader if none are present
    if (!this.geometry.hasAttribute('normal')) {
      this.material.flatShading = true;
    }
    // We default to a near-infinite attenuation distance. We can't directly
    // use "Infinity" as the glTF exporter converts that to zero for some
    // reason.
    this.material.attenuationDistance = 1e5;
  }

  override clone(): this {
    const mat = this.material.clone();
    const cloned = new IntrVisualMesh(this.geometry).copy(this) as this;
    cloned.material = mat;
    return cloned;
  }

  updateMaterialProperties(materialProperties: MaterialProperties) {
    // NOTE - it is important that we don't update the `flatShading`
    // field of the material here. It is set correctly in the constructor,
    // based on whether the geometry has normals or not.
    this.material.color = new three.Color(materialProperties.color);
    this.material.metalness = materialProperties.metalness;
    this.material.roughness = materialProperties.roughness;
    this.material.transmission = materialProperties.transmission;
  }

  enableHighlight(enabled: boolean) {
    this.material.enableHighlight(enabled);
  }

  enableColorHighlight(enabled: three.Color | false) {
    this.material.enableColorHighlight(enabled);
  }

  enableDiffHighlight(enabled: boolean) {
    this.material.enableDiffHighlight(enabled);
  }

  makeInvisible(enabled: boolean) {
    this.material.makeInvisible(enabled);
  }

  isInvisible(): boolean {
    return this.material.isInvisible();
  }
}

/**
 * Custom Points class that encapsulates custom rendering logic and BVH
 * acceleration.
 */
export class IntrVisualPoints extends three.Points {
  override raycast = acceleratedRaycast;
  override geometry: BVHBufferGeometry;
  override material = new materials.IntrVisualPointsMaterial();

  constructor(
    geometry: three.BufferGeometry,
    material?: three.Material | three.Material[],
  ) {
    const bvh = BVHBufferGeometry.createFrom(geometry);
    super(bvh);
    if (material instanceof three.PointsMaterial) {
      this.material.copy(material);
    }
    this.geometry = bvh;
  }

  override clone(): this {
    const clone = new IntrVisualPoints(this.geometry).copy(this);
    // Ensure materials have unique references
    clone.material = this.material.clone();
    return clone as this;
  }

  enableHighlight(enabled: boolean) {
    this.material.enableHighlight(enabled);
  }

  isInvisible(): boolean {
    return false;
  }
}

/**
 * An extended form of BufferGeometry that also builds an acceleration
 * structure that drastically improves ray casts.
 */
export class BVHBufferGeometry extends three.BufferGeometry {
  override readonly computeBoundsTree = computeBoundsTree;
  override readonly disposeBoundsTree = disposeBoundsTree;
  override boundsTree?: MeshBVH;

  static createFrom(geometry: three.BufferGeometry): BVHBufferGeometry {
    const patched =
      geometry instanceof BVHBufferGeometry
        ? geometry
        : new BVHBufferGeometry().copy(geometry);
    if (!patched.boundsTree) {
      patched.computeBoundsTree();
    }
    return patched;
  }

  override clone(): this {
    const cloned = new BVHBufferGeometry().copy(this);
    cloned.boundsTree = this.boundsTree;
    return cloned as this;
  }
}

class MarginMaterial extends three.LineBasicMaterial {
  shader?: three.ShaderLibShader;
  margin = 0;

  constructor(args: ConstructorParameters<typeof three.LineBasicMaterial>[0]) {
    super(args);
    this.onBeforeCompile = (shader: three.ShaderLibShader) => {
      this.patchShader(shader);
    };
  }

  setMargin(margin: number) {
    this.margin = margin;
    if (this.shader) {
      this.shader.uniforms['uMargin'].value = this.margin;
    }
  }

  // Patch the shader to avoid slowly recomputing the entire vertex buffer (host
  // side) every time we want to update the margin; instead simply let the gpu
  // do the extrusion. The idomatic way to extending materials in threejs
  // _appears_ to be string replacement of the templated shader:
  // https://threejs.org/docs/#api/en/materials/Material.onBeforeCompile
  private patchShader(shader: three.ShaderLibShader) {
    this.shader = shader;
    this.shader.uniforms['uMargin'] = {value: this.margin};

    // Append our margin uniform declaration
    this.shader.vertexShader = this.shader.vertexShader.replace(
      '#include <common>',
      ['#include <common>', 'uniform float uMargin;'].join('\n'),
    );

    // Because each line segment is a pair of points, we can use the vertex id
    // to conditionally extrude every other vertex along the normal by a
    // distance of <margin>.
    this.shader.vertexShader = this.shader.vertexShader.replace(
      '#include <begin_vertex>',
      [
        'vec3 transformed = (gl_VertexID % 2 == 1) ? vec3(position) + normal * max(uMargin, 0.0) : vec3(position);',
      ].join('\n'),
    );
  }
}

/**
 * Creates a mesh with extruded lines emanating from each triangle's barycenter.
 * @param mesh The mesh to show the margins of.
 */
export class MarginShape extends three.LineSegments {
  override material = new MarginMaterial({color: colors.PRIMARY_LIGHT});

  constructor(mesh: three.Mesh) {
    const margins = MarginShape.createMarginLines(mesh.geometry);
    super(margins);
  }

  private static createMarginLines(
    geometry: three.BufferGeometry,
  ): three.BufferGeometry {
    // After threejs r112, GL_TRIANGLES is the only supported topology.
    // Since we have triangles, make them un-indexed (i.e. no shared vertices)
    // and recompute the normals. This should give us flat shaded normals.
    const triangleSoup = geometry.toNonIndexed();
    triangleSoup.computeVertexNormals();
    const posAttr = triangleSoup.getAttribute('position');
    const normAttr = triangleSoup.getAttribute('normal');
    if (!posAttr || !normAttr || posAttr.count !== normAttr.count) {
      throw new Error('Invalid normals');
    }

    const barycenters: number[] = [];
    const normals: number[] = [];

    for (let i = 0; i < posAttr.count; i += 3) {
      // Compute the barycenter to extrude the line from
      const bary = new three.Vector3();
      bary.add(new three.Vector3().fromBufferAttribute(posAttr, i));
      bary.add(new three.Vector3().fromBufferAttribute(posAttr, i + 1));
      bary.add(new three.Vector3().fromBufferAttribute(posAttr, i + 2));
      bary.multiplyScalar(1 / 3);
      // Pick only the first normal as all three should be the same for each
      // triangle
      const normal = new three.Vector3().fromBufferAttribute(normAttr, i);
      // Place duplicate vertices at each barycenter, we will extrude the second
      // vertex of each pair.
      normals.push(normal.x, normal.y, normal.z, normal.x, normal.y, normal.z);
      barycenters.push(bary.x, bary.y, bary.z, bary.x, bary.y, bary.z);
    }

    const segments = new three.BufferGeometry();
    segments.setAttribute(
      'position',
      new three.BufferAttribute(new Float32Array(barycenters), 3),
    );
    segments.setAttribute(
      'normal',
      new three.BufferAttribute(new Float32Array(normals), 3),
    );

    return segments;
  }

  setMargin(margin: number) {
    this.material.setMargin(margin);
  }
}

/**
 * A base class for scene helper objects.
 */
export abstract class SceneHelper extends three.Group {
  abstract hoverText: string;
  abstract enableHighlight(enabled: boolean): void;
  isInvisible(): boolean {
    return false;
  }
}

/**
 * Create a hemisphere with given spec.
 * @param radius The radius of the sphere.
 *
 */
export function createHemisphere(radius: number): three.Mesh {
  const geom = new three.SphereGeometry(
    radius,
    32,
    32,
    0,
    Math.PI * 2,
    0,
    Math.PI / 2,
  );
  const mat = new three.MeshBasicMaterial({
    color: colors.PRIMARY_LIGHT,
    transparent: true,
    opacity: 0.5,
    side: three.DoubleSide,
  });
  return new three.Mesh(geom, mat);
}

/**
 * Non-default settings for creating a three.SphereGeometry.
 */
export interface SphereGeometryOptions {
  widthSegments?: number;
  heightSegments?: number;
  phiStart?: number;
  phiLength?: number;
  thetaStart?: number;
  thetaLength?: number;
}

/**
 * Create a ellipsoid with given spec.
 * @param radii The radii of the ellipsoid.
 * @param material The optional material applied to the mesh.
 * @param geometryOptions The optional sphere geometry parameters applied to the
 * mesh.
 */
export function createEllipsoid(
  radii: three.Vector3,
  material?: three.Material,
  geometryOptions?: SphereGeometryOptions,
): three.Mesh {
  const geometry = new three.SphereGeometry(
    /*radius=*/ 1.0,
    geometryOptions?.widthSegments ?? 32,
    geometryOptions?.heightSegments ?? 32,
    geometryOptions?.phiStart ?? 0,
    geometryOptions?.phiLength ?? Math.PI * 2,
    geometryOptions?.thetaStart ?? 0,
    geometryOptions?.thetaLength ?? Math.PI * 2,
  );
  geometry.scale(radii.x, radii.y, radii.z);
  return new three.Mesh(geometry, material);
}

/**
 * Create a sphere with given spec.
 * @param radius The radius of the sphere.
 * @param material The optional material applied to the mesh.
 * @param geometryOptions The optional sphere geometry parameters applied to the
 * mesh.
 */
export function createSphere(
  radius: number,
  material?: three.Material,
  geometryOptions?: SphereGeometryOptions,
): three.Mesh {
  return createEllipsoid(
    new three.Vector3(radius, radius, radius),
    material,
    geometryOptions,
  );
}

/**
 * Create a cylinder (centered at the origin) with given spec.
 * @param radius The radius of the sphere.
 * @param length The length of the cylinder.
 * @param material The optional material applied to the mesh.
 */
export function createCylinder(
  radius: number,
  length: number,
  material?: three.Material,
): three.Mesh {
  // Convert from one "up" convention to the one used in the geo service
  const quat = new three.Quaternion().setFromUnitVectors(
    new three.Vector3(0, 1, 0),
    new three.Vector3(0, 0, 1),
  );
  const geometry = new three.CylinderGeometry(
    radius,
    radius,
    length,
    64,
    1,
  ).applyQuaternion(quat);
  return new three.Mesh(geometry, material);
}

/**
 * Create a cylinder (centered at the origin) with given spec.
 * @param dimensions Vector containing the desired dimensions fo the box.
 * @param material The optional material applied to the mesh.
 */
export function createBox(
  dimensions: three.Vector3,
  material?: three.Material,
): three.Mesh {
  const geometry = new three.BoxGeometry(
    dimensions.x,
    dimensions.y,
    dimensions.z,
  );
  return new three.Mesh(geometry, material);
}

/**
 * A frame represented by three axes in red/x, green/y, blue/z.
 * Also adds a scaled up, invisible representation of each axis to make
 * click-picking of frames easier.
 */
export class FrameAxes extends three.Object3D {
  // A static instance that wraps a singleton of the actual GPU data that avoids re-instantiating GPU resources.
  private static readonly inst = FrameAxes.init();

  // FrameAxes should not be constructed directly, but instead created by cloning.
  private constructor() {
    super();
  }

  override clone(): this {
    const clone = super.clone();
    // Ensure material objects are unique.
    clone.traverse((n: three.Object3D) => {
      if (isMesh(n)) {
        if (Array.isArray(n.material)) {
          n.material = n.material.map((m) => m.clone());
        } else {
          n.material = n.material.clone();
        }
      }
    });
    return clone;
  }

  // Use separate init logic so that cloning doesn't re-create GPU resources in the constructor.
  private static init(): FrameAxes {
    const inst = new FrameAxes();
    const diameter = 0.05 * FRAME_AXIS_LENGTH;
    const red = createCapsule(FRAME_AXIS_LENGTH, diameter, colors.X_AXIS);
    const green = createCapsule(FRAME_AXIS_LENGTH, diameter, colors.Y_AXIS);
    const blue = createCapsule(FRAME_AXIS_LENGTH, diameter, colors.Z_AXIS);

    const pointAtX = new three.Vector3(1, 0, 0);
    red.lookAt(pointAtX);
    red.updateMatrixWorld(true);

    const pointAtY = new three.Vector3(0, 1, 0);
    green.lookAt(pointAtY);
    green.updateMatrixWorld(true);

    inst.add(red, green, blue);

    // Create the collision geometry
    const collX = createCapsule(FRAME_AXIS_LENGTH, diameter * 3, colors.X_AXIS);
    const collY = createCapsule(FRAME_AXIS_LENGTH, diameter * 3, colors.Y_AXIS);
    const collZ = createCapsule(FRAME_AXIS_LENGTH, diameter * 3, colors.Z_AXIS);
    collX.lookAt(pointAtX);
    collX.updateMatrixWorld(true);
    collY.lookAt(pointAtY);
    collY.updateMatrixWorld(true);
    collX.visible = false;
    collY.visible = false;
    collZ.visible = false;
    inst.add(collX, collY, collZ);
    return inst;
  }

  enableHighlight(enabled: boolean) {
    this.traverse((n: three.Object3D) => {
      if (isMesh(n) && n.material instanceof materials.IntrBasicMaterial) {
        n.material.enableHighlight(enabled);
      }
    });
  }

  enableDimming(disabled: boolean) {
    this.traverse((n: three.Object3D) => {
      if (isMesh(n) && n.material instanceof materials.IntrBasicMaterial) {
        n.material.enableDimming(disabled);
      }
    });
  }

  makeInvisible(enabled: boolean) {
    this.traverse((n: three.Object3D) => {
      if (isMesh(n) && n.material instanceof materials.IntrBasicMaterial) {
        n.material.makeInvisible(enabled);
      }
    });
  }

  isInvisible(): boolean {
    let invisible = false;
    this.traverse((n: three.Object3D) => {
      if (isMesh(n) && n.material instanceof materials.IntrBasicMaterial) {
        invisible ||= n.material.isInvisible();
      }
    });
    return invisible;
  }

  static create(): FrameAxes {
    // Cloning prevents allocating/uploading new vertex data to the GPU, but instead re-uses existing buffers.
    return FrameAxes.inst.clone();
  }
}

/**
 * Create a ring geometry with given spec.
 * @param radius The radius of the ring.
 * @param segments The number of segments in the ring.
 */
export function createRingGeometry(
  radius: number,
  segments: number,
): three.BufferGeometry {
  const geometry = new three.BufferGeometry();
  const vertices = [];
  const theta = (2 * Math.PI) / segments;
  for (let i = 0; i < segments; i++) {
    const x = radius * Math.cos(i * theta);
    const y = radius * Math.sin(i * theta);
    const z = 0;
    vertices.push(x, y, z);
  }
  geometry.setAttribute(
    'position',
    new three.Float32BufferAttribute(vertices, 3),
  );
  return geometry;
}

/**
 * Creates a frustum BufferGeometry given the frustum parameters.
 * The frustum is centered at the origin and extends along the positive Z axis.
 * @param minZ The minimum Z distance of the frustum.
 * @param maxZ The maximum Z distance of the frustum.
 * @param xAngle The x axis field of view angle of the frustum in radians.
 * @param yAngle The y axis field of view angle of the frustum in radians.
 */
export function frustumBufferGeometry(
  minZ: number,
  maxZ: number,
  xAngle: number,
  yAngle: number,
): three.BufferGeometry {
  const tanX = Math.tan(xAngle * 0.5);
  const tanY = Math.tan(yAngle * 0.5);
  const minX = minZ * tanX;
  const maxX = maxZ * tanX;
  const minY = minZ * tanY;
  const maxY = maxZ * tanY;
  const vertices = new Float32Array([
    // Min plane
    -minX,
    -minY,
    minZ,
    minX,
    -minY,
    minZ,
    minX,
    minY,
    minZ,
    -minX,
    minY,
    minZ,
    // Max plane
    -maxX,
    -maxY,
    maxZ,
    maxX,
    -maxY,
    maxZ,
    maxX,
    maxY,
    maxZ,
    -maxX,
    maxY,
    maxZ,
  ]);
  const indices = new Uint16Array([
    // MinZ face
    0, 1, 2, 0, 2, 3,
    // MaxZ face
    4, 7, 6, 4, 6, 5,
    // Sides
    0, 1, 5, 0, 5, 4, 1, 2, 6, 1, 6, 5, 2, 3, 7, 2, 7, 6, 3, 0, 4, 3, 4, 7,
  ]);
  const frustumGeometry = new three.BufferGeometry();
  frustumGeometry.setIndex(new three.BufferAttribute(indices, 1));
  frustumGeometry.setAttribute(
    'position',
    new three.BufferAttribute(vertices, 3),
  );
  frustumGeometry.computeVertexNormals();
  return frustumGeometry;
}

/**
 * Create a capsule with given spec. The capsule's central axis is aligned with the z Axis.
 * @param length The length of the rod of the capsule.
 * @param diameter The diameter of the rod of the capsule.
 * @param color The color of the capsule as a threejs Color object.
 * @param alpha If not specified, object will be opaque, otherwise it means
 *     the opacity.
 */
export function createCapsule(
  length: number,
  diameter: number,
  color: three.Color,
  alpha?: number,
): three.Object3D {
  const capsule = new three.Object3D();
  const baseGeometry = new three.CapsuleGeometry(
    0.5 * diameter,
    length,
    6,
    16,
  ).rotateX(Math.PI / 2);
  const material = new materials.IntrBasicMaterial({
    color,
    transparent: !!alpha,
    opacity: alpha ?? 1,
    toneMapped: false,
  });
  const base = new three.Mesh(baseGeometry, material);
  base.translateZ(length / 2);
  capsule.add(base);
  capsule.updateMatrixWorld(true);
  return capsule;
}

/**
 * Create a torus object formed by 4 subobjects of the geometry type torus
 * occupying PI/2 of the circumference. Subobjects will alternate between full
 * opacity and an alternating opacity value passed in as input. To indicate
 * the rotation direction, there is 2 arrows in opposite sides of the circle.
 * @param radius The diameter of the torus.
 * @param width The diameter of the cylinder.
 * @param color The color of the torus as a hex number.
 * @param alpha If not specified, object will be opaque, otherwise it means
 *     the opacity.
 * @param alternateAlpha If specified, alternate rings will have this opacity
 *     value instead
 * @param toneMapped Whether to render the object with the scene's default
 *     color scheme if true or with a flat shading if false
 */
export function createPartialRotationRing(
  radius: number,
  width: number,
  color: number | three.Color,
  alpha?: number,
  alternateAlpha?: number,
  toneMapped?: boolean,
): three.Object3D {
  const rotRing = new three.Object3D();
  const ringGeo1 = new three.TorusGeometry(
    0.5 * radius,
    0.5 * width,
    16,
    32,
    Math.PI / 2,
  );
  const ringGeo2 = ringGeo1.clone().rotateZ(Math.PI / 2);
  const ringGeo3 = ringGeo2.clone().rotateZ(Math.PI / 2);
  const ringGeo4 = ringGeo3.clone().rotateZ(Math.PI / 2);

  const material = new three.MeshBasicMaterial({
    color,
    transparent: !!alpha,
    opacity: alpha ?? 1,
    toneMapped,
    depthTest: false,
  });
  // Avoid z-fighting by disabling depth for the ring
  const alternateRingMaterial = material.clone();
  alternateRingMaterial.setValues({
    opacity: alternateAlpha ?? 1,
  });
  const ring1 = new three.Mesh(ringGeo1, material);
  const ring2 = new three.Mesh(ringGeo2, alternateRingMaterial);
  const ring3 = new three.Mesh(ringGeo3, material);
  const ring4 = new three.Mesh(ringGeo4, alternateRingMaterial);
  const originPoint = createSphere(2 * width, material);
  const arrowGeometry = new three.CylinderGeometry(0, 2 * width, 5 * width, 32);
  const arrow1 = new three.Mesh(arrowGeometry, material).rotateZ(Math.PI / 2);
  arrow1.translateX(radius / 2);
  arrow1.renderOrder = Infinity;
  const arrow2 = new three.Mesh(arrowGeometry, material).rotateZ(-Math.PI / 2);
  arrow2.translateX(radius / 2);
  arrow2.renderOrder = Infinity;

  rotRing.add(arrow1);
  rotRing.add(arrow2);
  rotRing.add(ring1);
  rotRing.add(ring2);
  rotRing.add(ring3);
  rotRing.add(ring4);
  rotRing.add(originPoint);
  rotRing.updateMatrixWorld(true);
  return rotRing;
}

/**
 * Create a line with a circle handle at start at end points.
 * @param startPoint 3d world space location where the line starts.
 * @param endPoint 3d world space location where the line ends.
 * @param color Color of line and handles as a hex number.
 * @param lineWidth Width of the line in pixels.
 * @param handleScale Scale of the handles. These are created as threejs
 *     sprites
 * which have a default size of 1.0.
 * @returns threejs objects for line, start point sprite, and end point sprite.
 */
export function createLineWithHandles(
  startPoint: three.Vector3,
  endPoint: three.Vector3,
  color: number,
  lineWidth: number,
  handleScale: number,
): [three.Line, three.Sprite, three.Sprite] {
  const colorObj = new three.Color(color);
  const lineGeometry = new three.BufferGeometry().setFromPoints([
    startPoint,
    endPoint,
  ]);
  const lineMaterial = new three.LineBasicMaterial({
    color: colorObj,
    linewidth: lineWidth,
    depthTest: false,
    depthWrite: true,
  });

  const line = new three.Line(lineGeometry, lineMaterial);
  const startHandle = new three.Sprite(
    createCircleSpriteMaterial(colorObj, undefined, false, false, false),
  );
  startHandle.scale.setScalar(handleScale);
  const endHandle = new three.Sprite(
    createCircleSpriteMaterial(colorObj, undefined, false, false, false),
  );
  endHandle.scale.setScalar(handleScale);

  line.renderOrder = Infinity;
  line.add(startHandle);
  line.add(endHandle);

  return [line, startHandle, endHandle];
}

/**
 * Create two opposite facing lines with given spec.
 * @param length The length of the rod of the lines.
 * @param width The diameter of the rod of the lines.
 * @param alpha If not specified, object will be opaque, otherwise it means
 *     the opacity.
 */
export function createOppositeLines(
  length: number,
  width: number,
  alpha?: number,
): three.Object3D {
  const lineGroup = new three.Object3D();
  const lineGeometry1 = new three.CylinderGeometry(
    0.5 * width,
    0.5 * width,
    length,
    16,
  ).rotateX(Math.PI / 2);
  const lineGeometry2 = new three.CylinderGeometry(
    0.5 * width,
    0.5 * width,
    -length,
    16,
  ).rotateX(Math.PI / 2);
  const lineMaterial = new three.MeshBasicMaterial({
    transparent: !!alpha,
    opacity: alpha ?? 1,
    depthTest: false,
    toneMapped: false,
  });

  const line1 = new three.Mesh(lineGeometry1, lineMaterial);
  line1.translateZ(length / 2);
  const line2 = new three.Mesh(lineGeometry2, lineMaterial);
  line2.translateZ(-length / 2);

  lineGroup.add(line1);
  lineGroup.add(line2);

  return lineGroup;
}

/**
 * Apply a material to a threejs node.
 */
export function applyMaterialToNode(
  node: three.Object3D,
  materialProperties: MaterialProperties,
) {
  node.traverse((child: three.Object3D) => {
    if (child instanceof IntrVisualMesh) {
      child.updateMaterialProperties(materialProperties);
    }
  });
}
