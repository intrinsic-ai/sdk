// Copyright 2023 Intrinsic Innovation LLC

// Package shared provides data types that client tooling uses as well for static typed api boundaries.
package shared

import "encoding/json"

// ConfigureData is the data type used during the configuration push by inctl.
type ConfigureData struct {
	Config      []byte `json:"config"`
	Hostname    string `json:"hostname"`
	Role        string `json:"role"`
	Cluster     string `json:"cluster"`
	Private     bool   `json:"private"`
	Region      string `json:"region"`
	Replace     bool   `json:"replace"`
	AutoUpdate  bool   `json:"auto_update"`
	DisplayName string `json:"display_name"`
	Location    string `json:"location"`
	// CreatedByTest is only used for automated testing, and contains the ID of the test that is
	// registering this device. It is used to label the resources (in particular the Robot CR) so we
	// can clean it up.
	CreatedByTest string `json:"created_by_test"`
}

// TokenPlaceholder is used to modify the config by string replacement.
const TokenPlaceholder = "INTRINSIC_BOOTSTRAP_TOKEN_PLACEHOLDER"

// DeviceInfo is the data type used to upload the key from a device to the install registry and is
// reported to the devicemanager on claim
type DeviceInfo struct {
	Key       string `json:"key"`
	HasGPU    bool   `json:"has_gpu"`
	CanDoReal bool   `json:"can_do_real"`
	LastErr   string `json:"last_error"`
	// FullID is used for the internal communication to tell the full id of the device registered via
	// partial id.
	FullID string `json:"full_id"`
	// Version provides the version running on the OS for guiding functionality that needs minimum
	// versions of the installation image.
	Version string `json:"version"`
}

// Nameservers sets DNS servers and search domains.
type Nameservers struct {
	// Search is a list of DNS search domains.
	Search []string `json:"search" jsonschema:"example=lab.intrinsic.ai"`

	// Addresses is a list of DNS servers.
	Addresses []string `json:"addresses" jsonschema:"format=ipv4"`
}

// These constants must match the enum in clustermanager_api.proto.
const (
	// EtherTypeIP is the IP protocol (realtime or nonrealtime).
	EtherTypeIP = 0
	// EtherTypeEtherCAT is the EtherCAT protocol.
	EtherTypeEtherCAT = 1
)

// Interface represents a network interface configuration.
type Interface struct {
	// DHCP4 enables or disables DHCP on the interface.
	DHCP4 bool `json:"dhcp4"`

	// Gateway4 specifies the default gateway, if DHCP4 is disabled.
	Gateway4 string `json:"gateway4" jsonschema:"format=ipv4"`

	// NOT IMPLEMENTED: DHCP6 enables or disables DHCP on the interface.
	DHCP6 *bool `json:"dhcp6"`

	// NOT IMPLEMENTED: Gateway6 specifies the default gateway, if DHCP6 is disabled.
	Gateway6 string `json:"gateway6" jsonschema:"format=ipv6"`

	// MTU is the maximum transfer unit of the device, in bytes. If omitted,
	// the system will choose a default.
	MTU int64 `json:"mtu" jsonschema:"example=9000"`

	// Nameservers sets DNS servers and search domains.
	Nameservers Nameservers `json:"nameservers"`

	// Addresses specifies the IP addresses. It is required if DHCP4 is
	// disabled. If DHCP4 is enabled it can be optionally used for additional
	// addresses.
	Addresses []string `json:"addresses" validate:"required_without=DHCP4,omitempty,min=1" jsonschema:"format=ipv4"`

	// Realtime identifies this interface to be used for realtime communication
	// with the robot.
	Realtime bool `json:"realtime"`

	// EtherType specifies the protocol used on this interface.
	EtherType int64 `json:"ether_type"`

	// DisplayName is a pretty name set by the user.
	DisplayName string `json:"display_name"`
}

// String implements fmt.Stringer for logging purposes.
func (i Interface) String() string {
	r, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(r)
}

// Status represents the current OS status. It is similar to config.Config but
// contains the current status instead of the wanted status.
type Status struct {
	NodeName      string                     `json:"nodeName"`
	Hostname      string                     `json:"hostname"`
	Network       map[string]StatusInterface `json:"network"`
	BuildID       string                     `json:"buildId"`
	ImageType     string                     `json:"imageType"`
	Board         string                     `json:"board"`
	ActiveCopy    string                     `json:"activeCopy"`
	OEMVars       map[string]string          `json:"oemVars"`
	NetworkIssues []string                   `json:"networkIssues"`
}

// StatusInterface represents a network interface.
type StatusInterface struct {
	Up              bool     `json:"up"`
	MacAddress      string   `json:"hwaddr"`
	MTU             int      `json:"mtu"`
	IPAddress       []string `json:"addresses"`
	Speed           int      `json:"speed,omitempty"`
	Realtime        bool     `json:"realtime"`
	HasCarrier      bool     `json:"carrier"`
	HasDefaultRoute bool     `json:"has_default_route"`
	DisplayName     string   `json:"display_name"`
	DefaultGateway  string   `json:"default_gateway"`
	// If empty / not present, this is an older OS that doesn't report compatibility.
	SupportedEtherType []int64 `json:"supported_ether_types"`
}

// PingCommand allows to trigger an ICMP ping from the device.
// This can be used to test (L3) connectivity to networked peripherals.
type PingCommand struct {
	Target string `json:"target"`
	// Maximum wait time in milliseconds
	Duration int `json:"duration"`
}

// PingResponse represents the response from the ping command.
type PingResponse struct {
	Success bool `json:"success"`
}
