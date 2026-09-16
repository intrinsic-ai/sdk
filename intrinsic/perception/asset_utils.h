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

#ifndef INTRINSIC_PERCEPTION_ASSET_UTILS_H_
#define INTRINSIC_PERCEPTION_ASSET_UTILS_H_

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

#endif  // INTRINSIC_PERCEPTION_ASSET_UTILS_H_
