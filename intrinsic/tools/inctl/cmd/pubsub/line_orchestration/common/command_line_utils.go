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
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	KeyHubEndpoint       = "hub-endpoint"
	KeySpokeEndpoints    = "spoke-endpoint"
	KeyHubServiceVersion = "hub-service-version"
	KeyForceLocalOnly    = "force-local-only"

	// The `cluster` flag in `hub-service-create / delete` commands is deprecated
	// because it can only contain a cluster id. `hub-service-create` and `delete`
	// commands need additional information (whether that cluster is on a local
	// network or in the cloud). This information can be provided in the `hub-endpoint`
	// flag.
	KeyClusterDeprecated = "cluster"

	// KeyRetainServiceAsset is the command line flag that specifies whether to
	// retain service asset in a solution.
	KeyRetainServiceAsset = "retain-service-asset"

	// KeyIgnoreOnpremErrors is the command line flag that specifies whether to
	// ignore errors that occur in onprem clusters.
	KeyIgnoreOnpremErrors = "ignore-onprem-errors"

	EndpointSpecSeparator     = "@"
	LocalEndpointDesignation  = "local"
	RemoteEndpointDesignation = "remote"
)

// GetHubEndpoint computes the hub endpoint spec.
//
// It is intended to preserve backwards compatibility with previous versions of
// the hub-service-create command, where the `--cluster` flag was used to specify
// the hub. Newer versions use the `--hub-endpoint` flag. Unlike `--cluster`, that
// flag allows `@local` and `@remote` suffixes. This, in turn, enables users to
// unambiguously specify whether the hub is on a local network or a remote server.
func GetHubEndpoint(out io.Writer, cluster, hubEndpoint string) (string, error) {
	if len(hubEndpoint) != 0 {
		return hubEndpoint, nil
	}

	if len(cluster) != 0 {
		fmt.Fprintf(out, "WARNING: the --cluster flag is deprecated, use --%s instead.\n", KeyHubEndpoint)
		if strings.HasPrefix(cluster, "vmp-") {
			fmt.Fprintf(out, "Assuming that %q is a remote server.\n", cluster)
			return fmt.Sprintf("%s%s%s", cluster, EndpointSpecSeparator, RemoteEndpointDesignation), nil
		}

		fmt.Fprintf(out, "Assuming that %q is on the local network.\n", cluster)
		return fmt.Sprintf("%s%s%s", cluster, EndpointSpecSeparator, LocalEndpointDesignation), nil
	}

	return "", errors.New("neither --cluster nor --hub-endpoint is specified")
}
