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

#ifndef INTRINSIC_ICON_COMMON_SLOT_PART_MAP_H_
#define INTRINSIC_ICON_COMMON_SLOT_PART_MAP_H_

#include <string>

#include "absl/container/btree_map.h"
#include "absl/container/flat_hash_map.h"
#include "intrinsic/icon/proto/v1/types.pb.h"

namespace intrinsic::icon {

// A SlotPartMap defines a mapping from the slot names used by an Action to
// Application Layer Part names. Uses btree_map since that, unlike
// flat_hash_map, is an ordered container and has equality operators and
// absl::Hash support.
using SlotPartMap = absl::btree_map<std::string, std::string>;

SlotPartMap SlotPartMapFromProto(
    const intrinsic_proto::icon::v1::SlotPartMap& proto);
intrinsic_proto::icon::v1::SlotPartMap ToProto(const SlotPartMap& part_map);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_COMMON_SLOT_PART_MAP_H_
