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

#include "intrinsic/icon/hal/interfaces/payload_command_utils.h"

#include "flatbuffers/buffer.h"
#include "flatbuffers/detached_buffer.h"
#include "flatbuffers/flatbuffer_builder.h"
#include "intrinsic/icon/hal/interfaces/payload_command.fbs.h"
#include "intrinsic/icon/hal/interfaces/robot_payload.fbs.h"
#include "intrinsic/icon/hal/interfaces/robot_payload_utils.h"

namespace intrinsic_fbs {

flatbuffers::DetachedBuffer BuildPayloadCommand() {
  flatbuffers::FlatBufferBuilder fbb;
  fbb.ForceDefaults(true);

  flatbuffers::Offset<OptionalRobotPayload> payload =
      AddOptionalRobotPayload(fbb);

  PayloadCommandBuilder builder(fbb);
  builder.add_full_payload(payload);

  auto payload_command = builder.Finish();

  fbb.Finish(payload_command);
  return fbb.Release();
}

}  // namespace intrinsic_fbs
