package testwriter

import (
	"fmt"
	"strconv"
	"sync"

	"gitlab.alipay-inc.com/sofa-go/sofa-writer-go/dsn"
)

var writers sync.Map

type TestWriter struct {
	sync.Mutex
	path    string
	discard bool
	b       []byte
}

func DelAll() {
	writers.Range(func(k, v interface{}) bool {
		writers.Delete(k)
		return true
	})
}

func Get(path string) (*TestWriter, bool) {
	i, ok := writers.Load(path)
	if !ok {
		return nil, ok
	}

	return i.(*TestWriter), true
}

func Del(path string) {
	writers.Delete(path)
}

func New(d *dsn.DSN) (*TestWriter, string, error) {
	dd := d.GetQuery("discard")
	if dd == "" {
		dd = "false"
	}

	discard, err := strconv.ParseBool(dd)
	if err != nil {
		return nil, "", err
	}

	dd = d.GetQuery("trace")
	if dd == "" {
		dd = "false"
	}

	l := d.GetQuery("level")

	trace, err := strconv.ParseBool(dd)
	if err != nil {
		return nil, "", err
	}

	_, ok := writers.Load(d.GetPath())
	if ok {
		return nil, "", fmt.Errorf("duplicate test writer: %s", d.GetPath())
	}

	w := &TestWriter{
		path:    d.GetPath(),
		discard: discard,
	}

	if trace {
		writers.Store(d.GetPath(), w)
	}

	return w, l, nil
}

func (tw *TestWriter) GetPath() string {
	tw.Lock()
	defer tw.Unlock()
	return tw.path
}

func (tw *TestWriter) GetBuffer() []byte {
	tw.Lock()
	defer tw.Unlock()
	return tw.b
}

func (tw *TestWriter) Write(p []byte) (int, error) {
	tw.Lock()
	defer tw.Unlock()
	if tw.discard {
		return len(p), nil
	}
	tw.b = append(tw.b, p...)
	return len(p), nil
}
