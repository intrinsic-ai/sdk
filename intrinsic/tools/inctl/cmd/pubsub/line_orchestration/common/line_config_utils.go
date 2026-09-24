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

package lineconfigutils

import (
	"context"

	"intrinsic/platform/pubsub/connect/common/lineidutils"
	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"

	lineconfigstoragepb "intrinsic/platform/pubsub/connect/cloud/proto/line_configuration_storage/v1/line_configuration_storage_go_proto"
	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"
	lineconfigpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/line_configuration_go_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FetchLineConfiguration retrieves line configuration from the cloud.
func FetchLineConfiguration(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	lineID string) (*lineconfigpb.LineConfiguration, error) {
	client := lineconfigstoragepb.NewLineConfigurationStorageServiceClient(conn)
	req := &lineconfigstoragepb.GetLineConfigurationRequest{
		Name: lineidutils.ConvertIDToName(lineID),
	}
	lineConfig, err := client.GetLineConfiguration(ctx, req)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, common.ErrLineConfigNotFound
		case codes.Unimplemented:
			return nil, common.ErrConfigStorageNotImplemented
		}

		return nil, err
	}

	return lineConfig, nil
}

// saveLineConfiguration persists configuration of the factory line in the cloud.
func SaveLineConfiguration(
	ctx context.Context,
	conn grpc.ClientConnInterface,
	lineID string,
	hubEndpointSpec *endpointpb.EndpointSpec,
	spokeEndpointSpecs []*endpointpb.EndpointSpec) error {
	client := lineconfigstoragepb.NewLineConfigurationStorageServiceClient(conn)
	req := &lineconfigstoragepb.UpdateLineConfigurationRequest{
		LineConfiguration: &lineconfigpb.LineConfiguration{
			Name:           lineidutils.ConvertIDToName(lineID),
			HubEndpoint:    hubEndpointSpec,
			SpokeEndpoints: spokeEndpointSpecs,
		},
		AllowMissing: true,
	}
	_, err := client.UpdateLineConfiguration(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return common.ErrConfigStorageNotImplemented
		}
		return err
	}

	return nil
}
