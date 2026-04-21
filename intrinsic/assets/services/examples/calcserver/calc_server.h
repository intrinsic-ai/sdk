// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_CALC_SERVER_H_
#define INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_CALC_SERVER_H_

#include <functional>

#include "absl/status/statusor.h"
#include "google/protobuf/any.pb.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.grpc.pb.h"
#include "intrinsic/assets/services/examples/calcserver/calc_server.pb.h"
#include "intrinsic/resources/proto/runtime_context.pb.h"

namespace intrinsic::services {

// Fetches data payloads for resolved dependencies.
using DataPayloadFetcher = std::function<absl::StatusOr<google::protobuf::Any>(
    const intrinsic_proto::assets::v1::ResolvedDependency&, absl::string_view)>;

// Performs basic calculator operations.
class CalculatorServiceImpl
    : public intrinsic_proto::services::Calculator::Service {
 public:
  explicit CalculatorServiceImpl(
      const intrinsic_proto::services::CalculatorConfig& config,
      DataPayloadFetcher fetcher)
      : config_(config), fetcher_(fetcher) {}

  grpc::Status Calculate(
      grpc::ServerContext* context,
      const intrinsic_proto::services::CalculatorRequest* request,
      intrinsic_proto::services::CalculatorResponse* response) override;

  grpc::Status Convert(
      grpc::ServerContext* context,
      const intrinsic_proto::services::ConvertRequest* request,
      intrinsic_proto::services::ConvertResponse* response) override;

 private:
  absl::StatusOr<intrinsic_proto::services::ConversionDataset>
  GetConversionDataset();

  intrinsic_proto::services::CalculatorConfig config_;
  DataPayloadFetcher fetcher_;
};

}  // namespace intrinsic::services

#endif  // INTRINSIC_ASSETS_SERVICES_EXAMPLES_CALCSERVER_CALC_SERVER_H_
