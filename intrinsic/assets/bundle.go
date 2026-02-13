// Copyright 2023 Intrinsic Innovation LLC

// Package bundle contains utils for working with Asset bundles.
package bundle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"intrinsic/assets/data/databundle"
	"intrinsic/assets/hardware_devices/hardwaredevicebundle"
	"intrinsic/assets/imageutils"
	"intrinsic/assets/ioutils"
	"intrinsic/assets/processes/processbundle"
	"intrinsic/assets/processes/processvalidate"
	"intrinsic/assets/services/servicebundle"
	"intrinsic/skills/skillbundle"

	"github.com/google/safearchive/tar"
	"google.golang.org/protobuf/proto"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	processassetpb "intrinsic/assets/processes/proto/process_asset_go_proto"
	assettagpb "intrinsic/assets/proto/asset_tag_go_proto"
	assettypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
)

const (
	dataAssetFileName = "data_asset.binpb"
	hardwareDeviceManifestFileName = "hardware_device_manifest.binpb"
	processManifestFileName = "process_manifest.binpb"
	serviceManifestPathInTar = "service_manifest.binarypb"
	skillManifestPathInTar = "skill_manifest.binpb"
)

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
	errNoValidTypeDetected   = errors.New("no recognized manifest detected")
	errMultipleTypesDetected = errors.New("invalid bundle")
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
		processManifestFileName:        bundleTypeProcess,
		serviceManifestPathInTar:       bundleTypeService,
		skillManifestPathInTar:         bundleTypeSkill,
	}

	var bt bundleType
	var found int
	if err := ioutils.WalkTarFile(ctx, tar.NewReader(f), ioutils.WithFallbackHandler(func(_ context.Context, path string, _ io.Reader) error {
		if val, ok := lookup[path]; ok {
			found++
			bt = val
		}
		return nil
	})); err != nil {
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

// Processor provides a way to process bundles of arbitrary types.  The
// processors are specific to a particular target (i.e. cluster or catalog) and
// should be for use across many bundles.
type Processor struct {
	imageutils.ImageProcessor
	// ProcessReferencedData is the databundle.ReferencedDataProcessor to use for Data assets (see
	// ReadDataBundle).
	ProcessReferencedData databundle.ReferencedDataProcessor
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

type processedDataBundle struct {
	da *dapb.DataAsset
}

func (b processedDataBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Data{
			Data: cloneOf(b.da),
		},
	}
}

func (b processedDataBundle) Release(details VersionDetails) *acpb.Asset {
	da := cloneOf(b.da)
	m := cloneOf(da.GetMetadata())
	m.IdVersion.Version = details.Version
	m.ReleaseNotes = details.ReleaseNotes
	return &acpb.Asset{
		Metadata:        m,
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_DataSpecificDeploymentData{
				DataSpecificDeploymentData: &acpb.Asset_DataDeploymentData{
					Data: da,
				},
			},
		},
	}
}

type processedHardwareDeviceBundle struct {
	manifest *hdmpb.ProcessedHardwareDeviceManifest
}

func (b processedHardwareDeviceBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_HardwareDevice{
			HardwareDevice: cloneOf(b.manifest),
		},
	}
}

func (b processedHardwareDeviceBundle) Release(details VersionDetails) *acpb.Asset {
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

type processedProcessBundle struct {
	pa *processassetpb.ProcessAsset
}

func (b processedProcessBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Process{
			Process: cloneOf(b.pa),
		},
	}
}

func (b processedProcessBundle) Release(details VersionDetails) *acpb.Asset {
	pa := cloneOf(b.pa)
	m := cloneOf(pa.GetMetadata())
	m.IdVersion.Version = details.Version
	m.ReleaseNotes = details.ReleaseNotes
	processvalidate.FillBackwardsCompatibleVersion(pa, details.Version)

	return &acpb.Asset{
		Metadata:        m,
		ReleaseMetadata: details.ReleaseMetadata,
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_ProcessSpecificDeploymentData{
				ProcessSpecificDeploymentData: &acpb.Asset_ProcessDeploymentData{
					Process: pa,
				},
			},
		},
	}
}

type processedServiceBundle struct {
	manifest *smpb.ProcessedServiceManifest
}

func (b processedServiceBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Service{
			Service: cloneOf(b.manifest),
		},
	}
}

func (b processedServiceBundle) Release(details VersionDetails) *acpb.Asset {
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

type processedSkillBundle struct {
	manifest *psmpb.ProcessedSkillManifest
}

func (b processedSkillBundle) Install() *iapb.CreateInstalledAssetRequest_Asset {
	return &iapb.CreateInstalledAssetRequest_Asset{
		Variant: &iapb.CreateInstalledAssetRequest_Asset_Skill{
			Skill: cloneOf(b.manifest),
		},
	}
}

func (b processedSkillBundle) Release(details VersionDetails) *acpb.Asset {
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
func (p *Processor) Process(ctx context.Context, path string) (ProcessedBundle, error) {
	bundleType, err := detectBundleType(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect bundle type: %w", err)
	}
	switch bundleType {
	case bundleTypeData:
		da, err := databundle.Process(ctx, path,
			databundle.WithReadOptions(databundle.WithProcessReferencedData(p.ProcessReferencedData)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to process Data Asset bundle: %w", err)
		}
		return processedDataBundle{da}, nil
	case bundleTypeHardwareDevice:
		assetInliner := hardwaredevicebundle.NewLocalAssetInliner(hardwaredevicebundle.LocalAssetInlinerOptions{
			ImageProcessor:          p.ImageProcessor,
			ProcessReferencedData:   p.ProcessReferencedData,
		})

		localAssetsDir, err := os.MkdirTemp("", "local-assets")
		if err != nil {
			return nil, fmt.Errorf("failed create temporary directory for local Assets: %w", err)
		}
		defer os.RemoveAll(localAssetsDir)

		hardwareDevice, err := hardwaredevicebundle.Process(ctx, path,
			hardwaredevicebundle.WithProcessAsset(assetInliner.Process),
			hardwaredevicebundle.WithReadOptions(hardwaredevicebundle.WithExtractLocalAssetsDir(localAssetsDir)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to process HardwareDevice bundle: %w", err)
		}
		return &processedHardwareDeviceBundle{hardwareDevice}, nil
	case bundleTypeProcess:
		process, err := processbundle.Process(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to process Process bundle: %w", err)
		}
		return processedProcessBundle{process}, nil
	case bundleTypeService:
		service, err := servicebundle.Process(ctx, path,
			servicebundle.WithImageProcessor(p.ImageProcessor),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to process Service bundle: %w", err)
		}
		return processedServiceBundle{service}, nil
	case bundleTypeSkill:
		skill, err := skillbundle.Process(ctx, path,
			skillbundle.WithImageProcessor(p.ImageProcessor),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to process Skill bundle: %w", err)
		}
		return processedSkillBundle{skill}, nil
	default:
		return nil, fmt.Errorf("unknown bundle type: %v", bundleType)
	}
}

// cloneOf clones a proto message while using generics to avoid a cast.
func cloneOf[M proto.Message](m M) M {
	return proto.Clone(m).(M)
}
