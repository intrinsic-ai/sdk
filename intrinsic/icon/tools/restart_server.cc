// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include <iostream>
#include <memory>
#include <ostream>
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/cc_client/client.h"
#include "intrinsic/icon/release/portable/init_intrinsic.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_macros.h"

ABSL_FLAG(std::string, server, "xfa.lan:17080", "Address of the ICON Server");
ABSL_FLAG(std::string, instance, "robot_controller",
          "Name of the ICON service/resource instance.");

const char* UsageString() {
  return R"(
Usage: restart_server [--server=<addr>]

Attempts to restart the entire ICON server.
The server waits until all sessions are ended, stops all hardware devices and
shuts down, then it starts again after a delay.
The ICON server receives the same signal as when the application is restarted.

This tool is not meant for typical operation, but only to apply exceptional
config changes in developer experiments (i.e. restarting ICON after adding
a plugin).

)";
}

namespace {

absl::Status Run(const intrinsic::ConnectionParams& connection_params) {
  INTR_ASSIGN_OR_RETURN(auto icon_channel,
                        intrinsic::Channel::MakeFromAddress(connection_params));
  return intrinsic::icon::Client(icon_channel).RestartServer();
}

}  // namespace

int main(int argc, char** argv) {
  InitIntrinsic(UsageString(), argc, argv);
  QCHECK_OK(Run(intrinsic::ConnectionParams::ResourceInstance(
      absl::GetFlag(FLAGS_instance), absl::GetFlag(FLAGS_server))));
  std::cout << "Requested ICON server restart." << std::endl;
  return 0;
}
