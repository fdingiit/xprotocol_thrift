package sofabolt

import (
	"io"

	"gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/javaobject"
	"gitlab.alipay-inc.com/sofa-go/sofa-hessian-go/sofahessian"
)

type Request struct {
	// nolint
	noCopy    noCopy
	command   Command
	tbconn    javaobject.TBRemotingConnectionRequest
	tbconnbuf []byte
}

func (c *Request) GetTBRemotingConnection() *javaobject.TBRemotingConnectionRequest {
	return &c.tbconn
}

func (c *Request) Read(ro *ReadOption, r io.Reader) (int, error) {
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
		c.tbconn.Ctx = &javaobject.TBRemotingRequestContext{}
		return n, nil // discard the error
	}

	return n, nil
}

func (c *Request) Write(wo *WriteOption, b []byte) ([]byte, error) {
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

func (c *Request) Reset() {
	c.tbconn.Reset()
	c.tbconnbuf = c.tbconnbuf[:0]
	proto := c.command.GetProto()
	typ := c.command.GetType()
	cmdcode := c.command.GetCMDCode()
	c.command.Reset()
	c.command.SetProto(proto)
	c.command.SetType(typ)
	c.command.SetCMDCode(cmdcode)
}

func (c *Request) SetProto(p Proto) *Request       { c.command.SetProto(p); return c }
func (c *Request) SetVer1(v Version) *Request      { c.command.SetVer1(v); return c }
func (c *Request) SetType(t Type) *Request         { c.command.SetType(t); return c }
func (c *Request) SetCMDCode(cmd CMDCode) *Request { c.command.SetCMDCode(cmd); return c }
func (c *Request) SetVer2(v uint8) *Request        { c.command.SetVer2(v); return c }
func (c *Request) SetRequestID(id uint32) *Request {
	if c.command.proto == ProtoTBRemoting {
		c.tbconn.Ctx.ID = int64(id)
		return c
	}

	c.command.SetRequestID(id)
	return c
}
func (c *Request) SetCodec(codec Codec) *Request { c.command.SetCodec(codec); return c }
func (c *Request) SetSwitc(s uint8) *Request     { c.command.SetSwitc(s); return c }
func (c *Request) SetTimeout(t uint32) *Request  { c.command.SetTimeout(t); return c }
func (c *Request) SetStatus(s Status) *Request   { c.command.SetStatus(s); return c }
func (c *Request) SetConnection(connection []byte) *Request {
	c.command.SetConnection(connection)
	return c
}
func (c *Request) SetClass(class []byte) *Request       { c.command.SetClass(class); return c }
func (c *Request) SetClassString(class string) *Request { c.command.SetClassString(class); return c }
func (c *Request) SetContent(content []byte) *Request   { c.command.SetContent(content); return c }
func (c *Request) SetContentString(content string) *Request {
	c.command.SetContentString(content)
	return c
}
func (c *Request) CopyTo(d *Request) *Request { c.command.CopyTo(&d.command); return c }

func (c *Request) String() string        { return c.command.String() }
func (c *Request) GetConnection() []byte { return c.command.GetConnection() }
func (c *Request) GetProto() Proto       { return c.command.GetProto() }
func (c *Request) GetVer1() Version      { return c.command.GetVer1() }
func (c *Request) GetType() Type         { return c.command.GetType() }
func (c *Request) GetCMDCode() CMDCode   { return c.command.GetCMDCode() }
func (c *Request) GetVer2() uint8        { return c.command.GetVer2() }
func (c *Request) GetRequestID() uint32 {
	if c.command.proto == ProtoTBRemoting {
		return uint32(c.tbconn.Ctx.ID)
	}
	return c.command.GetRequestID()
}
func (c *Request) GetCodec() Codec        { return c.command.GetCodec() }
func (c *Request) GetSwitc() uint8        { return c.command.GetSwitc() }
func (c *Request) GetTimeout() uint32     { return c.command.GetTimeout() }
func (c *Request) GetStatus() Status      { return c.command.GetStatus() }
func (c *Request) GetClass() []byte       { return c.command.GetClass() }
func (c *Request) GetHeaders() *SimpleMap { return c.command.GetHeaders() }
func (c *Request) GetContent() []byte     { return c.command.GetContent() }
func (c *Request) Size() int              { return c.command.Size() }
