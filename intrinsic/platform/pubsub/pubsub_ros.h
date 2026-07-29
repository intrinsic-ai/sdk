// Copyright 2023 Intrinsic Innovation LLC

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
