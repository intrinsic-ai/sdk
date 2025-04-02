// Copyright 2023 Intrinsic Innovation LLC

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
