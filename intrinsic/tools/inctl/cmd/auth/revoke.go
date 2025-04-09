// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/util/orgutil"
)

const (
	keyRevokeAll = "all"
)

var (
	revokeCmdFlags = cmdutils.NewCmdFlagsWithViper(viperLocal)
)

var revokeCmd = &cobra.Command{
	Use:     "revoke",
	Aliases: []string{"ls"},
	Short:   "Removes local credentials",
	Long:    "Remove selected local credentials. Credentials are currently not revoked on server.",
	Args:    cobra.NoArgs,
	RunE:    revokeCredentialsE,
}

func revokeCredentialsE(cmd *cobra.Command, _ []string) error {
	isRevokeAll := revokeCmdFlags.GetBool(keyRevokeAll)
	credName, isOrg := getConfigurationName()
	if !isRevokeAll && credName == "" {
		return fmt.Errorf("either --%s or --%s needs to be specified", orgutil.KeyOrganization, keyRevokeAll)
	}

	isBatch := revokeCmdFlags.GetBool(keyBatch)

	rw := newReadWriterForCmd(cmd)
	if credName == "" && isRevokeAll {
		if !isBatch {
			resp, err := userPrompt(rw, "Are you sure you want to remove all authorizations?", 1, "yes", "NO")
			if err != nil {
				// this error means something terrible happened with terminal, aborting is really only option
				return fmt.Errorf("cannot continue: %w", err)
			}
			if resp != "yes" {
				return fmt.Errorf("aborted by user")
			}
		}
		return authStore.RemoveAllKnownCredentials()
	}
	if !isBatch {
		prompt := fmt.Sprintf("Are you sure you want to revoke all credentials for '%s'", credName)
		resp, err := userPrompt(rw, prompt, 1, "yes", "NO")
		if err != nil {
			return err
		} else if resp != "yes" {
			return fmt.Errorf("aborted by user")
		}
	}
	if isOrg {
		return authStore.RemoveOrganization(credName)
	}
	return authStore.RemoveConfiguration(credName)
}

func getConfigurationName() (name string, isOrg bool) {
	if orgName := revokeCmdFlags.GetFlagOrganization(); orgName != "" {
		return orgName, true
	}
	if projectName := revokeCmdFlags.GetFlagProject(); projectName != "" {
		return projectName, false
	}
	return "", false
}

func init() {
	authCmd.AddCommand(revokeCmd)

	revokeCmdFlags.SetCommand(revokeCmd)
	revokeCmdFlags.OptionalBool(keyRevokeAll, false, fmt.Sprintf("Revokes all existing credentials. If --%s is omitted, removes all known credentials", orgutil.KeyOrganization))
	revokeCmdFlags.OptionalBool(keyBatch, false, "Suppresses command prompts and assume Yes or default as an answer. Use with shell scripts.")
}
