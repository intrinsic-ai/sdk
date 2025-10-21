// Copyright 2023 Intrinsic Innovation LLC

// Package intrinsic provides initialization functionality for Golang binaries.
package intrinsic

import (
	"flag"
)

// Init is the entry point for our Golang binaries. It parses command line flags and performs
// other common initialization.
func Init() {
	// Avoids logging command-line args and version information, which we get via
	// other means in our stack.
	flag.Set("silent_init", "true")
	// When manually running a binary, we want to see the logs on stderr.
	// When running in k8s, we want to see them in container logs and disable writing them to a file.
	flag.Set("logtostderr", "true")

	flag.Parse()
}
