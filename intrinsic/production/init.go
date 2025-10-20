// Copyright 2023 Intrinsic Innovation LLC

// Package intrinsic provides initialization functionality for Golang binaries.
package intrinsic

import (
	"flag"
)

// Init is the entry point for our Golang binaries. It parses command line flags and performs
// other common initialization.
func Init() {
	// When manually running a binary, we want to see the logs on stderr.
	// When running in k8s, we want to see them in container logs and disable writing them to a file.
	flag.Set("logtostderr", "true")

	flag.Parse()
}
