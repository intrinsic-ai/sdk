// Copyright 2023 Intrinsic Innovation LLC

// Package tagutils provides utilities for AssetTags.
package tagutils

import (
	"fmt"
	"maps"
	"slices"
	"strings"

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

// AllAssetTags returns all AssetTags except ASSET_TAG_UNSPECIFIED.
func AllAssetTags() []atagpb.AssetTag {
	allTags := make([]atagpb.AssetTag, 0, len(atagpb.AssetTag_name)-1)
	for tag := range atagpb.AssetTag_name {
		if tag != int32(atagpb.AssetTag_ASSET_TAG_UNSPECIFIED) {
			allTags = append(allTags, atagpb.AssetTag(tag))
		}
	}
	return allTags
}

// AssetTagDisplayName returns the display name of the given AssetTag, or "Unknown" if the tag is
// not known.
func AssetTagDisplayName(tag atagpb.AssetTag) string {
	metadata, err := AssetTagMetadataForTag(tag)
	if err != nil {
		return "Unknown"
	}
	return metadata.GetDisplayName()
}

// AssetTagFromDisplayName returns the AssetTag for the given display name, or
// ASSET_TAG_UNSPECIFIED if the display name is not known.
func AssetTagFromDisplayName(name string) atagpb.AssetTag {
	lcName := strings.ReplaceAll(strings.ToLower(name), " ", "_")
	allMetadata, err := AllAssetTagMetadata()
	if err != nil {
		return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED
	}
	for _, tagMetadata := range allMetadata {
		if strings.ReplaceAll(strings.ToLower(tagMetadata.GetDisplayName()), " ", "_") == lcName {
			return tagMetadata.GetAssetTag()
		}
	}
	return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED
}

// AssetTagName returns the name of the given AssetTag.
func AssetTagName(tag atagpb.AssetTag) string {
	return tag.String()
}

// AssetTagFromName returns an AssetTag, given its name, or ASSET_TAG_UNSPECIFIED if the name is not
// known.
func AssetTagFromName(name string) atagpb.AssetTag {
	value, ok := atagpb.AssetTag_value[name]
	if !ok {
		return atagpb.AssetTag_ASSET_TAG_UNSPECIFIED
	}

	return atagpb.AssetTag(value)
}

// AssetTagsForTypes returns a list of AssetTags that apply to any of the specified AssetTypes.
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

// AssetTagsForType returns a list of AssetTags that apply to the specified AssetType.
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

// AllAssetTagMetadata returns a list of all AssetTagMetadata.
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

// AssetTagMetadataForTag returns the AssetTagMetadata for the specified AssetTag.
func AssetTagMetadataForTag(tag atagpb.AssetTag) (*atagpb.AssetTagMetadata, error) {
	allMetadata, err := AllAssetTagMetadata()
	if err != nil {
		return nil, err
	}
	for _, tagMetadata := range allMetadata {
		if tagMetadata.GetAssetTag() == tag {
			return tagMetadata, nil
		}
	}
	return nil, fmt.Errorf("no metadata found for tag %v", tag)
}

// AssetTagMetadataForType returns a list of AssetTagMetadata for the specified AssetType.
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
