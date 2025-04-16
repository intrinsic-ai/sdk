// Copyright 2023 Intrinsic Innovation LLC

package descriptor

import (
	"testing"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMergeFileDescriptorSets(t *testing.T) {
	tests := []struct {
		desc    string
		fdss    []*dpb.FileDescriptorSet
		options []MergeFileDescriptorSetsOption
		want    *dpb.FileDescriptorSet
		wantErr bool
	}{
		{
			desc: "empty",
			fdss: []*dpb.FileDescriptorSet{},
			want: &dpb.FileDescriptorSet{},
		},
		{
			desc: "single",
			fdss: []*dpb.FileDescriptorSet{
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
			},
			want: &dpb.FileDescriptorSet{
				File: []*dpb.FileDescriptorProto{
					&dpb.FileDescriptorProto{
						Name: proto.String("file1"),
					},
				},
			},
		},
		{
			desc: "multiple",
			fdss: []*dpb.FileDescriptorSet{
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file2"),
						},
					},
				},
			},
			want: &dpb.FileDescriptorSet{
				File: []*dpb.FileDescriptorProto{
					&dpb.FileDescriptorProto{
						Name: proto.String("file1"),
					},
					&dpb.FileDescriptorProto{
						Name: proto.String("file2"),
					},
				},
			},
		},
		{
			desc: "duplicate",
			fdss: []*dpb.FileDescriptorSet{
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
			},
			want: &dpb.FileDescriptorSet{
				File: []*dpb.FileDescriptorProto{
					&dpb.FileDescriptorProto{
						Name: proto.String("file1"),
					},
				},
			},
		},
		{
			desc: "duplicate with different contents",
			fdss: []*dpb.FileDescriptorSet{
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name:    proto.String("file1"),
							Package: proto.String("package1"),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			desc: "duplicate with different contents, with keys",
			fdss: []*dpb.FileDescriptorSet{
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name: proto.String("file1"),
						},
					},
				},
				&dpb.FileDescriptorSet{
					File: []*dpb.FileDescriptorProto{
						&dpb.FileDescriptorProto{
							Name:    proto.String("file1"),
							Package: proto.String("package1"),
						},
					},
				},
			},
			options: []MergeFileDescriptorSetsOption{
				WithKeys([]string{"key1", "key2"}),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := MergeFileDescriptorSets(tc.fdss, tc.options...)
			if tc.wantErr != (err != nil) {
				t.Errorf("MergeFileDescriptorSets() returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("MergeFileDescriptorSets() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
