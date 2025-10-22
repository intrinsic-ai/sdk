// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_EXPONENTIATE_H_
#define INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_EXPONENTIATE_H_

#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.grpc.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"

namespace intrinsic::services {

class ExponentiateServiceImpl
    : public intrinsic_proto::services::CustomCalculation::Service {
 public:
  explicit ExponentiateServiceImpl() = default;

  grpc::Status Calculate(
      grpc::ServerContext* context,
      const intrinsic_proto::services::CustomCalculateRequest* request,
      intrinsic_proto::services::CalculatorResponse* response) override;
};

}  // namespace intrinsic::services

#endif  // INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_EXPONENTIATE_H_
