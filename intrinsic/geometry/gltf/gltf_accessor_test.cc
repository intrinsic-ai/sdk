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

#include "intrinsic/geometry/gltf/gltf_accessor.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <limits>

#include "absl/status/status.h"
#include "absl/status/status_matchers.h"
#include "internal/testing.h"
#include "intrinsic/geometry/gltf/gltf_buffer_test_util.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::absl_testing::IsOk;
using ::absl_testing::StatusIs;
using ::testing::DoubleEq;
using ::testing::Eq;
using ::testing::IsFalse;
using ::testing::IsTrue;

TEST(GltfAccessorTest, GetsSliceOfStridedAccessor) {
  tinygltf::Model model = MakeQuadModel();
  constexpr size_t kStride = 16;
  model.bufferViews[0].byteStride = kStride;
  model.accessors[0].count = 3;  // The 48 bytes hold 3 elements at stride 16.
  ASSERT_OK_AND_ASSIGN(const GltfAccessorSlice s,
                       GetGltfAccessorSlice(model, 0));
  EXPECT_THAT(s.data, Eq(model.buffers[0].data.data()));
  EXPECT_THAT(s.stride, Eq(kStride));
  EXPECT_THAT(s.elem_size, Eq(3 * sizeof(float)));
  EXPECT_THAT(s.count, Eq(3));
}

TEST(GltfAccessorTest, RejectsInvalidAccessors) {
  const tinygltf::Model model = MakeQuadModel();
  ASSERT_THAT(GetGltfAccessorSlice(model, 0), IsOk());
  for (int index : {-1, 1}) {
    EXPECT_THAT(GetGltfAccessorSlice(model, index),
                StatusIs(absl::StatusCode::kInvalidArgument))
        << "accessor " << index;
  }

  tinygltf::Model sparse = model;
  sparse.accessors[0].sparse.isSparse = true;
  EXPECT_THAT(GetGltfAccessorSlice(sparse, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));
  // The overload taking an accessor reads its dense data.
  EXPECT_THAT(GetGltfAccessorSlice(sparse, sparse.accessors[0]), IsOk());

  tinygltf::Model missing_view = model;
  missing_view.accessors[0].bufferView = 1;
  EXPECT_THAT(GetGltfAccessorSlice(missing_view, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model missing_buffer = model;
  missing_buffer.bufferViews[0].buffer = 1;
  EXPECT_THAT(GetGltfAccessorSlice(missing_buffer, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model unknown_type = model;
  unknown_type.accessors[0].componentType = 0;
  EXPECT_THAT(GetGltfAccessorSlice(unknown_type, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Strides shorter than an element would make elements overlap.
  tinygltf::Model short_stride = model;
  short_stride.bufferViews[0].byteStride = sizeof(float);
  EXPECT_THAT(GetGltfAccessorSlice(short_stride, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfAccessorTest, RejectsOutOfBoundsAccessors) {
  const tinygltf::Model model = MakeQuadModel();

  tinygltf::Model too_many = model;
  too_many.accessors[0].count = 5;
  EXPECT_THAT(GetGltfAccessorSlice(too_many, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model offset_past_view = model;
  offset_past_view.accessors[0].byteOffset = sizeof(float);
  EXPECT_THAT(GetGltfAccessorSlice(offset_past_view, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  tinygltf::Model view_past_buffer = model;
  view_past_buffer.bufferViews[0].byteLength += 1;
  EXPECT_THAT(GetGltfAccessorSlice(view_past_buffer, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));

  // Large values that would overflow if the checks used additions.
  tinygltf::Model huge_count = model;
  huge_count.accessors[0].count = std::numeric_limits<size_t>::max();
  EXPECT_THAT(GetGltfAccessorSlice(huge_count, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));
  tinygltf::Model huge_offset = model;
  huge_offset.bufferViews[0].byteOffset = std::numeric_limits<size_t>::max();
  EXPECT_THAT(GetGltfAccessorSlice(huge_offset, 0),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfAccessorTest, RejectsMatricesWithPaddedColumns) {
  tinygltf::Model model;
  model.buffers.emplace_back();
  constexpr std::array<uint8_t, 8> kMat2 = {1, 2, 0, 0, 3, 4, 0, 0};
  const int acc =
      AddAccessor(model, kMat2, TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE,
                  TINYGLTF_TYPE_MAT2, 1);
  EXPECT_THAT(GetGltfAccessorSlice(model, acc),
              StatusIs(absl::StatusCode::kInvalidArgument));
}

TEST(GltfAccessorTest, RecognizesIndexAccessors) {
  tinygltf::Accessor acc;
  acc.type = TINYGLTF_TYPE_SCALAR;
  for (int component_type : {TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE,
                             TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT,
                             TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT}) {
    acc.componentType = component_type;
    EXPECT_THAT(IsGltfIndexAccessor(acc), IsTrue()) << component_type;
  }
  for (int component_type :
       {TINYGLTF_COMPONENT_TYPE_BYTE, TINYGLTF_COMPONENT_TYPE_SHORT,
        TINYGLTF_COMPONENT_TYPE_INT, TINYGLTF_COMPONENT_TYPE_FLOAT}) {
    acc.componentType = component_type;
    EXPECT_THAT(IsGltfIndexAccessor(acc), IsFalse()) << component_type;
  }
  acc.componentType = TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT;
  acc.type = TINYGLTF_TYPE_VEC2;
  EXPECT_THAT(IsGltfIndexAccessor(acc), IsFalse());
}

TEST(GltfAccessorTest, ReadsIndicesOfEveryType) {
  // Each index is stored in the first bytes of a 4-byte element.
  constexpr std::array<uint8_t, 8> kBytes = {7,    0,    0,    0,
                                             0xFF, 0xFF, 0xFF, 0x7F};
  const GltfAccessorSlice s{
      .data = kBytes.data(), .stride = 4, .elem_size = 4, .count = 2};
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE, 1),
              Eq(0xFFu));
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT, 1),
              Eq(0xFFFFu));
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT, 1),
              Eq(0x7FFFFFFFu));
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT, 0),
              Eq(7u));
  // Other types, and indices past the slice, read nothing and return 0.
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_FLOAT, 1), Eq(0u));
  EXPECT_THAT(ReadGltfIndex(s, TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT, 2),
              Eq(0u));
}

TEST(GltfAccessorTest, ReadsRawComponentsOfEveryType) {
  // Component 1 of each element, as stored (ignoring `normalized`).
  constexpr std::array<int8_t, 2> kI8 = {0, -5};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kI8.data()),
                           TINYGLTF_COMPONENT_TYPE_BYTE, 1),
      DoubleEq(-5));
  constexpr std::array<uint8_t, 2> kU8 = {0, 200};
  EXPECT_THAT(ReadGltfRawComponent(kU8.data(),
                                   TINYGLTF_COMPONENT_TYPE_UNSIGNED_BYTE, 1),
              DoubleEq(200));
  constexpr std::array<int16_t, 2> kI16 = {0, -300};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kI16.data()),
                           TINYGLTF_COMPONENT_TYPE_SHORT, 1),
      DoubleEq(-300));
  constexpr std::array<uint16_t, 2> kU16 = {0, 60000};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kU16.data()),
                           TINYGLTF_COMPONENT_TYPE_UNSIGNED_SHORT, 1),
      DoubleEq(60000));
  constexpr std::array<int32_t, 2> kI32 = {0, -70000};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kI32.data()),
                           TINYGLTF_COMPONENT_TYPE_INT, 1),
      DoubleEq(-70000));
  constexpr std::array<uint32_t, 2> kU32 = {0, 4000000000u};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kU32.data()),
                           TINYGLTF_COMPONENT_TYPE_UNSIGNED_INT, 1),
      DoubleEq(4000000000.0));
  constexpr std::array<float, 2> kF32 = {0.0f, 1.5f};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kF32.data()),
                           TINYGLTF_COMPONENT_TYPE_FLOAT, 1),
      DoubleEq(1.5));
  constexpr std::array<double, 2> kF64 = {0.0, 2.25};
  EXPECT_THAT(
      ReadGltfRawComponent(reinterpret_cast<const unsigned char*>(kF64.data()),
                           TINYGLTF_COMPONENT_TYPE_DOUBLE, 1),
      DoubleEq(2.25));
  // Unknown types read nothing and return NaN.
  EXPECT_THAT(std::isnan(ReadGltfRawComponent(kU8.data(), 0, 1)), IsTrue());
}

}  // namespace
}  // namespace intrinsic::geo
