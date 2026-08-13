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
