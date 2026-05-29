// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/middleware/imw.h"
#include "intrinsic/middleware/zenoh/imw_zenoh.h"
#include "intrinsic/middleware/zenoh/imw_zenoh_data_callback_context.h"
#include "intrinsic/middleware/zenoh/imw_zenoh_query_context.h"
#include "intrinsic/middleware/zenoh/imw_zenoh_queryable_context.h"

namespace intrinsic {

static IMWZenoh* g_imw_zenoh_singleton = nullptr;
ABSL_CONST_INIT absl::Mutex IMWZenoh::init_fini_mutex_(absl::kConstInit);
static int g_imw_init_refcount = 0;

void IMWZenoh::static_data_callback(z_loaned_sample_t* sample,
                                    void* untyped_context) {
  IMWZenohDataCallbackContext* zenoh_context =
      static_cast<IMWZenohDataCallbackContext*>(untyped_context);

  zenoh_context->get_imw_zenoh_instance()->data_callback(
      zenoh_context->get_subscription_keyexpr(), sample);
}

void IMWZenoh::static_closure_drop(void* untyped_context) {
  // This callback is run by Zenoh during a call to z_undeclare_subscriber().
  // The intent is to use it to free the user_context memory that was allocated
  // on the heap and passed to Zenoh as part of the call to
  // declare_subscriber().
  IMWZenohDataCallbackContext* zenoh_context =
      static_cast<IMWZenohDataCallbackContext*>(untyped_context);
  delete zenoh_context;
}

void IMWZenoh::static_queryable_callback(z_loaned_query_t* query,
                                         void* untyped_context) {
  IMWZenohQueryableContext* context =
      static_cast<IMWZenohQueryableContext*>(untyped_context);
  context->get_imw_zenoh_instance()->queryable_callback(
      context->get_queryable_keyexpr(), query);
}

void IMWZenoh::static_queryable_drop(void* untyped_context) {
  // This callback is run by Zenoh during a call to z_undeclare_queryable().
  // The intent is to use it to free the user_context memory that was allocated
  // on the heap and passed to Zenoh as part of the call to
  // declare_queryable().
  IMWZenohQueryableContext* typed_context =
      static_cast<IMWZenohQueryableContext*>(untyped_context);
  delete typed_context;
}

void IMWZenoh::static_query_callback(z_loaned_reply_t* reply,
                                     void* untyped_context) {
  IMWZenohQueryContext* context =
      static_cast<IMWZenohQueryContext*>(untyped_context);
  context->imw_zenoh_instance_->query_callback(
      context->keyexpr_, context->callback_, reply, context->user_context_,
      &context->options_);
}

void IMWZenoh::static_query_drop(void* untyped_context) {
  IMWZenohQueryContext* typed_context =
      static_cast<IMWZenohQueryContext*>(untyped_context);
  if (typed_context != nullptr && typed_context->on_done_ != nullptr) {
    typed_context->on_done_(typed_context->keyexpr_,
                            typed_context->user_context_);
  }
  delete typed_context;
}

imw_ret_t imw_init(const char* config) {
  absl::MutexLock lock(&IMWZenoh::init_fini_mutex_);
  g_imw_init_refcount++;
  if (g_imw_zenoh_singleton == nullptr) g_imw_zenoh_singleton = new IMWZenoh();

  return g_imw_zenoh_singleton->create_session(config);
}

imw_ret_t imw_fini() {
  absl::MutexLock lock(&IMWZenoh::init_fini_mutex_);
  g_imw_init_refcount--;
  if (g_imw_zenoh_singleton == nullptr) return IMW_OK;
  if (g_imw_init_refcount > 0) return IMW_OK;

  const imw_ret_t rc = g_imw_zenoh_singleton->destroy_session();
  if (rc != IMW_OK) return rc;

  delete g_imw_zenoh_singleton;
  g_imw_zenoh_singleton = nullptr;

  return IMW_OK;
}

imw_ret_t imw_create_publisher(const char* keyexpr, const char* qos) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;

  return g_imw_zenoh_singleton->create_publisher(keyexpr, qos);
}

imw_ret_t imw_destroy_publisher(const char* keyexpr) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;

  return g_imw_zenoh_singleton->destroy_publisher(keyexpr);
}

imw_ret_t imw_publish(const char* keyexpr, const void* bytes,
                      const size_t bytes_len) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;

  return g_imw_zenoh_singleton->publish(keyexpr, bytes, bytes_len);
}

imw_ret_t imw_create_subscription(const char* keyexpr,
                                  imw_subscription_callback_fn* callback,
                                  const char* qos, void* user_context) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;

  return g_imw_zenoh_singleton->create_subscription(keyexpr, callback, qos,
                                                    user_context);
}

imw_ret_t imw_destroy_subscription(const char* keyexpr,
                                   imw_subscription_callback_fn* callback,
                                   const void* user_context) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;

  return g_imw_zenoh_singleton->destroy_subscription(keyexpr, callback,
                                                     user_context);
}

int imw_keyexpr_includes(const char* left, const char* right) {
  return IMWZenoh::keyexpr_includes(left, right);
}

int imw_keyexpr_intersects(const char* left, const char* right) {
  return IMWZenoh::keyexpr_intersects(left, right);
}

int imw_keyexpr_is_canon(const char* keyexpr) {
  return IMWZenoh::keyexpr_is_canon(keyexpr);
}

imw_ret_t imw_queryable_reply(const void* query_context, const char* keyexpr,
                              const void* reply_bytes,
                              const size_t reply_bytes_len) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->queryable_reply(query_context, keyexpr,
                                                reply_bytes, reply_bytes_len);
}

imw_ret_t imw_create_queryable(const char* keyexpr,
                               imw_queryable_callback_fn* callback,
                               void* user_context,
                               imw_queryable_options_t* options) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->create_queryable(keyexpr, callback,
                                                 user_context, options);
}

imw_ret_t imw_destroy_queryable(const char* keyexpr,
                                imw_queryable_callback_fn* callback,
                                void* user_context) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->destroy_queryable(keyexpr, callback,
                                                  user_context);
}

imw_ret_t imw_query(const char* keyexpr, imw_query_callback_fn* callback,
                    imw_query_on_done_callback_fn* on_done,
                    const void* query_payload, const size_t query_payload_len,
                    void* user_context, imw_query_options_t* options) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->query(keyexpr, callback, on_done, query_payload,
                                      query_payload_len, user_context, options);
}

imw_ret_t imw_set(const char* keyexpr, const void* bytes,
                  const size_t bytes_len) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->set(keyexpr, bytes, bytes_len);
}

imw_ret_t imw_delete_keyexpr(const char* keyexpr) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->delete_keyexpr(keyexpr);
}

const char* const imw_version() { return IMWZenoh::version(); }

imw_ret_t imw_publisher_has_matching_subscribers(const char* keyexpr,
                                                 bool* has_matching) {
  if (g_imw_zenoh_singleton == nullptr) return IMW_NOT_INITIALIZED;
  return g_imw_zenoh_singleton->publisher_has_matching_subscribers(
      keyexpr, has_matching);
}

}  // namespace intrinsic
