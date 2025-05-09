// Copyright 2023 Intrinsic Innovation LLC

// Package auth provides authorization client and client side library.
package auth

import (
	"bufio"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/printer"
)

const (
	keyProjectShort = "p"
	keyAlias        = "alias"
	keyBatch        = "batch"
)

// Can be overwridden/injected in tests.
var authStore = auth.NewStore()

func setPrinterFromOutputFlag(command *cobra.Command, args []string) (err error) {
	if out, err := printer.NewPrinter(root.FlagOutput); err == nil {
		command.SetOut(out)
	}
	return
}

var (
	viperLocal = viper.New()
)

var authCmd = orgutil.WrapCmd(&cobra.Command{
	Use:   "auth",
	Short: "Manages user authorization",
	Long:  "Manages user authorization for accessing solutions in the project.",
	// Catching common typos and potential alternatives
	SuggestFor:         []string{"ath", "uath", "auht", "user", "credentials"},
	PersistentPreRunE:  setPrinterFromOutputFlag,
	DisableFlagParsing: true,
}, viperLocal)

var (
	userYesNoPositiveDefOpt = []string{"Y", "n"}
	userYesNoPositiveDefIdx = 0
	userYesNoNegativeDefOpt = []string{"y", "N"}
	userYesNoNegativeDefIdx = 1
)

func userPrompt(rw *bufio.ReadWriter, prompt string, defaultOption int, options ...string) (string, error) {
	if len(options) > 0 {
		prompt += " [" + strings.Join(options, "/") + "]"
	} else {
		defaultOption = -1 // we just mark options as no default just in case here
	}
	prompt += ": "
	if _, err := rw.WriteString(prompt); err != nil {
		return "", err
	}
	rw.Flush() // print out buffer content before we request user input

	response, err := rw.ReadString('\n')
	if err != nil {
		return "", err
	}
	response = strings.TrimSpace(response)
	if response == "" && defaultOption > -1 {
		response = options[defaultOption]
	}
	return response, nil
}

func newReadWriterForCmd(cmd *cobra.Command) *bufio.ReadWriter {
	return bufio.NewReadWriter(
		bufio.NewReader(cmd.InOrStdin()),
		bufio.NewWriter(cmd.OutOrStdout()))
}

func init() {
	root.RootCmd.AddCommand(authCmd)
	authCmd.Flags().MarkHidden(orgutil.KeyProject)
}
