// Copyright 2023 Intrinsic Innovation LLC

// Package hardwaredevicevalidate provides utils for validating HardwareDevices.
package hardwaredevicevalidate

import (
	"fmt"
	"os"

	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	agpb "intrinsic/assets/proto/v1/asset_graph_go_proto"
	rpb "intrinsic/assets/proto/v1/reference_go_proto"
)

// VerifyCatalogAssetsExist is a function that verifies that referenced catalog Assets exist.
type VerifyCatalogAssetsExist func(assets []*rpb.CatalogAsset) error

// VerifyLocalAssetsExist is a function that verifies that referenced local Assets exist.
type VerifyLocalAssetsExist func(assets []*rpb.LocalAsset) error

// VerifyLocalAssetsExistOnDisk verifies that local Assets exist on disk.
func VerifyLocalAssetsExistOnDisk(assets []*rpb.LocalAsset) error {
	for _, asset := range assets {
		if _, err := os.Stat(asset.GetBundlePath()); err != nil {
			return fmt.Errorf("Asset %s has invalid bundle path %q: %w", idutils.IDFromProtoUnchecked(asset.GetId()), asset.GetBundlePath(), err)
		}
	}
	return nil
}

type hardwareDeviceManifestOptions struct {
	verifyCatalogAssetsExist VerifyCatalogAssetsExist
	verifyLocalAssetsExist   VerifyLocalAssetsExist
}

// HardwareDeviceManifestOption is an option for validating a HardwareDeviceManifest.
type HardwareDeviceManifestOption func(*hardwareDeviceManifestOptions)

// WithVerifyCatalogAssetsExist provides a function that verifies that referenced catalog Assets
// exist.
func WithVerifyCatalogAssetsExist(f VerifyCatalogAssetsExist) HardwareDeviceManifestOption {
	return func(opts *hardwareDeviceManifestOptions) {
		opts.verifyCatalogAssetsExist = f
	}
}

// WithVerifyLocalAssetsExist provides a function that verifies that referenced local Assets exist.
func WithVerifyLocalAssetsExist(f VerifyLocalAssetsExist) HardwareDeviceManifestOption {
	return func(opts *hardwareDeviceManifestOptions) {
		opts.verifyLocalAssetsExist = f
	}
}

// HardwareDeviceManifest validates a HardwareDeviceManifest.
//
// The following validation cannot be done on catalog reference nodes, since we don't read their
// metadata:
// - Verify that the specified Asset type is actually what is stored in the catalog.
// - Verify that configuration edges have matching source and target nodes.
func HardwareDeviceManifest(m *hdmpb.HardwareDeviceManifest, options ...HardwareDeviceManifestOption) error {
	opts := &hardwareDeviceManifestOptions{
		verifyLocalAssetsExist: VerifyLocalAssetsExistOnDisk,
	}
	for _, opt := range options {
		opt(opts)
	}

	if m == nil {
		return fmt.Errorf("HardwareDeviceManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid HardwareDeviceManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	// Validate individual Assets.
	assetInfoMap, err := validateAssets(m.GetAssets(), opts)
	if err != nil {
		return fmt.Errorf("invalid Assets for %q: %w", id, err)
	}

	if err := validateGraph(m.GetGraph(), assetInfoMap); err != nil {
		return fmt.Errorf("invalid graph for %q: %w", id, err)
	}

	return nil
}

type processedHardwareDeviceManifestOptions struct {
	verifyCatalogAssetsExist VerifyCatalogAssetsExist
}

// ProcessedHardwareDeviceManifestOption is an option for validating a ProcessedHardwareDeviceManifest.
type ProcessedHardwareDeviceManifestOption func(*processedHardwareDeviceManifestOptions)

// WithVerifyProcessedCatalogAssetsExist provides a function that verifies that referenced catalog
// Assets exist.
func WithVerifyProcessedCatalogAssetsExist(f VerifyCatalogAssetsExist) ProcessedHardwareDeviceManifestOption {
	return func(opts *processedHardwareDeviceManifestOptions) {
		opts.verifyCatalogAssetsExist = f
	}
}

// ProcessedHardwareDeviceManifest validates a ProcessedHardwareDeviceManifest.
//
// The following validation cannot be done on catalog reference nodes, since we don't read their
// metadata:
// - Verify that the specified Asset type is actually what is stored in the catalog.
// - Verify that configuration edges have matching source and target nodes.
func ProcessedHardwareDeviceManifest(pm *hdmpb.ProcessedHardwareDeviceManifest, options ...ProcessedHardwareDeviceManifestOption) error {
	opts := &processedHardwareDeviceManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if pm == nil {
		return fmt.Errorf("ProcessedHardwareDeviceManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(pm.GetMetadata()); err != nil {
		return fmt.Errorf("invalid ProcessedHardwareDeviceManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(pm.GetMetadata().GetId())

	// Validate individual Assets.
	assetInfoMap, err := validateProcessedAssets(pm.GetAssets(), opts)
	if err != nil {
		return fmt.Errorf("invalid processed Assets for %q: %w", id, err)
	}

	if err := validateGraph(pm.GetGraph(), assetInfoMap); err != nil {
		return fmt.Errorf("invalid processed graph for %q: %w", id, err)
	}

	return nil
}

type assetInfo struct {
	assetType atpb.AssetType
}

func validateGraph(g *agpb.AssetGraph, ai map[string]*assetInfo) error {
	// Verify that node constraints are satisfied.
	numSceneObjects := 0
	numServices := 0
	node2AssetInfo := make(map[string]*assetInfo)
	referencedAssets := make(map[string]struct{})
	for name, node := range g.GetNodes() {
		info, ok := ai[node.GetAsset()]
		if !ok {
			return fmt.Errorf("node %q refers to %q, which is not a specified Asset", name, node.GetAsset())
		}

		switch info.assetType {
		case atpb.AssetType_ASSET_TYPE_SCENE_OBJECT:
			numSceneObjects++
		case atpb.AssetType_ASSET_TYPE_SERVICE:
			numServices++
		case atpb.AssetType_ASSET_TYPE_DATA:
			// OK.
		default:
			return fmt.Errorf("node %q has unsupported Asset type %v", name, info.assetType)
		}

		node2AssetInfo[name] = info
		referencedAssets[node.GetAsset()] = struct{}{}
	}
	if numSceneObjects != 1 {
		return fmt.Errorf("must contain exactly one SceneObject, found %d", numSceneObjects)
	}
	if numServices == 0 {
		return fmt.Errorf("must contain at least one Service")
	}
	if numServices > 1 {
		return fmt.Errorf("must contain at most one Service, found %d", numServices)
	}

	// Verify that all Assets were referenced.
	for key := range ai {
		if _, ok := referencedAssets[key]; !ok {
			return fmt.Errorf("Asset %q is not referenced in the graph", key)
		}
	}

	// Verify that edge constraints are satisfied.
	dataConfigNodes := make(map[string]struct{})
	for _, edge := range g.GetEdges() {
		sourceNodeInfo, ok := node2AssetInfo[edge.GetSource()]
		if !ok {
			return fmt.Errorf("edge %q -> %q has unknown source node %q", edge.GetSource(), edge.GetTarget(), edge.GetSource())
		}
		targetNodeInfo, ok := node2AssetInfo[edge.GetTarget()]
		if !ok {
			return fmt.Errorf("edge %q -> %q has unknown target node %q", edge.GetSource(), edge.GetTarget(), edge.GetTarget())
		}

		sourceNodeType := sourceNodeInfo.assetType
		targetNodeType := targetNodeInfo.assetType

		switch edge.GetEdgeType().(type) {
		case *agpb.AssetEdge_Configures:
			if sourceNodeType != atpb.AssetType_ASSET_TYPE_DATA {
				return fmt.Errorf("edge %q -> %q cannot be configuration (source node must be Data, got %v)", edge.GetSource(), edge.GetTarget(), sourceNodeType)
			}
			if targetNodeType != atpb.AssetType_ASSET_TYPE_SERVICE {
				return fmt.Errorf("edge %q -> %q cannot be configuration (target node must be a Service, got %v)", edge.GetSource(), edge.GetTarget(), targetNodeType)
			}
		default:
			return fmt.Errorf("edge %q -> %q has unsupported edge type %v", edge.GetSource(), edge.GetTarget(), edge.GetEdgeType())
		}

		dataConfigNodes[edge.GetSource()] = struct{}{}
	}

	// Verify that all Data Assets were used in a configuration edge. We may eventually want to allow
	// Data Assets that are not used in a configuration edge (e.g., to enable the HardwareDevice to
	// provide some kind of data to other Assets). But we need this constraint for now, because the
	// ProcessedHardwareDeviceManifest will eventually need to be converted into a
	// ProcessedResourceManifest, which cannot represent data.
	for name, node := range g.GetNodes() {
		info, ok := ai[node.GetAsset()]
		if !ok {
			return fmt.Errorf("node %q refers to %s, which is not a specified Asset (while validating data nodes)", name, node.GetAsset())
		}
		if info.assetType == atpb.AssetType_ASSET_TYPE_DATA {
			if _, ok := dataConfigNodes[name]; !ok {
				return fmt.Errorf("Data Asset node %q is not used in a configuration edge", name)
			}
		}
	}

	return nil
}

func validateAssets(assets map[string]*hdmpb.HardwareDeviceManifest_Asset, opts *hardwareDeviceManifestOptions) (map[string]*assetInfo, error) {
	ai := make(map[string]*assetInfo)
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
				return nil, fmt.Errorf("Asset %q has no version", idutils.IDFromProtoUnchecked(id))
			}
			catalogAssets = append(catalogAssets, asset.GetCatalog())
		case *hdmpb.HardwareDeviceManifest_Asset_Local:
			assetType = asset.GetLocal().GetAssetType()
			id = asset.GetLocal().GetId()
			bundlePath := asset.GetLocal().GetBundlePath()
			if bundlePath == "" {
				return nil, fmt.Errorf("Asset %q has no bundle path", idutils.IDFromProtoUnchecked(id))
			}
			localAssets = append(localAssets, asset.GetLocal())
		default:
			return nil, fmt.Errorf("Asset has unknown variant %T", asset.Variant)
		}

		// By convention, we enforce that the key is the asset ID.
		idString := idutils.IDFromProtoUnchecked(id)
		if idString != key {
			return nil, fmt.Errorf("Asset ID %s does not match key %q", idString, key)
		}

		ai[key] = &assetInfo{
			assetType: assetType,
		}
	}

	// Check for the existence of catalog assets.
	if opts.verifyCatalogAssetsExist != nil {
		if err := opts.verifyCatalogAssetsExist(catalogAssets); err != nil {
			return nil, fmt.Errorf("failed to verify catalog Assets: %w", err)
		}
	}
	// Check for the existence of local assets.
	if opts.verifyLocalAssetsExist != nil {
		if err := opts.verifyLocalAssetsExist(localAssets); err != nil {
			return nil, fmt.Errorf("failed to verify local Assets: %w", err)
		}
	}

	return ai, nil
}

func validateProcessedAssets(assets map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset, opts *processedHardwareDeviceManifestOptions) (map[string]*assetInfo, error) {
	ai := make(map[string]*assetInfo)
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
				return nil, fmt.Errorf("Asset %q has no version", idutils.IDFromProtoUnchecked(id))
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
			return nil, fmt.Errorf("Asset has unknown variant %T", asset.Variant)
		}

		// By convention, we enforce that the key is the Asset ID.
		idString := idutils.IDFromProtoUnchecked(id)
		if idString != key {
			return nil, fmt.Errorf("Asset ID %s does not match key %q", idString, key)
		}

		ai[key] = &assetInfo{
			assetType: assetType,
		}
	}

	// Check for the existence of catalog Assets.
	if opts.verifyCatalogAssetsExist != nil {
		if err := opts.verifyCatalogAssetsExist(catalogAssets); err != nil {
			return nil, fmt.Errorf("failed to verify catalog Assets: %w", err)
		}
	}

	return ai, nil
}
