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
// type URLs. Type URLs that fail to resolve from the proto registry may be
// resolved using fallback resolvers.
type ProtoRegistryResolver struct {
	// Proto registry client for resolving Intrinsic type URLs.
	ctx           context.Context
	protoRegistry protoregistrypb.ProtoRegistryClient

	// Resolvers to try in case resolution via the proto registry fails. Fallback
	// resolvers will be tried in order.
	fallbackResolvers []Resolver

	// A cache by type URL to avoid redundant queries to the proto registry.
	typesCache map[string]*protoregistry.Types
}

// NewProtoRegistryResolver creates a new [ProtoRegistryResolver] which uses the
// given context and proto registry client to resolve type URLs. Type URLs that
// cannot be resolved using the proto registry will be resolved using the given
// fallback resolvers (tried in the given order).
//
// Example choices for `fallbackResolvers“:
//   - [][Resolver]{}: All lookups for type URLs that are not available in the
//     proto registry will fail.
//   - [][Resolver]{protoregistry.GlobalTypes}: All type URLs that are not
//     available in the proto registry will be looked up using the
//     default/compiled-in proto pool.
//
// The given context must have all the required metadata (org headers etc.)
// for sending requests to the proto registry.
//
// The returned resolver uses a local cache that prevents redundant queries to
// the proto registry during its lifetime. The returned resolver is NOT
// thread-safe and should not be used by multiple go-routines simultaneously.
func NewProtoRegistryResolver(ctx context.Context, protoRegistry protoregistrypb.ProtoRegistryClient, fallbackResolvers []Resolver) *ProtoRegistryResolver {
	return &ProtoRegistryResolver{
		ctx:               ctx,
		protoRegistry:     protoRegistry,
		fallbackResolvers: fallbackResolvers,
		typesCache:        make(map[string]*protoregistry.Types),
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

func (r *ProtoRegistryResolver) resolveViaProtoRegistry(typeURL string) (protoreflect.MessageType, error) {
	types, ok := r.typesCache[typeURL]
	if !ok {
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

func (r *ProtoRegistryResolver) FindMessageByURL(typeURL string) (protoreflect.MessageType, error) {
	if strings.HasPrefix(typeURL, typeurl.IntrinsicPrefix) {
		return r.resolveViaProtoRegistry(typeURL)
	}

	if strings.HasPrefix(typeURL, typeurl.DefaultPrefix) {
		// Explicitly enable fallback resolution for the default prefix. To make
		// this possible, only return here if there is NO error. If there is any
		// error, fall through to fallback resolvers.
		if msgType, err := r.resolveViaProtoRegistry(typeURL); err == nil { // if NO error
			return msgType, nil
		}
	}

	for _, resolver := range r.fallbackResolvers {
		msgType, err := resolver.FindMessageByURL(typeURL)
		if err == protoregistry.NotFound {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("fallback resolver lookup failed: %w", err)
		}
		return msgType, nil
	}

	return nil, protoregistry.NotFound
}
