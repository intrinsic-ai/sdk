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

#ifndef INTRINSIC_ICON_EQUIPMENT_CHANNEL_FACTORY_H_
#define INTRINSIC_ICON_EQUIPMENT_CHANNEL_FACTORY_H_

#include <memory>

#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "intrinsic/util/grpc/channel_interface.h"
#include "intrinsic/util/grpc/connection_params.h"

namespace intrinsic {
namespace icon {

// ChannelFactory instantiates a intrinsic::Channel given a gRPC address.
//
// This is needed for skill implementations, which expect to connect to ICON
// using a gRPC address from icon::proto::EquipmentConfig. In order for skill
// implementations to be testable, a ChannelFactory is injected and used. This
// gives us a level of indirection, allowing tests to use a FakeChannelFactory
// and normal binaries to use a DefaultChannelFactory.
//
// Non-skill code should typically use intrinsic::Channel and icon::ChannelFake
// directly. Their base class, ChannelInterface, provides a level of
// abstraction suitable for most normal usage and testing. It is only when you
// _must_ create a Channel from a gRPC address (and want to be able to test such
// code against a ChannelFake) that this additional level of indirection is
// needed.
class ChannelFactory {
 public:
  virtual ~ChannelFactory() = default;

  virtual absl::StatusOr<std::shared_ptr<ChannelInterface>> MakeChannel(
      const ConnectionParams& params, absl::Duration timeout) const = 0;

  absl::StatusOr<std::shared_ptr<ChannelInterface>> MakeChannel(
      const ConnectionParams& params) const;
};

// DefaultChannelFactory creates a Channel that connects to a gRPC address.
class DefaultChannelFactory : public ChannelFactory {
  absl::StatusOr<std::shared_ptr<ChannelInterface>> MakeChannel(
      const ConnectionParams& params, absl::Duration timeout) const override;
};

}  // namespace icon
}  // namespace intrinsic

#endif  // INTRINSIC_ICON_EQUIPMENT_CHANNEL_FACTORY_H_
