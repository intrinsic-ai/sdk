// Copyright 2023 Intrinsic Innovation LLC

package bundle

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	apb "google.golang.org/protobuf/types/known/anypb"
	"intrinsic/assets/bundleio"
	_ "intrinsic/tools/inbuild/cmd/service/test_data/example_service_go_proto" // Needed to resolve the proto message
	"intrinsic/util/path_resolver/pathresolver"
)

const (
	exampleManifestPath          = "intrinsic/tools/inbuild/cmd/service/test_data/example_service.manifest.pbtxt"
	exampleFileDescriptorSetPath = "intrinsic/tools/inbuild/cmd/service/test_data/example_service_proto-descriptor-set.proto.bin"
	exampleImagePath             = "intrinsic/tools/inbuild/cmd/service/test_data/example_service_py_image.tar"
	exampleDefaultConfigPath     = "intrinsic/tools/inbuild/cmd/service/test_data/example_service_default_config.textproto"
)

type bundleCheck func(t *testing.T)

func checkManifestHasID(t *testing.T, bundlePath string, wantPackage string, wantName string) bundleCheck {
	return func(t *testing.T) {
		t.Helper()
		manifest, err := bundleio.ReadServiceManifest(bundlePath)
		if err != nil {
			t.Fatalf("bundleio.ReadServiceManifest(%q) failed: %v", bundlePath, err)
		}
		if got := manifest.GetMetadata().GetId().GetPackage(); got != wantPackage {
			t.Errorf("manifest.GetMetadata().GetId().GetPackage() = %q, want %q", got, wantPackage)
		}
		if got := manifest.GetMetadata().GetId().GetName(); got != wantName {
			t.Errorf("manifest.GetMetadata().GetId().GetName() = %q, want %q", got, wantName)
		}
	}
}

func checkBundleHasDefaultConfig(t *testing.T, bundlePath string, wantConfigPath string) bundleCheck {
	return func(t *testing.T) {
		t.Helper()
		wantConfigBytes, err := os.ReadFile(wantConfigPath)
		if err != nil {
			t.Fatalf("os.ReadFile(%q) failed: %v", wantConfigPath, err)
		}
		wantConfig := &apb.Any{}
		if err := prototext.Unmarshal(wantConfigBytes, wantConfig); err != nil {
			t.Fatalf("pprototextroto.Unmarshal(%q) failed: %v", wantConfigPath, err)
		}
		_, gotContents, err := bundleio.ReadService(bundlePath)
		if err != nil {
			t.Fatalf("bundleio.ReadService(%q) failed: %v", bundlePath, err)
		}
		got := &apb.Any{}
		if err := proto.Unmarshal(gotContents["default_config.binarypb"], got); err != nil {
			t.Fatalf("proto.Unmarshal(gotContents[\"default_config.binarypb\"], got) failed: %v", err)
		}

		if !cmp.Equal(got, wantConfig, protocmp.Transform()) {
			t.Errorf("contents[\"default_config.binpb\"] = %v, want %v", got, wantConfig)
		}
	}
}

func mustResolve(t *testing.T, path string) string {
	t.Helper()
	resolved, err := pathresolver.ResolveRunfilesPath(path)
	if err != nil {
		t.Fatalf("pathresolver.ResolveRunfilesPath(%q) failed: %v", path, err)
	}
	return resolved
}

func TestBundleCreate(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		checks    []bundleCheck
		wantError string
	}{
		{
			name: "create service",
			args: []string{
				"--manifest", mustResolve(t, exampleManifestPath),
				"--file_descriptor_set", mustResolve(t, exampleFileDescriptorSetPath),
				"--oci_image", mustResolve(t, exampleImagePath),
			},
			checks: []bundleCheck{
				checkManifestHasID(t, "service.bundle.tar", "com.example", "my_service"),
			},
		},
		{
			name: "create service output elsewhere",
			args: []string{
				"--manifest", mustResolve(t, exampleManifestPath),
				"--file_descriptor_set", mustResolve(t, exampleFileDescriptorSetPath),
				"--oci_image", mustResolve(t, exampleImagePath),
				"--output", "my_service.bundle.tar",
			},
			checks: []bundleCheck{
				checkManifestHasID(t, "my_service.bundle.tar", "com.example", "my_service"),
			},
		},
		{
			name: "create with default config",
			args: []string{
				"--manifest", mustResolve(t, exampleManifestPath),
				"--file_descriptor_set", mustResolve(t, exampleFileDescriptorSetPath),
				"--oci_image", mustResolve(t, exampleImagePath),
				"--default_config", mustResolve(t, exampleDefaultConfigPath),
			},
			checks: []bundleCheck{
				checkManifestHasID(t, "service.bundle.tar", "com.example", "my_service"),
				checkBundleHasDefaultConfig(t, "service.bundle.tar", mustResolve(t, exampleDefaultConfigPath)),
			},
		},
		{
			name: "create with no config",
			args: []string{
				"--manifest", mustResolve(t, exampleManifestPath),
				"--oci_image", mustResolve(t, exampleImagePath),
			},
			checks: []bundleCheck{
				checkManifestHasID(t, "service.bundle.tar", "com.example", "my_service"),
			},
		},
		{
			name: "create with config but no file descriptor sets",
			args: []string{
				"--manifest", mustResolve(t, exampleManifestPath),
				"--oci_image", mustResolve(t, exampleImagePath),
				"--default_config", mustResolve(t, exampleDefaultConfigPath),
			},
			wantError: "--file_descriptor_set is required when --default_config is used",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure we're in a writable directory.
			os.Chdir(t.TempDir())
			// Prevent state leaking between tests.
			resetBundleCommand()

			BundleCmd.SetArgs(tt.args)

			err := BundleCmd.Execute()
			if tt.wantError != "" {
				if err == nil {
					t.Fatalf("Want BundleCmd.Execute() to err; got no err")
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("Want error to contain %v, got %v", tt.wantError, err)
				}
			}
			if tt.wantError == "" && err != nil {
				t.Fatalf("BundleCmd.Execute() failed: %v", err)
			}

			for _, check := range tt.checks {
				check(t)
			}
		})
	}
}
