// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_LOCKSTEP_SHARED_MEMORY_LOCKSTEP_H_
#define INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_LOCKSTEP_SHARED_MEMORY_LOCKSTEP_H_

#include <utility>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/util/thread/lockstep.h"

namespace intrinsic::icon {

// SharedMemoryLockstep is a Lockstep synchronization primitive that is stored
// in shared memory. This can be used for synchronization across process
// boundaries.
class SharedMemoryLockstep {
 public:
  // Null SharedMemoryLockstep. Dereferencing this will check-fail. This allows
  // value semantics / move to work.
  SharedMemoryLockstep() : memory_segment_(), lockstep_(nullptr) {}

  // Creates a SharedMemoryLockstep from a Lockstep memory segment. Prefer to
  // use CreateSharedMemoryLockstep or GetSharedMemoryLockstep instead.
  explicit SharedMemoryLockstep(ReadWriteMemorySegment<Lockstep>&& segment)
      : memory_segment_(std::move(segment)),
        lockstep_(&memory_segment_.GetValue()) {}

  // Returns true if the lockstep is attached to two instances.
  bool Connected() const;

  // Obtains the underlying shared memory Lockstep object that can be used for
  // synchronization. Returns nullptr if this is null (default-constructed).
  Lockstep* GetLockstep() { return lockstep_; }

  // Dereferencing returns the underlying Lockstep object. Check-fails if this
  // is null (default-constructed).
  Lockstep* operator*() {
    CHECK(lockstep_ != nullptr) << "null SharedMemoryLockstep dereferenced";
    return lockstep_;
  }
  const Lockstep* operator*() const {
    CHECK(lockstep_ != nullptr) << "null SharedMemoryLockstep dereferenced";
    return lockstep_;
  }
  Lockstep* operator->() {
    CHECK(lockstep_ != nullptr) << "null SharedMemoryLockstep dereferenced";
    return lockstep_;
  }
  const Lockstep* operator->() const {
    CHECK(lockstep_ != nullptr) << "null SharedMemoryLockstep dereferenced";
    return lockstep_;
  }

 private:
  // Hold onto the memory segment, since it is refcounted.
  ReadWriteMemorySegment<Lockstep> memory_segment_;
  // Raw pointer into the memory segment, for convenience.
  Lockstep* lockstep_;
};

// Creates a SharedMemoryLockstep whose shared memory is managed by `manager`
// and is stored in a segment named `memory_name`. The `manager` must outlive
// the returned SharedMemoryLockstep.
absl::StatusOr<SharedMemoryLockstep> CreateSharedMemoryLockstep(
    SharedMemoryManager& manager, absl::string_view memory_name);

// Convenience method for tests.
// Returns a SharedMemoryLockstep in a shared memory segment named
// `memory_name` that was created on `manager`.
// The `manager` must outlive the returned SharedMemoryLockstep.
absl::StatusOr<SharedMemoryLockstep> GetSharedMemoryLockstep(
    const SharedMemoryManager& manager, absl::string_view memory_name);

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_INTERPROCESS_SHARED_MEMORY_LOCKSTEP_SHARED_MEMORY_LOCKSTEP_H_
