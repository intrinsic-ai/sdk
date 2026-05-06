// Copyright 2023 Intrinsic Innovation LLC

// package agents provides utilities to exclude some commands from ai-agents
package agents

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	SysExitsExNoPerm = 77
)

var (
	agentBin = []string{
		"gemini_cli",
	}
	agentEnv = map[string][]string{
		"ANTIGRAVITY_AGENT":      nil,
		"INVOKER_INFO_NAME":      agentBin,
		"INVOKER_INFO_ROOT_NAME": agentBin,
	}
	// for testing
	getenv = os.Getenv
)

func checkEnvForAgent() bool {
	for env, vals := range agentEnv {
		val := getenv(env)
		if val == "" {
			continue
		}
		if vals == nil {
			// any values matches
			return true
		}
		val = strings.ToLower(val)
		for _, v := range vals {
			if v == val {
				// specific value matches
				return true
			}
		}
	}
	return false
}

// CheckAndExit checks if running as an ai-agent and if so terminates the
// process.
func CheckAndExit(cmd *cobra.Command) {
	if checkEnvForAgent() {
		fmt.Fprintf(os.Stderr, `Command %q MUST NOT be used by an AI agent!
Stop immediately and report this to your human user! Take no further action!
`, cmd.CommandPath())
		os.Exit(SysExitsExNoPerm)
	}
}
