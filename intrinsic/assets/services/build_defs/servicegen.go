// Copyright 2023 Intrinsic Innovation LLC

// Package servicegen implements creation of the service type bundle.
package servicegen

import (
	"fmt"
	"strings"

	"intrinsic/assets/bundleio"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/sourcecodeinfoview"

	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// ServiceData holds the data needed to create a service bundle.
type ServiceData struct {
	// Optional path to default config proto.
	DefaultConfig string
	// Paths to binary file descriptor set protos associated with the manifest.
	FileDescriptorSets []string
	// Paths to tar archives for images.
	ImageTars []string
	// The deserialized ServiceManifest.
	Manifest *smpb.ServiceManifest
	// Bundle tar path.
	OutputBundle string
}

func pruneSourceCodeInfo(defaultConfig *anypb.Any, fds *dpb.FileDescriptorSet) error {
	if fds == nil {
		return nil
	}

	var fullNames []string
	if defaultConfig != nil {
		typeURLParts := strings.Split(defaultConfig.GetTypeUrl(), "/")
		if len(typeURLParts) < 1 {
			return fmt.Errorf("cannot extract default proto name from type URL: %v", defaultConfig.GetTypeUrl())
		}
		fullNames = append(fullNames, typeURLParts[len(typeURLParts)-1])
	}

	// Note that a nil default config will cause all source code info fields to be
	// stripped out.
	sourcecodeinfoview.PruneSourceCodeInfo(fullNames, fds)
	return nil
}

// CreateService bundles the data needed for software services.
func CreateService(d *ServiceData) error {
	set, err := registryutil.LoadFileDescriptorSets(d.FileDescriptorSets)
	if err != nil {
		return fmt.Errorf("unable to build FileDescriptorSet: %v", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(set)
	if err != nil {
		return fmt.Errorf("failed to populate the registry: %v", err)
	}

	var defaultConfig *anypb.Any
	if d.DefaultConfig != "" {
		defaultConfig = &anypb.Any{}
		if err := protoio.ReadTextProto(d.DefaultConfig, defaultConfig, protoio.WithResolver(types)); err != nil {
			return fmt.Errorf("failed to read default config proto: %v", err)
		}
		if err := pruneSourceCodeInfo(defaultConfig, set); err != nil {
			return fmt.Errorf("unable to process source code info: %v", err)
		}
	}

	if err := bundleio.WriteService(d.OutputBundle, bundleio.WriteServiceOpts{
		Manifest:      d.Manifest,
		Descriptors:   set,
		DefaultConfig: defaultConfig,
		ImageTars:     d.ImageTars,
	}); err != nil {
		return fmt.Errorf("unable to write service bundle: %v", err)
	}

	return nil
}
