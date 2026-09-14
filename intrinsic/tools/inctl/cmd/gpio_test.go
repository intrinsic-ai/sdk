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

package gpio

import (
	"context"
	"errors"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	codespb "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	gpiogrpcpb "intrinsic/hardware/gpio/v1/gpio_service_go_proto"
	pb "intrinsic/hardware/gpio/v1/gpio_service_go_proto"
	signalpb "intrinsic/hardware/gpio/v1/signal_go_proto"
)

// mockGPIOServer implements gpiogrpcpb.GPIOServiceServer
type mockGPIOServer struct {
	gpiogrpcpb.UnimplementedGPIOServiceServer
	lastMetadata metadata.MD

	signalDescriptions []*signalpb.SignalDescription
	signalValues       *signalpb.SignalValueSet

	failInitialSession bool
}

func (s *mockGPIOServer) GetSignalDescriptions(ctx context.Context, req *pb.GetSignalDescriptionsRequest) (*pb.GetSignalDescriptionsResponse, error) {
	s.lastMetadata, _ = metadata.FromIncomingContext(ctx)
	return &pb.GetSignalDescriptionsResponse{SignalDescriptions: s.signalDescriptions}, nil
}

func (s *mockGPIOServer) ReadSignals(ctx context.Context, req *pb.ReadSignalsRequest) (*pb.ReadSignalsResponse, error) {
	s.lastMetadata, _ = metadata.FromIncomingContext(ctx)
	return &pb.ReadSignalsResponse{SignalValues: s.signalValues}, nil
}

func (s *mockGPIOServer) OpenWriteSession(stream gpiogrpcpb.GPIOService_OpenWriteSessionServer) error {
	s.lastMetadata, _ = metadata.FromIncomingContext(stream.Context())
	// Receive initial session data
	initReq, err := stream.Recv()
	if err != nil {
		return err
	}
	if initReq.GetInitialSessionData() == nil {
		return status.Errorf(codespb.InvalidArgument, "expected initial session data")
	}

	if s.failInitialSession {
		return stream.Send(&pb.OpenWriteSessionResponse{
			Status: status.New(codespb.PermissionDenied, "session denied").Proto(),
		})
	}

	// Send OK response
	if err := stream.Send(&pb.OpenWriteSessionResponse{
		Status: status.New(codespb.OK, "").Proto(),
	}); err != nil {
		return err
	}

	// Receive write signals request
	writeReq, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	if writeReq.GetWriteSignals() == nil {
		return status.Errorf(codespb.InvalidArgument, "expected write signals")
	}
	// Send write response
	return stream.Send(&pb.OpenWriteSessionResponse{
		Status: status.New(codespb.OK, "").Proto(),
	})
}

func TestGPIOCommands(t *testing.T) {
	tmpDir := t.TempDir()
	socketPath := filepath.Join(tmpDir, "gpio.sock")

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("failed to listen on %s: %v", socketPath, err)
	}

	server := grpc.NewServer()
	mock := &mockGPIOServer{
		signalDescriptions: []*signalpb.SignalDescription{
			{
				SignalName: "test_signal",
				CanRead:    true,
				CanWrite:   true,
				Type:       signalpb.SignalType_SIGNAL_TYPE_BOOL,
			},
		},
		signalValues: &signalpb.SignalValueSet{
			Values: map[string]*signalpb.SignalValue{
				"test_signal": {
					Value: &signalpb.SignalValue_BoolValue{BoolValue: true},
				},
			},
		},
	}
	gpiogrpcpb.RegisterGPIOServiceServer(server, mock)

	go server.Serve(lis)
	defer server.Stop()

	// Override dialCluster to connect to mock UDS server.
	// This exercises the REAL makeGPIOClient in gpio.go, including all alias resolution!
	var lastDialOpts ClientOptions
	originalDialCluster := dialCluster
	dialCluster = func(ctx context.Context, opts ClientOptions) (context.Context, *grpc.ClientConn, error) {
		lastDialOpts = opts
		conn, err := grpc.DialContext(ctx, "unix://"+socketPath, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, nil, err
		}
		return ctx, conn, nil
	}
	t.Cleanup(func() { dialCluster = originalDialCluster })

	tests := []struct {
		name           string
		run            func(ctx context.Context) error
		wantErr        bool
		verifyMetadata func(t *testing.T, md metadata.MD)
	}{
		{
			name: "list-signals",
			run: func(ctx context.Context) error {
				return listSignals(ctx, ClientOptions{})
			},
			verifyMetadata: func(t *testing.T, md metadata.MD) {
				if values := md.Get("x-resource-instance-name"); len(values) != 0 {
					t.Errorf("expected no x-resource-instance-name header, got %v", values)
				}
			},
		},
		{
			name: "read-signals",
			run: func(ctx context.Context) error {
				return readSignals(ctx, ClientOptions{}, []string{"test_signal"})
			},
		},
		{
			name: "read-signals empty",
			run: func(ctx context.Context) error {
				oldValues := mock.signalValues
				mock.signalValues = &signalpb.SignalValueSet{}
				defer func() { mock.signalValues = oldValues }()
				return readSignals(ctx, ClientOptions{}, []string{"nonexistent"})
			},
		},
		{
			name: "write-signal bool true",
			run: func(ctx context.Context) error {
				return writeSignal(ctx, ClientOptions{}, "test_signal", "true")
			},
		},
		{
			name: "write-signal bool false",
			run: func(ctx context.Context) error {
				return writeSignal(ctx, ClientOptions{}, "test_signal", "false")
			},
		},
		{
			name: "write-signal pbtxt bool",
			run: func(ctx context.Context) error {
				return writeSignal(ctx, ClientOptions{}, "test_signal", "bool_value: true")
			},
		},
		{
			name: "write-signal invalid value",
			run: func(ctx context.Context) error {
				return writeSignal(ctx, ClientOptions{}, "test_signal", "not-a-valid-value")
			},
			wantErr: true,
		},
		{
			name: "write-signal session failure",
			run: func(ctx context.Context) error {
				mock.failInitialSession = true
				defer func() { mock.failInitialSession = false }()
				return writeSignal(ctx, ClientOptions{}, "test_signal", "true")
			},
			wantErr: true,
		},
		{
			name: "list-signals with instance_name",
			run: func(ctx context.Context) error {
				return listSignals(ctx, ClientOptions{InstanceName: "my-gpio"})
			},
			verifyMetadata: func(t *testing.T, md metadata.MD) {
				if values := md.Get("x-resource-instance-name"); len(values) == 0 || values[0] != "my-gpio" {
					t.Errorf("expected x-resource-instance-name header to be 'my-gpio', got %v", values)
				}
			},
		},
		{
			name: "list-signals with address",
			run: func(ctx context.Context) error {
				err := listSignals(ctx, ClientOptions{Address: "explicit-address:1234"})
				if err != nil {
					return err
				}
				if got := lastDialOpts.Address; got != "explicit-address:1234" {
					t.Errorf("expected dialed address to be %q, got %q", "explicit-address:1234", got)
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock.lastMetadata = nil

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := tc.run(ctx)
			if (err != nil) != tc.wantErr {
				t.Fatalf("run() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}

			if mock.lastMetadata == nil {
				t.Errorf("expected gRPC call, but mock server received no metadata")
			}

			if tc.verifyMetadata != nil {
				tc.verifyMetadata(t, mock.lastMetadata)
			}
		})
	}
}

func TestMakeGPIOClientDialError(t *testing.T) {
	originalDialCluster := dialCluster
	dialCluster = func(ctx context.Context, opts ClientOptions) (context.Context, *grpc.ClientConn, error) {
		return nil, nil, errors.New("dial failed")
	}
	defer func() { dialCluster = originalDialCluster }()

	ctx := context.Background()
	_, client, err := makeGPIOClient(ctx, ClientOptions{})
	if err == nil {
		client.Close()
		t.Fatal("expected error when dialCluster fails, got nil")
	}
	if !strings.Contains(err.Error(), "dial failed") {
		t.Errorf("expected error to contain 'dial failed', got %v", err)
	}
}

func TestTryConvertSignalValue(t *testing.T) {
	tests := []struct {
		input       string
		wantErr     bool
		checkOutput func(t *testing.T, val *signalpb.SignalValue)
	}{
		{
			input: "true",
			checkOutput: func(t *testing.T, val *signalpb.SignalValue) {
				if b, ok := val.Value.(*signalpb.SignalValue_BoolValue); !ok || !b.BoolValue {
					t.Errorf("expected true bool, got %v", val)
				}
			},
		},
		{
			input: "TRUE",
			checkOutput: func(t *testing.T, val *signalpb.SignalValue) {
				if b, ok := val.Value.(*signalpb.SignalValue_BoolValue); !ok || !b.BoolValue {
					t.Errorf("expected true bool, got %v", val)
				}
			},
		},
		{
			input: "false",
			checkOutput: func(t *testing.T, val *signalpb.SignalValue) {
				if b, ok := val.Value.(*signalpb.SignalValue_BoolValue); !ok || b.BoolValue {
					t.Errorf("expected false bool, got %v", val)
				}
			},
		},
		{
			input: "bool_value: true",
			checkOutput: func(t *testing.T, val *signalpb.SignalValue) {
				if b, ok := val.Value.(*signalpb.SignalValue_BoolValue); !ok || !b.BoolValue {
					t.Errorf("expected true bool from pbtxt, got %v", val)
				}
			},
		},
		{
			input: "int_value: 42",
			checkOutput: func(t *testing.T, val *signalpb.SignalValue) {
				if i, ok := val.Value.(*signalpb.SignalValue_IntValue); !ok || i.IntValue != 42 {
					t.Errorf("expected int_value: 42, got %v", val)
				}
			},
		},
		{
			input:   "not_a_valid_value",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			val, err := tryConvertSignalValue(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("tryConvertSignalValue(%q) error = %v, wantErr = %v", tc.input, err, tc.wantErr)
			}
			if err == nil && tc.checkOutput != nil {
				tc.checkOutput(t, val)
			}
		})
	}
}

func TestDeprecatedFlagDefinitions(t *testing.T) {
	// Test that server and resource_instance_name flags exist on gpioCmd
	serverFlag := gpioCmd.PersistentFlags().Lookup("server")
	if serverFlag == nil {
		t.Fatal("expected 'server' persistent flag to exist")
	}
	if serverFlag.Deprecated == "" {
		t.Error("expected 'server' persistent flag to be marked deprecated")
	}

	resFlag := gpioCmd.PersistentFlags().Lookup("resource_instance_name")
	if resFlag == nil {
		t.Fatal("expected 'resource_instance_name' persistent flag to exist")
	}
	if resFlag.Deprecated == "" {
		t.Error("expected 'resource_instance_name' persistent flag to be marked deprecated")
	}

	instFlag := gpioCmd.PersistentFlags().Lookup("instance_name")
	if instFlag == nil {
		t.Fatal("expected 'instance_name' persistent flag to exist")
	}

	addrFlag := gpioCmd.PersistentFlags().Lookup("address")
	if addrFlag == nil {
		t.Fatal("expected 'address' persistent flag to exist")
	}
}

func resetGPIOFlags(t *testing.T) {
	t.Helper()
	flagInstanceName = ""
	initGPIOFlags()
}

func TestDeprecatedFlagParsing(t *testing.T) {
	t.Cleanup(func() { resetGPIOFlags(t) })

	tests := []struct {
		name             string
		args             []string
		wantAddress      string
		wantInstanceName string
	}{
		{
			name:        "server flag sets address",
			args:        []string{"--server", "my-server:1234"},
			wantAddress: "my-server:1234",
		},
		{
			name:             "resource_instance_name flag sets instance_name",
			args:             []string{"--resource_instance_name", "my-instance"},
			wantInstanceName: "my-instance",
		},
		{
			name:        "address flag takes precedence over server flag (server first)",
			args:        []string{"--server", "ignored:1234", "--address", "explicit:5678"},
			wantAddress: "explicit:5678",
		},
		{
			name:        "address flag takes precedence over server flag (address first)",
			args:        []string{"--address", "explicit:5678", "--server", "ignored:1234"},
			wantAddress: "explicit:5678",
		},
		{
			name:             "instance_name flag takes precedence over resource_instance_name flag (resource_instance_name first)",
			args:             []string{"--resource_instance_name", "ignored-inst", "--instance_name", "explicit-inst"},
			wantInstanceName: "explicit-inst",
		},
		{
			name:             "instance_name flag takes precedence over resource_instance_name flag (instance_name first)",
			args:             []string{"--instance_name", "explicit-inst", "--resource_instance_name", "ignored-inst"},
			wantInstanceName: "explicit-inst",
		},
		{
			name:        "multiple server flags take the last value",
			args:        []string{"--server", "first:1234", "--server", "second:5678"},
			wantAddress: "second:5678",
		},
		{
			name:             "multiple resource_instance_name flags take the last value",
			args:             []string{"--resource_instance_name", "first-inst", "--resource_instance_name", "second-inst"},
			wantInstanceName: "second-inst",
		},
		{
			name:        "address flag takes precedence over server flag (interleaved)",
			args:        []string{"--server", "first:1234", "--address", "explicit:5678", "--server", "second:9999"},
			wantAddress: "explicit:5678",
		},
		{
			name:             "instance_name flag takes precedence over resource_instance_name flag (interleaved)",
			args:             []string{"--resource_instance_name", "first-inst", "--instance_name", "explicit-inst", "--resource_instance_name", "second-inst"},
			wantInstanceName: "explicit-inst",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetGPIOFlags(t)

			if err := gpioCmd.ParseFlags(tc.args); err != nil {
				t.Fatalf("ParseFlags(%v) error = %v", tc.args, err)
			}

			opts := clientOptionsFromFlags()
			if got := flags.GetString("address"); got != tc.wantAddress {
				t.Errorf("flags.GetString(\"address\") = %q, want %q", got, tc.wantAddress)
			}
			if got := opts.Address; got != tc.wantAddress {
				t.Errorf("clientOptionsFromFlags().Address = %q, want %q", got, tc.wantAddress)
			}
			if got := flagInstanceName; got != tc.wantInstanceName {
				t.Errorf("flagInstanceName = %q, want %q", got, tc.wantInstanceName)
			}
			if got := opts.InstanceName; got != tc.wantInstanceName {
				t.Errorf("clientOptionsFromFlags().InstanceName = %q, want %q", got, tc.wantInstanceName)
			}
		})
	}
}
