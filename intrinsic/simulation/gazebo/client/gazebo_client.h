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

#ifndef INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_H_
#define INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_H_

#include <memory>
#include <string_view>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "gz/transport/Node.hh"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/connect/cc/grpc/channel.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.grpc.pb.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.pb.h"

namespace intrinsic {
namespace simulation {

// GazeboClient manages a connection to Gazebo simulation.
// The connection is established by querying the `GazeboService` gRPC end-point
// that is advertised by a Gazebo simulation server running within an Intrinsic
// application.
// Subsequently, a gz-tranport node is connected to the simulation server to
// communicate with the server over pub/sub.
class GazeboClient {
 public:
  // Create a client with the specified simulation server address.
  // If client creation is unsuccessful, the following status error is
  // returned:
  //   * kUnavailable: The server is unavailable.
  static absl::StatusOr<std::shared_ptr<GazeboClient>> Create(
      std::string_view sim_server_address,
      absl::Duration connection_timeout =
          intrinsic::connect::kGrpcClientConnectDefaultTimeout);

  // Create a client with an Asset's resolved dependency on the `GazeboService`
  // gRPC interface.
  //
  // The gRPC interface should be listed in the passed `ResolvedDependency`
  // proto.
  // Asset sim implementations that use this method should list the following
  // dependency in their configuration protobuf message definition.
  //
  // ```
  // optional intrinsic_proto.assets.v1.ResolvedDependency gazebo_simulator = 1
  //     [(intrinsic_proto.assets.field_metadata).dependency = {
  //       requires: "grpc://intrinsic_proto.simulation.v1.GazeboService",
  //     }];
  // ```
  //
  // For more details on leveraging service -> service dependencies, see
  // https://flowstate.intrinsic.ai/docs/assets/asset_dependencies/interacting_with_other_assets/#service---service-dependencies
  //
  // If client creation is unsuccessful, the following status
  // errors are returned:
  //   * kNotFound: The listed interfaces in `resolved_deps` does not include
  //     the `GazeboService` gRPC interface.
  //   * kUnavailable/kUnimplemented: The server is not available. A
  //     simulator Service Asset may not be installed in the cluster.
  static absl::StatusOr<std::shared_ptr<GazeboClient>>
  CreateFromResolvedDependency(
      const intrinsic_proto::assets::v1::ResolvedDependency& resolved_deps,
      absl::Duration connection_timeout =
          intrinsic::connect::kGrpcClientConnectDefaultTimeout);

  // Call the `ListTopics` gRPC endpoint on `GazeboService`.
  // Returns kUnavailable error if the server is not reachable.
  // Params:
  // object_name: Optional name of world object. Only topics associated with
  //   this object will be returned. Returns empty if the object does not exist
  //   in the Gazebo server or if it doesn't advertise any topics.
  // entity_type_filter: Optional filter to list only Model plugin-related or
  //   Sensor-related topics.
  absl::StatusOr<intrinsic_proto::simulation::v1::Topics> ListTopics(
      std::string_view object_name = "",
      intrinsic_proto::simulation::v1::TopicEntityTypeFilter
          entity_type_filter = intrinsic_proto::simulation::v1::ALL) const;

  // Get a reference to the transport node that manages pub/sub connections
  // to simulation.
  gz::transport::Node& Node();

 private:
  static absl::StatusOr<std::shared_ptr<GazeboClient>> CreateFromChannel(
      std::shared_ptr<::grpc::Channel> channel,
      absl::Duration connection_timeout);

  explicit GazeboClient(
      gz::transport::NodeOptions node_options, const std::string& sim_ip,
      std::unique_ptr<
          intrinsic_proto::simulation::v1::GazeboService::StubInterface>
          gazebo_service_stub);
  gz::transport::Node node_;
  std::unique_ptr<intrinsic_proto::simulation::v1::GazeboService::StubInterface>
      gazebo_service_stub_;
};

}  // namespace simulation
}  // namespace intrinsic

#endif  // INTRINSIC_SIMULATION_GAZEBO_CLIENT_GAZEBO_CLIENT_H_
