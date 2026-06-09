// Copyright 2023 Intrinsic Innovation LLC

// Package referenceddata provides utils for working with ReferencedData.
package referenceddata

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"intrinsic/util/proto/walkmessages"

	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	acpb "intrinsic/assets/catalog/proto/v1/asset_catalog_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
)

// ReferenceType is an enum for the type of reference in a ReferencedData.
type ReferenceType int

const (
	// FileReferenceType is a file reference.
	FileReferenceType ReferenceType = iota
	// CASReferenceType is a CAS reference.
	CASReferenceType
	// InlinedReferenceType is an inlined reference.
	InlinedReferenceType
)

const (
	defaultChunkSize                      = 1024 * 1024
	InlineReferenceFileSizeThresholdBytes = 1024 * 1024
)

// ReferencedData represents a reference to data (e.g., in a file or the cloud).
type ReferencedData struct {
	modified bool
	rd       *rdpb.ReferencedData
	refType  ReferenceType
}

// FromProto creates a new ReferencedData from a proto.
func FromProto(rd *rdpb.ReferencedData) *ReferencedData {
	ref := &ReferencedData{
		rd: &rdpb.ReferencedData{
			Digest:        rd.GetDigest(),
			SourceProject: rd.SourceProject,
		},
	}

	if rd.GetReference() != "" {
		var reference string
		ref.refType, reference, ref.modified = parseReference(rd.GetReference())
		ref.rd.Data = &rdpb.ReferencedData_Reference{
			Reference: reference,
		}
	} else {
		ref.refType = InlinedReferenceType
		if rd.GetInlined() != nil {
			ref.rd.Data = &rdpb.ReferencedData_Inlined{
				Inlined: rd.GetInlined(),
			}
		}
	}

	return ref
}

// Digest returns the digest in the ReferencedData.
func (ref *ReferencedData) Digest() string {
	return ref.rd.GetDigest()
}

// Inlined returns the inlined data in the ReferencedData.
func (ref *ReferencedData) Inlined() []byte {
	return ref.rd.GetInlined()
}

// Modified returns whether the ReferencedData has been modified.
func (ref *ReferencedData) Modified() bool {
	return ref.modified
}

// Reference returns the reference in the ReferencedData.
func (ref *ReferencedData) Reference() string {
	return ref.rd.GetReference()
}

// SourceProject returns the source project in the ReferencedData.
func (ref *ReferencedData) SourceProject() string {
	return ref.rd.GetSourceProject()
}

// Type returns the type of reference in the ReferencedData.
func (ref *ReferencedData) Type() ReferenceType {
	return ref.refType
}

// Equal returns whether the two ReferencedData are equal.
func (ref *ReferencedData) Equal(other *ReferencedData) bool {
	return proto.Equal(ref.rd, other.rd)
}

// SetBaseDir sets the base directory for relative file references.
func (ref *ReferencedData) SetBaseDir(baseDir string) *ReferencedData {
	if ref.Type() == FileReferenceType && !filepath.IsAbs(ref.rd.GetReference()) {
		ref.rd.Data = &rdpb.ReferencedData_Reference{
			Reference: filepath.Join(baseDir, ref.rd.GetReference()),
		}
		ref.modified = true
	}
	return ref
}

// SetDigest sets the digest in the ReferencedData.
func (ref *ReferencedData) SetDigest(digest string) *ReferencedData {
	ref.rd.Digest = digest
	ref.modified = true
	return ref
}

// SetSourceProject sets the source project in the ReferencedData.
func (ref *ReferencedData) SetSourceProject(sourceProject string) *ReferencedData {
	ref.rd.SourceProject = proto.String(sourceProject)
	ref.modified = true
	return ref
}

// SetInlined sets the inlined data in the ReferencedData.
func (ref *ReferencedData) SetInlined(inlined []byte) *ReferencedData {
	ref.rd.Data = &rdpb.ReferencedData_Inlined{
		Inlined: inlined,
	}
	ref.refType = InlinedReferenceType
	ref.modified = true
	return ref
}

// SetReference sets the reference in the ReferencedData.
func (ref *ReferencedData) SetReference(reference string) *ReferencedData {
	ref.refType, reference, _ = parseReference(reference)
	ref.rd.Data = &rdpb.ReferencedData_Reference{
		Reference: reference,
	}
	ref.modified = true
	return ref
}

// Copy returns a shallow copy of the ReferencedData.
//
// The modified flag is reset.
func (ref *ReferencedData) Copy() *ReferencedData {
	return FromProto(ref.rd)
}

// Merge merges the data from another ReferencedData into this one.
func (ref *ReferencedData) Merge(other *ReferencedData) *ReferencedData {
	if other.Digest() != "" && ref.Digest() != other.Digest() {
		ref.SetDigest(other.Digest())
	}
	if other.SourceProject() != "" && ref.SourceProject() != other.SourceProject() {
		ref.SetSourceProject(other.SourceProject())
	}
	switch other.Type() {
	case InlinedReferenceType:
		if ref.Type() != InlinedReferenceType || !bytes.Equal(ref.Inlined(), other.Inlined()) {
			ref.SetInlined(other.Inlined())
		}
	default:
		if ref.Type() != other.Type() || ref.Reference() != other.Reference() {
			ref.SetReference(other.Reference())
		}
	}

	return ref
}

// Replace replaces the data in the ReferencedData with the data from another ReferencedData.
//
// The underlying proto is not cloned, so modifying it will modify both this and the other ref.
func (ref *ReferencedData) Replace(other *ReferencedData) *ReferencedData {
	if !ref.Equal(other) {
		ref.rd = other.rd
		ref.refType = other.refType
		ref.modified = true
	}
	return ref
}

// ToProto converts the ReferencedData to a proto.
func (ref *ReferencedData) ToProto() *rdpb.ReferencedData {
	return proto.Clone(ref.rd).(*rdpb.ReferencedData)
}

// WalkUnique walks through a proto message, processing any unique ReferencedData it finds.
//
// Note that all inlined ReferencedData are considered unique.
//
// The function does not walk into Any protos.
//
// If the specified processing function modifies a ReferencedData, it replaces all instances of the
// original value in the message that is being walked.
//
// The input message may be mutated, and the processed message is returned.
func WalkUnique(msg proto.Message, f func(*ReferencedData) error) (proto.Message, error) {
	// Records references that were visited.
	type visited struct {
		ref      *ReferencedData
		returned proto.Message
	}
	visitedReferencedData := map[string]visited{}

	return walkmessages.Recursively(msg, func(msg proto.Message) (proto.Message, bool, error) {
		rd, ok, err := ToReferencedData(msg)
		if err != nil {
			return nil, false, fmt.Errorf("cannot convert message to ReferencedData: %w", err)
		}
		if !ok { // Not ReferencedData. Walk into the message.
			return nil, true, nil
		}

		ref := FromProto(rd)

		// If the reference has already been visited, just verify that the digest is the same.
		refKey := ref.Reference()
		if v, ok := visitedReferencedData[refKey]; ok {
			if v.ref.Digest() != "" && ref.Digest() != "" && v.ref.Digest() != ref.Digest() {
				return nil, false, fmt.Errorf("digest mismatch for reference %q: %q != %q", ref.Reference(), v.ref.Digest(), ref.Digest())
			}
			return v.returned, false, nil
		}

		if err := f(ref); err != nil {
			return nil, false, fmt.Errorf("cannot process ReferencedData: %w", err)
		}

		var msgOut proto.Message
		if ref.Modified() {
			if msgOut, err = fromReferencedData(ref.rd, msg); err != nil {
				return nil, false, fmt.Errorf("cannot convert ReferencedData to target message: %w", err)
			}
		}

		// For non-inlined ReferencedData, record the reference and the returned value.
		if refKey != "" {
			visitedReferencedData[refKey] = visited{
				ref:      ref,
				returned: msgOut,
			}
		}

		return msgOut, false, nil
	})
}

// ToReferencedData tries to convert the input message to ReferencedData, returning whether it
// succeeded.
func ToReferencedData(msg proto.Message) (*rdpb.ReferencedData, bool, error) {
	// Test whether message can be cast to ReferencedData.
	if rd, ok := msg.(*rdpb.ReferencedData); ok {
		return rd, true, nil
	}

	// Test whether message is a dynamicpb version of ReferencedData.
	rd := &rdpb.ReferencedData{}
	if msg.ProtoReflect().Descriptor().FullName() == rd.ProtoReflect().Descriptor().FullName() {
		// Convert to ReferencedData via Any.
		if msgAny, err := anypb.New(msg); err != nil {
			return nil, false, err
		} else if err := msgAny.UnmarshalTo(rd); err != nil {
			return nil, false, err
		}

		return rd, true, nil
	}

	return nil, false, nil
}

// fromReferencedData converts a ReferencedData to a target type.
func fromReferencedData(rd *rdpb.ReferencedData, target proto.Message) (proto.Message, error) {
	switch target.(type) {
	case *rdpb.ReferencedData:
		return rd, nil
	default:
		if rdAny, err := anypb.New(rd); err != nil {
			return nil, err
		} else if err := rdAny.UnmarshalTo(target); err != nil {
			return nil, err
		}
		return target, nil
	}
}

func parseReference(reference string) (ReferenceType, string, bool) {
	if strings.HasPrefix(reference, "intcas://") {
		return CASReferenceType, reference, false
	}
	if strings.HasPrefix(reference, "file://") {
		return FileReferenceType, reference[7:], true
	}
	return FileReferenceType, reference, false
}

// Reader represents ReferencedData and additional fields for reading the data it references.
type Reader struct {
	// Ref is a reference to the data.
	Ref *ReferencedData
	// Reader can be used to read the referenced data.
	Reader io.Reader
	// Size is the size of the referenced data, in bytes.
	Size int64
}

// ProcessOption is an option for Process.
type ProcessOption func(*Reader)

// WithReader specifies a custom io.Reader to provide for the referenced data along with the size
// of the reference data, in bytes.
func WithReader(r io.Reader, sz int64) ProcessOption {
	return func(rdr *Reader) {
		rdr.Reader = r
		rdr.Size = sz
	}
}

// Process processes a ReferencedData using the specified Processor.
//
// The reference is modified in place to update it to its processed form.
func Process(ctx context.Context, ref *ReferencedData, processor Processor, options ...ProcessOption) error {
	rdr := &Reader{
		Ref: ref,
	}
	for _, opt := range options {
		opt(rdr)
	}

	// Construct a default reader, if possible.
	if rdr.Reader == nil {
		switch ref.Type() {
		case FileReferenceType:
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
		case CASReferenceType:
			if processor.NeedsReaderFor(CASReferenceType) {
				return fmt.Errorf("CAS references cannot be read. got: %v", ref.Reference())
			}
		case InlinedReferenceType:
			rdr.Reader = bytes.NewReader(ref.Inlined())
			rdr.Size = int64(len(ref.Inlined()))
		default:
			return fmt.Errorf("unknown reference type: %d", ref.Type())
		}
	}

	return processor.Process(ctx, rdr)
}

// Processor is an interface for processing ReferencedData via a Reader.
type Processor interface {
	// NeedsReaderFor returns true if the processor needs an io.Reader for the given reference type.
	NeedsReaderFor(ReferenceType) bool

	// Process processes a ReferencedData. The processor should modify the ReferencedData in place.
	Process(context.Context, *Reader) error
}

type noOpProcessor struct{}

// NeedsReaderFor returns false for all reference types.
func (p *noOpProcessor) NeedsReaderFor(rt ReferenceType) bool {
	return false
}

// Process does not modify the given reference data.
func (p *noOpProcessor) Process(ctx context.Context, rdr *Reader) error {
	return nil
}

// NoOpProcessor returns a Processor that does nothing.
//
// This processor is only valid for dry runs.
func NoOpProcessor() Processor {
	return &noOpProcessor{}
}

type inlineProcessor struct{}

// NeedsReaderFor returns true for all reference types.
func (p *inlineProcessor) NeedsReaderFor(rt ReferenceType) bool {
	return true
}

// Process inlines the data referenced by the given ReferencedData.
func (p *inlineProcessor) Process(ctx context.Context, rdr *Reader) error {
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

// InlineProcessor returns a Processor that inlines the data referenced by the given ReferencedData.
//
// NOTE: This processor will read _all_ referenced data into memory and should not be used when
// large data may be referenced.
func InlineProcessor() Processor {
	return &inlineProcessor{}
}

type catalogProcessor struct {
	acClient  acpb.AssetCatalogClient
	chunkSize int
}

// CatalogProcessorOption is an option for CatalogProcessor.
type CatalogProcessorOption func(*catalogProcessor)

// WithChunkSize sets the chunk size for CatalogProcessor.
func WithChunkSize(size int) CatalogProcessorOption {
	return func(opts *catalogProcessor) {
		opts.chunkSize = size
	}
}

// NeedsReaderFor returns true for file references.
func (p *catalogProcessor) NeedsReaderFor(rt ReferenceType) bool {
	return rt == FileReferenceType
}

// Process prepares the given ReferencedData for inclusion in an Asset that will be released to the
// AssetCatalog.
func (p *catalogProcessor) Process(ctx context.Context, rdr *Reader) error {
	if rdr.Ref.Reference() != "" {
		log.Infof("Preparing reference %v", rdr.Ref.Reference())
	}

	stream, err := p.acClient.PrepareReferencedData(ctx)
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
	if rdr.Ref.Type() == FileReferenceType {
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
	rdr.Ref.Replace(FromProto(response.GetReferencedData()))

	return nil
}

// CatalogProcessor returns a Processor that prepares ReferencedData for inclusion in an Asset that
// will be released to the AssetCatalog.
func CatalogProcessor(client acpb.AssetCatalogClient, options ...CatalogProcessorOption) Processor {
	p := &catalogProcessor{
		acClient:  client,
		chunkSize: defaultChunkSize,
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

type solutionProcessor struct{}

func (p *solutionProcessor) NeedsReaderFor(rt ReferenceType) bool {
	return rt == FileReferenceType
}

func (p *solutionProcessor) Process(ctx context.Context, rdr *Reader) error {
	switch rdr.Ref.Type() {
	case FileReferenceType:
		// If the file is below the size threshold, inline it. Otherwise, upload it to CAS.
		if rdr.Size <= InlineReferenceFileSizeThresholdBytes {
			b, err := io.ReadAll(rdr.Reader)
			if err != nil {
				return fmt.Errorf("failed to read data file: %w", err)
			}
			rdr.Ref.SetInlined(b)
		} else {
			return fmt.Errorf("file upload is not supported: %v", rdr.Ref)
		}
	case CASReferenceType:
	case InlinedReferenceType:
		// Nothing to do.
	default:
		return fmt.Errorf("unknown referenced data: %v", rdr.Ref)
	}

	return nil
}

// SolutionProcessor returns a Processor that prepares ReferencedData for inclusion in a Solution.
//
// File references below a size threshold are inlined. Otherwise, they are uploaded to CAS.
func SolutionProcessor() Processor {
	return &solutionProcessor{}
}
