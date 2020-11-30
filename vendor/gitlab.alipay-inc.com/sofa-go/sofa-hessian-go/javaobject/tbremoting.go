package javaobject

type TBRemotingRequestContext struct {
	ID   int64                        `hessian:"id"`
	This *TBRemotingConnectionRequest `hessian:"this$0"`
}

func (c *TBRemotingRequestContext) Reset() {
	c.ID = 0
	c.This = nil
}

func (c *TBRemotingRequestContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest$RequestContext"
}

type TBRemotingConnectionRequest struct {
	Ctx                *TBRemotingRequestContext `hessian:"ctx"`
	FromAppKey         string                    `hessian:"fromAppKey"`
	ToAppKey           string                    `hessian:"toAppKey"`
	EncryptedToken     []byte                    `hessian:"encryptedToken"`
	ApplicationRequest interface{}               `hessian:"-"`
}

func (c *TBRemotingConnectionRequest) Reset() {
	if c.Ctx != nil {
		c.Ctx.Reset()
	}
	c.FromAppKey = ""
	c.ToAppKey = ""
	c.EncryptedToken = c.EncryptedToken[:0]
	c.ApplicationRequest = nil
}

func (c *TBRemotingConnectionRequest) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionRequest"
}

type TBRemotingConnectionResponseContext struct {
	ID int64 `hessian:"id"`
}

func (c *TBRemotingConnectionResponseContext) Reset() {
	c.ID = 0
}

func (c *TBRemotingConnectionResponseContext) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse$ResponseContext"
}

type TBRemotingConnectionResponse struct {
	Ctx        *TBRemotingConnectionResponseContext
	Host       string `hessian:"host"`
	Result     int32  `hessian:"result"`
	ErrorMsg   string `hessian:"errorMsg"`
	ErrorStack string `hessian:"errorStack"`
	FromAppKey string `hessian:"fromAppKey"`
	ToAppKey   string `hessian:"toAppKey"`
}

func (c *TBRemotingConnectionResponse) Reset() {
	if c.Ctx != nil {
		c.Ctx.Reset()
	}
}

func (c *TBRemotingConnectionResponse) GetJavaClassName() string {
	return "com.taobao.remoting.impl.ConnectionResponse"
}
