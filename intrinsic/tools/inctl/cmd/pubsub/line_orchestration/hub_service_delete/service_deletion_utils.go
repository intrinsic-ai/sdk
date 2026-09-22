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

package servicedeletionutils

import (
	"context"
	"fmt"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
)

const (
	KeyUninstallServiceAsset = "uninstall-service-asset"
	KeyIgnoreOnpremErrors    = "ignore-onprem-errors"
)

// ServiceDeleteCmdRunner handles execution of commands that delete service assets
// used for line orchestration.
type ServiceDeleteCmdRunner struct {
	common.CmdRunnerBase

	ShouldUninstallServiceAsset bool
}

// run implements the core logic of the command:
//   - Deletes all instances of the service.
//   - Optionally uninstalls the service asset.
func (r *ServiceDeleteCmdRunner) Run(ctx context.Context) error {
	fmt.Fprintf(
		r.OutputWriter,
		"Deleting existing instances of the %v service in the current solution.\n",
		r.ServiceName)
	numInstances, err := r.DeleteExistingServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to delete existing instances of the %v service: %w",
			r.ServiceName, err)
	}
	fmt.Fprintf(r.OutputWriter, "%v instances have been deleted.\n", numInstances)

	if !r.ShouldUninstallServiceAsset {
		fmt.Fprintf(
			r.OutputWriter,
			"The %v option is disabled, won't try to uninstall the %v service asset.\n",
			KeyUninstallServiceAsset,
			r.ServiceName)
		return nil
	}

	shouldUninstall := false
	if numInstances > 0 {
		// If there were instances of the service, then the service asset
		// must be installed. No need to check.
		shouldUninstall = true
	} else {
		// There were no instances of the service, but it is still possible that
		// the service asset is installed. Checking it here.
		fmt.Fprintf(
			r.OutputWriter,
			"Checking if the %v service asset is installed.\n",
			r.ServiceName)
		currentVersion, err := r.GetInstalledServiceAssetVersion(ctx)
		if err != nil {
			return fmt.Errorf(
				"failed to determine whether the %v service is installed: %w",
				r.ServiceName,
				err)
		}
		if len(currentVersion) != 0 {
			fmt.Fprintf(
				r.OutputWriter,
				"%v service, version %v, is installed. Will uninstall.\n", r.ServiceName, currentVersion)
			shouldUninstall = true
		} else {
			fmt.Fprintf(
				r.OutputWriter,
				"%v service asset is not installed, nothing else to do.\n",
				r.ServiceName)
		}
	}

	if shouldUninstall {
		fmt.Fprintf(r.OutputWriter, "Uninstalling %v service asset.\n", r.ServiceName)
		if err = r.UninstallServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to uninstall %v service asset: %w", r.ServiceName, err)
		}
		fmt.Fprintf(
			r.OutputWriter,
			"The %v service asset has been uninstalled.\n",
			r.ServiceName)
	}

	return nil
}
