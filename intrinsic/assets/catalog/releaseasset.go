// Copyright 2023 Intrinsic Innovation LLC

// Package releaseasset provides utils for releasing Assets to the AssetCatalog.
package releaseasset

import (
	"context"
	"fmt"

	"intrinsic/assets/bundle"
	"intrinsic/assets/data/databundle"
	"intrinsic/assets/idutils"
	"intrinsic/assets/imagetransfer"
	"intrinsic/assets/services/bundleimages"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	rmpb "intrinsic/assets/catalog/proto/v1/release_metadata_go_proto"
)

const (
)

// Printer is a function that prints a formatted status message about the release.
type Printer func(format string, a ...any)

type fromBundleOptions struct {
	acClient        acpb.AssetCatalogClient
	dryRun          bool
	flagDefault     bool
	flagOrgPrivate  bool
	ignoreExisting  bool
	imageTransferer imagetransfer.Transferer
	printer         Printer
	registry        string
	releaseNotes    string
	version         string
}

// FromBundleOption is an option for FromBundle.
type FromBundleOption func(*fromBundleOptions)

// WithAssetCatalogClient specifies the client to use for catalog requests.
func WithAssetCatalogClient(acc acpb.AssetCatalogClient) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.acClient = acc
	}
}

// WithConnection specifies the connection to use for all gRPC clients.
func WithConnection(conn *grpc.ClientConn) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.acClient = acpb.NewAssetCatalogClient(conn)
	}
}

// WithDryRun specifies whether to do a dry run of the release.
func WithDryRun(dryRun bool) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.dryRun = dryRun
	}
}

// WithFlagDefault specifies whether the released Asset should be the default version.
func WithFlagDefault(flagDefault bool) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.flagDefault = flagDefault
	}
}

// WithFlagOrgPrivate specifies whether the Asset should be org-private.
func WithFlagOrgPrivate(flagOrgPrivate bool) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.flagOrgPrivate = flagOrgPrivate
	}
}

// WithIgnoreExisting specifies whether to ignore errors if the Asset version already exists in the
// catalog.
func WithIgnoreExisting(ignoreExisting bool) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.ignoreExisting = ignoreExisting
	}
}

// WithImageTransferer specifies the image transferer to use.
func WithImageTransferer(transferer imagetransfer.Transferer) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.imageTransferer = transferer
	}
}

func WithPrinter(printer Printer) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.printer = printer
	}
}

// WithRegistry specifies the artifact registry to use.
func WithRegistry(registry string) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.registry = registry
	}
}

// WithReleaseNotes specifies release notes to include with the Asset.
func WithReleaseNotes(releaseNotes string) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.releaseNotes = releaseNotes
	}
}

// WithVersion specifies the version of the Asset.
func WithVersion(version string) FromBundleOption {
	return func(opts *fromBundleOptions) {
		opts.version = version
	}
}

func FromBundle(ctx context.Context, path string, options ...FromBundleOption) error {
	opts := &fromBundleOptions{
		printer: nullPrinter,
	}
	for _, opt := range options {
		opt(opts)
	}

	if opts.acClient == nil {
		return fmt.Errorf("acClient must not be nil")
	}
	if opts.imageTransferer == nil {
		return fmt.Errorf("transferer must not be nil")
	}
	if opts.registry == "" {
		return fmt.Errorf("registry must not be empty")
	}
	if opts.version == "" {
		return fmt.Errorf("version must not be empty")
	}

	referencedDataProcessor := databundle.NoOpReferencedData()
	if !opts.dryRun {
		referencedDataProcessor = databundle.ToCatalogReferencedData(ctx, databundle.WithACClient(opts.acClient))
	}

	processor := bundle.Processor{
		ImageProcessor: bundleimages.CreateImageProcessor(bundleimages.RegistryOptions{
			Transferer: opts.imageTransferer,
			URI:        opts.registry,
		}),
		ProcessReferencedData:   referencedDataProcessor,
	}

	processedBundle, err := processor.Process(ctx, path)
	if err != nil {
		return fmt.Errorf("unable to process Asset: %w", err)
	}

	asset := processedBundle.Release(bundle.VersionDetails{
		Version:      opts.version,
		ReleaseNotes: opts.releaseNotes,
		ReleaseMetadata: &rmpb.ReleaseMetadata{
			Default:    opts.flagDefault,
			OrgPrivate: opts.flagOrgPrivate,
		},
	})

	idVersion := idutils.IDVersionFromProtoUnchecked(asset.GetMetadata().GetIdVersion())
	opts.printer("Releasing Asset %q to the catalog", idVersion)

	if opts.dryRun {
		opts.printer("Skipping release: dry-run")
		return nil
	}

	if _, err := opts.acClient.CreateAsset(ctx, &acpb.CreateAssetRequest{
		Asset: asset,
	}); err != nil {
		if s, ok := status.FromError(err); ok && opts.ignoreExisting && s.Code() == codes.AlreadyExists {
			opts.printer("Skipping release: Asset already exists in the catalog")
			return nil
		}
		return fmt.Errorf("could not release the Asset: %w", err)
	}

	opts.printer("Finished releasing the Asset")

	return nil
}

func nullPrinter(format string, a ...any) {
}
