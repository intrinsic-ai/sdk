// Copyright 2023 Intrinsic Innovation LLC

// Package report provides a mechanism to treat errors as warnings.
package report

import (
	"errors"
)

// warningRule defines a function that returns true if an error should be
// collected as a warning rather than returned as an error.
type warningRule func(error) bool

// Report aggregates errors that are classified as warnings.
//
// Warning rules can be configured in the report. Errors matching these rules
// are accumulated as warnings in the report when an error is added.
// The caller can then inspect the report afterwards to check for warnings.
type Report struct {
	rules    []warningRule
	warnings []error
}

// Option is a functional option to configure the Report.
type Option func(*Report)

// AsWarning creates a rule that treats specific error values (e.g., sentinel errors)
// as warnings using errors.Is under the hood.
//
// Usage example:
//
//	r := report.New(report.AsWarning(ErrSentinel))
func AsWarning(target error) Option {
	return func(r *Report) {
		r.rules = append(r.rules, func(err error) bool {
			return errors.Is(err, target)
		})
	}
}

// AsWarningIf allows users to provide arbitrary logic to determine
// if an error should be treated as a warning.
//
// Usage example:
//
//	r := report.New(report.AsWarningIf(func(err error) bool {
//		return err.Error() == "warning"
//	}))
func AsWarningIf(fn func(error) bool) Option {
	return func(r *Report) {
		r.rules = append(r.rules, fn)
	}
}

// AsWarningIfType creates a rule that treats all errors of type T as warnings using errors.As.
//
// Usage example:
//
//	r := report.New(report.AsWarningIfType[*CustomError]())
func AsWarningIfType[T error]() Option {
	return func(r *Report) {
		r.rules = append(r.rules, func(err error) bool {
			var t T
			return errors.As(err, &t)
		})
	}
}

// AsWarningIfTypeAnd creates a generic rule that checks if an error is of a specific type T.
// If it is, it executes the provided closure. If the closure
// returns true, the error is collected as a warning.
//
// Usage example:
//
//	r := report.New(report.AsWarningIfTypeAnd(func(err *CustomError) bool {
//		return err.Code == 404
//	}))
func AsWarningIfTypeAnd[T error](fn func(T) bool) Option {
	return func(r *Report) {
		r.rules = append(r.rules, func(err error) bool {
			if target, ok := errors.AsType[T](err); ok {
				return fn(target)
			}
			return false
		})
	}
}

// New creates a new Report with the provided warning rules.
func New(opts ...Option) *Report {
	r := &Report{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Add evaluates an error. If the error is nil, it returns nil.
// If the error matches a defined warning rule, it collects the error
// into the internal warnings list and returns nil.
// If it does not match any warning rule, it returns the error directly.
func (r *Report) Add(err error) error {
	if err == nil {
		return nil
	}

	for _, rule := range r.rules {
		if rule(err) {
			// Matched a rule, collect as a warning and return nil.
			r.warnings = append(r.warnings, err)
			return nil
		}
	}

	// Error didn't match any rule; it is an error.
	return err
}

// Warnings returns the slice of all collected warning errors.
func (r *Report) Warnings() []error {
	return r.warnings
}
