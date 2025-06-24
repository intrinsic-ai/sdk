// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_GOLANG_PUBSUB_C_H_
#define INTRINSIC_PLATFORM_PUBSUB_GOLANG_PUBSUB_C_H_

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#include "absl/base/attributes.h"

#ifdef __cplusplus
extern "C" {
#endif

ABSL_ATTRIBUTE_UNUSED void* NewZenohHandle();
ABSL_ATTRIBUTE_UNUSED void DestroyZenohHandle(void* handle);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwInit(void* handle, const char* config);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwFini(void* handle);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwCreatePublisher(void* handle,
                                                        const char* keyexpr,
                                                        const char* qos);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwDestroyPublisher(void* handle,
                                                         const char* keyexpr);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwPublish(void* handle,
                                                const char* keyexpr,
                                                const void* bytes,
                                                size_t bytes_len);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwPublisherHasMatchingSubscribers(
    void* handle, const char* keyexpr, bool* has_matching);

typedef void (*zenoh_handle_imw_subscription_callback_fn)(
    const char* keyexpr, const void* bytes, const size_t bytes_len,
    void* user_context);

ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwCreateSubscription(
    void* handle, const char* keyexpr,
    zenoh_handle_imw_subscription_callback_fn callback, const char* qos,
    void* user_context);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwDestroySubscription(
    void* handle, const char* keyexpr,
    zenoh_handle_imw_subscription_callback_fn callback, void* user_context);

typedef void (*zenoh_handle_imw_queryable_callback_fn)(
    const char* keyexpr, const void* query_bytes, const size_t query_bytes_len,
    const void* query_context, void* user_context);
typedef void (*zenoh_handle_imw_query_callback_fn)(
    const char* keyexpr, const void* response_bytes,
    const size_t response_bytes_len, void* user_context);
typedef void (*zenoh_handle_imw_query_on_done_fn)(const char* keyexpr,
                                                  void* userContext);

// KV stuff now...
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwSet(void* handle, const char* keyexpr,
                                            const void* bytes,
                                            size_t bytes_len);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwDeleteKeyExpr(void* handle,
                                                      const char* keyexpr);

ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwCreateQueryable(
    void* handle, const char* keyexpr,
    zenoh_handle_imw_queryable_callback_fn callback, void* user_context,
    bool is_ros_service);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwDestroyQueryable(
    void* handle, const char* keyexpr,
    zenoh_handle_imw_queryable_callback_fn callback, void* user_context);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwQueryableReply(
    void* handle, const void* query_context, const char* keyexpr,
    const void* reply_bytes, size_t reply_bytes_len);
ABSL_ATTRIBUTE_UNUSED int ZenohHandleImwQuery(
    void* handle, const char* keyexpr,
    zenoh_handle_imw_query_callback_fn callback,
    zenoh_handle_imw_query_on_done_fn on_done, const void* query_payload,
    size_t query_payload_len, void* user_context, uint64_t timeout_ms,
    bool call_ros_service);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // INTRINSIC_PLATFORM_PUBSUB_GOLANG_PUBSUB_C_H_
