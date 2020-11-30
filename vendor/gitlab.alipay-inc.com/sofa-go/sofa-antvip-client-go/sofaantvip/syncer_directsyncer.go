package sofaantvip

import (
	"time"
)

type DirectSyncer struct {
	domains map[string]VipDomain
}

func NewDirectSyncer(domains map[string]VipDomain) *DirectSyncer {
	return &DirectSyncer{
		domains: domains,
	}
}

func (ds *DirectSyncer) GetVipDomainList(domains map[string]string,
	polling bool, timeout time.Duration) ([]VipDomain, error) {
	vipdomains := make([]VipDomain, 0, len(domains))
	for k := range domains {
		vd, ok := ds.domains[k]
		if !ok {
			continue
		}
		vipdomains = append(vipdomains, vd)
	}
	return vipdomains, nil
}
