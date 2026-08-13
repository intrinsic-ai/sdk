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

#ifndef INTRINSIC_KINEMATICS_TYPES_DYNAMIC_LIMITS_CHECK_MODE_H_
#define INTRINSIC_KINEMATICS_TYPES_DYNAMIC_LIMITS_CHECK_MODE_H_

#include <cstdint>

#include "absl/status/statusor.h"
#include "intrinsic/kinematics/types/dynamic_limits_check_mode.pb.h"

namespace intrinsic {

// Enum with available limits types to be checked for second order variables
// such as joint accelerations or joint torques. The default behavior is to
// check joint accelerations `kCheckJointAcceleration`, but can be
// skipped with `kCheckNone`.
enum class DynamicLimitsCheckMode : int8_t {
  kCheckJointAcceleration,  // Checks joint acceleration limits.
  kCheckNone,  // No checks for second order variables such as torques or
               // joint accelerations limits.
};

intrinsic_proto::DynamicLimitsCheckMode ToProto(
    const DynamicLimitsCheckMode& dynamic_limits_check_mode);

// Maps a proto DynamicLimitsCheckMode to its equivalent enum class. As the
// proto has an additional UNSPECIFIED field compared to the c++ enum class,
// this one gets mapped to the default value of `kCheckJointAcceleration`.
// Returns an error if the input type is not known.
absl::StatusOr<DynamicLimitsCheckMode> FromProto(
    const intrinsic_proto::DynamicLimitsCheckMode&
        dynamic_limits_check_mode_proto);

}  // namespace intrinsic

#endif  // INTRINSIC_KINEMATICS_TYPES_DYNAMIC_LIMITS_CHECK_MODE_H_
