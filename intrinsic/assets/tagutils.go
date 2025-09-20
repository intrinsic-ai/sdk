// Copyright 2023 Intrinsic Innovation LLC

// Package tagutils provides utilities for asset tags.
package tagutils

import (
	"fmt"
	"maps"
	"slices"

	"google.golang.org/protobuf/encoding/prototext"

	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"

	_ "embed"
)

//go:embed proto/asset_tags.textproto
var assetTags []byte

// allAssetTagMetadata is the cached return value of AllAssetTagMetadata.
var allAssetTagMetadata []*atagpb.AssetTagMetadata

// allAssetTagMetadataComputed indicates whether allAssetTagMetadata has been computed.
var allAssetTagMetadataComputed bool

// AssetTagFromName returns an AssetTag, given its name.
func AssetTagFromName(name string) (atagpb.AssetTag, error) {
	value, ok := atagpb.AssetTag_value[name]
	if !ok {
		return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, fmt.Errorf("unknown asset tag: %q", name)
	}

	return atagpb.AssetTag(value), nil
}

// AssetTagsForTypes returns a list of asset tags that apply to any of the specified asset types.
func AssetTagsForTypes(assetTypes []atypepb.AssetType) ([]atagpb.AssetTag, error) {
	tagsMap := make(map[atagpb.AssetTag]struct{})
	for _, assetType := range assetTypes {
		tagsForType, err := AssetTagsForType(assetType)
		if err != nil {
			return nil, err
		}
		for _, tag := range tagsForType {
			tagsMap[tag] = struct{}{}
		}
	}

	tags := slices.Collect(maps.Keys(tagsMap))
	slices.Sort(tags)
	return tags, nil
}

// AssetTagsForType returns a list of asset tags that apply to the specified asset type.
func AssetTagsForType(assetType atypepb.AssetType) ([]atagpb.AssetTag, error) {
	metadata, err := AssetTagMetadataForType(assetType)
	if err != nil {
		return nil, err
	}

	tags := make([]atagpb.AssetTag, len(metadata))
	for idx, tagMetadata := range metadata {
		tags[idx] = tagMetadata.GetAssetTag()
	}

	return tags, nil
}

// AssetTagMetadataForType returns a list of asset tag metadata for the specified asset type.
func AssetTagMetadataForType(assetType atypepb.AssetType) ([]*atagpb.AssetTagMetadata, error) {
	allMetadata, err := AllAssetTagMetadata()
	if err != nil {
		return nil, err
	}

	metadata := make([]*atagpb.AssetTagMetadata, 0, len(allMetadata))
	for _, tagMetadata := range allMetadata {
		if slices.Contains(tagMetadata.GetApplicableAssetTypes(), assetType) {
			metadata = append(metadata, tagMetadata)
		}
	}

	return metadata, nil
}

// AllAssetTagMetadata returns a list of all asset tag metadata.
func AllAssetTagMetadata() ([]*atagpb.AssetTagMetadata, error) {
	// Return cached value if it has been computed.
	if allAssetTagMetadataComputed {
		return allAssetTagMetadata, nil
	}

	metadataSet := &atagpb.AssetTagMetadataSet{}
	if err := prototext.Unmarshal(assetTags, metadataSet); err != nil {
		return nil, err
	}

	allAssetTagMetadata = metadataSet.GetTags()
	allAssetTagMetadataComputed = true

	return allAssetTagMetadata, nil
}
