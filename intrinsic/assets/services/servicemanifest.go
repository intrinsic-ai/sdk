// Copyright 2023 Intrinsic Innovation LLC

// Package servicemanifest contains tools for working with ServiceManifest.
package servicemanifest

import (
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	"intrinsic/util/go/validate"
	"intrinsic/util/proto/names"

	anypb "google.golang.org/protobuf/types/known/anypb"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	svpb "intrinsic/assets/services/proto/service_volume_go_proto"
)

var (
	missingServiceAllowlist = []string{
	}

)

// ValidateServiceManifestOptions contains options for validating a ServiceManifest.
type ValidateServiceManifestOptions struct {
	files         *protoregistry.Files
	defaultConfig *anypb.Any
}

// ValidateServiceManifestOption is an option for validating a ServiceManifest.
type ValidateServiceManifestOption func(*ValidateServiceManifestOptions)

// WithFiles adds the proto files to the validation options.
func WithFiles(files *protoregistry.Files) ValidateServiceManifestOption {
	return func(opts *ValidateServiceManifestOptions) {
		opts.files = files
	}
}

// WithDefaultConfig adds the Service's default config to the validation options.
//
// If specified, ValidateServiceManifest verifies that the default config is of the type specified
// in the manifest.
func WithDefaultConfig(defaultConfig *anypb.Any) ValidateServiceManifestOption {
	return func(opts *ValidateServiceManifestOptions) {
		opts.defaultConfig = defaultConfig
	}
}

// ValidateServiceManifest verifies that a ServiceManifest is consistent and valid.
func ValidateServiceManifest(m *smpb.ServiceManifest, options ...ValidateServiceManifestOption) error {
	opts := &ValidateServiceManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	// Collect the Service's pod specs (verifying that at least a sim spec is specified).
	servicePodSpecs := map[string]*smpb.ServicePodSpec{}
	if m.GetServiceDef() != nil {
		if m.GetServiceDef().GetRealSpec() != nil {
			servicePodSpecs["real"] = m.GetServiceDef().GetRealSpec()
		}

		if m.GetServiceDef().GetSimSpec() == nil {
			return fmt.Errorf("a sim_spec must be specified if a service_def is provided for Service %q", id)
		}
		servicePodSpecs["sim"] = m.GetServiceDef().GetSimSpec()
	}

	// Validate the Service's images.
	expectedImagePaths := map[string]struct{}{}
	for _, spec := range servicePodSpecs {
		if name := spec.GetImage().GetArchiveFilename(); name != "" {
			expectedImagePaths[name] = struct{}{}
		}
		for _, container := range spec.GetExtraImages() {
			expectedImagePaths[container.GetArchiveFilename()] = struct{}{}
		}
	}

	for p := range expectedImagePaths {
		if !slices.Contains(m.GetAssets().GetImageFilenames(), p) {
			return fmt.Errorf("image %q in the manifest for Service %q is not listed in the service assets", p, id)
		}
	}
	for _, p := range m.GetAssets().GetImageFilenames() {
		if _, ok := expectedImagePaths[p]; !ok {
			return fmt.Errorf("image %q in the service assets for Service %q is not used in the manifest", p, id)
		}
	}

	// Validate the Service's volumes.
	for podType, spec := range servicePodSpecs {
		if err := validateServicePodSpecVolumes(spec); err != nil {
			return fmt.Errorf("invalid volumes in the %s spec for Service %q: %w", podType, id, err)
		}
	}

	// Verify that the Service's config message is in the file descriptor set.
	if m.GetServiceDef().GetConfigMessageFullName() != "" {
		if opts.files == nil {
			return fmt.Errorf("config message %q specified, but no descriptors provided", m.GetServiceDef().GetConfigMessageFullName())
		}
		if _, err := opts.files.FindDescriptorByName(protoreflect.FullName(m.GetServiceDef().GetConfigMessageFullName())); err != nil {
			return fmt.Errorf("could not find config message %q in provided descriptors for Service %q: %w", m.GetServiceDef().GetConfigMessageFullName(), id, err)
		}

		// If a default config is specified, verify that it is of the specified config message type.
		if opts.defaultConfig != nil {
			if string(opts.defaultConfig.MessageName()) != m.GetServiceDef().GetConfigMessageFullName() {
				return fmt.Errorf("default config for Service %q is of type %q, but manifest specifies config type %q", id, opts.defaultConfig.MessageName(), m.GetServiceDef().GetConfigMessageFullName())
			}
		}
	}

	// Validate the Service's proto prefixes.
	for _, prefix := range m.GetServiceDef().GetServiceProtoPrefixes() {
		if err := names.ValidateProtoPrefix(prefix); err != nil {
			return fmt.Errorf("service proto prefix %q is not valid for Service %q: %w", prefix, id, err)
		}

		strippedPrefix := strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
		name := protoreflect.FullName(strippedPrefix)
		inAllowlist := slices.Contains(missingServiceAllowlist, strippedPrefix)
		if opts.files == nil {
			if !inAllowlist {
				return fmt.Errorf("service proto prefix %q specified, but no descriptors provided", prefix)
			}
		} else if _, err := opts.files.FindDescriptorByName(name); err == protoregistry.NotFound {
			if !inAllowlist {
				return fmt.Errorf("could not find service proto prefix %q in provided descriptors for Service %q: %w", prefix, id, err)
			}
		} else if err != nil {
			return fmt.Errorf("checking against the file descriptor set failed unexpectedly: %w", err)
		}
	}
	return nil
}

// ValidateVolume verifies that a Volume is valid.
func ValidateVolume(volume *svpb.Volume) error {
	if err := validate.DNSLabel(volume.GetName()); err != nil {
		return fmt.Errorf("invalid volume name %q: %w", volume.GetName(), err)
	}

	switch volume.GetSource().(type) {
	case *svpb.Volume_HostPath:
		if err := validate.UserString(volume.GetHostPath().GetPath()); err != nil {
			return fmt.Errorf("invalid host path %q: %w", volume.GetHostPath().GetPath(), err)
		}
	case *svpb.Volume_EmptyDir:
	case nil:
		return fmt.Errorf("volume %q did not specify a source", volume.GetName())
	default:
		return fmt.Errorf("unsupported volume source type %T for volume %q", volume.GetSource(), volume.GetName())
	}

	return nil
}

// ValidateVolumeMount verifies that a VolumeMount is valid.
func ValidateVolumeMount(mount *svpb.VolumeMount) error {
	if err := validate.DNSLabel(mount.GetName()); err != nil {
		return fmt.Errorf("invalid volume name %q: %w", mount.GetName(), err)
	}
	if err := validate.UserString(mount.GetMountPath()); err != nil {
		return fmt.Errorf("invalid mount path %q: %w", mount.GetMountPath(), err)
	}
	return nil
}

func validateServicePodSpecVolumes(spec *smpb.ServicePodSpec) error {
	// Validate the defined volumes.
	volumeNames := map[string]struct{}{}
	for _, volume := range spec.GetSettings().GetVolumes() {
		if err := ValidateVolume(volume); err != nil {
			return err
		}

		name := volume.GetName()
		if _, ok := volumeNames[name]; ok {
			return fmt.Errorf("volume %q is specified multiple times", name)
		}
		volumeNames[name] = struct{}{}
	}

	// Validate the volume mounts.
	for _, mount := range spec.GetImage().GetSettings().GetVolumeMounts() {
		if err := ValidateVolumeMount(mount); err != nil {
			return err
		}

		if _, ok := volumeNames[mount.GetName()]; !ok {
			return fmt.Errorf("volume mount references non-existent volume %q", mount.GetName())
		}
	}

	return nil
}

func setDifference(slice1, slice2 []string) []string {
	var difference []string
	for _, val := range slice1 {
		if !slices.Contains(slice2, val) {
			difference = append(difference, val)
		}
	}
	return difference
}
