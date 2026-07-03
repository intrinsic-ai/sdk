// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/time/deadline_timeout.h"

#include "absl/time/clock.h"
#include "absl/time/time.h"

namespace intrinsic {

absl::Duration ToTimeout(absl::Time deadline) {
  return std::max(deadline - absl::Now(), absl::ZeroDuration());
}

absl::Time ToDeadline(absl::Duration timeout) { return absl::Now() + timeout; }

}  // namespace intrinsic
