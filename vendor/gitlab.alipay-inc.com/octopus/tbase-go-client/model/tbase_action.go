package model

import (
	"time"

	"gitlab.alipay-inc.com/octopus/radix"
)

var NoKeyCmds = map[string]bool{
	"SENTINEL":     true,
	"CLUSTER":      true,
	"READONLY":     true,
	"READWRITE":    true,
	"ASKING":       true,
	"AUTH":         true,
	"ECHO":         true,
	"PING":         true,
	"QUIT":         true,
	"SELECT":       true,
	"SWAPDB":       true,
	"KEYS":         true,
	"MIGRATE":      true,
	"OBJECT":       true,
	"RANDOMKEY":    true,
	"WAIT":         true,
	"SCAN":         true,
	"EVAL":         true,
	"EVALSHA":      true,
	"SCRIPT":       true,
	"BGREWRITEAOF": true,
	"BGSAVE":       true,
	"CLIENT":       true,
	"COMMAND":      true,
	"CONFIG":       true,
	"DBSIZE":       true,
	"DEBUG":        true,
	"FLUSHALL":     true,
	"FLUSHDB":      true,
	"INFO":         true,
	"LASTSAVE":     true,
	"MONITOR":      true,
	"ROLE":         true,
	"SAVE":         true,
	"SHUTDOWN":     true,
	"SLAVEOF":      true,
	"SLOWLOG":      true,
	"SYNC":         true,
	"TIME":         true,
	"DISCARD":      true,
	"EXEC":         true,
	"MULTI":        true,
	"UNWATCH":      true,
	"WATCH":        true,
}

var SingKeyCmds = map[string]bool{
	"SET":              true,
	"GET":              true,
	"SETTSEX":          true,
	"GETTSEX":          true,
	"SETEX":            true,
	"TTL":              true,
	"DEL":              true,
	"SETNX":            true,
	"PSETEX":           true,
	"GETSET":           true,
	"STRLEN":           true,
	"APPEND":           true,
	"SETRANGE":         true,
	"GETRANGE":         true,
	"INCR":             true,
	"INCRBY":           true,
	"INCRBYFLOAT":      true,
	"DECR":             true,
	"DECRBY":           true,
	"HSET":             true,
	"HSETNX":           true,
	"HGET":             true,
	"HEXISTS":          true,
	"HDEL":             true,
	"HSTRLEN":          true,
	"HINCRBY":          true,
	"HINCRBYFLOAT":     true,
	"HMSET":            true,
	"HMGET":            true,
	"HKEYS":            true,
	"HVALS":            true,
	"HGETALL":          true,
	"HSCAN":            true,
	"LPUSH":            true,
	"LPUSHX":           true,
	"RPUSH":            true,
	"RPUSHX":           true,
	"LPOP":             true,
	"RPOP":             true,
	"RPOPLPUSH":        true,
	"LREM":             true,
	"LLEN":             true,
	"LINDEX":           true,
	"LINSERT":          true,
	"LSET":             true,
	"LRANGE":           true,
	"LTRIM":            true,
	"BLPOP":            true,
	"BRPOP":            true,
	"BRPOPLPUSH":       true,
	"SADD":             true,
	"SISMEMBER":        true,
	"SPOP":             true,
	"SRANDMEMBER":      true,
	"SREM":             true,
	"SMOVE":            true,
	"SCARD":            true,
	"SMEMBERS":         true,
	"SSCAN":            true,
	"SINTER":           true,
	"SINTERSTORE":      true,
	"SUNION":           true,
	"SUNIONSTORE":      true,
	"SDIFF":            true,
	"SDIFFSTORE":       true,
	"ZADD":             true,
	"ZSCORE":           true,
	"ZINCRBY":          true,
	"ZCARD":            true,
	"ZCOUNT":           true,
	"ZRANGE":           true,
	"ZREVRANGE":        true,
	"ZRANGEBYSCORE":    true,
	"ZREVRANGEBYSCORE": true,
	"ZRANK":            true,
	"ZREVRANK":         true,
	"ZREM":             true,
	"ZREMRANGEBYRANK":  true,
	"ZREMRANGEBYSCORE": true,
	"ZRANGEBYLEX":      true,
	"ZLEXCOUNT":        true,
	"ZREMRANGEBYLEX":   true,
	"ZSCAN":            true,
	"EXPIRE":           true,
	"EXPIREAT":         true,
	"PTTL":             true,
	"SRATELIMITER":     true,
	"RATELIMITER":      true,
}

var MultiKeyCmds = map[string]bool{
	"MSET":   true,
	"MSETNX": true,
	"MGET":   true,
}

type TBaseAction struct {
	Attempts       int
	SubmitTime     int64
	Timeout        int64
	EndpointString string

	Rcv  interface{}
	Cmd  string
	Args []string

	Flat     bool
	FlatKey  [1]string // use array to avoid allocation in Keys
	FlatArgs []interface{}
}

func (ta *TBaseAction) DecrementAttempts() {
	ta.Attempts--
}

func (ta *TBaseAction) Elapsed() int64 {
	return time.Now().UnixNano() - ta.SubmitTime
}

func (ta *TBaseAction) GetInnerAction() radix.CmdAction {
	if ta.Flat {
		return radix.FlatCmd(ta.Rcv, ta.Cmd, ta.FlatKey[0], ta.FlatArgs)
	} else {
		return radix.Cmd(ta.Rcv, ta.Cmd, ta.Args...)
	}
}

func (ta *TBaseAction) IsExpired() bool {
	return ta.RemainTime() <= 0
}

func (ta *TBaseAction) RemainTime() int64 {
	return ta.timeoutTime() - time.Now().UnixNano()
}

func (ta *TBaseAction) timeoutTime() int64 {
	return ta.SubmitTime + ta.Timeout
}
