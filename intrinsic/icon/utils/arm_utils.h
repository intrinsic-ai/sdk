// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_ARM_UTILS_H_
#define INTRINSIC_ICON_UTILS_ARM_UTILS_H_

#include <optional>

#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/part_status.pb.h"

namespace intrinsic::icon {

// Returns the current joint position extracted from `part_status`. For
// position-controlled robots, this returns the position command from the last
// cycle. Returns the sensed position if this is not available.
eigenmath::VectorXd GetCommandedOrSensedJointPosition(
    const intrinsic_proto::icon::PartStatus& part_status);

// Extracts the actually sensed (!) joint position from a PartStatus, returning
// it as a VectorXd.
//
// NOTE: only use the sensed joint position if you are 100%  sure about what
// you're doing. For most high-level use-cases like motion planning, prefer
// using `GetPositionCommandedLastCycle()` or simply use `GetJointPosition()`,
// which facilitates smooth concatenation of commanded trajectories and motion
// planning cache hits.
eigenmath::VectorXd GetSensedJointPosition(
    const intrinsic_proto::icon::PartStatus& part_status);

// Extracts and returns the previously commanded joint position from a
// PartStatus. Returns absl::nullopt if the part status is missing the
// previously commanded joint position for any joints.
std::optional<eigenmath::VectorXd> GetPositionCommandedLastCycle(
    const intrinsic_proto::icon::PartStatus& part_status);

// Extracts the sensed joint velocity from a PartStatus, returning it as a
// VectorXd.
eigenmath::VectorXd GetSensedJointVelocity(
    const intrinsic_proto::icon::PartStatus& part_status);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_ARM_UTILS_H_
