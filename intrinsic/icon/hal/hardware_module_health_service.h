// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_HARDWARE_MODULE_HEALTH_SERVICE_H_
#define INTRINSIC_ICON_HAL_HARDWARE_MODULE_HEALTH_SERVICE_H_

#include <memory>
#include <utility>

#include "absl/base/thread_annotations.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/proto/v1/service_state.grpc.pb.h"
#include "intrinsic/assets/services/proto/v1/service_state.pb.h"
#include "intrinsic/icon/hal/hardware_module_runtime.h"
#include "intrinsic/icon/hal/hardware_module_util.h"

namespace intrinsic::icon {

// Implementation of the ServiceState service for all hardware module
// instances. It queries the current hardware module state from the given
// hardware module. Calls to ClearFaults() while in fatal and init failures
// trigger a signal to restart the process. Other actions such as `Enable()` are
// prohibited.
class HardwareModuleHealthService
    : public intrinsic_proto::services::v1::ServiceState::Service {
 public:
  explicit HardwareModuleHealthService(
      std::weak_ptr<SharedPromiseWrapper<HardwareModuleExitCode>>
          exit_code_promise)
      : hardware_module_exit_code_promise_(std::move(exit_code_promise)) {}

  ~HardwareModuleHealthService() override;

  // ServiceState implementation.
  ::grpc::Status GetState(
      grpc::ServerContext* context,
      const intrinsic_proto::services::v1::GetStateRequest* request,
      intrinsic_proto::services::v1::SelfState* response) override;

  ::grpc::Status Enable(
      grpc::ServerContext* context,
      const intrinsic_proto::services::v1::EnableRequest* request,
      intrinsic_proto::services::v1::EnableResponse* response) override;

  ::grpc::Status Disable(
      grpc::ServerContext* context,
      const intrinsic_proto::services::v1::DisableRequest* request,
      intrinsic_proto::services::v1::DisableResponse* response) override;

  // Sets the hardware module runtime. `hardware_module_runtime` must outlive
  // this class instance.
  void SetHardwareModuleRuntime(
      HardwareModuleRuntime* hardware_module_runtime) {
    absl::MutexLock lock(&mutex_);
    hardware_module_runtime_ = hardware_module_runtime;
  }

  // When activated, `CheckHealth()` reports the given `latched_init_fault`
  // until shutdown.
  void ActivateLameDuckMode(absl::Status latched_init_fault) {
    absl::MutexLock lock(&mutex_);
    latched_init_fault_ = latched_init_fault;
  }

 private:
  void NotifyWithExitCode(HardwareModuleExitCode exit_code)
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  absl::Mutex mutex_;
  absl::Status latched_init_fault_ ABSL_GUARDED_BY(mutex_);
  HardwareModuleRuntime* hardware_module_runtime_ ABSL_GUARDED_BY(mutex_) =
      nullptr;
  std::weak_ptr<SharedPromiseWrapper<HardwareModuleExitCode>>
      hardware_module_exit_code_promise_ ABSL_GUARDED_BY(mutex_);
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_HARDWARE_MODULE_HEALTH_SERVICE_H_
