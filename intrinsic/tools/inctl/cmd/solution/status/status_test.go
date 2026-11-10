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

package status

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"intrinsic/testing/grpctest"
	"intrinsic/tools/inctl/util/printer"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	sspb "intrinsic/frontend/solution_service/proto/solution_service_go_proto"
	statuspb "intrinsic/frontend/solution_service/proto/status_go_proto"
)

type fakeSolutionServiceServer struct {
	sspb.UnimplementedSolutionServiceServer
	response *statuspb.Status
	err      error
}

func (s *fakeSolutionServiceServer) GetStatus(
	ctx context.Context,
	req *sspb.GetStatusRequest,
) (*statuspb.Status, error) {
	return s.response, s.err
}

func startSolutionStatusServer(t *testing.T, resp *statuspb.Status, err error) string {
	t.Helper()
	srv := grpc.NewServer()
	sspb.RegisterSolutionServiceServer(srv, &fakeSolutionServiceServer{
		response: resp,
		err:      err,
	})
	return grpctest.StartServerT(t, srv)
}

func TestStatusTextOutput(t *testing.T) {
	wantStatus := &statuspb.Status{
		DisplayName:     "Test Solution",
		Name:            "sol-123",
		State:           statuspb.Status_READY,
		Simulated:       true,
		PlatformVersion: "20260401.RC00",
		ClusterName:     "local-cluster",
	}
	addr := startSolutionStatusServer(t, wantStatus, nil)

	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--address=" + addr})

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext() unexpected error: %v", err)
	}

	got := out.String()
	for _, wantSubstr := range []string{
		"Solution:   Test Solution",
		"ID:         sol-123",
		"State:      READY",
		"Simulation: true",
		"Platform:   20260401.RC00",
		"Cluster:    local-cluster",
	} {
		if !strings.Contains(got, wantSubstr) {
			t.Errorf("output missing %q; got:\n%s", wantSubstr, got)
		}
	}
}

func TestStatusJSONOutput(t *testing.T) {
	wantStatus := &statuspb.Status{
		DisplayName:     "Test Solution",
		Name:            "sol-123",
		State:           statuspb.Status_PLATFORM_READY,
		Simulated:       false,
		PlatformVersion: "20260401.RC00",
		ClusterName:     "local-cluster",
	}
	addr := startSolutionStatusServer(t, wantStatus, nil)

	cmd := NewCommand()
	cmd.Flags().String(printer.KeyOutput, printer.TextOutputFormat, "Output format")
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--address=" + addr, "--output=json"})

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext() unexpected error: %v", err)
	}

	gotStatus := &statuspb.Status{}
	if err := protojson.Unmarshal(bytes.TrimSpace(out.Bytes()), gotStatus); err != nil {
		t.Fatalf("protojson.Unmarshal(%q) failed: %v", out.String(), err)
	}
	if diff := cmp.Diff(wantStatus, gotStatus, protocmp.Transform()); diff != "" {
		t.Errorf("GetStatus JSON mismatch (-want +got):\n%s", diff)
	}
}

func TestStatusError(t *testing.T) {
	addr := startSolutionStatusServer(t, nil, fmt.Errorf("backend unavailable"))

	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--address=" + addr})

	if err := cmd.ExecuteContext(context.Background()); err == nil {
		t.Fatal("ExecuteContext() succeeded, want error")
	}
}
