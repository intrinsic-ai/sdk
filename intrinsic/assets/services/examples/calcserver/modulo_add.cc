// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/services/examples/calcserver/modulo_add.h"

#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"

namespace intrinsic::services {

grpc::Status ModuloAddServiceImpl::Calculate(
    grpc::ServerContext* context,
    const intrinsic_proto::services::CustomCalculateRequest* request,
    intrinsic_proto::services::CalculatorResponse* response) {
  response->set_result((request->x() + request->y()) % config_.modulus());
  return grpc::Status::OK;
}

}  // namespace intrinsic::services
