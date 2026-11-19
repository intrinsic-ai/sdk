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

#include "intrinsic/geometry/gltf/gltf_buffer_reader.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cstdint>
#include <limits>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "internal/testing.h"
#include "intrinsic/geometry/gltf/gltf_accessor.h"
#include "intrinsic/geometry/gltf/gltf_buffer_test_util.h"
#include "intrinsic/geometry/gltf/gltf_buffer_writer.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::ElementsAre;
using ::testing::FloatNear;
using ::testing::IsEmpty;
using ::testing::SizeIs;

TEST(GltfBufferReaderTest, ReadsNormalizedAttributesAndRejectsNonFiniteValues) {
  tinygltf::Model model;
  model.buffers.emplace_back();

  const std::vector<uint8_t> colors = {255, 128, 0, 0, 0, 255, 64, 0};
  ASSERT_OK_AND_ASSIGN(
      const int col_bv,
      AppendGltfBufferView(model, 0, colors.data(), colors.size(),
                           TINYGLTF_TARGET_ARRAY_BUFFER, 4));
  tinygltf::Accessor col_acc;
  col_acc.bufferView = col_bv;
  col_acc.componentType = TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE;
  col_acc.normalized = true;
  col_acc.count = 2;
  col_acc.type = TINYGLTF_TYPE_VEC3;
  model.accessors.push_back(col_acc);

  tinygltf::Primitive prim;
  prim.attributes["COLOR_0"] = 0;

  ASSERT_OK_AND_ASSIGN(const std::vector<float> rgba,
                       ReadGltfAttributeFloats(model, prim, "COLOR_0", 4));
  ASSERT_THAT(rgba, SizeIs(8));
  EXPECT_THAT(rgba[0], FloatNear(1.0f, 1e-3f));
  EXPECT_THAT(rgba[1], FloatNear(128.0f / 255.0f, 1e-3f));
  EXPECT_THAT(rgba[2], FloatNear(0.0f, 1e-3f));
  EXPECT_THAT(rgba[3], FloatNear(1.0f, 1e-6f));

  // Missing attribute returns empty vector.
  ASSERT_OK_AND_ASSIGN(const std::vector<float> missing,
                       ReadGltfAttributeFloats(model, prim, "TEXCOORD_0", 2));
  EXPECT_THAT(missing, IsEmpty());

  // Non-finite float attribute is rejected.
  const std::vector<float> bad_pos = {
      0.0f, std::numeric_limits<float>::quiet_NaN(), 0.0f};
  ASSERT_OK_AND_ASSIGN(const int pos_bv,
                       AppendGltfBufferView(model, 0, bad_pos.data(),
                                            bad_pos.size() * sizeof(float),
                                            TINYGLTF_TARGET_ARRAY_BUFFER));
  tinygltf::Accessor pos_acc;
  pos_acc.bufferView = pos_bv;
  pos_acc.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  pos_acc.count = 1;
  pos_acc.type = TINYGLTF_TYPE_VEC3;
  model.accessors.push_back(pos_acc);
  prim.attributes["POSITION"] = 1;
  EXPECT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfBufferReaderTest, ReadsTriangleIndicesOfEveryIndexType) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<uint8_t> u8 = {0, 1, 2, 2, 1, 3};
  const std::vector<uint16_t> u16 = {0, 1, 2, 2, 1, 3};
  const std::vector<uint32_t> u32 = {0, 1, 2, 2, 1, 3};
  for (int acc :
       {AddAccessor(model, u8, TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE,
                    TINYGLTF_TYPE_SCALAR, 6),
        AddAccessor(model, u16, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                    TINYGLTF_TYPE_SCALAR, 6),
        AddAccessor(model, u32, TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT,
                    TINYGLTF_TYPE_SCALAR, 6)}) {
    tinygltf::Primitive prim;
    prim.indices = acc;
    ASSERT_OK_AND_ASSIGN(const std::vector<unsigned int> indices,
                         ReadGltfTriangleIndices(model, prim));
    EXPECT_THAT(indices, ElementsAre(0, 1, 2, 2, 1, 3)) << "accessor " << acc;
  }
}

TEST(GltfBufferReaderTest, GeneratesIndicesForNonIndexedPrimitives) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;  // prim.indices == -1
  const std::vector<float> positions = {
      0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f, 0.0f, 1.0f, 0.0f,
      0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f, 2.0f, 0.0f, 0.0f,
  };
  prim.attributes["POSITION"] = AddAccessor(
      model, positions, TINYGLTF_COMPONENT_TYPE_FLOAT, TINYGLTF_TYPE_VEC3, 6);
  ASSERT_OK_AND_ASSIGN(const std::vector<unsigned int> indices,
                       ReadGltfTriangleIndices(model, prim));
  EXPECT_THAT(indices, ElementsAre(0, 1, 2, 3, 4, 5));
}

// Out-of-range indices are rejected by ValidateGltfBuffers() before reading;
// see ValidateGltfBuffersCatchesBadOffsetsAndIndices.
TEST(GltfBufferReaderTest, RejectsNonTriangleIndexCounts) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;

  // Non-multiple-of-3 index count is rejected.
  const std::vector<uint16_t> non_multiple_of_3 = {0, 1};
  prim.indices = AddAccessor(model, non_multiple_of_3,
                             TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                             TINYGLTF_TYPE_SCALAR, 2);
  EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Non-indexed primitive whose vertex count is not a multiple of 3 (4
  // vertices in kQuadPositions) is rejected.
  prim.indices = -1;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Non-indexed primitive whose POSITION accessor doesn't exist.
  const int quad = prim.attributes["POSITION"];
  prim.attributes["POSITION"] = static_cast<int>(model.accessors.size());
  EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Non-indexed primitive whose POSITION claims far more vertices than its
  // data holds; the count must not be trusted for allocating indices.
  prim.attributes["POSITION"] = quad;
  model.accessors[quad].count = 1'000'000'000'000'000;
  EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Neither indices nor POSITION.
  prim.attributes.clear();
  EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfBufferReaderTest, RejectsNonIndexComponentTypes) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<int16_t> signed_indices = {0, 1, 2};
  const std::vector<float> float_indices = {0.0f, 1.0f, 2.0f};
  for (int acc :
       {AddAccessor(model, signed_indices, TINYGLTF_COMPONENT_TYPE_SHORT,
                    TINYGLTF_TYPE_SCALAR, 3),
        AddAccessor(model, float_indices, TINYGLTF_COMPONENT_TYPE_FLOAT,
                    TINYGLTF_TYPE_SCALAR, 3)}) {
    tinygltf::Primitive prim;
    prim.indices = acc;
    EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "accessor " << acc;
  }
}

TEST(GltfBufferReaderTest, ReadsUnsignedShortNormalizedAttributes) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<uint16_t> uvs = {0, 65535, 32768, 16384};
  tinygltf::Primitive prim;
  prim.attributes["TEXCOORD_0"] =
      AddAccessor(model, uvs, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                  TINYGLTF_TYPE_VEC2, 2, /*normalized=*/true);
  ASSERT_OK_AND_ASSIGN(const std::vector<float> values,
                       ReadGltfAttributeFloats(model, prim, "TEXCOORD_0", 2));
  ASSERT_THAT(values, SizeIs(4));
  EXPECT_THAT(values[0], FloatNear(0.0f, 1e-6f));
  EXPECT_THAT(values[1], FloatNear(1.0f, 1e-6f));
  EXPECT_THAT(values[2], FloatNear(32768.0f / 65535.0f, 1e-6f));
  EXPECT_THAT(values[3], FloatNear(16384.0f / 65535.0f, 1e-6f));
}

TEST(GltfBufferReaderTest, RejectsOutOfBoundsAndSparseAccessors) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  // Claims 5 positions but the bufferView only holds 4.
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 5);
  EXPECT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3),
              StatusIs(absl::StatusCode::kInvalidArgument));

  model.accessors[prim.attributes["POSITION"]].count = 4;
  model.accessors[prim.attributes["POSITION"]].sparse.isSparse = true;
  EXPECT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfBufferReaderTest, RejectsNormalizedFloatAttributes) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4, /*normalized=*/true);
  EXPECT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfBufferReaderTest, RejectsUnsupportedTypeOfEmptyAccessor) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<double> positions = {0.0, 0.0, 0.0};
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, positions, TINYGLTF_COMPONENT_TYPE_DOUBLE,
                  TINYGLTF_TYPE_VEC3, /*count=*/0);
  EXPECT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfBufferReaderTest, RejectsOutOfRangeOutputComponents) {
  tinygltf::Model model = MakeQuadModel();
  const tinygltf::Primitive& prim = model.meshes[0].primitives[0];
  ASSERT_THAT(ReadGltfAttributeFloats(model, prim, "POSITION", 3), IsOk());
  for (int out_components : {0, 5}) {
    EXPECT_THAT(
        ReadGltfAttributeFloats(model, prim, "POSITION", out_components),
        StatusIs(absl::StatusCode::kInvalidArgument))
        << out_components << " components";
  }
}

TEST(GltfBufferReaderTest, ReadsIndicesOnlyOfTriangleLists) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 3);
  for (int mode : {kGltfModeNotSet, TINYGLTF_MODE_TRIANGLES}) {
    prim.mode = mode;
    EXPECT_THAT(ReadGltfTriangleIndices(model, prim), IsOk())
        << "mode " << mode;
  }
  for (int mode : {TINYGLTF_MODE_POINTS, TINYGLTF_MODE_LINE,
                   TINYGLTF_MODE_TRIANGLE_STRIP, TINYGLTF_MODE_TRIANGLE_FAN}) {
    prim.mode = mode;
    EXPECT_THAT(ReadGltfTriangleIndices(model, prim),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "mode " << mode;
  }
}

TEST(GltfBufferReaderTest, ReadsSignedNormalizedAttributes) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<int8_t> normals = {127, -127, -128, 0};
  const std::vector<int16_t> uvs = {32767, -32768};
  tinygltf::Primitive prim;
  prim.attributes["NORMAL"] =
      AddAccessor(model, normals, TINYGLTF_COMPONENT_TYPE_BYTE,
                  TINYGLTF_TYPE_VEC4, 1, true);
  prim.attributes["TEXCOORD_0"] = AddAccessor(
      model, uvs, TINYGLTF_COMPONENT_TYPE_SHORT, TINYGLTF_TYPE_VEC2, 1, true);
  ASSERT_OK_AND_ASSIGN(const std::vector<float> n,
                       ReadGltfAttributeFloats(model, prim, "NORMAL", 4));
  // -128 clamps to -1, as glTF specifies.
  EXPECT_THAT(n, ElementsAre(1.0f, -1.0f, -1.0f, 0.0f));
  ASSERT_OK_AND_ASSIGN(const std::vector<float> uv,
                       ReadGltfAttributeFloats(model, prim, "TEXCOORD_0", 2));
  EXPECT_THAT(uv, ElementsAre(1.0f, -1.0f));
}

}  // namespace
}  // namespace intrinsic::geo
