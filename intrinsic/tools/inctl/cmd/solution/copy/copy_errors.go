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

package copy

import (
	"fmt"

	"intrinsic/util/status/extstatus"

	"google.golang.org/grpc/codes"
)

const (
	// copyNotConfirmedErrorCode indicates that the user did not confirm the copy at the prompt.
	copyNotConfirmedErrorCode = iota + 10000
	// descriptionWithoutTitleErrorCode indicates that --description was set without --title.
	descriptionWithoutTitleErrorCode
	// emptyTargetErrorCode indicates that the target flag was passed with an empty value.
	emptyTargetErrorCode
)

// componentName identifies extended status errors from `inctl solution copy`.
const componentName = "ai.intrinsic.inctl.solution.copy"

// copyNotConfirmedError returns the error for a copy that was not confirmed at the prompt,
// including when stdin is not interactive.
func copyNotConfirmedError() error {
	return withReport(extstatus.NewError(componentName, copyNotConfirmedErrorCode,
		extstatus.WithTitle("Copy not confirmed"),
		extstatus.WithGrpcCode(codes.Canceled),
		extstatus.WithUserMessage("The copy was not started because it was not confirmed."),
		extstatus.WithUserInstructions(`Answer "y" at the prompt, or pass --yes to start the copy without confirmation.`),
	))
}

// descriptionWithoutTitleError returns the error for --description set without --title. The
// service replaces the whole commit message when the title is empty (b/507928991), so the
// description would be silently lost.
func descriptionWithoutTitleError() error {
	return withReport(extstatus.NewError(componentName, descriptionWithoutTitleErrorCode,
		extstatus.WithTitle("--description requires --title"),
		extstatus.WithGrpcCode(codes.InvalidArgument),
		extstatus.WithUserMessage("The commit description is only kept when a commit title is set."),
		extstatus.WithUserInstructions("Pass --title as well, or drop --description to let the service generate the commit message."),
	))
}

// emptyTargetError returns the error for a --target-branch-id or --target-branch-display-name
// passed with an empty value, which the flag groups alone do not reject.
func emptyTargetError() error {
	return withReport(extstatus.NewError(componentName, emptyTargetErrorCode,
		extstatus.WithTitle("Empty copy target"),
		extstatus.WithGrpcCode(codes.InvalidArgument),
		extstatus.WithUserMessage("--target-branch-id and --target-branch-display-name must not be empty."),
		extstatus.WithUserInstructions("Pass a non-empty value to --target-branch-id or --target-branch-display-name."),
	))
}

// withReport adds the full extended status report to err's message, because the inctl root
// command only prints err.Error(), which for extended status errors is just the title. The
// returned error still unwraps to err.
func withReport(err error) error {
	es, ok := extstatus.FromError(err)
	if !ok {
		return err
	}
	return fmt.Errorf("%w\n%v", err, es)
}
