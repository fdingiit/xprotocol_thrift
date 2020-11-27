package forwarder

import (
	"net"
	"net/http"
	"time"
)

func HTTPTransport() *http.Transport {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
			// Enable TCP keepalives to detect broken connections
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 5,
		// to avoid connections sitting idle in the pool indefinitely
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return transport
}

type Logger interface {
	Debugf(template string, args ...interface{})
	Debug(args ...interface{})
	Infof(template string, args ...interface{})
	Info(args ...interface{})
	Warnf(template string, args ...interface{})
	Warn(args ...interface{})
	Errorf(template string, args ...interface{})
	Error(args ...interface{})
	Panicf(template string, args ...interface{})
	Panic(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatal(args ...interface{})
}
