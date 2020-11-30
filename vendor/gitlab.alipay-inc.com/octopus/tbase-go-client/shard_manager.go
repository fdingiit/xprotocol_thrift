package tbasego

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
	map_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/map.helper"

	rpc_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/rpc.helper"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"gitlab.alipay-inc.com/octopus/tbase-go-client/model"
	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	alibaba_coordinator_proto "gitlab.alipay-inc.com/octopus/tbase-go-client/pb/coordinator"
	convertor_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/convertor.helper"
	ip_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/ip.helper"
	message_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/message.helper"
	murmur_hash_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/murmur.hash.helper"
	pb_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/pb.helper"
	uuid_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/uuid.helper"
	"mosn.io/pkg/log"
)

const (
	// !!! DON'T CHANGE !!!
	// !!! DON'T CHANGE !!!
	// !!! DON'T CHANGE !!!
	DEFAULT_HASH_SEED = 0x9747b28c
)

var shardManagerCounter int32 = 0

type ShardManager struct {
	InstanceName              string
	ConnectionInfo            *model.ConnectionInfo
	CachedCoordinatorList     []string
	NextCoordinatorIndex      int
	coordinatorListUpdateTime int64
	IsShutdown                bool
	Version                   int32
	EndpointListIndex         int
	Endpoints                 *map[string]string
	endPointsPtr              *unsafe.Pointer //can't be visited outside shard manager
	endpointPointerAddress    int64
	Shard2EndpointId          *[]string
	shard2EndpointIdPtr       *unsafe.Pointer //can't be visited outside shard manager
	Shard2Endpoint            *[]string
	shard2EndpointPtr         *unsafe.Pointer //can't be visited outside shard manager
	Servers                   []string
	LastForceRefreshTime      int64
	RefreshLayoutChannel      chan *model.RefreshLayoutParam
	RefreshClientPoolChannel  chan *model.OldAndNewEndpoints
	wg                        *sync.WaitGroup
}

func NewShardManager(connInfo *model.ConnectionInfo) (*ShardManager, error) {

	// this is for ShardManager independent use
	tbase_log.InitTBaseLogger()

	instanceName := fmt.Sprintf("[SHARD_MANAGER]-(%v.%v#%v)", connInfo.Cluster, connInfo.Tenant, atomic.AddInt32(&shardManagerCounter, 1))
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] init shard manager %v with connection string %v", instanceName, connInfo.ToString())
	}

	shardManager := &ShardManager{
		InstanceName:              instanceName,
		ConnectionInfo:            connInfo,
		CachedCoordinatorList:     nil,
		NextCoordinatorIndex:      0,
		coordinatorListUpdateTime: 0,
		IsShutdown:                false,
		Version:                   -1,
		EndpointListIndex:         0,
		Servers:                   connInfo.Servers,
		LastForceRefreshTime:      0,
		RefreshLayoutChannel:      make(chan *model.RefreshLayoutParam, connInfo.MaxQueueSize),
		RefreshClientPoolChannel:  make(chan *model.OldAndNewEndpoints, 1),
		wg: new(sync.WaitGroup),
	}

	shardManager.wg.Add(1)
	shardManager.endPointsPtr = (*unsafe.Pointer)(unsafe.Pointer(&shardManager.Endpoints))
	shardManager.shard2EndpointIdPtr = (*unsafe.Pointer)(unsafe.Pointer(&shardManager.Shard2EndpointId))
	shardManager.shard2EndpointPtr = (*unsafe.Pointer)(unsafe.Pointer(&shardManager.Shard2Endpoint))

	shardManager.Servers = convertor_helper.Shuffle(shardManager.Servers)
	err := shardManager.refreshInternal(true)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get cluster layout error. error: %v", instanceName, err)
		return nil, err
	}
	shardManager.LastForceRefreshTime = time.Now().UnixNano()

	go func() {
		defer shardManager.wg.Done()
		shardManager.StartClusterLayoutUpdater()
	}()

	return shardManager, nil
}

func (shardManager *ShardManager) StartClusterLayoutUpdater() {
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v start cluster layout updater", shardManager.InstanceName)
	}
	for !shardManager.IsShutdown {
		select {

		// refresh layout periodically
		case <-time.After(time.Duration(shardManager.ConnectionInfo.LayoutRefreshInterval) * time.Millisecond):
			if shardManager.IsShutdown {
				break
			}

			// discard the error
			shardManager.DoUpdateClusterLayout(&model.RefreshLayoutParam{ForceRefresh: false, LastForceRefreshTime: atomic.LoadInt64(&shardManager.LastForceRefreshTime)})

		// refresh actively
		case refreshLayoutParam, ok := <-shardManager.RefreshLayoutChannel:
			if shardManager.IsShutdown {
				break
			}

			if !ok {
				tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v refresh cluster layout error. error: channel is closed. ", shardManager.InstanceName)
				return
			}

			// discard the error
			shardManager.DoUpdateClusterLayout(refreshLayoutParam)

		}
	}
}

func (shardManager *ShardManager) DoUpdateClusterLayout(param *model.RefreshLayoutParam) error {
	forceRefresh := param.ForceRefresh
	lastForceRefreshTime := param.LastForceRefreshTime

	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v try to update cluster layout, forceRefresh: %v", shardManager.InstanceName, forceRefresh)
	}
	if forceRefresh {
		elapsed := time.Now().UnixNano() - lastForceRefreshTime
		if elapsed >= shardManager.ConnectionInfo.MinimalLayoutRefreshTimespan {
			err := shardManager.refreshInternal(true)
			atomic.StoreInt64(&shardManager.LastForceRefreshTime, time.Now().UnixNano())
			return err
		} else {
			tbase_log.TBaseLogger.Warnf("[SHARD_MANAGER] %v update cluster layout request too often, "+
				"last update time is %v", shardManager.InstanceName, lastForceRefreshTime)
			return nil
		}
	} else {
		err := shardManager.refreshInternal(false)
		return err
	}
}

func (shardManager *ShardManager) GetClusterLayout(server string) (*alibaba_coordinator_proto.GetLayoutResponse, error) {
	requestId, err := uuid_helper.GenerateUUID()
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v generate uuid error. error: %v", shardManager.InstanceName, err)
		return nil, err
	}

	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v try get cluster layout from coordinator server %v. request params: coordinator Servers count=%v, requestId=%v, cluster=%v, tenant=%v, Version=%v",
			shardManager.InstanceName, server, shardManager.NextCoordinatorIndex, requestId, shardManager.ConnectionInfo.Cluster, shardManager.ConnectionInfo.Tenant, shardManager.Version)
	}

	getClusterLayoutRequest := pb_helper.ConstructGetClusterLayoutRequestMessage(requestId, shardManager.ConnectionInfo.Cluster, shardManager.ConnectionInfo.Tenant, shardManager.Version)
	requestBuf, err := message_helper.Encode(getClusterLayoutRequest)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v encode pb message to byte buffer error. error: %v", shardManager.InstanceName, err)
		return nil, err
	}

	timeout := time.Duration(shardManager.ConnectionInfo.CoordinatorTimeout) * time.Millisecond
	conn, err := rpc_helper.CreateConnectionAndSendRequest(requestId, server, timeout, requestBuf)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v create connection to server %v with timeout %v and send request error. error: %v",
			shardManager.InstanceName, server, timeout, err)
		return nil, err
	}
	defer conn.Close()

	getLayoutResponse, err := rpc_helper.GetResponse(conn, requestId, alibaba_coordinator_proto.E_GetLayoutResponse_Response)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get response error. dest server: %v, error: %v", shardManager.InstanceName, server, err)
	}

	if result, ok := getLayoutResponse.(*alibaba_coordinator_proto.GetLayoutResponse); ok {
		return result, nil
	} else {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v can't convert to \"GetLayoutResponse\". raw data: %v", shardManager.InstanceName, getLayoutResponse)
		return nil, error2.NewTBaseClientInternalError("can't convert to \"GetLayoutResponse\"")
	}
}

func (shardManager *ShardManager) RebuildLayout(layoutResponse *alibaba_coordinator_proto.GetLayoutResponse,
	coordinator string) error {
	if layoutResponse == nil {
		tbase_log.TBaseLogger.Errorf("[SHARD-MANAGER] %v try rebuild cluster layout from coordinator: %v error, "+
			"layout response is nil", coordinator, shardManager.InstanceName)
		return error2.NewTBaseClientInternalError("layout response is nil")
	}

	if len(layoutResponse.Layout) <= 0 {
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[SHARD-MANAGER] %v client layout version is the same as coordinator's. "+
				"coordinator: %v", shardManager.InstanceName, coordinator)
		}
		return error2.NewTBaseClientInternalError("client layout version is the same as server's")
	}

	unCompressedLayout, err := snappy.Decode(nil, layoutResponse.Layout)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v uncompress layout from coordinator %v error %v",
			shardManager.InstanceName, coordinator, err)
		return error2.NewTBaseClientInternalError(fmt.Sprintf("uncompress layout bytes error, error: %v", err))
	}

	layout := &alibaba_coordinator_proto.Layout{}
	err = proto.Unmarshal(unCompressedLayout, layout)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v bytes to \"Layout\" unmarshal error. "+
			"error: %s, coordinator: %v", shardManager.InstanceName, err, coordinator)
		return err
	}
	if layout == nil || layout.LayoutInfo == nil || layout.LayoutInfo.Version == nil {
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[SHARD-MANAGER] %v unmarshaled layout is empty, "+
				"local version %v is newest. coordinator: %v", shardManager.InstanceName, shardManager.Version, coordinator)
		}
		return error2.NewTBaseClientInternalError(fmt.Sprintf("unmarshaled layout is empty, local version %v is newest", shardManager.Version))
	}

	newVersion := layout.LayoutInfo.GetVersion()
	if newVersion <= shardManager.Version {
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[SHARD-MANAGER] %v local layout Version %v is newer "+
				"than version %v from coordinator %v , response discard.", shardManager.InstanceName,
				shardManager.Version, newVersion, coordinator)
		}
		return error2.NewTBaseClientInternalError(fmt.Sprintf("local layout Version %v is newer than Version %v", shardManager.Version, newVersion))
	}

	shardToEpsList := layout.ShardToEpsList
	newEndPoints, err := shardManager.buildEndpoints(shardToEpsList)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD-MANAGER] %v build endpoints error. coordinator: %v, error: %v",
			coordinator, shardManager.InstanceName, err)
		return err
	}
	shardNum := layout.GetLayoutInfo().GetNumShard()
	if shardNum <= 0 {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v layout %v's shard num field is invalid. coordinator: %v",
			shardManager.InstanceName, layout, coordinator)
		return error2.NewTBaseClientInternalError("layout 's shard num field is invalid")
	}
	newShard2Endpoint := make([]string, 0)
	newShard2EndpointId := make([]string, 0)
	/*
	* sharedIdMap
	* Reduce inuse objects count of pb unmarshal.
	* map[string]string just for less searching.
	 */
	sharedIdMap := make(map[string]string)

	for shardId, shardToEps := range shardToEpsList {
		if shardToEps.Eps == nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the shard %v miss endpoint list field. coordinator: %v",
				shardManager.InstanceName, shardId, coordinator)
			return error2.NewTBaseClientInternalError(fmt.Sprintf("the shard %v miss endpoint list field", shardId))
		}

		eps := shardToEps.Eps
		if len(eps.EpList) <= 0 {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the shard %v's endpoint list is empty. coordinator: %v",
				shardManager.InstanceName, shardId, coordinator)
			return error2.NewTBaseClientInternalError(fmt.Sprintf("the shard %v's endpoint list is empty", shardId))
		}

		endPoint := eps.GetEpList()[shardManager.EndpointListIndex]
		if endPoint.Id == nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the shard %v's endpoint miss id field. coordinator: %v",
				shardManager.InstanceName, shardId, coordinator)
			return error2.NewTBaseClientInternalError(fmt.Sprintf("the shard %v's endpoint miss id field", shardId))
		}

		endPointId := endPoint.GetId()
		address := newEndPoints[endPointId]
		if len(address) <= 0 {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the shard %v not matched endpoint %v coordinator: %v",
				shardManager.InstanceName, shardId, endPoint.GetId(), coordinator)
			return error2.NewTBaseClientInternalError(fmt.Sprintf("the shard %v not matched endpoint %v", shardId, endPointId))
		}
		newShard2Endpoint = append(newShard2Endpoint, address)
		if id, ok := sharedIdMap[endPointId]; ok {
			newShard2EndpointId = append(newShard2EndpointId, id) //append id, can't append endPointId
		} else {
			sharedIdMap[endPointId] = endPointId
			newShard2EndpointId = append(newShard2EndpointId, endPointId)
		}
	}

	if len(newShard2EndpointId) <= 0 {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the new layout is empty, oldVersion=%v, newVersion=%v, coordinator: %v",
			shardManager.InstanceName, shardManager.Version, newVersion, coordinator)
		return error2.NewTBaseClientInternalError("the new layout is empty")
	}

	if len(newShard2EndpointId) != int(shardNum) {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v the new layout shard num not match, expect=%v, actual=%v, coordinator: %v",
			shardManager.InstanceName, shardNum, len(newShard2EndpointId), coordinator)
		return error2.NewTBaseClientInternalError(fmt.Sprintf("the new layout shard num not match, expect=%v, actual=%v", shardNum, len(newShard2EndpointId)))
	}

	oldEndpoints := map_helper.GetValuesMap(shardManager.Endpoints)
	atomic.SwapPointer(shardManager.endPointsPtr, unsafe.Pointer(&newEndPoints))
	atomic.SwapPointer(shardManager.shard2EndpointIdPtr, unsafe.Pointer(&newShard2EndpointId))
	atomic.SwapPointer(shardManager.shard2EndpointPtr, unsafe.Pointer(&newShard2Endpoint))
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v rebuild cluster layout succeed, oldVersion=%v, newVersion=%v, "+
			"shardNum=%v, endpointNum=%v, coordinator: %v",
			shardManager.InstanceName, shardManager.Version, newVersion, shardNum, len(newEndPoints), coordinator)
	}
	shardManager.Version = newVersion
	shardManager.RefreshClientPoolChannel <- &model.OldAndNewEndpoints{OldEndpoints: oldEndpoints, NewEndpoints: map_helper.GetValuesMap(&newEndPoints)}

	return nil
}

func (shardManager *ShardManager) buildEndpoints(shardToEps []*alibaba_coordinator_proto.ShardToEps) (map[string]string, error) {
	newEndpoints := make(map[string]string)
	for shardId, shardToEps := range shardToEps {
		if shardToEps == nil {
			continue
		}

		eps := shardToEps.Eps
		if len(eps.EpList) <= 0 {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v shard %v doesn't have eps", shardManager.InstanceName, shardId)
			continue
		}

		if shardManager.EndpointListIndex >= 1 && len(eps.EpList) < 2 {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v EndpointListIndex: %v shardId: %v doesn't have slave", shardManager.InstanceName, shardManager.EndpointListIndex, shardId)
			return nil, error2.NewTBaseClientInternalError(fmt.Sprintf("EndpointListIndex: %v shardId: %v doesn't have slave", shardManager.EndpointListIndex, shardId))
		}

		endPoint := eps.EpList[shardManager.EndpointListIndex]
		if endPoint.Id == nil {
			tbase_log.TBaseLogger.Warnf("[SHARD_MANAGER] %v EndpointListIndex: %v shardId: %v endpoint: %v doesn't have a Id", shardManager.InstanceName, shardManager.EndpointListIndex, shardId, endPoint)
			continue
		}

		if endPoint.Ip != nil && endPoint.Port != nil {
			newEndpoints[endPoint.GetId()] = ip_helper.Uint32ip(endPoint.GetIp()) + ":" + convertor_helper.Int32ToString(endPoint.GetPort())
		}
	}

	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v build endpoint success. newEndpoints %v", shardManager.InstanceName, newEndpoints)
	}
	return newEndpoints, nil
}

func (shardManager *ShardManager) refreshInternal(isForceRefresh bool) error {
	if shardManager.IsShutdown {
		return error2.NewTBaseClientInternalError("shard manager already shutdown")
	}

	coordinatorListForceRefreshed := false
	var coordinatorRpcInvokeSucceed bool
	var layoutResponse *alibaba_coordinator_proto.GetLayoutResponse
	var getClusterLayoutErr error
	var rebuildClusterLayoutErr error

	for {
		coordinators, err := shardManager.getCoordinatorList(shardManager.Servers)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v can't get coordinator list. isForceRefresh: %v, error: %v", shardManager.InstanceName, isForceRefresh, err)
			return err
		}
		coordinatorRpcInvokeSucceed = false

		for i := 0; i < len(coordinators); i++ {
			if shardManager.NextCoordinatorIndex >= len(coordinators) {
				shardManager.NextCoordinatorIndex = 0
			}

			index := shardManager.NextCoordinatorIndex
			shardManager.NextCoordinatorIndex = shardManager.NextCoordinatorIndex + 1
			coordinator := coordinators[index]

			layoutResponse, getClusterLayoutErr = shardManager.GetClusterLayout(coordinator)
			if getClusterLayoutErr == nil {
				coordinatorRpcInvokeSucceed = true
				rebuildClusterLayoutErr = shardManager.RebuildLayout(layoutResponse, coordinator)
				if rebuildClusterLayoutErr == nil {
					return nil
				} else if !isForceRefresh {
					return rebuildClusterLayoutErr
				}
			}
		}

		if !coordinatorRpcInvokeSucceed && !coordinatorListForceRefreshed {
			err := shardManager.RefreshCoordinatorList(shardManager.Servers)
			if err != nil {
				tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v refresh coordinator list error, error: %v", shardManager.InstanceName, err)
				return err
			}
			coordinatorListForceRefreshed = true
		} else {
			break
		}
	}

	if getClusterLayoutErr != nil {
		return getClusterLayoutErr
	} else if rebuildClusterLayoutErr != nil {
		return rebuildClusterLayoutErr
	} else {
		return nil
	}
}

func (shardManager *ShardManager) getCoordinatorList(servers []string) ([]string, error) {
	if shardManager.CachedCoordinatorList == nil || (time.Now().UnixNano()-shardManager.coordinatorListUpdateTime) > int64(shardManager.ConnectionInfo.CoordinatorListCacheTime) {
		err := shardManager.RefreshCoordinatorList(servers)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v can not get coordinator list. error: %v", shardManager.InstanceName, err)
			return nil, err
		}
	}

	return shardManager.CachedCoordinatorList, nil
}

func (shardManager *ShardManager) RefreshCoordinatorList(servers []string) error {
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v start to refresh coordinator list, configured coordinator list: %s", shardManager.InstanceName, servers)
	}
	rpcInvokeSucceed := false
	errMap := make(map[string]error)
	for _, server := range servers {
		coordinatorResp, err := shardManager.GetCoordinatorStat(server)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get coordinator stat from server %v error, error: %v", shardManager.InstanceName, server, err)
			errMap[server] = err
			continue
		}
		if coordinatorResp == nil {
			tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get coordinator stat from server %v error, error: %v", shardManager.InstanceName, server, err)
			errMap[server] = error2.NewTBaseClientInternalError("get coordinator stat from server error, response is null")
			continue
		}
		rpcInvokeSucceed = true

		if coordinatorResp.GetCoordinators() != nil && len(coordinatorResp.GetCoordinators()) > 0 {
			coordinatorList := make([]string, len(coordinatorResp.GetCoordinators()))
			for i, address := range coordinatorResp.GetCoordinators() {
				coordinatorList[i] = ip_helper.CombineIpAndPort(address.GetIp(), address.GetPort())
			}

			shardManager.CachedCoordinatorList = convertor_helper.Shuffle(coordinatorList)
			shardManager.coordinatorListUpdateTime = time.Now().UnixNano()
			shardManager.NextCoordinatorIndex = 0

			if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
				tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] %v from coordinator %v take configured coordinator list: %v", shardManager.InstanceName, server, shardManager.CachedCoordinatorList)
			}
			return nil
		}

	}

	if rpcInvokeSucceed && shardManager.CachedCoordinatorList == nil {
		shardManager.CachedCoordinatorList = servers
		tbase_log.TBaseLogger.Errorf("[SHARD-MANAGER] %v misconfigured or old version coordinator, using connection string configured coordinators", shardManager.InstanceName)
		return nil
	} else {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get coordinator stat from configured vip servers error", shardManager.InstanceName)
		return rebuildError(errMap)
	}
}

func (shardManager *ShardManager) GetCoordinatorStat(server string) (*alibaba_coordinator_proto.GetCoordinatorStatResponse, error) {

	requestId, err := uuid_helper.GenerateUUID()
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v generate uuid error. error: %v", shardManager.InstanceName, err)
		return nil, err
	}

	getCoordinatorStatRequest := pb_helper.ConstructGetCoordinatorStatRequestMessage(requestId)
	requestBuf, err := message_helper.Encode(getCoordinatorStatRequest)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v encode pb message to byte buffer error. pb message: %v, error: %v", shardManager.InstanceName, getCoordinatorStatRequest, err)
		return nil, err
	}

	timeout := time.Duration(shardManager.ConnectionInfo.CoordinatorTimeout) * time.Millisecond
	conn, err := rpc_helper.CreateConnectionAndSendRequest(requestId, server, timeout, requestBuf)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v create connection to server %v with timeout %v and send request error. error: %v",
			shardManager.InstanceName, server, timeout, err)
		return nil, err
	}
	defer conn.Close()

	getCoordinatorStatResponse, err := rpc_helper.GetResponse(conn, requestId, alibaba_coordinator_proto.E_GetCoordinatorStatResponse_Response)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get coordinator response error. dest server: %v, error: %v", shardManager.InstanceName, server, err)
		return nil, err
	}

	if result, ok := getCoordinatorStatResponse.(*alibaba_coordinator_proto.GetCoordinatorStatResponse); ok {
		return result, nil
	} else {
		tbase_log.TBaseLogger.Errorf("can't convert response to \"GetCoordinatorStatResponse\"")
		return nil, error2.NewTBaseClientInternalError("can't convert response to \"GetCoordinatorStatResponse\"")
	}
}

func (shardManager *ShardManager) GetShardId(key []byte) (int, error) {
	hashCode := murmur_hash_helper.MurmurHash2(key, murmur_hash_helper.ConvertToInt32(DEFAULT_HASH_SEED))
	if len(*shardManager.Shard2Endpoint) <= 0 {
		return 0, error2.NewTBaseClientInternalError("local layout table is empty")
	}
	return int(math.Abs(float64(hashCode))) % len(*shardManager.Shard2Endpoint), nil
}

func (shardManager *ShardManager) GetEndpointByKey(key []byte) (string, error) {
	if shardManager.IsShutdown {
		return "", error2.NewTBaseClientInternalError("shard manager already shutdown")
	}

	shardId, err := shardManager.GetShardId(key)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v get shardId error, error: %v", shardManager.InstanceName, err)
		return "", err
	}
	endPoint := (*shardManager.Shard2Endpoint)[shardId]
	if len(endPoint) <= 0 {
		tbase_log.TBaseLogger.Errorf("[SHARD_MANAGER] %v shard id %v not exist", shardManager.InstanceName, shardId)
		return "", error2.NewTBaseClientInternalError(fmt.Sprintf("shard id %v not exist", shardId))
	}
	return endPoint, nil
}

func (shardManager *ShardManager) Refresh() {
	if len(shardManager.RefreshLayoutChannel) < 1 {
		shardManager.RefreshLayoutChannel <- &model.RefreshLayoutParam{ForceRefresh: true, LastForceRefreshTime: atomic.LoadInt64(&shardManager.LastForceRefreshTime)}
	}
}

func (shardManager *ShardManager) GetEndPoints() ([]string, error) {
	if shardManager.IsShutdown {
		return nil, error2.NewTBaseClientInternalError("shard manager already shutdown")
	}
	values := make([]string, 0, len(*shardManager.Endpoints))
	for _, v := range *shardManager.Endpoints {
		values = append(values, v)
	}
	return values, nil
}

func (shardManager *ShardManager) Close() {
	if !shardManager.IsShutdown {
		shardManager.IsShutdown = true
		// RefreshLayoutChannel loop will visit RefreshClientPoolChannel
		// so we need to close RefreshClientPoolChannel until RefreshLayoutChannel is closed
		// and RefreshLayoutChannel goroutine loop is finished
		close(shardManager.RefreshLayoutChannel)
		shardManager.wg.Wait()
		close(shardManager.RefreshClientPoolChannel)
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[SHARD_MANAGER] shard manager %v is closed. ", shardManager.InstanceName)
		}
	}
}

func rebuildError(errMap map[string]error) error {
	errMsg := "get coordinator stat from configured vip servers error. "
	for server, err := range errMap {
		errMsg = errMsg + fmt.Sprintf("server: %v, error: %v; ", server, err)
	}
	return error2.NewTBaseClientInternalError(errMsg)
}
