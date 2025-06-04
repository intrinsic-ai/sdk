// Copyright 2023 Intrinsic Innovation LLC

// Package volumemount provides a Service that mounts a volume.
package volumemount

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	vmpb "intrinsic/assets/services/examples/volume_mount/proto/v1/volume_mount_go_grpc_proto"
)

// Service implements VolumeMountService.
type Service struct {
	config        *vmpb.VolumeMountConfig
	rootMountPath string
}

// ServiceOptions contains options for NewService.
type ServiceOptions struct {
	Config        *vmpb.VolumeMountConfig
	RootMountPath string
}

// NewService creates a new Service.
func NewService(opts *ServiceOptions) *Service {
	return &Service{
		config:        opts.Config,
		rootMountPath: opts.RootMountPath,
	}
}

func (s *Service) ListDir(ctx context.Context, req *vmpb.ListDirRequest) (*vmpb.ListDirResponse, error) {
	path := s.toAbsPath(req.GetPath())

	var entries []*vmpb.ListDirResponse_Entry
	if req.GetRecursive() {
		if err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err == nil && p != path {
				entries = append(entries, &vmpb.ListDirResponse_Entry{
					Path:        s.toRelPath(p),
					IsDirectory: d.IsDir(),
				})
			}
			return nil
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list directory %q recursively: %v", path, err)
		}
	} else {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list directory %q: %v", path, err)
		}
		for _, dirEntry := range dirEntries {
			entries = append(entries, &vmpb.ListDirResponse_Entry{
				Path:        s.toRelPath(filepath.Join(path, dirEntry.Name())),
				IsDirectory: dirEntry.IsDir(),
			})
		}
	}

	return &vmpb.ListDirResponse{
		Entries: entries,
	}, nil
}

func (s *Service) ReadFile(ctx context.Context, req *vmpb.ReadFileRequest) (*vmpb.ReadFileResponse, error) {
	path := s.toAbsPath(req.GetPath())

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file %q: %v", path, err)
	}
	return &vmpb.ReadFileResponse{
		Contents: contents,
	}, nil
}

func (s *Service) WriteFile(ctx context.Context, req *vmpb.WriteFileRequest) (*vmpb.WriteFileResponse, error) {
	path := s.toAbsPath(req.GetPath())

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create directory %q: %v", path, err)
	}
	if err := os.WriteFile(path, req.GetContents(), 0644); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file %q: %v", path, err)
	}
	return &vmpb.WriteFileResponse{}, nil
}

// WriteInitialFiles writes the initial files to the root mount path.
func (s *Service) WriteInitialFiles(ctx context.Context) error {
	for path, contents := range s.config.GetInitialFiles() {
		if _, err := s.WriteFile(ctx, &vmpb.WriteFileRequest{
			Path:     path,
			Contents: contents,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) toAbsPath(path string) string {
	return filepath.Join(s.rootMountPath, path)
}

func (s *Service) toRelPath(path string) string {
	return strings.TrimPrefix(path, s.rootMountPath)
}
