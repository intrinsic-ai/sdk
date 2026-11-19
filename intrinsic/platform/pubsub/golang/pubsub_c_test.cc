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

#include "intrinsic/platform/pubsub/golang/pubsub_c.h"

#include <gtest/gtest.h>

#include <cstddef>
#include <cstdint>
#include <string>
#include <thread>
#include <vector>

#include "absl/base/thread_annotations.h"
#include "absl/synchronization/mutex.h"
#include "intrinsic/middleware/imw.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

namespace {

using ::intrinsic::IMW_ENCODING_APPLICATION_OCTET_STREAM;
using ::intrinsic::IMW_ENCODING_APPLICATION_PROTOBUF;
using ::intrinsic::imw_encoding_t;
using ::intrinsic::IMW_ENCODING_UNKNOWN;
using ::intrinsic::IMW_ENCODING_VARINT_PREFIXED_BYTES;
using ::intrinsic::IMW_ERROR;
using ::intrinsic::IMW_MESSAGE_KIND_DELETE;
using ::intrinsic::IMW_MESSAGE_KIND_PUT;
using ::intrinsic::imw_message_metadata_t;
using ::intrinsic::IMW_OK;
using ::intrinsic::imw_query_message_callback_fn;
using ::intrinsic::imw_query_on_done_callback_fn;
using ::intrinsic::imw_query_options_t;
using ::intrinsic::imw_ret_t;
using ::intrinsic::imw_subscription_message_callback_fn;
using ::intrinsic::ZenohHandle;
using ::testing::Test;

struct RecordedMessage {
  std::string keyexpr;
  std::string payload;
  int message_kind = 0;
  int encoding = 0;
  uintptr_t user_context = 0;
  bool operator==(const RecordedMessage&) const = default;
};

struct RecordedOnDone {
  std::string keyexpr;
  uintptr_t user_context = 0;
  bool operator==(const RecordedOnDone&) const = default;
};

struct StubState {
  absl::Mutex mu;
  imw_ret_t create_sub_ret = IMW_OK;
  imw_ret_t destroy_sub_ret = IMW_OK;
  imw_subscription_message_callback_fn* sub_callback = nullptr;
  void* sub_user_context = nullptr;
  int destroy_sub_calls = 0;

  imw_ret_t query_messages_ret = IMW_OK;
  bool query_invoke_message_sync = false;
  bool query_invoke_on_done_sync = false;
  bool query_spawn_concurrent_thread = false;
  std::string query_sync_keyexpr = "kv_store/sync_key";
  std::string query_sync_payload = "sync_payload";
  std::thread concurrent_thread;

  int delete_keyexpr_calls ABSL_GUARDED_BY(mu) = 0;
  std::string deleted_keyexpr ABSL_GUARDED_BY(mu);
  std::string last_query_keyexpr ABSL_GUARDED_BY(mu);
  uint64_t last_query_timeout_ms ABSL_GUARDED_BY(mu) = 0;
  std::vector<RecordedMessage> messages ABSL_GUARDED_BY(mu);
  std::vector<RecordedOnDone> query_on_dones ABSL_GUARDED_BY(mu);
  std::vector<imw_encoding_t> set_encodings ABSL_GUARDED_BY(mu);
  std::vector<imw_encoding_t> reply_encodings ABSL_GUARDED_BY(mu);
};

StubState* g_stub_state = nullptr;

class ScopedStubState {
 public:
  explicit ScopedStubState(StubState* state) { g_stub_state = state; }
  ScopedStubState(const ScopedStubState&) = delete;
  ScopedStubState& operator=(const ScopedStubState&) = delete;
  ~ScopedStubState() {
    if (g_stub_state != nullptr && g_stub_state->concurrent_thread.joinable()) {
      g_stub_state->concurrent_thread.join();
    }
    g_stub_state = nullptr;
  }
};

imw_ret_t StubQueryMessages(const char* query_keyexpr,
                            imw_query_message_callback_fn* callback,
                            imw_query_on_done_callback_fn* on_done,
                            const void* /*query_payload*/,
                            const size_t /*query_payload_len*/,
                            void* user_context, imw_query_options_t* options) {
  imw_ret_t ret = IMW_OK;
  bool message_sync = false, on_done_sync = false, spawn_concurrent = false;
  std::string keyexpr, payload;
  {
    absl::MutexLock lock(&g_stub_state->mu);
    g_stub_state->last_query_keyexpr =
        query_keyexpr != nullptr ? query_keyexpr : "";
    g_stub_state->last_query_timeout_ms =
        options != nullptr ? options->timeout_ms : 0;
    ret = g_stub_state->query_messages_ret;
    message_sync = g_stub_state->query_invoke_message_sync;
    on_done_sync = g_stub_state->query_invoke_on_done_sync;
    spawn_concurrent = g_stub_state->query_spawn_concurrent_thread;
    keyexpr = g_stub_state->query_sync_keyexpr;
    payload = g_stub_state->query_sync_payload;
  }

  constexpr imw_message_metadata_t kProtoMeta{
      .kind = IMW_MESSAGE_KIND_PUT,
      .encoding = IMW_ENCODING_APPLICATION_PROTOBUF};
  if (message_sync) {
    callback(keyexpr.c_str(), payload.data(), payload.size(), &kProtoMeta,
             user_context);
  }
  if (on_done_sync) {
    on_done(keyexpr.c_str(), user_context);
  }
  if (spawn_concurrent) {
    absl::MutexLock lock(&g_stub_state->mu);
    g_stub_state->concurrent_thread = std::thread(
        [callback, on_done, user_context, keyexpr, payload, kProtoMeta]() {
          callback(keyexpr.c_str(), payload.data(), payload.size(), &kProtoMeta,
                   user_context);
          on_done(keyexpr.c_str(), user_context);
        });
  }
  return ret;
}

ZenohHandle MakeStubZenohHandle() {
  ZenohHandle handle{};
  handle.imw_create_message_subscription =
      [](const char*, imw_subscription_message_callback_fn* cb, const char*,
         void* user_ctx) {
        absl::MutexLock lock(&g_stub_state->mu);
        g_stub_state->sub_callback = cb;
        g_stub_state->sub_user_context = user_ctx;
        return g_stub_state->create_sub_ret;
      };
  handle.imw_destroy_message_subscription =
      [](const char*, imw_subscription_message_callback_fn* cb,
         const void* user_ctx) {
        absl::MutexLock lock(&g_stub_state->mu);
        ++g_stub_state->destroy_sub_calls;
        EXPECT_EQ(cb, g_stub_state->sub_callback);
        EXPECT_EQ(user_ctx, g_stub_state->sub_user_context);
        return g_stub_state->destroy_sub_ret;
      };
  handle.imw_query_messages = StubQueryMessages;
  handle.imw_delete_keyexpr = [](const char* keyexpr) {
    absl::MutexLock lock(&g_stub_state->mu);
    ++g_stub_state->delete_keyexpr_calls;
    g_stub_state->deleted_keyexpr = keyexpr != nullptr ? keyexpr : "";
    return IMW_OK;
  };
  handle.imw_set_with_encoding = [](const char*, const void*, const size_t,
                                    imw_encoding_t enc) {
    absl::MutexLock lock(&g_stub_state->mu);
    g_stub_state->set_encodings.push_back(enc);
    return IMW_OK;
  };
  handle.imw_queryable_reply_with_encoding = [](const void*, const char*,
                                                const void*, const size_t,
                                                imw_encoding_t enc) {
    absl::MutexLock lock(&g_stub_state->mu);
    g_stub_state->reply_encodings.push_back(enc);
    return IMW_OK;
  };
  return handle;
}

void RecordMessage(const char* keyexpr, const void* bytes,
                   const size_t bytes_len, int message_kind, int encoding,
                   uintptr_t user_context) {
  absl::MutexLock lock(&g_stub_state->mu);
  g_stub_state->messages.push_back(RecordedMessage{
      .keyexpr = keyexpr != nullptr ? keyexpr : "",
      .payload = (bytes != nullptr && bytes_len > 0)
                     ? std::string(static_cast<const char*>(bytes), bytes_len)
                     : "",
      .message_kind = message_kind,
      .encoding = encoding,
      .user_context = user_context,
  });
}

void OnQueryDone(const char* keyexpr, uintptr_t user_context) {
  absl::MutexLock lock(&g_stub_state->mu);
  g_stub_state->query_on_dones.push_back(
      RecordedOnDone{.keyexpr = keyexpr != nullptr ? keyexpr : "",
                     .user_context = user_context});
}

class PubsubCTest : public Test {
 protected:
  int QueryMessages(const char* keyexpr, uintptr_t ctx = 0x42U) {
    return ZenohHandleImwQueryMessages(&zenoh_handle_, keyexpr, RecordMessage,
                                       OnQueryDone, ctx, /*timeout_ms=*/100);
  }

  StubState state_;
  ScopedStubState scoped_{&state_};
  ZenohHandle zenoh_handle_ = MakeStubZenohHandle();
};

TEST_F(PubsubCTest,
       CreateAndDestroyMessageSubscriptionForwardsMetadataAndNullFallback) {
  constexpr uintptr_t kUserContext = 0x123456789ABCDEF0ULL;
  void* sub_ctx = reinterpret_cast<void*>(static_cast<uintptr_t>(0x1U));
  state_.create_sub_ret = IMW_ERROR;
  EXPECT_EQ(ZenohHandleImwCreateMessageSubscription(
                &zenoh_handle_, "kv_store/fail", RecordMessage, "",
                kUserContext, &sub_ctx),
            IMW_ERROR);
  EXPECT_EQ(sub_ctx, nullptr);

  state_.create_sub_ret = IMW_OK;
  ASSERT_EQ(ZenohHandleImwCreateMessageSubscription(
                &zenoh_handle_, "kv_store/sub", RecordMessage, "", kUserContext,
                &sub_ctx),
            IMW_OK);
  ASSERT_NE(sub_ctx, nullptr);
  ASSERT_NE(state_.sub_callback, nullptr);

  const char delete_payload[] = "deleted";
  const imw_message_metadata_t delete_meta{
      .kind = IMW_MESSAGE_KIND_DELETE,
      .encoding = IMW_ENCODING_APPLICATION_OCTET_STREAM,
  };
  state_.sub_callback("kv_store/sub", delete_payload,
                      sizeof(delete_payload) - 1, &delete_meta,
                      state_.sub_user_context);
  const imw_message_metadata_t unknown_meta{
      .kind = IMW_MESSAGE_KIND_PUT,
      .encoding = IMW_ENCODING_UNKNOWN,
  };
  state_.sub_callback("kv_store/sub", "{}", 2, &unknown_meta,
                      state_.sub_user_context);
  state_.sub_callback("kv_store/sub", nullptr, 0, nullptr,
                      state_.sub_user_context);

  {
    absl::MutexLock lock(&state_.mu);
    EXPECT_EQ(
        state_.messages,
        (std::vector<RecordedMessage>{
            {"kv_store/sub", "deleted", ZENOH_HANDLE_IMW_MESSAGE_KIND_DELETE,
             ZENOH_HANDLE_IMW_ENCODING_APPLICATION_OCTET_STREAM, kUserContext},
            {"kv_store/sub", "{}", ZENOH_HANDLE_IMW_MESSAGE_KIND_PUT,
             ZENOH_HANDLE_IMW_ENCODING_UNKNOWN, kUserContext},
            {"kv_store/sub", "", ZENOH_HANDLE_IMW_MESSAGE_KIND_PUT,
             ZENOH_HANDLE_IMW_ENCODING_VARINT_PREFIXED_BYTES, kUserContext},
        }));
  }

  state_.destroy_sub_ret = IMW_ERROR;
  EXPECT_EQ(ZenohHandleImwDestroyMessageSubscription(&zenoh_handle_,
                                                     "kv_store/sub", sub_ctx),
            IMW_ERROR);
  state_.destroy_sub_ret = IMW_OK;
  EXPECT_EQ(ZenohHandleImwDestroyMessageSubscription(&zenoh_handle_,
                                                     "kv_store/sub", sub_ctx),
            IMW_OK);
  EXPECT_EQ(state_.destroy_sub_calls, 2);
}

// MessageQueryContext single-free on error paths is verified by ASan/LSan
// alongside query_on_dones.empty().
TEST_F(PubsubCTest,
       QueryMessagesErrorWithoutOnDoneFreesContextWithoutCallback) {
  state_.query_messages_ret = IMW_ERROR;
  EXPECT_EQ(QueryMessages("kv_store/err", 0xCAFEBABEULL), IMW_ERROR);
  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.last_query_keyexpr, "kv_store/err");
  EXPECT_EQ(state_.last_query_timeout_ms, 100ULL);
  EXPECT_TRUE(state_.query_on_dones.empty());
}

TEST_F(PubsubCTest,
       QueryMessagesSynchronousDropOnErrorSuppressesOnDoneAndFreesOnce) {
  state_.query_messages_ret = IMW_ERROR;
  state_.query_invoke_on_done_sync = true;
  EXPECT_EQ(QueryMessages("kv_store/sync_err", 0xDEADBEEFULL), IMW_ERROR);
  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.last_query_keyexpr, "kv_store/sync_err");
  EXPECT_EQ(state_.last_query_timeout_ms, 100ULL);
  EXPECT_TRUE(state_.query_on_dones.empty());
}

TEST_F(PubsubCTest, QueryMessagesSynchronousDropOnOkReplaysOnDone) {
  constexpr uintptr_t kUserContext = 0x1020304050607080ULL;
  state_.query_invoke_message_sync = true;
  state_.query_invoke_on_done_sync = true;
  EXPECT_EQ(QueryMessages("kv_store/sync_key", kUserContext), IMW_OK);

  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.last_query_keyexpr, "kv_store/sync_key");
  EXPECT_EQ(state_.last_query_timeout_ms, 100ULL);
  EXPECT_EQ(state_.messages,
            (std::vector<RecordedMessage>{
                {"kv_store/sync_key", "sync_payload",
                 ZENOH_HANDLE_IMW_MESSAGE_KIND_PUT,
                 ZENOH_HANDLE_IMW_ENCODING_APPLICATION_PROTOBUF, kUserContext},
            }));
  EXPECT_EQ(state_.query_on_dones,
            (std::vector<RecordedOnDone>{{"kv_store/sync_key", kUserContext}}));
}

TEST_F(PubsubCTest, QueryMessagesConcurrentOnDoneThreadRacesCleanly) {
  constexpr uintptr_t kUserContext = 0x8877665544332211ULL;
  state_.query_spawn_concurrent_thread = true;
  state_.query_sync_keyexpr = "kv_store/concurrent_key";
  state_.query_sync_payload = "concurrent_payload";

  EXPECT_EQ(QueryMessages("kv_store/concurrent_key", kUserContext), IMW_OK);
  if (state_.concurrent_thread.joinable()) {
    state_.concurrent_thread.join();
  }

  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.last_query_keyexpr, "kv_store/concurrent_key");
  EXPECT_EQ(state_.last_query_timeout_ms, 100ULL);
  EXPECT_EQ(state_.messages,
            (std::vector<RecordedMessage>{
                {"kv_store/concurrent_key", "concurrent_payload",
                 ZENOH_HANDLE_IMW_MESSAGE_KIND_PUT,
                 ZENOH_HANDLE_IMW_ENCODING_APPLICATION_PROTOBUF, kUserContext},
            }));
  EXPECT_EQ(
      state_.query_on_dones,
      (std::vector<RecordedOnDone>{{"kv_store/concurrent_key", kUserContext}}));
}

TEST_F(PubsubCTest, EncodingBridgeValidatesNullHandlesAndUnknownEncodings) {
  const char payload[] = "abc";
  struct Case {
    void* handle;
    int encoding;
    imw_ret_t want;
  };
  const Case cases[] = {
      {nullptr, ZENOH_HANDLE_IMW_ENCODING_VARINT_PREFIXED_BYTES, IMW_ERROR},
      {&zenoh_handle_, -1, IMW_ERROR},
      {&zenoh_handle_, ZENOH_HANDLE_IMW_ENCODING_UNKNOWN, IMW_ERROR},
      {&zenoh_handle_, ZENOH_HANDLE_IMW_ENCODING_VARINT_PREFIXED_BYTES, IMW_OK},
      {&zenoh_handle_, ZENOH_HANDLE_IMW_ENCODING_APPLICATION_OCTET_STREAM,
       IMW_OK},
      {&zenoh_handle_, ZENOH_HANDLE_IMW_ENCODING_APPLICATION_PROTOBUF, IMW_OK},
  };
  for (const Case& tc : cases) {
    EXPECT_EQ(ZenohHandleImwSetWithEncoding(tc.handle, "kv_store/k", payload,
                                            sizeof(payload) - 1, tc.encoding),
              tc.want);
    EXPECT_EQ(ZenohHandleImwQueryableReplyWithEncoding(
                  tc.handle, nullptr, "kv_store/k", payload,
                  sizeof(payload) - 1, tc.encoding),
              tc.want);
  }

  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.set_encodings,
            (std::vector<imw_encoding_t>{IMW_ENCODING_VARINT_PREFIXED_BYTES,
                                         IMW_ENCODING_APPLICATION_OCTET_STREAM,
                                         IMW_ENCODING_APPLICATION_PROTOBUF}));
  EXPECT_EQ(state_.reply_encodings, state_.set_encodings);
}

TEST_F(PubsubCTest, MessageBridgeValidatesNullArguments) {
  void* sub_ctx = reinterpret_cast<void*>(static_cast<uintptr_t>(0x99U));
  int dummy_ctx = 0;

  EXPECT_EQ(ZenohHandleImwCreateMessageSubscription(
                nullptr, "kv_store/sub", RecordMessage, "", 0x42U, &sub_ctx),
            IMW_ERROR);
  EXPECT_EQ(ZenohHandleImwCreateMessageSubscription(
                &zenoh_handle_, "kv_store/sub", nullptr, "", 0x42U, &sub_ctx),
            IMW_ERROR);
  EXPECT_EQ(
      ZenohHandleImwCreateMessageSubscription(
          &zenoh_handle_, "kv_store/sub", RecordMessage, "", 0x42U, nullptr),
      IMW_ERROR);
  EXPECT_EQ(state_.sub_callback, nullptr);

  EXPECT_EQ(ZenohHandleImwDestroyMessageSubscription(nullptr, "kv_store/sub",
                                                     &dummy_ctx),
            IMW_ERROR);
  EXPECT_EQ(ZenohHandleImwDestroyMessageSubscription(&zenoh_handle_,
                                                     "kv_store/sub", nullptr),
            IMW_ERROR);
  EXPECT_EQ(state_.destroy_sub_calls, 0);

  EXPECT_EQ(ZenohHandleImwQueryMessages(nullptr, "kv_store/q", RecordMessage,
                                        OnQueryDone, 0x42U, 100),
            IMW_ERROR);
  EXPECT_EQ(QueryMessages(nullptr), IMW_ERROR);
  EXPECT_EQ(ZenohHandleImwQueryMessages(&zenoh_handle_, "kv_store/q", nullptr,
                                        OnQueryDone, 0x42U, 100),
            IMW_ERROR);
  absl::MutexLock lock(&state_.mu);
  EXPECT_TRUE(state_.query_on_dones.empty());
}

TEST_F(PubsubCTest, DeleteKeyExprAndNullQueryOnDoneAreHandledCleanly) {
  EXPECT_EQ(ZenohHandleImwDeleteKeyExpr(&zenoh_handle_, "kv_store/del_key"),
            IMW_OK);

  state_.query_invoke_message_sync = true;
  state_.query_invoke_on_done_sync = true;
  EXPECT_EQ(ZenohHandleImwQueryMessages(&zenoh_handle_, "kv_store/null_done",
                                        RecordMessage, nullptr, 0x55U, 250),
            IMW_OK);
  state_.query_invoke_message_sync = false;
  state_.query_invoke_on_done_sync = false;
  state_.query_messages_ret = IMW_ERROR;
  EXPECT_EQ(ZenohHandleImwQueryMessages(&zenoh_handle_, "kv_store/null_done_e",
                                        RecordMessage, nullptr, 0x66U, 350),
            IMW_ERROR);

  absl::MutexLock lock(&state_.mu);
  EXPECT_EQ(state_.delete_keyexpr_calls, 1);
  EXPECT_EQ(state_.deleted_keyexpr, "kv_store/del_key");
  EXPECT_EQ(state_.last_query_keyexpr, "kv_store/null_done_e");
  EXPECT_EQ(state_.last_query_timeout_ms, 350ULL);
  EXPECT_TRUE(state_.query_on_dones.empty());
  EXPECT_EQ(state_.messages,
            (std::vector<RecordedMessage>{
                {"kv_store/sync_key", "sync_payload",
                 ZENOH_HANDLE_IMW_MESSAGE_KIND_PUT,
                 ZENOH_HANDLE_IMW_ENCODING_APPLICATION_PROTOBUF, 0x55U},
            }));
}

}  // namespace
