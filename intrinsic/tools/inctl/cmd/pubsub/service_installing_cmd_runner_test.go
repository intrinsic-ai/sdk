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

package pubsub

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	pubsubtesting "intrinsic/tools/inctl/cmd/pubsub/testing"
)

func setupTestCatalogServer(t *testing.T) (*pubsubtesting.FakeAssetCatalogServer, string) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen on local port: %v", err)
	}
	catalogAddr := lis.Addr().String()
	s := grpc.NewServer()
	catalogServer := pubsubtesting.NewFakeAssetCatalogServer()
	acgrpcpb.RegisterAssetCatalogServer(s, catalogServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			// Ignore error
		}
	}()

	t.Cleanup(func() {
		s.Stop()
		lis.Close()
	})

	return catalogServer, catalogAddr
}

func TestGetDefaultVersion(t *testing.T) {
	tests := []struct {
		name              string
		setupFakeServer   func(s *pubsubtesting.FakeAssetCatalogServer)
		expectedResult    string
		expectErr         bool
		expectErrContains string
	}{
		{
			name: "Fetch Default Version Success",
			setupFakeServer: func(s *pubsubtesting.FakeAssetCatalogServer) {
				s.GetAssetFn = func(ctx context.Context, req *acgrpcpb.GetAssetRequest) (*acgrpcpb.Asset, error) {
					return &acgrpcpb.Asset{
						Metadata: &metadatapb.Metadata{
							IdVersion: &idpb.IdVersion{
								Version: "4.5.6",
							},
						},
					}, nil
				}
			},
			expectedResult: "4.5.6",
			expectErr:      false,
		},
		{
			name: "Fetch Default Version GetAsset Error",
			setupFakeServer: func(s *pubsubtesting.FakeAssetCatalogServer) {
				s.GetAssetFn = func(ctx context.Context, req *acgrpcpb.GetAssetRequest) (*acgrpcpb.Asset, error) {
					return nil, errors.New("backend error")
				}
			},
			expectErr:         true,
			expectErrContains: "failed to fetch " + hubServiceName + " from catalog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalogServer, catalogAddr := setupTestCatalogServer(t)
			if tt.setupFakeServer != nil {
				tt.setupFakeServer(catalogServer)
			}

			var buf bytes.Buffer
			ctx := context.Background()

			conn, err := grpc.Dial(catalogAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("failed to connect to the catalog server at %v: %v", catalogAddr, err)
			}
			defer conn.Close()

			result, err := getDefaultVersionFromCatalog(ctx, conn, &buf, hubServicePackage, hubServiceName)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectErrContains)
				}
				if !strings.Contains(err.Error(), tt.expectErrContains) {
					t.Errorf("expected error to contain %q, got %v", tt.expectErrContains, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("expected result %q, got %q", tt.expectedResult, result)
				}
			}
		})
	}
}
