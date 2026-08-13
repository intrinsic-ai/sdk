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

#ifndef INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LIB_H_
#define INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LIB_H_

#include <memory>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/util/grpc/channel_interface.h"

namespace intrinsic::icon::examples {

// First moves all joints to a position offset from the joint range center, then
// switches to the stop action and then moves the joints to the center of the
// joint range.
// Controls the part defined by `part_name` using the provided `icon_channel`.
absl::Status RunJointMove(
    absl::string_view part_name,
    std::shared_ptr<intrinsic::ChannelInterface> icon_channel);
}  // namespace intrinsic::icon::examples

#endif  // INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LIB_H_
