// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/assets/dependencies/utils.h"

#include <memory>
#include <string>
#include <vector>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "grpcpp/create_channel.h"
#include "grpcpp/security/credentials.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic::assets::dependencies {

namespace {

using ::intrinsic_proto::assets::v1::ResolvedDependency;

char kIngressAddress[] =
    "istio-ingressgateway.app-ingress.svc.cluster.local:80";

absl::StatusOr<const ResolvedDependency::Interface*> FindInterface(
    const ResolvedDependency& dep, absl::string_view iface) {
  const auto it = dep.interfaces().find(std::string(iface));
  if (it == dep.interfaces().end()) {
    std::string explanation;
    if (dep.interfaces().empty()) {
      explanation = "no interfaces provided";
    } else {
      std::vector<std::string> keys;
      keys.reserve(dep.interfaces().size());
      for (const auto& [key, _] : dep.interfaces()) {
        keys.push_back(key);
      }
      explanation = absl::StrCat("got interfaces: ", absl::StrJoin(keys, ", "));
    }
    return absl::NotFoundError(
        absl::StrCat("Interface not found in resolved dependency (want ", iface,
                     ", ", explanation, ")"));
  }
  return &it->second;
}

std::unique_ptr<intrinsic_proto::data::v1::DataAssets::StubInterface>
MakeDefaultDataAssetsClient() {
  return intrinsic_proto::data::v1::DataAssets::NewStub(::grpc::CreateChannel(
      kIngressAddress,
      grpc::InsecureChannelCredentials()));  // NOLINT(insecure)
}

}  // namespace

absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect(
    grpc::ClientContext& context, const ResolvedDependency& dep,
    absl::string_view iface) {
  INTR_ASSIGN_OR_RETURN(const auto* iface_proto, FindInterface(dep, iface));
  if (!iface_proto->has_grpc_connection()) {
    return absl::InvalidArgumentError(absl::StrCat(
        "Interface is not gRPC or no connection information is available: ",
        iface));
  }

  // Add any needed metadata to the context.
  for (const auto& metadata : iface_proto->grpc_connection().metadata()) {
    context.AddMetadata(metadata.key(), metadata.value());
  }

  return ::grpc::CreateChannel(
      iface_proto->grpc_connection().address(),
      grpc::InsecureChannelCredentials());  // NOLINT(insecure)
}

absl::StatusOr<google::protobuf::Any> GetDataPayload(
    const ResolvedDependency& dep, absl::string_view iface,
    intrinsic_proto::data::v1::DataAssets::StubInterface* data_assets_client) {
  INTR_ASSIGN_OR_RETURN(const auto* iface_proto, FindInterface(dep, iface));
  if (!iface_proto->has_data()) {
    return absl::InvalidArgumentError(
        absl::StrCat("Interface is not data or no data dependency information "
                     "is available: ",
                     iface));
  }

  std::unique_ptr<intrinsic_proto::data::v1::DataAssets::StubInterface>
      default_data_assets_client;
  if (data_assets_client == nullptr) {
    default_data_assets_client = MakeDefaultDataAssetsClient();
    data_assets_client = default_data_assets_client.get();
  }

  // Get the DataAsset proto from the DataAssets service.
  intrinsic_proto::data::v1::GetDataAssetRequest request;
  *request.mutable_id() = iface_proto->data().id();
  intrinsic_proto::data::v1::DataAsset da;
  grpc::ClientContext context;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(data_assets_client->GetDataAsset(&context, request, &da)));

  return da.data();
}

}  // namespace intrinsic::assets::dependencies
