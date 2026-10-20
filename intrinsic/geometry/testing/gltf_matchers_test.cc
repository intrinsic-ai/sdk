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

#include "intrinsic/geometry/testing/gltf_matchers.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <array>
#include <cstring>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/material_properties.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/internal/util/gltf_helpers.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "intrinsic/geometry/shapes/box.h"
#include "tiny_gltf.h"

namespace intrinsic::testing {
namespace {

using ::intrinsic::geo::MaterialProperties;
using ::testing::_;
using ::testing::Contains;
using ::testing::DescribeMatcher;
using ::testing::ElementsAre;
using ::testing::Eq;
using ::testing::ExplainMatchResult;
using ::testing::Ge;
using ::testing::HasSubstr;
using ::testing::IsEmpty;
using ::testing::Ne;
using ::testing::Not;
using ::testing::SizeIs;
using ::testing::StartsWith;
using ::testing::StringMatchResultListener;

MaterialProperties RedMaterial() {
  return {
      .color = {1.0, 0.0, 0.0},
      .metalness = 0.25,
      .roughness = 0.75,
      .transmission = 0.5,
  };
}

// Returns the explanation `matcher` gives for `value`.
template <typename M, typename T>
std::string Explain(const M& matcher, const T& value) {
  StringMatchResultListener listener;
  ExplainMatchResult(matcher, value, &listener);
  return listener.str();
}

tinygltf::Material MakeMaterial(const MaterialProperties& props,
                                double alpha = 1.0,
                                const std::string& name = "material") {
  tinygltf::Material material;
  material.name = name;
  material.pbrMetallicRoughness.baseColorFactor = {
      props.color.x(), props.color.y(), props.color.z(), alpha};
  material.pbrMetallicRoughness.metallicFactor = props.metalness;
  material.pbrMetallicRoughness.roughnessFactor = props.roughness;
  material.extensions[std::string(geo::kKhrMaterialsTransmission)] =
      tinygltf::Value(
          tinygltf::Value::Object{{std::string(geo::kTransmissionFactorKey),
                                   tinygltf::Value(props.transmission)}});
  return material;
}

// Returns a model with one triangle primitive per entry of
// `primitive_materials`, each referencing the given material index, and the
// given `materials`.
tinygltf::Model MakeModel(std::vector<tinygltf::Material> materials,
                          const std::vector<int>& primitive_materials = {0}) {
  constexpr std::array<float, 9> kVertices = {1.0f, 0.0f, 0.0f, 0.0f, 1.0f,
                                              0.0f, 0.0f, 0.0f, 1.0f};
  tinygltf::Model model;
  model.asset.version = "2.0";
  model.materials = std::move(materials);

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

  if (!primitive_materials.empty()) {
    tinygltf::Mesh& mesh = model.meshes.emplace_back();
    for (const int material : primitive_materials) {
      tinygltf::Primitive& primitive = mesh.primitives.emplace_back();
      primitive.attributes["POSITION"] = 0;
      primitive.material = material;
      primitive.mode = TINYGLTF_MODE_TRIANGLES;
    }
    model.nodes.emplace_back().mesh = 0;
    model.scenes.emplace_back().nodes = {0};
    model.defaultScene = 0;
  }
  return model;
}

std::string ToGlb(const tinygltf::Model& model) {
  absl::StatusOr<std::string> glb = geo::SerializeGltfToGlb(model);
  EXPECT_TRUE(glb.ok()) << glb.status();
  return glb.value_or("");
}

std::string RedGlb() { return ToGlb(MakeModel({MakeMaterial(RedMaterial())})); }

TEST(GltfMatchersSourceTest, AcceptsAllSourceTypes) {
  const std::string glb = RedGlb();
  const auto matcher = HasMaterialProperties(RedMaterial());

  EXPECT_THAT(glb, matcher);
  EXPECT_THAT(absl::string_view(glb), matcher);
  EXPECT_THAT(geo::Renderable(glb), matcher);
  EXPECT_THAT(std::make_shared<geo::Renderable>(glb), matcher);
  EXPECT_THAT(std::make_shared<const geo::Renderable>(glb), matcher);
  EXPECT_THAT(geo::Geometry(geo::ExactGeometry(
                                geo::Box(eigenmath::Vector3d(1.0, 1.0, 1.0))),
                            std::make_shared<const geo::Renderable>(glb)),
              matcher);
  EXPECT_THAT(MakeModel({MakeMaterial(RedMaterial())}), matcher);
  EXPECT_THAT(MakeMaterial(RedMaterial()), matcher);
}

TEST(GltfMatchersSourceTest, ExplainsUnusableSources) {
  const auto matcher = HasGltfMaterialsCount(1);
  EXPECT_THAT(Explain(matcher, std::string()), Eq("GLB data is empty"));
  EXPECT_THAT(Explain(matcher, std::string("not a glb")),
              HasSubstr("Invalid GLB magic bytes"));
  EXPECT_THAT(Explain(matcher, std::shared_ptr<const geo::Renderable>()),
              Eq("Renderable is null"));
  EXPECT_THAT(Explain(matcher, geo::Geometry()),
              Eq("Geometry has no renderable"));

  EXPECT_THAT(std::string(), Not(matcher));
  EXPECT_THAT(std::shared_ptr<geo::Renderable>(), Not(matcher));
  EXPECT_THAT(geo::Geometry(), Not(matcher));
}

TEST(HasMaterialPropertiesTest, MatchesWithinTolerance) {
  MaterialProperties close = RedMaterial();
  close.metalness += 5e-5;
  EXPECT_THAT(MakeMaterial(RedMaterial()), HasMaterialProperties(close));

  close.metalness += 1e-3;
  EXPECT_THAT(MakeMaterial(RedMaterial()), Not(HasMaterialProperties(close)));
  EXPECT_THAT(MakeMaterial(RedMaterial()), HasMaterialProperties(close, 1e-2));
}

TEST(HasMaterialPropertiesTest, AcceptsProto) {
  const intrinsic_proto::geometry::v1::MaterialProperties proto =
      RedMaterial().ToProto();
  EXPECT_THAT(RedGlb(), HasMaterialProperties(proto));
  EXPECT_THAT(RedGlb(), HasMaterialProperties(proto, 1e-6));
}

TEST(HasMaterialPropertiesTest, ExplainsEachMismatchedProperty) {
  const tinygltf::Material material = MakeMaterial(RedMaterial());

  MaterialProperties expected = RedMaterial();
  expected.color = {0.0, 1.0, 0.0};
  EXPECT_THAT(Explain(HasMaterialProperties(expected), material),
              HasSubstr("baseColorFactor"));

  expected = RedMaterial();
  expected.metalness = 0.0;
  EXPECT_THAT(Explain(HasMaterialProperties(expected), material),
              HasSubstr("metallicFactor"));

  expected = RedMaterial();
  expected.roughness = 0.0;
  EXPECT_THAT(Explain(HasMaterialProperties(expected), material),
              HasSubstr("roughnessFactor"));

  expected = RedMaterial();
  expected.transmission = 0.0;
  EXPECT_THAT(Explain(HasMaterialProperties(expected), material),
              HasSubstr("transmission"));
}

TEST(HasMaterialPropertiesTest, RequiresOpaqueAlpha) {
  EXPECT_THAT(MakeMaterial(RedMaterial(), /*alpha=*/0.5),
              Not(HasMaterialProperties(RedMaterial())));
}

TEST(HasMaterialPropertiesTest, TreatsMissingTransmissionAsZero) {
  tinygltf::Material material = MakeMaterial(MaterialProperties{});
  material.extensions.clear();
  EXPECT_THAT(material, HasMaterialProperties(MaterialProperties{}));
}

TEST(HasMaterialPropertiesTest, ExplainsModelWithoutMaterials) {
  EXPECT_THAT(Explain(HasMaterialProperties(RedMaterial()), MakeModel({})),
              Eq("Model has no materials"));
}

TEST(HasGltfNameTest, MatchesNamedElements) {
  tinygltf::Node node;
  node.name = "base_link";
  tinygltf::Scene scene;
  scene.name = "scene";

  EXPECT_THAT(node, HasGltfName("base_link"));
  EXPECT_THAT(node, HasGltfName(StartsWith("base")));
  EXPECT_THAT(node, Not(HasGltfName("other")));
  EXPECT_THAT(scene, HasGltfName("scene"));
  EXPECT_THAT(MakeMaterial(RedMaterial(), 1.0, "red"), HasGltfName("red"));
}

TEST(HasFirstMaterialNameTest, MatchesFirstMaterialName) {
  const std::string glb =
      ToGlb(MakeModel({MakeMaterial(RedMaterial(), 1.0, "a"),
                       MakeMaterial(RedMaterial(), 1.0, "b")}));
  EXPECT_THAT(glb, HasFirstMaterialName("a"));
  EXPECT_THAT(glb, Not(HasFirstMaterialName("b")));
}

TEST(AllPrimitivesHaveMaterialTest, MatchesWhenAllPrimitivesShareMaterial) {
  const auto materials = std::vector<tinygltf::Material>{
      MakeMaterial(RedMaterial()), MakeMaterial(RedMaterial())};
  EXPECT_THAT(MakeModel(materials, {1, 1}), AllPrimitivesHaveMaterial(1));
  EXPECT_THAT(MakeModel(materials, {0, 1}), Not(AllPrimitivesHaveMaterial(0)));
  EXPECT_THAT(
      Explain(AllPrimitivesHaveMaterial(0), MakeModel(materials, {0, 1})),
      Eq("Mesh 0 primitive 1 has material 1 != expected 0"));
}

TEST(AllPrimitivesHaveMaterialTest, ExplainsModelWithoutPrimitives) {
  const tinygltf::Model model =
      MakeModel({MakeMaterial(RedMaterial())}, /*primitive_materials=*/{});
  EXPECT_THAT(Explain(AllPrimitivesHaveMaterial(0), model),
              Eq("Model has no mesh primitives"));
}

TEST(HasGltfMaterialsCountTest, MatchesCount) {
  const std::string glb = ToGlb(
      MakeModel({MakeMaterial(RedMaterial()), MakeMaterial(RedMaterial())}));
  EXPECT_THAT(glb, HasGltfMaterialsCount(2));
  EXPECT_THAT(glb, HasGltfMaterialsCount(Ge(1)));
  EXPECT_THAT(glb, Not(HasGltfMaterialsCount(1)));
}

TEST(HasGltfBaseColorTest, MatchesFirstMaterialBaseColor) {
  const tinygltf::Material material = MakeMaterial(RedMaterial(), 0.5);
  EXPECT_THAT(material, HasGltfBaseColor(1.0, 0.0, 0.0, 0.5));
  EXPECT_THAT(material, Not(HasGltfBaseColor(1.0, 0.0, 0.0)));
  EXPECT_THAT(material, HasGltfBaseColor(1.0, 0.0, 0.0, 0.6, /*eps=*/0.2));
}

TEST(HasGltfBaseColorTest, ExplainsShortBaseColorFactor) {
  tinygltf::Material material = MakeMaterial(RedMaterial());
  material.pbrMetallicRoughness.baseColorFactor = {1.0, 0.0};
  EXPECT_THAT(material, Not(HasGltfBaseColor(1.0, 0.0, 0.0)));
  EXPECT_THAT(Explain(HasGltfBaseColor(1.0, 0.0, 0.0), material),
              HasSubstr("baseColorFactor (1.000000, 0.000000, 0.000000, "
                        "0.000000)"));
}

TEST(HasGltfBaseColorFactorTest, AppliesMatcherToBaseColorFactor) {
  const tinygltf::Material material = MakeMaterial(RedMaterial());
  EXPECT_THAT(material,
              HasGltfBaseColorFactor(ElementsAre(1.0, 0.0, 0.0, 1.0)));
  EXPECT_THAT(material, HasGltfBaseColorFactor(
                            ElementsAre(Ne(0.5), Ne(0.5), Ne(0.5), _)));
  EXPECT_THAT(material, Not(HasGltfBaseColorFactor(Contains(0.5))));
  EXPECT_THAT(RedGlb(), HasGltfBaseColorFactor(ElementsAre(1.0, _, _, _)));
}

TEST(AllMaterialsHaveBaseColorTest, ChecksEveryMaterial) {
  const tinygltf::Material red = MakeMaterial(RedMaterial(), 1.0, "red");
  MaterialProperties green_props = RedMaterial();
  green_props.color = {0.0, 1.0, 0.0};
  const tinygltf::Material green = MakeMaterial(green_props, 1.0, "green");

  EXPECT_THAT(MakeModel({red, red}), AllMaterialsHaveBaseColor(1.0, 0.0, 0.0));
  EXPECT_THAT(MakeModel({red, green}),
              Not(AllMaterialsHaveBaseColor(1.0, 0.0, 0.0)));
  EXPECT_THAT(Explain(AllMaterialsHaveBaseColor(1.0, 0.0, 0.0),
                      MakeModel({red, green})),
              HasSubstr("(at material index 1, name: \"green\")"));
  EXPECT_THAT(Explain(AllMaterialsHaveBaseColor(1.0, 0.0, 0.0), MakeModel({})),
              Eq("Model has no materials"));
}

TEST(GltfExtensionMatchersTest, MatchExtensionLists) {
  tinygltf::Model model = MakeModel({MakeMaterial(RedMaterial())});
  model.extensionsUsed = {"KHR_mesh_quantization", "EXT_meshopt_compression"};
  model.extensionsRequired = {"KHR_mesh_quantization"};
  const std::string glb = ToGlb(model);

  EXPECT_THAT(glb, HasGltfExtensionUsed("EXT_meshopt_compression"));
  EXPECT_THAT(glb, HasGltfExtensionRequired("KHR_mesh_quantization"));
  EXPECT_THAT(glb, Not(HasGltfExtensionRequired("EXT_meshopt_compression")));
  EXPECT_THAT(glb, HasGltfExtensionsUsed(SizeIs(2)));
  EXPECT_THAT(MakeModel({}), HasGltfExtensionsUsed(IsEmpty()));
  EXPECT_THAT(MakeModel({}), HasGltfExtensionsRequired(IsEmpty()));
}

TEST(GltfExtensionMatchersTest, DescribesNegation) {
  const auto matcher = HasGltfExtensionsUsed(IsEmpty());
  EXPECT_THAT(DescribeMatcher<tinygltf::Model>(matcher),
              Eq("has glTF extensionsUsed that is empty"));
  EXPECT_THAT(DescribeMatcher<tinygltf::Model>(matcher, /*negation=*/true),
              Eq("has glTF extensionsUsed that isn't empty"));
}

TEST(GltfMatchersDescriptionTest, DescribesNames) {
  EXPECT_THAT(DescribeMatcher<tinygltf::Node>(HasGltfName("a")),
              Eq("has glTF name that is equal to \"a\""));
  EXPECT_THAT(
      DescribeMatcher<tinygltf::Node>(HasGltfName("a"), /*negation=*/true),
      Eq("has glTF name that isn't equal to \"a\""));
  EXPECT_THAT(
      DescribeMatcher<std::string>(HasFirstMaterialName(StartsWith("a")),
                                   /*negation=*/true),
      Eq("has first material name that doesn't start with \"a\""));
}

TEST(GltfMatchersDescriptionTest, DescribesMaterialsCount) {
  EXPECT_THAT(DescribeMatcher<std::string>(HasGltfMaterialsCount(1)),
              Eq("has materials count that is equal to 1"));
  EXPECT_THAT(DescribeMatcher<std::string>(HasGltfMaterialsCount(Ge(2)),
                                           /*negation=*/true),
              Eq("has materials count that isn't >= 2"));
}

TEST(GltfMatchersDescriptionTest, DescribesValueMatchers) {
  EXPECT_THAT(DescribeMatcher<std::string>(HasGltfBaseColor(1, 0, 0)),
              StartsWith("has base color [1.000000, 0.000000"));
  EXPECT_THAT(DescribeMatcher<std::string>(HasGltfBaseColor(1, 0, 0),
                                           /*negation=*/true),
              StartsWith("doesn't have base color [1.000000, 0.000000"));
  EXPECT_THAT(DescribeMatcher<std::string>(AllMaterialsHaveBaseColor(1, 0, 0)),
              StartsWith("all materials have base color"));
  EXPECT_THAT(DescribeMatcher<std::string>(AllMaterialsHaveBaseColor(1, 0, 0),
                                           /*negation=*/true),
              StartsWith("not all materials have base color"));
  EXPECT_THAT(DescribeMatcher<std::string>(AllPrimitivesHaveMaterial(0)),
              Eq("all primitives reference material index 0"));
  EXPECT_THAT(DescribeMatcher<std::string>(AllPrimitivesHaveMaterial(0),
                                           /*negation=*/true),
              Eq("not all primitives reference material index 0"));
  EXPECT_THAT(
      DescribeMatcher<std::string>(HasMaterialProperties(RedMaterial())),
      StartsWith("has material properties (color=[1.000000"));
  EXPECT_THAT(DescribeMatcher<std::string>(HasMaterialProperties(RedMaterial()),
                                           /*negation=*/true),
              StartsWith("doesn't have material properties (color=[1.000000"));
  EXPECT_THAT(
      DescribeMatcher<std::string>(HasGltfBaseColorFactor(Contains(0.5)),
                                   /*negation=*/true),
      Eq("has first material baseColorFactor that doesn't contain any element "
         "that is equal to 0.5"));
}

}  // namespace
}  // namespace intrinsic::testing
