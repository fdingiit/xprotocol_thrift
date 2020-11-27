package sofadrm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

const (
	HMACSHA256           = "HmacSHA256"
	DefaultCacheInterval = 5 * time.Minute
)

func errstr(err error) string {
	if err == nil {
		return "nil"
	}

	return err.Error()
}

func protostr(message fmt.Stringer) string {
	if message == nil {
		return "nil"
	}

	msg := message.String()
	if len(msg) > 1024 {
		return fmt.Sprintf("%s...more(%d)", msg[:1024], len(msg))
	}
	return msg
}

func merror(err error, errs ...error) error {
	noerr := true
	if err == nil {
		for _, err := range errs {
			if err != nil {
				noerr = false
			}
		}
	} else {
		noerr = false
	}

	if noerr {
		return nil
	}

	merr := multierror.Append(err, errs...)
	merr.ErrorFormat = merrorfmt
	return merr
}

func merrorfmt(errs []error) string {
	var b strings.Builder
	for i, err := range errs {
		_, _ = b.WriteString("#")
		_, _ = b.WriteString(strconv.Itoa(i))
		_, _ = b.WriteString(err.Error())
	}
	return b.String()
}

var ErrNoAvailableIPv4Addrs = errors.New("IPv4 address not available")

// AvailableIPv4Addrs returns a list of IPv4 addresses bound to ifaceName.
func AvailableIPv4Addrs(ifaceName string) ([]string, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ip := ipnet.IP.To4(); ip != nil {
			ips = append(ips, ip.String())
		}
	}
	if len(ips) == 0 {
		return nil, ErrNoAvailableIPv4Addrs
	}
	return ips, nil
}

func doHMACSHA256Base64(key, plaintext string) string {
	hash := hmac.New(sha256.New, []byte(key))
	// nolint
	hash.Write([]byte(plaintext))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func getSignatureMap(accesskey, secretkey, instanceID string) map[string]string {
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	cacheTime := int64(DefaultCacheInterval / time.Millisecond)
	timestamp := strconv.FormatInt(ts/cacheTime*cacheTime, 10)
	plaintext := fmt.Sprintf("%s%s", instanceID, timestamp)
	return map[string]string{
		"!AccessKey": accesskey,
		"!Algothrim": HMACSHA256,
		"!Signature": doHMACSHA256Base64(secretkey, plaintext),
		"!Timestamp": timestamp,
	}
}
