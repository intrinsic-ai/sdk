// Copyright 2023 Intrinsic Innovation LLC

package metadatautils

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	datamanifestpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	hardwaremanifestpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	processmanifestpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sceneobjectmanifestpb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	servicemanifestpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	"intrinsic/assets/typeutils"
	skillmanifestpb "intrinsic/skills/proto/skill_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

var manifestMetadata = map[atypepb.AssetType]ManifestMetadata{
	atypepb.AssetType_ASSET_TYPE_DATA: &datamanifestpb.DataManifest_Metadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_data",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Data",
		Documentation: &documentationpb.Documentation{
			Description: "Test Data Description",
		},
	},
	atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE: &hardwaremanifestpb.HardwareDeviceMetadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_hardware_device",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Hardware Device",
		Documentation: &documentationpb.Documentation{
			Description: "Test Hardware Device Description",
		},
	},
	atypepb.AssetType_ASSET_TYPE_PROCESS: &processmanifestpb.ProcessMetadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_process",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Process",
		Documentation: &documentationpb.Documentation{
			Description: "Test Process Description",
		},
	},
	atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT: &sceneobjectmanifestpb.SceneObjectMetadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_scene_object",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Scene Object",
		Documentation: &documentationpb.Documentation{
			Description: "Test Scene Object Description",
		},
	},
	atypepb.AssetType_ASSET_TYPE_SERVICE: &servicemanifestpb.ServiceMetadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_service",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Service",
		Documentation: &documentationpb.Documentation{
			Description: "Test Service Description",
		},
	},
	atypepb.AssetType_ASSET_TYPE_SKILL: &skillmanifestpb.SkillManifest{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "test_skill",
		},
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		DisplayName: "Test Skill",
		Documentation: &documentationpb.Documentation{
			Description: "Test Skill Description",
		},
	},
}

func metadataWithAllFields() *metadatapb.Metadata {
	return &metadatapb.Metadata{
		IdVersion: &idpb.IdVersion{
			Id: &idpb.Id{
				Package: "ai.intrinsic",
				Name:    "test_service",
			},
			Version: "1.2.3",
		},
		DisplayName: "Test Service",
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		Documentation: &documentationpb.Documentation{
			Description: "Test Service Description",
		},
		ReleaseNotes: "Test Service Release Notes",
		UpdateTime: &tpb.Timestamp{
			Seconds: 1711177200,
			Nanos:   0,
		},
		AssetType: atypepb.AssetType_ASSET_TYPE_SERVICE,
		Provides: []*metadatapb.Interface{
			{
				Uri: "grpc://intrinsic_proto.test.MyService",
			},
		},
		FileDescriptorSet: &dpb.FileDescriptorSet{},
	}
}

func metadataNoOutputOnlyFields() *metadatapb.Metadata {
	m := metadataWithAllFields()
	m.FileDescriptorSet = nil
	m.Provides = nil

	return m
}

func metadataInAsset() *metadatapb.Metadata {
	m := metadataWithAllFields()
	m.FileDescriptorSet = nil
	m.Provides = nil
	m.ReleaseNotes = ""
	m.UpdateTime = nil
	m.GetIdVersion().Version = ""

	return m
}

func TestValidateMetadata(t *testing.T) {
	tests := []struct {
		name          string
		m             *metadatapb.Metadata
		opts          []ValidateMetadataOption
		wantErrorCode codes.Code
	}{
		{
			name: "valid",
			m:    metadataWithAllFields(),
		},
		{
			name: "invalid version",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Version = "bob"
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "missing version",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Version = ""
				return m
			}(),
		},
		{
			name: "valid for catalog",
			m:    metadataWithAllFields(),
			opts: []ValidateMetadataOption{WithCatalogOptions()},
		},
		{
			name: "missing version for catalog",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Version = ""
				return m
			}(),
			opts:          []ValidateMetadataOption{WithCatalogOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no update time for catalog",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.UpdateTime = nil
				return m
			}(),
			opts:          []ValidateMetadataOption{WithCatalogOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "valid for in asset",
			m:    metadataInAsset(),
			opts: []ValidateMetadataOption{WithInAssetOptions()},
		},
		{
			name: "in asset with file descriptor set",
			m: func() *metadatapb.Metadata {
				m := metadataInAsset()
				m.FileDescriptorSet = &dpb.FileDescriptorSet{}
				return m
			}(),
			opts:          []ValidateMetadataOption{WithInAssetOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "in asset with provides",
			m: func() *metadatapb.Metadata {
				m := metadataInAsset()
				m.Provides = []*metadatapb.Interface{{
					Uri: "grpc://intrinsic_proto.test.MyService",
				}}
				return m
			}(),
			opts:          []ValidateMetadataOption{WithInAssetOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "in asset with release notes",
			m: func() *metadatapb.Metadata {
				m := metadataInAsset()
				m.ReleaseNotes = "i am released!"
				return m
			}(),
			opts:          []ValidateMetadataOption{WithInAssetOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "in asset with update time",
			m: func() *metadatapb.Metadata {
				m := metadataInAsset()
				m.UpdateTime = &tpb.Timestamp{}
				return m
			}(),
			opts:          []ValidateMetadataOption{WithInAssetOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "in asset with version",
			m: func() *metadatapb.Metadata {
				m := metadataInAsset()
				m.GetIdVersion().Version = "1.0.0"
				return m
			}(),
			opts:          []ValidateMetadataOption{WithInAssetOptions()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "invalid name",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Id.Name = "_invalid_name"
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "invalid package",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Id.Package = "_invalid_package"
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no display name",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.DisplayName = ""
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no vendor",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.Vendor.DisplayName = ""
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no asset type",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.AssetType = atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "wrong asset type",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.AssetType = atypepb.AssetType_ASSET_TYPE_SERVICE
				return m
			}(),
			opts:          []ValidateMetadataOption{WithAssetType(atypepb.AssetType_ASSET_TYPE_PROCESS)},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "invalid asset tag",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.AssetType = atypepb.AssetType_ASSET_TYPE_SERVICE
				m.AssetTag = atagpb.AssetTag_ASSET_TAG_SUBPROCESS
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "invalid provides interface",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.Provides = []*metadatapb.Interface{{Uri: "invalid"}}
				return m
			}(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no output-only fields required and absent",
			m:    metadataNoOutputOnlyFields(),
			opts: []ValidateMetadataOption{WithNoOutputOnlyFields()},
		},
		{
			name: "no output-only fields required and file descriptor set present",
			m: func() *metadatapb.Metadata {
				m := metadataNoOutputOnlyFields()
				m.FileDescriptorSet = &dpb.FileDescriptorSet{}
				return m
			}(),
			opts:          []ValidateMetadataOption{WithNoOutputOnlyFields()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "no output-only fields required and provides present",
			m: func() *metadatapb.Metadata {
				m := metadataNoOutputOnlyFields()
				m.Provides = []*metadatapb.Interface{{Uri: "grpc://intrinsic_proto.test.MyService"}}
				return m
			}(),
			opts:          []ValidateMetadataOption{WithNoOutputOnlyFields()},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "name too long",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.AssetType = atypepb.AssetType_ASSET_TYPE_SERVICE
				m.IdVersion.Id.Name = strings.Repeat("a", NameCharLength[atypepb.AssetType_ASSET_TYPE_SERVICE]+1)
				return m
			}(),
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name: "display name too long",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.DisplayName = strings.Repeat("a", DisplayNameCharLength+1)
				return m
			}(),
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name: "version too long",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.IdVersion.Version = fmt.Sprintf("1.0.0+%s", strings.Repeat("a", VersionCharLength+1))
				return m
			}(),
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name: "description too long",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.Documentation.Description = strings.Repeat("a", DescriptionCharLength+1)
				return m
			}(),
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name: "release notes too long",
			m: func() *metadatapb.Metadata {
				m := metadataWithAllFields()
				m.ReleaseNotes = strings.Repeat("a", RelNotesCharLength+1)
				return m
			}(),
			wantErrorCode: codes.ResourceExhausted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMetadata(tc.m, tc.opts...)
			if tc.wantErrorCode != codes.OK {
				if err == nil {
					t.Errorf("ValidateMetadata(%v) = nil, want error", tc.m)
				} else if s, ok := status.FromError(err); !ok {
					t.Errorf("Could not get status from ValidateMetadata()")
				} else if s.Code() != tc.wantErrorCode {
					t.Errorf("ValidateMetadata(%v) returned %q, want %v", tc.m, err, tc.wantErrorCode)
				}
			} else if err != nil {
				t.Errorf("ValidateMetadata(%v) = %v, want no error", tc.m, err)
			}
		})
	}
}

func TestValidateManifestMetadata(t *testing.T) {
	mSkill := manifestMetadata[atypepb.AssetType_ASSET_TYPE_SKILL]
	mInvalidName := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mInvalidName.Id.Name = "_invalid_name"
	mInvalidPackage := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mInvalidPackage.Id.Package = "_invalid_package"
	mNoDisplayName := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mNoDisplayName.DisplayName = ""
	mNoVendor := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mNoVendor.Vendor.DisplayName = ""
	mNameTooLong := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mNameTooLong.Id.Name = strings.Repeat("a", NameCharLength[atypepb.AssetType_ASSET_TYPE_SKILL]+1)
	mDisplayNameTooLong := proto.Clone(mSkill).(*skillmanifestpb.SkillManifest)
	mDisplayNameTooLong.DisplayName = strings.Repeat("a", DisplayNameCharLength+1)

	mService := manifestMetadata[atypepb.AssetType_ASSET_TYPE_SERVICE]
	mInvalidAssetTag := proto.Clone(mService).(*servicemanifestpb.ServiceMetadata)
	mInvalidAssetTag.AssetTag = atagpb.AssetTag_ASSET_TAG_SUBPROCESS

	tests := []struct {
		name          string
		m             ManifestMetadata
		wantErrorCode codes.Code
	}{
		{
			name: "valid data",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_DATA],
		},
		{
			name: "valid hardware device",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE],
		},
		{
			name: "valid process",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_PROCESS],
		},
		{
			name: "valid scene object",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT],
		},
		{
			name: "valid service",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_SERVICE],
		},
		{
			name: "valid skill",
			m:    manifestMetadata[atypepb.AssetType_ASSET_TYPE_SKILL],
		},
		{
			name:          "invalid name",
			m:             mInvalidName,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "invalid package",
			m:             mInvalidPackage,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "no display name",
			m:             mNoDisplayName,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "no vendor",
			m:             mNoVendor,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "invalid asset tag",
			m:             mInvalidAssetTag,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "name too long",
			m:             mNameTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name:          "display name too long",
			m:             mDisplayNameTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateManifestMetadata(tc.m)
			if tc.wantErrorCode != codes.OK {
				if err == nil {
					t.Errorf("ValidateManifestMetadata(%v) = nil, want error", tc.m)
				} else if s, ok := status.FromError(err); !ok {
					t.Errorf("Could not get status from ValidateManifestMetadata()")
				} else if s.Code() != tc.wantErrorCode {
					t.Errorf("ValidateManifestMetadata(%v) returned %v, want %v", tc.m, s.Code(), tc.wantErrorCode)
				}
			} else if err != nil {
				t.Errorf("ValidateManifestMetadata(%v) = %v, want no error", tc.m, err)
			}
		})
	}
}

func TestValidateManifestMetadataSupportsAllAssetTypes(t *testing.T) {
	for _, at := range typeutils.AllAssetTypes() {
		t.Run(at.String(), func(t *testing.T) {
			m, ok := manifestMetadata[at]
			if !ok {
				t.Errorf("No test manifest metadata for asset type %v", at)
			}
			if err := ValidateManifestMetadata(m); err != nil {
				t.Errorf("ValidateManifestMetadata(%v) = %v, want no error", m, err)
			}
		})
	}
}

func testMetadataNoOutputOnlyFields(t *testing.T) {
	m := metadataNoOutputOnlyFields()
	if err := ValidateMetadata(m,
		WithNoOutputOnlyFields(),
	); err != nil {
		t.Errorf("ValidateMetadata(%v) = %v for metadataNoOutputOnlyFields, want no error", m, err)
	}
}

func TestToInputMetadata(t *testing.T) {
	m := metadataWithAllFields()
	mInput := ToInputMetadata(m)
	if err := ValidateMetadata(mInput,
		WithNoOutputOnlyFields(),
	); err != nil {
		t.Errorf("ValidateMetadata(%v) = %v after calling ToInputMetadata, want no error", m, err)
	}
}

func TestMetadataInAsset(t *testing.T) {
	m := metadataInAsset()
	if err := ValidateMetadata(m,
		WithInAssetOptions(),
	); err != nil {
		t.Errorf("ValidateMetadata(%v) = %v for metadataInAsset, want no error", m, err)
	}
}

func TestToInAssetMetadata(t *testing.T) {
	m := metadataWithAllFields()
	mInAsset := ToInAssetMetadata(m)
	if err := ValidateMetadata(mInAsset,
		WithInAssetOptions(),
	); err != nil {
		t.Errorf("ValidateMetadata(%v) = %v after calling ToInAssetMetadata, want no error", m, err)
	}
}
