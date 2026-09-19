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

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_GET_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_GET_CONTEXT_H_

#include <string>

#include "intrinsic/middleware/imw.h"

namespace intrinsic {

class IMWZenoh;

class IMWLivelinessGetContext {
 public:
  IMWLivelinessGetContext(IMWZenoh* imw_zenoh_instance, const char* keyexpr,
                          imw_liveliness_get_callback_fn* callback,
                          imw_liveliness_get_on_done_callback_fn* on_done,
                          void* user_context)
      : imw_zenoh_instance_(imw_zenoh_instance),
        keyexpr_(keyexpr),
        callback_(callback),
        on_done_(on_done),
        user_context_(user_context) {}

  IMWZenoh* imw_zenoh_instance_;
  std::string keyexpr_;
  imw_liveliness_get_callback_fn* callback_;
  imw_liveliness_get_on_done_callback_fn* on_done_;
  void* user_context_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_GET_CONTEXT_H_
