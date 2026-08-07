// Copyright 2023 Intrinsic Innovation LLC

package fieldbehavior

import (
	"testing"

	"intrinsic/util/proto/testing/prototestutil"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	testpb "intrinsic/util/proto/testing/field_behavior_test_go_proto"
)

func TestClearOutputOnly(t *testing.T) {
	tests := []struct {
		name  string
		given proto.Message
		want  proto.Message
	}{
		{
			name:  "nil_message",
			given: nil,
			want:  nil,
		},
		{
			name: "clean_message_unmodified",
			given: &testpb.FieldBehaviorTestMessage{
				RegularField: "regular",
				AnyField:     prototestutil.MustWrapInAny(t, &testpb.FieldBehaviorTestSubMessage{RegularField: "inside_any"}),
				Nested: &testpb.FieldBehaviorTestSubMessage{
					RegularField: "nested_regular",
				},
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{RegularField: "rep_regular"},
				},
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"key": {RegularField: "map_regular"},
				},
			},
			want: &testpb.FieldBehaviorTestMessage{
				RegularField: "regular",
				AnyField:     prototestutil.MustWrapInAny(t, &testpb.FieldBehaviorTestSubMessage{RegularField: "inside_any"}),
				Nested: &testpb.FieldBehaviorTestSubMessage{
					RegularField: "nested_regular",
				},
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{RegularField: "rep_regular"},
				},
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"key": {RegularField: "map_regular"},
				},
			},
		},
		{
			name: "top_level_output_only_fields",
			given: &testpb.FieldBehaviorTestMessage{
				RegularField:                "regular",
				OutputOnlyField:             proto.String("should_be_cleared"),
				RepeatedOutputOnly:          []string{"a", "b"},
				MapOutputOnly:               map[string]string{"k": "v"},
				OutputOnlyBool:              true,
				OutputOnlyUint32:            proto.Uint32(42),
				NonOptionalOutputOnlyString: "should_also_be_cleared",
			},
			want: &testpb.FieldBehaviorTestMessage{
				RegularField: "regular",
			},
		},
		{
			name: "nested_output_only_fields",
			given: &testpb.FieldBehaviorTestMessage{
				Nested: &testpb.FieldBehaviorTestSubMessage{
					RegularField:    "nested_regular",
					OutputOnlyField: proto.String("nested_clear"),
				},
			},
			want: &testpb.FieldBehaviorTestMessage{
				Nested: &testpb.FieldBehaviorTestSubMessage{
					RegularField: "nested_regular",
				},
			},
		},
		{
			name: "repeated_nested_output_only_fields",
			given: &testpb.FieldBehaviorTestMessage{
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{
						RegularField:    "rep_regular_0",
						OutputOnlyField: proto.String("rep_clear_0"),
					},
					{
						RegularField:    "rep_regular_1",
						OutputOnlyField: proto.String("rep_clear_1"),
					},
				},
			},
			want: &testpb.FieldBehaviorTestMessage{
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{
						RegularField: "rep_regular_0",
					},
					{
						RegularField: "rep_regular_1",
					},
				},
			},
		},
		{
			name: "map_nested_output_only_fields",
			given: &testpb.FieldBehaviorTestMessage{
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"entry": {
						RegularField:    "map_regular",
						OutputOnlyField: proto.String("map_clear"),
					},
				},
			},
			want: &testpb.FieldBehaviorTestMessage{
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"entry": {
						RegularField: "map_regular",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := ClearOutputOnly(tc.given); err != nil {
				t.Fatalf("ClearOutputOnly() returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, tc.given, protocmp.Transform()); diff != "" {
				t.Errorf("ClearOutputOnly() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateNoOutputOnly(t *testing.T) {
	tests := []struct {
		name    string
		msg     proto.Message
		wantErr error
	}{
		{
			name: "valid_clean_message",
			msg: &testpb.FieldBehaviorTestMessage{
				RegularField: "hello",
				Nested: &testpb.FieldBehaviorTestSubMessage{
					RegularField: "world",
				},
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{RegularField: "item1"},
				},
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"k1": {RegularField: "v1"},
				},
			},
		},
		{
			name: "nil_message",
			msg:  nil,
		},
		{
			name: "non_optional_output_only_fields_set_to_default_passes",
			msg: &testpb.FieldBehaviorTestMessage{
				RegularField:                "hello",
				OutputOnlyBool:              false,
				NonOptionalOutputOnlyString: "",
				RepeatedOutputOnly:          []string{},
				MapOutputOnly:               map[string]string{},
			},
			wantErr: nil,
		},
		{
			name: "optional_output_only_string_set_to_empty_default_fails",
			msg: &testpb.FieldBehaviorTestMessage{
				OutputOnlyField: proto.String(""),
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.output_only_field"},
		},
		{
			name: "optional_output_only_uint32_set_to_zero_default_fails",
			msg: &testpb.FieldBehaviorTestMessage{
				OutputOnlyUint32: proto.Uint32(0),
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.output_only_uint32"},
		},
		{
			name: "nested_optional_output_only_string_set_to_empty_default_fails",
			msg: &testpb.FieldBehaviorTestMessage{
				Nested: &testpb.FieldBehaviorTestSubMessage{
					OutputOnlyField: proto.String(""),
				},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.nested.output_only_field"},
		},
		{
			name: "top_level_output_only_string",
			msg: &testpb.FieldBehaviorTestMessage{
				OutputOnlyField: proto.String("set"),
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.output_only_field"},
		},
		{
			name: "top_level_non_optional_output_only_string_non_default_fails",
			msg: &testpb.FieldBehaviorTestMessage{
				NonOptionalOutputOnlyString: "custom",
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.non_optional_output_only_string"},
		},
		{
			name: "top_level_output_only_bool_non_default_fails",
			msg: &testpb.FieldBehaviorTestMessage{
				OutputOnlyBool: true,
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.output_only_bool"},
		},
		{
			name: "top_level_output_only_uint32",
			msg: &testpb.FieldBehaviorTestMessage{
				OutputOnlyUint32: proto.Uint32(10),
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.output_only_uint32"},
		},
		{
			name: "top_level_repeated_output_only",
			msg: &testpb.FieldBehaviorTestMessage{
				RepeatedOutputOnly: []string{"invalid"},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.repeated_output_only"},
		},
		{
			name: "top_level_map_output_only",
			msg: &testpb.FieldBehaviorTestMessage{
				MapOutputOnly: map[string]string{"foo": "bar"},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.map_output_only"},
		},
		{
			name: "nested_output_only",
			msg: &testpb.FieldBehaviorTestMessage{
				Nested: &testpb.FieldBehaviorTestSubMessage{
					OutputOnlyField: proto.String("set"),
				},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.nested.output_only_field"},
		},
		{
			name: "repeated_nested_output_only",
			msg: &testpb.FieldBehaviorTestMessage{
				RepeatedNested: []*testpb.FieldBehaviorTestSubMessage{
					{RegularField: "ok"},
					{OutputOnlyField: proto.String("set")},
				},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.repeated_nested[1].output_only_field"},
		},
		{
			name: "map_nested_output_only",
			msg: &testpb.FieldBehaviorTestMessage{
				MapNested: map[string]*testpb.FieldBehaviorTestSubMessage{
					"my_key": {OutputOnlyField: proto.String("set")},
				},
			},
			wantErr: &OutputOnlyFieldError{FieldPath: "FieldBehaviorTestMessage.map_nested[my_key].output_only_field"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateNoOutputOnly(tc.msg)
			if diff := cmp.Diff(tc.wantErr, err); diff != "" {
				t.Errorf("ValidateNoOutputOnly() diff (-want +got):\n%s", diff)
			}
		})
	}
}
