// Copyright 2023 Intrinsic Innovation LLC

// The `list_parts` binary is a tool that lists available robot parts from an
// Icon Application Layer Service.
//
//
#include <iostream>
#include <ostream>
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/cc_client/client.h"
#include "intrinsic/icon/release/portable/init_xfa.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_FLAG(std::string, server, "xfa.lan:17080",
          "Address of the ICON Application Layer Server");
ABSL_FLAG(std::string, instance, "robot_controller",
          "Name of the ICON service/resource instance.");

namespace {

absl::Status Run(const intrinsic::ConnectionParams& connection_params) {
  INTR_ASSIGN_OR_RETURN(auto icon_channel,
                        intrinsic::Channel::Make(connection_params));
  INTR_ASSIGN_OR_RETURN(auto parts,
                        intrinsic::icon::Client(icon_channel).ListParts());
  for (const auto& part_name : parts) {
    std::cout << part_name << std::endl;
  }

  return absl::OkStatus();
}

}  // namespace

int main(int argc, char** argv) {
  InitXfa(argv[0], argc, argv);

  QCHECK_OK(Run(intrinsic::ConnectionParams::ResourceInstance(
      absl::GetFlag(FLAGS_instance), absl::GetFlag(FLAGS_server))));

  return 0;
}
