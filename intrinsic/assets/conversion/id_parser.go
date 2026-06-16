// Copyright 2023 Intrinsic Innovation LLC

// Package idparser implements helpers for extracting asset IDs from assets.
package idparser

import (
	"errors"
	"fmt"

	"intrinsic/assets/idutils"
	idpb "intrinsic/assets/proto/id_go_proto"
	assetpb "intrinsic/assets/proto/v1/asset_go_proto"
	processedassetpb "intrinsic/assets/proto/v1/processed_asset_go_proto"
	applicationpb "intrinsic/config/proto/application_go_proto"
)

var (
	ErrAssetNil            = errors.New("asset is nil")
	ErrNoSource            = errors.New("asset source not set")
	ErrNoVariant           = errors.New("asset variant not set")
	ErrNoLocalVariant      = errors.New("local asset has no processed variant")
	ErrUnknownLocalVariant = errors.New("unknown local asset variant")
	ErrMissingID           = errors.New("asset is missing identifying ID")
	ErrMissingIDVersion    = errors.New("asset is missing identifying ID version")
	ErrInvalidID           = errors.New("asset ID is empty or invalid")
	ErrInvalidVersion      = errors.New("asset version is invalid")
)

// AssetIDVersionFromAsset returns the Asset id version for the given Asset.
// It returns an empty version if the given Asset version is not set.
func AssetIDVersionFromAsset(asset *assetpb.Asset) (*idpb.IdVersion, error) {
	if asset == nil {
		return nil, ErrAssetNil
	}

	var idVersion *idpb.IdVersion

	switch src := asset.GetSource().(type) {
	case *assetpb.Asset_Catalog:
		idVersion = src.Catalog.GetIdVersion()
	case *assetpb.Asset_Local:
		id, err := AssetIDFromAsset(asset)
		if err != nil {
			return nil, err
		}
		idVersion = &idpb.IdVersion{
			Id: id,
		}
	default:
		return nil, ErrNoSource
	}

	if idVersion == nil {
		return nil, ErrMissingIDVersion
	}
	if err := idutils.ValidateIDProto(idVersion.GetId()); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	if idVersion.GetVersion() != "" {
		if err := idutils.ValidateVersion(idVersion.GetVersion()); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidVersion, err)
		}
	}

	return idVersion, nil
}

// AssetIDFromAsset returns the Asset id for the given Asset.
func AssetIDFromAsset(asset *assetpb.Asset) (*idpb.Id, error) {
	if asset == nil {
		return nil, ErrAssetNil
	}

	var id *idpb.Id

	switch src := asset.GetSource().(type) {
	case *assetpb.Asset_Catalog:
		id = src.Catalog.GetIdVersion().GetId()
	case *assetpb.Asset_Local:
		processed := src.Local
		if processed == nil {
			return nil, ErrNoLocalVariant
		}

		switch v := processed.GetVariant().(type) {
		case *processedassetpb.ProcessedAsset_Data:
			id = v.Data.GetMetadata().GetIdVersion().GetId()
		case *processedassetpb.ProcessedAsset_HardwareDevice:
			id = v.HardwareDevice.GetMetadata().GetId()
		case *processedassetpb.ProcessedAsset_Process:
			id = v.Process.GetMetadata().GetIdVersion().GetId()
		case *processedassetpb.ProcessedAsset_SceneObject:
			id = v.SceneObject.GetMetadata().GetId()
		case *processedassetpb.ProcessedAsset_Service:
			id = v.Service.GetMetadata().GetId()
		case *processedassetpb.ProcessedAsset_Skill:
			id = v.Skill.GetMetadata().GetId()
		default:
			return nil, ErrUnknownLocalVariant
		}

	default:
		return nil, ErrNoSource
	}

	if id == nil {
		return nil, ErrMissingID
	}
	if err := idutils.ValidateIDProto(id); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	return id, nil
}

// AssetIDFromApplicationAsset return the asset id for the given application asset.
func AssetIDFromApplicationAsset(asset *applicationpb.Application_Asset) (*idpb.Id, error) {
	if asset == nil {
		return nil, ErrAssetNil
	}

	var id *idpb.Id

	switch asset.GetVariant().(type) {
	case *applicationpb.Application_Asset_Catalog:
		id = asset.GetCatalog().GetId()
	case *applicationpb.Application_Asset_Data:
		id = asset.GetData().GetMetadata().GetIdVersion().GetId()
	case *applicationpb.Application_Asset_HardwareDevice:
		id = asset.GetHardwareDevice().GetMetadata().GetId()
	case *applicationpb.Application_Asset_Process:
		id = asset.GetProcess().GetMetadata().GetIdVersion().GetId()
	case *applicationpb.Application_Asset_SceneObject:
		id = asset.GetSceneObject().GetMetadata().GetId()
	case *applicationpb.Application_Asset_Service:
		id = asset.GetService().GetMetadata().GetId()
	case *applicationpb.Application_Asset_Skill:
		id = asset.GetSkill().GetMetadata().GetId()
	default:
		return nil, ErrNoVariant
	}

	if id == nil {
		return nil, ErrMissingID
	}
	if err := idutils.ValidateIDProto(id); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	return id, nil
}
