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

#ifndef INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_H_
#define INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_H_

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/proto/robot_payload.pb.h"
#include "intrinsic/world/robot_payload/robot_payload_base.h"

namespace intrinsic {

// Dynamic payload of a robot.
class RobotPayload : public RobotPayloadBase {
 public:
  // Creates a Payload from the given parameters. Fails if parameters are
  // invalid.
  static absl::StatusOr<RobotPayload> Create(
      double mass_kg, const Pose3d& tip_t_cog,
      const eigenmath::Matrix3d& inertia);

  // Sets the mass of the robot payload. Unit is kg.
  absl::Status SetMass(double mass_kg);

  // Set the center of gravity.
  absl::Status SetTipTCog(const Pose3d& tip_t_cog);

  // Sets the inertia matrix expressed about the payloads center of mass. Unit
  // is kg*m^2.
  absl::Status SetInertia(const eigenmath::Matrix3d& inertia);
};

// Converts a payload proto to a Payload. Fails for invalid parameters. Sets
// empty mass to 0 kg, empty tip_t_cog to the identity transform and an empty
// inertia to a 3x3 zero matrix.
absl::StatusOr<RobotPayload> FromProto(
    const intrinsic_proto::world::RobotPayload& proto);

// Converts a Payload to a payload proto.
intrinsic_proto::world::RobotPayload ToProto(const RobotPayload& payload);

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_H_
