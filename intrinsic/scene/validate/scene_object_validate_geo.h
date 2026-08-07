// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATE_GEO_H_
#define INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATE_GEO_H_

#include "absl/status/status.h"
#include "intrinsic/geometry/storage/geometry_deserializer.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"

namespace intrinsic {

namespace scene_object {

// Returns OK if referenced geometries in the scene object can be loaded.
absl::Status ValidateReferencedGeos(
    const intrinsic_proto::scene_object::v1::SceneObject& object,
    const GeometryDeserializer& geolib);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_VALIDATE_SCENE_OBJECT_VALIDATE_GEO_H_
