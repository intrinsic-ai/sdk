// Copyright 2023 Intrinsic Innovation LLC

// Package processfix contains utils that adapt Processes to meet the requirements of the latest
// platform version.
package processfix

import (
	papb "intrinsic/assets/processes/proto/process_asset_go_proto"
	pmpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
)

// Manifest updates a ProcessManifest to meet the requirements of the latest platform version.
func Manifest(manifest *pmpb.ProcessManifest) error {
	return nil
}

// DataAsset updates a ProcessAsset to meet the requirements of the latest platform version.
func ProcessAsset(pa *papb.ProcessAsset) error {
	// The metadata in an Asset definition shouldn't specify a version.
	if pa.GetMetadata().GetIdVersion().GetVersion() != "" {
		pa.GetMetadata().GetIdVersion().Version = ""
	}
	return nil
}
