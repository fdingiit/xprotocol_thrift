package resource

import (
	"encoding/json"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/dispatcher"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/resource/codec"
	"io/ioutil"

	"mosn.io/mosn/pkg/admin/server"
	"net/http"
)

var rest = &restResource{
	baseResource: baseResource{
		codec:      codec.YamlCodec(),
		dispatcher: dispatcher.DispatcherInstance(),
	},
}

func RestResource() *restResource {
	return rest
}

type restResource struct {
	baseResource
	// ip:port
	addr string
}

func (b *restResource) Start(addr string) {
	log.StartLogger().Infof("[mosng][resource][start] restResource start")
	server.RegisterAdminHandleFunc("/rest/local", localStart)
	server.RegisterAdminHandleFunc("/rest/apply", applyYaml)
}

func (b *restResource) Stop() {
	// no-op
}

type resp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
}

func applyYaml(w http.ResponseWriter, r *http.Request) {
	if con, err := ioutil.ReadAll(r.Body); err == nil {
		o, err := rest.codec.Decode(con)
		if err != nil {
			log.ConfigLogger().Errorf("parse request body err: %v", err)
			bytes, _ := json.Marshal(&resp{
				Success: false,
				Msg:     "parse request body err",
			})
			_, _ = w.Write(bytes)
			return
		}

		log.ConfigLogger().Infof("parse request body success: %s", string(con))

		if _, ok := rest.dispatcher.Dispatch(event.Update, o); ok {

			bytes, _ := json.Marshal(&resp{
				Success: true,
			})
			_, _ = w.Write(bytes)
		}
	}
}
