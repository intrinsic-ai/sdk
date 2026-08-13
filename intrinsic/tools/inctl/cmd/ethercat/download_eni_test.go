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

package ethercat

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fakedataassets "intrinsic/assets/data/fakedataassets"
	"intrinsic/assets/idutils"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	enipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/eni_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

func TestDownloadEni(t *testing.T) {
	ctx := context.Background()
	assetID := "my.package.my_eni"
	unknownAssetID := "my.unknown.package.my_eni"
	idProto, err := idutils.IDProtoFrom("my.package", "my_eni")
	if err != nil {
		t.Fatalf("Failed to parse asset ID %q: %v", assetID, err)
	}
	eniData := "<eni_data>"
	eniProto := &enipb.Eni{Data: eniData}
	eniAny, err := anypb.New(eniProto)
	if err != nil {
		t.Fatalf("Failed to marshal ENI proto: %v", err)
	}

	tests := []struct {
		name          string
		assetID       string
		fakeDataAsset *dapb.DataAsset
		outputFile    func() string
		wantContent   string
		wantErr       bool
		wantErrMsg    string
	}{
		{
			name:       "asset id invalid",
			assetID:    "invalid-id",
			wantErr:    true,
			wantErrMsg: "not a valid id",
		},
		{
			name:       "unknown data asset error",
			assetID:    unknownAssetID,
			wantErr:    true,
			wantErrMsg: "data asset not found",
		},
		{
			name:    "success no output file",
			assetID: assetID,
			fakeDataAsset: &dapb.DataAsset{
				Metadata: &mpb.Metadata{IdVersion: &idpb.IdVersion{Id: idProto}},
				Data:     eniAny,
			},
			wantContent: eniData,
		},
		{
			name:    "success with output file",
			assetID: assetID,
			fakeDataAsset: &dapb.DataAsset{
				Metadata: &mpb.Metadata{IdVersion: &idpb.IdVersion{Id: idProto}},
				Data:     eniAny,
			},
			outputFile: func() string {
				return filepath.Join(t.TempDir(), "eni_out.xml")
			},
			wantContent: eniData,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var fakeAssets []*dapb.DataAsset
			if tc.fakeDataAsset != nil {
				fakeAssets = append(fakeAssets, tc.fakeDataAsset)
			}
			fake := fakedataassets.StartServer(ctx, t, fakedataassets.WithDataAssets(fakeAssets))

			outFile := ""
			if tc.outputFile != nil {
				outFile = tc.outputFile()
			}

			var outBuf bytes.Buffer
			err := downloadEni(ctx, fake.Client, tc.assetID, outFile, &outBuf)

			if tc.wantErr {
				if err == nil {
					t.Errorf("downloadEni() succeeded, want error")
				}
				if !strings.Contains(err.Error(), tc.wantErrMsg) {
					t.Errorf("downloadEni() returned unexpected error: `%v` expected to contain %q", err, tc.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("downloadEni() failed: %v", err)
			}

			if outFile != "" {
				content, err := os.ReadFile(outFile)
				if err != nil {
					t.Fatalf("Failed to read output file %q: %v", outFile, err)
				}
				if string(content) != tc.wantContent {
					t.Errorf("downloadEni() wrote %q to output file, want %q", string(content), tc.wantContent)
				}
			} else {
				// The logic adds a newline, so we need to trim it for comparison.
				got := strings.TrimSpace(outBuf.String())
				if got != tc.wantContent {
					t.Errorf("downloadEni() wrote %q to stdout, want %q", got, tc.wantContent)
				}
			}
		})
	}
}
