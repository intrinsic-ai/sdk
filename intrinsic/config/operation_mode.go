// Copyright 2023 Intrinsic Innovation LLC

// Package operationmode provides features related to operation mode.
package operationmode

import opmodepb "intrinsic/config/proto/operation_mode_go_proto"

const (
	// Real is the command line string to start in real mode.
	Real = "real"
	// Sim is the command line string to start in sim mode.
	Sim = "sim"
)

// FromString transforms an operation mode string into a proto enum value.
func FromString(mode string) opmodepb.OperationMode {
	switch mode {
	case Real:
		return opmodepb.OperationMode_REAL_HARDWARE
	case Sim:
		return opmodepb.OperationMode_SIMULATION
	default:
		return opmodepb.OperationMode_OPERATION_MODE_UNSPECIFIED
	}
}
