// Copyright 2023 Intrinsic Innovation LLC

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
