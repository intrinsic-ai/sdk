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
