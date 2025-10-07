// Copyright 2023 Intrinsic Innovation LLC

package bundleio

import (
	"fmt"
	"io"
	"os"

	"github.com/google/safearchive/tar"
	"intrinsic/assets/processes/processutil"
	"intrinsic/util/archive/tartooling"

	processmanifestpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
)

const (
	// ProcessManifestFileName is the name of the file in a Process asset .tar bundle
	// that contains the ProcessManifest binary proto.
	ProcessManifestFileName = "process_manifest.binpb"
)

// WriteProcessManifest writes a Process asset .tar bundle file to the given path.
func WriteProcessManifest(manifest *processmanifestpb.ProcessManifest, path string) error {
	if manifest == nil {
		return fmt.Errorf("Process manifest must not be nil")
	}

	err := processutil.ValidateProcessManifest(manifest)
	if err != nil {
		return fmt.Errorf("invalid Process manifest: %w", err)
	}

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tarWriter := tar.NewWriter(out)

	if err := tartooling.AddBinaryProto(manifest, tarWriter, ProcessManifestFileName); err != nil {
		return fmt.Errorf("cannot write Process manifest to bundle: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("cannot close tar writer: %w", err)
	}

	return nil
}

// ReadProcessManifest reads a ProcessManifest from a bundle (see WriteProcessManifest).
func ReadProcessManifest(path string) (*processmanifestpb.ProcessManifest, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}
	defer file.Close()

	// Read single file from the bundle.
	tarReader := tar.NewReader(file)
	header, err := tarReader.Next()
	if err != nil {
		return nil, fmt.Errorf("could not read first entry of Process asset bundle: %w", err)
	}
	if header.Typeflag != tar.TypeReg {
		return nil, fmt.Errorf("unexpected entry type in Process asset bundle: %v", header.Typeflag)
	}
	if header.Name != ProcessManifestFileName {
		return nil, fmt.Errorf("unexpected file in Process asset bundle: %v", header.Name)
	}

	manifest := &processmanifestpb.ProcessManifest{}
	if err = readBinaryProto(tarReader, manifest); err != nil {
		return nil, fmt.Errorf("error reading ProcessManifest proto in bundle: %w", err)
	}

	if err := processutil.ValidateProcessManifest(manifest); err != nil {
		return nil, fmt.Errorf("invalid Process asset: %w", err)
	}

	// Fail if there are other files in the bundle.
	header, err = tarReader.Next()
	if err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("error reading second entry from Process asset bundle: %w", err)
		}
		return nil, fmt.Errorf("unexpected second entry in Process asset bundle: %v", header.Name)
	}

	return manifest, nil
}
