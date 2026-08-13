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

#include "intrinsic/util/thread/thread_pool.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>

#include "absl/status/status.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/stop_token.h"

namespace intrinsic {
namespace {

using ::absl_testing::StatusIs;

TEST(ThreadPoolTest, ScheduleWorks) {
  absl::Notification notification;
  ThreadPool pool(/*num_threads=*/1);
  EXPECT_OK(pool.Schedule([&]() { notification.Notify(); }));
  EXPECT_TRUE(notification.WaitForNotificationWithTimeout(absl::Seconds(1)));
}

TEST(ThreadPoolTest, ScheduleWorksWithBoundedQueue) {
  absl::Notification started;
  absl::Notification can_stop;
  ThreadPool pool(/*num_threads=*/1, /*max_queue_size=*/1);
  EXPECT_OK(pool.Schedule([&]() {
    started.Notify();
    // Block until we notify, so we can force the third Schedule() to fail.
    can_stop.WaitForNotification();
  }));
  started.WaitForNotification();
  EXPECT_EQ(pool.GetSize(), 0);

  // Now the queue size is 0, so the third Schedule() should fail.
  EXPECT_OK(pool.Schedule([&]() {}));
  EXPECT_EQ(pool.GetSize(), 1);
  EXPECT_THAT(pool.Schedule([&]() {}),
              StatusIs(absl::StatusCode::kResourceExhausted));
  EXPECT_EQ(pool.GetSize(), 1);
  can_stop.Notify();
}

TEST(ThreadPoolTest, ScheduledWorkIsCancelledOnDestruction) {
  auto pool = std::make_unique<ThreadPool>(/*num_threads=*/1);
  absl::Notification started_notification;
  absl::Notification notification;
  EXPECT_OK(pool->Schedule([&](StopToken stop_token) {
    started_notification.Notify();
    while (!stop_token.stop_requested()) {
      absl::SleepFor(absl::Milliseconds(100));
    }
    notification.Notify();
  }));
  constexpr absl::Duration timeout = absl::Seconds(1);
  // Ensure that the scheduled work has started before we destroy the pool.
  EXPECT_TRUE(started_notification.WaitForNotificationWithTimeout(timeout));
  pool.reset();

  // Check that the pool destruction has cancelled the worker, i.e. requested
  // a stop through the StopToken.
  EXPECT_TRUE(notification.HasBeenNotified());
}

TEST(ThreadPoolTest, GetSizeAndCapacityWork) {
  ThreadPool pool(/*num_threads=*/1, /*max_queue_size=*/3);
  EXPECT_EQ(pool.GetSize(), 0);

  absl::Notification started;
  absl::Notification can_stop;
  absl::Notification finished;
  EXPECT_OK(pool.Schedule([&]() {
    started.Notify();
    can_stop.WaitForNotification();
  }));
  started.WaitForNotification();  // Wait for dequeue.
  EXPECT_EQ(pool.GetSize(), 0);   // Was dequeued.
  EXPECT_OK(pool.Schedule([&]() {}));
  EXPECT_EQ(pool.GetSize(), 1);  // Pending can_stop invocable to complete.
  EXPECT_OK(pool.Schedule([&]() {}));
  EXPECT_EQ(pool.GetSize(), 2);
  EXPECT_OK(pool.Schedule([&]() { finished.Notify(); }));
  EXPECT_EQ(pool.GetSize(), 3);
  can_stop.Notify();
  finished.WaitForNotification();
  EXPECT_EQ(pool.GetSize(), 0);
}

}  // namespace

}  // namespace intrinsic
