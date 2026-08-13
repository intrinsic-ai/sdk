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

#ifndef INTRINSIC_ICON_HAL_HARDWARE_MODULE_UTIL_H_
#define INTRINSIC_ICON_HAL_HARDWARE_MODULE_UTIL_H_

#include <cstddef>
#include <future>  // NOLINT(build/c++11)
#include <optional>
#include <string>

#include "absl/base/thread_annotations.h"
#include "absl/status/status.h"
#include "absl/strings/str_format.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state.fbs.h"

namespace intrinsic_fbs {
// The maximum length of a fault reason string. To be the same length as defined
// in the flatbuffer.
constexpr size_t kMaxFaultReasonLength = 256;

template <typename Sink>
void AbslStringify(Sink& sink, StateCode e) {
  absl::Format(&sink, "%s", EnumNameStateCode(e));
}

}  // namespace intrinsic_fbs

namespace intrinsic::icon {

enum class TransitionGuardResult { kNoOp, kAllowed, kProhibited };

// Returns whether the transition from `from` to `to` is allowed, prohibited or
// a no-op (e.g. from MotionEnabled to MotionEnabling).
TransitionGuardResult HardwareModuleTransitionGuard(
    intrinsic_fbs::StateCode from, intrinsic_fbs::StateCode to);

// Exit codes that a HWM process uses to indicate special results to
// the caller.
enum class HardwareModuleExitCode {
  // HWM shutdown normally.
  kNormalShutdown = 0,

  // HWM reported a fatal fault during initialization.
  kRestartRequested = 110,

  // HWM reported a fatal fault during initialization.
  kFatalFaultDuringInit = 111,

  // HWM reported a fatal fault during execution.
  kFatalFaultDuringExec = 112,
};

// A thin wrapper around a std::promise that allows us to set the value of the
// promise from multiple threads using a mutex.
template <typename T>
class SharedPromiseWrapper {
 public:
  // Sets the value of the promise. If the promise is already set, returns an
  // error. This function is thread-safe.
  absl::Status SetValue(const T& value) {
    absl::MutexLock lock(&mutex_);
    if (!promise_.has_value()) {
      return absl::FailedPreconditionError("Promise is already set.");
    }
    promise_->set_value(value);
    promise_.reset();
    return absl::OkStatus();
  }

  // Returns a shared future. This function can be called multiple times. Each
  // copy of the future can be used in different threads concurrently.
  std::shared_future<T> GetSharedFuture() { return shared_future_; }

  // Returns true if the promise has been set already.
  bool HasBeenSet() const {
    absl::MutexLock lock(&mutex_);
    return !promise_.has_value();
  }

 private:
  mutable absl::Mutex mutex_;
  // std::optional so that we can reset it after use to avoid exceptions.
  std::optional<std::promise<T>> promise_ ABSL_GUARDED_BY(mutex_) =
      std::promise<T>();
  std::shared_future<T> shared_future_ = promise_->get_future().share();
};

// Returns a string that can be used to visualize the state machine of the
// hardware module using graphviz in DOT format.
std::string CreateDotGraphvizStateMachineString();

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_HARDWARE_MODULE_UTIL_H_
