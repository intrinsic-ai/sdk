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

#include "intrinsic/scene/scene_object/root_entity.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "internal/testing.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/util/proto/parse_text_proto.h"

namespace intrinsic {
namespace scene_object {
namespace {

using ::absl_testing::IsOkAndHolds;
using ::absl_testing::StatusIs;
using ::intrinsic::ParseTextProtoOrDie;
using ::intrinsic_proto::scene_object::v1::SceneObject;
using ::testing::HasSubstr;

TEST(FindRootEntityTest, ReturnsRootEntity) {
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "object"
    entities { name: "child" parent_name: "root" }
    entities { name: "root" }
    entities { name: "grandchild" parent_name: "child" }
  )pb");

  EXPECT_THAT(FindRootEntity(scene_object),
              IsOkAndHolds(&scene_object.entities(1)));
}

TEST(FindRootEntityTest, DoesNotValidateSceneObject) {
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "object"
    entities { name: "root" }
    entities { name: "orphan" parent_name: "missing" }
  )pb");

  EXPECT_THAT(FindRootEntity(scene_object),
              IsOkAndHolds(&scene_object.entities(0)));
}

TEST(FindRootEntityTest, FailsWithoutRootEntity) {
  const SceneObject empty = ParseTextProtoOrDie(R"pb(name: "empty")pb");
  EXPECT_THAT(FindRootEntity(empty),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("No root entity")));

  const SceneObject cycle = ParseTextProtoOrDie(R"pb(
    name: "cycle"
    entities { name: "a" parent_name: "b" }
    entities { name: "b" parent_name: "a" }
  )pb");
  EXPECT_THAT(FindRootEntity(cycle),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("No root entity")));
}

TEST(FindRootEntityTest, FailsWithMultipleRootEntities) {
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "object"
    entities { name: "root1" }
    entities { name: "root2" }
  )pb");

  EXPECT_THAT(FindRootEntity(scene_object),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Multiple root entities")));
}

TEST(FindRootEntityTest, FailsWithUnnamedRootEntity) {
  const SceneObject scene_object = ParseTextProtoOrDie(R"pb(
    name: "object"
    entities {}
  )pb");

  EXPECT_THAT(
      FindRootEntity(scene_object),
      StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("empty name")));
}

}  // namespace
}  // namespace scene_object
}  // namespace intrinsic
