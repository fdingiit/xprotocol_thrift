package mist

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"gitlab.alipay-inc.com/infrasec/api/mist/types"
	"strconv"
	"strings"
)

func ParseMistSdkConfig(jsonConf string) (config *types.MISTSdkConfig, err error) {
	config = new(types.MISTSdkConfig)

	// parse rules
	var un jsonpb.Unmarshaler
	un.AllowUnknownFields = true
	if err := un.Unmarshal(strings.NewReader(jsonConf), config); err != nil {
		return nil, err
	}

	if config.GetAppName() == "" {
		return nil, fmt.Errorf("[ParseMistSdkConfig] parseing mist sdk configuration failed, err: missing app_name")
	}

	if config.GetVersion() == "" {
		return nil, fmt.Errorf("[ParseMistSdkConfig] parseing mist sdk configuration failed, err: missing version field")
	}

	_, err = CheckVersionValid(config.GetAppName(), config.GetVersion())
	if err != nil {
		return nil, err
	}

	return config, nil
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
