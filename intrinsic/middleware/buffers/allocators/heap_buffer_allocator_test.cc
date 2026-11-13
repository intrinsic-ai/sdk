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

#include "intrinsic/middleware/buffers/allocators/heap_buffer_allocator.h"

#include <cstddef>
#include <cstdint>
#include <numeric>
#include <type_traits>

#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "intrinsic/middleware/buffers/allocators/buffer_allocator.h"
#include "intrinsic/middleware/buffers/mutable_buffer.h"
#include "internal/testing.h"

namespace intrinsic {
namespace {

using ::testing::ElementsAre;

TEST(HeapBufferAllocatorTest, EnforcesClosedBaseClassAndSingletonContract) {
  static_assert(!std::is_copy_constructible_v<BufferAllocator>);
  static_assert(!std::is_move_constructible_v<BufferAllocator>);
  static_assert(!std::is_destructible_v<BufferAllocator>);
  static_assert(std::is_base_of_v<BufferAllocator, HeapBufferAllocator>);
  static_assert(std::is_trivially_destructible_v<HeapBufferAllocator>);

  EXPECT_EQ(&HeapBufferAllocator::Instance(), &HeapBufferAllocator::Instance());
}

TEST(HeapBufferAllocatorTest, AllocateReturns256ByteAlignedWritableBuffer) {
  const BufferAllocator& allocator = HeapBufferAllocator::Instance();
  ASSERT_OK_AND_ASSIGN(MutableBuffer buffer, allocator.Allocate(/*size=*/4));
  ASSERT_TRUE(buffer.valid());
  EXPECT_FALSE(buffer.is_shm());
  EXPECT_EQ(buffer.size(), 4);
  EXPECT_EQ(buffer.alignment(), MutableBuffer::kAlignment);
  ASSERT_NE(buffer.data(), nullptr);
  EXPECT_EQ(
      reinterpret_cast<uintptr_t>(buffer.data()) % MutableBuffer::kAlignment,
      0);

  std::iota(buffer.begin(), buffer.end(), uint8_t{10});
  EXPECT_THAT(buffer, ElementsAre(10, 11, 12, 13));
}

TEST(HeapBufferAllocatorTest, AllocateZeroSizeSharesAlignedSentinel) {
  const HeapBufferAllocator custom_allocator;
  ASSERT_OK_AND_ASSIGN(MutableBuffer first,
                       custom_allocator.Allocate(/*size=*/0));
  ASSERT_OK_AND_ASSIGN(MutableBuffer second,
                       HeapBufferAllocator::Instance().Allocate(/*size=*/0));

  EXPECT_TRUE(first.valid());
  EXPECT_FALSE(first.is_shm());
  EXPECT_TRUE(first.empty());
  EXPECT_EQ(first.alignment(), MutableBuffer::kAlignment);
  ASSERT_NE(first.data(), nullptr);
  EXPECT_EQ(
      reinterpret_cast<uintptr_t>(first.data()) % MutableBuffer::kAlignment, 0);
  EXPECT_EQ(second.data(), first.data());
}

}  // namespace
}  // namespace intrinsic
