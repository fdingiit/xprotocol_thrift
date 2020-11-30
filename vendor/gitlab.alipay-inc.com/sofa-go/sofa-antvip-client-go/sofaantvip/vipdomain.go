package sofaantvip

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip/protobuf"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

type RealNode struct {
	Ip                  string            `json:"ip"`
	Weight              int32             `json:"weight"`
	HealthCheckPort     int32             `json:"healthCheckPort"`
	Zone                string            `json:"zone"`
	Available           bool              `json:"available"`
	RoundTripTime       int64             `json:"roundTripTime"`
	Reason              string            `json:"reason"`
	LastHealthCheckTime int64             `json:"lastHealthCheckTime"`
	DataCenter          string            `json:"dataCenter"`
	Falling             bool              `json:"falling"`
	Labels              map[string]string `json:"labels"`
}

func (rn *RealNode) String() string {
	key := "ip=%s weight=%d healthcheck_port=%d zone=%s available=%T " +
		"datacenter=%s falling=%T"
	return fmt.Sprintf(key, rn.Ip, rn.Weight, rn.HealthCheckPort,
		rn.Zone, rn.Available, rn.DataCenter, rn.Falling)
}

func (rn *RealNode) FromProtobuf(msg *protobuf.VipDomainMsg_RealNodeMsg) {
	rn.Ip = msg.GetIp()
	rn.Weight = msg.GetWeight()
	rn.Available = msg.GetAvailable()
	rn.Zone = msg.GetZone()
	rn.HealthCheckPort = msg.GetHealthCheckPort()
}

type VipDomain struct {
	Name                   string            `json:"name"`
	ProtectThreshold       int32             `json:"protectThreshold"`
	HealthCheckType        string            `json:"healthCheckType"`
	HealthCheckDefaultPort int32             `json:"healthCheckDefaultPort"`
	HealthCheckTimeout     int               `json:"healthCheckTimeout"`
	HealthCheckInterval    int               `json:"healthCheckInterval"`
	HealthCheckRaise       int               `json:"healthCheckRaise"`
	HealthCheckFall        int               `json:"healthCheckFall"`
	HealthCheckEnable      bool              `json:"healthCheckEnable"`
	HealthCheckPayload     map[string]string `json:"healthCheckPayload"`
	Version                int64             `json:"version"`
	RealNodes              []RealNode        `json:"realNodes"`
	Labels                 map[string]string `json:"labels"`
	IsDeleted              bool              `json:"-"`
	ZoneInfoList           *ZoneInfoList     `json:"-"`

	weightedAvailableRealServers []RealServer
	availableRealServers         []RealServer
	allRealServer                []RealServer
	weightedAllRealServers       []RealServer
	idcWeightedAvailableRsMap    map[string][]RealServer
	cityWeightedAvailableRsMap   map[string][]RealServer
	idcWeightedCntMap            map[string]int
}

func (v *VipDomain) MarshalLogObject(oe sofalogger.ObjectEncoder) error {
	oe.AddString("name", v.Name)
	oe.AddInt32("protect_threshold", v.ProtectThreshold)
	oe.AddInt64("version", v.Version)
	oe.AddBool("is_deleted", v.IsDeleted)
	for i := range v.RealNodes {
		n := v.RealNodes[i]
		oe.AddString(fmt.Sprintf("nodes[%d]", i), n.String())
	}
	return nil
}

func (v *VipDomain) Polyfill(defaultPort int32) {
	for i := range v.RealNodes {
		healthCheckPort := v.RealNodes[i].HealthCheckPort
		if healthCheckPort == 0 {
			v.RealNodes[i].HealthCheckPort = defaultPort
		}
	}
}

func (v *VipDomain) ChecksumAlipay() string {
	var result int32 = 0

	result += round(int32(v.Version), 0)
	result += round(hashCodeForAlipayRealNodes(v.RealNodes), 1)

	return strconv.Itoa(int(result))
}

func (v *VipDomain) ChecksumCloud() string {
	var result int32 = 1

	result = (Prime * result) + hashCodeForString(v.Name)
	result = (Prime * result) + v.ProtectThreshold
	result = (Prime * result) + v.HealthCheckDefaultPort
	result = (Prime * result) + hashCodeForRealNodes(v.RealNodes)

	return strconv.Itoa(int(result))
}

func (v *VipDomain) GetZoneInfoList() (*ZoneInfoList, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	return v.ZoneInfoList, true
}

func (v *VipDomain) FromProtobuf(msg *protobuf.VipDomainMsg, zi *ZoneInfoList) {
	var realNodes []RealNode

	if len(msg.GetRealNodes()) > 0 {
		for _, realNodeMsg := range msg.GetRealNodes() {
			var realnode RealNode
			realnode.FromProtobuf(realNodeMsg)
			realNodes = append(realNodes, realnode)
		}
	}

	v.ZoneInfoList = zi
	v.Name = msg.GetName()
	v.ProtectThreshold = msg.GetProtectThreshold()
	v.HealthCheckDefaultPort = msg.GetHealthCheckDefaultPort()
	v.Version = msg.GetVersion()
	v.RealNodes = realNodes
	v.idcWeightedAvailableRsMap = make(map[string][]RealServer, 16)
	v.cityWeightedAvailableRsMap = make(map[string][]RealServer, 16)
	v.idcWeightedCntMap = make(map[string]int, 16)
	v.buildWeightedRealServers(realNodes)

	for i := range v.RealNodes {
		rs := v.RealNodes[i]
		if zi == nil {
			continue
		}
		zone := strings.ToLower(rs.Zone)
		zi, ok := zi.GetFromZone(zone)
		if !ok {
			continue
		}

		v.buildLocalLDCRealserverMap(zi.zone, zi.city)
	}
}

func (v *VipDomain) buildWeightedRealServers(realNodes []RealNode) {
	zi := v.ZoneInfoList
	realServers := make([]RealServer, 0, len(realNodes))
	for i := range realNodes {
		var (
			idc string
			rn  = realNodes[i]
		)
		if zi != nil {
			zi, ok := zi.GetFromZone(rn.Zone)
			if ok {
				idc = zi.idc
			}
		}
		rs := NewRealServerWithIDC(rn, idc)
		realServers = append(realServers, rs)
	}

	totalCount := len(realNodes)
	availableCount := 0

	for i := range realServers {
		rs := realServers[i]
		if rs.IsAvailable() {
			availableCount++
			v.weightedAvailableRealServers = addWithWeightRepetition(v.weightedAvailableRealServers, rs)
			v.availableRealServers = append(v.availableRealServers, rs)
		}
		v.weightedAllRealServers = append(v.weightedAllRealServers, rs)
		v.allRealServer = append(v.allRealServer, rs)
	}

	protectThreshold := int(v.ProtectThreshold)
	if protectThreshold > 0 {
		leastAvailableCount := math.Ceil(float64(protectThreshold) * float64(totalCount) / 100.0)
		if float64(availableCount) < leastAvailableCount {
			unavailableNodes := v.getUnavailableList(realServers, int(leastAvailableCount)-availableCount)
			for _, rn := range unavailableNodes {
				v.weightedAvailableRealServers = addWithWeightRepetition(v.weightedAvailableRealServers, rn)
				v.availableRealServers = append(v.availableRealServers, rn)
			}
		}
	} else {
		if availableCount <= 0 {
			for _, rs := range realServers {
				v.weightedAvailableRealServers = addWithWeightRepetition(v.weightedAvailableRealServers, rs)
				v.availableRealServers = append(v.availableRealServers, rs)
			}
		}
	}
}

func (v *VipDomain) getUnavailableList(rs []RealServer, count int) []RealServer {
	list := make([]RealServer, 0)
	for _, s := range rs {
		if !s.IsAvailable() {
			list = append(list, s)
		}
		if len(list) >= count {
			break
		}
	}
	return list
}

func (v *VipDomain) getLocalIDCRealServerList(idc string, allowCrossIdc, allowCrossCity bool) ([]RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	z, ok := v.ZoneInfoList.GetFromIDC(idc)
	if !ok {
		return nil, false
	}

	city := z.city

	list, ok := v.idcWeightedAvailableRsMap[idc]
	if !ok {
		// TODO:
		_ = ok
	}

	weight, ok := v.idcWeightedCntMap[idc]
	if !ok {
		weight = 0
	}

	var percent float64
	if weight > 0 {
		percent = float64(len(list)) * 100.0 / float64(weight)
	}

	if len(list) > 0 {
		if !allowCrossIdc {
			return list, true
		}

		// TODO: implement the getIdcDisasterProtect
		if percent >= 0 { // getIdcDisasterProtect
			return list, true
		}

		return v.getLocalCityRealServerList(city, allowCrossCity)
	}

	return nil, false
}

func (v *VipDomain) getLocalCityRealServerList(city string, allowCrossCity bool) ([]RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	list, _ := v.cityWeightedAvailableRsMap[city]

	if len(city) > 0 {
		return list, false
	}

	if allowCrossCity {
		// TODO: support drm control
		return v.weightedAvailableRealServers, true
	}

	return nil, false
}

func (v *VipDomain) GetAvailableRealAddress(defaultPort int) []string {
	servers := v.GetRealServers()
	result := make([]string, 0, len(servers))
	noAvailable := make([]string, 0, len(servers))

	for _, server := range servers {
		srv := fmt.Sprintf("%s:%d", server.GetIp(), defaultPort)
		if server.IsAvailable() {
			result = append(result, srv)
		} else {
			noAvailable = append(noAvailable, srv)
		}
	}

	length := float64(len(servers))
	needCount := float64(v.ProtectThreshold) / float64(100) * length
	for _, no := range noAvailable {
		if float64(len(result)) >= needCount {
			return result
		}
		result = append(result, no)
	}
	return result
}

func (v *VipDomain) GetLocalCityRealServerFromZone(zone string, allowCrossCity bool) (*RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	zi, ok := v.ZoneInfoList.GetFromZone(zone)
	if !ok {
		return nil, false
	}

	return v.GetLocalCityRealServer(zi.city, allowCrossCity)
}

func (v *VipDomain) GetLocalCityRealServer(city string, allowCrossCity bool) (*RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	list, ok := v.getLocalCityRealServerList(city, allowCrossCity)
	if !ok {
		return nil, false
	}

	if len(list) == 0 {
		return nil, false
	}

	return &list[rand.Intn(len(list))], true
}

func (v *VipDomain) GetLocalIDCRealServerFromZone(zone string, allowCrossIDC, allowCrossCity bool) (*RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	zi, ok := v.ZoneInfoList.GetFromZone(zone)
	if !ok {
		return nil, false
	}

	return v.GetLocalIDCRealServer(zi.idc, allowCrossIDC, allowCrossCity)
}

func (v *VipDomain) GetLocalIDCRealServer(idc string, allowCrossIDC, allowCrossCity bool) (*RealServer, bool) {
	if v.ZoneInfoList == nil {
		return nil, false
	}

	list, ok := v.getLocalIDCRealServerList(idc, allowCrossIDC, allowCrossCity)
	if !ok {
		return nil, false
	}

	if len(list) == 0 {
		return nil, false
	}

	return &list[rand.Intn(len(list))], true
}

func (v *VipDomain) GetRealServers() []RealServer {
	realServers := make([]RealServer, 0, len(v.RealNodes))
	for i := range v.RealNodes {
		realServers = append(realServers, NewRealServer(v.RealNodes[i]))
	}
	return realServers
}

func (v *VipDomain) buildLocalLDCRealserverMap(idc, city string) []RealServer {
	zi := v.ZoneInfoList
	if zi == nil {
		return nil
	}

	list := make([]RealServer, 0, 16)
	citylist := make([]RealServer, 0, 16)
	idcweight := 0

	for i := range v.RealNodes {
		r := v.RealNodes[i]
		z, ok := zi.GetFromZone(r.Zone)
		if !ok {
			continue
		}
		rscity := z.city
		rsidc := z.idc

		if strings.EqualFold(idc, rsidc) {
			idcweight += int(r.Weight)
		}

		if !r.Available {
			continue
		}

		if strings.EqualFold(idc, rsidc) {
			list = addWithWeightRepetition(list, NewRealServer(r))
			citylist = addWithWeightRepetition(citylist, NewRealServer(r))
		} else if strings.EqualFold(city, rscity) {
			citylist = addWithWeightRepetition(citylist, NewRealServer(r))
		}
	}

	v.idcWeightedAvailableRsMap[idc] = list
	v.idcWeightedCntMap[idc] = idcweight
	v.cityWeightedAvailableRsMap[city] = citylist

	return list
}

func addWithWeightRepetition(list []RealServer, rs RealServer) []RealServer {
	for i := 0; i < int(rs.weight); i++ {
		list = append(list, rs)
	}
	return list
}
