// Copyright 2023 Intrinsic Innovation LLC

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
