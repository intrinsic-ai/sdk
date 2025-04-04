// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_ASYNC_BUFFER_H_
#define INTRINSIC_ICON_UTILS_ASYNC_BUFFER_H_

#include <atomic>
#include <cstdint>
#include <memory>
#include <utility>

#include "absl/log/check.h"

namespace intrinsic {

// A real-time safe producer/consumer buffer container.
// Writes by one producer thread and reads by one consumer thread may be
// concurrent, anything else is not thread-safe.
// Both reads and writes are real-time safe, and do not copy the data.
//
// AsyncBuffer implements a triple buffered single producer/single consumer
// container.
//
// Terminology:
// \li Active Buffer: the buffer that is being read by the consumer
// \li Free Buffer: a buffer ready to be updated by the producer
// \li Mailbox Buffer: a buffer not currently owned by either the producer or
// the consumer.  The mailbox buffer may be "full" (because it has recently been
// committed by the producer) or "empty" (because it was recently fetched by the
// consumer)
//
// A producer updates a buffer by:
// \code
// Buffer* free_buff = async.GetFreeBuffer();
// free_buff->update();
// async.CommitFreeBuffer();
// \endcode
//
// A consumer can get the most up to date buffer by:
// \code
// Buffer* active = async.GetActiveBuffer();
// \endcode
//
// The act of committing the free buffer atomically swaps the free buffer with
// the mailbox buffer.  After a commit operation, the mailbox buffer is
// considered to be "full" and will be the buffer fetched by the consumer the
// next time GetActiveBuffer() is called.
//
// The act of getting the active buffer either just gets the current active
// buffer if the mailbox is empty, or atomically swaps the current active buffer
// with the mailbox buffer, and returns the new active buffer if the mailbox is
// full.  In either case, the mailbox is considered to be empty after a call to
// GetActiveBuffer.
template <typename T>
class AsyncBuffer {
 public:
  // Creates an AsyncBuffer with internally allocated storage.
  template <typename... InitArgs>
  explicit AsyncBuffer(const InitArgs&... args);

  // Returns the active buffer to the consumer. This buffer is guaranteed not
  // to be modified until the next call to GetActiveBuffer() by the consumer.
  // This may be the same buffer which was returned to last call to
  // GetActiveBuffer().
  //
  // Returns true if the mailbox buffer was full and the returned pointer
  // points to it. Returns false if the mailbox buffer was empty and the
  // returned pointer points to the buffer that was already active at the time
  // of this call.
  //
  // If nothing has been committed yet, *buffer points to a default constructed
  // (potentially uninitialized) `T` afterwards.
  bool GetActiveBuffer(T** buffer);

  // Commits the free buffer by swapping it with the mailbox buffer.
  //
  // Each call to this function must be preceded by a call to GetFreeBuffer()
  // or else it has no effect. Returns true if this call was preceded by a call
  // to GetFreeBuffer(), false otherwise.
  bool CommitFreeBuffer();

  // Returns a free buffer to the producer. This buffer can be modified until
  // the producer calls CommitFreeBuffer().
  T* GetFreeBuffer();

 private:
  // Indicates which is the active/mailbox/free buffer in `buffers_`.
  struct State {
    bool IsConsistent() const {
      return active_index < 3 && mailbox_index < 3 && free_index < 3 &&
             (active_index != mailbox_index) && (mailbox_index != free_index) &&
             (free_index != active_index);
    }
    uint8_t active_index = 0;
    uint8_t mailbox_index = 1;
    uint8_t free_index = 2;
    bool mailbox_full = false;
  };

  std::atomic<State> state_;
  static_assert(decltype(state_)::is_always_lock_free);
  bool free_buffer_checked_out_ = false;
  // Buffers: active, mailbox, free.
  std::unique_ptr<T> buffers_[3];
};

template <typename T>
template <typename... InitArgs>
inline AsyncBuffer<T>::AsyncBuffer(const InitArgs&... args) {
  for (std::unique_ptr<T>& buf : buffers_) {
    buf = std::make_unique<T>(args...);
  }
}

// Implicitly, the AsyncBuffer `state_` is a state machine.  It manages 3
// buffers which, at any instant in time, are distributed across 3 different
// slots (Active, Mailbox and Free).
//
// There are 3! (6) different orders for the buffers to exist in the slots.  In
// addition, the mailbox slot can be considered to be either "full" or "empty".
// This makes the total number of states for the system 2 * 3! = 12.
//
// Three operations need to be supported by the AsyncBuffer.
// \li 1) Getting the free buffer.
// \li 2) Committing the free buffer.
// \li 3) Getting the active buffer.
//
// Operation #1 (getting the free buffer) does not permute the state of the
// system, it only needs to look up which buffer is the free buffer based on
// the state of the system.
//
// Operation #2 (committing the free buffer) always permutes the state of the
// system.  It always swaps the buffers in the Free and Mailbox slots, and it
// always causes the Mailbox to become full.
//
// Operation #3 (fetching the active buffer) permutes the state of the system if
// the mailbox is full (swapping Active and Mailbox and returning the new
// Active), but not if the mailbox is empty.

template <typename T>
inline T* AsyncBuffer<T>::GetFreeBuffer() {
  free_buffer_checked_out_ = true;
  return buffers_[state_.load().free_index].get();
}

template <typename T>
inline bool AsyncBuffer<T>::GetActiveBuffer(T** buffer) {
  State current;
  State next;
  do {
    current = state_.load();
    // If empty, nothing changes. If full, it gets empty and the active buffer
    // swaps with the mailbox. (Never touches the free buffer. Mailbox is empty
    // afterwards.)
    next = current;
    if (next.mailbox_full) {
      std::swap(next.active_index, next.mailbox_index);
      next.mailbox_full = false;
    }
  } while (!state_.compare_exchange_strong(current, next));
  DCHECK_EQ(current.free_index, next.free_index);
  DCHECK_EQ(next.mailbox_full, false);
  DCHECK(next.IsConsistent());

  *buffer = buffers_[next.active_index].get();
  return current.mailbox_full;
}

template <typename T>
inline bool AsyncBuffer<T>::CommitFreeBuffer() {
  if (!free_buffer_checked_out_) {
    return false;
  }

  State current;
  State next;
  do {
    current = state_.load();
    // Same behavior whether full or empty, mailbox is always full afterwards.
    // Free buffer swaps with mailbox. (Never touches the active buffer. Mailbox
    // is full afterwards.)
    next = {.active_index = current.active_index,
            .mailbox_index = current.free_index,
            .free_index = current.mailbox_index,
            .mailbox_full = true};
  } while (!state_.compare_exchange_strong(current, next));
  DCHECK_EQ(current.active_index, next.active_index);
  DCHECK_EQ(next.mailbox_full, true);
  DCHECK(next.IsConsistent());
  free_buffer_checked_out_ = false;
  return true;
}

}  // namespace intrinsic

#endif  // INTRINSIC_ICON_UTILS_ASYNC_BUFFER_H_
