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

#include "intrinsic/geometry/gltf/gltf_validation.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cstddef>
#include <cstdint>
#include <cstring>
#include <limits>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "absl/status/statusor.h"
#include "internal/testing.h"
#include "intrinsic/geometry/gltf/gltf_accessor.h"
#include "intrinsic/geometry/gltf/gltf_buffer_test_util.h"
#include "intrinsic/geometry/gltf/gltf_buffer_writer.h"
#include "intrinsic/util/status/status_macros.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::IsEmpty;

TEST(GltfValidationTest, ValidateGltfBuffersCatchesBadOffsetsAndIndices) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  const std::vector<uint16_t> indices = {0, 1, 2};
  prim.indices =
      AddAccessor(model, indices, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                  TINYGLTF_TYPE_SCALAR, 3);
  tinygltf::Mesh mesh;
  mesh.primitives.push_back(prim);
  model.meshes.push_back(mesh);
  EXPECT_THAT(ValidateGltfBuffers(model), IsOk());

  // An index that points past the last vertex.
  tinygltf::Model bad_index = model;
  const tinygltf::BufferView& view =
      bad_index.bufferViews[bad_index.accessors[prim.indices].bufferView];
  constexpr uint16_t kTooBig = 9;
  std::memcpy(&bad_index.buffers[0].data[view.byteOffset], &kTooBig,
              sizeof(kTooBig));
  EXPECT_THAT(ValidateGltfBuffers(bad_index),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A buffer view that runs past the end of its buffer.
  tinygltf::Model bad_view = model;
  bad_view.bufferViews[0].byteLength = 1 << 20;
  EXPECT_THAT(ValidateGltfBuffers(bad_view),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Indices stored as 12-byte float triples instead of unsigned scalars.
  tinygltf::Model bad_index_type = model;
  const std::vector<float> float_indices = {0.0f, 1.0f, 2.0f};
  bad_index_type.meshes[0].primitives[0].indices =
      AddAccessor(bad_index_type, float_indices, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 1);
  EXPECT_THAT(ValidateGltfBuffers(bad_index_type),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // POSITION refers to an accessor that doesn't exist.
  tinygltf::Model bad_position = model;
  bad_position.meshes[0].primitives[0].attributes["POSITION"] = 99;
  EXPECT_THAT(ValidateGltfBuffers(bad_position),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Another attribute or the indices refer to an accessor that doesn't exist.
  tinygltf::Model bad_normal = model;
  bad_normal.meshes[0].primitives[0].attributes["NORMAL"] = 999;
  EXPECT_THAT(ValidateGltfBuffers(bad_normal),
              StatusIs(absl::StatusCode::kInvalidArgument));
  tinygltf::Model bad_indices = model;
  bad_indices.meshes[0].primitives[0].attributes.erase("POSITION");
  bad_indices.meshes[0].primitives[0].indices = 999;
  EXPECT_THAT(ValidateGltfBuffers(bad_indices),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A NORMAL accessor with fewer vertices than POSITION.
  tinygltf::Model short_normal = model;
  const std::vector<float> normals = {0, 0, 1, 0, 0, 1};
  short_normal.meshes[0].primitives[0].attributes["NORMAL"] =
      AddAccessor(short_normal, normals, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 2);
  EXPECT_THAT(ValidateGltfBuffers(short_normal),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// The output rules ValidateGltfBuffers() checks before serializing.
TEST(GltfValidationTest, ValidateGltfBuffersRejectsInvalidOutput) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 4);
  model.meshes.emplace_back().primitives.push_back(prim);
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  // An empty buffer (GLB readers reject an empty BIN chunk).
  tinygltf::Model empty_buffer = model;
  empty_buffer.buffers.emplace_back();
  EXPECT_THAT(ValidateGltfBuffers(empty_buffer),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // An accessor with count 0, e.g. after every triangle was removed.
  tinygltf::Model empty_accessor = model;
  empty_accessor.accessors[0].count = 0;
  EXPECT_THAT(ValidateGltfBuffers(empty_accessor),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Non-finite POSITION bounds, e.g. from NaN positions.
  tinygltf::Model nan_bounds = model;
  nan_bounds.accessors[0].minValues = {0, 0, 0};
  nan_bounds.accessors[0].maxValues = {
      1, std::numeric_limits<double>::infinity(), 0};
  EXPECT_THAT(ValidateGltfBuffers(nan_bounds),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A skinned node whose mesh was removed.
  tinygltf::Model skin_without_mesh = model;
  skin_without_mesh.skins.emplace_back();
  skin_without_mesh.nodes.emplace_back().skin = 0;
  EXPECT_THAT(ValidateGltfBuffers(skin_without_mesh),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A MAT2 with byte components, whose padded columns aren't supported.
  tinygltf::Model padded_matrix = model;
  const std::vector<uint8_t> mat2 = {1, 2, 0, 0, 3, 4, 0, 0};
  AddAccessor(padded_matrix, mat2, TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE,
              TINYGLTF_TYPE_MAT2, 1);
  EXPECT_THAT(ValidateGltfBuffers(padded_matrix),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Removing every triangle is rejected by WriteRemappedPrimitiveBuffers(), and
// the empty accessors it would give are rejected by ValidateGltfBuffers().
TEST(GltfValidationTest, ValidationRejectsRemapToNothing) {
  tinygltf::Model model = MakeQuadModel();
  EXPECT_THAT(
      WriteRemappedPrimitiveBuffers(model, model.meshes[0].primitives[0],
                                    /*buffer_index=*/0,
                                    /*source_vertices=*/{},
                                    /*indices=*/{}),
      StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model empty = MakeQuadModel();
  empty.accessors[empty.meshes[0].primitives[0].attributes.at("POSITION")]
      .count = 0;
  EXPECT_THAT(ValidateGltfBuffers(empty),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// NaN positions are rejected by ValidateGltfBuffers() (and by
// WriteRemappedPrimitiveBuffers(), since they would make the bounds wrong).
TEST(GltfValidationTest, ValidationRejectsNanPositionBounds) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  const float nan = std::numeric_limits<float>::quiet_NaN();
  const std::vector<float> positions = {nan, 0, 0, 1, 0, 0, 0, 1, 0};
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] = AddAccessor(
      model, positions, TINYGLTF_COMPONENT_TYPE_FLOAT, TINYGLTF_TYPE_VEC3, 3);
  model.meshes.emplace_back().primitives.push_back(prim);
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Compacting a model with no data leaves an empty buffer
// (which GLB readers reject); ValidateGltfBuffers() catches it.
TEST(GltfValidationTest, ValidationRejectsCompactedModelWithoutData) {
  tinygltf::Model model;
  ASSERT_THAT(CompactGlbBuffers(model), IsOk());
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Draco-compressed primitives have accessors without a
// bufferView, which ValidateGltfBuffers() rejects before CompactGlbBuffers()
// could leave the extension pointing at the wrong bufferView.
TEST(GltfValidationTest, ValidationRejectsAccessorsWithoutBufferView) {
  tinygltf::Model model = MakeQuadModel();
  model.accessors[0].bufferView = -1;
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// glTF requires POSITION to be a VEC3; anything else would be read past its
// elements.
TEST(GltfValidationTest, ValidationRejectsNonVec3Position) {
  tinygltf::Model model = MakeQuadModel();
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());
  model.accessors[0].type = TINYGLTF_TYPE_SCALAR;
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// POSITION bounds may be missing, but if present they need 3 values.
TEST(GltfValidationTest, ValidationChecksPositionBoundsSize) {
  tinygltf::Model model = MakeQuadModel();
  ASSERT_THAT(model.accessors[0].minValues, IsEmpty());
  EXPECT_THAT(ValidateGltfBuffers(model), IsOk());
  model.accessors[0].minValues = {0.0, 0.0, 0.0};
  model.accessors[0].maxValues = {1.0, 1.0, 0.0};
  EXPECT_THAT(ValidateGltfBuffers(model), IsOk());
  model.accessors[0].maxValues = {1.0, 1.0};
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Returns a model with one triangle whose POSITION has `component_type`
// components of type `T`, all 0.
template <typename T>
tinygltf::Model MakeTriangleModel(int component_type) {
  constexpr size_t kVertices = 3;
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<T> positions(kVertices * kGltfPositionComponents, T{0});
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] = AddAccessor(model, positions, component_type,
                                            TINYGLTF_TYPE_VEC3, kVertices);
  model.meshes.emplace_back().primitives.push_back(prim);
  return model;
}

TEST(GltfValidationTest, ValidationChecksPositionComponentType) {
  // Floats, and 8/16-bit integers (KHR_mesh_quantization) are supported.
  EXPECT_THAT(ValidateGltfBuffers(
                  MakeTriangleModel<float>(TINYGLTF_COMPONENT_TYPE_FLOAT)),
              IsOk());
  EXPECT_THAT(ValidateGltfBuffers(MakeTriangleModel<uint16_t>(
                  TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT)),
              IsOk());
  EXPECT_THAT(ValidateGltfBuffers(
                  MakeTriangleModel<int8_t>(TINYGLTF_COMPONENT_TYPE_BYTE)),
              IsOk());
  // Doubles and 32-bit integers aren't.
  EXPECT_THAT(ValidateGltfBuffers(
                  MakeTriangleModel<double>(TINYGLTF_COMPONENT_TYPE_DOUBLE)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(ValidateGltfBuffers(
                  MakeTriangleModel<int32_t>(TINYGLTF_COMPONENT_TYPE_INT)),
              StatusIs(absl::StatusCode::kInvalidArgument));
  EXPECT_THAT(ValidateGltfBuffers(MakeTriangleModel<uint32_t>(
                  TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT)),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// glTF forbids the largest value of the index type (e.g. 65535 for 16-bit
// indices), even when there are enough vertices for it.
TEST(GltfValidationTest, ValidationRejectsLargestIndexValue) {
  constexpr size_t kVertices = 65536;
  tinygltf::Model model;
  model.buffers.emplace_back();
  const std::vector<float> positions(kVertices * kGltfPositionComponents, 0.0f);
  tinygltf::Primitive prim;
  prim.attributes["POSITION"] =
      AddAccessor(model, positions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, kVertices);
  const std::vector<uint16_t> valid = {0, 1, 65534};
  prim.indices =
      AddAccessor(model, valid, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                  TINYGLTF_TYPE_SCALAR, valid.size());
  model.meshes.emplace_back().primitives.push_back(prim);
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  const std::vector<uint16_t> restart = {0, 1, 65535};
  model.meshes[0].primitives[0].indices =
      AddAccessor(model, restart, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                  TINYGLTF_TYPE_SCALAR, restart.size());
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfValidationTest, ValidationRequiresPosition) {
  tinygltf::Model model = MakeQuadModel();
  tinygltf::Primitive& prim = model.meshes[0].primitives[0];
  // Reuses the POSITION data, which has the same vertex count.
  prim.attributes["NORMAL"] = prim.attributes.at("POSITION");
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());
  prim.attributes.erase("POSITION");
  EXPECT_THAT(ValidateGltfBuffers(model),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfValidationTest, ValidationChecksMorphTargets) {
  tinygltf::Model model = MakeQuadModel();
  const int position = model.meshes[0].primitives[0].attributes.at("POSITION");
  model.meshes[0].primitives[0].targets.push_back({{"POSITION", position}});
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  // A morph target accessor that doesn't exist.
  tinygltf::Model missing = model;
  missing.meshes[0].primitives[0].targets[0]["POSITION"] =
      static_cast<int>(missing.accessors.size());
  EXPECT_THAT(ValidateGltfBuffers(missing),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A morph target with 3 vertices for a 4-vertex primitive.
  tinygltf::Model fewer_vertices = model;
  fewer_vertices.meshes[0].primitives[0].targets[0]["POSITION"] =
      AddAccessor(fewer_vertices, kQuadPositions, TINYGLTF_COMPONENT_TYPE_FLOAT,
                  TINYGLTF_TYPE_VEC3, 3);
  EXPECT_THAT(ValidateGltfBuffers(fewer_vertices),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfValidationTest, ValidationChecksNodeMeshAndSkin) {
  tinygltf::Model model = MakeQuadModel();
  model.skins.emplace_back();
  tinygltf::Node& node = model.nodes.emplace_back();
  node.mesh = 0;
  node.skin = 0;
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  for (int mesh : {1, -2}) {
    tinygltf::Model bad_mesh = model;
    bad_mesh.nodes[0].mesh = mesh;
    EXPECT_THAT(ValidateGltfBuffers(bad_mesh),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "mesh " << mesh;
  }
  for (int skin : {1, -2}) {
    tinygltf::Model bad_skin = model;
    bad_skin.nodes[0].skin = skin;
    EXPECT_THAT(ValidateGltfBuffers(bad_skin),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "skin " << skin;
  }
}

TEST(GltfValidationTest, ValidationChecksSkinAnimationAndInstancingAccessors) {
  tinygltf::Model model = MakeQuadModel();
  const int missing_accessor = static_cast<int>(model.accessors.size());
  model.skins.emplace_back().inverseBindMatrices = 0;
  tinygltf::AnimationSampler& sampler =
      model.animations.emplace_back().samplers.emplace_back();
  sampler.input = 0;
  sampler.output = 0;
  model.nodes.emplace_back().extensions["EXT_mesh_gpu_instancing"] =
      tinygltf::Value(tinygltf::Value::Object{
          {"attributes", tinygltf::Value(tinygltf::Value::Object{
                             {"TRANSLATION", tinygltf::Value(0)}})}});
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  tinygltf::Model bad_skin = model;
  bad_skin.skins[0].inverseBindMatrices = missing_accessor;
  EXPECT_THAT(ValidateGltfBuffers(bad_skin),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model bad_sampler = model;
  bad_sampler.animations[0].samplers[0].output = missing_accessor;
  EXPECT_THAT(ValidateGltfBuffers(bad_sampler),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model bad_instancing = model;
  bad_instancing.nodes[0].extensions["EXT_mesh_gpu_instancing"] =
      tinygltf::Value(tinygltf::Value::Object{
          {"attributes",
           tinygltf::Value(tinygltf::Value::Object{
               {"TRANSLATION", tinygltf::Value(missing_accessor)}})}});
  EXPECT_THAT(ValidateGltfBuffers(bad_instancing),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model no_instancing_attributes = model;
  no_instancing_attributes.nodes[0].extensions["EXT_mesh_gpu_instancing"] =
      tinygltf::Value(tinygltf::Value::Object{});
  EXPECT_THAT(ValidateGltfBuffers(no_instancing_attributes),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

// Adds a sparse float VEC3 accessor of 4 elements without base data (as morph
// targets often are) that replaces element 3 with `value`, and returns its
// index.
absl::StatusOr<int> AddSparseAccessor(tinygltf::Model& model,
                                      const std::array<float, 3>& value) {
  constexpr std::array<uint16_t, 1> kReplaced = {3};
  tinygltf::Accessor acc;
  acc.componentType = TINYGLTF_COMPONENT_TYPE_FLOAT;
  acc.type = TINYGLTF_TYPE_VEC3;
  acc.count = 4;
  acc.sparse.isSparse = true;
  acc.sparse.count = static_cast<int>(kReplaced.size());
  acc.sparse.indices.componentType = TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT;
  INTR_ASSIGN_OR_RETURN(
      acc.sparse.indices.bufferView,
      AppendGltfBufferView(model, /*buffer_index=*/0, kReplaced.data(),
                           sizeof(kReplaced), /*target=*/0));
  INTR_ASSIGN_OR_RETURN(acc.sparse.values.bufferView,
                        AppendGltfBufferView(model, /*buffer_index=*/0,
                                             value.data(), sizeof(value),
                                             /*target=*/0));
  model.accessors.push_back(std::move(acc));
  return static_cast<int>(model.accessors.size() - 1);
}

TEST(GltfValidationTest, ValidationChecksSparseAccessors) {
  tinygltf::Model model = MakeQuadModel();
  ASSERT_OK_AND_ASSIGN(const int sparse,
                       AddSparseAccessor(model, {0.0f, 0.0f, 5.0f}));
  ASSERT_THAT(ValidateGltfBuffers(model), IsOk());

  // Base data from the quad's 4 positions is also allowed.
  tinygltf::Model with_base = model;
  with_base.accessors[sparse].bufferView = 0;
  EXPECT_THAT(ValidateGltfBuffers(with_base), IsOk());

  // Base data past the end of its bufferView.
  tinygltf::Model base_past_end = with_base;
  base_past_end.accessors[sparse].byteOffset = sizeof(float);
  EXPECT_THAT(ValidateGltfBuffers(base_past_end),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // No replaced elements, or more than the accessor has.
  for (int count : {0, 5}) {
    tinygltf::Model bad_count = model;
    bad_count.accessors[sparse].sparse.count = count;
    EXPECT_THAT(ValidateGltfBuffers(bad_count),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "count " << count;
  }

  // Signed sparse indices.
  tinygltf::Model signed_indices = model;
  signed_indices.accessors[sparse].sparse.indices.componentType =
      TINYGLTF_COMPONENT_TYPE_SHORT;
  EXPECT_THAT(ValidateGltfBuffers(signed_indices),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // A sparse index (3) past the accessor's count.
  tinygltf::Model index_out_of_range = model;
  index_out_of_range.accessors[sparse].count = 3;
  EXPECT_THAT(ValidateGltfBuffers(index_out_of_range),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Sparse values past the end of their bufferView.
  tinygltf::Model values_past_end = model;
  values_past_end.accessors[sparse].sparse.values.byteOffset = sizeof(float);
  EXPECT_THAT(ValidateGltfBuffers(values_past_end),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfValidationTest, ValidationRejectsNanSparsePositionValues) {
  tinygltf::Model finite = MakeQuadModel();
  ASSERT_OK_AND_ASSIGN(const int finite_sparse,
                       AddSparseAccessor(finite, {0.0f, 0.0f, 5.0f}));
  finite.meshes[0].primitives[0].attributes["POSITION"] = finite_sparse;
  EXPECT_THAT(ValidateGltfBuffers(finite), IsOk());

  tinygltf::Model nan = MakeQuadModel();
  ASSERT_OK_AND_ASSIGN(
      const int nan_sparse,
      AddSparseAccessor(nan,
                        {0.0f, std::numeric_limits<float>::quiet_NaN(), 5.0f}));
  nan.meshes[0].primitives[0].attributes["POSITION"] = nan_sparse;
  EXPECT_THAT(ValidateGltfBuffers(nan),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

}  // namespace
}  // namespace intrinsic::geo
