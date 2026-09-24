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

package hubservicedelete

import (
	"context"
	"errors"
	"fmt"
	"io"

	"intrinsic/platform/pubsub/connect/common/lineidutils"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	lineconfigutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/line_config_utils"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"

	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"

	"google.golang.org/grpc"
)

// HubServiceDeleteRunner implements high-level logic of the hub-service-delete command.
//
// Its primary responsibility is to determine the topology of the network (local-only or mixed),
// and delete all constituent parts of that network.
type HubServiceDeleteRunner struct {
	servicedeletionutils.ServiceDeleter

	HubEndpoint    string
	ForceLocalOnly bool

	// dialCloudCluster connects to the cloud-robotics cluster.
	// The default implementation connects to the real cloud-robotics cluster.
	// Unit tests can provide an implementation that connects to a fake server.
	DialCloudCluster func(context.Context) (*grpc.ClientConn, error)
}

func (r *HubServiceDeleteRunner) deleteLineConfig(
	ctx context.Context,
	conn *grpc.ClientConn,
	lineId string) error {
	client := lineconfigstoragepb.NewLineConfigurationStorageServiceClient(conn)
	_, err := client.DeleteLineConfiguration(ctx, &lineconfigstoragepb.DeleteLineConfigurationRequest{
		Name: lineidutils.ConvertIDToName(lineId),
	})
	if err != nil {
		return fmt.Errorf("failed to delete configuration for line %q: %w", lineId, err)
	}

	return nil
}

// Run executes the hub-service-delete command logic to delete the line orchestration network.
func (r *HubServiceDeleteRunner) Run(ctx context.Context, out io.Writer) error {
	hubEndpointSpec, err := common.ParseEndpointSpec(r.HubEndpoint)
	if err != nil {
		return err
	}

	if r.ForceLocalOnly {
		if common.IsRemoteEndpoint(hubEndpointSpec) {
			return fmt.Errorf(
				"The '--%v' flag is present, but the hub endpoint is remote.",
				common.KeyForceLocalOnly)
		}
		fmt.Fprintf(
			out,
			`
The '--%v' flag is present.
Deleting the line orchestration network assuming that it is local-only.
`,
			common.KeyForceLocalOnly)
		return r.DeleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			hubEndpointSpec.WorkcellName,
			common.HubServicePackage,
			common.HubServiceName)
	}

	fmt.Fprintf(out, "Connecting to the cloud API endpoint\n")
	conn, err := r.DialCloudCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to the cloud: %w", err)
	}
	defer conn.Close()

	lineID, err := common.MakeLineID(r.OrgID, hubEndpointSpec.WorkcellName)
	if err != nil {
		fmt.Fprintf(out, "Cannot create a valid line id from %q and %q\n", r.OrgID, hubEndpointSpec.WorkcellName)
		return err
	}

	lineConfig, err := lineconfigutils.FetchLineConfiguration(ctx, conn, lineID)
	if err != nil {
		if errors.Is(err, common.ErrLineConfigNotFound) {
			fmt.Fprintf(
				out,
				`
------------------------------------------------------------------
Configuration for line %q not found.

It is possible that a network for this line had been created
before we started saving line configurations in the cloud.

If that's the case, then that network is local-only.
To delete it, try running this command with the '--%v' flag.
------------------------------------------------------------------
`,
				lineID,
				common.KeyForceLocalOnly,
			)
			return err
		} else if errors.Is(err, common.ErrConfigStorageNotImplemented) {
			common.ExplainConfigStorageNotImplementedError(out, "hub-service-delete", common.KeyForceLocalOnly, "delete")
			return err
		} else {
			return fmt.Errorf("failed to fetch line configuration: %w", err)
		}
	}

	if common.RemoteEndpointsExist(lineConfig.GetSpokeEndpoints(), lineConfig.GetHubEndpoint()) {
		fmt.Fprintf(
			out,
			"Line %q has been identified as mixed (containing local and remote clusters, e.g. VMs).\n",
			lineID)
		if err = r.DeleteMixedLineOrchestrationNetwork(ctx, out, conn, lineID, lineConfig); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(out, "Line %q has been identified as local-only.\n", lineID)
		err = r.DeleteOnpremServiceAndMaybeIgnoreErrors(
			ctx,
			out,
			hubEndpointSpec.WorkcellName,
			common.HubServicePackage,
			common.HubServiceName)
		if err != nil {
			return fmt.Errorf(
				"failed to delete the relay router from %v: %v\n",
				hubEndpointSpec.WorkcellName,
				err)
		}
	}

	fmt.Fprintf(out, "Deleting line configuration stored in the cloud.\n")
	if err := r.deleteLineConfig(ctx, conn, lineID); err != nil {
		return err
	}

	fmt.Fprintf(
		out,
		"Line orchestration network with the hub at %q (line id %q) has been deleted.\n",
		hubEndpointSpec.WorkcellName,
		lineID)
	return nil
}
