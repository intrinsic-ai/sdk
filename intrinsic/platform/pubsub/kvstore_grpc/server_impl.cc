// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/kvstore_grpc/server_impl.h"

#include <string>
#include <utility>
#include <vector>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "grpcpp/server_context.h"
#include "grpcpp/support/status.h"
#include "intrinsic/platform/pubsub/kvstore.h"
#include "intrinsic/platform/pubsub/kvstore_grpc/kvstore.pb.h"
#include "intrinsic/platform/pubsub/pubsub.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros_grpc.h"

namespace intrinsic::kvstore {

absl::Status KVStoreServerImpl::Init() {
  // Wait for the kvstore to be ready.
  while (true) {
    absl::StatusOr<KeyValueStore> kvstore = pubsub_.KeyValueStore();
    absl::Status status = kvstore.status();
    if (kvstore.ok()) {
      status = kvstore->Set("grpc_kvstore_ready", google::protobuf::Any(),
                            /*high_consistency=*/true);
      if (status.ok()) {
        kvstore_ = *std::move(kvstore);
        break;
      }
    }

    LOG(INFO) << "Waiting for kvstore to be ready: " << status;
    absl::SleepFor(absl::Milliseconds(500));
  }
  return absl::OkStatus();
}

grpc::Status KVStoreServerImpl::Get(
    grpc::ServerContext* context,
    const intrinsic_proto::kvstore::GetRequest* request,
    intrinsic_proto::kvstore::GetResponse* response) {
  LOG(INFO) << "Getting key: " << request->key();
  INTR_ASSIGN_OR_RETURN_GRPC(
      google::protobuf::Any ret,
      kvstore_->Get<google::protobuf::Any>(request->key(), kDefaultGetTimeout));
  *response->mutable_value() = ret;
  return ToGrpcStatus(absl::OkStatus());
}

grpc::Status KVStoreServerImpl::Set(
    grpc::ServerContext* context,
    const intrinsic_proto::kvstore::SetRequest* request,
    intrinsic_proto::kvstore::SetResponse* response) {
  LOG(INFO) << "Setting key: " << request->key();
  return ToGrpcStatus(kvstore_->Set(request->key(), request->value(),
                                    /*high_consistency=*/true));
}

grpc::Status KVStoreServerImpl::Delete(
    grpc::ServerContext* context,
    const intrinsic_proto::kvstore::DeleteRequest* request,
    intrinsic_proto::kvstore::DeleteResponse* response) {
  return ToGrpcStatus(kvstore_->Delete(request->key()));
}

grpc::Status KVStoreServerImpl::List(
    grpc::ServerContext* context,
    const intrinsic_proto::kvstore::ListRequest* request,
    intrinsic_proto::kvstore::ListResponse* response) {
  INTR_ASSIGN_OR_RETURN_GRPC(std::vector<std::string> keys,
                             kvstore_->ListAllKeys());
  response->mutable_keys()->Add(keys.begin(), keys.end());
  return ToGrpcStatus(absl::OkStatus());
}

}  // namespace intrinsic::kvstore
