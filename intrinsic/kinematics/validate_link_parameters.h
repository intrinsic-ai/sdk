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

#ifndef INTRINSIC_KINEMATICS_VALIDATE_LINK_PARAMETERS_H_
#define INTRINSIC_KINEMATICS_VALIDATE_LINK_PARAMETERS_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::kinematics {

// Validates that the link mass is positive.
icon::RealtimeStatus ValidateMass(double mass_kg);
// Validates that the link inertia expressed at the center of gravity is
// positive definite (symmetric and with positive eigenvalues) and that its
// eigenvalues fulfill the triangle inequalities.
icon::RealtimeStatus ValidateInertia(const eigenmath::Matrix3d& inertia);
}  // namespace intrinsic::kinematics

#endif  // INTRINSIC_KINEMATICS_VALIDATE_LINK_PARAMETERS_H_
