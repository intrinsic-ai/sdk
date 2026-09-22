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

package common

import (
	"fmt"
	"slices"
	"strings"

	"intrinsic/platform/pubsub/connect/common/lineidutils"

	endpointpb "intrinsic/platform/pubsub/connect/common/proto/line_configuration/v1/endpoint_spec_go_proto"
)

// IsRemoteEndpoint checks whether the endpoint spec points to a remote endpoint.
func IsRemoteEndpoint(endpointSpec *endpointpb.EndpointSpec) bool {
	return (endpointSpec.GetRemote() != nil)
}

// RemoteEndpointsExist checks whether any endpoint specified in the command line
// points to a remote endpoint.
func RemoteEndpointsExist(spokeEndpointSpecs []*endpointpb.EndpointSpec, hubEndpointSpec *endpointpb.EndpointSpec) bool {
	return IsRemoteEndpoint(hubEndpointSpec) || slices.ContainsFunc(spokeEndpointSpecs, IsRemoteEndpoint)
}

// MakeLineID creates a line identifier based on the organization id and the hub cluster id.
//
// The resulting identifier consists of the organization id and the cluster id, separated by a hyphen.
// All underscores are replaced with hyphens.
//
// Examples:
//
//	MakeLineID("my_org", "vmp-123") = "my-org-vmp-123"
func MakeLineID(orgID, clusterID string) (string, error) {
	lineID := fmt.Sprintf(
		"%s-%s",
		strings.ReplaceAll(orgID, "_", "-"),
		strings.ReplaceAll(clusterID, "_", "-"))
	if err := lineidutils.Validate(lineID); err != nil {
		return "", err
	}
	return lineID, nil
}

// parseEndpointSpec creates an EndpointSpec based on a spoke-endpoint or hub-endpoint command line flag.
func ParseEndpointSpec(flagValue string) (*endpointpb.EndpointSpec, error) {
	parts := strings.Split(flagValue, EndpointSpecSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf(
			"Failed to parse %q. Each endpoint spec should consist of two parts separated by %v",
			flagValue, EndpointSpecSeparator)
	}

	result := &endpointpb.EndpointSpec{
		WorkcellName: parts[0],
	}

	switch parts[1] {
	case LocalEndpointDesignation:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Local{
			Local: &endpointpb.LocalConnectionSpec{},
		}
	case RemoteEndpointDesignation:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Remote{
			Remote: &endpointpb.RemoteConnectionSpec{},
		}
	default:
		result.ConnectionSpec = &endpointpb.EndpointSpec_Url{Url: parts[1]}
	}

	return result, nil
}

// parseEndpointSpecs parses all spoke endpoint specs.
func ParseEndpointSpecs(endpointSpecStrings []string) ([]*endpointpb.EndpointSpec, error) {
	result := make([]*endpointpb.EndpointSpec, 0, len(endpointSpecStrings))
	for _, spec := range endpointSpecStrings {
		endpointSpec, err := ParseEndpointSpec(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to parse endpoint spec %q: %w", spec, err)
		}
		result = append(result, endpointSpec)
	}

	return result, nil
}
