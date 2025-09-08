// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_STATUS_STATUS_CONVERSION_GRPC_H_
#define INTRINSIC_UTIL_STATUS_STATUS_CONVERSION_GRPC_H_

#include "absl/status/status.h"
#include "google/rpc/status.pb.h"
#include "grpcpp/support/status.h"

namespace intrinsic {

grpc::Status ToGrpcStatus(const absl::Status& status);
grpc::Status ToGrpcStatus(const google::rpc::Status& status);
absl::Status ToAbslStatus(const grpc::Status& status);

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_STATUS_STATUS_CONVERSION_GRPC_H_
