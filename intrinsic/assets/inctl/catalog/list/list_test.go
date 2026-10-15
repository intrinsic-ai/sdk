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

package list

import (
	"bytes"
	"context"
	"strings"
	"testing"

	catalogview "intrinsic/assets/inctl/catalog/view"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
)

type fakeAssetCatalogClient struct {
	acgrpcpb.AssetCatalogClient
	listAssetsFunc func(ctx context.Context, in *acpb.ListAssetsRequest, opts ...grpc.CallOption) (*acpb.ListAssetsResponse, error)
}

func (f *fakeAssetCatalogClient) ListAssets(ctx context.Context, in *acpb.ListAssetsRequest, opts ...grpc.CallOption) (*acpb.ListAssetsResponse, error) {
	if f.listAssetsFunc != nil {
		return f.listAssetsFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func makeTestAsset(pkg, name, version string) *acpb.Asset {
	return &acpb.Asset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Package: pkg,
					Name:    name,
				},
				Version: version,
			},
		},
	}
}

func TestListAndPrintAssets(t *testing.T) {
	ctx := context.Background()

	assetB := makeTestAsset("ai.intrinsic", "b_skill", "0.0.1")
	assetA := makeTestAsset("ai.intrinsic", "a_skill", "0.0.1")

	tests := []struct {
		name         string
		view         viewpb.AssetViewType
		cols         []string
		filter       *acpb.ListAssetsRequest_AssetFilter
		clientFunc   func(ctx context.Context, in *acpb.ListAssetsRequest, opts ...grpc.CallOption) (*acpb.ListAssetsResponse, error)
		wantContains []string
		wantErr      bool
	}{
		{
			name: "list all default assets sorted",
			view: viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			cols: []string{"id", "version"},
			filter: &acpb.ListAssetsRequest_AssetFilter{
				OnlyDefault: proto.Bool(true),
			},
			clientFunc: func(_ context.Context, in *acpb.ListAssetsRequest, _ ...grpc.CallOption) (*acpb.ListAssetsResponse, error) {
				if !in.GetStrictFilter().GetOnlyDefault() {
					return nil, status.Errorf(codes.InvalidArgument, "expected only_default to be true")
				}
				return &acpb.ListAssetsResponse{
					Assets: []*acpb.Asset{assetA, assetB},
				}, nil
			},
			wantContains: []string{
				"ID", "VERSION",
				"ai.intrinsic.a_skill", "0.0.1",
				"ai.intrinsic.b_skill", "0.0.1",
			},
		},
		{
			name: "list versions of specific asset",
			view: viewpb.AssetViewType_ASSET_VIEW_TYPE_VERSIONS,
			cols: []string{"version", "default", "org_private", "update_time"},
			filter: &acpb.ListAssetsRequest_AssetFilter{
				Id: proto.String("ai.intrinsic.a_skill"),
			},
			clientFunc: func(_ context.Context, in *acpb.ListAssetsRequest, _ ...grpc.CallOption) (*acpb.ListAssetsResponse, error) {
				if in.GetStrictFilter().GetId() != "ai.intrinsic.a_skill" {
					return nil, status.Errorf(codes.InvalidArgument, "unexpected id: %s", in.GetStrictFilter().GetId())
				}
				return &acpb.ListAssetsResponse{
					Assets: []*acpb.Asset{assetA},
				}, nil
			},
			wantContains: []string{
				"VERSION", "DEFAULT", "ORG PRIVATE", "UPDATE TIME",
				"0.0.1",
			},
		},
		{
			name: "client error",
			view: viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC,
			clientFunc: func(_ context.Context, _ *acpb.ListAssetsRequest, _ ...grpc.CallOption) (*acpb.ListAssetsResponse, error) {
				return nil, status.Error(codes.Internal, "internal server error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fakeAssetCatalogClient{
				listAssetsFunc: tc.clientFunc,
			}
			var buf bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&buf)
			prtr, err := printer.NewPrinterOfType(printer.OutputTypeTAB, cmd, catalogview.TableConfig(tc.cols...))
			if err != nil {
				t.Fatalf("failed to create printer: %v", err)
			}

			err = listAndPrintAssets(ctx, fakeClient, tc.view, tc.filter, prtr)
			if (err != nil) != tc.wantErr {
				t.Errorf("listAndPrintAssets() err = %v, wantErr = %v", err, tc.wantErr)
				return
			}
			if err == nil {
				for _, substr := range tc.wantContains {
					if !strings.Contains(buf.String(), substr) {
						t.Errorf("buf.String() = %q, want containing %q", buf.String(), substr)
					}
				}
			}
		})
	}
}
