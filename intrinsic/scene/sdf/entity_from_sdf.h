// Copyright 2023 Intrinsic Innovation LLC


#ifndef INTRINSIC_SCENE_SDF_ENTITY_FROM_SDF_H_
#define INTRINSIC_SCENE_SDF_ENTITY_FROM_SDF_H_

#include <optional>
#include <string>
#include <vector>

#include "absl/status/statusor.h"
#include "intrinsic/geometry/storage/geometry_serializer.h"
#include "intrinsic/scene/proto/v1/entity.pb.h"
#include "intrinsic/scene/proto/v1/scene_object.pb.h"
#include "intrinsic/scene/sdf/sdf_path_resolver.h"
#include "sdf/Frame.hh"
#include "sdf/Joint.hh"
#include "sdf/Link.hh"
#include "sdf/Model.hh"
#include "sdf/Sensor.hh"
#include "sdf/World.hh"

namespace intrinsic {
namespace scene_object {

// Converts a SDFormat Frame to an optional Intrinsic Scene Object Entity with
// entity_type frame.
absl::StatusOr<std::optional<intrinsic_proto::scene_object::v1::Entity>>
EntityFromSdfFrame(const ::sdf::Frame& sdf_frame);

// Converts a SDFormat Link to a collection of Intrinsic Scene Object Entities.
// Each <link> converts to one Entity with entity_type link and multiple
// Entities with entity_type sensor. A `geometry_serializer` is required to
// store the geometry in callers's preferred storage.
absl::StatusOr<std::vector<intrinsic_proto::scene_object::v1::Entity>>
EntitiesFromSdfLink(const ::sdf::Link& sdf_link,
                    const sdf::UriResolver& uri_resolver,
                    GeometrySerializer& geometry_serializer);

// Converts a SDFormat Joint to a collection of Intrinsic Scene Object Entities.
// Each <joint> converts to one Entity with entity_type joint and multiple
// Entities with entity_type sensor.
struct CreateJointEntitiesResult {
  intrinsic_proto::scene_object::v1::Entity joint_entity;
  std::vector<intrinsic_proto::scene_object::v1::Entity> sensor_entities;
  std::string joint_child_name;
};
absl::StatusOr<CreateJointEntitiesResult> EntitiesFromSdfJoint(
    const ::sdf::Joint& sdf_joint);

// Converts a SDFormat Sensor to an Intrinsic Scene Object Entity with
// entity_type sensor.
absl::StatusOr<intrinsic_proto::scene_object::v1::Entity> EntityFromSdfSensor(
    const ::sdf::Sensor& sdf_sensor);

}  // namespace scene_object
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_SDF_ENTITY_FROM_SDF_H_
