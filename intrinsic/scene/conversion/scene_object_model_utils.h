// Copyright 2023 Intrinsic Innovation LLC

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
