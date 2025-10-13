// Copyright 2023 Intrinsic Innovation LLC

package operationmode

import (
	"testing"

	opmodepb "intrinsic/config/proto/operation_mode_go_proto"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		mode string
		want opmodepb.OperationMode
	}{
		{
			mode: "real",
			want: opmodepb.OperationMode_REAL_HARDWARE,
		},
		{
			mode: "sim",
			want: opmodepb.OperationMode_SIMULATION,
		},
		{
			mode: "foo",
			want: opmodepb.OperationMode_OPERATION_MODE_UNSPECIFIED,
		},
		{
			mode: "",
			want: opmodepb.OperationMode_OPERATION_MODE_UNSPECIFIED,
		},
	}

	for _, tc := range tests {
		t.Run(tc.mode, func(t *testing.T) {
			got := FromString(tc.mode)
			if got != tc.want {
				t.Errorf("operationmode.FromString(%q) = %v, want %v", tc.mode, got, tc.want)
			}
		})
	}
}
