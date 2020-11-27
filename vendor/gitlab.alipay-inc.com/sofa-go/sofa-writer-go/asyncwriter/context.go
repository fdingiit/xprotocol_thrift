package asyncwriter

import (
	"io"
	"net"
	"time"
)

type context struct {
	option     *Option
	writer     io.Writer
	buffer     []byte
	buffers    [][]byte
	buffersp   []*[]byte
	conn       net.Conn
	netbuffers net.Buffers
}

func (ctx *context) reset() {
	ctx.option = nil
	ctx.writer = nil
	ctx.buffer = ctx.buffer[:0]
	ctx.buffers = ctx.buffers[:0]
	ctx.buffersp = ctx.buffersp[:0]
	ctx.conn = nil
	ctx.netbuffers = ctx.netbuffers[:0]
}

func (ctx *context) Flush() (int, error) {
	if len(ctx.buffer) > 0 {
		if ctx.option.timeout != 0 {
			if ctx.conn != nil { // can set timeout
				if err := ctx.conn.SetWriteDeadline(time.Now().Add(ctx.option.timeout)); err != nil {
					return 0, err
				}
			}
		}

		n, err := ctx.writer.Write(ctx.buffer)
		if err != nil {
			return 0, err
		}
		ctx.buffer = ctx.buffer[:0]
		return n, nil
	}
	return 0, nil
}
