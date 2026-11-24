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
 * @file Conversions between native math poses and protobuf-es Pose protos.
 */
import {create} from '@bufbuild/protobuf';
import {Pose, type ReadonlyPose} from '@intrinsic-ai/core/math';
import {
  type Pose as PoseProto,
  PoseSchema,
} from '@intrinsic-ai/types/math/proto/pose_pb';
import * as three from 'three';

/**
 * Converts a native pose into a protobuf-es Pose proto. The rotation is
 * normalized.
 */
export function poseToProto(pose: ReadonlyPose): PoseProto {
  const position = pose.getTranslation();
  const orientation = pose.getRotation().normalize();
  return create(PoseSchema, {
    position: {x: position.x, y: position.y, z: position.z},
    orientation: {
      x: orientation.x,
      y: orientation.y,
      z: orientation.z,
      w: orientation.w,
    },
  });
}

/**
 * Converts a transform matrix into a protobuf-es Pose proto. Any scale in the
 * matrix is discarded and the rotation is normalized.
 */
export function matrixToPoseProto(mat: three.Matrix4): PoseProto {
  const position = new three.Vector3();
  const rotation = new three.Quaternion();
  mat.decompose(position, rotation, new three.Vector3());
  return poseToProto(new Pose(position, rotation));
}

/**
 * Converts a protobuf-es Pose proto into a transform matrix. A missing position
 * is the origin and a missing orientation is the identity rotation.
 */
export function poseProtoToMatrix(proto: PoseProto | undefined): three.Matrix4 {
  const position = proto?.position;
  const orientation = proto?.orientation;
  return new three.Matrix4().compose(
    new three.Vector3(position?.x ?? 0, position?.y ?? 0, position?.z ?? 0),
    orientation
      ? new three.Quaternion(
          orientation.x,
          orientation.y,
          orientation.z,
          orientation.w,
        )
      : new three.Quaternion(),
    new three.Vector3(1, 1, 1),
  );
}
