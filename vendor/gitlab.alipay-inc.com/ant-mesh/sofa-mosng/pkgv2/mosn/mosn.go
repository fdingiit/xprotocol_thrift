package mosn

import (
	"encoding/json"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/router"
	"mosn.io/mosn/pkg/server"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
)

func ConvertInit() {
	event.EventListenerManagerInstance().RegisterList([]event.ResourceEventListener{
		event.ResourceEventListenerFuncs{
			Type:       api.GATEWAY,
			AddFunc:    addOrUpdateListener,
			UpdateFunc: addOrUpdateListener,
		},
		event.ResourceEventListenerFuncs{
			Type:       api.SERVICE,
			AddFunc:    addOrUpdateCluster,
			UpdateFunc: addOrUpdateCluster,
		},
		event.ResourceEventListenerFuncs{
			Type:       api.ROUTER,
			AddFunc:    addOrUpdateRouter,
			UpdateFunc: addOrUpdateRouter,
		},
		event.ResourceEventListenerFuncs{
			Type:       api.ROUTER_GROUP,
			AddFunc:    addOrUpdateRouterGroup,
			UpdateFunc: addOrUpdateRouterGroup,
		},
	})
}

func addOrUpdateListener(obj api.Object) (error, bool) {
	mosnListener := convertListenerConfig(obj.(*v1.Gateway))
	addOrUpdateMosnListener(mosnListener)
	return nil, true
}

func addOrUpdateCluster(s api.Object) (error, bool) {
	cluster := convertClustersConfig(s.(*v1.GatewayService))

	addOrUpdateMosnCluster(cluster)

	return nil, true
}

func addOrUpdateRouter(o api.Object) (error, bool) {

	//router = o.(*v1.Router)

	// todo
	return nil, true
}

func addOrUpdateRouterGroup(o api.Object) (error, bool) {
	mosnRouter := convertRouterGroupConfig(o.(*v1.RouterGroup))
	addOrUpdateMosnRouter(mosnRouter)
	return nil, true
}

func addOrUpdateMosnListener(mosnListener *v2.Listener) {

	if mosnListener == nil {
		return
	}

	listenerAdapter := server.GetListenerAdapterInstance()
	if listenerAdapter == nil {
		// if listenerAdapter is nil, return directly
		log.ConfigLogger().Errorf("listenerAdapter is nil and hasn't been initiated at this time")
		return
	}
	log.ConfigLogger().Infof("listenerAdapter.AddOrUpdateListener called, with mosn listener:%s", mosnListener.Name)

	if err := listenerAdapter.AddOrUpdateListener("", mosnListener); err == nil {
		log.ConfigLogger().Infof("AddOrUpdateListener success,listener address = %s", mosnListener.Addr.String())
	} else {
		log.ConfigLogger().Errorf("AddOrUpdateListener failure,listener address = %s, msg = %s ",
			mosnListener.Addr.String(), err.Error())
	}
}

func addOrUpdateMosnCluster(cluster *v2.Cluster) {
	var err error
	if cluster.ClusterType == v2.EDS_CLUSTER {
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerClusterAddOrUpdate(*cluster)
	} else {
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerClusterAndHostsAddOrUpdate(*cluster, cluster.Hosts)
	}

	if err != nil {
		log.ConfigLogger().Errorf("OnUpdateClusters failed,cluster name = %s, error: %v", cluster.Name, err.Error())

	} else {
		log.ConfigLogger().Infof("OnUpdateClusters success,cluster name = %s", cluster.Name)
	}
}

func addOrUpdateMosnRouter(mosnRouter *v2.RouterConfiguration) {
	if mosnRouter == nil {
		return
	}

	if routersMngIns := router.GetRoutersMangerInstance(); routersMngIns == nil {
		log.ConfigLogger().Errorf("xds OnAddOrUpdateRouters error: router manager in nil")
	} else {
		if jsonStr, err := json.Marshal(mosnRouter); err == nil {
			log.ConfigLogger().Infof("raw router config: %s", string(jsonStr))
		}

		log.ConfigLogger().Infof("mosnRouter config: %v", mosnRouter)
		routersMngIns.AddOrUpdateRouters(mosnRouter)
	}
}
