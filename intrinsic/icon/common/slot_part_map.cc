// Copyright 2023 Intrinsic Innovation LLC

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
