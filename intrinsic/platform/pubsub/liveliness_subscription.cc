// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/liveliness_subscription.h"

#include <memory>
#include <string>
#include <utility>

#include "absl/strings/string_view.h"
#include "intrinsic/platform/pubsub/liveliness_subscription_data.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

namespace intrinsic {

LivelinessSubscription::LivelinessSubscription() = default;

LivelinessSubscription::LivelinessSubscription(
    absl::string_view key_expression,
    std::unique_ptr<LivelinessSubscriptionData> subscription_data)
    : key_expression_(key_expression),
      subscription_data_(std::move(subscription_data)) {}

LivelinessSubscription::LivelinessSubscription(LivelinessSubscription&&) =
    default;

LivelinessSubscription& LivelinessSubscription::operator=(
    LivelinessSubscription&& other) {
  if (!key_expression_.empty()) {
    Zenoh().imw_destroy_subscription(
        subscription_data_->key_expression.c_str(), zenoh_static_callback,
        subscription_data_->callback_functor.get());
  }
  key_expression_ = std::move(other.key_expression_);
  subscription_data_ = std::move(other.subscription_data_);
  return *this;
}

LivelinessSubscription::~LivelinessSubscription() { Unsubscribe(); }

void LivelinessSubscription::Unsubscribe() {
  if (!key_expression_.empty()) {
    Zenoh().imw_destroy_liveliness_subscription(
        subscription_data_->key_expression.c_str(),
        zenoh_static_liveliness_callback,
        subscription_data_->callback_functor.get());
    key_expression_.clear();
  }
}

}  // namespace intrinsic
