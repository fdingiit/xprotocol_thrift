package dsn

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/rsyslogwriter"
)

const (
	MaxSizeKey            = "maxsize"
	MaxBackupsKey         = "maxbackups"
	MaxAgeKey             = "maxage"
	CompressKey           = "compress"
	RsyslogAppNameKey     = "rsyslog_appname"
	RsyslogSeverityKey    = "rsyslog_severity"
	RsyslogFacilityKey    = "rsyslog_facility"
	AsyncKey              = "async"
	AsyncBatchKey         = "async_batch"
	AsyncBlockKey         = "async_block"
	AsyncFlushIntervalKey = "async_flush_interval"
)

type DSNList struct {
	d []*DSN
}

type DSN struct {
	u *url.URL
}

func NewDSN(dsn string) (*DSN, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	return &DSN{u: u}, nil
}

func (d *DSN) String() string {
	return d.u.String()
}

func (d *DSN) GetScheme() string {
	return d.u.Scheme
}

func (d *DSN) GetPath() string {
	return d.u.Path
}

func (d *DSN) GetHost() string {
	return d.u.Host
}

func (d *DSN) GetQuery(key string) string {
	return d.u.Query().Get(key)
}

func NewDSNList(dsnlist string, sep string) (*DSNList, error) {
	dd := make([]*DSN, 0, 10)
	s := strings.Split(dsnlist, sep)
	for i := range s {
		d, err := NewDSN(s[i])
		if err != nil {
			return nil, err
		}
		dd = append(dd, d)
	}

	return &DSNList{
		d: dd,
	}, nil
}

func (d *DSNList) Get() []*DSN {
	return d.d
}

func ParseBool(s string, def bool) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return b
}

func ParseInt64(s string, def int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return i64
}

func ParseDuration(s string, d time.Duration) time.Duration {
	t, err := time.ParseDuration(s)
	if err != nil {
		return d
	}
	return t
}

func ParseSeverity(s string, d rsyslogwriter.Severity) rsyslogwriter.Severity {
	severity, err := rsyslogwriter.ParseSeverity(s)
	if err != nil {
		return d
	}
	return severity
}

func ParseFacility(s string, d rsyslogwriter.Facility) rsyslogwriter.Facility {
	facility, err := rsyslogwriter.ParseFacility(s)
	if err != nil {
		return d
	}
	return facility
}
