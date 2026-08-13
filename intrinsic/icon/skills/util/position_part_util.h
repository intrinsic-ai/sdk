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
