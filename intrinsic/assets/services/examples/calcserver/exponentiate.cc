// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/services/examples/calcserver/exponentiate.h"

#include <cmath>

#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"

namespace intrinsic::services {

grpc::Status ExponentiateServiceImpl::Calculate(
    grpc::ServerContext* context,
    const intrinsic_proto::services::CustomCalculateRequest* request,
    intrinsic_proto::services::CalculatorResponse* response) {
  response->set_result(std::pow(request->x(), request->y()));
  return grpc::Status::OK;
}

}  // namespace intrinsic::services
