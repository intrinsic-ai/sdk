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

package updatereleasemetadata

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
)

type fakeAssetCatalogClient struct {
	acgrpcpb.AssetCatalogClient
	updateReleaseMetadataFunc func(ctx context.Context, in *acpb.UpdateReleaseMetadataRequest, opts ...grpc.CallOption) (*rmpb.ReleaseMetadata, error)
}

func (f *fakeAssetCatalogClient) UpdateReleaseMetadata(ctx context.Context, in *acpb.UpdateReleaseMetadataRequest, opts ...grpc.CallOption) (*rmpb.ReleaseMetadata, error) {
	if f.updateReleaseMetadataFunc != nil {
		return f.updateReleaseMetadataFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func TestUpdateAndPrintReleaseMetadata(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		target      string
		args        []string
		stdinInput  string
		clientFunc  func(ctx context.Context, in *acpb.UpdateReleaseMetadataRequest, opts ...grpc.CallOption) (*rmpb.ReleaseMetadata, error)
		wantContain string
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:   "update default flag",
			target: "ai.intrinsic.my_skill.0.0.1",
			args:   []string{"--default=true"},
			clientFunc: func(_ context.Context, in *acpb.UpdateReleaseMetadataRequest, _ ...grpc.CallOption) (*rmpb.ReleaseMetadata, error) {
				if len(in.GetUpdateMask().GetPaths()) != 1 || in.GetUpdateMask().GetPaths()[0] != "default" {
					return nil, status.Errorf(codes.InvalidArgument, "unexpected update mask: %v", in.GetUpdateMask())
				}
				if !in.GetReleaseMetadata().GetDefault() {
					return nil, status.Errorf(codes.InvalidArgument, "expected default to be true")
				}
				return &rmpb.ReleaseMetadata{Default: true}, nil
			},
			wantContain: "default:true",
		},
		{
			name:       "update org_private with prompt confirmed",
			target:     "ai.intrinsic.my_skill.0.0.1",
			args:       []string{"--org_private=false"},
			stdinInput: "y\n",
			clientFunc: func(_ context.Context, in *acpb.UpdateReleaseMetadataRequest, _ ...grpc.CallOption) (*rmpb.ReleaseMetadata, error) {
				if len(in.GetUpdateMask().GetPaths()) != 1 || in.GetUpdateMask().GetPaths()[0] != "org_private" {
					return nil, status.Errorf(codes.InvalidArgument, "unexpected update mask: %v", in.GetUpdateMask())
				}
				if in.GetReleaseMetadata().GetOrgPrivate() {
					return nil, status.Errorf(codes.InvalidArgument, "expected org_private to be false")
				}
				return &rmpb.ReleaseMetadata{Default: true}, nil
			},
			wantContain: "default:true",
		},
		{
			name:       "update org_private with prompt aborted",
			target:     "ai.intrinsic.my_skill.0.0.1",
			args:       []string{"--org_private=false"},
			stdinInput: "n\n",
			wantErr:    true,
			wantErrMsg: "aborted",
		},
		{
			name:   "update org_private with skip prompts",
			target: "ai.intrinsic.my_skill.0.0.1",
			args:   []string{"--org_private=false", "--skip_prompts=true"},
			clientFunc: func(_ context.Context, in *acpb.UpdateReleaseMetadataRequest, _ ...grpc.CallOption) (*rmpb.ReleaseMetadata, error) {
				return &rmpb.ReleaseMetadata{Default: true}, nil
			},
			wantContain: "default:true",
		},
		{
			name:    "invalid id_version",
			target:  "invalid-target!",
			wantErr: true,
		},
		{
			name:   "client error",
			target: "ai.intrinsic.my_skill.0.0.1",
			args:   []string{"--default=true"},
			clientFunc: func(_ context.Context, _ *acpb.UpdateReleaseMetadataRequest, _ ...grpc.CallOption) (*rmpb.ReleaseMetadata, error) {
				return nil, status.Error(codes.PermissionDenied, "permission denied")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			flags := cmdutils.NewCmdFlags()
			flags.SetCommand(cmd)
			flags.AddFlagDefault("asset")
			flags.AddFlagOrgPrivate()
			flags.AddFlagSkipPrompts()

			if len(tc.args) > 0 {
				if err := cmd.ParseFlags(tc.args); err != nil {
					t.Fatalf("failed to parse flags: %v", err)
				}
			}

			fakeClient := &fakeAssetCatalogClient{
				updateReleaseMetadataFunc: tc.clientFunc,
			}
			var stdoutBuf bytes.Buffer
			stdinReader := strings.NewReader(tc.stdinInput)
			prtr, err := printer.NewPrinterWithWriter(printer.TextOutputFormat, &stdoutBuf)
			if err != nil {
				t.Fatalf("failed to create printer: %v", err)
			}

			err = updateAndPrintReleaseMetadata(ctx, fakeClient, tc.target, flags, stdinReader, &stdoutBuf, prtr)
			if (err != nil) != tc.wantErr {
				t.Errorf("updateAndPrintReleaseMetadata() err = %v, wantErr = %v", err, tc.wantErr)
				return
			}
			if tc.wantErrMsg != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErrMsg)) {
				t.Errorf("expected error message containing %q, got %v", tc.wantErrMsg, err)
			}
			if err == nil && tc.wantContain != "" && !strings.Contains(stdoutBuf.String(), tc.wantContain) {
				t.Errorf("output = %q, want containing %q", stdoutBuf.String(), tc.wantContain)
			}
		})
	}
}
