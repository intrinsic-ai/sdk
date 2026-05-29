// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

#include <dlfcn.h>

#include <cstdlib>
#include <string>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "intrinsic/middleware/imw.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_helpers.h"

namespace intrinsic {

void zenoh_static_callback(const char* keyexpr, const void* blob,
                           const size_t blob_len, void* fptr) {
  (*static_cast<imw_callback_functor_t*>(fptr))(keyexpr, blob, blob_len);
}

void zenoh_query_static_callback(const char* keyexpr, const void* blob,
                                 const size_t blob_len, void* fptr) {
  QueryContext* query_context = static_cast<QueryContext*>(fptr);
  (*(query_context->callback))(keyexpr, blob, blob_len);
}

void zenoh_query_static_on_done(const char* keyexpr, void* fptr) {
  QueryContext* query_context = static_cast<QueryContext*>(fptr);
  (*(query_context->on_done))(keyexpr);
}

ZenohHandle* ZenohHandle::CreateZenohHandle() {
  auto* zenoh = new ZenohHandle();
  zenoh->Initialize();
  return zenoh;
}

void ZenohHandle::Initialize() {
  this->imw_init = ::intrinsic::imw_init;
  this->imw_fini = ::intrinsic::imw_fini;
  this->imw_create_publisher = ::intrinsic::imw_create_publisher;
  this->imw_destroy_publisher = ::intrinsic::imw_destroy_publisher;
  this->imw_publish = ::intrinsic::imw_publish;
  this->imw_publisher_has_matching_subscribers =
      ::intrinsic::imw_publisher_has_matching_subscribers;
  this->imw_create_subscription = ::intrinsic::imw_create_subscription;
  this->imw_destroy_subscription = ::intrinsic::imw_destroy_subscription;
  this->imw_keyexpr_includes = ::intrinsic::imw_keyexpr_includes;
  this->imw_keyexpr_intersects = ::intrinsic::imw_keyexpr_intersects;
  this->imw_keyexpr_is_canon = ::intrinsic::imw_keyexpr_is_canon;
  this->imw_version = ::intrinsic::imw_version;
  this->imw_create_queryable = ::intrinsic::imw_create_queryable;
  this->imw_destroy_queryable = ::intrinsic::imw_destroy_queryable;
  this->imw_queryable_reply = ::intrinsic::imw_queryable_reply;
  this->imw_query = ::intrinsic::imw_query;
  this->imw_set = ::intrinsic::imw_set;
  this->imw_delete_keyexpr = ::intrinsic::imw_delete_keyexpr;
}

absl::StatusOr<std::string> ZenohHandle::add_topic_prefix(
    absl::string_view topic) {
  if (topic.empty()) {
    return absl::InvalidArgumentError("Empty topic string");
  }

  absl::string_view topic_without_leading_slash = topic;
  if (topic[0] == '/') {
    topic_without_leading_slash.remove_prefix(1);
  }

  if (absl::StartsWith(topic_without_leading_slash, "interipc_ps")) {
    return std::string(topic_without_leading_slash);
  }

  return absl::StrCat("in/", topic_without_leading_slash);
}

absl::StatusOr<std::string> ZenohHandle::add_key_prefix(
    absl::string_view key, absl::string_view key_prefix) {
  if (key.empty()) {
    return absl::InvalidArgumentError("Empty key string");
  } else if (key[0] == '/') {
    return absl::InvalidArgumentError("Key can't start with /");
  } else {
    return absl::StrCat(key_prefix, "/", key);
  }
}

absl::StatusOr<std::string> ZenohHandle::remove_topic_prefix(
    absl::string_view topic) {
  if (topic.length() < 3) {
    return absl::InvalidArgumentError("Topic string too short");
  }
  topic.remove_prefix(2);
  return std::string(topic);
}

const ZenohHandle& Zenoh() {
  static auto* zenoh_handle = ZenohHandle::CreateZenohHandle();
  return *zenoh_handle;
}

}  // namespace intrinsic
