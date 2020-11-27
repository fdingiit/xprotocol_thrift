package error

import (
	"fmt"
	"io"
	"net"
	"strings"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/model/ecode"
	error_code "gitlab.alipay-inc.com/octopus/tbase-go-client/model/ecode"
)

type TBaseError struct {
	errorCode error_code.TBaseErrorCode
	Message   string
}

func (tError *TBaseError) Error() string {
	return fmt.Sprintf("error code: %v, Message: %v", tError.errorCode, tError.Message)
}

type TBaseClientTimeoutError struct {
	TBaseError
}

type TBaseClientDataError struct {
	TBaseError
}

type TBaseClientInternalError struct {
	TBaseError
}

type TBaseClientConnectionError struct {
	TBaseError
}

type TBaseClientHotKeyDataError struct {
	TBaseError
}

type TBaseClientMovedDataError struct {
	TBaseError
}

type TBaseClientReadOnlyError struct {
	TBaseError
}

type TBaseClientLoadingError struct {
	TBaseError
}

type TBaseClientHandShakeError struct {
	TBaseError
}

type TBaseClientLagError struct {
	TBaseError
}

type TBaseClientOverLoadError struct {
	TBaseError
}

type TBaseClientIllegalArgumentsError struct {
	TBaseError
}

type TBaseClientKeyLengthOverFlowError struct {
	TBaseError
}

func NewTBaseClientTimeoutError(message string) *TBaseClientTimeoutError {
	return &TBaseClientTimeoutError{TBaseError: TBaseError{errorCode: error_code.TIMEOUT, Message: message}}
}

func NewTBaseClientDataError(message string) *TBaseClientDataError {
	return &TBaseClientDataError{TBaseError: TBaseError{errorCode: error_code.SERVER_ERROR, Message: message}}
}

func NewTBaseClientInternalError(message string) *TBaseClientInternalError {
	return &TBaseClientInternalError{TBaseError: TBaseError{errorCode: error_code.INTERNAL_ERROR, Message: message}}
}

func NewTBaseClientConnectionError(message string) *TBaseClientConnectionError {
	return &TBaseClientConnectionError{TBaseError: TBaseError{errorCode: error_code.CONNECTION_ERROR, Message: message}}
}

func NewTBaseClientHotKeyDataError(message string) *TBaseClientHotKeyDataError {
	return &TBaseClientHotKeyDataError{TBaseError: TBaseError{errorCode: error_code.HOTKEY, Message: message}}
}

func NewTBaseClientMovedDataError(message string) *TBaseClientMovedDataError {
	return &TBaseClientMovedDataError{TBaseError: TBaseError{errorCode: error_code.MOVED_DATA, Message: message}}
}

func NewTBaseClientReadOnlyError(message string) *TBaseClientReadOnlyError {
	return &TBaseClientReadOnlyError{TBaseError: TBaseError{errorCode: error_code.READONLY, Message: message}}
}

func NewTBaseClientLoadingError(message string) *TBaseClientLoadingError {
	return &TBaseClientLoadingError{TBaseError: TBaseError{errorCode: error_code.LOADING, Message: message}}
}

func NewTBaseClientHandShakeError(message string) *TBaseClientHandShakeError {
	return &TBaseClientHandShakeError{TBaseError: TBaseError{errorCode: error_code.HANDSHAKE, Message: message}}
}

func NewTBaseClientLagError(message string) *TBaseClientLagError {
	return &TBaseClientLagError{TBaseError: TBaseError{errorCode: error_code.TIMEOUT, Message: message}}
}

func NewTBaseClientOverloadError(message string) *TBaseClientOverLoadError {
	return &TBaseClientOverLoadError{TBaseError: TBaseError{errorCode: error_code.OVERLOAD, Message: message}}
}

func NewTBaseClientIllegalArgumentsError(message string) *TBaseClientIllegalArgumentsError {
	return &TBaseClientIllegalArgumentsError{TBaseError: TBaseError{errorCode: error_code.ILLEGAL_ARGUMENTS, Message: message}}
}

func NewTBaseClientKeyLengthOverFlowError(message string) *TBaseClientKeyLengthOverFlowError {
	return &TBaseClientKeyLengthOverFlowError{TBaseError: TBaseError{errorCode: error_code.KEY_LENGTH_OVERFLOW, Message: message}}
}

func IsConnectionErr(err error) bool {
	if netError, _ := err.(net.Error); netError != nil {
		return true
	}

	if strings.Contains(err.Error(), io.EOF.Error()) ||
		strings.Contains(err.Error(), io.ErrShortBuffer.Error()) ||
		strings.Contains(err.Error(), io.ErrShortWrite.Error()) ||
		strings.Contains(err.Error(), io.ErrNoProgress.Error()) ||
		strings.Contains(err.Error(), io.ErrUnexpectedEOF.Error()) ||
		strings.Contains(err.Error(), io.ErrClosedPipe.Error()) {
		return true
	}

	if strings.Contains(err.Error(), ecode.ERROR_NETWORK_UNREACHABLE) ||
		strings.Contains(err.Error(), ecode.ERROR_CONNECTION_RESET_BY_PEER_ERROR) ||
		strings.Contains(err.Error(), ecode.ERROR_BROKEN_PIPE) ||
		strings.Contains(err.Error(), ecode.ERROR_USE_CLOSED_CONN) ||
		strings.Contains(err.Error(), ecode.ERROR_IO_TIMEOUT) ||
		strings.Contains(err.Error(), ecode.ERROR_CONNECTION_REFUSED) {
		return true
	}

	return false
}

var COMMAND_TIMEOUT_ERROR = NewTBaseClientTimeoutError(error_code.ERROR_COMMAND_TIMEOUT)
var CONNECTION_NOT_USABLE_ERROR = NewTBaseClientConnectionError(error_code.ERROR_CONNECTION_NOT_USBALE)
var QUEUE_FULL_ERROR = NewTBaseClientInternalError(error_code.ERROR_QUEUE_FULL)
