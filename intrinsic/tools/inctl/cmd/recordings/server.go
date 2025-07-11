// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"
	"crypto/tls"

	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/auth/auth"

	grpcpb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

func newConn(ctx context.Context) (*grpc.ClientConn, error) {
	project := cmdFlags.GetString(cmdutils.KeyProject)
	addr := "www.endpoints." + project + ".cloud.goog:443"

	cfg, err := auth.NewStore().GetConfiguration(project)
	if err != nil {
		return nil, err
	}
	creds, err := cfg.GetDefaultCredentials()
	if err != nil {
		return nil, err
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(creds),
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}
	conn, err := grpc.NewClient(addr, grpcOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.Dial(%q)", addr)
	}
	return conn, nil
}

func newBagPackagerClient(ctx context.Context) (grpcpb.BagPackagerClient, error) {
	conn, err := newConn(ctx)
	if err != nil {
		return nil, err
	}
	return grpcpb.NewBagPackagerClient(conn), nil
}
