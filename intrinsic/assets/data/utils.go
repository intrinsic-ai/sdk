// Copyright 2023 Intrinsic Innovation LLC

// Package utils contains utils for working with Data assets.
package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"intrinsic/assets/idutils"
	"intrinsic/util/proto/registryutil"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	rdpb "intrinsic/assets/data/proto/v1/referenced_data_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
)

var (
	highwayHashKey = bytes.Repeat([]byte{0}, 32)
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

// ValidateDataAssetOptions contains options for a call to ValidateDataAsset.
type ValidateDataAssetOptions struct {
	// DisallowFileReferences indicates whether the Data asset must not contain ReferencedData with
	// file references.
	DisallowFileReferences bool
	// RequiresVersion indicates whether the Data asset must specify a version.
	RequiresVersion bool
}

// ValidateDataAssetOption is a functional option for ValidateDataAsset.
type ValidateDataAssetOption func(*ValidateDataAssetOptions)

// WithDisallowFileReferences sets the DisallowFileReferences option.
func WithDisallowFileReferences(disallowFileReferences bool) ValidateDataAssetOption {
	return func(opts *ValidateDataAssetOptions) {
		opts.DisallowFileReferences = disallowFileReferences
	}
}

// WithRequiresVersion sets the RequiresVersion option.
func WithRequiresVersion(requiresVersion bool) ValidateDataAssetOption {
	return func(opts *ValidateDataAssetOptions) {
		opts.RequiresVersion = requiresVersion
	}
}

// ValidateDataAsset validates a DataAsset.
//
// For required metadata, see
// https://github.com/intrinsic-ai/sdk/blob/main/intrinsic/assets/proto/metadata.proto.
func ValidateDataAsset(data *dapb.DataAsset, options ...ValidateDataAssetOption) error {
	opts := &ValidateDataAssetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	m := data.GetMetadata()

	if opts.RequiresVersion {
		if err := idutils.ValidateIDVersionProto(m.GetIdVersion()); err != nil {
			return fmt.Errorf("invalid id_version: %w", err)
		}
	} else if err := idutils.ValidateIDProto(m.GetIdVersion().GetId()); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if m.GetAssetType() != atpb.AssetType_ASSET_TYPE_DATA {
		return fmt.Errorf("asset_type must be ASSET_TYPE_DATA (got: %s)", m.GetAssetType())
	}
	if m.GetDisplayName() == "" {
		return fmt.Errorf("display_name must be specified")
	}
	if m.GetVendor().GetDisplayName() == "" {
		return fmt.Errorf("vendor.display_name must be specified")
	}

	if opts.DisallowFileReferences {
		if payload, err := ExtractPayload(data); err != nil {
			return err
		} else if _, err := WalkUniqueReferencedData(payload, func(ref *ReferencedDataExt) error {
			if ref.Type() == FileReferenceType {
				return fmt.Errorf("file references are not allowed (got: %q)", ref.Reference())
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// ValidateReferencedData validates a ReferencedData.
//
// Validation includes:
// - If specified, compare the digest against the referenced data.
func ValidateReferencedData(ref *ReferencedDataExt) error {
	// Validate the digest against the referenced data.
	if ref.Digest() != "" {
		// Get a reader for the data.
		var reader io.Reader
		switch ref.Type() {
		case FileReferenceType:
			file, err := os.Open(ref.Reference())
			if err != nil {
				return fmt.Errorf("cannot open referenced file %q: %w", ref.Reference(), err)
			}
			defer file.Close()
			reader = file
		case CASReferenceType:
			return fmt.Errorf("cannot validate digest of CAS reference %q", ref.Reference())
		case InlinedReferenceType:
			reader = bytes.NewReader(ref.Inlined())
		default:
			return fmt.Errorf("unknown reference type: %d", ref.Type())
		}

		// Test the digest.
		if parsed, err := ParseDigest(ref.Digest()); err != nil {
			return fmt.Errorf("cannot parse digest %q: %w", ref.Digest(), err)
		} else if gotDigest, err := Digest(reader, WithAlgorithm(parsed.Algorithm)); err != nil {
			return fmt.Errorf("cannot compute digest: %w", err)
		} else if gotDigest != parsed.Digest {
			return fmt.Errorf("digest mismatch: got %q, want %q", gotDigest, parsed.Digest)
		}
	}

	return nil
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

	return walkProtoMessages(msg, func(msg proto.Message) (proto.Message, bool, error) {
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

// FProcessMessage is a function that takes a message and processes it.
//
// Used in walkProtoMessages. If a non-nil message is returned, the return value replaces the
// original value in the message that is being walked.
//
// The function also indicated whether to enter into the message recursively.
type fProcessMessage func(proto.Message) (proto.Message, bool, error)

// walkProtoMessages walks through a proto message, executing a function for each message it finds.
//
// The function returns whether to enter into the message recursively.
//
// The input message may be mutated, and the processed message is returned.
func walkProtoMessages(msg proto.Message, f fProcessMessage) (proto.Message, error) {
	msgOut, shouldEnter, err := f(msg)
	if err != nil {
		return nil, err
	}
	if msgOut == nil { // No changes made. Use original message.
		msgOut = msg
	}
	if !shouldEnter {
		return msgOut, nil
	}

	msgOutR := msgOut.ProtoReflect()
	for i := 0; i < msgOutR.Descriptor().Fields().Len(); i++ {
		field := msgOutR.Descriptor().Fields().Get(i)

		// Skip unspecified fields.
		if !msgOutR.Has(field) {
			continue
		}

		// Skip non-message/group types.
		if field.Kind() != protoreflect.MessageKind && field.Kind() != protoreflect.GroupKind {
			continue
		}

		valueR := msgOutR.Get(field)
		if field.IsList() { // Walk through lists.
			for i := 0; i < valueR.List().Len(); i++ {
				msgItem := valueR.List().Get(i).Message().Interface()
				if msgItemOut, err := walkProtoMessages(msgItem, f); err != nil {
					return nil, err
				} else if msgItemOut != nil { // Item was changed; update the parent.
					valueR.List().Set(i, protoreflect.ValueOfMessage(msgItemOut.ProtoReflect()))
				}
			}
		} else if field.IsMap() { // Walk through maps.
			var err error
			valueR.Map().Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
				msgItem := value.Message().Interface()
				var msgItemOut proto.Message
				if msgItemOut, err = walkProtoMessages(msgItem, f); err != nil {
					return false
				} else if msgItemOut != nil { // Item was changed; update the parent.
					valueR.Map().Set(key, protoreflect.ValueOfMessage(msgItemOut.ProtoReflect()))
				}
				return true
			})
			if err != nil {
				return nil, err
			}
		} else if valueROut, err := walkProtoMessages(valueR.Message().Interface(), f); err != nil {
			return nil, err
		} else if valueROut != nil { // Field was changed; update the parent.
			msgOutR.Set(field, protoreflect.ValueOfMessage(valueROut.ProtoReflect()))
		}
	}

	return msgOut, nil
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
