// Copyright 2023 Intrinsic Innovation LLC

#ifndef INCODE_MIDDLEWARE_IMW_H_
#define INCODE_MIDDLEWARE_IMW_H_

#include <stddef.h>  // for size_t
#include <stdint.h>  // for uint64_t

namespace intrinsic {

typedef enum imw_ret {
  IMW_OK = 0,
  IMW_ERROR = 1,
  IMW_NOT_INITIALIZED = 2,
  IMW_UNDEFINED = 3,
} imw_ret_t;

typedef void imw_subscription_callback_fn(const char* keyexpr,
                                          const void* bytes,
                                          const size_t bytes_len,
                                          void* user_context);

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

imw_ret_t imw_init(const char* config);
imw_ret_t imw_fini();
imw_ret_t imw_create_publisher(const char* keyexpr, const char* qos);
imw_ret_t imw_destroy_publisher(const char* keyexpr);
imw_ret_t imw_publish(const char* keyexpr, const void* bytes,
                      const size_t bytes_len);
imw_ret_t imw_create_subscription(const char* keyexpr,
                                  imw_subscription_callback_fn* callback,
                                  const char* qos, void* user_context);
imw_ret_t imw_publisher_has_matching_subscribers(const char* keyexpr,
                                                 bool* has_matching);
imw_ret_t imw_destroy_subscription(const char* keyexpr,
                                   imw_subscription_callback_fn* callback,
                                   const void* user_context);
int imw_keyexpr_includes(const char* left, const char* right);
int imw_keyexpr_intersects(const char* left, const char* right);
int imw_keyexpr_is_canon(const char* keyexpr);
imw_ret_t imw_create_queryable(const char* keyexpr,
                               imw_queryable_callback_fn* callback,
                               void* user_context,
                               imw_queryable_options_t* options);
imw_ret_t imw_destroy_queryable(const char* keyexpr,
                                imw_queryable_callback_fn* callback,
                                void* user_context);
imw_ret_t imw_queryable_reply(const void* query_context, const char* keyexpr,
                              const void* reply_bytes,
                              const size_t reply_bytes_len);
imw_ret_t imw_query(const char* keyexpr, imw_query_callback_fn* callback,
                    imw_query_on_done_callback_fn* on_done,
                    const void* query_payload, const size_t query_payload_len,
                    void* user_context, imw_query_options_t* options);
imw_ret_t imw_set(const char* keyexpr, const void* bytes,
                  const size_t bytes_len);
imw_ret_t imw_delete_keyexpr(const char* keyexpr);
const char* const imw_version();

}  // namespace intrinsic

#endif  // INCODE_MIDDLEWARE_IMW_H_
