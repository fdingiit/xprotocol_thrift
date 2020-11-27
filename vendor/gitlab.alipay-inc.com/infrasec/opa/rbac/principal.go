package rbac

import (
	"fmt"
	"gitlab.alipay-inc.com/infrasec/opa/constants"
	"net"
	"reflect"
	"strconv"

	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"mosn.io/api"
	"mosn.io/mosn/pkg/mtls"
)

type InheritPrincipal interface {
	isInheritPrincipal()
	// A policy matches if and only if at least one of InheritPermission.Match return true
	// AND at least one of InheritPrincipal.Match return true
	Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool
}

func (*PrincipalAny) isInheritPrincipal()         {}
func (*PrincipalSourceIp) isInheritPrincipal()    {}
func (*PrincipalHeader) isInheritPrincipal()      {}
func (*PrincipalAndIds) isInheritPrincipal()      {}
func (*PrincipalOrIds) isInheritPrincipal()       {}
func (*PrincipalAppIdentity) isInheritPrincipal() {}

// PrincipalConf_Any
type PrincipalAny struct {
	Any bool
}

func NewPrincipalAny(principal *rbactypes.PrincipalConf_Any) (*PrincipalAny, error) {
	return &PrincipalAny{
		Any: principal.Any,
	}, nil
}

func (principal *PrincipalAny) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	return principal.Any
}

// PrincipalConf_SourceIp
type PrincipalSourceIp struct {
	CidrRange   *net.IPNet
	InvertMatch bool
}

func NewPrincipalSourceIp(principal *rbactypes.PrincipalConf_SourceIp) (*PrincipalSourceIp, error) {
	addressPrefix := principal.SourceIp.GetAddressPrefix()
	prefixLen := principal.SourceIp.GetPrefixLen()
	_, ipNet, err := net.ParseCIDR(addressPrefix + "/" + strconv.Itoa(int(prefixLen)))
	if err != nil {
		return nil, err
	}

	inheritPrincipal := &PrincipalSourceIp{
		CidrRange:   ipNet,
		InvertMatch: principal.SourceIp.GetInvertMatch(),
	}
	return inheritPrincipal, nil
}

func (principal *PrincipalSourceIp) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	if cb == nil || cb.Connection() == nil {
		return false
	}
	remoteAddr := cb.Connection().RemoteAddr()
	addr, err := net.ResolveTCPAddr(remoteAddr.Network(), remoteAddr.String())
	if err != nil {
		//fmt.Errorf("[PrincipalSourceIp.Match] failed to parse remote address in rbac filter, err: %v", err)
		return false
	}
	isMatch := principal.CidrRange.Contains(addr.IP)
	// InvertMatch xor isMatch
	return isMatch != principal.InvertMatch
}

// PrincipalConf_Header
type PrincipalHeader struct {
	Target      string
	Matcher     HeaderMatcher
	InvertMatch bool
}

func NewPrincipalHeader(principal *rbactypes.PrincipalConf_Header) (*PrincipalHeader, error) {
	headerMatcher, err := NewHeaderMatcher(principal.Header)
	if err != nil {
		return nil, err
	}

	inheritPrincipal := &PrincipalHeader{}
	inheritPrincipal.Target = principal.Header.GetName()
	inheritPrincipal.InvertMatch = principal.Header.GetInvertMatch()
	inheritPrincipal.Matcher = headerMatcher
	return inheritPrincipal, nil
}

func (principal *PrincipalHeader) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	targetValue, found := headers.Get(principal.Target)

	// HeaderMatcherPresentMatch is a little special
	if matcher, ok := principal.Matcher.(*HeaderMatcherPresentMatch); ok {
		// HeaderMatcherPresentMatch matches if and only if header found and PresentMatch is true
		isMatch := found && matcher.PresentMatch
		return principal.InvertMatch != isMatch
	}

	// return false when targetValue is not found, except matcher is `HeaderMatcherPresentMatch`
	if !found {
		return false
	}

	isMatch := principal.Matcher.Equal(targetValue)
	// principal.InvertMatch xor isMatch
	return principal.InvertMatch != isMatch
}

// PrincipalConf_AndIds
type PrincipalAndIds struct {
	AndIds []InheritPrincipal
}

func NewPrincipalAndIds(principal *rbactypes.PrincipalConf_AndIds) (*PrincipalAndIds, error) {
	inheritPrincipal := &PrincipalAndIds{}
	inheritPrincipal.AndIds = make([]InheritPrincipal, len(principal.AndIds.GetIds()))
	for idx, subPrincipal := range principal.AndIds.GetIds() {
		if subInheritPrincipal, err := NewInheritPrincipal(subPrincipal); err != nil {
			return nil, err
		} else {
			inheritPrincipal.AndIds[idx] = subInheritPrincipal
		}
	}
	return inheritPrincipal, nil
}

func (principal *PrincipalAndIds) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	for _, ids := range principal.AndIds {
		if isMatch := ids.Match(cb, headers); isMatch {
			continue
		}
		return false
	}
	return true
}

// PrincipalConf_OrIds
type PrincipalOrIds struct {
	OrIds []InheritPrincipal
}

func NewPrincipalOrIds(principal *rbactypes.PrincipalConf_OrIds) (*PrincipalOrIds, error) {
	inheritPrincipal := &PrincipalOrIds{}
	inheritPrincipal.OrIds = make([]InheritPrincipal, len(principal.OrIds.GetIds()))
	for idx, subPrincipal := range principal.OrIds.GetIds() {
		if subInheritPrincipal, err := NewInheritPrincipal(subPrincipal); err != nil {
			return nil, err
		} else {
			inheritPrincipal.OrIds[idx] = subInheritPrincipal
		}
	}
	return inheritPrincipal, nil
}

func (principal *PrincipalOrIds) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	for _, ids := range principal.OrIds {
		if isMatch := ids.Match(cb, headers); isMatch {
			return true
		}
	}
	return false
}

// PrincipalConf_AppIdentity
type PrincipalAppIdentity struct {
	AppName        StringMatcher
	AppNameKeyList []string
	InvertMatch    bool
	MtlsMatch      bool
	StrictMatch    bool
}

func NewPrincipalAppIdentity(principal *rbactypes.PrincipalConf_AppIdentity) (*PrincipalAppIdentity, error) {
	appNameMatcher, err := NewStringMatcher(principal.AppIdentity.GetAppName())
	if err != nil {
		return nil, fmt.Errorf("[NewPrincipalAppIdentity] failed to parse app name matcher, err: %v", err)
	}

	inheritPrincipal := &PrincipalAppIdentity{}
	inheritPrincipal.AppName = appNameMatcher
	inheritPrincipal.InvertMatch = principal.AppIdentity.GetInvertMatch()
	inheritPrincipal.MtlsMatch = principal.AppIdentity.GetMtlsMatch()
	inheritPrincipal.StrictMatch = principal.AppIdentity.GetStrictMatch()

	// default match ["mist_trust_identity", "rpc_trace_context.sofaCallerApp"]
	inheritPrincipal.AppNameKeyList = []string{constants.SOFARPC_ROUTER_HEADER_TRUST_IDENTITY, constants.SOFARPC_ROUTER_HEADER_CALLER_APP_KEY}
	if principal.AppIdentity.GetAppNameKeyList() != nil && len(principal.AppIdentity.GetAppNameKeyList()) > 0 {
		inheritPrincipal.AppNameKeyList = principal.AppIdentity.GetAppNameKeyList()
	}

	return inheritPrincipal, nil
}

func (principal *PrincipalAppIdentity) Match(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	if principal.MtlsMatch {
		if cb != nil && cb.Connection() != nil && cb.Connection().RawConn() != nil {
			conn := cb.Connection().RawConn()
			if tlsConn, ok := conn.(*mtls.TLSConn); ok {
				// mTLS handle
				cert := tlsConn.ConnectionState().PeerCertificates[0]
				// check app name
				hit := principal.AppName.Equal(cert.Subject.CommonName)
				return hit != principal.InvertMatch
			}
		}
		if principal.StrictMatch {
			return principal.StrictMatch == principal.InvertMatch
		}
	}
	// fetch appname from header
	appname, ok := getValueFromHeaderMapWithKeyList(principal.AppNameKeyList, headers)

	// if appname is not found in the header, return false
	if !ok {
		//fmt.Errorf("[PrincipalAppIdentity.Match] failed to parse app name in rbac filter")
		return principal.StrictMatch == principal.InvertMatch
	}

	// check app name
	hit := principal.AppName.Equal(appname)
	return hit != principal.InvertMatch
}

// Receive the rbactypes.PrincipalConf input and convert it to mosn rbac principal
func NewInheritPrincipal(principal *rbactypes.PrincipalConf) (InheritPrincipal, error) {
	switch principal.GetIdentifier().(type) {
	case *rbactypes.PrincipalConf_Any:
		return NewPrincipalAny(principal.GetIdentifier().(*rbactypes.PrincipalConf_Any))
	case *rbactypes.PrincipalConf_SourceIp:
		return NewPrincipalSourceIp(principal.GetIdentifier().(*rbactypes.PrincipalConf_SourceIp))
	case *rbactypes.PrincipalConf_Header:
		return NewPrincipalHeader(principal.GetIdentifier().(*rbactypes.PrincipalConf_Header))
	case *rbactypes.PrincipalConf_AndIds:
		return NewPrincipalAndIds(principal.GetIdentifier().(*rbactypes.PrincipalConf_AndIds))
	case *rbactypes.PrincipalConf_OrIds:
		return NewPrincipalOrIds(principal.GetIdentifier().(*rbactypes.PrincipalConf_OrIds))
	case *rbactypes.PrincipalConf_AppIdentity:
		return NewPrincipalAppIdentity(principal.GetIdentifier().(*rbactypes.PrincipalConf_AppIdentity))
	default:
		return nil, fmt.Errorf("[NewInheritPrincipal] not supported Principal.Identifier type found, detail: %v",
			reflect.TypeOf(principal.GetIdentifier()))
	}
}
