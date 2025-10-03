// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_REALTIME_STATUS_H_
#define INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_REALTIME_STATUS_H_

#include "absl/status/status.h"
#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::icon {

// These helpers convert back and forth between ICON/Abseil C++ Status values
// types and their C API equivalent.
//
// A static_assert in convert_c_types.cc ensures that the maximum length for
// messages in both C and C++ RealtimeStatus is the same.

// Truncates the message in `status` to at most
// kIntrinsicIconRealtimeStatusMaxMessageLength characters.
IntrinsicIconRealtimeStatus FromAbslStatus(const absl::Status& status);

absl::Status ToAbslStatus(const IntrinsicIconRealtimeStatus& status);

IntrinsicIconRealtimeStatus FromRealtimeStatus(const RealtimeStatus& status);

RealtimeStatus ToRealtimeStatus(const IntrinsicIconRealtimeStatus& status);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_CONTROL_C_API_CONVERT_C_REALTIME_STATUS_H_
