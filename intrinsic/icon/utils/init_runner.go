// Copyright 2023 Intrinsic Innovation LLC

// Provides a generic runner for ICON and HWM binaries with signal forwarding and restart logic.
package main

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	log "github.com/golang/glog"
)

// RunnerOptions defines the configuration for the runner.
type RunnerOptions struct {
	Args             []string
	RestartExitCodes []int
	IgnoredSignals   map[os.Signal]bool
	ShutdownSignals  []os.Signal
	CommandModifier  func(*exec.Cmd)
}

// runMain executes the command in a loop, handles signals and restarts based on exit codes.
func runMain(opts RunnerOptions) int {
	// We use a large buffer to ensure we don't miss any signals.
	sigs := make(chan os.Signal, 100)
	defer close(sigs)

	// If no specific shutdown signals are provided, use defaults.
	if len(opts.ShutdownSignals) == 0 {
		opts.ShutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP}
	}

	signal.Notify(sigs)
	defer signal.Reset()

	var shutdownRequested atomic.Bool
	var cmdMutex sync.Mutex
	var activeCmd *exec.Cmd

	go func() {
		// Propagate useful signals from the wrapper to the subprocess.
		for sig := range sigs {
			if opts.IgnoredSignals[sig] {
				continue
			}

			isShutdown := false
			for _, s := range opts.ShutdownSignals {
				if s == sig {
					isShutdown = true
					break
				}
			}

			if isShutdown {
				shutdownRequested.Store(true)
			}

			cmdMutex.Lock()
			cmd := activeCmd
			cmdMutex.Unlock()

			if cmd != nil && cmd.Process != nil {
				log.Infof("Forwarding signal %q to pgid %d", sig.String(), cmd.Process.Pid)
				// Negation is intentional: Send the signal to every process in the
				// process group.
				err := syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
				if err != nil {
					log.Infof("Failed forwarding signal %q with error: %v", sig.String(), err)
				}
			}
		}
	}()

	for {
		cmd := exec.Command(opts.Args[0], opts.Args[1:]...)
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		// Set the child's pgid equal to its pid, so that we can send signals to
		// the whole process group.
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if opts.CommandModifier != nil {
			opts.CommandModifier(cmd)
		}

		err := cmd.Start()
		if err == nil {
			cmdMutex.Lock()
			activeCmd = cmd
			cmdMutex.Unlock()

			err = cmd.Wait()

			cmdMutex.Lock()
			activeCmd = nil
			cmdMutex.Unlock()
		}

		if err == nil {
			log.Info("Process finished without error")
			return 0
		}

		if shutdownRequested.Load() {
			log.Info("Shutdown requested. Exiting.")
			return 0
		}

		if cmd.ProcessState == nil {
			log.Errorf("Process failed to start or wait: %v", err)
			return 1
		}

		exitCode := cmd.ProcessState.ExitCode()
		isRestartCode := false
		for _, code := range opts.RestartExitCodes {
			if code == exitCode {
				isRestartCode = true
				break
			}
		}

		if isRestartCode {
			log.Errorf("Process exited with code %d. Restarting.", exitCode)
			continue
		}

		log.Errorf("Process exited with unhandled code %d. Exiting.", exitCode)
		return exitCode
	}
}
