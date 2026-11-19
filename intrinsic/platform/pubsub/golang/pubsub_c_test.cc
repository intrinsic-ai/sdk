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

struct StubState {
  absl::Mutex mu;
  imw_ret_t create_sub_ret = IMW_OK;
  imw_ret_t destroy_sub_ret = IMW_OK;
  imw_subscription_message_callback_fn* sub_callback = nullptr;
  void* sub_user_context = nullptr;
  int destroy_sub_calls = 0;

  std::vector<RecordedMessage> messages ABSL_GUARDED_BY(mu);
  std::vector<imw_encoding_t> set_encodings ABSL_GUARDED_BY(mu);
  std::vector<imw_encoding_t> reply_encodings ABSL_GUARDED_BY(mu);
};

StubState* g_stub_state = nullptr;

class ScopedStubState {
 public:
  explicit ScopedStubState(StubState* state) { g_stub_state = state; }
  ScopedStubState(const ScopedStubState&) = delete;
  ScopedStubState& operator=(const ScopedStubState&) = delete;
  ~ScopedStubState() { g_stub_state = nullptr; }
};

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

class PubsubCTest : public Test {
 protected:
  StubState state_;
  ScopedStubState scoped_{&state_};
  ZenohHandle zenoh_handle_ = MakeStubZenohHandle();
};

TEST_F(PubsubCTest,
       CreateAndDestroyMessageSubscriptionForwardsMetadataAndNullFallback) {
  constexpr uintptr_t kUserContext = 0x123456789ABCDEF0ULL;
  void* sub_ctx = nullptr;
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
}

}  // namespace
