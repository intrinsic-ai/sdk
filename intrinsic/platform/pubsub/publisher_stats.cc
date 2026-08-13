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

#include "intrinsic/platform/pubsub/publisher_stats.h"

#include <cstdint>
#include <limits>

#include "absl/container/flat_hash_map.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"

namespace intrinsic::internal {

int PublisherStats::GetCount(absl::string_view topic) {
  absl::MutexLock lock(&mu_);
  return counts_[topic];
}

void PublisherStats::Increment(absl::string_view topic) {
  absl::MutexLock lock(&mu_);
  if (counts_[topic] == std::numeric_limits<int64_t>::max()) {
    counts_[topic] = 0;
  }
  counts_[topic]++;
}

void PublisherStats::Reset() {
  absl::MutexLock lock(&mu_);
  counts_.clear();
}

void ResetMessagesPublished() { return PublisherStats::Singleton().Reset(); }

int MessagesPublished(absl::string_view topic) {
  return PublisherStats::Singleton().GetCount(topic);
}

}  // namespace intrinsic::internal
