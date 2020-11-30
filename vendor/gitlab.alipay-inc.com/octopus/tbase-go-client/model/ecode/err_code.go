package ecode

type TBaseErrorCode int

const (
	ILLEGAL_ARGUMENTS   TBaseErrorCode = -3982
	SERVER_ERROR        TBaseErrorCode = -8888
	INTERNAL_ERROR      TBaseErrorCode = -2
	TIMEOUT             TBaseErrorCode = -1
	CONNECTION_ERROR    TBaseErrorCode = -6003
	HOTKEY              TBaseErrorCode = -6004
	MOVED_DATA          TBaseErrorCode = -3967
	READONLY            TBaseErrorCode = -3986
	LOADING             TBaseErrorCode = -18
	HANDSHAKE           TBaseErrorCode = -6002
	OVERLOAD            TBaseErrorCode = -6000
	KEY_LENGTH_OVERFLOW TBaseErrorCode = -5
)

const (
	ERROR_MOVED                          = "MOVED"
	ERROR_LOADING                        = "LOADING"
	ERROR_HANDSHAKE_ERR                  = "HANDSHAKE_ERR"
	ERROR_HOTKEY                         = "HOTKEY"
	ERROR_READONLY                       = "READONLY"
	ERROR_LAG                            = "LAG"
	ERROR_OVERLOAD                       = "OVERLOAD"
	ERROR_DATA_ERROR                     = "ERR"
	ERROR_QUEUE_FULL                     = "queue is full"
	ERROR_CONNECTION_NOT_USBALE          = "connection not usable"
	ERROR_COMMAND_TIMEOUT                = "command timeout"
	ERROR_NETWORK_UNREACHABLE            = "connect: network is unreachable"
	ERROR_CONNECTION_RESET_BY_PEER_ERROR = "reset by peer"
	ERROR_BROKEN_PIPE                    = "broken pipe"
	ERROR_USE_CLOSED_CONN                = "use of closed"
	ERROR_IO_TIMEOUT                     = "i/o timeout"
	ERROR_CONNECTION_REFUSED             = "connection refused"
)
