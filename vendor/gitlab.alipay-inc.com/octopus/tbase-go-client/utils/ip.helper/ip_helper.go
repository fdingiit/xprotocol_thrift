package ip_helper

import (
	"encoding/binary"
	"net"
	"strconv"
)

func Uint32ip(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip.String()
}

func CombineIpAndPort(address string, port uint32) string {
	return net.JoinHostPort(address, strconv.FormatUint(uint64(port), 10))
}
