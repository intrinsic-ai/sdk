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

#ifndef INTRINSIC_PERCEPTION_CALIBRATION_SERVICE_UTILS_H_
#define INTRINSIC_PERCEPTION_CALIBRATION_SERVICE_UTILS_H_

#include <memory>

#include "absl/status/statusor.h"
#include "grpcpp/client_context.h"
#include "intrinsic/perception/proto/v1/calibration_service.grpc.pb.h"
#include "intrinsic/skills/cc/equipment_pack.h"

namespace intrinsic::perception {

// Creates a CalibrationService gRPC stub using the "calibration_service"
// equipment handle from the provided EquipmentPack.
absl::StatusOr<std::unique_ptr<
    intrinsic_proto::perception::v1::CalibrationService::StubInterface>>
CreateCalibrationServiceStub(const skills::EquipmentPack& equipment);

// Creates a grpc::ClientContext with deadline and instance metadata for
// communicating with the CalibrationService based on the provided
// EquipmentPack.
absl::StatusOr<std::unique_ptr<grpc::ClientContext>>
CreateCalibrationServiceContext(const skills::EquipmentPack& equipment);

}  // namespace intrinsic::perception

#endif  // INTRINSIC_PERCEPTION_CALIBRATION_SERVICE_UTILS_H_
