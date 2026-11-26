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

#include <atomic>
#include <cstddef>
#include <cstdint>
#include <iterator>
#include <limits>
#include <memory>
#include <numeric>
#include <optional>
#include <span>
#include <thread>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/functional/any_invocable.h"
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

TEST(ImmutableBufferTest,
     AdoptImmutableBufferPreservesPointerAlignmentAndRunsDeleterOnLastDrop) {
  alignas(ImmutableBuffer::kAlignment)
      uint8_t aligned_bytes[ImmutableBuffer::kAlignment * 2] = {};
  std::iota(std::begin(aligned_bytes), std::end(aligned_bytes), uint8_t{1});

  // 256-byte-aligned external buffer: preserves pointer identity and alignment,
  // and runs deleter only after all copies and slices are destroyed.
  int aligned_deleter_calls = 0;
  ASSERT_OK_AND_ASSIGN(
      ImmutableBuffer adopted_aligned,
      AdoptImmutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          [&aligned_deleter_calls]() { ++aligned_deleter_calls; }));
  ASSERT_THAT(adopted_aligned, IsValidView(aligned_bytes, sizeof(aligned_bytes),
                                           ImmutableBuffer::kAlignment));
  EXPECT_THAT(adopted_aligned, ElementsAreArray(aligned_bytes));

  ImmutableBuffer copy = adopted_aligned;
  ImmutableBuffer slice = adopted_aligned.Slice(/*offset=*/16, /*length=*/32);
  adopted_aligned = ImmutableBuffer();
  copy = ImmutableBuffer();
  EXPECT_EQ(aligned_deleter_calls, 0);
  ASSERT_THAT(slice, IsValidView(aligned_bytes + 16, /*expected_size=*/32,
                                 /*expected_alignment=*/16));
  EXPECT_EQ(slice.data()[0], aligned_bytes[16]);

  slice = ImmutableBuffer();
  EXPECT_EQ(aligned_deleter_calls, 1);

  // 16-byte-aligned external buffer with a move-only rvalue-qualified deleter.
  int unaligned_deleter_calls = 0;
  struct MoveOnlyDeleter {
    std::unique_ptr<int> token;
    int* calls = nullptr;
    void operator()() && {
      EXPECT_NE(token, nullptr);
      ++(*calls);
    }
  };
  absl::AnyInvocable<void() &&> move_only_deleter = MoveOnlyDeleter{
      .token = std::make_unique<int>(42), .calls = &unaligned_deleter_calls};
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer adopted_unaligned,
                       AdoptImmutableBuffer(aligned_bytes + 16, /*size=*/32,
                                            std::move(move_only_deleter)));
  EXPECT_THAT(adopted_unaligned,
              IsValidView(aligned_bytes + 16, /*expected_size=*/32,
                          /*expected_alignment=*/16));
  ImmutableBuffer unaligned_copy = adopted_unaligned;
  adopted_unaligned = ImmutableBuffer();
  EXPECT_EQ(unaligned_deleter_calls, 0);
  unaligned_copy = ImmutableBuffer();
  EXPECT_EQ(unaligned_deleter_calls, 1);

  // Odd alignment reports 1; slicing at offset 255 restores 256-byte alignment,
  // while slicing at offset 256 remains 1-byte aligned.
  ASSERT_OK_AND_ASSIGN(
      ImmutableBuffer odd_aligned,
      AdoptImmutableBuffer(aligned_bytes + 1,
                           /*size=*/ImmutableBuffer::kAlignment + 1,
                           /*deleter=*/nullptr));
  EXPECT_THAT(odd_aligned,
              IsValidView(aligned_bytes + 1,
                          /*expected_size=*/ImmutableBuffer::kAlignment + 1,
                          /*expected_alignment=*/1));
  EXPECT_THAT(odd_aligned.Slice(/*offset=*/255, /*length=*/2),
              IsValidView(aligned_bytes + ImmutableBuffer::kAlignment,
                          /*expected_size=*/2, ImmutableBuffer::kAlignment));
  EXPECT_THAT(odd_aligned.Slice(/*offset=*/256, /*length=*/1),
              IsValidView(aligned_bytes + ImmutableBuffer::kAlignment + 1,
                          /*expected_size=*/1, /*expected_alignment=*/1));

  // Over-aligned (> 256B) external pointers cap `alignment()` at `kAlignment`.
  alignas(512) uint8_t over_aligned_bytes[16] = {};
  EXPECT_THAT(
      AdoptImmutableBuffer(over_aligned_bytes, sizeof(over_aligned_bytes),
                           /*deleter=*/nullptr),
      IsOkAndHolds(IsValidView(over_aligned_bytes, sizeof(over_aligned_bytes),
                               ImmutableBuffer::kAlignment)));

  // Null deleter is allowed for externally pinned or static storage.
  EXPECT_THAT(
      AdoptImmutableBuffer(aligned_bytes, /*size=*/8, /*deleter=*/nullptr),
      IsOkAndHolds(IsValidView(aligned_bytes, /*expected_size=*/8,
                               ImmutableBuffer::kAlignment)));
}

TEST(ImmutableBufferTest,
     AdoptImmutableBufferUnfreezeOrCopyCopiesForBothUnalignedAnd256BAligned) {
  alignas(ImmutableBuffer::kAlignment)
      uint8_t aligned_bytes[ImmutableBuffer::kAlignment * 2] = {};
  std::iota(std::begin(aligned_bytes), std::end(aligned_bytes), uint8_t{10});

  // 1) < 256B-aligned adopted buffer must copy into a 256B-aligned buffer and
  // invoke its deleter.
  int unaligned_deleter_calls = 0;
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer unaligned_adopted,
                       AdoptImmutableBuffer(aligned_bytes + 16, /*size=*/32,
                                            [&unaligned_deleter_calls]() {
                                              ++unaligned_deleter_calls;
                                            }));
  EXPECT_EQ(unaligned_adopted.alignment(), 16);
  EXPECT_FALSE(buffer_internal::PrivateAccessor::TryUnfreeze(unaligned_adopted)
                   .has_value());
  EXPECT_EQ(unaligned_deleter_calls, 0);

  ASSERT_OK_AND_ASSIGN(MutableBuffer unaligned_copy,
                       UnfreezeOrCopy(std::move(unaligned_adopted)));
  EXPECT_THAT(unaligned_adopted, IsEmptyAndInvalid());
  EXPECT_EQ(unaligned_deleter_calls, 1);
  ASSERT_THAT(unaligned_copy, IsValidAllocatedBuffer(/*expected_size=*/32));
  EXPECT_NE(unaligned_copy.data(), aligned_bytes + 16);
  EXPECT_THAT(unaligned_copy,
              ElementsAreArray(aligned_bytes + 16, /*count=*/32));
  unaligned_copy.data()[0] = 0xFF;
  EXPECT_EQ(aligned_bytes[16], 26);

  // 2) Even a uniquely-owned 256B-aligned adopted buffer is read-only
  // (`IsWritable() == false`) and must copy rather than const_cast caller data.
  int aligned_deleter_calls = 0;
  ASSERT_OK_AND_ASSIGN(
      ImmutableBuffer aligned_adopted,
      AdoptImmutableBuffer(
          aligned_bytes, sizeof(aligned_bytes),
          [&aligned_deleter_calls]() { ++aligned_deleter_calls; }));
  EXPECT_EQ(aligned_adopted.alignment(), ImmutableBuffer::kAlignment);
  EXPECT_FALSE(buffer_internal::PrivateAccessor::TryUnfreeze(aligned_adopted)
                   .has_value());
  EXPECT_THAT(aligned_adopted, IsValidView(aligned_bytes, sizeof(aligned_bytes),
                                           ImmutableBuffer::kAlignment));
  EXPECT_EQ(aligned_deleter_calls, 0);

  ASSERT_OK_AND_ASSIGN(MutableBuffer aligned_copy,
                       UnfreezeOrCopy(std::move(aligned_adopted)));
  EXPECT_THAT(aligned_adopted, IsEmptyAndInvalid());
  EXPECT_EQ(aligned_deleter_calls, 1);
  ASSERT_THAT(aligned_copy, IsValidAllocatedBuffer(sizeof(aligned_bytes)));
  EXPECT_NE(aligned_copy.data(), aligned_bytes);
  EXPECT_THAT(aligned_copy, ElementsAreArray(aligned_bytes));
  aligned_copy.data()[0] = 0xEE;
  EXPECT_EQ(aligned_bytes[0], 10);
}

TEST(ImmutableBufferTest,
     AdoptImmutableBufferHandlesNullDataErrorsAndEmptyAdoption) {
  int error_deleter_calls = 0;
  EXPECT_THAT(AdoptImmutableBuffer(
                  /*data=*/nullptr, /*size=*/16,
                  [&error_deleter_calls]() { ++error_deleter_calls; }),
              StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("null")));
  EXPECT_EQ(error_deleter_calls, 1);

  EXPECT_THAT(
      AdoptImmutableBuffer(/*data=*/nullptr, /*size=*/16, /*deleter=*/nullptr),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("null")));

  int empty_null_deleter_calls = 0;
  ASSERT_OK_AND_ASSIGN(
      ImmutableBuffer empty_from_null,
      AdoptImmutableBuffer(
          /*data=*/nullptr, /*size=*/0,
          [&empty_null_deleter_calls]() { ++empty_null_deleter_calls; }));
  EXPECT_THAT(empty_from_null, IsValidAllocatedBuffer(/*expected_size=*/0));
  const uint8_t* const empty_dummy_ptr = empty_from_null.data();
  EXPECT_EQ(empty_null_deleter_calls, 0);
  {
    ImmutableBuffer empty_slice =
        empty_from_null.Slice(/*offset=*/0, /*length=*/0);
    EXPECT_THAT(empty_slice,
                IsValidView(empty_dummy_ptr,
                            /*expected_size=*/0, ImmutableBuffer::kAlignment));
    empty_from_null = ImmutableBuffer();
    EXPECT_EQ(empty_null_deleter_calls, 0);
  }
  EXPECT_EQ(empty_null_deleter_calls, 1);

  const uint8_t non_null_byte = 7;
  int empty_non_null_deleter_calls = 0;
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer empty_from_non_null,
                       AdoptImmutableBuffer(&non_null_byte, /*size=*/0,
                                            [&empty_non_null_deleter_calls]() {
                                              ++empty_non_null_deleter_calls;
                                            }));
  EXPECT_THAT(empty_from_non_null,
              IsValidView(empty_dummy_ptr, /*expected_size=*/0,
                          ImmutableBuffer::kAlignment));
  EXPECT_NE(empty_from_non_null.data(), &non_null_byte);
  EXPECT_EQ(empty_non_null_deleter_calls, 0);
  EXPECT_FALSE(
      buffer_internal::PrivateAccessor::TryUnfreeze(empty_from_non_null)
          .has_value());
  ASSERT_OK_AND_ASSIGN(MutableBuffer unfrozen_empty,
                       UnfreezeOrCopy(std::move(empty_from_non_null)));
  EXPECT_THAT(empty_from_non_null, IsEmptyAndInvalid());
  EXPECT_EQ(empty_non_null_deleter_calls, 1);
  EXPECT_THAT(unfrozen_empty, IsValidAllocatedBuffer(/*expected_size=*/0));
  EXPECT_NE(unfrozen_empty.data(), empty_dummy_ptr);

  EXPECT_THAT(
      AdoptImmutableBuffer(/*data=*/nullptr, /*size=*/0, /*deleter=*/nullptr),
      IsOkAndHolds(IsValidView(empty_dummy_ptr, /*expected_size=*/0,
                               ImmutableBuffer::kAlignment)));
}

TEST(ImmutableBufferTest,
     AdoptImmutableBufferAssignmentReleasesOldBlockAfterReadingRhs) {
  struct Node {
    ImmutableBuffer next;
  };
  alignas(ImmutableBuffer::kAlignment) uint8_t first_bytes[16] = {1, 2, 3, 4};
  alignas(ImmutableBuffer::kAlignment) uint8_t second_bytes[16] = {5, 6, 7, 8};

  // Copy-assignment: LHS deleter destroys the object owning RHS.
  int copy_first_deleted = 0;
  int copy_second_deleted = 0;
  auto copy_node = std::make_unique<Node>();
  ASSERT_OK_AND_ASSIGN(copy_node->next,
                       AdoptImmutableBuffer(second_bytes, sizeof(second_bytes),
                                            [&copy_second_deleted]() {
                                              ++copy_second_deleted;
                                            }));
  const Node* const copy_node_ptr = copy_node.get();
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer copy_current,
                       AdoptImmutableBuffer(first_bytes, sizeof(first_bytes),
                                            [owned = std::move(copy_node),
                                             &copy_first_deleted]() mutable {
                                              owned.reset();
                                              ++copy_first_deleted;
                                            }));

  copy_current = copy_node_ptr->next;
  EXPECT_EQ(copy_first_deleted, 1);
  EXPECT_EQ(copy_second_deleted, 0);
  EXPECT_THAT(copy_current, IsValidView(second_bytes, sizeof(second_bytes),
                                        ImmutableBuffer::kAlignment));
  EXPECT_THAT(copy_current, ElementsAreArray(second_bytes));
  copy_current = ImmutableBuffer();
  EXPECT_EQ(copy_second_deleted, 1);

  // Move-assignment: LHS deleter destroys the object owning RHS.
  int move_first_deleted = 0;
  int move_second_deleted = 0;
  auto move_node = std::make_unique<Node>();
  ASSERT_OK_AND_ASSIGN(move_node->next,
                       AdoptImmutableBuffer(second_bytes, sizeof(second_bytes),
                                            [&move_second_deleted]() {
                                              ++move_second_deleted;
                                            }));
  Node* const move_node_ptr = move_node.get();
  ASSERT_OK_AND_ASSIGN(ImmutableBuffer move_current,
                       AdoptImmutableBuffer(first_bytes, sizeof(first_bytes),
                                            [owned = std::move(move_node),
                                             &move_first_deleted]() mutable {
                                              owned.reset();
                                              ++move_first_deleted;
                                            }));

  move_current = std::move(move_node_ptr->next);
  EXPECT_EQ(move_first_deleted, 1);
  EXPECT_EQ(move_second_deleted, 0);
  EXPECT_THAT(move_current, IsValidView(second_bytes, sizeof(second_bytes),
                                        ImmutableBuffer::kAlignment));
  EXPECT_THAT(move_current, ElementsAreArray(second_bytes));
  move_current = ImmutableBuffer();
  EXPECT_EQ(move_second_deleted, 1);
}

TEST(ImmutableBufferTest,
     AdoptImmutableBufferConcurrentCopySliceAndDropInvokesDeleterOnce) {
  const std::vector<uint8_t> source = {10, 20, 30, 40, 50, 60, 70, 80};
  std::atomic<int> deleter_calls{0};
  ASSERT_OK_AND_ASSIGN(
      ImmutableBuffer shared,
      AdoptImmutableBuffer(source.data(), source.size(), [&deleter_calls]() {
        deleter_calls.fetch_add(1, std::memory_order_relaxed);
      }));

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
  EXPECT_EQ(deleter_calls.load(std::memory_order_acquire), 1);
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

TEST(BufferControlBlockTest,
     MutableBufferMoveAssignmentReleasesOldBlockAfterMovingRhs) {
  class CallbackControlBlock final
      : public buffer_internal::BufferControlBlock {
   public:
    explicit CallbackControlBlock(
        absl::AnyInvocable<void() &&> on_destroy) noexcept
        : on_destroy_(std::move(on_destroy)) {}

    ~CallbackControlBlock() override {
      if (on_destroy_) {
        std::move(on_destroy_)();
      }
    }

    [[nodiscard]] bool IsWritable() const noexcept override { return true; }

   private:
    absl::AnyInvocable<void() &&> on_destroy_;
  };

  struct Node {
    MutableBuffer next;
  };
  alignas(MutableBuffer::kAlignment) uint8_t first_bytes[16] = {1, 2, 3, 4};
  alignas(MutableBuffer::kAlignment) uint8_t second_bytes[16] = {5, 6, 7, 8};

  int first_deleted = 0;
  int second_deleted = 0;
  auto node = std::make_unique<Node>();
  node->next = buffer_internal::PrivateAccessor::CreateMutableBuffer(
      second_bytes, sizeof(second_bytes),
      new TrackingControlBlock(&second_deleted), /*is_shm=*/false);
  Node* const node_ptr = node.get();

  MutableBuffer current = buffer_internal::PrivateAccessor::CreateMutableBuffer(
      first_bytes, sizeof(first_bytes),
      new CallbackControlBlock(
          [owned = std::move(node), &first_deleted]() mutable {
            owned.reset();
            ++first_deleted;
          }),
      /*is_shm=*/false);

  current = std::move(node_ptr->next);
  EXPECT_EQ(first_deleted, 1);
  EXPECT_EQ(second_deleted, 0);
  EXPECT_THAT(current, IsValidView(second_bytes, sizeof(second_bytes),
                                   MutableBuffer::kAlignment));
  EXPECT_THAT(current, ElementsAreArray(second_bytes));

  current = MutableBuffer();
  EXPECT_EQ(second_deleted, 1);
}

}  // namespace
}  // namespace intrinsic
