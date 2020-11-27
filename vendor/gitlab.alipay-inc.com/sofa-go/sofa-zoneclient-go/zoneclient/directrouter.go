package zoneclient

import "gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

type DirectRouter struct {
}

func NewDirectRouter() *DirectRouter {
	return &DirectRouter{}
}

func (dr *DirectRouter) GetZoneNames() []string {
	return nil
}

func (dr *DirectRouter) QueryZoneInfo(zoneName string) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) GetDefaultGzone(inGray int32) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) GetDefaultCzone() (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) GetSrcZone(zoneName string) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) Route(rv model.RouterValue) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) ElasticRoute(rv model.RouterValue) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) IsRpcElastic(interfaze, version, uniqueId, method string,
	isMark bool) (model.BizElastic, error) {
	return model.BizElastic{}, nil
}

func (dr *DirectRouter) IsAppElastic(appName string, isMark bool) (bool, error) {
	return false, nil
}

func (dr *DirectRouter) IsGzRpc(service string) (bool, error) {
	return false, nil
}

func (dr *DirectRouter) QueryAllZoneInfo() ([]model.ZoneInfo, error) {
	return nil, nil
}

func (dr *DirectRouter) GetCurZone() (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (dr *DirectRouter) GetDisasterStatus(uid string, isElastic bool) (model.DisasterStatusEnum, error) {
	return model.NORMAL, nil
}

func (dr *DirectRouter) GetUidZone(uid string) (model.UidZone, error) {
	return model.UidZone{}, nil
}

func (dr *DirectRouter) GetProdUidZone(uid string) (model.UidZone, error) {
	return model.UidZone{}, nil
}

func (dr *DirectRouter) IsElastic(uid string) (bool, error) {
	return false, nil
}

func (dr *DirectRouter) GetElasticStatus(isMark bool) (model.ElasticStatusEnum, error) {
	return model.ElasticStatusEnum(""), nil
}

func (al *DirectRouter) QueryOriginZone(zoneName, flowType string) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (al *DirectRouter) QueryDefaultCzone(flowType string) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}

func (al *DirectRouter) IsUidElastic(uid string, flowType string) (bool, error) {
	return false, nil
}

func (al *DirectRouter) GetElasticLabel(eid, interfaze, version, uniqueId, method string, isMark bool) (string, error) {
	return "", nil
}
func (al *DirectRouter) GetCloselyRzone(filterElasticZone bool, flowType string) (model.ZoneInfo, error) {
	return model.ZoneInfo{}, nil
}
