// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_GEOMETRY_API_APPLY_TRANSFORM_H_
#define INTRINSIC_GEOMETRY_API_APPLY_TRANSFORM_H_

#include "absl/status/statusor.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/geometry/api/geometry.h"
#include "intrinsic/math/pose3.h"

namespace intrinsic {

// These modifiers change the underlying representation.

// Transfers the representation of `geo` into reference space. The returned
// Geometry uses the following formula for all of the points of the mesh.
//
// new_point = ref_t_geo * old_point
absl::StatusOr<Geometry> ApplyTransform(const Geometry& geo,
                                        const Pose3d& ref_t_geo);

// Transfers the representation of `geo` into reference space. The returned
// Geometry uses the following formula for all of the points of the mesh.
//
// new_point = ref_t_geo * old_point
absl::StatusOr<Geometry> ApplyTransform(const Geometry& geo,
                                        const eigenmath::Matrix4d& ref_t_geo);

}  // namespace intrinsic

#endif  // INTRINSIC_GEOMETRY_API_APPLY_TRANSFORM_H_
