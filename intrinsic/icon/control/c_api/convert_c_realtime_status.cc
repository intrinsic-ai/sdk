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

#include "intrinsic/icon/control/c_api/convert_c_realtime_status.h"

#include <algorithm>
#include <cstring>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/control/c_api/c_realtime_status.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::icon {
namespace {
static_assert(
    kIntrinsicIconRealtimeStatusMaxMessageLength ==
        RealtimeStatus::kMaxMessageLength,
    "C++ RealtimeStatus and C IntrinsicIconRealtimeStatus have different "
    "maximum message lengths. This breaks the ICON C API!");
}

IntrinsicIconRealtimeStatus FromAbslStatus(const absl::Status& status) {
  IntrinsicIconRealtimeStatus status_out;
  status_out.status_code = static_cast<int>(status.code());
  if (!status.ok()) {
    // Since `status.message()` is a string_view, it may not be null-terminated,
    // so we cannot use (safe)strncpy. Instead, we memcpy the contents of
    // `status.message()`, truncating at the end of `status_out.message`, and
    // set `status_out.size` accordingly.
    //
    // Limit the number of characters we copy to prevent writing into invalid
    // memory.
    status_out.size =
        std::min(sizeof(status_out.message), status.message().size());
    memcpy(status_out.message, status.message().data(),
           /*n=*/status_out.size);
  }
  return status_out;
}

absl::Status ToAbslStatus(const IntrinsicIconRealtimeStatus& status) {
  absl::StatusCode code = static_cast<absl::StatusCode>(status.status_code);
  if (code == absl::StatusCode::kOk) {
    return absl::OkStatus();
  } else {
    // There's nothing preventing a caller from setting `status.size` to a value
    // that's greater than the size of `status.message`, so limit the
    // string_view to avoid reading from memory we don't own.
    return absl::Status(
        code, absl::string_view(status.message,
                                std::min(sizeof(status.message), status.size)));
  }
}

IntrinsicIconRealtimeStatus FromRealtimeStatus(const RealtimeStatus& status) {
  IntrinsicIconRealtimeStatus status_out{
      .status_code = static_cast<int>(status.code()),
      .message = "",
  };
  if (!status.ok()) {
    // Since `status.message()` is a string_view, it may not be null-terminated,
    // so we cannot use (safe)strncpy. Instead, we memcpy the contents of
    // `status.message()`, truncating at the end of `status_out.message`, and
    // set `status_out.size` accordingly.
    //
    // Limit the number of characters we copy to prevent writing into invalid
    // memory.
    status_out.size =
        std::min(sizeof(status_out.message), status.message().size());
    memcpy(status_out.message, status.message().data(),
           /*n=*/status_out.size);
  }
  return status_out;
}

RealtimeStatus ToRealtimeStatus(const IntrinsicIconRealtimeStatus& status) {
  absl::StatusCode code = static_cast<absl::StatusCode>(status.status_code);
  if (code == absl::StatusCode::kOk) {
    return icon::OkStatus();
  } else {
    return icon::RealtimeStatus(
        code, absl::string_view(status.message,
                                std::min(sizeof(status.message), status.size)));
  }
}

}  // namespace intrinsic::icon
