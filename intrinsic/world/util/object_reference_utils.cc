// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/world/util/object_reference_utils.h"

#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {

bool IsEmptyObjectReference(
    const intrinsic_proto::world::ObjectReference& ref) {
  if (ref.has_by_name()) {
    return ref.by_name().object_name().empty();
  } else if (ref.has_id()) {
    return ref.id().empty();
  }
  return true;
}

}  // namespace intrinsic
