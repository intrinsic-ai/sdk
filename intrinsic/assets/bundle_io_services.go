// Copyright 2023 Intrinsic Innovation LLC

package bundleio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/safearchive/tar"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/services/servicemanifest"
	"intrinsic/util/archive/tartooling"
	"intrinsic/util/proto/registryutil"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
)

const (
	serviceManifestPathInTar = "service_manifest.binarypb"
)

// makeOnlyServiceManifestHandlers returns a map of handlers that only pull out
// the service manifest from the tar file into the returned proto.  Can be used
// with a fallback handler.
func makeOnlyServiceManifestHandlers() (*smpb.ServiceManifest, map[string]handler) {
	manifest := new(smpb.ServiceManifest)
	handlers := map[string]handler{
		serviceManifestPathInTar: makeBinaryProtoHandler(manifest),
	}
	return manifest, handlers
}

// makeServiceAssetHandlers returns handlers for all assets listed in the
// service manifest.  This will be at most:
// * An handler that ignores the manifest
// * A binary proto handler for the default configuration file
// * A binary proto handler for the file descriptor set file
// * A handler that wraps opts.ImageProcessor to be called on every image
func makeServiceAssetHandlers(manifest *smpb.ServiceManifest, opts ProcessServiceOpts) (*smpb.ProcessedServiceAssets, map[string]handler) {
	handlers := map[string]handler{
		serviceManifestPathInTar: ignoreHandler, // already read this.
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
		handlers[*p] = makeBinaryProtoHandler(processedAssets.DefaultConfiguration)
	}
	if p := manifest.GetAssets().ParameterDescriptorFilename; p != nil {
		processedAssets.FileDescriptorSet = new(descriptorpb.FileDescriptorSet)
		handlers[*p] = makeBinaryProtoHandler(processedAssets.FileDescriptorSet)
	}
	for _, p := range manifest.GetAssets().GetImageFilenames() {
		if opts.ImageProcessor == nil {
			handlers[p] = ignoreHandler
		} else {
			handlers[p] = func(r io.Reader) error {
				img, err := opts.ImageProcessor(manifest.GetMetadata().GetId(), p, r)
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

// ReadService reads the service bundle archive from path. It returns the
// service manifest and a mapping between bundle filenames and their contents.
func ReadService(path string) (*smpb.ServiceManifest, map[string][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	m, handlers := makeOnlyServiceManifestHandlers()
	inlined, fallback := makeCollectInlinedFallbackHandler()
	if err := walkTarFile(tar.NewReader(f), handlers, fallback); err != nil {
		return nil, nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}
	return m, inlined, nil
}

// ReadServiceManifest reads the bundle archive from path. It returns only
// service manifest.
func ReadServiceManifest(path string) (*smpb.ServiceManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	m, handlers := makeOnlyServiceManifestHandlers()
	if err := walkTarFile(tar.NewReader(f), handlers, nil); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}
	return m, nil
}

// ProcessServiceOpts contains the necessary handlers to generate a processed
// service manifest.
type ProcessServiceOpts struct {
	ImageProcessor
}

// ProcessService creates a processed manifest from a bundle on disk using the
// provided processing functions.  It avoids doing any validation except for
// that required to transform the specified files in the bundle into their
// processed variants.
func ProcessService(path string, opts ProcessServiceOpts) (*smpb.ProcessedServiceManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	// Read the manifest and then reset the file once we have the information
	// about the bundle we're going to process.
	manifest, handlers := makeOnlyServiceManifestHandlers()
	if err := walkTarFile(tar.NewReader(f), handlers, nil); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("could not seek in %q: %v", path, err)
	}

	// Initialize handlers for when we walk through the file again now that we
	// know what we're looking for, but error on unexpected files this time.
	processedAssets, handlers := makeServiceAssetHandlers(manifest, opts)
	if err := walkTarFile(tar.NewReader(f), handlers, alwaysErrorAsUnexpected); err != nil {
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

// ValidateService checks that the assets of a service bundle are all
// contained within the inlined file map.
func ValidateService(manifest *smpb.ServiceManifest, inlinedFiles map[string][]byte) error {
	files := make([]string, 0, len(inlinedFiles))
	usedFiles := make(map[string]bool)
	for f := range inlinedFiles {
		files = append(files, f)
		usedFiles[f] = true
	}
	fileNames := strings.Join(files, ", ")
	// Check that every defined asset is in the inlined filemap.
	assets := map[string]string{
		"default configuration file": manifest.GetAssets().GetDefaultConfigurationFilename(),
		"parameter descriptor file":  manifest.GetAssets().GetParameterDescriptorFilename(),
		"image tar":                  manifest.GetServiceDef().GetRealSpec().GetImage().GetArchiveFilename(),
		"simulation image tar":       manifest.GetServiceDef().GetSimSpec().GetImage().GetArchiveFilename(),
	}
	for desc, path := range assets {
		if path != "" {
			if _, ok := inlinedFiles[path]; !ok {
				return fmt.Errorf("the service manifest's %s %q is not in the bundle. files are %s", desc, path, fileNames)
			}
			delete(usedFiles, path)
		}
	}
	for _, path := range manifest.GetAssets().GetImageFilenames() {
		if _, ok := inlinedFiles[path]; !ok {
			return fmt.Errorf("the service manifest's image file %q is not in the bundle. files are %s", path, fileNames)
		}
		delete(usedFiles, path)
	}
	if len(usedFiles) > 0 {
		files := make([]string, 0, len(usedFiles))
		for f := range usedFiles {
			files = append(files, f)
		}
		fileNames := strings.Join(files, ", ")
		return fmt.Errorf("found unexpected files in the archive: %s", fileNames)
	}
	return nil
}

// WriteServiceOpts provides the details to construct a service bundle.
type WriteServiceOpts struct {
	Manifest      *smpb.ServiceManifest
	Descriptors   *descriptorpb.FileDescriptorSet
	DefaultConfig *anypb.Any
	ImageTars     []string
}

// WriteService creates a tar archive at the specified path with the details
// given in opts.  Only the manifest is required and its assets field will be
// overwritten with what is placed in the archive based on ops.
func WriteService(path string, opts WriteServiceOpts) error {
	if opts.Manifest == nil {
		return fmt.Errorf("opts.Manifest must not be nil")
	}
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tw := tar.NewWriter(out)

	opts.Manifest.Assets = new(smpb.ServiceAssets)
	if opts.Descriptors != nil {
		descriptorName := "descriptors-transitive-descriptor-set.proto.bin"
		opts.Manifest.Assets.ParameterDescriptorFilename = &descriptorName
		if err := tartooling.AddBinaryProto(opts.Descriptors, tw, descriptorName); err != nil {
			return fmt.Errorf("unable to write FileDescriptorSet to bundle: %v", err)
		}
	}
	if opts.DefaultConfig != nil {
		configName := "default_config.binarypb"
		opts.Manifest.Assets.DefaultConfigurationFilename = &configName
		if err := tartooling.AddBinaryProto(opts.DefaultConfig, tw, configName); err != nil {
			return fmt.Errorf("unable to write default config to bundle: %v", err)
		}
	}
	for _, path := range opts.ImageTars {
		base := filepath.Base(path)
		opts.Manifest.Assets.ImageFilenames = append(opts.Manifest.Assets.ImageFilenames, base)
		if err := tartooling.AddFile(path, tw, base); err != nil {
			return fmt.Errorf("unable to write %q to bundle: %v", path, err)
		}
	}

	var files *protoregistry.Files
	if opts.Descriptors != nil {
		files, err = protodesc.NewFiles(opts.Descriptors)
		if err != nil {
			return fmt.Errorf("failed to create proto files: %v", err)
		}
	}
	if err := servicemanifest.ValidateServiceManifest(opts.Manifest,
		servicemanifest.WithFiles(files),
		servicemanifest.WithDefaultConfig(opts.DefaultConfig),
	); err != nil {
		return fmt.Errorf("invalid manifest: %v", err)
	}

	// Now we can write the manifest, since assets have been completed.
	if err := tartooling.AddBinaryProto(opts.Manifest, tw, "service_manifest.binarypb"); err != nil {
		return fmt.Errorf("unable to write FileDescriptorSet to bundle: %v", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
