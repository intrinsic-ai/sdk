// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_
#define INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_

#include <memory>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "grpcpp/client_context.h"
#include "grpcpp/grpcpp.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"

namespace intrinsic::assets::dependencies {

// Creates a gRPC channel for communicating with the provider of the specified
// interface.
//
// The context will be updated with any needed metadata for communicating with
// the provider.
absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect(
    grpc::ClientContext& context,
    const intrinsic_proto::assets::v1::ResolvedDependency& dep,
    absl::string_view iface);

}  // namespace intrinsic::assets::dependencies

#endif  // INTRINSIC_ASSETS_DEPENDENCIES_UTILS_H_
