package dispatcher

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/valid"
)

const (
	commonLogFmtStr = "[dispatcher] [%s:%v] error: %v"
)

type dispatcher struct {
	validationManager    valid.ValidationManager
	eventListenerManager event.EventListenerManager
	store                config.Store
}

func (d *dispatcher) SetEventListenerManager(e event.EventListenerManager) {
	d.eventListenerManager = e
}

func (d *dispatcher) SetValidationManager(v valid.ValidationManager) {
	d.validationManager = v
}

func (d *dispatcher) SetStore(s config.Store) {
	d.store = s
}

func (d *dispatcher) Dispatch(e event.EventType, o api.Object) (error, bool) {
	defer func() {
		if err := recover(); err != nil {
			log.ConfigLogger().Errorf("[Dispatch] recover success. err %v", err)
		}
	}()
	// GW_CONFIG 需要分开处理
	if o.Type() == api.GW_CONFIG {
		gateway := o.(*v1.GwConfig)

		for _, rg := range gateway.GatewayMetadatas {
			if err, ok := d.doDispatch(e, rg); !ok {
				log.ConfigLogger().Errorf(commonLogFmtStr, e, gateway.GatewayMetadatas, err)
				return err, false
			}
		}

		if err, ok := d.doDispatch(e, gateway.Extension); !ok {
			log.ConfigLogger().Errorf(commonLogFmtStr, e, gateway.Extension, err)
			return err, false
		}

		for _, rg := range gateway.GatewayServices {
			if err, ok := d.doDispatch(e, rg); !ok {
				log.ConfigLogger().Errorf(commonLogFmtStr, e, gateway.GatewayServices, err)
				return err, false
			}
		}

		// filter 创建依赖 config
		if err, ok := d.doDispatch(e, gateway.GlobalFilters); !ok {
			log.ConfigLogger().Errorf(commonLogFmtStr, e, gateway.GlobalFilters, err)
			return err, false
		}

		for _, rg := range gateway.FilterChains {
			if err, ok := d.doDispatch(e, rg); !ok {
				log.ConfigLogger().Errorf(commonLogFmtStr, e, rg, err)
				return err, false
			}
		}

		// router 创建 pipeline 依赖 service、globalFilter 和 filterChain
		for _, rg := range gateway.RouterGroups {
			if err, ok := d.doDispatch(e, rg); !ok {
				log.ConfigLogger().Errorf(commonLogFmtStr, e, rg, err)
				return err, false
			}
		}

		// server 创建依赖 router
		for _, rg := range gateway.Gateways {
			if err, ok := d.doDispatch(e, rg); !ok {
				log.ConfigLogger().Errorf(commonLogFmtStr, e, rg, err)
				return err, false
			}
		}

	} else {
		if err, ok := d.doDispatch(e, o); !ok {
			log.ConfigLogger().Errorf(commonLogFmtStr, e, o, err)
			return err, false
		}
	}

	return nil, true
}

func (d *dispatcher) doDispatch(e event.EventType, o api.Object) (error, bool) {

	// valid
	if err, ok := d.validationManager.DoValid(o); !ok {
		return err, false
	}

	// diff
	diffEvent := d.store.Diff(o)

	// listener
	for _, e := range diffEvent {
		if err, ok := d.eventListenerManager.DoEvent(e.EventType, e.Object); !ok {
			return err, false
		}
	}

	// store
	d.store.Store(o)

	return nil, true
}

var (
	dispatcherIns config.Dispatcher = &dispatcher{
		validationManager:    valid.ValidationManagerInstance(),
		eventListenerManager: event.EventListenerManagerInstance(),
		store:                config.StoreInstance(),
	}
)

func DispatcherInstance() config.Dispatcher {
	return dispatcherIns
}
