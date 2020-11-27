package rbac

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"strconv"
	"strings"

	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"mosn.io/api"
)

func getValueFromHeaderMapWithKeyList(keyList []string, headers api.HeaderMap) (string, bool) {
	for _, key := range keyList {
		value, ok := headers.Get(key)
		if ok && len(value) >= 0 {
			return value, ok
		}
	}
	return "", false
}

// parse rbac filter config to RBAC struct
func ParseRbacFilterConfig(jsonConf string) (filterConfig *rbactypes.RBACFilterConf, err error) {
	filterConfig = new(rbactypes.RBACFilterConf)

	// parse rules
	var un jsonpb.Unmarshaler
	un.AllowUnknownFields = true
	if err := un.Unmarshal(strings.NewReader(jsonConf), filterConfig); err != nil {
		return nil, err
	}

	if filterConfig.GetAppName() == "" {
		return nil, fmt.Errorf("[ParseRbacFilterConfig] parseing rbac filter configuration failed, err: missing app_name")
	}

	if filterConfig.GetVersion() == "" {
		return nil, fmt.Errorf("[ParseRbacFilterConfig] parseing rbac filter configuration failed, err: missing version field")
	}

	_, err = CheckVersionValid(filterConfig.GetAppName(), filterConfig.GetVersion())
	if err != nil {
		return nil, err
	}

	return filterConfig, nil
}

func CheckVersionValid(appname string, version string) (numericalVersion uint64, err error) {
	versionSeparator := "-"
	// version must in this format: ${appname}-${YYYYMMDD}-${numerical_version}
	items := strings.Split(version, versionSeparator)
	if len(items) != 3 {
		return 0, fmt.Errorf("[CheckVersionValid] rbac version must in format of ${appname}-${YYYYMMDD}-${numerical_version}, received: %s", version)
	}

	// appname check
	if items[0] != appname {
		return 0, fmt.Errorf("[CheckVersionValid] rbac version is not valid, expected appname: %s, received: %s", appname, items[0])
	}

	numericalVersion, err = strconv.ParseUint(items[1]+items[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[CheckVersionValid] rbac version must in format of ${appname}-${YYYYMMDD}-${numerical_version}, received: %s", version)
	}

	return numericalVersion, nil
}
