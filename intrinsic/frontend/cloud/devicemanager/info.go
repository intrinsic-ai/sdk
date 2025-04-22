// Copyright 2023 Intrinsic Innovation LLC

package device

import (
	"log/slog"
	"time"

	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
	"intrinsic/kubernetes/inversion/inversion"
)

// UpdateProgress contains the information regarding the update
// progress if the Node's platform is updating.
type UpdateProgress struct {
	// CurrentStep contains the string that describes
	// the current update step.
	CurrentStep string `json:"currentStep"`
	// CurrentStepProgress is a number [0, 100] showing the percentage of
	// the current update step.
	// This variable is optional, however when it's not given, a value `0`
	// would also be reasonable.
	CurrentStepProgress int `json:"currentStepProgress,omitempty"`
	// TotalUpdateProgress is a number [0, 100] showing the percentage of
	// the whole update progress aggregated from all update steps involved.
	// This variable is optional, however when it's not given, a value `0`
	// would also be reasonable.
	TotalUpdateProgress int `json:"totalUpdateProgress,omitempty"`
}

// IPCInfo contains update information about an IPC Node in a cluster
type IPCInfo struct {
	Name       string `json:"name,omitempty"`
	LastSeenTS string `json:"lastSeenTS,omitempty"`
	// OSVersion is in API format
	OSVersion      string          `json:"osVersion,omitempty"`
	ISControl      bool            `json:"isControl,omitempty"`
	UpdateProgress *UpdateProgress `json:"updateProgress,omitempty"`
}

// LastSeen returns when the IPC was last seen (parsing the heartbeat string timestamp into a time.Time)
func (i *IPCInfo) LastSeen() (time.Time, error) {
	return time.Parse(time.RFC3339, i.LastSeenTS)
}

// Info contains update information about a cluster
type Info struct {
	Cluster     string `json:"cluster,omitempty"`
	State       string `json:"state,omitempty"`
	OSState     string `json:"osState,omitempty"`
	BaseState   string `json:"baseState,omitempty"`
	BaseManager string `json:"baseManager,omitempty"`
	Mode        string `json:"mode,omitempty"`
	// CurrentBase is in API format
	CurrentBase string `json:"currentBase,omitempty"`
	// TargetBase is in API format
	TargetBase string `json:"targetBase,omitempty"`
	// CurrentOS is in API format
	CurrentOS string `json:"currentOS,omitempty"`
	// TargetOS is in API format
	TargetOS string `json:"targetOS,omitempty"`
	// RollbackOS is in API format
	RollbackOS string `json:"rollbackOS,omitempty"`
	// RollbackBase is in API format
	RollbackBase   string          `json:"rollbackBase,omitempty"`
	LastSeenTS     string          `json:"lastSeenTS,omitempty"`
	Nodes          []*IPCInfo      `json:"nodes,omitempty"`
	UpdateProgress *UpdateProgress `json:"updateProgress,omitempty"`
}

// LastSeen returns when the control plane was last seen
func (i *Info) LastSeen() (time.Time, error) {
	return time.Parse(time.RFC3339, i.LastSeenTS)
}

// RollbackAvailable reports whether a rollback is available according to this info object
func (i *Info) RollbackAvailable() bool {
	if i.BaseManager != "inversion" {
		// No rollback available for non-inversion base clusters
		return false
	}
	return i.RollbackOS != "" && i.RollbackBase != ""
}

// OSUpdateDone returns true when the OS is in a deployed state
func (i *Info) OSUpdateDone() bool {
	return i.OSState == "Deployed"
}

// UpdateDone returns true when the system is in a deployed state
func (i *Info) UpdateDone() bool {
	return i.State == "Deployed"
}

func parseCurrentStep(currentStepStr string) clustermanagerpb.UpdateProgressStep {
	switch currentStepStr {
	case "Completed":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_COMPLETED
	case "NewOSUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_NEW_OS_UPDATE
	case "PendingOSUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_PENDING_OS_UPDATE
	case "DownloadOSUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_DOWNLOAD_OS_UPDATE
	case "CopyOSUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_COPY_OS_UPDATE
	case "ApplyOSUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_APPLY_OS_UPDATE
	case "NewBaseUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_NEW_BASE_UPDATE
	case "PendingBaseUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_PENDING_BASE_UPDATE
	case "DownloadBaseUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_DOWNLOAD_BASE_UPDATE
	case "ApplyBaseUpdate":
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_APPLY_BASE_UPDATE
	default:
		slog.Warn("Invalid UpdateProgressStep: %s", currentStepStr)
		return clustermanagerpb.UpdateProgressStep_UPDATE_PROGRESS_UNSPECIFIED
	}
}

func convertToPbUpdateProgressFormat(up *UpdateProgress) *clustermanagerpb.UpdateProgress {
	if up == nil {
		return nil
	}

	return &clustermanagerpb.UpdateProgress{
		CurrentStep:         parseCurrentStep(up.CurrentStep),
		CurrentStepProgress: int32(up.CurrentStepProgress),
		TotalUpdateProgress: int32(up.TotalUpdateProgress),
	}
}

func convertToDeviceUpdateProgressFormat(up *inversion.UpdateProgress) *UpdateProgress {
	if up == nil {
		return nil
	}

	return &UpdateProgress{
		CurrentStep:         up.CurrentStep,
		CurrentStepProgress: up.CurrentStepProgress,
		TotalUpdateProgress: up.TotalUpdateProgress,
	}
}
