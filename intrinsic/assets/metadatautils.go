// Copyright 2023 Intrinsic Innovation LLC

// Package metadatautils contains utilities for asset metadata.
package metadatautils

import (
	"slices"

	"intrinsic/assets/idutils"
	"intrinsic/assets/interfaceutils"
	"intrinsic/assets/tagutils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

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
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
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
var NameCharLength = map[atypepb.AssetType]int{
	atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT:    80,
	atypepb.AssetType_ASSET_TYPE_SERVICE:         80,
	atypepb.AssetType_ASSET_TYPE_SKILL:           45,
	atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE: 80,
	atypepb.AssetType_ASSET_TYPE_DATA:            80,
	atypepb.AssetType_ASSET_TYPE_PROCESS:         80,
}

// ManifestMetadata is an interface for an asset manifest proto used to specify metadata.
type ManifestMetadata interface {
	proto.Message

	GetDisplayName() string
	GetDocumentation() *documentationpb.Documentation
	GetId() *idpb.Id
	GetVendor() *vendorpb.Vendor
}

type metadataWithTag interface {
	GetAssetTag() atagpb.AssetTag
}

type metadataWithTags interface {
	GetAssetTags() []atagpb.AssetTag
}

type validateMetadataOptions struct {
	specifiesFileDescriptorSet *bool
	specifiesProvides          *bool
	specifiesReleaseNotes      *bool
	specifiesUpdateTime        *bool
	specifiesVersion           *bool
	requiredAssetType          atypepb.AssetType
}

// ValidateMetadataOption is a functional option for ValidateMetadata.
type ValidateMetadataOption func(*validateMetadataOptions)

// WithAssetType requires that the metadata has the given Asset type.
func WithAssetType(at atypepb.AssetType) ValidateMetadataOption {
	return func(opts *validateMetadataOptions) {
		opts.requiredAssetType = at
	}
}

// WithNoOutputOnlyFields requires that the metadata has no output-only fields.
func WithNoOutputOnlyFields() ValidateMetadataOption {
	return func(opts *validateMetadataOptions) {
		opts.specifiesFileDescriptorSet = proto.Bool(false)
		opts.specifiesProvides = proto.Bool(false)
	}
}

// WithInAssetOptions adds options for validating metadata that are represented within an Asset
// definition (e.g., see Data and Process Assets).
func WithInAssetOptions() ValidateMetadataOption {
	return func(opts *validateMetadataOptions) {
		WithNoOutputOnlyFields()(opts)
		opts.specifiesReleaseNotes = proto.Bool(false)
		opts.specifiesUpdateTime = proto.Bool(false)
		opts.specifiesVersion = proto.Bool(false)
	}
}

// WithInAssetOptionsIgnoreVersion adds options for validating metadata that are represented within
// an Asset definition (e.g., see Data and Process Assets).
//
// It ignores the presence of a version, to support transition from metadata that previously
// required a version.
func WithInAssetOptionsIgnoreVersion() ValidateMetadataOption {
	return func(opts *validateMetadataOptions) {
		WithNoOutputOnlyFields()(opts)
		opts.specifiesReleaseNotes = proto.Bool(false)
		opts.specifiesUpdateTime = proto.Bool(false)
	}
}

// WithCatalogOptions adds options for validating metadata for use in the AssetCatalog.
func WithCatalogOptions() ValidateMetadataOption {
	return func(opts *validateMetadataOptions) {
		opts.specifiesUpdateTime = proto.Bool(true)
		opts.specifiesVersion = proto.Bool(true)
	}
}

// ValidateMetadata validates asset metadata.
func ValidateMetadata(m *metadatapb.Metadata, options ...ValidateMetadataOption) error {
	opts := &validateMetadataOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if m.GetIdVersion().GetVersion() == "" {
		if err := idutils.ValidateIDProto(m.GetIdVersion().GetId()); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
		}
	} else if err := idutils.ValidateIDVersionProto(m.GetIdVersion()); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid id version: %v", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetIdVersion().GetId())

	if err := validateFieldPresence(m.GetFileDescriptorSet(), opts.specifiesFileDescriptorSet, "file descriptor set", id); err != nil {
		return err
	}
	if err := validateArrayPresence(m.GetProvides(), opts.specifiesProvides, "provides", id); err != nil {
		return err
	}
	if err := validateFieldPresence(m.GetReleaseNotes(), opts.specifiesReleaseNotes, "release notes", id); err != nil {
		return err
	}
	if err := validateFieldPresence(m.GetUpdateTime(), opts.specifiesUpdateTime, "update time", id); err != nil {
		return err
	}
	if err := validateFieldPresence(m.GetIdVersion().GetVersion(), opts.specifiesVersion, "version", id); err != nil {
		return err
	}

	if err := validateAssetType(m.GetAssetType(), opts.requiredAssetType, id); err != nil {
		return err
	}
	if err := validateName(m.GetIdVersion().GetId().GetName(), m.GetAssetType(), id); err != nil {
		return err
	}
	if err := validateDisplayName(m.GetDisplayName(), id); err != nil {
		return err
	}
	if err := validateAssetTag(m.GetAssetTag(), m.GetAssetType(), id); err != nil {
		return err
	}
	if err := validateVendor(m.GetVendor().GetDisplayName(), id); err != nil {
		return err
	}
	if err := validateVersion(m.GetIdVersion().GetVersion(), id); err != nil {
		return err
	}
	if err := validateDescription(m.GetDocumentation().GetDescription(), id); err != nil {
		return err
	}
	if err := validateRelNotes(m.GetReleaseNotes(), id); err != nil {
		return err
	}

	if sz := proto.Size(m); sz > MaxMetadataSize {
		return status.Errorf(codes.ResourceExhausted, "metadata size too large for %q: %d bytes > max %d bytes (Try reducing size of release notes and/or documentation.)", id, sz, MaxMetadataSize)
	}

	// Validate provides interfaces.
	for _, pi := range m.GetProvides() {
		if err := interfaceutils.ValidateInterfaceName(pi.GetUri()); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid provides interface %q for %q: %v", pi.GetUri(), id, err)
		}
	}

	return nil
}

// ValidateManifestMetadata validates asset manifest metadata.
func ValidateManifestMetadata(m ManifestMetadata) error {
	metadata, err := metadataFromManifestMetadata(m)
	if err != nil {
		return err
	}

	return ValidateMetadata(metadata)
}

// ToInputMetadata returns a clone of the input Metadata with output-only fields stripped.
func ToInputMetadata(m *metadatapb.Metadata) *metadatapb.Metadata {
	mOut := proto.Clone(m).(*metadatapb.Metadata)
	StripNonInputMetadata(mOut)

	return mOut
}

// StripNonInputMetadata strips the output-only fields from the specified metadata.
func StripNonInputMetadata(m *metadatapb.Metadata) {
	m.FileDescriptorSet = nil
	m.Provides = nil
}

// ToInAssetMetadata returns a clone of the input Metadata that is suitable for representation
// within an Asset definition (e.g., see Data and Process Assets).
func ToInAssetMetadata(m *metadatapb.Metadata) *metadatapb.Metadata {
	mOut := proto.Clone(m).(*metadatapb.Metadata)
	StripNonInAssetMetadata(mOut)

	return mOut
}

// StripNonInAssetMetadata strips from the input Metadata fields that don't apply when represented
// within an Asset definition.
func StripNonInAssetMetadata(m *metadatapb.Metadata) {
	StripNonInputMetadata(m)
	m.ReleaseNotes = ""
	m.UpdateTime = nil
	if m.GetIdVersion() != nil {
		m.GetIdVersion().Version = ""
	}
}

func validateFieldPresence[T comparable](x T, specifies *bool, name string, id string) error {
	var z T
	if specifies != nil {
		if *specifies {
			if x == z {
				return status.Errorf(codes.InvalidArgument, "required %s not specified for %q", name, id)
			}
		} else if x != z {
			return status.Errorf(codes.InvalidArgument, "disallowed %s specified for %q", name, id)
		}
	}
	return nil
}

func validateArrayPresence[T any](x []T, specifies *bool, name string, id string) error {
	if specifies != nil && !*specifies && len(x) > 0 {
		return status.Errorf(codes.InvalidArgument, "disallowed %s specified for %q", name, id)
	}
	return nil
}

func validateAssetType(at atypepb.AssetType, required atypepb.AssetType, id string) error {
	if at == atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
		return status.Errorf(codes.InvalidArgument, "no asset type specified for %q", id)
	}
	if required != atypepb.AssetType_ASSET_TYPE_UNSPECIFIED && at != required {
		return status.Errorf(codes.InvalidArgument, "invalid asset type for %q: required %v but got %v", id, required.String(), at.String())
	}

	return nil
}

func validateName(name string, at atypepb.AssetType, id string) error {
	maxLength, exists := NameCharLength[at]
	if !exists {
		return status.Errorf(codes.Internal, "unsupported asset type: %v", at)
	}
	if name == "" {
		return status.Errorf(codes.InvalidArgument, "no name specified for %q", id)
	}
	if nameLen := len(name); nameLen > maxLength {
		return status.Errorf(codes.ResourceExhausted, "name too long for %q: %q length %d > max %d", id, name, nameLen, maxLength)
	}
	return nil
}

func validateDisplayName(displayName string, id string) error {
	if displayName == "" {
		return status.Errorf(codes.InvalidArgument, "no display name specified for %q", id)
	}
	if nameLen := len(displayName); nameLen > DisplayNameCharLength {
		return status.Errorf(codes.ResourceExhausted, "display name too long for %q: %q length %d > max %d", id, displayName, nameLen, DisplayNameCharLength)
	}
	return nil
}

func validateAssetTag(atag atagpb.AssetTag, at atypepb.AssetType, id string) error {
	if atag == atagpb.AssetTag_ASSET_TAG_UNSPECIFIED {
		return nil
	}

	applicableTags, err := tagutils.AssetTagsForType(at)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot get asset tags for asset type %q: %v", at.String(), err)
	}
	if !slices.Contains(applicableTags, atag) {
		return status.Errorf(codes.InvalidArgument, "invalid asset tag for %q: tag %q is not applicable to asset type %q", id, atag.String(), at.String())
	}

	return nil
}

func validateVendor(vendor string, id string) error {
	if vendor == "" {
		return status.Errorf(codes.InvalidArgument, "no vendor specified for %q", id)
	}

	return nil
}

func validateVersion(version string, id string) error {
	if versionLen := len(version); versionLen > VersionCharLength {
		return status.Errorf(codes.ResourceExhausted, "version too long for %q: %q length %d > max %d", id, version, versionLen, VersionCharLength)
	}
	return nil
}

func validateDescription(description string, id string) error {
	if descriptionLen := len(description); descriptionLen > DescriptionCharLength {
		return status.Errorf(codes.ResourceExhausted, "description too long for %q: %q length %d > max %d", id, description, descriptionLen, DescriptionCharLength)
	}
	return nil
}

func validateRelNotes(relnotes string, id string) error {
	if relnotesLen := len(relnotes); relnotesLen > RelNotesCharLength {
		return status.Errorf(codes.ResourceExhausted, "release notes too long for %q: %q length %d > max %d", id, relnotes, relnotesLen, RelNotesCharLength)
	}
	return nil
}

func metadataFromManifestMetadata(m ManifestMetadata) (*metadatapb.Metadata, error) {
	var at atypepb.AssetType
	switch mt := m.(type) {
	case *skillmanifestpb.SkillManifest, *psmpb.SkillMetadata:
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
		return nil, status.Errorf(codes.Internal, "unsupported manifest type: %T", mt)
	}

	tag, err := tagFromManifestMetadata(m)
	if err != nil {
		return nil, err
	}

	return &metadatapb.Metadata{
		AssetTag:      tag,
		AssetType:     at,
		DisplayName:   m.GetDisplayName(),
		Documentation: m.GetDocumentation(),
		IdVersion: &idpb.IdVersion{
			Id: m.GetId(), // ID only, since manifests do not specify versions.
		},
		Vendor: m.GetVendor(),
	}, nil
}

func tagFromManifestMetadata(m any) (atagpb.AssetTag, error) {
	// Some metadata has a `tag` field.
	if mWithTag, ok := m.(metadataWithTag); ok {
		return mWithTag.GetAssetTag(), nil
	}

	// Some metadata has a `tags` field.
	if mWithTags, ok := m.(metadataWithTags); ok {
		if len(mWithTags.GetAssetTags()) > 0 {
			return mWithTags.GetAssetTags()[0], nil
		}

		return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, nil
	}

	// Some metadata cannot represent tags.
	if !metadataHasTag(m) {
		return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, nil
	}

	return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, status.Errorf(codes.Internal, "cannot get tag from unsupported manifest type for %T", m)
}

func metadataHasTag(m any) bool {
	switch m.(type) {
	case *datamanifestpb.DataManifest_Metadata, *skillmanifestpb.SkillManifest, *psmpb.SkillMetadata:
		return false
	}

	return true
}
