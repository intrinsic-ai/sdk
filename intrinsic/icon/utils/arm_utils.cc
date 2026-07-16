// Copyright 2023 Intrinsic Innovation LLC

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
