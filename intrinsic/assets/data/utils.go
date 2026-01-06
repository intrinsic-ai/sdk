// Copyright 2023 Intrinsic Innovation LLC

// Package utils contains utils for working with Data assets.
package utils

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/walkmessages"

	"google.golang.org/protobuf/proto"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

var highwayHashKey = bytes.Repeat([]byte{0}, 32)

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

// ReferencedDataExt represents a ReferencedData with extended functionality.
type ReferencedDataExt struct {
	modified bool
	rd       *rdpb.ReferencedData
	refType  ReferenceType
}

// NewReferencedDataExt creates a new ReferencedDataExt.
func NewReferencedDataExt(rd *rdpb.ReferencedData) *ReferencedDataExt {
	ref := &ReferencedDataExt{
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
func (ref *ReferencedDataExt) Digest() string {
	return ref.rd.GetDigest()
}

// Inlined returns the inlined data in the ReferencedData.
func (ref *ReferencedDataExt) Inlined() []byte {
	return ref.rd.GetInlined()
}

// Modified returns whether the ReferencedData has been modified.
func (ref *ReferencedDataExt) Modified() bool {
	return ref.modified
}

// Reference returns the reference in the ReferencedData.
func (ref *ReferencedDataExt) Reference() string {
	return ref.rd.GetReference()
}

// SourceProject returns the source project in the ReferencedData.
func (ref *ReferencedDataExt) SourceProject() string {
	return ref.rd.GetSourceProject()
}

// Type returns the type of reference in the ReferencedData.
func (ref *ReferencedDataExt) Type() ReferenceType {
	return ref.refType
}

// Equal returns whether the two ReferencedDataExt are equal.
func (ref *ReferencedDataExt) Equal(other *ReferencedDataExt) bool {
	return proto.Equal(ref.rd, other.rd)
}

// SetBaseDir sets the base directory for relative file references.
func (ref *ReferencedDataExt) SetBaseDir(baseDir string) *ReferencedDataExt {
	if ref.Type() == FileReferenceType && !filepath.IsAbs(ref.rd.GetReference()) {
		ref.rd.Data = &rdpb.ReferencedData_Reference{
			Reference: filepath.Join(baseDir, ref.rd.GetReference()),
		}
		ref.modified = true
	}
	return ref
}

// SetDigest sets the digest in the ReferencedData.
func (ref *ReferencedDataExt) SetDigest(digest string) *ReferencedDataExt {
	ref.rd.Digest = digest
	ref.modified = true
	return ref
}

// SetSourceProject sets the source project in the ReferencedData.
func (ref *ReferencedDataExt) SetSourceProject(sourceProject string) *ReferencedDataExt {
	ref.rd.SourceProject = proto.String(sourceProject)
	ref.modified = true
	return ref
}

// SetInlined sets the inlined data in the ReferencedData.
func (ref *ReferencedDataExt) SetInlined(inlined []byte) *ReferencedDataExt {
	ref.rd.Data = &rdpb.ReferencedData_Inlined{
		Inlined: inlined,
	}
	ref.refType = InlinedReferenceType
	ref.modified = true
	return ref
}

// SetReference sets the reference in the ReferencedData.
func (ref *ReferencedDataExt) SetReference(reference string) *ReferencedDataExt {
	ref.refType, reference, _ = parseReference(reference)
	ref.rd.Data = &rdpb.ReferencedData_Reference{
		Reference: reference,
	}
	ref.modified = true
	return ref
}

// Copy returns a shallow copy of the ReferencedDataExt.
//
// The modified flag is reset.
func (ref *ReferencedDataExt) Copy() *ReferencedDataExt {
	return NewReferencedDataExt(ref.rd)
}

// Merge merges the data from another ReferencedDataExt into this one.
func (ref *ReferencedDataExt) Merge(other *ReferencedDataExt) *ReferencedDataExt {
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

// Replace replaces the data in the ReferencedDataExt with the data from another ReferencedDataExt.
//
// The underlying proto is not cloned, so modifying it will modify both this and the other ref.
func (ref *ReferencedDataExt) Replace(other *ReferencedDataExt) *ReferencedDataExt {
	if !ref.Equal(other) {
		ref.rd = other.rd
		ref.refType = other.refType
		ref.modified = true
	}
	return ref
}

// ToProto converts the ReferencedDataExt to a ReferencedData proto.
func (ref *ReferencedDataExt) ToProto() *rdpb.ReferencedData {
	return proto.Clone(ref.rd).(*rdpb.ReferencedData)
}

// ExtractPayload extracts a Data asset payload to a dynamicpb analog of the target payload type.
func ExtractPayload(data *dapb.DataAsset) (proto.Message, error) {
	payloadAny := data.GetData()

	types, err := registryutil.NewTypesFromFileDescriptorSet(data.GetFileDescriptorSet())
	if err != nil {
		return nil, fmt.Errorf("cannot populate registry: %w", err)
	}

	msgType, err := types.FindMessageByName(payloadAny.MessageName())
	if err != nil {
		return nil, fmt.Errorf("cannot find message for Data payload: %w", err)
	}

	payload := msgType.New().Interface()
	if err := payloadAny.UnmarshalTo(payload); err != nil {
		return nil, fmt.Errorf("cannot unmarshal data payload: %w", err)
	}

	return payload, nil
}

// ReferencedDataProcessor is a function that takes a ReferencedDataExt as input and processes it.
//
// Used in WalkUniqueReferencedData. If the ReferencedData is modified, it replaces the original value in
// the message that is being walked.
type ReferencedDataProcessor func(*ReferencedDataExt) error

// WalkUniqueReferencedData walks through a proto message, processing any unique ReferencedData it
// finds. (Note that all inlined ReferencedData are considered unique.)
//
// The function does not walk into Any protos.
//
// If the specified processor modifies a ReferencedData, it replaces all instances of the original
// value in the message that is being walked.
//
// The input message may be mutated, and the processed message is returned.
func WalkUniqueReferencedData(msg proto.Message, f ReferencedDataProcessor) (proto.Message, error) {
	// Records references that were visited.
	type visited struct {
		ref      *ReferencedDataExt
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

		ref := NewReferencedDataExt(rd)

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

// HashAlgorithm is an enum for the supported digest hash algorithms.
type HashAlgorithm string

const (
	// HighwayHash128 is the HighwayHash-128 algorithm.
	HighwayHash128 HashAlgorithm = "highwayhash128"
)

// DigestOptions contains options for a call to Digest.
type DigestOptions struct {
	// Algorithm is the hashing algorithm to use.
	Algorithm HashAlgorithm
}

// DigestOption is a functional option for DigestFrom*.
type DigestOption func(*DigestOptions)

// WithAlgorithm sets the Algorithm option.
func WithAlgorithm(algorithm HashAlgorithm) DigestOption {
	return func(opts *DigestOptions) {
		opts.Algorithm = algorithm
	}
}

// Digest creates a digest for the specified data.
func Digest(reader io.Reader, options ...DigestOption) (string, error) {
	opts := DigestOptions{
		Algorithm: HighwayHash128,
	}
	for _, opt := range options {
		opt(&opts)
	}

	switch opts.Algorithm {
	default:
		return "", fmt.Errorf("unknown hash algorithm: %v", opts.Algorithm)
	}
}

// ParsedDigest is a parsed digest.
type ParsedDigest struct {
	Algorithm HashAlgorithm
	Digest    string
	Hash      string
}

// ParseDigest parses a digest.
func ParseDigest(digest string) (*ParsedDigest, error) {
	parts := strings.Split(digest, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid digest: %q", digest)
	}
	return &ParsedDigest{
		Algorithm: HashAlgorithm(parts[0]),
		Digest:    digest,
		Hash:      parts[1],
	}, nil
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
