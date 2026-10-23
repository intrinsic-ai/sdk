// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package data provides conversion utilities for Data Assets.
package data

import (
	"fmt"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
)

// ManifestFromAsset extracts from the given DataAsset a DataManifest which is suitable for
// creating a corresponding bundle from it.
func ManifestFromAsset(da *dapb.DataAsset) (*dmpb.DataManifest, error) {
	if da == nil {
		return nil, fmt.Errorf("DataAsset must not be nil")
	}

	var metadata *dmpb.DataManifest_Metadata
	if m := da.GetMetadata(); m != nil {
		metadata = &dmpb.DataManifest_Metadata{
			Id:            m.GetIdVersion().GetId(),
			DisplayName:   m.GetDisplayName(),
			Documentation: m.GetDocumentation(),
			Vendor:        m.GetVendor(),
		}
	}

	return &dmpb.DataManifest{
		Metadata: metadata,
		Data:     da.GetData(),
	}, nil
}
