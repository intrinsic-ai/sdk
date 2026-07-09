// Copyright 2023 Intrinsic Innovation LLC

// Package errors provides structured validation errors and a warning reporting mechanism.
package errors

import (
	"errors"
	"fmt"
)

// Code represents a unique identifier for a validation error.
type Code uint32

// Error codes defined here can also be used as status codes for ExtendedStatus errors.
//
// A starting error code of 2000 is arbitrarily chosen here. Codes greater than 10000
// are reserved for errors not originating as part of the system (or platform).
const (
	// CodeErrUnknownDependency is used when a node has a dependency on a non-existent node.
	CodeErrUnknownDependency Code = iota + 2000

	// CodeErrInterfaceMismatch is used when a node has a dependency on an Asset that does not satisfy the dependency.
	CodeErrInterfaceMismatch
)

// Error is a structured validation error that implements the error interface.
type Error interface {
	error
	// Code returns the error code.
	//
	// This is used by [Report] to specify validation errors that should instead be logged
	// as warnings.
	Code() Code
}

// defaultError is a basic implementation of the Error interface that has an
// error code and wraps an underlying error.
type defaultError struct {
	code Code
	err  error
}

// Error returns the error message.
func (e *defaultError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// Code returns the error code.
func (e *defaultError) Code() Code {
	return e.code
}

// Unwrap returns the underlying wrapped error.
func (e *defaultError) Unwrap() error {
	return e.err
}

// Is returns true if target is a defaultError with the same code and message.
func (e *defaultError) Is(target error) bool {
	t, ok := target.(*defaultError)
	if !ok {
		return false
	}
	if e.code != t.code {
		return false
	}
	if errors.Is(e.err, t.err) {
		return true
	}
	if e.err == nil || t.err == nil {
		return false
	}
	return e.err.Error() == t.err.Error()
}

// NewError creates a new [Error] with the specified code and error.
func NewError(code Code, err error) Error {
	return &defaultError{code: code, err: err}
}

// Report aggregates validation warnings.
//
// Error codes that should be treated as warnings can be configured in the report.
// Errors matching these codes will not cause validation to fail, but will be
// accumulated as warnings in the report. The caller can then inspect the report
// afterwards to check for warnings.
type Report struct {
	warningCodes map[Code]bool
	warnings     []Error
}

// ReportOption is an option for creating a validation [Report].
type ReportOption func(*Report)

// WithWarningCodes sets the error codes that should be treated as warnings in the [Report].
func WithWarningCodes(codes ...Code) ReportOption {
	return func(r *Report) {
		for _, c := range codes {
			r.warningCodes[c] = true
		}
	}
}

// NewReport creates a new validation report where errors with the given error
// codes are treated as warnings.
func NewReport(opts ...ReportOption) *Report {
	r := &Report{
		warningCodes: make(map[Code]bool),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Add integrates an error into the validation [Report].
//
//   - All errors recognized as warnings are appended to the report warnings and nil is returned.
//   - All other errors are returned directly to fail-fast.
func (r *Report) Add(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := errors.AsType[Error](err); ok {
		if r.warningCodes[e.Code()] {
			r.warnings = append(r.warnings, e)
			return nil
		}
	}
	return err
}

// Warnings returns all accumulated warnings.
func (r *Report) Warnings() []Error {
	return r.warnings
}

type unknownDependencyError struct {
	source        string
	target        string
	sourceAssetID string
	err           error
	code          Code
}

// UnknownDependencyErrorDetails contains additional optional details that can be provided for
// an [Error] specific to unknown dependency validation.
type UnknownDependencyErrorDetails struct {
	// SourceAssetID is the Asset ID of the dependent Asset or Asset Instance.
	SourceAssetID string
}

func (e *unknownDependencyError) Error() string {
	return e.err.Error()
}

func (e *unknownDependencyError) Code() Code {
	return e.code
}

func (e *unknownDependencyError) Unwrap() error {
	return e.err
}

// NewUnknownDependencyError creates a new [Error] for unknown dependency validation.
//
// `source` is the identifier of the dependent Asset or Asset Instance.
// For Services and HardwareDevices, this is the instance name. For Processes, this is
// the Asset ID.
//
// `target` is the identifier of the Asset or Asset Instance being depended on.
// If the depended on entity is an Asset Instance of a Service, HardwareDevice or
// SceneObject, this is the instance name. If it is a DataAsset, this is the Asset ID.
func NewUnknownDependencyError(source, target string, err error, opts *UnknownDependencyErrorDetails) Error {
	return &unknownDependencyError{
		source:        source,
		target:        target,
		sourceAssetID: opts.SourceAssetID,
		err:           err,
		code:          CodeErrUnknownDependency,
	}
}

type interfaceMismatchError struct {
	source        string
	target        string
	sourceAssetID string
	targetAssetID string
	err           error
	code          Code
}

// InterfaceMismatchErrorDetails contains additional optional details that can be provided for
// an [Error] specific to interface mismatch validation.
type InterfaceMismatchErrorDetails struct {
	// SourceAssetID is the Asset ID of the dependent Asset or Asset Instance.
	SourceAssetID string
	// TargetAssetID is the Asset ID of the Asset or Asset Instance being depended on.
	TargetAssetID string
}

func (e *interfaceMismatchError) Code() Code {
	return e.code
}

func (e *interfaceMismatchError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("%q has an interface that is incompatible with %q", e.source, e.target)
	}
	return fmt.Sprintf("%q has an interface that is incompatible with %q: %v", e.source, e.target, e.err)
}

// NewInterfaceMismatchError creates a new [Error] for interface mismatch.
//
// `source` is the identifier of the dependent Asset or Asset Instance.
// For Services and HardwareDevices, this is the instance name. For Processes, this is
// the Asset ID.
//
// `target` is the identifier of the Asset or Asset Instance being depended on.
// If the depended on entity is an Asset Instance of a Service, HardwareDevice or
// SceneObject, this is the instance name. If it is a DataAsset, this is the Asset ID.
func NewInterfaceMismatchError(source, target string, err error, opts *InterfaceMismatchErrorDetails) Error {
	return &interfaceMismatchError{
		source:        source,
		target:        target,
		err:           err,
		code:          CodeErrInterfaceMismatch,
		sourceAssetID: opts.SourceAssetID,
		targetAssetID: opts.TargetAssetID,
	}
}
