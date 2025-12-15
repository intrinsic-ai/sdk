// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/control/joint_acceleration_command.h"

#include <cstddef>
#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::icon {

RealtimeStatusOr<JointAccelerationCommand> JointAccelerationCommand::Create(
    const eigenmath::VectorNd& acceleration,
    const std::optional<eigenmath::VectorNd>& torque) {
  if (torque.has_value() && torque->size() != acceleration.size()) {
    return InvalidArgumentError(RealtimeStatus::StrCat(
        "Torque has ", torque->size(),
        " values, but acceleration setpoints have ", acceleration.size()));
  }

  return JointAccelerationCommand(acceleration, torque);
}

const eigenmath::VectorNd& JointAccelerationCommand::acceleration() const {
  return acceleration_;
}

const std::optional<eigenmath::VectorNd>& JointAccelerationCommand::torque()
    const {
  return torque_;
}

size_t JointAccelerationCommand::Size() const { return acceleration_.size(); }

JointAccelerationCommand::JointAccelerationCommand(
    const eigenmath::VectorNd& acceleration,
    const std::optional<eigenmath::VectorNd>& torque)
    : acceleration_(acceleration), torque_(torque) {}

}  // namespace intrinsic::icon
