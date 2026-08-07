// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/conversion/sdf/joint_conversion.h"

#include <utility>

#include "absl/status/statusor.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/scene/sdf/sdf_util.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/component/kinematics_component.h"
#include "intrinsic/world/conversion/sdf/kinematics_component_from_sdf.h"
#include "intrinsic/world/conversion/sdf/sensor_conversion.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"
#include "sdf/Joint.hh"
#include "sdf/Model.hh"
#include "sdf/SemanticPose.hh"

namespace intrinsic {
namespace sdf {

absl::StatusOr<ParseJointResult> ParseJoint(const ::sdf::Joint& sdf_joint) {
  ParseJointResult result;

  result.name = sdf_joint.Name();
  result.parent_name = sdf_joint.ParentName();
  result.child_name = sdf_joint.ChildName();
  INTR_ASSIGN_OR_RETURN(result.kinematics_component,
                        KinematicsComponentFromSdfJoint(sdf_joint));

  INTR_ASSIGN_OR_RETURN(
      result.child_t_joint,
      ParseSemanticPose(sdf_joint.SemanticPose(), sdf_joint.ChildName()));
  result.kinematics_component->SetOutboardTChild(
      result.child_t_joint.inverse());
  if (sdf_joint.ParentName() != "world") {
    INTR_ASSIGN_OR_RETURN(
        Pose3d parent_t_joint,
        ParseSemanticPose(sdf_joint.SemanticPose(), sdf_joint.ParentName()));
    result.parent_t_child = parent_t_joint * result.child_t_joint.inverse();
    result.kinematics_component->SetParentTInboard(parent_t_joint);
  }

  // Parse sensors.
  const auto sensor_count = sdf_joint.SensorCount();
  for (auto i = 0; i < sensor_count; ++i) {
    INTR_ASSIGN_OR_RETURN(ParseSensorResult parse_sensor_result,
                          ParseSensor(*sdf_joint.SensorByIndex(i)));
    result.sensors.push_back(std::move(parse_sensor_result));
  }

  // Then we parse the pose of the joint.
  return result;
}
}  // namespace sdf
}  // namespace intrinsic
