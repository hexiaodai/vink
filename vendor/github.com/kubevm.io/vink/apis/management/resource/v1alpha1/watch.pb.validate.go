// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: management/resource/v1alpha1/watch.proto

package v1alpha1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"

	types "github.com/kubevm.io/vink/apis/types"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort

	_ = types.ResourceType(0)
)

// Validate checks the field values on WatchRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WatchRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WatchRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WatchRequestMultiError, or
// nil if none found.
func (m *WatchRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *WatchRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ResourceType

	if all {
		switch v := interface{}(m.GetOptions()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WatchRequestValidationError{
					field:  "Options",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WatchRequestValidationError{
					field:  "Options",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOptions()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WatchRequestValidationError{
				field:  "Options",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return WatchRequestMultiError(errors)
	}

	return nil
}

// WatchRequestMultiError is an error wrapping multiple validation errors
// returned by WatchRequest.ValidateAll() if the designated constraints aren't met.
type WatchRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WatchRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WatchRequestMultiError) AllErrors() []error { return m }

// WatchRequestValidationError is the validation error returned by
// WatchRequest.Validate if the designated constraints aren't met.
type WatchRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WatchRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WatchRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WatchRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WatchRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WatchRequestValidationError) ErrorName() string { return "WatchRequestValidationError" }

// Error satisfies the builtin error interface
func (e WatchRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWatchRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WatchRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WatchRequestValidationError{}

// Validate checks the field values on WatchResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WatchResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WatchResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WatchResponseMultiError, or
// nil if none found.
func (m *WatchResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *WatchResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for EventType

	if len(errors) > 0 {
		return WatchResponseMultiError(errors)
	}

	return nil
}

// WatchResponseMultiError is an error wrapping multiple validation errors
// returned by WatchResponse.ValidateAll() if the designated constraints
// aren't met.
type WatchResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WatchResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WatchResponseMultiError) AllErrors() []error { return m }

// WatchResponseValidationError is the validation error returned by
// WatchResponse.Validate if the designated constraints aren't met.
type WatchResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WatchResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WatchResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WatchResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WatchResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WatchResponseValidationError) ErrorName() string { return "WatchResponseValidationError" }

// Error satisfies the builtin error interface
func (e WatchResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWatchResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WatchResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WatchResponseValidationError{}

// Validate checks the field values on WatchOptions with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WatchOptions) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WatchOptions with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WatchOptionsMultiError, or
// nil if none found.
func (m *WatchOptions) ValidateAll() error {
	return m.validate(true)
}

func (m *WatchOptions) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetFieldSelectorGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WatchOptionsValidationError{
					field:  "FieldSelectorGroup",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WatchOptionsValidationError{
					field:  "FieldSelectorGroup",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFieldSelectorGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WatchOptionsValidationError{
				field:  "FieldSelectorGroup",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return WatchOptionsMultiError(errors)
	}

	return nil
}

// WatchOptionsMultiError is an error wrapping multiple validation errors
// returned by WatchOptions.ValidateAll() if the designated constraints aren't met.
type WatchOptionsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WatchOptionsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WatchOptionsMultiError) AllErrors() []error { return m }

// WatchOptionsValidationError is the validation error returned by
// WatchOptions.Validate if the designated constraints aren't met.
type WatchOptionsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WatchOptionsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WatchOptionsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WatchOptionsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WatchOptionsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WatchOptionsValidationError) ErrorName() string { return "WatchOptionsValidationError" }

// Error satisfies the builtin error interface
func (e WatchOptionsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWatchOptions.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WatchOptionsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WatchOptionsValidationError{}
