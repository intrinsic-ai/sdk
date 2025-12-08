// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/examples/calculate_skill.h"

#include <memory>

#include "absl/log/log.h"
#include "absl/status/statusor.h"
#include "google/protobuf/message.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "intrinsic/assets/dependencies/utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.grpc.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/examples/calculate_skill.pb.h"
#include "intrinsic/skills/proto/skill_service.pb.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace skills {

namespace {

char kCalculatorInterface[] = "grpc://intrinsic_proto.services.Calculator";

}  // namespace

absl::StatusOr<std::unique_ptr<google::protobuf::Message>>
CalculateSkill::Execute(const ExecuteRequest& request,
                        ExecuteContext& context) {
  // Get the Skill's parameters.
  INTR_ASSIGN_OR_RETURN(
      auto params, request.params<intrinsic_proto::skills::CalculateParams>());

  LOG(INFO) << "Calculating " << params.operation() << " with x: " << params.x()
            << " and y: " << params.y();

  // Connect to the Calculator service.
  ::grpc::ClientContext ctx;
  INTR_ASSIGN_OR_RETURN(std::shared_ptr<grpc::Channel> channel,
                        assets::dependencies::Connect(ctx, params.calculator(),
                                                      kCalculatorInterface));
  auto stub = intrinsic_proto::services::Calculator::NewStub(channel);

  // Create the request.
  intrinsic_proto::services::CalculatorRequest calculator_request;
  calculator_request.set_operation(params.operation());
  calculator_request.set_x(params.x());
  calculator_request.set_y(params.y());

  // Call the Calculator service.
  LOG(INFO) << "Calling the Calculator service";
  intrinsic_proto::services::CalculatorResponse response;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(stub->Calculate(&ctx, calculator_request, &response)));
  LOG(INFO) << "Result: " << response.result();

  // Construct the return value.
  auto return_value =
      std::make_unique<intrinsic_proto::skills::CalculateResult>();
  return_value->set_result(response.result());

  return return_value;
}

}  // namespace skills
}  // namespace intrinsic
