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

#ifndef INTRINSIC_PLATFORM_PUBSUB_PYTHON_GIL_AWARE_PUBSUB_H_
#define INTRINSIC_PLATFORM_PUBSUB_PYTHON_GIL_AWARE_PUBSUB_H_

#include "absl/strings/string_view.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

namespace intrinsic::pubsub {

// Subclass of PubSub that can be used in Python (via PyBind).
//
// GilAwarePubSub overrides certain helper methods of the base class to make
// sure that GIL is acquired and released when necessary to avoid deadlocks and
// memory management issues.
class GilAwarePubSub : public PubSub {
 public:
  GilAwarePubSub() : PubSub() {}
  explicit GilAwarePubSub(absl::string_view participant_name)
      : PubSub(participant_name) {}
  explicit GilAwarePubSub(absl::string_view participant_name,
                          absl::string_view config)
      : PubSub(participant_name, config) {}

 private:
  // Releases GIL while waiting for liveliness tokens.
  //
  // It unblocks other threads that need GIL. In particular, it unblocks threads
  // where callbacks are running. Unblocking those threads allows other features
  // that rely on callbacks (e.g. PubSub subscriptions) to run concurrently
  // with `LivelinessGetAllSynchronous`.
  void WaitForLivelinessTokens(std::function<void()> wait_fn) override;
};

}  // namespace intrinsic::pubsub

#endif  // INTRINSIC_PLATFORM_PUBSUB_PYTHON_GIL_AWARE_PUBSUB_H_
