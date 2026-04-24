// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_INTERFACES_IO_CONTROLLER_UTILS_H_
#define INTRINSIC_ICON_HAL_INTERFACES_IO_CONTROLLER_UTILS_H_

#include <list>
#include <string>
#include <vector>

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffers.h"

namespace intrinsic_fbs {

// Creates a detached flatbuffer that stores a message defined as AIOStatus.
flatbuffers::DetachedBuffer BuildAIOStatus(
    const std::vector<std::string>& descriptions);

// Creates a detached flatbuffer that stores a message defined as AIOCommand.
flatbuffers::DetachedBuffer BuildAIOCommand(
    const std::vector<std::string>& descriptions);

// Creates a detached flatbuffer that stores a message defined as DIOStatus.
flatbuffers::DetachedBuffer BuildDIOStatus(
    const std::vector<std::string>& descriptions);

// Creates a detached flatbuffer that stores a message defined as DIOCommand.
flatbuffers::DetachedBuffer BuildDIOCommand(
    const std::vector<std::string>& descriptions);

}  // namespace intrinsic_fbs

#endif  // INTRINSIC_ICON_HAL_INTERFACES_IO_CONTROLLER_UTILS_H_
