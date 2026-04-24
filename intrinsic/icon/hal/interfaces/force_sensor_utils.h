// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_FORCE_SENSOR_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_FORCE_SENSOR_UTILS_H_

#include "flatbuffers/detached_buffer.h"

namespace intrinsic_fbs {

// Creates a detached flatbuffer that stores a message defined as
// ForceSensorStatus.
flatbuffers::DetachedBuffer CreateForceSensorStatusBuffer();

// Creates a detached flatbuffer that stores a message defined as
// ForceSensorStatus.
flatbuffers::DetachedBuffer CreateForceSensorCommandBuffer();

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_FORCE_SENSOR_UTILS_H_
