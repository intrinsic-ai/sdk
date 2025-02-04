// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_REALTIME_CLOCK_H_
#define INTRINSIC_ICON_HAL_REALTIME_CLOCK_H_

#include <stdint.h>

#include <memory>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/control/realtime_clock_interface.h"
#include "intrinsic/icon/interprocess/shared_memory_lockstep/shared_memory_lockstep.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/icon/utils/clock.h"
#include "intrinsic/icon/utils/realtime_status.h"

namespace intrinsic::icon {

static constexpr absl::string_view kRealtimeClockLockstepInterfaceName =
    "realtime_clock_lockstep";
static constexpr absl::string_view kRealtimeClockUpdateInterfaceName =
    "realtime_clock_update";

// Payload for clock updates; gets stored in shared memory.
struct RealtimeClockUpdate {
  // Cycle start time in nanoseconds since the epoch.
  int64_t cycle_start_nanoseconds;
};

// RealtimeClock is an implementation of RealtimeClockInterface used by
// hardware modules to drive the realtime clock. It talks with the ICON server
// over shared memory.
class RealtimeClock : public RealtimeClockInterface {
 public:
  // Creates a RealtimeClock using memory segments with names specified by
  // `kRealtimeClockLockstepInterfaceName` and
  // `kRealtimeClockUpdateInterfaceName` by registering the respective segment
  // on `shm_manager`.
  static absl::StatusOr<std::unique_ptr<RealtimeClock>> Create(
      SharedMemoryManager& shm_manager);

  // This class is non-moveable and non-copyable to ensure that custom
  // destructor logic only ever runs once.
  RealtimeClock(const RealtimeClock& other) = delete;
  RealtimeClock& operator=(const RealtimeClock& other) = delete;

  // Signals to the ICON server that a real time cycle should begin. Blocks
  // until the cycle's update logic has completed; that is, blocks until
  // ApplyCommand has completed for all hardware modules. It is the caller's
  // responsibility to further wait until the next cycle's start time before
  // calling this again.
  //
  // The current_timestamp is considered the start time for the cycle.
  // Returns a deadline exceeded error in case of the deadline has expired.
  // Don't assume that the realtime cycle has been completed in case of such an
  // error. Use `Reset` to recover from such a situation!
  RealtimeStatus TickBlockingWithDeadline(intrinsic::Time current_timestamp,
                                          absl::Time deadline) override;

  // Resets the clock to its state after initialization, i.e. ready to call
  // TickBlockingWithTimeout.
  // Returns a deadline exceeded error on timeout.
  RealtimeStatus Reset(absl::Duration timeout) override;

  ~RealtimeClock() override;

 private:
  // The provided `lockstep` object synchronizes the callsite with the ICON
  // server's realtime update loop. The provided `realtime_clock_update`
  // communicates the cycle start time.
  RealtimeClock(
      SharedMemoryLockstep lockstep,
      ReadWriteMemorySegment<RealtimeClockUpdate> realtime_clock_update);

  SharedMemoryLockstep lockstep_;
  ReadWriteMemorySegment<RealtimeClockUpdate> update_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_REALTIME_CLOCK_H_
