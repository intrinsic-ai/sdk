// Copyright 2023 Intrinsic Innovation LLC

package ethercat

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"intrinsic/assets/cmdutils"
	"intrinsic/testing/grpctest"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	enipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/eni_go_proto"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

func TestUploadEniEmptyArgs(t *testing.T) {
	cmd := getUploadEniCommand()
	cmd.SetArgs([]string{})
	// Mute command usage output during tests
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))

	err := cmd.Execute()
	if err == nil {
		t.Errorf("Expected error for empty args, got nil")
	} else if !contains(err.Error(), "accepts 2 arg(s), received 0") {
		t.Errorf("Unexpected error for empty args: %v", err)
	}
}

var (
	minimalEni = `<?xml version="1.0" encoding="UTF-8"?>
	<EtherCATConfig Version="1.5">
		<Config>
			<Master>
				<Info>
					<Name>MyMaster</Name>
					<Destination>001122334455</Destination>
					<Source>AABBCCDDEEFF</Source>
				</Info>
			</Master>
		</Config>
	</EtherCATConfig>`

	minimalEniMissingSourceChild = `<?xml version="1.0" encoding="UTF-8"?>
	<EtherCATConfig Version="1.5">
		<Config>
			<Master>
				<Info>
					<Name>MyMaster</Name>
					<Destination>001122334455</Destination>
				</Info>
			</Master>
		</Config>
	</EtherCATConfig>`
)

type uploadEniTestCase struct {
	desc             string
	fileContent      string
	assetID          string
	displayName      string
	vendor           string
	wantDisplayName  string
	wantVendor       string
	wantErrMsg       string
	createAssetErr   bool
	waitOperationErr bool
}

func TestUploadEni(t *testing.T) {
	ctx := context.Background()

	tests := []uploadEniTestCase{
		{
			desc:            "success with defaults",
			fileContent:     minimalEni,
			assetID:         "my.eni.asset",
			wantDisplayName: "test.eni",
			wantVendor:      "undefined",
		},
		{
			desc:            "success with custom metadata",
			fileContent:     minimalEni,
			assetID:         "my.eni.asset",
			displayName:     "Custom Display Name",
			vendor:          "Custom Vendor",
			wantDisplayName: "Custom Display Name",
			wantVendor:      "Custom Vendor",
		},
		{
			desc:        "malformed xml",
			fileContent: "<EtherCATConfig><unclosed>",
			assetID:     "my.eni.asset",
			wantErrMsg:  "ENI file validation failed: failed to parse ENI file",
		},
		{
			desc:           "create asset error",
			fileContent:    minimalEni,
			assetID:        "my.eni.asset",
			createAssetErr: true,
			wantErrMsg:     "failed to install asset",
		},
		{
			desc:             "wait operation error",
			fileContent:      minimalEni,
			assetID:          "my.eni.asset",
			waitOperationErr: true,
			wantErrMsg:       "installation failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tmpDir := t.TempDir()
			eniPath := filepath.Join(tmpDir, "test.eni")
			if err := os.WriteFile(eniPath, []byte(tc.fileContent), 0644); err != nil {
				t.Fatalf("Failed to write test ENI file: %v", err)
			}

			mockServer := &mockInstalledAssetsServer{
				returnCreateErr: tc.createAssetErr,
				returnWaitErr:   tc.waitOperationErr,
				installedIDs:    make(map[string]bool),
			}
			srv := grpc.NewServer()
			iagrpcpb.RegisterInstalledAssetsServer(srv, mockServer)
			lrogrpcpb.RegisterOperationsServer(srv, mockServer)
			addr := grpctest.StartServerT(t, srv)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("Failed to dial mock server: %v", err)
			}
			defer conn.Close()

			oldDial := clientutilsDialClusterFromInctl
			clientutilsDialClusterFromInctl = func(ctx context.Context, flags *cmdutils.CmdFlags) (context.Context, *grpc.ClientConn, string, error) {
				return ctx, conn, addr, nil
			}
			defer func() { clientutilsDialClusterFromInctl = oldDial }()

			viper.Set("org", "intrinsic@test-org")
			cmd := getUploadEniCommand()
			flags := cmdutils.NewCmdFlags()
			flags.SetCommand(cmd)
			cmd.SetContext(ctx)

			if tc.displayName != "" {
				cmd.Flags().Set("display_name", tc.displayName)
			}
			if tc.vendor != "" {
				cmd.Flags().Set("vendor", tc.vendor)
			}

			err = runUploadENI(cmd, []string{eniPath, tc.assetID}, flags)

			if tc.wantErrMsg != "" {
				if err == nil {
					t.Errorf("runUploadENI() succeeded, want error containing %q", tc.wantErrMsg)
				} else if !contains(err.Error(), tc.wantErrMsg) {
					t.Errorf("runUploadENI() returned error %v, want error containing %q", err, tc.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("runUploadENI() failed: %v", err)
			}

			if len(mockServer.receivedRequests) != 1 {
				t.Fatalf("runUploadENI() sent %d requests, want 1", len(mockServer.receivedRequests))
			}

			req := mockServer.receivedRequests[0]
			da := req.GetAsset().GetData()
			if da.GetMetadata().GetDisplayName() != tc.wantDisplayName {
				t.Errorf("DisplayName = %q, want %q", da.GetMetadata().GetDisplayName(), tc.wantDisplayName)
			}
			if da.GetMetadata().GetVendor().GetDisplayName() != tc.wantVendor {
				t.Errorf("Vendor = %q, want %q", da.GetMetadata().GetVendor().GetDisplayName(), tc.wantVendor)
			}

			gotEni := &enipb.Eni{}
			if err := da.GetData().UnmarshalTo(gotEni); err != nil {
				t.Fatalf("could not unmarshal got ENI: %v", err)
			}
			if gotEni.GetData() != tc.fileContent {
				t.Errorf("ENI data mismatch, got %q, want %q", gotEni.GetData(), tc.fileContent)
			}

			idProto := da.GetMetadata().GetIdVersion().GetId()
			gotID := fmt.Sprintf("%s.%s", idProto.GetPackage(), idProto.GetName())
			if gotID != tc.assetID {
				t.Errorf("AssetID = %q, want %q", gotID, tc.assetID)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
