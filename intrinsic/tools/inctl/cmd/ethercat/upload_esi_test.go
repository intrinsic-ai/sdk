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
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"intrinsic/assets/cmdutils"
	"intrinsic/testing/grpctest"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"

	iagrpcpb "intrinsic/assets/proto/installed_assets_go_proto"
	iapb "intrinsic/assets/proto/installed_assets_go_proto"
	esipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/esi_go_proto"

	lrogrpcpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

// mockInstalledAssetsServer is a mock implementation of the InstalledAssetsServer and OperationsServer.
type mockInstalledAssetsServer struct {
	iagrpcpb.UnimplementedInstalledAssetsServer
	receivedRequests []*iapb.CreateInstalledAssetRequest
	lrogrpcpb.UnimplementedOperationsServer
	returnCreateErr bool
	returnWaitErr   bool
	installedIDs    map[string]bool
}

// CreateInstalledAsset is a mock implementation of the CreateInstalledAsset RPC.
func (s *mockInstalledAssetsServer) CreateInstalledAsset(ctx context.Context, req *iapb.CreateInstalledAssetRequest) (*lropb.Operation, error) {
	s.receivedRequests = append(s.receivedRequests, req)
	if s.returnCreateErr {
		return nil, fmt.Errorf("mock create error")
	}
	idProto := req.GetAsset().GetData().GetMetadata().GetIdVersion().GetId()
	id := fmt.Sprintf("%s.%s", idProto.GetPackage(), idProto.GetName())
	if s.installedIDs[id] {
		return nil, fmt.Errorf("asset with id %q already exists", id)
	}
	s.installedIDs[id] = true
	return &lropb.Operation{Name: "testOp"}, nil
}

// WaitOperation is a mock implementation of the WaitOperation RPC.
func (s *mockInstalledAssetsServer) WaitOperation(ctx context.Context, req *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	if s.returnWaitErr {
		return nil, fmt.Errorf("mock wait error")
	}
	return &lropb.Operation{Name: "testOp", Done: true}, nil
}

func createTestFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create directory for %q: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %q: %v", name, err)
		}
	}
}

func TestUploadESIBundles(t *testing.T) {
	ctx := context.Background()
	simpleInfo := `<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>`
	infoWithDictRef := `<EtherCATInfo><Vendor><Name>VendorA</Name></Vendor><Descriptions><Devices><Device><Profile><DictionaryFile>dict.xml</DictionaryFile></Profile></Device></Devices></Descriptions></EtherCATInfo>`
	dict := `<EtherCATDict></EtherCATDict>`
	infoWithSubDirRef := `<EtherCATInfo><Vendor><Name>VendorC</Name></Vendor><InfoReference>sub/module.xml</InfoReference></EtherCATInfo>`
	module := `<EtherCATModule></EtherCATModule>`
	info1WithRef := `<EtherCATInfo><Vendor><Name>VendorE</Name></Vendor><InfoReference>info2.xml</InfoReference></EtherCATInfo>`
	info2 := `<EtherCATInfo><Vendor><Name>VendorF</Name></Vendor></EtherCATInfo>`
	malformedXML := `<EtherCATInfo><Vendor><Name>TestVendor</Name>`
	iso88591Info := `<?xml version="1.0" encoding="ISO-8859-1"?>
	<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>`
	iso88591InfoDecoded := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n&#x9;<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>"
	infoMissingVendor := `<EtherCATInfo><Vendor></Vendor></EtherCATInfo>`

	tests := []struct {
		desc             string
		files            map[string]string
		overrideIDs      []string
		wantBundles      map[string]*esipb.EsiBundle
		wantErr          error
		wantReqCount     int
		createAssetErr   bool
		waitOperationErr bool
	}{
		{
			desc: "single file bundle",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			wantBundles: map[string]*esipb.EsiBundle{
				"ai.intrinsic.ethercat.esi.testvendor.simple_xml": {
					Files: map[string]*esipb.Esi{"simple.xml": {Data: simpleInfo}},
				},
			},
			wantReqCount: 1,
		},
		{
			desc: "single file bundle with override-id",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			overrideIDs: []string{"simple.xml=my.package.simple_xml"},
			wantBundles: map[string]*esipb.EsiBundle{
				"my.package.simple_xml": {
					Files: map[string]*esipb.Esi{"simple.xml": {Data: simpleInfo}},
				},
			},
			wantReqCount: 1,
		},
		{
			desc: "override-id with unknown file",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			overrideIDs:  []string{"unknown.xml=my.package.id"},
			wantErr:      errOverrideNotPrimaryFile,
			wantReqCount: 0,
		},
		{
			desc: "override-id with non-primary file",
			files: map[string]string{
				"info_dict.xml": infoWithDictRef,
				"dict.xml":      dict,
			},
			overrideIDs:  []string{"dict.xml=my.package.id"},
			wantErr:      errOverrideNotPrimaryFile,
			wantReqCount: 0,
		},
		{
			desc: "override-id with invalid format",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			overrideIDs:  []string{"simple.xmlmy.package.id"},
			wantErr:      errInvalidOverrideFormat,
			wantReqCount: 0,
		},
		{
			desc: "override-id with invalid id value format",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			overrideIDs:  []string{"simple.xml=invalidformat"},
			wantErr:      errInvalidIDString,
			wantReqCount: 0,
		},
		{
			desc: "multi-file bundle",
			files: map[string]string{
				"info_dict.xml": infoWithDictRef,
				"dict.xml":      dict,
			},
			wantBundles: map[string]*esipb.EsiBundle{
				"ai.intrinsic.ethercat.esi.vendora.info_dict_xml": {
					Files: map[string]*esipb.Esi{
						"info_dict.xml": {Data: infoWithDictRef},
						"dict.xml":      {Data: dict},
					},
				},
			},
			wantReqCount: 1,
		},
		{
			desc: "multi-file bundle with subdirectory",
			files: map[string]string{
				"info_subdir.xml": infoWithSubDirRef,
				"sub/module.xml":  module,
			},
			wantBundles: map[string]*esipb.EsiBundle{
				"ai.intrinsic.ethercat.esi.vendorc.info_subdir_xml": {
					Files: map[string]*esipb.Esi{
						"info_subdir.xml": {Data: infoWithSubDirRef},
						"sub/module.xml":  {Data: module},
					},
				},
			},
			wantReqCount: 1,
		},
		{
			desc: "bundle with two info files - error",
			files: map[string]string{
				"info1_ref.xml": info1WithRef,
				"info2.xml":     info2,
			},
			wantErr:      errMultipleInfoFiles,
			wantReqCount: 0,
		},
		{
			desc: "two distinct bundles, one multi-file",
			files: map[string]string{
				"simple.xml":    simpleInfo,
				"info_dict.xml": infoWithDictRef,
				"dict.xml":      dict,
			},
			wantBundles: map[string]*esipb.EsiBundle{
				"ai.intrinsic.ethercat.esi.testvendor.simple_xml": {
					Files: map[string]*esipb.Esi{"simple.xml": {Data: simpleInfo}},
				},
				"ai.intrinsic.ethercat.esi.vendora.info_dict_xml": {
					Files: map[string]*esipb.Esi{
						"info_dict.xml": {Data: infoWithDictRef},
						"dict.xml":      {Data: dict},
					},
				},
			},
			wantReqCount: 2,
		},
		{
			desc: "malformed xml",
			files: map[string]string{
				"malformed.xml": malformedXML,
			},
			wantErr:      errMalformedXML,
			wantReqCount: 0,
		},
		{
			desc: "create asset error",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			wantErr:        errInstallAssetFailed,
			wantReqCount:   1,
			createAssetErr: true,
		},
		{
			desc: "wait operation error",
			files: map[string]string{
				"simple.xml": simpleInfo,
			},
			wantErr:          errWaitForInstallationFailed,
			wantReqCount:     1,
			waitOperationErr: true,
		},
		{
			desc: "non-utf8 encoding",
			files: map[string]string{
				"iso88591.xml": iso88591Info,
			},
			wantBundles: map[string]*esipb.EsiBundle{
				"ai.intrinsic.ethercat.esi.testvendor.iso88591_xml": {
					Files: map[string]*esipb.Esi{"iso88591.xml": {Data: iso88591InfoDecoded}},
				},
			},
			wantReqCount: 1,
		},
		{
			desc: "missing vendor name",
			files: map[string]string{
				"missing_vendor.xml": infoMissingVendor,
			},
			wantErr:      errVendorNameMissing,
			wantReqCount: 0,
		},
		{
			desc: "two identical bundles in different folders",
			files: map[string]string{
				"folder1/simple.xml": simpleInfo,
				"folder2/simple.xml": simpleInfo,
			},
			wantErr:      errInstallAssetFailed,
			wantReqCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tmpDir := t.TempDir()
			createTestFiles(t, tmpDir, tc.files)

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
			cmd := &cobra.Command{}
			flags := cmdutils.NewCmdFlags()
			flags.SetCommand(cmd)
			flags.AddFlagsAddressClusterSolution()
			flags.AddFlagPolicy("data")
			flags.AddFlagsProjectOrg()
			cmd.Flags().StringSlice("override_id", []string{}, "Override generated asset ID for a bundle. Format: <primary_file_basename>=<package.name>")
			cmd.SetContext(ctx)
			if len(tc.overrideIDs) > 0 {
				if err := cmd.Flags().Set("override_id", strings.Join(tc.overrideIDs, ",")); err != nil {
					t.Fatalf("Failed to set override_id flag: %v", err)
				}
			}

			err = runUploadESI(cmd, []string{tmpDir}, flags)

			if tc.wantErr != nil {
				if err == nil {
					t.Errorf("runUploadESI() succeeded, want error %v", tc.wantErr)
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("runUploadESI() returned error %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("runUploadESI() failed: %v", err)
			}

			if len(mockServer.receivedRequests) != tc.wantReqCount {
				t.Fatalf("runUploadESI() sent %d requests, want %d", len(mockServer.receivedRequests), tc.wantReqCount)
			}

			for _, req := range mockServer.receivedRequests {
				idProto := req.GetAsset().GetData().GetMetadata().GetIdVersion().GetId()
				id := fmt.Sprintf("%s.%s", idProto.GetPackage(), idProto.GetName())
				wantBundle, ok := tc.wantBundles[id]
				if !ok {
					t.Errorf("runUploadESI() sent unexpected request for bundle %q", id)
					continue
				}

				gotBundle := &esipb.EsiBundle{}
				if err := req.GetAsset().GetData().GetData().UnmarshalTo(gotBundle); err != nil {
					t.Fatalf("could not unmarshal got bundle: %v", err)
				}
				if diff := cmp.Diff(wantBundle, gotBundle, protocmp.Transform()); diff != "" {
					t.Errorf("runUploadESI() sent unexpected bundle data for %q, diff (-want +got):\n%s", id, diff)
				}
			}
		})
	}
}

func TestCommonAncestor(t *testing.T) {
	tests := []struct {
		paths []string
		want  string
	}{
		{[]string{"/a/b/c", "/a/b/d"}, "/a/b"},
		{[]string{"/a/b/c"}, "/a/b"},
		{[]string{}, ""},
		{[]string{"/tmp/a/b", "/home/c/d"}, ""},
		{[]string{"/a/b/c", "/a/b/c/d"}, "/a/b"},
	}

	for _, tc := range tests {
		got := commonAncestor(tc.paths)
		if got != tc.want {
			t.Errorf("commonAncestor(%v) = %q, want %q", tc.paths, got, tc.want)
		}
	}
}

func TestToName(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"SomeVendor", "somevendor", false},
		{"some.vendor", "some_vendor", false},
		{"ESI Vendor", "esi_vendor", false},
		{"123Vendor", "a123vendor", false},
		{"Vendor@Home", "vendor_home", false},
		{"", "", true},
	}

	for _, tc := range tests {
		got, err := toName(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("toName(%q) succeeded, want error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("toName(%q) returned an unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("toName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestToPackageName(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"SomeVendor", "ai.intrinsic.ethercat.esi.somevendor", false},
		{"", "", true},
	}

	for _, tc := range tests {
		got, err := toPackageName(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("toPackageName(%q) succeeded, want error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("toPackageName(%q) returned an unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("toPackageName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestDecodeToUTF8XML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name: "iso88591",
			input: `<?xml version="1.0" encoding="ISO-8859-1"?>
	<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>`,
			want: "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n&#x9;<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>",
		},
		{
			name:  "utf8_no_decl",
			input: `<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>`,
			want:  `<EtherCATInfo><Vendor><Name>TestVendor</Name></Vendor></EtherCATInfo>`,
		},
		{
			name:    "malformed",
			input:   `<EtherCATInfo><Vendor>`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeToUTF8XML(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("decodeToUTF8XML() succeeded, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeToUTF8XML() returned an unexpected error: %v", err)
			}

			if strings.TrimSpace(got) != strings.TrimSpace(tc.want) {
				t.Errorf("decodeToUTF8XML() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		desc       string
		files      map[string]string
		inputIsDir bool // Pass directory to validateInput if true
		wantErr    bool
	}{
		{
			desc: "valid file",
			files: map[string]string{
				"valid.xml": "<EtherCATInfo></EtherCATInfo>",
			},
			wantErr: false,
		},
		{
			desc: "no files with invalid extension",
			files: map[string]string{
				"a.txt": "abc",
			},
			wantErr: true,
		},
		{
			desc:    "non-existent path",
			files:   map[string]string{},
			wantErr: true,
		},
		{
			desc: "mixed files in dir",
			files: map[string]string{
				"valid.xml": "<EtherCATInfo></EtherCATInfo>",
				"b.txt":     "abc",
			},
			inputIsDir: true,
			wantErr:    false,
		},
		{
			desc:       "empty dir",
			files:      map[string]string{},
			inputIsDir: true,
			wantErr:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tmpDir := t.TempDir()
			createTestFiles(t, tmpDir, tc.files)
			var args []string
			if tc.inputIsDir {
				args = []string{tmpDir}
			} else {
				for file := range tc.files {
					args = append(args, filepath.Join(tmpDir, file))
				}
				// If no files are specified for file list mode, test with a non-existent path
				if len(tc.files) == 0 {
					args = append(args, filepath.Join(tmpDir, "nonexistent"))
				}
			}

			_, _, err := validateInput(args)
			if tc.wantErr {
				if err == nil {
					t.Errorf("validateInput() succeeded, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("validateInput() failed: %v", err)
			}
		})
	}
}
