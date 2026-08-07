// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_CONVERSION_SDF_PHYSICS_COMPONENT_FROM_SDF_H_
#define INTRINSIC_WORLD_CONVERSION_SDF_PHYSICS_COMPONENT_FROM_SDF_H_

#include <memory>

#include "absl/status/statusor.h"
#include "intrinsic/world/component/physics_component.h"
#include "intrinsic/world/proto/physics_component.pb.h"
#include "sdf/Link.hh"

namespace intrinsic {
namespace sdf {

absl::StatusOr<std::unique_ptr<PhysicsComponent>> PhysicsComponentFromSdfLink(
    const ::sdf::Link& sdf_link);

}  // namespace sdf
}  // namespace intrinsic
#endif  // INTRINSIC_WORLD_CONVERSION_SDF_PHYSICS_COMPONENT_FROM_SDF_H_
