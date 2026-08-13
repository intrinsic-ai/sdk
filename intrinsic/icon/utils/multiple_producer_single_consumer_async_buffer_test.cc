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

#include "intrinsic/icon/utils/multiple_producer_single_consumer_async_buffer.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <atomic>
#include <memory>
#include <utility>
#include <vector>

#include "absl/synchronization/notification.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic {
namespace {

TEST(MultipleProducerSingleConsumerAsyncBuffer, UpdateBufferWithConstRef) {
  MultipleProducerSingleConsumerAsyncBuffer<int> buffer;
  const int expected_value = 42;
  EXPECT_OK(buffer.UpdateBuffer(expected_value));

  int* value;
  EXPECT_TRUE(buffer.GetActiveBuffer(&value));
  EXPECT_EQ(*value, expected_value);
}

TEST(MultipleProducerSingleConsumerAsyncBuffer, UpdateBufferWithRValue) {
  MultipleProducerSingleConsumerAsyncBuffer<std::unique_ptr<int>> buffer;
  auto expected_value = std::make_unique<int>(42);
  auto* ptr = expected_value.get();
  EXPECT_OK(buffer.UpdateBuffer(std::move(expected_value)));

  std::unique_ptr<int>* value;
  EXPECT_TRUE(buffer.GetActiveBuffer(&value));
  EXPECT_EQ(value->get(), ptr);
  EXPECT_EQ(**value, 42);
}

TEST(MultipleProducerSingleConsumerAsyncBuffer, SingleProducerSingleConsumer) {
  MultipleProducerSingleConsumerAsyncBuffer<int> buffer;
  int expected_value = 42;
  EXPECT_OK(buffer.UpdateBuffer(
      [&](int& buffer_value) { buffer_value = expected_value; }));

  int* value;
  EXPECT_TRUE(buffer.GetActiveBuffer(&value));
  EXPECT_EQ(*value, expected_value);
}

TEST(MultipleProducerSingleConsumerAsyncBuffer,
     MultipleProducerSingleConsumer) {
  MultipleProducerSingleConsumerAsyncBuffer<int> buffer;

  int num_threads = 5;
  std::vector<std::unique_ptr<Thread>> threads;
  // For a pod value, the behavior is the same as with an atomic. Therefore,
  // we compare with the value stored in the atomic.
  std::atomic<int> latest_value = 0;
  absl::Notification start_threads;

  for (int i = 0; i < num_threads; ++i) {
    threads.push_back(std::make_unique<Thread>([&, i]() {
      start_threads.WaitForNotification();
      EXPECT_OK(buffer.UpdateBuffer([&, i](int& buffer_value) {
        // Update buffer holds the lock while this callback is executed.
        // Therefore, other threads cannot come between updating the control
        // value (the atomic) and the actual value in the buffer. This test
        // makes sure of this.
        latest_value = i;
        buffer_value = latest_value;
      }));
    }));
  }
  start_threads.Notify();

  threads.clear();
  int* value;
  EXPECT_TRUE(buffer.GetActiveBuffer(&value));
  EXPECT_EQ(*value, latest_value);
}

}  // namespace
}  // namespace intrinsic
