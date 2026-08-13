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

#ifndef INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_
#define INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_

#include <cstddef>
#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::icon {

// Represents a set of command parameters for joint acceleration control.
class JointAccelerationCommand {
 public:
  // Default constructor to make this play nice with containers and StatusOr.
  JointAccelerationCommand() : acceleration_(0) {}

  explicit JointAccelerationCommand(const eigenmath::VectorNd& acceleration)
      : acceleration_(acceleration) {}

  // Builds a JointAccelerationCommand object.
  //
  // Returns InvalidArgument if any of the two vectors' sizes don't match up.
  static RealtimeStatusOr<JointAccelerationCommand> Create(
      const eigenmath::VectorNd& acceleration,
      const std::optional<eigenmath::VectorNd>& torque = std::nullopt);

  const eigenmath::VectorNd& acceleration() const;

  const std::optional<eigenmath::VectorNd>& torque() const;

  size_t Size() const;

 private:
  JointAccelerationCommand(const eigenmath::VectorNd& acceleration,
                           const std::optional<eigenmath::VectorNd>& torque);

  eigenmath::VectorNd acceleration_;
  std::optional<eigenmath::VectorNd> torque_ = std::nullopt;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_JOINT_ACCELERATION_COMMAND_H_
