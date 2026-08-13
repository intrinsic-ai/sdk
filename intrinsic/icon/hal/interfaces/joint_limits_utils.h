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

#ifndef INTRINSIC_ICON_HAL_INTERFACES_JOINT_LIMITS_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_JOINT_LIMITS_UTILS_H_

#include <cstdint>

#include "absl/status/status.h"
#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/hardware_interface_handle.h"
#include "intrinsic/icon/hal/interfaces/joint_limits.fbs.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/kinematics/types/joint_limits.h"
#include "intrinsic/kinematics/types/joint_limits.pb.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildJointLimits(uint32_t num_dof);

}  // namespace intrinsic_fbs

namespace intrinsic::icon {

// Parses a JointLimits protobuf into a JointLimits hardware interface handle.
// Returns kInvalidArgument if the number of joints in the non-empty protobuf
// fields are not equal to the ones in the flatbuffer handle. Expects all
// non-empty JointLimits proto fields to have the same size and each flatbuffer
// field to match that size.
absl::Status ParseProtoJointLimits(
    const intrinsic_proto::JointLimits& pb_limits,
    intrinsic::icon::MutableHardwareInterfaceHandle<intrinsic_fbs::JointLimits>&
        fb_limits);

// Copies a JointLimits struct to a JointLimits flatbuffer. Fails if the number
// of joints in the struct does not match the size of the flatbuffer.
RealtimeStatus CopyTo(const JointLimits& limits,
                      intrinsic_fbs::JointLimits& fb_limits);

// Copies a JointLimits flatbuffer to a JointLimits struct.
RealtimeStatus CopyTo(const intrinsic_fbs::JointLimits& fb_limits,
                      JointLimits& limits);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_INTERFACES_JOINT_LIMITS_UTILS_H_
