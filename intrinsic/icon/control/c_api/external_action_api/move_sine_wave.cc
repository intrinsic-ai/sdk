// Copyright 2023 Intrinsic Innovation LLC

#include <memory>
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/cc_client/condition.h"
#include "intrinsic/icon/cc_client/session.h"
#include "intrinsic/icon/common/id_types.h"
#include "intrinsic/icon/control/c_api/external_action_api/sine_wave_action.pb.h"
#include "intrinsic/icon/control/c_api/external_action_api/sine_wave_plugin_action.h"
#include "intrinsic/icon/release/portable/init_xfa.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_FLAG(std::string, server, "xfa.lan:17080",
          "Address of the ICON Application Layer Server");
ABSL_FLAG(std::string, instance, "robot_controller",
          "Name of the ICON service/resource instance.");
ABSL_FLAG(std::string, part, "arm", "Part to control.");

const char kUsage[] =
    "Performs a sine wave motion. This motion is only available when the sine "
    "wave plugin has been loaded.";

namespace intrinsic::icon {
namespace {

constexpr int kDof = 6;
constexpr double kCycleDuration = 4.0f;

absl::Status Main(const ConnectionParams& connection_params,
                  absl::string_view part_name) {
  if (connection_params.address.empty()) {
    return absl::FailedPreconditionError("`--server` must not be empty.");
  }
  if (part_name.empty()) {
    return absl::FailedPreconditionError("`--part` must not be empty.");
  }

  INTR_ASSIGN_OR_RETURN(auto icon_channel, Channel::Make(connection_params));
  INTR_ASSIGN_OR_RETURN(std::unique_ptr<Session> session,
                        Session::Start(icon_channel, {std::string(part_name)}));
  LOG(INFO) << "Created session";
  SineWavePluginAction::ParameterProto params;
  for (int i = 0; i < kDof; ++i) {
    auto* joint_params = params.add_joints();
    joint_params->set_amplitude_rad(static_cast<double>(i) * 0.1f);
    joint_params->set_frequency_hz(1.0f / kCycleDuration);
  }
  ActionDescriptor sine_move = ActionDescriptor(SineWavePluginAction::kName,
                                                ActionInstanceId(1), part_name)
                                   .WithFixedParams(params);
  ReactionHandle timed_out;
  sine_move.WithReaction(
      ReactionDescriptor(IsGreaterThanOrEqual(
                             SineWavePluginAction::kStateVariableTimeSinceStart,
                             2 * kCycleDuration))
          .WithHandle(timed_out));
  LOG(INFO) << "AddAction. Parameters: " << params;
  INTR_ASSIGN_OR_RETURN(auto action, session->AddAction(sine_move));
  LOG(INFO) << "StartAction";
  INTR_RETURN_IF_ERROR(session->StartAction(action));
  LOG(INFO) << "RunWatcherLoop";
  auto status = session->RunWatcherLoopUntilReaction(timed_out);
  if (status.ok()) {
    return absl::OkStatus();
  }
  LOG(WARNING) << "Session ended early: " << status.message();
  return status;
}

}  // namespace
}  // namespace intrinsic::icon

int main(int argc, char** argv) {
  InitXfa(kUsage, argc, argv);
  QCHECK_OK(intrinsic::icon::Main(
      intrinsic::ConnectionParams::ResourceInstance(
          absl::GetFlag(FLAGS_instance), absl::GetFlag(FLAGS_server)),
      absl::GetFlag(FLAGS_part)));
  return 0;
}
