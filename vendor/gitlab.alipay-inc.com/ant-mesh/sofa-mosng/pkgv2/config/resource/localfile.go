package resource

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/dispatcher"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/resource/codec"
	"io/ioutil"
	"net/http"
)

var local = &localResource{
	baseResource: baseResource{
		codec:      codec.YamlCodec(),
		dispatcher: dispatcher.DispatcherInstance(),
	},
}

func LocalResource() *localResource {
	return local
}

func localStart(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	path := vars.Get("path")
	if path == "" {
		return
	}

	LocalResource().Start(path)

	bytes, _ := json.Marshal(&resp{
		Success: true,
	})
	_, _ = w.Write(bytes)
}

type localResource struct {
	baseResource
	filePath string
	md5      string
}

func (b *localResource) Start(filePath string) {
	log.ConfigLogger().Infof("[mosng][localResource][start] local start, filePath = %s", filePath)

	b.filePath = filePath
	bytes, _ := ioutil.ReadFile(filePath)
	b.resolve(bytes)
}

func (b *localResource) StartWithBytes(bytes []byte) {
	log.ConfigLogger().Infof("[mosng][localResource][start] local start, bytes = %s", bytes)

	b.resolve(bytes)
}

func (b *localResource) Stop() {
	// no-op
}

func (b *localResource) listen(filePath string) {
	bytes, _ := ioutil.ReadFile(filePath)
	b.resolve(bytes)
}

func (b *localResource) resolve(bytes []byte) {
	g, err := b.codec.Decode(bytes)

	if err != nil {
		log.ConfigLogger().Errorf("[mosng][localResource][resolve] resolve error: %v", err)
		return
	}

	log.ConfigLogger().Infof("[mosng][localResource][resolve] resolve success, config = %v", g)

	if _, ok := b.dispatcher.Dispatch(event.Update, g); !ok {
		log.ConfigLogger().Errorf("[mosng][localResource][resolve] dispatcher error: %v", err)
	}
}

func (b *localResource) checkMd5(bytes []byte) bool {
	sum := md5.Sum(bytes)

	md5str := fmt.Sprintf("%x", sum)

	if md5str == b.md5 {
		return true
	}

	b.md5 = md5str

	return false
}
