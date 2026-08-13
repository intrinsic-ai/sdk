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

#ifndef INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_
#define INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_

#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {

// Returns true if the object reference is empty (i.e. has neither a name nor an
// id, or the name or id is empty).
bool IsEmptyObjectReference(const intrinsic_proto::world::ObjectReference& ref);

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_UTIL_OBJECT_REFERENCE_UTILS_H_
