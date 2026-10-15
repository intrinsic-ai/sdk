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

package view

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"intrinsic/tools/inctl/util/printer"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
	assettagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
)

func sampleAsset() *acpb.Asset {
	return &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Package: "ai.intrinsic",
					Name:    "my_skill",
				},
				Version: "1.2.3",
			},
			DisplayName: "My Skill",
			AssetType:   atypepb.AssetType_ASSET_TYPE_SKILL,
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic AI",
			},
			AssetTag:   assettagpb.AssetTag_ASSET_TAG_CAMERA,
			UpdateTime: timestamppb.New(time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC)),
			Documentation: &documentationpb.Documentation{
				Description: "A skill for object detection.",
			},
			ReleaseNotes: "- Initial release\n- Performance improvements",
			Provides: []*metadatapb.Interface{
				{Uri: "grpc://intrinsic_proto.services.Detector"},
			},
		},
		ReleaseMetadata: &rmpb.ReleaseMetadata{
			Default:    true,
			OrgPrivate: false,
		},
	}
}

func TestViewOfNil(t *testing.T) {
	v := Of(nil)
	if s := v.String(); s != "" {
		t.Errorf("String() = %q, want empty", s)
	}

	b, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() err = %v", err)
	}
	if string(b) != "null" {
		t.Errorf("MarshalJSON() = %q, want %q", string(b), "null")
	}

	tab := v.Tabulated([]string{"id", "version"})
	if len(tab) != 2 || tab[0] != "" || tab[1] != "" {
		t.Errorf("Tabulated() = %v, want empty columns", tab)
	}
}

func TestViewMarshalJSON(t *testing.T) {
	asset := sampleAsset()
	v := Of(asset)

	got, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() err = %v", err)
	}

	want, err := protojson.Marshal(asset)
	if err != nil {
		t.Fatalf("protojson.Marshal() err = %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("MarshalJSON() diff (-want +got):\n%s", cmp.Diff(string(want), string(got)))
	}
}

func TestViewTabulated(t *testing.T) {
	asset := sampleAsset()
	v := Of(asset)

	columns := []string{
		"id",
		"version",
		"display_name",
		"name",
		"default",
		"org_private",
		"update_time",
		"unknown_col",
	}

	got := v.Tabulated(columns)
	want := []string{
		"ai.intrinsic.my_skill",
		"1.2.3",
		"My Skill",
		"My Skill",
		"true",
		"false",
		"2026-10-07T10:00:00Z",
		"",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Tabulated() diff (-want +got):\n%s", diff)
	}
}

func TestHeaderName(t *testing.T) {
	tests := []struct {
		col  string
		want string
	}{
		{"id", "ID"},
		{"version", "VERSION"},
		{"default", "DEFAULT"},
		{"org_private", "ORG PRIVATE"},
		{"update_time", "UPDATE TIME"},
		{"display_name", "DISPLAY NAME"},
		{"display-name", "DISPLAY NAME"},
		{"custom_field", "CUSTOM FIELD"},
	}

	for _, tc := range tests {
		t.Run(tc.col, func(t *testing.T) {
			if got := headerName(tc.col); got != tc.want {
				t.Errorf("headerName(%q) = %q, want %q", tc.col, got, tc.want)
			}
		})
	}
}

func TestStringMinimalAsset(t *testing.T) {
	asset := &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Package: "ai.intrinsic",
					Name:    "simple_asset",
				},
				Version: "0.0.1",
			},
		},
	}

	v := Of(asset)
	got := v.String()

	want := `ID:           ai.intrinsic.simple_asset
Version:      0.0.1`

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestStringFullAsset(t *testing.T) {
	asset := sampleAsset()
	v := Of(asset)
	got := v.String()

	mustContain := []string{
		"ID:           ai.intrinsic.my_skill",
		"Version:      1.2.3",
		"Display Name: My Skill",
		"Type:         Skill",
		"Vendor:       Intrinsic AI",
		"Asset Tag:    Camera",
		"Update Time:  2026-10-07T10:00:00Z",
		"Description:  A skill for object detection.",
		"Release Notes:\n  - Initial release\n  - Performance improvements",
		"Provides:\n  - grpc://intrinsic_proto.services.Detector",
		"Release:\n  Default:     true\n  Org Private: false",
	}

	for _, substr := range mustContain {
		if !strings.Contains(got, substr) {
			t.Errorf("String() missing expected substring %q\nFull output:\n%s", substr, got)
		}
	}
}

func TestTableConfigOutput(t *testing.T) {
	asset := sampleAsset()
	var buf bytes.Buffer

	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	prtr, err := printer.NewPrinterOfType(printer.OutputTypeTAB, cmd, TableConfig("id", "version", "display_name"))
	if err != nil {
		t.Fatalf("failed to create printer: %v", err)
	}

	prtr.Print(Of(asset))
	printer.Flush(prtr)

	out := buf.String()
	// Must contain headers
	if !strings.Contains(out, "ID") || !strings.Contains(out, "VERSION") || !strings.Contains(out, "DISPLAY NAME") {
		t.Errorf("output missing headers, got:\n%s", out)
	}
	// Must contain row values
	if !strings.Contains(out, "ai.intrinsic.my_skill") || !strings.Contains(out, "1.2.3") || !strings.Contains(out, "My Skill") {
		t.Errorf("output missing row data, got:\n%s", out)
	}
}
