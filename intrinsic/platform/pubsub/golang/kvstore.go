// Copyright 2023 Intrinsic Innovation LLC

// Package kvstore is a wrapper around the C++ kvstore class provided bypubsub.
package kvstore

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

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

	// Delete removes a value from the store for the given key.
	Delete(key string) error
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
