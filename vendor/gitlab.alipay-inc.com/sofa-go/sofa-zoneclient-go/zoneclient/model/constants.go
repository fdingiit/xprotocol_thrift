package model

import (
	"os"
	"os/user"
	"path"
	"runtime"
)

var User_Define_Zone_Mng_Host = os.Getenv("ZONE_MNG")

const (
	HttpServerAddress  = "http://zonemng-pool"
	HttpPort           = 80
	HttpZoneInfoUrl    = "/rest/queryConfig.json?key=zoneInfo&zone="
	HttpElasticRuleUrl = "/rest/elasticRule?version="
)

var (
	MosnConfPath        = "/home/admin/mosn/conf/"
	ZoneInfoFilePath    = FetchZoneClientHome() + "routeRule.json"
	ElasticRuleFilePath = FetchZoneClientHome() + "elasticRule.json"
)

const (
	ROUTE_OUT_GRAY = -1
	ROUTE_IN_GRAY  = 1
	ROUTE_DEFAULT  = 0
)

type DisasterStatusEnum string

// disaster status
const (
	NORMAL DisasterStatusEnum = "NORMAL"
	REMOTE DisasterStatusEnum = "REMOTE"
	LOCAL  DisasterStatusEnum = "LOCAL"
)

// disaster type
type DrCommonCode string

const (
	Normal DrCommonCode = "normal"
	Ldr    DrCommonCode = "ldr"   // local disaster rule
	Rdr    DrCommonCode = "rdr"   // remote disaster rule
	Group  DrCommonCode = "group" // dr rule level
	Zone   DrCommonCode = "zone"  // dr rule level
)

const CONVERTUID_STRING = "ABCDEFGHIJ"

type ElasticSubRuleStatus string

// elastic sub rule status
const (
	VALID   ElasticSubRuleStatus = "2"
	PRESS   ElasticSubRuleStatus = "1"
	INVALID ElasticSubRuleStatus = "0"
)

type ElasticStatusEnum string

// global elastic status
const (
	ES_NORMAL  ElasticStatusEnum = "NORMAL"
	ES_ELASTIC ElasticStatusEnum = "ELASTIC"
	ES_BACK    ElasticStatusEnum = "BACK" // deprecated
)

type ZoneStatusEnum string

// zone running status
const (
	ZS_BUILDING ZoneStatusEnum = "BUILDING"
	ZS_RUNNING  ZoneStatusEnum = "RUNNING"
)

type FlowTypeEnum string

// flow rule type, mark flow is 'PRESS', un-mark flow is 'ONLINE'
const (
	FT_ONLINE FlowTypeEnum = "ONLINE"
	FT_PRESS  FlowTypeEnum = "PRESS"
)

type UidTypeEnum string

// uid elastic type
const (
	ELASTIC_UID    UidTypeEnum = "ELASTIC"    // elastic uid
	NO_ELASTIC_UID UidTypeEnum = "NO_ELASTIC" // un-elastic uid
)

const (
	DEFAULT_EID    = "-1"
	NO_ELASTIC_EID = "-2"
)

// drm resource info
const (
	DRM_DOMAIN               string = "routeClient"
	ZONE_INFO                string = "zoneInfo"
	ELASTIC_RULE             string = "elasticRuleVersion"
	ZONE_COLOR               string = "zoneColor"
	WHITE_LIST_RPC_LOAD_TEST string = "whiteListRPCLoadTest"
)

// elastic label
const (
	ELASTIC_DEFAULT   = "D"
	ELASTIC_UNDEFAULT = "T"
	NO_ELASTIC        = "F"
	NO_CONSISTENT     = "N"
)

func FetchZoneClientHome() string {
	if u, err := user.Current(); err == nil {
		if runtime.GOOS == "darwin" {
			MosnConfPath = path.Join(u.HomeDir, "conf")
		} else if runtime.GOOS == "windows" {
			MosnConfPath = path.Join(u.HomeDir, "conf")
		}
	}

	return MosnConfPath
}
