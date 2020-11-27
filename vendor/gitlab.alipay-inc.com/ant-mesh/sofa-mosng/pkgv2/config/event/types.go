package event

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
)

type EventType string

type ResourceEventListener interface {
	ResourceType() api.Type
	OnAdd(obj api.Object) (error, bool)
	OnUpdate(obj api.Object) (error, bool)
	OnDelete(obj api.Object) (error, bool)
}

type ResourceEventListenerFuncs struct {
	Type       api.Type
	AddFunc    func(obj api.Object) (error, bool)
	UpdateFunc func(obj api.Object) (error, bool)
	DeleteFunc func(obj api.Object) (error, bool)
}

func (r ResourceEventListenerFuncs) ResourceType() api.Type {
	return r.Type
}

// OnAdd calls AddFunc if it's not nil.
func (r ResourceEventListenerFuncs) OnAdd(obj api.Object) (error, bool) {
	if r.AddFunc != nil {
		return r.AddFunc(obj)
	}

	return nil, true
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r ResourceEventListenerFuncs) OnUpdate(obj api.Object) (error, bool) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(obj)
	}

	return nil, true
}

// OnDelete calls DeleteFunc if it's not nil.
func (r ResourceEventListenerFuncs) OnDelete(obj api.Object) (error, bool) {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(obj)
	}

	return nil, true
}

type DifferEvent struct {
	Object    api.Object
	EventType EventType
}

type Differ interface {
	Diff(old api.Object) []DifferEvent
}

type EventListenerManager interface {
	Register(ResourceEventListener)
	RegisterList([]ResourceEventListener)
	DoEvent(EventType, api.Object) (error, bool)
}

const (
	Add    EventType = "add"
	Update EventType = "update"
	Delete EventType = "delete"
)
