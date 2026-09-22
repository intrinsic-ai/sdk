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

// Package lineidutils provides utilities for sanitizing factory line IDs.
package lineidutils

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	lineNamePrefix = "factoryLines"

	// Max length or a resource name allowed by Kubernetes is 253 characters.
	// So the line id can be at least 237 characters long so that
	// `app-<line id>-line-router` fits into 253 characters.
	maxLineIDLength = 237
)

var (
	validLineIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
	lineNamePattern    = regexp.MustCompile(fmt.Sprintf(`%v/([a-z0-9-]+)`, lineNamePrefix))
)

func Validate(lineID string) error {
	if len(lineID) == 0 {
		return errors.New("line id is empty")
	}

	if len(lineID) > maxLineIDLength {
		return fmt.Errorf(
			"line id %q is more than %d characters long",
			lineID,
			maxLineIDLength)
	}

	if !validLineIDPattern.MatchString(lineID) {
		return fmt.Errorf("%q contains invalid characters", lineID)
	}

	return nil
}

// ConvertNameToID converts AIP-122 compliant line name
// to a line id.
func ConvertNameToID(lineName string) (string, error) {
	matches := lineNamePattern.FindStringSubmatch(lineName)

	if len(matches) != 2 {
		return "", fmt.Errorf("failed to convert name %q to line id", lineName)
	}

	return matches[1], nil
}

// ConvertIDToName converts the given line id to an
// AIP-122 compliant name.
func ConvertIDToName(lineID string) string {
	return lineNamePrefix + "/" + lineID
}
