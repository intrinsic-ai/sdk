// Copyright 2023 Intrinsic Innovation LLC

// Package typeutils provides utilities for asset types.
package typeutils

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	atypepb "intrinsic/assets/proto/asset_type_go_proto"
)

var (
	customAssetTypeToName = map[atypepb.AssetType]string{
		atypepb.AssetType_ASSET_TYPE_UNSPECIFIED: "asset",
	}
	regexEnumName = regexp.MustCompile("ASSET_TYPE_(?P<asset_type>[A-Za-z0-9_]+)$")
	regexEnumNameGroups       = regexEnumName.SubexpNames()
	regexEnumNameAssetTypeIdx = slices.Index(regexEnumNameGroups, "asset_type")
)

// AllAssetTypes returns all asset types.
func AllAssetTypes() []atypepb.AssetType {
	var assetTypes []atypepb.AssetType
	for _, i := range atypepb.AssetType_value {
		assetType := atypepb.AssetType(i)
		if assetType == atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
			continue
		}
		assetTypes = append(assetTypes, assetType)
	}
	return assetTypes
}

// AssetTypesWithInstances returns the asset types that have instances.
func AssetTypesWithInstances() []atypepb.AssetType {
	return []atypepb.AssetType{
		atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE,
		atypepb.AssetType_ASSET_TYPE_SERVICE,
		atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT,
	}
}

// UniqueAssetTypes returns the unique asset types from a list.
//
// Order is arbitrary, and the unspecified asset type is ignored.
func UniqueAssetTypes(assetTypes []atypepb.AssetType) []atypepb.AssetType {
	assetTypesMap := make(map[atypepb.AssetType]bool)
	for _, assetType := range assetTypes {
		if assetType != atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
			assetTypesMap[assetType] = true
		}
	}
	uniqueAssetTypes := make([]atypepb.AssetType, len(assetTypesMap))
	i := 0
	for assetType := range assetTypesMap {
		uniqueAssetTypes[i] = assetType
		i++
	}
	return uniqueAssetTypes
}

// NameFromAssetType returns the name of an asset type.
func NameFromAssetType(a atypepb.AssetType) string {
	if name, ok := customAssetTypeToName[a]; ok {
		return name
	}
	if submatches := regexEnumName.FindStringSubmatch(a.String()); submatches != nil {
		submatch := submatches[regexEnumNameAssetTypeIdx]
		return strings.ReplaceAll(strings.ToLower(submatch), "_", " ")
	}
	return NameFromAssetType(atypepb.AssetType_ASSET_TYPE_UNSPECIFIED)
}

// AssetTypeFromName returns the asset type enum from a string.
func AssetTypeFromName(name string) (atypepb.AssetType, error) {
	customNameToType := map[string]atypepb.AssetType{}
	for k, v := range customAssetTypeToName {
		customNameToType[v] = k
	}
	if assetType, ok := customNameToType[name]; ok {
		return assetType, nil
	}

	if name == "" {
		return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED, fmt.Errorf("asset type name is empty")
	}

	return EnumFromString(fmt.Sprintf("ASSET_TYPE_%s", strings.ToUpper(strings.ReplaceAll(name, " ", "_"))))
}

// IntToAssetType converts a raw int to an AssetType enum.
// It returns an error if the integer does not match any enum value.
func IntToAssetType(i int32) (atypepb.AssetType, error) {
	if _, ok := atypepb.AssetType_name[i]; !ok {
		return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED, fmt.Errorf("asset type enum int %d is invalid", i)
	}
	return atypepb.AssetType(i), nil
}

// EnumFromString returns an AssetType enum from a string.
func EnumFromString(t string) (atypepb.AssetType, error) {
	if i, exists := atypepb.AssetType_value[t]; exists {
		return atypepb.AssetType(i), nil
	}
	return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED, fmt.Errorf("unknown asset type: %q", t)
}
