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
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
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
			err := ValidateMetadata(tc.m)
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
