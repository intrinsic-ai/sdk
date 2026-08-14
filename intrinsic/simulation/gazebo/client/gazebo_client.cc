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

#include "intrinsic/simulation/gazebo/client/gazebo_client.h"

#include <memory>
#include <string>
#include <string_view>
#include <utility>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/time/time.h"
#include "gz/transport/Node.hh"
#include "gz/transport/NodeOptions.hh"
#include "gz/transport/TopicUtils.hh"
#include "intrinsic/assets/dependencies/utils.h"
#include "intrinsic/assets/proto/v1/resolved_dependency.pb.h"
#include "intrinsic/connect/cc/grpc/channel.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.grpc.pb.h"
#include "intrinsic/simulation/gazebo/proto/v1/gazebo_service.pb.h"
#include "intrinsic/util/grpc/channel.h"
#include "intrinsic/util/grpc/connection_params.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"

namespace intrinsic {
namespace simulation {

namespace {

inline constexpr char kDefaultGzNodePartition[] = "intrinsic_sim";

absl::StatusOr<gz::transport::NodeOptions> CreateNodeOptions(
    std::string_view partition) {
  gz::transport::NodeOptions node_opts;
  node_opts.SetPartition(std::string(partition));

  // Sanity check to make sure the values are valid
  if (!gz::transport::TopicUtils::IsValidNamespace(node_opts.NameSpace()) ||
      !gz::transport::TopicUtils::IsValidPartition(node_opts.Partition())) {
    return absl::InternalError(
        "Node option contains invalid namespace or partition");
  }

  return node_opts;
}

}  // namespace

absl::StatusOr<std::shared_ptr<GazeboClient>> GazeboClient::Create(
    std::string_view sim_server_address, absl::Duration connection_timeout) {
  if (sim_server_address.empty()) {
    return absl::InvalidArgumentError("Sim server address cannot be empty.");
  }
  INTR_ASSIGN_OR_RETURN(auto channel_obj,
                        Channel::MakeFromAddress(
                            ConnectionParams::NoIngress(sim_server_address)));
  return CreateFromChannel(channel_obj->GetChannel(), connection_timeout);
}

absl::StatusOr<std::shared_ptr<GazeboClient>>
GazeboClient::CreateFromResolvedDependency(
    const intrinsic_proto::assets::v1::ResolvedDependency& resolved_deps,
    absl::Duration connection_timeout) {
  // Connect to the Gazebo service hosted by the remote server.
  INTR_ASSIGN_OR_RETURN(
      std::shared_ptr<::grpc::Channel> channel,
      assets::dependencies::Connect(
          resolved_deps,
          /*iface=*/"grpc://intrinsic_proto.simulation.v1.GazeboService",
          connect::UnlimitedMessageSizeGrpcChannelArgs()));

  return CreateFromChannel(std::move(channel), connection_timeout);
}

absl::StatusOr<std::shared_ptr<GazeboClient>> GazeboClient::CreateFromChannel(
    std::shared_ptr<::grpc::Channel> channel,
    absl::Duration connection_timeout) {
  INTR_RETURN_IF_ERROR(
      connect::WaitForChannelReady(channel, connection_timeout));

  std::unique_ptr<intrinsic_proto::simulation::v1::GazeboService::Stub>
      gazebo_service_stub(
          intrinsic_proto::simulation::v1::GazeboService::NewStub(channel));

  intrinsic_proto::simulation::v1::GetIPRequest request;
  intrinsic_proto::simulation::v1::GetIPResponse response;
  ::grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(ToAbslStatus(
      gazebo_service_stub->GetIPForTransportRelay(&ctx, request, &response)));

  INTR_ASSIGN_OR_RETURN(
      gz::transport::NodeOptions node_opts,
      CreateNodeOptions(response.partition().empty() ? kDefaultGzNodePartition
                                                     : response.partition()));

  return std::shared_ptr<GazeboClient>(
      new GazeboClient(std::move(node_opts), response.ip_address(),
                       std::move(gazebo_service_stub)));
}

GazeboClient::GazeboClient(
    gz::transport::NodeOptions node_options, const std::string& sim_ip,
    std::unique_ptr<
        intrinsic_proto::simulation::v1::GazeboService::StubInterface>
        gazebo_service_stub)
    : node_(std::move(node_options)),
      gazebo_service_stub_(std::move(gazebo_service_stub)) {
  // By default, gz-transport does not support communication of nodes on
  // different local networks, e.g. when the sim connection client and the
  // sim server have different subnets. We need to explicitly tell the
  // gz-tranport node to relay traffic to the specified sim server ip.
  node_.AddGlobalRelay(sim_ip);
}

absl::StatusOr<intrinsic_proto::simulation::v1::Topics>
GazeboClient::ListTopics(std::string_view object_name,
                         intrinsic_proto::simulation::v1::TopicEntityTypeFilter
                             entity_type_filter) const {
  intrinsic_proto::simulation::v1::ListTopicsRequest request;
  if (!object_name.empty()) {
    request.set_object_name(std::string(object_name));
  }
  request.set_entity_type_filter(entity_type_filter);
  intrinsic_proto::simulation::v1::Topics response;
  ::grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(gazebo_service_stub_->ListTopics(&ctx, request, &response)));
  return response;
}

gz::transport::Node& GazeboClient::Node() { return node_; }

}  // namespace simulation
}  // namespace intrinsic
