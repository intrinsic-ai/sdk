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

#include "intrinsic/util/time/time.h"

#include <gtest/gtest.h>
#include <stdint.h>

#include <chrono>  // NOLINT(build/c++11)
#include <ratio>   // NOLINT(build/c++11)

#include "absl/time/time.h"

namespace intrinsic {
namespace {

TEST(TimeSteadyTest, ConstructionDestruction) {
  TimeSteady time_steady;
  (void)time_steady;
}

TEST(TimeSteadyTest, InfiniteFuture) {
  TimeSteady time_steady = TimeSteady::InfiniteFuture();
  (void)time_steady;
}

TEST(TimeSteadyTest, InfinitePast) {
  TimeSteady time_steady = TimeSteady::InfinitePast();
  (void)time_steady;
}

TEST(TimeSteadyTest, Now) {
  TimeSteady time_steady = TimeSteady::Now();
  (void)time_steady;
}

class TimeSteadyFiniteMathTest : public ::testing::Test {
 public:
  TimeSteadyFiniteMathTest() : duration_finite_(absl::Microseconds(1)) {}

  ~TimeSteadyFiniteMathTest() override = default;

 protected:
  const absl::Duration duration_finite_;
  const TimeSteady time_steady_finite_;
};

TEST_F(TimeSteadyFiniteMathTest, UnsaturatedAddAndSubtractNonNegativeOperands) {
  // test that add and subtract perform the expected inverse operations.
  const TimeSteady time_steady = time_steady_finite_ + duration_finite_;
  EXPECT_NE(time_steady, time_steady_finite_);
  EXPECT_EQ(time_steady - duration_finite_, time_steady_finite_);
  EXPECT_EQ(time_steady - time_steady_finite_, duration_finite_);
}

TEST_F(TimeSteadyFiniteMathTest, PlusEqualsPlusEquivalence) {
  TimeSteady time_steady = time_steady_finite_;
  time_steady += duration_finite_;
  EXPECT_EQ(time_steady, time_steady_finite_ + duration_finite_);
}

TEST_F(TimeSteadyFiniteMathTest, MinusEqualsMinusEquivalence) {
  TimeSteady time_steady = time_steady_finite_ + duration_finite_;
  TimeSteady expected_result = time_steady - duration_finite_;
  time_steady -= duration_finite_;
  EXPECT_EQ(time_steady, expected_result);
}

TEST_F(TimeSteadyFiniteMathTest, UnsaturatedAddAndSubtractResultInThePast) {
  const TimeSteady time_steady_past = time_steady_finite_ - duration_finite_;
  EXPECT_LT(time_steady_past, time_steady_finite_);
  EXPECT_EQ(time_steady_finite_ - time_steady_past, duration_finite_);
}

TEST_F(TimeSteadyFiniteMathTest, UnsaturatedAddAndSubtractNegativeDuration) {
  const absl::Duration neg_duration = -absl::Microseconds(1);
  const TimeSteady time_steady_future = time_steady_finite_ - neg_duration;
  EXPECT_GT(time_steady_future, time_steady_finite_);
  EXPECT_EQ(time_steady_finite_ - time_steady_future, neg_duration);
}

class TimeSteadySaturatedMathTest : public ::testing::Test {
 public:
  TimeSteadySaturatedMathTest()
      : time_steady_max_(TimeSteady::InfiniteFuture()),
        time_steady_min_(TimeSteady::InfinitePast()) {}

 protected:
  const TimeSteady time_steady_max_;
  const TimeSteady time_steady_min_;
};

TEST_F(TimeSteadySaturatedMathTest, AddFiniteDurationToInfiniteTimeSteady) {
  const absl::Duration duration_finite_positive(absl::Nanoseconds(1));
  const absl::Duration duration_finite_negative(-absl::Nanoseconds(1));

  EXPECT_EQ(time_steady_max_ + duration_finite_positive, time_steady_max_);
  EXPECT_EQ(time_steady_max_ + duration_finite_negative, time_steady_max_);
  EXPECT_EQ(time_steady_min_ + duration_finite_positive, time_steady_min_);
  EXPECT_EQ(time_steady_min_ + duration_finite_negative, time_steady_min_);
}

TEST_F(TimeSteadySaturatedMathTest,
       SubtractFiniteDurationFromInfiniteTimeSteady) {
  const absl::Duration duration_finite_positive(absl::Nanoseconds(1));
  const absl::Duration duration_finite_negative(-absl::Nanoseconds(1));

  EXPECT_EQ(time_steady_max_ - duration_finite_positive, time_steady_max_);
  EXPECT_EQ(time_steady_max_ - duration_finite_negative, time_steady_max_);
  EXPECT_EQ(time_steady_min_ - duration_finite_positive, time_steady_min_);
  EXPECT_EQ(time_steady_min_ - duration_finite_negative, time_steady_min_);
}

TEST_F(TimeSteadySaturatedMathTest, AddInfiniteDurationToInfiniteTimeSteady) {
  const absl::Duration duration_infinite_positive = absl::InfiniteDuration();
  const absl::Duration duration_infinite_negative = -absl::InfiniteDuration();

  EXPECT_EQ(time_steady_min_ + duration_infinite_positive, time_steady_min_);
  EXPECT_EQ(time_steady_max_ + duration_infinite_positive, time_steady_max_);
  EXPECT_EQ(time_steady_max_ + duration_infinite_negative, time_steady_max_);
  EXPECT_EQ(time_steady_min_ + duration_infinite_negative, time_steady_min_);
}

TEST_F(TimeSteadySaturatedMathTest,
       SubtractInfiniteDurationFromInfiniteTimeSteady) {
  const absl::Duration duration_infinite_positive = absl::InfiniteDuration();
  const absl::Duration duration_infinite_negative = -absl::InfiniteDuration();

  EXPECT_EQ(time_steady_min_ - duration_infinite_positive, time_steady_min_);
  EXPECT_EQ(time_steady_max_ - duration_infinite_positive, time_steady_max_);
  EXPECT_EQ(time_steady_max_ - duration_infinite_negative, time_steady_max_);
  EXPECT_EQ(time_steady_min_ - duration_infinite_negative, time_steady_min_);
}

TEST_F(TimeSteadySaturatedMathTest, TimeSteadySubtraction) {
  const TimeSteady time_steady_finite;
  const absl::Duration duration_infinite_positive = absl::InfiniteDuration();
  const absl::Duration duration_infinite_negative = -absl::InfiniteDuration();

  EXPECT_EQ(time_steady_max_ - time_steady_max_, duration_infinite_positive);
  EXPECT_EQ(time_steady_min_ - time_steady_min_, duration_infinite_negative);
  EXPECT_EQ(time_steady_max_ - time_steady_min_, duration_infinite_positive);
  EXPECT_EQ(time_steady_min_ - time_steady_max_, duration_infinite_negative);
  EXPECT_EQ(time_steady_max_ - time_steady_finite, duration_infinite_positive);
  EXPECT_EQ(time_steady_min_ - time_steady_finite, duration_infinite_negative);
  EXPECT_EQ(time_steady_finite - time_steady_max_, duration_infinite_negative);
  EXPECT_EQ(time_steady_finite - time_steady_min_, duration_infinite_positive);
}

class TimeSteadySaturatedMathResultTest : public ::testing::Test {
 public:
  TimeSteadySaturatedMathResultTest()
      : duration_large_positive_(absl::FromChrono(
            std::chrono::nanoseconds::max() - std::chrono::nanoseconds(1))),
        duration_large_negative_(absl::FromChrono(
            std::chrono::nanoseconds::min() + std::chrono::nanoseconds(1))),
        time_steady_large_positive_(TimeSteady() + duration_large_positive_),
        time_steady_large_negative_(TimeSteady() + duration_large_negative_) {
    // ensure that we are actually using non-infinite durations since we're
    // relying on absl::FromChrono to form these in the expected way
    EXPECT_NE(duration_large_positive_, absl::InfiniteDuration());
    EXPECT_NE(duration_large_negative_, -absl::InfiniteDuration());
  }

 protected:
  const absl::Duration duration_large_positive_;
  const absl::Duration duration_large_negative_;
  const TimeSteady time_steady_large_positive_;
  const TimeSteady time_steady_large_negative_;
};

TEST_F(TimeSteadySaturatedMathResultTest,
       InfiniteResultFromNonInfiniteOperandAddition) {
  EXPECT_EQ(time_steady_large_positive_ + duration_large_positive_,
            TimeSteady::InfiniteFuture());
  EXPECT_EQ(time_steady_large_negative_ + duration_large_negative_,
            TimeSteady::InfinitePast());
}

TEST_F(TimeSteadySaturatedMathResultTest,
       InfiniteResultFromNonInfiniteOperandSubtraction) {
  EXPECT_EQ(time_steady_large_negative_ - duration_large_positive_,
            TimeSteady::InfinitePast());
  EXPECT_EQ(time_steady_large_positive_ - duration_large_negative_,
            TimeSteady::InfiniteFuture());
}

TEST_F(TimeSteadySaturatedMathResultTest,
       InfiniteResultFromNonInfiniteTimeSteadySubtraction) {
  EXPECT_EQ(time_steady_large_positive_ - time_steady_large_negative_,
            absl::InfiniteDuration());
  EXPECT_EQ(time_steady_large_negative_ - time_steady_large_positive_,
            -absl::InfiniteDuration());
}

TEST(TimeSteadyTest, ConstexprTest) {
  constexpr TimeSteady kTimeSteadyInfiniteFuture = TimeSteady::InfiniteFuture();
  constexpr TimeSteady kTimeSteadyInfinitePast = TimeSteady::InfinitePast();
  constexpr TimeSteady kTimeSteadyDefault;
  constexpr TimeSteady kTimeSteadyCopy = kTimeSteadyDefault;
  constexpr TimeSteady kTimeSteadyCopyConstructed(kTimeSteadyDefault);

  static_assert(kTimeSteadyCopyConstructed == kTimeSteadyDefault, "");
  static_assert(kTimeSteadyCopy == kTimeSteadyDefault, "");
  static_assert(kTimeSteadyInfiniteFuture == kTimeSteadyInfiniteFuture, "");
  static_assert(kTimeSteadyInfiniteFuture != kTimeSteadyInfinitePast, "");
  static_assert(kTimeSteadyInfinitePast < kTimeSteadyInfiniteFuture, "");
  static_assert(kTimeSteadyInfinitePast <= kTimeSteadyInfiniteFuture, "");
  static_assert(kTimeSteadyInfiniteFuture > kTimeSteadyInfinitePast, "");
  static_assert(kTimeSteadyInfiniteFuture >= kTimeSteadyInfinitePast, "");
}

}  // namespace

TEST(TimeSteadyTest, Comparison) {
  // Relies on private constructor, since this lets us compare known
  // values without needing to rely on other operators of the class we are
  // defining, and doesn't force us to expose std::chrono to users of the
  // class.
  TimeSteady t_0_nanoeconds_past_epoch_(
      std::chrono::time_point<std::chrono::steady_clock,
                              std::chrono::duration<int64_t, std::nano>>(
          std::chrono::duration<int64_t, std::nano>(0)));
  TimeSteady t_1_nanoeconds_past_epoch_(
      std::chrono::time_point<std::chrono::steady_clock,
                              std::chrono::duration<int64_t, std::nano>>(
          std::chrono::duration<int64_t, std::nano>(1)));
  TimeSteady t_2_nanoeconds_past_epoch_(
      std::chrono::time_point<std::chrono::steady_clock,
                              std::chrono::duration<int64_t, std::nano>>(
          std::chrono::duration<int64_t, std::nano>(2)));

  EXPECT_EQ(t_0_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_EQ(t_1_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_EQ(t_2_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);

  EXPECT_NE(t_0_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_NE(t_0_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
  EXPECT_NE(t_1_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);

  EXPECT_GT(t_1_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_GT(t_2_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_GT(t_2_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);

  EXPECT_LT(t_0_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_LT(t_0_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
  EXPECT_LT(t_1_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);

  EXPECT_GE(t_0_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_GE(t_1_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_GE(t_2_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
  EXPECT_GE(t_1_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_GE(t_2_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_GE(t_2_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);

  EXPECT_LE(t_0_nanoeconds_past_epoch_, t_0_nanoeconds_past_epoch_);
  EXPECT_LE(t_1_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_LE(t_2_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
  EXPECT_LE(t_0_nanoeconds_past_epoch_, t_1_nanoeconds_past_epoch_);
  EXPECT_LE(t_0_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
  EXPECT_LE(t_1_nanoeconds_past_epoch_, t_2_nanoeconds_past_epoch_);
}

}  // namespace intrinsic
