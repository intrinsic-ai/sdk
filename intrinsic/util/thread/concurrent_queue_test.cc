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

#include "intrinsic/util/thread/concurrent_queue.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "internal/testing.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;

TEST(ConcurrentQueueTest, EnqueueAndDequeueWithInfiniteTimeout) {
  ConcurrentQueue<int> queue(10);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_OK(queue.Enqueue(3));
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(1));
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(2));
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(3));
}

TEST(ConcurrentQueueTest, EnqueueAndDequeueWithZeroTimeout) {
  ConcurrentQueue<int> queue(10);
  EXPECT_OK(queue.Enqueue(1, absl::ZeroDuration()));
  EXPECT_OK(queue.Enqueue(2, absl::ZeroDuration()));
  EXPECT_OK(queue.Enqueue(3, absl::ZeroDuration()));
  EXPECT_THAT(queue.Dequeue(absl::ZeroDuration()), IsOkAndHolds(1));
  EXPECT_THAT(queue.Dequeue(absl::ZeroDuration()), IsOkAndHolds(2));
  EXPECT_THAT(queue.Dequeue(absl::ZeroDuration()), IsOkAndHolds(3));
}

TEST(ConcurrentQueueTest, DequeuingFailsOnTooShortTimeout) {
  ConcurrentQueue<int> queue(10);
  Thread producer_thread([&]() {
    absl::SleepFor(absl::Milliseconds(500));
    EXPECT_OK(queue.Enqueue(1));
  });
  EXPECT_THAT(queue.Dequeue(absl::Milliseconds(50)),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
}

TEST(ConcurrentQueueTest, DequeuingSucceedsOnGoodTimeout) {
  ConcurrentQueue<int> queue(10);
  Thread producer_thread([&]() {
    absl::SleepFor(absl::Milliseconds(500));
    EXPECT_OK(queue.Enqueue(1));
  });
  EXPECT_THAT(queue.Dequeue(absl::Milliseconds(550)), IsOkAndHolds(1));
}

TEST(ConcurrentQueueTest, EnqueuingFailsOnFullQueue) {
  ConcurrentQueue<int> queue(3);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_OK(queue.Enqueue(3));
  EXPECT_THAT(queue.Enqueue(4, absl::Milliseconds(50)),
              StatusIs(absl::StatusCode::kResourceExhausted));
}

TEST(ConcurrentQueueTest, EnqueuingSucceedsOnGoodTimeout) {
  ConcurrentQueue<int> queue(3);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_OK(queue.Enqueue(3));
  Thread consumer_thread([&]() {
    absl::SleepFor(absl::Milliseconds(50));
    EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(1));
  });
  EXPECT_OK(queue.Enqueue(4, absl::Milliseconds(100)));
}

TEST(ConcurrentQueueTest, EmptyQueueIsEmpty) {
  ConcurrentQueue<int> queue(3);
  EXPECT_TRUE(queue.IsEmpty());
}

TEST(ConcurrentQueueTest, EnqueuingFailsOnClosedQueue) {
  ConcurrentQueue<int> queue(3);
  queue.Close();
  EXPECT_THAT(queue.Enqueue(1), StatusIs(absl::StatusCode::kUnavailable));
}

TEST(ConcurrentQueueTest, EnqueuingFailsOnClosedQueueWhichContainsData) {
  ConcurrentQueue<int> queue(3);
  EXPECT_OK(queue.Enqueue(1));
  queue.Close();
  EXPECT_THAT(queue.Enqueue(1), StatusIs(absl::StatusCode::kUnavailable));
}

TEST(ConcurrentQueueTest, DequeuingWorksAndFailsOnClosedEmptyQueue) {
  ConcurrentQueue<int> queue(3);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_OK(queue.Enqueue(3));
  queue.Close();
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(1));
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(2));
  EXPECT_THAT(queue.Dequeue(), IsOkAndHolds(3));
  EXPECT_THAT(queue.Dequeue(), StatusIs(absl::StatusCode::kUnavailable));
}

TEST(ConcurrentQueueTest, ElementsCanBeEnqueuedOnUnboundedQueue) {
  ConcurrentQueue<int> queue(0);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_OK(queue.Enqueue(3));
  EXPECT_OK(queue.Enqueue(4));
}

TEST(ConcurrentQueueTest, GetSizeReturnsCorrectSizeBounded) {
  ConcurrentQueue<int> queue(3);
  EXPECT_EQ(queue.GetSize(), 0);
  EXPECT_OK(queue.Enqueue(1));
  EXPECT_EQ(queue.GetSize(), 1);
  EXPECT_OK(queue.Enqueue(2));
  EXPECT_EQ(queue.GetSize(), 2);
  EXPECT_OK(queue.Enqueue(3));
  EXPECT_EQ(queue.GetSize(), 3);
}

TEST(ConcurrentQueueTest, GetSizeReturnsCorrectSizeUnbounded) {
  ConcurrentQueue<int> queue(0);
  EXPECT_EQ(queue.GetSize(), 0);
  for (int i = 0; i < 10; ++i) {
    EXPECT_OK(queue.Enqueue(i));
    EXPECT_EQ(queue.GetSize(), i + 1);
  }
}

TEST(ConcurrentQueueTest, ClearResultsInEmptyQueue) {
  ConcurrentQueue<int> queue(0);
  constexpr int kItemCount = 10;
  for (int i = 0; i < kItemCount; ++i) {
    EXPECT_OK(queue.Enqueue(i));
  }
  EXPECT_EQ(queue.GetSize(), kItemCount);
  queue.Clear();
  EXPECT_EQ(queue.GetSize(), 0);
  queue.Clear();  // Confirm that it works on an empty queue too.
  EXPECT_EQ(queue.GetSize(), 0);
}

}  // namespace
}  // namespace intrinsic
