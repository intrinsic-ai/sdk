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

#include "intrinsic/geometry/gltf/gltf_buffer_writer.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cstdint>
#include <cstring>
#include <limits>
#include <numeric>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "internal/testing.h"
#include "intrinsic/geometry/gltf/gltf_accessor.h"
#include "intrinsic/geometry/gltf/gltf_buffer_reader.h"
#include "intrinsic/geometry/gltf/gltf_buffer_test_util.h"
#include "intrinsic/geometry/gltf/gltf_io.h"
#include "intrinsic/geometry/gltf/gltf_validation.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::ElementsAre;
using ::testing::Eq;
using ::testing::FloatNear;
using ::testing::IsEmpty;
using ::testing::SizeIs;

// Returns 0, 1, ..., n - 1.
std::vector<unsigned int> Iota(unsigned int n) {
  std::vector<unsigned int> values(n);
  std::iota(values.begin(), values.end(), 0u);
  return values;
}

// Returns `n` distinct float VEC3 positions, (i, 0, 0) for vertex i.
std::vector<float> PositionsAlongX(unsigned int n) {
  std::vector<float> positions(n * kGltfPositionComponents, 0.0f);
  for (unsigned int v = 0; v < n; ++v) {
    positions[v * kGltfPositionComponents] = static_cast<float>(v);
  }
  return positions;
}

TEST(GltfBufferWriterTest,
     WritesFourByteAlignedAttributesAndCompactsMeshesAndExtensions) {
  tinygltf::Model model;
  model.buffers.resize(2);

  const std::vector<float> pos = {
      0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 0.0f, 0.0f, 1.0f, 0.0f, 9.0f, 9.0f, 9.0f,
  };
  const std::vector<uint8_t> rgb = {10, 20, 30, 0, 40, 50, 60, 0,
                                    70, 80, 90, 0, 99, 99, 99, 0};
  ASSERT_OK_AND_ASSIGN(
      const int pos_bv,
      AppendGltfBufferView(model, 0, pos.data(), pos.size() * sizeof(float),
                           TINYGLTF_TARGET_ARRAY_BUFFER));
  ASSERT_OK_AND_ASSIGN(const int col_bv,
                       AppendGltfBufferView(model, 0, rgb.data(), rgb.size(),
                                            TINYGLTF_TARGET_ARRAY_BUFFER, 4));

  tinygltf::Accessor pos_acc;
  pos_acc.bufferView = pos_bv;
  pos_acc.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  pos_acc.count = 4;
  pos_acc.type = TINYGLTF_TYPE_VEC3;
  model.accessors.push_back(pos_acc);

  tinygltf::Accessor col_acc;
  col_acc.bufferView = col_bv;
  col_acc.componentType = TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE;
  col_acc.normalized = true;
  col_acc.count = 4;
  col_acc.type = TINYGLTF_TYPE_VEC3;
  model.accessors.push_back(col_acc);

  // Dead accessor + live GPU-instancing accessor.
  model.accessors.push_back(pos_acc);  // index 2 (dead)
  const std::vector<float> inst = {1.0f, 2.0f, 3.0f};
  ASSERT_OK_AND_ASSIGN(
      const int inst_bv,
      AppendGltfBufferView(model, 0, inst.data(), sizeof(float) * 3, 0));
  tinygltf::Accessor inst_acc;
  inst_acc.bufferView = inst_bv;
  inst_acc.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  inst_acc.count = 1;
  inst_acc.type = TINYGLTF_TYPE_VEC3;
  model.accessors.push_back(inst_acc);  // index 3 (referenced by extension)

  tinygltf::Primitive prim;
  prim.attributes["POSITION"] = 0;
  prim.attributes["COLOR_0"] = 1;
  // Keep vertices 0..2 and drop vertex 3.
  const std::vector<unsigned int> sources = {0, 1, 2};
  const std::vector<unsigned int> simp_idx = {0, 1, 2};
  ASSERT_THAT(
      WriteRemappedPrimitiveBuffers(model, prim,
                                    /*buffer_index=*/1, sources, simp_idx),
      IsOk());

  tinygltf::Mesh empty_mesh;
  tinygltf::Mesh valid_mesh;
  valid_mesh.primitives.push_back(prim);
  model.meshes = {empty_mesh, valid_mesh};

  tinygltf::Node n0;
  n0.mesh = 0;  // points to empty_mesh -> should become -1
  tinygltf::Node n1;
  n1.mesh = 1;  // points to valid_mesh -> should shift to 0
  n1.extensions["EXT_mesh_gpu_instancing"] =
      tinygltf::Value(tinygltf::Value::Object{
          {"attributes", tinygltf::Value(tinygltf::Value::Object{
                             {"TRANSLATION", tinygltf::Value(3)},
                         })}});
  model.nodes = {n0, n1};

  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  EXPECT_THAT(model.buffers, SizeIs(1));
  ASSERT_THAT(model.meshes, SizeIs(1));
  EXPECT_THAT(model.nodes[0].mesh, Eq(-1));
  EXPECT_THAT(model.nodes[1].mesh, Eq(0));

  // Sub-4-byte VEC3 UNSIGNED_BYTE COLOR_0 has byteStride == 4.
  const int out_col_acc =
      model.meshes[0].primitives[0].attributes.at("COLOR_0");
  EXPECT_THAT(
      model.bufferViews[model.accessors[out_col_acc].bufferView].byteStride,
      Eq(4));

  // POSITION min/max excludes dropped vertex 3 (9, 9, 9).
  const int out_pos_acc =
      model.meshes[0].primitives[0].attributes.at("POSITION");
  EXPECT_THAT(model.accessors[out_pos_acc].maxValues,
              ElementsAre(1.0, 1.0, 0.0));

  // Only 4 accessors remain (POSITION, COLOR_0, indices, instancing
  // TRANSLATION).
  EXPECT_THAT(model.accessors, SizeIs(4));
}

TEST(GltfBufferWriterTest, AppendRejectsInvalidBufferIndex) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  constexpr std::array<uint8_t, 4> kData = {1, 2, 3, 4};
  for (int buffer_index : {-1, 1}) {
    EXPECT_THAT(AppendGltfBufferView(model, buffer_index, kData.data(),
                                     kData.size(), /*target=*/0),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "buffer " << buffer_index;
  }
  EXPECT_THAT(model.bufferViews, IsEmpty());
}

TEST(GltfBufferWriterTest, AppendRejectsDataFromTheSameBuffer) {
  tinygltf::Model model = MakeQuadModel();
  const std::vector<unsigned char>& bin = model.buffers[0].data;
  EXPECT_THAT(AppendGltfBufferView(model, /*buffer_index=*/0, bin.data(),
                                   bin.size(), /*target=*/0),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // A copy is fine.
  const std::vector<unsigned char> copy = bin;
  EXPECT_THAT(AppendGltfBufferView(model, /*buffer_index=*/0, copy.data(),
                                   copy.size(), /*target=*/0),
              IsOk());
}

// CompactGlbBuffers keeps data that is only referenced from sparse accessors
// (common for morph targets) and embedded images, with valid offsets.
TEST(GltfBufferWriterTest, CompactKeepsSparseAccessorDataAndImages) {
  tinygltf::Model model;
  model.buffers.resize(2);
  // Garbage in buffer 0 that compaction should drop.
  const std::vector<float> garbage(64, 7.0f);
  ASSERT_THAT(AppendGltfBufferView(model, 0, garbage.data(),
                                   garbage.size() * sizeof(float), 0),
              IsOk());

  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);

  // Morph target: a sparse accessor moving vertex 3 by (0, 0, 5).
  const std::vector<uint16_t> sparse_indices = {3};
  const std::vector<float> sparse_values = {0.0f, 0.0f, 5.0f};
  tinygltf::Accessor target;
  target.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  target.type = TINYGLTF_TYPE_VEC3;
  target.count = 4;
  target.sparse.isSparse = true;
  target.sparse.count = 1;
  target.sparse.indices.componentType = TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT;
  ASSERT_OK_AND_ASSIGN(target.sparse.indices.bufferView,
                       AppendGltfBufferView(model, 1, sparse_indices.data(),
                                            sparse_indices.size() * 2, 0));
  ASSERT_OK_AND_ASSIGN(
      target.sparse.values.bufferView,
      AppendGltfBufferView(model, 1, sparse_values.data(),
                           sparse_values.size() * sizeof(float), 0));
  model.accessors.push_back(target);
  prim.targets.push_back({{"POSITION", 1}});

  const std::vector<unsigned char> png = {0x89, 'P', 'N', 'G', 1, 2, 3};
  tinygltf::Image image;
  ASSERT_OK_AND_ASSIGN(
      image.bufferView,
      AppendGltfBufferView(model, 1, png.data(), png.size(), 0));
  image.mimeType = "image/png";
  model.images.push_back(image);

  tinygltf::Mesh mesh;
  mesh.primitives.push_back(prim);
  model.meshes.push_back(mesh);

  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  ASSERT_THAT(model.buffers, SizeIs(1));
  // POSITION, sparse indices, sparse values, image; the garbage is gone.
  EXPECT_THAT(model.bufferViews, SizeIs(4));

  const auto bytes_of = [&model](int bv) {
    const tinygltf::BufferView& view = model.bufferViews[bv];
    return std::vector<unsigned char>(
        model.buffers[0].data.begin() + view.byteOffset,
        model.buffers[0].data.begin() + view.byteOffset + view.byteLength);
  };
  const tinygltf::Accessor& out_target =
      model.accessors[model.meshes[0].primitives[0].targets[0].at("POSITION")];
  ASSERT_TRUE(out_target.sparse.isSparse);
  const std::vector<unsigned char> out_indices =
      bytes_of(out_target.sparse.indices.bufferView);
  const std::vector<unsigned char> out_values =
      bytes_of(out_target.sparse.values.bufferView);
  uint16_t index = 0;
  ASSERT_THAT(out_indices, SizeIs(sizeof(index)));
  std::memcpy(&index, out_indices.data(), sizeof(index));
  EXPECT_THAT(index, Eq(3));
  float value[3] = {};
  ASSERT_THAT(out_values, SizeIs(sizeof(value)));
  std::memcpy(value, out_values.data(), sizeof(value));
  EXPECT_THAT(value[2], FloatNear(5.0f, 1e-6f));
  EXPECT_THAT(bytes_of(model.images[0].bufferView), Eq(png));
}

TEST(GltfBufferWriterTest, WritesSmallIndicesAs16Bit) {
  tinygltf::Model model;
  model.buffers.resize(2);
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  const std::vector<unsigned int> sources = {0, 1, 2, 3};
  ASSERT_THAT(WriteRemappedPrimitiveBuffers(
                  model, prim,
                  /*buffer_index=*/1, sources,
                  std::vector<unsigned int>{0, 1, 2, 2, 1, 3}),
              IsOk());
  const tinygltf::Accessor& indices = model.accessors[prim.indices];
  EXPECT_THAT(indices.componentType,
              Eq(TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT));
  ASSERT_OK_AND_ASSIGN(const std::vector<unsigned int> read,
                       ReadGltfTriangleIndices(model, prim));
  EXPECT_THAT(read, ElementsAre(0, 1, 2, 2, 1, 3));
}

// POSITION bounds hold the raw stored values, even for normalized integers.
TEST(GltfBufferWriterTest, ComputesRawBoundsForIntegerPositions) {
  tinygltf::Model model;
  model.buffers.resize(2);
  const std::vector<uint16_t> positions = {
      10, 20, 30, 40, 50, 60, 0, 0, 0, 0};  // last value pads to 4 bytes
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, positions, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                  TINYGLTF_TYPE_VEC3, 3, /*normalized=*/true);
  const std::vector<unsigned int> sources = {0, 1};
  ASSERT_THAT(WriteRemappedPrimitiveBuffers(model, prim,
                                            /*buffer_index=*/1, sources,
                                            std::vector<unsigned int>{0, 1, 0}),
              IsOk());
  const tinygltf::Accessor& out = model.accessors[prim.attributes["POSITION"]];
  EXPECT_THAT(out.minValues, ElementsAre(10.0, 20.0, 30.0));
  EXPECT_THAT(out.maxValues, ElementsAre(40.0, 50.0, 60.0));
}

TEST(GltfBufferWriterTest, WriteRemappedRejectsBadInputs) {
  tinygltf::Model model;
  model.buffers.resize(2);
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  const std::vector<unsigned int> indices = {0, 1, 2};
  const auto write = [&model](tinygltf::Primitive p, int scratch,
                              std::vector<unsigned int> sources,
                              std::vector<unsigned int> idx) {
    tinygltf::Model m = model;
    return WriteRemappedPrimitiveBuffers(m, p, scratch, sources, idx);
  };
  EXPECT_THAT(write(prim, 1, {0, 1, 2}, indices), IsOk());
  // Source vertex past the end of the 4-vertex input.
  EXPECT_THAT(write(prim, 1, {0, 1, 4}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // Index past the 3 new vertices.
  EXPECT_THAT(write(prim, 1, {0, 1, 2}, {0, 1, 3}),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // No indices, or an index count that isn't a multiple of 3.
  EXPECT_THAT(write(prim, 1, {0, 1, 2}, {}),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(write(prim, 1, {0, 1, 2}, {0, 1}),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // Scratch buffer that doesn't exist.
  EXPECT_THAT(write(prim, 5, {0, 1, 2}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // POSITION that isn't VEC3.
  tinygltf::Primitive vec2_prim;
  vec2_prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC2, 4);
  EXPECT_THAT(write(vec2_prim, 1, {0, 1, 2}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // Non-finite POSITION value.
  const std::vector<float> nan_positions = {
      std::numeric_limits<float>::quiet_NaN(),
      0.0f,
      0.0f,
      1.0f,
      0.0f,
      0.0f,
      0.0f,
      1.0f,
      0.0f};
  tinygltf::Primitive nan_prim;
  nan_prim.attributes["POSITION"] =
      AddAccessor(model, nan_positions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 3);
  EXPECT_THAT(write(nan_prim, 1, {0, 1, 2}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // Primitive that isn't a triangle list.
  tinygltf::Primitive lines_prim = prim;
  lines_prim.mode = TINYGLTF_MODE_LINE;
  EXPECT_THAT(write(lines_prim, 1, {0, 1, 2}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // Primitive with a morph target.
  tinygltf::Primitive morph_prim = prim;
  morph_prim.targets.push_back({{"POSITION", prim.attributes["POSITION"]}});
  EXPECT_THAT(write(morph_prim, 1, {0, 1, 2}, indices),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// A source vertex can be used more than once, e.g. to split a vertex.
TEST(GltfBufferWriterTest, WriteRemappedCanDuplicateVertices) {
  tinygltf::Model model;
  model.buffers.resize(2);
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  ASSERT_THAT(WriteRemappedPrimitiveBuffers(model, prim, /*buffer_index=*/1,
                                            std::vector<unsigned int>{3, 1, 3},
                                            std::vector<unsigned int>{0, 1, 2}),
              IsOk());
  ASSERT_OK_AND_ASSIGN(const std::vector<float> positions,
                       ReadGltfAttributeFloats(model, prim, "POSITION", 3));
  EXPECT_THAT(positions, ElementsAre(1, 1, 0, 1, 0, 0, 1, 1, 0));
}

// CompactGlbBuffers() keeps bufferView names.
TEST(GltfBufferWriterTest, CompactKeepsBufferViewNames) {
  tinygltf::Model model = MakeQuadModel();
  model.bufferViews[0].name = "positions";
  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  ASSERT_THAT(model.bufferViews, SizeIs(1));
  EXPECT_THAT(model.bufferViews[0].name, Eq("positions"));
}

TEST(GltfBufferWriterTest, CompactRejectsOutOfBoundsBufferViews) {
  // A view past the end of its buffer.
  tinygltf::Model past_end = MakeQuadModel();
  past_end.bufferViews[0].byteLength += 1;
  EXPECT_THAT(CompactGlbBuffers(past_end),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // A view of a buffer that doesn't exist.
  tinygltf::Model missing_buffer = MakeQuadModel();
  missing_buffer.bufferViews[0].buffer = 1;
  EXPECT_THAT(CompactGlbBuffers(missing_buffer),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Remap, compact, validate, write and read back, as the simplifier does.
TEST(GltfBufferWriterTest, RemappedModelRoundTripsThroughGlb) {
  tinygltf::Model model = MakeQuadModel();
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());
  ASSERT_THAT(WriteRemappedPrimitiveBuffers(
                  model, model.meshes[0].primitives[0], /*buffer_index=*/0,
                  std::vector<unsigned int>{0, 1, 2},
                  std::vector<unsigned int>{0, 1, 2}),
              IsOk());
  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());
  ASSERT_OK_AND_ASSIGN(const std::string glb, SerializeGltfToGlb(model));
  ASSERT_OK_AND_ASSIGN(const tinygltf::Model reloaded, ParseGlb(glb));
  ASSERT_THAT(ValidateGltfBuffers(reloaded), IsOk());
  const tinygltf::Primitive& prim = reloaded.meshes[0].primitives[0];
  EXPECT_THAT(reloaded.accessors[prim.attributes.at("POSITION")].count, Eq(3));
  EXPECT_THAT(reloaded.accessors[prim.indices].count, Eq(3));
}

// Indices that don't fit below 65535 are written as 32-bit.
TEST(GltfBufferWriterTest, WritesLargeIndicesAs32Bit) {
  constexpr unsigned int kVertices = 65536;
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, PositionsAlongX(kVertices),
                  TINYGLTF_COMPONENT_TYPE_FLOAT, TINYGLTF_TYPE_VEC3, kVertices);
  ASSERT_THAT(WriteRemappedPrimitiveBuffers(
                  model, prim, /*buffer_index=*/0, Iota(kVertices),
                  std::vector<unsigned int>{0, 1, kVertices - 1}),
              IsOk());
  EXPECT_THAT(model.accessors[prim.indices].componentType,
              Eq(TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT));
  ASSERT_OK_AND_ASSIGN(const std::vector<unsigned int> read,
                       ReadGltfTriangleIndices(model, prim));
  EXPECT_THAT(read, ElementsAre(0, 1, kVertices - 1));
}

// A node whose mesh is removed also loses its skin and morph weights, which
// glTF only allows with a mesh.
TEST(GltfBufferWriterTest, CompactClearsSkinOfNodesThatLoseTheirMesh) {
  tinygltf::Model model = MakeQuadModel();
  model.meshes[0].primitives.clear();
  tinygltf::Node& node = model.nodes.emplace_back();
  node.mesh = 0;
  node.skin = 0;
  node.weights = {0.5};
  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  EXPECT_THAT(model.nodes[0].mesh, Eq(-1));
  EXPECT_THAT(model.nodes[0].skin, Eq(-1));
  EXPECT_THAT(model.nodes[0].weights, IsEmpty());
}

// Accessors used outside mesh primitives (skins, animations, GPU instancing)
// survive compaction and are renumbered.
TEST(GltfBufferWriterTest, CompactKeepsSkinAnimationAndInstancingAccessors) {
  tinygltf::Model model = MakeQuadModel();
  const std::vector<float> unused(4, 1.0f);
  AddAccessor(model, unused, TINYGLTF_COMPONENT_TYPE_FLOAT,
              TINYGLTF_TYPE_SCALAR, 4);  // Dropped by compaction.
  const std::vector<float> matrices(16, 2.0f);
  const std::vector<float> times = {0.0f};
  const std::vector<float> translations = {3.0f, 0.0f, 0.0f};
  tinygltf::Skin& skin = model.skins.emplace_back();
  skin.inverseBindMatrices = AddAccessor(
      model, matrices, TINYGLTF_COMPONENT_TYPE_FLOAT, TINYGLTF_TYPE_MAT4, 1);
  tinygltf::AnimationSampler& sampler =
      model.animations.emplace_back().samplers.emplace_back();
  sampler.input = AddAccessor(model, times, TINYGLTF_COMPONENT_TYPE_FLOAT,
                              TINYGLTF_TYPE_SCALAR, 1);
  sampler.output =
      AddAccessor(model, translations, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 1);
  const int instancing =
      AddAccessor(model, translations, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 1);
  tinygltf::Value::Object attributes;
  attributes["TRANSLATION"] = tinygltf::Value(instancing);
  tinygltf::Value::Object ext;
  ext["attributes"] = tinygltf::Value(attributes);
  model.nodes.emplace_back().extensions["EXT_mesh_gpu_instancing"] =
      tinygltf::Value(ext);

  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  // POSITION, inverse bind matrices, sampler input and output, instancing.
  ASSERT_THAT(model.accessors, SizeIs(5));
  EXPECT_THAT(model.accessors[model.skins[0].inverseBindMatrices].type,
              Eq(TINYGLTF_TYPE_MAT4));
  EXPECT_THAT(model.accessors[model.animations[0].samplers[0].input].type,
              Eq(TINYGLTF_TYPE_SCALAR));
  EXPECT_THAT(model.accessors[model.animations[0].samplers[0].output].type,
              Eq(TINYGLTF_TYPE_VEC3));
  const int new_instancing = model.nodes[0]
                                 .extensions.at("EXT_mesh_gpu_instancing")
                                 .Get("attributes")
                                 .Get("TRANSLATION")
                                 .GetNumberAsInt();
  EXPECT_THAT(new_instancing, Eq(4));
  EXPECT_THAT(ValidateGltfBuffers(model), IsOk());
}

}  // namespace
}  // namespace intrinsic::geo
