// Copyright 2023 Intrinsic Innovation LLC

// Package defaults provides default values for VM pool commands.
package defaults

const (
	// Tier is the default tier to use when inctl vm pool is used without specifying a tier.
	Tier = "ad_hoc_single"
	// HardwareTemplate is the default hardware template to use when inctl vm pool is used without specifying a hardware template.
	HardwareTemplate = "default"
)
