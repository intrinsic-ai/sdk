// Copyright 2023 Intrinsic Innovation LLC

// Package logs defines a command for working with various logs.
package logs

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/prototext"
	"intrinsic/assets/cmdutils"
	"intrinsic/assets/idutils"
	srvpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	sklpb "intrinsic/skills/proto/skill_manifest_go_proto"
	"intrinsic/tools/inctl/cmd/root"
)

const (
	keyFollow        = "follow"
	keyPrefixType    = "prefix_type"
	keyPrefixID      = "prefix_id"
	keySinceSec      = "since"
	keyTailLines     = "tail"
	keyTimestamps    = "timestamps"
	keyTypeService   = "service"
	keyTypeSkill     = "skill"
	keyTypeResource  = "resource"
	keyHiddenDebug   = "debug"
	keyOnpremAddress = "onprem_address"
)

const (
	// The port through which the simulation service can be reached.
	ingressPort = 17080
)

var (
	flagContext      string
	flagUseLocalhost bool
)

var (
	showLogs = &cobra.Command{
		Use:        "logs",
		Aliases:    []string{"slogs"},
		Example:    "inctl logs --org ORGANIZATION --solution SOLUTION-ID --follow --service NAME",
		Short:      "Prints logs from the solution",
		Long:       "Prints resource logs (skill or service) from the instance running in given solution.",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"ID"},
		RunE:       runLogsCmd,
	}

	localViper = viper.New()
	cmdFlags   = cmdutils.NewCmdFlagsWithViper(localViper)
)

func runLogsCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmd.Help()
	}
	target := args[0]

	verboseDebug = cmdFlags.GetBool(keyHiddenDebug)
	verboseOut = cmd.OutOrStderr()

	params := &cmdParams{
		follow:        cmdFlags.GetBool(keyFollow),
		timestamps:    cmdFlags.GetBool(keyTimestamps),
		tailLines:     cmdFlags.GetInt(keyTailLines),
		projectName:   cmdFlags.GetString(cmdutils.KeyProject),
		context:       cmdFlags.GetString(cmdutils.KeyContext),
		solution:      cmdFlags.GetString(cmdutils.KeySolution),
		org:           cmdFlags.GetFlagOrganization(),
		onpremAddress: cmdFlags.GetString(keyOnpremAddress),
	}

	var err error
	if params.resourceType, err = getResourceType(); err != nil {
		return err
	}

	if params.resourceID, err = getResourceID(params.resourceType, target); err != nil {
		return err
	}

	ctx, cancelFx := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer cancelFx()

	err = readLogsFromSolution(ctx, params, cmd.OutOrStdout())
	if err != nil {
		cmd.PrintErrln("Error reading logs. Issue is non-transient and cannot be handled automatically. Please run command again.")
		cmd.PrintErrf("Details: %s\n", err)
		os.Exit(1) // we are doing custom exit, as we are doing custom error handling
	}
	return nil
}

func getResourceID(resType resourceType, target string) (string, error) {
	if strings.HasSuffix(target, ".textproto") {
		file, err := os.Open(target)
		if err != nil {
			return "", fmt.Errorf("cannot open manifest file: %w", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("cannot read manifest file: %w", err)
		}

		switch resType {
		case rtService:
			var manifest srvpb.ServiceManifest
			if err := prototext.Unmarshal(content, &manifest); err != nil {
				return "", fmt.Errorf("cannot parse manifest: %w", err)
			}
			return idutils.IDFrom(manifest.Metadata.Id.Package, manifest.Metadata.Id.Name)
		case rtSkill:
			var manifest sklpb.SkillManifest
			if err := prototext.Unmarshal(content, &manifest); err != nil {
				return "", fmt.Errorf("cannot parse manifest: %w", err)
			}
			return idutils.IDFrom(manifest.Id.Package, manifest.Id.Name)
		default:
			return "", fmt.Errorf("unexpected type %d", resType)
		}
	}
	// We didn't really get a file, so we will treat it as ID
	k8sNormalized := target
	if resType != rtSkill {
		// for the non-skill resources, we need to normalize labels
		k8sNormalized = strings.ReplaceAll(target, "_", "-")
		k8sNormalized = strings.ReplaceAll(k8sNormalized, ".", "-")
	}
	return k8sNormalized, nil
}

func getResourceType() (resourceType, error) {
	if cmdFlags.IsSet(keyTypeSkill) {
		return rtSkill, nil
	}
	if cmdFlags.IsSet(keyTypeService) {
		return rtService, nil
	}
	// todo: make sure resource is mentioned in error internally.
	return -1, fmt.Errorf("resource type for target not set, needs --%s or --%s", keyTypeSkill, keyTypeService)
}

func init() {
	root.RootCmd.AddCommand(showLogs)
	cmdFlags.SetCommand(showLogs)

	// inctl logs --(org|project) --solution [--address] --follow --(service|skill|resource) (manifest|id)

	cmdFlags.AddFlagsProjectOrg()

	cmdFlags.OptionalEnvString(cmdutils.KeySolution, "", "Solution ID from which logs will be read.")
	cmdFlags.OptionalEnvString(cmdutils.KeyContext, "", fmt.Sprintf("The Kubernetes cluster to use or localhost if used with --%s", cmdutils.KeyAddress))
	cmdFlags.MarkHidden(cmdutils.KeyContext)
	cmdFlags.AddFlagAddress()
	cmdFlags.OptionalString(cmdutils.KeyTimeout, "300s", "Maximum time to wait to receive logs.")
	cmdFlags.OptionalBool(keyPrefixType, false, "Prefixes each log line with the type of origin as follows [srv] for service, [skl] for skill and [res] for sesource")
	cmdFlags.OptionalBool(keyPrefixID, false, "Prefixes each log line with the ID of origin in shortened form, e.g.: [ai.int.my_thing]")
	cmdFlags.OptionalBool(keyFollow, false, "Whether to follow the solution logs.")
	cmdFlags.OptionalBool(keyTimestamps, false, "Whether to include timestamps on each log line.")
	cmdFlags.OptionalInt(keyTailLines, 10, "The number of recent log lines to display. An input number less than 0 shows all log lines.")
	cmdFlags.OptionalString(keySinceSec, "", "Show logs starting since value. Value is either relative (e.g 10m) or \ndate time in RFC3339 format (e.g: 2006-01-02T15:04:05Z07:00)")

	cmdFlags.OptionalBool(keyTypeSkill, false, "Indicates logs source is the skill")
	cmdFlags.OptionalBool(keyTypeService, false, "Indicates logs source is the service")

	cmdFlags.OptionalBool(keyHiddenDebug, false, "Prints extensive debug messages")

	// For using the onprem address to fetch logs
	cmdFlags.OptionalString(keyOnpremAddress, "", "The onprem address (host:port) of the workcell. Used to circumvent the need of routing through the cloud, if the workcell is running in the same network as the inctl")

	cmdFlags.MarkHidden(cmdutils.KeyContext, cmdutils.KeyProject, keyTypeResource)
	showLogs.MarkFlagsMutuallyExclusive(keyTypeSkill, keyTypeService)

}
