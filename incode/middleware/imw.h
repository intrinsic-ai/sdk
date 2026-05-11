// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_IMW_H_
#define MIDDLEWARE_IMW_H_

#include <stddef.h>  // for size_t
#include <stdint.h>  // for uint64_t

typedef enum imw_ret {
  IMW_OK = 0,
  IMW_ERROR = 1,
  IMW_NOT_INITIALIZED = 2,
  IMW_UNDEFINED = 3,
} imw_ret_t;

typedef imw_ret_t imw_init_fn(const char* config);

typedef imw_ret_t imw_fini_fn();

typedef imw_ret_t imw_create_publisher_fn(const char* keyexpr, const char* qos);

typedef imw_ret_t imw_destroy_publisher_fn(const char* keyexpr);

typedef imw_ret_t imw_publish_fn(const char* keyexpr, const void* bytes,
                                 const size_t bytes_len);

typedef void imw_subscription_callback_fn(const char* keyexpr,
                                          const void* bytes,
                                          const size_t bytes_len,
                                          void* user_context);

typedef imw_ret_t imw_publisher_has_matching_subscribers_fn(const char* keyexpr,
                                                            bool* has_matching);

typedef imw_ret_t imw_create_subscription_fn(
    const char* keyexpr, imw_subscription_callback_fn* callback,
    const char* qos, void* user_context);

typedef imw_ret_t imw_destroy_subscription_fn(
    const char* keyexpr, imw_subscription_callback_fn* callback,
    const void* user_context);

typedef int imw_intersects_fn(const char* left, const char* right);
typedef int imw_includes_fn(const char* left, const char* right);
typedef int imw_keyexpr_is_canon_fn(const char* keyexpr);

typedef imw_ret_t imw_queryable_reply_fn(const void* query_context,
                                         const char* keyexpr,
                                         const void* reply_bytes,
                                         const size_t reply_bytes_len);

typedef void imw_queryable_callback_fn(const char* keyexpr,
                                       const void* query_bytes,
                                       const size_t query_bytes_len,
                                       const void* query_context,
                                       void* user_context);

// This struct must be assignable, as it will be copied by assignment by
// imw_create_queryable() for later use.
struct imw_queryable_options_t {
  bool is_ros_service = false;
};

typedef imw_ret_t imw_create_queryable_fn(const char* keyexpr,
                                          imw_queryable_callback_fn* callback,
                                          void* user_context,
                                          imw_queryable_options_t* options);

typedef imw_ret_t imw_destroy_queryable_fn(const char* keyexpr,
                                           imw_queryable_callback_fn* callback,
                                           void* user_context);

typedef void imw_query_callback_fn(const char* keyexpr,
                                   const void* response_bytes,
                                   const size_t response_bytes_len,
                                   void* user_context);

typedef void imw_query_on_done_callback_fn(const char* keyexpr,
                                           void* user_context);

// This struct must be assignable, as it will be copied by assignment by
// imw_query() for later use.
struct imw_query_options_t {
  uint64_t timeout_ms = 0;
  bool call_ros_service = false;
};

typedef imw_ret_t imw_query_fn(const char* keyexpr,
                               imw_query_callback_fn* callback,
                               imw_query_on_done_callback_fn* on_done,
                               const void* query_payload,
                               const size_t query_payload_len,
                               void* user_context,
                               imw_query_options_t* options);

typedef const char* const imw_version_fn();

typedef imw_ret_t imw_set_fn(const char* keyexpr, const void* bytes,
                             const size_t bytes_len);

typedef imw_ret_t imw_delete_keyexpr_fn(const char* keyexpr,
                                        void* user_context);

#endif  // MIDDLEWARE_IMW_H_
