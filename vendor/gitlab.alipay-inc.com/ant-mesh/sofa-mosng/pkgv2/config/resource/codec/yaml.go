package codec

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gopkg.in/yaml.v2"
)

type yamlCodec struct {
}

func YamlCodec() *yamlCodec {
	return &yamlCodec{}
}

func (c yamlCodec) Decode(bytes []byte) (api.Object, error) {
	m := metadata{}

	// todo refactor
	err := yaml.Unmarshal(bytes, &m)
	if err != nil {
		log.ConfigLogger().Errorf("[gateway][codec][yaml] yaml decode metadata failed: %s", err.Error())
		return nil, err
	}

	out, err := yaml.Marshal(m.Spec)
	if err != nil {
		log.ConfigLogger().Errorf("[gateway][codec][yaml] yaml encode spec failed: %s", err.Error())
		return nil, err
	}

	var object api.Object
	if m.Version == api.V1Version || m.Version == api.K8S_BETA1_VERSION {
		object = v1.NewObject(api.GetType(m.Kind))
	} else {
		panic("api version not support " + m.Version)
	}

	err = yaml.Unmarshal(out, object)
	if err != nil {
		log.ConfigLogger().Errorf("[gateway][codec][yaml] yaml decode %s failed: %s", m.Kind, err.Error())
		return nil, err
	}

	log.ConfigLogger().Infof("[gateway][codec][yaml] yaml decode success")
	return object, err
}

func (c yamlCodec) Encode(o api.Object) ([]byte, error) {
	return nil, nil
}
