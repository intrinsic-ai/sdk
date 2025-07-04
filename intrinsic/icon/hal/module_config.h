// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_HAL_MODULE_CONFIG_H_
#define INTRINSIC_ICON_HAL_MODULE_CONFIG_H_

#include <string>

#include "absl/container/flat_hash_set.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "intrinsic/icon/control/realtime_clock_interface.h"
#include "intrinsic/icon/hal/proto/hardware_module_config.pb.h"
#include "intrinsic/icon/hardware_modules/sim_bus/sim_bus_hardware_module.pb.h"
#include "intrinsic/icon/utils/realtime_guard.h"
#include "intrinsic/util/proto/any.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::icon {
namespace internal {

// Registers a string in the global set of proto types. This is exposed for
// templating only.
bool RegisterProtoTypes(absl::string_view type);

// A placeholder object for static registration of a proto name in the
// registered config proto set. Evaluating kRegisterProtoConfig<T> in a context
// that is not compiled out registers T with in the set.
template <typename T>
const bool kRegisterProtoConfig =
    internal::RegisterProtoTypes(T::descriptor()->name());

}  // namespace internal

// Returns a set of all of the type names of protos that GetConfig<T> has been
// used with. This is primarily used for static analysis of config types.
absl::flat_hash_set<std::string>& GetRegisteredConfigProtoTypes();

// A context object representing the state that a Hardware Module is initialized
// with.
class ModuleConfig {
 public:
  ModuleConfig() = default;

  // shared_memory_namespace is a prefix for all shared memory to make tests
  // hermetic. realtime_clock will be non-null if this hardware module drives
  // the clock. In that case, realtime_clock must outlive this ModuleConfig
  // object and any HardwareModule objects that are passed this ModuleConfig
  // during Init. Thread options can be provided in case the hardware module
  // needs to create its own threads using the same settings as ICON.
  explicit ModuleConfig(
      const intrinsic_proto::icon::HardwareModuleConfig& config,
      absl::string_view shared_memory_namespace,
      RealtimeClockInterface* realtime_clock,
      const ThreadOptions& icon_thread_options = ThreadOptions());

  // Returns the module config, typed to T. If the supplied
  // HardwareModuleConfig's module_config cannot be parsed as T, an error is
  // returned.
  template <typename T>
  absl::StatusOr<T> GetConfig() const;

  // Returns the simulation module config.
  intrinsic_proto::icon::SimBusModuleConfig GetSimulationConfig() const;

  // Returns the name of the module.
  // This will be the same as the context name, unless the module config has a
  // name set.
  const std::string& GetName() const;
  // Returns the name of the context in which the module is running.
  const std::string& GetContextName() const;

  // Returns thread options used by ICON.
  ThreadOptions GetIconThreadOptions() const;

  // Obtains the realtime clock, which can be used to tick the control layer.
  // Returns nullptr if this hardware module is not configured to drive to the
  // clock.
  RealtimeClockInterface* GetRealtimeClock() const;

  absl::string_view GetSimulationServerAddress() const;

  absl::string_view GetSharedMemoryNamespace() const;

  bool DisableMallocGuard() const { return config_.disable_malloc_guard(); }

  // Returns the control period as defined in the underlying
  // HardwareModuleConfig.
  //
  // Returns NotFoundError if the config does not have a control period *or*
  // control frequency value.
  // Returns InvalidArgumentError if the control period / frequency is invalid
  // (less than or equal to zero, or infinite).
  absl::StatusOr<absl::Duration> GetControlPeriod() const;

 private:
  intrinsic_proto::icon::HardwareModuleConfig config_;
  std::string shared_memory_namespace_;

  // Raw pointer to externally-owned realtime clock. See constructor's comments
  // above for lifetime requirements.
  RealtimeClockInterface* realtime_clock_;

  // Thread options used by ICON.
  ThreadOptions icon_thread_options_;
};

template <typename T>
absl::StatusOr<T> ModuleConfig::GetConfig() const {
  INTRINSIC_ASSERT_NON_REALTIME();

  // Register the proto for use in static analysis. The purpose of the call
  // below is that kRegisterProtoConfig is statically evaluated, which adds T to
  // a global set of registered configs.
  (void)internal::kRegisterProtoConfig<T>;

  return UnpackAny<T>(config_.module_config());
}

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_HAL_MODULE_CONFIG_H_
