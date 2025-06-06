// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/thread/util.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <functional>
#include <string>

#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/util/testing/gtest_wrapper.h"
#include "intrinsic/util/thread/thread.h"
#include "ortools/base/helpers.h"
#include "ortools/base/options.h"
#include "ortools/base/path.h"

namespace intrinsic {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::testing::ContainerEq;
using ::testing::HasSubstr;
using ::testing::TempDir;

// rcu_nocbs contains a single CPU
constexpr absl::string_view kSingleCPUCommandline =
    R"(init=/usr/lib/systemd/systemd boot=local rootwait ro noresume loglevel=7 
    console=tty1 console=ttyS0,115200 apparmor=0 virtio_net.napi_tx=1 
    systemd.unified_cgroup_hierarchy=true csm.disabled=1 
    loadpin.exclude=kernel-module modules-load=loadpin_trigger 
    module.sig_enforce=1 i915.modeset=1 efi=runtime processor.max_cstate=0 
    idle=poll isolcpus=5 nohz=on nohz_full=5 rcu_nocbs=5 rcu_nocb_poll 
    nowatchdog pcie_aspm=off   dm_verity.error_behavior=3 dm_verity.max_bios=-1 
    dm_verity.dev_wait=1 root=/dev/dm-0 dm="1 vroot none ro 1,0 4077568 verity 
    payload=PARTLABEL=IROOT-B hashtree=PARTLABEL=IROOT-B hashstart=4077568 
    alg=sha256 
    root_hexdigest=a578c32661d6e56b94f8dd6352c2ff64d07ba76b272c306cbd745b28f0bb1870 
    salt=1536f7fcc5721837477e891e54a99ff4d0f466830b74c7a7e10f824da62118e5")";

// rcu_nocbs contains the list of CPUs
constexpr absl::string_view kMultiCPUCommandline =
    R"(init=/usr/lib/systemd/systemd boot=local rootwait ro noresume loglevel=7 
    console=tty1 console=ttyS0,115200 apparmor=0 virtio_net.napi_tx=1 
    systemd.unified_cgroup_hierarchy=true csm.disabled=1 
    loadpin.exclude=kernel-module modules-load=loadpin_trigger 
    module.sig_enforce=1 i915.modeset=1 efi=runtime processor.max_cstate=0 
    idle=poll isolcpus=5 nohz=on nohz_full=5 rcu_nocbs=0-2,7,12-14,19-18 
    rcu_nocb_poll nowatchdog pcie_aspm=off   dm_verity.error_behavior=3 
    dm_verity.max_bios=-1 dm_verity.dev_wait=1 root=/dev/dm-0 dm="1 vroot none 
    ro 1,0 4077568 verity payload=PARTLABEL=IROOT-B hashtree=PARTLABEL=IROOT-B 
    hashstart=4077568 alg=sha256 
    root_hexdigest=a578c32661d6e56b94f8dd6352c2ff64d07ba76b272c306cbd745b28f0bb1870 
    salt=1536f7fcc5721837477e891e54a99ff4d0f466830b74c7a7e10f824da62118e5")";

TEST(ReadCpuAffinitySetFromCommandLine, FailsForInvalidPath) {
  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine("tmp/IdoNotExistHopefully5234"),
              StatusIs(absl::StatusCode::kNotFound));
}

TEST(ReadCpuAffinitySetFromCommandLine, WorksSingle) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, kSingleCPUCommandline, file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              IsOkAndHolds(ContainerEq(absl::flat_hash_set<int>{5})));
}

TEST(ReadCpuAffinitySetFromCommandLine, WorksGroup) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, kMultiCPUCommandline, file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              IsOkAndHolds(ContainerEq(
                  absl::flat_hash_set<int>{0, 1, 2, 7, 12, 13, 14, 18, 19})));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsWithoutRcuNocbs) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "", file::Defaults()));

  EXPECT_THAT(
      ReadCpuAffinitySetFromCommandLine(path),
      StatusIs(absl::StatusCode::kFailedPrecondition, HasSubstr("rcu_nocbs")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForDuplicatesWithRange) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=1,1-2", file::Defaults()));

  EXPECT_THAT(
      ReadCpuAffinitySetFromCommandLine(path),
      StatusIs(absl::StatusCode::kFailedPrecondition, HasSubstr("Duplicate")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForDuplicateSingleEntries) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=1,1,1,2", file::Defaults()));

  EXPECT_THAT(
      ReadCpuAffinitySetFromCommandLine(path),
      StatusIs(absl::StatusCode::kFailedPrecondition, HasSubstr("Duplicate")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForInvalidRangeFormat) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=1--2", file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              StatusIs(absl::StatusCode::kFailedPrecondition,
                       HasSubstr("Expected Format")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForInvalidRangeFormat2) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=1-", file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              StatusIs(absl::StatusCode::kFailedPrecondition,
                       HasSubstr("Expected Format")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForNegativeRangeFormat) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=-1-2", file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              StatusIs(absl::StatusCode::kFailedPrecondition,
                       HasSubstr("Failed to parse")));
}

TEST(ReadCpuAffinitySetFromCommandLine, FailsForNegativeCPU) {
  const std::string path = file::JoinPath(TempDir(), "cmdline");
  ASSERT_OK(file::SetContents(path, "rcu_nocbs=-1", file::Defaults()));

  EXPECT_THAT(ReadCpuAffinitySetFromCommandLine(path),
              StatusIs(absl::StatusCode::kFailedPrecondition,
                       HasSubstr("Expected Format")));
}

class WaitForNotificationWithInterruptTest : public ::testing::Test {
 public:
  void SetUp() override {
    should_quit_ = false;
    timed_out_ = false;
    was_notified_ = false;
  }

  void StartWaiting(absl::Notification& notification) {
    thread_ = Thread([this, &notification]() {
      absl::Time start_time = absl::Now();
      this->was_notified_ = WaitForNotificationWithInterrupt(
          notification, std::function<bool()>([this, &start_time]() {
            timed_out_ = absl::Now() >= start_time + timeout_;
            return this->should_quit_ || timed_out_;
          }));
    });
  }

  Thread thread_;
  absl::Duration timeout_ = absl::Seconds(60);
  bool timed_out_;
  bool should_quit_;
  bool was_notified_;
};

TEST_F(WaitForNotificationWithInterruptTest, ReturnsWhenNotified) {
  absl::Notification notification;

  StartWaiting(notification);
  notification.Notify();
  thread_.join();

  EXPECT_TRUE(was_notified_);
  EXPECT_FALSE(timed_out_);
}

TEST_F(WaitForNotificationWithInterruptTest, CanBeInterrupted) {
  absl::Notification notification;

  StartWaiting(notification);
  should_quit_ = true;
  thread_.join();

  EXPECT_FALSE(was_notified_);
  EXPECT_FALSE(timed_out_);
}

class WaitForNotificationWithDeadlineAndInterruptTest : public ::testing::Test {
 public:
  void SetUp() override {
    should_quit_ = false;
    was_notified_ = false;
  }

  void StartWaiting(absl::Notification& notification, absl::Duration timeout) {
    thread_ = Thread([this, &notification, timeout]() {
      absl::Time start_time = absl::Now();
      this->was_notified_ = WaitForNotificationWithDeadlineAndInterrupt(
          notification, absl::Now() + timeout,
          std::function<bool()>([this]() { return this->should_quit_; }));
      duration_ = absl::Now() - start_time;
    });
  }

  Thread thread_;
  bool should_quit_;
  bool was_notified_;
  absl::Duration duration_;
};

TEST_F(WaitForNotificationWithDeadlineAndInterruptTest, ReturnsWhenNotified) {
  absl::Notification notification;
  absl::Duration timeout = absl::Seconds(60);

  StartWaiting(notification, timeout);
  notification.Notify();
  thread_.join();

  EXPECT_TRUE(was_notified_);
  // Just verify that the function returned before the deadline was reached.
  EXPECT_LT(duration_, timeout / 6);
}

TEST_F(WaitForNotificationWithDeadlineAndInterruptTest, TimesOut) {
  absl::Notification notification;
  absl::Duration timeout = absl::Milliseconds(100);

  StartWaiting(notification, timeout);
  thread_.join();

  EXPECT_FALSE(was_notified_);
  EXPECT_GE(duration_, timeout);
}

TEST_F(WaitForNotificationWithDeadlineAndInterruptTest, CanBeInterrupted) {
  absl::Notification notification;
  absl::Duration timeout = absl::Seconds(60);

  StartWaiting(notification, timeout);
  should_quit_ = true;
  thread_.join();

  EXPECT_FALSE(was_notified_);
  // Just verify that the function returned before the deadline was reached.
  EXPECT_LT(duration_, timeout / 6);
}

}  // namespace
}  // namespace intrinsic
