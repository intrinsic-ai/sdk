// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/services/examples/calcserver/calc_server.h"

#include <cstdint>
#include <memory>
#include <string>

#include "absl/log/log.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/dependencies/utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.grpc.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"
#include "intrinsic/util/status/status_macros_grpc.h"

namespace intrinsic::services {

namespace {

char kCustomCalculationInterface[] =
    "grpc://intrinsic_proto.services.CustomCalculation";

}  // namespace

grpc::Status CalculatorServiceImpl::Calculate(
    grpc::ServerContext* context,
    const intrinsic_proto::services::CalculatorRequest* request,
    intrinsic_proto::services::CalculatorResponse* response) {
  int64_t result;

  int64_t a, b;
  if (config_.reverse_order()) {
    a = request->y();
    b = request->x();
  } else {
    a = request->x();
    b = request->y();
  }

  switch (request->operation()) {
    case intrinsic_proto::services::CALCULATOR_OPERATION_ADD:
      result = a + b;
      LOG(INFO) << a << " + " << b << " = " << result;
      break;
    case intrinsic_proto::services::CALCULATOR_OPERATION_MULTIPLY:
      result = a * b;
      LOG(INFO) << a << " * " << b << " = " << result;
      break;
    case intrinsic_proto::services::CALCULATOR_OPERATION_SUBTRACT:
      result = a - b;
      LOG(INFO) << a << " - " << b << " = " << result;
      break;
    case intrinsic_proto::services::CALCULATOR_OPERATION_DIVIDE:
      if (b == 0) {
        LOG(INFO) << "Cannot divide by 0 (" << a << " / " << b << ")";
        return grpc::Status(grpc::StatusCode::INVALID_ARGUMENT,
                            "Cannot divide by 0");
      }
      result = a / b;
      LOG(INFO) << a << " / " << b << " = " << result;
      break;
    case intrinsic_proto::services::CALCULATOR_OPERATION_CUSTOM: {
      // Connect to the CustomCalculation service.
      ::grpc::ClientContext ctx;
      INTR_ASSIGN_OR_RETURN_GRPC(
          std::shared_ptr<grpc::Channel> channel,
          assets::dependencies::Connect(ctx, config_.custom_calculation(),
                                        kCustomCalculationInterface));
      auto stub =
          intrinsic_proto::services::CustomCalculation::NewStub(channel);

      // Create the request.
      intrinsic_proto::services::CustomCalculateRequest custom_request;
      custom_request.set_x(a);
      custom_request.set_y(b);

      // Call the CustomCalculation service.
      intrinsic_proto::services::CalculatorResponse custom_response;
      INTR_RETURN_IF_ERROR_GRPC(
          stub->Calculate(&ctx, custom_request, &custom_response));
      result = custom_response.result();
      LOG(INFO) << "Custom(" << a << ", " << b << ") = " << result;

      break;
    }
    default:
      return grpc::Status(grpc::StatusCode::INVALID_ARGUMENT,
                          "Invalid operation");
  }
  response->set_result(result);
  return grpc::Status::OK;
}

}  // namespace intrinsic::services
