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

import {MaterialProperties} from './material_properties';

/**
 * The special name denoting a named geometry is visual geometry. See:
 * geometry_component.proto
 */
export const VISUAL_GEOMETRY_KEY = 'Intrinsic_Visual';
/**
 * The special name denoting a named geometry is collision geometry. See:
 * geometry_component.proto
 */
export const COLLISION_GEOMETRY_KEY = 'Intrinsic_Collision';

/**
 * Types of geometry that may exist on Link entities.
 */
export enum GeometryType {
  Visual = 'Intrinsic_Visual',
  Collision = 'Intrinsic_Collision',
}

/**
 * A primitive shape.
 */
export class Primitive {
  isBox(): this is Box {
    return this instanceof Box;
  }

  isCapsule(): this is Capsule {
    return this instanceof Capsule;
  }

  isCylinder(): this is Cylinder {
    return this instanceof Cylinder;
  }

  isEllipsoid(): this is Ellipsoid {
    return this instanceof Ellipsoid;
  }

  isFrustum(): this is Frustum {
    return this instanceof Frustum;
  }

  isSphere(): this is Sphere {
    return this instanceof Sphere;
  }
}

/**
 * A box primitive shape.
 */
export class Box extends Primitive {
  constructor(public size: three.Vector3) {
    super();
  }
}

/**
 * A capsule primitive shape.
 */
export class Capsule extends Primitive {
  constructor(
    public length: number,
    public radius: number,
  ) {
    super();
  }
}

/**
 * A cylinder primitive shape.
 */
export class Cylinder extends Primitive {
  constructor(
    public length: number,
    public radius: number,
  ) {
    super();
  }
}

/**
 * An ellipsoid primitive shape.
 */
export class Ellipsoid extends Primitive {
  constructor(public radii: three.Vector3) {
    super();
  }
}

/**
 * A frustum primitive shape.
 */
export class Frustum extends Primitive {
  constructor(
    public xAngle: number,
    public yAngle: number,
    public minZDistance: number,
    public maxZDistance: number,
  ) {
    super();
  }
}

/**
 * A sphere primitive shape.
 */
export class Sphere extends Primitive {
  constructor(public radius: number) {
    super();
  }
}

/**
 *  A remote shape is a shape whose geometry and renderable is queryable from CAS.
 */
export declare interface RemoteShape {
  // The reference string to retrieve the exact geometry from CAS.
  exactGeometryRef: string;

  // The reference string to retrieve the renderable glTF data from CAS.
  renderableRef: string;

  // Whether or not the renderable was user-generated. This is generally
  // irrelevant to the frontend except if we do geometry processing, we need
  // to re-inform the backend.
  keepRenderable: boolean;
}

/**
 *  A renderable is a shape whose geometry and renderable is stored in memory.
 */
export declare interface Renderable {
  // For serialization back to Renderable proto.
  getGlbBytes(): Uint8Array<ArrayBuffer>;
}

/**
 *  A transformed primitive shape is a primitive shape with a transformation.
 */
export declare interface TransformedPrimitiveShape {
  primitive: Primitive;
  refTShape: three.Matrix4;
}

/**
 *  A set of transformed primitive shapes.
 */
export declare interface TransformedPrimitiveShapeSet {
  primitives: TransformedPrimitiveShape[];
}

/**
 *  A triangle mesh is a collection of vertices and face indices.
 */
export declare interface TriangleMesh {
  vertices: Float64Array;
  faces: Int32Array;
}

/**
 *  A point cloud is a collection of points and their normals.
 */
export declare interface PointCloud {
  points: Float64Array;
  normals: Float64Array;
}

/**
 *  A local shape stores geometry and rendering data in memory.
 */
export declare interface LocalShape {
  // The computational geometry representation of this shape
  exactGeometry: TransformedPrimitiveShapeSet | TriangleMesh | PointCloud;
  // An optional renderable representation of this shape.
  renderable?: Renderable;
}

/**
 * The different types of shapes that we support to represent the geometry.
 *
 * Currently only remote shapes are supported.
 */
export type GeometryShape = RemoteShape | LocalShape;

/**
 * Returns true if the shape is a remote shape.
 */
export function isRemoteShape(shape: GeometryShape): shape is RemoteShape {
  return 'renderableRef' in shape;
}

/**
 * A logical representation of geometry used by the world. Each geometry has a
 * local transform that defines its offset from the origin of its parent link,
 * and a shape which is used to define its geometry.
 */
export class Geometry {
  constructor(
    private readonly refTShape: three.Matrix4,
    private readonly shape: GeometryShape,
    private readonly materialOverrides?: MaterialProperties,
  ) {}

  // The affine matrix representing the transform from this geometry to the
  // origin of the link that contains it.
  getRefTShape(): three.Matrix4 {
    return this.refTShape;
  }

  // The shape used by this geometry.
  getShape(): GeometryShape {
    return this.shape;
  }

  // Returns true if the shape stored by this geometry needs to be fetched from
  // CAS.
  hasRemoteShape(): this is {getShape: () => RemoteShape} & this {
    return isRemoteShape(this.shape);
  }

  // Returns the optional material property override for this geometry.
  getMaterialOverrides(): MaterialProperties | undefined {
    return this.materialOverrides;
  }
}

/**
 * A specification for removing a geometry with the given name and/or type.
 */
export interface GeometryToRemove {
  /**
   * The type of geometry to remove.
   * If left undefined, will remove both collision and visual geometries matching `geometryName`.
   */
  readonly type?: GeometryType;

  /**
   * The name of the geometry to remove. If undefined or empty, will remove
   * all geometries of the given type.
   */
  readonly geometryName?: string;
}

/**
 * A specification for either updating or creating a geometry with the given
 * name. In other words, the resulting geometry in the geometry component for
 * link will have the properties specified in `geometry`.
 */
export interface GeometryToSet {
  /**
   * The type of geometry to set.
   */
  readonly type: GeometryType;

  /**
   * The geometry's name. May refer to an existing geometry, in which case
   * it will be replaced by that provided in `geometry`.
   */
  readonly geometryName: string;

  /**
   * The geometry to replace for the given name. If left undefined, it will
   * be treated the same as removing this geometry.
   */
  readonly geometry?: Geometry;
}

/**
 * TS equivalent of the GeometryUpdate proto from scene_object_updates.proto.
 *
 * Note that the order of updates is:
 * - First, all `geometriesToRemove` updates are applied
 * - Next, all `geometriesToSet` updates are applied
 */
export class GeometryUpdate {
  constructor(
    readonly geometriesToRemove: GeometryToRemove[] = [],
    readonly geometriesToSet: GeometryToSet[] = [],
  ) {}

  addGeometryToRemove(geometryToRemove: GeometryToRemove): this {
    this.geometriesToRemove.push(geometryToRemove);
    return this;
  }

  addGeometryToSet(geometryToSet: GeometryToSet): this {
    this.geometriesToSet.push(geometryToSet);
    return this;
  }
}
