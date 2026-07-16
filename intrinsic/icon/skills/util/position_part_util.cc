// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/skills/util/position_part_util.h"

#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/equipment/icon_equipment.pb.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/objects/object_world_client.h"
#include "intrinsic/world/objects/world_object.h"
#include "intrinsic/world/util/object_reference_utils.h"

namespace intrinsic::skills {

absl::StatusOr<ArmPartInformation> GetArmPartInformation(
    const intrinsic_proto::icon::Icon2PositionPart& position_part,
    const world::ObjectWorldClient& world,
    std::optional<intrinsic_proto::world::ObjectReference> arm_part_object) {
  if (arm_part_object && IsEmptyObjectReference(*arm_part_object)) {
    arm_part_object = std::nullopt;
  }

  // If multiple parts are specified, we error unless the user has provided us
  // with a part name to disambiguate.
  std::optional<std::string> arm_part_name;

  // If the part name was not given, then get it from the map if it is
  // unambiguous, i.e., only one part is provided.
  if (!arm_part_object) {
    if (position_part.object_names().size() > 1) {
      return absl::InvalidArgumentError(
          "No arm part object reference was given, and the given ICON instance "
          "has multiple. Please specify an arm part object reference to select "
          "one part to use.");
    }
    if (position_part.object_names().empty()) {
      return absl::InvalidArgumentError(
          "No arm part object reference was given, and none was provided in "
          "the map.");
    }
    arm_part_object.emplace();
    arm_part_object->mutable_by_name()->set_object_name(
        position_part.object_names().begin()->second);
    arm_part_name = position_part.object_names().begin()->first;
  }
  if (!arm_part_object.has_value()) {
    // This shouldn't happen unless a coding mistake leaves the arm part object
    // reference unset.
    return absl::InternalError("Failed to deduce arm part name.");
  }

  // If we don't have the arm_part_name yet (because the user provided the part
  // as an object reference), try to do the lookup.
  if (!arm_part_name) {
    // First, we turn the ObjectReference into a name in case it was somehow
    // provided as an id.
    INTR_ASSIGN_OR_RETURN(const world::WorldObject world_object,
                          world.GetObject(*arm_part_object));
    const std::string object_name = world_object.Name().value();
    std::vector<absl::string_view> object_names;
    object_names.reserve(position_part.object_names().size());

    // Now do a look up against the position parts.
    for (const auto& [current_part_name, current_object_name] :
         position_part.object_names()) {
      object_names.push_back(current_object_name);
      if (current_object_name == object_name) {
        arm_part_name = current_part_name;
        break;
      }
    }

    // If we didn't find a match, the object the user provided is not a position
    // part and we error.
    if (!arm_part_name) {
      return absl::InvalidArgumentError(absl::StrCat(
          "Could not find a position part corresponding to object with name '",
          object_name, "'. Please check you have the right object and ICON ",
          "instance selected. Objects corresponding to position parts for the ",
          "selected ICON instance are: ", absl::StrJoin(object_names, ", "),
          "."));
    }
  }

  // Both of these will be set at this point.
  if (!arm_part_name || !arm_part_object) {
    return absl::InternalError(
        "Could not deduce arm_part_name or arm_part_object. Please report "
        "this!");
  }

  return ArmPartInformation{.name = std::move(*arm_part_name),
                            .object = std::move(*arm_part_object)};
}

}  // namespace intrinsic::skills
