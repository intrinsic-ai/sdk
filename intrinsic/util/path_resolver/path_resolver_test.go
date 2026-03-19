// Copyright 2023 Intrinsic Innovation LLC

package path_resolver_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"intrinsic/util/path_resolver/pathresolver"
)

const (
	sourcePath = "intrinsic/util/path_resolver/path_resolver_test.go"
	sourceGlob = "intrinsic/*/path_resolver/path_resolver_test.g?"
)

func testThisFileContent(t *testing.T, content []byte) {
	t.Helper()

	if !strings.Contains(string(content), "Shane was here") {
		t.Errorf("Did not find expected content: %v", string(content))
	}
}

func TestResolveRunfilesPath(t *testing.T) {
	path, err := pathresolver.ResolveRunfilesPath(sourcePath)
	if err != nil {
		t.Fatalf("Unable to get location of %v: %v", sourcePath, err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Expected file to exist at %v, got error: %v", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Unable to read %v: %v", path, err)
	}

	testThisFileContent(t, content)
}

func TestResolveRunfilesPath_InvalidPath(t *testing.T) {
	nonExistentFile := "non_existent_file.txt"
	path, err := pathresolver.ResolveRunfilesPath(nonExistentFile)
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}
	if !strings.HasSuffix(path, nonExistentFile) {
		t.Errorf("Expected path to end with %q, got: %v", nonExistentFile, path)
	}

	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file not to exist at %v, but it does", path)
	}
}

func TestResolveRunfilesOrLocalPath(t *testing.T) {
	// 1. Valid path
	path, err := pathresolver.ResolveRunfilesOrLocalPath(sourcePath)
	if err != nil {
		t.Fatalf("Unable to resolve %v: %v", sourcePath, err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Expected file to exist at %v, got error: %v", path, err)
	}

	// 2. Invalid path - should still return a path ending with the filename, but not exist
	invalidPath := "some/totally/fake/path.txt"
	path, err = pathresolver.ResolveRunfilesOrLocalPath(invalidPath)
	if err != nil {
		t.Fatalf("Expected no error for non-existent path mapping, got: %v", err)
	}
	if !strings.HasSuffix(path, invalidPath) {
		t.Errorf("Expected path to end with %q, got: %v", invalidPath, path)
	}
	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file not to exist at %v, but it does", path)
	}
}

func TestGlob(t *testing.T) {
	rootFs, err := pathresolver.ResolveRunfilesFs()
	if err != nil {
		t.Fatalf("Unable to get runfiles fs: %v", err)
	}
	paths, err := fs.Glob(rootFs, sourceGlob)
	if err != nil {
		t.Fatalf("Unable to glob %v: %v", sourceGlob, err)
	}

	if len(paths) != 1 {
		t.Fatalf("len(paths) = %d , want 1: %v", len(paths), paths)
	}

	content, err := fs.ReadFile(rootFs, paths[0])
	if err != nil {
		t.Fatalf("Unable to read %v: %v", paths[0], err)
	}

	testThisFileContent(t, content)
}
