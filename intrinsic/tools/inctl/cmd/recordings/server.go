// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"

	"intrinsic/tools/inctl/auth/auth"

	grpcpb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

func newBagPackagerClient(ctx context.Context) (grpcpb.BagPackagerClient, error) {
	conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(localViper))
	if err != nil {
		return nil, err
	}
	return grpcpb.NewBagPackagerClient(conn), nil
}
