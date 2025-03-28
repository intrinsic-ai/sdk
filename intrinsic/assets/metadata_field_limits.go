// Copyright 2023 Intrinsic Innovation LLC

// Package metadatafieldlimits defines size restrictions on various fields for asset metadata.
package metadatafieldlimits

import (
	"fmt"

	atypepb "intrinsic/assets/proto/asset_type_go_proto"
)

const (
	// DisplayNameCharLength sets character limits for the asset display name.
	DisplayNameCharLength = 80
	// VersionCharLength sets character limits for asset versions.
	// Semver does not have restrictions on character limits for versions, and leaves it to best
	// judgement: https://semver.org/#does-semver-have-a-size-limit-on-the-version-string. An
	// arbitrary limit is set considering potential lengths of prerelease label and build metadata.
	VersionCharLength = 128
	// DescriptionCharLength sets character limits for asset descriptions.
	// The character limits for description were set based on the longest existing asset description
	// at the time of writing.
	DescriptionCharLength = 2400
	// RelNotesCharLength sets character limits for release notes.
	// The character limits for release notes were set arbitrarily based on the length of
	// description.
	RelNotesCharLength = 2400
)

var (
	// NameCharLength sets character limits for asset names as a function of their type.
	// Skill IDs are used as the DNS names in k8s, and are formatted as "<package_name>.<name>".
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#rfc-1035-label-names
	// restricts label names to 63 characters. Keeping this as the maximum allowed upper limit, a
	// conservative limit is chosen for skill names to account for package names, other prefixes and
	// special symbols that maybe be added when converting IDs to k8s labels.
	//
	// Note: Skill IDs are used for DNS names but asset names are what is checked here for character
	// limits. We cannot exactly limit the ID character limit, because they could have nested
	// packages. And for the very same reason, we cannot restrict character limits for packages.
	NameCharLength = map[atypepb.AssetType]int{
		atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT:    80,
		atypepb.AssetType_ASSET_TYPE_SERVICE:         80,
		atypepb.AssetType_ASSET_TYPE_SKILL:           45,
		atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE: 80,
		atypepb.AssetType_ASSET_TYPE_DATA:            80,
	}
)

// ValidateNameLength validates the length of asset names.
func ValidateNameLength(name string, at atypepb.AssetType) error {
	maxLength, exists := NameCharLength[at]
	if !exists {
		return fmt.Errorf("unsupported asset type: %v", at)
	}
	if nameLen := len(name); nameLen > maxLength {
		return fmt.Errorf("name %q exceeds character limit: length %d > max %d", name, nameLen, maxLength)
	}
	return nil
}

// ValidateDisplayNameLength validates the length of asset display names.
func ValidateDisplayNameLength(name string) error {
	if nameLen := len(name); nameLen > DisplayNameCharLength {
		return fmt.Errorf("display name %q exceeds character limit: length %d > max %d", name, nameLen, DisplayNameCharLength)
	}
	return nil
}

// ValidateVersionLength validates the length of asset versions.
func ValidateVersionLength(version string) error {
	if versionLen := len(version); versionLen > VersionCharLength {
		return fmt.Errorf("version exceeds character limit: length %d > max %d", versionLen, VersionCharLength)
	}
	return nil
}

// ValidateDescriptionLength validates the length of asset names.
func ValidateDescriptionLength(description string) error {
	if descriptionLen := len(description); descriptionLen > DescriptionCharLength {
		return fmt.Errorf("description exceeds character limit: length %d > max %d", descriptionLen, DescriptionCharLength)
	}
	return nil
}

// ValidateRelNotesLength validates the length of asset names.
func ValidateRelNotesLength(relnotes string) error {
	if relnotesLen := len(relnotes); relnotesLen > RelNotesCharLength {
		return fmt.Errorf("release notes exceeds character limit: length %d > max %d", relnotesLen, RelNotesCharLength)
	}
	return nil
}
