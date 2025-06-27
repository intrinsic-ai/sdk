// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/common/buffers/rt_queue_buffer.h"

#include <gtest/gtest.h>

#include <cstddef>

namespace intrinsic {
namespace internal {
namespace {

TEST(RtQueueBufferTest, ConstructDestruct) { RtQueueBuffer<int> queue(10); }

TEST(RtQueueBufferTest, CapacityIsCorrect) {
  size_t constexpr kCapacity = 10;
  RtQueueBuffer<int> queue(kCapacity);
  EXPECT_EQ(queue.Capacity(), kCapacity);
}

TEST(RtQueueBufferTest, EmptyReturnsTrueForEmptyQeuue) {
  RtQueueBuffer<int> queue(10);
  EXPECT_TRUE(queue.Empty());
  EXPECT_EQ(queue.Size(), 0);
}

TEST(RtQueueBufferTest, PrepareInsertReturnsNullptrWhenFull) {
  RtQueueBuffer<int> queue(1);
  EXPECT_NE(queue.PrepareInsert(), nullptr);
  queue.FinishInsert();
  EXPECT_EQ(queue.PrepareInsert(), nullptr);
}

TEST(RtQueueBufferTest, SizeIsCorrectAfterInsertAndRemove) {
  RtQueueBuffer<int> queue(10);
  for (int i = 0; i < 10; ++i) {
    (void)queue.PrepareInsert();
    queue.FinishInsert();
    EXPECT_EQ(queue.Size(), i + 1);
  }

  for (int i = 10; i > 0; --i) {
    (void)queue.Front();
    queue.DropFront();
    EXPECT_EQ(queue.Size(), i - 1);
  }
}

TEST(RtQueueBufferTest, FullReportsFullWhenCapacityReached) {
  RtQueueBuffer<int> queue(2);
  EXPECT_NE(queue.PrepareInsert(), nullptr);
  queue.FinishInsert();
  EXPECT_NE(queue.PrepareInsert(), nullptr);
  queue.FinishInsert();
  EXPECT_TRUE(queue.Full());
}

TEST(RtQueueBufferTest, KeepFrontMaintainsFrontOfQueue) {
  RtQueueBuffer<int> queue(2);
  constexpr int kExpectedResult = 2;
  {
    int* item = queue.PrepareInsert();
    *item = kExpectedResult;
    queue.FinishInsert();
  }

  {
    int* item = queue.PrepareInsert();
    *item = kExpectedResult + 1;  // some different value
    queue.FinishInsert();
  }

  int* result = queue.Front();
  EXPECT_EQ(*result, kExpectedResult);
  queue.KeepFront();
  result = queue.Front();
  EXPECT_EQ(*result, kExpectedResult);
}

TEST(RtQueueBufferTest, DropFrontMovesFrontToNextValue) {
  RtQueueBuffer<int> queue(2);
  constexpr int kExpectedResult1 = 1;
  constexpr int kExpectedResult2 = 2;
  {
    int* item = queue.PrepareInsert();
    *item = kExpectedResult1;
    queue.FinishInsert();
  }

  {
    int* item = queue.PrepareInsert();
    *item = kExpectedResult2;
    queue.FinishInsert();
  }

  int* result = queue.Front();
  EXPECT_EQ(*result, kExpectedResult1);
  queue.DropFront();
  result = queue.Front();
  EXPECT_EQ(*result, kExpectedResult2);
}

TEST(RtQueueBufferTest, InitElementsInitializesPreparedElements) {
  RtQueueBuffer<int> queue(10);
  int n = 0;
  queue.InitElements([&n](int* item) { *item = n++; });
  for (int count = 0; count < queue.Capacity(); ++count) {
    int* item = queue.PrepareInsert();
    EXPECT_EQ(*item, count);
    queue.FinishInsert();
  }
}

TEST(RtQueueBufferTest, ConstructWithInitInitializesPreparedElements) {
  int n = 0;
  RtQueueBuffer<int> queue(10, [&n](int* item) { *item = n++; });
  for (int count = 0; count < queue.Capacity(); ++count) {
    int* item = queue.PrepareInsert();
    EXPECT_EQ(*item, count);
    queue.FinishInsert();
  }
}

}  // namespace
}  // namespace internal
}  // namespace intrinsic
