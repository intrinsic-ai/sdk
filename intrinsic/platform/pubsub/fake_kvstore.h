// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_FAKE_KVSTORE_H_
#define INTRINSIC_PLATFORM_PUBSUB_FAKE_KVSTORE_H_

#include <string>

#include "absl/base/thread_annotations.h"
#include "absl/container/node_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "google/protobuf/any.pb.h"
#include "intrinsic/platform/pubsub/kvstore.h"

namespace intrinsic {

// A fake implementation of KeyValueStore for use in unit tests.
//
// This class stores key-value pairs in an in-memory `absl::node_hash_map`.
// It avoids dependencies on Zenoh or real PubSub infrastructure, making it
// fast and suitable for testing components that require a KeyValueStore.
//
// Example usage:
//   std::unique_ptr<KeyValueStore> kv_store =
//    std::make_unique<FakeKeyValueStore>();
//   kv_store->Set("my_key", my_proto_value);
//   // Pass kv_store to the class under test.
class FakeKeyValueStore : public KeyValueStore {
 public:
  FakeKeyValueStore() : KeyValueStore(/*prefix_override=*/std::nullopt) {}
  ~FakeKeyValueStore() override = default;

  using KeyValueStore::Set;

  absl::Status Set(
      absl::string_view key, const google::protobuf::Any& value,
      std::optional<bool> high_consistency = std::nullopt) override;

  absl::Status Delete(absl::string_view key) override;

 protected:
  absl::StatusOr<google::protobuf::Any> GetAny(absl::string_view key,
                                               absl::Duration timeout) override;

 private:
  absl::Mutex mutex_;
  absl::node_hash_map<std::string, google::protobuf::Any> store_
      ABSL_GUARDED_BY(mutex_);
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_FAKE_KVSTORE_H_
