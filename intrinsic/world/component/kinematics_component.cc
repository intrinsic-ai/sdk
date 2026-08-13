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

#include "intrinsic/world/component/kinematics_component.h"

#include <algorithm>
#include <limits>
#include <memory>
#include <utility>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/proto/vector3.pb.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/scene/constants.h"
#include "intrinsic/util/proto_time.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/labels.h"
#include "intrinsic/world/proto/kinematics_component.pb.h"

namespace intrinsic {

using eigenmath::Vector3d;

namespace {

struct Limits {
  double effort_limit = std::numeric_limits<double>::infinity();
  double velocity_limit = std::numeric_limits<double>::infinity();
  double acceleration_limit = std::numeric_limits<double>::infinity();
  double jerk_limit = std::numeric_limits<double>::infinity();
  double raw_value_fixed_lower_limit = -std::numeric_limits<double>::infinity();
  double raw_value_fixed_upper_limit = std::numeric_limits<double>::infinity();

  intrinsic_proto::world::KinematicsComponent::Limits ToProto() const {
    intrinsic_proto::world::KinematicsComponent::Limits proto;
#define KINEMATICS_COMPONENT_SET_LIMIT_FIELD(field) \
  do {                                              \
    if (field##_limit > 0) {                        \
      proto.set_##field(field##_limit);             \
    }                                               \
  } while (false)

    KINEMATICS_COMPONENT_SET_LIMIT_FIELD(effort);
    KINEMATICS_COMPONENT_SET_LIMIT_FIELD(velocity);
    KINEMATICS_COMPONENT_SET_LIMIT_FIELD(acceleration);
    KINEMATICS_COMPONENT_SET_LIMIT_FIELD(jerk);
#undef KINEMATICS_COMPONENT_SET_LIMIT_FIELD

    proto.mutable_fixed_limits()->set_lower(raw_value_fixed_lower_limit);
    proto.mutable_fixed_limits()->set_upper(raw_value_fixed_upper_limit);
    return proto;
  }

  static absl::StatusOr<Limits> FromProto(
      const intrinsic_proto::world::KinematicsComponent::Limits&
          system_limits) {
    Limits limits;
#define KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD(field)  \
  do {                                                         \
    if (system_limits.field() < 0) {                           \
      return InvalidArgumentErrorBuilder()                     \
             << #field " limit should be non-negative! (got: " \
             << system_limits.field() << ")";                  \
    } else {                                                   \
      limits.field##_limit = system_limits.field();            \
    }                                                          \
  } while (false)

    KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD(effort);
    KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD(velocity);
    KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD(acceleration);
    KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD(jerk);
#undef KINEMATICS_COMPONENT_CHECK_AND_SET_LIMIT_FIELD

    if (system_limits.has_fixed_limits()) {
      limits.raw_value_fixed_lower_limit = system_limits.fixed_limits().lower();
      limits.raw_value_fixed_upper_limit = system_limits.fixed_limits().upper();
      if (limits.raw_value_fixed_lower_limit >
          limits.raw_value_fixed_upper_limit) {
        return absl::InvalidArgumentError(
            "fixed lower limit is greater than fixed upper limit");
      }
    }

    return limits;
  }
};

class KinematicsComponentImpl : public KinematicsComponent {
 public:
  KinematicsComponentImpl();
  KinematicsComponentImpl(
      const Pose3d& parent_t_inboard, const Pose3d& outboard_t_child,
      intrinsic_proto::world::KinematicsComponent::MotionType motion_type,
      const Vector3d& axis, double raw_value,
      std::optional<absl::Time> timestamp, Limits system_limits,
      Limits application_limits,
      double damping, double friction);

  absl::StatusOr<intrinsic_proto::world::KinematicsComponent> ToProto()
      const override;
  std::unique_ptr<KinematicsComponent> Clone() const override;

  const Pose3d& GetParentTInboard() const override { return parent_t_inboard_; }
  void SetParentTInboard(const Pose3d& parent_t_inboard) override {
    parent_t_inboard_ = parent_t_inboard;
  }

  const Pose3d& GetOutboardTChild() const override { return outboard_t_child_; }

  void SetOutboardTChild(const Pose3d& outboard_t_child) override {
    outboard_t_child_ = outboard_t_child;
  }

  intrinsic_proto::world::KinematicsComponent::MotionType GetMotionType()
      const override {
    return motion_type_;
  }
  void SetMotionType(intrinsic_proto::world::KinematicsComponent::MotionType
                         motion_type) override {
    motion_type_ = motion_type;
  }
  const Vector3d& GetAxis() const override { return axis_; }
  void SetAxis(const eigenmath::Vector3d& axis) override { axis_ = axis; }
  double GetRawValue() const override { return raw_value_; }
  void SetRawValue(double value, std::optional<absl::Time> timestamp) override {
    raw_value_ = value;
    timestamp_ = timestamp;
  }
  std::optional<absl::Time> GetTimestamp() const override { return timestamp_; }
  double GetSystemEffortLimit() const override {
    return system_limits_.effort_limit;
  }
  absl::Status SetSystemEffortLimit(double limit, bool enforce_limits) override;
  double GetApplicationEffortLimit() const override {
    return application_limits_.effort_limit;
  }
  absl::Status SetApplicationEffortLimit(double limit,
                                         bool enforce_limits) override;
  double GetSystemVelocityLimit() const override {
    return system_limits_.velocity_limit;
  }
  absl::Status SetSystemVelocityLimit(double limit,
                                      bool enforce_limits) override;
  double GetApplicationVelocityLimit() const override {
    return application_limits_.velocity_limit;
  }
  absl::Status SetApplicationVelocityLimit(double limit,
                                           bool enforce_limits) override;
  double GetSystemAccelerationLimit() const override {
    return system_limits_.acceleration_limit;
  }
  absl::Status SetSystemAccelerationLimit(double limit,
                                          bool enforce_limits) override;
  double GetApplicationAccelerationLimit() const override {
    return application_limits_.acceleration_limit;
  }
  absl::Status SetApplicationAccelerationLimit(double limit,
                                               bool enforce_limits) override;
  double GetSystemJerkLimit() const override {
    return system_limits_.jerk_limit;
  }
  absl::Status SetSystemJerkLimit(double limit, bool enforce_limits) override;
  double GetApplicationJerkLimit() const override {
    return application_limits_.jerk_limit;
  }
  absl::Status SetApplicationJerkLimit(double limit,
                                       bool enforce_limits) override;
  std::pair<double, double> GetSystemRawValueFixedLimits() const override;
  absl::Status SetSystemRawValueFixedLimits(double lower, double upper,
                                            bool enforce_limits) override;
  std::pair<double, double> GetApplicationRawValueFixedLimits() const override;
  absl::Status SetApplicationRawValueFixedLimits(double lower, double upper,
                                                 bool enforce_limits) override;
  double GetDamping() const override { return damping_; }
  void SetDamping(double value) override { damping_ = value; }
  double GetFriction() const override { return friction_; }
  void SetFriction(double value) override { friction_ = value; }

  absl::Status UpdateFromProto(
      const intrinsic_proto::world::KinematicsComponent& proto) override;

 private:
  Pose3d parent_t_inboard_;
  Pose3d outboard_t_child_;
  intrinsic_proto::world::KinematicsComponent::MotionType motion_type_;
  Vector3d axis_;
  double raw_value_{0};
  std::optional<absl::Time> timestamp_;
  Limits system_limits_;
  Limits application_limits_;
  double damping_{0};
  double friction_{0};
};

KinematicsComponentImpl::KinematicsComponentImpl()
    : motion_type_(
          intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_FIXED),
      axis_(Vector3d::UnitZ()) {}

KinematicsComponentImpl::KinematicsComponentImpl(
    const Pose3d& parent_t_inboard, const Pose3d& outboard_t_child,
    intrinsic_proto::world::KinematicsComponent::MotionType motion_type,
    const Vector3d& axis, double raw_value, std::optional<absl::Time> timestamp,
    Limits system_limits, Limits application_limits,
    double damping, double friction)
    : parent_t_inboard_(parent_t_inboard),
      outboard_t_child_(outboard_t_child),
      motion_type_(motion_type),
      axis_(axis),
      raw_value_(raw_value),
      timestamp_(timestamp),
      system_limits_(system_limits),
      application_limits_(application_limits),
      damping_(damping),
      friction_(friction) {}

absl::StatusOr<intrinsic_proto::world::KinematicsComponent>
KinematicsComponentImpl::ToProto() const {
  intrinsic_proto::world::KinematicsComponent proto;
  *proto.mutable_parent_t_inboard() = intrinsic::ToProto(parent_t_inboard_);
  *proto.mutable_outboard_t_child() = intrinsic::ToProto(outboard_t_child_);
  proto.set_motion_type(motion_type_);
  if (!axis_.isApprox(Vector3d::UnitZ())) {
    *proto.mutable_axis() = ToVectorProto(axis_);
  }
  proto.set_raw_value(raw_value_);
  *proto.mutable_system_limits() = system_limits_.ToProto();
  *proto.mutable_application_limits() = application_limits_.ToProto();
  proto.set_damping(damping_);
  proto.set_friction(friction_);
  return proto;
}

std::unique_ptr<KinematicsComponent> KinematicsComponentImpl::Clone() const {
  auto ret = std::make_unique<KinematicsComponentImpl>();
  ret->parent_t_inboard_ = parent_t_inboard_;
  ret->outboard_t_child_ = outboard_t_child_;
  ret->motion_type_ = motion_type_;
  ret->axis_ = axis_;
  ret->raw_value_ = raw_value_;
  ret->timestamp_ = timestamp_;
  ret->system_limits_ = system_limits_;
  ret->application_limits_ = application_limits_;
  ret->damping_ = damping_;
  ret->friction_ = friction_;
  return std::move(ret);
}

#define KINEMATICS_IMPL_LIMIT_SETTERS(LimitName, limit_name)                  \
  absl::Status KinematicsComponentImpl::SetSystem##LimitName##Limit(          \
      double limit, bool enforce_limits) {                                    \
    if (limit < 0.0) {                                                        \
      return ::intrinsic::InvalidArgumentErrorBuilder()                       \
             << "limit must be non-negative, was " << limit;                  \
    }                                                                         \
                                                                              \
    const double min_limit = application_limits_.limit_name##_limit /         \
                             scene_object::kMaxApplicationLimitsMultiplier;   \
    if (enforce_limits && limit < min_limit) {                                \
      if (limit >= application_limits_.limit_name##_limit) {                  \
        LOG(WARNING)                                                          \
            << "Setting system " #limit_name " limit below "                  \
            << absl::StrFormat(                                               \
                   "%2f percent",                                             \
                   (1.0 / scene_object::kMaxApplicationLimitsMultiplier))     \
            << " of the application " #limit_name " limit "                   \
            << application_limits_.limit_name##_limit << " "                  \
            << "is unsupported. Clamping -- this will be an error in the "    \
            << "future.";                                                     \
        limit = min_limit;                                                    \
      } else {                                                                \
        return ::intrinsic::InvalidArgumentErrorBuilder()                     \
               << "new system " #limit_name " limit " << limit                \
               << " is less than "                                            \
               << absl::StrFormat(                                            \
                      "%2f percent",                                          \
                      (1.0 / scene_object::kMaxApplicationLimitsMultiplier))  \
               << " of the application " #limit_name " limit "                \
               << application_limits_.limit_name##_limit << " "               \
               << absl::StrFormat("(%.2f)", min_limit);                       \
      }                                                                       \
    }                                                                         \
                                                                              \
    system_limits_.limit_name##_limit = limit;                                \
    return absl::OkStatus();                                                  \
  }                                                                           \
                                                                              \
  absl::Status KinematicsComponentImpl::SetApplication##LimitName##Limit(     \
      double limit, bool enforce_limits) {                                    \
    if (limit < 0.0) {                                                        \
      return ::intrinsic::InvalidArgumentErrorBuilder()                       \
             << "limit must be non-negative, was " << limit;                  \
    }                                                                         \
    const double max_limit = scene_object::kMaxApplicationLimitsMultiplier *  \
                             system_limits_.limit_name##_limit;               \
    if (enforce_limits && limit > max_limit) {                                \
      if (limit <= system_limits_.limit_name##_limit) {                       \
        LOG(WARNING)                                                          \
            << "Setting application " #limit_name " limit above "             \
            << absl::StrFormat("%2f percent",                                 \
                               scene_object::kMaxApplicationLimitsMultiplier) \
            << " of the application " #limit_name " limit "                   \
            << application_limits_.limit_name##_limit << " "                  \
            << "is unsupported. Clamping -- this will be an error in the "    \
            << "future.";                                                     \
        limit = max_limit;                                                    \
      } else {                                                                \
        return ::intrinsic::InvalidArgumentErrorBuilder()                     \
               << "new application " #limit_name " limit " << limit           \
               << " is more than "                                            \
               << absl::StrFormat(                                            \
                      "%2f percent",                                          \
                      scene_object::kMaxApplicationLimitsMultiplier)          \
               << " of the system " #limit_name " limit "                     \
               << system_limits_.limit_name##_limit << " "                    \
               << absl::StrFormat("(%.2f)", max_limit);                       \
      }                                                                       \
    }                                                                         \
                                                                              \
    application_limits_.limit_name##_limit = limit;                           \
    return absl::OkStatus();                                                  \
  }

KINEMATICS_IMPL_LIMIT_SETTERS(Velocity, velocity)
KINEMATICS_IMPL_LIMIT_SETTERS(Acceleration, acceleration)
KINEMATICS_IMPL_LIMIT_SETTERS(Jerk, jerk)

#undef KINEMATICS_IMPL_LIMIT_SETTERS

absl::Status KinematicsComponentImpl::SetSystemEffortLimit(
    double limit, bool enforce_limits) {
  if (limit < 0.0) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "limit must be non-negative, was " << limit;
  }

  if (enforce_limits && limit < application_limits_.effort_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new system effort limit " << limit
           << " is less than the application effort limit "
           << application_limits_.effort_limit;
  }

  system_limits_.effort_limit = limit;
  return absl::OkStatus();
}

absl::Status KinematicsComponentImpl::SetApplicationEffortLimit(
    double limit, bool enforce_limits) {
  if (limit < 0.0) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "limit must be non-negative, was " << limit;
  }
  if (enforce_limits && limit > system_limits_.effort_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new application effort limit " << limit
           << " is more than the system effort limit "
           << system_limits_.effort_limit;
  }

  application_limits_.effort_limit = limit;
  return absl::OkStatus();
}

std::pair<double, double>
KinematicsComponentImpl::GetSystemRawValueFixedLimits() const {
  return std::make_pair(system_limits_.raw_value_fixed_lower_limit,
                        system_limits_.raw_value_fixed_upper_limit);
}

absl::Status KinematicsComponentImpl::SetSystemRawValueFixedLimits(
    double lower, double upper, bool enforce_limits) {
  if (lower > upper) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new lower limit " << lower << " is greater than upper limit "
           << upper;
  }
  if (enforce_limits &&
      lower > application_limits_.raw_value_fixed_lower_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new lower system limit " << lower
           << " is more than the lower application limit "
           << application_limits_.raw_value_fixed_lower_limit;
  }
  if (enforce_limits &&
      upper < application_limits_.raw_value_fixed_upper_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new upper system limit " << upper
           << " is less than the upper application limit "
           << application_limits_.raw_value_fixed_upper_limit;
  }

  system_limits_.raw_value_fixed_lower_limit = lower;
  system_limits_.raw_value_fixed_upper_limit = upper;
  return absl::OkStatus();
}

std::pair<double, double>
KinematicsComponentImpl::GetApplicationRawValueFixedLimits() const {
  return std::make_pair(application_limits_.raw_value_fixed_lower_limit,
                        application_limits_.raw_value_fixed_upper_limit);
}

absl::Status KinematicsComponentImpl::SetApplicationRawValueFixedLimits(
    double lower, double upper, bool enforce_limits) {
  if (lower > upper) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new lower limit " << lower << " is greater than upper limit "
           << upper;
  }
  if (enforce_limits && lower < system_limits_.raw_value_fixed_lower_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new lower application limit " << lower
           << " is less than the lower system limit "
           << system_limits_.raw_value_fixed_lower_limit;
  }
  if (enforce_limits && upper > system_limits_.raw_value_fixed_upper_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "new upper application limit " << upper
           << " is more than the upper system limit "
           << system_limits_.raw_value_fixed_upper_limit;
  }

  application_limits_.raw_value_fixed_lower_limit = lower;
  application_limits_.raw_value_fixed_upper_limit = upper;
  return absl::OkStatus();
}

absl::Status KinematicsComponentImpl::UpdateFromProto(
    const intrinsic_proto::world::KinematicsComponent& proto) {
  INTR_ASSIGN_OR_RETURN(const auto component,
                        KinematicsComponent::FromProto(proto));
  const auto* raw_component =
      static_cast<const KinematicsComponentImpl*>(component.get());

  parent_t_inboard_ = raw_component->parent_t_inboard_;
  outboard_t_child_ = raw_component->outboard_t_child_;
  motion_type_ = raw_component->motion_type_;
  axis_ = raw_component->axis_;
  raw_value_ = raw_component->raw_value_;
  timestamp_ = raw_component->timestamp_;
  system_limits_ = raw_component->system_limits_;
  application_limits_ = raw_component->application_limits_;
  damping_ = raw_component->damping_;
  friction_ = raw_component->friction_;

  return absl::OkStatus();
}

}  // namespace

std::unique_ptr<KinematicsComponent> KinematicsComponent::Create() {
  return std::make_unique<KinematicsComponentImpl>();
}

absl::StatusOr<std::unique_ptr<KinematicsComponent>>
KinematicsComponent::FromProto(
    const intrinsic_proto::world::KinematicsComponent& proto) {
  Pose3d parent_t_inboard;
  if (proto.has_parent_t_inboard()) {
    INTR_ASSIGN_OR_RETURN(
        parent_t_inboard, intrinsic_proto::FromProto(proto.parent_t_inboard()),
        _ << "Failed to parse parent_t_inboard from KinematicsComponent");
  }

  Pose3d outboard_t_child;
  if (proto.has_outboard_t_child()) {
    INTR_ASSIGN_OR_RETURN(outboard_t_child,
                          intrinsic_proto::FromProto(proto.outboard_t_child()));
  }

  if (proto.motion_type() ==
      intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_UNDEFINED) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "Motion type must be specified!";
  }

  Vector3d axis = Vector3d::UnitZ();
  if (proto.has_axis() &&
      proto.motion_type() !=
          intrinsic_proto::world::KinematicsComponent::MOTION_TYPE_FIXED) {
    axis = intrinsic_proto::FromProto(proto.axis());
    if (axis.isApprox(Vector3d::Zero())) {
      return absl::InvalidArgumentError("axis is approximately zero");
    }
    if (!axis.allFinite()) {
      return absl::InvalidArgumentError("axis is not finite");
    }
    axis.normalize();
  }

  Limits application_limits;
  Limits system_limits;
  if (proto.has_system_limits()) {
    INTR_ASSIGN_OR_RETURN(system_limits,
                          Limits::FromProto(proto.system_limits()),
                          _ << "While parsing system_limits");
  }
  if (proto.has_application_limits()) {
    INTR_ASSIGN_OR_RETURN(application_limits,
                          Limits::FromProto(proto.application_limits()),
                          _ << "While parsing application_limits");
    if (!proto.application_limits().has_fixed_limits() &&
        proto.has_system_limits()) {
      application_limits.raw_value_fixed_lower_limit =
          system_limits.raw_value_fixed_lower_limit;
      application_limits.raw_value_fixed_upper_limit =
          system_limits.raw_value_fixed_upper_limit;
    }

    // If we got infinity limits for velocity or effort, fix them here.
    if (application_limits.velocity_limit ==
        std::numeric_limits<double>::infinity()) {
      application_limits.velocity_limit = system_limits.velocity_limit;
    }
    if (application_limits.effort_limit ==
        std::numeric_limits<double>::infinity()) {
      application_limits.effort_limit = system_limits.effort_limit;
    }
    if (application_limits.acceleration_limit ==
        std::numeric_limits<double>::infinity()) {
      application_limits.acceleration_limit = system_limits.acceleration_limit;
    }
    if (application_limits.jerk_limit ==
        std::numeric_limits<double>::infinity()) {
      application_limits.jerk_limit = system_limits.jerk_limit;
    }
  } else if (proto.has_system_limits()) {
    // If we do not have application limits from the proto use the system limits
    // from the proto
    application_limits = system_limits;
  }

  if (system_limits.raw_value_fixed_lower_limit >
      application_limits.raw_value_fixed_lower_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System position lower ("
           << system_limits.raw_value_fixed_lower_limit
           << ") limit is more than the application position lower ("
           << application_limits.raw_value_fixed_lower_limit << ") limit";
  }

  if (system_limits.raw_value_fixed_upper_limit <
      application_limits.raw_value_fixed_upper_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System position upper ("
           << system_limits.raw_value_fixed_upper_limit
           << ") limit is less than the application position upper ("
           << application_limits.raw_value_fixed_upper_limit << ") limit";
  }

  if (system_limits.velocity_limit < application_limits.velocity_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System velocity (" << system_limits.velocity_limit
           << ") limit is less than the application veolicty ("
           << application_limits.velocity_limit << ") limit";
  }
  if (system_limits.acceleration_limit <
      application_limits.acceleration_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System acceleration (" << system_limits.acceleration_limit
           << ") limit is less than the application acceleration ("
           << application_limits.acceleration_limit << ") limit";
  }
  if (system_limits.jerk_limit < application_limits.jerk_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System jerk (" << system_limits.jerk_limit
           << ") limit is less than the application jerk ("
           << application_limits.jerk_limit << ") limit";
  }
  if (system_limits.effort_limit < application_limits.effort_limit) {
    return ::intrinsic::InvalidArgumentErrorBuilder()
           << "System effort (" << system_limits.effort_limit
           << ") limit is less than the application effort ("
           << application_limits.effort_limit << ") limit";
  }

  // Check if the raw value is within application limits.
  double raw_value = 0.0;
  if (proto.has_raw_value()) {
    const double u = application_limits.raw_value_fixed_upper_limit;
    const double l = application_limits.raw_value_fixed_lower_limit;

    raw_value = proto.raw_value();
    // We do the comparison this way to also guard against NaNs.
    if (!(raw_value <= u && raw_value >= l)) {
      LOG(ERROR) << "Clamping raw joint value " << raw_value
                 << " which exceeds application limits (" << l << ", " << u
                 << ").";
      raw_value = std::clamp(raw_value, l, u);
    }
  } else {
    const double u = application_limits.raw_value_fixed_upper_limit;
    const double l = application_limits.raw_value_fixed_lower_limit;

    if (l <= 0 && 0 <= u) {
      raw_value = 0.0;
    } else {
      // Do arithmetic this way to prevent overflow.
      raw_value = l + (u - l) * 0.5;
      LOG(WARNING) << "Raw value in kinematics component proto was unset, but "
                   << "the default value of 0 would exceed the application "
                   << "limits (" << l << ", " << u << ") . Defaulting to the "
                   << "midpoint: " << raw_value;
    }
  }

  std::optional<absl::Time> timestamp;
  return std::make_unique<KinematicsComponentImpl>(
      parent_t_inboard, outboard_t_child, proto.motion_type(), axis, raw_value,
      timestamp, system_limits, application_limits,
      proto.damping(), proto.friction());
}

}  // namespace intrinsic
