// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
