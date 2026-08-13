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

namespace intrinsic {

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

}  // namespace intrinsic
