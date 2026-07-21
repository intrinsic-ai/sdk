// Copyright 2023 Intrinsic Innovation LLC

package inbuild_data_test

import (
	"testing"

	"intrinsic/assets/data/databundle"
	"intrinsic/assets/data/utils"
	"intrinsic/assets/referenceddata"
	"intrinsic/util/testing/testio"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
)

const (
	mappedBundlePath = "intrinsic/tools/inbuild/integration_tests/inbuild_data_mapped.bundle.tar"
)

func TestInbuildDataBundle_Mapped(t *testing.T) {
	rp := testio.MustCreateRunfilePath(t, mappedBundlePath)

	bundle, err := databundle.ReadFile(t.Context(), rp,
		databundle.WithReferencedDataProcessor(referenceddata.InlineProcessor()),
	)
	if err != nil {
		t.Fatalf("databundle.ReadFile returned unexpected error: %v", err)
	}

	gotPayload, err := utils.ExtractPayload(bundle.Data)
	if err != nil {
		t.Fatalf("utils.ExtractPayload returned unexpected error: %v", err)
	}

	rdStruct := &rdspb.ReferencedDataStruct{}
	if err := bundle.Data.GetData().UnmarshalTo(rdStruct); err != nil {
		t.Fatalf("failed to unmarshal Data to ReferencedDataStruct: %v (extracted: %v)", err, gotPayload)
	}

	wantPayload := &rdspb.ReferencedDataStruct{
		Fields: map[string]*rdspb.Value{
			"item": {
				Kind: &rdspb.Value_StringValue{
					StringValue: "hello mapped",
				},
			},
			"file_ref": {
				Kind: &rdspb.Value_ReferencedDataValue{
					ReferencedDataValue: &rdpb.ReferencedData{
						Data: &rdpb.ReferencedData_Inlined{
							Inlined: []byte("{\n  \"message\": \"Some mapped file content.\"\n}\n"),
						},
					},
				},
			},
		},
	}

	opts := []cmp.Option{
		protocmp.Transform(),
		protocmp.IgnoreFields(&rdpb.ReferencedData{}, "digest"),
	}
	if diff := cmp.Diff(wantPayload, rdStruct, opts...); diff != "" {
		t.Errorf("Data payload mismatch. Got %v want %v. Diff (-want +got):\n%s", rdStruct, wantPayload, diff)
	}
}
