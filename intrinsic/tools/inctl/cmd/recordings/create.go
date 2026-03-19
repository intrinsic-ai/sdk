// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/color"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/tools/inctl/util/promptutil"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bagmetadatapb "intrinsic/logging/proto/bag_metadata_go_proto"
	bagpackagerpb "intrinsic/logging/proto/bag_packager_service_go_proto"
	loggerpb "intrinsic/logging/proto/logger_service_go_proto"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	flagDescription  string
	flagQuiet        bool
	flagSkipGenerate bool

	// Data inclusion flags
	flagIncludeAllData         bool
	flagTextLogs               bool
	flagFlowstateData          bool
	flagSceneData              bool
	flagRobotData              bool
	flagPerceptionData         bool
	flagDebugData              bool
	flagAdditionalEventSources []string
)
var eventSourceWorkflowMap = map[string][]string{
	"include_text_logs": {
		"/text-log-out",
		"/asset-text-log-out",
	},
	"include_scene_data": {
		"/assets/.*/markers",
		"/tf",
		"trace_cube_.*",
	},
	"include_robot_data": {
		"motion_planner_service.PlanTrajectory.debug_data",
		"/icon/.*/robot_status",
	},
	"include_perception_data": {
		"perception.*",
	},
	"include_flowstate_data": {
		"/flowstate_events",
		"executive.*",
	},
	"include_debug_data": {
		"/system/metrics",
		"error_report",
	},
}

// Limits following the create_recording skill.
const (
	// mixedDataLimit represents the maximum allowed duration for recordings containing
	// any data other than text logs or flowstate data.
	mixedDataLimit = 10 * time.Minute

	// textOrFlowstateDataLimit represents the maximum allowed duration for recordings
	// containing only text logs or flowstate data.
	textOrFlowstateDataLimit = 24 * time.Hour
)

var createParams = viper.New()

type CreateCmdRunner struct {
	NewDataLoggerClient  func(cmd *cobra.Command) (loggerpb.DataLoggerClient, io.Closer, error)
	NewBagPackagerClient func(cmd *cobra.Command) (bagpackagerpb.BagPackagerClient, error)
	PromptYesNo          func(cmd *cobra.Command, prompt string, defaultBehavior promptutil.DefaultBehavior, invalidInputBehavior promptutil.InvalidInputBehavior) (bool, error)
	RunGenerate          func(cmd *cobra.Command, createResponseBagID string) error
}

// parseAndValidateRecordingTimestamps parses the start and end timestamp flags.
//
// If neither is provided, it defaults to the last 5 minutes.
// If only one is provided, or if the format is invalid, it returns an error.
func parseAndValidateRecordingTimestamps(cmd *cobra.Command) (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	var err error

	if flagStartTimestamp != "" && flagEndTimestamp != "" {
		startTime, err = time.Parse(time.RFC3339, flagStartTimestamp)
		if err != nil {
			return startTime, endTime, fmt.Errorf("invalid start timestamp: %s", flagStartTimestamp)
		}
		endTime, err = time.Parse(time.RFC3339, flagEndTimestamp)
		if err != nil {
			return startTime, endTime, fmt.Errorf("invalid end timestamp: %s", flagEndTimestamp)
		}
	} else if flagStartTimestamp == "" && flagEndTimestamp == "" {
		color.C.Yellow().Fprintf(cmd.ErrOrStderr(), "Warning: No --start_timestamp or --end_timestamp provided, defaulting to the last 5 minutes.\n")
		endTime = time.Now()
		startTime = endTime.Add(-5 * time.Minute)
	} else {
		return startTime, endTime, fmt.Errorf("must supply BOTH start_timestamp and end_timestamp, or NEITHER to default to the last 5 minutes")
	}
	return startTime, endTime, nil
}

// validateRecordingDuration checks if the requested recording duration exceeds limits.
//
// Mixed data recordings are limited to 10 minutes, while text/flowstate-only recordings can be up to 24 hours.
func validateRecordingDuration(duration time.Duration) error {
	isOnlyTextLogsOrFlowstate := (flagTextLogs || flagFlowstateData) &&
		!flagSceneData && !flagRobotData && !flagPerceptionData && !flagDebugData && len(flagAdditionalEventSources) == 0

	if !isOnlyTextLogsOrFlowstate && duration > mixedDataLimit {
		return fmt.Errorf("recording duration of %v exceeds the 10-minute limit for mixed data. Please specify a shorter duration or reduce the requested data", duration)
	} else if isOnlyTextLogsOrFlowstate && duration > textOrFlowstateDataLimit {
		return fmt.Errorf("recording duration of %v exceeds the 24-hour limit for text/flowstate data. Please specify a shorter duration", duration)
	}
	return nil
}

// determineEventSourcesToRecord inspects the provided flags to build a deduplicated list
//
// of all regular expressions matching the event sources to record. If no specific sources
// are requested, it prompts the user to confirm recording all data.
func determineEventSourcesToRecord(cmd *cobra.Command, runner *CreateCmdRunner, flagNames []string) ([]string, []string, error) {
	includedEventSourcesFromFlags := make(map[string]bool)
	includeFlagsProvided := false

	for _, flagName := range append(flagNames, "additional_event_sources") {
		if cmd.Flags().Changed(flagName) {
			includeFlagsProvided = true
			break
		}
	}

	if !includeFlagsProvided {
		if flagIncludeAllData {
			flagAdditionalEventSources = append(flagAdditionalEventSources, ".*")
		} else {
			confirm := false
			var err error
			if !flagQuiet {
				confirm, err = runner.PromptYesNo(cmd, "No data requested to be included. Do you want to create a recording with ALL data (.*) instead?", promptutil.DefaultNo, promptutil.ReemitPromptOnInvalidInput)
				if err != nil {
					return nil, nil, err
				}
			} else {
				// Quiet mode defaults to DefaultNo.
				confirm = false
			}
			if !confirm {
				if flagQuiet {
					return nil, nil, fmt.Errorf("no data was requested to be included")
				}
				return nil, nil, fmt.Errorf("aborted")
			}
			flagAdditionalEventSources = append(flagAdditionalEventSources, ".*")
		}
	}

	var includedFlags []string
	for _, flagName := range flagNames {
		b, err := cmd.Flags().GetBool(flagName)
		if err == nil && b {
			includedFlags = append(includedFlags, strings.TrimPrefix(flagName, "include_"))
			for _, v := range eventSourceWorkflowMap[flagName] {
				includedEventSourcesFromFlags[v] = true
			}
		}
	}
	for _, source := range flagAdditionalEventSources {
		includedEventSourcesFromFlags[source] = true
	}

	var finalEventSources []string
	for source := range includedEventSourcesFromFlags {
		finalEventSources = append(finalEventSources, source)
	}

	if len(finalEventSources) == 0 {
		return nil, nil, fmt.Errorf("no event sources mapped to record. Provide valid flags")
	}

	return finalEventSources, includedFlags, nil
}

// generateRecordingDescription creates a human-readable description for the recording
// if the user did not explicitly provide one. It summarizes the requested data sources.
func generateRecordingDescription(includedFlags []string) string {
	desc := flagDescription
	if desc == "" {
		displayFlags := make([]string, len(includedFlags))
		copy(displayFlags, includedFlags)
		if len(flagAdditionalEventSources) > 0 {
			allData := false
			for _, v := range flagAdditionalEventSources {
				if v == ".*" {
					allData = true
					break
				}
			}
			if allData {
				displayFlags = append(displayFlags, "all eligible event sources")
			} else {
				displayFlags = append(displayFlags, "additional_event_sources")
			}
		}
		if len(displayFlags) == 0 {
			displayFlags = append(displayFlags, "all eligible event sources")
		}

		sort.Strings(displayFlags)
		desc = fmt.Sprintf("CLI generated recording at %s containing: %s", time.Now().Format(time.RFC3339), strings.Join(displayFlags, ", "))
	}
	return desc
}

// printRecordingMetadata displays a summary of the requested recording properties to the user.
// This includes requested duration, start/end times, description, flags, and itemized source metadata.
func printRecordingMetadata(out io.Writer, duration time.Duration, startTime time.Time, endTime time.Time, desc string, includedFlags []string, bag *bagmetadatapb.BagMetadata) {
	fmt.Fprintf(out, "Recording metadata:\n")
	fmt.Fprintf(out, "  Requested duration: %v\n", duration)
	fmt.Fprintf(out, "  Requested start time:         %v\n", startTime.Format(time.RFC3339))
	fmt.Fprintf(out, "  Requested end time:           %v\n", endTime.Format(time.RFC3339))
	fmt.Fprintf(out, "  Description:        %s\n", desc)
	if len(includedFlags) > 0 {
		fmt.Fprintf(out, "  Included flags:     %s\n", strings.Join(includedFlags, ", "))
	}
	if len(flagAdditionalEventSources) > 0 {
		fmt.Fprintf(out, "  Custom sources:     %s\n", strings.Join(flagAdditionalEventSources, ", "))
	}
	if bag != nil && len(bag.GetEventSources()) > 0 {
		fmt.Fprintf(out, "  Created bag contains the following event sources:\n")
		for _, source := range bag.GetEventSources() {
			fmt.Fprintf(out, "    - %s: %d items, %d bytes\n",
				source.GetEventSourceWithTypeHints().GetEventSource(),
				source.GetNumLogItems(),
				source.GetNumBytes())
		}
	}
	fmt.Fprintf(out, "\n")
}

// sendCreateLocalRecordingRequest dispatches the creation request to the remote DataLogger service,
// applying an exponential backoff retry in case the workcell is momentarily unavailable.
func sendCreateLocalRecordingRequest(ctx context.Context, out io.Writer, client loggerpb.DataLoggerClient, req *loggerpb.CreateLocalRecordingRequest) (*loggerpb.CreateLocalRecordingResponse, error) {
	var resp *loggerpb.CreateLocalRecordingResponse
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	boff := backoff.WithContext(backoff.NewExponentialBackOff(), ctxWithTimeout)
	err := backoff.RetryNotify(func() error {
		var rpcErr error
		resp, rpcErr = client.CreateLocalRecording(ctxWithTimeout, req)
		if rpcErr != nil {
			if status.Code(rpcErr) == codes.Unavailable {
				return rpcErr // Retryable error (e.g., workcell offline momentarily)
			}
			return backoff.Permanent(rpcErr)
		}
		return nil
	}, boff, func(err error, d time.Duration) {
		color.C.Yellow().Fprintf(out, "Warning: failed to request recording, workcell may be offline. Retrying in %v...\n", d.Truncate(time.Millisecond))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create local recording on workcell: %v", err)
	}
	return resp, nil
}

// Poll until recording reaches a terminal state.
func (r *CreateCmdRunner) pollUploadProgress(ctx context.Context, out io.Writer, client bagpackagerpb.BagPackagerClient, bagID string) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Fprintf(out, "Waiting for recording %s to be uploaded from the workcell...", bagID)

	req := &bagpackagerpb.GetBagRequest{
		BagId: bagID,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, err := client.GetBag(ctx, req)
			if err != nil {
				if status.Code(err) == codes.NotFound || strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "failed to get bag record") {
					fmt.Fprintf(out, "\nWaiting for workcell to register recording with cloud...")
					continue
				}
				return err
			}

			bagMeta := resp.GetBag().GetBagMetadata()
			bagStatus := bagMeta.GetStatus().GetStatus()

			pct := 0.0
			if bagMeta.GetTotalBytes() > 0 {
				pct = float64(bagMeta.GetTotalUploadedBytes()) / float64(bagMeta.GetTotalBytes()) * 100.0
			}

			fmt.Fprintf(out, "\nUpload progress: %.1f%% (%d / %d items, %dMB / %dMB) - Status: %s",
				pct,
				bagMeta.GetTotalUploadedLogItems(), bagMeta.GetTotalLogItems(),
				bagMeta.GetTotalUploadedBytes()/(1024*1024), bagMeta.GetTotalBytes()/(1024*1024),
				bagStatus.String())

			if bagStatus >= bagmetadatapb.BagStatus_UPLOADED {
				fmt.Fprintln(out, "\n\nUpload finished!\n")
				if bagStatus == bagmetadatapb.BagStatus_FAILED {
					color.C.Yellow().Fprintf(out, "Warning: Bag upload finished with terminal failure state: %s. Reason: %s\n", bagStatus.String(), bagMeta.GetStatus().GetReason())
				} else if bagStatus == bagmetadatapb.BagStatus_UNCOMPLETABLE {
					color.C.Yellow().Fprintf(out, "Note: Bag upload finished in state: %s. There might be missing data, but the recording can still be generated. Reason: %s\n", bagStatus.String(), bagMeta.GetStatus().GetReason())
				}

				return nil
			}
		}
	}
}

func (r *CreateCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	startTime, endTime, err := parseAndValidateRecordingTimestamps(cmd)
	if err != nil {
		return err
	}

	duration := endTime.Sub(startTime)
	if err := validateRecordingDuration(duration); err != nil {
		return err
	}

	flagNames := []string{"include_text_logs", "include_flowstate_data", "include_scene_data", "include_robot_data", "include_perception_data", "include_debug_data"}
	finalEventSources, includedFlags, err := determineEventSourcesToRecord(cmd, r, flagNames)
	if err != nil {
		if err.Error() == "aborted" {
			// This causes a non-zero exit status as requested
			return fmt.Errorf("aborted")
		}
		return err
	}

	loggerClient, closer, err := r.NewDataLoggerClient(cmd)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	desc := generateRecordingDescription(includedFlags)
	req := &loggerpb.CreateLocalRecordingRequest{
		StartTime:            tpb.New(startTime),
		EndTime:              tpb.New(endTime),
		Description:          desc,
		EventSourcesToRecord: finalEventSources,
	}

	fmt.Fprintf(out, "\nRequesting workcell %q to create local recording...\n", flagWorkcellName)

	resp, err := sendCreateLocalRecordingRequest(ctx, out, loggerClient, req)
	if err != nil {
		return err
	}

	bagID := resp.GetBag().GetBagId()
	color.C.BlueBackground().White().Fprintf(out, "\nSuccessfully created recording on %q - Recording ID: %s", flagWorkcellName, bagID)
	fmt.Fprint(out, "\n\n")

	printRecordingMetadata(out, duration, startTime, endTime, desc, includedFlags, resp.GetBag())

	bagPackagerClient, err := r.NewBagPackagerClient(cmd)
	if err != nil {
		return err
	}

	fullOrgName := orgutil.QualifiedOrg(createParams.GetString(orgutil.KeyProject), createParams.GetString(orgutil.KeyOrganization))

	fmt.Fprintln(out, "\nWaiting for bag upload completion... (Press Ctrl+C to abort)\n")
	err = r.pollUploadProgress(ctx, out, bagPackagerClient, bagID)
	if err != nil {
		return err
	}

	// Prompt to generate recording.
	// Logic:
	// - If --skip_generate is true, do not prompt and do not generate.
	// - If --skip_generate is false AND quiet mode is on, silently generate.
	// - Otherwise, prompt the user for permission to generate.
	generate := false
	if flagSkipGenerate {
		generate = false
	} else if flagQuiet {
		generate = true
	} else {
		generate, err = r.PromptYesNo(cmd, "Do you want to generate the recording now?", promptutil.DefaultYes, promptutil.ReemitPromptOnInvalidInput)
		if err != nil {
			return err
		}
	}

	generateCmdStr := fmt.Sprintf("inctl recordings generate --recording_id %s --org %s", bagID, fullOrgName)
	if generate {
		fmt.Fprintf(out, "\nRunning: %s\n\n", generateCmdStr)

		// Copy all global flags/authentication settings bound to the create command into the generate command's params.
		// This ensures credentials like --env and --api-key correctly propagate since we bypass root flag parsing.
		for k, v := range createParams.AllSettings() {
			generateParam.Set(k, v)
		}
		// generateParam.Set(orgutil.KeyOrganization, createParams.GetString(orgutil.KeyOrganization))
		return r.RunGenerate(cmd, bagID)
	} else {
		fmt.Fprintf(out, "\nYou can generate it later by running:\n  %s\n", generateCmdStr)
	}

	return nil
}

func NewCreateCmd(runner *CreateCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &CreateCmdRunner{
			NewDataLoggerClient: func(cmd *cobra.Command) (loggerpb.DataLoggerClient, io.Closer, error) {
				conn, err := auth.NewCloudConnection(cmd.Context(), auth.WithFlagValues(createParams), auth.WithCluster(flagWorkcellName))
				if err != nil {
					return nil, nil, err
				}
				return loggerpb.NewDataLoggerClient(conn), conn, nil
			},
			NewBagPackagerClient: func(cmd *cobra.Command) (bagpackagerpb.BagPackagerClient, error) {
				return newBagPackagerClient(cmd.Context(), createParams)
			},
			PromptYesNo: promptutil.PromptYesNo,
		}
	}

	if runner.RunGenerate == nil {
		// This is a helper function for calling `inctl recordings generate` that runs when RunE is invoked.
		runner.RunGenerate = func(cmd *cobra.Command, createResponseBagID string) error {
			genCmd := NewGenerateCmd(nil)
			flagBagID = createResponseBagID // Set flag on the generate command invocation.

			// Pass the parent command's context so downstream RPCs don't panic with nil contexts.
			genCmd.SetContext(cmd.Context())
			return genCmd.RunE(genCmd, []string{})
		}
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new recording directly on the workcell",
		Long:  "Creates a new recording directly on the workcell using the DataLogger API. You can specify exact times and content groupings.\n\nLimits:\n  - Mixed data (anything beyond text-logs or flowstate-data) will fail if it exceeds a maximum of 10 minutes.\n  - Text/flowstate-only data will fail if it exceeds a maximum of 24 hours.",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
		Example: `  # Create a default recording (all data) for the last 5 minutes
  inctl recordings create --workcell my-workcell --org my-org

  # Create a recording containing everything, with a custom timeframe in UTC (answer yes on prompt)
  inctl recordings create --workcell my-workcell --org my-org \
    --start_timestamp 2024-08-20T12:00:00Z \
    --end_timestamp 2024-08-20T12:05:00Z

  # Create a recording with a custom timeframe in a specific timezone (e.g. PST, -07:00)
  inctl recordings create --workcell my-workcell --org my-org \
    --start_timestamp 2024-08-20T12:00:00-07:00 \
    --end_timestamp 2024-08-20T12:05:00-07:00

  # Record only text logs plus a specific custom event source
  inctl recordings create --workcell my-workcell --include_text_logs --additional_event_sources "^/my/custom/topic$"

  # Record multiple custom data sources using multiple flags or regex
  inctl recordings create --workcell my-workcell \
    --additional_event_sources "^/my/custom/topic1$" \
    --additional_event_sources "^/my/custom/topic2/.*"`,
	}

	flags := createCmd.Flags()
	flags.StringVar(&flagWorkcellName, "workcell", "", "The Kubernetes cluster to use.")
	flags.StringVar(&flagStartTimestamp, "start_timestamp", "", "Start timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z (UTC) or 2024-08-20T12:00:00-07:00 (UTC-7)")
	flags.StringVar(&flagEndTimestamp, "end_timestamp", "", "End timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z (UTC) or 2024-08-20T12:00:00-07:00 (UTC-7)")
	flags.StringVar(&flagDescription, "description", "", "A human-readable description for the recording.")
	flags.BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress interactive prompts. Prompts will resolve to their default values (e.g. deny missing-data prompt, accept generate prompt).")
	flags.BoolVar(&flagSkipGenerate, "skip_generate", false, "If set, unconditionally skips prompting for and generating the recording after creation.")

	// Data inclusion flags.
	flags.BoolVar(&flagIncludeAllData, "include_all_data", false, "Include all eligible event sources (.*). Use this flag to suppress the interactive prompt and intentionally record everything.")
	flags.BoolVar(&flagTextLogs, "include_text_logs", false, "Include text logs (/text-log-out, /asset-text-log-out).")
	flags.BoolVar(&flagFlowstateData, "include_flowstate_data", false, "Include Flowstate execution data.")
	flags.BoolVar(&flagSceneData, "include_scene_data", false, "Include scene and TF data.")
	flags.BoolVar(&flagRobotData, "include_robot_data", false, "Include robot statuses and trajectory plans.")
	flags.BoolVar(&flagPerceptionData, "include_perception_data", false, "Include perception data.")
	flags.BoolVar(&flagDebugData, "include_debug_data", false, "Include system metrics and error reports.")
	flags.StringSliceVar(&flagAdditionalEventSources, "additional_event_sources", []string{}, "Custom RE2 regex patterns of event sources to record.")

	createCmd.MarkFlagRequired("workcell")

	return orgutil.WrapCmd(createCmd, createParams, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))
}

func init() {
	RecordingsCmd.AddCommand(NewCreateCmd(nil))
}
