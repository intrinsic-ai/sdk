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

#include "intrinsic/world/objects/object_world_client_utils.h"

#include <optional>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/objects/frame.h"
#include "intrinsic/world/objects/object_world_client.h"
#include "intrinsic/world/objects/object_world_ids.h"
#include "intrinsic/world/objects/transform_node.h"
#include "intrinsic/world/objects/world_object.h"

namespace intrinsic::world {

absl::StatusOr<bool> IsObjectAncestorOfNode(
    const world::WorldObject& object, const world::TransformNode& node,
    const world::ObjectWorldClient& world) {
  if (object.Id() == RootObjectId()) {
    return true;
  }

  std::optional<world::WorldObject> current_obj =
      world::WorldObject::FromTransformNode(node);

  // If 'node' is a frame, get its parent object.
  if (!current_obj) {
    if (std::optional<world::Frame> frame =
            world::Frame::FromTransformNode(node);
        frame) {
      INTR_ASSIGN_OR_RETURN(current_obj, world.GetObject(frame->ObjectId()));
    } else {
      return absl::InvalidArgumentError("Unknown type of TransformNode.");
    }
  }

  // Traverse upwards util we find the object in question or stop at root.
  while (current_obj->Id() != RootObjectId()) {
    if (current_obj->Id() == object.Id()) {
      return true;
    }
    INTR_ASSIGN_OR_RETURN(current_obj,
                          world.GetObject(current_obj->ParentId()));
  }

  return false;
}

}  // namespace intrinsic::world
