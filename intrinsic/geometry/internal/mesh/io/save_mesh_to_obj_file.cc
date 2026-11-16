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

#include "intrinsic/geometry/internal/mesh/io/save_mesh_to_obj_file.h"

#include <fstream>
#include <ostream>
#include <string>

#include "Eigen/Core"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/geometry/internal/mesh/mesh.h"
#include "intrinsic/util/full_precision.h"
#include "intrinsic/util/status/status_macros.h"
#include "ortools/base/file.h"
#include "ortools/base/options.h"

namespace intrinsic::geo {

absl::Status SaveMeshToObjFile(absl::string_view filename, const Mesh& mesh) {
  std::ofstream ostream(static_cast<const std::string>(filename));

  for (const Mesh::Vertex& v : mesh.vertices()) {
    ostream << "v " << FullPrecision(v[0]) << " " << FullPrecision(v[1]) << " "
            << FullPrecision(v[2]) << std::endl;
  }
  // Offset +1 because .obj format starts counting at 1.
  for (const Mesh::Face& f : mesh.faces()) {
    ostream << "f " << f[0] + 1 << " " << f[1] + 1 << " " << f[2] + 1
            << std::endl;
  }
  ostream.close();
  return absl::OkStatus();
}

}  // namespace intrinsic::geo
