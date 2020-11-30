package sofaantvip

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var zeroVipServer VipServer

type VipServer struct {
	Host     string `json:"host"`
	HostName string `json:"hostName"`
	Weight   int32  `json:"weight"`
	IsLocal  bool   `json:"isLocal"`
}

type vipServers struct {
	created  time.Time
	checksum string
	servers  []VipServer
}

func newVipServers(servers []VipServer) vipServers {
	names := make([]string, 0, len(servers))
	for i := range servers {
		names = append(names, servers[i].Host)
	}

	return vipServers{
		servers:  servers,
		checksum: checksumStringSlice(names),
	}
}

func newVipServersFromCSV(csv string) (*vipServers, error) {
	servers, checksum, err := parseCommaSeparatedVipServers(csv, true)
	if err != nil {
		return nil, err
	}

	return &vipServers{
		created:  time.Now(),
		checksum: checksum,
		servers:  servers,
	}, nil
}

func (vs *vipServers) getRandomServer() (VipServer, bool) {
	servers := vs.servers
	if len(servers) > 0 {
		return servers[rand.Intn(len(servers))], true
	}
	return zeroVipServer, false
}

func (vs *vipServers) getChecksum() string { return vs.checksum }

func (vs *vipServers) getServers() []VipServer {
	return vs.servers
}
