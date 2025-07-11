// Copyright 2023 Intrinsic Innovation LLC

package volumemount

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	vmpb "intrinsic/assets/services/examples/volume_mount/proto/v1/volume_mount_go_grpc_proto"
)

func TestListDir(t *testing.T) {
	tests := []struct {
		name      string
		contents  map[string][]byte
		path      string
		recursive bool
		want      *vmpb.ListDirResponse
	}{
		{
			name: "empty",
			path: "/",
			want: &vmpb.ListDirResponse{},
		},
		{
			name: "single file",
			contents: map[string][]byte{
				"hello.txt": []byte("Hello world!"),
			},
			path: "/",
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/hello.txt",
						IsDirectory: false,
					},
				},
			},
		},
		{
			name: "list root with subdirectory and non-recursive",
			contents: map[string][]byte{
				"hello.txt":          []byte("Hello world!"),
				"subdir/goodbye.txt": []byte("Goodbye world!"),
			},
			path: "/",
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/hello.txt",
						IsDirectory: false,
					},
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir",
						IsDirectory: true,
					},
				},
			},
		},
		{
			name: "list root with subdirectory and recursive",
			contents: map[string][]byte{
				"hello.txt":          []byte("Hello world!"),
				"subdir/goodbye.txt": []byte("Goodbye world!"),
			},
			path:      "/",
			recursive: true,
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/hello.txt",
						IsDirectory: false,
					},
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir",
						IsDirectory: true,
					},
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir/goodbye.txt",
						IsDirectory: false,
					},
				},
			},
		},
		{
			name: "list subdirectory non-recursive",
			contents: map[string][]byte{
				"hello.txt":          []byte("Hello world!"),
				"subdir/goodbye.txt": []byte("Goodbye world!"),
			},
			path:      "subdir",
			recursive: false,
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir/goodbye.txt",
						IsDirectory: false,
					},
				},
			},
		},
		{
			name: "list subdirectory recursive",
			contents: map[string][]byte{
				"hello.txt":          []byte("Hello world!"),
				"subdir/goodbye.txt": []byte("Goodbye world!"),
			},
			path:      "subdir",
			recursive: true,
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir/goodbye.txt",
						IsDirectory: false,
					},
				},
			},
		},
		{
			name: "list subdirectory with leading slash",
			contents: map[string][]byte{
				"subdir/hello.txt": []byte("Hello world!"),
			},
			path:      "/subdir",
			recursive: false,
			want: &vmpb.ListDirResponse{
				Entries: []*vmpb.ListDirResponse_Entry{
					&vmpb.ListDirResponse_Entry{
						Path:        "/subdir/hello.txt",
						IsDirectory: false,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			s := NewService(&ServiceOptions{
				Config: &vmpb.VolumeMountConfig{
					InitialFiles: tc.contents,
				},
				RootMountPath: t.TempDir(),
			})
			if err := s.WriteInitialFiles(ctx); err != nil {
				t.Fatalf("WriteInitialFiles() failed: %v", err)
			}

			if got, err := s.ListDir(ctx, &vmpb.ListDirRequest{
				Path:      tc.path,
				Recursive: tc.recursive,
			}); err != nil {
				t.Fatalf("ListDir() failed: %v", err)
			} else if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ListDir() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteThenReadFile(t *testing.T) {
	tests := []struct {
		name            string
		initialContents map[string][]byte
		path            string
		contents        []byte
	}{
		{
			name:     "write new file",
			path:     "hello.txt",
			contents: []byte("Hello world!"),
		},
		{
			name: "overwrite file",
			initialContents: map[string][]byte{
				"hello.txt": []byte("Hello world!"),
			},
			path:     "hello.txt",
			contents: []byte("Bon joure le monde!"),
		},
		{
			name:     "file in subdirectory",
			path:     "subdir/hello.txt",
			contents: []byte("Hello world!"),
		},
		{
			name:     "file in subdirectory with leading slash",
			path:     "/subdir/hello.txt",
			contents: []byte("Hello world!"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			s := NewService(&ServiceOptions{
				Config: &vmpb.VolumeMountConfig{
					InitialFiles: tc.initialContents,
				},
				RootMountPath: t.TempDir(),
			})
			if err := s.WriteInitialFiles(ctx); err != nil {
				t.Fatalf("WriteInitialFiles() failed: %v", err)
			}

			if _, err := s.WriteFile(ctx, &vmpb.WriteFileRequest{
				Path:     tc.path,
				Contents: tc.contents,
			}); err != nil {
				t.Fatalf("WriteFile() failed: %v", err)
			}
			if got, err := s.ReadFile(ctx, &vmpb.ReadFileRequest{
				Path: tc.path,
			}); err != nil {
				t.Fatalf("ReadFile() failed: %v", err)
			} else if diff := cmp.Diff(tc.contents, got.GetContents()); diff != "" {
				t.Errorf("ReadFile() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}
