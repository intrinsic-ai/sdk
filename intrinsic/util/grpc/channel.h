// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_GRPC_CHANNEL_H_
#define INTRINSIC_UTIL_GRPC_CHANNEL_H_

#include <memory>
#include <string>
#include <string_view>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "grpcpp/channel.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/grpc/grpc.h"

namespace intrinsic {

// A channel to an Intrinsic gRPC service at the specified address.
class Channel : public ChannelInterface {
 public:
  // Creates a channel to an Intrinsic gRPC service based on the provided
  // connection parameters.  `timeout` specifies the maximum amount of time to
  // wait for a response from the server before giving up on creating a channel.
  static absl::StatusOr<std::shared_ptr<Channel>> MakeFromAddress(
      const ConnectionParams& params,
      absl::Duration timeout = kGrpcClientConnectDefaultTimeout);

  // Constructs a Channel with given connection parameters.
  explicit Channel(std::shared_ptr<grpc::Channel> channel,
                   std::string_view instance_name = "",
                   std::string_view header = "x-resource-instance-name");

  std::shared_ptr<grpc::Channel> GetChannel() const override;

  ClientContextFactory GetClientContextFactory() const override;

 private:
  std::shared_ptr<grpc::Channel> channel_;

  // The ingress instance name.  This determines which VirtualService in
  // kubernetes is targeted.  If empty, the header information is not added to
  // the gRPC connection.
  const std::string instance_name_;
  // The header to be used when establishing a gRPC connection to the ingress.
  // The header's value will be instance_name.
  const std::string header_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_GRPC_CHANNEL_H_
