// Copyright 2023 Intrinsic Innovation LLC

// Prepares environment variables and retries for Hardware Module binaries.
package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"flag"
	log "github.com/golang/glog"
	intrinsic "intrinsic/production/intrinsic"
)

const (
	usage = "usage: init_hwm -- path/to/executable [args [...]]"
)

var (
	// We avoid bringing these constants in via cgo because it leads to bad packaging interactions with pkg_tar.
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kRestartRequested
	hardwareModuleRestartRequested = 110
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kFatalFaultDuringInit
	hardwareModuleFatalFaultDuringInit = 111
	// Corresponds to intrinsic::icon::HardwareModuleExitCode::kFatalFaultDuringExec
	hardwareModuleFatalFaultDuringExec = 112
	shutdownRequested = false
)

func main() {
	flag.Set("alsologtostderr", "true")
	intrinsic.Init()
	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs)
	defer signal.Reset()

	args := flag.Args()
	if len(args) == 0 {
		log.Exitf("Bad invocation: No args given\n%s", usage)
	}

	for {
		cmd := exec.Command(args[0], args[1:]...)

		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		// Set the child's pgid equal to its pid, so that we can send signals to
		// the whole process group (see godoc/pkg/syscall#SysProcAttr).
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		go func() {
			// Propagate useful signals from the wrapper to the subprocess. This allows for
			// e.g. SIGTERM to be properly handled instead of the subprocess being immediately terminated.
			for sig := range sigs {
				if sig != syscall.SIGCHLD && cmd != nil && cmd.Process != nil {
					if sig == syscall.SIGTERM || sig == syscall.SIGINT || sig == syscall.SIGKILL {
						shutdownRequested = true
					}
					// Negation is intentional: Send the signal to every process in the
					// process group, see the manpage for kill(2).
					err := syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))

					if err != nil {
						log.Infof("Failed forwarding signal %q to pgid %d with error: %v", sig.String(), cmd.Process.Pid, err)
					}
				}
			}
		}()
		log.Info("Starting Hardware Module: ", args)
		if err := cmd.Run(); err == nil {
			log.Info("Finished Hardware Module without error")
			return
		}
		// Ignore exit failures after a shutdown request.
		if shutdownRequested {
			log.Info("Shutdown requested. Exiting init wrapper.")
			return
		}

		exitCode := cmd.ProcessState.ExitCode()
		switch exitCode {
		case hardwareModuleRestartRequested:
			log.Error("Hardware Module requested a restart. Restarting.")
			continue
		case hardwareModuleFatalFaultDuringInit:
			log.Error("Hardware Module faulted during initialization. Restarting.")
			continue
		case hardwareModuleFatalFaultDuringExec:
			log.Error("Hardware Module faulted during execution. Restarting.")
			continue
		default:
			log.Error("Hardware Module returned an unhandled status. Exiting. Exit code:", exitCode)
			os.Exit(cmd.ProcessState.ExitCode())
		}
	}
}
