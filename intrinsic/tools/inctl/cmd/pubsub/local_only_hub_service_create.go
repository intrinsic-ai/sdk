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

package pubsub

import (
	"context"

	endpointpb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/endpoint_spec_go_proto"
	relayrouterpb "intrinsic/platform/pubsub/connect/onprem/relay_router_service/relay_router_service_go_proto"
)

// LocalOnlyHubServiceCreateCmdRunner installs the relay service on the a workcell.
//
// This service supports a local-only line orchestration network, where that workcell is the hub.
type LocalOnlyHubServiceCreateCmdRunner struct {
	ServiceInstallingCmdRunner

	spokeEndpointSpecs []*endpointpb.EndpointSpec
}

// run creates a configuration proto for the relay service,
// and triggers installation of that service.
func (r *LocalOnlyHubServiceCreateCmdRunner) run(ctx context.Context) error {
	config := &relayrouterpb.RelayRouterServiceConfig{
		HubWorkcellName:    r.clusterId,
		SpokeEndpointSpecs: r.spokeEndpointSpecs,
	}
	return r.updateInstalledServiceInstances(ctx, config)
}
