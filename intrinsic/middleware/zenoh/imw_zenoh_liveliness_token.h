// Copyright 2023 Intrinsic Innovation LLC

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
