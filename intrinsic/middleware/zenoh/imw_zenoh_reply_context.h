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

#ifndef MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_
#define MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_

#include <string>

struct imw_queryable_options_t;
struct z_loaned_query_t;

namespace intrinsic {

class IMWZenohReplyContext {
 public:
  IMWZenohReplyContext(const z_loaned_query_t* query,
                       const imw_queryable_options_t* options)
      : query_(query), options_(options) {}
  ~IMWZenohReplyContext() = default;

  // Move constructor is fine
  IMWZenohReplyContext(IMWZenohReplyContext&& other) = default;

  // No need for copy constructors in the intended usage of this class.
  IMWZenohReplyContext(const IMWZenohReplyContext&) = delete;
  IMWZenohReplyContext& operator=(const IMWZenohReplyContext&) = delete;

  const z_loaned_query_t* query_;
  const imw_queryable_options_t* options_;
};

}  // namespace intrinsic

#endif  // MIDDLEWARE_ZENOH_IMW_ZENOH_REPLY_CONTEXT_H_
