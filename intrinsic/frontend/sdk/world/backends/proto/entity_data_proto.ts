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
 * @file Conversions from protobuf-es world entity component protos to native
 * entity data.
 */
import type {GeometryComponent} from '@intrinsic-ai/types/world/proto/geometry_component_pb';
import {
  type KinematicsComponent,
  type KinematicsComponent_Limits,
  KinematicsComponent_MotionType,
} from '@intrinsic-ai/types/world/proto/kinematics_component_pb';
import type {
  SensorComponent,
  SensorComponent_Camera,
  SensorComponent_DepthCamera,
} from '@intrinsic-ai/types/world/proto/sensor_component_pb';
import * as three from 'three';

import {
  type Clip,
  type EntityData,
  type JointLimits as EntityJointLimits,
  type GeometryData,
  type JointEntityData,
  type LinkEntityData,
  MotionType,
  type Sensor,
  type SensorEntityData,
} from '../../model/entity';
import {transformedGeometryProtoToGeometry} from './geometry_proto';
import {poseProtoToMatrix} from './pose_proto';

function motionTypeFromProto(
  motionType: KinematicsComponent_MotionType,
  entityName: string,
): MotionType {
  switch (motionType) {
    case KinematicsComponent_MotionType.UNDEFINED:
      return MotionType.UNDEFINED;
    case KinematicsComponent_MotionType.REVOLUTE:
      return MotionType.REVOLUTE;
    case KinematicsComponent_MotionType.PRISMATIC:
      return MotionType.PRISMATIC;
    case KinematicsComponent_MotionType.FIXED:
      return MotionType.FIXED;
    default:
      console.warn(
        `Internal error parsing motion type for entity '${entityName}'.` +
          `Defaulting to REVOLUTE`,
      );
      return MotionType.REVOLUTE;
  }
}

function limitsFromProto(
  proto: KinematicsComponent_Limits | undefined,
): EntityJointLimits {
  const fixedLimits =
    proto?.rawValueLimitsType.case === 'fixedLimits'
      ? proto.rawValueLimitsType.value
      : undefined;
  return {
    velocity: proto?.velocity ?? 0,
    acceleration: proto?.acceleration ?? 0,
    jerk: proto?.jerk ?? 0,
    effort: proto?.effort ?? 0,
    fixed: fixedLimits
      ? {upper: fixedLimits.upper, lower: fixedLimits.lower}
      : undefined,
  };
}

/**
 * Converts a KinematicsComponent proto into joint entity data, on top of the
 * entity's `baseData`.
 */
export function kinematicsComponentToEntityData(
  baseData: EntityData,
  kinematics: KinematicsComponent,
): JointEntityData {
  const axis = kinematics.axis
    ? new three.Vector3(kinematics.axis.x, kinematics.axis.y, kinematics.axis.z)
    : new three.Vector3(0, 0, 1);
  return {
    ...baseData,
    motionType: motionTypeFromProto(kinematics.motionType, baseData.name),
    axis,
    rawValue: kinematics.rawValue ?? 0,
    systemLimits: limitsFromProto(kinematics.systemLimits),
    applicationLimits: limitsFromProto(kinematics.applicationLimits),
    parentTInboard: poseProtoToMatrix(kinematics.parentTInboard),
    outboardTChild: poseProtoToMatrix(kinematics.outboardTChild),
  };
}

/**
 * Converts a GeometryComponent proto into link entity data, on top of the
 * entity's `baseData`.
 */
export function geometryComponentToEntityData(
  baseData: EntityData,
  geometryComponent: GeometryComponent,
): LinkEntityData {
  const geometries = new Map<string, Map<string, GeometryData>>();
  for (const [setName, geometrySet] of Object.entries(
    geometryComponent.namedGeometries,
  )) {
    const geometryDataMap = new Map<string, GeometryData>();
    for (const [name, geometryProto] of Object.entries(
      geometrySet.namedGeometries,
    )) {
      const geometry = transformedGeometryProtoToGeometry(geometryProto);
      geometryDataMap.set(name, {
        refTShape: geometry.getRefTShape(),
        shape: geometry.getShape(),
        materialOverrides: geometry.getMaterialOverrides(),
      });
    }
    geometries.set(setName, geometryDataMap);
  }
  return {...baseData, geometries};
}

function cameraToClip(
  camera: SensorComponent_Camera | SensorComponent_DepthCamera,
): Clip | undefined {
  const clip = camera.properties?.clip;
  if (clip && clip.near > 0 && clip.far > clip.near) {
    return {near: clip.near, far: clip.far};
  }
  return undefined;
}

/**
 * Converts a SensorComponent proto into sensor entity data, on top of the
 * entity's `baseData`.
 */
export function sensorComponentToEntityData(
  baseData: EntityData,
  sensorComponent: SensorComponent,
): SensorEntityData {
  let sensor: Sensor | undefined;
  switch (sensorComponent.typeOneof.case) {
    case 'camera':
      sensor = {
        type: 'camera',
        clip: cameraToClip(sensorComponent.typeOneof.value),
      };
      break;
    case 'depthCamera':
      sensor = {
        type: 'depth',
        clip: cameraToClip(sensorComponent.typeOneof.value),
      };
      break;
    default:
      sensor = undefined;
  }
  return {...baseData, sensor};
}
