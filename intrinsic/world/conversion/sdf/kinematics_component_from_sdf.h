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

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_KINEMATICS_COMPONENT_FROM_SDF_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_KINEMATICS_COMPONENT_FROM_SDF_H_

#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/world/component/kinematics_component.h"
#include "sdf/Joint.hh"
#include "sdf/Model.hh"

namespace intrinsic {
namespace sdf {
// Creates a kinematics component without specifying the joint's inbound pose.
// Information on `sdf_joint`'s parent <model> element is required for reasoning
// about that information.
absl::StatusOr<std::unique_ptr<KinematicsComponent>>
KinematicsComponentFromSdfJoint(const ::sdf::Joint& sdf_joint);
}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_CONVERSION_SDF_KINEMATICS_COMPONENT_FROM_SDF_H_
