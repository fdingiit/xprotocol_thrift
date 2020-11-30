package event

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
)

var (
	managerInstance EventListenerManager = &eventListenerManager{
		make(map[api.Type][]ResourceEventListener),
	}
)

func EventListenerManagerInstance() EventListenerManager {
	return managerInstance
}

type eventListenerManager struct {
	// 初始化写，运行时读, 不需要锁
	lMap map[api.Type][]ResourceEventListener
}

func (e *eventListenerManager) Register(l ResourceEventListener) {
	e.lMap[l.ResourceType()] = append(e.lMap[l.ResourceType()], l)
}

func (e *eventListenerManager) RegisterList(ls []ResourceEventListener) {
	for _, l := range ls {
		e.Register(l)
	}
}

func (e *eventListenerManager) DoEvent(et EventType, o api.Object) (error, bool) {
	for _, l := range e.lMap[o.Type()] {
		var err error
		var ok bool
		switch et {
		case Add:
			err, ok = l.OnAdd(o)
		case Update:
			err, ok = l.OnUpdate(o)
		case Delete:
			err, ok = l.OnDelete(o)
		}

		if !ok {
			return err, false
		}
	}

	return nil, true
}
