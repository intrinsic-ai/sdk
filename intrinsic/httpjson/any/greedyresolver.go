// Copyright 2023 Intrinsic Innovation LLC

// Package greedyresolver resolves protobuf names and urls to MessageType using a greedy algorithm.
package greedyresolver

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// GreedyResolver implements protoregistry.MessageTypeResolver, returning the first result found from its list of resolvers.
type GreedyResolver struct {
	resolvers []protoregistry.MessageTypeResolver
}

// NewGreedyResolver returns a new GreedyResolver instance.
func NewGreedyResolver(resolvers []protoregistry.MessageTypeResolver) *GreedyResolver {
	return &GreedyResolver{
		resolvers: resolvers,
	}
}

// FindMessageByName looks up a message by its full name.
// E.g., "google.protobuf.Any"
//
// This returns (nil, NotFound) if not found.
func (r *GreedyResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	for _, resolver := range r.resolvers {
		msg, err := resolver.FindMessageByName(message)
		if err == nil {
			return msg, nil
		}
	}
	return nil, protoregistry.NotFound
}

// FindMessageByURL looks up a message by a URL identifier.
// See documentation on google.protobuf.Any.type_url for the URL format.
//
// This returns (nil, NotFound) if not found.
func (r *GreedyResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	for _, resolver := range r.resolvers {
		msg, err := resolver.FindMessageByURL(url)
		if err == nil {
			return msg, nil
		}
	}
	return nil, protoregistry.NotFound
}
