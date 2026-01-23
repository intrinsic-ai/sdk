// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides testing utils for Data Assets.
package utils

import (
	"testing"

	"intrinsic/util/proto/descriptor"

	"google.golang.org/protobuf/proto"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type makeDataManifestOptions struct {
	metadata *dmpb.DataManifest_Metadata
	payload  *anypb.Any
}

// MakeDataManifestOption is an option for MakeDataManifest.
type MakeDataManifestOption func(*makeDataManifestOptions)

// WithMetadata specifies the metadata to use for the DataManifest.
func WithMetadata(metadata *dmpb.DataManifest_Metadata) MakeDataManifestOption {
	return func(opts *makeDataManifestOptions) {
		opts.metadata = metadata
	}
}

// WithPayload converts the specified data payload to an Any and uses it as the DataManifest's
// payload.
func WithPayload(t *testing.T, payload proto.Message) MakeDataManifestOption {
	payloadAny, err := anypb.New(payload)
	if err != nil {
		t.Fatalf("anypb.New(%v) returned unexpected error: %v", payload, err)
	}
	return WithPayloadAny(payloadAny)
}

// WithPayloadAny specifies the data payload Any to use for the DataManifest.
func WithPayloadAny(payload *anypb.Any) MakeDataManifestOption {
	return func(opts *makeDataManifestOptions) {
		opts.payload = payload
	}
}

// MakeDataManifest makes a DataManifest for testing.
func MakeDataManifest(t *testing.T, options ...MakeDataManifestOption) *dmpb.DataManifest {
	t.Helper()

	opts := &makeDataManifestOptions{
		metadata: &dmpb.DataManifest_Metadata{
			Id: &idpb.Id{
				Name:    "some_data_asset",
				Package: "package.some",
			},
			DisplayName: "Some Data Asset",
			Documentation: &documentationpb.Documentation{
				Description: "Some documentation",
			},
			Vendor: &vpb.Vendor{
				DisplayName: "Some Company",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	return &dmpb.DataManifest{
		Metadata: opts.metadata,
		Data:     opts.payload,
	}
}

type makeDataAssetOptions struct {
	fileDescriptorSet *dpb.FileDescriptorSet
	metadata          *mpb.Metadata
	payload           *anypb.Any
}

// MakeDataAssetOption is an option for MakeDataAsset.
type MakeDataAssetOption func(*makeDataAssetOptions)

// WithDataAssetFileDescriptorSet specifies the FileDescriptorSet to use for the DataAsset.
func WithDataAssetFileDescriptorSet(fds *dpb.FileDescriptorSet) MakeDataAssetOption {
	return func(opts *makeDataAssetOptions) {
		opts.fileDescriptorSet = fds
	}
}

// WithDataAssetMetadata specifies the Metadata to use for the DataAsset.
func WithDataAssetMetadata(metadata *mpb.Metadata) MakeDataAssetOption {
	return func(opts *makeDataAssetOptions) {
		opts.metadata = metadata
	}
}

// WithDataAssetPayload specifies the data payload to use for the DataAsset and also sets the
// corresponding FileDescriptorSet.
func WithDataAssetPayload(t *testing.T, payload proto.Message) MakeDataAssetOption {
	payloadAny, err := anypb.New(payload)
	if err != nil {
		t.Fatalf("anypb.New(%v) returned unexpected error: %v", payload, err)
	}
	return func(opts *makeDataAssetOptions) {
		WithDataAssetFileDescriptorSet(descriptor.FileDescriptorSetFrom(payload))(opts)
		WithDataAssetPayloadAny(payloadAny)(opts)
	}
}

// WithDataAssetPayloadAny specifies the data payload Any to use for the DataAsset.
func WithDataAssetPayloadAny(payload *anypb.Any) MakeDataAssetOption {
	return func(opts *makeDataAssetOptions) {
		opts.payload = payload
	}
}

// MakeDataAsset makes a DataAsset for testing.
func MakeDataAsset(t *testing.T, options ...MakeDataAssetOption) *dapb.DataAsset {
	t.Helper()

	defaultPayload := &rdspb.ReferencedDataStruct{
		Fields: map[string]*rdspb.Value{
			"foo": {
				Kind: &rdspb.Value_StringValue{
					StringValue: "bar",
				},
			},
		},
	}
	defaultPayloadAny, err := anypb.New(defaultPayload)
	if err != nil {
		t.Fatalf("anypb.New(%v) returned unexpected error: %v", defaultPayload, err)
	}

	opts := &makeDataAssetOptions{
		fileDescriptorSet: descriptor.FileDescriptorSetFrom(defaultPayload),
		metadata: &mpb.Metadata{
			AssetType:   atpb.AssetType_ASSET_TYPE_DATA,
			DisplayName: "Some Data Asset",
			Documentation: &documentationpb.Documentation{
				Description: "Some documentation",
			},
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Name:    "some_data_asset",
					Package: "package.some",
				},
			},
			Vendor: &vpb.Vendor{
				DisplayName: "Some Company",
			},
		},
		payload: defaultPayloadAny,
	}
	for _, opt := range options {
		opt(opts)
	}

	return &dapb.DataAsset{
		FileDescriptorSet: opts.fileDescriptorSet,
		Metadata:          opts.metadata,
		Data:              opts.payload,
	}
}
