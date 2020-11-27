package sofaantvip

import (
	"time"
)

type Syncer interface {
	// GetVipDomainList gets a list of real servers from antvip server via the requested domains.
	// domains holds the interested domain defined as the following:
	// map[string]string{
	// 	"example.com": "N",
	// 	"a.example.com": "checksum",
	// }
	GetVipDomainList(domains map[string]string, polling bool, timeout time.Duration) ([]VipDomain, error)
}
