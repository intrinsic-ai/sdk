// Copyright 2023 Intrinsic Innovation LLC

package protoregistryclient

import (
	"context"
	"fmt"
	"strings"

	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/typeurl"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	protoregistrypb "intrinsic/proto_tools/proto/proto_registry_go_proto"
)

// Resolver is the interface required to customize the marshalling and
// unmarshalling of nested Any protos to text-format and JSON.
type Resolver interface {
	protoregistry.ExtensionTypeResolver
	protoregistry.MessageTypeResolver
}

// ProtoRegistryResolver is a resolver which uses the proto registry to resolve
// type URLs starting with "type.intrinsic.ai/". All other type URLs will be
// resolved using one or more given default resolvers.
type ProtoRegistryResolver struct {
	// Proto registry client for resolving Intrinsic type URLs.
	ctx           context.Context
	protoRegistry protoregistrypb.ProtoRegistryClient

	// Resolvers for non-Intrinsic type URLs (tried in order).
	defaultResolvers []Resolver

	// A cache by type URL to avoid redundant queries to the proto registry.
	typesCache map[string]*protoregistry.Types
}

// NewProtoRegistryResolver creates a new ProtoRegistryResolver which uses the
// given context and proto registry client to resolve type URLs starting with
// "type.intrinsic.ai/". All other type URLs will be resolved using the given
// default resolvers (tried in the given order).
//
// Example choices for 'defaultResolvers':
//   - "[]Resolver{}": All lookups for non-Intrinsic type URLs will fail.
//   - "[]Resolver{protoregistry.GlobalTypes}": All non-Intrinsic type URLs will
//     be looked up using the default/compiled-in proto pool.
//
// The given context must have all the required metadata (org headers etc.)
// for sending requests to the proto registry.
//
// The returned resolver uses a local cache that prevents redundant queries to
// the proto registry during its lifetime. The returned resolver is NOT
// thread-safe and should not be used by multiple go-routines simultaneously.
func NewProtoRegistryResolver(ctx context.Context, protoRegistry protoregistrypb.ProtoRegistryClient, defaultResolvers []Resolver) *ProtoRegistryResolver {
	return &ProtoRegistryResolver{
		ctx:              ctx,
		protoRegistry:    protoRegistry,
		defaultResolvers: defaultResolvers,
		typesCache:       make(map[string]*protoregistry.Types),
	}
}

func (r *ProtoRegistryResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	return nil, protoregistry.NotFound
}

func (r *ProtoRegistryResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	return nil, protoregistry.NotFound
}

func (r *ProtoRegistryResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return nil, protoregistry.NotFound
}

func (r *ProtoRegistryResolver) FindMessageByURL(typeURL string) (protoreflect.MessageType, error) {
	if !strings.HasPrefix(typeURL, typeurl.IntrinsicPrefix) {
		for _, resolver := range r.defaultResolvers {
			msgType, err := resolver.FindMessageByURL(typeURL)
			if err == protoregistry.NotFound {
				continue
			} else if err != nil {
				return nil, fmt.Errorf("default resolver lookup failed: %w", err)
			}
			return msgType, nil
		}
		return nil, protoregistry.NotFound
	}

	// Handle type URL starting with type.intrinsic.ai/
	types, ok := r.typesCache[typeURL]

	if !ok {
		// Cache miss, query proto registry
		req := &protoregistrypb.GetNamedFileDescriptorSetRequest{
			IdentifierType: &protoregistrypb.GetNamedFileDescriptorSetRequest_TypeUrl{
				TypeUrl: typeURL,
			},
		}

		result, err := r.protoRegistry.GetNamedFileDescriptorSet(r.ctx, req)
		if err != nil {
			return nil, fmt.Errorf("request to proto registry failed: %w", err)
		}

		types, err = registryutil.NewTypesFromFileDescriptorSet(result.GetFileDescriptorSet())
		if err != nil {
			return nil, fmt.Errorf("failed creating types from file descriptor set: %w", err)
		}

		r.typesCache[typeURL] = types
	}

	return types.FindMessageByURL(typeURL)
}
