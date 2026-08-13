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

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_TOKEN_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_TOKEN_H_

#include "zenoh.h"  // NOLINT(build/include_subdir)

namespace intrinsic {

class IMWZenohLivelinessToken {
 public:
  explicit IMWZenohLivelinessToken(z_owned_liveliness_token_t token)
      : token_(token) {}
  IMWZenohLivelinessToken(const IMWZenohLivelinessToken& other) = delete;
  IMWZenohLivelinessToken(IMWZenohLivelinessToken&& other) = delete;
  IMWZenohLivelinessToken& operator=(const IMWZenohLivelinessToken&) = delete;
  IMWZenohLivelinessToken& operator=(IMWZenohLivelinessToken&&) = delete;

  z_owned_liveliness_token_t& get_token() { return token_; }

 private:
  z_owned_liveliness_token_t token_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_LIVELINESS_TOKEN_H_
