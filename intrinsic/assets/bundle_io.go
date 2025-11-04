// Copyright 2023 Intrinsic Innovation LLC

// Package bundleio contains a function that reads a bundle archive file.
package bundleio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/safearchive/tar"
	"google.golang.org/protobuf/proto"
	"intrinsic/assets/processes/processutil"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	processassetpb "intrinsic/assets/processes/proto/process_asset_go_proto"
	assettagpb "intrinsic/assets/proto/asset_tag_go_proto"
	assettypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_grpc_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	ipb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
)

// cloneOf clones a proto message while using generics to avoid a cast.
func cloneOf[M proto.Message](m M) M {
	return proto.Clone(m).(M)
}

type handler func(ctx context.Context, r io.Reader) error
type fallbackHandler func(context.Context, string, io.Reader) error

// ImageProcessor is a closure that pushes an image and returns the resulting
// pointer to the container registry.  It is provided the id of the bundle being
// processed as well as the name of the specific image.  It is expected to
// upload the image and produce a usable image spec.  The reader points to an
// image archive.  This may be invoked multiple times.  Images are ignored if it
// is not specified.
type ImageProcessor func(ctx context.Context, idProto *idpb.Id, filename string, r io.Reader) (*ipb.Image, error)

// walkTarFile walks through a tar file and invokes handlers on specific
// filenames.  fallback can be nil.  Returns an error if all handlers in
// handlers are not invoked.  It ignores all non-regular files.
func walkTarFile(ctx context.Context, t *tar.Reader, handlers map[string]handler, fallback fallbackHandler) error {
	for len(handlers) > 0 || fallback != nil {
		hdr, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("getting next file failed: %v", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		n := hdr.Name
		if h, ok := handlers[n]; ok {
			delete(handlers, n)
			if err := h(ctx, t); err != nil {
				return fmt.Errorf("error processing file %q: %v", n, err)
			}
		} else if fallback != nil {
			if err := fallback(ctx, n, t); err != nil {
				return fmt.Errorf("error processing file %q: %v", n, err)
			}
		}
	}
	if len(handlers) != 0 {
		keys := make([]string, 0, len(handlers))
		for k := range handlers {
			keys = append(keys, k)
		}
		return fmt.Errorf("missing expected files %s", keys)
	}
	return nil
}

// ignoreHandler is a function that can be used as a handler to ignore specific
// files.
func ignoreHandler(ctx context.Context, r io.Reader) error {
	return nil
}

// alwaysErrorAsUnexpected can be used as a fallback handler that will always
// trigger an unexpected file error.  This forces all files to be handled
// explicitly.
func alwaysErrorAsUnexpected(ctx context.Context, n string, r io.Reader) error {
	return fmt.Errorf("unexpected file %q", n)
}

// makeBinaryProtoHandler creates a handler that reads a binary proto file and
// unmarshals it into a file.  The proto must not be nil.
func makeBinaryProtoHandler(p proto.Message) handler {
	return func(ctx context.Context, r io.Reader) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading: %v", err)
		}
		if err := proto.Unmarshal(b, p); err != nil {
			return fmt.Errorf("error parsing proto: %v", err)
		}
		return nil
	}
}

// makeCollectInlinedFallbackHandler constructs a default handler that collects
// all of the unknown files and reads their bytes into a map.  The key of the
// map is the filename, and the value is the file contents.
func makeCollectInlinedFallbackHandler() (map[string][]byte, fallbackHandler) {
	inlined := map[string][]byte{}
	fallback := func(ctx context.Context, n string, r io.Reader) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading: %v", err)
		}
		inlined[n] = b
		return nil
	}
	return inlined, fallback
}

// readBinaryProto reads a binary proto from a reader and unmarshals it into a proto.
func readBinaryProto(r io.Reader, p proto.Message) error {
	if b, err := io.ReadAll(r); err != nil {
		return fmt.Errorf("error reading: %v", err)
	} else if err := proto.Unmarshal(b, p); err != nil {
		return fmt.Errorf("error parsing proto: %v", err)
	}

	return nil
}

// bundleType is used to return the type of a bundle file.
type bundleType int

// The different bundle types that can be detected from a file.
const (
	bundleTypeData bundleType = iota
	bundleTypeHardwareDevice
	bundleTypeProcess
	bundleTypeService
	bundleTypeSkill
)

var (
	errNoValidTypeDetected       = errors.New("no recognized manifest detected")
	errMultipleTypesDetected     = errors.New("invalid bundle")
)

// detectBundleType will return the type of bundle a file represents.  It does
// not do any validation of the particular file, just provides an indication
// what sort of processing should be done on the file.
func detectBundleType(ctx context.Context, path string) (bundleType, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open %q: %v", path, err)
	}
	defer f.Close()

	lookup := map[string]bundleType{
		dataAssetFileName:              bundleTypeData,
		hardwareDeviceManifestFileName: bundleTypeHardwareDevice,
		ProcessManifestFileName:        bundleTypeProcess,
		serviceManifestPathInTar:       bundleTypeService,
		skillManifestPathInTar:         bundleTypeSkill,
	}

	var bt bundleType
	var found int
	if err := walkTarFile(ctx, tar.NewReader(f), map[string]handler{}, func(_ context.Context, path string, _ io.Reader) error {
		if val, ok := lookup[path]; ok {
			found++
			bt = val
		}
		return nil
	}); err != nil {
		return bt, err
	}
	switch found {
	case 0:
		return 0, errNoValidTypeDetected
	case 1:
		return bt, nil
	default:
		return 0, errMultipleTypesDetected
	}
}

// BundleProcessor provides a way to process bundles of arbitrary types.  The
// processors are specific to a particular target (i.e. cluster or catalog) and
// should be for use across many bundles.
type BundleProcessor struct {
	ImageProcessor
	// ProcessReferencedData is the ReferencedDataProcessor to use for Data assets (see
	// ReadDataAsset).
	ProcessReferencedData ReferencedDataProcessor
}

// VersionDetails provides the specific details about a version when it is
// released to the catalog.
type VersionDetails struct {
	Version         string
	ReleaseNotes    string
	ReleaseMetadata *rmpb.ReleaseMetadata
}

// ProcessedBundle is a bundle that has been processed and can be viewed as a
// message for use in different outbound requests.
type ProcessedBundle interface {
	Install() *iapb.CreateInstalledAssetRequest_Asset
	Release(VersionDetails) *acpb.Asset
}

type dataBundle struct {
	manifest *dapb.DataAsset
}

func (b dataBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Data{
			Data: cloneOf(b.manifest),
		},
	}
}

func (b dataBundle) Release(details VersionDetails) *acpb.Asset {
	data := cloneOf(b.manifest)
	data.Metadata.IdVersion.Version = details.Version
	data.Metadata.ReleaseNotes = details.ReleaseNotes
	return &acpb.Asset{
		Metadata:        data.GetMetadata(),
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_DataSpecificDeploymentData{
				DataSpecificDeploymentData: &acpb.Asset_DataDeploymentData{
					Data: data,
				},
			},
		},
	}
}

type hardwareDeviceBundle struct {
	manifest *hdmpb.ProcessedHardwareDeviceManifest
}

func (b hardwareDeviceBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_HardwareDevice{
			HardwareDevice: cloneOf(b.manifest),
		},
	}
}

func (b hardwareDeviceBundle) Release(details VersionDetails) *acpb.Asset {
	manifest := cloneOf(b.manifest)

	// Take the first tag, if one exists.  Validation can be done later on the
	// deployment data.
	var tag assettagpb.AssetTag
	if len(manifest.GetMetadata().GetAssetTags()) > 1 {
		tag = manifest.GetMetadata().GetAssetTags()[0]
	}
	return &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id:      manifest.GetMetadata().GetId(),
				Version: details.Version,
			},
			AssetType:     assettypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE,
			AssetTag:      tag,
			DisplayName:   manifest.GetMetadata().GetDisplayName(),
			Documentation: manifest.GetMetadata().GetDocumentation(),
			Vendor:        manifest.GetMetadata().GetVendor(),
			ReleaseNotes:  details.ReleaseNotes,
		},
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_HardwareDeviceSpecificDeploymentData{
				HardwareDeviceSpecificDeploymentData: &acpb.Asset_HardwareDeviceDeploymentData{
					Manifest: manifest,
				},
			},
		},
	}
}

type processBundle struct {
	manifest *processassetpb.ProcessAsset
}

func (b processBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Process{
			Process: cloneOf(b.manifest),
		},
	}
}

func (b processBundle) Release(details VersionDetails) *acpb.Asset {
	manifest := cloneOf(b.manifest)
	manifest.Metadata.IdVersion.Version = details.Version
	manifest.Metadata.ReleaseNotes = details.ReleaseNotes

	// We have just updated the version in the metadata, sync the version-related
	// fields in Skill metadata in the BehaviorTree.
	processutil.FillInSkillIDVersionFromAssetMetadata(manifest.GetBehaviorTree().GetDescription(), manifest.GetMetadata())
	return &acpb.Asset{
		Metadata:        manifest.GetMetadata(),
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_ProcessSpecificDeploymentData{
				ProcessSpecificDeploymentData: &acpb.Asset_ProcessDeploymentData{
					Process: manifest,
				},
			},
		},
	}
}

type serviceBundle struct {
	manifest *smpb.ProcessedServiceManifest
}

func (b serviceBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Service{
			Service: cloneOf(b.manifest),
		},
	}
}

func (b serviceBundle) Release(details VersionDetails) *acpb.Asset {
	manifest := cloneOf(b.manifest)
	return &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id:      manifest.GetMetadata().GetId(),
				Version: details.Version,
			},
			AssetType:     assettypepb.AssetType_ASSET_TYPE_SERVICE,
			DisplayName:   manifest.GetMetadata().GetDisplayName(),
			Documentation: manifest.GetMetadata().GetDocumentation(),
			Vendor:        manifest.GetMetadata().GetVendor(),
			AssetTag:      manifest.GetMetadata().GetAssetTag(),
			ReleaseNotes:  details.ReleaseNotes,
		},
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_ServiceSpecificDeploymentData{
				ServiceSpecificDeploymentData: &acpb.Asset_ServiceDeploymentData{
					Manifest: manifest,
				},
			},
		},
	}
}

type skillBundle struct {
	manifest *psmpb.ProcessedSkillManifest
}

func (b skillBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Skill{
			Skill: cloneOf(b.manifest),
		},
	}
}

func (b skillBundle) Release(details VersionDetails) *acpb.Asset {
	manifest := cloneOf(b.manifest)
	return &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id:      manifest.GetMetadata().GetId(),
				Version: details.Version,
			},
			AssetType:     assettypepb.AssetType_ASSET_TYPE_SKILL,
			DisplayName:   manifest.GetMetadata().GetDisplayName(),
			Documentation: manifest.GetMetadata().GetDocumentation(),
			Vendor:        manifest.GetMetadata().GetVendor(),
			ReleaseNotes:  details.ReleaseNotes,
		},
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_SkillSpecificDeploymentData{
				SkillSpecificDeploymentData: &acpb.Asset_SkillDeploymentData{
					Manifest: manifest,
				},
			},
		},
	}
}

// Process auto-detects a bundle type and processes it to be sent to an
// appropriate target.
func (p *BundleProcessor) Process(ctx context.Context, path string) (ProcessedBundle, error) {
	bundleType, err := detectBundleType(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to detect bundle type: %w", err)
	}
	switch bundleType {
	case bundleTypeData:
		data, err := ReadDataAsset(ctx, path, WithProcessReferencedData(p.ProcessReferencedData))
		if err != nil {
			return nil, fmt.Errorf("unable to read data asset bundle: %w", err)
		}
		return dataBundle{data}, nil
	case bundleTypeHardwareDevice:
		assetInliner := NewLocalAssetInliner(LocalAssetInlinerOptions{
			ImageProcessor:          p.ImageProcessor,
			ProcessReferencedData:   p.ProcessReferencedData,
		})

		localAssetsDir, err := os.MkdirTemp("", "local-assets")
		if err != nil {
			return nil, fmt.Errorf("could not create temporary directory for local assets: %w", err)
		}
		defer os.RemoveAll(localAssetsDir)

		hardwareDevice, err := ProcessHardwareDevice(ctx, path,
			WithProcessAsset(assetInliner.Process),
			WithReadOptions(WithExtractLocalAssetsDir(localAssetsDir)),
		)
		if err != nil {
			return nil, fmt.Errorf("could not process HardwareDevice bundle: %w", err)
		}
		return &hardwareDeviceBundle{hardwareDevice}, nil
	case bundleTypeProcess:
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("could not open Process asset target %q: %w", path, err)
		}
		defer f.Close()
		process, err := ProcessProcessAsset(f)
		if err != nil {
			return nil, fmt.Errorf("unable to process process bundle: %w", err)
		}
		return processBundle{process}, nil
	case bundleTypeService:
		service, err := ProcessService(ctx, path, ProcessServiceOpts{
			ImageProcessor: p.ImageProcessor,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to process service bundle: %w", err)
		}
		return serviceBundle{service}, nil
	case bundleTypeSkill:
		skill, err := ProcessSkill(ctx, path, ProcessSkillOpts{
			ImageProcessor: p.ImageProcessor,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to process skill bundle: %w", err)
		}
		return skillBundle{skill}, nil
	default:
		return nil, fmt.Errorf("unable to detect bundle type: %w", err)
	}
}
