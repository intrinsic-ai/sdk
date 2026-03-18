// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	bagmetadatapb "intrinsic/logging/proto/bag_metadata_go_proto"
	bagpackagerpb "intrinsic/logging/proto/bag_packager_service_go_proto"
	loggerpb "intrinsic/logging/proto/logger_service_go_proto"
	"intrinsic/tools/inctl/util/promptutil"
)

type mockDataLoggerClient struct {
	loggerpb.DataLoggerClient
	CreateLocalRecordingFunc func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error)
}

func (m *mockDataLoggerClient) CreateLocalRecording(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
	if m.CreateLocalRecordingFunc != nil {
		return m.CreateLocalRecordingFunc(ctx, in, opts...)
	}
	return nil, errors.New("mock CreateLocalRecording should not be called directly")
}

type mockBagPackagerClientForCreate struct {
	bagpackagerpb.BagPackagerClient
	GetBagFunc      func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error)
	GenerateBagFunc func(ctx context.Context, in *bagpackagerpb.GenerateBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GenerateBagResponse, error)
}

func (m *mockBagPackagerClientForCreate) GetBag(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
	if m.GetBagFunc != nil {
		return m.GetBagFunc(ctx, in, opts...)
	}
	return nil, errors.New("mock GetBag should not be called directly")
}

func (m *mockBagPackagerClientForCreate) GenerateBag(ctx context.Context, in *bagpackagerpb.GenerateBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GenerateBagResponse, error) {
	if m.GenerateBagFunc != nil {
		return m.GenerateBagFunc(ctx, in, opts...)
	}
	return nil, errors.New("mock GenerateBag should not be called directly")
}

func TestCreateRecordingE(t *testing.T) {
	const (
		testOrg      = "test-org"
		testWorkcell = "test-workcell"
		testBagID    = "test-bag-id"
	)

	tests := []struct {
		name                           string
		args                           []string
		createFunc                     func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error)
		promptConfirmRecordAllResponse bool // Controls the answer to the missing data warning prompt when no flags are given.
		promptGenerateResponse         bool // Controls the answer to the "Do you want to generate the recording now?" prompt, which appears after upload completes.
		wantErr                        string
		wantOut                        string
		wantEventSourcesRegex          []string
		wantDescRegex                  string
		mockBagStatus                  bagmetadatapb.BagStatus_BagStatusEnum
	}{
		{
			name:                           "Successfully creates a recording with explicit flags and declines generate",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene", "--include_robot_data"},
			promptConfirmRecordAllResponse: true,
			promptGenerateResponse:         false,
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:               testBagID,
			wantEventSourcesRegex: []string{"/assets/.*/markers", "motion_planner_service.PlanTrajectory.debug_data"},
			wantDescRegex:         "CLI generated recording at .* containing: robot_data, scene",
			mockBagStatus:         bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:                           "Successfully creates a recording and accepts generate prompt",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene"},
			promptConfirmRecordAllResponse: true,
			promptGenerateResponse:         true,
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:               "Running: inctl recordings generate",
			wantEventSourcesRegex: []string{"/assets/.*/markers"},
			wantDescRegex:         "CLI generated recording at .* containing: scene",
			mockBagStatus:         bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:                           "Prompts and records all data when no flags are provided and user confirms",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg},
			promptConfirmRecordAllResponse: true,
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:               testBagID,
			wantEventSourcesRegex: []string{".*"},
			wantDescRegex:         "CLI generated recording at .* containing: all eligible event sources",
			mockBagStatus:         bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:                           "Aborts when no flags are provided and user declines confirmation",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg},
			promptConfirmRecordAllResponse: false,
			wantErr:                        "aborted",
			mockBagStatus:                  bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:    "Errors when mixed data exceeds 10 minutes",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene", "--start_timestamp", "2024-08-20T12:00:00Z", "--end_timestamp", "2024-08-20T13:00:00Z"},
			wantErr: "exceeds the 10-minute limit for mixed data",
		},
		{
			name:    "Errors when text/flowstate data exceeds 24 hours",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--include_text_logs", "--start_timestamp", "2024-08-20T12:00:00Z", "--end_timestamp", "2024-08-23T12:00:00Z"},
			wantErr: "exceeds the 24-hour limit for text/flowstate data",
		},
		{
			name:    "Errors when workcell flag is missing",
			args:    []string{"--org", testOrg},
			wantErr: "required flag(s) \"workcell\" not set",
		},
		{
			name:    "Errors when only start_timestamp is provided",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--start_timestamp", "2024-08-20T12:00:00Z"},
			wantErr: "must supply BOTH start_timestamp and end_timestamp, or NEITHER",
		},
		{
			name:                           "Errors when --quiet is provided without any data flags",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--quiet"},
			promptConfirmRecordAllResponse: false, // Doesn't matter, shouldn't be called
			promptGenerateResponse:         false, // Doesn't matter, shouldn't be called
			wantErr:                        "no data was requested to be included",
			mockBagStatus:                  bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:                           "Skips the missing-data prompt when --include_all_data is provided",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--include_all_data"},
			promptConfirmRecordAllResponse: false, // Doesn't matter, shouldn't be called
			promptGenerateResponse:         false, // Shouldn't be called if we say no
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:               testBagID,
			wantEventSourcesRegex: []string{".*"},
			mockBagStatus:         bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:    "Errors when only end_timestamp is provided",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--end_timestamp", "2024-08-20T12:00:00Z"},
			wantErr: "must supply BOTH start_timestamp and end_timestamp, or NEITHER",
		},
		{
			name:    "Errors when start_timestamp is invalid",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--start_timestamp", "invalid-date", "--end_timestamp", "2024-08-20T12:00:00Z"},
			wantErr: "invalid start timestamp: invalid-date",
		},
		{
			name:    "Errors when end_timestamp is invalid",
			args:    []string{"--workcell", testWorkcell, "--org", testOrg, "--start_timestamp", "2024-08-20T12:00:00Z", "--end_timestamp", "invalid-date"},
			wantErr: "invalid end timestamp: invalid-date",
		},
		{
			name: "Errors when creation backend fails",
			args: []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene"},
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return nil, errors.New("backend failed")
			},
			wantErr: "failed to create local recording on workcell: backend failed",
		},
		{
			name: "Successfully creates recording with explicit description",
			args: []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene", "--description", "my awesome recording"},
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:       testBagID,
			wantDescRegex: "my awesome recording",
			mockBagStatus: bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name: "Successfully passes multiple custom additional event sources",
			args: []string{"--workcell", testWorkcell, "--org", testOrg, "--additional_event_sources", "^/my/custom/topic1$", "--additional_event_sources", "^/my/custom/topic2/.*"},
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:               testBagID,
			wantEventSourcesRegex: []string{"^/my/custom/topic1$", "^/my/custom/topic2/.*"},
			wantDescRegex:         "CLI generated recording at .* containing: additional_event_sources",
			mockBagStatus:         bagmetadatapb.BagStatus_COMPLETED,
		},
		{
			name:                           "Warns when upload finishes with FAILED status",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene"},
			promptConfirmRecordAllResponse: true,
			promptGenerateResponse:         false,
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:       "Warning: Bag upload finished with terminal failure state: FAILED.",
			mockBagStatus: bagmetadatapb.BagStatus_FAILED,
		},
		{
			name:                           "Succeeds silently when upload finishes with UNCOMPLETABLE status",
			args:                           []string{"--workcell", testWorkcell, "--org", testOrg, "--include_scene"},
			promptConfirmRecordAllResponse: true,
			promptGenerateResponse:         false,
			createFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
				return &loggerpb.CreateLocalRecordingResponse{
					Bag: &bagmetadatapb.BagMetadata{
						BagId: testBagID,
					},
				}, nil
			},
			wantOut:       "Note: Bag upload finished in state: UNCOMPLETABLE. There might be missing data, but the recording can still be generated.",
			mockBagStatus: bagmetadatapb.BagStatus_UNCOMPLETABLE,
		},
	}

	// We disable the org check in tests because the test environment does not have
	// the necessary home directory configuration.
	originalCheckOrgExists := checkOrgExists
	checkOrgExists = false
	t.Cleanup(func() { checkOrgExists = originalCheckOrgExists })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Cleanup global flags
			originalFlagWorkcell := flagWorkcellName
			originalFlagStart := flagStartTimestamp
			originalFlagEnd := flagEndTimestamp
			t.Cleanup(func() {
				flagWorkcellName = originalFlagWorkcell
				flagStartTimestamp = originalFlagStart
				flagEndTimestamp = originalFlagEnd
				flagTextLogs = false
				flagFlowstateData = false
				flagScene = false
				flagRobotData = false
				flagPerception = false
				flagDebugData = false
				flagAdditionalEventSources = []string{}
				flagQuiet = false
				flagIncludeAllData = false
				flagSkipGenerate = false
			})
			flagWorkcellName = ""
			flagStartTimestamp = ""
			flagEndTimestamp = ""

			mockClient := &mockDataLoggerClient{
				CreateLocalRecordingFunc: func(ctx context.Context, in *loggerpb.CreateLocalRecordingRequest, opts ...grpc.CallOption) (*loggerpb.CreateLocalRecordingResponse, error) {
					if len(tc.wantEventSourcesRegex) > 0 {
						var sources []string
						for _, s := range in.EventSourcesToRecord {
							sources = append(sources, s)
						}
						sourceStr := strings.Join(sources, ",")
						// Validate that the request contains the expected event sources
						for _, r := range tc.wantEventSourcesRegex {
							assert.Contains(t, sourceStr, r)
						}
					}
					// If the test provides a custom mock function, run it.
					if tc.createFunc != nil {
						return tc.createFunc(ctx, in, opts...)
					}
					return nil, nil
				},
			}
			runner := &CreateCmdRunner{
				NewDataLoggerClient: func(cmd *cobra.Command) (loggerpb.DataLoggerClient, io.Closer, error) {
					return mockClient, nil, nil
				},
				NewBagPackagerClient: func(cmd *cobra.Command) (bagpackagerpb.BagPackagerClient, error) {
					// Provides a dummy BagPackagerClient that always returns a COMPLETED bag,
					// or the status provided by the test case.
					status := tc.mockBagStatus
					if status == bagmetadatapb.BagStatus_UNSET {
						status = bagmetadatapb.BagStatus_COMPLETED
					}
					return &mockBagPackagerClientForCreate{
						GetBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
							return &bagpackagerpb.GetBagResponse{
								Bag: &bagpackagerpb.BagRecord{
									BagMetadata: &bagmetadatapb.BagMetadata{
										Status: &bagmetadatapb.BagStatus{
											Status: status,
										},
									},
								},
							}, nil
						},
					}, nil
				},
				PromptYesNo: func(cmd *cobra.Command, prompt string, defaultBehavior promptutil.DefaultBehavior, invalidInputBehavior promptutil.InvalidInputBehavior) (bool, error) {
					// Route the prompts to the specific test case fields to simulate user interaction
					if strings.Contains(prompt, "ALL data (.*)") {
						return tc.promptConfirmRecordAllResponse, nil
					}
					if strings.Contains(prompt, "generate the recording now?") {
						return tc.promptGenerateResponse, nil
					}
					return false, nil
				},
				RunGenerate: func(cmd *cobra.Command, createResponseBagID string) error {
					assert.Equal(t, testBagID, createResponseBagID)
					genRunner := &GenerateCmdRunner{
						NewClient: func(cmd *cobra.Command) (bagpackagerpb.BagPackagerClient, error) {
							return &mockBagPackagerClientForCreate{
								GetBagFunc: func(ctx context.Context, in *bagpackagerpb.GetBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GetBagResponse, error) {
									return &bagpackagerpb.GetBagResponse{
										Bag: &bagpackagerpb.BagRecord{
											BagMetadata: &bagmetadatapb.BagMetadata{
												Status: &bagmetadatapb.BagStatus{
													Status: bagmetadatapb.BagStatus_COMPLETED,
												},
											},
											// BagFile intentionally omitted so GenerateBag gets invoked without failing early
										},
									}, nil
								},
								GenerateBagFunc: func(ctx context.Context, in *bagpackagerpb.GenerateBagRequest, opts ...grpc.CallOption) (*bagpackagerpb.GenerateBagResponse, error) {
									assert.Equal(t, testBagID, in.GetBagId())
									return &bagpackagerpb.GenerateBagResponse{}, nil
								},
							}, nil
						},
					}
					genCmd := NewGenerateCmd(genRunner)
					// Use exactly the same flags as the real implementation to ensure test coverage on the command interface.
					genCmd.SetArgs([]string{
						"--recording_id", createResponseBagID,
						"--org", testOrg,
					})
					// Capture output so it doesn't pollute test logs
					genCmd.SetOut(io.Discard)
					genCmd.SetErr(io.Discard)
					return genCmd.Execute()
				},
			}

			var out bytes.Buffer
			rootCmd := &cobra.Command{Use: "inctl"}
			recordingsCmd := &cobra.Command{Use: "recordings"}
			createCmd := NewCreateCmd(runner)
			recordingsCmd.AddCommand(createCmd)

			rootCmd.AddCommand(recordingsCmd)
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs(append([]string{"recordings", "create"}, tc.args...))

			err := rootCmd.Execute()

			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				assert.NoError(t, err)
				if tc.wantOut != "" {
					assert.Contains(t, out.String(), tc.wantOut)
				}
				if tc.wantDescRegex != "" {
					assert.Regexp(t, tc.wantDescRegex, out.String())
				}
			}
		})
	}
}
