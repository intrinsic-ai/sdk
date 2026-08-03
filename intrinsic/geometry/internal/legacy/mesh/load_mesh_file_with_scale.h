// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"

namespace intrinsic {
namespace geometry_legacy {

// Load a mesh object from a file, scale according to the given vector.
absl::StatusOr<Mesh> LoadMeshFileWithScale(absl::string_view path,
                                           const eigenmath::Vector3d& scale);

}  // namespace geometry_legacy
}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_LOAD_MESH_FILE_WITH_SCALE_H_
