// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_CONVERSION_USER_DATA_H_
#define INTRINSIC_SCENE_CONVERSION_USER_DATA_H_

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/descriptor.pb.h"
#include "google/protobuf/dynamic_message.h"
#include "google/protobuf/map.h"

namespace intrinsic::scene_object {

// Parses a textproto string representing SceneObject user_data.
//
// The string should be a textproto representation of a message containing a
// `user_data` map field. The `fds` must be provided as file descriptors for
// protobuf types stored in `Any`.
absl::StatusOr<google::protobuf::Map<std::string, google::protobuf::Any>>
UserDataFromString(const std::string& textproto,
                   const google::protobuf::FileDescriptorSet& fds);

// Converts SceneObject user_data to a textproto string.
//
// The `fds` must be provided as file descriptors for types stored in `Any` to
// allow expanding them to textproto.
absl::StatusOr<std::string> UserDataToString(
    const google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const google::protobuf::FileDescriptorSet& fds);

// Returns OK if the user_data can be parsed using the provided file descriptor
// set.
absl::Status UserDataCanBeParsed(
    const google::protobuf::Map<std::string, google::protobuf::Any>& user_data,
    const google::protobuf::FileDescriptorSet& fds);

}  // namespace intrinsic::scene_object

#endif  // INTRINSIC_SCENE_CONVERSION_USER_DATA_H_
