// Copyright 2023 Intrinsic Innovation LLC

package assetviews

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	ipb "intrinsic/assets/proto/id_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
)

func mustToBase64(t *testing.T, m proto.Message) string {
	t.Helper()
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatalf("proto.Marshal() returned error: %v", err)
	}
	return base64.StdEncoding.EncodeToString(data)
}

func TestFromAsset(t *testing.T) {
	tests := []struct {
		desc                    string
		asset                   Asset
		options                 []Option
		wantStringOutput        string
		wantDescriptionFromJSON *assetDescription
	}{
		{
			desc: "AssetCatalog Service",
			asset: &acpb.Asset{
				Metadata: &metadatapb.Metadata{
					AssetType:    atypepb.AssetType_ASSET_TYPE_SERVICE,
					DisplayName:  "test_display_name",
					IdVersion:    &ipb.IdVersion{Id: &ipb.Id{Name: "test_id", Package: "test.package"}, Version: "1.0.0"},
					ReleaseNotes: "test_release_notes",
				},
			},
			wantStringOutput: "test.package.test_id.1.0.0",
			wantDescriptionFromJSON: &assetDescription{
				Name:         "test_id",
				PackageName:  "test.package",
				Version:      "1.0.0",
				ID:           "test.package.test_id",
				IDVersion:    "test.package.test_id.1.0.0",
				ReleaseNotes: "test_release_notes",
			},
		},
		{
			desc: "InstalledAssets Service",
			asset: &iapb.InstalledAsset{
				Metadata: &metadatapb.Metadata{
					AssetType:    atypepb.AssetType_ASSET_TYPE_SERVICE,
					DisplayName:  "test_display_name",
					IdVersion:    &ipb.IdVersion{Id: &ipb.Id{Name: "test_id", Package: "test.package"}, Version: "1.0.0"},
					ReleaseNotes: "test_release_notes",
				},
			},
			wantStringOutput: "test.package.test_id.1.0.0",
			wantDescriptionFromJSON: &assetDescription{
				Name:         "test_id",
				PackageName:  "test.package",
				Version:      "1.0.0",
				ID:           "test.package.test_id",
				IDVersion:    "test.package.test_id.1.0.0",
				ReleaseNotes: "test_release_notes",
			},
		},
		{
			desc: "AssetCatalog Service with id view type",
			asset: &acpb.Asset{
				Metadata: &metadatapb.Metadata{
					AssetType:    atypepb.AssetType_ASSET_TYPE_SERVICE,
					DisplayName:  "test_display_name",
					IdVersion:    &ipb.IdVersion{Id: &ipb.Id{Name: "test_id", Package: "test.package"}, Version: "1.0.0"},
					ReleaseNotes: "test_release_notes",
				},
			},
			options:          []Option{WithTextViewType(AssetTextViewTypeID)},
			wantStringOutput: "test.package.test_id",
			wantDescriptionFromJSON: &assetDescription{
				Name:         "test_id",
				PackageName:  "test.package",
				Version:      "1.0.0",
				ID:           "test.package.test_id",
				IDVersion:    "test.package.test_id.1.0.0",
				ReleaseNotes: "test_release_notes",
			},
		},
		{
			desc: "AssetCatalog Service with id version proto base64 view type",
			asset: &acpb.Asset{
				Metadata: &metadatapb.Metadata{
					AssetType:    atypepb.AssetType_ASSET_TYPE_SERVICE,
					DisplayName:  "test_display_name",
					IdVersion:    &ipb.IdVersion{Id: &ipb.Id{Name: "test_id", Package: "test.package"}, Version: "1.0.0"},
					ReleaseNotes: "test_release_notes",
				},
			},
			options:          []Option{WithTextViewType(AssetTextViewTypeIDVersionProtoBase64)},
			wantStringOutput: mustToBase64(t, &ipb.IdVersion{Id: &ipb.Id{Name: "test_id", Package: "test.package"}, Version: "1.0.0"}),
			wantDescriptionFromJSON: &assetDescription{
				Name:         "test_id",
				PackageName:  "test.package",
				Version:      "1.0.0",
				ID:           "test.package.test_id",
				IDVersion:    "test.package.test_id.1.0.0",
				ReleaseNotes: "test_release_notes",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			view := FromAsset(tc.asset, tc.options...)
			if diff := cmp.Diff(tc.wantStringOutput, view.String()); diff != "" {
				t.Errorf("FromAsset() String view returned diff (-want +got):\n%s", diff)
			}

			encoded, err := json.Marshal(view)
			if err != nil {
				t.Fatalf("json.Marshal() returned error: %v", err)
			}
			var gotDescription assetDescription
			if err := json.Unmarshal(encoded, &gotDescription); err != nil {
				t.Fatalf("json.Unmarshal() returned error: %v", err)
			}

			if diff := cmp.Diff(tc.wantDescriptionFromJSON, &gotDescription); diff != "" {
				t.Errorf("FromAsset() JSON description returned diff (-want +got):\n%s", diff)
			}
		})
	}
}
