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
 * @file Conversions between native geometries and protobuf-es geometry protos.
 */
import {type MessageInitShape, create} from '@bufbuild/protobuf';
import type {ExactGeometrySchema} from '@intrinsic-ai/types/geometry/proto/v1/exact_geometry_pb';
import {
  type GeometricTransform as GeometricTransformProto,
  GeometricTransformSchema,
} from '@intrinsic-ai/types/geometry/proto/v1/geometric_transform_pb';
import {
  type Geometry as GeometryProto,
  GeometrySchema,
} from '@intrinsic-ai/types/geometry/proto/v1/geometry_pb';
import {
  type GeometryStorageRefs as GeometryStorageRefsProto,
  GeometryStorageRefsSchema,
} from '@intrinsic-ai/types/geometry/proto/v1/geometry_storage_refs_pb';
import {
  type InlineGeometry as InlineGeometryProto,
  InlineGeometrySchema,
} from '@intrinsic-ai/types/geometry/proto/v1/inline_geometry_pb';
import {
  type MaterialProperties as MaterialPropertiesProto,
  MaterialPropertiesSchema,
} from '@intrinsic-ai/types/geometry/proto/v1/material_pb';
import type {
  PrimitiveShape as PrimitiveShapeProto,
  PrimitiveShapeSchema,
} from '@intrinsic-ai/types/geometry/proto/v1/primitive_shape_pb';
import {
  type TransformedGeometry as TransformedGeometryProto,
  TransformedGeometrySchema,
} from '@intrinsic-ai/types/geometry/proto/v1/transformed_geometry_pb';
import {
  type TransformedPrimitiveShape as TransformedPrimitiveShapeProto,
  TransformedPrimitiveShapeSchema,
} from '@intrinsic-ai/types/geometry/proto/v1/transformed_primitive_shape_pb';
import {
  type Matrixd as MatrixdProto,
  MatrixdSchema,
} from '@intrinsic-ai/types/math/proto/matrix_pb';
import type {Vector3 as Vector3Proto} from '@intrinsic-ai/types/math/proto/vector3_pb';
import {
  GeometryUpdate_GeometryType as GeometryTypeProto,
  type GeometryUpdate as GeometryUpdateProto,
  GeometryUpdateSchema,
} from '@intrinsic-ai/types/scene/proto/v1/scene_object_updates_pb';
import * as three from 'three';

import {
  Box,
  Capsule,
  Cylinder,
  Ellipsoid,
  Frustum,
  Geometry,
  type GeometryShape,
  GeometryType,
  type GeometryUpdate,
  type LocalShape,
  type PointCloud,
  type Primitive,
  type RemoteShape,
  type Renderable,
  Sphere,
  type TransformedPrimitiveShape,
  type TransformedPrimitiveShapeSet,
  type TriangleMesh,
} from '../../model/geometry';
import type {MaterialProperties} from '../../model/material_properties';

/**
 * A renderable that wraps a glb bytes buffer.
 */
export class GltfRenderable implements Renderable {
  constructor(private readonly glbBytes: Uint8Array<ArrayBuffer>) {}

  getGlbBytes(): Uint8Array<ArrayBuffer> {
    return this.glbBytes;
  }
}

/**
 * Converts a 4x4 Matrixd proto into a transform matrix. Throws if the proto
 * doesn't have 4x4 values.
 */
export function matrixdProtoToMatrix(proto: MatrixdProto): three.Matrix4 {
  if (proto.rows !== 4 || proto.cols !== 4) {
    throw new Error(
      `Unexpected matrix dimensions: ${proto.cols} x ${proto.rows}`,
    );
  }
  if (proto.values.length !== 16) {
    throw new Error(
      `Unexpected number of elements. Expected 16, got ${proto.values.length}.`,
    );
  }
  return new three.Matrix4().fromArray(proto.values);
}

/**
 * Converts a 3x3 Matrixd proto into a matrix. Throws if the proto doesn't have
 * 3x3 values.
 */
export function matrixdProtoToMatrix3(proto: MatrixdProto): three.Matrix3 {
  if (proto.rows !== 3 || proto.cols !== 3) {
    throw new Error(
      `Unexpected matrix dimensions: ${proto.cols} x ${proto.rows}`,
    );
  }
  if (proto.values.length !== 9) {
    throw new Error(
      `Unexpected number of elements. Expected 9, got ${proto.values.length}.`,
    );
  }
  return new three.Matrix3().fromArray(proto.values);
}

/** Converts a transform matrix into a 4x4 Matrixd proto. */
export function matrixToMatrixdProto(mat: three.Matrix4): MatrixdProto {
  return create(MatrixdSchema, {rows: 4, cols: 4, values: [...mat.elements]});
}

/** Converts a GeometricTransform proto into a transform matrix. */
export function geometricTransformProtoToMatrix(
  proto: GeometricTransformProto,
): three.Matrix4 {
  switch (proto.data.case) {
    case 'matrix4d':
      return matrixdProtoToMatrix(proto.data.value);
    case 'trs': {
      const trs = proto.data.value;
      const translation = vector3ProtoToVector(trs.translation);
      const rotation = new three.Quaternion();
      const rpy = trs.rotationRpy;
      if (rpy) {
        rotation.setFromEuler(new three.Euler(rpy.r, rpy.p, rpy.y, 'ZYX'));
      }
      const scale = trs.scale
        ? vector3ProtoToVector(trs.scale)
        : new three.Vector3(1, 1, 1);
      return new three.Matrix4().compose(translation, rotation, scale);
    }
    default:
      throw new Error(
        `Unsupported geometric transform data case: ${proto.data.case}`,
      );
  }
}

/** Converts a transform matrix into a GeometricTransform proto. */
export function matrixToGeometricTransformProto(
  mat: three.Matrix4,
): GeometricTransformProto {
  return create(GeometricTransformSchema, {
    data: {case: 'matrix4d', value: matrixToMatrixdProto(mat)},
  });
}

function vector3ProtoToVector(proto: Vector3Proto | undefined): three.Vector3 {
  return new three.Vector3(proto?.x ?? 0, proto?.y ?? 0, proto?.z ?? 0);
}

function vectorToVector3(vector: three.Vector3): {
  x: number;
  y: number;
  z: number;
} {
  return {x: vector.x, y: vector.y, z: vector.z};
}

/** Converts a GeometryStorageRefs proto into a remote shape. */
export function geometryRefsProtoToRemoteShape(
  proto: GeometryStorageRefsProto,
): RemoteShape {
  return {
    renderableRef: proto.renderableRef,
    exactGeometryRef: proto.exactGeometryRef,
    keepRenderable: proto.keepRenderable,
  };
}

/** Converts a remote shape into a GeometryStorageRefs proto. */
export function remoteShapeToGeometryRefsProto(
  shape: RemoteShape,
): GeometryStorageRefsProto {
  return create(GeometryStorageRefsSchema, {
    renderableRef: shape.renderableRef,
    exactGeometryRef: shape.exactGeometryRef ?? '',
    keepRenderable: shape.keepRenderable,
  });
}

function primitiveToProto(
  primitive: Primitive,
): MessageInitShape<typeof PrimitiveShapeSchema>['shape'] {
  if (primitive.isBox()) {
    return {case: 'box', value: {size: vectorToVector3(primitive.size)}};
  }
  if (primitive.isCapsule()) {
    return {
      case: 'capsule',
      value: {length: primitive.length, radius: primitive.radius},
    };
  }
  if (primitive.isCylinder()) {
    return {
      case: 'cylinder',
      value: {length: primitive.length, radius: primitive.radius},
    };
  }
  if (primitive.isEllipsoid()) {
    return {
      case: 'ellipsoid',
      value: {radii: vectorToVector3(primitive.radii)},
    };
  }
  if (primitive.isFrustum()) {
    return {
      case: 'frustum',
      value: {
        xAngle: primitive.xAngle,
        yAngle: primitive.yAngle,
        minZDistance: primitive.minZDistance,
        maxZDistance: primitive.maxZDistance,
      },
    };
  }
  if (primitive.isSphere()) {
    return {case: 'sphere', value: {radius: primitive.radius}};
  }
  throw new Error(`Primitive ${primitive} can't be converted to a proto!`);
}

function primitiveFromProto(proto: PrimitiveShapeProto | undefined): Primitive {
  const shape = proto?.shape;
  switch (shape?.case) {
    case 'box':
      return new Box(vector3ProtoToVector(shape.value.size));
    case 'capsule':
      return new Capsule(shape.value.length, shape.value.radius);
    case 'cylinder':
      return new Cylinder(shape.value.length, shape.value.radius);
    case 'ellipsoid':
      return new Ellipsoid(vector3ProtoToVector(shape.value.radii));
    case 'frustum':
      return new Frustum(
        shape.value.xAngle,
        shape.value.yAngle,
        shape.value.minZDistance,
        shape.value.maxZDistance,
      );
    case 'sphere':
      return new Sphere(shape.value.radius);
    default:
      throw new Error(
        `Can't convert TransformedPrimitiveShape proto to TransformedPrimitiveShape, shape is not set.`,
      );
  }
}

function transformedPrimitiveShapeToProto(
  shape: TransformedPrimitiveShape,
): TransformedPrimitiveShapeProto {
  return create(TransformedPrimitiveShapeSchema, {
    shape: {shape: primitiveToProto(shape.primitive)},
    refTShape: matrixToGeometricTransformProto(shape.refTShape),
  });
}

function transformedPrimitiveShapeFromProto(
  proto: TransformedPrimitiveShapeProto,
): TransformedPrimitiveShape {
  return {
    primitive: primitiveFromProto(proto.shape),
    refTShape: proto.refTShape
      ? geometricTransformProtoToMatrix(proto.refTShape)
      : new three.Matrix4(),
  };
}

function exactGeometryToProto(
  exactGeometry: LocalShape['exactGeometry'],
): MessageInitShape<typeof ExactGeometrySchema>['data'] {
  if ('primitives' in exactGeometry) {
    return {
      case: 'primitiveSet',
      value: {
        primitives: exactGeometry.primitives.map(
          transformedPrimitiveShapeToProto,
        ),
      },
    };
  }
  if ('vertices' in exactGeometry) {
    return {
      case: 'triangleMesh',
      value: {
        vertices: Array.from(exactGeometry.vertices),
        faces: Array.from(exactGeometry.faces),
      },
    };
  }
  if ('points' in exactGeometry) {
    return {
      case: 'pointCloud',
      value: {
        points: Array.from(exactGeometry.points),
        normals: Array.from(exactGeometry.normals),
      },
    };
  }
  throw new Error(
    `Exact geometry ${exactGeometry} can't be converted to a proto!`,
  );
}

/** Converts a local shape into an InlineGeometry proto. */
export function localShapeToInlineGeometryProto(
  shape: LocalShape,
): InlineGeometryProto {
  return create(InlineGeometrySchema, {
    exactGeometry: {data: exactGeometryToProto(shape.exactGeometry)},
    renderableType: shape.renderable
      ? {case: 'renderable', value: {glbBytes: shape.renderable.getGlbBytes()}}
      : undefined,
  });
}

/** Converts an InlineGeometry proto into a local shape. */
export function inlineGeometryProtoToLocalShape(
  proto: InlineGeometryProto,
): LocalShape {
  const data = proto.exactGeometry?.data;
  let exactGeometry: TransformedPrimitiveShapeSet | TriangleMesh | PointCloud;
  switch (data?.case) {
    case 'primitiveSet':
      exactGeometry = {
        primitives: data.value.primitives.map(
          transformedPrimitiveShapeFromProto,
        ),
      };
      break;
    case 'triangleMesh':
      exactGeometry = {
        vertices: new Float64Array(data.value.vertices),
        faces: new Int32Array(data.value.faces),
      };
      break;
    case 'pointCloud':
      exactGeometry = {
        points: new Float64Array(data.value.points),
        normals: new Float64Array(data.value.normals),
      };
      break;
    default:
      throw new Error(`InlineGeometry proto can't be converted to LocalShape!`);
  }

  // Both a renderable and a generated renderable are rendered the same way.
  const renderableProto = proto.renderableType.value;
  return {
    exactGeometry,
    renderable: renderableProto
      ? new GltfRenderable(toArrayBufferBytes(renderableProto.glbBytes))
      : undefined,
  };
}

/** Copies `bytes` if they aren't backed by a plain ArrayBuffer. */
function toArrayBufferBytes(bytes: Uint8Array): Uint8Array<ArrayBuffer> {
  return bytes.buffer instanceof ArrayBuffer
    ? (bytes as Uint8Array<ArrayBuffer>)
    : new Uint8Array(bytes);
}

/** Converts a geometry shape into a Geometry proto. */
export function geometryShapeToGeometryProto(
  shape: GeometryShape,
): GeometryProto {
  if ('exactGeometry' in shape) {
    return create(GeometrySchema, {
      data: {
        case: 'inlineGeometryData',
        value: localShapeToInlineGeometryProto(shape),
      },
    });
  }
  if ('exactGeometryRef' in shape) {
    return create(GeometrySchema, {
      data: {case: 'geoRef', value: remoteShapeToGeometryRefsProto(shape)},
    });
  }
  throw new Error(`GeometryShape can't be converted to a Geometry proto!`);
}

/** Converts a Geometry proto into a geometry shape. */
export function geometryProtoToGeometryShape(
  proto: GeometryProto,
): GeometryShape {
  switch (proto.data.case) {
    case 'inlineGeometryData':
      return inlineGeometryProtoToLocalShape(proto.data.value);
    case 'geoRef':
      return geometryRefsProtoToRemoteShape(proto.data.value);
    default:
      throw new Error(`Geometry proto can't be converted to GeometryShape!`);
  }
}

function materialPropertiesFromProto(
  proto: MaterialPropertiesProto,
): MaterialProperties {
  const color = new three.Color(
    proto.baseColor?.red ?? 0,
    proto.baseColor?.green ?? 0,
    proto.baseColor?.blue ?? 0,
  );
  return {
    color: `#${color.getHexString()}`,
    metalness: proto.metalness ?? 0,
    roughness: proto.roughness ?? 0,
    transmission: proto.transmission ?? 0,
  };
}

function materialPropertiesToProto(
  materialProperties: MaterialProperties,
): MaterialPropertiesProto {
  const color = new three.Color(materialProperties.color);
  return create(MaterialPropertiesSchema, {
    baseColor: {red: color.r, green: color.g, blue: color.b},
    metalness: materialProperties.metalness,
    roughness: materialProperties.roughness,
    transmission: materialProperties.transmission,
  });
}

/** Converts a TransformedGeometry proto into a native geometry. */
export function transformedGeometryProtoToGeometry(
  proto: TransformedGeometryProto,
): Geometry {
  const geometryProto = proto.geometry;
  if (!geometryProto) {
    throw new Error(
      `TransformedGeometry proto can't be converted to Geometry, geometry is not set.`,
    );
  }
  const shape = geometryProtoToGeometryShape(geometryProto);
  const materialOverrides = geometryProto.materialOverrides
    ? materialPropertiesFromProto(geometryProto.materialOverrides)
    : undefined;
  const refTShape = proto.refTShape
    ? geometricTransformProtoToMatrix(proto.refTShape)
    : new three.Matrix4();
  return new Geometry(refTShape, shape, materialOverrides);
}

/** Converts a native geometry into a TransformedGeometry proto. */
export function geometryToTransformedGeometryProto(
  geometry: Geometry,
): TransformedGeometryProto {
  const geometryProto = geometryShapeToGeometryProto(geometry.getShape());
  const materialOverrides = geometry.getMaterialOverrides();
  if (materialOverrides) {
    geometryProto.materialOverrides =
      materialPropertiesToProto(materialOverrides);
  }
  return create(TransformedGeometrySchema, {
    geometry: geometryProto,
    refTShape: matrixToGeometricTransformProto(geometry.getRefTShape()),
  });
}

function geometryTypeToProto(type: GeometryType): GeometryTypeProto {
  switch (type) {
    case GeometryType.Visual:
      return GeometryTypeProto.VISUAL;
    case GeometryType.Collision:
      return GeometryTypeProto.COLLISION;
  }
}

/**
 * Converts a native geometry update of the entity named `entityName` into a
 * GeometryUpdate proto.
 */
export function geometryUpdateToProto(
  geometryUpdate: GeometryUpdate,
  entityName: string,
): GeometryUpdateProto {
  return create(GeometryUpdateSchema, {
    entityName,
    geometriesToRemove: geometryUpdate.geometriesToRemove.map((r) => ({
      type:
        r.type !== undefined
          ? geometryTypeToProto(r.type)
          : GeometryTypeProto.UNSPECIFIED,
      geometryName: r.geometryName ?? '',
    })),
    geometriesToSet: geometryUpdate.geometriesToSet.map((s) => ({
      type: geometryTypeToProto(s.type),
      geometryName: s.geometryName,
      geometry:
        s.geometry !== undefined
          ? geometryToTransformedGeometryProto(s.geometry)
          : undefined,
    })),
  });
}
