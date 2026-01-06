// Copyright 2023 Intrinsic Innovation LLC

// Package servicegen implements creation of a Service Asset bundle.
package servicegen

import (
	"fmt"
	"strings"

	"intrinsic/assets/services/servicebundle"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/sourcecodeinfoview"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// CreateServiceBundleOptions provides the data needed to create a Service Asset bundle.
type CreateServiceBundleOptions struct {
	DefaultConfigPath      string
	FileDescriptorSetPaths []string
	ImageTarPaths          []string
	ManifestPath           string
	OutputBundlePath       string
}

func pruneSourceCodeInfo(defaultConfig *anypb.Any, fds *dpb.FileDescriptorSet) error {
	if fds == nil {
		return nil
	}

	var fullNames []string
	if defaultConfig != nil {
		typeURLParts := strings.Split(defaultConfig.GetTypeUrl(), "/")
		if len(typeURLParts) < 1 {
			return fmt.Errorf("failed to extract default proto name from type URL: %v", defaultConfig.GetTypeUrl())
		}
		fullNames = append(fullNames, typeURLParts[len(typeURLParts)-1])
	}

	// Note that a nil default config will cause all source code info fields to be
	// stripped out.
	sourcecodeinfoview.PruneSourceCodeInfo(fullNames, fds)
	return nil
}

// CreateServiceBundle creates a Service Asset bundle on disk.
func CreateServiceBundle(opts *CreateServiceBundleOptions) error {
	m := &smpb.ServiceManifest{}
	if err := protoio.ReadTextProto(opts.ManifestPath, m); err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	fds, err := registryutil.LoadFileDescriptorSets(opts.FileDescriptorSetPaths)
	if err != nil {
		return fmt.Errorf("unable to build FileDescriptorSet: %w", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(fds)
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %w", err)
	}

	var defaultConfig *anypb.Any
	if opts.DefaultConfigPath != "" {
		defaultConfig = &anypb.Any{}
		if err := protoio.ReadTextProto(opts.DefaultConfigPath, defaultConfig, protoio.WithResolver(types)); err != nil {
			return fmt.Errorf("failed to read default config proto: %w", err)
		}
		if err := pruneSourceCodeInfo(defaultConfig, fds); err != nil {
			return fmt.Errorf("failed to prune source code info: %w", err)
		}
	}

	if err := servicebundle.Write(m, opts.OutputBundlePath,
		servicebundle.WithFileDescriptorSet(fds),
		servicebundle.WithDefaultConfig(defaultConfig),
		servicebundle.WithImageTarPaths(opts.ImageTarPaths),
	); err != nil {
		return fmt.Errorf("failed to write Service Asset bundle: %w", err)
	}

	return nil
}
