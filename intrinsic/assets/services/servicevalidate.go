// Copyright 2023 Intrinsic Innovation LLC

// Package servicevalidate provides utils for validating Services.
package servicevalidate

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	deputils "intrinsic/assets/dependencies/utils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	"intrinsic/util/go/validate"
	"intrinsic/util/proto/names"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	svpb "intrinsic/assets/services/proto/service_volume_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

var (
	missingServiceAllowlist = []string{
	}

	errContainsSkillAnnotations = errors.New("config message for the Service must not contain Skill-specific dependency annotations")
)

type serviceManifestOptions struct {
	files         *protoregistry.Files
	defaultConfig *anypb.Any
}

// ServiceManifestOption is an option for validating a ServiceManifest.
type ServiceManifestOption func(*serviceManifestOptions)

// WithDefaultConfig adds the Service's default config to the validation options.
//
// If specified, Validate verifies that the default config is of the type specified in the manifest.
func WithDefaultConfig(defaultConfig *anypb.Any) ServiceManifestOption {
	return func(opts *serviceManifestOptions) {
		opts.defaultConfig = defaultConfig
	}
}

// WithFiles provides a Files for validating proto messages.
func WithFiles(files *protoregistry.Files) ServiceManifestOption {
	return func(opts *serviceManifestOptions) {
		opts.files = files
	}
}

// ServiceManifest validates a ServiceManifest.
func ServiceManifest(m *smpb.ServiceManifest, options ...ServiceManifestOption) error {
	opts := &serviceManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if m == nil {
		return fmt.Errorf("ServiceManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid ServiceManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	expectedImagePaths, err := validateServiceDef(m.GetServiceDef(), opts.files)
	if err != nil {
		return fmt.Errorf("invalid service_def for Service %q: %w", id, err)
	}

	imagePaths := m.GetAssets().GetImageFilenames()
	for p := range expectedImagePaths {
		if !slices.Contains(imagePaths, p) {
			return fmt.Errorf("image %q in the manifest for Service %q is not listed in its assets", p, id)
		}
	}
	for _, p := range imagePaths {
		if _, ok := expectedImagePaths[p]; !ok {
			return fmt.Errorf("image %q in the assets for Service %q is not used in the manifest", p, id)
		}
	}

	if err := validateServiceConfig(
		m.GetServiceDef().GetConfigMessageFullName(),
		opts.defaultConfig,
		m.GetAssets().GetDefaultConfigurationFilename() != "",
		opts.files,
	); err != nil {
		return fmt.Errorf("invalid service config for Service %q: %w", id, err)
	}

	return nil
}

// ProcessedServiceManifest validates a ProcessedServiceManifest.
func ProcessedServiceManifest(m *smpb.ProcessedServiceManifest) error {
	if m == nil {
		return fmt.Errorf("ProcessedServiceManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid ProcessedServiceManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	var files *protoregistry.Files
	if fds := m.GetAssets().GetFileDescriptorSet(); fds != nil {
		var err error
		files, err = protodesc.NewFiles(fds)
		if err != nil {
			return fmt.Errorf("failed to populate the registry: %w", err)
		}
	}

	expectedImagePaths, err := validateServiceDef(m.GetServiceDef(), files)
	if err != nil {
		return fmt.Errorf("invalid service_def for Service %q: %w", id, err)
	}

	imagePaths := slices.Collect(maps.Keys(m.GetAssets().GetImages()))
	for p := range expectedImagePaths {
		if !slices.Contains(imagePaths, p) {
			return fmt.Errorf("image %q in the manifest for Service %q is not listed in its assets", p, id)
		}
	}
	for _, p := range imagePaths {
		if _, ok := expectedImagePaths[p]; !ok {
			return fmt.Errorf("image %q in the assets for Service %q is not used in the manifest", p, id)
		}
	}

	defaultConfig := m.GetAssets().GetDefaultConfiguration()
	if err := validateServiceConfig(m.GetServiceDef().GetConfigMessageFullName(), defaultConfig, defaultConfig != nil, files); err != nil {
		return fmt.Errorf("invalid service config for Service %q: %w", id, err)
	}

	return nil
}

// Volume validates a Volume.
func Volume(volume *svpb.Volume) error {
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

// VolumeMount validates a VolumeMount.
func VolumeMount(mount *svpb.VolumeMount) error {
	if err := validate.DNSLabel(mount.GetName()); err != nil {
		return fmt.Errorf("invalid volume name %q: %w", mount.GetName(), err)
	}
	if err := validate.UserString(mount.GetMountPath()); err != nil {
		return fmt.Errorf("invalid mount path %q: %w", mount.GetMountPath(), err)
	}
	return nil
}

func validateServiceDef(sd *smpb.ServiceDef, files *protoregistry.Files) (map[string]struct{}, error) {
	// Collect the Service's pod specs (verifying that at least a sim spec is specified).
	servicePodSpecs := map[string]*smpb.ServicePodSpec{}
	if sd != nil {
		if sd.GetRealSpec() != nil {
			servicePodSpecs["real"] = sd.GetRealSpec()
		}

		if sd.GetSimSpec() == nil {
			return nil, fmt.Errorf("a sim_spec must be specified if a service_def is provided")
		}
		servicePodSpecs["sim"] = sd.GetSimSpec()

		// Validate the Service's proto prefixes.
		for _, prefix := range sd.GetServiceProtoPrefixes() {
			if err := names.ValidateProtoPrefix(prefix); err != nil {
				return nil, fmt.Errorf("service proto prefix %q is not valid: %w", prefix, err)
			}

			strippedPrefix := strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
			inAllowlist := slices.Contains(missingServiceAllowlist, strippedPrefix)
			if files == nil {
				if !inAllowlist {
					return nil, fmt.Errorf("service proto prefix %q specified, but no descriptors provided", prefix)
				}
			} else if _, err := files.FindDescriptorByName(protoreflect.FullName(strippedPrefix)); err == protoregistry.NotFound {
				if !inAllowlist {
					return nil, fmt.Errorf("could not find service proto prefix %q in provided descriptors: %w", prefix, err)
				}
			} else if err != nil {
				return nil, fmt.Errorf("checking against the file descriptor set failed unexpectedly: %w", err)
			}
		}

		if sd.GetServiceInspectionConfig() != nil {
			config := sd.GetServiceInspectionConfig()
			if config.GetDataProtoMessageFullName() == "" {
				return nil, fmt.Errorf("inspection config is present but data_proto_message_full_name is empty")
			}
			// Validate the inspection proto message to be in the FileDescriptorSet.
			if files == nil {
				return nil, fmt.Errorf("inspection data proto message %q specified, but no descriptors provided", config.GetDataProtoMessageFullName())
			}
			if _, err := files.FindDescriptorByName(protoreflect.FullName(config.GetDataProtoMessageFullName())); err != nil {
				return nil, fmt.Errorf("could not find inspection data proto message %q in provided descriptors: %w", config.GetDataProtoMessageFullName(), err)
			}
		}

	}

	// Validate the Service's volumes.
	for podType, spec := range servicePodSpecs {
		if err := validateServicePodSpecVolumes(spec); err != nil {
			return nil, fmt.Errorf("invalid volumes in the %s spec: %w", podType, err)
		}
	}

	// Collect the Service's image paths for future validation.
	expectedImagePaths := map[string]struct{}{}
	for _, spec := range servicePodSpecs {
		if name := spec.GetImage().GetArchiveFilename(); name != "" {
			expectedImagePaths[name] = struct{}{}
		}
		for _, container := range spec.GetExtraImages() {
			expectedImagePaths[container.GetArchiveFilename()] = struct{}{}
		}
	}

	return expectedImagePaths, nil
}

func validateServicePodSpecVolumes(spec *smpb.ServicePodSpec) error {
	// Validate the defined volumes.
	volumeNames := map[string]struct{}{}
	for _, volume := range spec.GetSettings().GetVolumes() {
		if err := Volume(volume); err != nil {
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
		if err := VolumeMount(mount); err != nil {
			return err
		}

		if _, ok := volumeNames[mount.GetName()]; !ok {
			return fmt.Errorf("volume mount references non-existent volume %q", mount.GetName())
		}
	}

	return nil
}

func validateServiceConfig(configMessageFullName string, defaultConfig *anypb.Any, defaultConfigNeeded bool, files *protoregistry.Files) error {
	defaultConfigProvided := defaultConfig != nil
	if defaultConfigNeeded && !defaultConfigProvided {
		return fmt.Errorf("default config is needed but not provided")
	}
	if !defaultConfigNeeded && defaultConfigProvided {
		return fmt.Errorf("no default config is needed, but one was provided")
	}
	if configMessageFullName != "" {
		if defaultConfigNeeded {
			// Verify that the default config is of the specified config message type.
			if string(defaultConfig.MessageName()) != configMessageFullName {
				return fmt.Errorf("default config is of type %q, but manifest specifies config type %q", defaultConfig.MessageName(), configMessageFullName)
			}
		}
	} else if defaultConfigNeeded {
		defaultConfigMessageName := string(defaultConfig.MessageName())
		if defaultConfigMessageName == "" {
			return fmt.Errorf("default config cannot be an empty Any message; omit it instead")
		}
		configMessageFullName = defaultConfigMessageName
	}

	// Verify that the Service's config message, if any, is in the file descriptor set.
	if configMessageFullName != "" {
		if files == nil {
			return fmt.Errorf("config message specified (%q), but no descriptors provided", configMessageFullName)
		}
		d, err := files.FindDescriptorByName(protoreflect.FullName(configMessageFullName))
		if err != nil {
			return fmt.Errorf("could not find config message %q in provided descriptors: %w", configMessageFullName, err)
		}
		if md, ok := d.(protoreflect.MessageDescriptor); !ok {
			return fmt.Errorf("config message %q is not a message", configMessageFullName)
		} else if hasSkillAnnotations := deputils.HasResolvedDependency(md, deputils.WithSkillAnnotations()); hasSkillAnnotations {
			return errContainsSkillAnnotations
		}
	}

	return nil
}
