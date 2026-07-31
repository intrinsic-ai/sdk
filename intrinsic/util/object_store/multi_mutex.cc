// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/object_store/multi_mutex.h"

#include <cstdint>
#include <memory>
#include <string>
#include <utility>

#include "absl/log/check.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"

namespace intrinsic {

MultiMutex::Lock::Lock(MultiMutex* owner, absl::string_view id)
    : owner_(owner), id_(id) {}

MultiMutex::Lock::~Lock() { owner_->Release(id_); }

MultiMutex::Lock MultiMutex::Acquire(absl::string_view id) {
  absl::Mutex* mutex;
  {
    MutexMap::accessor accessor;
    if (!mutexes_.find(accessor, std::string(id))) {
      mutexes_.insert(
          accessor,
          std::make_pair(std::string(id),
                         std::make_pair(0, std::make_unique<absl::Mutex>())));
    }
    int64_t& counter = accessor->second.first;
    mutex = accessor->second.second.get();
    ++counter;
  }
  // We grab the lock outside of the reservation to prevent deadlock. If another
  // thread currently holds the lock, it will need to enter the Reservation in
  // the Release method (see below) in order to unlock it. The counter variable
  // ensures that the mutex exists and is not destroyed prematurely.
  mutex->lock();
  return Lock(this, id);
}

void MultiMutex::Release(absl::string_view id) {
  MutexMap::accessor accessor;
  CHECK(mutexes_.find(accessor, std::string(id)))
      << "Release called without corresponding Acquire (id=" << id << ")";
  int64_t& counter = accessor->second.first;
  absl::Mutex* mutex = accessor->second.second.get();
  CHECK_GT(counter, 0) << "Release called with invalid counter variable. This "
                          "shouldn't be possible (id="
                       << id << ")";
  --counter;
  mutex->unlock();
  if (counter == 0) {
    mutexes_.erase(accessor);
  }
}

}  // namespace intrinsic
