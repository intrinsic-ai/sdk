// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_QUEUE_MULTI_WRITER_H_
#define INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_QUEUE_MULTI_WRITER_H_

#include "absl/base/thread_annotations.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/icon/utils/realtime_guard.h"
#include "intrinsic/platform/common/buffers/rt_queue.h"

namespace intrinsic {

// Wrapper around RealtimeQueue<T>::Writer that makes it thread safe for
// multiple concurrent writers. In doing this, we drop the realtime safety of
// the write operation. Any readers of the RealtimeQueue are of course still
// realtime safe.
template <class T>
class RealtimeQueueMultiWriter {
 public:
  // Not copyable or moveable (Copying would undermine the thread safety, and
  // moving would make existing references to a MultiWriter invalid).
  RealtimeQueueMultiWriter(const RealtimeQueueMultiWriter&) = delete;
  RealtimeQueueMultiWriter(RealtimeQueueMultiWriter&&) = delete;

  // Wraps `writer` to provide thread-safe insertion into a RealtimeQueue.
  // `writer` must outlive this RealtimeQueueMultiWriter.
  explicit RealtimeQueueMultiWriter(typename RealtimeQueue<T>::Writer& writer)
      : writer_(writer) {}

  // Moves `value` into the underlying RealtimeQueue if possible.
  // Not realtime safe, but thread safe.
  // Returns ResourceExhausted if the underlying RealtimeQueue is full.
  absl::Status Insert(T&& value) ABSL_LOCKS_EXCLUDED(mutex_) {
    INTRINSIC_ASSERT_NON_REALTIME();
    absl::MutexLock l(&mutex_);
    // RealtimeQueue::Writer::Insert() makes a copy, negating any performance
    // benefits we get from accepting rvalues.
    T* queue_item = writer_.PrepareInsert();
    if (queue_item == nullptr) {
      return absl::ResourceExhaustedError("RealtimeQueue capacity exhausted");
    }
    *queue_item = std::move(value);
    writer_.FinishInsert();
    return absl::OkStatus();
  }

 private:
  absl::Mutex mutex_;
  typename RealtimeQueue<T>::Writer& writer_ ABSL_GUARDED_BY(mutex_);
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_COMMON_BUFFERS_RT_QUEUE_MULTI_WRITER_H_
