// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/service/discovery/v2/mist_svid_sign.proto

package envoy_service_discovery_v2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _mist_svid_sign_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on MistSvidSignRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MistSvidSignRequest) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetExt() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MistSvidSignRequestValidationError{
					field:  fmt.Sprintf("Ext[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MistSvidSignRequestValidationError is the validation error returned by
// MistSvidSignRequest.Validate if the designated constraints aren't met.
type MistSvidSignRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MistSvidSignRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MistSvidSignRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MistSvidSignRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MistSvidSignRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MistSvidSignRequestValidationError) ErrorName() string {
	return "MistSvidSignRequestValidationError"
}

// Error satisfies the builtin error interface
func (e MistSvidSignRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMistSvidSignRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MistSvidSignRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MistSvidSignRequestValidationError{}

// Validate checks the field values on MistSvidSignResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MistSvidSignResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Svid

	for idx, item := range m.GetExt() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MistSvidSignResponseValidationError{
					field:  fmt.Sprintf("Ext[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MistSvidSignResponseValidationError is the validation error returned by
// MistSvidSignResponse.Validate if the designated constraints aren't met.
type MistSvidSignResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MistSvidSignResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MistSvidSignResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MistSvidSignResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MistSvidSignResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MistSvidSignResponseValidationError) ErrorName() string {
	return "MistSvidSignResponseValidationError"
}

// Error satisfies the builtin error interface
func (e MistSvidSignResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMistSvidSignResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MistSvidSignResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MistSvidSignResponseValidationError{}
