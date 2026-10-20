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

#ifndef INTRINSIC_GEOMETRY_TESTING_GLTF_MATCHERS_H_
#define INTRINSIC_GEOMETRY_TESTING_GLTF_MATCHERS_H_

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <cmath>
#include <concepts>
#include <cstddef>
#include <memory>
#include <string>
#include <type_traits>
#include <vector>

#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/material_properties.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/internal/util/gltf_helpers.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"
#include "tiny_gltf.h"

namespace intrinsic::testing {

namespace internal {

template <typename T>
concept GltfNamedElement = requires(const std::decay_t<T>& v) {
  { v.name } -> std::convertible_to<absl::string_view>;
};

template <typename T>
concept GltfByteSource =
    std::convertible_to<std::decay_t<T>, absl::string_view> ||
    std::same_as<std::decay_t<T>, std::shared_ptr<const geo::Renderable>> ||
    std::same_as<std::decay_t<T>, std::shared_ptr<geo::Renderable>> ||
    std::same_as<std::decay_t<T>, geo::Renderable> ||
    std::same_as<std::decay_t<T>, geo::Geometry>;

template <typename T>
concept GltfModelSource =
    std::same_as<std::decay_t<T>, tinygltf::Model> || GltfByteSource<T>;

template <typename T>
concept GltfMaterialSource =
    std::same_as<std::decay_t<T>, tinygltf::Material> || GltfModelSource<T>;

template <GltfModelSource T, typename Fn>
bool WithModel(const T& arg, ::testing::MatchResultListener* listener,
               Fn&& fn) {
  using Decayed = std::decay_t<T>;
  if constexpr (std::same_as<Decayed, tinygltf::Model>) {
    return fn(arg);
  } else {
    auto parse_and_run = [&](absl::string_view glb_bytes) -> bool {
      if (glb_bytes.empty()) {
        *listener << "GLB data is empty";
        return false;
      }
      geo::GltfParseDiagnostics diagnostics;
      absl::StatusOr<tinygltf::Model> model =
          geo::ParseGlb(glb_bytes, &diagnostics);
      if (!model.ok()) {
        *listener << model.status().message();
        return false;
      }
      if (!diagnostics.errors.empty()) {
        *listener << "glTF parse reported errors: " << diagnostics.errors;
        return false;
      }
      return fn(*model);
    };

    if constexpr (std::convertible_to<Decayed, absl::string_view>) {
      return parse_and_run(absl::string_view(arg));
    } else if constexpr (std::same_as<Decayed,
                                      std::shared_ptr<const geo::Renderable>> ||
                         std::same_as<Decayed,
                                      std::shared_ptr<geo::Renderable>>) {
      if (arg == nullptr) {
        *listener << "Renderable is null";
        return false;
      }
      const std::string glb = arg->GetGLBString();
      return parse_and_run(glb);
    } else if constexpr (std::same_as<Decayed, geo::Renderable>) {
      const std::string glb = arg.GetGLBString();
      return parse_and_run(glb);
    } else if constexpr (std::same_as<Decayed, geo::Geometry>) {
      if (arg.GetRenderable() == nullptr) {
        *listener << "Geometry has no renderable";
        return false;
      }
      const std::string glb = arg.GetRenderable()->GetGLBString();
      return parse_and_run(glb);
    }
  }
}

template <GltfMaterialSource T, typename Fn>
bool WithFirstMaterial(const T& arg, ::testing::MatchResultListener* listener,
                       Fn&& fn) {
  using Decayed = std::decay_t<T>;
  if constexpr (std::same_as<Decayed, tinygltf::Material>) {
    return fn(arg);
  } else {
    return WithModel(arg, listener, [&](const tinygltf::Model& model) {
      if (model.materials.empty()) {
        *listener << "Model has no materials";
        return false;
      }
      return fn(model.materials[0]);
    });
  }
}

inline bool MatchBaseColor(const tinygltf::Material& mat, double r, double g,
                           double b, double a, double eps,
                           ::testing::MatchResultListener* listener) {
  const auto& pbr = mat.pbrMetallicRoughness;
  if (pbr.baseColorFactor.size() < 4 ||
      std::abs(pbr.baseColorFactor[0] - r) > eps ||
      std::abs(pbr.baseColorFactor[1] - g) > eps ||
      std::abs(pbr.baseColorFactor[2] - b) > eps ||
      std::abs(pbr.baseColorFactor[3] - a) > eps) {
    *listener << absl::StrFormat(
        "Material baseColorFactor (%f, %f, %f, %f) != expected (%f, %f, %f, "
        "%f)",
        pbr.baseColorFactor.size() > 0 ? pbr.baseColorFactor[0] : 0.0,
        pbr.baseColorFactor.size() > 1 ? pbr.baseColorFactor[1] : 0.0,
        pbr.baseColorFactor.size() > 2 ? pbr.baseColorFactor[2] : 0.0,
        pbr.baseColorFactor.size() > 3 ? pbr.baseColorFactor[3] : 0.0, r, g, b,
        a);
    return false;
  }
  return true;
}

inline bool MatchMaterial(const tinygltf::Material& mat,
                          const geo::MaterialProperties& expected_material,
                          double eps,
                          ::testing::MatchResultListener* listener) {
  if (!MatchBaseColor(mat, expected_material.color.x(),
                      expected_material.color.y(), expected_material.color.z(),
                      /*a=*/1.0, eps, listener)) {
    return false;
  }

  const auto& pbr = mat.pbrMetallicRoughness;

  if (std::abs(pbr.metallicFactor - expected_material.metalness) > eps) {
    *listener << absl::StrFormat("Material metallicFactor %f != expected %f",
                                 pbr.metallicFactor,
                                 expected_material.metalness);
    return false;
  }

  if (std::abs(pbr.roughnessFactor - expected_material.roughness) > eps) {
    *listener << absl::StrFormat("Material roughnessFactor %f != expected %f",
                                 pbr.roughnessFactor,
                                 expected_material.roughness);
    return false;
  }

  const double transmission = geo::GetGltfTransmissionFactor(mat);
  if (std::abs(transmission - expected_material.transmission) > eps) {
    *listener << absl::StrFormat("Material transmission %f != expected %f",
                                 transmission, expected_material.transmission);
    return false;
  }

  return true;
}

}  // namespace internal

// Polymorphic matcher that verifies a GLB (passed as string, Renderable,
// shared_ptr<Renderable>, Geometry, Model) or a tinygltf::Material has the
// specified MaterialProperties.
MATCHER_P2(HasMaterialProperties, expected_material, eps,
           absl::StrFormat(
               "%s material properties (color=[%f, %f, %f], metalness=%f, "
               "roughness=%f, transmission=%f)",
               negation ? "doesn't have" : "has", expected_material.color.x(),
               expected_material.color.y(), expected_material.color.z(),
               expected_material.metalness, expected_material.roughness,
               expected_material.transmission)) {
  return internal::WithFirstMaterial(
      arg, result_listener, [&](const tinygltf::Material& mat) {
        return internal::MatchMaterial(mat, expected_material, eps,
                                       result_listener);
      });
}

// Overload for default tolerance eps = 1e-4.
inline auto HasMaterialProperties(
    const geo::MaterialProperties& expected_material) {
  return HasMaterialProperties(expected_material, 1e-4);
}

// Overload taking proto MaterialProperties.
inline auto HasMaterialProperties(
    const intrinsic_proto::geometry::v1::MaterialProperties& proto,
    double eps = 1e-4) {
  return HasMaterialProperties(geo::MaterialProperties::FromProto(proto), eps);
}

// Matcher that checks the glTF name of a named glTF element (Material, Node,
// Mesh, Scene).
MATCHER_P(HasGltfName, name_matcher,
          absl::StrCat("has glTF name that ",
                       ::testing::DescribeMatcher<std::string>(name_matcher,
                                                               negation))) {
  static_assert(
      internal::GltfNamedElement<decltype(arg)>,
      "HasGltfName requires a named glTF element (Material, Node, Mesh, "
      "Scene); use HasFirstMaterialName for Model / Renderable / GLB sources.");
  const ::testing::Matcher<std::string> matcher = name_matcher;
  return matcher.MatchAndExplain(arg.name, result_listener);
}

// Matcher that checks the name of the first material in a Material, Model,
// Renderable, Geometry, or GLB byte string.
MATCHER_P(HasFirstMaterialName, name_matcher,
          absl::StrCat("has first material name that ",
                       ::testing::DescribeMatcher<std::string>(name_matcher,
                                                               negation))) {
  const ::testing::Matcher<std::string> matcher = name_matcher;
  return internal::WithFirstMaterial(
      arg, result_listener, [&](const tinygltf::Material& mat) {
        return matcher.MatchAndExplain(mat.name, result_listener);
      });
}

// Matcher that verifies all mesh primitives in a GLB reference the specified
// material index.
MATCHER_P(AllPrimitivesHaveMaterial, expected_material_index,
          absl::StrCat(negation ? "not all" : "all",
                       " primitives reference material index ",
                       expected_material_index)) {
  return internal::WithModel(
      arg, result_listener, [&](const tinygltf::Model& model) {
        bool has_primitives = false;
        for (size_t i = 0; i < model.meshes.size(); ++i) {
          const auto& primitives = model.meshes[i].primitives;
          for (size_t j = 0; j < primitives.size(); ++j) {
            has_primitives = true;
            if (primitives[j].material != expected_material_index) {
              *result_listener << absl::StrFormat(
                  "Mesh %d primitive %d has material %d != expected %d", i, j,
                  primitives[j].material, expected_material_index);
              return false;
            }
          }
        }
        if (!has_primitives) {
          *result_listener << "Model has no mesh primitives";
          return false;
        }
        return true;
      });
}

// Matcher that checks the number of materials in the GLB.
MATCHER_P(HasGltfMaterialsCount, matcher,
          absl::StrCat("has materials count that ",
                       ::testing::DescribeMatcher<int>(matcher, negation))) {
  const ::testing::Matcher<int> count_matcher = matcher;
  return internal::WithModel(
      arg, result_listener, [&](const tinygltf::Model& model) {
        return count_matcher.MatchAndExplain(
            static_cast<int>(model.materials.size()), result_listener);
      });
}

// Matcher that checks base color factor RGBA values.
MATCHER_P5(HasGltfBaseColor, r, g, b, a, eps,
           absl::StrFormat("%s base color [%f, %f, %f, %f]",
                           negation ? "doesn't have" : "has", r, g, b, a)) {
  return internal::WithFirstMaterial(
      arg, result_listener, [&](const tinygltf::Material& mat) {
        return internal::MatchBaseColor(mat, r, g, b, a, eps, result_listener);
      });
}

inline auto HasGltfBaseColor(double r, double g, double b, double a = 1.0) {
  return HasGltfBaseColor(r, g, b, a, 1e-4);
}

// Matcher that applies `matcher` to the RGBA baseColorFactor, as
// std::vector<double>, of the first material. Useful for per-channel checks
// that HasGltfBaseColor() cannot express, e.g.
// HasGltfBaseColorFactor(ElementsAre(Ne(0.5), Ne(0.5), Ne(0.5), _)).
MATCHER_P(HasGltfBaseColorFactor, matcher,
          absl::StrCat("has first material baseColorFactor that ",
                       ::testing::DescribeMatcher<std::vector<double>>(
                           matcher, negation))) {
  return internal::WithFirstMaterial(
      arg, result_listener, [&](const tinygltf::Material& mat) {
        return ::testing::ExplainMatchResult(
            matcher, mat.pbrMetallicRoughness.baseColorFactor, result_listener);
      });
}

// Matcher that checks that all materials in the glTF model have the specified
// base color factor RGBA values.
MATCHER_P5(AllMaterialsHaveBaseColor, r, g, b, a, eps,
           absl::StrFormat("%s materials have base color [%f, %f, %f, %f]",
                           negation ? "not all" : "all", r, g, b, a)) {
  return internal::WithModel(
      arg, result_listener, [&](const tinygltf::Model& model) {
        if (model.materials.empty()) {
          *result_listener << "Model has no materials";
          return false;
        }
        for (size_t i = 0; i < model.materials.size(); ++i) {
          if (!internal::MatchBaseColor(model.materials[i], r, g, b, a, eps,
                                        result_listener)) {
            *result_listener
                << absl::StrFormat(" (at material index %d, name: \"%s\")", i,
                                   model.materials[i].name);
            return false;
          }
        }
        return true;
      });
}

inline auto AllMaterialsHaveBaseColor(double r, double g, double b,
                                      double a = 1.0) {
  return AllMaterialsHaveBaseColor(r, g, b, a, 1e-4);
}

// Matchers for glTF extensions. `matcher` is applied to the list of extension
// names as std::vector<std::string>.
MATCHER_P(HasGltfExtensionsUsed, matcher,
          absl::StrCat("has glTF extensionsUsed that ",
                       ::testing::DescribeMatcher<std::vector<std::string>>(
                           matcher, negation))) {
  return internal::WithModel(
      arg, result_listener, [&](const tinygltf::Model& model) {
        return ::testing::ExplainMatchResult(matcher, model.extensionsUsed,
                                             result_listener);
      });
}

MATCHER_P(HasGltfExtensionsRequired, matcher,
          absl::StrCat("has glTF extensionsRequired that ",
                       ::testing::DescribeMatcher<std::vector<std::string>>(
                           matcher, negation))) {
  return internal::WithModel(
      arg, result_listener, [&](const tinygltf::Model& model) {
        return ::testing::ExplainMatchResult(matcher, model.extensionsRequired,
                                             result_listener);
      });
}

inline auto HasGltfExtensionUsed(absl::string_view extension_name) {
  return HasGltfExtensionsUsed(
      ::testing::Contains(std::string(extension_name)));
}

inline auto HasGltfExtensionRequired(absl::string_view extension_name) {
  return HasGltfExtensionsRequired(
      ::testing::Contains(std::string(extension_name)));
}

}  // namespace intrinsic::testing

#endif  // INTRINSIC_GEOMETRY_TESTING_GLTF_MATCHERS_H_
