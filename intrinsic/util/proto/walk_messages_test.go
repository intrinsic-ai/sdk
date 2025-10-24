// Copyright 2023 Intrinsic Innovation LLC

package walkmessages

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	pb "intrinsic/util/proto/testing/test_message_go_proto"
)

func TestRecursively(t *testing.T) {
	testCases := []struct {
		name        string
		inMsg       proto.Message
		processFunc fProcessMessage
		wantMsg     proto.Message
		wantErr     bool
		numCalls    int
	}{
		{
			name: "traversal_and_modification",
			inMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
				RepeatedNested: []*pb.NestedMessage{
					{Value: "repeated1"},
					{Value: "repeated2"},
				},
				MapNested: map[string]*pb.NestedMessage{
					"key1": {Value: "map1"},
					"key2": {Value: "map2"},
				},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				switch msg := m.(type) {
				case *pb.TestMessage:
					newMsg := proto.Clone(msg).(*pb.TestMessage)
					newMsg.Name = "modified_root"
					return newMsg, true, nil
				case *pb.NestedMessage:
					newMsg := proto.Clone(msg).(*pb.NestedMessage)
					newMsg.Value = "modified_" + newMsg.Value
					return newMsg, true, nil
				}
				return nil, true, nil
			},
			wantMsg: &pb.TestMessage{
				Name:   "modified_root",
				Nested: &pb.NestedMessage{Value: "modified_nested"},
				RepeatedNested: []*pb.NestedMessage{
					{Value: "modified_repeated1"},
					{Value: "modified_repeated2"},
				},
				MapNested: map[string]*pb.NestedMessage{
					"key1": {Value: "modified_map1"},
					"key2": {Value: "modified_map2"},
				},
			},
			numCalls: 6,
		},
		{
			name: "no_recursion",
			inMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				return nil, false, nil
			},
			wantMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
			},
			numCalls: 1,
		},
		{
			name: "error_propagation",
			inMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				if _, ok := m.(*pb.NestedMessage); ok {
					return nil, false, errors.New("test error")
				}
				return nil, true, nil
			},
			wantErr: true,
		},
		{
			name:  "nil_input",
			inMsg: nil,
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				if m == nil {
					return nil, false, nil
				}
				return nil, false, fmt.Errorf("function called with non-nil message: %v", m)
			},
			wantMsg: nil,
		},
		{
			name: "oneof_nested_message",
			inMsg: &pb.TestMessage{
				OneofField: &pb.TestMessage_OneofNested{
					OneofNested: &pb.NestedMessage{Value: "oneof_nested"},
				},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				if msg, ok := m.(*pb.NestedMessage); ok {
					newMsg := proto.Clone(msg).(*pb.NestedMessage)
					newMsg.Value = "modified_" + newMsg.Value
					return newMsg, true, nil
				}
				return nil, true, nil
			},
			wantMsg: &pb.TestMessage{
				OneofField: &pb.TestMessage_OneofNested{
					OneofNested: &pb.NestedMessage{Value: "modified_oneof_nested"},
				},
			},
			numCalls: 2,
		},
		{
			name: "oneof_primitive_type",
			inMsg: &pb.TestMessage{
				OneofField: &pb.TestMessage_OneofString{
					OneofString: "oneof_string",
				},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				return nil, true, nil
			},
			wantMsg: &pb.TestMessage{
				OneofField: &pb.TestMessage_OneofString{
					OneofString: "oneof_string",
				},
			},
			numCalls: 1,
		},
		{
			name: "skip_subtree",
			inMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				if _, ok := m.(*pb.TestMessage); ok {
					return nil, false, nil // Don't recurse into TestMessage
				}
				return nil, true, errors.New("should not have recursed into subtree")
			},
			wantMsg: &pb.TestMessage{
				Name:   "root",
				Nested: &pb.NestedMessage{Value: "nested"},
			},
			numCalls: 1,
		},
		{
			name: "map_with_non_message_value",
			inMsg: &pb.TestMessage{
				MapToNonMessage: map[int32]float32{
					1: 1.0,
				},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				return nil, true, nil
			},
			wantMsg: &pb.TestMessage{
				MapToNonMessage: map[int32]float32{
					1: 1.0,
				},
			},
			numCalls: 1,
		},
		{
			name: "repeated_with_non_message_value",
			inMsg: &pb.TestMessage{
				RepeatedNonMessage: []float32{1.0, 2.0},
			},
			processFunc: func(m proto.Message) (proto.Message, bool, error) {
				return nil, true, nil
			},
			wantMsg: &pb.TestMessage{
				RepeatedNonMessage: []float32{1.0, 2.0},
			},
			numCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var numCalls int
			countingFunc := func(m proto.Message) (proto.Message, bool, error) {
				numCalls++
				return tc.processFunc(m)
			}

			gotMsg, err := Recursively(tc.inMsg, countingFunc)

			if (err != nil) != tc.wantErr {
				t.Errorf("Recursively() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if diff := cmp.Diff(tc.wantMsg, gotMsg, protocmp.Transform()); diff != "" {
				t.Errorf("Recursively() returned diff (-want +got):\n%s", diff)
			}

			if tc.numCalls > 0 && numCalls != tc.numCalls {
				t.Errorf("processFunc was called %d times, want %d", numCalls, tc.numCalls)
			}
		})
	}
}
