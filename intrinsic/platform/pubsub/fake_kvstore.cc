// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/fake_kvstore.h"

#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"

namespace intrinsic {

absl::Status FakeKeyValueStore::Set(absl::string_view key,
                                    const google::protobuf::Any& value,
                                    std::optional<bool> high_consistency) {
  absl::MutexLock lock(&mutex_);
  store_[std::string(key)] = value;
  return absl::OkStatus();
}

absl::Status FakeKeyValueStore::Delete(absl::string_view key) {
  absl::MutexLock lock(&mutex_);
  if (!store_.contains(key)) {
    return absl::NotFoundError(absl::StrCat("Key not found: ", key));
  }
  store_.erase(key);
  return absl::OkStatus();
}

absl::StatusOr<google::protobuf::Any> FakeKeyValueStore::GetAny(
    absl::string_view key, absl::Duration timeout) {
  absl::MutexLock lock(&mutex_);
  auto it = store_.find(key);
  if (it == store_.end()) {
    return absl::NotFoundError(absl::StrCat("Key not found: ", key));
  }
  return it->second;
}

}  // namespace intrinsic
