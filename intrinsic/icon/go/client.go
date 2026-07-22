// Copyright 2023 Intrinsic Innovation LLC

// Package icon provides a Go client library for ICON, an industrial robot
// control framework.
//
// Initialize a client connection to the server over GRPC using:
//
//	client, err := InitClient(addr, opts)
//	if err != nil {
//		// error handling
//	}
//	defer client.Close()
package icon

import (
	"context"

	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpcpb "intrinsic/icon/proto/v1/service_go_proto"
	pb "intrinsic/icon/proto/v1/service_go_proto"
	typespb "intrinsic/icon/proto/v1/types_go_proto"
	contextpb "intrinsic/logging/proto/context_go_proto"

	epb "google.golang.org/protobuf/types/known/emptypb"
)

// resourceInstanceHeaderKey is an outgoing metadata key used to select
// among multiple resource instances that are routed through a single address.
var resourceInstanceHeaderKey = "x-resource-instance-name"

// HardwareGroup is a group of hardware modulesto disable.
type HardwareGroup int

// A hardware group to disable, all or only operational hardware modules.
const (
	AllHardware HardwareGroup = iota
	OperationalHardwareOnly
)

// BoolPartPropertyValue creates a PartPropertyValue proto that holds a boolean value.
func BoolPartPropertyValue(v bool) *pb.PartPropertyValue {
	return &pb.PartPropertyValue{Value: &pb.PartPropertyValue_BoolValue{BoolValue: v}}
}

// Float64PartPropertyValue creates a PartPropertyValue proto that holds a float64 value.
func Float64PartPropertyValue(v float64) *pb.PartPropertyValue {
	return &pb.PartPropertyValue{Value: &pb.PartPropertyValue_DoubleValue{DoubleValue: v}}
}

// The Client interface provides the top-level ICON API functionality.
type Client interface {
	// ActionSignatureByName gets details of an action type, by name.
	ActionSignatureByName(ctx context.Context, name string) (*typespb.ActionSignature, error)

	// ActionSignatures lists details of all available actions.
	ActionSignatures(ctx context.Context) ([]*typespb.ActionSignature, error)

	// ClearFaults clears all faults and returns the server to an enabled state.
	// Returns OK if faults were successfully cleared and the server is enabled.
	//
	// NOTE: Clearing faults is something the user does directly. DO NOT call this
	// from library code automatically to make things more convenient, ESPECIALLY
	// not in connection with re-enabling the server afterwards! Human users must
	// be able to rely on the robot to stay still unless they explicitly clear the
	// fault(s) and enable it again.
	//
	// Some classes of faults (internal server errors, or issues that have a
	// physical root cause) may require additional server- or hardware-specific
	// mitigation before ClearFaults can successfully clear the fault.
	ClearFaults(ctx context.Context) error

	// Close closes the client's connection to the server.
	// Subsequent calls using this Client object will return an error.
	Close() error

	// CompatibleParts lists the parts that are compatible with all listed
	// action types. If no actionTypes are listed, returns all parts.
	CompatibleParts(ctx context.Context, actionTypes ...string) ([]string, error)

	// Config returns a list of part configs (for example, the number of
	// DOFs for a robot arm) and the server config, which are fixed
	// properties for the lifetime of the server.
	Config(ctx context.Context) ([]*typespb.PartConfig, *typespb.ServerConfig, error)

	// Disable disables the specified hardware group on the server.
	// Cancels sessions that depend on that hardware group.
	//
	// NOTE: Disabling a server is something the user does directly. DO NOT call
	// this from library code automatically to make things more convenient. Human
	// users must be able to rely on the robot to stay enabled unless they
	// explicitly disable it (or the robot encounters a fault).
	// If the operational state of the server is already DISABLED, then this does
	// nothing (successfully). Returns an error if the server is faulted.
	//
	// With `group` set to `AllHardware`, all hardware will be disabled.
	// With `group` set to `OperationalHardwareOnly`, parts that only use hardware
	// modules that are configured with
	// `IconMainConfig.hardware_config.cell_control_hardware` will be skipped,
	// keeping them enabled if they are enabled.
	// One use case is to integrate cell-level control where
	// operational robot hardware can be paused such that automatic
	// mode is not needed, while still reading/writing input/output on a fieldbus
	// hardware module for cell-level control.
	Disable(ctx context.Context, group HardwareGroup) error

	// Enable enables all parts on the server, which performs all steps necessary
	// to ge the parts ready to receive commands.
	//
	// NOTE: Enabling a server is something the user does directly. DO NOT call
	// this from library code automatically to make things more convenient. Human
	// users must be able to rely on the robot to stay still unless they enable
	// it.
	//
	// If the operational state of the server is already ENABLED, then this does
	// nothing (successfully). Returns an error if the server is faulted.
	Enable(ctx context.Context) error

	// IsActionCompatible reports whether actions of type actionType are compatible
	// with a part.
	IsActionCompatible(ctx context.Context, actionType string, part string) (bool, error)

	// Returns the summarized status of the server, including a fault reason if the server is faulted.
	OperationalStatus(ctx context.Context) (*typespb.OperationalStatus, error)

	// Returns the status of cell control hardware, which is marked with
	// `IconMainConfig.hardware_config.cell_control_hardware`.
	// Cell control hardware is a group of hardware modules that does not inherit
	// faults from operational hardware, and only gets disabled when a
	// any cell control hardware module faults (or when `Disable` is called).
	CellControlHardwareStatus(ctx context.Context) (*typespb.OperationalStatus, error)

	// Returns the current speed override value (a number between 0 and 1)
	GetSpeedOverride(ctx context.Context) (float64, error)

	// Sets the speed override to the given value (must be a number between 0 and 1).
	// Compatible ICON actions will attempt to scale down their speed according to this value.
	SetSpeedOverride(ctx context.Context, newSpeedOverride float64) error

	// Parts lists all available parts.
	Parts(ctx context.Context) ([]string, error)

	// StartSession starts a session which can be used to command robot motion.
	// The ICON logs are tagged with the provided `logContext`, which can be nil.
	StartSession(ctx context.Context, parts []string, logContext *contextpb.Context) (*Session, error)

	// Status returns the server status, including a list of part statuses,
	// representing a snapshot of the robot's state.
	Status(ctx context.Context) (*pb.GetStatusResponse, error)

	// Calls GetPartProperties on the ICON server and returns
	// the values of all part properties.
	//
	// Returns a GetPartPropertiesResponse proto that contains
	// * The control timestamp at the time the properties were reported
	// * The wall time at the time the properties were reported
	// * A map from part name to a map from property name to value.
	//	For instance: {'robot': {'motor_0_current_amps': 2.0}}
	GetPartProperties(ctx context.Context) (*pb.GetPartPropertiesResponse, error)

	// Sets part properties.
	//
	// Check the output of get_part_properties to learn the available properties
	// and their types.
	//
	// `part_properties` is a map from part name to a map from property name to
	// value. For instance: {'robot': {'internal_controller_p_value': icon.Float64PartPropertyValue(1.23)}}
	SetPartProperties(ctx context.Context, newPartProperties map[string]map[string]*pb.PartPropertyValue) error

	// Restarts the ICON server.
	//
	// This is mostly useful to force ICON to reload its configuration, including kinematics models
	// and limits.
	//
	// NOTE: Restarting the server also ends all Sessions, and shuts down any hardware safely.
	RestartServer(ctx context.Context) error
}

// clientOptions configure an InitClient call. clientOptions are set by the
// ClientOption values passed to InitClient.
type clientOptions struct {
	serverInstanceHeaderValue string
	dialOptions               []grpc.DialOption
}

// ClientOption configures the client.
type ClientOption interface {
	apply(*clientOptions)
}

// funcClientOption wraps a function that modifies clientOptions into an
// implementation of the ClientOption interface.
type funcClientOption struct {
	f func(*clientOptions)
}

func (fco *funcClientOption) apply(co *clientOptions) {
	fco.f(co)
}

func newFuncClientOption(f func(*clientOptions)) *funcClientOption {
	return &funcClientOption{
		f: f,
	}
}

// WithDialOptions sets the gRPC dial options to use when connecting to the server.
func WithDialOptions(do ...grpc.DialOption) ClientOption {
	return newFuncClientOption(func(co *clientOptions) {
		co.dialOptions = do
	})
}

// WithServerInstanceHeaderValue sets the server instance name to send in the metadata
// of all outgoing requests. Under certain network configurations, this can be
// used to select among multiple ICON server instances that are routed through
// a single address. If name is non-empty, an HTTP header metadata field
// "x-resource-instance-name" will be added with a value of name, unless another
// value is provided using the WithServerInstanceHeaderKey option.
func WithServerInstanceHeaderValue(name string) ClientOption {
	return newFuncClientOption(func(co *clientOptions) {
		co.serverInstanceHeaderValue = name
	})
}

// grpcClient is a connection to the ICON server, providing ICON API methods.
type grpcClient struct {
	conn                      *grpc.ClientConn
	client                    grpcpb.IconApiClient
	serverInstanceHeaderValue string
}

// InitClient connects to the ICON server at `addr` and initializes a Client object.
func InitClient(addr string, opts ...ClientOption) (Client, error) {
	copts := clientOptions{}
	for _, opt := range opts {
		opt.apply(&copts)
	}
	conn, err := grpc.Dial(addr, copts.dialOptions...)
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		conn:                      conn,
		client:                    grpcpb.NewIconApiClient(conn),
		serverInstanceHeaderValue: copts.serverInstanceHeaderValue,
	}, nil
}

// InitClientFromConn initializes a Client object from an existing gRPC connection.
func InitClientFromConn(conn *grpc.ClientConn, opts ...ClientOption) Client {
	copts := clientOptions{}
	for _, opt := range opts {
		opt.apply(&copts)
	}
	return &grpcClient{
		conn:                      conn,
		client:                    grpcpb.NewIconApiClient(conn),
		serverInstanceHeaderValue: copts.serverInstanceHeaderValue,
	}
}

// SetResourceInstanceHeaderOutgoingMetadata sets the x-resource-instance-name header in the outgoing metadata.
func SetResourceInstanceHeaderOutgoingMetadata(ctx context.Context, resourceInstanceHeaderValue string) context.Context {
	if resourceInstanceHeaderValue == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, resourceInstanceHeaderKey, resourceInstanceHeaderValue)
}

func (c *grpcClient) addOutgoingMetadata(ctx context.Context) context.Context {
	return SetResourceInstanceHeaderOutgoingMetadata(ctx, c.serverInstanceHeaderValue)
}

func (c *grpcClient) Close() error {
	err := c.conn.Close()
	// Log the error, since this may have been called with `defer client.Close()`.
	if err != nil {
		log.Errorf("Failed to close ICON client connection: %v", err)
	}
	return err
}

func (c *grpcClient) ActionSignatureByName(ctx context.Context, name string) (*typespb.ActionSignature, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.GetActionSignatureByName(ctx, &pb.GetActionSignatureByNameRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return resp.ActionSignature, nil
}

func (c *grpcClient) Config(ctx context.Context) ([]*typespb.PartConfig, *typespb.ServerConfig, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.GetConfig(ctx, &pb.GetConfigRequest{})
	if err != nil {
		return nil, nil, err
	}
	return resp.PartConfigs, resp.ServerConfig, nil
}

func (c *grpcClient) Status(ctx context.Context) (*pb.GetStatusResponse, error) {
	ctx = c.addOutgoingMetadata(ctx)
	return c.client.GetStatus(ctx, &pb.GetStatusRequest{})
}

func (c *grpcClient) IsActionCompatible(ctx context.Context, actionType string, part string) (bool, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.IsActionCompatible(ctx, &pb.IsActionCompatibleRequest{
		ActionTypeName: actionType,
		SlotData:       &pb.IsActionCompatibleRequest_PartName{PartName: part},
	})
	if err != nil {
		return false, err
	}
	return resp.IsCompatible, nil
}

func (c *grpcClient) ActionSignatures(ctx context.Context) ([]*typespb.ActionSignature, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.ListActionSignatures(ctx, &pb.ListActionSignaturesRequest{})
	if err != nil {
		return nil, err
	}
	return resp.ActionSignatures, nil
}

func (c *grpcClient) CompatibleParts(ctx context.Context, actionTypes ...string) ([]string, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.ListCompatibleParts(ctx, &pb.ListCompatiblePartsRequest{ActionTypeNames: actionTypes})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}

func (c *grpcClient) Parts(ctx context.Context) ([]string, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.ListParts(ctx, &pb.ListPartsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}

func (c *grpcClient) StartSession(ctx context.Context, parts []string, logContext *contextpb.Context) (*Session, error) {
	ctx = c.addOutgoingMetadata(ctx)
	return newSession(ctx, c, parts, logContext)
}

func (c *grpcClient) Enable(ctx context.Context) error {
	ctx = c.addOutgoingMetadata(ctx)
	_, err := c.client.Enable(ctx, &pb.EnableRequest{})
	return err
}

func (c *grpcClient) Disable(ctx context.Context, group HardwareGroup) error {
	ctx = c.addOutgoingMetadata(ctx)
	req := &pb.DisableRequest{}
	if group == OperationalHardwareOnly {
		req.Group = pb.DisableRequest_OPERATIONAL_HARDWARE_ONLY
	} else {
		req.Group = pb.DisableRequest_ALL_HARDWARE
	}

	_, err := c.client.Disable(ctx, req)
	return err
}

func (c *grpcClient) ClearFaults(ctx context.Context) error {
	ctx = c.addOutgoingMetadata(ctx)
	_, err := c.client.ClearFaults(ctx, &pb.ClearFaultsRequest{})
	return err
}

func (c *grpcClient) OperationalStatus(ctx context.Context) (*typespb.OperationalStatus, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.GetOperationalStatus(ctx, &pb.GetOperationalStatusRequest{})
	if err != nil {
		return nil, err
	}
	return resp.OperationalStatus, nil
}

func (c *grpcClient) CellControlHardwareStatus(ctx context.Context) (*typespb.OperationalStatus, error) {
	ctx = c.addOutgoingMetadata(ctx)
	req := &pb.GetOperationalStatusRequest{}
	resp, err := c.client.GetOperationalStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.CellControlHardwareStatus, nil
}

func (c *grpcClient) GetSpeedOverride(ctx context.Context) (float64, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.GetSpeedOverride(ctx, &pb.GetSpeedOverrideRequest{})
	if err != nil {
		return 1.0, err
	}
	return resp.GetOverrideFactor(), nil
}

func (c *grpcClient) SetSpeedOverride(ctx context.Context, newSpeedOverride float64) error {
	ctx = c.addOutgoingMetadata(ctx)
	_, err := c.client.SetSpeedOverride(ctx, &pb.SetSpeedOverrideRequest{OverrideFactor: newSpeedOverride})
	return err
}

func (c *grpcClient) GetPartProperties(ctx context.Context) (*pb.GetPartPropertiesResponse, error) {
	ctx = c.addOutgoingMetadata(ctx)
	resp, err := c.client.GetPartProperties(ctx, &pb.GetPartPropertiesRequest{})
	if err != nil {
		return &pb.GetPartPropertiesResponse{}, err
	}
	return resp, nil
}

func (c *grpcClient) SetPartProperties(ctx context.Context, newPartProperties map[string]map[string]*pb.PartPropertyValue) error {
	ctx = c.addOutgoingMetadata(ctx)
	partsMap := map[string]*pb.PartPropertyValues{}
	for partName, properties := range newPartProperties {
		partsMap[partName] = &pb.PartPropertyValues{PropertyValuesByName: properties}
	}
	_, err := c.client.SetPartProperties(ctx, &pb.SetPartPropertiesRequest{PartPropertiesByPartName: partsMap})
	return err
}

func (c *grpcClient) RestartServer(ctx context.Context) error {
	ctx = c.addOutgoingMetadata(ctx)
	_, err := c.client.RestartServer(ctx, &epb.Empty{})
	return err
}
