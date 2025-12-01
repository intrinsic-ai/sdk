// Copyright 2023 Intrinsic Innovation LLC

package tagutils

import (
	"slices"
	"testing"

	"intrinsic/assets/typeutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	atagpb "intrinsic/assets/proto/asset_tag_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
)

func TestAllAssetTags(t *testing.T) {
	allTags := AllAssetTags()
	if len(allTags) != len(atagpb.AssetTag_value)-1 {
		t.Errorf("AllAssetTags() returned %d items, want %d", len(allTags), len(atagpb.AssetTag_value)-1)
	}
	for _, tag := range allTags {
		if tag == atagpb.AssetTag_ASSET_TAG_UNSPECIFIED {
			t.Errorf("AllAssetTags() returned ASSET_TAG_UNSPECIFIED")
		}
	}
}

func TestAllAssetTagsReturnsTagsInEnumOrder(t *testing.T) {
	allTags := AllAssetTags()
	for i, tag := range allTags {
		if tag != atagpb.AssetTag(i+1) {
			t.Errorf("AllAssetTags()[%d] == %v, want %v", i, tag, atagpb.AssetTag(i+1))
		}
	}
}

func TestAssetTagDisplayName(t *testing.T) {
	tests := []struct {
		name            string
		tag             atagpb.AssetTag
		wantDisplayName string
	}{
		{
			name:            "camera",
			tag:             atagpb.AssetTag_ASSET_TAG_CAMERA,
			wantDisplayName: "Camera",
		},
		{
			name:            "gripper",
			tag:             atagpb.AssetTag_ASSET_TAG_GRIPPER,
			wantDisplayName: "Gripper",
		},
		{
			name:            "subprocess",
			tag:             atagpb.AssetTag_ASSET_TAG_SUBPROCESS,
			wantDisplayName: "Subprocess",
		},
		{
			name:            "arm",
			tag:             atagpb.AssetTag_ASSET_TAG_ARM,
			wantDisplayName: "Arm",
		},
		{
			name:            "unspecified",
			tag:             atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
			wantDisplayName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotDisplayName := AssetTagDisplayName(tc.tag)
			if gotDisplayName != tc.wantDisplayName {
				t.Errorf("AssetTagDisplayName(%v) = %v, want %v", tc.tag, gotDisplayName, tc.wantDisplayName)
			}
		})
	}
}

func TestAssetTagFromDisplayName(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		wantTag     atagpb.AssetTag
	}{
		{
			name:        "camera",
			displayName: "Camera",
			wantTag:     atagpb.AssetTag_ASSET_TAG_CAMERA,
		},
		{
			name:        "lower case camera",
			displayName: "camera",
			wantTag:     atagpb.AssetTag_ASSET_TAG_CAMERA,
		},
		{
			name:        "gripper",
			displayName: "Gripper",
			wantTag:     atagpb.AssetTag_ASSET_TAG_GRIPPER,
		},
		{
			name:        "subprocess",
			displayName: "Subprocess",
			wantTag:     atagpb.AssetTag_ASSET_TAG_SUBPROCESS,
		},
		{
			name:        "arm",
			displayName: "Arm",
			wantTag:     atagpb.AssetTag_ASSET_TAG_ARM,
		},
		{
			name:        "unspecified",
			displayName: "Unspecified",
			wantTag:     atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
		},
		{
			name:        "paperclip",
			displayName: "Paperclip",
			wantTag:     atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTag := AssetTagFromDisplayName(tc.displayName)
			if gotTag != tc.wantTag {
				t.Errorf("AssetTagFromDisplayName(%v) = %v, want %v", tc.displayName, gotTag, tc.wantTag)
			}
		})
	}
}

func TestAssetTagName(t *testing.T) {
	tests := []struct {
		name string
		tag  atagpb.AssetTag
		want string
	}{
		{
			name: "camera",
			tag:  atagpb.AssetTag_ASSET_TAG_CAMERA,
			want: "ASSET_TAG_CAMERA",
		},
		{
			name: "gripper",
			tag:  atagpb.AssetTag_ASSET_TAG_GRIPPER,
			want: "ASSET_TAG_GRIPPER",
		},
		{
			name: "subprocess",
			tag:  atagpb.AssetTag_ASSET_TAG_SUBPROCESS,
			want: "ASSET_TAG_SUBPROCESS",
		},
		{
			name: "arm",
			tag:  atagpb.AssetTag_ASSET_TAG_ARM,
			want: "ASSET_TAG_ARM",
		},
		{
			name: "unspecified",
			tag:  atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
			want: "ASSET_TAG_UNSPECIFIED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AssetTagName(tc.tag)
			if got != tc.want {
				t.Errorf("AssetTagName(%v) = %v, want %v", tc.tag, got, tc.want)
			}
		})
	}
}

func TestAssetTagFromName(t *testing.T) {
	tests := []struct {
		name    string
		wantTag atagpb.AssetTag
	}{
		{
			name:    "ASSET_TAG_CAMERA",
			wantTag: atagpb.AssetTag_ASSET_TAG_CAMERA,
		},
		{
			name:    "ASSET_TAG_UNSPECIFIED",
			wantTag: atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
		},
		{
			name:    "ASSET_TAG_PAPERCLIP",
			wantTag: atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTag := AssetTagFromName(tc.name)
			if gotTag != tc.wantTag {
				t.Errorf("AssetTagFromName(%v) = %v, want %v", tc.name, gotTag, tc.wantTag)
			}
		})
	}
}

func TestAssetTagsForTypes(t *testing.T) {
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
			wantTags: []atagpb.AssetTag{
				atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
				atagpb.AssetTag_ASSET_TAG_CAMERA,
				atagpb.AssetTag_ASSET_TAG_GRIPPER,
				atagpb.AssetTag_ASSET_TAG_ARM,
			},
		},
		{
			name:       "service and scene object",
			assetTypes: []atypepb.AssetType{atypepb.AssetType_ASSET_TYPE_SERVICE, atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT},
			wantTags: []atagpb.AssetTag{
				atagpb.AssetTag_ASSET_TAG_UNSPECIFIED,
				atagpb.AssetTag_ASSET_TAG_CAMERA,
				atagpb.AssetTag_ASSET_TAG_GRIPPER,
				atagpb.AssetTag_ASSET_TAG_ARM,
			},
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

func TestAssetTagsForType(t *testing.T) {
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

func TestAssetTagMetadataForTag(t *testing.T) {
	for tagName, tagValue := range atagpb.AssetTag_value {
		t.Run(tagName, func(t *testing.T) {
			tag := atagpb.AssetTag(tagValue)
			metadata, err := AssetTagMetadataForTag(tag)
			if err != nil {
				t.Fatalf("AssetTagMetadataForTag(%v) failed: %v", tag, err)
			}
			if metadata.GetAssetTag() != tag {
				t.Errorf("AssetTagMetadataForTag(%v) = %v, want %v", tag, metadata.GetAssetTag(), tag)
			}
		})
	}
}

func TestAssetTagMetadataForType(t *testing.T) {
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

func TestAllAssetTagMetadata(t *testing.T) {
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

			if tagMetadata.GetAssetTag() != atagpb.AssetTag_ASSET_TAG_UNSPECIFIED && tagMetadata.GetDisplayName() == "" {
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
