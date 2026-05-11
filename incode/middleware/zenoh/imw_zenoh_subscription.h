// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_SUBSCRIPTION_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_SUBSCRIPTION_H_

#include <string>
#include <utility>
#include <vector>

#include "absl/log/log.h"
#include "absl/synchronization/mutex.h"
#include "incode/middleware/imw.h"
#include "zenoh.h"  // NOLINT(build/include_subdir)

namespace intrinsic {

class IMWZenohSubscription {
 public:
  IMWZenohSubscription(const std::string& keyexpr,
                       imw_subscription_callback_fn* callback,
                       z_owned_subscriber_t zenoh_sub, void* user_context);

  IMWZenohSubscription(IMWZenohSubscription&& other) = delete;
  IMWZenohSubscription(const IMWZenohSubscription&) = delete;
  IMWZenohSubscription& operator=(const IMWZenohSubscription&) = delete;

  ~IMWZenohSubscription();

  z_owned_subscriber_t& get_zenoh_sub() { return zenoh_sub_; }
  const std::string& get_keyexpr() const { return keyexpr_; }
  void add_callback(imw_subscription_callback_fn* fptr, void* user_context);
  bool remove_callback(imw_subscription_callback_fn* fptr,
                       const void* user_context);
  bool clear_callbacks();
  bool is_empty() const { return callbacks_.empty(); }
  void invoke_callbacks(const char* keyexpr, const void* blob,
                        const size_t blob_len);
  struct Statistics {
    size_t n_messages;
    size_t n_bytes;
  };
  Statistics get_statistics();

 private:
  std::string keyexpr_;
  std::vector<std::pair<imw_subscription_callback_fn*, void*>> callbacks_;
  z_owned_subscriber_t zenoh_sub_;
  absl::Mutex mutex_;  // protects internal data structures
  size_t n_bytes;
  size_t n_messages;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_SUBSCRIPTION_H_
