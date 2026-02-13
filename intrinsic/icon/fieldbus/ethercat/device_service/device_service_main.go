// Copyright 2023 Intrinsic Innovation LLC

// Package main provides an EtherCAT device service, that loads the requested ESI files from the
// data assets server and returns the configuration of the device with these ESI files.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	"intrinsic/icon/fieldbus/ethercat/device_service/deviceservice"
	"intrinsic/production/intrinsic"

	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	dagrpcpb "intrinsic/assets/data/proto/v1/data_assets_go_proto"
	ssgrpcpb "intrinsic/assets/services/proto/v1/service_state_go_proto"
	sspb "intrinsic/assets/services/proto/v1/service_state_go_proto"
	dscpb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_config_go_proto"
	dsgrpcpb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_go_proto"
	dspb "intrinsic/icon/fieldbus/ethercat/device_service/v1/device_service_go_proto"
	rtpb "intrinsic/resources/proto/runtime_context_go_proto"
	espb "intrinsic/util/status/extended_status_go_proto"
)

var (
	// http server
	runtimeContext = flag.String(
		"runtime_context_file",
		"/etc/intrinsic/runtime_config.pb",
		"Path to resource instance's runtime context file.")
	// Address of the workcell cluster service.
	workcellClusterServiceAddress = flag.String("workcell_cluster_service", "istio-ingressgateway.app-ingress.svc.cluster.local:80", "The workcell cluster service address")
)

func main() {
	intrinsic.Init()
	ctx := context.Background()
	protoData, err := os.ReadFile(*runtimeContext)
	if err != nil {
		log.ErrorContextf(ctx, "Failed to open file: %q with err: %v", *runtimeContext, err)
		log.ExitContext(ctx, err)
	}

	rtCtx := &rtpb.RuntimeContext{}
	if err := proto.Unmarshal(protoData, rtCtx); err != nil {
		log.ErrorContextf(ctx, "Failed to unmarshal file: %v", err)
		log.ExitContext(ctx, err)
	}

	serviceConfig := &dscpb.DeviceServiceConfig{}
	if err := rtCtx.GetConfig().UnmarshalTo(serviceConfig); err != nil {
		log.ErrorContextf(ctx, "Failed to unmarshal config: %v", err)
		log.ExitContext(ctx, err)
	}

	log.InfoContextf(ctx, "Got proto: %v", rtCtx)

	address := fmt.Sprintf(":%d", rtCtx.GetPort())
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.ExitContextf(ctx, "Server failed to listen at %q: %v", address, err)
	}
	log.InfoContextf(ctx, "Server is now listening at %q", address)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	log.InfoContextf(ctx, "Connecting to the workcell cluster service %q", *workcellClusterServiceAddress)
	workcellConn, err := grpc.NewClient(
		*workcellClusterServiceAddress,
		opts...,
	)
	if err != nil {
		log.ExitContextf(ctx, "Failed to establish connection to the workcell cluster service %q: %v", *workcellClusterServiceAddress, err)
	}
	defer workcellConn.Close() // Close the connection on exit
	daClient := dagrpcpb.NewDataAssetsClient(workcellConn)

	manager := newServiceManager(ctx, serviceConfig, daClient)

	grpcServer := grpc.NewServer()
	dsgrpcpb.RegisterDeviceServiceServer(grpcServer, manager)
	ssgrpcpb.RegisterServiceStateServer(grpcServer, manager)
	reflection.Register(grpcServer)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.ErrorContextf(ctx, "Server failed to serve: %v", err)
		log.ExitContext(ctx, err)
	}
}

// serviceManager manages the lifecycle of the [deviceservice.DeviceService] and implements the
// gRPC interfaces for the [deviceservice.DeviceService] and the ServiceState.
type serviceManager struct {
	dsgrpcpb.UnimplementedDeviceServiceServer
	ssgrpcpb.UnimplementedServiceStateServer

	config        *dscpb.DeviceServiceConfig
	daClient      dagrpcpb.DataAssetsClient
	deviceService *deviceservice.DeviceService

	mu    sync.Mutex // Protects state
	state *sspb.SelfState
}

// newServiceManager creates a new serviceManager.
func newServiceManager(ctx context.Context, config *dscpb.DeviceServiceConfig, daClient dagrpcpb.DataAssetsClient) *serviceManager {
	if config == nil || daClient == nil {
		state := &sspb.SelfState{
			StateCode: sspb.SelfState_STATE_CODE_ERROR,
			ExtendedStatus: &espb.ExtendedStatus{
				Severity: espb.ExtendedStatus_ERROR,
				Title:    "Invalid service configuration: Config and data assets client must not be nil",
			},
		}
		return &serviceManager{
			config:   config,
			daClient: daClient,
			state:    state,
		}
	}

	ds, err := deviceservice.NewDeviceService(ctx, config, daClient)
	if err != nil {
		log.ErrorContextf(ctx, "Failed to create device service: %v", err)
		state := &sspb.SelfState{
			StateCode: sspb.SelfState_STATE_CODE_ERROR,
			ExtendedStatus: &espb.ExtendedStatus{
				Severity: espb.ExtendedStatus_ERROR,
				Title:    "Failed to create device service: " + err.Error(),
			},
		}
		return &serviceManager{
			config:   config,
			daClient: daClient,
			state:    state,
		}
	}
	return &serviceManager{
		config:        config,
		daClient:      daClient,
		deviceService: ds,
		state:         &sspb.SelfState{StateCode: sspb.SelfState_STATE_CODE_ENABLED},
	}
}

// GetConfiguration proxies the request to the underlying DeviceService.
func (m *serviceManager) GetConfiguration(ctx context.Context, req *dspb.GetConfigurationRequest) (*dspb.GetConfigurationResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.InfoContextf(ctx, "GetConfiguration called when in state %v", m.state.StateCode)

	if m.state.StateCode != sspb.SelfState_STATE_CODE_ENABLED || m.deviceService == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "service is not enabled")
	}
	return m.deviceService.GetConfiguration(ctx, req)
}

// GetState returns the current state of the service.
func (m *serviceManager) GetState(ctx context.Context, req *sspb.GetStateRequest) (*sspb.SelfState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Return a copy to prevent external modification of the internal state.
	return proto.Clone(m.state).(*sspb.SelfState), nil
}

// Enable creates and enables the DeviceService.
func (m *serviceManager) Enable(ctx context.Context, req *sspb.EnableRequest) (*sspb.EnableResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.InfoContextf(ctx, "Enable called when in state %v", m.state.StateCode)

	if m.state.StateCode == sspb.SelfState_STATE_CODE_ENABLED {
		return &sspb.EnableResponse{}, nil
	}

	ds, err := deviceservice.NewDeviceService(ctx, m.config, m.daClient)
	if err != nil {
		m.state = &sspb.SelfState{
			StateCode: sspb.SelfState_STATE_CODE_ERROR,
			ExtendedStatus: &espb.ExtendedStatus{
				Severity: espb.ExtendedStatus_ERROR,
				Title:    "Failed to create device service: " + err.Error(),
			},
		}
		return nil, fmt.Errorf("failed to create device service: %w", err)
	}

	m.deviceService = ds
	m.state = &sspb.SelfState{StateCode: sspb.SelfState_STATE_CODE_ENABLED}
	return &sspb.EnableResponse{}, nil
}

// disables the [serviceManager] and erases the [deviceservice.DeviceService] contained in it.
func (m *serviceManager) Disable(ctx context.Context, req *sspb.DisableRequest) (*sspb.DisableResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.InfoContextf(ctx, "Disable called when in state %v", m.state.StateCode)

	m.deviceService = nil
	m.state = &sspb.SelfState{StateCode: sspb.SelfState_STATE_CODE_DISABLED}
	return &sspb.DisableResponse{}, nil
}
