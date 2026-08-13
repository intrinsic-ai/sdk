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


#ifndef INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_VIEW_H_
#define INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_VIEW_H_

#include "intrinsic/icon/flatbuffers/transform_types.fbs.h"
#include "intrinsic/math/twist.h"

namespace intrinsic_fbs {

// Utility helper function that allows to interpret intrinsic_fbs::Wrench as
// Wrench type.
const ::intrinsic::Wrench* FromSchema(const ::intrinsic_fbs::Wrench& wrench);
::intrinsic::Wrench* FromSchema(::intrinsic_fbs::Wrench* wrench);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_FLATBUFFERS_CONTROL_TYPES_VIEW_H_
