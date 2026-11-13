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

#include <cstddef>
#include <cstdint>
#include <limits>
#include <numeric>
#include <span>
#include <thread>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "intrinsic/middleware/buffers/allocators/heap_buffer_allocator.h"
#include "intrinsic/middleware/buffers/buffer_factory.h"
#include "intrinsic/middleware/buffers/buffer_internal.h"
#include "intrinsic/middleware/buffers/immutable_buffer.h"
#include "intrinsic/middleware/buffers/mutable_buffer.h"
#include "internal/testing.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::ElementsAre;
using ::testing::ElementsAreArray;
using ::testing::HasSubstr;
using ::testing::MatchResultListener;

template <typename BufferT>
bool MatchBufferState(const BufferT& arg, bool expected_valid,
                      const uint8_t* expected_ptr, size_t expected_size,
                      size_t expected_alignment,
                      MatchResultListener* result_listener) {
  const uint8_t* const expected_end =
      expected_valid ? expected_ptr + expected_size : nullptr;
  const bool misaligned =
      expected_valid &&
      (expected_ptr == nullptr ||
       reinterpret_cast<uintptr_t>(arg.data()) % expected_alignment != 0);
  if (arg.valid() != expected_valid || arg.is_shm() ||
      arg.data() != expected_ptr || misaligned || arg.size() != expected_size ||
      arg.empty() != (expected_size == 0) ||
      arg.alignment() != expected_alignment || arg.begin() != expected_ptr ||
      arg.end() != expected_end || arg.cbegin() != expected_ptr ||
      arg.cend() != expected_end) {
    *result_listener << "where valid=" << arg.valid()
                     << ", is_shm=" << arg.is_shm()
                     << ", data=" << static_cast<const void*>(arg.data())
                     << " (expected " << static_cast<const void*>(expected_ptr)
                     << "), size=" << arg.size() << " (expected "
                     << expected_size << "), alignment=" << arg.alignment()
                     << " (expected " << expected_alignment << ")";
    return false;
  }
  return true;
}

MATCHER(IsEmptyAndInvalid, "is an empty, invalid buffer") {
  return MatchBufferState(arg, /*expected_valid=*/false,
                          /*expected_ptr=*/nullptr, /*expected_size=*/0,
                          /*expected_alignment=*/0, result_listener);
}

MATCHER_P3(IsValidView, expected_data, expected_size, expected_alignment,
           "is a valid buffer view with expected data, size, and alignment") {
  return MatchBufferState(
      arg, /*expected_valid=*/true, static_cast<const uint8_t*>(expected_data),
      static_cast<size_t>(expected_size),
      static_cast<size_t>(expected_alignment), result_listener);
}

MATCHER_P(IsValidAllocatedBuffer, expected_size,
          "is a valid root-aligned buffer of the expected size") {
  using BufferT = std::decay_t<decltype(arg)>;
  return MatchBufferState(arg, /*expected_valid=*/true, arg.data(),
                          static_cast<size_t>(expected_size),
                          BufferT::kAlignment, result_listener);
}

absl::StatusOr<MutableBuffer> AllocateIotaBuffer(size_t size,
                                                 uint8_t start_value = 0) {
  absl::StatusOr<MutableBuffer> buffer = AllocateMutableBuffer(size);
  if (buffer.ok()) {
    std::iota(buffer->begin(), buffer->end(), start_value);
  }
  return buffer;
}

TEST(MutableBufferTest, DefaultConstructedBufferIsInvalidAndEmpty) {
  MutableBuffer buffer;
  EXPECT_THAT(buffer, IsEmptyAndInvalid());
}

TEST(MutableBufferTest, AllocateHeapReturns256ByteAlignedWritableBuffer) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buffer,
                       AllocateMutableBuffer(/*size=*/37));
  ASSERT_THAT(buffer, IsValidAllocatedBuffer(/*expected_size=*/37));

  const HeapBufferAllocator custom_heap_allocator;
  EXPECT_THAT(AllocateMutableBuffer(/*size=*/16, custom_heap_allocator),
              IsOkAndHolds(IsValidAllocatedBuffer(/*expected_size=*/16)));

  std::iota(buffer.begin(), buffer.end(), uint8_t{10});
  std::span<const uint8_t> view(buffer);
  ASSERT_EQ(view.size(), 37);
  EXPECT_EQ(view.front(), 10);
  EXPECT_EQ(view.back(), 46);

  const MutableBuffer& const_buffer = buffer;
  EXPECT_THAT(const_buffer, IsValidView(buffer.data(), /*expected_size=*/37,
                                        MutableBuffer::kAlignment));
}

TEST(MutableBufferTest, AllocateZeroSizeHeapBufferRemainsValidAndNonNull) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer buffer, AllocateMutableBuffer(/*size=*/0));
  ASSERT_THAT(buffer, IsValidAllocatedBuffer(/*expected_size=*/0));

  ASSERT_OK_AND_ASSIGN(MutableBuffer second_zero,
                       AllocateMutableBuffer(/*size=*/0));
  EXPECT_THAT(second_zero, IsValidView(buffer.data(), /*expected_size=*/0,
                                       MutableBuffer::kAlignment));

  EXPECT_THAT(FreezeMutableBuffer(std::move(second_zero)),
              IsOkAndHolds(IsValidView(buffer.data(), /*expected_size=*/0,
                                       ImmutableBuffer::kAlignment)));
  EXPECT_THAT(second_zero, IsEmptyAndInvalid());
}

TEST(MutableBufferTest, MoveConstructionAndMoveAssignmentTransferOwnership) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer first,
                       AllocateIotaBuffer(/*size=*/4, /*start_value=*/1));
  uint8_t* const original_ptr = first.data();

  MutableBuffer moved(std::move(first));
  EXPECT_THAT(first, IsEmptyAndInvalid());
  EXPECT_THAT(moved, IsValidView(original_ptr, /*expected_size=*/4,
                                 MutableBuffer::kAlignment));
  EXPECT_THAT(moved, ElementsAre(1, 2, 3, 4));

  ASSERT_OK_AND_ASSIGN(MutableBuffer second,
                       AllocateIotaBuffer(/*size=*/2, /*start_value=*/8));

  second = std::move(moved);
  EXPECT_THAT(moved, IsEmptyAndInvalid());
  EXPECT_THAT(second, IsValidView(original_ptr, /*expected_size=*/4,
                                  MutableBuffer::kAlignment));
  EXPECT_THAT(second, ElementsAre(1, 2, 3, 4));

  MutableBuffer& second_ref = second;
  second = std::move(second_ref);
  EXPECT_THAT(second, IsValidView(original_ptr, /*expected_size=*/4,
                                  MutableBuffer::kAlignment));

  second = MutableBuffer();
  EXPECT_THAT(second, IsEmptyAndInvalid());
}

TEST(ImmutableBufferTest, DefaultConstructedBufferIsInvalidAndEmpty) {
  ImmutableBuffer buffer;
  EXPECT_THAT(buffer, IsEmptyAndInvalid());
}

TEST(ImmutableBufferTest, CopyToImmutableBufferAllocatesAlignedReadOnlyBuffer) {
  const std::vector<uint8_t> source = {5, 10, 15, 20, 25};
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer buffer,
                       CopyToImmutableBuffer(source.data(), source.size()));
  EXPECT_THAT(buffer, IsValidAllocatedBuffer(source.size()));
  EXPECT_NE(buffer.data(), source.data());
  EXPECT_THAT(buffer, ElementsAreArray(source));

  const HeapBufferAllocator custom_heap_allocator;
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer custom_buffer,
                       CopyToImmutableBuffer(source.data(), source.size(),
                                             custom_heap_allocator));
  EXPECT_THAT(custom_buffer, IsValidAllocatedBuffer(source.size()));
  EXPECT_THAT(custom_buffer, ElementsAreArray(source));
}

TEST(ImmutableBufferTest, CopyToImmutableBufferZeroBytesAndNullDataHandling) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer empty_buffer,
                       CopyToImmutableBuffer(/*data=*/nullptr, /*size=*/0));
  ASSERT_THAT(empty_buffer, IsValidAllocatedBuffer(/*expected_size=*/0));

  const uint8_t non_null_byte = 7;
  EXPECT_THAT(CopyToImmutableBuffer(&non_null_byte, /*size=*/0),
              IsOkAndHolds(IsValidView(empty_buffer.data(), /*expected_size=*/0,
                                       ImmutableBuffer::kAlignment)));
  EXPECT_THAT(empty_buffer.Slice(/*offset=*/0, /*length=*/0),
              IsValidView(empty_buffer.data(), /*expected_size=*/0,
                          ImmutableBuffer::kAlignment));
  EXPECT_THAT(CopyToImmutableBuffer(/*data=*/nullptr, /*size=*/8),
              StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("null")));
}

TEST(ImmutableBufferTest, FreezeTransfersOwnershipWithoutCopyingBytes) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/3, /*start_value=*/42));
  const uint8_t* const original_ptr = writable.data();

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer frozen,
                       FreezeMutableBuffer(std::move(writable)));
  EXPECT_THAT(writable, IsEmptyAndInvalid());
  EXPECT_THAT(frozen, IsValidView(original_ptr, /*expected_size=*/3,
                                  ImmutableBuffer::kAlignment));
  EXPECT_THAT(frozen, ElementsAre(42, 43, 44));

  EXPECT_THAT(
      FreezeMutableBuffer(std::move(writable)),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("invalid")));
}

TEST(ImmutableBufferTest, SlicePreservesBackingStorageAndComputesAlignment) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/512, /*start_value=*/0));
  const uint8_t* const root_ptr = writable.data();

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  EXPECT_THAT(root, IsValidView(root_ptr, /*expected_size=*/512,
                                ImmutableBuffer::kAlignment));
  EXPECT_THAT(root.Slice(/*offset=*/256, /*length=*/64),
              IsValidView(root_ptr + 256, /*expected_size=*/64,
                          ImmutableBuffer::kAlignment));

  ImmutableBuffer unaligned_slice = root.Slice(/*offset=*/24, /*length=*/4);
  EXPECT_THAT(unaligned_slice, IsValidView(root_ptr + 24, /*expected_size=*/4,
                                           /*expected_alignment=*/8));
  EXPECT_THAT(unaligned_slice, ElementsAre(24, 25, 26, 27));

  ImmutableBuffer chained_slice =
      unaligned_slice.Slice(/*offset=*/2, /*length=*/2);
  EXPECT_THAT(chained_slice, IsValidView(root_ptr + 26, /*expected_size=*/2,
                                         /*expected_alignment=*/2));
  EXPECT_THAT(chained_slice, ElementsAre(26, 27));

  // Assigning a slice back to the parent buffer variable preserves the backing
  // allocation while updating the view range.
  root = unaligned_slice;
  EXPECT_THAT(root, IsValidView(root_ptr + 24, /*expected_size=*/4,
                                /*expected_alignment=*/8));
  EXPECT_THAT(root, ElementsAre(24, 25, 26, 27));
}

TEST(ImmutableBufferTest, CopyAndMoveAssignmentManageReferenceCount) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/16, /*start_value=*/10));
  const uint8_t* const root_ptr = writable.data();

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer copy = root;
  EXPECT_THAT(copy, IsValidView(root_ptr, /*expected_size=*/16,
                                ImmutableBuffer::kAlignment));

  ImmutableBuffer surviving_slice = root.Slice(/*offset=*/4, /*length=*/4);

  ImmutableBuffer moved(std::move(copy));
  EXPECT_THAT(copy, IsEmptyAndInvalid());
  EXPECT_THAT(moved, IsValidView(root_ptr, /*expected_size=*/16,
                                 ImmutableBuffer::kAlignment));

  root = std::move(moved);
  EXPECT_THAT(moved, IsEmptyAndInvalid());
  EXPECT_THAT(root, IsValidView(root_ptr, /*expected_size=*/16,
                                ImmutableBuffer::kAlignment));

  const ImmutableBuffer& const_root_ref = root;
  root = const_root_ref;
  EXPECT_THAT(root, IsValidView(root_ptr, /*expected_size=*/16,
                                ImmutableBuffer::kAlignment));

  ImmutableBuffer& root_ref = root;
  root = std::move(root_ref);
  EXPECT_THAT(root, IsValidView(root_ptr, /*expected_size=*/16,
                                ImmutableBuffer::kAlignment));

  ImmutableBuffer invalid;
  EXPECT_THAT(ImmutableBuffer(invalid), IsEmptyAndInvalid());
  root = invalid;
  EXPECT_THAT(root, IsEmptyAndInvalid());
  root = surviving_slice;
  ASSERT_TRUE(root.valid());
  root = ImmutableBuffer();
  EXPECT_THAT(root, IsEmptyAndInvalid());
  EXPECT_THAT(surviving_slice, IsValidView(root_ptr + 4, /*expected_size=*/4,
                                           /*expected_alignment=*/4));
  EXPECT_THAT(surviving_slice, ElementsAre(14, 15, 16, 17));
}

TEST(ImmutableBufferTest, SliceRejectsOutOfBoundsAndOverflowArguments) {
  const std::vector<uint8_t> source = {1, 2, 3, 4};
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer buffer,
                       CopyToImmutableBuffer(source.data(), source.size()));
  EXPECT_THAT(buffer.Slice(/*offset=*/5, /*length=*/0), IsEmptyAndInvalid());
  EXPECT_THAT(buffer.Slice(/*offset=*/2, /*length=*/3), IsEmptyAndInvalid());
  EXPECT_THAT(buffer.Slice(/*offset=*/2, std::numeric_limits<size_t>::max()),
              IsEmptyAndInvalid());
  EXPECT_THAT(buffer.Slice(/*offset=*/4, /*length=*/0),
              IsValidView(buffer.data() + 4, /*expected_size=*/0,
                          /*expected_alignment=*/4));

  ImmutableBuffer invalid_buffer;
  EXPECT_THAT(invalid_buffer.Slice(/*offset=*/0, /*length=*/0),
              IsEmptyAndInvalid());
}

TEST(ImmutableBufferTest, ConcurrentCopySliceAndDestructionAreThreadSafe) {
  const std::vector<uint8_t> source = {10, 20, 30, 40, 50, 60, 70, 80};
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer shared,
                       CopyToImmutableBuffer(source.data(), source.size()));

  std::vector<std::thread> threads;
  threads.reserve(8);
  for (int t = 0; t < 8; ++t) {
    threads.emplace_back([shared]() {
      for (int i = 0; i < 500; ++i) {
        ImmutableBuffer copy = shared;
        ImmutableBuffer slice = copy.Slice(/*offset=*/2, /*length=*/4);
        EXPECT_THAT(slice, ElementsAre(30, 40, 50, 60));
      }
    });
  }
  shared = ImmutableBuffer();
  for (std::thread& worker : threads) {
    worker.join();
  }
}

class TrackingControlBlock final : public buffer_internal::BufferControlBlock {
 public:
  explicit TrackingControlBlock(int* destroyed_count) noexcept
      : destroyed_count_(destroyed_count) {}

  ~TrackingControlBlock() override { ++(*destroyed_count_); }

 private:
  int* destroyed_count_;
};

TEST(BufferControlBlockTest,
     PolymorphicSubclassRunsDerivedDestructorOnLastRelease) {
  static_assert(
      std::has_virtual_destructor_v<buffer_internal::BufferControlBlock>);
  static_assert(!std::is_constructible_v<buffer_internal::BufferControlBlock>);
  static_assert(
      !std::is_copy_constructible_v<buffer_internal::BufferControlBlock>);

  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment] = {};
  int destroyed_count = 0;

  {
    MutableBuffer unfrozen =
        buffer_internal::PrivateAccessor::CreateMutableBuffer(
            aligned_bytes, sizeof(aligned_bytes),
            new TrackingControlBlock(&destroyed_count), /*is_shm=*/false);
    EXPECT_EQ(destroyed_count, 0);
  }
  EXPECT_EQ(destroyed_count, 1);

  MutableBuffer writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count), /*is_shm=*/false);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer copy = root;
  ImmutableBuffer slice = root.Slice(/*offset=*/0, /*length=*/64);

  root = ImmutableBuffer();
  copy = ImmutableBuffer();
  EXPECT_EQ(destroyed_count, 1);

  slice = ImmutableBuffer();
  EXPECT_EQ(destroyed_count, 2);
}

}  // namespace
}  // namespace intrinsic
