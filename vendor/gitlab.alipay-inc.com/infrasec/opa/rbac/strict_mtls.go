package rbac

import (
	"fmt"
	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"gitlab.alipay-inc.com/infrasec/opa/constants"
	"mosn.io/api"
	"mosn.io/mosn/pkg/mtls"
)

type InheritStrictMTLS struct {
	Enable             bool
	ServiceNameKeyList []string
	CheckServiceList   []StringMatcher
}

// Receive the rbactypes.StrictMTLSConf input and convert it to InheritStrictMTLS
func NewInheritStrictMTLS(strictMTLS *rbactypes.StrictMTLSConf) (*InheritStrictMTLS, error) {
	inheritStrictMTLS := new(InheritStrictMTLS)

	// fill EnableStrictMTLS, false by default
	inheritStrictMTLS.Enable = strictMTLS.GetEnable()

	// default match ["service"]
	inheritStrictMTLS.ServiceNameKeyList = []string{constants.SOFARPC_ROUTER_SERVICE_MATCH_KEY}
	if strictMTLS.GetServiceNameKeyList() != nil && len(strictMTLS.GetServiceNameKeyList()) > 0 {
		inheritStrictMTLS.ServiceNameKeyList = strictMTLS.GetServiceNameKeyList()
	}

	// fill CheckServiceList
	inheritStrictMTLS.CheckServiceList = make([]StringMatcher, len(strictMTLS.GetCheckServiceList()))
	for idx, matcherConf := range strictMTLS.GetCheckServiceList() {
		matcher, err := NewStringMatcher(matcherConf)
		if err != nil {
			return nil, fmt.Errorf("[NewInheritStrictMTLS] failed to parse strict mtls policy, err: %v", err)
		}
		inheritStrictMTLS.CheckServiceList[idx] = matcher
	}
	return inheritStrictMTLS, nil
}

func (strictMTLS *InheritStrictMTLS) Allowed(cb api.StreamReceiverFilterHandler, headers api.HeaderMap) (allowed bool, serviceName string) {
	if strictMTLS.Enable == false {
		return true, ""
	}

	// Check conn is mTLS or not
	if cb == nil || cb.Connection() == nil || cb.Connection().RawConn() == nil {
		return true, ""
	}
	conn := cb.Connection().RawConn()
	if _, ok := conn.(*mtls.TLSConn); ok {
		return true, ""
	}
	serviceName, ok := getValueFromHeaderMapWithKeyList(strictMTLS.ServiceNameKeyList, headers)
	// if service name is not found in the header, return true
	if !ok {
		//fmt.Errorf("[InheritStrictMTLS.Allowed] failed to parse service name in rbac filter")
		return true, ""
	}
	for _, matcher := range strictMTLS.CheckServiceList {
		if matcher.Equal(serviceName) {
			// Service need in mTLS and conn is not mTLS
			return false, serviceName
		}
	}

	return true, ""
}
