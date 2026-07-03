// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_TIME_DEADLINE_TIMEOUT_H_
#define INTRINSIC_UTIL_TIME_DEADLINE_TIMEOUT_H_

#include "absl/time/time.h"

namespace intrinsic {

// Converts a deadline to a timeout.
absl::Duration ToTimeout(absl::Time deadline);

// Converts a timeout to a deadline.
absl::Time ToDeadline(absl::Duration timeout);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_TIME_DEADLINE_TIMEOUT_H_
