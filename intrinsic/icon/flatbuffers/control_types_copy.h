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


#ifndef INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_COPY_H_
#define INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_COPY_H_

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/flatbuffers/control_types_view.h"
#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/icon/flatbuffers/transform_types.h"
#include "intrinsic/icon/flatbuffers/transform_view.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic_fbs {

// Copies the raw data memory from Pose3d to memory in
// Transform.
void CopyTo(Transform* transform, const intrinsic::Pose3d& pose);

// Copies the raw data memory from Transform to memory in Pose3d.
void CopyTo(intrinsic::Pose3d* pose, const Transform& transform);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_COPY_H_
