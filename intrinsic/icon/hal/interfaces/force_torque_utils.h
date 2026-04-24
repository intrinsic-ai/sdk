// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_FORCE_TORQUE_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_FORCE_TORQUE_UTILS_H_

#include "flatbuffers/detached_buffer.h"
#include "intrinsic/icon/hal/interfaces/force_sensor.fbs.h"
#include "intrinsic/icon/hal/interfaces/force_torque.fbs.h"
#include "intrinsic/icon/utils/fixed_string.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer CreateFbsForceTorqueStatus();
flatbuffers::DetachedBuffer CreateFbsForceTorqueCommand();

constexpr int kMaxFaultLength = 512;

// Returns a FixedStr describing the canonical error code encoded in
// `status_code`.
intrinsic::icon::FixedString<kMaxFaultLength> ToFixedString(
    intrinsic_fbs::ForceSensorStatusCode status_code);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_FORCE_TORQUE_UTILS_H_
