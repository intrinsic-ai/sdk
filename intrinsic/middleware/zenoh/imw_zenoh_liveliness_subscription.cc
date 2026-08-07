// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/middleware/zenoh/imw_zenoh_liveliness_subscription.h"

#include <string>
#include <utility>  // for std::move

#include "absl/log/log.h"

using std::string;

namespace intrinsic {

IMWZenohLivelinessSubscription::IMWZenohLivelinessSubscription(
    const std::string& keyexpr, imw_liveliness_callback_fn* callback,
    z_owned_subscriber_t zenoh_sub, void* user_context)
    : keyexpr_(keyexpr), zenoh_sub_(zenoh_sub), n_bytes(0), n_messages(0) {
  callbacks_.push_back(std::make_pair(callback, user_context));
}

IMWZenohLivelinessSubscription::~IMWZenohLivelinessSubscription() {
  // by the time this destructor fires, the subscription must have
  // previously been released by the owning IMWZenoh object that has
  // access to the Zenoh session. We'll just sanity-check it here.
  if (z_internal_check(zenoh_sub_)) {
    LOG(ERROR) << "~IMWZenohLivelinessSubscription had a valid sub!";
  }
}

void IMWZenohLivelinessSubscription::add_callback(
    imw_liveliness_callback_fn* fptr, void* user_context) {
  absl::MutexLock lock(&mutex_);
  callbacks_.push_back(std::make_pair(fptr, user_context));
}

bool IMWZenohLivelinessSubscription::clear_callbacks() {
  mutex_.ForgetDeadlockInfo();
  absl::MutexLock lock(&mutex_);
  callbacks_.clear();
  return true;
}

bool IMWZenohLivelinessSubscription::remove_callback(
    imw_liveliness_callback_fn* fptr, const void* user_context) {
  // The deadlock cycle detector does not like this, but I am currently
  // convinced this is safe.
  mutex_.ForgetDeadlockInfo();
  absl::MutexLock lock(&mutex_);
  for (auto it = callbacks_.begin(); it != callbacks_.end(); ++it) {
    if (it->first == fptr && it->second == user_context) {
      callbacks_.erase(it);
      return true;
    }
  }
  LOG(ERROR) << "Tried to remove callback to " << keyexpr_
             << " but could not find it";
  return false;
}

void IMWZenohLivelinessSubscription::invoke_callbacks(const char* keyexpr,
                                                      bool alive) {
  // People writing callbacks probably are assuming they are non-reentrant,
  // so we'll use a mutex to make sure that's true!
  absl::MutexLock lock(&mutex_);
  n_messages++;
  for (auto it = callbacks_.begin(); it != callbacks_.end(); ++it) {
    it->first(keyexpr, alive, it->second);
  }
}

IMWZenohLivelinessSubscription::Statistics
IMWZenohLivelinessSubscription::get_statistics() {
  // Lock the mutex to avoid tearing if a messages arrives at this very
  // instant, causing a statistics field update as we are reading it.
  absl::MutexLock lock(&mutex_);
  Statistics s;
  s.n_messages = n_messages;
  s.n_bytes = n_bytes;
  return s;
}

}  // namespace intrinsic
