// Copyright 2023 Intrinsic Innovation LLC

// Package protoregistryresolver resolves protobuf names and urls to MessageType using the ProtoRegsitry service.
package protoregistryresolver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	protoregistrygrpcpb "intrinsic/proto_tools/proto/proto_registry_go_proto"
	"intrinsic/util/proto/registryutil"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// ProtoRegistry service only resolves types that start with this URL
const ProtoRegistryTypePrefix = "type.intrinsic.ai/"

// ProtoRegistryResolver implements protoregistry.MessageTypeResolver, returning types retrieved from the ProtoRegistry service.
type ProtoRegistryResolver struct {
	client    protoregistrygrpcpb.ProtoRegistryClient
	typeCache map[string]*protoregistry.Types
}

func (r *ProtoRegistryResolver) descriptorFromProtoRegistry(url string) (*descriptorpb.FileDescriptorSet, error) {
	if !strings.HasPrefix(url, ProtoRegistryTypePrefix) {
		return nil, protoregistry.NotFound
	}
	req := &protoregistrygrpcpb.GetNamedFileDescriptorSetRequest{
		IdentifierType: &protoregistrygrpcpb.GetNamedFileDescriptorSetRequest_TypeUrl{
			TypeUrl: string(url),
		},
	}
	result, err := r.client.GetNamedFileDescriptorSet(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		switch st.Code() {
		case codes.NotFound:
			return nil, protoregistry.NotFound
		default:
			return nil, err
		}
	}
	return result.FileDescriptorSet, nil
}

func (r *ProtoRegistryResolver) addDescriptorToCache(url string, fileDescriptorSet *descriptorpb.FileDescriptorSet) (*protoregistry.Types, error) {
	files := new(protoregistry.Files)
	types := new(protoregistry.Types)

	for _, fdProto := range fileDescriptorSet.GetFile() {
		// NewFile checks dependencies. We pass 'files' as the resolver so it can
		// resolve dependencies that were processed in previous iterations of this loop.
		file, err := protodesc.NewFile(fdProto, files)
		if err != nil {
			return nil, fmt.Errorf("failed to create file descriptor: %w", err)
		}

		err = files.RegisterFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to register file: %w", err)
		}
	}

	if err := registryutil.PopulateTypesFromFiles(types, files); err != nil {
		return nil, fmt.Errorf("failed to populate types: %w", err)
	}

	// Insert into cache (No lock)
	r.typeCache[url] = types

	return types, nil
}

// NewProtoRegistryResolver returns a new ProtoRegistryResolver instance.
func NewProtoRegistryResolver(protoRegistryAddress string) (*ProtoRegistryResolver, error) {
	if protoRegistryAddress == "" {
		return nil, errors.New("protoRegistryAddress must not be empty")
	}
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	grpcClient, err := grpc.NewClient(protoRegistryAddress, options...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient(%q) failed: %w", protoRegistryAddress, err)
	}

	return &ProtoRegistryResolver{
		client:    protoregistrygrpcpb.NewProtoRegistryClient(grpcClient),
		typeCache: make(map[string]*protoregistry.Types),
	}, nil
}

// FindMessageByName looks up a message by its full name.
// E.g., "google.protobuf.Any"
//
// This always returns (nil, NotFound) because the ProtoRegistry services only supports types with a specific URL prefix.
func (r *ProtoRegistryResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return nil, protoregistry.NotFound
}

// FindMessageByURL looks up a message by a URL identifier.
// See documentation on google.protobuf.Any.type_url for the URL format.
//
// This returns (nil, NotFound) if not found.
func (r *ProtoRegistryResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	// Check Cache (No lock)
	if cachedType, found := r.typeCache[url]; found {
		return cachedType.FindMessageByURL(url)
	}

	// Not in cache, fetch from registry
	fds, err := r.descriptorFromProtoRegistry(url)
	if err != nil {
		return nil, err
	}

	// Add to cache
	cachedType, err := r.addDescriptorToCache(url, fds)
	if err != nil {
		return nil, err
	}

	return cachedType.FindMessageByURL(url)
}
