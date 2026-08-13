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

#ifndef INTRINSIC_ICON_CONTROL_SLOT_TYPES_H_
#define INTRINSIC_ICON_CONTROL_SLOT_TYPES_H_

#include <cstdint>

#include "intrinsic/icon/proto/v1/types.pb.h"
#include "ortools/base/strong_int.h"

namespace intrinsic::icon {

DEFINE_STRONG_INT_TYPE(RealtimeSlotId, int64_t);

struct SlotInfo {
  // Contains things like the supported FeatureInterfaceTypes, and (depending on
  // which types are supported) number of DoFs, maximum limits, etc.
  //
  // Actions may need this information in their Non-RT initialization routine.
  intrinsic_proto::icon::v1::PartConfig config;
  // Action Factories should save this ID and use it to access this slot from
  // realtime functions.
  RealtimeSlotId slot_id;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_SLOT_TYPES_H_
