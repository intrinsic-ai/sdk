// Copyright 2023 Intrinsic Innovation LLC

package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// Helper for testing. This function can be run by calling this binary with `-test.run=TestHelperProcess`
func TestHelperProcess(t *testing.T) {
	// If we are being called as the subprocess.
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		helperMain()
		return
	}
}

func helperMain() {
	arg := os.Getenv("HELPER_ARG")
	switch arg {
	case "exit0":
		os.Exit(0)
	case "exit110":
		// hardwareModuleRestartRequested
		os.Exit(110)
	case "wait_signal":
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)
		<-sigs
		os.Exit(0)
	default:
		os.Exit(1)
	}
}

func TestRunSuccess(t *testing.T) {
	opts := RunnerOptions{
		Args: []string{os.Args[0], "-test.run=TestHelperProcess"},
		CommandModifier: func(cmd *exec.Cmd) {
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HELPER_ARG=exit0")
		},
	}
	got := runMain(opts)
	if got != 0 {
		t.Errorf("runMain() = %d; want 0", got)
	}
}

func TestRunRestart(t *testing.T) {
	opts := RunnerOptions{
		Args:             []string{os.Args[0], "-test.run=TestHelperProcess"},
		RestartExitCodes: []int{110},
		CommandModifier: func(cmd *exec.Cmd) {
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HELPER_ARG=exit110")
		},
	}

	// We'll send ourselves a SIGTERM after a short delay to stop the loop.
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	got := runMain(opts)
	if got != 0 {
		t.Errorf("runMain() = %d; want 0 (after shutdown)", got)
	}
}

func TestSignalForwarding(t *testing.T) {
	opts := RunnerOptions{
		Args: []string{os.Args[0], "-test.run=TestHelperProcess"},
		CommandModifier: func(cmd *exec.Cmd) {
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HELPER_ARG=wait_signal")
		},
	}

	// Send SIGTERM to ourselves, which should be forwarded to the subprocess.
	go func() {
		// Wait until the subprocess is likely started.
		time.Sleep(200 * time.Millisecond)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	got := runMain(opts)
	if got != 0 {
		t.Errorf("runMain() = %d; want 0 (after signal forwarding)", got)
	}
}
