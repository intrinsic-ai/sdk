// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_COMPONENT_H_
#define INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_COMPONENT_H_

#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"
#include "intrinsic/icon/utils/realtime_status_or.h"

namespace intrinsic::fieldbus {

// Interface for fieldbus bus devices.
class BusComponent {
 public:
  virtual ~BusComponent() = default;

  // Triggers a single real-time read of the device.
  // Returns a `RequestStatus` if okay, or a `RealtimeStatus` if in error.
  // If not in error, BusComponents are expected to return kProcessing if the
  // request requires further processing and kDone otherwise. When `type` is
  // `kNormalOperation` bus devices are expected to always return `kDone`.
  virtual intrinsic::icon::RealtimeStatusOr<RequestStatus> CyclicRead(
      RequestType type) = 0;

  // Triggers a single real-time write of the device.
  // Returns a `RequestStatus` if okay, or a `RealtimeStatus` if in error.
  // If not in error, BusComponents are expected to return kProcessing if the
  // request requires further processing and kDone otherwise. When `type` is
  // `kNormalOperation` bus devices are expected to always return `kDone`.
  virtual intrinsic::icon::RealtimeStatusOr<RequestStatus> CyclicWrite(
      RequestType type) = 0;
};

}  // namespace intrinsic::fieldbus

#endif  // INTRINSIC_ICON_HAL_LIB_FIELDBUS_BUS_COMPONENT_H_
