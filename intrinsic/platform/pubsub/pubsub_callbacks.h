// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_PUBSUB_CALLBACKS_H_
#define INTRINSIC_PLATFORM_PUBSUB_PUBSUB_CALLBACKS_H_

#include <functional>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"

namespace intrinsic {

// The following two callbacks are defined to be used asynchronously when a
// message arrives on a topic that a subscription was requested to.
//
// SubscriptionOkCallback is called whenever a valid message is received. This
// callback returns typed message.
template <typename T>
using SubscriptionOkCallback = std::function<void(const T& message)>;

template <typename T>
using SubscriptionOkExpandedCallback =
    std::function<void(absl::string_view topic, const T& message)>;

// Called when a key is deleted from a key-value store.
using DeletionCallback = std::function<void(std::string_view key)>;

// SubscriptionErrorCallback is called whenever a message is received on a
// subscribed topic but the message could not be parsed or converted to the
// desired type. This function returns the raw value of the received packet and
// status error indicating the problem.
using SubscriptionErrorCallback =
    std::function<void(absl::string_view packet, absl::Status error)>;

using SubscriptionErrorExpandedCallback = std::function<void(
    absl::string_view topic, absl::string_view packet, absl::Status error)>;

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_PUBSUB_CALLBACKS_H_
