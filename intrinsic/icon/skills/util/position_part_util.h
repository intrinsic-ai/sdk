// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_SKILLS_UTIL_POSITION_PART_UTIL_H_
#define INTRINSIC_ICON_SKILLS_UTIL_POSITION_PART_UTIL_H_

#include <optional>
#include <string>

#include "absl/status/statusor.h"
#include "intrinsic/icon/equipment/icon_equipment.pb.h"
#include "intrinsic/world/objects/object_world_client.h"
#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic::skills {

struct ArmPartInformation {
  std::string name;  // The name of the "arm" part as configured by ICON.
  intrinsic_proto::world::ObjectReference object;
};

// Gets the arm part object and ICON name for the position part.
//
// The arm part object reference can optionally be provided directly. If not
// provided, it will be deduced from the position part. An error will be
// returned if no arm part object reference is provided directly or the
// information in the position part doesn't specify a single arm part.
absl::StatusOr<ArmPartInformation> GetArmPartInformation(
    const intrinsic_proto::icon::Icon2PositionPart& position_part,
    const world::ObjectWorldClient& world,
    std::optional<intrinsic_proto::world::ObjectReference> arm_part_object);

}  // namespace intrinsic::skills

#endif  // INTRINSIC_ICON_SKILLS_UTIL_POSITION_PART_UTIL_H_
