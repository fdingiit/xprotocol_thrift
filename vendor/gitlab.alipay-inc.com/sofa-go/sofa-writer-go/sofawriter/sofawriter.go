package sofawriter

import (
	"errors"
	"io"
	"io/ioutil"

	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/asyncwriter"
	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/dsn"
	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/rollingwriter"
	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/rsyslogwriter"
	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/testwriter"
)

type Writer struct {
	dsn     *dsn.DSN
	dsnlist *dsn.DSNList
	w       io.Writer
}

func New(writers ...io.Writer) *Writer {
	return &Writer{
		w: io.MultiWriter(writers...),
	}
}

func (w *Writer) GetDSNList() *dsn.DSNList {
	return w.dsnlist
}

func (w *Writer) GetDSN() *dsn.DSN {
	return w.dsn
}

func (w *Writer) Close() error {
	if rw, ok := w.w.(io.Closer); ok {
		return rw.Close()
	}
	return nil
}

func (w *Writer) Write(p []byte) (int, error) { return w.w.Write(p) }

func NewFromDSNString(d string) (*Writer, error) {
	n, err := dsn.NewDSN(d)
	if err != nil {
		return nil, err
	}

	return NewFromDSN(n)
}

func NewFromDSN(d *dsn.DSN) (*Writer, error) {
	w, err := newWriter(d)
	if err != nil {
		return nil, err
	}

	return &Writer{
		dsn: d,
		w:   w,
	}, nil
}

func NewFromDSNList(dsnlist *dsn.DSNList) (*Writer, error) {
	var w io.Writer
	dl := dsnlist.Get()
	if len(dl) == 0 {
		w = ioutil.Discard
	} else if len(dl) > 1 {
		writers := make([]io.Writer, 0, len(dl))
		for i := range dl {
			nw, err := newWriter(dl[i])
			if err != nil {
				return nil, err
			}
			writers = append(writers, nw)
		}
		w = io.MultiWriter(writers...)

	} else {
		nw, err := newWriter(dl[0])
		if err != nil {
			return nil, err
		}
		w = nw
	}

	return &Writer{
		dsnlist: dsnlist,
		w:       w,
	}, nil
}

func newWriter(d *dsn.DSN) (io.Writer, error) {
	var w io.Writer
	switch d.GetScheme() {
	case "", "file", "unix":
		option := rollingwriter.NewOption()
		option.SetMaxSize(int(dsn.ParseInt64(d.GetQuery(dsn.MaxSizeKey), 0)))
		option.SetMaxAge(int(dsn.ParseInt64(d.GetQuery(dsn.MaxAgeKey), 0)))
		option.SetMaxBackups(int(dsn.ParseInt64(d.GetQuery(dsn.MaxBackupsKey), 0)))
		rw := rollingwriter.New(d.GetPath(), option)
		w = rw

	case "rsyslog", "syslog":
		option := rsyslogwriter.NewOption()
		option.SetServer(d.GetHost())
		if ak := d.GetQuery(dsn.RsyslogAppNameKey); ak != "" {
			option.SetAppname(ak)
		}

		option.SetSeverity(dsn.ParseSeverity(d.GetQuery(dsn.RsyslogSeverityKey), rsyslogwriter.INFO))
		option.SetFacility(dsn.ParseFacility(d.GetQuery(dsn.RsyslogFacilityKey), rsyslogwriter.USER))

		rw, err := rsyslogwriter.New(option)
		if err != nil {
			return nil, err
		}
		w = rw
	case "test":
		tw, _, err := testwriter.New(d)
		if err != nil {
			return nil, err
		}

		w = tw

	default:
		return nil, errors.New("unknown scheme type")
	}

	async := d.GetQuery(dsn.AsyncKey)
	if len(async) > 0 {
		option := asyncwriter.NewOption().
			SetBatch(int(
				dsn.ParseInt64(d.GetQuery(dsn.AsyncBatchKey), 0)),
			).
			SetFlushInterval(dsn.ParseDuration(d.GetQuery(dsn.AsyncFlushIntervalKey), 0))

		if dsn.ParseBool(d.GetQuery(dsn.AsyncBlockKey), false) {
			option.AllowBlockForever()
		}
		var err error
		w, err = asyncwriter.New(w, asyncwriter.WithAsyncWriterOption(option))
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}
