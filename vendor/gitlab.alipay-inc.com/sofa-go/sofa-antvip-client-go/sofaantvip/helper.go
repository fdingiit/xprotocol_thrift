package sofaantvip

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const (
	Prime int32 = 31
)

var ipPatternRegex *regexp.Regexp
var domainPatternRegex *regexp.Regexp

func init() {
	ipPatternRegex, _ = regexp.Compile("^((((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\\:[0-9]+)?\\,?)+)$")
	domainPatternRegex, _ = regexp.Compile("^(([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(\\:[0-9]+)?\\,?)+)$")
}

func prettyDomains(domains map[string]string) string {
	var s strings.Builder
	for k := range domains {
		s.WriteString(k)
		s.WriteString("&")
	}

	w := s.String()
	if len(w) > 256 {
		return w[:256]
	}
	return w
}

func hashCodeForBool(ok bool) int32 {
	if ok {
		return 1231
	}
	return 1237
}

func hashCodeForRealNode(node RealNode) int32 {
	var result int32 = 1

	result = (Prime * result) + hashCodeForBool(node.Available)
	result = (Prime * result) + hashCodeForString(node.Ip)
	result = (Prime * result) + node.Weight
	result = (Prime * result) + node.HealthCheckPort

	return result
}

func hashCodeForString(s string) int32 {
	if s == "" {
		return 0
	}

	var h int32 = 0
	for i := 0; i < len(s); i++ {
		h = Prime*h + int32(s[i])
	}
	return h
}

func checksumStringSlice(nameList []string) string {
	if len(nameList) == 0 {
		return "N"
	}

	var result int32 = 1
	for _, name := range nameList {
		result = (Prime * result) + hashCodeForString(name)
	}

	return strconv.Itoa(int(result))
}

func hashCodeForRealNodes(nodes []RealNode) int32 {
	if len(nodes) == 0 {
		return 0
	}

	var result int32 = 1

	for _, node := range nodes {
		result = (Prime * result) + hashCodeForRealNode(node)
	}

	return result
}

func hashCodeForAlipayRealNodes(nodes []RealNode) int32 {
	if len(nodes) == 0 {
		return 0
	}

	var result int32 = 0

	for _, node := range nodes {
		result = (Prime * result) + hashCodeForAlipayRealNode(node)
	}

	return result
}

func hashCodeForAlipayRealNode(node RealNode) int32 {
	var result int32 = 0

	result += round(hashCodeForString(node.Ip), 0)
	result += round(hashCodeForBool(node.Available), 1)

	return result
}

func parseCommaSeparatedVipServers(csv string, checksum bool) (vipservers []VipServer, sum string, err error) {
	err = checkServers(csv)
	if err != nil {
		return nil, "", err
	}

	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil, "", errors.New("sofaantvip: nil comma separated servers")
	}

	servers := strings.Split(csv, ",")
	vipservers = make([]VipServer, 0, len(servers))

	for _, server := range servers {
		server = strings.TrimSpace(server)
		if strings.Index(server, ":") > 0 {
			server = strings.Split(server, ":")[0]
		}
		vipservers = append(vipservers, VipServer{
			Host: server,
		})
	}

	if checksum {
		return vipservers, checksumStringSlice(servers), nil
	}

	return vipservers, "", nil
}

func checkServers(servers string) error {

	var ipResult = false
	var domainResult = false

	if ipPatternRegex == nil && domainPatternRegex == nil {
		return nil
	}

	if ipPatternRegex != nil {
		ipResult = ipPatternRegex.MatchString(servers)
	}

	if domainPatternRegex != nil {
		domainResult = domainPatternRegex.MatchString(servers)
	}

	if ipResult || domainResult {
		return nil
	} else {
		return errors.New(fmt.Sprintf("antvip server list is illegal, result = %s", servers))
	}
}

func doHMACSha256AndBase64(secretkey []byte, in []byte) string {
	h := hmac.New(sha256.New, secretkey)
	// nolint
	h.Write(in)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func errstring(err error) string {
	if err == nil {
		return "nil"
	}

	return err.Error()
}

func round(factor int32, round int32) int32 {
	var result int32 = 1
	var i int32 = 0
	for ; i < round; i++ {
		result *= Prime
	}
	return factor * result
}

func merror(err error, errs ...error) error {
	merr := multierror.Append(err, errs...)
	merr.ErrorFormat = merrorfmt
	return merr
}

func merrorfmt(errs []error) string {
	var b strings.Builder
	for i, err := range errs {
		_, _ = b.WriteString("#")
		_, _ = b.WriteString(strconv.Itoa(i))
		_, _ = b.WriteString(err.Error())
	}
	return b.String()
}
