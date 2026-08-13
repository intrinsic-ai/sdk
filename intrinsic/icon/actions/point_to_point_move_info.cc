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

#include "intrinsic/icon/actions/point_to_point_move_info.h"

#include "absl/types/span.h"
#include "intrinsic/icon/proto/joint_space.pb.h"
#include "intrinsic/kinematics/types/joint_limits.h"

namespace intrinsic {
namespace icon {

PointToPointMoveInfo::FixedParams CreatePointToPointMoveFixedParams(
    absl::Span<const double> goal_position,
    absl::Span<const double> goal_velocity) {
  PointToPointMoveInfo::FixedParams fixed_params;
  *fixed_params.mutable_goal_position()->mutable_joints() = {
      goal_position.begin(), goal_position.end()};
  *fixed_params.mutable_goal_velocity()->mutable_joints() = {
      goal_velocity.begin(), goal_velocity.end()};
  return fixed_params;
}

PointToPointMoveInfo::FixedParams CreatePointToPointMoveFixedParams(
    absl::Span<const double> goal_position,
    absl::Span<const double> goal_velocity, const JointLimits& joint_limits) {
  PointToPointMoveInfo::FixedParams fixed_params;
  *fixed_params.mutable_goal_position()->mutable_joints() = {
      goal_position.begin(), goal_position.end()};
  *fixed_params.mutable_goal_velocity()->mutable_joints() = {
      goal_velocity.begin(), goal_velocity.end()};
  *fixed_params.mutable_joint_limits() = ToProto(joint_limits);
  return fixed_params;
}

}  // namespace icon
}  // namespace intrinsic
