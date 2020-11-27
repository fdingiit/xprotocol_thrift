package extension

import (
	"github.com/yuin/gopher-lua"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/filter"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

type LuaScriptEngine struct {
}

func LoadFilterExt(fe v1.FilterExtension) {
	// Register lua filter
	filter.Register(fe.Name, &LuaFilterFactory{fe.File})
}

// Lua Filter Factory
type LuaFilterFactory struct {
	file string
}

func (f *LuaFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	codeToShare, err := CompileLua(f.file)
	if err == nil {
		//todo
	}
	return &LuaFilter{base: filter, filterCode: codeToShare}
}

// Lua Proxy Filter
type LuaFilter struct {
	base *v1.Filter

	filterCode *lua.FunctionProto
}

func (f *LuaFilter) Name() string {
	return f.base.Name
}

func (f *LuaFilter) Priority() int64 {
	return f.base.Priority
}

func (f *LuaFilter) ShouldIn(ctx types.Context) bool {
	return true
}

func (f *LuaFilter) InBound(ctx types.Context) (types.FilterStatus, error) {

	ret := RunFilter("InBound", ctx, f.filterCode)

	switch ret {
	case 1:
		return types.Success, nil
	default:
		return types.Error, nil
	}

}

func (f *LuaFilter) ShouldOut(ctx types.Context) bool {
	return true
}

func (f *LuaFilter) OutBound(ctx types.Context) (types.FilterStatus, error) {

	ret := RunFilter("OutBound", ctx, f.filterCode)

	switch ret {
	case 1:
		return types.Success, nil
	default:
		return types.Error, nil
	}
}

func RunFilter(method string, ctx types.Context, fp *lua.FunctionProto) int {

	// get LState from pool
	L := luaPool.Get()
	defer luaPool.Put(L)

	SetContext(L, ctx)

	// get filter
	if err := DoCompiledFile(L, fp); err != nil {
		//todo
	}

	luaFilter := L.ToTable(-1)

	// new lua context
	luaCtx := newLuaContext(L)

	L.SetGlobal("mosng", luaCtx)

	filterFunc := luaFilter.RawGet(lua.LString(method)).(*lua.LFunction)
	if err := L.CallByParam(lua.P{
		Fn:      filterFunc,
		NRet:    1,
		Protect: true,
	}, lua.LNil); err != nil {
		panic(err)
	}
	return L.ToInt(-1)
}

func newLuaContext(L *lua.LState) *lua.LTable {
	// lua request
	luaRequest := L.NewTable()
	luaRequest.RawSetString("getHeaders", L.NewFunction(luaGetRequestHeaders))
	luaRequest.RawSetString("getHeader", L.NewFunction(luaGetRequestHeader))
	luaRequest.RawSetString("addHeader", L.NewFunction(luaAddRequestHeaders))
	luaRequest.RawSetString("delHeader", L.NewFunction(luaDelRequestHeaders))
	luaRequest.RawSetString("getData", L.NewFunction(luaGetRequestDataBuf))
	luaRequest.RawSetString("setData", L.NewFunction(luaSetRequestDataBuf))

	// lua response
	luaResponse := L.NewTable()
	luaResponse.RawSetString("getHeaders", L.NewFunction(luaGetResponseHeaders))
	luaResponse.RawSetString("getHeader", L.NewFunction(luaGetResponseHeader))
	luaResponse.RawSetString("addHeader", L.NewFunction(luaAddResponseHeaders))
	luaResponse.RawSetString("delHeader", L.NewFunction(luaDelResponseHeaders))
	luaResponse.RawSetString("getData", L.NewFunction(luaGetResponseDataBuf))
	luaResponse.RawSetString("setData", L.NewFunction(luaSetResponseDataBuf))

	// lua ctx
	luaCtx := L.NewTable()
	luaCtx.RawSetString("request", luaRequest)
	luaCtx.RawSetString("response", luaResponse)

	return luaCtx
}

// lua to go function

// request function
//GetHeaders() types.HeaderMap
//todo mutli value in one key
func luaGetRequestHeaders(L *lua.LState) int {
	headers := L.NewTable()
	ctx := GetContext(L)
	ctx.Request().GetHeaders().Range(func(key, value string) bool {
		headers.RawSetString(key, lua.LString(value))
		return true
	})
	L.Push(headers)
	return 1
}

//GetHeader(key string) string
//todo mutli value in one key
//todo param check
func luaGetRequestHeader(L *lua.LState) int {
	key := L.ToString(-1)

	ctx := GetContext(L)
	if value, ok := ctx.Request().GetHeaders().Get(key); ok {
		L.Push(lua.LString(value))
	} else {
		L.Push(lua.LNil)
	}

	return 1
}

//SetHeader(key, value string)
func luaAddRequestHeaders(L *lua.LState) int {
	if L.GetTop() != 2 {
		// "expecting exactly 2 arguments");
		//return err
	}
	key := L.ToString(-2)
	value := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Request().SetHeader(key, value)
	return 0
}

//DelHeader(key string)
func luaDelRequestHeaders(L *lua.LState) int {
	key := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Response().DelHeader(key)
	return 0
}

//GetDataBuf() types.IoBuffer
func luaGetRequestDataBuf(L *lua.LState) int {
	ctx := GetContext(L)
	b := ctx.Request().GetDataBytes()
	L.Push(lua.LString(b))
	return 1
}

//SetDataBuf(buf types.IoBuffer)
func luaSetRequestDataBuf(L *lua.LState) int {
	data := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Request().SetDataBytes([]byte(data))
	return 0
}

// response function
//GetHeaders() types.HeaderMap
func luaGetResponseHeaders(L *lua.LState) int {
	headers := L.NewTable()
	ctx := GetContext(L)
	ctx.Response().GetHeaders().Range(func(key, value string) bool {
		headers.RawSetString(key, lua.LString(value))
		return true
	})
	L.Push(headers)
	return 1
}

func luaGetResponseHeader(L *lua.LState) int {
	key := L.ToString(-1)
	ctx := GetContext(L)
	if value, ok := ctx.Request().GetHeaders().Get(key); ok {
		L.Push(lua.LString(value))
	} else {
		L.Push(lua.LNil)
	}

	return 1
}

//DelHeader(key string)
func luaDelResponseHeaders(L *lua.LState) int {
	key := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Response().DelHeader(key)
	return 0
}

//SetHeader(key, value string)
func luaAddResponseHeaders(L *lua.LState) int {
	if L.GetTop() != 2 {
		// "expecting exactly 2 arguments");
		//return err
	}
	key := L.ToString(-2)
	value := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Response().SetHeader(key, value)
	return 0
}

//GetDataBuf() types.IoBuffer
func luaGetResponseDataBuf(L *lua.LState) int {
	ctx := GetContext(L)
	b := ctx.Response().GetDataBytes()
	L.Push(lua.LString(b))
	return 1
}

//SetDataBuf(buf types.IoBuffer)
func luaSetResponseDataBuf(L *lua.LState) int {
	data := L.ToString(-1)
	ctx := GetContext(L)
	ctx.Response().SetDataBytes([]byte(data))
	return 0
}

func SetContext(L *lua.LState, ctx types.Context) {
	ud := L.NewUserData()
	ud.Value = ctx

	L.Push(ud)
	L.SetGlobal(LuaContextKey, ud)
}

func GetContext(L *lua.LState) types.Context {
	v := L.GetGlobal(LuaContextKey)
	if ud, ok := v.(*lua.LUserData); ok {
		if ctx, ok := ud.Value.(types.Context); ok {
			return ctx
		}
	}
	return nil
}
