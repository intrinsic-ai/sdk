// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_HARDWARE_MODULE_INIT_CONTEXT_H_
#define INTRINSIC_ICON_HAL_HARDWARE_MODULE_INIT_CONTEXT_H_

#include <string>

#include "absl/base/attributes.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/time/time.h"
#include "grpcpp/server_builder.h"
#include "intrinsic/icon/hal/hardware_interface_registry.h"
#include "intrinsic/icon/hal/module_config.h"

namespace intrinsic::icon {

// Maximum frequency of inspection data publishing is 5 Hz.
const absl::Duration kMinInspectionDataPublishPeriod = absl::Seconds(1.0 / 5.0);
const absl::Duration kMaxInspectionDataPublishPeriod = absl::Seconds(3);

// Provides configuration and allows the hardware module to set some
// configuration on the hardware module runtime or register services. This class
// is passed during initialization of a Hardware Module to the hardware module
// Init() function and provides functionality such as
// - access to the module configuration
// - access to the interface registry
// - the ability to register a gRPC service
class HardwareModuleInitContext {
 public:
  HardwareModuleInitContext(HardwareInterfaceRegistry& interface_registry
                                ABSL_ATTRIBUTE_LIFETIME_BOUND,
                            grpc::ServerBuilder& server_builder
                                ABSL_ATTRIBUTE_LIFETIME_BOUND,
                            const ModuleConfig& config)
      : interface_registry_(interface_registry),
        server_builder_(server_builder),
        module_config_(config) {}
  ~HardwareModuleInitContext() = default;
  // Delete copy and move constructors since this class contains temporary
  // objects which are deleted after the hardware module is initialized.
  // Hardware modules should not be able to copy or move this class.
  HardwareModuleInitContext(const HardwareModuleInitContext&) = delete;
  HardwareModuleInitContext& operator=(const HardwareModuleInitContext&) =
      delete;
  HardwareModuleInitContext& operator=(HardwareModuleInitContext&&) = delete;

  // Returns the interface registry for this Hardware Module to register
  // interfaces.
  HardwareInterfaceRegistry& GetInterfaceRegistry() const {
    return interface_registry_;
  }

  // Returns the config for this Hardware Module.
  const ModuleConfig& GetModuleConfig() const { return module_config_; }

  // Registers a gRPC service with the hardware module runtime. The runtime
  // makes this service available to external components some time after the
  // hardware module's `Init()` function returns.
  //
  // Attention: `service` must live until the `Shutdown()` of the Hardware
  // Module is called!
  //
  // The gRPC service will still be served even if
  // HardwareModuleInterface::Init() returns an error.
  //
  // The gRPC service will run on a port that is reachable from external
  // components such as the frontend.
  void RegisterGrpcService(grpc::Service& service) {
    server_builder_.RegisterService(&service);
  }

  // Enables gathering cycle time metrics.
  // Call this during Init() of a hardware module, when `cycle_duration` is
  // known, gather cycle time metrics when the robot is enabled.
  // Logs warnings when the cycle time is exceeded, or a single operation took
  // too long when `log_cycle_time_warnings` is true.
  void EnableCycleTimeMetrics(absl::Duration cycle_duration,
                              bool log_cycle_time_warnings) {
    cycle_duration_for_cycle_time_metrics_ = cycle_duration;
    log_cycle_time_warnings_ = log_cycle_time_warnings;
  }

  // Returns true if cycle time warnings should be logged.
  bool AreCycleTimeWarningsEnabled() const { return log_cycle_time_warnings_; }
  // Returns the cycle duration, or ZeroDuration if not set.
  absl::Duration GetCycleDurationForCycleTimeMetrics() const {
    return cycle_duration_for_cycle_time_metrics_;
  }

  // Returns the asset instance name/service instance name. If it is not set
  // (as it is the case for unit tests), it returns the module name, which can
  // be configured in the hardware module config.
  std::string GetAssetInstanceName() const {
    return !module_config_.GetContextName().empty()
               ? module_config_.GetContextName()
               : module_config_.GetName();
  }

  // Sets the interval at which inspection data is published.
  // The interval must be within the range [kMinInspectionDataPublishPeriod,
  // kMaxInspectionDataPublishPeriod].
  absl::Status SetInspectionDataPublishInterval(absl::Duration interval) {
    if (interval < kMinInspectionDataPublishPeriod) {
      return absl::InvalidArgumentError(
          absl::StrCat("Inspection data publish interval must be at least ",
                       absl::FormatDuration(kMinInspectionDataPublishPeriod)));
    }
    if (interval > kMaxInspectionDataPublishPeriod) {
      return absl::InvalidArgumentError(
          absl::StrCat("Inspection data publish interval must be at most ",
                       absl::FormatDuration(kMaxInspectionDataPublishPeriod)));
    }
    inspection_data_publishing_period_ = interval;
    return absl::OkStatus();
  }
  absl::Duration GetInspectionDataPublishPeriod() const {
    return inspection_data_publishing_period_;
  }

 private:
  HardwareInterfaceRegistry& interface_registry_;
  grpc::ServerBuilder& server_builder_;
  const ModuleConfig module_config_;
  absl::Duration cycle_duration_for_cycle_time_metrics_ = absl::ZeroDuration();
  bool log_cycle_time_warnings_ = false;
  // HardwareModuleInterface::ProvideInspectionData() is called and at which the
  // inspection data is published.
  absl::Duration inspection_data_publishing_period_ = absl::Seconds(1.0 / 3.0);
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_HARDWARE_MODULE_INIT_CONTEXT_H_
