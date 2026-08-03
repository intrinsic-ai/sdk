// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_APPLY_MATERIAL_PROPERTIES_H_
#define INTRINSIC_GEOMETRY_API_APPLY_MATERIAL_PROPERTIES_H_

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "assimp/scene.h"
#include "intrinsic/geometry/api/material.h"
#include "intrinsic/geometry/api/renderable.h"
#include "intrinsic/geometry/proto/v1/material.pb.h"

namespace intrinsic {

absl::StatusOr<std::shared_ptr<const Renderable>> ApplyMaterialProperties(
    std::shared_ptr<const Renderable> geo,
    const intrinsic_proto::geometry::v1::MaterialProperties&
        material_properties);

absl::StatusOr<std::string> ApplyMaterialPropertiesToGlb(
    absl::string_view glb_bytes,
    const intrinsic_proto::geometry::v1::MaterialProperties&
        material_properties);

absl::Status ApplyMaterialToAiScene(aiScene& scene, const Material& material);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_APPLY_MATERIAL_PROPERTIES_H_
