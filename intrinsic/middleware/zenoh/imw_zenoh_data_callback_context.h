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
