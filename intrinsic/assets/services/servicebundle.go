// Copyright 2023 Intrinsic Innovation LLC

// Package databundle provides utils for working with Service bundles.
package servicebundle

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"intrinsic/assets/imageutils"
	"intrinsic/assets/ioutils"
	"intrinsic/assets/services/servicemanifest"
	"intrinsic/util/archive/tartooling"
	"intrinsic/util/proto/registryutil"

	"github.com/google/safearchive/tar"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	serviceManifestPathInTar = "service_manifest.binarypb"
)

// ServiceBundle represents a Service Asset bundle.
type ServiceBundle struct {
	Manifest *smpb.ServiceManifest
	Files    map[string][]byte
}

type readServiceBundleOptions struct {
	readFiles bool
}

// ReadServiceBundleOption is a function option for ReadServiceBundle.
type ReadServiceBundleOption func(*readServiceBundleOptions)

// WithReadFiles specifies whether to read additional files when reading the bundle.
func WithReadFiles(b bool) ReadServiceBundleOption {
	return func(opts *readServiceBundleOptions) {
		opts.readFiles = b
	}
}

// ReadServiceBundle reads a Service Asset bundle from disk (see WriteServiceBundle).
func ReadServiceBundle(ctx context.Context, path string, options ...ReadServiceBundleOption) (*ServiceBundle, error) {
	opts := &readServiceBundleOptions{}
	for _, opt := range options {
		opt(opts)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	m, handlers := makeOnlyServiceManifestHandlers()
	walkTarOpts := []ioutils.WalkTarFileOption{
		ioutils.WithHandlers(handlers),
	}

	var inlined map[string][]byte
	if opts.readFiles {
		var fallback ioutils.WalkTarFileFallbackHandler
		inlined, fallback = ioutils.MakeCollectInlinedFallbackHandler()
		walkTarOpts = append(walkTarOpts, ioutils.WithFallbackHandler(fallback))
	}

	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f), walkTarOpts...); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}

	return &ServiceBundle{
		Manifest: m,
		Files:    inlined,
	}, nil
}

// ProcessServiceOpts contains the necessary handlers to generate a processed
// service manifest.
type ProcessServiceOpts struct {
	imageutils.ImageProcessor
}

// ProcessService creates a processed manifest from a bundle on disk using the
// provided processing functions.  It avoids doing any validation except for
// that required to transform the specified files in the bundle into their
// processed variants.
func ProcessService(ctx context.Context, path string, opts ProcessServiceOpts) (*smpb.ProcessedServiceManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	// Read the manifest and then reset the file once we have the information
	// about the bundle we're going to process.
	manifest, handlers := makeOnlyServiceManifestHandlers()
	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f), ioutils.WithHandlers(handlers)); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("could not seek in %q: %v", path, err)
	}

	// Initialize handlers for when we walk through the file again now that we
	// know what we're looking for, but error on unexpected files this time.
	processedAssets, handlers := makeServiceAssetHandlers(manifest, opts)
	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f),
		ioutils.WithHandlers(handlers),
		ioutils.WithFallbackHandler(ioutils.AlwaysErrorAsUnexpected),
	); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}

	m := &smpb.ProcessedServiceManifest{
		Metadata:   manifest.GetMetadata(),
		ServiceDef: manifest.GetServiceDef(),
		Assets:     processedAssets,
	}

	if m.GetServiceDef().GetConfigMessageFullName() != "" {
		// Generate an empty default config if none was provided.
		if m.GetAssets().GetDefaultConfiguration() == nil {
			types, err := registryutil.NewTypesFromFileDescriptorSet(m.GetAssets().GetFileDescriptorSet())
			if err != nil {
				return nil, fmt.Errorf("failed to populate the registry: %v", err)
			}
			msgType, err := types.FindMessageByName(protoreflect.FullName(m.GetServiceDef().GetConfigMessageFullName()))
			if err != nil {
				return nil, fmt.Errorf("failed to find config message %q: %v", m.GetServiceDef().GetConfigMessageFullName(), err)
			}
			defaultConfig, err := anypb.New(msgType.New().Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to create default config: %v", err)
			}
			m.GetAssets().DefaultConfiguration = defaultConfig
		}
	} else if m.GetAssets().GetDefaultConfiguration() != nil {
		// Derive config message name from the default config, if specified.
		m.GetServiceDef().ConfigMessageFullName = string(m.GetAssets().GetDefaultConfiguration().MessageName())
	}

	return m, nil
}

// WriteServiceBundleOptions provides options for writing a Service Asset bundle.
type WriteServiceBundleOptions struct {
	DefaultConfig *anypb.Any
	Descriptors   *descriptorpb.FileDescriptorSet
	ImageTarPaths []string
}

// WriteServiceBundle writes a Service .tar bundle file at the specified path.
func WriteServiceBundle(m *smpb.ServiceManifest, path string, opts *WriteServiceBundleOptions) error {
	if m == nil {
		return fmt.Errorf("ServiceManifest must not be nil")
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tw := tar.NewWriter(out)

	m.Assets = new(smpb.ServiceAssets)
	if opts.Descriptors != nil {
		descriptorName := "descriptors-transitive-descriptor-set.proto.bin"
		m.Assets.ParameterDescriptorFilename = &descriptorName
		if err := tartooling.AddBinaryProto(opts.Descriptors, tw, descriptorName); err != nil {
			return fmt.Errorf("failed to write FileDescriptorSet to bundle: %w", err)
		}
	}
	if opts.DefaultConfig != nil {
		configName := "default_config.binarypb"
		m.Assets.DefaultConfigurationFilename = &configName
		if err := tartooling.AddBinaryProto(opts.DefaultConfig, tw, configName); err != nil {
			return fmt.Errorf("failed to write default config to bundle: %w", err)
		}
	}
	for _, path := range opts.ImageTarPaths {
		base := filepath.Base(path)
		m.Assets.ImageFilenames = append(m.Assets.ImageFilenames, base)
		if err := tartooling.AddFile(path, tw, base); err != nil {
			return fmt.Errorf("failed to write %q to bundle: %w", path, err)
		}
	}

	var files *protoregistry.Files
	if opts.Descriptors != nil {
		files, err = protodesc.NewFiles(opts.Descriptors)
		if err != nil {
			return fmt.Errorf("failed to create proto files: %w", err)
		}
	}
	if err := servicemanifest.ValidateServiceManifest(m,
		servicemanifest.WithFiles(files),
		servicemanifest.WithDefaultConfig(opts.DefaultConfig),
	); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	// Now we can write the manifest, since assets have been completed.
	if err := tartooling.AddBinaryProto(m, tw, "service_manifest.binarypb"); err != nil {
		return fmt.Errorf("failed to write ServiceManifest to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}

// makeOnlyServiceManifestHandlers returns a map of handlers that only pull out
// the service manifest from the tar file into the returned proto.  Can be used
// with a fallback handler.
func makeOnlyServiceManifestHandlers() (*smpb.ServiceManifest, map[string]ioutils.WalkTarFileHandler) {
	manifest := new(smpb.ServiceManifest)
	handlers := map[string]ioutils.WalkTarFileHandler{
		serviceManifestPathInTar: ioutils.MakeBinaryProtoHandler(manifest),
	}
	return manifest, handlers
}

// makeServiceAssetHandlers returns handlers for all assets listed in the
// service manifest.  This will be at most:
// * An handler that ignores the manifest
// * A binary proto handler for the default configuration file
// * A binary proto handler for the file descriptor set file
// * A handler that wraps opts.ImageProcessor to be called on every image
func makeServiceAssetHandlers(manifest *smpb.ServiceManifest, opts ProcessServiceOpts) (*smpb.ProcessedServiceAssets, map[string]ioutils.WalkTarFileHandler) {
	handlers := map[string]ioutils.WalkTarFileHandler{
		serviceManifestPathInTar: ioutils.IgnoreHandler, // already read this.
	}
	// Don't generate an empty assets message if there wasn't one to begin
	// with.  This is a slightly odd state, but Process is not doing validation
	// of the manifest.  This also protects against nil access of
	// manifest.GetAssets().{MemberVariable}, which is required for checking
	// the "optional" piece of "optional string" fields in this version of the
	// golang proto API.
	if manifest.GetAssets() == nil {
		return nil, handlers
	}

	processedAssets := new(smpb.ProcessedServiceAssets)
	if p := manifest.GetAssets().DefaultConfigurationFilename; p != nil {
		processedAssets.DefaultConfiguration = new(anypb.Any)
		handlers[*p] = ioutils.MakeBinaryProtoHandler(processedAssets.DefaultConfiguration)
	}
	if p := manifest.GetAssets().ParameterDescriptorFilename; p != nil {
		processedAssets.FileDescriptorSet = new(descriptorpb.FileDescriptorSet)
		handlers[*p] = ioutils.MakeBinaryProtoHandler(processedAssets.FileDescriptorSet)
	}
	for _, p := range manifest.GetAssets().GetImageFilenames() {
		if opts.ImageProcessor == nil {
			handlers[p] = ioutils.IgnoreHandler
		} else {
			handlers[p] = func(ctx context.Context, r io.Reader) error {
				img, err := opts.ImageProcessor(ctx, manifest.GetMetadata().GetId(), p, r)
				if err != nil {
					return fmt.Errorf("error processing image: %v", err)
				}
				if processedAssets.Images == nil {
					processedAssets.Images = make(map[string]*ipb.Image)
				}
				processedAssets.Images[p] = img
				return nil
			}
		}
	}
	return processedAssets, handlers
}
