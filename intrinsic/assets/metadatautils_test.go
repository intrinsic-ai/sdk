// Copyright 2023 Intrinsic Innovation LLC

package metadatautils

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
	datamanifestpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	hardwaremanifestpb "intrinsic/assets/hardware_devices/proto/v1/hardware_device_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sceneobjectmanifestpb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	servicemanifestpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	skillmanifestpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

func TestValidateMetadata(t *testing.T) {
	m := &metadatapb.Metadata{
		DisplayName: "Test Skill",
		Documentation: &documentationpb.Documentation{
			Description: "Test Skill Description",
		},
		IdVersion: &idpb.IdVersion{
			Id: &idpb.Id{
				Package: "ai.intrinsic",
				Name:    "test_skill",
			},
			Version: "1.2.3",
		},
		ReleaseNotes: "Test Skill Release Notes",
		Vendor: &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		},
		AssetType: atypepb.AssetType_ASSET_TYPE_SKILL,
		UpdateTime: &tpb.Timestamp{
			Seconds: 1711177200,
			Nanos:   0,
		},
	}
	mInvalidName := proto.Clone(m).(*metadatapb.Metadata)
	mInvalidName.IdVersion.Id.Name = "_invalid_name"
	mInvalidPackage := proto.Clone(m).(*metadatapb.Metadata)
	mInvalidPackage.IdVersion.Id.Package = "_invalid_package"
	mInvalidVersion := proto.Clone(m).(*metadatapb.Metadata)
	mInvalidVersion.IdVersion.Version = "bob"
	mMissingVersion := proto.Clone(m).(*metadatapb.Metadata)
	mMissingVersion.IdVersion.Version = ""
	mNoDisplayName := proto.Clone(m).(*metadatapb.Metadata)
	mNoDisplayName.DisplayName = ""
	mNoVendor := proto.Clone(m).(*metadatapb.Metadata)
	mNoVendor.Vendor.DisplayName = ""
	mNoUpdateTime := proto.Clone(m).(*metadatapb.Metadata)
	mNoUpdateTime.UpdateTime = nil
	mNoAssetType := proto.Clone(m).(*metadatapb.Metadata)
	mNoAssetType.AssetType = atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
	mNameTooLong := proto.Clone(m).(*metadatapb.Metadata)
	mNameTooLong.IdVersion.Id.Name = strings.Repeat("a", NameCharLength[atypepb.AssetType_ASSET_TYPE_SKILL]+1)
	mDisplayNameTooLong := proto.Clone(m).(*metadatapb.Metadata)
	mDisplayNameTooLong.DisplayName = strings.Repeat("a", DisplayNameCharLength+1)
	mVersionTooLong := proto.Clone(m).(*metadatapb.Metadata)
	mVersionTooLong.IdVersion.Version = fmt.Sprintf("1.0.0+%s", strings.Repeat("a", VersionCharLength+1))
	mDescriptionTooLong := proto.Clone(m).(*metadatapb.Metadata)
	mDescriptionTooLong.Documentation.Description = strings.Repeat("a", DescriptionCharLength+1)
	mReleaseNotesTooLong := proto.Clone(m).(*metadatapb.Metadata)
	mReleaseNotesTooLong.ReleaseNotes = strings.Repeat("a", RelNotesCharLength+1)

	tests := []struct {
		name          string
		m             *metadatapb.Metadata
		opts          []ValidateMetadataOption
		wantErrorCode codes.Code
	}{
		{
			name: "valid",
			m:    m,
		},
		{
			name: "valid with version required",
			m:    m,
			opts: []ValidateMetadataOption{WithRequireVersion(true)},
		},
		{
			name: "valid with version not required",
			m:    mMissingVersion,
			opts: []ValidateMetadataOption{WithRequireVersion(false)},
		},
		{
			name:          "missing version",
			m:             mMissingVersion,
			opts:          []ValidateMetadataOption{WithRequireVersion(true)},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "missing version for catalog",
			m:             mMissingVersion,
			opts:          WithCatalogOptions(),
			wantErrorCode: codes.InvalidArgument,
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
			name:          "invalid version",
			m:             mInvalidVersion,
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
			name:          "no update time",
			m:             mNoUpdateTime,
			opts:          []ValidateMetadataOption{WithRequireUpdateTime(true)},
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "no update time for catalog",
			m:             mNoUpdateTime,
			opts:          WithCatalogOptions(),
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name:          "no asset type",
			m:             mNoAssetType,
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
		{
			name:          "version too long",
			m:             mVersionTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name:          "description too long",
			m:             mDescriptionTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name:          "release notes too long",
			m:             mReleaseNotesTooLong,
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
					t.Errorf("ValidateMetadata(%v) returned %v, want %v", tc.m, s.Code(), tc.wantErrorCode)
				}
			} else if err != nil {
				t.Errorf("ValidateMetadata(%v) = %v, want error", tc.m, err)
			}
		})
	}
}

func TestValidateManifestMetadata(t *testing.T) {
	m := &skillmanifestpb.SkillManifest{
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
	}
	mInvalidName := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mInvalidName.Id.Name = "_invalid_name"
	mInvalidPackage := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mInvalidPackage.Id.Package = "_invalid_package"
	mNoDisplayName := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mNoDisplayName.DisplayName = ""
	mNoVendor := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mNoVendor.Vendor.DisplayName = ""
	mNameTooLong := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mNameTooLong.Id.Name = strings.Repeat("a", NameCharLength[atypepb.AssetType_ASSET_TYPE_SKILL]+1)
	mDisplayNameTooLong := proto.Clone(m).(*skillmanifestpb.SkillManifest)
	mDisplayNameTooLong.DisplayName = strings.Repeat("a", DisplayNameCharLength+1)

	mService := &servicemanifestpb.ServiceMetadata{
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
	}
	mData := &datamanifestpb.DataManifest_Metadata{
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
	}
	mSceneObject := &sceneobjectmanifestpb.SceneObjectMetadata{
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
	}
	mHardwareDevice := &hardwaremanifestpb.HardwareDeviceMetadata{
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
	}

	tests := []struct {
		name          string
		m             ManifestMetadata
		wantErrorCode codes.Code
	}{
		{
			name: "valid",
			m:    m,
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
			name:          "name too long",
			m:             mNameTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name:          "display name too long",
			m:             mDisplayNameTooLong,
			wantErrorCode: codes.ResourceExhausted,
		},
		{
			name: "valid service",
			m:    mService,
		},
		{
			name: "valid data",
			m:    mData,
		},
		{
			name: "valid scene object",
			m:    mSceneObject,
		},
		{
			name: "valid hardware device",
			m:    mHardwareDevice,
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
				t.Errorf("ValidateManifestMetadata(%v) = %v, want error", tc.m, err)
			}
		})
	}
}
