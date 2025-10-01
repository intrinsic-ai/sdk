// Copyright 2023 Intrinsic Innovation LLC

#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/examples/adio_lib.h"
#include "intrinsic/icon/release/portable/init_intrinsic.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_FLAG(std::string, server, "xfa.lan:17080",
          "Address of the ICON Application Layer Server");
ABSL_FLAG(std::string, instance, "robot_controller",
          "Name of the ICON service/resource instance.");
ABSL_FLAG(std::string, part, "adio", "Part to control.");

ABSL_FLAG(std::string, output_block, "outputs",
          "Name of the output_block to set bits.");

const char kUsage[] =
    "Sequentially sets all bits of 'output_block' to '1' and then clears them "
    "again. Only sets the two lowest bits if unable to determine the size of "
    "the output block.";

namespace {
absl::Status Run(const intrinsic::ConnectionParams& connection_params,
                 absl::string_view part_name,
                 absl::string_view output_block_name) {
  if (connection_params.address.empty()) {
    return absl::FailedPreconditionError("`--server` must not be empty.");
  }
  if (part_name.empty()) {
    return absl::FailedPreconditionError("`--part` must not be empty.");
  }

  INTR_ASSIGN_OR_RETURN(auto icon_channel,
                        intrinsic::Channel::MakeFromAddress(connection_params));

  return intrinsic::icon::examples::ExampleSetDigitalOutput(
      part_name, output_block_name, icon_channel);
}
}  // namespace

int main(int argc, char** argv) {
  InitXfa(kUsage, argc, argv);
  QCHECK_OK(Run(intrinsic::ConnectionParams::ResourceInstance(
                    absl::GetFlag(FLAGS_instance), absl::GetFlag(FLAGS_server)),
                absl::GetFlag(FLAGS_part), absl::GetFlag(FLAGS_output_block)));
  return 0;
}
