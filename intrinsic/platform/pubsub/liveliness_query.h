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

#ifndef INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_QUERY_H_
#define INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_QUERY_H_

#include <functional>
#include <memory>

#include "absl/strings/string_view.h"
#include "intrinsic/platform/pubsub/zenoh_util/liveliness_get_context.h"

namespace intrinsic {

// Represents a query started by the `LivelinessGet` API.
class LivelinessQuery {
 public:
  LivelinessQuery(std::function<void(absl::string_view)> callback,
                  std::function<void(absl::string_view)> on_done)
      : context_(std::make_unique<LivelinessGetContext>(callback, on_done)) {}

  LivelinessGetContext* GetContext() { return context_.get(); }

 private:
  std::unique_ptr<LivelinessGetContext> context_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_QUERY_H_
