// Copyright 2023 Intrinsic Innovation LLC

package bundleio

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/google/safearchive/tar"
	"intrinsic/assets/hardware_devices/hardwaredevicemanifest"
	"intrinsic/assets/idutils"
	"intrinsic/util/archive/tartooling"

	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
)

const (
	hardwareDeviceManifestFileName = "hardware_device_manifest.binpb"
)

var (
	tarBundlePathRegex = regexp.MustCompile(`^assets/(?P<key>[^/]+)\.bundle\.tar$`)
)

// writeHardwareDeviceOptions contains options for a call to WriteHardwareDevice.
type writeHardwareDeviceOptions struct {
}

// WriteHardwareDeviceOption is a functional option for WriteHardwareDevice.
type WriteHardwareDeviceOption func(*writeHardwareDeviceOptions)

// WriteHardwareDevice writes a HardwareDevice asset .tar bundle file to the specified path.
//
// The bundles of local assets are included in the HardwareDevice .tar bundle.
//
// The input manifest may be mutated by this function.
func WriteHardwareDevice(hdm *hdmpb.HardwareDeviceManifest, path string, options ...WriteHardwareDeviceOption) error {
	opts := &writeHardwareDeviceOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if hdm == nil {
		return fmt.Errorf("HardwareDeviceManifest must not be nil")
	}

	if err := hardwaredevicemanifest.ValidateHardwareDeviceManifest(hdm); err != nil {
		return fmt.Errorf("invalid HardwareDeviceManifest: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer f.Close()
	tw := tar.NewWriter(f)

	// Save local assets into the bundle and update their paths in the manifest.
	for key, asset := range hdm.GetAssets() {
		switch asset.GetVariant().(type) {
		case *hdmpb.HardwareDeviceManifest_Asset_Local:
			tarPath := tarBundlePathFrom(key)
			if err := tartooling.AddFile(asset.GetLocal().GetBundlePath(), tw, tarPath); err != nil {
				return fmt.Errorf("cannot add local asset %s to bundle: %w", key, err)
			}
			asset.GetLocal().BundlePath = tarPath
		}
	}

	if err := tartooling.AddBinaryProto(hdm, tw, hardwareDeviceManifestFileName); err != nil {
		return fmt.Errorf("cannot write HardwareDeviceManifest to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("cannot close tar writer: %w", err)
	}

	return nil
}

// readHardwareDeviceOptions contains options for a call to ReadHardwareDevice.
type readHardwareDeviceOptions struct {
	// extractLocalAssetsDir is the directory to which to extract local asset bundles.
	extractLocalAssetsDir string
}

// ReadHardwareDeviceOption is a functional option for ReadHardwareDevice.
type ReadHardwareDeviceOption func(*readHardwareDeviceOptions)

// WithExtractLocalAssetsDir provides a directory to which to extract local asset bundles.
//
// If specified, local asset bundles will be extracted to this directory, and bundle paths updated
// in the returned manifest. The directory must already exist.
func WithExtractLocalAssetsDir(dir string) ReadHardwareDeviceOption {
	return func(opts *readHardwareDeviceOptions) {
		opts.extractLocalAssetsDir = dir
	}
}

// ReadHardwareDevice reads a HardwareDevice asset from a .tar bundle (see WriteHardwareDevice).
func ReadHardwareDevice(p string, options ...ReadHardwareDeviceOption) (*hdmpb.HardwareDeviceManifest, error) {
	opts := &readHardwareDeviceOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Open the tar file for reading.
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", p, err)
	}
	defer f.Close()

	// Process the .tar bundle.
	var hdm *hdmpb.HardwareDeviceManifest
	extractedBundlePaths := map[string]string{}
	manifestHandler := func(r io.Reader) error {
		hdm = &hdmpb.HardwareDeviceManifest{}
		if err := readBinaryProto(r, hdm); err != nil {
			return fmt.Errorf("error reading HardwareDeviceManifest: %w", err)
		}
		return nil
	}
	fallbackHandler := func(n string, r io.Reader) error {
		if key, ok := tryExtractAssetKey(n); ok {
			// Ignore local asset bundles if no extraction directory was provided.
			if opts.extractLocalAssetsDir == "" {
				return nil
			}
			if _, ok := extractedBundlePaths[key]; ok {
				return fmt.Errorf("duplicate local asset bundle %q", key)
			}
			extractedBundlePaths[key] = path.Join(opts.extractLocalAssetsDir, fmt.Sprintf("%s.bundle.tar", key))
			if err := writeToFile(r, extractedBundlePaths[key]); err != nil {
				return fmt.Errorf("error writing local asset bundle %q: %w", key, err)
			}
			return nil
		}
		return fmt.Errorf("unexpected file %q", n)
	}
	if err := walkTarFile(
		tar.NewReader(f),
		map[string]handler{
			hardwareDeviceManifestFileName: manifestHandler,
		},
		fallbackHandler,
	); err != nil {
		return nil, fmt.Errorf("error processing tar file %q: %w", p, err)
	}

	// Replace local asset bundle paths with their extracted paths.
	if opts.extractLocalAssetsDir != "" {
		for key, asset := range hdm.GetAssets() {
			switch asset.GetVariant().(type) {
			case *hdmpb.HardwareDeviceManifest_Asset_Local:
				extractedPath, ok := extractedBundlePaths[key]
				if !ok {
					return nil, fmt.Errorf("extracted bundle path for local asset %s not found", key)
				}
				asset.GetLocal().BundlePath = extractedPath
			}
		}
	}

	return hdm, nil
}

// AssetProcessor is a function that processes a single asset.
type AssetProcessor func(*hdmpb.HardwareDeviceManifest_Asset) (*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, error)

// PassThrough is an AssetProcessor that passes asset catalog references through unchanged.
func PassThrough(a *hdmpb.HardwareDeviceManifest_Asset) (*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, error) {
	switch a.GetVariant().(type) {
	case *hdmpb.HardwareDeviceManifest_Asset_Catalog:
		return &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog{
				Catalog: a.GetCatalog(),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported asset type: %T", a.GetVariant())
	}
}

// LocalAssetInlinerOptions contains options for LocalAssetInliner.
type LocalAssetInlinerOptions struct {
	ImageProcessor
}

// LocalAssetInliner processes local assets in a HardwareDevice by inlining them.
//
// Its Process method can be provided as an AssetProcessor to ProcessHardwareDevice.
type LocalAssetInliner struct {
	opts LocalAssetInlinerOptions
}

// Process is an AssetProcessor that processes a local asset bundle.
func (p *LocalAssetInliner) Process(a *hdmpb.HardwareDeviceManifest_Asset) (*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, error) {
	switch a.GetVariant().(type) {
	case *hdmpb.HardwareDeviceManifest_Asset_Local:
		switch at := a.GetLocal().GetAssetType(); at {
		case atpb.AssetType_ASSET_TYPE_SERVICE:
			psm, err := ProcessService(a.GetLocal().GetBundlePath(), ProcessServiceOpts{
				ImageProcessor: p.opts.ImageProcessor,
			})
			if err != nil {
				return nil, fmt.Errorf("error processing Service %s: %w", idutils.IDFromProtoUnchecked(a.GetLocal().GetId()), err)
			}
			return &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
				Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service{
					Service: psm,
				},
			}, nil
		default:
			return nil, fmt.Errorf("unsupported local asset type: %s", at)
		}
	default:
		return PassThrough(a)
	}
}

// NewLocalAssetInliner creates a LocalAssetInliner with the given options.
func NewLocalAssetInliner(opts LocalAssetInlinerOptions) *LocalAssetInliner {
	return &LocalAssetInliner{opts: opts}
}

// processHardwareDeviceOptions contains options for a call to ProcessHardwareDevice.
type processHardwareDeviceOptions struct {
	// processAsset is a function that processes a single asset.
	processAsset AssetProcessor
	// readOptions are options to pass to ReadHardwareDevice.
	readOptions []ReadHardwareDeviceOption
}

// ProcessHardwareDeviceOption is a functional option for ProcessHardwareDevice.
type ProcessHardwareDeviceOption func(*processHardwareDeviceOptions)

// WithProcessAsset provides a function to process a single asset.
//
// If unspecified, a default processor will be used.
func WithProcessAsset(f AssetProcessor) ProcessHardwareDeviceOption {
	return func(opts *processHardwareDeviceOptions) {
		opts.processAsset = f
	}
}

// WithReadOptions provides options to pass to ReadHardwareDevice.
func WithReadOptions(options ...ReadHardwareDeviceOption) ProcessHardwareDeviceOption {
	return func(opts *processHardwareDeviceOptions) {
		opts.readOptions = options
	}
}

// ProcessHardwareDevice creates a processed manifest from a bundle on disk.
//
// It assumes that the bundle has already been validated.
func ProcessHardwareDevice(path string, options ...ProcessHardwareDeviceOption) (*hdmpb.ProcessedHardwareDeviceManifest, error) {
	opts := &processHardwareDeviceOptions{
		processAsset: PassThrough,
	}
	for _, opt := range options {
		opt(opts)
	}

	hdm, err := ReadHardwareDevice(path, opts.readOptions...)
	if err != nil {
		return nil, fmt.Errorf("error reading HardwareDeviceManifest: %w", err)
	}

	// Process each asset.
	processedAssets := make(map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, len(hdm.GetAssets()))
	for key, asset := range hdm.GetAssets() {
		processedAsset, err := opts.processAsset(asset)
		if err != nil {
			return nil, err
		}
		switch processedAsset.GetVariant().(type) {
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog:
			idVersion := processedAsset.GetCatalog().GetIdVersion()
			if idVersion.GetVersion() == "" {
				return nil, fmt.Errorf("asset %s does not specify a version", idutils.IDFromProtoUnchecked(idVersion.GetId()))
			}
		}
		processedAssets[key] = processedAsset
	}

	return &hdmpb.ProcessedHardwareDeviceManifest{
		Metadata: hdm.GetMetadata(),
		Assets:   processedAssets,
		Graph:    hdm.GetGraph(),
	}, nil
}

// tarBundlePathFrom returns the in-tar path for a local asset bundle with the given key.
func tarBundlePathFrom(key string) string {
	return fmt.Sprintf("assets/%s.bundle.tar", key)
}

// tryExtractAssetKey returns the key of the asset bundle at the given in-tar path, or the empty
// string if the path is not a valid in-tar asset bundle path.
//
// The second return value is true if the path is a valid in-tar asset bundle path.
func tryExtractAssetKey(path string) (string, bool) {
	submatches := tarBundlePathRegex.FindStringSubmatch(path)
	if submatches == nil {
		return "", false
	}
	return submatches[1], true
}

// writeToFile reads a file from r and writes it to the specified path.
func writeToFile(r io.Reader, path string) error {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy file to %q: %w", path, err)
	}
	return nil
}
