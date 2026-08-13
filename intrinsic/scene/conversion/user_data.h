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
