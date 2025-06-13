// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/hardware_module_main_util.h"

#include <dirent.h>
#include <fcntl.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <unistd.h>

#include <future>  // NOLINT
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/base/nullability.h"
#include "absl/container/flat_hash_set.h"
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "grpc/grpc.h"
#include "grpcpp/security/server_credentials.h"
#include "grpcpp/server.h"
#include "grpcpp/server_builder.h"
#include "intrinsic/assets/services/proto/v1/service_state.grpc.pb.h"
#include "intrinsic/icon/hal/hardware_module_health_service.h"
#include "intrinsic/icon/hal/hardware_module_runtime.h"
#include "intrinsic/icon/hal/hardware_module_util.h"
#include "intrinsic/icon/hal/proto/hardware_module_config.pb.h"
#include "intrinsic/icon/hal/realtime_clock.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/icon/release/file_helpers.h"
#include "intrinsic/icon/utils/clock.h"
#include "intrinsic/icon/utils/duration.h"
#include "intrinsic/icon/utils/shutdown_signals.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"
#include "intrinsic/util/proto/any.h"
#include "intrinsic/util/proto/get_text_proto.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/thread_options.h"
#include "intrinsic/util/thread/util.h"

namespace intrinsic::icon {

namespace {

absl::StatusOr<HardwareModuleMainConfig> LoadHardwareModuleConfig(
    intrinsic_proto::config::RuntimeContext&& context) {
  INTR_ASSIGN_OR_RETURN(
      auto module_config,
      intrinsic::UnpackAny<intrinsic_proto::icon::HardwareModuleConfig>(
          context.config()),
      _ << "Unpacking module config");

  if (!module_config.name().empty()) {
    LOG(INFO) << "Explicit hardware module name '" << module_config.name()
              << "' specified. Consider removing the name field from the "
                 "hardware module config.";
  } else {
    module_config.set_name(context.name());
  }
  // Always set the context name, even if the module config has a name.
  module_config.set_context_name(context.name());

  module_config.set_simulation_server_address(
      context.simulation_server_address());

  // override the realtime flag based upon what mode we're running in.
  const bool use_realtime_scheduling = [&context]() {
    switch (context.level()) {
      case intrinsic_proto::config::RuntimeContext::REALITY:
        return true;
      case intrinsic_proto::config::RuntimeContext::PHYSICS_SIM:
        return false;
      case intrinsic_proto::config::RuntimeContext::UNSPECIFIED:
      default:
        LOG(WARNING) << "Received unexpected runtime context level of "
                     << context.level()
                     << ".  Running with realtime priority disabled.";
        return false;
    }
  }();
  return HardwareModuleMainConfig{std::move(context), module_config,
                                  use_realtime_scheduling};
}

}  // namespace

absl::StatusOr<HardwareModuleMainConfig> LoadConfig(
    absl::string_view module_config_file,
    absl::string_view runtime_context_file, bool use_realtime_scheduling) {
  if (module_config_file.empty() && runtime_context_file.empty()) {
    return absl::InvalidArgumentError(
        "Either runtime context file or module config file must be set");
  }
  if (!module_config_file.empty()) {
    LOG(INFO) << "Not running as a resource. Loading textproto from "
              << module_config_file;
    intrinsic_proto::icon::HardwareModuleConfig module_config;
    INTR_RETURN_IF_ERROR(
        intrinsic::GetTextProto(module_config_file, module_config));
    return HardwareModuleMainConfig{
        .module_config = module_config,
        .use_realtime_scheduling = use_realtime_scheduling};
  }

  // Shall never fail based on user configuration.
  INTR_ASSIGN_OR_RETURN(
      auto runtime_context,
      intrinsic::GetBinaryProto<intrinsic_proto::config::RuntimeContext>(
          runtime_context_file));

  LOG(INFO) << "Running as a resource. Loading runtime context from binary "
               "proto from "
            << runtime_context_file;
  // Could fail due to user configuration.
  return LoadHardwareModuleConfig(std::move(runtime_context));
}

absl::StatusOr<HardwareModuleRtSchedulingData> SetupRtScheduling(
    const intrinsic_proto::icon::HardwareModuleConfig& module_config,
    SharedMemoryManager& shm_manager, bool use_realtime_scheduling,
    std::optional<int> realtime_core, bool disable_malloc_guard) {
  std::unique_ptr<intrinsic::icon::RealtimeClock> realtime_clock = nullptr;
  if (module_config.drives_realtime_clock()) {
    INTR_ASSIGN_OR_RETURN(realtime_clock,
                          intrinsic::icon::RealtimeClock::Create(shm_manager));
  }

  absl::StatusOr<absl::flat_hash_set<int>> affinity_set =
      absl::FailedPreconditionError("Did not read Affinity set.");

  if (!module_config.realtime_cores().empty()) {
    LOG(INFO) << "Reading realtime core from proto config.";
    affinity_set =
        absl::flat_hash_set<int>{module_config.realtime_cores().begin(),
                                 module_config.realtime_cores().end()};
  } else if (realtime_core.has_value()) {
    LOG(INFO) << "Reading realtime core from flag.";
    affinity_set = absl::flat_hash_set<int>{*realtime_core};
  } else {
    LOG(INFO) << "Reading realtime core from /proc/cmdline";
    affinity_set = intrinsic::ReadCpuAffinitySetFromCommandLine();
  }

  intrinsic::ThreadOptions server_thread_options;
  if (use_realtime_scheduling) {
    LOG(INFO) << "Configuring hardware module with RT options.";
    // A realtime config without affinity set is not valid.
    INTR_RETURN_IF_ERROR(affinity_set.status());
    LOG(INFO) << "Realtime cores are: " << absl::StrJoin(*affinity_set, ", ");
    server_thread_options =
        intrinsic::ThreadOptions()
            .SetRealtimeHighPriorityAndScheduler()
            .SetAffinity({affinity_set->begin(), affinity_set->end()});
  }
  return HardwareModuleRtSchedulingData{
      std::move(realtime_clock), server_thread_options,
      affinity_set.value_or(absl::flat_hash_set<int>{})};
}

namespace {

// Simple convenience function to check if `exit_code_future` was set before
// reaching deadline.
// Returns true, if `exit_code_future` has a value set.
bool ExitCodeFutureHasValue(
    std::shared_future<HardwareModuleExitCode>& exit_code_future,
    Clock::time_point deadline) {
  if (!exit_code_future.valid()) {
    return false;
  }
  return (exit_code_future.wait_until(deadline) == std::future_status::ready);
}

std::optional<HardwareModuleExitCode> GetExitCodeFromFuture(
    std::shared_future<HardwareModuleExitCode>& exit_code_future,
    Clock::time_point deadline) {
  if (!ExitCodeFutureHasValue(exit_code_future, deadline)) {
    return std::nullopt;
  }
  return exit_code_future.get();
}

std::optional<HardwareModuleExitCode> GetExitCodeFromFuture(
    std::shared_future<HardwareModuleExitCode>& exit_code_future,
    Duration timeout = Milliseconds(0)) {
  return GetExitCodeFromFuture(exit_code_future, Clock::Now() + timeout);
}

}  // namespace

absl::StatusOr<std::optional<HardwareModuleExitCode>>
RunRuntimeWithGrpcServerAndWaitForShutdown(
    const absl::StatusOr<HardwareModuleMainConfig>& main_config,
    const std::shared_ptr<SharedPromiseWrapper<HardwareModuleExitCode>>&
        exit_code_promise,
    absl::StatusOr</*absl_nonnull*/
                   std::unique_ptr<intrinsic::icon::HardwareModuleRuntime>>&
        runtime,
    std::optional<int> cli_grpc_server_port,
    const std::vector<int>& cpu_affinity) {
  absl::Status hwm_run_error;
  grpc::ServerBuilder server_builder;
  server_builder.AddChannelArgument(GRPC_ARG_ALLOW_REUSEPORT, 0);
  server_builder.AddChannelArgument(GRPC_ARG_MAX_METADATA_SIZE,
                                    16 * 1024);  // Set to 16KB
  std::shared_future<HardwareModuleExitCode> exit_code_future =
      exit_code_promise->GetSharedFuture();

  if (runtime.ok()) {
    INTR_RETURN_IF_ERROR(main_config.status())
        << "Runtime OK but config not OK - this is a bug: "
        << main_config.status();
    LOG(INFO) << "PUBLIC: Starting hardware module "
              << main_config->module_config.name();
    auto status = runtime.value()->Run(
        server_builder, main_config->use_realtime_scheduling, cpu_affinity);
    if (!status.ok()) {
      LOG(ERROR) << "PUBLIC: Error running hardware module: "
                 << status.message();
      hwm_run_error = status;
    }
  }

  intrinsic::icon::HardwareModuleHealthService health_service(
      exit_code_promise);

  std::optional<int> grpc_server_port = std::nullopt;
  if (main_config.ok() && main_config->runtime_context.has_value()) {
    grpc_server_port = main_config->runtime_context->port();
    LOG(INFO) << "Health Service port: " << *grpc_server_port;
  } else if (cli_grpc_server_port.has_value()) {
    grpc_server_port = cli_grpc_server_port;
    LOG(WARNING) << "No runtime context provided. Using grpc port "
                 << *grpc_server_port << " from command line ";
  }

  // Check if the start up of the HWM failed. If not, expose error via
  // ServiceState service.
  if (main_config.ok()) {
    if (runtime.ok()) {
      health_service.SetHardwareModuleRuntime(runtime.value().get());
    }
    if (!runtime.ok() || !runtime.value()->IsStarted()) {
      auto status = !runtime.ok() ? runtime.status() : hwm_run_error;
      LOG(INFO) << "Starting lame duck mode due to init error: " << status;
      health_service.ActivateLameDuckMode(status);
    } else {
      LOG(INFO) << "Hardware Module Runtime started.";
    }
  } else {
    health_service.ActivateLameDuckMode(absl::FailedPreconditionError(
        absl::StrCat("Failed to load hardware module config: ",
                     main_config.status().message())));
  }

  // grpc server variable needs to live until shutdown, but must be destroyed
  // before calling HardwareModuleRuntime::Shutdown(). Create server after
  // setting init faults on HealthService, so that init faults are immediately
  // available.
  std::unique_ptr<::grpc::Server> grpc_server;
  if (grpc_server_port.has_value()) {
    std::string address = absl::StrCat("[::]:", *grpc_server_port);
    server_builder.AddListeningPort(
        address,
        ::grpc::InsecureServerCredentials()  // NOLINT (insecure)
    );

    server_builder.RegisterService(
        static_cast<intrinsic_proto::services::v1::ServiceState::Service*>(
            &health_service));

    grpc_server = server_builder.BuildAndStart();
    LOG(INFO) << "gRPC server started on port " << *grpc_server_port;
  } else {
    LOG(WARNING) << "No gRPC port provided. Will not start gRPC server.";
  }

  LOG(INFO) << "Running until receiving shutdown signal.";
  // The poll loop here is necessary because we can't call runtime->Stop()
  // from the signal handler and a future isn't signal-safe
  // either.
  const Duration kPollShutdownSignalEvery = intrinsic::Seconds(0.2);
  std::optional<HardwareModuleExitCode> exit_code;
  while (IsShutdownRequested() == ShutdownType::kNotRequested) {
    auto next_check_deadline = Clock::Now() + kPollShutdownSignalEvery;
    if (ExitCodeFutureHasValue(exit_code_future, next_check_deadline)) {
      exit_code = GetExitCodeFromFuture(exit_code_future);
      break;
    }
  }
  LOG(INFO) << "Shutdown signal received";

  return exit_code;
}
}  // namespace intrinsic::icon
