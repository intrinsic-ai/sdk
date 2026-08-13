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

#include "intrinsic/util/thread/thread_options.h"

#include <sched.h>

#include <string>
#include <vector>

#include "absl/strings/string_view.h"

namespace intrinsic {

// The settings are platform-dependent on Linux.
ThreadOptions& ThreadOptions::SetRealtimeHighPriorityAndScheduler() {
  priority_ = kHighRealtimePriority;
  policy_ = SCHED_FIFO;
  return *this;
}

ThreadOptions& ThreadOptions::SetRealtimeLowPriorityAndScheduler() {
  priority_ = kLowRealtimePriority;
  policy_ = SCHED_FIFO;
  return *this;
}

ThreadOptions& ThreadOptions::SetNormalPriorityAndScheduler() {
  priority_ = 0;
  policy_ = SCHED_OTHER;
  return *this;
}

ThreadOptions& ThreadOptions::SetPriority(int priority) {
  priority_ = priority;
  return *this;
}

ThreadOptions& ThreadOptions::SetSchedulePolicy(int policy) {
  policy_ = policy;
  return *this;
}

ThreadOptions& ThreadOptions::SetAffinity(const std::vector<int>& cpus) {
  cpus_ = cpus;
  return *this;
}

ThreadOptions& ThreadOptions::SetName(absl::string_view name) {
  name_ = std::string(name);
  return *this;
}

}  // namespace intrinsic
