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

#ifndef INTRINSIC_ICON_HAL_INTERFACES_JOINT_COMMAND_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_JOINT_COMMAND_UTILS_H_

#include <cstdint>

#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/interfaces/joint_command.fbs.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildJointPositionCommand(uint32_t num_dof);

flatbuffers::DetachedBuffer BuildJointVelocityCommand(uint32_t num_dof);

flatbuffers::DetachedBuffer BuildJointTorqueCommand(uint32_t num_dof);

flatbuffers::DetachedBuffer BuildJointAccelerationAndTorqueCommand(
    uint32_t num_dof);

flatbuffers::DetachedBuffer BuildHandGuidingCommand();

intrinsic::icon::RealtimeStatus CopyTo(const JointPositionCommand& src,
                                       JointPositionCommand& dest);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_JOINT_COMMAND_UTILS_H_
