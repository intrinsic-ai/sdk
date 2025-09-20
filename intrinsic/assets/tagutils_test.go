// Copyright 2023 Intrinsic Innovation LLC

package tagutils

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"intrinsic/assets/typeutils"

	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
)

func TestAssetTagFromName(t *testing.T) {
	tests := []struct {
		name    string
		wantTag atagpb.AssetTag
		wantErr bool
	}{
		{
			name:    "ASSET_TAG_CAMERA",
			wantTag: atagpb.AssetTag_ASSET_TAG_CAMERA,
		},
		{
			name:    "ASSET_TAG_PAPERCLIP",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTag, err := AssetTagFromName(tc.name)
			if tc.wantErr {
				if err == nil {
					t.Errorf("AssetTagFromName(%v) = %v, want error", tc.name, gotTag)
				}
			} else if err != nil {
				t.Errorf("AssetTagFromName(%v) failed: %v", tc.name, err)
			} else if gotTag != tc.wantTag {
				t.Errorf("AssetTagFromName(%v) = %v, want %v", tc.name, gotTag, tc.wantTag)
			}
		})
	}
}

func TestAssetTagsForTypesReturnsCorrectTags(t *testing.T) {
	allAssetTagMetadata, err := AllAssetTagMetadata()
	if err != nil {
		t.Fatalf("AllAssetTagMetadata() failed: %v", err)
	}
	var allAssetTags []atagpb.AssetTag
	for _, metadata := range allAssetTagMetadata {
		allAssetTags = append(allAssetTags, metadata.GetAssetTag())
	}

	tests := []struct {
		name       string
		assetTypes []atypepb.AssetType
		wantTags   []atagpb.AssetTag
	}{
		{
			name:       "scene object",
			assetTypes: []atypepb.AssetType{atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT},
			wantTags:   []atagpb.AssetTag{atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, atagpb.AssetTag_ASSET_TAG_CAMERA, atagpb.AssetTag_ASSET_TAG_GRIPPER},
		},
		{
			name:       "service and scene object",
			assetTypes: []atypepb.AssetType{atypepb.AssetType_ASSET_TYPE_SERVICE, atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT},
			wantTags:   []atagpb.AssetTag{atagpb.AssetTag_ASSET_TAG_UNSPECIFIED, atagpb.AssetTag_ASSET_TAG_CAMERA, atagpb.AssetTag_ASSET_TAG_GRIPPER},
		},
		{
			name:       "skill",
			assetTypes: []atypepb.AssetType{atypepb.AssetType_ASSET_TYPE_SKILL},
			wantTags:   []atagpb.AssetTag{atagpb.AssetTag_ASSET_TAG_UNSPECIFIED},
		},
		{
			name:       "all asset types",
			assetTypes: typeutils.AllAssetTypes(),
			wantTags:   allAssetTags,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTags, err := AssetTagsForTypes(tc.assetTypes)
			if err != nil {
				t.Fatalf("AssetTagsForTypes(%v) failed: %v", tc.assetTypes, err)
			}

			if diff := cmp.Diff(tc.wantTags, gotTags, cmpopts.SortSlices(func(a, b atagpb.AssetTag) bool { return a < b })); diff != "" {
				t.Errorf("AssetTagsForTypes(%v) returned diff (-want +got):\n%s", tc.assetTypes, diff)
			}
		})
	}
}

func TestAssetTagsForTypeReturnsCorrectTags(t *testing.T) {
	allMetadata, err := AllAssetTagMetadata()
	if err != nil {
		t.Fatalf("AllAssetTagMetadata() failed: %v", err)
	}

	// Compute the expected tags per asset type.
	wantTags := make(map[atypepb.AssetType]map[atagpb.AssetTag]bool)
	for _, typeValue := range atypepb.AssetType_value {
		wantTags[atypepb.AssetType(typeValue)] = make(map[atagpb.AssetTag]bool)
	}
	for _, tagMetadata := range allMetadata {
		for _, assetType := range tagMetadata.GetApplicableAssetTypes() {
			wantTags[assetType][tagMetadata.GetAssetTag()] = true
		}
	}

	for typeName, typeValue := range atypepb.AssetType_value {
		t.Run(typeName, func(t *testing.T) {
			assetType := atypepb.AssetType(typeValue)
			typeAssetTags, err := AssetTagsForType(assetType)
			if err != nil {
				t.Fatalf("AssetTagsForType(%v) failed: %v", atypepb.AssetType(typeValue), err)
			}

			gotTags := make(map[atagpb.AssetTag]bool)
			for _, tag := range typeAssetTags {
				gotTags[tag] = true
			}

			if !cmp.Equal(gotTags, wantTags[assetType]) {
				t.Errorf("AssetTagsForType(%v) = %v, want %v", assetType, gotTags, wantTags[assetType])
			}
		})
	}
}

func TestAssetTagMetadataForTypeReturnsCorrectTags(t *testing.T) {
	allMetadata, err := AllAssetTagMetadata()
	if err != nil {
		t.Fatalf("AllAssetTagMetadata() failed: %v", err)
	}

	// Compute the expected tags per asset type.
	wantTags := make(map[atypepb.AssetType]map[atagpb.AssetTag]bool)
	for _, typeValue := range atypepb.AssetType_value {
		wantTags[atypepb.AssetType(typeValue)] = make(map[atagpb.AssetTag]bool)
	}
	for _, tagMetadata := range allMetadata {
		for _, assetType := range tagMetadata.GetApplicableAssetTypes() {
			wantTags[assetType][tagMetadata.GetAssetTag()] = true
		}
	}

	for typeName, typeValue := range atypepb.AssetType_value {
		t.Run(typeName, func(t *testing.T) {
			assetType := atypepb.AssetType(typeValue)
			metadata, err := AssetTagMetadataForType(assetType)
			if err != nil {
				t.Fatalf("AssetTagMetadataForType(%v) failed: %v", assetType, err)
			}

			gotTags := make(map[atagpb.AssetTag]bool)
			for _, tagMetadata := range metadata {
				gotTags[tagMetadata.GetAssetTag()] = true
			}

			if !cmp.Equal(gotTags, wantTags[assetType]) {
				t.Errorf("AssetTagMetadataForType(%v) = %v, want %v", assetType, gotTags, wantTags[assetType])
			}
		})
	}
}

func TestAllAssetTagMetadataReturnsCorrectTags(t *testing.T) {
	metadata, err := AllAssetTagMetadata()
	if err != nil {
		t.Fatalf("AllAssetTagMetadata() failed: %v", err)
	}

	if len(metadata) != len(atagpb.AssetTag_value) {
		t.Fatalf("AllAssetTagMetadata() returned %d items, want %d", len(metadata), len(atagpb.AssetTag_value))
	}

	for tagName, tagValue := range atagpb.AssetTag_value {
		t.Run(tagName, func(t *testing.T) {
			idx := int(tagValue)
			tagMetadata := metadata[idx]

			wantTag := atagpb.AssetTag(idx)
			if wantTag != tagMetadata.GetAssetTag() {
				t.Errorf("Tag at AllAssetTagMetadata()[%d] == %v, want %v", idx, tagMetadata.GetAssetTag(), wantTag)
			}

			if tagMetadata.GetDisplayName() == "" {
				t.Errorf("Tag at AllAssetTagMetadata()[%d] has empty display name", idx)
			}

			gotTypes := make(map[atypepb.AssetType]bool)
			for _, assetType := range tagMetadata.GetApplicableAssetTypes() {
				if _, ok := gotTypes[assetType]; ok {
					t.Errorf("Applicable asset type %v duplicated for tag %v", assetType, tagMetadata.GetAssetTag())
				}
				gotTypes[assetType] = true
			}
		})
	}
}

func TestUnspecifiedAssetTagAppliesToAllAssetTypes(t *testing.T) {
	metadata, err := AllAssetTagMetadata()
	if err != nil {
		t.Fatalf("AllAssetTagMetadata() failed: %v", err)
	}

	unspecifiedTagTypes := metadata[atagpb.AssetTag_ASSET_TAG_UNSPECIFIED].GetApplicableAssetTypes()
	for _, typeValue := range atypepb.AssetType_value {
		if !slices.Contains(unspecifiedTagTypes, atypepb.AssetType(typeValue)) {
			t.Errorf("unspecified asset tag does not apply to asset type %v", atypepb.AssetType(typeValue))
		}
	}
}
