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

package data

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	docpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestManifestFromAsset(t *testing.T) {
	payload, err := anypb.New(&emptypb.Empty{})
	if err != nil {
		t.Fatalf("anypb.New() failed: %v", err)
	}

	tests := []struct {
		desc    string
		da      *dapb.DataAsset
		want    *dmpb.DataManifest
		wantErr bool
	}{
		{
			desc: "populated DataAsset",
			da: &dapb.DataAsset{
				Metadata: &metadatapb.Metadata{
					AssetType:   atpb.AssetType_ASSET_TYPE_DATA,
					DisplayName: "Some Data asset",
					Documentation: &docpb.Documentation{
						Description: "Some description",
					},
					IdVersion: &idpb.IdVersion{
						Id: &idpb.Id{
							Name:    "data_asset",
							Package: "ai.intrinsic",
						},
					},
					Vendor: &vpb.Vendor{
						DisplayName: "Intrinsic",
					},
				},
				Data: payload,
			},
			want: &dmpb.DataManifest{
				Metadata: &dmpb.DataManifest_Metadata{
					Id: &idpb.Id{
						Name:    "data_asset",
						Package: "ai.intrinsic",
					},
					DisplayName: "Some Data asset",
					Documentation: &docpb.Documentation{
						Description: "Some description",
					},
					Vendor: &vpb.Vendor{
						DisplayName: "Intrinsic",
					},
				},
				Data: payload,
			},
		},
		{
			desc: "empty DataAsset",
			da:   &dapb.DataAsset{},
			want: &dmpb.DataManifest{},
		},
		{
			desc:    "nil DataAsset",
			da:      nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ManifestFromAsset(tc.da)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ManifestFromAsset() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ManifestFromAsset() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
