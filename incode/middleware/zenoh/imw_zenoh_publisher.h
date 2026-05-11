// Copyright 2023 Intrinsic Innovation LLC

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_PUBLISHER_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_PUBLISHER_H_

#include <string>
#include <utility>
#include <vector>

#include "absl/synchronization/mutex.h"
#include "incode/middleware/imw.h"
#include "zenoh.h"  // NOLINT(build/include_subdir)

namespace intrinsic {

class IMWZenohPublisher {
 public:
  IMWZenohPublisher(const std::string& keyexpr, z_owned_publisher_t zenoh_pub,
                    z_owned_keyexpr_t zenoh_keyexpr);

  IMWZenohPublisher(IMWZenohPublisher&& other) = delete;
  IMWZenohPublisher(const IMWZenohPublisher&) = delete;
  IMWZenohPublisher& operator=(const IMWZenohPublisher&) = delete;

  ~IMWZenohPublisher();

  z_owned_publisher_t& get_zenoh_pub() { return zenoh_pub_; }
  const std::string& get_keyexpr() const { return keyexpr_; }
  z_owned_keyexpr_t& get_zenoh_keyexpr() { return zenoh_keyexpr_; }
  void record_message_size(size_t msg_len);
  size_t get_n_messages() const { return n_messages; }
  size_t get_n_bytes() const { return n_bytes; }

 public:
  std::atomic<bool> marked_for_deletion_;

 private:
  std::string keyexpr_;
  z_owned_publisher_t zenoh_pub_;
  z_owned_keyexpr_t zenoh_keyexpr_;
  size_t n_bytes;
  size_t n_messages;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_PUBLISHER_H_
