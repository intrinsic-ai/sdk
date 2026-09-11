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

	relayservicepb "intrinsic/platform/pubsub/connect/onprem/onprem_to_line_router_relay_service/onprem_to_line_router_relay_service_go_proto"
)

// OnpremToLineRouterHubServiceCreateCmdRunner installs the relay service on a workcell or VM.
//
// That service connects the workcell / VM to the line-level cloud router.
type OnpremToLineRouterHubServiceCreateCmdRunner struct {
	ServiceInstallingCmdRunner

	orgId  string
	lineId string
}

func (r *OnpremToLineRouterHubServiceCreateCmdRunner) run(ctx context.Context) error {
	config := &relayservicepb.OnpremToLineRouterRelayServiceConfig{
		OrgId:  r.orgId,
		LineId: r.lineId,
	}

	return r.updateInstalledServiceInstances(ctx, config)
}
