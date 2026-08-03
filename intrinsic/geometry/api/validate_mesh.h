// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_
#define INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_

#include <string>

#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/eigenmath/types.h"

namespace intrinsic {

// Validates a mesh file (STL, GLB, DAE, etc.) to ensure it is valid and all its
// dependencies (e.g. textures) are present.
absl::Status ValidateMeshFile(absl::string_view filename,
                              const eigenmath::Vector3d& scale);

// Validates mesh data provided as a buffer.
// `extension` should be the file extension (e.g. "stl", "glb", "dae") to help
// identifying the format.
absl::Status ValidateMeshData(absl::string_view mesh_data,
                              absl::string_view extension,
                              const eigenmath::Vector3d& scale);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_VALIDATE_MESH_H_
