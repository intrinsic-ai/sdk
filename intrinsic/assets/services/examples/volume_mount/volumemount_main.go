// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// main is the entry point for the volume mount service.
package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"intrinsic/assets/services/config"
	"intrinsic/assets/services/examples/volume_mount/volumemount"
	"intrinsic/production/intrinsic"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

	vmgrpcpb "intrinsic/assets/services/examples/volume_mount/proto/v1/volume_mount_go_proto"
	vmpb "intrinsic/assets/services/examples/volume_mount/proto/v1/volume_mount_go_proto"
)

func main() {
	intrinsic.Init()

	rc, err := config.LoadRuntimeContext()
	if err != nil {
		log.Exitf("Failed to load runtime context: %v", err)
	}
	log.Infof("Service Name: %q", rc.GetName())

	config := &vmpb.VolumeMountConfig{}
	if err := rc.GetConfig().UnmarshalTo(config); err != nil {
		log.Exitf("Failed to unpack config: %v", err)
	}

	address := fmt.Sprintf(":%d", rc.GetPort())
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Exitf("Server failed to listen at %q: %v", address, err)
	}
	log.Infof("Server is now listening at %q", address)

	service := volumemount.NewService(&volumemount.ServiceOptions{
		Config:        config,
		RootMountPath: "/volumes",
	})

	s := grpc.NewServer()
	vmgrpcpb.RegisterVolumeMountServiceServer(s, service)

	if err := service.WriteInitialFiles(context.Background()); err != nil {
		log.Exitf("Failed to write initial files: %v", err)
	}

	response, err := service.ListDir(context.Background(), &vmpb.ListDirRequest{
		Path:      "/",
		Recursive: true,
	})
	if err != nil {
		log.Exitf("Failed to list contents of mounted volumes: %v", err)
	}
	var paths []string
	for _, entry := range response.GetEntries() {
		paths = append(paths, entry.GetPath())
	}
	log.Infof("Initial contents of mounted volumes:\n%v", strings.Join(paths, "\n"))

	s.Serve(lis)
}
