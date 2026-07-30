// Copyright 2023 Intrinsic Innovation LLC

package serialization_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"intrinsic/httpjson/serialization"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestRequestFormat(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        serialization.Format
	}{
		{"DefaultJSON", "", serialization.FormatJSON},
		{"ExplicitJSON", "application/json", serialization.FormatJSON},
		{"ProtobufJSON", "application/protobuf+json", serialization.FormatJSON},
		{"ProtobufBinary", "application/protobuf", serialization.FormatProto},
		{"XProtobufBinary", "application/x-protobuf", serialization.FormatProto},
		{"OctetStream", "application/octet-stream", serialization.FormatProto},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/test", nil)
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			if got := serialization.RequestFormat(req); got != tc.want {
				t.Errorf("RequestFormat() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResponseFormat(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		accept      string
		want        serialization.Format
	}{
		{"DefaultJSON_NoHeaders", "", "", serialization.FormatJSON},
		{"ExplicitAcceptProto", "", "application/protobuf", serialization.FormatProto},
		{"ExplicitAcceptJSON", "", "application/json", serialization.FormatJSON},
		{"FallbackToContentTypeWhenAcceptWildcard", "application/protobuf", "*/*", serialization.FormatProto},
		{"FallbackToContentTypeWhenAcceptEmpty", "application/x-protobuf", "", serialization.FormatProto},
		{"AcceptOverridesContentType", "application/json", "application/protobuf", serialization.FormatProto},
		{"AcceptJSONOverridesProtoContentType", "application/protobuf", "application/json", serialization.FormatJSON},
		{"MultiValuedAcceptProtoFirst", "application/json", "application/protobuf;q=1.0, application/json;q=0.5", serialization.FormatProto},
		{"MultiValuedAcceptServerPreference", "application/json", "application/json, application/protobuf", serialization.FormatProto},
		{"MultiValuedAcceptWithWildcardFallback", "application/x-protobuf", "text/html, */*;q=0.8", serialization.FormatProto},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/test", nil)
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			if tc.accept != "" {
				req.Header.Set("Accept", tc.accept)
			}
			if got := serialization.ResponseFormat(req); got != tc.want {
				t.Errorf("ResponseFormat() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestAnyWithoutDescriptors proves that binary protobuf serialization allows forwarding
// google.protobuf.Any messages without requiring file descriptor sets or type resolvers,
// whereas JSON serialization fails without a type resolver.
func TestAnyWithoutDescriptors(t *testing.T) {
	unknownAny := &anypb.Any{
		TypeUrl: "type.intrinsic.ai/intrinsic.unknown.v1.UnregisteredMessage",
		Value:   []byte{0x0a, 0x04, 0x74, 0x65, 0x73, 0x74}, // random proto field bytes
	}

	tests := []struct {
		name      string
		marshaler runtime.Marshaler
		wantError bool
	}{
		{
			name:      "BinaryProto_SucceedsWithoutDescriptors",
			marshaler: &runtime.ProtoMarshaller{},
			wantError: false,
		},
		{
			name:      "JSON_FailsWithoutDescriptors",
			marshaler: &runtime.JSONPb{}, // No resolvers configured
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.marshaler.Marshal(unknownAny)
			if (err != nil) != tc.wantError {
				t.Fatalf("Marshal() error = %v, wantError %v", err, tc.wantError)
			}
			if tc.wantError {
				return
			}

			decoded := &anypb.Any{}
			if err := tc.marshaler.Unmarshal(data, decoded); err != nil {
				t.Fatalf("Unmarshal() unexpectedly failed: %v", err)
			}

			if !proto.Equal(unknownAny, decoded) {
				t.Errorf("Decoded Any mismatch.\nGot: %v\nWant: %v", decoded, unknownAny)
			}
		})
	}
}

// TestGatewayMux_ContentNegotiation proves that when runtime.ServeMux is configured with our
// multi-format marshaler options, HTTP handlers seamlessly negotiate binary protobuf requests and responses
// containing unknown Any messages without needing descriptor sets.
func TestGatewayMux_ContentNegotiation(t *testing.T) {
	mux := runtime.NewServeMux(serialization.NewGatewayMarshalerOptions(nil)...)

	// Register a simple custom HTTP handler on the mux that echoes back any pb.Any payload.
	// In standard generated gRPC-gateway stubs, runtime.ServeMux automatically selects the marshaler
	// based on the request Content-Type and Accept headers via runtime.MarshalerForRequest.
	err := mux.HandlePath("POST", "/test/any", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		marshaler, _ := runtime.MarshalerForRequest(mux, r)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqMsg := &anypb.Any{}
		if err := marshaler.Unmarshal(body, reqMsg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Echo the request message back using the negotiated marshaler
		respBytes, err := marshaler.Marshal(reqMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", marshaler.ContentType(reqMsg))
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	})
	if err != nil {
		t.Fatalf("Failed to register path on mux: %v", err)
	}

	server := httptest.NewServer(mux)
	defer server.Close()

	unknownAny := &anypb.Any{
		TypeUrl: "type.intrinsic.ai/intrinsic.unknown.v1.ZeroDescriptorMessage",
		Value:   []byte{0x12, 0x34, 0x56, 0x78, 0x90},
	}
	binaryData, _ := proto.Marshal(unknownAny)

	t.Run("HTTP_BinaryProtobuf_EchoesSuccessfully", func(t *testing.T) {
		req, _ := http.NewRequest("POST", server.URL+"/test/any", bytes.NewReader(binaryData))
		req.Header.Set("Content-Type", "application/x-protobuf")
		req.Header.Set("Accept", "application/x-protobuf")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("HTTP POST failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Expected 200 OK, got %d: %s", resp.StatusCode, string(body))
		}

		if gotCT := resp.Header.Get("Content-Type"); !strings.Contains(gotCT, "protobuf") {
			t.Errorf("Expected protobuf Content-Type header, got %q", gotCT)
		}

		respBody, _ := io.ReadAll(resp.Body)
		gotMsg := &anypb.Any{}
		if err := proto.Unmarshal(respBody, gotMsg); err != nil {
			t.Fatalf("Failed to decode response proto: %v", err)
		}

		if !proto.Equal(unknownAny, gotMsg) {
			t.Errorf("Echoed proto mismatch.\nGot: %v\nWant: %v", gotMsg, unknownAny)
		}
	})
}
