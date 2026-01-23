// Copyright 2023 Intrinsic Innovation LLC

// Package utils contains utilities for testing HardwareDevices.
package utils

import (
	"context"
	"log"
	"testing"

	datatestutils "intrinsic/assets/data/testing/utils"
	"intrinsic/assets/idutils"
	sceneobjecttestutils "intrinsic/assets/scene_objects/testing/utils"
	servicetestutils "intrinsic/assets/services/testing/utils"
	"intrinsic/util/proto/names"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	hdmpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	agpb "intrinsic/assets/proto/v1/asset_graph_go_proto"
	rpb "intrinsic/assets/proto/v1/reference_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
)

// GraphSpecVariant specifies a specific variant of a graph spec.
type GraphSpecVariant int

const (
	// GraphSpecVariantLocal has only local assets.
	GraphSpecVariantLocal GraphSpecVariant = iota
	// GraphSpecVariantCatalog has only catalog assets.
	GraphSpecVariantCatalog
	// GraphSpecVariantLocalAndCatalog has both local and catalog assets.
	GraphSpecVariantLocalAndCatalog
	// GraphSpecVariantDataConfiguresService has a Data asset that configures a Service asset.
	GraphSpecVariantDataConfiguresService
)

// GraphSpec is a spec specifying assets and a graph for a test ProcessedHardwareDeviceManifest.
type GraphSpec struct {
	Assets map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset
	Graph  *agpb.AssetGraph
}

type makeHardwareDeviceManifestOptions struct {
	assets   map[string]*hdmpb.HardwareDeviceManifest_Asset
	graph    *agpb.AssetGraph
	metadata *hdmpb.HardwareDeviceMetadata
}

// MakeHardwareDeviceManifestOption is an option for MakeHardwareDeviceManifest.
type MakeHardwareDeviceManifestOption func(*testing.T, *makeHardwareDeviceManifestOptions)

// WithAssets specifies the Assets to use in the HardwareDeviceManifest.
func WithAssets(assets map[string]*hdmpb.HardwareDeviceManifest_Asset) MakeHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeHardwareDeviceManifestOptions) {
		opts.assets = assets
	}
}

// WithGraph specifies the graph to use in the HardwareDeviceManifest.
func WithGraph(graph *agpb.AssetGraph) MakeHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeHardwareDeviceManifestOptions) {
		opts.graph = graph
	}
}

// WithGraphSpecVariant specifies the assets and graph of the HardwareDeviceManifest based on a
// preset variant.
func WithGraphSpecVariant(variant GraphSpecVariant, options ...MakeGraphSpecVariantOption) MakeHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeHardwareDeviceManifestOptions) {
		t.Helper()

		spec := MakeGraphSpecVariant(t, variant, options...)
		opts.assets = AssetsFromProcessedAssets(t, spec.Assets)
		opts.graph = spec.Graph
	}
}

// WithMetadata specifies the metadata to use in the HardwareDeviceManifest.
func WithMetadata(metadata *hdmpb.HardwareDeviceMetadata) MakeHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeHardwareDeviceManifestOptions) {
		opts.metadata = metadata
	}
}

// MakeHardwareDeviceManifest makes a HardwareDeviceManifest for testing.
func MakeHardwareDeviceManifest(t *testing.T, options ...MakeHardwareDeviceManifestOption) *hdmpb.HardwareDeviceManifest {
	t.Helper()

	opts := &makeHardwareDeviceManifestOptions{
		metadata: &hdmpb.HardwareDeviceMetadata{
			Id: &idpb.Id{
				Name:    "some_hardware_device",
				Package: "package.some",
			},
			DisplayName: "Some HardwareDevice",
			Vendor: &vpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	WithGraphSpecVariant(GraphSpecVariantCatalog)(t, opts)
	for _, opt := range options {
		opt(t, opts)
	}

	return &hdmpb.HardwareDeviceManifest{
		Metadata: opts.metadata,
		Assets:   opts.assets,
		Graph:    opts.graph,
	}
}

type makeProcessedHardwareDeviceManifestOptions struct {
	assets   map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset
	graph    *agpb.AssetGraph
	metadata *hdmpb.HardwareDeviceMetadata
}

// MakeProcessedHardwareDeviceManifestOption is an option for MakeProcessedHardwareDeviceManifest.
type MakeProcessedHardwareDeviceManifestOption func(*testing.T, *makeProcessedHardwareDeviceManifestOptions)

// WithProcessedAssets specifies the Assets to use in the ProcessedHardwareDeviceManifest.
func WithProcessedAssets(assets map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset) MakeProcessedHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeProcessedHardwareDeviceManifestOptions) {
		opts.assets = assets
	}
}

// WithProcessedGraph specifies the graph to use in the ProcessedHardwareDeviceManifest.
func WithProcessedGraph(graph *agpb.AssetGraph) MakeProcessedHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeProcessedHardwareDeviceManifestOptions) {
		opts.graph = graph
	}
}

// WithProcessedGraphSpecVariant specifies the assets and graph of the
// ProcessedHardwareDeviceManifest based on a preset variant.
func WithProcessedGraphSpecVariant(variant GraphSpecVariant, options ...MakeGraphSpecVariantOption) MakeProcessedHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeProcessedHardwareDeviceManifestOptions) {
		t.Helper()

		spec := MakeGraphSpecVariant(t, variant, options...)
		opts.assets = spec.Assets
		opts.graph = spec.Graph
	}
}

// WithProcessedMetadata specifies the metadata to use in the ProcessedHardwareDeviceManifest.
func WithProcessedMetadata(metadata *hdmpb.HardwareDeviceMetadata) MakeProcessedHardwareDeviceManifestOption {
	return func(t *testing.T, opts *makeProcessedHardwareDeviceManifestOptions) {
		opts.metadata = metadata
	}
}

// MakeProcessedHardwareDeviceManifest makes a ProcessedHardwareDeviceManifest for testing.
func MakeProcessedHardwareDeviceManifest(t *testing.T, options ...MakeProcessedHardwareDeviceManifestOption) *hdmpb.ProcessedHardwareDeviceManifest {
	t.Helper()

	opts := &makeProcessedHardwareDeviceManifestOptions{
		metadata: &hdmpb.HardwareDeviceMetadata{
			Id: &idpb.Id{
				Name:    "some_hardware_device",
				Package: "package.some",
			},
			DisplayName: "Some HardwareDevice",
			Vendor: &vpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	WithProcessedGraphSpecVariant(GraphSpecVariantLocal)(t, opts)
	for _, opt := range options {
		opt(t, opts)
	}

	return &hdmpb.ProcessedHardwareDeviceManifest{
		Metadata: opts.metadata,
		Assets:   opts.assets,
		Graph:    opts.graph,
	}
}

type makeGraphSpecVariantOptions struct {
	assetCatalogClient acgrpcpb.AssetCatalogClient
	context            context.Context
	data               *dapb.DataAsset
	sceneObject        *sompb.ProcessedSceneObjectManifest
	service            *smpb.ProcessedServiceManifest
	version            string
}

// MakeGraphSpecVariantOption is an option for MakeGraphSpecVariant.
type MakeGraphSpecVariantOption func(*makeGraphSpecVariantOptions)

// WithAssetCatalogClient sets the AssetCatalogClient to use for ensuring that any catalog
// references in the returned variant are valid.
func WithAssetCatalogClient(ctx context.Context, client acgrpcpb.AssetCatalogClient) MakeGraphSpecVariantOption {
	return func(opts *makeGraphSpecVariantOptions) {
		opts.assetCatalogClient = client
		opts.context = ctx
	}
}

// WithGraphSpecData specifies the DataAsset to use in the graph.
func WithGraphSpecData(data *dapb.DataAsset) MakeGraphSpecVariantOption {
	return func(opts *makeGraphSpecVariantOptions) {
		opts.data = data
	}
}

// WithGraphSpecSceneObject specifies the ProcessedSceneObjectManifest to use in the graph.
func WithGraphSpecSceneObject(sceneObject *sompb.ProcessedSceneObjectManifest) MakeGraphSpecVariantOption {
	return func(opts *makeGraphSpecVariantOptions) {
		opts.sceneObject = sceneObject
	}
}

// WithGraphSpecService specifies the ProcessedServiceManifest to use in the graph.
func WithGraphSpecService(service *smpb.ProcessedServiceManifest) MakeGraphSpecVariantOption {
	return func(opts *makeGraphSpecVariantOptions) {
		opts.service = service
	}
}

// WithGraphSpecVersion specifies the version to use for catalog Assets in the graph.
func WithGraphSpecVersion(version string) MakeGraphSpecVariantOption {
	return func(opts *makeGraphSpecVariantOptions) {
		opts.version = version
	}
}

// MakeGraphSpecVariant makes a graph spec for testing.
func MakeGraphSpecVariant(t *testing.T, variant GraphSpecVariant, options ...MakeGraphSpecVariantOption) *GraphSpec {
	t.Helper()

	opts := &makeGraphSpecVariantOptions{
		context: context.Background(),
		version: "1.0.0",
	}
	for _, opt := range options {
		opt(opts)
	}
	if opts.data == nil {
		opts.data = datatestutils.MakeDataAsset(t)
	}
	if opts.sceneObject == nil {
		opts.sceneObject = sceneobjecttestutils.MakeProcessedSceneObjectManifest(t)
	}
	if opts.service == nil {
		name, err := names.AnyToProtoName(opts.data.GetData())
		if err != nil {
			t.Fatalf("could not get name from Any: %v", err)
		}
		opts.service = servicetestutils.MakeProcessedServiceManifest(t,
			servicetestutils.WithProcessedConfigMessageFullName(name),
			servicetestutils.WithFileDescriptorSet(opts.data.GetFileDescriptorSet()),
		)
	}

	catalogSceneObject := &rpb.CatalogAsset{
		AssetType: atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT,
		IdVersion: &idpb.IdVersion{
			Id:      opts.sceneObject.GetMetadata().GetId(),
			Version: opts.version,
		},
	}
	catalogSceneObjectAsset := &acpb.Asset{
		Metadata: &mpb.Metadata{
			AssetType:   catalogSceneObject.GetAssetType(),
			DisplayName: opts.sceneObject.GetMetadata().GetDisplayName(),
			IdVersion:   catalogSceneObject.GetIdVersion(),
			Vendor:      opts.sceneObject.GetMetadata().GetVendor(),
		},
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_SceneObjectSpecificDeploymentData{
				SceneObjectSpecificDeploymentData: &acpb.Asset_SceneObjectDeploymentData{
					Manifest: opts.sceneObject,
				},
			},
		},
	}
	catalogService := &rpb.CatalogAsset{
		AssetType: atypepb.AssetType_ASSET_TYPE_SERVICE,
		IdVersion: &idpb.IdVersion{
			Id:      opts.service.GetMetadata().GetId(),
			Version: opts.version,
		},
	}
	catalogServiceAsset := &acpb.Asset{
		Metadata: &mpb.Metadata{
			AssetType:   catalogService.GetAssetType(),
			IdVersion:   catalogService.GetIdVersion(),
			DisplayName: opts.service.GetMetadata().GetDisplayName(),
			Vendor:      opts.service.GetMetadata().GetVendor(),
		},
		DeploymentData: &acpb.Asset_AssetDeploymentData{
			AssetSpecificDeploymentData: &acpb.Asset_AssetDeploymentData_ServiceSpecificDeploymentData{
				ServiceSpecificDeploymentData: &acpb.Asset_ServiceDeploymentData{
					Manifest: opts.service,
				},
			},
		},
	}
	sceneObjectID := idutils.IDFromProtoUnchecked(opts.sceneObject.GetMetadata().GetId())
	serviceID := idutils.IDFromProtoUnchecked(opts.service.GetMetadata().GetId())
	dataID := idutils.IDFromProtoUnchecked(opts.data.GetMetadata().GetIdVersion().GetId())

	var catalogAssets []*acpb.Asset

	spec := &GraphSpec{
		Assets: map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{},
		Graph:  &agpb.AssetGraph{},
	}
	switch variant {
	case GraphSpecVariantLocal:
		spec.Assets[sceneObjectID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject{
				SceneObject: opts.sceneObject,
			},
		}
		spec.Assets[serviceID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service{
				Service: opts.service,
			},
		}
		spec.Graph = &agpb.AssetGraph{
			Nodes: map[string]*agpb.AssetNode{
				"a_scene_object": {Asset: sceneObjectID},
				"a_service":      {Asset: serviceID},
			},
		}
	case GraphSpecVariantCatalog:
		spec.Assets[sceneObjectID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog{
				Catalog: catalogSceneObject,
			},
		}
		spec.Assets[serviceID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog{
				Catalog: catalogService,
			},
		}
		spec.Graph = &agpb.AssetGraph{
			Nodes: map[string]*agpb.AssetNode{
				"a_scene_object": {Asset: sceneObjectID},
				"a_service":      {Asset: serviceID},
			},
		}
		catalogAssets = append(catalogAssets, catalogSceneObjectAsset, catalogServiceAsset)
	case GraphSpecVariantLocalAndCatalog:
		spec.Assets[sceneObjectID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog{
				Catalog: catalogSceneObject,
			},
		}
		spec.Assets[serviceID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service{
				Service: opts.service,
			},
		}
		spec.Graph = &agpb.AssetGraph{
			Nodes: map[string]*agpb.AssetNode{
				"a_scene_object": {Asset: sceneObjectID},
				"a_service":      {Asset: serviceID},
			},
		}
		catalogAssets = append(catalogAssets, catalogSceneObjectAsset)
	case GraphSpecVariantDataConfiguresService:
		spec.Assets[sceneObjectID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject{
				SceneObject: opts.sceneObject,
			},
		}
		spec.Assets[serviceID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service{
				Service: opts.service,
			},
		}
		spec.Assets[dataID] = &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset{
			Variant: &hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Data{
				Data: opts.data,
			},
		}
		spec.Graph = &agpb.AssetGraph{
			Nodes: map[string]*agpb.AssetNode{
				"a_scene_object": {Asset: sceneObjectID},
				"a_service":      {Asset: serviceID},
				"a_data":         {Asset: dataID},
			},
			Edges: []*agpb.AssetEdge{
				{
					Source: "a_data",
					Target: "a_service",
					EdgeType: &agpb.AssetEdge_Configures{
						Configures: &agpb.AssetEdge_Configuration{},
					},
				},
			},
		}
	default:
		t.Fatalf("MakeGraphSpecVariant received unsupported variant %v", variant)
	}

	if opts.assetCatalogClient != nil {
		for _, asset := range catalogAssets {
			if _, err := opts.assetCatalogClient.CreateAsset(opts.context, &acpb.CreateAssetRequest{Asset: asset}); err != nil {
				// Ignore the error if the asset already exists.
				if s, ok := status.FromError(err); !ok || s.Code() != codes.AlreadyExists {
					t.Fatalf("Failed to create asset %q: %v", idutils.IDFromProtoUnchecked(asset.GetMetadata().GetIdVersion().GetId()), err)
				}
			}
		}
	}

	return spec
}

// AssetsFromProcessedAssets devices a map of *hdmpb.HardwareDeviceManifest_Asset from a map of
// *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset.
func AssetsFromProcessedAssets(t *testing.T, pas map[string]*hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset) map[string]*hdmpb.HardwareDeviceManifest_Asset {
	assets := map[string]*hdmpb.HardwareDeviceManifest_Asset{}
	for key, pa := range pas {
		assets[key] = assetFromProcessedAsset(t, pa)
	}

	return assets
}

func assetFromProcessedAsset(t *testing.T, pa *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset) *hdmpb.HardwareDeviceManifest_Asset {
	asset := &hdmpb.HardwareDeviceManifest_Asset{}
	switch pa.GetVariant().(type) {
	case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Service:
		asset.Variant = &hdmpb.HardwareDeviceManifest_Asset_Local{
			Local: &rpb.LocalAsset{
				AssetType:  atypepb.AssetType_ASSET_TYPE_SERVICE,
				Id:         pa.GetService().GetMetadata().GetId(),
				BundlePath: "some_service_bundle_path.tar",
			},
		}
	case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Catalog:
		asset.Variant = &hdmpb.HardwareDeviceManifest_Asset_Catalog{
			Catalog: pa.GetCatalog(),
		}
	case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_Data:
		asset.Variant = &hdmpb.HardwareDeviceManifest_Asset_Local{
			Local: &rpb.LocalAsset{
				AssetType:  atypepb.AssetType_ASSET_TYPE_DATA,
				Id:         pa.GetData().GetMetadata().GetIdVersion().GetId(),
				BundlePath: "some_data_bundle_path.tar",
			},
		}
	case *hdmpb.ProcessedHardwareDeviceManifest_ProcessedAsset_SceneObject:
		asset.Variant = &hdmpb.HardwareDeviceManifest_Asset_Local{
			Local: &rpb.LocalAsset{
				AssetType:  atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT,
				Id:         pa.GetSceneObject().GetMetadata().GetId(),
				BundlePath: "some_scene_object_bundle_path.tar",
			},
		}
	default:
		log.Fatalf("unknown asset type in HardwareDevice: %v", pa.GetVariant())
	}

	return asset
}
