// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_COMMON_BUFFERS_REALTIME_WRITE_QUEUE_H_
#define INTRINSIC_PLATFORM_COMMON_BUFFERS_REALTIME_WRITE_QUEUE_H_

#include <cstddef>
#include <cstdint>
#include <cstring>

#include "absl/base/attributes.h"
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "intrinsic/icon/interprocess/binary_futex.h"
#include "intrinsic/icon/utils/realtime_status.h"
#include "intrinsic/platform/common/buffers/rt_queue_buffer.h"

namespace intrinsic {

// Possible results of a read.
enum class ReadResult {
  // An item was consumed.
  kConsumed,
  // The writer has closed the queue and all items have been consumed.
  kClosed,
  // The deadline was reached.
  kDeadlineExceeded
};

// Implements a single-producer single-consumer thread-safe queue with real-time
// safe non-blocking writes and non-real-time safe blocking reads. In
// particular, concurrent writes are unsafe, as are concurrent reads. Reads may
// be concurrent to writes.
//
// Example usage:
//
// RealtimeWriteQueue<int> queue;
// Thread writer([&queue]() {
//   for (int num = 0; num < 10; num++) {
//     queue.Writer().Write(num)
//   }
//   queue.Writer().Close();
// });
// Thread reader([&queue]() {
//   int num = 0;
//   while (queue.Reader().Read(num) == ReadResult::kConsumed) {
//   }
// });
// reader.join();
// writer.join();
template <typename T>
class RealtimeWriteQueue {
 public:
  class NonRtReader {
   public:
    // Blocks until either
    // (a) the next item has been read into *item; returns kConsumed.
    // (b) Close() has been called on the corresponding writer *and* all values
    //     have been consumed; returns kClosed.
    // (c) The deadline is reached; returns kDeadlineExceeded.
    // If the deadline is in the past, this function will not block, but it will
    // return kConsumed or kClosed if possible.
    ABSL_MUST_USE_RESULT ReadResult ReadWithTimeout(T& item,
                                                    absl::Time deadline);

    ABSL_MUST_USE_RESULT ReadResult Read(T& item) {
      return ReadWithTimeout(item, absl::InfiniteFuture());
    }

    // Returns true when the buffer is empty.
    bool Empty() const { return buffer_.Empty(); }

   private:
    friend RealtimeWriteQueue;
    explicit NonRtReader(internal::RtQueueBuffer<T>& buffer,
                         icon::BinaryFutex& notification);

    internal::RtQueueBuffer<T>& buffer_;
    icon::BinaryFutex& notification_;

    bool closed_ = false;
    uint64_t count_available_ = 0;
  };

  class RtWriter {
   public:
    // Returns true if the write succeeded. A write can fail if the queue is
    // full. It is invalid to call Write() after calling Close().
    ABSL_MUST_USE_RESULT bool Write(const T& item);

    // Marks the queue as 'closed', further attempts to Write() to the queue
    // are invalid.
    void Close();

    // Returns true if the queue is 'closed'.
    bool Closed() const;

   private:
    friend RealtimeWriteQueue;
    explicit RtWriter(internal::RtQueueBuffer<T>& buffer,
                      icon::BinaryFutex& notification);
    internal::RtQueueBuffer<T>& buffer_;
    icon::BinaryFutex& notification_;

    bool closed_ = false;
  };

  static constexpr size_t kDefaultBufferCapacity = 100;

  RealtimeWriteQueue();
  explicit RealtimeWriteQueue(size_t capacity);

  NonRtReader& Reader() { return reader_; }
  RtWriter& Writer() { return writer_; }

 private:
  internal::RtQueueBuffer<T> buffer_;
  icon::BinaryFutex notification_ =
      icon::BinaryFutex(/*posted=*/false, /*private_futex=*/true);

  NonRtReader reader_;
  RtWriter writer_;
};

template <typename T>
ReadResult RealtimeWriteQueue<T>::NonRtReader::ReadWithTimeout(
    T& item, absl::Time deadline) {
  if (count_available_ == 0 && closed_) {
    return ReadResult::kClosed;
  }

  if (count_available_ == 0) {
    // Block until items are ready.
    icon::RealtimeStatus status = notification_.WaitUntil(deadline);
    if (status.code() == absl::StatusCode::kAborted) {
      closed_ = true;
    }
    // Only read new items after WaitUntil, because WaitUntil resets to
    // BinaryFutex::kReady and tells the writer to resume notifications.
    count_available_ = buffer_.Size();
  }

  if (count_available_ == 0) {
    if (closed_) {
      return ReadResult::kClosed;
    }
    // Treat other errors as temporary.
    return ReadResult::kDeadlineExceeded;
  }

  T* element = buffer_.Front();
  // This is guaranteed to never happen, since we check that
  // `count_available_` is not zero before calling `buffer_.Front()`, and during
  // Writer::Write() the count is incremented after inserting to the queue.
  CHECK(element != nullptr) << "Attempted to read when no count_available_";
  item = *element;
  buffer_.DropFront();

  count_available_--;
  return ReadResult::kConsumed;
}

template <typename T>
RealtimeWriteQueue<T>::NonRtReader::NonRtReader(
    internal::RtQueueBuffer<T>& buffer, icon::BinaryFutex& notification)
    : buffer_(buffer), notification_(notification) {}

template <typename T>
bool RealtimeWriteQueue<T>::RtWriter::Write(const T& item) {
  CHECK(!closed_) << "Invalid to Write() after Close()ing the queue";
  T* element = buffer_.PrepareInsert();
  if (element == nullptr) {
    return false;
  }
  *element = item;
  buffer_.FinishInsert();
  return notification_.Post().ok();
}

template <typename T>
void RealtimeWriteQueue<T>::RtWriter::Close() {
  closed_ = true;
  notification_.Close();
}

template <typename T>
bool RealtimeWriteQueue<T>::RtWriter::Closed() const {
  return closed_;
}

template <typename T>
RealtimeWriteQueue<T>::RtWriter::RtWriter(internal::RtQueueBuffer<T>& buffer,
                                          icon::BinaryFutex& notification)
    : buffer_(buffer), notification_(notification) {}

template <typename T>
RealtimeWriteQueue<T>::RealtimeWriteQueue()
    : buffer_(kDefaultBufferCapacity),
      reader_(buffer_, notification_),
      writer_(buffer_, notification_) {}

template <typename T>
RealtimeWriteQueue<T>::RealtimeWriteQueue(size_t capacity)
    : buffer_(capacity),
      reader_(buffer_, notification_),
      writer_(buffer_, notification_) {}

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_COMMON_BUFFERS_REALTIME_WRITE_QUEUE_H_
