// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/conversion/sdf/kinematics_component_from_sdf.h"

#include <cmath>
#include <limits>
#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "gz/math/Vector3.hh"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/scene/constants.h"
#include "intrinsic/scene/sdf/custom_tags.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/simulation/gazebo/type_conversion.h"
#include "intrinsic/util/status/ret_check.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/component/kinematics_component.h"
#include "intrinsic/world/hashing/hashing.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"
#include "sdf/Joint.hh"
#include "sdf/JointAxis.hh"

namespace intrinsic {
namespace sdf {

namespace {

constexpr const char* kLimitElement = "limit";

absl::StatusOr<intrinsic_proto::world::KinematicsComponent::MotionType>
IntrMotionTypeFromSdfJointType(::sdf::JointType sdf_joint_type) {
  using ::intrinsic_proto::world::KinematicsComponent;
  static const auto* kMotionTypeFromJointType =
      new WorldHashMap<::sdf::JointType, KinematicsComponent::MotionType>(
          {{::sdf::JointType::REVOLUTE,
            KinematicsComponent::MOTION_TYPE_REVOLUTE},
           {::sdf::JointType::PRISMATIC,
            KinematicsComponent::MOTION_TYPE_PRISMATIC},
           {::sdf::JointType::FIXED, KinematicsComponent::MOTION_TYPE_FIXED},
           // Ball joints are handled as fixed joints for now.
           {::sdf::JointType::BALL, KinematicsComponent::MOTION_TYPE_FIXED}});
  if (!kMotionTypeFromJointType->contains(sdf_joint_type)) {
    return absl::UnimplementedError(absl::Substitute(
        "Joint type does not have Intrinsic MotionType equivalent."));
  }
  return kMotionTypeFromJointType->at(sdf_joint_type);
}

}  // namespace

absl::StatusOr<std::unique_ptr<KinematicsComponent>>
KinematicsComponentFromSdfJoint(const ::sdf::Joint& sdf_joint) {
  std::unique_ptr<KinematicsComponent> kinematics_component =
      KinematicsComponent::Create();

  INTR_ASSIGN_OR_RETURN(
      auto motion_type, IntrMotionTypeFromSdfJointType(sdf_joint.Type()),
      _ << "Joint type " << sdf_joint.Element()->Get<std::string>("type")
        << " is not supported.");
  kinematics_component->SetMotionType(motion_type);
  // If motion type is fixed, we don't need the rest of the joint.
  if (motion_type ==
      ::intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_FIXED) {
    return kinematics_component;
  }

  // Get the JointAxis, or error if not found.
  ::sdf::JointAxis sdf_joint_axis;
  if (auto axis_ptr = sdf_joint.Axis(); axis_ptr) {
    sdf_joint_axis = *axis_ptr;
  } else {
    return absl::InvalidArgumentError(
        absl::Substitute("Could not find an axis for joint '$0'. Please "
                         "provide the <axis> element in the SDF file.",
                         sdf_joint.Name()));
  }

  // Note - the ResolveXyz method gives the axis value with respect to the
  // joint's frame (aka joint's inboard frame). This is the same convention we
  // use in our KinematicsComponent proto (the joint axis is with respect to the
  // joint frame). In some legacy SDF files the <axis> value is specified with
  // respect to the joint's parent space, but the SDF library handles that
  // conversion automatically here. See:
  // https://sdformat.org/tutorials/specification/spec_model_kinematics/#jointaxis
  ::gz::math::Vector3d xyz_in_joint_frame;
  if (auto errors = sdf_joint_axis.ResolveXyz(xyz_in_joint_frame);
      !errors.empty()) {
    return absl::InvalidArgumentError(
        absl::Substitute("Joint axis xyz of joint '$0' failed to resolve to "
                         "the joint frame. Errors are:\n$1",
                         sdf_joint.Name(),
                         absl::StrJoin(errors, "\n", absl::StreamFormatter())));
  }
  const eigenmath::Vector3d intr_xyz = GzToIntrinsic(xyz_in_joint_frame);
  kinematics_component->SetAxis(intr_xyz);

  // The underlying SDF parser interprets missing "lower" and "upper" elements
  // as -1e16 and 1e16, respectively. We will instead convert these values to
  // infinity.
  constexpr double kInfinityRoundingThreshold = 9.9e15;

  double lower_position = sdf_joint_axis.Lower();
  if (lower_position < -kInfinityRoundingThreshold) {
    lower_position = -std::numeric_limits<double>::infinity();
  }
  double upper_position = sdf_joint_axis.Upper();
  if (upper_position > kInfinityRoundingThreshold) {
    upper_position = std::numeric_limits<double>::infinity();
  }
  INTR_RETURN_IF_ERROR(kinematics_component->SetApplicationRawValueFixedLimits(
      lower_position, upper_position, /*enforce_limits=*/false));
  INTR_RETURN_IF_ERROR(kinematics_component->SetSystemRawValueFixedLimits(
      lower_position, upper_position, /*enforce_limits=*/true));

  // Process optional <velocity>. Default is -1, meaning unenforced/infinity.
  const double velocity = sdf_joint_axis.MaxVelocity();
  if (std::isfinite(velocity) && velocity >= 0) {
    // For an unlimited kinematics_component just created, we don't need to
    // enforce limit on the first set of a system limit.
    INTR_RETURN_IF_ERROR(kinematics_component->SetSystemVelocityLimit(
        velocity, /*enforce_limits=*/false));
    INTR_RET_CHECK_OK(kinematics_component->SetApplicationVelocityLimit(
        velocity * scene_object::kMaxApplicationLimitsMultiplier,
        /*enforce_limits=*/true))
        << "Setting application velocity limit to a fraction "
        << scene_object::kMaxApplicationLimitsMultiplier
        << " of system limit should never fail.";
  }

  // Process optional <effort>. Default is -1, meaning unenforced/infinity.
  const double effort = sdf_joint_axis.Effort();
  if (std::isfinite(effort) && effort >= 0) {
    INTR_RETURN_IF_ERROR(kinematics_component->SetApplicationEffortLimit(
        effort, /*enforce_limits=*/false));
    INTR_RETURN_IF_ERROR(kinematics_component->SetSystemEffortLimit(
        effort, /*enforce_limits=*/true));
  }

  // Use raw element to process Intrinsic custom elements.
  auto sdf_joint_axis_element = sdf_joint_axis.Element();
  if (!sdf_joint_axis_element->HasElement(kLimitElement)) {
    return absl::InvalidArgumentError(
        absl::Substitute("SDF joint <axis> requires a <limit> element. No "
                         "<limit> element found in sdf <axis> element:\n $0",
                         sdf_joint_axis_element->ToString("")));
  }

  // Process optional <intrinsic:acceleration>.
  auto limit_element = sdf_joint_axis_element->GetElement(kLimitElement);
  if (limit_element->HasElement(std::string(kAccelerationCustomElement))) {
    INTR_ASSIGN_OR_RETURN(
        auto acceleration,
        ParseChildAs<double>(limit_element,
                             std::string(kAccelerationCustomElement)));
    // For an unlimited kinematics_component just created, we don't need to
    // enforce limit on the first set of a system limit.
    INTR_RETURN_IF_ERROR(kinematics_component->SetSystemAccelerationLimit(
        acceleration, /*enforce_limits=*/false));
    INTR_RET_CHECK_OK(kinematics_component->SetApplicationAccelerationLimit(
        acceleration * scene_object::kMaxApplicationLimitsMultiplier,
        /*enforce_limits=*/true))
        << "Setting application acceleration limit to a fraction "
        << scene_object::kMaxApplicationLimitsMultiplier
        << " of system limit should never fail.";
  }

  // Process optional <intrinsic:jerk>.
  if (limit_element->HasElement(std::string(kJerkCustomElement))) {
    INTR_ASSIGN_OR_RETURN(
        double jerk,
        ParseChildAs<double>(limit_element, std::string(kJerkCustomElement)));
    // For an unlimited kinematics_component just created, we don't need to
    // enforce limit on the first set of a system limit.
    INTR_RETURN_IF_ERROR(kinematics_component->SetSystemJerkLimit(
        jerk, /*enforce_limits=*/false));
    INTR_RETURN_IF_ERROR(kinematics_component->SetApplicationJerkLimit(
        jerk * scene_object::kMaxApplicationLimitsMultiplier,
        /*enforce_limits=*/true))
        << "Setting application jerk limit to a fraction "
        << scene_object::kMaxApplicationLimitsMultiplier
        << " of system limit should never fail.";
  }

  // Process dynamics
  kinematics_component->SetDamping(sdf_joint_axis.Damping());
  kinematics_component->SetFriction(sdf_joint_axis.Friction());

  return kinematics_component;
}

}  // namespace sdf
}  // namespace intrinsic
