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

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SLOT_MAP_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SLOT_MAP_H_

#include <stdint.h>

#include "intrinsic/icon/control/c_api/c_feature_interfaces.h"

#ifdef __cplusplus
extern "C" {
#endif

struct IntrinsicIconRealtimeSlotMap;

struct IntrinsicIconRealtimeSlotMapVtable {
  IntrinsicIconFeatureInterfacesForSlot (
      *get_mutable_feature_interfaces_for_slot)(
      IntrinsicIconRealtimeSlotMap* self, uint64_t slot_id);
  IntrinsicIconConstFeatureInterfacesForSlot (*get_feature_interfaces_for_slot)(
      const IntrinsicIconRealtimeSlotMap* self, uint64_t slot_id);
};

#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SLOT_MAP_H_
