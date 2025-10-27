// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/time/clock.h"

#include "absl/time/clock.h"
#include "absl/time/time.h"

namespace intrinsic {

absl::Time RealClock::Now() const { return absl::Now(); }

ClockInterface* RealClock::GetInstance() {
  static ClockInterface* kInstance = new RealClock();
  return kInstance;
}

}  // namespace intrinsic
