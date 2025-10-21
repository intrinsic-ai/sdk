// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/dependencies/utils.h"

#include <memory>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "grpcpp/create_channel.h"
#include "grpcpp/security/credentials.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"

namespace intrinsic::assets::dependencies {

namespace {

using ::intrinsic_proto::assets::v1::ResolvedDependency;

}  // namespace

absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect(
    grpc::ClientContext& context, const ResolvedDependency& dep,
    absl::string_view iface) {
  const auto& interfaces = dep.interfaces();
  const auto it = interfaces.find(std::string(iface));
  if (it == interfaces.end()) {
    return absl::NotFoundError(
        absl::StrCat("Interface ", iface, " not found in resolved dependency"));
  }

  const auto& iface_proto = it->second;
  if (!iface_proto.has_grpc_connection()) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Interface ", iface,
        " is not gRPC or no connection information is available."));
  }

  // Add any needed metadata to the context.
  for (const auto& metadata : iface_proto.grpc_connection().metadata()) {
    context.AddMetadata(metadata.key(), metadata.value());
  }

  return ::grpc::CreateChannel(
      iface_proto.grpc_connection().address(),
      grpc::InsecureChannelCredentials());  // NOLINT(insecure)
}

}  // namespace intrinsic::assets::dependencies
