// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ASSETS_INSTANCES_CONNECT_CONNECT_H_
#define INTRINSIC_ASSETS_INSTANCES_CONNECT_CONNECT_H_

#include <memory>

#include "absl/status/statusor.h"
#include "grpcpp/channel.h"
#include "grpcpp/support/channel_arguments.h"
#include "intrinsic/assets/proto/v1/grpc_connection.pb.h"

namespace intrinsic::assets::instances::connect {

// Creates a gRPC channel for communicating with the provider specified by
// the GrpcConnection message.
absl::StatusOr<std::shared_ptr<grpc::Channel>> Connect(
    const intrinsic_proto::assets::v1::GrpcConnection& connection,
    const ::grpc::ChannelArguments& channel_args = ::grpc::ChannelArguments());

}  // namespace intrinsic::assets::instances::connect

#endif  // INTRINSIC_ASSETS_INSTANCES_CONNECT_CONNECT_H_
