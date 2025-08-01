// Copyright 2023 Intrinsic Innovation LLC

// Package metadatautils contains utilities for asset metadata.
package metadatautils

import (
	"fmt"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"intrinsic/assets/idutils"
	"intrinsic/assets/tagutils"

	datamanifestpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	pmpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sceneobjectmanifestpb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	servicemanifestpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	skillmanifestpb "intrinsic/skills/proto/skill_manifest_go_proto"
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
	// MaxMetadataSize is the maximum size of a Metadata message.
	MaxMetadataSize = 750 * 1024 // 750 kiB
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
		atypepb.AssetType_ASSET_TYPE_PROCESS:         80,
	}
)

// ManifestMetadata is an interface for an asset manifest proto used to specify metadata.
type ManifestMetadata interface {
	GetDisplayName() string
	GetDocumentation() *documentationpb.Documentation
	GetId() *idpb.Id
	GetVendor() *vendorpb.Vendor
}

// ValidateMetadataOptions contains options for ValidateMetadata.
type ValidateMetadataOptions struct {
	requireUpdateTime bool
	requireVersion    bool
}

// ValidateMetadataOption is a functional option for ValidateMetadata.
type ValidateMetadataOption func(*ValidateMetadataOptions)

// WithRequireUpdateTime requires that the metadata has an update time.
func WithRequireUpdateTime(requireUpdateTime bool) ValidateMetadataOption {
	return func(opts *ValidateMetadataOptions) {
		opts.requireUpdateTime = requireUpdateTime
	}
}

// WithRequireVersion requires that the metadata has a version.
func WithRequireVersion(requireVersion bool) ValidateMetadataOption {
	return func(opts *ValidateMetadataOptions) {
		opts.requireVersion = requireVersion
	}
}

// WithCatalogOptions adds options for validating metadata for use in the catalog.
func WithCatalogOptions() []ValidateMetadataOption {
	return []ValidateMetadataOption{
		WithRequireUpdateTime(true),
		WithRequireVersion(true),
	}
}

// ValidateMetadata validates asset metadata.
func ValidateMetadata(m *metadatapb.Metadata, options ...ValidateMetadataOption) error {
	opts := &ValidateMetadataOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.requireVersion || m.GetIdVersion().GetVersion() != "" {
		if err := idutils.ValidateIDVersionProto(m.GetIdVersion()); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id version: %v", err)
		}
	} else if err := idutils.ValidateIDProto(m.GetIdVersion().GetId()); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetIdVersion().GetId())

	if m.GetDisplayName() == "" {
		return status.Errorf(codes.InvalidArgument, "no display name specified for %q", id)
	}
	if m.GetVendor().GetDisplayName() == "" {
		return status.Errorf(codes.InvalidArgument, "no vendor specified for %q", id)
	}
	if opts.requireUpdateTime && m.GetUpdateTime() == nil {
		return status.Errorf(codes.InvalidArgument, "no update time specified for %q", id)
	}
	if m.GetAssetType() == atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
		return status.Errorf(codes.InvalidArgument, "no asset type specified for %q", id)
	}
	if m.GetAssetTag() != atagpb.AssetTag_ASSET_TAG_UNSPECIFIED {
		applicableTags, err := tagutils.AssetTagsForType(m.GetAssetType())
		if err != nil {
			return status.Errorf(
				codes.Internal,
				"cannot get asset tags for asset type %q: %v",
				atypepb.AssetType_name[int32(m.GetAssetType())],
				err)
		}
		if !slices.Contains(applicableTags, m.GetAssetTag()) {
			return status.Errorf(
				codes.InvalidArgument,
				"asset tag %q is not applicable to %q asset %q",
				atagpb.AssetTag_name[int32(m.GetAssetTag())],
				atypepb.AssetType_name[int32(m.GetAssetType())],
				id)
		}
	}

	// Validate metadata size limits.
	if err := ValidateNameLength(m.GetIdVersion().GetId().GetName(), m.GetAssetType()); err != nil {
		return status.Errorf(codes.ResourceExhausted, "invalid name for %q: %v", id, err)
	}
	if err := ValidateDisplayNameLength(m.GetDisplayName()); err != nil {
		return status.Errorf(codes.ResourceExhausted, "invalid display name for %q: %v", id, err)
	}
	if err := ValidateVersionLength(m.GetIdVersion().GetVersion()); err != nil {
		return status.Errorf(codes.ResourceExhausted, "invalid version for %q: %v", id, err)
	}
	if err := ValidateDescriptionLength(m.GetDocumentation().GetDescription()); err != nil {
		return status.Errorf(codes.ResourceExhausted, "invalid description for %q: %v", id, err)
	}
	if err := ValidateRelNotesLength(m.GetReleaseNotes()); err != nil {
		return status.Errorf(codes.ResourceExhausted, "invalid release notes for %q: %v", id, err)
	}
	if sz := proto.Size(m); sz > MaxMetadataSize {
		return status.Errorf(codes.ResourceExhausted, "metadata size of %q is too large: %d bytes > max %d bytes (Try reducing size of release notes and/or documentation.)", id, sz, MaxMetadataSize)
	}

	return nil
}

// ValidateManifestMetadata validates asset manifest metadata.
func ValidateManifestMetadata(m ManifestMetadata) error {
	var at atypepb.AssetType
	switch m := m.(type) {
	case *skillmanifestpb.SkillManifest:
		at = atypepb.AssetType_ASSET_TYPE_SKILL
	case *servicemanifestpb.ServiceMetadata:
		at = atypepb.AssetType_ASSET_TYPE_SERVICE
	case *sceneobjectmanifestpb.SceneObjectMetadata:
		at = atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT
	case *datamanifestpb.DataManifest_Metadata:
		at = atypepb.AssetType_ASSET_TYPE_DATA
	case *hdmpb.HardwareDeviceMetadata:
		at = atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE
	case *pmpb.ProcessMetadata:
		at = atypepb.AssetType_ASSET_TYPE_PROCESS
	default:
		return fmt.Errorf("unsupported manifest type: %T", m)
	}

	// Some metadata protos have an asset tag field, some don't.
	type metadataWithTag interface{ GetAssetTag() atagpb.AssetTag }
	tag := atagpb.AssetTag_ASSET_TAG_UNSPECIFIED
	if mWithTag, ok := m.(metadataWithTag); ok {
		tag = mWithTag.GetAssetTag()
	}

	return ValidateMetadata(&metadatapb.Metadata{
		AssetType:     at,
		DisplayName:   m.GetDisplayName(),
		Documentation: m.GetDocumentation(),
		IdVersion: &idpb.IdVersion{
			Id: m.GetId(), // ID only, since manifests do not specify versions.
		},
		Vendor:   m.GetVendor(),
		AssetTag: tag,
	})
}

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
