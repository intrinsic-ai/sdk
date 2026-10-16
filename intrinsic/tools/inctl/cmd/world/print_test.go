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

package world

import (
	"bytes"
	"context"
	"net"
	"strings"
	"sync"
	"testing"

	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"

	owspb "intrinsic/world/proto/object_world_service_go_proto"
)

// fakeObjectWorldService serves a fixed set of worlds and records every
// GetWorld request it receives.
type fakeObjectWorldService struct {
	owspb.UnimplementedObjectWorldServiceServer

	worlds map[string]*owspb.World

	mu       sync.Mutex
	requests []*owspb.GetWorldRequest
}

func (s *fakeObjectWorldService) GetWorld(_ context.Context, req *owspb.GetWorldRequest) (*owspb.World, error) {
	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

	w, ok := s.worlds[req.GetWorldId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "world %q not found", req.GetWorldId())
	}
	return w, nil
}

func (s *fakeObjectWorldService) gotRequests() []*owspb.GetWorldRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.requests
}

// startFakeServer starts an in-process gRPC server backed by srv and points
// dialCluster at it for the duration of the test.
func startFakeServer(t *testing.T, srv owspb.ObjectWorldServiceServer) {
	t.Helper()

	lis := bufconn.Listen(1 << 20)
	s := grpc.NewServer()
	owspb.RegisterObjectWorldServiceServer(s, srv)
	// Serve returns once s.Stop is called during cleanup.
	go s.Serve(lis)
	t.Cleanup(s.Stop)

	origDial := dialCluster
	dialCluster = func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
		// the address is a dummy and never resolved
		conn, err := grpc.NewClient("passthrough:///bufnet",
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return lis.DialContext(ctx)
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		return ctx, conn, err
	}
	t.Cleanup(func() { dialCluster = origDial })
}

// runPrint executes `inctl world print` with the given arguments and returns
// its stdout.
//
// The command runs standalone rather than through root.RootCmd so that the
// organization pre-run of the world command, which resolves local
// credentials, is skipped. It reuses the RunE and the flag definitions of the
// real command, including the global --output flag.
func runPrint(t *testing.T, args ...string) (string, error) {
	t.Helper()

	// Flag values are bound to package variables and persist across executions
	// in the same process, so they are reset to their defaults first.
	flagWorldID = defaultWorldID
	root.FlagOutput = printer.TextOutputFormat

	cmd := &cobra.Command{
		Use:           printCmd.Use,
		Args:          printCmd.Args,
		RunE:          printCmd.RunE,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().AddFlagSet(printCmd.Flags())
	cmd.Flags().AddFlag(root.RootCmd.PersistentFlags().Lookup(printer.KeyOutput))

	var out bytes.Buffer
	cmd.SetOut(&out)
	// A nil slice makes cobra fall back to os.Args, so a non-nil one is passed.
	cmd.SetArgs(append([]string{}, args...))
	err := cmd.ExecuteContext(t.Context())
	return out.String(), err
}

func testWorld(id string) *owspb.World {
	return &owspb.World{
		WorldMetadata: &owspb.WorldMetadata{Id: id},
		Objects: []*owspb.Object{
			{WorldId: id, Id: "root", Name: "root"},
		},
	}
}

func TestPrintRequestsFullView(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantWorldID string
	}{
		{
			name:        "default_world_id",
			args:        []string{},
			wantWorldID: defaultWorldID,
		},
		{
			name:        "explicit_world_id",
			args:        []string{"--world_id=init_world"},
			wantWorldID: "init_world",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeObjectWorldService{
				worlds: map[string]*owspb.World{tc.wantWorldID: testWorld(tc.wantWorldID)},
			}
			startFakeServer(t, fake)

			if _, err := runPrint(t, tc.args...); err != nil {
				t.Fatalf("inctl world print %v returned error: %v", tc.args, err)
			}

			want := []*owspb.GetWorldRequest{{
				WorldId: tc.wantWorldID,
				View:    owspb.WorldView_WORLD_VIEW_FULL,
			}}
			if diff := cmp.Diff(want, fake.gotRequests(), protocmp.Transform()); diff != "" {
				t.Errorf("inctl world print %v sent unexpected GetWorld requests (-want +got):\n%s", tc.args, diff)
			}
		})
	}
}

func TestPrintJSONOutputIsStable(t *testing.T) {
	startFakeServer(t, &fakeObjectWorldService{
		worlds: map[string]*owspb.World{defaultWorldID: testWorld(defaultWorldID)},
	})

	const want = `{"worldMetadata":{"id":"world"},"objects":[{"worldId":"world","id":"root","name":"root"}]}` + "\n"
	for _, args := range [][]string{
		{"--output=json"},
		{"-o", "json"},
	} {
		got, err := runPrint(t, args...)
		if err != nil {
			t.Fatalf("inctl world print %v returned error: %v", args, err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("inctl world print %v printed unexpected output (-want +got):\n%s", args, diff)
		}
	}
}

func TestPrintWorldNotFound(t *testing.T) {
	startFakeServer(t, &fakeObjectWorldService{})

	out, err := runPrint(t, "--world_id=missing")
	if got := status.Code(err); got != codes.NotFound {
		t.Errorf("inctl world print --world_id=missing returned code %v, want %v (error: %v)", got, codes.NotFound, err)
	}
	if err != nil && !strings.Contains(err.Error(), `"missing"`) {
		t.Errorf("inctl world print --world_id=missing error = %q, want it to mention the world id", err)
	}
	if out != "" {
		t.Errorf("inctl world print --world_id=missing printed %q, want no output", out)
	}
}
