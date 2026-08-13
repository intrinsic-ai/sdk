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

#ifndef INTRINSIC_ICON_CONTROL_PARTS_REALTIME_ROBOT_PAYLOAD_H_
#define INTRINSIC_ICON_CONTROL_PARTS_REALTIME_ROBOT_PAYLOAD_H_

#include <cstddef>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/fixed_string.h"
#include "intrinsic/icon/utils/realtime_status_or.h"
#include "intrinsic/kinematics/types/to_fixed_string.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/proto/robot_payload.pb.h"
#include "intrinsic/world/robot_payload/robot_payload_base.h"

namespace intrinsic::icon {

// Number of characters needed to print a RobotPayload. See ToFixedString.
constexpr size_t kRobotPayloadStrSize =
    kPose3dStrSize + 7 * kDoubleStrSize + 66;

// Payload of a robot. Provides a real-time safe factory method.
class RealtimeRobotPayload : public RobotPayloadBase {
 public:
  // Creates a Payload from the given parameters. Fails if parameters are
  // invalid. The inertia matrix is expressed about the payload's center of
  // mass.
  // The mass and the inertia can be zero.
  static RealtimeStatusOr<RealtimeRobotPayload> Create(
      double mass_kg, const Pose3d& tip_t_cog,
      const eigenmath::Matrix3d& inertia);
};

FixedString<kRobotPayloadStrSize> ToFixedString(
    const RobotPayloadBase& payload);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_PARTS_REALTIME_ROBOT_PAYLOAD_H_
