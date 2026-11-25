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

#include "intrinsic/geometry/internal/mesh/io/check_gltf_binary_header.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstdint>
#include <string>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/strings/string_view.h"
namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::HasSubstr;

// b/568616441: a 93-byte glTF 1.0 binary whose header declares a 3.7 GB scene.
constexpr absl::string_view kGltfV1HugeSceneLength(
    "glTF\001\000\000\000\376\377\000\000\376c\377\336\000\000\000"
    "\000\000\000QQQQQQH\303HHHHH\254HHHHHHHH\321HHHHaskHHHHH\224"
    "\357I\246HHHHHHHHHHHHHHHHHHHHHHHHHLHHHHHH",
    93);

// A glTF 2.0 binary header followed by a JSON chunk header with the given
// little-endian chunk length, then `json`.
std::string GltfV2(uint32_t json_length, const std::string& json) {
  std::string out("glTF\002\000\000\000\000\000\000\000", 12);
  for (int i = 0; i < 4; ++i) {
    out.push_back(static_cast<char>((json_length >> (8 * i)) & 0xFF));
  }
  out.append("JSON");
  out.append(json);
  return out;
}

TEST(CheckGltfBinaryJsonLength,
     ReturnsInvalidArgumentWhenGltfV1SceneLengthExceedsContent) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(kGltfV1HugeSceneLength, "gltf"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Binary glTF JSON length 3741279230")));
}

TEST(CheckGltfBinaryJsonLength, ReturnsInvalidArgumentWhenExtensionIsGlb) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(kGltfV1HugeSceneLength, "glb"),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(CheckGltfBinaryJsonLength,
     ReturnsInvalidArgumentWhenExtensionIsUpperCase) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(kGltfV1HugeSceneLength, "GLB"),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(CheckGltfBinaryJsonLength,
     ReturnsInvalidArgumentWhenGltfV2JsonChunkLengthExceedsContent) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(GltfV2(0xF0000000u, "{}"), "glb"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Binary glTF JSON length 4026531840")));
}

TEST(CheckGltfBinaryJsonLength, ReturnsInvalidArgumentWhenExtensionIsVrm) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(GltfV2(0xF0000000u, "{}"), "vrm"),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(CheckGltfBinaryJsonLength,
     ReturnsInvalidArgumentWhenExtensionHasLeadingDot) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(kGltfV1HugeSceneLength, ".glb"),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(CheckGltfBinaryJsonLength, ReturnsOkWhenJsonLengthMatchesContent) {
  const std::string glb = GltfV2(2, "{}");

  EXPECT_THAT(CheckGltfBinaryJsonLength(glb, "glb"), IsOk());
}

TEST(CheckGltfBinaryJsonLength, ReturnsOkWhenExtensionIsNotGltf) {
  EXPECT_THAT(CheckGltfBinaryJsonLength(kGltfV1HugeSceneLength, "stl"), IsOk());
}

TEST(CheckGltfBinaryJsonLength, ReturnsOkWhenContentIsNotBinaryGltf) {
  EXPECT_THAT(
      CheckGltfBinaryJsonLength(R"({"asset": {"version": "2.0"}})", "gltf"),
      IsOk());
}

TEST(CheckGltfBinaryJsonLength, ReturnsOkWhenContentIsShorterThanHeader) {
  EXPECT_THAT(CheckGltfBinaryJsonLength("glTF\001", "glb"), IsOk());
}

}  // namespace
}  // namespace intrinsic::geo
