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

#include "intrinsic/geometry/internal/mesh/io/load_ai_scene_from_buffer.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::StatusIs;
using ::testing::HasSubstr;
TEST(LoadAiSceneFromBuffer,
     ReturnsInvalidArgumentWhenGltfJsonLengthExceedsBuffer) {
  // b/568616441: a glTF 1.0 binary header that declares a 3.7 GB scene.
  const std::string gltf(
      "glTF\001\000\000\000\376\377\000\000\376c\377\336\000\000\000\000"
      "{}",
      22);

  EXPECT_THAT(LoadAiSceneFromBuffer(gltf, "glb"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Binary glTF JSON length 3741279230")));
}

}  // namespace
}  // namespace intrinsic::geo
