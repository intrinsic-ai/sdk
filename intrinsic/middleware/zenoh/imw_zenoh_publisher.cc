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

#include "intrinsic/middleware/zenoh/imw_zenoh_publisher.h"

#include <string>
#include <utility>  // for std::move

using std::string;

namespace intrinsic {

IMWZenohPublisher::IMWZenohPublisher(const std::string& keyexpr,
                                     z_owned_publisher_t zenoh_pub,
                                     z_owned_keyexpr_t zenoh_keyexpr)
    : marked_for_deletion_(false),
      keyexpr_(keyexpr),
      zenoh_pub_(zenoh_pub),
      zenoh_keyexpr_(zenoh_keyexpr),
      n_bytes(0),
      n_messages(0) {}

IMWZenohPublisher::~IMWZenohPublisher() {
  if (z_internal_check(zenoh_keyexpr_)) {
    z_drop(z_move(zenoh_keyexpr_));
  }
  // by the time this destructor fires, the publisher must have
  // previously been released by the owning IMWZenoh object that has
  // access to the Zenoh session. We'll just sanity-check it here.
  if (z_internal_check(zenoh_pub_)) {
    fprintf(stderr, "ERROR! ~IMWZenohPublisher had a valid pub!\n");
  }
}

void IMWZenohPublisher::record_message_size(size_t msg_len) {
  n_messages++;
  n_bytes += msg_len;
}

}  // namespace intrinsic
