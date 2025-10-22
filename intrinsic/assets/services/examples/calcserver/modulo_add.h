// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_MODULO_ADD_H_
#define INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_MODULO_ADD_H_

#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.grpc.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/assets/services/examples/calcserver/modulo_add.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"

namespace intrinsic::services {

class ModuloAddServiceImpl
    : public intrinsic_proto::services::CustomCalculation::Service {
 public:
  explicit ModuloAddServiceImpl(
      const intrinsic_proto::services::ModuloAddServiceConfig& config)
      : config_(config) {}

  grpc::Status Calculate(
      grpc::ServerContext* context,
      const intrinsic_proto::services::CustomCalculateRequest* request,
      intrinsic_proto::services::CalculatorResponse* response) override;

 private:
  intrinsic_proto::services::ModuloAddServiceConfig config_;
};

}  // namespace intrinsic::services

#endif  // INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_MODULO_ADD_H_
