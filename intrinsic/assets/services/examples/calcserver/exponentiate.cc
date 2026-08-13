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
