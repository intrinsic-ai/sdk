// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/module_config.h"

#include <string>

#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/control/realtime_clock_interface.h"
#include "intrinsic/icon/hal/proto/hardware_module_config.pb.h"
#include "intrinsic/icon/hardware_modules/sim_bus/sim_bus_hardware_module.pb.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::icon {
namespace internal {

bool RegisterProtoTypes(absl::string_view type) {
  GetRegisteredConfigProtoTypes().emplace(type);
  return true;
}

}  // namespace internal

ModuleConfig::ModuleConfig(
    const intrinsic_proto::icon::HardwareModuleConfig& config,
    absl::string_view shared_memory_namespace,
    RealtimeClockInterface* realtime_clock,
    const ThreadOptions& icon_thread_options)
    : config_(config),
      shared_memory_namespace_(shared_memory_namespace),
      realtime_clock_(realtime_clock),
      icon_thread_options_(icon_thread_options) {}

intrinsic_proto::icon::SimBusModuleConfig ModuleConfig::GetSimulationConfig()
    const {
  return config_.simulation_module_config();
}

const std::string& ModuleConfig::GetName() const { return config_.name(); }

ThreadOptions ModuleConfig::GetIconThreadOptions() const {
  return icon_thread_options_;
}

RealtimeClockInterface* ModuleConfig::GetRealtimeClock() const {
  return realtime_clock_;
}

absl::flat_hash_set<std::string>& GetRegisteredConfigProtoTypes() {
  static auto* proto_types = new absl::flat_hash_set<std::string>();
  return *proto_types;
}

absl::string_view ModuleConfig::GetSimulationServerAddress() const {
  return config_.simulation_server_address();
}
absl::string_view ModuleConfig::GetSharedMemoryNamespace() const {
  return shared_memory_namespace_;
}

absl::StatusOr<absl::Duration> ModuleConfig::GetControlPeriod() const {
  switch (config_.control_rate_case()) {
    case intrinsic_proto::icon::HardwareModuleConfig::CONTROL_RATE_NOT_SET:
      break;
    case intrinsic_proto::icon::HardwareModuleConfig::kControlFrequencyHz: {
      if (config_.control_frequency_hz() <= 0) {
        return absl::InvalidArgumentError(
            absl::StrCat("Control frequency must be positive, but is ",
                         config_.control_frequency_hz()));
      }
      return absl::Seconds(1.0 / config_.control_frequency_hz());
    }
    case intrinsic_proto::icon::HardwareModuleConfig::kControlPeriodNs: {
      if (config_.control_period_ns() <= 0) {
        return absl::InvalidArgumentError(
            absl::StrCat("Control period must be positive, but is ",
                         config_.control_period_ns()));
      }
      return absl::Nanoseconds(config_.control_period_ns());
    }
  }
  return absl::NotFoundError(
      "No control frequency or period in HardwareModuleConfig");
}

}  // namespace intrinsic::icon
