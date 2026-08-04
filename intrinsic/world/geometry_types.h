// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_GEOMETRY_TYPES_H_
#define INTRINSIC_WORLD_GEOMETRY_TYPES_H_

namespace intrinsic {

// If specified this set of geometry will be used to do generic collision
// checking against other objects.
constexpr char kKindCollisionGeometry[] = "Intrinsic_Collision";

// If specified this set of geometry will be used for visualization purposes
// instead of the default canonical set.
constexpr char kKindVisualGeometry[] = "Intrinsic_Visual";

}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_GEOMETRY_TYPES_H_
