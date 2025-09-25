// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/util/grpc/channel.h"

#include <memory>
#include <string_view>
#include <utility>

#include "absl/status/statusor.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "grpcpp/channel.h"
#include "grpcpp/client_context.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/grpc/grpc.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {

absl::StatusOr<std::shared_ptr<Channel>> Channel::MakeFromAddress(
    const ConnectionParams& params, absl::Duration timeout) {
  // Set the max message size to unlimited to allow longer trajectories.
  // Please check with the motion team before changing the value (see
  // b/275280379).
  INTR_ASSIGN_OR_RETURN(
      std::shared_ptr<grpc::Channel> channel,
      CreateClientChannel(params.address, absl::Now() + timeout,
                          UnlimitedMessageSizeGrpcChannelArgs()));
  return std::shared_ptr<Channel>(
      new Channel(channel, params.instance_name, params.header));
}

std::shared_ptr<grpc::Channel> Channel::GetChannel() const { return channel_; }

ClientContextFactory Channel::GetClientContextFactory() const {
  return [header = header_, instance_name = instance_name_]() {
    auto context = std::make_unique<::grpc::ClientContext>();
    ConfigureClientContext(context.get());
    if (!header.empty() && !instance_name.empty()) {
      context->AddMetadata(header, instance_name);
    }
    return context;
  };
}

Channel::Channel(std::shared_ptr<grpc::Channel> channel,
                 std::string_view instance_name, std::string_view header)
    : channel_(std::move(channel)),
      instance_name_(instance_name),
      header_(header) {}

}  // namespace intrinsic
