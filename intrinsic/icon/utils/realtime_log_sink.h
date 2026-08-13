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

#ifndef INTRINSIC_ICON_UTILS_REALTIME_LOG_SINK_H_
#define INTRINSIC_ICON_UTILS_REALTIME_LOG_SINK_H_

#include <cstddef>

#include "intrinsic/icon/utils/log_sink.h"
#include "intrinsic/platform/common/buffers/realtime_write_queue.h"

namespace intrinsic::icon {

// A real-time safe log sink that writes to std::cerr.
// When there are multiple threads, each should create a thread-local object.
// Messages are buffered by a single, global non-RT thread and written out
// message by message.
class RealtimeLogSink : public LogSinkInterface {
 public:
  // Not RT safe.
  RealtimeLogSink();

  // Not RT safe.
  // Blocks until the log buffer has been written.
  ~RealtimeLogSink() override;

  // RT safe.
  // Not thread-safe, but concurrent use is allowed when each thread uses a
  // separate RealtimeLogSink object.
  // Messages reaching or exceeding kMessageMaxSize will be truncated.
  // If the buffer is full, messages may be dropped.
  void Log(const LogEntry& entry) override;

 private:
  RealtimeWriteQueue<LogEntry>::RtWriter* writer_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_REALTIME_LOG_SINK_H_
