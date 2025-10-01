// Copyright 2023 Intrinsic Innovation LLC

#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/examples/joint_then_cart_move_lib.h"
#include "intrinsic/icon/release/portable/init_intrinsic.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_FLAG(std::string, server, "xfa.lan:17080",
          "Address of the ICON Application Layer Server");
ABSL_FLAG(std::string, instance, "robot_controller",
          "Name of the ICON service/resource instance.");
ABSL_FLAG(std::string, part, "arm", "Part to control.");

const char kUsage[] =
    "Initially moves all joints into to a fixed position near the center of "
    "the joint position ranges. Then, performs a small Cartesian move in "
    "positive x direction.";

namespace {
absl::Status Run(const intrinsic::ConnectionParams& connection_params,
                 absl::string_view part_name) {
  if (connection_params.address.empty()) {
    return absl::FailedPreconditionError("`--server` must not be empty.");
  }
  if (part_name.empty()) {
    return absl::FailedPreconditionError("`--part` must not be empty.");
  }

  INTR_ASSIGN_OR_RETURN(auto icon_channel,
                        intrinsic::Channel::MakeFromAddress(connection_params));

  return intrinsic::icon::examples::JointThenCartMove(part_name, icon_channel);
}

}  // namespace

int main(int argc, char** argv) {
  InitXfa(kUsage, argc, argv);
  QCHECK_OK(Run(intrinsic::ConnectionParams::ResourceInstance(
                    absl::GetFlag(FLAGS_instance), absl::GetFlag(FLAGS_server)),
                absl::GetFlag(FLAGS_part)));
  return 0;
}
