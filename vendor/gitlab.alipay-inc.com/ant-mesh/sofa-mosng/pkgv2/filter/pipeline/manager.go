package pipeline

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"sync"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/filter"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	event.EventListenerManagerInstance().RegisterList([]event.ResourceEventListener{
		event.ResourceEventListenerFuncs{
			Type:       api.GLOBAL_FILTER,
			AddFunc:    GetPipelineManagerInstance().AddOrUpdateGlobalFilter,
			UpdateFunc: GetPipelineManagerInstance().AddOrUpdateGlobalFilter,
			DeleteFunc: GetPipelineManagerInstance().DeleteGlobalFilter,
		},
		event.ResourceEventListenerFuncs{
			Type:       api.FILTER_CHAIN,
			AddFunc:    GetPipelineManagerInstance().AddOrUpdateFilterChain,
			UpdateFunc: GetPipelineManagerInstance().AddOrUpdateFilterChain,
			DeleteFunc: GetPipelineManagerInstance().DeleteChainFilter,
		},
	})
}

type Manager struct {
	// todo sync.Map?
	gMux           sync.RWMutex
	globalPipeline []types.GatewayFilter

	cMux           sync.RWMutex
	pipelineChains map[string][]types.GatewayFilter
}

func (pm *Manager) AddOrUpdateGlobalFilter(gf api.Object) (error, bool) {
	pm.gMux.Lock()
	defer pm.gMux.Unlock()

	if gf == nil {
		return nil, true
	}

	log.ConfigLogger().Infof("[gateway][pipeline][AddOrUpdateGlobalFilter] start add GlobalFilter")

	var gp []types.GatewayFilter

	for _, fConf := range gf.(v1.GlobalFilter).Filters {
		if gatewayFilter, e := filter.New(fConf); e == nil {
			gp = append(gp, gatewayFilter)
		} else {
			return errors.Errorf("[gateway][pipeline][AddOrUpdateGlobalFilter] build filter %s err, %v", fConf.Name, e), false
		}
	}

	pm.globalPipeline = gp

	return nil, true
}

func (pm *Manager) DeleteGlobalFilter(o api.Object) (error, bool) {
	pm.gMux.Lock()
	defer pm.gMux.Unlock()

	// todo
	return nil, true
}

func (pm *Manager) GetGlobalPipeline() []types.GatewayFilter {
	pm.gMux.RLock()
	defer pm.gMux.RUnlock()

	return pm.globalPipeline
}

func (pm *Manager) AddOrUpdateFilterChain(o api.Object) (error, bool) {
	pm.cMux.Lock()
	defer pm.cMux.Unlock()

	chain := o.(*v1.FilterChain)

	log.ConfigLogger().Infof("[gateway][pipeline][AddOrUpdateFilterChain] start add FilterChain %s", chain.ChainName)

	var filters []types.GatewayFilter

	for _, fConf := range chain.Filters {
		if gatewayFilter, e := filter.New(fConf); e == nil {
			filters = append(filters, gatewayFilter)
		} else {
			return errors.Errorf("[gateway][pipeline][AddOrUpdateFilterChain] build filter %s err, %v", fConf.Name, e), false
		}
	}

	pm.pipelineChains[chain.ChainName] = filters

	return nil, true
}

func (pm *Manager) DeleteChainFilter(o api.Object) (error, bool) {
	pm.cMux.Lock()
	defer pm.cMux.Unlock()

	// todo delete
	return nil, true
}

func (pm *Manager) GetPipelineChains(chains []string) (fs []types.GatewayFilter) {
	pm.cMux.RLock()
	defer pm.cMux.RUnlock()

	if chains == nil {
		return nil
	}

	for _, c := range chains {
		fs = append(fs, pm.pipelineChains[c]...)
	}

	return
}

func (pm *Manager) BuildPipeline(rConf *v1.Router, sConf *v1.GatewayService) types.Pipeline {
	// router's filterChain
	// todo
	rc := pm.GetPipelineChains(rConf.FilterChains)

	// router's service's filterChain
	sc := pm.GetPipelineChains(sConf.FilterChains)

	// global's filterChain
	gp := pm.GetGlobalPipeline()

	// build pipeline with router'filter rc service's filter sc and gp
	pipeline := NewPipeline(rConf.Filters, rc, sConf.Filters, sc, gp)

	// sort filter
	pipeline.Sort()

	return pipeline
}

var (
	once                    sync.Once
	pipelineManagerInstance *Manager
)

func GetPipelineManagerInstance() *Manager {
	once.Do(func() {
		pipelineManagerInstance = newPipelineManager()
	})

	return pipelineManagerInstance
}

func newPipelineManager() *Manager {
	return &Manager{
		globalPipeline: []types.GatewayFilter{},
		pipelineChains: make(map[string][]types.GatewayFilter),
	}
}
