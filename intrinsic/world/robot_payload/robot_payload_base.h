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

#ifndef INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_BASE_H_
#define INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_BASE_H_

#include <ostream>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/proto/robot_payload.pb.h"

namespace intrinsic {

// Base class of the payload of a robot. It is read-only and real-time
// safe.
class RobotPayloadBase {
 public:
  RobotPayloadBase();

  // Mass of the robot payload. Unit is kg.
  double mass() const { return mass_kg_; }

  // Center of gravity of the robot payload relative to the robot flange/tip
  // frame.
  const Pose3d& tip_t_cog() const { return tip_t_cog_; }

  // 3x3 symmetric inertia matrix of the robot payload expressed about the
  // payloads center of mass. Unit is kg*m^2.
  const eigenmath::Matrix3d& inertia() const { return inertia_in_cog_; }

  // Uses IsApprox to compare the payloads with `kStdError` as
  // precision.
  bool operator==(const RobotPayloadBase& other) const;

  // Returns true if the two payloads are almost equal given the precision.
  bool IsApprox(const intrinsic::RobotPayloadBase& other,
                double precision) const;

 protected:
  // This constructor is protected to force the use of the factory method in the
  // derived class. The factory method ensures that the payload is valid.
  RobotPayloadBase(double mass, const Pose3d& tip_t_cog,
                   const eigenmath::Matrix3d& inertia);

  double mass_kg_;
  Pose3d tip_t_cog_;
  eigenmath::Matrix3d inertia_in_cog_;
};

std::ostream& operator<<(std::ostream& os, const RobotPayloadBase& payload);

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_ROBOT_PAYLOAD_ROBOT_PAYLOAD_BASE_H_
