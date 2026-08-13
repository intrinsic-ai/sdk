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

#ifndef INTRINSIC_GEOMETRY_API_MATERIAL_H_
#define INTRINSIC_GEOMETRY_API_MATERIAL_H_

#include <array>

namespace intrinsic {

// Material parameters similar to those found in fixed-pipeline OpenGL.
// https://www.khronos.org/registry/OpenGL-Refpages/gl2.1/xhtml/glMaterial.xml
struct Material {
  std::array<float, 4> ambient = {0.2f, 0.2f, 0.2f, 1};
  std::array<float, 4> diffuse = {0.8f, 0.8f, 0.8f, 1};
  std::array<float, 4> specular = {0, 0, 0, 1};
  std::array<float, 4> emission = {0, 0, 0, 1};
  float shininess = 0;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_MATERIAL_H_
