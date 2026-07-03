// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/perception/calibration_service_utils.h"

#include <memory>
#include <string_view>
#include <utility>

#include "absl/status/statusor.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "grpcpp/client_context.h"
#include "intrinsic/connect/cc/grpc/channel.h"
#include "intrinsic/icon/release/grpc_time_support.h"
#include "intrinsic/perception/proto/v1/calibration_service.grpc.pb.h"
#include "intrinsic/skills/cc/equipment_pack.h"
#include "intrinsic/skills/cc/skill_utils.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_builder.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/time/deadline_timeout.h"

namespace intrinsic::perception {

namespace {
constexpr absl::Duration kMaxCreateRequestTime = absl::Seconds(600);
constexpr absl::Duration kMaxConnectTime = absl::Seconds(600);
constexpr std::string_view kCalibrationServiceEquipmentSlot =
    "calibration_service";
}  // namespace

absl::StatusOr<std::unique_ptr<
    intrinsic_proto::perception::v1::CalibrationService::StubInterface>>
CreateCalibrationServiceStub(const skills::EquipmentPack& equipment) {
  INTR_ASSIGN_OR_RETURN(const auto& calibration_service_handle,
                        equipment.GetHandle(kCalibrationServiceEquipmentSlot));
  INTR_ASSIGN_OR_RETURN(
      const intrinsic::ConnectionParams connection_config,
      skills::GetConnectionParamsFromHandle(calibration_service_handle));
  INTR_ASSIGN_OR_RETURN(
      auto channel, connect::CreateClientChannel(
                        connection_config.address, ToDeadline(kMaxConnectTime),
                        connect::UnlimitedMessageSizeGrpcChannelArgs(),
                        /*use_default_application_credentials=*/false,
                        connection_config.instance_name));
  auto stub =
      intrinsic_proto::perception::v1::CalibrationService::NewStub(channel);
  if (stub == nullptr) {
    return intrinsic::InternalErrorBuilder()
           << "Cannot connect to calibration server "
           << connection_config.address;
  }
  return std::move(stub);
}

absl::StatusOr<std::unique_ptr<grpc::ClientContext>>
CreateCalibrationServiceContext(const skills::EquipmentPack& equipment) {
  INTR_ASSIGN_OR_RETURN(const auto& calibration_service_handle,
                        equipment.GetHandle(kCalibrationServiceEquipmentSlot));
  INTR_ASSIGN_OR_RETURN(
      const intrinsic::ConnectionParams connection_config,
      skills::GetConnectionParamsFromHandle(calibration_service_handle));
  auto grpc_context = std::make_unique<grpc::ClientContext>();
  grpc_context->set_deadline(absl::Now() + kMaxCreateRequestTime);
  grpc_context->AddMetadata("x-resource-instance-name",
                            connection_config.instance_name);
  return std::move(grpc_context);
}

}  // namespace intrinsic::perception
