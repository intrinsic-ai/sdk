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
 * @file Conversions between native joint limits and protobuf-es JointLimits
 * protos.
 */
import {create} from '@bufbuild/protobuf';
import {
  type JointLimits as JointLimitsProto,
  type JointLimitsUpdate as JointLimitsUpdateProto,
  JointLimitsUpdateSchema,
  type RepeatedDouble,
  RepeatedDoubleSchema,
} from '@intrinsic-ai/types/kinematics/types/joint_limits_pb';

import {
  JointLimits,
  JointLimitsUpdate,
  type PerJointLimits,
  type PerJointLimitsUpdate,
} from '../../model/joint_limits';

/**
 * Converts a native joint limits update into a protobuf-es JointLimitsUpdate
 * proto. Each field must be set either for all joints or for none of them.
 */
export function jointLimitsUpdateToProto(
  update: JointLimitsUpdate,
): JointLimitsUpdateProto {
  return create(JointLimitsUpdateSchema, {
    minPosition: collectJointLimitField(update, 'minPosition'),
    maxPosition: collectJointLimitField(update, 'maxPosition'),
    maxVelocity: collectJointLimitField(update, 'maxVelocity'),
    maxAcceleration: collectJointLimitField(update, 'maxAcceleration'),
    maxJerk: collectJointLimitField(update, 'maxJerk'),
    maxEffort: collectJointLimitField(update, 'maxEffort'),
  });
}

/**
 * Converts a protobuf-es JointLimits proto into native joint limits. All set
 * fields must have the same number of values.
 */
export function jointLimitsFromProto(proto: JointLimitsProto): JointLimits {
  const minPosition = proto.minPosition?.values;
  const maxPosition = proto.maxPosition?.values;
  const maxVelocity = proto.maxVelocity?.values;
  const maxAcceleration = proto.maxAcceleration?.values;
  const maxJerk = proto.maxJerk?.values;
  const maxEffort = proto.maxEffort?.values;

  const allValues = [
    minPosition,
    maxPosition,
    maxVelocity,
    maxAcceleration,
    maxJerk,
    maxEffort,
  ];
  const numJoints = Math.max(0, ...allValues.map((v) => v?.length ?? 0));
  if (allValues.some((v) => v !== undefined && v.length !== numJoints)) {
    throw new Error('All set joint limit fields must have the same length');
  }
  const joints: PerJointLimits[] = [];
  for (let i = 0; i < numJoints; i++) {
    joints.push({
      minPosition: minPosition?.[i],
      maxPosition: maxPosition?.[i],
      maxVelocity: maxVelocity?.[i],
      maxAcceleration: maxAcceleration?.[i],
      maxJerk: maxJerk?.[i],
      maxEffort: maxEffort?.[i],
    });
  }
  return new JointLimits(joints);
}

function collectJointLimitField(
  update: JointLimitsUpdate,
  field: keyof PerJointLimitsUpdate,
): RepeatedDouble | undefined {
  const values: number[] = [];
  for (const joint of update.joints) {
    const value = joint[field];
    if (value !== undefined) {
      values.push(value);
    }
  }
  if (values.length === 0) {
    return undefined;
  }
  if (values.length !== update.joints.length) {
    throw new Error(
      `Joint limits update must set '${field}' for all joints or none ` +
        `(set for ${values.length} of ${update.joints.length})`,
    );
  }
  return create(RepeatedDoubleSchema, {values});
}
