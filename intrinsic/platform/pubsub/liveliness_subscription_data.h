// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_DATA_H_
#define INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_DATA_H_

#include <memory>
#include <string>

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

namespace intrinsic {

// Contains information about a subscription to liveliness notifications
// that should be shared with lower-level APIs (ZenohHandle and below).
struct LivelinessSubscriptionData {
  std::unique_ptr<imw_liveliness_callback_functor_t> callback_functor;
  std::string key_expression;
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_LIVELINESS_SUBSCRIPTION_DATA_H_
