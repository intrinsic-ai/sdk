// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_
#define INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_

#include <memory>
#include <optional>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "grpcpp/channel.h"
#include "grpcpp/support/channel_arguments.h"

namespace intrinsic::connect {

// Default timeout for the initial GRPC connection made by client libraries.
inline constexpr absl::Duration kGrpcClientConnectDefaultTimeout =
    absl::Seconds(5);

// Wait for a newly created channel to be connected
absl::Status WaitForChannelConnected(absl::string_view address,
                                     std::shared_ptr<::grpc::Channel> channel,
                                     absl::Time deadline = absl::Now());

// Get recommended default gRPC channel arguments.
::grpc::ChannelArguments DefaultGrpcChannelArgs();

// Get gRPC channel arguments with unlimited send/receive message size.
// This also includes all settings from DefaultGrpcChannelArgs(). This can be
// used for services that send large messages, e.g., the geometry service.
::grpc::ChannelArguments UnlimitedMessageSizeGrpcChannelArgs();

// Apply default configuration of our project and create a new channel.
absl::StatusOr<std::shared_ptr<::grpc::Channel>> CreateClientChannel(
    absl::string_view address, absl::Time deadline,
    const ::grpc::ChannelArguments& channel_args =
        connect::DefaultGrpcChannelArgs(),
    bool use_default_application_credentials = false,
    std::optional<std::string> server_instance_name = std::nullopt);

}  // namespace intrinsic::connect

#endif  // INTRINSIC_CONNECT_CC_GRPC_CHANNEL_H_
