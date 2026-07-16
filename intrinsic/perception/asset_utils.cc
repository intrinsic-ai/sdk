// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/perception/asset_utils.h"

#include "absl/base/nullability.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "grpcpp/client_context.h"
#include "intrinsic/assets/data/proto/v1/data_asset.pb.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/data/proto/v1/data_assets.pb.h"
#include "intrinsic/perception/proto/v1/perception_model.pb.h"
#include "intrinsic/util/status/status_macros_grpc.h"

namespace intrinsic::perception {
namespace {
constexpr char kDefaultAssetPackage[] = "ai.intrinsic";

}  // namespace

absl::StatusOr<intrinsic_proto::perception::v1::PerceptionModel>
GetPerceptionModelFromDataAsset(
    intrinsic_proto::assets::Id id,
    intrinsic_proto::data::v1::DataAssets::StubInterface& data_assets_stub,
    grpc::ClientContext* absl_nullable client_context) {
  if (id.package().empty()) {
    id.set_package(kDefaultAssetPackage);
  }
  intrinsic_proto::data::v1::GetDataAssetRequest get_request;
  intrinsic_proto::data::v1::DataAsset data_asset;
  *get_request.mutable_id() = id;

  grpc::ClientContext local_context;
  if (client_context == nullptr) client_context = &local_context;
  auto status = ToAbslStatus(
      data_assets_stub.GetDataAsset(client_context, get_request, &data_asset));
  if (!status.ok()) {
    return absl::InternalError(absl::StrFormat(
        "Failed to get data asset. Error: %s", status.message()));
  }

  intrinsic_proto::perception::v1::PerceptionModel perception_model;
  if (!data_asset.data().UnpackTo(&perception_model)) {
    return absl::InternalError(
        absl::StrFormat("Failed to unpack data asset to v1::perception model. "
                        "Data asset id: %s",
                        id.name()));
  }
  return perception_model;
}
}  // namespace intrinsic::perception
