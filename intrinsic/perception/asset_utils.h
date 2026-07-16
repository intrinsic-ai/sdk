// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PERCEPTION_PUBLIC_ASSET_UTILS_H_
#define INTRINSIC_PERCEPTION_PUBLIC_ASSET_UTILS_H_

#include "absl/base/nullability.h"
#include "absl/status/statusor.h"
#include "grpcpp/client_context.h"
#include "intrinsic/assets/data/proto/v1/data_assets.grpc.pb.h"
#include "intrinsic/assets/proto/id.pb.h"
#include "intrinsic/perception/proto/v1/perception_model.pb.h"

namespace intrinsic::perception {

absl::StatusOr<intrinsic_proto::perception::v1::PerceptionModel>
GetPerceptionModelFromDataAsset(
    intrinsic_proto::assets::Id id,
    intrinsic_proto::data::v1::DataAssets::StubInterface& data_assets_stub,
    grpc::ClientContext* absl_nullable client_context = nullptr);

}  // namespace intrinsic::perception

#endif  // INTRINSIC_PERCEPTION_PUBLIC_ASSET_UTILS_H_
