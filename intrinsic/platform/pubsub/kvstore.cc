// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/kvstore.h"

#include <cstddef>
#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/notification.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_helpers.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

namespace {
constexpr static absl::string_view kDefaultKeyPrefix = "kv_store";
constexpr static absl::Duration kHighConsistencyTimeout = absl::Seconds(30);
}  // namespace

KeyValueStore::KeyValueStore(std::optional<std::string> prefix_override)
    : key_prefix_(prefix_override.has_value() ? prefix_override.value()
                                              : kDefaultKeyPrefix) {}

absl::Status KeyValueStore::Set(absl::string_view key,
                                const google::protobuf::Any& value,
                                const NamespaceConfig& config,
                                std::optional<bool> high_consistency) {
  INTR_RETURN_IF_ERROR(intrinsic::ValidZenohKeyexpr(key));
  absl::StatusOr<std::string> prefixed_name =
      ZenohHandle::add_key_prefix(key, key_prefix_);
  if (!prefixed_name.ok()) {
    // Should not happen since ValidKeyexpr was called before this.
    return prefixed_name.status();
  }
  imw_ret_t ret =
      Zenoh().imw_set(prefixed_name->c_str(), value.SerializeAsString().c_str(),
                      value.ByteSizeLong());
  if (ret != IMW_OK) {
    return absl::InternalError(
        absl::StrFormat("Error setting a key, return code: %d", ret));
  }
  // If high consistency is set, we need to block until the key value is
  // committed.
  if (high_consistency.has_value() && *high_consistency) {
    absl::Time deadline = absl::Now() + kHighConsistencyTimeout;
    while (true) {
      auto set_result = GetAny(key, config, absl::Seconds(10));
      if (set_result.ok() &&
          set_result.value().SerializeAsString() == value.SerializeAsString()) {
        // Key value is committed.
        return absl::OkStatus();
      } else {
        if (!set_result.ok() &&
            set_result.status() != absl::NotFoundError("Key not found")) {
          return absl::InternalError(
              absl::StrFormat("Unexpected error, return code: %d", ret));
        }
        if (absl::Now() > deadline) {
          return absl::DeadlineExceededError(
              "Timeout waiting for high consistency");
        }
        // Small wait before retrying.
        absl::SleepFor(absl::Milliseconds(100));
      }
    }
  }
  return absl::OkStatus();
}

absl::StatusOr<google::protobuf::Any> KeyValueStore::GetAny(
    absl::string_view key, const NamespaceConfig& config,
    absl::Duration timeout) {
  INTR_RETURN_IF_ERROR(intrinsic::ValidZenohKey(key));
  INTR_ASSIGN_OR_RETURN(absl::StatusOr<std::string> prefixed_name,
                        ZenohHandle::add_key_prefix(key, key_prefix_));
  if (timeout < absl::ZeroDuration()) {
    return absl::InvalidArgumentError("Timeout must be zero or positive");
  }
  google::protobuf::Any value;
  absl::Notification notif;
  absl::Status lambda_status = absl::NotFoundError("Key not found");
  // We can capture all variables by reference because we wait for the
  // imw_query() to finish before before returning. This ensures that all local
  // variables will outlive the callback.
  auto reply_functor = std::make_unique<imw_callback_functor_t>(
      [&value, &lambda_status](const char* keyexpr, const void* response_bytes,
                               const size_t response_bytes_len) {
        bool ok = value.ParseFromString(absl::string_view(
            static_cast<const char*>(response_bytes), response_bytes_len));
        if (ok) {
          lambda_status = absl::OkStatus();
        } else {
          lambda_status = absl::InternalError("Failed to parse response");
        }
      });
  auto on_done_functor = std::make_unique<imw_on_done_functor_t>(
      [&notif](const char* unused_keyexpr) { notif.Notify(); });
  KVQuery query(std::move(reply_functor), std::move(on_done_functor));

  imw_query_options_t query_options{
      .timeout_ms = static_cast<uint64_t>(timeout / absl::Milliseconds(1))};
  imw_ret ret =
      Zenoh().imw_query(prefixed_name->c_str(), zenoh_query_static_callback,
                        zenoh_query_static_on_done, nullptr, 0,
                        query.GetContext(), &query_options);
  if (ret != IMW_OK) {
    return absl::InternalError(
        absl::StrFormat("Error getting a key, return code: %d", ret));
  }
  bool returned = notif.WaitForNotificationWithTimeout(timeout);
  if (!returned) {
    return absl::DeadlineExceededError("Timeout waiting for key");
  }
  if (!lambda_status.ok()) {
    return lambda_status;
  }

  return std::move(value);
}

absl::StatusOr<KVQuery> KeyValueStore::GetAll(absl::string_view keyexpr,
                                              const WildcardQueryConfig& config,
                                              KeyValueCallback callback,
                                              OnDoneCallback on_done) {
  INTR_RETURN_IF_ERROR(intrinsic::ValidZenohKey(keyexpr));
  INTR_ASSIGN_OR_RETURN(absl::StatusOr<std::string> prefixed_name,
                        ZenohHandle::add_key_prefix(keyexpr, key_prefix_));
  auto functor = std::make_unique<imw_callback_functor_t>(
      [callback = std::move(callback)](const char* key,
                                       const void* response_bytes,
                                       const size_t response_bytes_len) {
        auto value = std::make_unique<google::protobuf::Any>();
        value->ParseFromString(absl::string_view(
            static_cast<const char*>(response_bytes), response_bytes_len));
        callback(key, std::move(value));
      });
  auto on_done_functor = std::make_unique<imw_on_done_functor_t>(
      [on_done = std::move(on_done)](const char* keyexpr) {
        on_done(absl::string_view(keyexpr));
      });
  KVQuery query(std::move(functor), std::move(on_done_functor));
  imw_ret_t ret = Zenoh().imw_query(
      prefixed_name->c_str(), zenoh_query_static_callback,
      zenoh_query_static_on_done, nullptr, 0, query.GetContext(), nullptr);
  if (ret != IMW_OK) {
    return absl::InternalError(
        absl::StrFormat("Error getting a key, return code: %d", ret));
  }

  return std::move(query);
}

absl::StatusOr<std::vector<std::string>> KeyValueStore::ListAllKeys(
    absl::Duration timeout) {
  std::vector<std::string> keys;
  absl::string_view query_keyexpr = "**";
  INTR_RETURN_IF_ERROR(intrinsic::ValidZenohKey(query_keyexpr));
  INTR_ASSIGN_OR_RETURN(
      absl::StatusOr<std::string> prefixed_name,
      ZenohHandle::add_key_prefix(query_keyexpr, key_prefix_));
  absl::Notification notif;
  auto callback = std::make_unique<imw_callback_functor_t>(
      [&keys, &notif](const char* keyexpr, const void* unused_response_bytes,
                      const size_t unused_response_bytes_len) {
        if (notif.HasBeenNotified()) {
          return;
        }
        keys.push_back(keyexpr);
      });
  auto on_done_functor = std::make_unique<imw_on_done_functor_t>(
      [&notif](const char* unused_keyexpr) { notif.Notify(); });
  KVQuery query(std::move(callback), std::move(on_done_functor));
  imw_query_options_t query_options{
      .timeout_ms = static_cast<uint64_t>(timeout / absl::Milliseconds(1))};
  imw_ret ret =
      Zenoh().imw_query(prefixed_name->c_str(), zenoh_query_static_callback,
                        zenoh_query_static_on_done, nullptr, 0,
                        query.GetContext(), &query_options);
  if (ret != IMW_OK) {
    return absl::InternalError(
        absl::StrFormat("Error getting a key, return code: %d", ret));
  }
  notif.WaitForNotificationWithTimeout(timeout);
  return std::move(keys);
}

absl::Status KeyValueStore::Delete(absl::string_view key,
                                   const NamespaceConfig& config) {
  INTR_RETURN_IF_ERROR(intrinsic::ValidZenohKey(key));
  INTR_ASSIGN_OR_RETURN(absl::StatusOr<std::string> prefixed_name,
                        ZenohHandle::add_key_prefix(key, key_prefix_));
  imw_ret_t ret = Zenoh().imw_delete_keyexpr(prefixed_name->c_str());
  if (ret != IMW_OK) {
    return absl::InternalError(
        absl::StrFormat("Error deleting a key, return code: %d", ret));
  }
  return absl::OkStatus();
}

}  // namespace intrinsic
