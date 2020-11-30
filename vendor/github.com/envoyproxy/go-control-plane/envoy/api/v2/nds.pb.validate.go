// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/api/v2/nds.proto

package envoy_api_v2

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
var _nds_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on NodeSentry with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *NodeSentry) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Name

	// no validation rules for Namespace

	// no validation rules for InstanceIp

	// no validation rules for NodeName

	// no validation rules for NodeIp

	// no validation rules for ClusterName

	return nil
}

// NodeSentryValidationError is the validation error returned by
// NodeSentry.Validate if the designated constraints aren't met.
type NodeSentryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NodeSentryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NodeSentryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NodeSentryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NodeSentryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NodeSentryValidationError) ErrorName() string { return "NodeSentryValidationError" }

// Error satisfies the builtin error interface
func (e NodeSentryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNodeSentry.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NodeSentryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NodeSentryValidationError{}
