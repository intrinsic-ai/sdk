// Copyright 2023 Intrinsic Innovation LLC

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
