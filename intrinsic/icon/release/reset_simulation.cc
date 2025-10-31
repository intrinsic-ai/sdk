// Copyright 2023 Intrinsic Innovation LLC

#include <iostream>
#include <memory>
#include <ostream>
#include <string>

#include "absl/flags/flag.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "google/protobuf/empty.pb.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "intrinsic/icon/release/portable/init_intrinsic.h"
#include "intrinsic/simulation/service/proto/v1/simulation_service.grpc.pb.h"
#include "intrinsic/simulation/service/proto/v1/simulation_service.pb.h"
#include "intrinsic/util/grpc/grpc.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"

// Default to the ingress service's port.
ABSL_FLAG(std::string, address, "xfa.lan:17080",
          "Address of the Simulation Service");

const char* UsageString() {
  return R"(
Usage: reset_simulation --address=<addr>

Resets both the physics simulation and the ICON server.
Physics simulation is reset to a default position and any errors in the robot
control state are deleted.
This is useful for recovering from errors such as exceeding maximum limits or
hitting obstacles.
)";
}

namespace {

absl::Status ResetSimulation(absl::string_view address) {
  if (address.empty()) {
    return absl::FailedPreconditionError("You must provide --address=<addr>.");
  }

  INTR_ASSIGN_OR_RETURN(
      std::shared_ptr<grpc::Channel> channel,
      intrinsic::connect::CreateClientChannel(
          address,
          absl::Now() + intrinsic::connect::kGrpcClientConnectDefaultTimeout));
  std::unique_ptr<intrinsic_proto::simulation::v1::SimulationService::Stub>
      stub(
          intrinsic_proto::simulation::v1::SimulationService::NewStub(channel));
  if (stub == nullptr) {
    return absl::UnavailableError(
        "Could not create grpc stub to simulation server.");
  }
  grpc::ClientContext context;
  google::protobuf::Empty empty;
  intrinsic_proto::simulation::v1::ResetSimulationResponse response;
  std::cout << "Starting resetting simulation." << std::endl;
  INTR_RETURN_IF_ERROR(intrinsic::ToAbslStatus(stub->ResetSimulation(
      &context, intrinsic_proto::simulation::v1::ResetSimulationRequest(),
      &response)));
  std::cout << "Finished resetting simulation." << std::endl;
  return absl::OkStatus();
}

}  // namespace

int main(int argc, char** argv) {
  InitIntrinsic(UsageString(), argc, argv);
  QCHECK_OK(ResetSimulation(absl::GetFlag(FLAGS_address)));
  return 0;
}
