package rbac

import (
	"context"
	"fmt"
	"gitlab.alipay-inc.com/infrasec/opa/constants"
	"mosn.io/mosn/pkg/types"
	"net"
	"reflect"
	"strconv"

	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"mosn.io/api"
)

type InheritPermission interface {
	isInheritPermission()
	// A policy matches if and only if at least one of InheritPermission.Match return true
	// AND at least one of InheritPrincipal.Match return true
	Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool
}

func (*PermissionAny) isInheritPermission()             {}
func (*PermissionDestinationIp) isInheritPermission()   {}
func (*PermissionDestinationPort) isInheritPermission() {}
func (*PermissionHeader) isInheritPermission()          {}
func (*PermissionProviderAppname) isInheritPermission() {}
func (*PermissionAndRules) isInheritPermission()        {}
func (*PermissionOrRules) isInheritPermission()         {}
func (*PermissionRpcService) isInheritPermission()      {}

// PermissionConf_Any
type PermissionAny struct {
	Any bool
}

func NewPermissionAny(permission *rbactypes.PermissionConf_Any) (*PermissionAny, error) {
	return &PermissionAny{
		Any: permission.Any,
	}, nil
}

func (permission *PermissionAny) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	return permission.Any
}

// PermissionConf_DestinationIp
type PermissionDestinationIp struct {
	CidrRange   *net.IPNet
	InvertMatch bool
}

func NewPermissionDestinationIp(permission *rbactypes.PermissionConf_DestinationIp) (*PermissionDestinationIp, error) {
	addressPrefix := permission.DestinationIp.GetAddressPrefix()
	prefixLen := permission.DestinationIp.GetPrefixLen()
	_, ipNet, err := net.ParseCIDR(addressPrefix + "/" + strconv.Itoa(int(prefixLen)))
	if err != nil {
		return nil, err
	}
	inheritPermission := &PermissionDestinationIp{
		CidrRange:   ipNet,
		InvertMatch: permission.DestinationIp.GetInvertMatch(),
	}
	return inheritPermission, nil
}

func (permission *PermissionDestinationIp) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	if cb == nil || cb.Connection() == nil {
		return false
	}
	localAddr := cb.Connection().LocalAddr()
	addr, err := net.ResolveTCPAddr(localAddr.Network(), localAddr.String())
	if err != nil {
		//fmt.Errorf("[PermissionDestinationIp.Match] failed to parse local address in rbac filter, err: %v", err)
		return false
	}
	isMatch := permission.CidrRange.Contains(addr.IP)
	// InvertMatch xor isMatch
	return isMatch != permission.InvertMatch
}

// PermissionConf_DestinationPort
type PermissionDestinationPort struct {
	DestinationPort uint32
}

func NewPermissionDestinationPort(permission *rbactypes.PermissionConf_DestinationPort) (*PermissionDestinationPort, error) {
	return &PermissionDestinationPort{
		DestinationPort: permission.DestinationPort,
	}, nil
}

func (permission *PermissionDestinationPort) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	if cb == nil || cb.Connection() == nil {
		return false
	}
	localAddr := cb.Connection().LocalAddr()
	addr, err := net.ResolveTCPAddr(localAddr.Network(), localAddr.String())
	if err != nil {
		return false
	}
	return addr.Port == int(permission.DestinationPort)
}

// PermissionConf_Header
type PermissionHeader struct {
	Target      string
	Matcher     HeaderMatcher
	InvertMatch bool
}

func NewPermissionHeader(permission *rbactypes.PermissionConf_Header) (*PermissionHeader, error) {
	headerMatcher, err := NewHeaderMatcher(permission.Header)
	if err != nil {
		return nil, err
	}

	inheritPermission := &PermissionHeader{}
	inheritPermission.Target = permission.Header.GetName()
	inheritPermission.InvertMatch = permission.Header.GetInvertMatch()
	inheritPermission.Matcher = headerMatcher
	return inheritPermission, nil
}

func (permission *PermissionHeader) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	targetValue, found := headers.Get(permission.Target)

	// HeaderMatcherPresentMatch is a little special
	if matcher, ok := permission.Matcher.(*HeaderMatcherPresentMatch); ok {
		// HeaderMatcherPresentMatch matches if and only if header found and PresentMatch is true
		isMatch := found && matcher.PresentMatch
		return permission.InvertMatch != isMatch
	}

	// return false when targetValue is not found, except matcher is `HeaderMatcherPresentMatch`
	if !found {
		return false
	}

	isMatch := permission.Matcher.Equal(targetValue)
	// permission.InvertMatch xor isMatch
	return permission.InvertMatch != isMatch
}

// PermissionConf_ProviderAppname
type PermissionProviderAppname struct {
	ProviderAppname        StringMatcher
	ProviderAppnameKeyList []string
}

func NewPermissionProviderAppname(permission *rbactypes.PermissionConf_ProviderAppname) (*PermissionProviderAppname, error) {
	providerAppnameMatcher, err := NewStringMatcher(permission.ProviderAppname)
	if err != nil {
		return nil, fmt.Errorf("[NewPermissionProviderAppname] failed to provider_appname matcher, err: %v", err)
	}

	// default match ["provider_app_name_local"]
	return &PermissionProviderAppname{
		ProviderAppname:        providerAppnameMatcher,
		ProviderAppnameKeyList: []string{constants.RBAC_PROVIDER_APP_NAME_KEY},
	}, nil
}

func (permission *PermissionProviderAppname) Match(context context.Context, cb api.StreamReceiverFilterHandler, headers types.HeaderMap) bool {
	currentAppName, _ := getValueFromHeaderMapWithKeyList(permission.ProviderAppnameKeyList, headers)
	return permission.ProviderAppname.Equal(currentAppName)
}

// PermissionConf_AndRules
type PermissionAndRules struct {
	AndRules []InheritPermission
}

func NewPermissionAndRules(permission *rbactypes.PermissionConf_AndRules) (*PermissionAndRules, error) {
	inheritPermission := &PermissionAndRules{}
	inheritPermission.AndRules = make([]InheritPermission, len(permission.AndRules.GetRules()))
	for idx, subPermission := range permission.AndRules.GetRules() {
		if subInheritPermission, err := NewInheritPermission(subPermission); err != nil {
			return nil, err
		} else {
			inheritPermission.AndRules[idx] = subInheritPermission
		}
	}
	return inheritPermission, nil
}

func (permission *PermissionAndRules) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	for _, rule := range permission.AndRules {
		if isMatch := rule.Match(ctx, cb, headers); isMatch {
			continue
		}
		return false
	}
	return true
}

// PermissionConf_OrRules
type PermissionOrRules struct {
	OrRules []InheritPermission
}

func NewPermissionOrRules(permission *rbactypes.PermissionConf_OrRules) (*PermissionOrRules, error) {
	inheritPermission := &PermissionOrRules{}
	inheritPermission.OrRules = make([]InheritPermission, len(permission.OrRules.GetRules()))
	for idx, subPermission := range permission.OrRules.GetRules() {
		if subInheritPermission, err := NewInheritPermission(subPermission); err != nil {
			return nil, err
		} else {
			inheritPermission.OrRules[idx] = subInheritPermission
		}
	}
	return inheritPermission, nil
}

func (permission *PermissionOrRules) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	for _, rule := range permission.OrRules {
		if isMatch := rule.Match(ctx, cb, headers); isMatch {
			return true
		}
	}
	return false
}

// PermissionConf_RpcService
type PermissionRpcService struct {
	ServiceNameKeyList []string
	ServiceName        StringMatcher
	MethodNameKeyList  []string
	MethodName         StringMatcher
}

func NewPermissionRpcService(permission *rbactypes.PermissionConf_RpcService) (permissionRpcService *PermissionRpcService, err error) {
	inheritPermission := &PermissionRpcService{}

	if permission.RpcService.GetServiceName() != nil {
		serviceNameMatcher, err := NewStringMatcher(permission.RpcService.GetServiceName())
		if err != nil {
			return nil, fmt.Errorf("[NewPermissionRpcService] failed to parse server name matcher, err: %v", err)
		}
		inheritPermission.ServiceName = serviceNameMatcher
	}

	// methodNameMatcher is allowed to be nil
	if permission.RpcService.GetMethodName() != nil {
		methodNameMatcher, err := NewStringMatcher(permission.RpcService.GetMethodName())
		if err != nil {
			return nil, fmt.Errorf("[NewPermissionRpcService] failed to parse method name matcher, err: %v", err)
		}
		inheritPermission.MethodName = methodNameMatcher
	}

	// default match ["service"]
	inheritPermission.ServiceNameKeyList = []string{constants.SOFARPC_ROUTER_SERVICE_MATCH_KEY}
	if permission.RpcService.GetServiceNameKeyList() != nil && len(permission.RpcService.GetServiceNameKeyList()) > 0 {
		inheritPermission.ServiceNameKeyList = permission.RpcService.GetServiceNameKeyList()
	}

	// default match ["sofa_head_method_name"]
	inheritPermission.MethodNameKeyList = []string{constants.SOFARPC_ROUTER_HEADER_METHOD_KEY}
	if permission.RpcService.GetMethodNameKeyList() != nil && len(permission.RpcService.GetMethodNameKeyList()) > 0 {
		inheritPermission.MethodNameKeyList = permission.RpcService.GetMethodNameKeyList()
	}
	return inheritPermission, nil
}

func (permission *PermissionRpcService) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	serviceName, ok := getValueFromHeaderMapWithKeyList(permission.ServiceNameKeyList, headers)

	// if service name is not found in the header, return false
	if !ok {
		//fmt.Errorf("[PermissionRpcService.Match] failed to parse service name in rbac filter")
		return false
	}

	if permission.ServiceName == nil && permission.MethodName == nil {
		return false
	}

	var serviceNameHit, methodNameHit bool
	// check service name
	if permission.ServiceName != nil {
		serviceNameHit = permission.ServiceName.Equal(serviceName)
	} else {
		serviceNameHit = true
	}
	// if permission.MethodName is empty, skip the method name check
	if permission.MethodName == nil {
		return serviceNameHit
	}

	// check the method name
	methodName, ok := getValueFromHeaderMapWithKeyList(permission.MethodNameKeyList, headers)

	// if method name is not found in the header, return false
	if !ok {
		//fmt.Errorf("[PermissionRpcService.Match] failed to parse method name in rbac filter")
		return false
	}

	if permission.MethodName != nil {
		methodNameHit = permission.MethodName.Equal(methodName)
	} else {
		methodNameHit = true
	}
	return serviceNameHit == true && methodNameHit == true
}

// Receive the rbactypes.PermissionConf input and convert it to mosn rbac permission
func NewInheritPermission(permission *rbactypes.PermissionConf) (InheritPermission, error) {
	switch permission.GetRule().(type) {
	case *rbactypes.PermissionConf_Any:
		return NewPermissionAny(permission.GetRule().(*rbactypes.PermissionConf_Any))
	case *rbactypes.PermissionConf_DestinationIp:
		return NewPermissionDestinationIp(permission.GetRule().(*rbactypes.PermissionConf_DestinationIp))
	case *rbactypes.PermissionConf_DestinationPort:
		return NewPermissionDestinationPort(permission.GetRule().(*rbactypes.PermissionConf_DestinationPort))
	case *rbactypes.PermissionConf_Header:
		return NewPermissionHeader(permission.GetRule().(*rbactypes.PermissionConf_Header))
	case *rbactypes.PermissionConf_ProviderAppname:
		return NewPermissionProviderAppname(permission.GetRule().(*rbactypes.PermissionConf_ProviderAppname))
	case *rbactypes.PermissionConf_AndRules:
		return NewPermissionAndRules(permission.GetRule().(*rbactypes.PermissionConf_AndRules))
	case *rbactypes.PermissionConf_OrRules:
		return NewPermissionOrRules(permission.GetRule().(*rbactypes.PermissionConf_OrRules))
	case *rbactypes.PermissionConf_RpcService:
		return NewPermissionRpcService(permission.GetRule().(*rbactypes.PermissionConf_RpcService))
	default:
		return nil, fmt.Errorf("[NewInheritPermission] not supported Permission.Rule type found, detail: %v",
			reflect.TypeOf(permission.GetRule()))
	}
}
