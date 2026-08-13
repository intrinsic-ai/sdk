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

#ifndef INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_H_
#define INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_H_

#include <memory>
#include <string>

#include "absl/strings/string_view.h"

namespace intrinsic {

struct LivelinessSubscriptionData;

// Subscription to liveliness notifications.
class LivelinessSubscription {
 public:
  LivelinessSubscription();
  LivelinessSubscription(
      absl::string_view key_expression,
      std::unique_ptr<LivelinessSubscriptionData> subscription_data);

  ~LivelinessSubscription();

  LivelinessSubscription(const LivelinessSubscription&) = delete;
  LivelinessSubscription& operator=(const LivelinessSubscription&) = delete;
  LivelinessSubscription(LivelinessSubscription&&);
  LivelinessSubscription& operator=(LivelinessSubscription&&);

  absl::string_view KeyExpression() const { return key_expression_; }

  // To handle complex cases, such as when unsubscribing from a Python topic,
  // the shutdown sequence may need to be done in delicate ordering to avoid the
  // potential for deadlock. The Python GIL needs to be acquired when the
  // callback is deleted, but the GIL needs to be released when the actual
  // Zenoh subscription is destroyed, due to internal mutex contention in the
  // callback thread pool.  Exposing the Unsubscribe() function allows a pybind
  // holder type to do this in the correct order.
  void Unsubscribe();

 private:
  std::string key_expression_;
  std::unique_ptr<LivelinessSubscriptionData> subscription_data_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_H_
