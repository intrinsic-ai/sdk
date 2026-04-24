// Copyright 2023 Intrinsic Innovation LLC

// Package datafix contains utils that adapt Data Assets to meet the requirements of the latest
// platform version.
package datafix

import (
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
)

// Manifest updates a DataManifest to meet the requirements of the latest platform version.
func Manifest(manifest *dmpb.DataManifest) error {
	return nil
}

// DataAsset updates a DataAsset to meet the requirements of the latest platform version.
func DataAsset(da *dapb.DataAsset) error {
	// The metadata in an Asset definition shouldn't specify a version.
	if da.GetMetadata().GetIdVersion().GetVersion() != "" {
		da.GetMetadata().GetIdVersion().Version = ""
	}
	return nil
}
