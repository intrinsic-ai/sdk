// Copyright 2023 Intrinsic Innovation LLC

package inventoryserver_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	oah "intrinsic/httpjson/openapi/handlers"
	"intrinsic/httpjson/test/inventoryserver"

	pb "intrinsic/httpjson/test/inventory_service_go_proto"
)

const bufSize = 1024 * 1024

func TestHttpEndpoints(t *testing.T) {
	ctx := context.Background()
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	is := inventoryserver.NewInventoryServer()
	pb.RegisterInventoryServiceServer(srv, is)

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	pb.RegisterInventoryServiceHandlerServer(ctx, mux, is)

	// Call AddSku through HTTP/JSON endpoint
	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/skus/c001ca7cafe",
		strings.NewReader(`{"display_name": "Coffee Cat"}`))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Decode and check AddSku response
	if got, want := rr.Code, 200; got != want {
		t.Fatalf("Unexpected HTTP code. got: %v, want: %v", got, want)
	}
	bytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var gotResponse pb.Sku
	if err := protojson.Unmarshal(bytes, &gotResponse); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	wantResponse := &pb.Sku{
		SkuId:       "c001ca7cafe",
		DisplayName: "Coffee Cat",
	}
	if diff := cmp.Diff(wantResponse, gotResponse, protocmp.Transform()); diff != "" {
		t.Errorf("Unexpected response: (-want +got):\n%s", diff)
	}
}

func TestInventoryServer_ListSkus_Empty(t *testing.T) {
	server := inventoryserver.NewInventoryServer()
	req := &pb.ListSkusRequest{}
	resp, err := server.ListSkus(context.Background(), req)
	if err != nil {
		t.Fatalf("ListSkus failed: %v", err)
	}

	if len(resp.GetSkus()) != 0 {
		t.Errorf("Expected no skus, but got %d", len(resp.GetSkus()))
	}
}

func TestOpenAPI(t *testing.T) {
	mux := runtime.NewServeMux()
	mustRegisterOpenAPIHandler(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if got, want := rr.Code, 200; got != want {
		t.Fatalf("Unexpected HTTP code. got: %v, want: %v", got, want)
	}

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	openAPIContent := string(body)
	if got, want := openAPIContent, "/v1/skus"; !strings.Contains(got, want) {
		t.Errorf("OpenAPI spec does not contain expected string. Got %s, want %s", got, want)
	}
}

func mustRegisterOpenAPIHandler(t *testing.T, mux *runtime.ServeMux) {
	rlocationPath := "ai_intrinsic_sdks/intrinsic/httpjson/test/_inventory_service_openapi/openapi.yaml"
	handler, err := oah.MakeOpenAPIHandlerFromRunfiles(rlocationPath)
	if err != nil {
		t.Errorf("Failed to make OpenAPI Handler %v", err)
	}
	mux.HandlePath("GET", "/openapi.yaml", handler)
}
