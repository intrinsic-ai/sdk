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

#include "intrinsic/icon/hal/interfaces/icon_state_utils.h"

#include <cstdint>
#include <limits>

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/interfaces/icon_state.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildIconState() {
  flatbuffers::FlatBufferBuilder builder;
  // Initializes with `max`, because NaN is not supported for integers, so that
  // the IconState flatbuffer is invalid/inconsistent until it receives its
  // first update. Requires initializing to a different value than
  // intrinsic/icon/interprocess/shared_memory_manager/segment_header.h.
  builder.Finish(builder.CreateStruct(
      IconState(/*current_cycle=*/std::numeric_limits<uint64_t>::max())));
  return builder.Release();
}
}  // namespace intrinsic_fbs
