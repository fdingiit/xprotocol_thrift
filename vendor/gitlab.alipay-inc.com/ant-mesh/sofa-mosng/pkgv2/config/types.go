package config

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/valid"
)

type Resource interface {
	SetCodec(Codec)
	SetDispatcher(Dispatcher)
}

type Codec interface {
	Decode(bytes []byte) (api.Object, error)
	Encode(o api.Object) ([]byte, error)
}

type Dispatcher interface {
	Dispatch(event.EventType, api.Object) (error, bool)
	SetEventListenerManager(event.EventListenerManager)
	SetValidationManager(valid.ValidationManager)
	SetStore(Store)
}



type Store interface {
	Diff(api.Object) []event.DifferEvent
	Store(api.Object)
	Get(t api.Type, name string) api.Object
	GetAll(t api.Type) interface{}
}
