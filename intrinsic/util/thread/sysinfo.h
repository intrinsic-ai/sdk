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

#ifndef INTRINSIC_UTIL_THREAD_SYSINFO_H_
#define INTRINSIC_UTIL_THREAD_SYSINFO_H_

#include <sys/syscall.h>
#include <unistd.h>

#include <thread>  // NOLINT

#include "absl/time/time.h"

namespace intrinsic {

// Number of logical processors.
inline int NumCPUs() { return std::thread::hardware_concurrency(); }

// Return the total cpu time used by the current thread.
absl::Duration ThreadCPUUsage();

// Return the thread id of the current thread, as told by the system.
inline pid_t GetTID() { return static_cast<pid_t>(syscall(SYS_gettid)); }

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_THREAD_SYSINFO_H_
