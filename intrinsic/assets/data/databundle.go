// Copyright 2023 Intrinsic Innovation LLC

// Package databundle provides utils for working with Data Asset bundles.
package databundle

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"intrinsic/assets/data/datavalidate"
	"intrinsic/assets/data/utils"
	"intrinsic/assets/ioutils"
	"intrinsic/util/archive/tartooling"

	log "github.com/golang/glog"
	"github.com/google/safearchive/tar"

	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	dataAssetFileName = "data_asset.binpb"
	dataFileBaseDir = "data_files"

	// inlineReferenceFileSizeThresholdBytes is the file size threshold for inlining ReferencedData.
	// If a referenced file is <= than this threshold, it is inlined. Otherwise, it is uploaded to
	// CAS.
	inlineReferenceFileSizeThresholdBytes = 1024 * 1024

	defaultChunkSize = 1024 * 1024
)

// ReferencedDataReader contains a ReferencedData message and additional fields for reading it.
type ReferencedDataReader struct {
	// Ref is the referenced data.
	Ref *utils.ReferencedDataExt
	// Reader can be used to read the data referenced by the ReferencedData.
	Reader io.Reader
	// Size is the size of the referenced data, in bytes.
	Size int64
}

// ReferencedDataProcessor is an interface for processing ReferencedData.
type ReferencedDataProcessor interface {
	// NeedsReaderFor returns true if the processor needs an io.Reader for the given reference type.
	NeedsReaderFor(utils.ReferenceType) bool

	// Process is called for each ReferencedData value in the Data Asset. The processor should
	// modify the ReferencedData in place.
	Process(*ReferencedDataReader) error
}

type noOpReferencedData struct{}

// NeedsReaderFor returns false for all reference types.
func (p *noOpReferencedData) NeedsReaderFor(rt utils.ReferenceType) bool {
	return false
}

// Process does not modify the given reference data.
func (p *noOpReferencedData) Process(rdr *ReferencedDataReader) error {
	return nil
}

// NoOpReferencedData returns a ReferencedDataProcessor that does nothing.
//
// This processor is only valid for dry runs, but it ensures that referenced data are available.
func NoOpReferencedData() ReferencedDataProcessor {
	return &noOpReferencedData{}
}

type inlineReferencedData struct{}

// NeedsReaderFor returns true for all reference types.
func (p *inlineReferencedData) NeedsReaderFor(rt utils.ReferenceType) bool {
	return true
}

// Process inlines the data referenced by the given ReferencedData.
func (p *inlineReferencedData) Process(rdr *ReferencedDataReader) error {
	// Nothing to do for already inlined data.
	if rdr.Ref.Reference() == "" {
		return nil
	}

	if rdr.Reader == nil {
		return fmt.Errorf("no reader for referenced data: %v", rdr.Ref)
	}

	b, err := io.ReadAll(rdr.Reader)
	if err != nil {
		return fmt.Errorf("cannot read data: %w", err)
	}

	rdr.Ref.SetInlined(b)

	return nil
}

// InlineReferencedData returns a ReferencedDataProcessor that inlines the data referenced by the
// given ReferencedData.
//
// NOTE: This processor will read _all_ referenced data into memory and should not be used when
// large data may be referenced.
func InlineReferencedData() ReferencedDataProcessor {
	return &inlineReferencedData{}
}

type toCatalogReferencedData struct {
	acClient  acgrpcpb.AssetCatalogClient
	chunkSize int
	ctx       context.Context
}

// ToCatalogReferencedDataOption is a functional option for ToCatalogReferencedData.
type ToCatalogReferencedDataOption func(*toCatalogReferencedData)

// WithACClient sets the AssetCatalogClient for ToCatalogReferencedData.
func WithACClient(client acgrpcpb.AssetCatalogClient) ToCatalogReferencedDataOption {
	return func(opts *toCatalogReferencedData) {
		opts.acClient = client
	}
}

// WithChunkSize sets the chunk size for ToCatalogReferencedData.
func WithChunkSize(size int) ToCatalogReferencedDataOption {
	return func(opts *toCatalogReferencedData) {
		opts.chunkSize = size
	}
}

// NeedsReaderFor returns true for file references.
func (p *toCatalogReferencedData) NeedsReaderFor(rt utils.ReferenceType) bool {
	return rt == utils.FileReferenceType
}

// Process prepares the given ReferencedData for inclusion in an Asset that will be released to the
// AssetCatalog.
func (p *toCatalogReferencedData) Process(rdr *ReferencedDataReader) error {
	if rdr.Ref.Reference() != "" {
		log.Infof("Preparing reference %v", rdr.Ref.Reference())
	}

	stream, err := p.acClient.PrepareReferencedData(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to open PrepareReferencedData stream: %w", err)
	}

	// First send the referenced data.
	if err := stream.Send(&acgrpcpb.PrepareReferencedDataRequest{
		Data: &acgrpcpb.PrepareReferencedDataRequest_ReferencedData{
			ReferencedData: rdr.Ref.ToProto(),
		},
	}); err != nil {
		return fmt.Errorf("failed to send referenced data: %w", err)
	}

	// For file references, send the file data.
	if rdr.Ref.Type() == utils.FileReferenceType {
		log.Infof("Sending file data for %v", rdr.Ref.Reference())
		buf := make([]byte, p.chunkSize)
		for {
			n, err := rdr.Reader.Read(buf)
			if err != io.EOF && err != nil {
				return fmt.Errorf("failed to read data: %w", err)
			}
			if n > 0 {
				if err := stream.Send(&acgrpcpb.PrepareReferencedDataRequest{
					Data: &acgrpcpb.PrepareReferencedDataRequest_DataChunk{
						DataChunk: buf[:n],
					},
				}); err != nil {
					return fmt.Errorf("failed to send data chunk: %w", err)
				}
			}
			if err == io.EOF {
				break
			}
		}
	}

	// Close the stream and get the updated referenced data.
	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to close stream: %w", err)
	}

	// Replace the referenced data with the updated referenced data from the catalog.
	rdr.Ref.Replace(utils.NewReferencedDataExt(response.GetReferencedData()))

	return nil
}

// ToCatalogReferencedData returns a ReferencedDataProcessor that prepares the given ReferencedData
// for inclusion in an Asset that will be released to the AssetCatalog.
func ToCatalogReferencedData(ctx context.Context, options ...ToCatalogReferencedDataOption) *toCatalogReferencedData {
	p := &toCatalogReferencedData{
		ctx:       ctx,
		chunkSize: defaultChunkSize,
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

type toPortableReferencedData struct{}

func (p *toPortableReferencedData) NeedsReaderFor(rt utils.ReferenceType) bool {
	return rt == utils.FileReferenceType
}

func (p *toPortableReferencedData) Process(rdr *ReferencedDataReader) error {
	switch rdr.Ref.Type() {
	case utils.FileReferenceType:
		// If the file is below the size threshold, inline it. Otherwise, upload it to CAS.
		if rdr.Size <= inlineReferenceFileSizeThresholdBytes {
			b, err := io.ReadAll(rdr.Reader)
			if err != nil {
				return fmt.Errorf("failed to read data file: %w", err)
			}
			rdr.Ref.SetInlined(b)
		} else {
			return fmt.Errorf("file upload is not supported: %v", rdr.Ref)
		}
	case utils.CASReferenceType:
	case utils.InlinedReferenceType:
		// Nothing to do.
	default:
		return fmt.Errorf("unknown referenced data: %v", rdr.Ref)
	}

	return nil
}

// ToPortableReferencedData returns a ReferencedDataProcessor that converts the given ReferencedData
// to a portable form.
//
// File references below a size threshold are inlined. Otherwise, they are uploaded to CAS.
func ToPortableReferencedData() ReferencedDataProcessor {
	return &toPortableReferencedData{}
}

type writeOptions struct {
	excludedReferencedFilePaths []string
	expectedReferencedFilePaths []string
	writer                      io.Writer
}

// WriteOption is a functional option for Write.
type WriteOption func(*writeOptions)

// WithExcludedReferencedFilePaths provides a list of paths to files that should not be included in
// the .tar bundle.
//
// Relative paths must be relative to the output bundle's directory.
//
// These files are left as is and referenced by the Data Asset along with a digest to ensure the
// data are not modified.
func WithExcludedReferencedFilePaths(paths []string) WriteOption {
	return func(opts *writeOptions) {
		opts.excludedReferencedFilePaths = paths
	}
}

// WithExpectedReferencedFilePaths provides a list of paths to files that are expected to be
// referenced in the Data Asset.
//
// Relative paths must be relative to the output bundle's directory.
func WithExpectedReferencedFilePaths(paths []string) WriteOption {
	return func(opts *writeOptions) {
		opts.expectedReferencedFilePaths = paths
	}
}

// WithWriter specifies the Writer to use instead of creating one for the specified path.
func WithWriter(w io.Writer) WriteOption {
	return func(opts *writeOptions) {
		opts.writer = w
	}
}

// Write writes a Data Asset .tar bundle.
//
// Relative path references in the Data Asset must be relative to the output bundle's directory.
func Write(da *dapb.DataAsset, path string, options ...WriteOption) error {
	opts := &writeOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if da == nil {
		return fmt.Errorf("DataAsset must not be nil")
	}
	if da.GetMetadata().GetAssetType() == atpb.AssetType_ASSET_TYPE_UNSPECIFIED {
		da.Metadata.AssetType = atpb.AssetType_ASSET_TYPE_DATA
	}
	if err := datavalidate.DataAsset(da); err != nil {
		return fmt.Errorf("invalid DataAsset: %w", err)
	}

	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	baseDir := filepath.Dir(path)

	writer := opts.writer
	if writer == nil {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open %q for writing: %w", path, err)
		}
		defer f.Close()
		writer = f
	}

	tw := tar.NewWriter(writer)

	payload, err := utils.ExtractPayload(da)
	if err != nil {
		return fmt.Errorf("failed to extract data payload: %w", err)
	}

	excludedReferencedFilePaths := map[string]struct{}{}
	for _, path := range opts.excludedReferencedFilePaths {
		excludedReferencedFilePaths[path] = struct{}{}
	}

	expectedReferencedFilePaths := map[string]struct{}{}
	for _, path := range opts.expectedReferencedFilePaths {
		expectedReferencedFilePaths[path] = struct{}{}
	}

	// Records all file path references that were found.
	referencedFilePaths := map[string]struct{}{}

	// Walk through the data payload. For each ReferencedData:
	// - Validate it;
	// - If it is a file reference that is not excluded, copy it to the tar bundle;
	// - If it is a file reference that is excluded, ensure it has a digest.
	tarPaths := map[string]struct{}{} // Keeps track of used tar paths.
	payloadOut, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		refBase := ref.Copy().SetBaseDir(baseDir)

		if err := datavalidate.ReferencedData(refBase); err != nil {
			return fmt.Errorf("invalid ReferencedData: %w", err)
		}

		// Process file references. We either add the file to the tar bundle or ensure the reference has
		// a digest to guard against later modification to the file.
		if ref.Type() == utils.FileReferenceType {
			referencedFilePaths[ref.Reference()] = struct{}{}

			if _, ok := excludedReferencedFilePaths[ref.Reference()]; !ok { // Add to the tar bundle.
				inBundlePath := toUniqueTarPath(ref.Reference(), dataFileBaseDir, tarPaths)
				if err := tartooling.AddFile(refBase.Reference(), tw, inBundlePath); err != nil {
					return fmt.Errorf("failed to add data file to bundle: %w", err)
				}
				ref.SetReference(inBundlePath)
			} else if ref.Digest() == "" { // Keep the file external; ensure its reference has a digest.
				file, err := os.Open(refBase.Reference())
				if err != nil {
					return fmt.Errorf("failed to open referenced file %q: %w", refBase.Reference(), err)
				}
				defer file.Close()

			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk referenced data: %w", err)
	}

	// Verify that all excluded referenced file paths are actually referenced by the data payload.
	for path := range excludedReferencedFilePaths {
		if _, ok := referencedFilePaths[path]; !ok {
			return fmt.Errorf("excluded referenced file path %q is not referenced by the data payload. referenced: %v", path, referencedFilePaths)
		}
	}

	// Verify that all expected referenced file paths are actually referenced by the data payload.
	for path := range expectedReferencedFilePaths {
		if _, ok := referencedFilePaths[path]; !ok {
			return fmt.Errorf("expected referenced file path %q is not referenced by the data payload. referenced: %v", path, referencedFilePaths)
		}
	}

	payloadOutAny, err := anypb.New(payloadOut)
	if err != nil {
		return fmt.Errorf("failed to create Any proto for data payload: %w", err)
	}

	// Construct and add the bundle version of the Data Asset.
	daOut := &dapb.DataAsset{
		Data:              payloadOutAny,
		FileDescriptorSet: da.GetFileDescriptorSet(),
		Metadata:          da.GetMetadata(),
	}
	if err := tartooling.AddBinaryProto(daOut, tw, dataAssetFileName); err != nil {
		return fmt.Errorf("failed to write DataAsset to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	return nil
}

// DataBundle represents a Data Asset bundle.
type DataBundle struct {
	Data *dapb.DataAsset
}

type readOptions struct {
	processReferencedData ReferencedDataProcessor
	reader                io.Reader
}

// ReadOption is a functional option for Read.
type ReadOption func(*readOptions)

// WithProcessReferencedData specifies a ReferencedDataProcessor to call for each unique
// ReferencedData value in the Data Asset as it is read.
//
// (Note that all inlined ReferencedData are considered unique.)
//
// If a non-nil ReferencedData is returned by the processor, the return value replaces all of the
// matching ReferencedData values in the Data Asset.
func WithProcessReferencedData(f ReferencedDataProcessor) ReadOption {
	return func(opts *readOptions) {
		opts.processReferencedData = f
	}
}

// WithReader specifies the Reader to use instead of creating one for the specified path.
func WithReader(r io.Reader) ReadOption {
	return func(opts *readOptions) {
		opts.reader = r
	}
}

// Read reads a Data Asset bundle (see Write).
//
// Relative file references in the Data Asset must be relative to the bundle's directory.
func Read(ctx context.Context, path string, options ...ReadOption) (*DataBundle, error) {
	opts := &readOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}
	baseDir := filepath.Dir(path)

	reader := opts.reader
	if reader == nil {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %q for reading: %w", path, err)
		}
		defer f.Close()
		reader = f
	}

	tr := tar.NewReader(reader)

	// Read all files from the bundle, processing ReferencedData for in-tar data files as we go.
	var da *dapb.DataAsset
	processedReferences := map[string]*utils.ReferencedDataExt{}
	var unknownFilePaths []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to process tar file: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		tarPath := hdr.Name
		switch tarPath {
		case dataAssetFileName:
			da = &dapb.DataAsset{}
			if err := ioutils.ReadBinaryProto(tr, da); err != nil {
				return nil, fmt.Errorf("failed to read DataAsset: %w", err)
			}
			da.GetMetadata().AssetType = atpb.AssetType_ASSET_TYPE_DATA
		default:
			if strings.HasPrefix(tarPath, dataFileBaseDir) { // Found a referenced data file.
				// Process the reference for the in-tar file now, since we won't be able to pass the reader
				// to the processor below when we walk the payload.
				// Note that we don't validate the referenced data in this case.
				ref := utils.NewReferencedDataExt(&rdpb.ReferencedData{
					Data: &rdpb.ReferencedData_Reference{
						Reference: tarPath,
					},
				})
				if opts.processReferencedData != nil {
					if err := opts.processReferencedData.Process(&ReferencedDataReader{
						Ref:    ref,
						Reader: tr,
						Size:   hdr.Size,
					}); err != nil {
						return nil, fmt.Errorf("failed to process ReferencedData: %w", err)
					}
				}
				processedReferences[tarPath] = ref
			} else { // Unknown file. This will be an error below.
				unknownFilePaths = append(unknownFilePaths, tarPath)
			}
		}
	}

	if len(unknownFilePaths) > 0 {
		return nil, fmt.Errorf("unknown files: %v", unknownFilePaths)
	}
	if da == nil {
		return nil, fmt.Errorf("DataAsset not found in tar file")
	}

	// Process the payload's ReferencedData values.
	payload, err := utils.ExtractPayload(da)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data payload: %w", err)
	}
	payloadOut, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		// Check whether we already processed the reference during our pass through the .tar bundle.
		if pRef, ok := processedReferences[ref.Reference()]; ok {
			ref.Merge(pRef)
			return nil
		}

		refBase := ref.Copy().SetBaseDir(baseDir)

		// Validate the ReferencedData (e.g., verify its digest).
		if err := datavalidate.ReferencedData(refBase); err != nil {
			return fmt.Errorf("invalid ReferencedData: %w", err)
		}

		// Optionally process the ReferencedData.
		if opts.processReferencedData != nil {
			// Construct the ReferencedDataReader to pass to the processor.
			rdr := &ReferencedDataReader{
				Ref: ref,
			}
			switch ref.Type() {
			case utils.FileReferenceType:
				file, err := os.Open(ref.Reference())
				if err != nil {
					return fmt.Errorf("failed to open data file: %w", err)
				}
				defer file.Close()
				rdr.Reader = file

				fi, err := file.Stat()
				if err != nil {
					return fmt.Errorf("failed to stat data file: %w", err)
				}
				rdr.Size = fi.Size()
			case utils.CASReferenceType:
				if opts.processReferencedData.NeedsReaderFor(utils.CASReferenceType) {
					return fmt.Errorf("CAS references cannot be read. got: %v", ref.Reference())
				}
			case utils.InlinedReferenceType:
				rdr.Reader = bytes.NewReader(ref.Inlined())
				rdr.Size = int64(len(ref.Inlined()))
			default:
				return fmt.Errorf("unknown reference type: %d", ref.Type())
			}

			if err := opts.processReferencedData.Process(rdr); err != nil {
				return fmt.Errorf("failed to process ReferencedData: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk referenced data: %w", err)
	}

	if payloadOutAny, err := anypb.New(payloadOut); err != nil {
		return nil, fmt.Errorf("failed to create Any proto for data payload: %w", err)
	} else {
		da.Data = payloadOutAny
	}

	return &DataBundle{
		Data: da,
	}, nil
}

type processOptions struct {
	readOptions []ReadOption
}

// ProcessOption is a functional option for Process.
type ProcessOption func(*processOptions)

// WithReadOptions specifies the ReadOptions to use when reading the Data Asset.
func WithReadOptions(options ...ReadOption) ProcessOption {
	return func(opts *processOptions) {
		opts.readOptions = options
	}
}

// Process creates a processed Data Asset from a bundle.
func Process(ctx context.Context, path string, options ...ProcessOption) (*dapb.DataAsset, error) {
	opts := &processOptions{}
	for _, opt := range options {
		opt(opts)
	}

	bundle, err := Read(ctx, path, opts.readOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to read Data bundle: %w", err)
	}

	return bundle.Data, nil
}

// toUniqueTarPath generates a unique path in the tar bundle to which to save a data file.
//
// `tarPaths` keeps track of tar paths that have already been used.
func toUniqueTarPath(path string, tarDir string, tarPaths map[string]struct{}) string {
	fileName := filepath.Base(path)
	fileExt := filepath.Ext(fileName)
	filePre := strings.TrimSuffix(fileName, fileExt)
	suffix := ""
	for i := 1; ; i++ {
		tarPath := filepath.Join(tarDir, fmt.Sprintf("%s%s%s", filePre, suffix, fileExt))
		if _, ok := tarPaths[tarPath]; !ok {
			tarPaths[tarPath] = struct{}{}
			return tarPath
		}
		suffix = fmt.Sprintf("%03d", i)
	}
}
