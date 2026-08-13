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
