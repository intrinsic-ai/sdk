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

#include "intrinsic/icon/utils/bitset.h"

#include <gtest/gtest.h>

#include <cstdint>
#include <limits>
#include <type_traits>

namespace {

template <typename T>
class BitsetTest : public ::testing::Test {};

using MyTypes = ::testing::Types<bool, uint8_t, uint16_t, uint32_t, uint64_t>;
TYPED_TEST_SUITE(BitsetTest, MyTypes);

TYPED_TEST(BitsetTest, TestMinZeroMax) {
  using T = TypeParam;
  intrinsic::bitset<T> zero(static_cast<T>(0));
  intrinsic::bitset<T> min(std::numeric_limits<T>::min());
  intrinsic::bitset<T> max(std::numeric_limits<T>::max());

  EXPECT_EQ(intrinsic::GetValue<T>(zero), static_cast<T>(0));
  EXPECT_EQ(intrinsic::GetValue<T>(min), std::numeric_limits<T>::min());
  EXPECT_EQ(intrinsic::GetValue<T>(max), std::numeric_limits<T>::max());

  // Check that the internal representation corresponds to the value.
  EXPECT_EQ(zero.to_ullong(), static_cast<uint64_t>(0));
  EXPECT_EQ(min.to_ullong(),
            static_cast<uint64_t>(std::numeric_limits<T>::min()));
  EXPECT_EQ(max.to_ullong(),
            static_cast<uint64_t>(std::numeric_limits<T>::max()));
}

TYPED_TEST(BitsetTest, TestSizeMatchesExpectedSize) {
  using T = TypeParam;
  intrinsic::bitset<T> bitset;
  if constexpr (std::is_same_v<T, bool>) {
    EXPECT_EQ(bitset.size(), 1);
  } else {
    EXPECT_EQ(bitset.size(), 8 * sizeof(T));
  }
}

TYPED_TEST(BitsetTest, TestFromEnumAndToEnum) {
  using T = TypeParam;
  enum class MyEnum : T {
    kFoo = std::numeric_limits<T>::min(),
    kBar = 0,
    kBaz = std::numeric_limits<T>::max(),
  };

  intrinsic::from_enum_t<MyEnum> from_enum_foo =
      intrinsic::FromEnum(MyEnum::kFoo);
  EXPECT_EQ(intrinsic::GetValue<T>(from_enum_foo),
            std::numeric_limits<T>::min());
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(intrinsic::GetValue<T>(from_enum_foo)),
            MyEnum::kFoo);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(
                intrinsic::from_enum_t<MyEnum>(std::numeric_limits<T>::min())),
            MyEnum::kFoo);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(from_enum_foo), MyEnum::kFoo);

  intrinsic::from_enum_t<MyEnum> from_enum_bar =
      intrinsic::FromEnum(MyEnum::kBar);
  EXPECT_EQ(intrinsic::GetValue<T>(from_enum_bar), T{0});
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(intrinsic::GetValue<T>(from_enum_bar)),
            MyEnum::kBar);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(intrinsic::from_enum_t<MyEnum>(T{0})),
            MyEnum::kBar);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(from_enum_bar), MyEnum::kBar);

  intrinsic::from_enum_t<MyEnum> from_enum_baz =
      intrinsic::FromEnum(MyEnum::kBaz);
  EXPECT_EQ(intrinsic::GetValue<T>(from_enum_baz),
            std::numeric_limits<T>::max());
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(intrinsic::GetValue<T>(from_enum_baz)),
            MyEnum::kBaz);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(
                intrinsic::from_enum_t<MyEnum>(std::numeric_limits<T>::max())),
            MyEnum::kBaz);
  EXPECT_EQ(intrinsic::ToEnum<MyEnum>(from_enum_baz), MyEnum::kBaz);
}

}  // namespace
