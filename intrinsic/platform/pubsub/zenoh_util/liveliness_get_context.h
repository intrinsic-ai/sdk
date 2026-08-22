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

#ifndef INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_LIVELINESS_GET_CONTEXT_H_
#define INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_LIVELINESS_GET_CONTEXT_H_

#include <functional>

#include "absl/strings/string_view.h"

namespace intrinsic {

// Context for the LivelinessGet API call.
//
// Stores application-level callbacks that are invoked as liveliness tokens
// are being fetched by LivelinessGet.
struct LivelinessGetContext {
  LivelinessGetContext(std::function<void(absl::string_view)> callback,
                       std::function<void(absl::string_view)> on_done)
      : callback_(callback), on_done_(on_done) {}

  std::function<void(absl::string_view)> callback_;
  std::function<void(absl::string_view)> on_done_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_LIVELINESS_GET_CONTEXT_H_
