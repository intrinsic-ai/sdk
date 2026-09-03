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

#ifndef INTRINSIC_UTIL_THREAD_UTIL_H_
#define INTRINSIC_UTIL_THREAD_UTIL_H_

#include <sched.h>
#include <unistd.h>

#include <string>

#include "absl/base/attributes.h"
#include "absl/container/flat_hash_set.h"
#include "absl/functional/any_invocable.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/time.h"

namespace intrinsic {

// Reads the CPU affinity set from the `rcu_nocbs` parameter in `/proc/cmdline`.
// The parameter `path_for_testing` can be used to define a custom `cmdline` for
// testing.
// Returns FailedPreconditionError if the path is not valid,`rcu_nocbs` is not
// defined or on parsing errors.
absl::StatusOr<absl::flat_hash_set<int>> ReadCpuAffinitySetFromCommandLine(
    absl::string_view path_for_testing = "/proc/cmdline");

// Like Notification::WaitForNotification, but the user also provides a function
// that is polled periodically to determine whether to quit waiting. If the
// function returns false, we stop waiting and return the current value of
// notification.HasBeenNotified(). If it returns true, we keep waiting.
bool WaitForNotificationWithInterrupt(
    absl::Notification& notification, absl::AnyInvocable<bool()> should_quit,
    absl::Duration poll_interval = absl::Milliseconds(100));

// Like Notification::WaitForNotificationWithDeadline, but the user also
// provides a function that is polled periodically to determine whether to quit
// waiting. If the function returns false, we stop waiting and return the
// current value of notification.HasBeenNotified(). If it returns true, we keep
// waiting.
bool WaitForNotificationWithDeadlineAndInterrupt(
    absl::Notification& notification, absl::Time deadline,
    absl::AnyInvocable<bool()> should_quit,
    absl::Duration poll_interval = absl::Milliseconds(100));

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_THREAD_UTIL_H_
