// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PERCEPTION_PUBLIC_CALIBRATION_SERVICE_UTILS_H_
#define INTRINSIC_PERCEPTION_PUBLIC_CALIBRATION_SERVICE_UTILS_H_

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

#endif  // INTRINSIC_PERCEPTION_PUBLIC_CALIBRATION_SERVICE_UTILS_H_
