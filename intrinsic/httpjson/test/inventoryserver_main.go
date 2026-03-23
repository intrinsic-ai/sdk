// Copyright 2023 Intrinsic Innovation LLC

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"intrinsic/httpjson/test/inventoryserver"
	"intrinsic/util/proto/protoio"

	"google.golang.org/grpc"

	inventorypb "intrinsic/httpjson/test/inventory_service_go_proto"
	rcpb "intrinsic/resources/proto/runtime_context_go_proto"
)

const (
	defaultRuntimeContextPath = "/etc/intrinsic/runtime_config.pb"
	extensionTextProto        = ".textproto"
	extensionTxtPb            = ".txtpb"
)

var runtimeContextPath = flag.String("runtime_context_path", defaultRuntimeContextPath, "Path to the runtime context binary proto.")

func main() {
	flag.Parse()

	rc := new(rcpb.RuntimeContext)

	// Determine if we should parse as text or binary proto
	isTextProto := strings.HasSuffix(*runtimeContextPath, extensionTxtPb) ||
		strings.HasSuffix(*runtimeContextPath, extensionTextProto)

	if isTextProto {
		if err := protoio.ReadTextProto(*runtimeContextPath, rc); err != nil {
			log.Fatalf("Failed to read runtime context text proto: %v", err)
		}
	} else {
		if err := protoio.ReadBinaryProto(*runtimeContextPath, rc); err != nil {
			log.Fatalf("Failed to read runtime context binary proto: %v", err)
		}
	}

	address := fmt.Sprintf(":%d", rc.GetPort())
	log.Printf("Listening on address %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", rc.GetPort(), err)
	}

	grpcServer := grpc.NewServer()
	invServer := inventoryserver.NewInventoryServer()

	inventorypb.RegisterInventoryServiceServer(grpcServer, invServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
