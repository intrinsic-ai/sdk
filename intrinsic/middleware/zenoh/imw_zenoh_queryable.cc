// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/middleware/zenoh/imw_zenoh_queryable.h"

#include <string>
#include <utility>  // for std::move

#include "absl/log/log.h"
#include "intrinsic/middleware/zenoh/imw_zenoh_reply_context.h"

using std::string;

namespace intrinsic {

IMWZenohQueryable::IMWZenohQueryable(const std::string& keyexpr,
                                     imw_queryable_callback_fn* callback,
                                     z_owned_queryable_t zenoh_queryable,
                                     void* user_context,
                                     imw_queryable_options_t* options)
    : keyexpr_(keyexpr),
      callback_(callback),
      zenoh_queryable_(zenoh_queryable),
      user_context_(user_context) {
  if (options != nullptr) {
    options_ = *options;
  }
}

IMWZenohQueryable::~IMWZenohQueryable() {
  // by the time this destructor fires, the subscription must have
  // previously been released by the owning IMWZenoh object that has
  // access to the Zenoh session. We'll just sanity-check it here.
  if (z_internal_check(zenoh_queryable_)) {
    LOG(ERROR) << "~IMWZenohQueryable had a valid queryable! This is bad!";
  }
}

void IMWZenohQueryable::invoke(const char* keyexpr, const void* query_payload,
                               const size_t query_payload_len,
                               const z_loaned_query_t* query_context) {
  const IMWZenohReplyContext reply_context(query_context, &options_);
  // People writing callbacks probably are assuming they are non-reentrant,
  // so we'll use a mutex to make sure that's true!
  absl::MutexLock lock(&mutex_);
  callback_(keyexpr, query_payload, query_payload_len, &reply_context,
            user_context_);
}

}  // namespace intrinsic
