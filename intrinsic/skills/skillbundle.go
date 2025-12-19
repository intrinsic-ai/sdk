// Copyright 2023 Intrinsic Innovation LLC

// Package databundle provides utils for working with Skill bundles.
package skillbundle

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"intrinsic/assets/imageutils"
	"intrinsic/assets/ioutils"
	"intrinsic/skills/internal/skillmanifest"
	"intrinsic/util/archive/tartooling"
	"intrinsic/util/proto/registryutil"

	"github.com/google/safearchive/tar"
	"google.golang.org/protobuf/proto"

	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const (
	skillManifestPathInTar = "skill_manifest.binpb"
)

// SkillBundle represents a Skill Asset bundle.
type SkillBundle struct {
	Manifest *smpb.SkillManifest
	Files    map[string][]byte
}

type readSkillBundleOptions struct {
	readFiles bool
}

// ReadSkillBundleOption is a functional option for ReadSkillBundle
type ReadSkillBundleOption func(*readSkillBundleOptions)

// WithReadFiles specifies whether to read additional files when reading the bundle.
func WithReadFiles(b bool) ReadSkillBundleOption {
	return func(opts *readSkillBundleOptions) {
		opts.readFiles = b
	}
}

// ReadSkillBundle reads a Skill Asset bundle from disk (see WriteSkillBundle).
func ReadSkillBundle(ctx context.Context, path string, options ...ReadSkillBundleOption) (*SkillBundle, error) {
	opts := &readSkillBundleOptions{}
	for _, opt := range options {
		opt(opts)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	m, handlers := makeOnlySkillManifestHandlers()
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

	return &SkillBundle{
		Manifest: m,
		Files:    inlined,
	}, nil
}

// ProcessSkillOpts contains the necessary handlers to generate a processed
// skill manifest.
type ProcessSkillOpts struct {
	imageutils.ImageProcessor
}

// ProcessSkill creates a processed manifest from a bundle on disk using the
// provided processing functions. It avoids doing any validation except for
// that required to transform the specified files in the bundle into their
// processed variants.
func ProcessSkill(ctx context.Context, path string, opts ProcessSkillOpts) (*psmpb.ProcessedSkillManifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	// Read the manifest and then reset the file once we have the information
	// about the bundle we're going to process.
	manifest, handlers := makeOnlySkillManifestHandlers()
	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f), ioutils.WithHandlers(handlers)); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("could not seek in %q: %v", path, err)
	}

	// Initialize handlers for when we walk through the file again now that we
	// know what we're looking for, but error on unexpected files this time.
	processedAssets, handlers := makeSkillAssetHandlers(manifest, opts)
	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f),
		ioutils.WithHandlers(handlers),
		ioutils.WithFallbackHandler(ioutils.AlwaysErrorAsUnexpected),
	); err != nil {
		return nil, fmt.Errorf("error in tar file %q: %v", path, err)
	}

	psm := &psmpb.ProcessedSkillManifest{
		Assets: processedAssets,
	}
	m := &psmpb.SkillMetadata{
		Id:            manifest.GetId(),
		Vendor:        manifest.GetVendor(),
		Documentation: manifest.GetDocumentation(),
		DisplayName:   manifest.GetDisplayName(),
	}
	// We do this to avoid adding empty sub messages.
	if !proto.Equal(m, &psmpb.SkillMetadata{}) {
		psm.Metadata = m
	}
	d := &psmpb.SkillDetails{
		Options:       manifest.GetOptions(),
		Dependencies:  manifest.GetDependencies(),
		Parameter:     manifest.GetParameter(),
		ExecuteResult: manifest.GetReturnType(),
		StatusInfo:    manifest.GetStatusInfo(),
	}
	// We do this to avoid adding empty sub messages.
	if !proto.Equal(d, &psmpb.SkillDetails{}) {
		psm.Details = d
	}
	return psm, nil
}

// WriteSkillBundleOptions provides options for writing a Skill Asset bundle.
type WriteSkillBundleOptions struct {
	Descriptors  *descriptorpb.FileDescriptorSet
	ImageTarPath string
}

// WriteSkillBundle writes a Skill .tar bundle at the specified path.
func WriteSkillBundle(m *smpb.SkillManifest, path string, opts *WriteSkillBundleOptions) error {
	if m == nil {
		return fmt.Errorf("SkillManifest must not be nil")
	}

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tw := tar.NewWriter(out)

	m.Assets = &smpb.SkillAssets{}
	if opts.Descriptors != nil {
		descriptorName := "descriptors-transitive-descriptor-set.proto.bin"
		m.Assets.FileDescriptorSetFilename = &descriptorName
		if err := tartooling.AddBinaryProto(opts.Descriptors, tw, descriptorName); err != nil {
			return fmt.Errorf("failed to write FileDescriptorSet to bundle: %w", err)
		}
	}
	if opts.ImageTarPath != "" {
		base := filepath.Base(opts.ImageTarPath)
		m.Assets.DeploymentType = &smpb.SkillAssets_ImageFilename{
			ImageFilename: base,
		}
		if err := tartooling.AddFile(opts.ImageTarPath, tw, base); err != nil {
			return fmt.Errorf("unable to write %q to bundle: %w", path, err)
		}
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(opts.Descriptors)
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %w", err)
	}
	if err := skillmanifest.ValidateSkillManifest(m,
		skillmanifest.WithTypes(types),
		skillmanifest.WithIncompatibleDisallowManifestDependencies(false),
	); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	// Now we can write the manifest, since assets have been completed.
	if err := tartooling.AddBinaryProto(m, tw, skillManifestPathInTar); err != nil {
		return fmt.Errorf("failed to write SkillManifest to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}

// makeOnlySkillManifestHandlers returns a map of handlers that only pull out
// the skill manifest from the tar file into the returned proto. Can be used
// with a fallback handler.
func makeOnlySkillManifestHandlers() (*smpb.SkillManifest, map[string]ioutils.WalkTarFileHandler) {
	manifest := new(smpb.SkillManifest)
	handlers := map[string]ioutils.WalkTarFileHandler{
		skillManifestPathInTar: ioutils.MakeBinaryProtoHandler(manifest),
	}
	return manifest, handlers
}

// makeSkillAssetHandlers returns handlers for all assets listed in the
// skill manifest. This will be at most:
// * A handler that ignores the manifest
// * A binary proto handler for the file descriptor set file
// * A handler that wraps opts.ImageProcessor to be called on every image
// * A binary proto handler for the parameterized behavior tree file
func makeSkillAssetHandlers(manifest *smpb.SkillManifest, opts ProcessSkillOpts) (*psmpb.ProcessedSkillAssets, map[string]ioutils.WalkTarFileHandler) {
	handlers := map[string]ioutils.WalkTarFileHandler{
		skillManifestPathInTar: ioutils.IgnoreHandler, // already read this.
	}
	// Don't generate an empty assets message if there wasn't one to begin
	// with. This is a slightly odd state, but Process is not doing validation of
	// the manifest. This also protects against nil access of
	// manifest.GetAssets().{MemberVariable}, which is required for checking the
	// "optional" piece of "optional string" fields in this version of the golang
	// proto API.
	if manifest.GetAssets() == nil {
		return nil, handlers
	}

	processedAssets := &psmpb.ProcessedSkillAssets{}
	if p := manifest.GetAssets().FileDescriptorSetFilename; p != nil {
		processedAssets.FileDescriptorSet = new(descriptorpb.FileDescriptorSet)
		handlers[*p] = ioutils.MakeBinaryProtoHandler(processedAssets.FileDescriptorSet)
	}
	switch manifest.GetAssets().GetDeploymentType().(type) {
	case *smpb.SkillAssets_ImageFilename:
		p := manifest.GetAssets().GetImageFilename()
		handlers[p] = func(ctx context.Context, r io.Reader) error {
			img, err := opts.ImageProcessor(ctx, manifest.GetId(), p, r)
			if err != nil {
				return fmt.Errorf("error processing image: %v", err)
			}
			processedAssets.DeploymentType = &psmpb.ProcessedSkillAssets_Image{
				Image: img,
			}
			return nil
		}
	}
	return processedAssets, handlers
}
