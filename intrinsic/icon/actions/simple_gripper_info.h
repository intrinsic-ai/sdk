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

#ifndef INTRINSIC_ICON_ACTIONS_SIMPLE_GRIPPER_INFO_H_
#define INTRINSIC_ICON_ACTIONS_SIMPLE_GRIPPER_INFO_H_

#include "intrinsic/icon/actions/simple_gripper.pb.h"

namespace intrinsic::icon {

struct SimpleGripperActionInfo {
  static constexpr char kActionTypeName[] = "intrinsic.simple_gripper";
  static constexpr char kActionDescription[] =
      "Controls a simple binary state (open/closed) gripper.";

  static constexpr char kSentCommand[] =
      "intrinsic.simple_gripper.sent_command";
  static constexpr char kSentCommandDescription[] =
      "The action has sent the command to the GripperPart.";

  static constexpr char kGrasped[] = "intrinsic.grasped";
  static constexpr char kGraspedDescription[] =
      "Gripper is in the GRASPED state. The exact meaning depends on the "
      "part and Gripper setup.";

  static constexpr char kReleased[] = "intrinsic.released";
  static constexpr char kReleasedDescription[] =
      "Gripper is in the RELEASED state. The exact meaning depends on the "
      "part and Gripper setup.";

  static constexpr char kSlotName[] = "gripper";
  static constexpr char kSlotDescription[] =
      "A Part that implements the SimpleGripper Feature Interface.";

  using FixedParams =
      ::intrinsic_proto::icon::actions::proto::SimpleGripperFixedParams;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_ACTIONS_SIMPLE_GRIPPER_INFO_H_
