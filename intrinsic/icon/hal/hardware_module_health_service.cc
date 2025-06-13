// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/hardware_module_health_service.h"

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/proto/v1/service_state.pb.h"
#include "intrinsic/icon/hal/hardware_module_util.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state.fbs.h"
#include "intrinsic/icon/hal/interfaces/hardware_module_state_utils.h"
#include "intrinsic/util/status/status_macros_grpc.h"

namespace intrinsic::icon {
namespace {

intrinsic_proto::services::v1::State::StateCode GetServiceStateCode(
    const intrinsic_fbs::HardwareModuleState& fb_state) {
  ::intrinsic_proto::services::v1::State::StateCode state =
      intrinsic_proto::services::v1::State::STATE_CODE_UNSPECIFIED;
  switch (fb_state.code()) {
    case intrinsic_fbs::StateCode::kDeactivated:
    case intrinsic_fbs::StateCode::kDeactivating:
    case intrinsic_fbs::StateCode::kActivated:
    case intrinsic_fbs::StateCode::kActivating:
    case intrinsic_fbs::StateCode::kMotionDisabling:
    case intrinsic_fbs::StateCode::kMotionEnabling:
      state = intrinsic_proto::services::v1::State::STATE_CODE_DISABLED;
      break;
    case intrinsic_fbs::StateCode::kMotionEnabled:
      state = intrinsic_proto::services::v1::State::STATE_CODE_ENABLED;
      break;
    case intrinsic_fbs::StateCode::kFaulted:
    case intrinsic_fbs::StateCode::kClearingFaults:
    case intrinsic_fbs::StateCode::kInitFailed:
    case intrinsic_fbs::StateCode::kFatallyFaulted:
      state = intrinsic_proto::services::v1::State::STATE_CODE_ERROR;
      break;
    default:
      state = intrinsic_proto::services::v1::State::STATE_CODE_UNSPECIFIED;
      break;
  }
  return state;
}
}  // namespace

HardwareModuleHealthService::~HardwareModuleHealthService() {
  // We must always set a value on  the promise, otherwise std::~promise will
  // throw, which is bad since we don't allow exceptions.
  absl::MutexLock lock(&mutex_);
  NotifyWithExitCode(HardwareModuleExitCode::kNormalShutdown);
};

void HardwareModuleHealthService::NotifyWithExitCode(
    HardwareModuleExitCode exit_code) {
  if (auto promise_wrapper = hardware_module_exit_code_promise_.lock();
      promise_wrapper != nullptr) {
    if (promise_wrapper->HasBeenSet()) {
      return;
    }
    if (auto status = promise_wrapper->SetValue(exit_code); !status.ok()) {
      LOG(ERROR) << "Failed to set exit code: " << status;
    }
  }
}

::grpc::Status HardwareModuleHealthService::GetState(
    grpc::ServerContext* context,
    const intrinsic_proto::services::v1::GetStateRequest* request,
    intrinsic_proto::services::v1::State* response) {
  absl::MutexLock lock(&mutex_);

  if (!latched_init_fault_.ok()) {
    response->set_state_code(
        intrinsic_proto::services::v1::State::STATE_CODE_ERROR);
    response->mutable_extended_status()->set_title(
        "Hardware module is in init failure.");
    response->mutable_extended_status()->mutable_user_report()->set_message(
        latched_init_fault_.message());

  } else if (hardware_module_runtime_ == nullptr) {
    response->set_state_code(
        intrinsic_proto::services::v1::State::STATE_CODE_ERROR);
    response->mutable_extended_status()->set_title(
        "Creation of hardware module failed.");
    response->mutable_extended_status()->mutable_user_report()->set_message(
        "Try restarting the hardware module.");

  } else {
    INTR_ASSIGN_OR_RETURN_GRPC(
        auto fb_state, hardware_module_runtime_->GetHardwareModuleState());
    intrinsic_proto::services::v1::State::StateCode state_code =
        GetServiceStateCode(fb_state);
    response->set_state_code(state_code);
    absl::string_view message = GetMessage(&fb_state);
    if (!message.empty()) {
      response->mutable_extended_status()->mutable_user_report()->set_message(
          message);
    }
  }
  return ::grpc::Status::OK;
}

::grpc::Status HardwareModuleHealthService::Enable(
    grpc::ServerContext* context,
    const intrinsic_proto::services::v1::EnableRequest* request,
    intrinsic_proto::services::v1::EnableResponse* response) {
  absl::MutexLock lock(&mutex_);

  if (!latched_init_fault_.ok()) {
    NotifyWithExitCode(HardwareModuleExitCode::kFatalFaultDuringInit);
    return ::grpc::Status::OK;
  }
  if (hardware_module_runtime_ == nullptr) {
    return ToGrpcStatus(absl::FailedPreconditionError(
        "Cannot use enable to clear faults since there is no hardware module "
        "running."));
  }
  INTR_ASSIGN_OR_RETURN_GRPC(
      const intrinsic_fbs::HardwareModuleState state,
      hardware_module_runtime_->GetHardwareModuleState());

  if (state.code() == intrinsic_fbs::StateCode::kInitFailed) {
    NotifyWithExitCode(HardwareModuleExitCode::kFatalFaultDuringInit);
    return ::grpc::Status::OK;
  } else if (state.code() == intrinsic_fbs::StateCode::kFatallyFaulted) {
    NotifyWithExitCode(HardwareModuleExitCode::kFatalFaultDuringExec);
    return ::grpc::Status::OK;
  } else if (state.code() == intrinsic_fbs::StateCode::kFaulted ||
             state.code() == intrinsic_fbs::StateCode::kClearingFaults) {
    return ToGrpcStatus(absl::UnavailableError(
        "Cannot use enable to clear runtime faults directly on hardware "
        "modules. Clear the error on the realtime control service or in the "
        "robot control panel."));
  } else if (state.code() == intrinsic_fbs::StateCode::kMotionEnabled ||
             state.code() == intrinsic_fbs::StateCode::kMotionEnabling) {
    return ::grpc::Status::OK;
  }

  return ToGrpcStatus(absl::UnavailableError(
      "Cannot enable hardware module directly. Hardware modules are enabled "
      "automatically via the realtime control service when no hardware "
      "module has an error."));
}

::grpc::Status HardwareModuleHealthService::Disable(
    grpc::ServerContext* context,
    const intrinsic_proto::services::v1::DisableRequest* request,
    intrinsic_proto::services::v1::DisableResponse* response) {
  return ToGrpcStatus(absl::UnavailableError(
      "Cannot disable hardware module directly. They are disabled "
      "automatically when an error is detected."));
}

}  // namespace intrinsic::icon
