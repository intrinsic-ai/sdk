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

#include "intrinsic/geometry/gltf/gltf_io.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cstdint>
#include <cstring>
#include <sstream>
#include <string>

#include "absl/base/log_severity.h"
#include "absl/log/scoped_mock_log.h"
#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "internal/testing.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::_;
using ::testing::AllOf;
using ::testing::Eq;
using ::testing::Gt;
using ::testing::HasSubstr;
using ::testing::IsEmpty;
using ::testing::Not;

// Returns GLB bytes of a model with a single triangle mesh and one material,
// serialized directly with tinygltf.
std::string MakeTriangleGlb() {
  constexpr std::array<float, 9> kVertices = {1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
                                              0.0f, 0.0f, 0.0f, 1.0f};
  tinygltf::Model model;
  model.asset.version = "2.0";

  tinygltf::Buffer& buffer = model.buffers.emplace_back();
  buffer.data.resize(sizeof(kVertices));
  std::memcpy(buffer.data.data(), kVertices.data(), sizeof(kVertices));

  tinygltf::BufferView& view = model.bufferViews.emplace_back();
  view.buffer = 0;
  view.byteLength = sizeof(kVertices);
  view.target = TINYGLTF_TARGET_ARRAY_BUFFER;

  tinygltf::Accessor& accessor = model.accessors.emplace_back();
  accessor.bufferView = 0;
  accessor.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  accessor.count = 3;
  accessor.type = TINYGLTF_TYPE_VEC3;
  accessor.minValues = {0.0, 0.0, 0.0};
  accessor.maxValues = {1.0, 1.0, 1.0};

  model.materials.emplace_back().name = "material";

  tinygltf::Primitive& primitive =
      model.meshes.emplace_back().primitives.emplace_back();
  primitive.attributes["POSITION"] = 0;
  primitive.material = 0;
  primitive.mode = TINYGLTF_MODE_TRIANGLES;

  model.nodes.emplace_back().mesh = 0;
  model.scenes.emplace_back().nodes = {0};
  model.defaultScene = 0;

  std::ostringstream out;
  tinygltf::TinyGLTF writer;
  EXPECT_TRUE(writer.WriteGltfSceneToStream(&model, out,
                                            /*prettyPrint=*/false,
                                            /*writeBinary=*/true));
  return out.str();
}

// Wraps glTF JSON in a GLB container without a BIN chunk, so tests can
// describe models as JSON.
std::string JsonToGlb(absl::string_view json) {
  std::string chunk(json);
  chunk.resize((chunk.size() + 3) / 4 * 4, ' ');
  const auto u32 = [](uint32_t v) {
    return std::string(reinterpret_cast<const char*>(&v), sizeof(v));
  };
  constexpr uint32_t kGlbVersion = 2;
  constexpr uint32_t kHeaderAndChunkHeaderSize = 20;
  return absl::StrCat(
      "glTF", u32(kGlbVersion),
      u32(static_cast<uint32_t>(kHeaderAndChunkHeaderSize + chunk.size())),
      u32(static_cast<uint32_t>(chunk.size())), "JSON", chunk);
}

TEST(GltfIoTest, ParseGlbAndSerializeGltfToGlbRoundtrip) {
  const std::string glb = MakeTriangleGlb();
  ASSERT_OK_AND_ASSIGN(const tinygltf::Model model, ParseGlb(glb));
  EXPECT_THAT(model.meshes.size(), Gt(0));

  ASSERT_OK_AND_ASSIGN(const std::string serialized_glb,
                       SerializeGltfToGlb(model));
  EXPECT_THAT(serialized_glb, Not(IsEmpty()));

  ASSERT_OK_AND_ASSIGN(const tinygltf::Model reparsed_model,
                       ParseGlb(serialized_glb));
  EXPECT_THAT(reparsed_model.meshes.size(), Eq(model.meshes.size()));
  EXPECT_THAT(reparsed_model.materials.size(), Eq(model.materials.size()));
}

TEST(GltfIoTest, ParseGlbReportsEmptyDiagnosticsForValidGlb) {
  GltfParseDiagnostics diagnostics{.errors = "stale", .warnings = "stale"};
  EXPECT_THAT(ParseGlb(MakeTriangleGlb(), &diagnostics), IsOk());
  EXPECT_THAT(diagnostics.errors, IsEmpty());
  EXPECT_THAT(diagnostics.warnings, IsEmpty());
}

TEST(GltfIoTest, ParseGlbReportsErrorsInDiagnosticsOnFailure) {
  // Valid magic followed by a truncated header.
  GltfParseDiagnostics diagnostics;
  EXPECT_THAT(ParseGlb("glTF", &diagnostics),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(diagnostics.errors, Not(IsEmpty()));
}

TEST(GltfIoTest, ParseGlbReportsWarningsInDiagnostics) {
  // An unresolvable external image is reported as a warning, not an error.
  GltfParseDiagnostics diagnostics;
  EXPECT_THAT(ParseGlb(JsonToGlb(R"json({
    "asset": {"version": "2.0"},
    "images": [{"uri": "does_not_exist.png"}]
  })json"),
                       &diagnostics),
              IsOk());
  EXPECT_THAT(diagnostics.errors, IsEmpty());
  EXPECT_THAT(diagnostics.warnings, HasSubstr("does_not_exist.png"));
}

TEST(GltfIoTest, ParseGlbReportsJsonErrorsInDiagnostics) {
  GltfParseDiagnostics diagnostics;
  EXPECT_THAT(ParseGlb(JsonToGlb(R"json({"scenes": []})json"), &diagnostics),
              StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("asset")));
  EXPECT_THAT(diagnostics.errors, HasSubstr("asset"));
}

constexpr absl::string_view kGltfWithUnresolvableImage = R"json({
  "asset": {"version": "2.0"},
  "images": [{"uri": "does_not_exist.png"}]
})json";

TEST(GltfIoTest, ParseWithNullDiagnosticsLogsWarnings) {
  absl::ScopedMockLog log(absl::MockLogDefault::kIgnoreUnexpected);
  EXPECT_CALL(
      log, Log(absl::LogSeverity::kWarning, _, HasSubstr("does_not_exist.png")))
      .Times(1);
  log.StartCapturingLogs();
  EXPECT_THAT(ParseGlb(JsonToGlb(kGltfWithUnresolvableImage),
                       /*diagnostics=*/nullptr),
              IsOk());
}

TEST(GltfIoTest, ParseWithDiagnosticsDoesNotLogWarnings) {
  absl::ScopedMockLog log(absl::MockLogDefault::kIgnoreUnexpected);
  EXPECT_CALL(log, Log(absl::LogSeverity::kWarning, _, _)).Times(0);
  log.StartCapturingLogs();
  GltfParseDiagnostics diagnostics;
  EXPECT_THAT(ParseGlb(JsonToGlb(kGltfWithUnresolvableImage), &diagnostics),
              IsOk());
  EXPECT_THAT(diagnostics.warnings, HasSubstr("does_not_exist.png"));
}

TEST(GltfIoTest, ParseWithNullDiagnosticsHandlesAllPaths) {
  EXPECT_THAT(ParseGlb(MakeTriangleGlb(), /*diagnostics=*/nullptr), IsOk());
  EXPECT_THAT(ParseGlb("not a glb", /*diagnostics=*/nullptr),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Invalid GLB magic bytes")));
  EXPECT_THAT(ParseGlb("glTF", /*diagnostics=*/nullptr),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       HasSubstr("Failed to parse GLB")));
  EXPECT_THAT(ParseGlb(JsonToGlb(R"json({"scenes": []})json"),
                       /*diagnostics=*/nullptr),
              StatusIs(absl::StatusCode::kInvalidArgument, HasSubstr("asset")));
}

TEST(GltfIoTest, ParseGlbClearsDiagnosticsOnInvalidMagic) {
  GltfParseDiagnostics diagnostics{.errors = "stale", .warnings = "stale"};
  EXPECT_THAT(ParseGlb("not a glb", &diagnostics),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(diagnostics.errors, IsEmpty());
  EXPECT_THAT(diagnostics.warnings, IsEmpty());
}

TEST(GltfIoTest, ParseGlbRejectsJsonAndInvalidInput) {
  EXPECT_THAT(ParseGlb(R"json({"asset": {"version": "2.0"}})json"),
              StatusIs(absl::StatusCode::kInvalidArgument,
                       AllOf(HasSubstr("Failed to parse GLB"),
                             HasSubstr("Invalid GLB magic bytes"))));
  EXPECT_THAT(ParseGlb(""), StatusIs(absl::StatusCode::kInvalidArgument,
                                     HasSubstr("Failed to parse GLB")));
}

TEST(GltfIoTest, ParseGlbKeepsEncodedImageBytesWithoutDecoding) {
  // The image bytes are not a decodable image, which a decoding loader would
  // reject.
  constexpr absl::string_view kImageBytes = "fake";
  tinygltf::Model model;
  model.asset.version = "2.0";
  model.buffers.emplace_back().data.assign(kImageBytes.begin(),
                                           kImageBytes.end());
  tinygltf::BufferView& view = model.bufferViews.emplace_back();
  view.buffer = 0;
  view.byteLength = kImageBytes.size();
  tinygltf::Image& image = model.images.emplace_back();
  image.bufferView = 0;
  image.mimeType = "image/png";
  std::ostringstream out;
  tinygltf::TinyGLTF writer;
  ASSERT_TRUE(writer.WriteGltfSceneToStream(&model, out, /*prettyPrint=*/false,
                                            /*writeBinary=*/true));

  ASSERT_OK_AND_ASSIGN(const tinygltf::Model parsed, ParseGlb(out.str()));
  ASSERT_OK_AND_ASSIGN(const std::string serialized,
                       SerializeGltfToGlb(parsed));
  ASSERT_OK_AND_ASSIGN(const tinygltf::Model reparsed, ParseGlb(serialized));

  for (const tinygltf::Model& m : {parsed, reparsed}) {
    ASSERT_THAT(m.images.size(), Eq(1));
    EXPECT_THAT(m.images[0].image, IsEmpty());
    EXPECT_THAT(m.images[0].mimeType, Eq("image/png"));
    ASSERT_THAT(m.images[0].bufferView, Eq(0));
    const tinygltf::BufferView& bv = m.bufferViews[0];
    const auto begin = m.buffers[bv.buffer].data.begin() + bv.byteOffset;
    EXPECT_THAT(std::string(begin, begin + bv.byteLength), Eq(kImageBytes));
  }
}

}  // namespace
}  // namespace intrinsic::geo
