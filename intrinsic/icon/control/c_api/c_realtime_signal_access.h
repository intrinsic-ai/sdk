// Copyright 2023 Intrinsic Innovation LLC

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
