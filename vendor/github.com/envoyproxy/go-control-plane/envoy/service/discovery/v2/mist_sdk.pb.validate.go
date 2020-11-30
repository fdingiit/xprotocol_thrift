// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/service/discovery/v2/mist_sdk.proto

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
var _mist_sdk_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on MistSDKRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *MistSDKRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for AppName

	// no validation rules for AppContainerIp

	// no validation rules for Command

	for idx, item := range m.GetParameters() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MistSDKRequestValidationError{
					field:  fmt.Sprintf("Parameters[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MistSDKRequestValidationError is the validation error returned by
// MistSDKRequest.Validate if the designated constraints aren't met.
type MistSDKRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MistSDKRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MistSDKRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MistSDKRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MistSDKRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MistSDKRequestValidationError) ErrorName() string { return "MistSDKRequestValidationError" }

// Error satisfies the builtin error interface
func (e MistSDKRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMistSDKRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MistSDKRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MistSDKRequestValidationError{}

// Validate checks the field values on MistSDKResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *MistSDKResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for AppName

	// no validation rules for AppContainerIp

	// no validation rules for Command

	for idx, item := range m.GetResponse() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MistSDKResponseValidationError{
					field:  fmt.Sprintf("Response[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MistSDKResponseValidationError is the validation error returned by
// MistSDKResponse.Validate if the designated constraints aren't met.
type MistSDKResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MistSDKResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MistSDKResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MistSDKResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MistSDKResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MistSDKResponseValidationError) ErrorName() string { return "MistSDKResponseValidationError" }

// Error satisfies the builtin error interface
func (e MistSDKResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMistSDKResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MistSDKResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MistSDKResponseValidationError{}
