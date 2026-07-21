// Copyright 2023 Intrinsic Innovation LLC

package bundle

import (
	"os"
	"path/filepath"
	"testing"

	"intrinsic/assets/data/databundle"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
	"intrinsic/assets/referenceddata"
	cdpb "intrinsic/tools/inbuild/cmd/data/test_data/custom_data_go_proto"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/testing/testio"

	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

func TestBundleCommand_Validation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing manifest",
			args: []string{},
		},
		{
			name: "invalid replace_with_external_reference",
			args: []string{"--manifest=foo.textproto", "--replace_with_external_reference=bad_entry"},
		},
		{
			name: "invalid reference_to_path",
			args: []string{"--manifest=foo.textproto", "--reference_to_path=bad_entry"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetBundleCommand()
			BundleCmd.SetArgs(tc.args)
			if err := BundleCmd.Execute(); err == nil {
				t.Errorf("BundleCmd.Execute() returned nil, expected error")
			}
		})
	}
}

func mustWriteTextproto(t *testing.T, path string, msg proto.Message) {
	t.Helper()
	if err := protoio.WriteStableTextProto(path, msg); err != nil {
		t.Fatalf("failed to write textproto %q: %v", path, err)
	}
}

func TestBundleCommand_Execution(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "dist")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create dist dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "app.js"), []byte("console.log('test');"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	manifestPath := filepath.Join(tempDir, "manifest.textproto")
	manifest := &dmpb.DataManifest{
		Metadata: &dmpb.DataManifest_Metadata{
			Id: &idpb.Id{
				Package: "ai.intrinsic",
				Name:    "test_bundle",
			},
			DisplayName: "Test Bundle",
			Vendor: &vpb.Vendor{
				DisplayName: "Test Vendor",
			},
		},
	}
	mustWriteTextproto(t, manifestPath, manifest)

	rdStruct := &rdspb.ReferencedDataStruct{
		Fields: map[string]*rdspb.Value{
			"app.js": {
				Kind: &rdspb.Value_ReferencedDataValue{
					ReferencedDataValue: &rdpb.ReferencedData{
						Data: &rdpb.ReferencedData_Reference{
							Reference: "app.js",
						},
					},
				},
			},
		},
	}
	dataAny, err := anypb.New(rdStruct)
	if err != nil {
		t.Fatalf("anypb.New failed: %v", err)
	}
	manifestMapped := &dmpb.DataManifest{
		Metadata: manifest.Metadata,
		Data:     dataAny,
	}
	manifestMappedPath := filepath.Join(tempDir, "manifest_mapped.textproto")
	mustWriteTextproto(t, manifestMappedPath, manifestMapped)

	if err := os.WriteFile(filepath.Join(tempDir, "app.js"), []byte("console.log('test');"), 0644); err != nil {
		t.Fatalf("failed to write app.js in tempDir: %v", err)
	}

	customPayload := &cdpb.CustomDataPayload{
		MyData: &rdpb.ReferencedData{
			Data: &rdpb.ReferencedData_Reference{
				Reference: "app.js",
			},
		},
	}
	customDataAny, err := anypb.New(customPayload)
	if err != nil {
		t.Fatalf("anypb.New failed: %v", err)
	}
	customManifest := &dmpb.DataManifest{
		Metadata: manifest.Metadata,
		Data:     customDataAny,
	}
	customManifestPath := filepath.Join(tempDir, "custom_manifest.textproto")
	mustWriteTextproto(t, customManifestPath, customManifest)
	customFdsPath := testio.MustCreateRunfilePath(t, "intrinsic/tools/inbuild/cmd/data/test_data/custom_data_fds_transitive_set_sci.proto.bin")

	wantRdStruct := &rdspb.ReferencedDataStruct{
		Fields: map[string]*rdspb.Value{
			"app.js": {
				Kind: &rdspb.Value_ReferencedDataValue{
					ReferencedDataValue: &rdpb.ReferencedData{
						Data: &rdpb.ReferencedData_Inlined{
							Inlined: []byte("console.log('test');"),
						},
					},
				},
			},
		},
	}

	wantCustomPayload := &cdpb.CustomDataPayload{
		MyData: &rdpb.ReferencedData{
			Data: &rdpb.ReferencedData_Inlined{
				Inlined: []byte("console.log('test');"),
			},
		},
	}

	tests := []struct {
		name        string
		args        []string
		outputFile  string
		wantPayload proto.Message
	}{
		{
			name: "reference mapping (ReferencedDataStruct without fds)",
			args: []string{
				"--manifest=" + manifestMappedPath,
				"--reference_to_path=app.js=" + filepath.Join(sourceDir, "app.js"),
				"--output=" + filepath.Join(tempDir, "mapped.bundle.tar"),
			},
			outputFile:  filepath.Join(tempDir, "mapped.bundle.tar"),
			wantPayload: wantRdStruct,
		},
		{
			name: "automatic reference mapping (ReferencedDataStruct without fds)",
			args: []string{
				"--manifest=" + manifestMappedPath,
				"--output=" + filepath.Join(tempDir, "auto_mapped.bundle.tar"),
			},
			outputFile:  filepath.Join(tempDir, "auto_mapped.bundle.tar"),
			wantPayload: wantRdStruct,
		},
		{
			name: "custom payload with ReferencedData field (with fds)",
			args: []string{
				"--manifest=" + customManifestPath,
				"--file_descriptor_set=" + customFdsPath,
				"--output=" + filepath.Join(tempDir, "custom.bundle.tar"),
			},
			outputFile:  filepath.Join(tempDir, "custom.bundle.tar"),
			wantPayload: wantCustomPayload,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetBundleCommand()
			BundleCmd.SetArgs(tc.args)
			if err := BundleCmd.Execute(); err != nil {
				t.Fatalf("BundleCmd.Execute() returned unexpected error: %v", err)
			}

			bundle, err := databundle.ReadFile(t.Context(), tc.outputFile,
				databundle.WithReferencedDataProcessor(referenceddata.InlineProcessor()),
			)
			if err != nil {
				t.Fatalf("databundle.ReadFile returned unexpected error: %v", err)
			}

			gotPayload := tc.wantPayload.ProtoReflect().New().Interface()
			if err := bundle.Data.GetData().UnmarshalTo(gotPayload); err != nil {
				t.Fatalf("failed to unmarshal Data payload: %v", err)
			}
			if !proto.Equal(gotPayload, tc.wantPayload) {
				t.Errorf("got payload:\n%v\nwant:\n%v", gotPayload, tc.wantPayload)
			}
		})
	}
}
