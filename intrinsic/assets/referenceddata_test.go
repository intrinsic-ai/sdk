// Copyright 2023 Intrinsic Innovation LLC

package referenceddata

import (
	"fmt"
	"strings"
	"testing"

	"intrinsic/util/proto/descriptor"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
	"intrinsic/assets/data/utils"

	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestReferencedData(t *testing.T) {
	tests := []struct {
		desc               string
		ref                *ReferencedData
		wantModified       bool
		wantReferencedData *rdpb.ReferencedData
		wantSourceProject  string
		wantType           ReferenceType
	}{
		{
			desc:               "empty",
			ref:                FromProto(&rdpb.ReferencedData{}),
			wantModified:       false,
			wantReferencedData: &rdpb.ReferencedData{},
			wantType:           InlinedReferenceType,
		},
		{
			desc: "cas reference",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "intcas://foo",
				},
				SourceProject: proto.String("bar"),
			}),
			wantModified: false,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "intcas://foo",
				},
				SourceProject: proto.String("bar"),
			},
			wantSourceProject: "bar",
			wantType:          CASReferenceType,
		},
		{
			desc: "file reference",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "file:///foo",
				},
			}),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "/foo",
				},
			},
			wantType: FileReferenceType,
		},
		{
			desc: "file reference without protocol prefix",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
			}),
			wantModified: false,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
			},
			wantType: FileReferenceType,
		},
		{
			desc:               "set base dir empty",
			ref:                FromProto(&rdpb.ReferencedData{}).SetBaseDir("/"),
			wantModified:       false,
			wantReferencedData: &rdpb.ReferencedData{},
			wantType:           InlinedReferenceType,
		},
		{
			desc: "set base dir relative file",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
			}).SetBaseDir("/"),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "/foo",
				},
			},
			wantType: FileReferenceType,
		},
		{
			desc: "set base dir absolute file",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "/foo",
				},
			}).SetBaseDir("/"),
			wantModified: false,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "/foo",
				},
			},
			wantType: FileReferenceType,
		},
		{
			desc:         "set digest",
			ref:          FromProto(&rdpb.ReferencedData{}).SetDigest(fmt.Sprintf("%s:foo", utils.HighwayHash128)),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Digest: fmt.Sprintf("%s:foo", utils.HighwayHash128),
			},
			wantType: InlinedReferenceType,
		},
		{
			desc:         "set source project",
			ref:          FromProto(&rdpb.ReferencedData{}).SetSourceProject("foo"),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				SourceProject: proto.String("foo"),
			},
			wantType:          InlinedReferenceType,
			wantSourceProject: "foo",
		},
		{
			desc: "set reference",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
			}).SetReference("bar"),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "bar",
				},
			},
			wantType: FileReferenceType,
		},
		{
			desc: "set inlined",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
			}).SetInlined([]byte("bar")),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte("bar"),
				},
			},
			wantType: InlinedReferenceType,
		},
		{
			desc: "merge with changed data",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest: fmt.Sprintf("%s:foo", utils.HighwayHash128),
			}).Merge(FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte("bar"),
				},
			})),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte("bar"),
				},
				Digest: fmt.Sprintf("%s:foo", utils.HighwayHash128),
			},
			wantType: InlinedReferenceType,
		},
		{
			desc: "merge with unchanged data",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			}).Merge(FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			})),
			wantModified: false,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			},
			wantType:          FileReferenceType,
			wantSourceProject: "bar",
		},
		{
			desc: "empty data merged into file reference",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo.txt",
				},
			}).Merge(FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte{},
				},
			})),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte{},
				},
			},
			wantType: InlinedReferenceType,
		},
		{
			desc: "replace with changed data",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			}).Replace(FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte("bar"),
				},
				Digest: fmt.Sprintf("%s:bar", utils.HighwayHash128),
			})),
			wantModified: true,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Inlined{
					Inlined: []byte("bar"),
				},
				Digest: fmt.Sprintf("%s:bar", utils.HighwayHash128),
			},
			wantType: InlinedReferenceType,
		},
		{
			desc: "replace with unchanged data",
			ref: FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			}).Replace(FromProto(&rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			})),
			wantModified: false,
			wantReferencedData: &rdpb.ReferencedData{
				Data: &rdpb.ReferencedData_Reference{
					Reference: "foo",
				},
				Digest:        fmt.Sprintf("%s:foo", utils.HighwayHash128),
				SourceProject: proto.String("bar"),
			},
			wantType:          FileReferenceType,
			wantSourceProject: "bar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if diff := cmp.Diff(tc.wantReferencedData, tc.ref.rd, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected diff (-want +got): %v", diff)
			}
			if diff := cmp.Diff(tc.wantReferencedData, tc.ref.ToProto(), protocmp.Transform()); diff != "" {
				t.Errorf("ToProto() returned unexpected diff (-want +got): %v", diff)
			}
			if tc.ref.SourceProject() != tc.wantSourceProject {
				t.Errorf("SourceProject() returned an unexpected project, got: %v, want: %v", tc.ref.SourceProject(), tc.wantSourceProject)
			}
			if tc.ref.Modified() != tc.wantModified {
				t.Errorf("Modified() returned unexpected value, got: %v, want: %v", tc.ref.Modified(), tc.wantModified)
			}
			if tc.ref.Type() != tc.wantType {
				t.Errorf("Type() returned an unexpected type, got: %v, want: %v", tc.ref.Type(), tc.wantType)
			}
		})
	}
}

func TestReferencedDataCopy(t *testing.T) {
	ref := FromProto(&rdpb.ReferencedData{
		Data: &rdpb.ReferencedData_Reference{
			Reference: "foo",
		},
		Digest:        fmt.Sprintf("%s:bar", utils.HighwayHash128),
		SourceProject: proto.String("baz"),
	})

	refCopy := ref.Copy()
	if refCopy.Modified() {
		t.Errorf("Copy() returned a modified copy, want: false, got: true")
	}
	if refCopy.Reference() != "foo" {
		t.Errorf("Copy() did not copy reference, got: %v, want: %v", ref.Reference(), "foo")
	}
	wantDigest := fmt.Sprintf("%s:bar", utils.HighwayHash128)
	if refCopy.Digest() != wantDigest {
		t.Errorf("Copy() did not copy digest, got: %v, want: %v", ref.Digest(), wantDigest)
	}

	refCopy.SetDigest(fmt.Sprintf("%s:baz", utils.HighwayHash128))
	if ref.Digest() != fmt.Sprintf("%s:bar", utils.HighwayHash128) {
		t.Errorf("modifying copy digest modified the original")
	}
	refCopy.SetSourceProject("qux")
	if ref.SourceProject() != "baz" {
		t.Errorf("modifying copy source project modified the original")
	}
	refCopy.SetReference("baz")
	if ref.Reference() != "foo" {
		t.Errorf("modifying copy reference modified the original")
	}
}

func TestWalkUnique(t *testing.T) {
	emptyMsg := &emptypb.Empty{}
	makeStructMsg := func(ref1, ref2, ref3 string) *rdspb.ReferencedDataStruct {
		return &rdspb.ReferencedDataStruct{
			Fields: map[string]*rdspb.Value{
				"foo": {
					Kind: &rdspb.Value_ReferencedDataValue{
						ReferencedDataValue: &rdpb.ReferencedData{
							Data: &rdpb.ReferencedData_Reference{
								Reference: ref1,
							},
						},
					},
				},
				"foo_list": {
					Kind: &rdspb.Value_ListValue{
						ListValue: &rdspb.ListValue{
							Values: []*rdspb.Value{
								{
									Kind: &rdspb.Value_ReferencedDataValue{
										ReferencedDataValue: &rdpb.ReferencedData{
											Data: &rdpb.ReferencedData_Reference{
												Reference: ref2,
											},
										},
									},
								},
							},
						},
					},
				},
				"foo_map": {
					Kind: &rdspb.Value_StructValue{
						StructValue: &rdspb.ReferencedDataStruct{
							Fields: map[string]*rdspb.Value{
								"bar": {
									Kind: &rdspb.Value_ReferencedDataValue{
										ReferencedDataValue: &rdpb.ReferencedData{
											Data: &rdpb.ReferencedData_Reference{
												Reference: ref3,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	msg := makeStructMsg("not_bob", "also_not_bob", "still_not_bob")
	msgAny, err := anypb.New(msg)
	if err != nil {
		t.Fatalf("anypb.New(%v) failed: %v", msg, err)
	}
	msgExtracted, err := utils.ExtractPayload(&dapb.DataAsset{
		Data:              msgAny,
		FileDescriptorSet: descriptor.FileDescriptorSetFrom(msg),
	})
	if err != nil {
		t.Fatalf("ExtractPayload() failed: %v", err)
	}

	replaceWithBob := func(ref *ReferencedData) error {
		ref.SetReference("bob")
		return nil
	}
	failIfAlsoNotBob := func(ref *ReferencedData) error {
		if ref.Reference() == "also_not_bob" {
			return fmt.Errorf("also_not_bob is not bob")
		}
		return nil
	}

	tests := []struct {
		desc      string
		msg       proto.Message
		f         ReferencedDataProcessor
		wantMsg   proto.Message
		wantError string
	}{
		{
			desc:    "empty",
			msg:     emptyMsg,
			f:       func(*ReferencedData) error { return nil },
			wantMsg: emptyMsg,
		},
		{
			desc:    "struct replace with bob",
			msg:     makeStructMsg("not_bob", "also_not_bob", "still_not_bob"),
			f:       replaceWithBob,
			wantMsg: makeStructMsg("bob", "bob", "bob"),
		},
		{
			desc:      "struct fail if also not bob",
			msg:       makeStructMsg("not_bob", "also_not_bob", "still_not_bob"),
			f:         failIfAlsoNotBob,
			wantError: "also_not_bob is not bob",
		},
		{
			desc:    "extracted",
			msg:     msgExtracted,
			f:       replaceWithBob,
			wantMsg: makeStructMsg("bob", "bob", "bob"),
		},
		{
			desc:    "verify processed references are unique",
			msg:     makeStructMsg("not_bob", "not_bob", "also_not_bob"),
			f:       makeVerifyUniqueReferenceProcessor(),
			wantMsg: makeStructMsg("not_bob", "not_bob", "also_not_bob"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			msgOut, err := WalkUnique(tc.msg, tc.f)
			if tc.wantError != "" {
				if err == nil {
					t.Errorf("WalkUnique(%v, %v) did not return expected error, want: %v", tc.msg, tc.f, tc.wantError)
				} else if !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("WalkUnique(%v, %v) returned unexpected error, got: %v, want: %v", tc.msg, tc.f, err, tc.wantError)
				}
			} else if err != nil {
				t.Errorf("WalkUnique(%v, %v) returned unexpected error: %v", tc.msg, tc.f, err)
			} else if diff := cmp.Diff(msgOut, tc.wantMsg, protocmp.Transform()); diff != "" {
				t.Errorf("WalkUnique(%v, %v) returned unexpected difference (-want +got):\n%s", tc.msg, tc.f, diff)
			}
		})
	}
}

func makeVerifyUniqueReferenceProcessor() func(ref *ReferencedData) error {
	visitedReferences := map[string]struct{}{}
	return func(ref *ReferencedData) error {
		if _, ok := visitedReferences[ref.Reference()]; ok {
			return fmt.Errorf("reference %q already visited", ref.Reference())
		}
		visitedReferences[ref.Reference()] = struct{}{}

		return nil
	}
}
