// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_CONVERSION_OBJECT_PROPERTIES_CONVERSION_H_
#define INTRINSIC_SCENE_CONVERSION_OBJECT_PROPERTIES_CONVERSION_H_

#include "absl/status/statusor.h"
#include "intrinsic/kinematics/types/cartesian_limits.h"
#include "intrinsic/scene/proto/v1/object_properties.pb.h"

namespace intrinsic::scene_object {

intrinsic_proto::scene_object::v1::CartesianLimits ToProto(
    const CartesianLimits& limits);

// Converts a proto::CartesianLimits proto to a CartesianLimits.
// Reports InvalidArgumentError if the resulting CartesianLimits are not valid.
absl::StatusOr<CartesianLimits> FromProto(
    const intrinsic_proto::scene_object::v1::CartesianLimits& proto);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_CONVERSION_OBJECT_PROPERTIES_CONVERSION_H_
