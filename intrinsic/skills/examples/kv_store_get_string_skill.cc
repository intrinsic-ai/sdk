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

#include "intrinsic/skills/examples/kv_store_get_string_skill.h"

#include <memory>
#include <string>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "google/protobuf/message.h"
#include "google/protobuf/wrappers.pb.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "intrinsic/assets/dependencies/utils.h"
#include "intrinsic/assets/interface_utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/platform/pubsub/kvstore_grpc/kvstore.grpc.pb.h"
#include "intrinsic/platform/pubsub/kvstore_grpc/kvstore.pb.h"
#include "intrinsic/skills/cc/skill_interface.h"
#include "intrinsic/skills/examples/kv_store_get_string_skill.pb.h"
#include "intrinsic/skills/proto/skill_service.pb.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace skills {

namespace {

std::string KvStoreInterfaceUri() {
  return absl::StrCat(intrinsic::assets::kGrpcUriPrefix,
                      intrinsic_proto::kvstore::KVStore::service_full_name());
}

}  // namespace

absl::StatusOr<std::unique_ptr<google::protobuf::Message>>
KvStoreGetStringSkill::Execute(const ExecuteRequest& request,
                               ExecuteContext& context) {
  INTR_ASSIGN_OR_RETURN(
      auto params,
      request.params<intrinsic_proto::skills::KvStoreGetStringParams>());

  // Connect to the KVStore service.
  INTR_ASSIGN_OR_RETURN(std::shared_ptr<grpc::Channel> channel,
                        assets::dependencies::Connect(
                            params.intrinsic_runtime(), KvStoreInterfaceUri()));
  auto stub = intrinsic_proto::kvstore::KVStore::NewStub(channel);

  intrinsic_proto::kvstore::GetRequest get_request;
  get_request.set_key(params.key());

  grpc::ClientContext client_context;
  intrinsic_proto::kvstore::GetResponse get_response;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(stub->Get(&client_context, get_request, &get_response)));

  auto result =
      std::make_unique<intrinsic_proto::skills::KvStoreGetStringResult>();
  if (get_response.has_value()) {
    google::protobuf::StringValue string_value;
    if (!get_response.value().UnpackTo(&string_value)) {
      return absl::InvalidArgumentError(
          "Stored value is not a google.protobuf.StringValue");
    }
    LOG(INFO) << "Retrieved value for key '" << params.key()
              << "': " << string_value.value();
    *result->mutable_value() = string_value;
  }

  return result;
}

}  // namespace skills
}  // namespace intrinsic
