// Copyright 2023 Intrinsic Innovation LLC

// Package typeutils provides utilities for AssetTypes.
package typeutils

import (
	"fmt"
	"strings"

	atypepb "intrinsic/assets/proto/asset_type_go_proto"
)

type assetTypeInfo struct {
	CodeName          string
	DisplayName       string
	DisplayNamePlural string
	HasInstances      bool
}

var (
	allAssetTypeInfo = map[atypepb.AssetType]assetTypeInfo{
		atypepb.AssetType_ASSET_TYPE_UNSPECIFIED: assetTypeInfo{
			CodeName:          "unspecified",
			DisplayName:       "Asset",
			DisplayNamePlural: "Assets",
			HasInstances:      false,
		},
		atypepb.AssetType_ASSET_TYPE_SCENE_OBJECT: assetTypeInfo{
			CodeName:          "scene_object",
			DisplayName:       "SceneObject",
			DisplayNamePlural: "SceneObjects",
			HasInstances:      true,
		},
		atypepb.AssetType_ASSET_TYPE_SERVICE: assetTypeInfo{
			CodeName:          "service",
			DisplayName:       "Service",
			DisplayNamePlural: "Services",
			HasInstances:      true,
		},
		atypepb.AssetType_ASSET_TYPE_SKILL: assetTypeInfo{
			CodeName:          "skill",
			DisplayName:       "Skill",
			DisplayNamePlural: "Skills",
			HasInstances:      false,
		},
		atypepb.AssetType_ASSET_TYPE_HARDWARE_DEVICE: assetTypeInfo{
			CodeName:          "hardware_device",
			DisplayName:       "HardwareDevice",
			DisplayNamePlural: "HardwareDevices",
			HasInstances:      true,
		},
		atypepb.AssetType_ASSET_TYPE_DATA: assetTypeInfo{
			CodeName:          "data",
			DisplayName:       "Data",
			DisplayNamePlural: "Data",
			HasInstances:      false,
		},
		atypepb.AssetType_ASSET_TYPE_PROCESS: assetTypeInfo{
			CodeName:          "process",
			DisplayName:       "Process",
			DisplayNamePlural: "Processes",
			HasInstances:      false,
		},
	}
)

// AllAssetTypes returns all AssetTypes except ASSET_TYPE_UNSPECIFIED.
func AllAssetTypes() []atypepb.AssetType {
	var assetTypes []atypepb.AssetType
	for assetType := range allAssetTypeInfo {
		if assetType != atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
			assetTypes = append(assetTypes, assetType)
		}
	}
	return assetTypes
}

// AssetTypesWithInstances returns the AssetTypes that have instances.
func AssetTypesWithInstances() []atypepb.AssetType {
	var assetTypes []atypepb.AssetType
	for assetType := range allAssetTypeInfo {
		if allAssetTypeInfo[assetType].HasInstances {
			assetTypes = append(assetTypes, assetType)
		}
	}
	return assetTypes
}

// AssetTypeCodeName returns the code name of an AssetType.
func AssetTypeCodeName(a atypepb.AssetType) string {
	if info, ok := allAssetTypeInfo[a]; ok {
		return info.CodeName
	}
	return allAssetTypeInfo[atypepb.AssetType_ASSET_TYPE_UNSPECIFIED].CodeName
}

// AssetTypeFromCodeName returns the AssetType for the given code name, or
// ASSET_TYPE_UNSPECIFIED if the code name is not known.
func AssetTypeFromCodeName(name string) atypepb.AssetType {
	for assetType, info := range allAssetTypeInfo {
		if info.CodeName == name {
			return assetType
		}
	}
	return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
}

// AssetTypeDisplayName returns the display name of an AssetType.
func AssetTypeDisplayName(a atypepb.AssetType) string {
	if info, ok := allAssetTypeInfo[a]; ok {
		return info.DisplayName
	}
	return allAssetTypeInfo[atypepb.AssetType_ASSET_TYPE_UNSPECIFIED].DisplayName
}

// AssetTypeFromDisplayName returns the AssetType for the given display name, or
// ASSET_TYPE_UNSPECIFIED if the display name is not known.
func AssetTypeFromDisplayName(name string) atypepb.AssetType {
	for assetType, info := range allAssetTypeInfo {
		if info.DisplayName == name {
			return assetType
		}
	}
	return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
}

// AssetTypeDisplayNamePlural returns the plural display name of an AssetType.
func AssetTypeDisplayNamePlural(a atypepb.AssetType) string {
	if info, ok := allAssetTypeInfo[a]; ok {
		return info.DisplayNamePlural
	}
	return allAssetTypeInfo[atypepb.AssetType_ASSET_TYPE_UNSPECIFIED].DisplayNamePlural
}

// AssetTypeDisplayNamesPlural returns the plural display names for a list of AssetTypes.
func AssetTypeDisplayNamesPlural(types []atypepb.AssetType) string {
	displayNames := make([]string, 0, len(types))
	foundDisplayNames := map[string]struct{}{}
	for _, at := range types {
		displayName := AssetTypeDisplayNamePlural(at)
		if displayName == "Assets" {
			return "Assets"
		}
		if _, ok := foundDisplayNames[displayName]; !ok {
			foundDisplayNames[displayName] = struct{}{}
			displayNames = append(displayNames, displayName)
		}
	}
	switch len(displayNames) {
	case 0:
		return "Assets"
	case 1:
		return displayNames[0]
	case 2:
		return fmt.Sprintf("%s and %s", displayNames[0], displayNames[1])
	default:
		displayNames, last := displayNames[:len(displayNames)-1], displayNames[len(displayNames)-1]
		return fmt.Sprintf("%s, and %s", strings.Join(displayNames, ", "), last)
	}
}

// AssetTypeName returns the name of the given AssetType.
func AssetTypeName(t atypepb.AssetType) string {
	return t.String()
}

// AssetTypeFromName returns an AssetType, given its name, or ASSET_TYPE_UNSPECIFIED if the name is
// not known.
func AssetTypeFromName(t string) atypepb.AssetType {
	if i, ok := atypepb.AssetType_value[t]; ok {
		return atypepb.AssetType(i)
	}
	return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
}

// AssetTypeFromInt returns an AssetType enum from an integer.
func AssetTypeFromInt(i int32) atypepb.AssetType {
	if _, ok := atypepb.AssetType_name[i]; ok {
		return atypepb.AssetType(i)
	}
	return atypepb.AssetType_ASSET_TYPE_UNSPECIFIED
}

// UniqueAssetTypes returns the unique AssetTypes from a list.
//
// Order is arbitrary, and the unspecified AssetType is ignored.
func UniqueAssetTypes(assetTypes []atypepb.AssetType) []atypepb.AssetType {
	assetTypesMap := map[atypepb.AssetType]struct{}{}
	for _, assetType := range assetTypes {
		if assetType != atypepb.AssetType_ASSET_TYPE_UNSPECIFIED {
			assetTypesMap[assetType] = struct{}{}
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
