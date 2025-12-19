// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevicegen implements creation of a HardwareDevice Asset bundle.
package hardwaredevicegen

import (
	"fmt"
	"os"

	"intrinsic/assets/hardware_devices/hardwaredevicebundle"
	"intrinsic/assets/idutils"
	"intrinsic/util/proto/protoio"

	apb "intrinsic/assets/build_defs/asset_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	rpb "intrinsic/assets/proto/v1/reference_go_proto"
)

// CreateHardwareDeviceBundleOptions provides the data needed to create a HardwareDevice Asset
// bundle.
type CreateHardwareDeviceBundleOptions struct {
	// AssetCatalogRefInfoPaths are the paths to serialized AssetCatalogRefInfo protos of assets to
	// add to the manifest.
	AssetCatalogRefInfoPaths []string
	// AssetLocalBundlePaths are the paths to asset bundle .tar files that correspond to assets in
	// AssetLocalInfoPaths.
	AssetLocalInfoPaths []string
	// Manifest is the path to a HardwareDeviceManifest .textproto file.
	ManifestPath string
	// OutputBundlePath is the output path for the tar bundle.
	OutputBundlePath string
}

// CreateHardwareDeviceBundle creates a HardwareDevice Asset bundle on disk.
func CreateHardwareDeviceBundle(opts *CreateHardwareDeviceBundleOptions) error {
	m := &hdmpb.HardwareDeviceManifest{}
	if err := protoio.ReadTextProto(opts.ManifestPath, m); err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	assets := m.GetAssets()
	if assets == nil {
		assets = map[string]*hdmpb.HardwareDeviceManifest_Asset{}
	}
	for _, path := range opts.AssetCatalogRefInfoPaths {
		info := &apb.AssetCatalogRefInfo{}
		if err := protoio.ReadBinaryProto(path, info); err != nil {
			return fmt.Errorf("failed to read AssetCatalogRefInfo: %w", err)
		}
		id := idutils.IDFromProtoUnchecked(info.GetIdVersion().GetId())
		if _, ok := assets[id]; ok {
			return fmt.Errorf("asset %s already exists in manifest", id)
		}
		assets[id] = &hdmpb.HardwareDeviceManifest_Asset{
			Variant: &hdmpb.HardwareDeviceManifest_Asset_Catalog{
				Catalog: &rpb.CatalogAsset{
					AssetType: info.GetAssetType(),
					IdVersion: info.GetIdVersion(),
				},
			},
		}
	}
	for _, path := range opts.AssetLocalInfoPaths {
		info := &apb.AssetLocalInfo{}
		if err := protoio.ReadBinaryProto(path, info); err != nil {
			return fmt.Errorf("failed to read AssetLocalInfo: %w", err)
		}
		id := idutils.IDFromProtoUnchecked(info.GetId())
		if _, ok := assets[id]; ok {
			return fmt.Errorf("asset %s already exists in manifest", id)
		}
		if _, err := os.Stat(info.GetBundlePath()); err != nil {
			return fmt.Errorf("asset %s has invalid bundle path %q: %w", id, info.GetBundlePath(), err)
		}
		assets[id] = &hdmpb.HardwareDeviceManifest_Asset{
			Variant: &hdmpb.HardwareDeviceManifest_Asset_Local{
				Local: &rpb.LocalAsset{
					AssetType:  info.GetAssetType(),
					Id:         info.GetId(),
					BundlePath: info.GetBundlePath(),
				},
			},
		}
	}
	m.Assets = assets

	if err := hardwaredevicebundle.WriteHardwareDeviceBundle(m, opts.OutputBundlePath); err != nil {
		return fmt.Errorf("failed to write HardwareDevice Asset bundle: %w", err)
	}

	return nil
}
