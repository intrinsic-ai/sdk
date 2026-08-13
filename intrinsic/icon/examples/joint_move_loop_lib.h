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

#ifndef INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LOOP_LIB_H_
#define INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LOOP_LIB_H_

#include <memory>
#include <optional>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/examples/joint_move_positions.pb.h"
#include "intrinsic/util/grpc/channel_interface.h"

namespace intrinsic::icon::examples {

// Moves in a loop between two joint positions for the part `part_name` for
// `duration`. The first position is slightly off the center of the joint range
// and the second position is at the center of the joint range.
// The loop can optionally be parametrized via the `joint_move_positions`
// argument, specifying the two joint positions.

// A valid connection to an ICON server is passed in using the parameter
// `icon_channel`.
absl::Status RunJointMoveLoop(
    absl::string_view part_name, absl::Duration duration,
    std::shared_ptr<intrinsic::ChannelInterface> icon_channel,
    std::optional<intrinsic_proto::icon::JointMovePositions>
        joint_move_positions = std::nullopt);

}  // namespace intrinsic::icon::examples

#endif  // INTRINSIC_ICON_EXAMPLES_JOINT_MOVE_LOOP_LIB_H_
