// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_KVSTORE_H_
#define INTRINSIC_PLATFORM_PUBSUB_KVSTORE_H_

#include <functional>
#include <memory>
#include <optional>
#include <string>
#include <type_traits>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/flags/declare.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/message.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_DECLARE_FLAG(bool, use_replicated_kv_store);

namespace intrinsic {

constexpr absl::Duration kDefaultGetTimeout = absl::Seconds(10);
constexpr absl::Duration kDefaultAdminCloudCopyTimeout = absl::Seconds(20);
constexpr absl::string_view kDefaultKeyPrefix = "kv_store";
constexpr absl::string_view kReplicationPrefix = "kv_store_repl";

using KeyValueCallback = std::function<void(
    absl::string_view key, std::unique_ptr<google::protobuf::Any> value)>;

// Callback invoked when the KeyValueCallback is called for all keys that match.
// Make sure to keep this callback lightweight.
using OnDoneCallback = std::function<void(absl::string_view key)>;

class KVQuery {
 public:
  explicit KVQuery(std::unique_ptr<imw_callback_functor_t> callback,
                   std::unique_ptr<imw_on_done_functor_t> on_done)
      : callback_(std::move(callback)),
        on_done_(std::move(on_done)),
        context_(
            std::make_unique<QueryContext>(callback_.get(), on_done_.get())) {}

  QueryContext* GetContext() { return context_.get(); }

 private:
  std::unique_ptr<imw_callback_functor_t> callback_;
  std::unique_ptr<imw_on_done_functor_t> on_done_;
  std::unique_ptr<QueryContext> context_;
};

class KeyValueStore {
 public:
  friend class PubSub;

  // Sets the value for the given key. A key can't include any of the following
  // characters: /, *, ?, #, [ and ].
  absl::Status Set(absl::string_view key, const google::protobuf::Any& value,
                   std::optional<bool> high_consistency = std::nullopt);

  // Sets the value for the given key. A key can't include any of the following
  // characters: /, *, ?, #, [ and ].
  template <typename T>
  absl::Status Set(absl::string_view key, T&& value,
                   std::optional<bool> high_consistency = std::nullopt)
    requires(
        std::is_base_of_v<google::protobuf::Message, std::remove_cvref_t<T>> &&
        !std::is_same_v<google::protobuf::Any, std::remove_cvref_t<T>>)
  {
    google::protobuf::Any any;
    if (!any.PackFrom(std::forward<T>(value))) {
      return absl::InternalError(
          absl::StrCat("Failed to pack value for the key: ", key));
    }
    return Set(key, any, high_consistency);
  }

  template <typename T>
  absl::StatusOr<T> Get(absl::string_view key,
                        absl::Duration timeout = kDefaultGetTimeout)
    requires(std::is_base_of_v<google::protobuf::Message, T>)
  {
    INTR_ASSIGN_OR_RETURN(google::protobuf::Any any_value,
                          GetAny(key, timeout));
    if constexpr (std::is_same_v<google::protobuf::Any, T>) {
      return any_value;
    } else {
      T value;
      if (!any_value.UnpackTo(&value)) {
        return absl::InternalError(
            absl::StrCat("Failed to unpack value for the key: ", key));
      }
      return value;
    }
  }

  // For a given key and WildcardQueryConfig, the KeyValueCallback will be
  // invoked for each key that matches the expression. The caller is expected to
  // keep the Query object alive until the OnDoneCallback is called.
  absl::StatusOr<KVQuery> GetAll(
      absl::string_view key, KeyValueCallback callback,
      OnDoneCallback on_done = [](absl::string_view key) {});

  // Deletes the key from the KVStore.
  absl::Status Delete(absl::string_view key);

  // Lists all keys in the non replicated KVStore. Returns an error if called on
  // replicated KVStore. Essentially lists keys in kv_store/**
  absl::StatusOr<std::vector<std::string>> ListAllKeys(
      absl::Duration timeout = kDefaultGetTimeout);

  // Lists all keys in the global cloud KVStore key space.
  // Essentially lists keys in kv_store_repl/global/**
  absl::StatusOr<std::vector<std::string>> ListAllGlobalKeys(
      absl::Duration timeout = kDefaultGetTimeout);

  // Lists all keys in the onprem replicated KVStore. Essentially lists keys in
  // kv_store_repl/<workcell_name>/**
  absl::StatusOr<std::vector<std::string>> ListAllOnpremKeys(
      absl::string_view workcell_name,
      absl::Duration timeout = kDefaultGetTimeout);

  // Use this method to copy local key-value pairs to the cloud key value
  // store. For eg, you can use this method to copy kv_store/<current_ipc>/key
  // to kv_store_repl/<destination_ipc>/key. To use this method, the you must be
  // running on a cluster with credentials that allow cloud ingress access. The
  // timeout is not enforced for the entire duration of the copy, but rather for
  // the call to the cloud server.
  absl::Status AdminCloudCopy(absl::string_view source_key,
                              absl::string_view target_key,
                              absl::Duration timeout);

  // Same as GetAll, but does not need a callback. The tradeoff is less control.
  absl::StatusOr<absl::flat_hash_map<std::string, google::protobuf::Any>>
  GetAllSynchronous(absl::string_view keyexpr, absl::Duration timeout);

 private:
  explicit KeyValueStore(std::optional<std::string> prefix_override);

  absl::StatusOr<std::vector<std::string>> ExecuteList(
      absl::string_view keyexpr, absl::Duration timeout);

  absl::StatusOr<google::protobuf::Any> GetAny(absl::string_view key,
                                               absl::Duration timeout);

  std::string key_prefix_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_KVSTORE_H_
