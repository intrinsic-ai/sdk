// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gpio contains all commands for gpio handling.
package gpio

import (
	"context"
	"fmt"

	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/sshutil"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gpiogrpcpb "intrinsic/hardware/gpio/v1/gpio_service_go_proto"
)

var (
	flagServerAddress        string
	flagK8sContext           string
	flagResourceInstanceName string
)

type connectionManager interface {
	client() gpiogrpcpb.GPIOServiceClient
	close()
}

type direct struct {
	gpioC gpiogrpcpb.GPIOServiceClient
	conn  *grpc.ClientConn
}

func (d *direct) client() gpiogrpcpb.GPIOServiceClient { return d.gpioC }

func (d *direct) close() {
	if err := d.conn.Close(); err != nil {
		fmt.Printf("Failed to close gRPC connection: %v", err)
	}
}

type portForward struct {
	gpioC gpiogrpcpb.GPIOServiceClient
	cm    *sshutil.ConnectionManager
}

func (p *portForward) client() gpiogrpcpb.GPIOServiceClient { return p.gpioC }

func (p *portForward) close() {
	if err := p.cm.Connection.Close(); err != nil {
		fmt.Printf("Failed to close gRPC connection: %v", err)
	}
	p.cm.Close()
}

func makeDirect(serverAddress string) (*direct, error) {
	opts := []grpc.DialOption{
		// suppress go/nogo-check#disallowedfunction google3/third_party/golang/grpc/grpc.WithInsecure
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.NewClient(%q) failed", serverAddress)
	}
	return &direct{
		gpioC: gpiogrpcpb.NewGPIOServiceClient(conn),
		conn:  conn,
	}, nil
}

func makePortForward(ctx context.Context, k8sContext string) (*portForward, error) {
	p := &portForward{cm: new(sshutil.ConnectionManager)}
	if err := p.cm.InitPortForward(ctx, k8sContext, "app-ingress", "istio-ingressgateway", 80); err != nil {
		return nil, errors.Wrap(err, "cannot establish port forward and connect to the ingress")
	}
	p.gpioC = gpiogrpcpb.NewGPIOServiceClient(p.cm.Connection)
	return p, nil
}

func makeConnectionManager(ctx context.Context, serverAddress string, k8sContext string) (connectionManager, error) {
	if serverAddress != "" {
		return makeDirect(serverAddress)
	}
	if k8sContext != "" {
		return makePortForward(ctx, k8sContext)
	}
	return nil, errors.New("either --server or --context must be set")
}

var gpioCmd = cobrautil.ParentOfNestedSubcommands("gpio", "Introspect and operate GPIO")

func init() {
	gpioCmd.PersistentFlags().StringVar(&flagResourceInstanceName, "resource_instance_name", "", "name of the ICON resource instance to connect to")
	gpioCmd.PersistentFlags().StringVar(&flagServerAddress, "server", "", "address and port of the GPIO server to contact")
	gpioCmd.PersistentFlags().StringVar(&flagK8sContext, "context", "", "Kubernetes context to use for port-forwarding")
	gpioCmd.MarkPersistentFlagRequired("resource_instance_name")
	gpioCmd.MarkFlagsMutuallyExclusive("server", "context")
	root.RootCmd.AddCommand(gpioCmd)
}
