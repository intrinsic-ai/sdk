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

#ifndef INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_MESH_FROM_FILE_H_
#define INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_MESH_FROM_FILE_H_

#include <string>
#include <utility>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"

namespace intrinsic::geo {

// Reads a mesh from a file
// File type must be readable by assimp
absl::StatusOr<Mesh> LoadMeshFromFile(
    const std::string& filename,
    const eigenmath::Vector3d& scale = eigenmath::Vector3d::Ones());

inline Mesh LoadMeshFromFileOrDie(
    const std::string& filename,
    const eigenmath::Vector3d& scale = eigenmath::Vector3d::Ones()) {
  auto status_or_mesh = LoadMeshFromFile(filename, scale);
  if (!status_or_mesh.ok()) {
    LOG(ERROR) << "Failed to load mesh from: " << filename;
    LOG(ERROR) << "Returned status: " << status_or_mesh.status();
  }

  return std::move(status_or_mesh.value());
}

}  // namespace intrinsic::geo
#endif  // INTRINSIC_GEOMETRY_INTERNAL_LEGACY_MESH_IO_LOAD_MESH_FROM_FILE_H_
