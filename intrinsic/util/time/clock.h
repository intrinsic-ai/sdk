// Copyright 2023 Intrinsic Innovation LLC

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
