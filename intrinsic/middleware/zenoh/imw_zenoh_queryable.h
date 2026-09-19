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

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_H_

#include <string>
#include <utility>
#include <vector>

#include "absl/log/log.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/middleware/imw.h"
#include "zenoh.h"  // NOLINT(build/include_subdir)

struct z_loaned_query_t;

namespace intrinsic {

class IMWZenohQueryable {
 public:
  IMWZenohQueryable(const std::string& keyexpr,
                    imw_queryable_callback_fn* callback,
                    z_owned_queryable_t zenoh_queryable, void* user_context,
                    imw_queryable_options_t* options);

  IMWZenohQueryable(IMWZenohQueryable&& other) = delete;
  IMWZenohQueryable(const IMWZenohQueryable&) = delete;
  IMWZenohQueryable& operator=(const IMWZenohQueryable&) = delete;

  ~IMWZenohQueryable();

  z_owned_queryable_t& get_zenoh_queryable() { return zenoh_queryable_; }
  const std::string& get_keyexpr() const { return keyexpr_; }
  const void* get_user_context() const { return user_context_; }
  imw_queryable_callback_fn* get_callback() const { return callback_; }
  void invoke(const char* keyexpr, const void* query_payload,
              const size_t query_payload_len,
              const z_loaned_query_t* query_context);
  const imw_queryable_options_t& get_options() const { return options_; }

 private:
  std::string keyexpr_;
  imw_queryable_callback_fn* callback_;
  z_owned_queryable_t zenoh_queryable_;
  absl::Mutex mutex_;  // protects internal data structures
  void* user_context_;
  imw_queryable_options_t options_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_QUERYABLE_H_
