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

#include "intrinsic/icon/utils/inspection_publisher.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/assets/services/proto/v1/service_inspection.pb.h"
#include "intrinsic/icon/hal/proto/hardware_module_inspection.pb.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::absl_testing::StatusIs;
using ::intrinsic::testing::EqualsProto;

namespace intrinsic::icon {

namespace {

TEST(InspectionPublisherTest, CreateFailsWithEmptyServiceName) {
  EXPECT_THAT(InspectionPublisher::Create("", ""),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(InspectionPublisherTest, CreateFailsWithEmptyServiceInspectionTopic) {
  EXPECT_THAT(InspectionPublisher::Create("my_test_service", ""),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(InspectionPublisherTest, CreateAndPublish) {
  const std::string service_name = "my_test_service";
  const std::string service_inspection_topic =
      "/service_inspection/services/my_test_service";
  ASSERT_OK_AND_ASSIGN(
      auto publisher,
      InspectionPublisher::Create(service_name, service_inspection_topic));

  EXPECT_EQ(publisher.GetTopicName(), service_inspection_topic);

  // Create a subscriber to listen for the published message.
  PubSub pub_sub;
  absl::Notification msg_received;
  intrinsic_proto::services::v1::ServiceInspectionData received_message;
  ASSERT_OK_AND_ASSIGN(
      auto sub,
      pub_sub.CreateSubscription<
          intrinsic_proto::services::v1::ServiceInspectionData>(
          publisher.GetTopicName(), {},
          [&](const intrinsic_proto::services::v1::ServiceInspectionData&
                  message) {
            received_message = message;
            msg_received.Notify();
          }));

  // Prepare and publish a test message.
  intrinsic_proto::icon::v1::HardwareModuleInspectionData test_payload;
  test_payload.mutable_event_history()->add_events()->set_message("test event");
  const absl::Time publish_time = absl::Now();
  ASSERT_OK(publisher.Publish(test_payload, publish_time));

  // Wait for the message to be received.
  ASSERT_TRUE(msg_received.WaitForNotificationWithTimeout(absl::Seconds(5)));

  // Verify the received message.
  intrinsic_proto::icon::v1::HardwareModuleInspectionData received_payload;
  ASSERT_TRUE(received_message.data().UnpackTo(&received_payload));
  EXPECT_THAT(received_payload, EqualsProto(test_payload));

  const absl::Time received_timestamp =
      absl::FromUnixSeconds(received_message.timestamp().seconds()) +
      absl::Nanoseconds(received_message.timestamp().nanos());
  const absl::Time truncated_publish_time =
      absl::FromUnixNanos(absl::ToUnixNanos(publish_time));
  EXPECT_EQ(received_timestamp, truncated_publish_time);
}

}  // namespace

}  // namespace intrinsic::icon
