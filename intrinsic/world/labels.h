// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_LABELS_H_
#define INTRINSIC_WORLD_LABELS_H_

#include "absl/strings/string_view.h"
#include "intrinsic/util/string_type.h"

namespace intrinsic {

INTRINSIC_DEFINE_STRING_TYPE_AS(LabelId,
                                intrinsic::SharedPtrStringRepresentation);

// A label that a path planner can use to identify a dof. The current use intent
// is to have label such as "arm", "tool", that can be used to enable different
// path planning strategy for different dof categories. The consumer of these
// label is a path planner.
INTRINSIC_DEFINE_STRING_TYPE_AS(DofLabel,
                                intrinsic::SharedPtrStringRepresentation);

namespace labels {

inline const LabelId& ArmBase() {
  static auto* arm_base = new LabelId("arm_base");
  return *arm_base;
}

inline const LabelId& UrdfRoot() {
  static auto* urdf_root = new LabelId("urdf_root");
  return *urdf_root;
}

inline const LabelId& Tip() {
  static auto* tip = new LabelId("tip");
  return *tip;
}

inline const LabelId& ToolGripper() {
  static auto* tool_gripper = new LabelId("gripper");
  return *tool_gripper;
}

inline const LabelId& WorldOrigin() {
  static auto* world_origin = new LabelId("world_origin");
  return *world_origin;
}

// Do not use this label except in the internal implementation of SDF to World
// conversion. Consider this label reserved.
inline const LabelId& SdfModelRoot() {
  static auto* sdf_model_root = new LabelId("sdf_model_root");
  return *sdf_model_root;
}

inline const DofLabel& BasePart() {
  static auto* base_part = new DofLabel("base_part");
  return *base_part;
}

inline const DofLabel& ArmPart() {
  static auto* arm_part = new DofLabel("arm_part");
  return *arm_part;
}

inline const DofLabel& ToolPart() {
  static auto* tool_part = new DofLabel("tool_part");
  return *tool_part;
}

}  // namespace labels
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_LABELS_H_
