package zoneclient

import "gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

type ZoneClient struct {
	router Router
}

func New(router Router) *ZoneClient {
	return &ZoneClient{
		router: router,
	}
}

func (zc *ZoneClient) GetZoneNames() []string {
	return zc.router.GetZoneNames()
}

func (zc *ZoneClient) QueryZoneInfo(zoneName string) (model.ZoneInfo, error) {
	return zc.router.QueryZoneInfo(zoneName)
}

func (zc *ZoneClient) GetDefaultGzone() (model.ZoneInfo, error) {
	return zc.router.GetDefaultGzone(0)
}

func (zc *ZoneClient) GetDefaultCzone() (model.ZoneInfo, error) {
	return zc.router.GetDefaultCzone()
}

func (zc *ZoneClient) GetSrcZone(zoneName string) (model.ZoneInfo, error) {
	return zc.router.GetSrcZone(zoneName)
}

func (zc *ZoneClient) Route(rv model.RouterValue) (model.ZoneInfo, error) {
	return zc.router.Route(rv)
}

func (zc *ZoneClient) ElasticRoute(rv model.RouterValue) (model.ZoneInfo, error) {
	return zc.router.ElasticRoute(rv)
}

func (zc *ZoneClient) IsRpcElastic(interfaze, version, uniqueId, method string, isMark bool) (model.BizElastic, error) {
	return zc.router.IsRpcElastic(interfaze, version, uniqueId, method, isMark)
}

func (zc *ZoneClient) IsAppElastic(appName string, isMark bool) (bool, error) {
	return zc.router.IsAppElastic(appName, isMark)
}

func (zc *ZoneClient) IsGzRpc(service string) (bool, error) {
	return zc.router.IsGzRpc(service)
}

func (zc *ZoneClient) QueryAllZoneInfo() ([]model.ZoneInfo, error) {
	return zc.router.QueryAllZoneInfo()
}

func (zc *ZoneClient) GetCurZone() (model.ZoneInfo, error) {
	return zc.router.GetCurZone()
}

func (zc *ZoneClient) GetDisasterStatus(uid string, isElastic bool) (model.DisasterStatusEnum, error) {
	return zc.router.GetDisasterStatus(uid, isElastic)
}

func (zc *ZoneClient) GetUidZone(uid string) (model.UidZone, error) {
	return zc.router.GetUidZone(uid)
}

func (zc *ZoneClient) GetProdUidZone(uid string) (model.UidZone, error) {
	return zc.router.GetProdUidZone(uid)
}

func (zc *ZoneClient) IsElastic(uid string) (bool, error) {
	return zc.router.IsElastic(uid)
}

func (zc *ZoneClient) GetElasticStatus(isMark bool) (model.ElasticStatusEnum, error) {
	return zc.router.GetElasticStatus(isMark)
}

func (zc *ZoneClient) QueryOriginZone(zoneName, flowType string) (model.ZoneInfo, error) {
	return zc.router.QueryOriginZone(zoneName, flowType)
}

func (zc *ZoneClient) QueryDefaultCzone(flowType string) (model.ZoneInfo, error) {
	return zc.router.QueryDefaultCzone(flowType)
}

func (zc *ZoneClient) IsUidElastic(uid string, flowType string) (bool, error) {
	return zc.router.IsUidElastic(uid, flowType)
}

func (zc *ZoneClient) GetElasticLabel(eid, interfaze, version, uniqueId, method string, isMark bool) (string, error) {
	return zc.router.GetElasticLabel(eid, interfaze, version, uniqueId, method, isMark)
}
func (zc *ZoneClient) GetCloselyRzone(filterElasticZone bool, flowType string) (model.ZoneInfo, error) {
	return zc.router.GetCloselyRzone(filterElasticZone, flowType)
}
