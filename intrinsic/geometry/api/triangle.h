// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_TRIANGLE_H_
#define INTRINSIC_GEOMETRY_API_TRIANGLE_H_

#include "intrinsic/eigenmath/types.h"

namespace intrinsic {

// Triangle types used with Mesh types.
struct Triangle {
  eigenmath::Vector3d v0;
  eigenmath::Vector3d v1;
  eigenmath::Vector3d v2;
};

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_TRIANGLE_H_
