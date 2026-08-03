// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/geometry/internal/legacy/mesh/load_mesh_file_with_scale.h"

#include <string>

#include "absl/log/check.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/legacy/mesh/io/load_mesh_from_file.h"
#include "intrinsic/geometry/internal/legacy/mesh/mesh.h"
#include "intrinsic/scene/sdf/sdf_path_resolver.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace geometry_legacy {

absl::StatusOr<Mesh> LoadMeshFileWithScale(absl::string_view path,
                                           const eigenmath::Vector3d& scale) {
  INTR_ASSIGN_OR_RETURN(auto full_path,
                        sdf::SdfPathResolver(std::string(path)));
  return LoadMeshFromFile(full_path, scale);
}

}  // namespace geometry_legacy
}  // namespace intrinsic
