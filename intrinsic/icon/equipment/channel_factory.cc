// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/equipment/channel_factory.h"

#include <memory>

#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/grpc/grpc.h"

namespace intrinsic {
namespace icon {

absl::StatusOr<std::shared_ptr<ChannelInterface>> ChannelFactory::MakeChannel(
    const ConnectionParams& params) const {
  return MakeChannel(params, connect::kGrpcClientConnectDefaultTimeout);
}

absl::StatusOr<std::shared_ptr<ChannelInterface>>
DefaultChannelFactory::MakeChannel(const ConnectionParams& params,
                                   absl::Duration timeout) const {
  return Channel::MakeFromAddress(params, timeout);
}

}  // namespace icon
}  // namespace intrinsic
