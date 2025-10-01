// Copyright 2023 Intrinsic Innovation LLC

package bundleio

import (
	"fmt"
	"io"
	"os"

	"github.com/google/safearchive/tar"
	"intrinsic/assets/metadatautils"
	"intrinsic/assets/processes/processutil"
	"intrinsic/util/archive/tartooling"

	processassetpb "intrinsic/assets/processes/proto/process_asset_go_proto"
)

const (
	// ProcessAssetFileName is the name of the file in a Process asset .tar bundle
	// that contains the ProcessAsset binary proto.
	ProcessAssetFileName = "process_asset.binpb"
)

func bundleMetadataOptions() []metadatautils.ValidateMetadataOption {
	return []metadatautils.ValidateMetadataOption{
		metadatautils.WithRequireNoVersion(true),
		metadatautils.WithRequireNoOutputOnlyFields(),
	}
}

// WriteProcessAsset writes a Process asset .tar bundle file to the given path.
func WriteProcessAsset(asset *processassetpb.ProcessAsset, path string) error {
	if asset == nil {
		return fmt.Errorf("Process asset must not be nil")
	}

	err := processutil.ValidateProcessAsset(asset, bundleMetadataOptions()...)
	if err != nil {
		return fmt.Errorf("invalid Process asset: %w", err)
	}

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tarWriter := tar.NewWriter(out)

	if err := tartooling.AddBinaryProto(asset, tarWriter, ProcessAssetFileName); err != nil {
		return fmt.Errorf("cannot write Process asset to bundle: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("cannot close tar writer: %w", err)
	}

	return nil
}

// ReadProcessAsset reads a ProcessAsset from a bundle (see WriteProcessAsset).
func ReadProcessAsset(path string) (*processassetpb.ProcessAsset, error) {
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
	if header.Name != ProcessAssetFileName {
		return nil, fmt.Errorf("unexpected file in Process asset bundle: %v", header.Name)
	}

	asset := &processassetpb.ProcessAsset{}
	if err = readBinaryProto(tarReader, asset); err != nil {
		return nil, fmt.Errorf("error reading Process asset proto in bundle: %w", err)
	}

	if err := processutil.ValidateProcessAsset(asset, bundleMetadataOptions()...); err != nil {
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

	return asset, nil
}
