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

#ifndef INTRINSIC_UTIL_THREAD_LOCKABLE_MUTEX_H_
#define INTRINSIC_UTIL_THREAD_LOCKABLE_MUTEX_H_

#include "absl/base/attributes.h"
#include "absl/base/thread_annotations.h"
#include "absl/synchronization/mutex.h"

namespace intrinsic {

// This class wraps a absl::Mutex to make it compatible with the BasicLockable
// and Lockable named requirements (C++11 and later). This way, absl::Mutex can
// be used with std::lock and std::lock_guard. Especially for multi-lock
// acquisition the deadlock-free std::lock function is valuable.

class LockableMutex {
 public:
  explicit LockableMutex(absl::Mutex& mutex ABSL_ATTRIBUTE_LIFETIME_BOUND)
      : mutex_(mutex) {}

  LockableMutex(const LockableMutex&) = delete;
  LockableMutex& operator=(const LockableMutex&) = delete;

  void lock() ABSL_EXCLUSIVE_LOCK_FUNCTION(mutex_) { mutex_.Lock(); }
  void unlock() ABSL_UNLOCK_FUNCTION(mutex_) { mutex_.Unlock(); }
  bool try_lock() ABSL_EXCLUSIVE_TRYLOCK_FUNCTION(true, mutex_) {
    return mutex_.TryLock();
  }

 private:
  absl::Mutex& mutex_;
};

}  // namespace intrinsic
#endif  // INTRINSIC_UTIL_THREAD_LOCKABLE_MUTEX_H_
