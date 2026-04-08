// Copyright 2023 Intrinsic Innovation LLC


#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstdint>
#include <optional>

#include "intrinsic/icon/utils/log_internal.h"

namespace intrinsic::icon::internal {
namespace {

struct FakeTime {
  int64_t robot_time_ns = 0;
  int64_t wall_time_ns = 0;

  void GetTime(int64_t* robot_timestamp_ns, int64_t* wall_timestamp_ns) const {
    *robot_timestamp_ns = robot_time_ns;
    *wall_timestamp_ns = wall_time_ns;
  }
};

TEST(LogBackoffThrottlerTest, BackoffDoublesAndResets) {
  LogBackoffThrottler throttler;
  FakeTime fake_time;
  const int64_t kStartTime = 100 * LogBackoffThrottler::kMaxPeriodNanoseconds;
  fake_time.robot_time_ns = kStartTime;

  auto get_time = [&](int64_t* r, int64_t* w) { fake_time.GetTime(r, w); };

  const int64_t kP0 = LogBackoffThrottler::kInitialPeriodNanoseconds;
  const int64_t kMaxP = LogBackoffThrottler::kMaxPeriodNanoseconds;
  const int kFactor = LogBackoffThrottler::kBackoffFactor;

  // 1. First call should always log.
  auto res1 = throttler.Tick(get_time);
  ASSERT_TRUE(res1.has_value());
  EXPECT_EQ(res1->num_calls_merged, 1);

  // 2. Immediate second call should be throttled.
  // Next expected log is after kP0.
  fake_time.robot_time_ns = kStartTime + kP0 * 0.1;
  auto res2 = throttler.Tick(get_time);
  EXPECT_FALSE(res2.has_value());

  // 3. Call after kP0 should log and multiply period by factor.
  fake_time.robot_time_ns = kStartTime + kP0;
  auto res3 = throttler.Tick(get_time);
  ASSERT_TRUE(res3.has_value());
  EXPECT_EQ(res3->num_calls_merged, 2);

  // 4. Call after another 1.2*kP0 should be throttled (current period is
  // kFactor*kP0).
  fake_time.robot_time_ns = kStartTime + 2.2 * kP0;
  auto res4 = throttler.Tick(get_time);
  EXPECT_FALSE(res4.has_value());

  // 5. Call after another kP0 (total 2.2*kP0 since last log) should log.
  fake_time.robot_time_ns = kStartTime + 3.2 * kP0;
  auto res5 = throttler.Tick(get_time);
  ASSERT_TRUE(res5.has_value());

  // 6. Test Max Backoff: Skip ahead.
  // Reach kMaxP by repeatedly logging.
  int64_t current_p = kFactor * kFactor * kP0;
  int64_t last_log_time = fake_time.robot_time_ns;
  while (current_p < kMaxP) {
    fake_time.robot_time_ns = last_log_time + current_p;
    throttler.Tick(get_time);
    last_log_time = fake_time.robot_time_ns;
    current_p *= kFactor;
  }

  // Now period should be kMaxP.
  fake_time.robot_time_ns = last_log_time + kMaxP / 2;
  EXPECT_FALSE(throttler.Tick(get_time).has_value());

  fake_time.robot_time_ns =
      last_log_time + kMaxP;  // Total kMaxP total since last log.
  auto res_max = throttler.Tick(get_time);
  ASSERT_TRUE(res_max.has_value());
  last_log_time = fake_time.robot_time_ns;

  // 7. Test Inactivity Reset.
  // Period is currently kMaxP.
  // If we wait > kFactor * kMaxP, it should reset to kP0.
  fake_time.robot_time_ns = last_log_time + kFactor * kMaxP + kP0;
  throttler.Tick(get_time);  // This log will have delta > period * kFactor, so
                             // it should reset to kP0.
  last_log_time = fake_time.robot_time_ns;

  // If it reset to kP0, a log after 1.2*kP0 should work.
  fake_time.robot_time_ns = last_log_time + 1.2 * kP0;
  auto res_reset = throttler.Tick(get_time);
  EXPECT_TRUE(res_reset.has_value());
}

TEST(LogBackoffThrottlerTest, ManualReset) {
  LogBackoffThrottler throttler;
  FakeTime fake_time;
  fake_time.robot_time_ns = 100 * LogBackoffThrottler::kMaxPeriodNanoseconds;
  auto get_time = [&](int64_t* r, int64_t* w) { fake_time.GetTime(r, w); };

  const int64_t kStartTime = 100 * LogBackoffThrottler::kMaxPeriodNanoseconds;
  fake_time.robot_time_ns = kStartTime;
  const int64_t kP0 = LogBackoffThrottler::kInitialPeriodNanoseconds;
  const int kFactor = LogBackoffThrottler::kBackoffFactor;

  // 1. Initial log.
  throttler.Tick(get_time);

  // 2. Advance to double backoff.
  fake_time.robot_time_ns = kStartTime + kP0 * kFactor;
  throttler.Tick(get_time);  // Next period is kFactor*kP0.

  // 3. Next log would be throttled if reset=false.
  fake_time.robot_time_ns = kStartTime + kP0 * kFactor * 1.1;
  EXPECT_FALSE(throttler.Tick(get_time, /*reset=*/false).has_value());

  // 4. Force log and reset with reset=true.
  auto res = throttler.Tick(get_time, /*reset=*/true);
  ASSERT_TRUE(res.has_value());

  // 5. Verify it reset to kP0.
  fake_time.robot_time_ns = kStartTime + kP0 * kFactor * 1.1 + kP0;
  EXPECT_TRUE(throttler.Tick(get_time).has_value());
}

}  // namespace
}  // namespace intrinsic::icon::internal
