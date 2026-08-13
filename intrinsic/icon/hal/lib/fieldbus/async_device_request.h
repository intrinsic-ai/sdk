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

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_ASYNC_DEVICE_REQUEST_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_ASYNC_DEVICE_REQUEST_H_

#include <stdint.h>

#include <string>

#include "intrinsic/icon/utils/async_request.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::fieldbus {

// The type of request.
enum class RequestType : uint8_t {
  kNormalOperation,  // Normal operation.
  kActivate,         // The device shall be activated.
  kDeactivate,       // The device shall be deactivated.
  kEnableMotion,     // The device shall be enabled.
  kDisableMotion,    // The device shall be disabled.
  kClearFaults,      // The device shall clear any fault (if present).
};

// Returns a string representation of the request type.
std::string ToString(RequestType type);

// The status of the request.
enum class RequestStatus : uint8_t {
  kDone,        // The device didn't have to, or has completed the request.
  kProcessing,  // The device has not yet completed the request.
};

// Returns a string representation of the request status.
std::string ToString(RequestStatus status);

// Encapsulate request type so that we can specify a default value.
struct RequestData {
  RequestType request_type = RequestType::kNormalOperation;
};

// This movable request wraps the request type and a promise.
// Used as communication channel between a non-rt async call to a hardware
// module's `EnableMotion`, `ClearFaults`, `DisableMotion`, ... and the module's
// real-time loop.
using AsyncDeviceRequest =
    intrinsic::icon::AsyncRequest<RequestData, intrinsic::icon::RealtimeStatus>;

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_ASYNC_DEVICE_REQUEST_H_
