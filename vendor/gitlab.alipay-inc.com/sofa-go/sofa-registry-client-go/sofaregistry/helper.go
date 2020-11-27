package sofaregistry

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	registryproto "gitlab.alipay-inc.com/sofa-go/sofa-registry-proto-go/proto"
)

const (
	HMACSHA256           = "HmacSHA256"
	DefaultCacheInterval = 5 * time.Minute
)

func errstring(err error) string {
	if err == nil {
		return "ok"
	}

	return err.Error()
}

func prettyReceivedDataPb(r *registryproto.ReceivedDataPb) string {
	// nolint
	// DataId              string                  `protobuf:"bytes,1,opt,name=dataId" json:"dataId,omitempty"`
	// Group               string                  `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// InstanceId          string                  `protobuf:"bytes,3,opt,name=instanceId" json:"instanceId,omitempty"`
	// Segment             string                  `protobuf:"bytes,4,opt,name=segment" json:"segment,omitempty"`
	// Scope               string                  `protobuf:"bytes,5,opt,name=scope" json:"scope,omitempty"`
	// SubscriberRegistIds []string                `protobuf:"bytes,6,rep,name=subscriberRegistIds" json:"subscriberRegistIds,omitempty"`
	// Data                map[string]*DataBoxesPb `protobuf:"bytes,7,rep,name=data" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Version             int64                   `protobuf:"varint,8,opt,name=version" json:"version,omitempty"`
	// LocalZone           string                  `protobuf:"bytes,9,opt,name=localZone" json:"localZone,omitempty"`
	data := fmt.Sprintf("%v", r.Data)
	if len(data) > 512 {
		data = data[:512] + fmt.Sprintf("...more(%d)", len(data))
	}

	return fmt.Sprintf("(DataID:%s Group:%s InstanceID:%s Segment:%s Scope:%s"+
		" SubscriberRegistIds:%v Version:%d LocalZone:%s data:%s)",
		r.DataId, r.Group, r.InstanceId, r.Segment, r.Scope, r.SubscriberRegistIds, r.Version, r.LocalZone, data)
}

func dataList2DataBoxesPb(dl []string) []*registryproto.DataBoxPb {
	db := make([]*registryproto.DataBoxPb, 0, len(dl))
	for _, s := range dl {
		db = append(db, &registryproto.DataBoxPb{Data: s})
	}
	return db
}

func dataBoxesPb2DataList(db []*registryproto.DataBoxPb) []string {
	dl := make([]string, 0, len(db))
	for _, b := range db {
		dl = append(dl, b.Data)
	}
	return dl
}

func doHMACSHA256Base64(key, plaintext string) string {
	hash := hmac.New(sha256.New, []byte(key))
	// nolint
	hash.Write([]byte(plaintext))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func getSignatureMap(accesskey, secretkey, instanceID string) map[string]string {
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	cacheTime := int64(DefaultCacheInterval / time.Millisecond)
	timestamp := strconv.FormatInt(ts/cacheTime*cacheTime, 10)
	plaintext := fmt.Sprintf("%s%s", instanceID, timestamp)
	return map[string]string{
		"!AccessKey": accesskey,
		"!Algothrim": HMACSHA256,
		"!Signature": doHMACSHA256Base64(secretkey, plaintext),
		"!Timestamp": timestamp,
	}
}
