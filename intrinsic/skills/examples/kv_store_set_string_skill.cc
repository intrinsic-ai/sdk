// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/skills/examples/kv_store_set_string_skill.h"

#include <memory>
#include <string>

#include "absl/log/log.h"
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
#include "intrinsic/skills/examples/kv_store_set_string_skill.pb.h"
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
KvStoreSetStringSkill::Execute(const ExecuteRequest& request,
                               ExecuteContext& context) {
  INTR_ASSIGN_OR_RETURN(
      auto params,
      request.params<intrinsic_proto::skills::KvStoreSetStringParams>());

  // Connect to the KVStore service.
  INTR_ASSIGN_OR_RETURN(std::shared_ptr<grpc::Channel> channel,
                        assets::dependencies::Connect(
                            params.intrinsic_runtime(), KvStoreInterfaceUri()));
  auto stub = intrinsic_proto::kvstore::KVStore::NewStub(channel);

  intrinsic_proto::kvstore::SetRequest set_request;
  set_request.set_key(params.key());
  google::protobuf::StringValue string_value;
  string_value.set_value(params.value());
  set_request.mutable_value()->PackFrom(string_value);

  grpc::ClientContext client_context;
  intrinsic_proto::kvstore::SetResponse set_response;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(stub->Set(&client_context, set_request, &set_response)));
  LOG(INFO) << "Set value for key '" << params.key() << "': " << params.value();

  return std::make_unique<intrinsic_proto::skills::KvStoreSetStringResult>();
}

}  // namespace skills
}  // namespace intrinsic
