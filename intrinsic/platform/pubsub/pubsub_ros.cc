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

#include "intrinsic/platform/pubsub/pubsub_ros.h"

#include <cstddef>
#include <cstring>
#include <memory>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "intrinsic/middleware/imw.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"
#include "rclcpp/serialized_message.hpp"

namespace intrinsic {

template <>
absl::StatusOr<rclcpp::SerializedMessage> PubSub::CallOne(
    absl::string_view key, const rclcpp::SerializedMessage& serialized_request,
    const QueryOptions& query_options) {
  const absl::Span<const std::byte> request(
      (const std::byte*)(serialized_request.get_rcl_serialized_message()
                             .buffer),
      serialized_request.get_rcl_serialized_message().buffer_length);

  absl::Notification notification;
  rclcpp::SerializedMessage serialized_response;
  bool reply_functor_called = false;
  auto reply_functor = std::make_unique<imw_callback_functor_t>(
      [&serialized_response, &reply_functor_called](
          const char* keyexpr, const void* response_bytes,
          const size_t response_bytes_len) {
        serialized_response = rclcpp::SerializedMessage(response_bytes_len);
        std::memcpy(serialized_response.get_rcl_serialized_message().buffer,
                    response_bytes, response_bytes_len);
        serialized_response.get_rcl_serialized_message().buffer_length =
            response_bytes_len;
        reply_functor_called = true;
      });
  auto on_done_functor = std::make_unique<imw_on_done_functor_t>(
      [&notification](const char* unused_keyexpr) { notification.Notify(); });

  imw_query_options_t imw_query_options;
  if (query_options.timeout.has_value()) {
    imw_query_options.timeout_ms =
        *query_options.timeout / absl::Milliseconds(1);
  }
  imw_query_options.call_ros_service = true;

  QueryContext query_context;
  query_context.callback = reply_functor.get();
  query_context.on_done = on_done_functor.get();

  const void* request_data = request.empty() ? nullptr : request.data();
  imw_ret ret = Zenoh().imw_query(
      key.data(), zenoh_query_static_callback, zenoh_query_static_on_done,
      request_data, request.size(), &query_context, &imw_query_options);

  if (ret != IMW_OK) {
    return absl::InternalError("Error querying Zenoh");
  }
  absl::Duration timeout;
  if (query_options.timeout.has_value()) {
    timeout = query_options.timeout.value() + absl::Seconds(1);
  } else {
    timeout = absl::Seconds(10);
  }
  const bool callback_called =
      notification.WaitForNotificationWithTimeout(timeout);
  if (!callback_called) {
    // This branch should normally not be taken, because Zenoh should time
    // out internally before the notification timeout. If this timeout fires,
    // it is due to an internal Zenoh error of some sort, which is unexpected.
    return absl::DeadlineExceededError(absl::StrCat(
        "Unexpected timeout (", timeout, ") waiting for ROS service: ", key));
  } else if (!reply_functor_called) {
    return absl::DeadlineExceededError(
        absl::StrCat("Timeout (", timeout, ") waiting for ROS service: ", key));
  }
  return serialized_response;
}

}  // namespace intrinsic
