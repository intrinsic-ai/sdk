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

#include "intrinsic/geometry/gltf/gltf_materials.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "absl/strings/string_view.h"
#include "internal/testing.h"
#include "intrinsic/geometry/gltf/gltf_io.h"
#include "tiny_gltf.h"

namespace intrinsic::geo {
namespace {

using ::testing::DoubleEq;
using ::testing::ElementsAre;
using ::testing::Eq;

// Property key of the transmission factor in KHR_materials_transmission.
constexpr absl::string_view kTransmissionFactorKey = "transmissionFactor";

TEST(GltfMaterialsTest, GetGltfTransmissionFactorExtractsExpectedValue) {
  tinygltf::Material mat;
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.0));

  SetGltfTransmissionFactor(mat, 0.75);
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.75));

  // Overwrites an existing transmissionFactor in place.
  SetGltfTransmissionFactor(mat, 0.25);
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.25));

  mat.extensions[std::string(kKhrMaterialsTransmission)] =
      tinygltf::Value(tinygltf::Value::Object{
          {std::string(kTransmissionFactorKey), tinygltf::Value(1)}});
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(1.0));
}

TEST(GltfMaterialsTest, AddGltfExtensionUsedDeduplicatesEntries) {
  tinygltf::Model model;
  AddGltfExtensionUsed(model, kKhrMaterialsTransmission);
  AddGltfExtensionUsed(model, "KHR_mesh_quantization");
  AddGltfExtensionUsed(model, kKhrMaterialsTransmission);

  EXPECT_THAT(model.extensionsUsed,
              ElementsAre(std::string(kKhrMaterialsTransmission),
                          "KHR_mesh_quantization"));
}

TEST(GltfMaterialsTest, GetGltfTransmissionFactorHandlesMalformedExtension) {
  tinygltf::Material mat;
  mat.extensions[std::string(kKhrMaterialsTransmission)] = tinygltf::Value(
      tinygltf::Value::Object{{std::string(kTransmissionFactorKey),
                               tinygltf::Value(std::string("1.0"))}});
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.0));

  mat.extensions[std::string(kKhrMaterialsTransmission)] = tinygltf::Value();
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.0));
}

// glTF allows the extension without transmissionFactor, which defaults to 0.
TEST(GltfMaterialsTest, GetGltfTransmissionFactorDefaultsWhenFactorIsOmitted) {
  tinygltf::Material mat;
  mat.extensions[std::string(kKhrMaterialsTransmission)] =
      tinygltf::Value(tinygltf::Value::Object{});
  EXPECT_THAT(GetGltfTransmissionFactor(mat), DoubleEq(0.0));
}

class GltfTransmissionRoundtripTest : public ::testing::TestWithParam<double> {
};

// 1.0 also covers tinygltf writing integral doubles as integers ("1").
TEST_P(GltfTransmissionRoundtripTest, SurvivesSerialization) {
  const double transmission = GetParam();
  tinygltf::Model model;
  model.asset.version = "2.0";
  SetGltfTransmissionFactor(model.materials.emplace_back(), transmission);
  AddGltfExtensionUsed(model, kKhrMaterialsTransmission);
  ASSERT_OK_AND_ASSIGN(const std::string glb, SerializeGltfToGlb(model));
  ASSERT_OK_AND_ASSIGN(const tinygltf::Model reparsed, ParseGlb(glb));
  ASSERT_THAT(reparsed.materials.size(), Eq(1));
  EXPECT_THAT(GetGltfTransmissionFactor(reparsed.materials[0]),
              DoubleEq(transmission));
}

INSTANTIATE_TEST_SUITE_P(Factors, GltfTransmissionRoundtripTest,
                         ::testing::Values(0.0, 0.5, 1.0));

}  // namespace
}  // namespace intrinsic::geo
