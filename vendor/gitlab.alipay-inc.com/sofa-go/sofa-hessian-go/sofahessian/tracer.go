package sofahessian

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Tracer interface {
	OnTraceStart(event string)
	OnTraceStop(event string)
}

type DummyTracer struct{}

func NewDummyTracer() *DummyTracer { return &DummyTracer{} }

func (t *DummyTracer) OnTraceStart(event string) { _ = event }

func (t *DummyTracer) OnTraceStop(event string) { _ = event }

type StdoutTracer struct {
	depth int
}

func NewStdoutTracer() *StdoutTracer { return &StdoutTracer{} }

func (t *StdoutTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Printf("%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *StdoutTracer) OnTraceStop(event string) {
	fmt.Printf("%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}

type StderrTracer struct {
	depth int
}

func NewStderrTracer() *StderrTracer { return &StderrTracer{} }

func (t *StderrTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Fprintf(os.Stderr, "%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *StderrTracer) OnTraceStop(event string) {
	fmt.Fprintf(os.Stderr, "%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}

type WriterTracer struct {
	w     io.Writer
	depth int
}

func NewWriterTracer(w io.Writer) *WriterTracer {
	return &WriterTracer{
		w:     w,
		depth: 0,
	}
}

func (t *WriterTracer) OnTraceStart(event string) {
	t.depth++
	fmt.Fprintf(t.w, "%s<start %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
}

func (t *WriterTracer) OnTraceStop(event string) {
	fmt.Fprintf(t.w, "%s<stop %s>\n", strings.Repeat(" ", (t.depth-1)*4), event)
	t.depth--
}
