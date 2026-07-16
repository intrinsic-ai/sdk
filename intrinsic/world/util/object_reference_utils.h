// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_
#define INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_

#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {

// Returns true if the object reference is empty (i.e. has neither a name nor an
// id, or the name or id is empty).
bool IsEmptyObjectReference(const intrinsic_proto::world::ObjectReference& ref);

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_
