package zoneclient

import "gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

type Router interface {
	GetZoneNames() []string
	QueryZoneInfo(zoneName string) (model.ZoneInfo, error)
	GetDefaultGzone(inGray int32) (model.ZoneInfo, error)
	GetDefaultCzone() (model.ZoneInfo, error)
	GetSrcZone(zoneName string) (model.ZoneInfo, error)
	Route(rv model.RouterValue) (model.ZoneInfo, error)
	ElasticRoute(rv model.RouterValue) (model.ZoneInfo, error)
	IsRpcElastic(interfaze, version, uniqueId, method string, isMark bool) (model.BizElastic, error)
	IsAppElastic(appName string, isMark bool) (bool, error)
	IsGzRpc(service string) (bool, error)
	QueryAllZoneInfo() ([]model.ZoneInfo, error)
	GetCurZone() (model.ZoneInfo, error)
	GetDisasterStatus(uid string, isElastic bool) (model.DisasterStatusEnum, error)
	GetUidZone(uid string) (model.UidZone, error)
	GetProdUidZone(uid string) (model.UidZone, error)
	IsElastic(uid string) (bool, error)
	GetElasticStatus(isMark bool) (model.ElasticStatusEnum, error)
	QueryOriginZone(zoneName, flowType string) (model.ZoneInfo, error)
	QueryDefaultCzone(flowType string) (model.ZoneInfo, error)
	IsUidElastic(uid string, flowType string) (bool, error)
	GetElasticLabel(eid, interfaze, version, uniqueId, method string, isMark bool) (string, error)
	GetCloselyRzone(filterElasticZone bool, flowType string) (model.ZoneInfo, error)
}
