// Copyright 2023 Intrinsic Innovation LLC

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
