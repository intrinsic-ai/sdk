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

#ifndef INTRINSIC_UTIL_TIME_CLOCK_H_
#define INTRINSIC_UTIL_TIME_CLOCK_H_

#include "absl/time/time.h"

namespace intrinsic {

class ClockInterface {
 public:
  virtual ~ClockInterface() = default;
  virtual absl::Time Now() const = 0;
};

class RealClock : public ClockInterface {
 public:
  // Returns the global singleton instance of the real clock. This instance
  // is not owned by the caller and should not be deleted.
  static ClockInterface* GetInstance();

  RealClock() = default;
  ~RealClock() override = default;

  // Not copyable or movable.
  RealClock(const RealClock&) = delete;
  RealClock& operator=(const RealClock&) = delete;
  RealClock(RealClock&&) = delete;
  RealClock& operator=(RealClock&&) = delete;

  absl::Time Now() const override;
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_TIME_CLOCK_H_
