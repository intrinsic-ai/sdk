// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/utils/log.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "absl/strings/str_cat.h"
#include "intrinsic/icon/release/source_location.h"
#include "intrinsic/icon/utils/log_internal.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::icon {
namespace {

using ::testing::AllOf;
using ::testing::ElementsAre;
using ::testing::HasSubstr;
using ::testing::IsEmpty;
using ::testing::StartsWith;
using ::testing::StrEq;

class FakeLogger : public LogSinkInterface {
 public:
  void Log(const LogEntry& entry) override {
    messages_.emplace_back(entry.msg);
    char buffer[kLogMessageMaxSize] = {0};
    LogEntryFormatToBuffer(buffer, sizeof(buffer), entry);
    text_.emplace_back(std::string(buffer));
  }
  std::vector<std::string> text_;
  std::vector<std::string> messages_;
};

TEST(IconUtilsLogTest, LogsSeverityLineFilename) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  double d = 0.5;
  INTRINSIC_RT_LOG(INFO) << "dof:" << 3 << " d:" << d;
  auto location = INTRINSIC_LOC;
  std::string expected_line_number = absl::StrCat(":", location.line() - 1);
  EXPECT_THAT(logger->messages_, ElementsAre(StrEq("dof:3 d:0.5")));
  EXPECT_THAT(
      logger->text_,
      ElementsAre(AllOf(StartsWith("I"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("dof:3 d:0.5"))));
}

TEST(IconUtilsLogTest, Throttles) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  INTRINSIC_RT_LOG_THROTTLED(WARNING) << "logged1";
  auto location = INTRINSIC_LOC;
  std::string expected_line_number = absl::StrCat(":", location.line() - 1);
  EXPECT_THAT(
      logger->text_,
      ElementsAre(AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged1"))));
  INTRINSIC_RT_LOG_THROTTLED(WARNING) << "logged2";
  for (int i = 0; i < 2000; ++i) {
    INTRINSIC_RT_LOG_THROTTLED(INFO) << "rarely logged i:" << i;
  }
  EXPECT_THAT(
      logger->messages_,
      ElementsAre(StrEq("logged1"), StrEq("logged2"),
                  StartsWith("rarely logged i:0"),
                  StartsWith("rarely logged i:1000 (repeated 1000 times in ")));
}

TEST(IconUtilsLogTest, DefaultLoggerDoesNotAllocate) {
  GlobalLogContext::SetThreadLocalLogSink(nullptr);
  RtLogInitForThisThread();
  double d = 0.5;
  std::string s = "text";
  INTRINSIC_RT_LOG(INFO) << "dof:" << 3 << " d:" << d;
  INTRINSIC_RT_LOG(ERROR) << "error: " << s;
  INTRINSIC_RT_LOG_THROTTLED(WARNING) << "logged1";
  for (int i = 0; i < 2000; ++i) {
    INTRINSIC_RT_LOG_THROTTLED(WARNING) << "rarely logged i:" << i;
  }
}

TEST(IconUtilsLogTest, DefaultLoggerIsThreadSafe) {
  GlobalLogContext::SetThreadLocalLogSink(nullptr);
  RtLogInitForThisThread();
  auto worker = []() {
    for (int i = 0; i < 10; ++i) {
      std::string s = "text";
      INTRINSIC_RT_LOG(INFO) << " i:" << i;
      INTRINSIC_RT_LOG(ERROR) << "error: " << s;
      INTRINSIC_RT_LOG_THROTTLED(WARNING) << "logged1";
      for (int j = 0; j < 2000; ++j) {
        INTRINSIC_RT_LOG_THROTTLED(WARNING) << "rarely logged j:" << j;
      }
    }
  };
  intrinsic::Thread worker1_thread{worker};
  intrinsic::Thread worker2_thread{worker};
  worker1_thread.join();
  worker2_thread.join();
}

TEST(IconUtilsLogTest, LogFirstWorks) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  for (int i = 0; i < 10; ++i) {
    INTRINSIC_RT_LOG_FIRST(WARNING) << "logged " << i;
  }
  auto location = INTRINSIC_LOC;
  std::string expected_line_number = absl::StrCat(":", location.line() - 2);
  EXPECT_EQ(logger->messages_.size(), 1);
  EXPECT_THAT(
      logger->text_,
      ElementsAre(AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 0"))));
}

TEST(IconUtilsLogTest, LogFirstDoesNotAllocate) {
  GlobalLogContext::SetThreadLocalLogSink(nullptr);
  RtLogInitForThisThread();
  double d = 0.5;
  std::string s = "text";
  INTRINSIC_RT_LOG_FIRST(INFO) << "dof:" << 3 << " d:" << d;
  INTRINSIC_RT_LOG_FIRST(ERROR) << "error: " << s;
  INTRINSIC_RT_LOG_FIRST(WARNING) << "logged1";
  for (int i = 0; i < 2000; ++i) {
    INTRINSIC_RT_LOG_FIRST(WARNING) << "rarely logged i:" << i;
  }
}

TEST(IconUtilsLogTest, LogFirstNWorks) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  for (int i = 0; i < 10; ++i) {
    INTRINSIC_RT_LOG_FIRST_N(WARNING, 5) << "logged " << i;
  }

  EXPECT_EQ(logger->messages_.size(), 5);
  auto location = INTRINSIC_LOC;
  std::string expected_line_number = absl::StrCat(":", location.line() - 4);
  EXPECT_THAT(
      logger->text_,
      ElementsAre(AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 0")),
                  AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 1")),
                  AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 2")),
                  AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 3")),
                  AllOf(StartsWith("W"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("logged 4"))));
}

TEST(IconUtilsLogTest, LogFirstNDoesNotAllocate) {
  GlobalLogContext::SetThreadLocalLogSink(nullptr);
  RtLogInitForThisThread();
  double d = 0.5;
  std::string s = "text";
  INTRINSIC_RT_LOG_FIRST_N(INFO, 5) << "dof:" << 3 << " d:" << d;
  INTRINSIC_RT_LOG_FIRST_N(ERROR, 5) << "error: " << s;
  INTRINSIC_RT_LOG_FIRST_N(WARNING, 5) << "logged1";
  for (int i = 0; i < 2000; ++i) {
    INTRINSIC_RT_LOG_FIRST_N(WARNING, 5) << "rarely logged i:" << i;
  }
}

TEST(IconUtilsLogTest, LogBackoffDoesNotAllocate) {
  GlobalLogContext::SetThreadLocalLogSink(nullptr);
  RtLogInitForThisThread();
  double d = 0.5;
  std::string s = "text";
  INTRINSIC_RT_LOG_BACKOFF(INFO, false) << "dof:" << 3 << " d:" << d;
  INTRINSIC_RT_LOG_BACKOFF(ERROR, false) << "error: " << s;
  INTRINSIC_RT_LOG_BACKOFF(WARNING, false) << "logged1";
  for (int i = 0; i < 2000; ++i) {
    INTRINSIC_RT_LOG_BACKOFF(WARNING, false) << "rarely logged i:" << i;
  }
  INTRINSIC_RT_LOG_BACKOFF(WARNING, true) << "manual reset";
}

TEST(IconUtilsLogTest, LogIfWorksWhenTrue) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  double d = 0.5;
  INTRINSIC_RT_LOG_IF(INFO, true) << "dof:" << 3 << " d:" << d;
  auto location = INTRINSIC_LOC;
  std::string expected_line_number = absl::StrCat(":", location.line() - 1);
  EXPECT_THAT(logger->messages_, ElementsAre(StrEq("dof:3 d:0.5")));
  EXPECT_THAT(
      logger->text_,
      ElementsAre(AllOf(StartsWith("I"), HasSubstr(expected_line_number),
                        HasSubstr("log_test.cc"), HasSubstr("dof:3 d:0.5"))));
}

TEST(IconUtilsLogTest, LogIfDoesNotLogWhenFalse) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));
  EXPECT_THAT(logger->messages_, IsEmpty());
  double d = 0.5;
  INTRINSIC_RT_LOG_IF(INFO, false) << "dof:" << 3 << " d:" << d;
  EXPECT_THAT(logger->messages_, IsEmpty());
  EXPECT_THAT(logger->text_, IsEmpty());
}

static int64_t fake_robot_time_ns = 0;
static int64_t fake_wall_time_ns = 0;

void FakeTimeFunction(int64_t* robot_timestamp_ns, int64_t* wall_timestamp_ns) {
  *robot_timestamp_ns = fake_robot_time_ns;
  *wall_timestamp_ns = fake_wall_time_ns;
}

TEST(IconUtilsLogTest, Backoff) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));

  GlobalLogContext::SetTimeFunction(FakeTimeFunction);
  // Use a large enough start time.
  fake_robot_time_ns =
      100 * internal::LogBackoffThrottler::kMaxPeriodNanoseconds;

  const int64_t kP0 = internal::LogBackoffThrottler::kInitialPeriodNanoseconds;

  // Use a loop to ensure we hit the same call site (same static throttler).
  for (int i = 0; i < 3; ++i) {
    INTRINSIC_RT_LOG_BACKOFF(INFO, false) << "backoff_msg";
    if (i == 0) {
      // 1. First call should always log.
      EXPECT_THAT(logger->messages_, ElementsAre("backoff_msg"));
      fake_robot_time_ns += kP0 / 5;
    } else if (i == 1) {
      // 2. Immediate second call should be throttled.
      EXPECT_THAT(logger->messages_, ElementsAre("backoff_msg"));
      fake_robot_time_ns += kP0;
    } else {
      // 3. Call after kP0 should log and multiply period by factor.
      EXPECT_THAT(logger->messages_,
                  ElementsAre("backoff_msg",
                              StartsWith("backoff_msg (repeated 2 times")));
    }
  }

  GlobalLogContext::SetTimeFunction(nullptr);
}

TEST(IconUtilsLogTest, BackoffReset) {
  auto unique_logger = std::make_unique<FakeLogger>();
  auto* logger = unique_logger.get();
  GlobalLogContext::SetThreadLocalLogSink(std::move(unique_logger));

  GlobalLogContext::SetTimeFunction(FakeTimeFunction);
  fake_robot_time_ns =
      100 * internal::LogBackoffThrottler::kMaxPeriodNanoseconds;

  const int64_t kP0 = internal::LogBackoffThrottler::kInitialPeriodNanoseconds;
  const int kBackoffFactor = internal::LogBackoffThrottler::kBackoffFactor;
  const int kResetFactor = internal::LogBackoffThrottler::kResetFactor;

  // Use a loop to hit the same call site.
  for (int i = 0; i < 4; ++i) {
    INTRINSIC_RT_LOG_BACKOFF(INFO, false) << "reset_msg";
    if (i == 0) {
      fake_robot_time_ns += kP0;
    } else if (i == 1) {
      // Now period is kP0 * kBackoffFactor. Wait for inactivity reset.
      fake_robot_time_ns += kP0 * kBackoffFactor * kResetFactor + kP0;
    } else if (i == 2) {
      // Should log and reset period to kP0. Now wait for kP0.
      fake_robot_time_ns += kP0 * 1.1;
    } else {
      // Should log.
      EXPECT_EQ(logger->messages_.size(), 4);
    }
  }

  // Test manual reset via macro
  INTRINSIC_RT_LOG_BACKOFF(INFO, true) << "manual_reset";
  EXPECT_EQ(logger->messages_.size(), 5);
  EXPECT_THAT(logger->messages_.back(), StrEq("manual_reset"));

  GlobalLogContext::SetTimeFunction(nullptr);
}

}  // namespace
}  // namespace intrinsic::icon
