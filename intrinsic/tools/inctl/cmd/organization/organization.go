// Copyright 2023 Intrinsic Innovation LLC

// Package organization provides commands for viewing and managing your organizations.
package organization

import (
	"fmt"
	"slices"
	"strings"

	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var vipr = viper.New()

// organizationCmd is the `inctl organization` command.
var organizationCmd = cobrautil.ParentOfNestedSubcommands("organization", "Manage your Flowstate organizations.")

var (
	flagDebugRequests   bool
	flagName            string
	flagEmail           string
	flagInvitationToken string
	flagRoleCSV         string
	flagRole            string
	flagParent          string
	flagShowDeleted     bool
	flagOrgDisplayName  string
	flagYes             bool
)

func init() {
	organizationCmd.Aliases = []string{"organizations"}
	connectInit()
	root.RootCmd.AddCommand(organizationCmd)
}

const (
	orgPrefix = "organizations/"
)

func addPrefix(s string, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

func addPrefixes(s []string, prefix string) []string {
	ps := slices.Clone(s)
	for i := range ps {
		ps[i] = addPrefix(ps[i], prefix)
	}
	return ps
}

// checkOrgNotIntrinsic makes sure the organization is not a reserved name.
func checkOrgNotIntrinsic() error {
	if vipr.GetString(orgutil.KeyOrganization) == "intrinsic" {
		return fmt.Errorf("the current organization cannot be 'intrinsic' for this command")
	}
	return nil
}

// processOrgFlag parses the org flag. Errors if the organization is not given.
func processOrgFlag() (string, error) {
	org := vipr.GetString(orgutil.KeyOrganization)
	if org == "" {
		return "", fmt.Errorf("the --org flag is required for this command")
	}
	if err := checkOrgNotIntrinsic(); err != nil {
		return "", err
	}
	// The `inctl organization` command does not support @project.
	// This is to avoid confusion with the `--env` flag and enforce
	// the notion that `inctl organization` is a global command.
	if strings.Contains(org, "@") {
		return "", fmt.Errorf("`--org=<org>@<project>` syntax is not supported and not required by `inctl organization`")
	}
	return org, nil
}

// resolveOrgArgOrFlag resolves the organization from an optional positional argument or falls back to --org / Viper config.
func resolveOrgArgOrFlag(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		org := args[0]
		if strings.Contains(org, "@") {
			return "", fmt.Errorf("`<org>@<project>` syntax is not supported and not required by `inctl organization`")
		}
		return org, nil
	}
	return processOrgFlag()
}

func protoPrint(p proto.Message) {
	ms, err := protojson.MarshalOptions{
		Multiline:         true,
		UseProtoNames:     true,
		EmitUnpopulated:   true,
		EmitDefaultValues: true,
	}.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(ms))
}
