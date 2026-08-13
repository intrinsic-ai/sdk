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

#include "intrinsic/util/thread/stop_token.h"

#include <gtest/gtest.h>

namespace intrinsic {
namespace {

// StopSource tests

TEST(StopSourceTest, StopSourceWithNoState) {
  StopSource stop_source(detail::NoState);
  EXPECT_FALSE(stop_source.stop_possible());
  EXPECT_FALSE(stop_source.stop_requested());
}

TEST(StopSourceTest, StopSourceWithState) {
  StopSource stop_source;
  EXPECT_TRUE(stop_source.stop_possible());
  EXPECT_FALSE(stop_source.stop_requested());
}

TEST(StopSourceTest, StopSourceCanBeStopped) {
  StopSource stop_source;
  EXPECT_TRUE(stop_source.request_stop());
  EXPECT_TRUE(stop_source.stop_requested());
}

TEST(StopSourceTest, StopSourceWithoutStateCannotBeStopped) {
  StopSource stop_source(detail::NoState);
  EXPECT_FALSE(stop_source.stop_possible());
  EXPECT_FALSE(stop_source.request_stop());
}

// StopToken tests

TEST(StopTokenTest, StopSourceWithNoState) {
  StopToken stop_token;
  EXPECT_FALSE(stop_token.stop_possible());
  EXPECT_FALSE(stop_token.stop_requested());
}

TEST(StopTokenTest, StopSourceWithStateCanStop) {
  StopSource stop_source;
  StopToken stop_token = stop_source.get_token();
  EXPECT_TRUE(stop_token.stop_possible());
  EXPECT_FALSE(stop_token.stop_requested());
}

TEST(StopTokenTest, StopSourceWithStateStops) {
  StopSource stop_source;
  EXPECT_TRUE(stop_source.request_stop());
  StopToken stop_token = stop_source.get_token();
  EXPECT_TRUE(stop_token.stop_possible());
  EXPECT_TRUE(stop_token.stop_requested());
}

TEST(StopTokenTest, StopSourceWithStateStopsAfterTokenCreation) {
  StopSource stop_source;
  StopToken stop_token = stop_source.get_token();
  EXPECT_TRUE(stop_token.stop_possible());
  EXPECT_FALSE(stop_token.stop_requested());
  EXPECT_TRUE(stop_source.request_stop());
  EXPECT_TRUE(stop_token.stop_requested());
}

}  // namespace
}  // namespace intrinsic
