package sofabolt

import (
	"io"
	"sync"
	"time"

	bufiorw "github.com/detailyang/bufio-rw-go"
	"gitlab.alipay-inc.com/sofa-go/sofa-syncpool-go/bytespool"
)

var (
	b32Pool = sync.Pool{
		New: func() interface{} {
			var p [32]byte
			return &p
		},
	}

	bpool = bytespool.NewPool()

	requestPool = sync.Pool{
		New: func() interface{} {
			return &Request{}
		},
	}

	responsePool = sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}

	brPool    = sync.Pool{}
	bwPool    = sync.Pool{}
	timerPool sync.Pool
	ictxPool  = sync.Pool{
		New: func() interface{} {
			return &InvokeContext{
				doneCh: make(chan struct{}),
			}
		},
	}
)

func acquireInvokeContext(req *Request, res *Response, timeout time.Duration) *InvokeContext {
	ictx, ok := ictxPool.Get().(*InvokeContext)
	if !ok {
		panic("failed to type casting")
	}

	if len(ictx.errCh) == 1 || ictx.errCh == nil {
		ictx.errCh = make(chan error, 1) // sanity: allocate new error channel
	}

	ictx.req = req
	ictx.res = res
	ictx.created = time.Now()
	ictx.timeout = timeout

	return ictx
}

func releaseInvokeContext(ictx *InvokeContext) {
	ictxPool.Put(ictx)
}

func acquireB32() *[32]byte {
	return b32Pool.Get().(*[32]byte)
}

func releaseB32(p *[32]byte) {
	b32Pool.Put(p)
}

func AcquireRequest() *Request {
	req, ok := requestPool.Get().(*Request)
	if !ok {
		panic("failed to type casting")
	}
	req.SetProto(ProtoBOLTV1)
	req.SetCMDCode(CMDCodeBOLTRequest)
	req.SetType(TypeBOLTRequest)
	req.SetCodec(CodecHessian2)
	return req
}

func ReleaseRequest(di *Request) {
	di.Reset()
	requestPool.Put(di)
}

func AcquireResponse() *Response {
	res, ok := responsePool.Get().(*Response)
	if !ok {
		panic("failed to type casting")
	}
	res.SetProto(ProtoBOLTV1)
	res.SetCMDCode(CMDCodeBOLTResponse)
	res.SetType(TypeBOLTResponse)
	res.SetCodec(CodecHessian2)
	return res
}

func ReleaseResponse(di *Response) {
	di.Reset()
	responsePool.Put(di)
}

func acquireBufioWriter(w io.Writer) *bufiorw.Writer {
	i := bwPool.Get()
	if i == nil {
		return bufiorw.NewWriterSize(w, 8192)
	}

	bw, ok := i.(*bufiorw.Writer)
	if !ok {
		panic("failed to type casting")
	}
	bw.Reset(w)

	return bw
}

func releaseBufioWriter(bw *bufiorw.Writer) {
	bw.Reset(nil)
	bwPool.Put(bw)
}

func acquireBufioReader(r io.Reader) *bufiorw.Reader {
	i := brPool.Get()
	if i == nil {
		return bufiorw.NewReaderSize(r, 8192)
	}

	br, ok := i.(*bufiorw.Reader)
	if !ok {
		panic("failed to type casting")
	}
	br.Reset(r)

	return br
}

func releaseBufioReader(br *bufiorw.Reader) {
	br.Reset(nil)
	brPool.Put(br)
}

func initTimer(t *time.Timer, timeout time.Duration) *time.Timer {
	if t == nil {
		return time.NewTimer(timeout)
	}
	if t.Reset(timeout) {
		panic("BUG: active timer trapped into initTimer()")
	}
	return t
}

func stopTimer(t *time.Timer) {
	if !t.Stop() {
		// Collect possibly added time from the channel
		// if timer has been stopped and nobody collected its' value.
		select {
		case <-t.C:
		default:
		}
	}
}

func acquireTimer(timeout time.Duration) *time.Timer {
	v := timerPool.Get()
	if v == nil {
		return time.NewTimer(timeout)
	}
	t := v.(*time.Timer)
	initTimer(t, timeout)
	return t
}

func releaseTimer(t *time.Timer) {
	stopTimer(t)
	timerPool.Put(t)
}

func acquireBytes() *[]byte  { return bpool.Acquire() }
func releaseBytes(d *[]byte) { bpool.Release(d) }
