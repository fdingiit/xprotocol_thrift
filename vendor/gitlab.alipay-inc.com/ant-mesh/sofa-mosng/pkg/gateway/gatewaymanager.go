package gateway

import (
	"os"
	"path"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/utils"
)

var (
	gatewayManagerInstance GatewayManager
	instanceMutex          = sync.Mutex{}
	instanceInited         int32
	dumpMutex              = sync.Mutex{}

	gatewayConfig GatewayConfig
)

type gatewayManager struct {
	services sync.Map
	apps     sync.Map
	clusters sync.Map
	controls sync.Map

	xdsServiceMap     sync.Map
	xdsAppMap         sync.Map
	xdsClusterMap     sync.Map
	xdsGatewayRuleMap sync.Map

	filterConfig *GatewayFilterConfig
}

func GetGatewayManager() GatewayManager {
	if atomic.LoadInt32(&instanceInited) == 1 {
		return gatewayManagerInstance
	}

	instanceMutex.Lock()
	defer instanceMutex.Unlock()

	if instanceInited == 0 {
		gatewayManagerInstance = initGatewayManager()
		atomic.StoreInt32(&instanceInited, 1)
	}

	return gatewayManagerInstance
}

func initGatewayManager() GatewayManager {
	gatewayManagerInstance = &gatewayManager{
		services: sync.Map{},
		apps:     sync.Map{},
		clusters: sync.Map{},
		controls: sync.Map{},

		xdsServiceMap:     sync.Map{},
		xdsAppMap:         sync.Map{},
		xdsClusterMap:     sync.Map{},
		xdsGatewayRuleMap: sync.Map{},
	}

	return gatewayManagerInstance
}

func (g *gatewayManager) AddOrUpdateGateway(gateway GatewayConfig) bool {

	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[gateway][conf] start to AddOrUpdateGateway,config=%+v", gateway)
	}

	gatewayConfig = gateway

	for _, scf := range gateway.ServiceConfigs {
		g.AddOrUpdateService(scf)
	}

	for _, gcf := range gateway.GatewayRules {
		g.AddOrUpdateGatewayRule(gcf)
	}

	for _, acf := range gateway.AppConfigs {
		g.AddOrUpdateApp(acf)
	}

	for _, ccf := range gateway.ClusterConfigs {
		g.AddOrUpdateCluster(ccf)
	}

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][conf] success to AddOrUpdateGateway,config=%+v", gateway)
	}

	return true
}

func (g *gatewayManager) AddOrUpdateService(scf ServiceConfig) bool {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[gateway][conf][service] start to AddOrUpdateService, config=%+v", scf)
	}

	api := &ApiService{}
	api.SetServiceKey(scf.Id)
	api.SetStatus(scf.Status)
	api.SetTimeout(scf.Timeout)

	for _, rule := range scf.Rules {
		if parser, exist := serviceRuleParserFactory[rule.Name]; exist {
			if attr, err := parser(rule.Config); err == nil {
				api.SetAttribute(rule.Name, attr)
			}
		}
	}

	var upt Upstream
	if upstreamParser, exist := upstreamParserFactory[string(scf.Upstream.Protocol)]; exist {
		upt = upstreamParser.parser(scf.Upstream)
	}
	api.SetUpstream(upt)

	g.services.Store(scf.Id, api)
	g.xdsServiceMap.Store(scf.Id, scf)

	for _, f := range configFilterFactory {
		f.OnAddOrUpdateService(api)
	}

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][conf][service] success to AddOrUpdateService, config=%+v", scf)
	}
	return true
}

func (g *gatewayManager) AddOrUpdateApp(app AppConfig) bool {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[gateway][conf][app] start to AddOrUpdateApp,config=%+v", app)
	}

	for k, v := range app.Spec {
		spec, _ := g.unmarshalAppSpec(k, v)
		app.Spec[k] = spec
	}
	g.apps.Store(app.Id, &app)
	g.xdsAppMap.Store(app.Id, app)

	for _, f := range configFilterFactory {
		f.OnAddOrUpdateApp(app)
	}

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][conf][app] success to AddOrUpdateApp, config=%+v", app)
	}
	return true
}

func (g *gatewayManager) AddOrUpdateCluster(cluster ClusterConfig) bool {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[gateway][conf][cluster] start to AddOrUpdateCluster,config=%+v", cluster)
	}

	for k, v := range cluster.Spec {
		spec, _ := g.unmarshalClusterSpec(k, v)
		cluster.Spec[k] = spec
	}

	g.clusters.Store(cluster.Id, &cluster)
	g.xdsClusterMap.Store(cluster.Id, cluster)

	for _, f := range configFilterFactory {
		f.OnAddOrUpdateCluster(cluster)
	}

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][conf][cluster] success to AddOrUpdateCluster, config=%+v", cluster)
	}
	return true
}

func (g *gatewayManager) AddOrUpdateGatewayRule(control GatewayRule) bool {
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("[gateway][conf][rule] start to AddOrUpdateGatewayRule, config=%+v", control)
	}

	conf, _ := g.unmarshalGatewayRule(control.Name, control.Config)
	g.notifyConfig(control.Name, conf)
	g.controls.Store(control.Name, conf)

	g.xdsGatewayRuleMap.Store(control.Name, control)

	for _, f := range configFilterFactory {
		f.OnAddOrUpdateGatewayRule(control)
	}

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][conf][rule] success to AddOrUpdateGatewayRule, config=%+v", control)
	}
	return true
}

func (g *gatewayManager) GetService(id string) Service {
	if item, exist := g.services.Load(id); exist {
		if s, ok := item.(Service); ok {
			return s
		}
	}
	return nil
}

func (g *gatewayManager) GetApp(appName string) *AppConfig {
	if item, exist := g.apps.Load(appName); exist {
		if s, ok := item.(*AppConfig); ok {
			return s
		}
	}
	return nil
}

func (g *gatewayManager) GetCluster(clusterName string) *ClusterConfig {
	if item, exist := g.clusters.Load(clusterName); exist {
		if s, ok := item.(*ClusterConfig); ok {
			return s
		}
	}
	return nil
}

func (g *gatewayManager) GetGatewayRule(controlName string) interface{} {
	if item, exist := g.controls.Load(controlName); exist {
		return item
	}
	return nil
}

func (g *gatewayManager) RemoveService(service ServiceConfig) bool {
	g.services.Delete(service.Id)
	g.xdsServiceMap.Delete(service.Id)
	return true
}

func (g *gatewayManager) RemoveApp(app AppConfig) bool {
	g.apps.Delete(app.Name)
	g.xdsAppMap.Delete(app.Name)
	return true
}

func (g *gatewayManager) RemoveCluster(cluster ClusterConfig) bool {
	g.clusters.Delete(cluster.Name)
	g.xdsClusterMap.Delete(cluster.Name)
	return true
}

func (g *gatewayManager) RemoveGatewayRule(control GatewayRule) bool {
	g.controls.Delete(control.Name)
	g.xdsGatewayRuleMap.Delete(control.Name)
	return true
}

func (g *gatewayManager) unmarshalGatewayRule(name string, config interface{}) (interface{}, bool) {
	return unmarshalByType(gatewayRuleTypeFactory, name, config)
}

func (g *gatewayManager) unmarshalAppSpec(name string, config interface{}) (interface{}, bool) {
	return unmarshalByType(appSpecTypeFactory, name, config)
}

func (g *gatewayManager) unmarshalClusterSpec(name string, config interface{}) (interface{}, bool) {
	return unmarshalByType(clusterSpecTypeFactory, name, config)
}

func unmarshalByType(typeMap sync.Map, specName string, config interface{}) (interface{}, bool) {
	if parser, exist := typeMap.Load(specName); exist {
		if t, ok := parser.(reflect.Type); ok {
			spec := reflect.New(t).Interface()
			if _, ok := config.(string); !ok {
				config, _ = json.MarshalToString(config)
			}
			if err := json.UnmarshalFromString(config.(string), &spec); err == nil {
				return spec, true
			}
		}
	}
	return config, true
}

func (g *gatewayManager) notifyConfig(name string, config interface{}) {
	if item, exist := configListenerFactory.Load(name); exist {
		if listeners, ok := item.([]ConfigListener); ok {
			for _, listener := range listeners {
				listener.Update(name, config)
			}
		}
	}
}

func (g *gatewayManager) startDumpTimer() {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for t := range ticker.C {
			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][conf][dump] start to dump gateway config at %s", t.String())
			}
			g.Dump()
		}
	}()
}

func (g *gatewayManager) SetFilterConfig(conf *GatewayFilterConfig) {
	g.filterConfig = conf
}

func (g *gatewayManager) GetFilterConfig() *GatewayFilterConfig {
	return g.filterConfig
}

func (g *gatewayManager) Dump() {
	dumpMutex.Lock()
	defer dumpMutex.Unlock()

	if g.filterConfig == nil || len(g.filterConfig.ConfigPath) < 1 {
		if log.DefaultLogger.GetLogLevel() >= log.WARN {
			log.DefaultLogger.Warnf("[gateway][conf][dump] dump warn: gateway config file should be declared")
		}
		return
	}

	preConfig := gatewayConfig

	services := make([]ServiceConfig, 0)
	g.xdsServiceMap.Range(func(key, value interface{}) bool {
		if conf, ok := value.(ServiceConfig); ok {
			services = append(services, conf)
		}
		return true
	})

	apps := make([]AppConfig, 0)
	g.xdsAppMap.Range(func(key, value interface{}) bool {
		if conf, ok := value.(AppConfig); ok {
			apps = append(apps, conf)
		}
		return true
	})

	clusters := make([]ClusterConfig, 0)
	g.xdsClusterMap.Range(func(key, value interface{}) bool {
		if conf, ok := value.(ClusterConfig); ok {
			clusters = append(clusters, conf)
		}
		return true
	})

	controls := make([]GatewayRule, 0)
	g.xdsGatewayRuleMap.Range(func(key, value interface{}) bool {
		if conf, ok := value.(GatewayRule); ok {
			controls = append(controls, conf)
		}
		return true
	})

	gatewayConfig.ServiceConfigs = services
	gatewayConfig.AppConfigs = apps
	gatewayConfig.ClusterConfigs = clusters
	gatewayConfig.GatewayRules = controls

	if !reflect.DeepEqual(preConfig, gatewayConfig) {
		if log.DefaultLogger.GetLogLevel() >= log.INFO {
			log.DefaultLogger.Infof("[gateway][conf][dump] gateway config changed, now dump.")
		}

		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("[gateway][conf][dump] gateway config: %+v", gatewayConfig)
		}

		if data, err := json.Marshal(&gatewayConfig); err == nil {
			writeFile(g.filterConfig.ConfigPath, data)
		}
	} else {
		if log.DefaultLogger.GetLogLevel() >= log.INFO {
			log.DefaultLogger.Infof("[gateway][conf][dump] gateway config not changed.")
		}
	}
}

func writeFile(fileName string, fileData []byte) {
	dir := path.Dir(fileName)
	if isExist, _ := pathExists(dir); !isExist {
		wErr := os.MkdirAll(dir, 0755)

		if wErr != nil {
			if log.DefaultLogger.GetLogLevel() >= log.ERROR {
				log.DefaultLogger.Errorf("[gateway][conf] mkdir %s failed: %s", dir, wErr.Error())
			}
		}
	}

	utils.WriteFileSafety(fileName, fileData, 0644)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
