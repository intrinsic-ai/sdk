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

package get

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
)

type fakeAssetCatalogClient struct {
	acgrpcpb.AssetCatalogClient
	getAssetFunc func(ctx context.Context, in *acpb.GetAssetRequest, opts ...grpc.CallOption) (*acpb.Asset, error)
}

func (f *fakeAssetCatalogClient) GetAsset(ctx context.Context, in *acpb.GetAssetRequest, opts ...grpc.CallOption) (*acpb.Asset, error) {
	if f.getAssetFunc != nil {
		return f.getAssetFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func TestGetAndPrintAsset(t *testing.T) {
	ctx := context.Background()

	testAsset := &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Package: "ai.intrinsic",
					Name:    "my_skill",
				},
				Version: "0.0.1",
			},
		},
	}

	tests := []struct {
		name        string
		target      string
		view        viewpb.AssetViewType
		clientFunc  func(ctx context.Context, in *acpb.GetAssetRequest, opts ...grpc.CallOption) (*acpb.Asset, error)
		wantContain string
		wantErr     bool
	}{
		{
			name:   "success with id_version",
			target: "ai.intrinsic.my_skill.0.0.1",
			view:   viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			clientFunc: func(_ context.Context, in *acpb.GetAssetRequest, _ ...grpc.CallOption) (*acpb.Asset, error) {
				if in.GetIdVersion().GetVersion() != "0.0.1" {
					return nil, status.Errorf(codes.InvalidArgument, "unexpected version: %v", in.GetIdVersion())
				}
				return testAsset, nil
			},
			wantContain: "ai.intrinsic",
		},
		{
			name:   "success with id only",
			target: "ai.intrinsic.my_skill",
			view:   viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			clientFunc: func(_ context.Context, in *acpb.GetAssetRequest, _ ...grpc.CallOption) (*acpb.Asset, error) {
				if in.GetId().GetName() != "my_skill" {
					return nil, status.Errorf(codes.InvalidArgument, "unexpected id: %v", in.GetId())
				}
				return testAsset, nil
			},
			wantContain: "my_skill",
		},
		{
			name:    "invalid target",
			target:  "invalid-target-format!",
			view:    viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			wantErr: true,
		},
		{
			name:   "not found",
			target: "ai.intrinsic.unknown.0.0.1",
			view:   viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			clientFunc: func(_ context.Context, _ *acpb.GetAssetRequest, _ ...grpc.CallOption) (*acpb.Asset, error) {
				return nil, status.Error(codes.NotFound, "asset not found")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fakeAssetCatalogClient{
				getAssetFunc: tc.clientFunc,
			}
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&buf)
			prtr, err := printer.NewPrinterFromCommand(cmd)
			if err != nil {
				t.Fatalf("failed to create printer: %v", err)
			}
			defer printer.Flush(prtr)

			err = getAndPrintAsset(ctx, fakeClient, tc.target, tc.view, prtr)
			if (err != nil) != tc.wantErr {
				t.Errorf("getAndPrintAsset() err = %v, wantErr = %v", err, tc.wantErr)
				return
			}
			if err == nil && !strings.Contains(buf.String(), tc.wantContain) {
				t.Errorf("buf.String() = %q, want containing %q", buf.String(), tc.wantContain)
			}
		})
	}
}
