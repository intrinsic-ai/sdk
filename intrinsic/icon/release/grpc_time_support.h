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

#ifndef INTRINSIC_ICON_RELEASE_GRPC_TIME_SUPPORT_H_
#define INTRINSIC_ICON_RELEASE_GRPC_TIME_SUPPORT_H_

#include "absl/time/time.h"
#include "grpc/support/log.h"
#include "grpc/support/time.h"
#include "grpcpp/support/time.h"

namespace grpc {

template <>
class TimePoint<absl::Time> {
 public:
  explicit TimePoint(absl::Time time) : time_(TimeToGprTimespec(time)) {}

  gpr_timespec raw_time() const { return time_; }

 private:
  static gpr_timespec TimeToGprTimespec(absl::Time time);

  const gpr_timespec time_;
};

// Converts gpr timespec to absl::Time
absl::Time TimeFromGprTimespec(gpr_timespec time);

// Converts absl::Time to gpr timespec with clock type set to REALTIME
gpr_timespec GprTimeSpecFromTime(absl::Time time);

// Converts gpr timespec to absl::Duration. The gpr timespec clock type
// should be TIMESPAN.
absl::Duration DurationFromGprTimespec(gpr_timespec time);

// Converts absl::Duration to gpr timespec with clock type set to TIMESPAM
gpr_timespec GprTimeSpecFromDuration(absl::Duration duration);

// Creates an absolute-time deadline from now + dur.
gpr_timespec DeadlineFromDuration(absl::Duration dur);

}  // namespace grpc

#endif  // INTRINSIC_ICON_RELEASE_GRPC_TIME_SUPPORT_H_
