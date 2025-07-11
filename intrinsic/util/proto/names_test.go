// Copyright 2023 Intrinsic Innovation LLC

package names

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	diamondapb "intrinsic/util/proto/testing/diamond_a_go_proto"
)

func TestValidateProtoName(t *testing.T) {
	tests := []struct {
		name      string
		protoName string
		wantError bool
	}{
		{
			name:      "package with subpackage",
			protoName: "intrinsic_proto.my_package.MyMessage",
			wantError: false,
		},
		{
			name:      "no subpackage",
			protoName: "intrinsic_proto.MyMessage",
			wantError: false,
		},
		{
			name:      "package has leading underscore",
			protoName: "intrinsic_proto._my_package.MyMessage",
			wantError: false,
		},
		{
			name:      "package has trailing underscore",
			protoName: "intrinsic_proto.my_package_.MyMessage",
			wantError: false,
		},
		{
			name:      "message has leading underscore",
			protoName: "intrinsic_proto.my_package._MyMessage",
			wantError: false,
		},
		{
			name:      "message has trailing underscore",
			protoName: "intrinsic_proto.my_package.MyMessage_",
			wantError: false,
		},
		{
			name:      "empty",
			protoName: "",
			wantError: true,
		},
		{
			name:      "only message",
			protoName: "MyMessage",
			wantError: true,
		},
		{
			name:      "package starts with number",
			protoName: "1intrinsic_proto.MyMessage",
			wantError: true,
		},
		{
			name:      "message starts with number",
			protoName: "intrinsic_proto.1MyMessage",
			wantError: true,
		},
		{
			name:      "package contains invalid character",
			protoName: "intrinsic_proto.My!Package.MyMessage",
			wantError: true,
		},
		{
			name:      "message contains invalid character",
			protoName: "intrinsic_proto.My!Message",
			wantError: true,
		},
		{
			name:      "package starts with dot",
			protoName: ".intrinsic_proto.MyMessage",
			wantError: true,
		},
		{
			name:      "message starts with dot",
			protoName: "intrinsic_proto.MyMessage.",
			wantError: true,
		},
		{
			name:      "no message",
			protoName: "intrinsic_proto.",
			wantError: true,
		},
		{
			name:      "message ends with dot",
			protoName: "intrinsic_proto.MyMessage.",
			wantError: true,
		},
		{
			name:      "service prefix",
			protoName: "/intrinsic_proto.MyService/",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProtoName(tc.protoName)
			if tc.wantError != (err != nil) {
				t.Errorf("ValidateProtoName(%q) returned error %v, want error %v", tc.protoName, err, tc.wantError)
			}
			if tc.wantError && !errors.Is(err, ErrInvalidProtoName) {
				t.Errorf("ValidateProtoName(%q) returned error %v, want error %v", tc.protoName, err, ErrInvalidProtoName)
			}
		})
	}
}

func TestValidateProtoPrefix(t *testing.T) {
	tests := []struct {
		name        string
		protoPrefix string
		wantError   bool
	}{
		{
			name:        "package with subpackage",
			protoPrefix: "/intrinsic_proto.my_package.MyService/",
			wantError:   false,
		},
		{
			name:        "no subpackage",
			protoPrefix: "/intrinsic_proto.MyService/",
			wantError:   false,
		},
		{
			name:        "package has leading underscore",
			protoPrefix: "/intrinsic_proto._my_package.MyService/",
			wantError:   false,
		},
		{
			name:        "package has trailing underscore",
			protoPrefix: "/intrinsic_proto.my_package_.MyService/",
			wantError:   false,
		},
		{
			name:        "service has leading underscore",
			protoPrefix: "/intrinsic_proto.my_package._MyService/",
			wantError:   false,
		},
		{
			name:        "service has trailing underscore",
			protoPrefix: "/intrinsic_proto.my_package.MyService_/",
			wantError:   false,
		},
		{
			name:        "empty",
			protoPrefix: "",
			wantError:   true,
		},
		{
			name:        "empty with slashes",
			protoPrefix: "//",
			wantError:   true,
		},
		{
			name:        "only service",
			protoPrefix: "/MyService/",
			wantError:   true,
		},
		{
			name:        "package starts with number",
			protoPrefix: "/1intrinsic_proto.MyService/",
			wantError:   true,
		},
		{
			name:        "service starts with number",
			protoPrefix: "/intrinsic_proto.1MyService/",
			wantError:   true,
		},
		{
			name:        "package contains invalid character",
			protoPrefix: "/intrinsic_proto.My!Package.MyService/",
			wantError:   true,
		},
		{
			name:        "service contains invalid character",
			protoPrefix: "/intrinsic_proto.My!Service/",
			wantError:   true,
		},
		{
			name:        "package starts with dot",
			protoPrefix: "/.intrinsic_proto.MyService/",
			wantError:   true,
		},
		{
			name:        "service starts with dot",
			protoPrefix: "/intrinsic_proto.MyService./",
			wantError:   true,
		},
		{
			name:        "no service",
			protoPrefix: "/intrinsic_proto./",
			wantError:   true,
		},
		{
			name:        "service ends with dot",
			protoPrefix: "/intrinsic_proto.MyService./",
			wantError:   true,
		},
		{
			name:        "missing slashes",
			protoPrefix: "intrinsic_proto.MyService",
			wantError:   true,
		},
		{
			name:        "too many slashes",
			protoPrefix: "//intrinsic_proto.MyService//",
			wantError:   true,
		},
		{
			name:        "missing leading slash",
			protoPrefix: "intrinsic_proto.MyService/",
			wantError:   true,
		},
		{
			name:        "missing trailing slash",
			protoPrefix: "/intrinsic_proto.MyService",
			wantError:   true,
		},
		{
			name:        "too many leading slashs",
			protoPrefix: "//intrinsic_proto.MyService/",
			wantError:   true,
		},
		{
			name:        "too many trailing slashs",
			protoPrefix: "/intrinsic_proto.MyService//",
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProtoPrefix(tc.protoPrefix)
			if tc.wantError != (err != nil) {
				t.Errorf("ValidateProtoPrefix(%q) returned error %v, want error %v", tc.protoPrefix, err, tc.wantError)
			}
			if tc.wantError && !errors.Is(err, ErrInvalidProtoPrefix) {
				t.Errorf("ValidateProtoPrefix(%q) returned error %v, want error %v", tc.protoPrefix, err, ErrInvalidProtoPrefix)
			}
		})
	}
}

func TestAnyToProtoName(t *testing.T) {
	tests := []struct {
		name      string
		msg       proto.Message
		wantName  string
		wantError bool
	}{
		{
			name:     "empty",
			msg:      &emptypb.Empty{},
			wantName: "google.protobuf.Empty",
		},
		{
			name:     "diamond a",
			msg:      &diamondapb.A{Value: "foo"},
			wantName: "intrinsic_proto.test.A",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msgAny, err := anypb.New(tc.msg)
			if err != nil {
				t.Fatalf("anypb.New(%v) returned error %v, want nil", tc.msg, err)
			}
			gotName, err := AnyToProtoName(msgAny)
			if tc.wantError {
				if err == nil {
					t.Errorf("AnyToProtoName(%v) returned nil error, want non-nil", msgAny)
				}
			} else if err != nil {
				t.Errorf("AnyToProtoName(%v) returned error %v, want nil", msgAny, err)
			} else if gotName != tc.wantName {
				t.Errorf("AnyToProtoName(%v) returned name %q, want %q", msgAny, gotName, tc.wantName)
			} else if err := ValidateProtoName(gotName); err != nil {
				t.Errorf("AnyToProtoName(%v) returned invalid name %q: %v", msgAny, gotName, err)
			}
		})
	}
}
