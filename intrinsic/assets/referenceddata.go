// Copyright 2023 Intrinsic Innovation LLC

// Package referenceddata provides utils for working with ReferencedData.
package referenceddata

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"intrinsic/util/proto/walkmessages"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

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

// ReferencedDataProcessor is a function that takes a ReferencedData as input and processes it.
//
// Used in WalkUnique. If the ReferencedData is modified, it replaces the original value in the
// message that is being walked.
type ReferencedDataProcessor func(*ReferencedData) error

// WalkUnique walks through a proto message, processing any unique ReferencedData it finds.
//
// Note that all inlined ReferencedData are considered unique.
//
// The function does not walk into Any protos.
//
// If the specified processor modifies a ReferencedData, it replaces all instances of the original
// value in the message that is being walked.
//
// The input message may be mutated, and the processed message is returned.
func WalkUnique(msg proto.Message, f ReferencedDataProcessor) (proto.Message, error) {
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
