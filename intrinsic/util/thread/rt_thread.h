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

#ifndef INTRINSIC_UTIL_THREAD_RT_THREAD_H_
#define INTRINSIC_UTIL_THREAD_RT_THREAD_H_

#include <utility>

#include "absl/status/statusor.h"
#include "intrinsic/icon/utils/realtime_guard.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"
#include "intrinsic/util/thread/thread_utils.h"

namespace intrinsic {

// Creates a thread that is capable of running realtime code.
//
// The underlying thread is only truly an RT thread if the ThreadOptions
// specify it. For additional details on the threads which are created by this
// function, please see the ThreadOptions documentation.
// Note: When default constructed ThreadOptions are used consider to just use a
// plain intrinsic::Thread.
template <typename Function, typename... Args>
absl::StatusOr<Thread> CreateRealtimeCapableThread(const ThreadOptions& options,
                                                   Function&& f,
                                                   Args&&... args) {
  INTRINSIC_ASSERT_NON_REALTIME();
  return CreateThread(options, std::forward<Function>(f),
                      std::forward<Args>(args)...);
}

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_THREAD_RT_THREAD_H_
