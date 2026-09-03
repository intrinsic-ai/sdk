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

#ifndef INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_
#define INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_

#include "absl/status/status.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/world/proto/geometry_component.pb.h"

namespace intrinsic::scene_object {

// Applies a GeometryUpdate to a GeometryComponent.
//
// This function modifies `geometry_component` in-place based on the
// instructions in `update`.
absl::Status ApplyGeometryUpdate(
    const intrinsic_proto::scene_object::v1::GeometryUpdate& update,
    intrinsic_proto::world::GeometryComponent& geometry_component);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_UTIL_GEOMETRY_UPDATE_UTIL_H_
