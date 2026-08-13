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

#ifndef INTRINSIC_UTIL_TESTING_GTEST_WRAPPER_H_
#define INTRINSIC_UTIL_TESTING_GTEST_WRAPPER_H_

#ifndef ASSERT_OK
#define ASSERT_OK(expr) ASSERT_THAT(expr, ::absl_testing::IsOk())
#endif

#ifndef EXPECT_OK
#define EXPECT_OK(expr) EXPECT_THAT(expr, ::absl_testing::IsOk())
#endif

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status_matchers.h"
#include "internal/testing.h"
#include "intrinsic/util/testing/status_payload_matchers.h"
#include "protobuf-matchers/protocol-buffer-matchers.h"

namespace intrinsic {
namespace testing {

using ::intrinsic::testing::StatusHasGenericPayload;
using ::intrinsic::testing::StatusHasProtoPayload;
using ::protobuf_matchers::EqualsProto;
using ::protobuf_matchers::EquivToProto;
using ::protobuf_matchers::internal::ProtoCompare;
using ::protobuf_matchers::internal::ProtoComparison;
using ::protobuf_matchers::proto::Approximately;
using ::protobuf_matchers::proto::IgnoringFieldPaths;
using ::protobuf_matchers::proto::IgnoringFields;
using ::protobuf_matchers::proto::IgnoringRepeatedFieldOrdering;
using ::protobuf_matchers::proto::Partially;
using ::protobuf_matchers::proto::TreatingNaNsAsEqual;
using ::protobuf_matchers::proto::WhenDeserialized;
using ::protobuf_matchers::proto::WhenDeserializedAs;

}  // namespace testing
}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_TESTING_GTEST_WRAPPER_H_
