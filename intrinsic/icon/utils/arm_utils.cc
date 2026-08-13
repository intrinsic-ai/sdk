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

#include "intrinsic/icon/utils/arm_utils.h"

#include <cstddef>
#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/part_status.pb.h"

namespace intrinsic::icon {

eigenmath::VectorXd GetCommandedOrSensedJointPosition(
    const intrinsic_proto::icon::PartStatus& part_status) {
  return GetPositionCommandedLastCycle(part_status)
      .value_or(GetSensedJointPosition(part_status));
}

eigenmath::VectorXd GetSensedJointPosition(
    const intrinsic_proto::icon::PartStatus& part_status) {
  eigenmath::VectorXd result =
      eigenmath::VectorXd::Zero(part_status.joint_states().size());
  for (size_t i = 0; i < result.size(); ++i) {
    result(i) = part_status.joint_states(i).position_sensed();
  }
  return result;
}

std::optional<eigenmath::VectorXd> GetPositionCommandedLastCycle(
    const intrinsic_proto::icon::PartStatus& part_status) {
  eigenmath::VectorXd result =
      eigenmath::VectorXd::Zero(part_status.joint_states().size());
  if (result.size() == 0) {
    return std::nullopt;
  }
  for (size_t i = 0; i < result.size(); ++i) {
    if (!part_status.joint_states(i).has_position_commanded_last_cycle()) {
      return std::nullopt;
    }
    result(i) = part_status.joint_states(i).position_commanded_last_cycle();
  }
  return result;
}

eigenmath::VectorXd GetSensedJointVelocity(
    const intrinsic_proto::icon::PartStatus& part_status) {
  eigenmath::VectorXd result =
      eigenmath::VectorXd::Zero(part_status.joint_states().size());
  for (size_t i = 0; i < result.size(); ++i) {
    result(i) = part_status.joint_states(i).velocity_sensed();
  }
  return result;
}

}  // namespace intrinsic::icon
