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

#ifndef INTRINSIC_GEOMETRY_API_RENDERABLE_GENERATION_H_
#define INTRINSIC_GEOMETRY_API_RENDERABLE_GENERATION_H_

#include <memory>

#include "absl/container/flat_hash_map.h"
#include "absl/status/statusor.h"
#include "assimp/scene.h"
#include "intrinsic/geometry/api/exact_geometry.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/geometry/api/material_properties.h"
#include "intrinsic/geometry/api/renderable.h"

namespace intrinsic::geo {

// Recommended usage of these functions:
//
// 1. For mesh geometry, we should generate a Renderable. It will contain
// the mesh vertex data and optionally any material info. Note that the
// ExactGeometry proto will also contain the TriangleMesh vertex data, which
// is used by parts of the system like motion planning that only need the
// vertex data.
//
// 2. For primitives (like cubes, spheres), we should avoid generating a
// Renderable unless the primitive has a material. If it has a material,
// we should generate the Renderable, which will contain the shape data
// converted to a mesh, along with the material info.
//

// Returns the renderable representing this Geometry object. Optionally
// creating it if it was not previously available. This may be a lossy
// conversion.
absl::StatusOr<std::shared_ptr<const Renderable>> GetOrGenerateRenderable(
    const Geometry& geo);

// Returns the renderable representing this Geometry object. Optionally
// creating it if it was not previously available. This may be a lossy
// conversion.
//
// The difference between this and GetOrGenerateRenderable is that this
// function will use the material overrides stored within the Geometry to
// generate the renderable. If no material overrides are present, this is
// equivalent to calling GetOrGenerateRenderable.
absl::StatusOr<std::shared_ptr<const Renderable>>
GenerateRenderableWithMaterialOverrides(const Geometry& geo);

// Generates a renerable for the given ExactGeometry, with optional PBR Material
// Properties. This should be called for meshes and also for primitive shapes
// with materials. For primitive shapes without materials, there is no need to
// generate the renderable like this.
absl::StatusOr<std::shared_ptr<const Renderable>>
GenerateRenderableForExactGeometry(
    const ExactGeometry& exact_geometry,
    const std::optional<intrinsic_proto::geometry::v1::MaterialProperties>&
        material_properties);

// Ensures that the returned Geometry has a renderable representation. If the
// input already had a renderable no action is taken, and the input is simply
// returned, if the input did not have a renderable, a new one is generated and
// returned with an updated Geometry instance.
absl::StatusOr<Geometry> EnsureRenderableIsAvailable(Geometry geo);

// Helper class to generate a Renderable containing multiple meshes and
// materials. Use for cases where the Renderable should contain >1 mesh, like
// when importing USD assets where the main mesh may have sub-meshes.
class RenderableGenerator {
 public:
  RenderableGenerator() = default;

  // Adds a geometry with an optional material to the Renderable
  absl::Status AddGeometry(
      const ExactGeometry& exact_geometry,
      const std::optional<MaterialProperties>& material_properties);

  // Generate the final Renderable
  absl::StatusOr<std::shared_ptr<const Renderable>> Finish();

 private:
  void AddMesh(const Mesh& mesh,
               const std::optional<MaterialProperties>& material_properties);
  size_t GetOrAddMaterial(const MaterialProperties& material);

  std::vector<std::unique_ptr<aiMesh>> meshes_;
  std::vector<std::unique_ptr<aiMaterial>> materials_;
  absl::flat_hash_map<MaterialProperties, size_t> material_to_index_;
};

}  // namespace intrinsic::geo

namespace intrinsic {
using ::intrinsic::geo::EnsureRenderableIsAvailable;
using ::intrinsic::geo::GenerateRenderableWithMaterialOverrides;
using ::intrinsic::geo::RenderableGenerator;
}  // namespace intrinsic
#endif  // INTRINSIC_GEOMETRY_API_RENDERABLE_GENERATION_H_
