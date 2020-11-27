package model

import (
	"bufio"
	"fmt"
	"reflect"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	"gitlab.alipay-inc.com/octopus/radix/resp/resp2"
	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	convertor_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/convertor.helper"
)

var (
	arrayPrefix = []byte{'*'} // 42
	errPrefix   = []byte{'-'} //45
)

type CasValue struct {
	Value   interface{}
	Version int
}

// UnmarshalRESP unmarshal response from bufio.Reader to CasValue
//
// The prefix must be arrayPrefix and the length of the array must be 2
// The first one of the array is the Value and the second is the Version
func (c *CasValue) UnmarshalRESP(br *bufio.Reader) error {
	b, err := br.Peek(1)
	if err != nil {
		return err
	}
	prefix := b[0]

	_, err = br.Discard(1)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("discard 1 byte error, error: %v", err)
		return err
	}

	b, err = convertor_helper.BufferedBytesDelim(br)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("extract bytes out of \r\n error, error: %v", err)
		return err
	}

	if prefix == errPrefix[0] {
		return error2.NewTBaseClientDataError(string(b))
	}

	if prefix != arrayPrefix[0] {
		tbase_log.TBaseLogger.Errorf("expect array prefix for response, actually is %v", string(prefix))
		return error2.NewTBaseClientInternalError(fmt.Sprintf("expect array prefix for response, actually is %v", string(prefix)))
	}

	l, err := convertor_helper.ByteArrayToInt64(b)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("extract array length error, error: %v", err)
		return err
	} else if l != 2 {
		tbase_log.TBaseLogger.Errorf("extract array length is 2, actually is %v", l)
	}

	ai := resp2.Any{I: &c.Value}
	if err := ai.UnmarshalRESP(br); err != nil {
		tbase_log.TBaseLogger.Errorf("unmarshal resp to %v error, error: %v", reflect.TypeOf(c.Value), err)
		return err
	}

	ai = resp2.Any{I: &c.Version}
	if err := ai.UnmarshalRESP(br); err != nil {
		tbase_log.TBaseLogger.Errorf("unmarshal resp to %v error, error: %v", reflect.TypeOf(c.Value), err)
		return err
	}

	return nil
}
