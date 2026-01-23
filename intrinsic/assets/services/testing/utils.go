// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides testing utils for Services.
package utils

import (
	"testing"

	idpb "intrinsic/assets/proto/id_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	sipb "intrinsic/assets/services/proto/v1/service_inspection_go_proto" 
	imagepb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/types/known/anypb"
)

type makeServiceManifestOptions struct {
	configMessageFullName        string
	defaultConfigurationFilename *string
	imageFilenames               []string
	metadata                     *smpb.ServiceMetadata
	realSpec                     *smpb.ServicePodSpec
	serviceInspectionConfig      *sipb.ServiceInspectionConfig 
	serviceProtoPrefixes         []string
	simSpec                      *smpb.ServicePodSpec
}

// MakeServiceManifestOption is an option for MakeServiceManifest.
type MakeServiceManifestOption func(*makeServiceManifestOptions)

// WithConfigMessageFullName specifies the config message name to use in the ServiceManifest.
func WithConfigMessageFullName(name string) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.configMessageFullName = name
	}
}

// WithDefaultConfigurationFilename specifies the default config file name to use in the
// ServiceManifest.
func WithDefaultConfigurationFilename(filename string) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.defaultConfigurationFilename = &filename
	}
}

// WithImageFilenames specifies the image files names to use in the ServiceManifest.
func WithImageFilenames(filenames []string) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.imageFilenames = filenames
	}
}

// WithMetadata specifies the metadata to use in the ServiceManifest.
func WithMetadata(metadata *smpb.ServiceMetadata) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.metadata = metadata
	}
}

// WithRealSpec specifies the real spec to use in the ServiceManifest.
func WithRealSpec(spec *smpb.ServicePodSpec) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.realSpec = spec
	}
}


// WithServiceInspectionConfig specifies the service inspection config to use in the
// ServiceManifest.
func WithServiceInspectionConfig(config *sipb.ServiceInspectionConfig) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.serviceInspectionConfig = config
	}
}



// WithServiceProtoPrefixes specifies the service proto prefixes to use in the ServiceManifest.
func WithServiceProtoPrefixes(prefixes []string) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.serviceProtoPrefixes = prefixes
	}
}

// WithSimSpec specifies the sim spec to use in the ServiceManifest.
func WithSimSpec(spec *smpb.ServicePodSpec) MakeServiceManifestOption {
	return func(opts *makeServiceManifestOptions) {
		opts.simSpec = spec
	}
}

// MakeServiceManifest makes a ServiceManifest for testing.
func MakeServiceManifest(t *testing.T, options ...MakeServiceManifestOption) *smpb.ServiceManifest {
	t.Helper()

	opts := &makeServiceManifestOptions{
		metadata: &smpb.ServiceMetadata{
			Id: &idpb.Id{
				Name:    "some_service",
				Package: "package.some",
			},
			DisplayName: "Some Service",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}
	if opts.imageFilenames == nil && opts.realSpec == nil && opts.simSpec == nil {
		opts.imageFilenames = []string{"sim_image", "real_image", "extra_image"}
		opts.realSpec = &smpb.ServicePodSpec{
			Image: &smpb.ServiceImage{
				ArchiveFilename: "real_image",
			},
			ExtraImages: []*smpb.ServiceImage{
				{
					ArchiveFilename: "extra_image",
				},
			},
		}
		opts.simSpec = &smpb.ServicePodSpec{
			Image: &smpb.ServiceImage{
				ArchiveFilename: "sim_image",
			},
			ExtraImages: []*smpb.ServiceImage{
				{
					ArchiveFilename: "extra_image",
				},
			},
		}
	}

	return &smpb.ServiceManifest{
		Metadata: opts.metadata,
		ServiceDef: &smpb.ServiceDef{
			ConfigMessageFullName:   opts.configMessageFullName,
			ServiceInspectionConfig: opts.serviceInspectionConfig, 
			ServiceProtoPrefixes:    opts.serviceProtoPrefixes,
			RealSpec:                opts.realSpec,
			SimSpec:                 opts.simSpec,
		},
		Assets: &smpb.ServiceAssets{
			DefaultConfigurationFilename: opts.defaultConfigurationFilename,
			ImageFilenames:               opts.imageFilenames,
		},
	}
}

type makeProcessedServiceManifestOptions struct {
	configMessageFullName   string
	defaultConfiguration    *anypb.Any
	fileDescriptorSet       *dpb.FileDescriptorSet
	images                  map[string]*imagepb.Image
	metadata                *smpb.ServiceMetadata
	realSpec                *smpb.ServicePodSpec
	serviceInspectionConfig *sipb.ServiceInspectionConfig 
	serviceProtoPrefixes    []string
	simSpec                 *smpb.ServicePodSpec
}

// MakeProcessedServiceManifestOption is an option for MakeProcessedServiceManifest.
type MakeProcessedServiceManifestOption func(*makeProcessedServiceManifestOptions)

// WithProcessedConfigMessageFullName specifies the config message name to use in the
// ProcessedServiceManifest.
func WithProcessedConfigMessageFullName(name string) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.configMessageFullName = name
	}
}

// WithDefaultConfiguration specifies the default configuration to use in the
// ProcessedServiceManifest.
func WithDefaultConfiguration(config *anypb.Any) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.defaultConfiguration = config
	}
}

// WithFileDescriptorSet specifies the file descriptor set to use in the ProcessedServiceManifest.
func WithFileDescriptorSet(fds *dpb.FileDescriptorSet) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.fileDescriptorSet = fds
	}
}

// WithImages specifies the images to use in the ProcessedServiceManifest.
func WithImages(images map[string]*imagepb.Image) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.images = images
	}
}

// WithProcessedMetadata specifies the metadata to use in the ProcessedServiceManifest.
func WithProcessedMetadata(metadata *smpb.ServiceMetadata) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.metadata = metadata
	}
}

// WithProcessedRealSpec specifies the real spec to use in the ProcessedServiceManifest.
func WithProcessedRealSpec(spec *smpb.ServicePodSpec) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.realSpec = spec
	}
}


// WithProcessedServiceInspectionConfig specifies the service inspection config to use in the
// ProcessedServiceManifest.
func WithProcessedServiceInspectionConfig(config *sipb.ServiceInspectionConfig) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.serviceInspectionConfig = config
	}
}



// WithProcessedServiceProtoPrefixes specifies the service proto prefixes to use in the
// ProcessedServiceManifest.
func WithProcessedServiceProtoPrefixes(prefixes []string) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.serviceProtoPrefixes = prefixes
	}
}

// WithProcessedSimSpec specifies the sim spec to use in the ProcessedServiceManifest.
func WithProcessedSimSpec(spec *smpb.ServicePodSpec) MakeProcessedServiceManifestOption {
	return func(opts *makeProcessedServiceManifestOptions) {
		opts.simSpec = spec
	}
}

// MakeProcessedServiceManifest makes a ProcessedServiceManifest for testing.
func MakeProcessedServiceManifest(t *testing.T, options ...MakeProcessedServiceManifestOption) *smpb.ProcessedServiceManifest {
	t.Helper()

	opts := &makeProcessedServiceManifestOptions{
		metadata: &smpb.ServiceMetadata{
			Id: &idpb.Id{
				Name:    "some_service",
				Package: "package.some",
			},
			DisplayName: "Some Service",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}
	if opts.images == nil && opts.realSpec == nil && opts.simSpec == nil {
		opts.images = map[string]*imagepb.Image{
			"extra_image": {
				Registry: "gcr.io/test-project",
				Name:     "extra_image",
				Tag:      ":real",
			},
			"real_image": {
				Registry: "gcr.io/test-project",
				Name:     "real_image",
				Tag:      ":real",
			},
			"sim_image": {
				Registry: "gcr.io/test-project",
				Name:     "sim_image",
				Tag:      ":real",
			},
		}
		opts.realSpec = &smpb.ServicePodSpec{
			Image: &smpb.ServiceImage{
				ArchiveFilename: "real_image",
			},
			ExtraImages: []*smpb.ServiceImage{
				{
					ArchiveFilename: "extra_image",
				},
			},
		}
		opts.simSpec = &smpb.ServicePodSpec{
			Image: &smpb.ServiceImage{
				ArchiveFilename: "sim_image",
			},
			ExtraImages: []*smpb.ServiceImage{
				{
					ArchiveFilename: "extra_image",
				},
			},
		}
	}

	return &smpb.ProcessedServiceManifest{
		Metadata: opts.metadata,
		ServiceDef: &smpb.ServiceDef{
			ConfigMessageFullName:   opts.configMessageFullName,
			ServiceInspectionConfig: opts.serviceInspectionConfig, 
			ServiceProtoPrefixes:    opts.serviceProtoPrefixes,
			RealSpec:                opts.realSpec,
			SimSpec:                 opts.simSpec,
		},
		Assets: &smpb.ProcessedServiceAssets{
			DefaultConfiguration: opts.defaultConfiguration,
			FileDescriptorSet:    opts.fileDescriptorSet,
			Images:               opts.images,
		},
	}
}
