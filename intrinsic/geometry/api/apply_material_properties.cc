// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/api/apply_material_properties.h"

#include <memory>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "assimp/color4.h"
#include "assimp/material.h"
#include "assimp/scene.h"
#include "google/protobuf/wrappers.pb.h"
#include "google/type/color.pb.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "intrinsic/util/status/status_macros.h"
#include "tiny_gltf.h"

namespace intrinsic {

using ::intrinsic_proto::geometry::v1::MaterialProperties;

namespace {
constexpr double kDefaultColor[4]{0.5, 0.5, 0.5, 1.0};
constexpr double kDefaultMetalness = 1.0;
constexpr double kDefaultRoughness = 1.0;
constexpr double kDefaultTransmission = 0.0;
}  // namespace

absl::StatusOr<std::shared_ptr<const Renderable>> ApplyMaterialProperties(
    std::shared_ptr<const Renderable> geo,
    const intrinsic_proto::geometry::v1::MaterialProperties&
        material_properties) {
  std::string glb_bytes = geo->GetGLBString();
  INTR_ASSIGN_OR_RETURN(
      glb_bytes, ApplyMaterialPropertiesToGlb(glb_bytes, material_properties));
  return std::make_shared<Renderable>(std::move(glb_bytes));
}

absl::StatusOr<std::string> ApplyMaterialPropertiesToGlb(
    absl::string_view glb_bytes,
    const MaterialProperties& material_properties) {
  // Parse the original glTF model.
  tinygltf::TinyGLTF gltf_parser;
  tinygltf::Model model;
  std::string err;
  std::string warn;
  if (!gltf_parser.LoadBinaryFromMemory(
          &model, &err, &warn,
          reinterpret_cast<const unsigned char*>(glb_bytes.data()),
          glb_bytes.size())) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Failed to load glTF bytes to glTF model. Error(s): ", err));
  }

  // Initialize material properties with default values.
  std::vector<double> base_color_factor;
  if (material_properties.has_base_color()) {
    const auto& base_color = material_properties.base_color();
    if (base_color.red() < 0.0 || base_color.red() > 1.0 ||
        base_color.green() < 0.0 || base_color.green() > 1.0 ||
        base_color.blue() < 0.0 || base_color.blue() > 1.0) {
      return absl::InvalidArgumentError(absl::Substitute(
          "Base color must be in the range of [0, 1]. Got [r=$0, g=$1, b=$2]",
          base_color.red(), base_color.green(), base_color.blue()));
    }
    if (base_color.has_alpha() && base_color.alpha().value() != 1.0) {
      return absl::InvalidArgumentError(
          "Base color alpha is not supported. Please use transmission in "
          "material properties instead.");
    }
    base_color_factor = {base_color.red(), base_color.green(),
                         base_color.blue(), 1.0};
  } else {
    base_color_factor = std::vector<double>(kDefaultColor, kDefaultColor + 4);
  }
  const double metalness = material_properties.has_metalness()
                               ? material_properties.metalness()
                               : kDefaultMetalness;
  if (metalness < 0.0 || metalness > 1.0) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Metalness must be in the range of [0, 1]. Got $0", metalness));
  }
  const double transmission = material_properties.has_transmission()
                                  ? material_properties.transmission()
                                  : kDefaultTransmission;
  if (transmission < 0.0 || transmission > 1.0) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Transmission must be in the range of [0, 1]. Got $0", transmission));
  }
  const double roughness = material_properties.has_roughness()
                               ? material_properties.roughness()
                               : kDefaultRoughness;
  if (roughness < 0.0 || roughness > 1.0) {
    return absl::InvalidArgumentError(absl::Substitute(
        "Roughness must be in the range of [0, 1]. Got $0", roughness));
  }

  // Clears existing materials and creates a new material. Replace all mesh
  // primitives to use this new material.
  tinygltf::Material material;
  material.name = "user_specified_material";
  material.pbrMetallicRoughness.baseColorFactor = std::move(base_color_factor);
  material.pbrMetallicRoughness.metallicFactor = metalness;
  material.pbrMetallicRoughness.roughnessFactor = roughness;
  material.extensions["KHR_materials_transmission"] =
      tinygltf::Value(tinygltf::Value::Object{
          {"transmissionFactor", tinygltf::Value(transmission)}});
  model.materials = {material};

  for (auto& mesh : model.meshes) {
    for (auto& primitive : mesh.primitives) {
      primitive.material = 0;
    }
  }

  // Write the updated glTF model to a string and return it.
  std::ostringstream ss;
  gltf_parser.WriteGltfSceneToStream(&model, ss, /*prettyPrint=*/false,
                                     /*writeBinary=*/true);
  return ss.str();
}

absl::Status ApplyMaterialToAiScene(aiScene& scene, const Material& material) {
  aiColor4D ambient(material.ambient[0], material.ambient[1],
                    material.ambient[2], material.ambient[3]);
  aiColor4D diffuse(material.diffuse[0], material.diffuse[1],
                    material.diffuse[2], material.diffuse[3]);
  aiColor4D specular(material.specular[0], material.specular[1],
                     material.specular[2], material.specular[3]);
  aiColor4D emission(material.emission[0], material.emission[1],
                     material.emission[2], material.emission[3]);

  // If we didn't have any materials, create a new one.
  if (scene.mNumMaterials == 0) {
    scene.mMaterials = new aiMaterial*[1];
    scene.mMaterials[0] = new aiMaterial();
    scene.mNumMaterials = 1;
  }

  // Apply the overrides to all materials in the scene.
  for (int i = 0; i < scene.mNumMaterials; ++i) {
    scene.mMaterials[i]->Clear();
    scene.mMaterials[i]->AddProperty(&ambient, 1, AI_MATKEY_COLOR_AMBIENT);

    scene.mMaterials[i]->AddProperty(&diffuse, 1, AI_MATKEY_COLOR_DIFFUSE);

    scene.mMaterials[i]->AddProperty(&specular, 1, AI_MATKEY_COLOR_SPECULAR);

    scene.mMaterials[i]->AddProperty(&emission, 1, AI_MATKEY_COLOR_EMISSIVE);
    scene.mMaterials[i]->AddProperty(&material.shininess, 1,
                                     AI_MATKEY_SHININESS);
  }

  return absl::OkStatus();
}

}  // namespace intrinsic
