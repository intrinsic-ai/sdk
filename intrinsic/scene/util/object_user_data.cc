// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/scene/util/object_user_data.h"

#include <memory>
#include <string>

#include "absl/base/attributes.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/descriptor_database.h"
#include "google/protobuf/dynamic_message.h"
#include "google/protobuf/map.h"
#include "google/protobuf/message.h"
#include "intrinsic/scene/conversion/user_data.h"
#include "intrinsic/scene/proto/v1/scene_object_updates.pb.h"
#include "intrinsic/util/proto/descriptors.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace scene_object {
namespace {

using ::intrinsic_proto::scene_object::v1::UpdateUserData;

}  // namespace

absl::Status SceneObjectUserDataCanBeParsed(
    const intrinsic_proto::scene_object::v1::SceneObject& object,
    const google::protobuf::FileDescriptorSet& fds) {
  return intrinsic::scene_object::UserDataCanBeParsed(object.user_data(), fds);
}

absl::Status MergeSceneObjectUserData(
    google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const intrinsic_proto::scene_object::v1::UpdateUserData& update) {
  switch (update.policy()) {
    case UpdateUserData::POLICY_UNSPECIFIED:
      ABSL_FALLTHROUGH_INTENDED;
    case UpdateUserData::POLICY_INSERT:
      for (const auto& [k, _] : update.user_data()) {
        if (user_data.contains(k)) {
          return absl::InvalidArgumentError(
              absl::StrCat("User data already exists for key: ", k));
        }
      }
      ABSL_FALLTHROUGH_INTENDED;
    case UpdateUserData::POLICY_INSERT_OR_UPDATE:
      for (const auto& [k, v] : update.user_data()) {
        if (user_data.contains(k)) {
          user_data.at(k) = v;
        } else {
          user_data.insert({k, v});
        }
      }
      return absl::OkStatus();
    case UpdateUserData::POLICY_REMOVE:
      for (const auto& [k, _] : update.user_data()) {
        user_data.erase(k);
      }
      return absl::OkStatus();
    case UpdateUserData::POLICY_CLEAR_AND_REPLACE:
      user_data = update.user_data();
      return absl::OkStatus();
    case UpdateUserData::POLICY_CLEAR_ALL:
      user_data.clear();
      return absl::OkStatus();
    default:
      return absl::InvalidArgumentError(
          absl::StrCat("Unknown UpdateUserData policy: ", update.policy()));
  }
}

absl::Status MergeSceneObjectUserData(
    google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const google::protobuf::Map<std::string, google::protobuf::Any>&
        user_data_update,
    const intrinsic_proto::scene_object::v1::UpdateUserData::
        UpdateUserDataPolicy policy) {
  intrinsic_proto::scene_object::v1::UpdateUserData update;
  update.set_policy(policy);
  *update.mutable_user_data() = user_data_update;
  return MergeSceneObjectUserData(user_data, update);
}

}  // namespace scene_object
}  // namespace intrinsic
