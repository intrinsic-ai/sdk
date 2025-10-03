// Copyright 2023 Intrinsic Innovation LLC

// Package assetdescriptions contains utils for commands that list Assets.
package assetdescriptions

import (
	"encoding/json"
	"sort"
	"strings"

	"intrinsic/assets/idutils"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
)

// AssetTextViewType is the view type of an Asset when formatted for text output.
type AssetTextViewType int

const (
	// AssetTextViewTypeIDVersion is the view type for outputting Asset ID versions.
	AssetTextViewTypeIDVersion AssetTextViewType = iota
	// AssetTextViewTypeID is the view type for outputting Asset IDs.
	AssetTextViewTypeID
)

// AssetDescription describes an Asset.
//
// It has custom proto->json conversion to handle fields like the update timestamp.
type AssetDescription struct {
	Name         string `json:"name,omitempty"`
	Vendor       string `json:"vendor,omitempty"`
	PackageName  string `json:"packageName,omitempty"`
	Version      string `json:"version,omitempty"`
	UpdateTime   string `json:"updateTime,omitempty"`
	ID           string `json:"id,omitempty"`
	IDVersion    string `json:"idVersion,omitempty"`
	ReleaseNotes string `json:"releaseNotes,omitempty"`
	Description  string `json:"description,omitempty"`
}

// AssetDescriptions describes a list of Assets.
//
// It implements methods intended for use with printing outputs of Asset list commands.
type AssetDescriptions struct {
	Assets       []*AssetDescription `json:"assets"`
	TextViewType AssetTextViewType
}

// MarshalJSON marshals the underlying Asset descriptions.
func (ad *AssetDescriptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(ad.Assets)
}

// String returns a string with one Asset per line.
func (ad *AssetDescriptions) String() string {
	var lines []string
	for _, asset := range ad.Assets {
		var line string
		switch ad.TextViewType {
		case AssetTextViewTypeID:
			line = asset.ID
		default: // Default to AssetTextViewTypeIDVersion
			line = asset.IDVersion
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// Option is a functional option for AssetDescriptions.
type Option func(*AssetDescriptions)

// WithTextViewType sets the text view type of the AssetDescriptions.
func WithTextViewType(viewType AssetTextViewType) Option {
	return func(ad *AssetDescriptions) {
		ad.TextViewType = viewType
	}
}

// FromCatalogAssets returns AssetDescriptions from AssetCatalog Assets.
func FromCatalogAssets(assets []*acpb.Asset, options ...Option) (*AssetDescriptions, error) {
	descriptions := &AssetDescriptions{
		Assets: make([]*AssetDescription, len(assets)),
	}
	for _, option := range options {
		option(descriptions)
	}

	for i, asset := range assets {
		metadata := asset.GetMetadata()
		idVersion, err := idutils.IDVersionFromProto(metadata.GetIdVersion())
		if err != nil {
			return nil, err
		}
		ivp, err := idutils.NewIDVersionParts(idVersion)
		if err != nil {
			return nil, err
		}

		description := &AssetDescription{
			Name:         ivp.Name(),
			Vendor:       metadata.GetVendor().GetDisplayName(),
			PackageName:  ivp.Package(),
			Version:      ivp.Version(),
			ID:           ivp.ID(),
			IDVersion:    idVersion,
			ReleaseNotes: metadata.GetReleaseNotes(),
			Description:  metadata.GetDocumentation().GetDescription(),
		}
		if metadata.GetUpdateTime() != nil {
			description.UpdateTime = metadata.GetUpdateTime().AsTime().String()
		}
		descriptions.Assets[i] = description
	}

	return descriptions, nil
}
