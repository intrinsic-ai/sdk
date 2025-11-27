// Copyright 2023 Intrinsic Innovation LLC

// Package customer provides access to Flowstate features to create and manage organizations.
package customer

import (
	"fmt"
	"slices"
	"strings"

	"intrinsic/config/environments"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var vipr = viper.New()

// AdminCmd is the `inctl customer` command.
var customerCmd = cobrautil.ParentOfNestedSubcommands("customer", "Manage your Flowstate customers.")

var (
	flagCustomer        string
	flagDebugRequests   bool
	flagName            string
	flagEmail           string
	flagEnvironment     string
	flagInvitationToken string
	flagRoleCSV         string
	flagRole            string
)

func init() {
	customerCmd.PersistentFlags().StringVar(&flagEnvironment, orgutil.KeyEnvironment, environments.Prod, "The environment to use for the command.")
	customerCmd.PersistentFlags().BoolVar(&flagDebugRequests, "debug-requests", false, "If true, print the full request and response for each API call.")
	customerCmd = orgutil.WrapCmd(customerCmd, vipr)
	root.RootCmd.AddCommand(customerCmd)
}

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

func protoPrint(p proto.Message) {
	fmt.Println(p.ProtoReflect().Descriptor().Name())
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
