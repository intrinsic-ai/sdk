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

#ifndef INTRINSIC_WORLD_COMPONENT_KINEMATICS_COMPONENT_H_
#define INTRINSIC_WORLD_COMPONENT_KINEMATICS_COMPONENT_H_

#include <memory>
#include <optional>
#include <utility>

#include "absl/base/attributes.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/world/labels.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic {

// A component to hold information about the kinematics of a specific dof.
class KinematicsComponent {
 public:
  virtual ~KinematicsComponent() = default;

  // Returns a new KinematicsComponent instance.
  static std::unique_ptr<KinematicsComponent> Create();

  // Returns a new KinematicsComponent instance derived from the given proto.
  static absl::StatusOr<std::unique_ptr<KinematicsComponent>> FromProto(
      const intrinsic_proto::world::KinematicsComponent& proto);

  // Returns a new KinematicsComponent instance that is a copy of this one.
  virtual std::unique_ptr<KinematicsComponent> Clone() const = 0;

  // Returns a proto representation of this component.
  virtual absl::StatusOr<intrinsic_proto::world::KinematicsComponent> ToProto()
      const = 0;

  // Returns the joint's parent_t_inboard transform.
  virtual const Pose3d& GetParentTInboard() const = 0;

  // Sets the joint's parent_t_inboard transform.
  virtual void SetParentTInboard(const Pose3d& parent_t_inboard) = 0;
  ABSL_DEPRECATED(
      "Will be removed once child link's AttachmentComponent don't require "
      "identity pose.")
  virtual const Pose3d& GetOutboardTChild() const = 0;

  // Sets the joint's outboard_t_child transform.
  ABSL_DEPRECATED(
      "Will be removed once child link's AttachmentComponent don't require "
      "identity pose.")
  virtual void SetOutboardTChild(const Pose3d& outboard_t_child) = 0;

  // Returns the DoF's type of motion (e.g. revolute, prismatic, etc.). See
  // intrinsic/world/proto/kinematics_component.proto?q=%22enum+MotionType%22.
  virtual intrinsic_proto::world::KinematicsComponent::MotionType
  GetMotionType() const = 0;

  // Sets the DoF's type of motion (e.g. revolute, prismatic, etc.). See
  // intrinsic/world/proto/kinematics_component.proto?q=%22enum+MotionType%22.
  virtual void SetMotionType(
      intrinsic_proto::world::KinematicsComponent::MotionType motion_type) = 0;

  // Returns the DoF's axis in its inboard coordinate frame. The exact meaning
  // of this axis depends on the type returned by GetMotionType().
  virtual const eigenmath::Vector3d& GetAxis() const = 0;

  // Sets the DoF's axis in its inboard coordinate frame. The exact meaning of
  // this axis depends on the type returned by GetMotionType().
  virtual void SetAxis(const eigenmath::Vector3d& axis) = 0;

  // Returns the raw value of this DoF.
  virtual double GetRawValue() const = 0;

  // Sets the raw value of this DoF.
  virtual void SetRawValue(
      double value, std::optional<absl::Time> timestamp = std::nullopt) = 0;

  // Returns the timestamp of this DoF.
  virtual std::optional<absl::Time> GetTimestamp() const = 0;

  // Returns the DoF's system effort limit.
  virtual double GetSystemEffortLimit() const = 0;

  // Sets the DoF's system effort limit. The argument must be non-negative.
  virtual absl::Status SetSystemEffortLimit(double limit,
                                            bool enforce_limits) = 0;

  // Returns the DoF's application effort limit.
  virtual double GetApplicationEffortLimit() const = 0;

  // Sets the DoF's application effort limit. The argument must be non-negative.
  virtual absl::Status SetApplicationEffortLimit(double limit,
                                                 bool enforce_limits) = 0;

  // Returns the DoF's system velocity limit.
  virtual double GetSystemVelocityLimit() const = 0;

  // Sets the DoF's system velocity limit. The argument must be
  // non-negative.
  virtual absl::Status SetSystemVelocityLimit(double limit,
                                              bool enforce_limits) = 0;

  // Returns the DoF's application velocity limit.
  virtual double GetApplicationVelocityLimit() const = 0;

  // Sets the DoF's application velocity limit. The argument must be
  // non-negative.
  virtual absl::Status SetApplicationVelocityLimit(double limit,
                                                   bool enforce_limits) = 0;

  // Returns the DoF's system acceleration limit.
  virtual double GetSystemAccelerationLimit() const = 0;

  // Sets the DoF's system acceleration limit. The argument must be
  // non-negative.
  virtual absl::Status SetSystemAccelerationLimit(double limit,
                                                  bool enforce_limits) = 0;

  // Returns the DoF's application acceleration limit.
  virtual double GetApplicationAccelerationLimit() const = 0;

  // Sets the DoF's application acceleration limit. The argument must be
  // non-negative.
  virtual absl::Status SetApplicationAccelerationLimit(double limit,
                                                       bool enforce_limits) = 0;

  // Returns the DoF's system jerk limit.
  virtual double GetSystemJerkLimit() const = 0;

  // Sets the DoF's system jerk limit. The argument must be non-negative.
  virtual absl::Status SetSystemJerkLimit(double limit,
                                          bool enforce_limits) = 0;

  // Returns the DoF's application jerk limit.
  virtual double GetApplicationJerkLimit() const = 0;

  // Sets the DoF's application jerk limit. The argument must be non-negative.
  virtual absl::Status SetApplicationJerkLimit(double limit,
                                               bool enforce_limits) = 0;

  // Returns the DoF's system raw value limits as a pair of (lower limit,
  // upper limit).
  virtual std::pair<double, double> GetSystemRawValueFixedLimits() const = 0;

  // Sets the DoF's system raw value limits as a pair of fixed values.
  virtual absl::Status SetSystemRawValueFixedLimits(double lower, double upper,
                                                    bool enforce_limits) = 0;

  // Returns the DoF's application raw value limits as a pair of (lower limit,
  // upper limit).
  virtual std::pair<double, double> GetApplicationRawValueFixedLimits()
      const = 0;

  // Sets the DoF's application raw value limits as a pair of fixed values.
  virtual absl::Status SetApplicationRawValueFixedLimits(
      double lower, double upper, bool enforce_limits) = 0;

  // Returns damping coefficient of the joint.
  virtual double GetDamping() const = 0;

  // Sets damping coefficient of the joint.
  virtual void SetDamping(double value) = 0;

  // Returns friction value of the joint.
  virtual double GetFriction() const = 0;

  // Sets friction value of the joint.
  virtual void SetFriction(double value) = 0;

  // Updates the component based on the given proto. This will do a full
  // override, if something is missing from this proto it will override any
  // existing values with the defaults.
  virtual absl::Status UpdateFromProto(
      const intrinsic_proto::world::KinematicsComponent& proto) = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_COMPONENT_KINEMATICS_COMPONENT_H_
