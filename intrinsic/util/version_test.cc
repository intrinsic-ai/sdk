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

#include "intrinsic/util/version.h"

#include <gtest/gtest.h>

#include <functional>
#include <string_view>

namespace intrinsic {
namespace {

template <typename T>
class VersionTestEqual : public ::testing::Test {};

TYPED_TEST_SUITE_P(VersionTestEqual);

TYPED_TEST_P(VersionTestEqual, Split) {
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1.2..3")));
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1__2_3")));
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1/2///3")));
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1--2-3")));
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1,2,3")));
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1-2/3.0")));
}

TYPED_TEST_P(VersionTestEqual, Generic) {
  EXPECT_TRUE(TypeParam()(Version("1.1.13-alpha2"), Version("1.1.13-alpha02")));
  EXPECT_TRUE(TypeParam()(Version("1.0030"), Version("1.30")));
  EXPECT_TRUE(TypeParam()(Version("1.03"), Version("1.3")));
  EXPECT_TRUE(TypeParam()(Version("254.83"), Version("254.83.0")));
  EXPECT_TRUE(TypeParam()(Version("9.103"), Version("9.103.0.0")));
  EXPECT_TRUE(TypeParam()(Version("6.0"), Version("6")));
  EXPECT_TRUE(TypeParam()(Version("0"), Version("0000.000")));
  EXPECT_TRUE(TypeParam()(Version("Manufacturer"), Version(" Manufacturer ")));
  EXPECT_TRUE(TypeParam()(Version("   "), Version("")));
  EXPECT_TRUE(TypeParam()(Version(""), Version("")));
}

TYPED_TEST_P(VersionTestEqual, SemVer) {
  EXPECT_TRUE(TypeParam()(Version("3.2.0-beta"), Version("3.2.0-beta")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-alpha"), Version("4.2.0-alpha")));
  EXPECT_TRUE(TypeParam()(Version("1.2.0+bar"), Version("1.2.0+baz")));
}

REGISTER_TYPED_TEST_SUITE_P(VersionTestEqual, Split, Generic, SemVer);
using EqualTypes =
    ::testing::Types<std::equal_to<Version>, std::less_equal<Version>,
                     std::greater_equal<Version>>;
INSTANTIATE_TYPED_TEST_SUITE_P(VersionTestEqual, VersionTestEqual, EqualTypes);

template <typename T>
class VersionTestLess : public ::testing::Test {};

TYPED_TEST_SUITE_P(VersionTestLess);

TYPED_TEST_P(VersionTestLess, Generic) {
  EXPECT_TRUE(TypeParam()(Version("1.1.13"), Version("1.1.13a1")));
  EXPECT_TRUE(TypeParam()(Version("1.1.13"), Version("1.1.14-alpha12")));
  EXPECT_TRUE(TypeParam()(Version("0"), Version("20")));
  EXPECT_TRUE(TypeParam()(Version("2"), Version("20")));
  EXPECT_TRUE(TypeParam()(Version("3.1"), Version("31")));
  EXPECT_TRUE(TypeParam()(Version("4"), Version("21")));
  EXPECT_TRUE(TypeParam()(Version(" 3g6"), Version("12g")));
  EXPECT_TRUE(TypeParam()(Version("1.24++"), Version("2.6.kw")));
  EXPECT_TRUE(TypeParam()(Version("1.2t"), Version("2024.06.05")));
  EXPECT_TRUE(TypeParam()(Version("5.8"), Version("5.9.0")));
  EXPECT_TRUE(TypeParam()(Version("1"), Version("1.1")));
  EXPECT_TRUE(TypeParam()(Version(""), Version("0")));
  EXPECT_TRUE(TypeParam()(Version("14.0.BABA1E3C"), Version("14.0.D83C7207")));
  EXPECT_TRUE(TypeParam()(Version("13.1.794391F9"), Version("14.1.D83C7207")));
}
TYPED_TEST_P(VersionTestLess, SemVer) {
  EXPECT_TRUE(TypeParam()(Version("1.2.3"), Version("1.5.1")));
  EXPECT_TRUE(TypeParam()(Version("1.0.0-0.0.0"), Version("1.0.0")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta"), Version("4.2.0")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-alpha"), Version("4.2.0-beta")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta"), Version("4.2.0-beta.2")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta"), Version("4.2.0-beta.foo")));
  EXPECT_TRUE(TypeParam()(Version("1.0.0-beta.4"), Version("1.0.0-beta.-2")));
  EXPECT_TRUE(TypeParam()(Version("1.0.0-beta.-2"), Version("1.0.0-beta.-3")));
}

REGISTER_TYPED_TEST_SUITE_P(VersionTestLess, Generic, SemVer);
using LessTypes =
    ::testing::Types<std::not_equal_to<Version>, std::less<Version>,
                     std::less_equal<Version>>;
INSTANTIATE_TYPED_TEST_SUITE_P(VersionTestLess, VersionTestLess, LessTypes);

template <typename T>
class VersionTestGreater : public ::testing::Test {};

TYPED_TEST_SUITE_P(VersionTestGreater);

TYPED_TEST_P(VersionTestGreater, Generic) {
  EXPECT_TRUE(TypeParam()(Version("1.1.13a1"), Version("1.1.13")));
  EXPECT_TRUE(TypeParam()(Version("1.1.14-alpha12"), Version("1.1.13")));
  EXPECT_TRUE(TypeParam()(Version("20"), Version("0")));
  EXPECT_TRUE(TypeParam()(Version("20"), Version("2")));
  EXPECT_TRUE(TypeParam()(Version("31"), Version("3.1")));
  EXPECT_TRUE(TypeParam()(Version("21"), Version("4")));
  EXPECT_TRUE(TypeParam()(Version("12g"), Version(" 3g6")));
  EXPECT_TRUE(TypeParam()(Version("2.6.kw"), Version("1.24++")));
  EXPECT_TRUE(TypeParam()(Version("2024.06.05"), Version("1.2t")));
  EXPECT_TRUE(TypeParam()(Version("5.9.0"), Version("5.8")));
  EXPECT_TRUE(TypeParam()(Version("1.1"), Version("1")));
  EXPECT_TRUE(TypeParam()(Version("0"), Version("")));
  EXPECT_TRUE(TypeParam()(Version("14.0.D83C7207"), Version("14.0.BABA1E3C")));
  EXPECT_TRUE(TypeParam()(Version("14.1.D83C7207"), Version("13.1.794391F9")));
}

TYPED_TEST_P(VersionTestGreater, SemVer) {
  EXPECT_TRUE(TypeParam()(Version("2.2.3"), Version("1.5.1")));
  EXPECT_TRUE(TypeParam()(Version("2.2.3"), Version("2.2.2")));
  EXPECT_TRUE(TypeParam()(Version("1.3.0"), Version("1.1.4")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0"), Version("4.2.0-beta")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta.2"), Version("4.2.0-beta.1")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta2"), Version("4.2.0-beta1")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta.2"), Version("4.2.0-beta")));
  EXPECT_TRUE(TypeParam()(Version("4.2.0-beta.foo"), Version("4.2.0-beta")));
  EXPECT_TRUE(TypeParam()(Version("1.0.0-beta.-3"), Version("1.0.0-beta.5")));
}

REGISTER_TYPED_TEST_SUITE_P(VersionTestGreater, Generic, SemVer);
using GreaterTypes =
    ::testing::Types<std::not_equal_to<Version>, std::greater<Version>,
                     std::greater_equal<Version>>;
INSTANTIATE_TYPED_TEST_SUITE_P(VersionTestGreater, VersionTestGreater,
                               GreaterTypes);

}  // namespace
}  // namespace intrinsic
