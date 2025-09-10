// Copyright 2023 Intrinsic Innovation LLC

package bundleio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/golang/glog"
	"github.com/google/safearchive/tar"
	"intrinsic/assets/data/utils"
	"intrinsic/util/archive/tartooling"

	anypb "google.golang.org/protobuf/types/known/anypb"
	acgrpcpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_grpc_proto"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
)

const (
	dataAssetFileName = "data_asset.binpb"
	dataFileBaseDir   = "data_files"

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

// InlineReferencedData is a ReferencedDataProcessor that inlines the data referenced by the given
// ReferencedData.
//
// NOTE: This processor will read _all_ referenced data into memory and should not be used when
// large data may be referenced.
var InlineReferencedData = &inlineReferencedData{}

// ToCatalogReferencedDataProcessor is a ReferencedDataProcessor that prepares the given
// ReferencedData for inclusion in an Asset that will be released to the AssetCatalog.
type ToCatalogReferencedDataProcessor struct {
	acClient  acgrpcpb.AssetCatalogClient
	chunkSize int
	ctx       context.Context
}

// ToCatalogReferencedDataOption is a functional option for ToCatalogReferencedData.
type ToCatalogReferencedDataOption func(*ToCatalogReferencedDataProcessor)

// WithACClient sets the AssetCatalogClient for ToCatalogReferencedData.
func WithACClient(client acgrpcpb.AssetCatalogClient) ToCatalogReferencedDataOption {
	return func(opts *ToCatalogReferencedDataProcessor) {
		opts.acClient = client
	}
}

// WithChunkSize sets the chunk size for ToCatalogReferencedData.
func WithChunkSize(size int) ToCatalogReferencedDataOption {
	return func(opts *ToCatalogReferencedDataProcessor) {
		opts.chunkSize = size
	}
}

// NeedsReaderFor returns true for file references.
func (p *ToCatalogReferencedDataProcessor) NeedsReaderFor(rt utils.ReferenceType) bool {
	return rt == utils.FileReferenceType
}

// Process prepares the given ReferencedData for inclusion in an Asset that will be released to the
// AssetCatalog.
func (p *ToCatalogReferencedDataProcessor) Process(rdr *ReferencedDataReader) error {
	if rdr.Ref.Reference() != "" {
		log.Infof("Preparing reference %v", rdr.Ref.Reference())
	}

	stream, err := p.acClient.PrepareReferencedData(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to open PrepareReferencedData stream: %w", err)
	}

	// First send the referenced data.
	if err := stream.Send(&acpb.PrepareReferencedDataRequest{
		Data: &acpb.PrepareReferencedDataRequest_ReferencedData{
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
				if err := stream.Send(&acpb.PrepareReferencedDataRequest{
					Data: &acpb.PrepareReferencedDataRequest_DataChunk{
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
// for inclusion in an asset that will be released to the AssetCatalog.
func ToCatalogReferencedData(ctx context.Context, options ...ToCatalogReferencedDataOption) *ToCatalogReferencedDataProcessor {
	p := &ToCatalogReferencedDataProcessor{
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
				return fmt.Errorf("cannot read data file: %w", err)
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

// ToPortableReferencedData is a ReferencedDataProcessor that converts the given ReferencedData to a
// portable form.
//
// File references below a size threshold are inlined. Otherwise, they are uploaded to CAS.
var ToPortableReferencedData = &toPortableReferencedData{}

// WriteDataAssetOptions contains options for a call to WriteDataAsset.
type WriteDataAssetOptions struct {
	// ExcludedReferencedFilePaths is a list of paths to files that should not be included in the tar
	// bundle.
	//
	// Relative paths must be relative to the output bundle's directory.
	//
	// These files are left as is and referenced by the Data asset along with a digest to ensure the
	// data are not modified.
	ExcludedReferencedFilePaths []string
	// ExpectedReferencedFilePaths is a list of paths to files that are expected to be referenced in
	// the Data asset.
	//
	// Relative paths must be relative to the output bundle's directory.
	ExpectedReferencedFilePaths []string
}

// WriteDataAssetOption is a functional option for WriteDataAsset.
type WriteDataAssetOption func(*WriteDataAssetOptions)

// WithExcludedReferencedFilePaths sets the ExcludedReferencedFilePaths option.
func WithExcludedReferencedFilePaths(paths []string) WriteDataAssetOption {
	return func(opts *WriteDataAssetOptions) {
		opts.ExcludedReferencedFilePaths = paths
	}
}

// WithExpectedReferencedFilePaths sets the ExpectedReferencedFilePaths option.
func WithExpectedReferencedFilePaths(paths []string) WriteDataAssetOption {
	return func(opts *WriteDataAssetOptions) {
		opts.ExpectedReferencedFilePaths = paths
	}
}

// WriteDataAsset writes a Data asset .tar bundle file to the specified path.
//
// Relative path references in the Data asset must be relative to the output bundle's directory.
func WriteDataAsset(da *dapb.DataAsset, path string, options ...WriteDataAssetOption) error {
	opts := &WriteDataAssetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if da == nil {
		return fmt.Errorf("data asset must not be nil")
	}
	if da.GetMetadata().GetAssetType() == atpb.AssetType_ASSET_TYPE_UNSPECIFIED {
		da.Metadata.AssetType = atpb.AssetType_ASSET_TYPE_DATA
	}
	if err := utils.ValidateDataAsset(da); err != nil {
		return fmt.Errorf("invalid Data asset: %w", err)
	}

	baseDir := filepath.Dir(path)

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", path, err)
	}
	defer out.Close()
	tw := tar.NewWriter(out)

	payload, err := utils.ExtractPayload(da)
	if err != nil {
		return fmt.Errorf("cannot extract data payload: %w", err)
	}

	excludedReferencedFilePaths := map[string]struct{}{}
	for _, path := range opts.ExcludedReferencedFilePaths {
		excludedReferencedFilePaths[path] = struct{}{}
	}

	expectedReferencedFilePaths := map[string]struct{}{}
	for _, path := range opts.ExpectedReferencedFilePaths {
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

		if err := utils.ValidateReferencedData(refBase); err != nil {
			return fmt.Errorf("invalid ReferencedData: %w", err)
		}

		// Process file references. We either add the file to the tar bundle or ensure the reference has
		// a digest to guard against later modification to the file.
		if ref.Type() == utils.FileReferenceType {
			referencedFilePaths[ref.Reference()] = struct{}{}

			if _, ok := excludedReferencedFilePaths[ref.Reference()]; !ok { // Add to the tar bundle.
				inBundlePath := toUniqueTarPath(ref.Reference(), dataFileBaseDir, tarPaths)
				if err := tartooling.AddFile(refBase.Reference(), tw, inBundlePath); err != nil {
					return fmt.Errorf("cannot add data file to bundle: %w", err)
				}
				ref.SetReference(inBundlePath)
			} else if ref.Digest() == "" { // Keep the file external; ensure its reference has a digest.
				file, err := os.Open(refBase.Reference())
				if err != nil {
					return fmt.Errorf("cannot open referenced file %q: %w", refBase.Reference(), err)
				}
				defer file.Close()

			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot walk referenced data: %w", err)
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
		return fmt.Errorf("cannot create Any proto for data payload: %w", err)
	}

	// Construct and add the bundle version of the Data asset.
	daOut := &dapb.DataAsset{
		Data:              payloadOutAny,
		FileDescriptorSet: da.GetFileDescriptorSet(),
		Metadata:          da.GetMetadata(),
	}
	if err := tartooling.AddBinaryProto(daOut, tw, dataAssetFileName); err != nil {
		return fmt.Errorf("cannot write Data asset to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("cannot close tar writer: %w", err)
	}

	return nil
}

// ReadDataAssetOptions contains options for a call to ReadDataAsset.
type ReadDataAssetOptions struct {
	// ProcessReferencedData is an optional function that will be called for each unique
	// ReferencedData value in the Data asset as it is read. (Note that all inlined ReferencedData are
	// considered unique.)
	//
	// If a non-nil ReferencedData is returned, the return value replaces all of the matching
	// ReferencedData values in the Data asset.
	ProcessReferencedData ReferencedDataProcessor
}

// ReadDataAssetOption is a functional option for ReadDataAsset.
type ReadDataAssetOption func(*ReadDataAssetOptions)

// WithProcessReferencedData sets the ProcessReferencedData option.
func WithProcessReferencedData(f ReferencedDataProcessor) ReadDataAssetOption {
	return func(opts *ReadDataAssetOptions) {
		opts.ProcessReferencedData = f
	}
}

// ReadDataAsset reads a DataAsset from a bundle (see WriteDataAsset).
//
// Relative file references in the Data asset must be relative to the bundle's directory.
func ReadDataAsset(path string, options ...ReadDataAssetOption) (*dapb.DataAsset, error) {
	opts := &ReadDataAssetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	baseDir := filepath.Dir(path)

	// Open the tar file for reading.
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}
	defer f.Close()

	// Read all files from the bundle, processing ReferencedData for in-tar data files as we go.
	reader := tar.NewReader(f)
	var da *dapb.DataAsset
	processedReferences := map[string]*utils.ReferencedDataExt{}
	var unknownFilePaths []string
	for {
		hdr, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error processing tar file %q: %w", path, err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		tarPath := hdr.Name
		switch tarPath {
		case dataAssetFileName:
			da = &dapb.DataAsset{}
			if err := readBinaryProto(reader, da); err != nil {
				return nil, fmt.Errorf("error reading Data asset: %w", err)
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
				if opts.ProcessReferencedData != nil {
					if err := opts.ProcessReferencedData.Process(&ReferencedDataReader{
						Ref:    ref,
						Reader: reader,
						Size:   hdr.Size,
					}); err != nil {
						return nil, fmt.Errorf("error processing ReferencedData: %w", err)
					}
				}
				processedReferences[tarPath] = ref
			} else { // Unknown file. This will be an error below.
				unknownFilePaths = append(unknownFilePaths, tarPath)
			}
		}
	}

	if len(unknownFilePaths) > 0 {
		return nil, fmt.Errorf("unknown files in %q: %v", path, unknownFilePaths)
	}
	if da == nil {
		return nil, fmt.Errorf("Data asset not found in %q", path)
	}

	// Process the payload's ReferencedData values.
	payload, err := utils.ExtractPayload(da)
	if err != nil {
		return nil, fmt.Errorf("cannot extract data payload: %w", err)
	}
	payloadOut, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		// Check whether we already processed the reference during our pass through the .tar bundle.
		if pRef, ok := processedReferences[ref.Reference()]; ok {
			ref.Merge(pRef)
			return nil
		}

		refBase := ref.Copy().SetBaseDir(baseDir)

		// Validate the ReferencedData (e.g., verify its digest).
		if err := utils.ValidateReferencedData(refBase); err != nil {
			return fmt.Errorf("invalid ReferencedData: %w", err)
		}

		// Optionally process the ReferencedData.
		if opts.ProcessReferencedData != nil {
			// Construct the ReferencedDataReader to pass to the processor.
			rdr := &ReferencedDataReader{
				Ref: ref,
			}
			switch ref.Type() {
			case utils.FileReferenceType:
				file, err := os.Open(ref.Reference())
				if err != nil {
					return fmt.Errorf("cannot open data file: %w", err)
				}
				defer file.Close()
				rdr.Reader = file

				fi, err := file.Stat()
				if err != nil {
					return fmt.Errorf("cannot stat data file: %w", err)
				}
				rdr.Size = fi.Size()
			case utils.CASReferenceType:
				if opts.ProcessReferencedData.NeedsReaderFor(utils.CASReferenceType) {
					return fmt.Errorf("CAS references cannot be read. got: %v", ref.Reference())
				}
			case utils.InlinedReferenceType:
				rdr.Reader = bytes.NewReader(ref.Inlined())
				rdr.Size = int64(len(ref.Inlined()))
			default:
				return fmt.Errorf("unknown reference type: %d", ref.Type())
			}

			if err := opts.ProcessReferencedData.Process(rdr); err != nil {
				return fmt.Errorf("error calling ReferencedDataProcessor: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot walk referenced data while reading %q: %w", path, err)
	}

	if payloadOutAny, err := anypb.New(payloadOut); err != nil {
		return nil, fmt.Errorf("cannot create Any proto for data payload: %w", err)
	} else {
		da.Data = payloadOutAny
	}

	return da, nil
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
