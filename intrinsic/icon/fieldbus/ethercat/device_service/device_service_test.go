// Copyright 2023 Intrinsic Innovation LLC

package deviceservice

import (
	"context"
	"errors"
	"intrinsic/assets/data/fakedataassets"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	anypb "google.golang.org/protobuf/types/known/anypb"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dagrpcpb "intrinsic/assets/data/proto/v1/data_assets_go_grpc_proto"
	ipb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	rdpb "intrinsic/assets/proto/v1/resolved_dependency_go_proto"
	dscpb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_config_go_proto"
	dspb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_go_grpc_proto"
	esipb "intrinsic/icon/fieldbus/ethercat/device_service/v1/esi_go_proto"
)

// TestDeviceService tests the creation of the DeviceService and the GetConfiguration call.
func TestDeviceService(t *testing.T) {
	bundleID := &ipb.Id{
		Package: "intrinsic_proto.fieldbus.ethercat.test",
		Name:    "test_bundle_1",
	}
	iface := "data://" + esiBundleDataAssetProtoName
	config := &dscpb.DeviceServiceConfig{
		DeviceIdentifier: &dscpb.DeviceIdentifier{
			VendorId:    0x0001,
			ProductCode: 0x0002,
			Revision:    0x0003,
		},
		EsiBundle: &rdpb.ResolvedDependency{
			Interfaces: map[string]*rdpb.ResolvedDependency_Interface{
				iface: {
					Protocol: &rdpb.ResolvedDependency_Interface_Data_{
						Data: &rdpb.ResolvedDependency_Interface_Data{
							Id: bundleID,
						},
					},
				},
			},
		},
	}

	bundle := &esipb.EsiBundle{
		Files: map[string]*esipb.Esi{
			"path/to/esi/file.esi": {Data: "test esi data"},
		},
	}
	bundleAny, err := anypb.New(bundle)
	if err != nil {
		t.Fatalf("Failed to marshal ESI bundle: %v", err)
	}
	validDataAsset := &dapb.DataAsset{
		Metadata: &mpb.Metadata{IdVersion: &ipb.IdVersion{Id: bundleID}},
		Data:     bundleAny,
	}

	bundleAbsPath := &esipb.EsiBundle{
		Files: map[string]*esipb.Esi{
			"/path/to/esi/file.esi": {Data: "test esi data"},
		},
	}
	bundleAbsPathAny, err := anypb.New(bundleAbsPath)
	if err != nil {
		t.Fatalf("Failed to marshal ESI bundle with absolute path: %v", err)
	}
	absPathDataAsset := &dapb.DataAsset{
		Metadata: &mpb.Metadata{IdVersion: &ipb.IdVersion{Id: bundleID}},
		Data:     bundleAbsPathAny,
	}

	bundleDotDotPath := &esipb.EsiBundle{
		Files: map[string]*esipb.Esi{
			"../file.esi": {Data: "test esi data"},
		},
	}
	bundleDotDotPathAny, err := anypb.New(bundleDotDotPath)
	if err != nil {
		t.Fatalf("Failed to marshal ESI bundle with .. in path: %v", err)
	}
	dotDotPathDataAsset := &dapb.DataAsset{
		Metadata: &mpb.Metadata{IdVersion: &ipb.IdVersion{Id: bundleID}},
		Data:     bundleDotDotPathAny,
	}

	wrongProtoDataAsset := &dapb.DataAsset{
		Metadata: &mpb.Metadata{IdVersion: &ipb.IdVersion{Id: bundleID}},
		Data:     &anypb.Any{Value: []byte("invalid data"), TypeUrl: "type.googleapis.com/google.protobuf.Duration"},
	}

	type testArgs struct {
		config     *dscpb.DeviceServiceConfig
		dataAssets []*dapb.DataAsset
	}

	tests := []struct {
		desc         string
		testArgs     testArgs
		wantResponse *dspb.GetConfigurationResponse
		wantErr      error
	}{
		{
			desc: "valid config",
			testArgs: testArgs{
				config:     config,
				dataAssets: []*dapb.DataAsset{validDataAsset},
			},
			wantResponse: &dspb.GetConfigurationResponse{
				DeviceServiceConfig: config,
				EsiBundle:           bundle,
			},
		},
		{
			desc: "data asset not found",
			testArgs: testArgs{
				config:     config,
				dataAssets: []*dapb.DataAsset{},
			},
			wantErr: ErrEsiBundleNotFound,
		},
		{
			desc: "wrong proto type in data asset",
			testArgs: testArgs{
				config:     config,
				dataAssets: []*dapb.DataAsset{wrongProtoDataAsset},
			},
			wantErr: ErrEsiUnmarshal,
		},
		{
			desc: "bundle with absolute path",
			testArgs: testArgs{
				config:     config,
				dataAssets: []*dapb.DataAsset{absPathDataAsset},
			},
			wantErr: ErrEsiInvalidPath,
		},
		{
			desc: "bundle with .. at the start of the path",
			testArgs: testArgs{
				config:     config,
				dataAssets: []*dapb.DataAsset{dotDotPathDataAsset},
			},
			wantErr: ErrEsiInvalidPath,
		},
		{
			desc: "nil config",
			testArgs: testArgs{
				config: nil,
			},
			wantErr: ErrConfigNil,
		},
		{
			desc: "no device identifier",
			testArgs: testArgs{
				config: &dscpb.DeviceServiceConfig{},
			},
			wantErr: ErrNoDeviceIdentifier,
		},
		{
			desc: "no esi bundle data asset id in config",
			testArgs: testArgs{
				config: &dscpb.DeviceServiceConfig{
					DeviceIdentifier: &dscpb.DeviceIdentifier{},
				},
			},
			wantErr: ErrEsiBundleNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := context.Background()
			var client dagrpcpb.DataAssetsClient
			if tc.testArgs.dataAssets != nil {
				fakeDA := fakedataassets.StartServer(ctx, t, fakedataassets.WithDataAssets(tc.testArgs.dataAssets))
				client = fakeDA.Client
			}

			service, err := NewDeviceService(ctx, tc.testArgs.config, client)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("NewDeviceService(...) = nil error, want non-nil error matching %v", tc.wantErr)
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("NewDeviceService(...) = %v, want error to wrap %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("NewDeviceService(...) = %v, want nil error", err)
			}
			// If we expected an error, and got it, we are done with this test case.
			if tc.wantErr != nil {
				return
			}
			// If tc.wantErr was nil, err must also be nil at this point. Proceed to check response.
			if tc.wantResponse != nil {
				if service == nil {
					t.Fatalf("NewDeviceService() returned nil service, expected a valid service")
				}

				gotConfig, err := service.GetConfiguration(context.Background(), &dspb.GetConfigurationRequest{})
				if err != nil {
					t.Fatalf("GetConfiguration() returned unexpected error: %v", err)
				}

				if diff := cmp.Diff(tc.wantResponse, gotConfig, protocmp.Transform()); diff != "" {
					t.Errorf("GetConfiguration() returned diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
