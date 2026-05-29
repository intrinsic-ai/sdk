// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_DATA_CALLBACK_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_DATA_CALLBACK_CONTEXT_H_

#include <string>

namespace intrinsic {

class IMWZenoh;

class IMWZenohDataCallbackContext {
 public:
  IMWZenohDataCallbackContext(IMWZenoh* imw_zenoh_instance,
                              const std::string& subscription_keyexpr);
  ~IMWZenohDataCallbackContext();

  // Move constructor is fine
  IMWZenohDataCallbackContext(IMWZenohDataCallbackContext&& other) = default;

  // Because of the pointer to the singleton, let's just disable copy
  // constructors so we don't have to think any harder about if they're valid
  // or not.
  IMWZenohDataCallbackContext(const IMWZenohDataCallbackContext&) = delete;
  IMWZenohDataCallbackContext& operator=(const IMWZenohDataCallbackContext&) =
      delete;

  IMWZenoh* get_imw_zenoh_instance() { return imw_zenoh_instance_; }
  const std::string& get_subscription_keyexpr() {
    return subscription_keyexpr_;
  }

 private:
  IMWZenoh* imw_zenoh_instance_;
  std::string subscription_keyexpr_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_DATA_CALLBACK_CONTEXT_H_
