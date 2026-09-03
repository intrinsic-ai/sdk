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

#ifndef INTRINSIC_SCENE_CONVERSION_SCENE_OBJECT_MODEL_UTILS_H_
#define INTRINSIC_SCENE_CONVERSION_SCENE_OBJECT_MODEL_UTILS_H_

#include <string>
#include <vector>

#include "absl/container/flat_hash_set.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"

namespace intrinsic {
namespace scene_object {

// Returns the names of the non-fixed joints in the given entities.
absl::flat_hash_set<std::string> GetNonFixedJointNames(
    const std::vector<intrinsic_proto::scene_object::v1::Entity>& entities);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_CONVERSION_SCENE_OBJECT_MODEL_UTILS_H_
