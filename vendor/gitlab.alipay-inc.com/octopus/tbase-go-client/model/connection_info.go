package model

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
)

const (
	DEFAULT_CONNECTION_TIMEOUT              int   = 1000
	DEFAULT_REDIS_TIMEOUT                   int   = 200
	DEFAULT_COORDINATOR_TIMEOUT             int   = 1000 * 3
	DEFAULT_PORT                            int   = 6001
	DEFAULT_MAX_CONNECTIONS                 int   = 1
	DEFAULT_MAX_RETRIES                     int   = 2
	DEFAULT_LAYOUT_REFRESH_INTERVAL         int   = 1000 * 30
	DEFAULT_MINIMAL_LAYOUT_REFRESH_TIMESPAN int64 = 1000 * 5
	DEFAULT_FAILURES_TO_REFRESH             int32 = 3
	DEFAULT_FAILURE_DETECT_INTERVAL         int64 = 1000 * 60
	DEFAULT_MAX_QUEUE_SIZE                  int   = 4096

	// todo 该配置值的处理
	DEFAULT_MAX_KEY_SIZE           int   = 1024 * 1
	DEFAULT_MAX_VALUE_SIZE         int   = 1024 * 1024
	DEFAULT_MAX_KEY_COUNT          int   = 1024
	DEFAULT_WARMUP                 bool  = false
	DEFAULT_HANDSHAKE              bool  = false
	DEFAULT_INIT_TIMEOUT           int   = 1000 * 1
	DEFAULT_COORDINATOR_CACHE_TIME int64 = 1000 * 60 * 5
	DEFAULT_HOTKEY_EXPIRE_TIME     int   = 2000
	DEFAULT_SLAVE_MODE             bool  = false
)

type ConnectionInfo struct {
	Servers                      []string
	Cluster                      string
	Tenant                       string
	AppName                      string
	ConnectionTimeout            int //ms
	RedisTimeout                 int //ms
	CoordinatorTimeout           int //ms
	MaxConnections               int
	MaxRetries                   int
	LayoutRefreshInterval        int   //ms
	MinimalLayoutRefreshTimespan int64 //ns
	FailuresToRefresh            int32
	FailureDetectInterval        int64 //ns
	MaxQueueSize                 int
	MaxWorkers                   int
	MaxKeySize                   int
	MaxValueSize                 int
	MaxKeyCount                  int
	Warmup                       bool
	Handshake                    bool
	InitTimeout                  int
	CoordinatorListCacheTime     int64 //ns
	HotkeyExpireTime             int
	SlaveMode                    bool
	Password                     string
}

func Parse(connectionStr string) (*ConnectionInfo, error) {

	conInfo := NewConnectionInfo(nil, "", "", "", DEFAULT_CONNECTION_TIMEOUT, DEFAULT_REDIS_TIMEOUT, DEFAULT_COORDINATOR_TIMEOUT,
		DEFAULT_MAX_CONNECTIONS, DEFAULT_MAX_RETRIES, DEFAULT_LAYOUT_REFRESH_INTERVAL, DEFAULT_MINIMAL_LAYOUT_REFRESH_TIMESPAN, DEFAULT_FAILURES_TO_REFRESH, DEFAULT_FAILURE_DETECT_INTERVAL,
		DEFAULT_MAX_QUEUE_SIZE, runtime.GOMAXPROCS(0), DEFAULT_MAX_KEY_SIZE, DEFAULT_MAX_VALUE_SIZE, DEFAULT_MAX_KEY_COUNT, DEFAULT_WARMUP, DEFAULT_HANDSHAKE, DEFAULT_INIT_TIMEOUT,
		DEFAULT_COORDINATOR_CACHE_TIME, DEFAULT_HOTKEY_EXPIRE_TIME, DEFAULT_SLAVE_MODE, "")

	connectionStr = strings.TrimSuffix(connectionStr, ";")
	items := strings.Split(connectionStr, ";")

	for _, item := range items {
		keyValues := strings.Split(item, "=")
		if len(keyValues) != 2 {
			return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("expect \"key=value\" format, actually is %v", item))
		}
		key := keyValues[0]
		value := keyValues[1]
		if "servers" == key {
			conInfo.Servers = strings.Split(value, ",")
		} else if "cluster" == key {
			conInfo.Cluster = value
		} else if "tenant" == key {
			conInfo.Tenant = value
		} else if "appname" == key {
			conInfo.AppName = value
		} else if "connection timeout" == key {
			connectionTimeout, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"connection timeout\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			if connectionTimeout <= 0 {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'connection timeout' <= 0"))
			}
			conInfo.ConnectionTimeout = connectionTimeout
		} else if "redis timeout" == key {
			redisTimeout, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"redis timeout\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.RedisTimeout = redisTimeout
		} else if "coordinator timeout" == key {
			coordinatorTimeout, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"coordinator timeout\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.CoordinatorTimeout = coordinatorTimeout
		} else if "max connections" == key {
			maxConnections, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max connections\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxConnections = maxConnections
		} else if "max retries" == key {
			maxRetries, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max retries\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxRetries = maxRetries
		} else if "layout refresh interval" == key {
			layoutRefreshInterval, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"layout refresh interval\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.LayoutRefreshInterval = layoutRefreshInterval
		} else if "minimal layout refresh timespan" == key {
			minimalLayoutRefreshTimespan, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"minimal layout refresh timespan\" value to \"int64\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MinimalLayoutRefreshTimespan = minimalLayoutRefreshTimespan * 1e6
		} else if "failures to refresh" == key {
			failuresToRefresh, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"failures to refresh\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.FailuresToRefresh = int32(failuresToRefresh)
		} else if "failure detect interval" == key {
			failureDetectInterval, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"failure detect interval\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.FailureDetectInterval = failureDetectInterval * 1e6
		} else if "max queue size" == key {
			maxQueueSize, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max queue size\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxQueueSize = maxQueueSize
		} else if "max workers" == key {
			maxWorkers, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max workers\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxWorkers = maxWorkers
		} else if "max key size" == key {
			maxKeySize, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max key size\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxKeySize = maxKeySize
		} else if "max value size" == key {
			maxValueSize, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max value size\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxValueSize = maxValueSize
		} else if "max key count" == key {
			maxKeyCount, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"max key count\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.MaxKeyCount = maxKeyCount
		} else if "warmup" == key {
			warmup, err := strconv.ParseBool(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"warmup\" value to \"bool\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.Warmup = warmup
		} else if "handshake" == key {
			handshake, err := strconv.ParseBool(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"handshake\" value to \"bool\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.Handshake = handshake
		} else if "init timeout" == key {
			initTimeout, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"init timeout\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.InitTimeout = initTimeout
		} else if "coordinator list cache time" == key {
			coordinatorListCacheTime, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"coordinator list cache time\" value to \"int64\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.CoordinatorListCacheTime = coordinatorListCacheTime * 1e6
		} else if "hotkey expire time" == key {
			hotkeyExpireTime, err := strconv.Atoi(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"hotkey expire time\" value to \"int\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.HotkeyExpireTime = hotkeyExpireTime
		} else if "slavemode" == key {
			slaveMode, err := strconv.ParseBool(value)
			if err != nil {
				return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("convert \"slavemode\" value to \"bool\" error. raw value: %v, error: %v", value, err))
			}
			conInfo.SlaveMode = slaveMode
		} else if "password" == key {
			conInfo.Password = value
		} else {
			return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, unknown item %v=%v", key, value))
		}
	}

	if len(conInfo.Servers) <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'servers' is required"))
	}
	if len(conInfo.Cluster) <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'cluster' is required"))
	}
	if len(conInfo.Tenant) <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'tenant' is required"))
	}
	if conInfo.RedisTimeout <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'redis timeout' <= 0"))
	}
	if conInfo.CoordinatorTimeout <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'coordinator timeout' <= 0"))
	}
	if conInfo.MaxConnections <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max connections' <= 0"))
	}
	if conInfo.MaxRetries <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max retries' <= 0"))
	}
	if conInfo.LayoutRefreshInterval <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'layout refresh interval' <= 0"))
	}
	if conInfo.MinimalLayoutRefreshTimespan <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'minimal layout refresh timespan' <= 0"))
	}
	if conInfo.FailuresToRefresh <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'failures to refresh' <= 0"))
	}
	if conInfo.FailureDetectInterval <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'failure detect interval' <= 0"))
	}
	if conInfo.MaxQueueSize <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max queue size' <= 0"))
	}
	if conInfo.MaxWorkers <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max workers' <= 0"))
	}
	if conInfo.MaxValueSize <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max value size' <= 0"))
	}
	if conInfo.MaxKeyCount <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max key count' <= 0"))
	}
	if conInfo.MaxValueSize < conInfo.MaxKeySize {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'max key size' must large than 'max value size'"))
	}
	if conInfo.InitTimeout <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'init timeout' <= 0"))
	}
	if conInfo.CoordinatorListCacheTime <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'coordinator list cache time' <= 0"))
	}
	if conInfo.HotkeyExpireTime <= 0 {
		return nil, error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("parse connection string error, param 'hotkey expire time' <= 0"))
	}

	return conInfo, nil

}

func NewConnectionInfo(Servers []string, cluster string, tenant string, appName string, connectionTimeout int, redisTimeout int, coordinatorTimeout int,
	maxConnections int, maxRetries int, layoutRefreshInterval int, minimalLayoutRefreshTimespan int64, failuresToRefresh int32, failureDetectInterval int64,
	maxQueueSize int, maxWorkers int, maxKeySize int, maxValueSize int, maxKeyCount int, warmup bool, handshake bool, initTimeout int, coordinatorListCacheTime int64,
	hotkeyExpireTime int, slaveMode bool, password string) *ConnectionInfo {
	return &ConnectionInfo{
		Servers:                      Servers,
		Cluster:                      cluster,
		Tenant:                       tenant,
		AppName:                      appName,
		ConnectionTimeout:            connectionTimeout,
		RedisTimeout:                 redisTimeout,
		CoordinatorTimeout:           coordinatorTimeout,
		MaxConnections:               maxConnections,
		MaxRetries:                   maxRetries,
		LayoutRefreshInterval:        layoutRefreshInterval,
		MinimalLayoutRefreshTimespan: minimalLayoutRefreshTimespan * 1e6,
		FailuresToRefresh:            failuresToRefresh,
		FailureDetectInterval:        failureDetectInterval * 1e6,
		MaxQueueSize:                 maxQueueSize,
		MaxWorkers:                   maxWorkers,
		MaxKeySize:                   maxKeySize,
		MaxValueSize:                 maxValueSize,
		MaxKeyCount:                  maxKeyCount,
		Warmup:                       warmup,
		Handshake:                    handshake,
		InitTimeout:                  initTimeout,
		CoordinatorListCacheTime:     coordinatorListCacheTime * 1e6,
		HotkeyExpireTime:             hotkeyExpireTime,
		SlaveMode:                    slaveMode,
		Password:                     password,
	}
}

func (c *ConnectionInfo) ToString() string {
	connectionString := "servers=" + c.Servers[0]
	for i := 1; i < len(c.Servers); i++ {
		connectionString += fmt.Sprintf(",%v", c.Servers[i])
	}
	connectionString += fmt.Sprintf(";cluster=%v;tenant=%v;appname=%v;connection timeout=%v;"+
		"redis timeout=%v;coordinator timeout=%v;max connections=%v;max retries=%v;layout refresh interval=%v;"+
		"minimal layout refresh timespan=%v;failures to refresh=%v;failure detect interval=%v;max queue size=%v"+
		";max workers=%v;max key size=%v;max value size=%v;max key count=%v;warmup=%v;handshake=%v;init timeout=%v"+
		";coordinator list cache time=%v;slavemode=%v;hotkey expire time=%v;password=%v", c.Cluster, c.Tenant, c.AppName, c.ConnectionTimeout,
		c.RedisTimeout, c.ConnectionTimeout, c.MaxConnections, c.MaxRetries, c.LayoutRefreshInterval, c.MinimalLayoutRefreshTimespan, c.FailuresToRefresh,
		c.FailureDetectInterval, c.MaxQueueSize, c.MaxWorkers, c.MaxKeySize, c.MaxValueSize, c.MaxKeyCount, c.Warmup, c.Handshake, c.InitTimeout,
		c.CoordinatorListCacheTime, c.SlaveMode, c.HotkeyExpireTime, c.Password)
	return connectionString
}
