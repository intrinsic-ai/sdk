// Copyright 2023 Intrinsic Innovation LLC

// Confidential and proprietary property of Intrinsic, a secret project of X,
// The Moonshot Factory.  Please note the project associated with this code is
// not publicly disclosed.  Please do not use this code without first contacting
// intrinsic-tech@.
#include "intrinsic/util/thread/lockstep.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <atomic>

#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/icon/utils/realtime_status_matchers.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/thread.h"

using absl_testing::StatusIs;

namespace intrinsic {

// Timeout for StartOperationX calls.
static constexpr absl::Duration kLockTimeout = absl::Milliseconds(100);
static constexpr absl::Duration kLongLockTimeout = absl::Milliseconds(500);

TEST(LockstepTest, StartOperationAWithTimeout) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
}

TEST(LockstepTest, StartOperationAWithDeadline) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout));
}

TEST(LockstepTest, StartOperationACancelled) {
  Lockstep lockstep;
  lockstep.Cancel();
  constexpr int kNumIterations = 5;
  for (int i = 0; i < kNumIterations; i++) {
    EXPECT_THAT(lockstep.StartOperationAWithTimeout(kLockTimeout),
                StatusIs(absl::StatusCode::kAborted));
  }
}

TEST(LockstepTest, StartOperationBCancelled) {
  Lockstep lockstep;
  lockstep.Cancel();
  constexpr int kNumIterations = 5;
  for (int i = 0; i < kNumIterations; i++) {
    EXPECT_THAT(lockstep.StartOperationBWithTimeout(kLockTimeout),
                StatusIs(absl::StatusCode::kAborted));
  }
}

TEST(LockstepTest, MismatchedEndOperationA) {
  Lockstep lockstep;

  EXPECT_THAT(lockstep.EndOperationA(),
              StatusIs(absl::StatusCode::kFailedPrecondition));

  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());

  EXPECT_THAT(lockstep.EndOperationA(),
              StatusIs(absl::StatusCode::kFailedPrecondition));
}

TEST(LockstepTest, MismatchedEndOperationB) {
  Lockstep lockstep;
  EXPECT_THAT(lockstep.EndOperationB(),
              StatusIs(absl::StatusCode::kFailedPrecondition));

  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout));
  EXPECT_OK(lockstep.EndOperationB());

  EXPECT_THAT(lockstep.EndOperationB(),
              StatusIs(absl::StatusCode::kFailedPrecondition));
}

TEST(LockstepTest, ABABABSingleThread) {
  Lockstep lockstep;
  constexpr int kNumIterations = 25000;
  for (int i = 0; i < kNumIterations; i++) {
    EXPECT_OK(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout));
    EXPECT_OK(lockstep.EndOperationA());
    EXPECT_OK(lockstep.StartOperationBWithTimeout(kLockTimeout));
    EXPECT_OK(lockstep.EndOperationB());
  }
}

TEST(LockstepTest, ABABABMultiThread) {
  Lockstep lockstep;

  static constexpr int kNumIterations = 25000;
  std::atomic<int> a_count = 0;
  std::atomic<int> b_count = 0;

  // Kick off thread for Operation A.
  Thread operation_a_thread([&lockstep, &a_count, &b_count]() {
    for (int i = 0; i < kNumIterations; i++) {
      EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
      CHECK_EQ(a_count, b_count);
      a_count++;
      EXPECT_OK(lockstep.EndOperationA());
    }
  });

  // Kick off thread for Operation B.
  Thread operation_b_thread([&lockstep, &a_count, &b_count]() {
    for (int i = 0; i < kNumIterations; i++) {
      EXPECT_OK(
          lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout));
      CHECK_EQ(a_count, b_count + 1);
      b_count++;
      EXPECT_OK(lockstep.EndOperationB());
    }
  });

  operation_a_thread.join();
  operation_b_thread.join();

  EXPECT_EQ(a_count, kNumIterations);
  EXPECT_EQ(b_count, kNumIterations);
}

TEST(LockstepTest, StartOperationABlockThenCancel) {
  Lockstep lockstep;
  // Bring the lockstep to the state, where a call to `StartOperationA` will
  // have to wait.
  EXPECT_OK(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.StartOperationBWithTimeout(kLockTimeout));

  // Kick off thread that sleeps then Cancel()s the lockstep.
  Thread cancel_thread([&lockstep]() {
    absl::SleepFor(kLockTimeout);
    lockstep.Cancel();
  });

  // Blocking until `Cancel` is called.
  EXPECT_THAT(lockstep.StartOperationAWithTimeout(/*timeout=*/kLongLockTimeout),
              StatusIs(absl::StatusCode::kAborted));

  cancel_thread.join();
}

TEST(LockstepTest, StartOperationBBlockThenCancel) {
  Lockstep lockstep;

  // Kick off thread that sleeps then Cancel()s the lockstep.
  Thread cancel_thread([&lockstep]() {
    absl::SleepFor(kLockTimeout);
    lockstep.Cancel();
  });

  // Blocking until `Cancel` is called.
  EXPECT_THAT(lockstep.StartOperationBWithTimeout(/*timeout=*/kLongLockTimeout),
              StatusIs(absl::StatusCode::kAborted));

  cancel_thread.join();
}

TEST(LockstepTest, EndOperationOkWhenCancelledDuringOperationA) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  lockstep.Cancel();
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.EndOperationB());
}

TEST(LockstepTest, EndOperationOkWhenCancelledDuringOperationB) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout));
  lockstep.Cancel();
  EXPECT_OK(lockstep.EndOperationB());
  EXPECT_OK(lockstep.EndOperationA());
}

TEST(LockstepTest, StartOperationASucceedsAfterReset) {
  Lockstep lockstep;
  lockstep.Cancel();
  // We cannot start A nor B.
  EXPECT_THAT(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout),
              StatusIs(absl::StatusCode::kAborted));
  EXPECT_THAT(lockstep.StartOperationBWithTimeout(kLockTimeout),
              StatusIs(absl::StatusCode::kAborted));

  EXPECT_THAT(lockstep.Reset(), icon::RealtimeIsOk());
  // After reset, we can start A.
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
}

TEST(LockstepTest, StartOperationBFailsAfterReset) {
  Lockstep lockstep;
  lockstep.Cancel();
  // We cannot start A nor B.
  EXPECT_THAT(lockstep.StartOperationAWithTimeout(kLockTimeout),
              StatusIs(absl::StatusCode::kAborted));
  EXPECT_THAT(lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout),
              StatusIs(absl::StatusCode::kAborted));

  EXPECT_THAT(lockstep.Reset(), icon::RealtimeIsOk());
  // After reset, we can still not start B.
  EXPECT_THAT(lockstep.StartOperationBWithTimeout(kLockTimeout),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
}

TEST(LockstepTest, StartOperationASucceedsAfterResetMultithread) {
  Lockstep lockstep;
  static constexpr int kCyclesUntilCancel = 100;
  absl::Notification operation_cancelled;
  absl::Notification lockstep_reset;

  // Kick off thread for operation a, that will eventually cancel and later
  // reset the lockstep.
  Thread operation_a_thread(
      [&lockstep, &operation_cancelled, &lockstep_reset]() {
        for (int i = 0; i < kCyclesUntilCancel; i++) {
          EXPECT_OK(
              lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout));
          EXPECT_OK(lockstep.EndOperationA());
        }
        // Cancel the lockstep. This will make operation b fail.
        lockstep.Cancel();
        // Wait for `operation_b_thread` to signal, that it noticed the
        // cancellation.
        EXPECT_TRUE(
            operation_cancelled.WaitForNotificationWithTimeout(kLockTimeout));
        // Reset the lockstep.
        EXPECT_OK(lockstep.Reset());
        lockstep_reset.Notify();
        // Run one last cycle after `Reset`.
        EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
        EXPECT_OK(lockstep.EndOperationA());
      });

  // Kick off thread for operation b.
  Thread operation_b_thread([&lockstep, &operation_cancelled,
                             &lockstep_reset]() {
    auto status = icon::OkStatus();
    // Run until the `operation_a_thread` cancels the lockstep.
    while (status =
               lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout),
           status.ok()) {
      EXPECT_OK(lockstep.EndOperationB());
    }
    EXPECT_THAT(status, StatusIs(absl::StatusCode::kAborted));

    // Tell `operation_a_thread` that we've cancelled.
    operation_cancelled.Notify();
    // Wait for `operation_a_thread` to `Reset` the lockstep.
    EXPECT_TRUE(lockstep_reset.WaitForNotificationWithTimeout(kLockTimeout));
    // Run one last cycle after `Reset`.
    EXPECT_OK(lockstep.StartOperationBWithTimeout(kLockTimeout));
    EXPECT_OK(lockstep.EndOperationB());
  });

  operation_a_thread.join();
  operation_b_thread.join();
}

TEST(LockstepTest, ResetFailsWhenNotCancelled) {
  Lockstep lockstep;
  EXPECT_THAT(lockstep.Reset(kLockTimeout),
              StatusIs(absl::StatusCode::kFailedPrecondition));
}

TEST(LockstepTest, StartOperationBTimesOutWhenOperationAIsRunning) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_THAT(lockstep.StartOperationBWithDeadline(absl::Now() + kLockTimeout),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.StartOperationBWithTimeout(kLockTimeout));
}

TEST(LockstepTest, StartOperationBTimesOutWithoutOperationA) {
  Lockstep lockstep;
  EXPECT_THAT(lockstep.StartOperationBWithTimeout(kLockTimeout),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
}

TEST(LockstepTest, StartOperationATimesOutWhenOperationBIsRunning) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_OK(lockstep.StartOperationBWithTimeout(kLockTimeout));
  EXPECT_THAT(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
}

TEST(LockstepTest, StartOperationATimesOutWithoutOperationB) {
  Lockstep lockstep;
  EXPECT_OK(lockstep.StartOperationAWithTimeout(kLockTimeout));
  EXPECT_OK(lockstep.EndOperationA());
  EXPECT_THAT(lockstep.StartOperationAWithDeadline(absl::Now() + kLockTimeout),
              StatusIs(absl::StatusCode::kDeadlineExceeded));
}

}  // namespace intrinsic
