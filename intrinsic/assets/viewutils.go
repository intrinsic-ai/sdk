// Copyright 2023 Intrinsic Innovation LLC

// Package viewutils provides utilities for working with asset views.
package viewutils

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
	viewpb "intrinsic/assets/proto/view_go_proto"
)

var (
	metadataFieldName              = "metadata"
	assetSpecificMetadataOneOfName = "asset_specific_metadata"
	deploymentDataFieldName        = "deployment_data"

	errInvalidAssetType = errors.New("invalid asset type")
)

// Asset is an interface for constructing asset views.
type Asset interface {
	proto.Message

	GetMetadata() *metadatapb.Metadata
}

type customViewValue struct {
	// Field is the name of the field to set.
	Field string
	// Value is a function that returns the value to set.
	Value func() (proto.Message, error)
	// Views is the list of views for which the value should be set.
	Views []viewpb.AssetViewType
}

type assetToViewOptions struct {
	// CustomViewValues are custom values to set for some asset views.
	CustomViewValues []customViewValue
	// DeploymentData returns the deployment data for the asset view.
	//
	// This option is a function rather than a message so the value doesn't need to be computed unless
	// it is actually needed.
	//
	// If provided, the returned deployment data is used to populate the deployment data instead of
	// the deployment data from the input asset.
	DeploymentData func() (proto.Message, error)
	// FileDescriptorSet returns a FileDescriptorSet for the asset metadata.
	//
	// This option is a function rather than a message so the value doesn't need to be computed unless
	// it is actually needed.
	//
	// If provided, the returned FileDescriptorSet is used to populate the metadata instead of the
	// FileDescriptorSet from the input metadata.
	FileDescriptorSet func() (*dpb.FileDescriptorSet, error)
	// Metadata is the metadata to use for the asset view, instead of the input metadata.
	Metadata *metadatapb.Metadata
}

// AssetToViewOption is an option for AssetToView.
type AssetToViewOption func(*assetToViewOptions)

// WithCustomViewValue returns an AssetToViewOption that specifies a custom value to set for some
// asset views.
func WithCustomViewValue(field string, views []viewpb.AssetViewType, value func() (proto.Message, error)) AssetToViewOption {
	return func(opts *assetToViewOptions) {
		opts.CustomViewValues = append(opts.CustomViewValues, customViewValue{
			Field: field,
			Value: value,
			Views: views,
		})
	}
}

// WithDeploymentData returns an AssetToViewOption that specifies a function to compute the
// DeploymentData.
func WithDeploymentData(dd func() (proto.Message, error)) AssetToViewOption {
	return func(opts *assetToViewOptions) {
		opts.DeploymentData = dd
	}
}

// WithFileDescriptorSet returns an AssetToViewOption that specifies a function to compute the
// FileDescriptorSet.
func WithFileDescriptorSet(fds func() (*dpb.FileDescriptorSet, error)) AssetToViewOption {
	return func(opts *assetToViewOptions) {
		opts.FileDescriptorSet = fds
	}
}

// WithMetadata returns an AssetToViewOption that specifies the Metadata to use.
func WithMetadata(md *metadatapb.Metadata) AssetToViewOption {
	return func(opts *assetToViewOptions) {
		opts.Metadata = md
	}
}

// AssetToView returns the specified view of an Asset.
func AssetToView[T Asset](asset T, view viewpb.AssetViewType, options ...AssetToViewOption) (T, error) {
	opts := assetToViewOptions{}
	for _, option := range options {
		option(&opts)
	}
	if opts.Metadata == nil {
		opts.Metadata = proto.Clone(asset.GetMetadata()).(*metadatapb.Metadata)
	}

	assetView := newMessage(asset)
	if err := setFieldValue(assetView, metadataFieldName, &metadatapb.Metadata{}); err != nil {
		return assetView, err
	}

	// Fail early if the input type doesn't have the required fields.
	if asset.ProtoReflect().Descriptor().Oneofs().ByName(protoreflect.Name(assetSpecificMetadataOneOfName)) == nil {
		return assetView, fmt.Errorf("%w: %T does not have oneof %q", errInvalidAssetType, asset, assetSpecificMetadataOneOfName)
	}
	if asset.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(deploymentDataFieldName)) == nil {
		return assetView, fmt.Errorf("%w: %T does not have field %q", errInvalidAssetType, asset, deploymentDataFieldName)
	}
	for _, cv := range opts.CustomViewValues {
		if asset.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(cv.Field)) == nil {
			return assetView, fmt.Errorf("%w: %T does not have field %q", errInvalidAssetType, asset, cv.Field)
		}
	}

	switch view {
	case viewpb.AssetViewType_ASSET_VIEW_TYPE_BASIC:
		assetView.GetMetadata().AssetType = opts.Metadata.GetAssetType()
		assetView.GetMetadata().IdVersion = opts.Metadata.GetIdVersion()
	case viewpb.AssetViewType_ASSET_VIEW_TYPE_DETAIL:
		assetView.GetMetadata().AssetTag = opts.Metadata.GetAssetTag()
		assetView.GetMetadata().AssetType = opts.Metadata.GetAssetType()
		assetView.GetMetadata().DisplayName = opts.Metadata.GetDisplayName()
		assetView.GetMetadata().Documentation = opts.Metadata.GetDocumentation()
		assetView.GetMetadata().IdVersion = opts.Metadata.GetIdVersion()
		assetView.GetMetadata().Vendor = opts.Metadata.GetVendor()
		if err := copyOneOfFieldValue(asset, assetView, assetSpecificMetadataOneOfName); err != nil {
			return assetView, err
		}
	case viewpb.AssetViewType_ASSET_VIEW_TYPE_VERSIONS:
		assetView.GetMetadata().AssetType = opts.Metadata.GetAssetType()
		assetView.GetMetadata().IdVersion = opts.Metadata.GetIdVersion()
		assetView.GetMetadata().ReleaseNotes = opts.Metadata.GetReleaseNotes()
		assetView.GetMetadata().ReleaseTag = opts.Metadata.GetReleaseTag()
		assetView.GetMetadata().UpdateTime = opts.Metadata.GetUpdateTime()
		assetView.GetMetadata().Vendor = opts.Metadata.GetVendor()
	case viewpb.AssetViewType_ASSET_VIEW_TYPE_ALL_METADATA:
		if err := mergeAllMetadataToView(asset, assetView, opts); err != nil {
			return assetView, err
		}
	case viewpb.AssetViewType_ASSET_VIEW_TYPE_ALL:
		if err := mergeAllMetadataToView(asset, assetView, opts); err != nil {
			return assetView, err
		}
		if opts.DeploymentData != nil {
			dd, err := opts.DeploymentData()
			if err != nil {
				return assetView, err
			}
			if dd != nil {
				if err := setFieldValue(assetView, deploymentDataFieldName, dd); err != nil {
					return assetView, err
				}
			}
		} else if err := copyFieldValue(asset, assetView, deploymentDataFieldName); err != nil {
			return assetView, err
		}
	default:
		return assetView, status.Errorf(codes.InvalidArgument, "unsupported asset view type %v", view.String())
	}

	// Set custom view values.
	for _, cv := range opts.CustomViewValues {
		if slices.Contains(cv.Views, view) {
			if value, err := cv.Value(); err != nil {
				return assetView, err
			} else if err := setFieldValue(assetView, cv.Field, value); err != nil {
				return assetView, err
			}
		}
	}

	return assetView, nil
}

func mergeAllMetadataToView[T Asset](asset T, assetView T, opts assetToViewOptions) error {
	if opts.Metadata != nil {
		if err := setFieldValue(assetView, metadataFieldName, opts.Metadata); err != nil {
			return err
		}
	}

	if opts.FileDescriptorSet != nil {
		fds, err := opts.FileDescriptorSet()
		if err != nil {
			return err
		}
		assetView.GetMetadata().FileDescriptorSet = fds
	}

	if err := copyOneOfFieldValue(asset, assetView, assetSpecificMetadataOneOfName); err != nil {
		return err
	}

	return nil
}

// EnumFromString returns an AssetViewType enum from a string.
func EnumFromString(t string) (viewpb.AssetViewType, error) {
	if i, exists := viewpb.AssetViewType_value[t]; exists {
		return viewpb.AssetViewType(i), nil
	}
	return viewpb.AssetViewType_ASSET_VIEW_TYPE_UNSPECIFIED, fmt.Errorf("invalid view: %q", t)
}

// EnumFromShortString returns an AssetViewType enum from a short string. The
// empty string maps to ASSET_VIEW_TYPE_UNSPECIFIED.
func EnumFromShortString(t string) (viewpb.AssetViewType, error) {
	if t == "" {
		return viewpb.AssetViewType_ASSET_VIEW_TYPE_UNSPECIFIED, nil
	}
	return EnumFromString(fmt.Sprintf("ASSET_VIEW_TYPE_%s", strings.ToUpper(t)))
}

// ShortStringFromEnum returns a short string from an AssetViewType enum.
func ShortStringFromEnum(t viewpb.AssetViewType) string {
	return strings.ToLower(strings.TrimPrefix(t.String(), "ASSET_VIEW_TYPE_"))
}

func setFieldValue(m proto.Message, name string, value proto.Message) error {
	mR := m.ProtoReflect()

	fd := mR.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		return fmt.Errorf("field %q not found in message %v", name, mR.Descriptor().FullName())
	}

	valueR := protoreflect.ValueOf(value.ProtoReflect())
	if value == nil || !valueR.Message().IsValid() {
		mR.Clear(fd)
		return nil
	}

	mR.Set(fd, valueR)

	return nil
}

func copyFieldValue[T proto.Message](src T, dst T, name string) error {
	srcR := src.ProtoReflect()
	dstR := dst.ProtoReflect()

	fd := srcR.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		return fmt.Errorf("field %q not found in message %v", name, srcR.Descriptor().FullName())
	}

	if !srcR.Has(fd) {
		if dstR.Has(fd) {
			dstR.Clear(fd)
		}
		return nil
	}

	dstR.Set(fd, srcR.Get(fd))

	return nil
}

func copyOneOfFieldValue[T proto.Message](src T, dst T, name string) error {
	srcR := src.ProtoReflect()
	dstR := dst.ProtoReflect()

	od := srcR.Descriptor().Oneofs().ByName(protoreflect.Name(name))
	if od == nil {
		return fmt.Errorf("oneof %q not found in message %v", name, srcR.Descriptor().FullName())
	}

	active := srcR.WhichOneof(od)
	if active == nil {
		activeDst := dstR.WhichOneof(od)
		if activeDst != nil {
			dstR.Clear(activeDst)
		}
		return nil
	}

	fd := srcR.Descriptor().Fields().ByName(protoreflect.Name(active.Name()))
	if fd == nil {
		return fmt.Errorf("field %q not found in message %v", active.Name(), srcR.Descriptor().FullName())
	}
	dstR.Set(fd, srcR.Get(fd))

	return nil
}

func newMessage[T proto.Message](m T) T {
	return m.ProtoReflect().New().Interface().(T)
}
