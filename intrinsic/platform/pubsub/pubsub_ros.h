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

#ifndef INTRINSIC_PLATFORM_PUBSUB_PUBSUB_ROS_H_
#define INTRINSIC_PLATFORM_PUBSUB_PUBSUB_ROS_H_

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "rclcpp/serialized_message.hpp"

namespace intrinsic {

// Specialization for pre-serialized ROS messages.
template <>
absl::StatusOr<rclcpp::SerializedMessage> PubSub::CallOne(
    absl::string_view key, const rclcpp::SerializedMessage& serialized_request,
    const QueryOptions& query_options);

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_PUBSUB_ROS_H_
