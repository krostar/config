// Package trivialerr helps to define and use trivial (not critical) errors.
package trivialerr

import (
	"errors"
	"fmt"
)

// TrivialError defines the way a trivial error is defined.
type TrivialError interface {
	error
	IsTrivial() bool
}

// New creates a new trivial error using fmt.Errorf.
func New(msg string, args ...interface{}) error {
	return WrapIf(false, fmt.Errorf(msg, args...))
}

// Wrap wraps an existing error to make it trivial.
func Wrap(err error) error {
	return WrapIf(false, err)
}

// WrapIf wraps an error, as the Wrap function, but
// only if strict is false. Otherwise its returning
// the original error.
func WrapIf(strict bool, err error) error {
	if err == nil {
		return nil
	}
	if strict {
		return err
	}
	return trivialError{err: err}
}

// IsTrivial returns true if error implements IsTrivial,
// and if err.IsTrivial is true.
func IsTrivial(err error) bool {
	var trivErr TrivialError
	if errors.As(err, &trivErr) && trivErr.IsTrivial() {
		return true
	}
	return false
}

type trivialError struct {
	err error
}

// Error implements error.
func (te trivialError) Error() string { return te.err.Error() }

// Unwrap implements errors.Unwrap.
func (te trivialError) Unwrap() error { return te.err }

// IsTrivial implements IsTrivial.
func (te trivialError) IsTrivial() bool { return true }
