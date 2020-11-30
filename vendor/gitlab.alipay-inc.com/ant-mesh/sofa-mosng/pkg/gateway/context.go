package gateway

import (
	"context"
)

type GatewayContext struct {
	uniqueId     string
	request      GatewayRequest
	response     GatewayResponse
	service      Service
	attributes   map[string]interface{}
	currentIndex int
}

func NewGatewayContext(ctx context.Context, request GatewayRequest) *GatewayContext {
	uniqueId := request.RequestId()
	return &GatewayContext{
		uniqueId:   uniqueId,
		request:    request,
		response:   nil,
		attributes: make(map[string]interface{}),
	}
}

func BuildGatewayContext() *GatewayContext {
	return &GatewayContext{
		attributes: make(map[string]interface{}),
	}
}

func GetGatewayContext(ctx context.Context) *GatewayContext {
	return ctx.Value(GATEWAY_CONTEXT_NAME).(*GatewayContext)
}

func GetGatewayContextAttr(context context.Context) map[string]interface{} {
	return GetGatewayContext(context).attributes
}

func GetGatewayContextAttrValue(context context.Context, key string) interface{} {
	return GetGatewayContextAttr(context)[key]
}

func SetGatewayContextAttrValue(context context.Context, key string, value interface{}) {
	GetGatewayContextAttr(context)[key] = value
}

func (g *GatewayContext) UniqueId() string {
	return g.uniqueId
}

func (g *GatewayContext) SetUniqueId(uniqueId string) {
	g.uniqueId = uniqueId
}

func (g *GatewayContext) Request() GatewayRequest {
	return g.request
}

func (g *GatewayContext) SetRequest(request GatewayRequest) {
	g.request = request
}

func (g *GatewayContext) Response() GatewayResponse {
	return g.response
}

func (g *GatewayContext) SetResponse(response GatewayResponse) {
	g.response = response
}

func (g *GatewayContext) Service() Service {
	return g.service
}

func (g *GatewayContext) SetService(service Service) {
	g.service = service
}

func (g *GatewayContext) GetAttribute(key string) interface{} {
	return g.attributes[key]
}

func (g *GatewayContext) SetAttribute(key string, value interface{}) {
	if g.attributes == nil {
		g.attributes = make(map[string]interface{})
	}
	g.attributes[key] = value
}

func (g *GatewayContext) CurrentIndex() int {
	return g.currentIndex
}

func (g *GatewayContext) SetIndex(index int) {
	g.currentIndex = index
}
