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

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_handle.h"

#include <gtest/gtest.h>

#include <memory>
#include <string>
#include <vector>

#include "intrinsic/middleware/imw.h"

namespace intrinsic {
namespace {

struct CapturedMessage {
  bool called = false;
  std::string key;
  std::string payload;
  imw_message_metadata_t metadata{};

  imw_message_callback_functor_t MakeFunctor() {
    return [this](const char* message_key, const void* data, size_t len,
                  const imw_message_metadata_t* meta) {
      called = true;
      key = (message_key != nullptr) ? message_key : "";
      payload.assign(static_cast<const char*>(data), len);
      if (meta != nullptr) {
        metadata = *meta;
      }
    };
  }
};

}  // namespace

TEST(ZenohHandleTest, AddTopicPrefix) {
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("foo"), "in/foo");
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("/foo"), "in/foo");
  EXPECT_EQ(ZenohHandle::add_topic_prefix("").ok(), false);
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("interipc_something/tf"),
            "in/interipc_something/tf");
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("/interipc_something/tf"),
            "in/interipc_something/tf");
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("/interipc_ps/tf"),
            "interipc_ps/tf");
  EXPECT_EQ(*ZenohHandle::add_topic_prefix("interipc_ps/tf"), "interipc_ps/tf");
}

TEST(ZenohHandleTest, RemoveTopicPrefix) {
  EXPECT_EQ(*ZenohHandle::remove_topic_prefix("in/foo"), "/foo");
  EXPECT_EQ(*ZenohHandle::remove_topic_prefix("in/"), "/");
  EXPECT_EQ(ZenohHandle::remove_topic_prefix("").ok(), false);
}

TEST(ZenohHandleTest, AddKeyPrefix) {
  EXPECT_EQ(*ZenohHandle::add_key_prefix("my_key", "in/prefix"),
            "in/prefix/my_key");
  EXPECT_EQ(*ZenohHandle::add_key_prefix("nested/key", "root"),
            "root/nested/key");
  EXPECT_FALSE(ZenohHandle::add_key_prefix("", "in/prefix").ok());
  EXPECT_FALSE(ZenohHandle::add_key_prefix("/leading_slash", "in/prefix").ok());
}

TEST(ZenohHandleTest, EncodingFunctionPointersAreNonNull) {
  const ZenohHandle& zenoh = Zenoh();
  EXPECT_NE(zenoh.imw_create_message_subscription, nullptr);
  EXPECT_NE(zenoh.imw_destroy_message_subscription, nullptr);
  EXPECT_NE(zenoh.imw_queryable_reply_with_encoding, nullptr);
  EXPECT_NE(zenoh.imw_query_messages, nullptr);
  EXPECT_NE(zenoh.imw_set_with_encoding, nullptr);
  EXPECT_NE(zenoh.imw_set_with_encoding_and_deleter, nullptr);
}

TEST(ZenohHandleTest, ZenohStaticMessageCallbackDispatchesAllFields) {
  CapturedMessage captured;
  imw_message_callback_functor_t functor = captured.MakeFunctor();

  const char payload[] = "binary_message_data";
  imw_message_metadata_t meta{
      .kind = IMW_MESSAGE_KIND_PUT,
      .encoding = IMW_ENCODING_APPLICATION_OCTET_STREAM,
  };
  zenoh_static_message_callback("in/message_topic", payload,
                                sizeof(payload) - 1, &meta, &functor);

  ASSERT_TRUE(captured.called);
  EXPECT_EQ(captured.key, "in/message_topic");
  EXPECT_EQ(captured.payload, "binary_message_data");
  EXPECT_EQ(captured.metadata.kind, IMW_MESSAGE_KIND_PUT);
  EXPECT_EQ(captured.metadata.encoding, IMW_ENCODING_APPLICATION_OCTET_STREAM);
}

TEST(ZenohHandleTest, ZenohQueryMessageStaticCallbackDispatchesAllFields) {
  CapturedMessage captured;
  imw_message_callback_functor_t functor = captured.MakeFunctor();

  QueryContext ctx;
  ctx.message_callback = &functor;

  const char payload[] = "query_response_data";
  imw_message_metadata_t meta{
      .kind = IMW_MESSAGE_KIND_PUT,
      .encoding = IMW_ENCODING_APPLICATION_PROTOBUF,
  };
  zenoh_query_message_static_callback("in/query_reply", payload,
                                      sizeof(payload) - 1, &meta, &ctx);

  ASSERT_TRUE(captured.called);
  EXPECT_EQ(captured.key, "in/query_reply");
  EXPECT_EQ(captured.payload, "query_response_data");
  EXPECT_EQ(captured.metadata.kind, IMW_MESSAGE_KIND_PUT);
  EXPECT_EQ(captured.metadata.encoding, IMW_ENCODING_APPLICATION_PROTOBUF);

  // Null context is safe.
  zenoh_query_message_static_callback("in/null", nullptr, 0, nullptr, nullptr);
}

TEST(ZenohHandleTest, ZenohQueryStaticOnDoneRunsOnDoneThenOnDestroySelfDelete) {
  std::vector<std::string> execution_order;

  imw_on_done_functor_t on_done = [&](const char* keyexpr) {
    execution_order.push_back(std::string("on_done:") + keyexpr);
  };

  auto owned = std::make_unique<QueryContext>();
  QueryContext* context = owned.get();
  context->on_done = &on_done;
  context->on_destroy = [&owned, &execution_order]() {
    execution_order.push_back("on_destroy");
    owned.reset();
  };

  zenoh_query_static_on_done("in/query_done", context);
  EXPECT_EQ(owned, nullptr);

  ASSERT_EQ(execution_order.size(), 2);
  EXPECT_EQ(execution_order[0], "on_done:in/query_done");
  EXPECT_EQ(execution_order[1], "on_destroy");

  // Null context is safe.
  zenoh_query_static_on_done("in/query_done", nullptr);
}

TEST(ZenohHandleTest, ZenohQueryCallbacksDispatchAndSkipNullFunctors) {
  bool callback_called = false;
  std::string received_key;
  std::string received_payload;
  imw_callback_functor_t callback = [&](const char* key, const void* data,
                                        size_t len) {
    callback_called = true;
    received_key = (key != nullptr) ? key : "";
    received_payload.assign(static_cast<const char*>(data), len);
  };

  QueryContext ctx;
  ctx.callback = &callback;
  const char payload[] = "legacy_query_reply";
  zenoh_query_static_callback("in/legacy_query", payload, sizeof(payload) - 1,
                              &ctx);
  ASSERT_TRUE(callback_called);
  EXPECT_EQ(received_key, "in/legacy_query");
  EXPECT_EQ(received_payload, "legacy_query_reply");

  // Empty QueryContext (null callback / message_callback / on_done) and
  // nullptr.
  QueryContext empty_ctx;
  zenoh_query_static_callback("in/empty", payload, sizeof(payload) - 1,
                              &empty_ctx);
  zenoh_query_static_callback("in/null", nullptr, 0, nullptr);
  zenoh_query_message_static_callback("in/empty", nullptr, 0, nullptr,
                                      &empty_ctx);
  zenoh_query_static_on_done("in/empty", &empty_ctx);

  // Targetless imw_on_done_functor_t skips invocation while running on_destroy.
  imw_on_done_functor_t empty_on_done;
  int destroy_count = 0;
  QueryContext destroy_only_ctx;
  destroy_only_ctx.on_done = &empty_on_done;
  destroy_only_ctx.on_destroy = [&destroy_count]() { ++destroy_count; };
  zenoh_query_static_on_done("in/destroy_only", &destroy_only_ctx);
  EXPECT_EQ(destroy_count, 1);
}

}  // namespace intrinsic
