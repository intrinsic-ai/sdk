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

#include "intrinsic/world/hashing/hashing.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstddef>
#include <cstdint>
#include <limits>
#include <set>
#include <string>
#include <string_view>
#include <tuple>
#include <utility>
#include <vector>

#include "absl/hash/hash_testing.h"
#include "intrinsic/world/entity_id.h"
#include "intrinsic/world/labels.h"

namespace intrinsic {
namespace {

using ::testing::Eq;
using ::testing::Gt;
using ::testing::SizeIs;

// Adapts a value hashed with WorldHasher<T> for use with
// absl::VerifyTypeImplementsAbslHashCorrectly. Because WorldHasher is a custom
// Hash functor rather than an AbslHashValue implementation, this test-only
// wrapper feeds the WorldHasher output into SpyHashState to verify determinism,
// congruence with operator==, and absence of collisions across test cases.
template <typename T, typename Hasher = WorldHasher<T>>
struct WorldHashed {
  T value;

  friend bool operator==(const WorldHashed& lhs, const WorldHashed& rhs) {
    return lhs.value == rhs.value;
  }

  template <typename H>
  friend H AbslHashValue(H h, const WorldHashed& w) {
    return H::combine(std::move(h), Hasher{}(w.value));
  }
};

template <typename T>
WorldHashed<T> WrapWorldHash(T value) {
  return WorldHashed<T>{std::move(value)};
}

TEST(WorldHasherTest, EntityIdImplementsHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(kInvalidEntityId),
      WrapWorldHash(EntityId(0)),
      WrapWorldHash(EntityId(1)),
      WrapWorldHash(EntityId(2)),
      WrapWorldHash(EntityId((1u << 16) | 2u)),
      WrapWorldHash(EntityId((2u << 16) | 2u)),
      WrapWorldHash(EntityId(std::numeric_limits<uint32_t>::max())),
  }));
}

TEST(WorldHasherTest, TypedEntityIdImplementsHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(EntityHandle(0)),
      WrapWorldHash(EntityHandle(1)),
      WrapWorldHash(EntityHandle((1u << 16) | 1u)),
      WrapWorldHash(EntityHandle(std::numeric_limits<uint32_t>::max())),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(AttachmentEntityId(0)),
      WrapWorldHash(AttachmentEntityId(1)),
      WrapWorldHash(AttachmentEntityId((1u << 16) | 1u)),
      WrapWorldHash(AttachmentEntityId(std::numeric_limits<uint32_t>::max())),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(PhysicalEntityId(0)),
      WrapWorldHash(PhysicalEntityId(1)),
      WrapWorldHash(PhysicalEntityId((1u << 16) | 1u)),
      WrapWorldHash(PhysicalEntityId(std::numeric_limits<uint32_t>::max())),
  }));
}

TEST(WorldHasherTest, IntegralTypesImplementHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(false),
      WrapWorldHash(true),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(std::numeric_limits<int32_t>::min()),
      WrapWorldHash(int32_t{-1}),
      WrapWorldHash(int32_t{0}),
      WrapWorldHash(int32_t{1}),
      WrapWorldHash(std::numeric_limits<int32_t>::max()),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(uint64_t{0}),
      WrapWorldHash(uint64_t{1}),
      WrapWorldHash(uint64_t{1} << 32),
      WrapWorldHash(uint64_t{1} << 33),
      WrapWorldHash(uint64_t{1} << 63),
      WrapWorldHash(std::numeric_limits<uint64_t>::max() - 1),
      WrapWorldHash(std::numeric_limits<uint64_t>::max()),
  }));
}

TEST(WorldHasherTest, EnumTypesImplementHashCorrectly) {
  enum UnscopedEnum { kUnscopedA = -1, kUnscopedB = 0, kUnscopedC = 42 };
  enum class ScopedEnum : int64_t {
    kMin = std::numeric_limits<int64_t>::min(),
    kZero = 0,
    kOne = 1,
    kMax = std::numeric_limits<int64_t>::max(),
  };

  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(kUnscopedA),
      WrapWorldHash(kUnscopedB),
      WrapWorldHash(kUnscopedC),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(ScopedEnum::kMin),
      WrapWorldHash(ScopedEnum::kZero),
      WrapWorldHash(ScopedEnum::kOne),
      WrapWorldHash(ScopedEnum::kMax),
  }));
}

TEST(WorldHasherTest, PointerTypesImplementHashCorrectly) {
  int values[4] = {0, 1, 2, 3};
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(static_cast<int*>(nullptr)),
      WrapWorldHash(&values[0]),
      WrapWorldHash(&values[1]),
      WrapWorldHash(&values[2]),
      WrapWorldHash(&values[3]),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(static_cast<const int*>(nullptr)),
      WrapWorldHash(static_cast<const int*>(&values[0])),
      WrapWorldHash(static_cast<const int*>(&values[1])),
  }));
}

TEST(WorldHasherTest, StringHeterogeneousLookupImplementsHashCorrectly) {
  auto eq = [](const auto& lhs, const auto& rhs) {
    return std::string_view(lhs.value) == std::string_view(rhs.value);
  };
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly(
      std::make_tuple(
          WorldHashed<std::string, WorldHasher<std::string>>{""},
          WorldHashed<std::string_view, WorldHasher<std::string>>{""},
          WorldHashed<const char*, WorldHasher<std::string>>{""},
          WorldHashed<std::string, WorldHasher<std::string>>{"alpha"},
          WorldHashed<std::string_view, WorldHasher<std::string>>{"alpha"},
          WorldHashed<const char*, WorldHasher<std::string>>{"alpha"},
          WorldHashed<std::string, WorldHasher<std::string>>{"beta"},
          WorldHashed<std::string_view, WorldHasher<std::string>>{"beta"},
          WorldHashed<const char*, WorldHasher<std::string>>{"beta"}),
      eq));
}

TEST(WorldHasherTest, LabelIdImplementsHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(LabelId("")),
      WrapWorldHash(LabelId("label1")),
      WrapWorldHash(LabelId("label2")),
  }));
}

TEST(WorldHasherTest, PairAndTupleImplementHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(std::make_pair(EntityId(0), EntityId(0))),
      WrapWorldHash(std::make_pair(EntityId(1), EntityId(2))),
      WrapWorldHash(std::make_pair(EntityId(1), EntityId(2))),
      WrapWorldHash(std::make_pair(EntityId(2), EntityId(1))),
      WrapWorldHash(std::make_pair(EntityId(2), EntityId(2))),
  }));
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(std::make_tuple(EntityId(0), AttachmentEntityId(0))),
      WrapWorldHash(std::make_tuple(EntityId(1), AttachmentEntityId(2))),
      WrapWorldHash(std::make_tuple(EntityId(1), AttachmentEntityId(2))),
      WrapWorldHash(std::make_tuple(EntityId(2), AttachmentEntityId(1))),
  }));
}

TEST(WorldHasherTest, FlatHashSetImplementsHashCorrectly) {
  EXPECT_TRUE(absl::VerifyTypeImplementsAbslHashCorrectly({
      WrapWorldHash(WorldHashSet<EntityId>{}),
      WrapWorldHash(WorldHashSet<EntityId>{EntityId(1)}),
      WrapWorldHash(WorldHashSet<EntityId>{EntityId(2)}),
      WrapWorldHash(WorldHashSet<EntityId>{EntityId(1), EntityId(2)}),
      WrapWorldHash(WorldHashSet<EntityId>{EntityId(2), EntityId(1)}),
      WrapWorldHash(
          WorldHashSet<EntityId>{EntityId(1), EntityId(2), EntityId(3)}),
      WrapWorldHash(
          WorldHashSet<EntityId>{EntityId(3), EntityId(1), EntityId(2)}),
  }));
}

TEST(WorldHasherTest, MatchesFmix64ReferenceValues) {
  // Pins known MurmurHash3 fmix64 outputs to guard against accidental constant
  // typos and verify deterministic, cross-library hash stability (b/380032603).
  EXPECT_THAT(WorldHasher<uint64_t>{}(0ULL),
              Eq(static_cast<std::size_t>(0x0000000000000000ULL)));
  EXPECT_THAT(WorldHasher<uint64_t>{}(1ULL),
              Eq(static_cast<std::size_t>(0xb456bcfc34c2cb2cULL)));
  EXPECT_THAT(WorldHasher<uint64_t>{}(0x00010002ULL),
              Eq(static_cast<std::size_t>(0x392fdca3b6a8fa99ULL)));
  EXPECT_THAT(WorldHasher<uint64_t>{}(0xffffffffULL),
              Eq(static_cast<std::size_t>(0xcc71ecda2aa8bcc6ULL)));
  EXPECT_THAT(WorldHasher<uint64_t>{}(0xffffffffffffffffULL),
              Eq(static_cast<std::size_t>(0x64b5720b4b825f21ULL)));
  EXPECT_THAT(WorldHasher<EntityId>{}(EntityId((1u << 16) | 2u)),
              Eq(static_cast<std::size_t>(0x392fdca3b6a8fa99ULL)));
}

TEST(WorldHasherTest, TypedEntityIdMatchesEntityIdHash) {
  for (const uint32_t raw_id :
       {0u, 1u, (1u << 16) | 2u, std::numeric_limits<uint32_t>::max()}) {
    const EntityId id(raw_id);
    const std::size_t expected = WorldHasher<EntityId>{}(id);
    EXPECT_THAT(WorldHasher<EntityHandle>{}(EntityHandle(id)), Eq(expected));
    EXPECT_THAT(WorldHasher<AttachmentEntityId>{}(AttachmentEntityId(id)),
                Eq(expected));
    EXPECT_THAT(WorldHasher<PhysicalEntityId>{}(PhysicalEntityId(id)),
                Eq(expected));
  }
}

}  // namespace
}  // namespace intrinsic
