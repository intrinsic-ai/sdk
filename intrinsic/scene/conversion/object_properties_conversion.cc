// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/conversion/object_properties_conversion.h"

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "intrinsic/icon/proto/eigen_conversion.h"
#include "intrinsic/kinematics/types/cartesian_limits.h"
#include "intrinsic/scene/proto/v1/object_properties.pb.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::scene_object {

using ::intrinsic::icon::RepeatedDoubleToVector3d;
using ::intrinsic::icon::Vector3dToRepeatedDouble;

intrinsic_proto::scene_object::v1::CartesianLimits ToProto(
    const CartesianLimits& limits) {
  intrinsic_proto::scene_object::v1::CartesianLimits out;
  Vector3dToRepeatedDouble(limits.min_translational_position,
                           out.mutable_min_translational_position());
  Vector3dToRepeatedDouble(limits.max_translational_position,
                           out.mutable_max_translational_position());
  Vector3dToRepeatedDouble(limits.min_translational_velocity,
                           out.mutable_min_translational_velocity());
  Vector3dToRepeatedDouble(limits.max_translational_velocity,
                           out.mutable_max_translational_velocity());
  Vector3dToRepeatedDouble(limits.min_translational_acceleration,
                           out.mutable_min_translational_acceleration());
  Vector3dToRepeatedDouble(limits.max_translational_acceleration,
                           out.mutable_max_translational_acceleration());
  Vector3dToRepeatedDouble(limits.min_translational_jerk,
                           out.mutable_min_translational_jerk());
  Vector3dToRepeatedDouble(limits.max_translational_jerk,
                           out.mutable_max_translational_jerk());
  out.set_max_rotational_velocity(limits.max_rotational_velocity);
  out.set_max_rotational_acceleration(limits.max_rotational_acceleration);
  out.set_max_rotational_jerk(limits.max_rotational_jerk);
  return out;
}

absl::StatusOr<CartesianLimits> FromProto(
    const intrinsic_proto::scene_object::v1::CartesianLimits& proto) {
  CartesianLimits out;
  INTR_ASSIGN_OR_RETURN(
      out.min_translational_position,
      RepeatedDoubleToVector3d(proto.min_translational_position()));
  INTR_ASSIGN_OR_RETURN(
      out.max_translational_position,
      RepeatedDoubleToVector3d(proto.max_translational_position()));
  INTR_ASSIGN_OR_RETURN(
      out.min_translational_velocity,
      RepeatedDoubleToVector3d(proto.min_translational_velocity()));
  INTR_ASSIGN_OR_RETURN(
      out.max_translational_velocity,
      RepeatedDoubleToVector3d(proto.max_translational_velocity()));
  INTR_ASSIGN_OR_RETURN(
      out.min_translational_acceleration,
      RepeatedDoubleToVector3d(proto.min_translational_acceleration()));
  INTR_ASSIGN_OR_RETURN(
      out.max_translational_acceleration,
      RepeatedDoubleToVector3d(proto.max_translational_acceleration()));
  INTR_ASSIGN_OR_RETURN(
      out.min_translational_jerk,
      RepeatedDoubleToVector3d(proto.min_translational_jerk()));
  INTR_ASSIGN_OR_RETURN(
      out.max_translational_jerk,
      RepeatedDoubleToVector3d(proto.max_translational_jerk()));
  out.max_rotational_velocity = proto.max_rotational_velocity();
  out.max_rotational_acceleration = proto.max_rotational_acceleration();
  out.max_rotational_jerk = proto.max_rotational_jerk();
  if (!out.IsValid()) {
    return absl::InvalidArgumentError(
        absl::StrCat("Cartesian limits are invalid: ", proto));
  }
  return out;
}

}  // namespace intrinsic::scene_object
