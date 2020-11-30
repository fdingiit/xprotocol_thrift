package sofaantvip

type RealServer struct {
	ip              string
	idc             string
	weight          int32
	healthCheckPort int32
	available       bool
	zone            string
	dataCenter      string
	labels          map[string]string
}

func NewRealServerFromDomain(domain VipDomain, realNode RealNode) RealServer {
	labels := make(map[string]string)
	if domain.Labels != nil {
		for k, v := range domain.Labels {
			labels[k] = v
		}
	}

	if realNode.Labels != nil {
		for k, v := range realNode.Labels {
			labels[k] = v
		}
	}
	rs := NewRealServer(realNode)
	rs.labels = labels

	return rs
}

func NewRealServerWithIDC(realNode RealNode, idc string) RealServer {
	rs := NewRealServer(realNode)
	rs.idc = idc
	return rs
}

func NewRealServer(realNode RealNode) RealServer {
	return RealServer{
		ip:              realNode.Ip,
		weight:          realNode.Weight,
		healthCheckPort: realNode.HealthCheckPort,
		available:       realNode.Available,
		zone:            realNode.Zone,
		dataCenter:      realNode.DataCenter,
		labels:          realNode.Labels,
	}
}

func (rs *RealServer) GetLabels() map[string]string {
	return rs.labels
}

func (rs *RealServer) GetIp() string {
	return rs.ip
}

func (rs *RealServer) GetWeight() int32 {
	return rs.weight
}

func (rs *RealServer) GetHealthCheckPort() int32 {
	return rs.healthCheckPort
}

func (rs *RealServer) IsAvailable() bool {
	return rs.available
}

func (rs *RealServer) GetZone() string {
	return rs.zone
}

func (rs *RealServer) GetDataCenter() string {
	return rs.dataCenter
}
