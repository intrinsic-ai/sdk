// Copyright 2023 Intrinsic Innovation LLC

package device

// ModeRequest is a request message to update the mode field
type ModeRequest struct {
	Mode string `json:"mode"`
}

// ClusterProjectTargetResponse is the response to the cluster project target request
type ClusterProjectTargetResponse struct {
	OS   string `json:"os"`
	Base string `json:"base"`
}
