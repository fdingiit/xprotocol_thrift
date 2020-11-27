package sofaantvip

import (
	"strings"
	"sync"
)

type ZoneInfo struct {
	zone string
	idc  string
	city string
}

type ZoneInfoList struct {
	sync.RWMutex
	zm       map[string]*ZoneInfo
	im       map[string]*ZoneInfo
	raw      string
	checksum string
}

func NewZoneInfoList(s string) *ZoneInfoList {
	z := &ZoneInfoList{
		zm:  make(map[string]*ZoneInfo, 128),
		im:  make(map[string]*ZoneInfo, 128),
		raw: s,
	}
	z.Lock()
	defer z.Unlock()

	ss := strings.Split(s, ",")

	for i := range ss {
		zs := ss[i]
		zss := strings.Split(zs, ":")
		if len(zss) < 3 { // invalid format
			continue
		}

		zi := &ZoneInfo{
			zone: strings.ToLower(zss[0]),
			idc:  strings.ToLower(zss[1]),
			city: strings.ToLower(zss[2]),
		}

		if zi.zone != "" {
			z.zm[zi.zone] = zi
		}

		if zi.idc != "" {
			z.im[zi.idc] = zi
		}
	}

	z.checksum = checksumStringSlice(ss)

	return z
}

func (z *ZoneInfoList) GetFromZone(zone string) (zi *ZoneInfo, ok bool) {
	z.RLock()
	defer z.RUnlock()
	zi, ok = z.zm[strings.ToLower(zone)]
	return zi, ok
}

func (z *ZoneInfoList) GetFromIDC(idc string) (zi *ZoneInfo, ok bool) {
	z.RLock()
	defer z.RUnlock()
	zi, ok = z.im[strings.ToLower(idc)]
	return zi, ok
}

func (z *ZoneInfoList) ChecksumAlipay() string {
	z.RLock()
	defer z.RUnlock()
	return z.checksum
}
