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

#ifndef INTRINSIC_SCENE_UTIL_OBJECT_USER_DATA_H_
#define INTRINSIC_SCENE_UTIL_OBJECT_USER_DATA_H_

#include <string>

#include "absl/status/status.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/map.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/util/proto/descriptors.h"

namespace intrinsic {
namespace scene_object {

// Returns OK if the user_data can be parsed using the provided file descriptor
// set.
absl::Status SceneObjectUserDataCanBeParsed(
    const intrinsic_proto::scene_object::v1::SceneObject& object,
    const google::protobuf::FileDescriptorSet& fds);

// Merges the user_data from the update into the existing user_data.
absl::Status MergeSceneObjectUserData(
    google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const intrinsic_proto::scene_object::v1::UpdateUserData& update);

// Same as above except with unpacked arguments.
absl::Status MergeSceneObjectUserData(
    google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const google::protobuf::Map<std::string, google::protobuf::Any>&
        user_data_update,
    intrinsic_proto::scene_object::v1::UpdateUserData::UpdateUserDataPolicy
        policy);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_UTIL_OBJECT_USER_DATA_H_
