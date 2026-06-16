// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/cc/skill_utils.h"

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "intrinsic/resources/proto/resource_handle.pb.h"
#include "intrinsic/skills/proto/equipment.pb.h"
#include "intrinsic/skills/proto/skills.pb.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::skills {

absl::StatusOr<intrinsic::ConnectionParams> GetConnectionParamsFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle) {
  if (!handle.connection_info().has_grpc()) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "Resource handle \"%s\" does not specify grpc connection_info",
        handle.name()));
  }
  return intrinsic::ConnectionParams{
      .address =
          std::string(handle.connection_info().grpc().address()),  // NOLINT
      .instance_name = std::string(                                // NOLINT
          handle.connection_info().grpc().server_instance()),      // NOLINT
      .header =
          std::string(handle.connection_info().grpc().header()),  // NOLINT
  };
}

absl::StatusOr<std::shared_ptr<intrinsic::Channel>> CreateChannelFromHandle(
    const intrinsic_proto::resources::ResourceHandle& handle) {
  INTR_ASSIGN_OR_RETURN(const intrinsic::ConnectionParams connection_params,
                        GetConnectionParamsFromHandle(handle));

  return intrinsic::Channel::MakeFromAddress(connection_params);
}

intrinsic_proto::skills::Footprint CreateObjectReservationFootprint(
    absl::string_view object_name,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type) {
  intrinsic_proto::skills::Footprint footprint;
  AddObjectReservation(object_name, type, footprint);
  return footprint;
}

intrinsic_proto::skills::Footprint CreateObjectReservationFootprint(
    const intrinsic_proto::world::ObjectReferenceByName& object,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type) {
  intrinsic_proto::skills::Footprint footprint;
  auto* reservation = footprint.add_object_reservation();
  reservation->set_type(type);
  *reservation->mutable_object() = object;
  return footprint;
}

intrinsic_proto::skills::Footprint CreateUniverseLockFootprint() {
  intrinsic_proto::skills::Footprint footprint;
  footprint.set_lock_the_universe(true);
  return footprint;
}

void AddObjectReservation(
    absl::string_view object_name,
    intrinsic_proto::skills::ObjectWorldReservation::SharingType type,
    intrinsic_proto::skills::Footprint& footprint) {
  auto* reservation = footprint.add_object_reservation();
  reservation->set_type(type);
  reservation->mutable_object()->set_object_name(std::string(object_name));
}

void AddResourceReservation(
    absl::string_view resource_name,
    intrinsic_proto::skills::ResourceReservation::SharingType type,
    intrinsic_proto::skills::Footprint& footprint) {
  auto* reservation = footprint.add_resource_reservation();
  reservation->set_name(std::string(resource_name));
  reservation->set_type(type);
}

}  // namespace intrinsic::skills
