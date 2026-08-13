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

#include "intrinsic/icon/equipment/channel_factory.h"

#include <memory>

#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/connect/cc/grpc/channel.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"

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
