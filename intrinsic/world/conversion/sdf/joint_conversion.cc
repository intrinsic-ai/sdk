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
    double theta = result.kinematics_component->GetRawValue();
    Pose3d inboard_t_outboard = Pose3d::Identity();
    if (result.kinematics_component->GetMotionType() ==
        intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_REVOLUTE) {
      inboard_t_outboard.setQuaternion(
          eigenmath::Quaterniond(eigenmath::AngleAxisd(
              theta, result.kinematics_component->GetAxis())));
    } else if (result.kinematics_component->GetMotionType() ==
               intrinsic_proto::world::KinematicsComponent::
                   MOTION_TYPE_PRISMATIC) {
      inboard_t_outboard.translation() =
          theta * result.kinematics_component->GetAxis();
    }
    result.parent_t_child =
        parent_t_joint * inboard_t_outboard * result.child_t_joint.inverse();
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
