// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SKILLS_CC_SKILL_UTILS_H_
#define INTRINSIC_SKILLS_CC_SKILL_UTILS_H_

#include <memory>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/resources/proto/resource_handle.pb.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/footprint.pb.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"

namespace intrinsic {
namespace skills {

absl::StatusOr<ConnectionParams> GetConnectionParamsFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle);

// Creates client channel for communicating with equipment.
absl::StatusOr<std::shared_ptr<intrinsic::Channel>> CreateChannelFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle);

// Creates a footprint with a single object reservation.
intrinsic_proto::skills::Footprint CreateObjectReservationFootprint(
    absl::string_view object_name,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type);

// Creates a footprint with a single object reservation using a pre-constructed
// ObjectReferenceByName.
intrinsic_proto::skills::Footprint CreateObjectReservationFootprint(
    const intrinsic_proto::world::ObjectReferenceByName& object,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type);

// Creates a footprint with a universe lock.
intrinsic_proto::skills::Footprint CreateUniverseLockFootprint();

// Adds an object reservation to an existing footprint.
void AddObjectReservation(
    absl::string_view object_name,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type,
    intrinsic_proto::skills::Footprint& footprint);

// Adds a resource reservation to an existing footprint.
void AddResourceReservation(
    absl::string_view resource_name,
    intrinsic_proto::skills::ResourceReservation::SharingType type,
    intrinsic_proto::skills::Footprint& footprint);

}  // namespace skills
}  // namespace intrinsic

#endif  // INTRINSIC_SKILLS_CC_SKILL_UTILS_H_
