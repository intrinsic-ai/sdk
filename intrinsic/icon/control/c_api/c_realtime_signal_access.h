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

#ifndef INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SIGNAL_ACCESS_H_
#define INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SIGNAL_ACCESS_H_

#include <stdint.h>

#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/control/c_api/c_types.h"

#ifdef __cplusplus
extern "C" {
#endif

struct IntrinsicIconRealtimeSignalAccess;

struct IntrinsicIconRealtimeSignalAccessVtable {
  IntrinsicIconRealtimeStatus (*read_signal)(
      IntrinsicIconRealtimeSignalAccess* self, uint64_t id,
      IntrinsicIconSignalValue* signal_value);
};

#ifdef __cplusplus
}
#endif

#endif  // INTRINSIC_ICON_CONTROL_C_API_C_REALTIME_SIGNAL_ACCESS_H_
