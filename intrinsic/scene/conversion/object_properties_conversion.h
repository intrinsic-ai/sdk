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
