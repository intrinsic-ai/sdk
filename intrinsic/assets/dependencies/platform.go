// Copyright 2023 Intrinsic Innovation LLC

// Package platform contains utilities that list the interfaces an Asset provides to the
// platform. These interfaces can be used to determine whether an Asset is compatible with
// a given platform version.
package platform

import (
	"intrinsic/assets/idutils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	idpb "intrinsic/assets/proto/id_go_proto"
)

const (
	// RuntimeAssetID is the Asset ID of the "Platform as Asset".
	RuntimeAssetID = "ai.intrinsic.runtime"
	// RuntimeInstanceName is the Asset instance name of the "Platform as Asset".
	RuntimeInstanceName = "intrinsic_runtime"
)

// ValidateIDNotReserved validates that the given Asset ID is not reserved.
func ValidateIDNotReserved(id *idpb.Id) error {
	idStr, err := idutils.IDFromProto(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	// Verify that the Asset ID is not reserved.
	if idStr == RuntimeAssetID {
		return status.Errorf(codes.InvalidArgument, "invalid id %q is reserved", RuntimeAssetID)
	}
	return nil
}

// ValidateInstanceNameNotReserved validates that the given Asset instance name is not reserved.
func ValidateInstanceNameNotReserved(instanceName string) error {
	// Verify that the Asset instance name is not reserved.
	if instanceName == RuntimeInstanceName {
		return status.Errorf(codes.InvalidArgument, "invalid instance name %q is reserved", RuntimeInstanceName)
	}
	return nil
}
