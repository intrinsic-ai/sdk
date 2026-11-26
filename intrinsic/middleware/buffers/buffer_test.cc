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
#include <iterator>
#include <limits>
#include <numeric>
#include <optional>
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
#include "intrinsic/middleware/buffers/allocators/buffer_allocator.h"
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

template <typename AllocatorT>
concept CanUnfreezeOrCopyWith =
    requires(ImmutableBuffer&& buffer, AllocatorT&& allocator) {
      UnfreezeOrCopy(std::move(buffer), std::forward<AllocatorT>(allocator));
    };

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

TEST(ImmutableBufferTest, UnfreezeOrCopyReclaimsUniqueAlignedBufferInPlace) {
  static_assert(CanUnfreezeOrCopyWith<HeapBufferAllocator&>);
  static_assert(CanUnfreezeOrCopyWith<const HeapBufferAllocator&>);
  static_assert(CanUnfreezeOrCopyWith<const BufferAllocator&>);
  static_assert(!CanUnfreezeOrCopyWith<HeapBufferAllocator>);
  static_assert(!CanUnfreezeOrCopyWith<const HeapBufferAllocator&&>);
  static_assert(!CanUnfreezeOrCopyWith<const BufferAllocator&&>);

  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/64, /*start_value=*/10));
  const uint8_t* const original_ptr = writable.data();

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer frozen,
                       FreezeMutableBuffer(std::move(writable)));
  ASSERT_OK_AND_ASSIGN(MutableBuffer unfrozen,
                       UnfreezeOrCopy(std::move(frozen)));
  EXPECT_THAT(frozen, IsEmptyAndInvalid());
  ASSERT_THAT(unfrozen, IsValidView(original_ptr, /*expected_size=*/64,
                                    MutableBuffer::kAlignment));
  EXPECT_EQ(unfrozen.data()[0], 10);
  EXPECT_EQ(unfrozen.data()[63], 73);

  unfrozen.data()[0] = 200;
  EXPECT_EQ(unfrozen.data()[0], 200);
}

TEST(ImmutableBufferTest,
     UnfreezeOrCopyCopiesSharedBufferWithoutMutatingRemainingCopy) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/16, /*start_value=*/10));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer shared = root;

  ASSERT_OK_AND_ASSIGN(MutableBuffer copy, UnfreezeOrCopy(std::move(root)));
  EXPECT_THAT(root, IsEmptyAndInvalid());
  ASSERT_THAT(copy, IsValidAllocatedBuffer(/*expected_size=*/16));
  EXPECT_NE(copy.data(), shared.data());
  EXPECT_THAT(copy, ElementsAreArray(shared));

  copy.data()[0] = 255;
  EXPECT_EQ(shared.data()[0], 10);

  const HeapBufferAllocator custom_heap_allocator;
  ImmutableBuffer shared_copy = shared;
  ASSERT_OK_AND_ASSIGN(
      MutableBuffer custom_copy,
      UnfreezeOrCopy(std::move(shared_copy), custom_heap_allocator));
  EXPECT_THAT(shared_copy, IsEmptyAndInvalid());
  EXPECT_THAT(custom_copy, IsValidAllocatedBuffer(/*expected_size=*/16));
  EXPECT_NE(custom_copy.data(), shared.data());
  EXPECT_THAT(custom_copy, ElementsAreArray(shared));
}

TEST(ImmutableBufferTest, UnfreezeOrCopyCopiesUnalignedSliceIntoAlignedBuffer) {
  ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                       AllocateIotaBuffer(/*size=*/512, /*start_value=*/0));
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer unaligned_slice = root.Slice(/*offset=*/16, /*length=*/32);
  root = ImmutableBuffer();

  const uint8_t* const slice_ptr = unaligned_slice.data();
  EXPECT_EQ(unaligned_slice.alignment(), 16);

  ASSERT_OK_AND_ASSIGN(MutableBuffer copy,
                       UnfreezeOrCopy(std::move(unaligned_slice)));
  EXPECT_THAT(unaligned_slice, IsEmptyAndInvalid());
  EXPECT_THAT(copy, IsValidAllocatedBuffer(/*expected_size=*/32));
  EXPECT_NE(copy.data(), slice_ptr);
  std::vector<uint8_t> expected(32);
  std::iota(expected.begin(), expected.end(), uint8_t{16});
  EXPECT_THAT(copy, ElementsAreArray(expected));
}

TEST(ImmutableBufferTest, UnfreezeOrCopyHandlesEmptyValidAndInvalidBuffers) {
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer empty_unique,
                       CopyToImmutableBuffer(/*data=*/nullptr, /*size=*/0));
  const uint8_t* const empty_ptr = empty_unique.data();

  ASSERT_OK_AND_ASSIGN(MutableBuffer unfrozen_empty,
                       UnfreezeOrCopy(std::move(empty_unique)));
  EXPECT_THAT(empty_unique, IsEmptyAndInvalid());
  EXPECT_THAT(unfrozen_empty, IsValidView(empty_ptr, /*expected_size=*/0,
                                          MutableBuffer::kAlignment));

  ASSERT_OK_AND_ASSIGN(ImmutableBuffer empty_shared,
                       CopyToImmutableBuffer(/*data=*/nullptr, /*size=*/0));
  ImmutableBuffer empty_copy = empty_shared;
  ASSERT_OK_AND_ASSIGN(MutableBuffer copied_empty,
                       UnfreezeOrCopy(std::move(empty_copy)));
  EXPECT_THAT(empty_copy, IsEmptyAndInvalid());
  EXPECT_THAT(copied_empty, IsValidAllocatedBuffer(/*expected_size=*/0));
  EXPECT_THAT(empty_shared, IsValidAllocatedBuffer(/*expected_size=*/0));

  ImmutableBuffer invalid;
  EXPECT_THAT(
      UnfreezeOrCopy(std::move(invalid)),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("invalid")));
  EXPECT_THAT(invalid, IsEmptyAndInvalid());
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
  explicit TrackingControlBlock(int* destroyed_count,
                                bool writable = true) noexcept
      : destroyed_count_(destroyed_count), writable_(writable) {}

  ~TrackingControlBlock() override { ++(*destroyed_count_); }

  [[nodiscard]] bool IsWritable() const noexcept override { return writable_; }

 private:
  int* destroyed_count_;
  bool writable_;
};

TEST(BufferControlBlockTest,
     PolymorphicSubclassRunsDerivedDestructorOnLastRelease) {
  static_assert(
      std::has_virtual_destructor_v<buffer_internal::BufferControlBlock>);
  static_assert(std::is_abstract_v<buffer_internal::BufferControlBlock>);
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

TEST(BufferControlBlockTest,
     TryUnfreezeRejectsInvalidBufferWithoutModification) {
  ImmutableBuffer invalid;
  EXPECT_FALSE(
      buffer_internal::PrivateAccessor::TryUnfreeze(invalid).has_value());
  EXPECT_THAT(invalid, IsEmptyAndInvalid());
}

TEST(BufferControlBlockTest,
     TryUnfreezeRejectsSharedBufferWithoutModification) {
  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment * 2] = {};
  int destroyed_count = 0;

  MutableBuffer writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count), /*is_shm=*/false);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer shared = root;

  EXPECT_FALSE(buffer_internal::PrivateAccessor::TryUnfreeze(root).has_value());
  EXPECT_THAT(root, IsValidView(aligned_bytes, sizeof(aligned_bytes),
                                ImmutableBuffer::kAlignment));
  EXPECT_EQ(destroyed_count, 0);

  root = ImmutableBuffer();
  shared = ImmutableBuffer();
  EXPECT_EQ(destroyed_count, 1);
}

TEST(BufferControlBlockTest,
     TryUnfreezeRejectsUnalignedSliceWithoutModification) {
  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment * 2] = {};
  int destroyed_count = 0;

  MutableBuffer writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count), /*is_shm=*/false);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer unaligned_slice = root.Slice(/*offset=*/16, /*length=*/32);
  root = ImmutableBuffer();

  EXPECT_FALSE(buffer_internal::PrivateAccessor::TryUnfreeze(unaligned_slice)
                   .has_value());
  EXPECT_THAT(unaligned_slice,
              IsValidView(aligned_bytes + 16, /*expected_size=*/32,
                          /*expected_alignment=*/16));
  EXPECT_EQ(destroyed_count, 0);

  unaligned_slice = ImmutableBuffer();
  EXPECT_EQ(destroyed_count, 1);
}

TEST(BufferControlBlockTest,
     TryUnfreezeReclaimsAlignedSliceInPlaceAndTransfersControlBlock) {
  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment * 2] = {};
  int destroyed_count = 0;

  MutableBuffer writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count), /*is_shm=*/false);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer root,
                       FreezeMutableBuffer(std::move(writable)));
  ImmutableBuffer aligned_slice =
      root.Slice(/*offset=*/MutableBuffer::kAlignment, /*length=*/64);
  root = ImmutableBuffer();

  {
    std::optional<MutableBuffer> unfrozen =
        buffer_internal::PrivateAccessor::TryUnfreeze(aligned_slice);
    ASSERT_TRUE(unfrozen.has_value());
    EXPECT_THAT(aligned_slice, IsEmptyAndInvalid());
    EXPECT_THAT(*unfrozen,
                IsValidView(aligned_bytes + MutableBuffer::kAlignment,
                            /*expected_size=*/64, MutableBuffer::kAlignment));
    EXPECT_EQ(destroyed_count, 0);
  }
  EXPECT_EQ(destroyed_count, 1);
}

TEST(BufferControlBlockTest,
     TryUnfreezeRejectsShmBufferAndUnfreezeOrCopyCopies) {
  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment * 2] = {};
  std::iota(std::begin(aligned_bytes), std::end(aligned_bytes), uint8_t{0});
  int destroyed_count = 0;

  MutableBuffer shm_writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count), /*is_shm=*/true);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer shm_buffer,
                       FreezeMutableBuffer(std::move(shm_writable)));

  EXPECT_FALSE(
      buffer_internal::PrivateAccessor::TryUnfreeze(shm_buffer).has_value());
  EXPECT_TRUE(shm_buffer.valid());
  EXPECT_TRUE(shm_buffer.is_shm());
  EXPECT_EQ(shm_buffer.data(), aligned_bytes);
  EXPECT_EQ(shm_buffer.size(), sizeof(aligned_bytes));
  EXPECT_EQ(destroyed_count, 0);

  ASSERT_OK_AND_ASSIGN(MutableBuffer shm_copy,
                       UnfreezeOrCopy(std::move(shm_buffer)));
  EXPECT_THAT(shm_buffer, IsEmptyAndInvalid());
  EXPECT_EQ(destroyed_count, 1);
  EXPECT_THAT(shm_copy, IsValidAllocatedBuffer(sizeof(aligned_bytes)));
  EXPECT_NE(shm_copy.data(), aligned_bytes);
  EXPECT_THAT(shm_copy, ElementsAreArray(aligned_bytes));
}

TEST(BufferControlBlockTest,
     TryUnfreezeRejectsReadOnlyControlBlockAndUnfreezeOrCopyCopies) {
  alignas(MutableBuffer::kAlignment)
      uint8_t aligned_bytes[MutableBuffer::kAlignment] = {};
  std::iota(std::begin(aligned_bytes), std::end(aligned_bytes), uint8_t{1});
  int destroyed_count = 0;

  MutableBuffer writable =
      buffer_internal::PrivateAccessor::CreateMutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          new TrackingControlBlock(&destroyed_count, /*writable=*/false),
          /*is_shm=*/false);
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer read_only_buffer,
                       FreezeMutableBuffer(std::move(writable)));

  // Even when uniquely owned, non-SHM, and 256-byte-aligned, a control block
  // with IsWritable() == false must not be unfrozen in place and remains
  // unmodified when TryUnfreeze returns std::nullopt.
  EXPECT_FALSE(buffer_internal::PrivateAccessor::TryUnfreeze(read_only_buffer)
                   .has_value());
  EXPECT_THAT(read_only_buffer,
              IsValidView(aligned_bytes, sizeof(aligned_bytes),
                          ImmutableBuffer::kAlignment));
  EXPECT_EQ(destroyed_count, 0);

  ASSERT_OK_AND_ASSIGN(MutableBuffer copy,
                       UnfreezeOrCopy(std::move(read_only_buffer)));
  EXPECT_THAT(read_only_buffer, IsEmptyAndInvalid());
  EXPECT_EQ(destroyed_count, 1);
  ASSERT_THAT(copy, IsValidAllocatedBuffer(sizeof(aligned_bytes)));
  EXPECT_NE(copy.data(), aligned_bytes);
  EXPECT_THAT(copy, ElementsAreArray(aligned_bytes));

  copy.data()[0] = 0xFF;
  EXPECT_EQ(aligned_bytes[0], 1);
}

TEST(BufferControlBlockTest,
     TryUnfreezeSynchronizesWithConcurrentReaderRelease) {
  constexpr size_t kSize = 256;
  constexpr int kIterations = 64;
  for (int iter = 0; iter < kIterations; ++iter) {
    const uint8_t start = static_cast<uint8_t>(iter);
    ASSERT_OK_AND_ASSIGN(MutableBuffer writable,
                         AllocateIotaBuffer(kSize, start));
    const uint8_t* const original_ptr = writable.data();
    const uint32_t expected_sum =
        std::accumulate(writable.begin(), writable.end(), uint32_t{0});

    ASSERT_OK_AND_ASSIGN(ImmutableBuffer owner,
                         FreezeMutableBuffer(std::move(writable)));
    ImmutableBuffer reader = owner;
    uint32_t observed_sum = 0;

    std::thread worker([reader = std::move(reader), &observed_sum]() mutable {
      observed_sum = std::accumulate(reader.begin(), reader.end(), uint32_t{0});
      reader = ImmutableBuffer();
    });

    std::optional<MutableBuffer> unfrozen;
    while (!unfrozen.has_value()) {
      unfrozen = buffer_internal::PrivateAccessor::TryUnfreeze(owner);
      if (!unfrozen.has_value()) {
        std::this_thread::yield();
      }
    }

    // Mutate in place before joining `worker` so TSan verifies the
    // `Release()` (`memory_order_acq_rel`) -> `IsUnique()`
    // (`memory_order_acquire`) happens-before edge directly.
    if (unfrozen->valid() && unfrozen->size() == kSize) {
      std::iota(unfrozen->begin(), unfrozen->end(),
                static_cast<uint8_t>(~start));
    }

    worker.join();
    EXPECT_EQ(observed_sum, expected_sum);
    EXPECT_THAT(owner, IsEmptyAndInvalid());
    EXPECT_THAT(*unfrozen,
                IsValidView(original_ptr, kSize, MutableBuffer::kAlignment));
  }
}

}  // namespace
}  // namespace intrinsic
