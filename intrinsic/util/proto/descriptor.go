// Copyright 2023 Intrinsic Innovation LLC

// Package descriptor provides utilities for working with proto descriptors.
package descriptor

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// mergeFileDescriptorSetsOptions contains options for MergeFileDescriptorSets.
type mergeFileDescriptorSetsOptions struct {
	keys []string
}

// MergeFileDescriptorSetsOption is a functional option for MergeFileDescriptorSets.
type MergeFileDescriptorSetsOption func(*mergeFileDescriptorSetsOptions)

// WithKeys is an option for MergeFileDescriptorSets that specifies the keys of the assets to
// merge. These keys are used to provide context in error messages if duplicate FileDescriptorProtos
// are found.
//
// For example, when merging FileDescriptorSets for multiple assets, each asset's id can be used
// as the key.
func WithKeys(keys []string) MergeFileDescriptorSetsOption {
	return func(opts *mergeFileDescriptorSetsOptions) {
		opts.keys = keys
	}
}

func defaultMergedFileDescriptorSetsKeys(n int) []string {
	keys := make([]string, n)
	for i := range keys {
		keys[i] = fmt.Sprintf("FileDescriptorSet %d", i)
	}
	return keys
}

type fileDescriptorAndKey struct {
	f   *dpb.FileDescriptorProto
	key string
}

// MergeFileDescriptorSets merges a list of FileDescriptorSets into a single FileDescriptorSet.
//
// It returns an error if any duplicate FileDescriptorProtos are not identical.
func MergeFileDescriptorSets(fdss []*dpb.FileDescriptorSet, options ...MergeFileDescriptorSetsOption) (*dpb.FileDescriptorSet, error) {
	opts := mergeFileDescriptorSetsOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	if len(opts.keys) == 0 {
		opts.keys = defaultMergedFileDescriptorSetsKeys(len(fdss))
	}
	if len(opts.keys) != len(fdss) {
		return nil, fmt.Errorf("len(keys) (%d) != len(fdss) (%d)", len(opts.keys), len(fdss))
	}

	// Merge them into a single FileDescriptorSet.
	merged := &dpb.FileDescriptorSet{}
	seen := map[string]fileDescriptorAndKey{}
	for i, fds := range fdss {
		for _, f := range fds.GetFile() {
			if existing, ok := seen[f.GetName()]; !ok {
				merged.File = append(merged.File, f)
				// seen maps FileDescriptorProto.name to FileDescriptorProto and the key that referenced it.
				seen[f.GetName()] = fileDescriptorAndKey{f: f, key: opts.keys[i]}
			} else if !proto.Equal(existing.f, f) {
				// We currently require that any duplicate FileDescriptorProtos are identical. We could
				// potentially relax this constraint later by determining which of multiple compatible
				// FileDescriptorProtos we find is the "latest" one and using that.
				return nil, fmt.Errorf("duplicate FileDescriptorProto %q with different contents for %q and %q", f.GetName(), existing.key, opts.keys[i])
			}
		}
	}

	return merged, nil
}

// FileDescriptorSetFrom constructs a FileDescriptorSet for the specified message.
func FileDescriptorSetFrom(msg proto.Message) *dpb.FileDescriptorSet {
	return &dpb.FileDescriptorSet{
		File: collectDependencies(msg.ProtoReflect().Descriptor().ParentFile(), make(map[string]bool)),
	}
}

// collectDependencies collects all dependencies of the specified file descriptor.
//
// Only returns the dependencies that are not already in the seen map. Initial callers should pass
// in an empty map.
func collectDependencies(desc protoreflect.FileDescriptor, seen map[string]bool) []*dpb.FileDescriptorProto {
	if desc == nil || seen[desc.Path()] {
		return []*dpb.FileDescriptorProto{}
	}
	seen[desc.Path()] = true
	result := []*dpb.FileDescriptorProto{
		protodesc.ToFileDescriptorProto(desc),
	}

	fileImports := desc.Imports()
	for i := 0; i < fileImports.Len(); i++ {
		for _, dep := range collectDependencies(fileImports.Get(i), seen) {
			result = append(result, dep)
			seen[dep.GetName()] = true
		}
	}

	return result
}
