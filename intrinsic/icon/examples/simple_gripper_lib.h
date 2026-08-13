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

#ifndef INTRINSIC_ICON_EXAMPLES_SIMPLE_GRIPPER_LIB_H_
#define INTRINSIC_ICON_EXAMPLES_SIMPLE_GRIPPER_LIB_H_

#include <cstddef>
#include <memory>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/actions/simple_gripper_info.h"
#include "intrinsic/icon/cc_client/client.h"
#include "intrinsic/util/grpc/channel_interface.h"

namespace intrinsic::icon::examples {

// Requests the status for `part_name` and prints it to std::cout.
// Returns NotFoundError if no status for that part is returned.
absl::Status PrintPartStatus(absl::string_view part_name, Client& icon_client);

// Opens a session for `part_name`, sends the command defined by
// `action_parameters` and waits for the condition
// SimpleGripperActionInfo::kSentCommand.
absl::Status SendGripperCommand(
    absl::string_view part_name,
    const SimpleGripperActionInfo::FixedParams& action_parameters,
    std::shared_ptr<ChannelInterface> icon_channel);

// Sends a GRASP command to `part_name`, prints the part status to std::cout,
// then waits 10s, sends a RELEASE command and prints the part status to
// std::cout.
absl::Status ExampleGraspAndRelease(
    absl::string_view part_name,
    std::shared_ptr<ChannelInterface> icon_channel);

}  // namespace intrinsic::icon::examples

#endif  // INTRINSIC_ICON_EXAMPLES_SIMPLE_GRIPPER_LIB_H_
