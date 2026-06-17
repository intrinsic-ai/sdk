// Copyright 2023 Intrinsic Innovation LLC

// Package any resolves Any messages on a best-effort basis given only the type_url.
package any

import (
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Resolver combines MessageTypeResolver and ExtensionTypeResolver interfaces.
type Resolver interface {
	protoregistry.MessageTypeResolver
	protoregistry.ExtensionTypeResolver
}
