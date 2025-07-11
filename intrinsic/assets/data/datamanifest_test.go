// Copyright 2023 Intrinsic Innovation LLC

package datamanifest

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/registryutil"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
)

func TestValidateDataManifest(t *testing.T) {
	msg := &rdspb.ReferencedDataStruct{
		Fields: map[string]*rdspb.Value{
			"foo": &rdspb.Value{
				Kind: &rdspb.Value_StringValue{
					StringValue: "bar",
				},
			},
		},
	}
	msgAny := &anypb.Any{}
	if err := msgAny.MarshalFrom(msg); err != nil {
		t.Fatalf("cannot marshal message: %v", err)
	}
	msgFDS := descriptor.FileDescriptorSetFrom(msg)
	types, err := registryutil.NewTypesFromFileDescriptorSet(msgFDS)
	if err != nil {
		t.Fatalf("cannot populate registry types: %v", err)
	}

	m := &dmpb.DataManifest{
		Metadata: &dmpb.DataManifest_Metadata{
			Id: &idpb.Id{
				Name:    "data_asset",
				Package: "ai.intrinsic",
			},
			DisplayName: "Some Data asset",
			Vendor: &vpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
		Data: msgAny,
	}

	mInvalidName := proto.Clone(m).(*dmpb.DataManifest)
	mInvalidName.GetMetadata().GetId().Name = "_invalid_name"
	mInvalidPackage := proto.Clone(m).(*dmpb.DataManifest)
	mInvalidPackage.GetMetadata().GetId().Package = "_invalid_package"
	mNoDisplayName := proto.Clone(m).(*dmpb.DataManifest)
	mNoDisplayName.GetMetadata().Vendor.DisplayName = ""
	mNoVendor := proto.Clone(m).(*dmpb.DataManifest)
	mNoVendor.GetMetadata().Vendor = nil
	mNoData := proto.Clone(m).(*dmpb.DataManifest)
	mNoData.Data = nil

	tests := []struct {
		desc    string
		m       *dmpb.DataManifest
		opts    []ValidateDataManifestOption
		wantErr bool
	}{
		{
			desc: "valid",
			m:    m,
			opts: []ValidateDataManifestOption{
				WithTypes(types),
			},
		},
		{
			desc:    "invalid name",
			m:       mInvalidName,
			wantErr: true,
		},
		{
			desc:    "invalid package",
			m:       mInvalidPackage,
			wantErr: true,
		},
		{
			desc:    "no display name",
			m:       mNoDisplayName,
			wantErr: true,
		},
		{
			desc:    "no vendor",
			m:       mNoVendor,
			wantErr: true,
		},
		{
			desc:    "no data",
			m:       mNoData,
			wantErr: true,
		},
		{
			desc: "missing proto",
			m:    m,
			opts: []ValidateDataManifestOption{
				WithTypes(&protoregistry.Types{}),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := ValidateDataManifest(tc.m, tc.opts...)
			if tc.wantErr && err == nil {
				t.Error("ValidateDataManifest() succeeded, want error")
			} else if !tc.wantErr && err != nil {
				t.Errorf("ValidateDataManifest() failed: %v", err)
			}
		})
	}
}
