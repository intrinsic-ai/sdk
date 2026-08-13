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

// IsAgent returns whether the command is invoked by an AI agent.
func IsAgent() bool {
	return checkEnvForAgent()
}

// AgentName returns the name of the AI agent invoking the command.
func AgentName() string {
	return getenv("INVOKER_INFO_NAME")
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
