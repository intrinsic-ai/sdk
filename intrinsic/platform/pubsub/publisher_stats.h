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

#ifndef INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_STATS_H_
#define INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_STATS_H_

#include <string>

#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/debugging/leak_check.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"

namespace intrinsic::internal {

class PublisherStats {
 public:
  int GetCount(absl::string_view topic);

  void Increment(absl::string_view topic);

  void Reset();

  static PublisherStats& Singleton() {
    static PublisherStats* stats = absl::IgnoreLeak(new PublisherStats);
    return *stats;
  }

 private:
  absl::Mutex mu_;
  absl::flat_hash_map<std::string, int64_t> counts_ ABSL_GUARDED_BY(mu_);
};

// How many pubsub messages have been sent by this process for a given topic.
int MessagesPublished(absl::string_view topic);
// Reset the MessagePublished counters to 0.
void ResetMessagesPublished();

}  // namespace intrinsic::internal

#endif  // INTRINSIC_PLATFORM_PUBSUB_PUBLISHER_STATS_H_
