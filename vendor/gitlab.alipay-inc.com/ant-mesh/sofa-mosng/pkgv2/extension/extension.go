package extension

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
)

func init() {
	event.EventListenerManagerInstance().Register(event.ResourceEventListenerFuncs{
		Type:       api.EXTENSION,
		AddFunc:    AddOrUpdateFilterExtension,
		UpdateFunc: AddOrUpdateFilterExtension,
	})
}

func AddOrUpdateFilterExtension(o api.Object) (error, bool) {
	ext := o.(v1.Extension)
	for _, fe := range ext.Filters {
		LoadFilterExt(fe)
	}

	return nil, true
}
