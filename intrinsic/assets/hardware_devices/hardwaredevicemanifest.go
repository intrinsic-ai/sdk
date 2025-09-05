// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevicemanifest contains tools for working with HardwareDeviceManifest.
package hardwaredevicemanifest

import (
	"fmt"
	"os"

	"intrinsic/assets/idutils"

	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	agpb "intrinsic/assets/proto/v1/asset_graph_go_proto"
	rpb "intrinsic/assets/proto/v1/reference_go_proto"
)

// HardwareDeviceManifest abstracts HardwareDeviceManifest and ProcessedHardwareDeviceManifest.
type HardwareDeviceManifest interface {
	GetMetadata() *hdmpb.HardwareDeviceMetadata
	GetGraph() *agpb.AssetGraph
}

// VerifyCatalogAssetsExist is a function that verifies that referenced catalog assets exist.
type VerifyCatalogAssetsExist func(assets []*rpb.CatalogAsset) error

// VerifyLocalAssetsExist is a function that verifies that referenced local assets exist.
type VerifyLocalAssetsExist func(assets []*rpb.LocalAsset) error

// VerifyLocalAssetsExistOnDisk verifies that the local asset exists on disk.
func VerifyLocalAssetsExistOnDisk(assets []*rpb.LocalAsset) error {
	for _, asset := range assets {
		if _, err := os.Stat(asset.GetBundlePath()); err != nil {
			return fmt.Errorf("asset %s has invalid bundle path %q: %w", idutils.IDFromProtoUnchecked(asset.GetId()), asset.GetBundlePath(), err)
		}
	}
	return nil
}

// validateHardwareDeviceManifestOptions contains options for a call to ValidateHardwareDeviceManifest.
type validateHardwareDeviceManifestOptions struct {
	// verifyCatalogAssetsExist is a function that verifies that referenced catalog assets exist.
	verifyCatalogAssetsExist VerifyCatalogAssetsExist
	// verifyLocalAssetsExist is a function that verifies that referenced local assets exist.
	verifyLocalAssetsExist VerifyLocalAssetsExist
}

// ValidateHardwareDeviceManifestOption is a functional option for ValidateHardwareDeviceManifest.
type ValidateHardwareDeviceManifestOption func(*validateHardwareDeviceManifestOptions)

// WithVerifyCatalogAssetsExist provides a function that verifies that referenced catalog asset
// exist.
func WithVerifyCatalogAssetsExist(f VerifyCatalogAssetsExist) ValidateHardwareDeviceManifestOption {
	return func(opts *validateHardwareDeviceManifestOptions) {
		opts.verifyCatalogAssetsExist = f
	}
}

// WithVerifyLocalAssetsExist provides a function that verifies that referenced local assets exist.
func WithVerifyLocalAssetsExist(f VerifyLocalAssetsExist) ValidateHardwareDeviceManifestOption {
	return func(opts *validateHardwareDeviceManifestOptions) {
		opts.verifyLocalAssetsExist = f
	}
}

// ValidateHardwareDeviceManifest validates the given HardwareDeviceManifest.
//
// The following validation cannot be done on reference nodes, since we don't read their metadata:
// - Verify that the specified asset type is actually what is stored in the catalog.
// - Verify that configuration edges have matching source and target nodes.
func ValidateHardwareDeviceManifest(hdm HardwareDeviceManifest, options ...ValidateHardwareDeviceManifestOption) error {
	if hdm == nil {
		return fmt.Errorf("HardwareDeviceManifest must not be nil")
	}
	opts := &validateHardwareDeviceManifestOptions{
		verifyLocalAssetsExist: VerifyLocalAssetsExistOnDisk,
	}
	for _, opt := range options {
		opt(opts)
	}

	// Validate the metadata.
	if err := idutils.ValidateIDProto(hdm.GetMetadata().GetId()); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if hdm.GetMetadata().GetDisplayName() == "" {
		return fmt.Errorf("display_name must be specified")
	}
	if hdm.GetMetadata().GetVendor().GetDisplayName() == "" {
		return fmt.Errorf("vendor.display_name must be specified")
	}

	// Validate individual assets.
	var assetInfoMap map[string]*assetInfo
	var err error
	switch hdm.(type) {
	case *hdmpb.HardwareDeviceManifest:
		assetInfoMap, err = validateAssets(hdm.(*hdmpb.HardwareDeviceManifest).GetAssets(), opts)
	case *hdmpb.ProcessedHardwareDeviceManifest:
		assetInfoMap, err = validateProcessedAssets(hdm.(*hdmpb.ProcessedHardwareDeviceManifest).GetAssets(), opts)
	default:
		return fmt.Errorf("HardwareDeviceManifest has unknown type %T", hdm)
	}
	if err != nil {
		return err
	}

	g := hdm.GetGraph()

	// Verify that node constraints are satisfied.
	numSceneObjects := 0
	numServices := 0
	node2AssetInfo := make(map[string]*assetInfo)
	referencedAssets := make(map[string]struct{})
	for name, node := range g.GetNodes() {
		info, ok := assetInfoMap[node.GetAsset()]
		if !ok {
			return fmt.Errorf("node %q refers to %s, which is not a specified asset", name, node.GetAsset())
		}

		switch info.assetType {
		case atpb.AssetType_ASSET_TYPE_SCENE_OBJECT:
			numSceneObjects++
		case atpb.AssetType_ASSET_TYPE_SERVICE:
			numServices++
		case atpb.AssetType_ASSET_TYPE_DATA:
			// OK.
		default:
			return fmt.Errorf("HardwareDevice node %q has unsupported asset type %v", name, info.assetType)
		}

		node2AssetInfo[name] = info
		referencedAssets[node.GetAsset()] = struct{}{}
	}
	if numSceneObjects != 1 {
		return fmt.Errorf("HardwareDevice must contain exactly one SceneObject, found %d", numSceneObjects)
	}
	if numServices == 0 {
		return fmt.Errorf("HardwareDevice must contain at least one Service")
	}
	if numServices > 1 {
		return fmt.Errorf("HardwareDevice must contain at most one Service, found %d", numServices)
	}

	// Verify that all assets were referenced.
	for key := range assetInfoMap {
		if _, ok := referencedAssets[key]; !ok {
			return fmt.Errorf("asset %q is not referenced in the graph", key)
		}
	}

	// Verify that edge constraints are satisfied.
	dataConfigNodes := make(map[string]struct{})
	for _, edge := range g.GetEdges() {
		sourceNodeInfo, ok := node2AssetInfo[edge.GetSource()]
		if !ok {
			return fmt.Errorf("HardwareDevice edge %q -> %q has unknown source node %q", edge.GetSource(), edge.GetTarget(), edge.GetSource())
		}
		targetNodeInfo, ok := node2AssetInfo[edge.GetTarget()]
		if !ok {
			return fmt.Errorf("HardwareDevice edge %q -> %q has unknown target node %q", edge.GetSource(), edge.GetTarget(), edge.GetTarget())
		}

		sourceNodeType := sourceNodeInfo.assetType
		targetNodeType := targetNodeInfo.assetType

		switch edge.GetEdgeType().(type) {
		case *agpb.AssetEdge_Configures:
			if sourceNodeType != atpb.AssetType_ASSET_TYPE_DATA {
				return fmt.Errorf("HardwareDevice edge %q -> %q cannot be configuration (source node must be Data, got %v)", edge.GetSource(), edge.GetTarget(), sourceNodeType)
			}
			if targetNodeType != atpb.AssetType_ASSET_TYPE_SERVICE {
				return fmt.Errorf("HardwareDevice edge %q -> %q cannot be configuration (target node must be a Service, got %v)", edge.GetSource(), edge.GetTarget(), targetNodeType)
			}
		default:
			return fmt.Errorf("HardwareDevice edge %q -> %q has unsupported edge type %v", edge.GetSource(), edge.GetTarget(), edge.GetEdgeType())
		}

		dataConfigNodes[edge.GetSource()] = struct{}{}
	}

	// Verify that all Data assets were used in a configuration edge. We may eventually want to
	// allow Data assets that are not used in a configuration edge (e.g., to enable the HardwareDevice
	// to provide some kind of data to other assets). But we need this constraint for now, because the
	// ProcessedHardwareDeviceManifest will eventually need to be converted into a
	// ProcessedResourceManifest, which cannot represent data.
	for name, node := range g.GetNodes() {
		info, ok := assetInfoMap[node.GetAsset()]
		if !ok {
			return fmt.Errorf("node %q refers to %s, which is not a specified asset (while validating data nodes)", name, node.GetAsset())
		}
		if info.assetType == atpb.AssetType_ASSET_TYPE_DATA {
			if _, ok := dataConfigNodes[name]; !ok {
				return fmt.Errorf("Data asset node %q is not used in a configuration edge", name)
			}
		}
	}

	return nil
}

type assetInfo struct {
	assetType atpb.AssetType
}

func validateAssets(assets map[string]*hdmpb.HardwareDeviceManifest_Asset, opts *validateHardwareDeviceManifestOptions) (map[string]*assetInfo, error) {
	assetInfoMap := make(map[string]*assetInfo)
	var catalogAssets []*rpb.CatalogAsset
	var localAssets []*rpb.LocalAsset
	for key, asset := range assets {
		var assetType atpb.AssetType
		var id *idpb.Id
		switch asset.Variant.(type) {
		case *hdmpb.HardwareDeviceManifest_Asset_Catalog:
			assetType = asset.GetCatalog().GetAssetType()
			idVersion := asset.GetCatalog().GetIdVersion()
			id = idVersion.GetId()
			if idVersion.GetVersion() == "" {
				return nil, fmt.Errorf("asset %q has no version", idutils.IDFromProtoUnchecked(id))
			}
			catalogAssets = append(catalogAssets, asset.GetCatalog())
		case *hdmpb.HardwareDeviceManifest_Asset_Local:
			assetType = asset.GetLocal().GetAssetType()
			id = asset.GetLocal().GetId()
			bundlePath := asset.GetLocal().GetBundlePath()
			if bundlePath == "" {
				return nil, fmt.Errorf("asset %q has no bundle path", idutils.IDFromProtoUnchecked(id))
			}
			localAssets = append(localAssets, asset.GetLocal())
		default:
			return nil, fmt.Errorf("asset has unknown variant %T", asset.Variant)
		}

		// By convention, we enforce that the key is the asset ID.
		idString := idutils.IDFromProtoUnchecked(id)
		if idString != key {
			return nil, fmt.Errorf("asset ID %s does not match key %q", idString, key)
		}

		assetInfoMap[key] = &assetInfo{
			assetType: assetType,
		}
	}

	// Check for the existence of catalog assets.
	if opts.verifyCatalogAssetsExist != nil {
		if err := opts.verifyCatalogAssetsExist(catalogAssets); err != nil {
			return nil, err
		}
	}
	// Check for the existence of local assets.
	if opts.verifyLocalAssetsExist != nil {
		if err := opts.verifyLocalAssetsExist(localAssets); err != nil {
			return nil, err
		}
	}

	return assetInfoMap, nil
}

func validateProcessedAssets(assets map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, opts *validateHardwareDeviceManifestOptions) (map[string]*assetInfo, error) {
	assetInfoMap := make(map[string]*assetInfo)
	var catalogAssets []*rpb.CatalogAsset
	for key, asset := range assets {
		var assetType atpb.AssetType
		var id *idpb.Id
		switch asset.Variant.(type) {
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog:
			assetType = asset.GetCatalog().GetAssetType()
			idVersion := asset.GetCatalog().GetIdVersion()
			id = idVersion.GetId()
			if idVersion.GetVersion() == "" {
				return nil, fmt.Errorf("asset %q has no version", idutils.IDFromProtoUnchecked(id))
			}
			catalogAssets = append(catalogAssets, asset.GetCatalog())
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service:
			assetType = atpb.AssetType_ASSET_TYPE_SERVICE
			id = asset.GetService().GetMetadata().GetId()
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject:
			assetType = atpb.AssetType_ASSET_TYPE_SCENE_OBJECT
			id = asset.GetSceneObject().GetMetadata().GetId()
		case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Data:
			assetType = atpb.AssetType_ASSET_TYPE_DATA
			id = asset.GetData().GetMetadata().GetIdVersion().GetId()
		default:
			return nil, fmt.Errorf("asset has unknown variant %T", asset.Variant)
		}

		// By convention, we enforce that the key is the asset ID.
		idString := idutils.IDFromProtoUnchecked(id)
		if idString != key {
			return nil, fmt.Errorf("asset ID %s does not match key %q", idString, key)
		}

		assetInfoMap[key] = &assetInfo{
			assetType: assetType,
		}
	}

	// Check for the existence of catalog assets.
	if opts.verifyCatalogAssetsExist != nil {
		if err := opts.verifyCatalogAssetsExist(catalogAssets); err != nil {
			return nil, err
		}
	}

	return assetInfoMap, nil
}
