// Copyright 2023 Intrinsic Innovation LLC

// Package kvstore is a wrapper around the C++ kvstore class provided bypubsub.
package kvstore

import (
	"fmt"
	"path"
	"time"

	"intrinsic/platform/pubsub/golang/pubsubinterface"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

// MakeKey creates a key from an arbitrary number of strings, removing leading and
// trailing slashes, and joining them with a slash delimiter.
func MakeKey(parts ...string) string {
	return path.Join(parts...)
}

// KVStore represents an instance of a KVStore object. It provides an
// interface that allows getting and setting key/value pairs.
type KVStore interface {
	// Set sets the value for the given key. A key can't include any of the
	// following characters: /, *, ?, #, [ and ]. It will be namespaced according
	// to the settings provided in the NamespaceConfig.
	Set(key string, value proto.Message, highConsistency bool) error

	// Get returns the value for the given key. Wildcard queries are not supported
	// with this method, use the GetAll method instead. If wildcard values are
	// sent as part of the config, an error is returned. Timeout may be nil.
	Get(key string, timeout *time.Duration) (*anypb.Any, error)

	// GetAll invokes the valueCallback for a given key that matches the expression.
	GetAll(key string, valueCallback func(*anypb.Any), ondoneCallback func(string)) (KVQuery, error)

	// ListAllKeys invokes the keyCallback for a given key that matches the expression.
	ListAllKeys(key string, keyCallback func(string), ondoneCallback func(string)) (KVQuery, error)

	// Delete removes a value from the store for the given key.
	Delete(key string) error

	// Subscribe creates a subscription to changes in the KV store.
	//
	// It assumes that all values whose keys match the given key expression have
	// the same type.
	//
	// Parameters:
	// - keyExpression - key expression to subscribe to.
	// - config - subscription configuration.
	// - exemplar - empty proto of the same type as values in the KV store.
	// - msgCallback - callback that will be called when a value is added or updated.
	// - deletionCallback - callback that will be called when a value is deleted.
	// - errCallback - callback that will be called in case of type mismatch.
	Subscribe(
		keyExpression string, config pubsubinterface.TopicConfig,
		exemplar proto.Message,
		msgCallback func(string, proto.Message),
		deletionCallback func(string),
		errCallback func(string, *anypb.Any, error)) (pubsubinterface.Subscription, error)

	// SubscribeToRawValues creates a subscription to changes in the KV store.
	//
	// It doesn't make any assumptions about types of values that match the given
	// key expression. The calling code is responsible for extracting those values
	// from `anypb.Any` protos and checking their type.
	//
	// Parameters:
	// - keyExpression - key expression to subscribe to.
	// - config - subscription configuration.
	// - msgCallback - callback that will be called when a value is added or updated.
	// - deletionCallback - callback that will be called when a value is deleted.
	SubscribeToRawValues(
		keyExpression string, config pubsubinterface.TopicConfig,
		msgCallback func(string, *anypb.Any),
		deletionCallback func(string)) (pubsubinterface.Subscription, error)

	// GetWorkcellReplicationNamespace returns a string that corresponds to the
	// namespace needed for the workcell's replicated namespace. Using this
	// namespace prefix for keyexprs will make the KVStore use replicated storage.
	//
	// This function may return an error because the namespace may not be available
	// immediately. The calling code should check whether the returned error is
	// kvstore.ErrNotFound, and retry.
	GetWorkcellReplicationNamespace() (string, error)

	// GetGlobalReplicationNamespace returns a string that corresponds to the
	// namespace needed for the global replicated namespace. Using this
	// namespace prefix for keyexprs will make the KVStore use replicated storage
	// available to all workcells within an organization.
	GetGlobalReplicationNamespace() string
}

// KVQuery is a handle for a GetALl KVStore query. Keep the handle alive to
// continue the query.
type KVQuery interface {
	// Close closes out the query
	Close()
}

var (
	// ErrDeadlineExceeded is returned if a deadline is exceeded.
	ErrDeadlineExceeded = fmt.Errorf("deadline exceeded")

	// ErrNotFound is returned if Get does not find any values.
	ErrNotFound = fmt.Errorf("not found")
)
