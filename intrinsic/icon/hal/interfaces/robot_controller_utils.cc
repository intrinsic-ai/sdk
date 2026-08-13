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

#include "intrinsic/icon/hal/interfaces/robot_controller_utils.h"

#include <cmath>

#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/interfaces/robot_controller.fbs.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildRobotControllerStatus() {
  flatbuffers::FlatBufferBuilder builder;
  builder.Finish(builder.CreateStruct(
      RobotControllerStatus(/*speed_scaling=*/std::nan("1"))));
  return builder.Release();
}

}  // namespace intrinsic_fbs
