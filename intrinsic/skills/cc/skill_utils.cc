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

namespace {

constexpr char kAssetInstanceNameHeader[] = "x-resource-instance-name";

}  // namespace

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

absl::StatusOr<intrinsic::ConnectionParams>
GetConnectionParamsFromResolvedDependency(
    const intrinsic_proto::assets::v1::ResolvedDependency& dep,
    absl::string_view interface_uri) {
  const auto it = dep.interfaces().find(std::string(interface_uri));
  if (it == dep.interfaces().end()) {
    return absl::NotFoundError(absl::StrFormat(
        "Interface \"%s\" not found in ResolvedDependency \"%s\"",
        interface_uri, dep.name()));
  }
  const auto& iface = it->second;
  if (!iface.has_grpc() || !iface.grpc().has_connection()) {
    return absl::InvalidArgumentError(
        absl::StrFormat("Interface \"%s\" in ResolvedDependency \"%s\" does "
                        "not specify gRPC connection",
                        interface_uri, dep.name()));
  }

  const auto& grpc_conn = iface.grpc().connection();
  if (grpc_conn.address().empty()) {
    return absl::InvalidArgumentError(
        absl::StrFormat("gRPC connection for interface \"%s\" in "
                        "ResolvedDependency \"%s\" has empty address",
                        interface_uri, dep.name()));
  }

  if (grpc_conn.metadata_size() > 1) {
    return absl::InvalidArgumentError(absl::StrFormat(
        "gRPC connection for interface \"%s\" in ResolvedDependency \"%s\" has "
        "%d metadata entries; only at most one '%s' "
        "header is supported",
        interface_uri, dep.name(), grpc_conn.metadata_size(),
        kAssetInstanceNameHeader));
  }

  intrinsic::ConnectionParams params;
  params.address = grpc_conn.address();

  if (grpc_conn.metadata_size() == 1) {
    const auto& metadata = grpc_conn.metadata(0);
    if (metadata.key() != kAssetInstanceNameHeader) {
      return absl::InvalidArgumentError(absl::StrFormat(
          "gRPC connection for interface \"%s\" in ResolvedDependency \"%s\" "
          "specifies unsupported metadata key \"%s\"; only "
          "'%s' is supported",
          interface_uri, dep.name(), metadata.key(), kAssetInstanceNameHeader));
    }
    params.header = metadata.key();
    params.instance_name = metadata.value();
  }

  return params;
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
