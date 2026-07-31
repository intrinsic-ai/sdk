// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_OBJECT_STORE_MULTI_MUTEX_H_
#define INTRINSIC_UTIL_OBJECT_STORE_MULTI_MUTEX_H_

// Prevent conflicts with QT and concurrent hashmap.
// See macro definitions in qobjectdefs.h for details.
#include <cstdint>
#include <memory>
#include <string>
#include <utility>
#ifdef QOBJECTDEFS_H
#pragma push_macro("emit")
#undef emit
#pragma push_macro("slots")
#undef slots
#pragma push_macro("signals")
#undef signals
#endif

#include "absl/base/thread_annotations.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "tbb/concurrent_hash_map.h"

#ifdef QOBJECTDEFS_H
#pragma pop_macro("emit")
#pragma pop_macro("slots")
#pragma pop_macro("signals")
#endif

namespace intrinsic {

// The initial size of the concurrent hashmap.
constexpr int kMultiMutexInitialMapSize = 100;

// This class represents a dynamic collection of mutexes.
//
// It is used in order to synchronize access to a large collection of objects,
// where it does not make sense to permanently hold a mutex for each and every
// object, since the number of concurrent accesses to those objects is way
// smaller than the number of objects.
//
// The implementation is very simple and assumes that contention is rare, thus
// potentially waking up every blocked thread every time any mutex is acquired
// or released (which is not an issue if threads rarely block anyway).
//
// Example usage:
//
// ThreadUnsafeObject shared_object1;
// ThreadUnsafeObject shared_object2;
// MultiMutex mutexes;
//
// ... thread 1 ...
// {
//   Lock lock = mutexes.Acquire("1");
//   shared_object1.Use();
// }
//
//
// ... thread 2 ...
// {
//   Lock lock = mutexes.Acquire("1");
//   shared_object1.Use();
// }
//
// ... thread 3 ...
// {
//   Lock lock = mutexes.Acquire("2");
//   shared_object2.Use();
// }
//
// In the above example, threads 1/2 will each get exclusive access to
// shared_object1, possibly blocking each other, but never blocking, or being
// blocked by thread 3.
class ABSL_LOCKABLE MultiMutex {
 public:
  // An RAII object representing having an active mutex lock. Destroying this
  // object released the lock.
  class ABSL_SCOPED_LOCKABLE Lock {
   public:
    // Cannot be copied, can be move-constructed.
    Lock(const Lock&) = delete;
    Lock& operator=(const Lock&) = delete;
    Lock(Lock&&) = default;

    ~Lock();

   private:
    MultiMutex* const owner_;
    std::string id_;

    friend class MultiMutex;
    Lock(MultiMutex* owner, absl::string_view id);
  };

  using MutexMap = tbb::concurrent_hash_map<
      std::string, std::pair<int64_t, std::unique_ptr<absl::Mutex>>>;

  MultiMutex() : mutexes_(kMultiMutexInitialMapSize) {}

  // Acquire a lock of a named mutex given by id. This call may block if another
  // caller is currently holding a lock with the same id. Locks are NOT
  // recursive. Caller must make sure not to acquire a mutex with a given id if
  // that same id is already being locked by the same thread. It is OK for the
  // same thread to acquire two mutexes if their ids are different.
  Lock Acquire(absl::string_view id) ABSL_NO_THREAD_SAFETY_ANALYSIS;

 private:
  // A concurrent hashmap containing mutexes keyed by id. This data structure
  // allows us to avoid using a global lock around our object store. See
  // concurrent/containers/concurrent_hash_map.h for details.
  MutexMap mutexes_;

  void ABSL_UNLOCK_FUNCTION() Release(absl::string_view id);
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_OBJECT_STORE_MULTI_MUTEX_H_
