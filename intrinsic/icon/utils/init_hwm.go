// Copyright 2023 Intrinsic Innovation LLC

// Prepares environment variables and retries for Hardware Module binaries.
package main

import (
	"flag"
	"os"
	"syscall"

	"intrinsic/production/intrinsic"

	log "github.com/golang/glog"
)

const (
	usageHwm = "usage: init_hwm -- path/to/executable [args [...]]"
)

var (
	// We avoid bringing these constants in via cgo because it leads to bad packaging interactions with pkg_tar.
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kRestartRequested
	hardwareModuleRestartRequested = 110
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kFatalFaultDuringInit
	hardwareModuleFatalFaultDuringInit = 111
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kFatalFaultDuringExec
	hardwareModuleFatalFaultDuringExec = 112
)

func main() {
	os.Exit(run())
}

func run() int {
	intrinsic.Init()

	args := flag.Args()
	if len(args) == 0 {
		log.Errorf("Bad invocation: No args given\n%s", usageHwm)
		return 1
	}

	return runMain(RunnerOptions{
		Args: args,
		RestartExitCodes: []int{
			hardwareModuleRestartRequested,
			hardwareModuleFatalFaultDuringInit,
			hardwareModuleFatalFaultDuringExec,
		},
		IgnoredSignals: map[os.Signal]bool{
			syscall.SIGCHLD: true,
		},
	})
}
