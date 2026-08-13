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

#include "intrinsic/icon/common/slot_part_map.h"

#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

SlotPartMap SlotPartMapFromProto(
    const intrinsic_proto::icon::v1::SlotPartMap& proto) {
  SlotPartMap map;

  for (const auto& [slot_name, part_name] : proto.slot_name_to_part_name()) {
    map.emplace(slot_name, part_name);
  }

  return map;
}

intrinsic_proto::icon::v1::SlotPartMap ToProto(const SlotPartMap& part_map) {
  intrinsic_proto::icon::v1::SlotPartMap proto;
  for (const auto& [slot_name, part_name] : part_map) {
    proto.mutable_slot_name_to_part_name()->insert({slot_name, part_name});
  }
  return proto;
}

}  // namespace intrinsic::icon
