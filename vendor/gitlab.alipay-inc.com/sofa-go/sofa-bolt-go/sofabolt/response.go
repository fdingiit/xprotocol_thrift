package sofabolt

import (
	"io"

	"gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/javaobject"
	"gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/sofahessian"
)

type Response struct {
	// nolint
	noCopy    noCopy
	command   Command
	tbconn    javaobject.TBRemotingConnectionResponse
	tbconnbuf []byte
}

func (c *Response) GetTBRemotingConnection() *javaobject.TBRemotingConnectionResponse {
	return &c.tbconn
}

func (c *Response) Derive(b *Request) {
	c.command.proto = b.command.proto
	c.command.ver1 = b.command.ver1
	c.command.typ = TypeBOLTResponse
	c.command.cmdcode = CMDCodeBOLTResponse
	c.command.ver2 = b.command.ver2
	c.command.rid = b.command.rid
	c.command.codec = b.command.codec
	c.command.switc = b.command.switc
	c.command.timeout = b.command.timeout
	c.command.status = b.command.status
}

func (c *Response) Read(ro *ReadOption, r io.Reader) (int, error) {
	n, err := c.command.Read(ro, r)
	if err != nil {
		return n, err
	}

	if c.command.proto != ProtoTBRemoting {
		return n, err
	}

	dctx := sofahessian.AcquireHessianDecodeContext().
		SetClassRegistry(&trregistry)
	bbr := sofahessian.AcquireBytesBufioReader(c.command.GetConnection())
	err = sofahessian.DecodeObjectToHessian3V2(dctx, bbr.GetBufioReader(), &c.tbconn)
	sofahessian.ReleaseBytesBufioReader(bbr)
	sofahessian.ReleaseHessianDecodeContext(dctx)

	if err != nil {
		c.tbconn.Ctx = &javaobject.TBRemotingConnectionResponseContext{}
		return n, nil // discard the error
	}

	return n, nil
}

func (c *Response) Write(wo *WriteOption, b []byte) ([]byte, error) {
	if c.command.proto != ProtoTBRemoting {
		return c.command.Write(wo, b)
	}

	var err error

	ectx := sofahessian.AcquireHessianEncodeContext()
	c.tbconnbuf, err = sofahessian.EncodeObjectToHessian3V2(ectx, c.tbconnbuf[:0], c.tbconn)
	sofahessian.ReleaseHessianEncodeContext(ectx)
	if err != nil {
		return b, err
	}

	c.command.SetConnection(c.tbconnbuf)
	return c.command.Write(wo, b)
}

func (c *Response) Reset() {
	proto := c.command.GetProto()
	typ := c.command.GetType()
	cmdcode := c.command.GetCMDCode()
	c.command.Reset()
	c.command.SetProto(proto)
	c.command.SetType(typ)
	c.command.SetCMDCode(cmdcode)
}

func (c *Response) SetProto(p Proto) *Response       { c.command.SetProto(p); return c }
func (c *Response) SetVer1(v Version) *Response      { c.command.SetVer1(v); return c }
func (c *Response) SetType(t Type) *Response         { c.command.SetType(t); return c }
func (c *Response) SetCMDCode(cmd CMDCode) *Response { c.command.SetCMDCode(cmd); return c }
func (c *Response) SetVer2(v uint8) *Response        { c.command.SetVer2(v); return c }
func (c *Response) SetRequestID(id uint32) *Response {
	if c.command.proto == ProtoTBRemoting {
		c.tbconn.Ctx.ID = int64(id)
	}

	c.command.SetRequestID(id)
	return c
}
func (c *Response) SetCodec(codec Codec) *Response { c.command.SetCodec(codec); return c }
func (c *Response) SetSwitc(s uint8) *Response     { c.command.SetSwitc(s); return c }
func (c *Response) SetTimeout(t uint32) *Response  { c.command.SetTimeout(t); return c }
func (c *Response) SetStatus(s Status) *Response   { c.command.SetStatus(s); return c }
func (c *Response) SetConnection(content []byte) *Response {
	c.command.SetConnection(content)
	return c
}
func (c *Response) SetClass(class []byte) *Response { c.command.SetClass(class); return c }
func (c *Response) SetClassString(class string) *Response {
	c.command.SetClassString(class)
	return c
}
func (c *Response) SetContent(content []byte) *Response { c.command.SetContent(content); return c }
func (c *Response) SetContentString(content string) *Response {
	c.command.SetContentString(content)
	return c
}
func (c *Response) CopyTo(d *Response) *Response { c.command.CopyTo(&d.command); return c }

func (c *Response) String() string      { return c.command.String() }
func (c *Response) GetProto() Proto     { return c.command.GetProto() }
func (c *Response) GetVer1() Version    { return c.command.GetVer1() }
func (c *Response) GetType() Type       { return c.command.GetType() }
func (c *Response) GetCMDCode() CMDCode { return c.command.GetCMDCode() }

func (c *Response) GetRequestID() uint32 {
	if c.command.proto == ProtoTBRemoting {
		return uint32(c.tbconn.Ctx.ID)
	}
	return c.command.GetRequestID()
}

func (c *Response) GetVer2() uint8         { return c.command.GetVer2() }
func (c *Response) GetCodec() Codec        { return c.command.GetCodec() }
func (c *Response) GetSwitc() uint8        { return c.command.GetSwitc() }
func (c *Response) GetTimeout() uint32     { return c.command.GetTimeout() }
func (c *Response) GetStatus() Status      { return c.command.GetStatus() }
func (c *Response) GetClass() []byte       { return c.command.GetClass() }
func (c *Response) GetHeaders() *SimpleMap { return c.command.GetHeaders() }
func (c *Response) GetContent() []byte     { return c.command.GetContent() }
func (c *Response) Size() int              { return c.command.Size() }
