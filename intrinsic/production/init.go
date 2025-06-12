// Copyright 2023 Intrinsic Innovation LLC

// Package intrinsic provides initialization functionality for Golang binaries.
package intrinsic

import (
	"flag"
)

// Init is the entry point for Golang binaries. It parses command line flags and performs
// other common initialization.
func Init() {
	// Makes LOG(INFO) visible in container logs and disables writing them to a file.
	flag.Set("logtostderr", "true")

	flag.Parse()
	// Other calls can be added here.
}
