// Copyright 2023 Intrinsic Innovation LLC

// Package deviceservice implements a bare minimal DeviceService rpc service.
package deviceservice

import (
	"context"
	"errors"
	"fmt"
	"intrinsic/assets/dependencies/utils"
	"path"
	"strings"

	log "github.com/golang/glog"

	dagrpcpb "intrinsic/assets/data/proto/v1/data_assets_go_grpc_proto"
	dscpb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_config_go_proto"
	dspb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_go_grpc_proto"
	esipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/esi_go_proto"
)

var (
	// ErrConfigNil indicates that the configuration is nil.
	ErrConfigNil = errors.New("config must not be nil")
	// ErrNoEsiBundleID indicates that the configuration has no ESI bundle ID.
	ErrNoEsiBundleID = errors.New("config must have esi_bundle_data_asset_id")
	// ErrNoDeviceIdentifier indicates that the configuration has no device identifier.
	ErrNoDeviceIdentifier = errors.New("config must have a device identifier")
	// ErrEsiListMetadata indicates that listing ESI data asset metadata failed.
	ErrEsiListMetadata = errors.New("listing ESI data asset metadata")
	// ErrNoEsiDataAssets indicates that no ESI data assets were found.
	ErrNoEsiDataAssets = errors.New("no ESI data assets found")
	// ErrEsiBundleNotFound indicates that one or more requested ESI data assets were not found.
	ErrEsiBundleNotFound = errors.New("ESI data asset(s) not found")
	// ErrEsiGet indicates that getting an ESI data asset failed.
	ErrEsiGet = errors.New("getting ESI data asset")
	// ErrEsiUnmarshal indicates that unmarshalling an ESI data asset failed.
	ErrEsiUnmarshal = errors.New("unmarshalling ESI data asset")
	// ErrEsiInvalidPath indicates that an ESI bundle contains an invalid path.
	ErrEsiInvalidPath = errors.New("ESI bundle contains invalid path")

	esiBundleMsg                esipb.EsiBundle
	esiBundleDataAssetProtoName = string(esiBundleMsg.ProtoReflect().Descriptor().FullName())
)

// DeviceService implements the DeviceService rpc service.
type DeviceService struct {
	// config holds the user-provided service configuration of the DeviceService.
	config *dscpb.DeviceServiceConfig
	// esiBundle contains the ESI bundle loaded as a data asset.
	esiBundle *esipb.EsiBundle
}

// fetchESIBundle fetches the ESI bundle data asset specified in the config from the data asset
// service. It returns the ESI bundle, its asset ID, or an error if fetching or validation fails.
//
// Parameters:
//   - ctx: The context for the request.
//   - config: The DeviceServiceConfig containing the ESI bundle data asset ID.
//   - daClient: The DataAssetsClient to use for fetching the data asset.
//
// Returns:
//   - The fetched ESI bundle.
//   - An error under one of the following conditions:
//   - ErrEsiBundleNotFound: if the requested bundle is not found.
//   - ErrEsiUnmarshal: if unmarshalling fails.
//   - ErrEsiInvalidPath: if the bundle contains invalid paths.
func fetchESIBundle(ctx context.Context, config *dscpb.DeviceServiceConfig, daClient dagrpcpb.DataAssetsClient) (*esipb.EsiBundle, error) {

	iface := "data://" + esiBundleDataAssetProtoName
	anyProto, err := utils.GetDataPayload(ctx, config.GetEsiBundle(), iface, utils.WithDataAssetsClient(daClient))
	if err != nil {
		return nil, fmt.Errorf("get ESI bundle data asset for interface %q failed: %w: %w", iface, ErrEsiBundleNotFound, err)
	}

	esiBundle := &esipb.EsiBundle{}
	if err := anyProto.UnmarshalTo(esiBundle); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEsiUnmarshal, err)
	}

	if len(esiBundle.GetFiles()) == 0 {
		return nil, fmt.Errorf("ESI bundle has no files: %w", `ErrEsiUnmarshal`)
	}
	for p := range esiBundle.GetFiles() {
		if path.IsAbs(p) {
			return nil, fmt.Errorf("ESI bundle contains absolute path %q: %w: %w", p, ErrEsiInvalidPath, ErrEsiUnmarshal)
		}
		cleaned := path.Clean(p)
		if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
			return nil, fmt.Errorf("ESI bundle contains path %q resolving to outside bundle: %w: %w", p, ErrEsiInvalidPath, ErrEsiUnmarshal)
		}
	}

	log.InfoContextf(ctx, "Successfully loaded ESI bundle\n")
	return esiBundle, nil
}

// NewDeviceService creates a new DeviceService.
//
// Parameters:
//   - ctx: The context for the request.
//   - config: The DeviceServiceConfig containing service configuration.
//   - daClient: The DataAssetsClient to use for fetching ESI data assets.
//
// Returns:
//   - A new DeviceService instance.
//   - An error under one of the following conditions:
//   - ErrConfigNil: if config is nil.
//   - ErrNoDeviceIdentifier: if device identifier is missing.
//   - If fetching the ESI bundle fails.
func NewDeviceService(ctx context.Context, config *dscpb.DeviceServiceConfig, daClient dagrpcpb.DataAssetsClient) (*DeviceService, error) {
	if config == nil {
		return nil, ErrConfigNil
	}
	if config.GetDeviceIdentifier() == nil {
		return nil, ErrNoDeviceIdentifier
	}
	esiBundle, err := fetchESIBundle(ctx, config, daClient)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ESI bundle: %w", err)
	}
	log.InfoContextf(ctx, "Device service started with config: %v", config)
	return &DeviceService{config: config, esiBundle: esiBundle}, nil
}

// GetConfiguration returns the configuration of the DeviceService.
//
// Parameters:
//   - ctx: The context for the request.
//   - req: The GetConfigurationRequest.
//
// Returns:
//   - The GetConfigurationResponse containing the service configuration and ESI bundle Asset.
//   - An error if no ESI bundle is loaded.
func (s *DeviceService) GetConfiguration(ctx context.Context, req *dspb.GetConfigurationRequest) (*dspb.GetConfigurationResponse, error) {
	if s.esiBundle == nil {
		return nil, fmt.Errorf("no ESI bundle loaded")
	}
	// Assuming esipb.EsiBundle has a field named 'Esis' of type []*esipb.Esi
	return &dspb.GetConfigurationResponse{
		DeviceServiceConfig: s.config,
		EsiBundle:           s.esiBundle,
	}, nil
}
