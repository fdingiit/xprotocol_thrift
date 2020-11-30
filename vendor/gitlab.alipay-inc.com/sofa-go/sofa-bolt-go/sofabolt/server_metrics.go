package sofabolt

import "sync/atomic"

type ServerMetrics struct {
	numwrite           int64
	numread            int64
	commands           int64
	pendingcommands    int64
	connections        int64
	pendingconnections int64
}

func (sm *ServerMetrics) GetBytesRead() int64 {
	return atomic.LoadInt64(&sm.numread)
}

func (sm *ServerMetrics) GetBytesWrite() int64 {
	return atomic.LoadInt64(&sm.numwrite)
}

func (sm *ServerMetrics) GetCommands() int64 {
	return atomic.LoadInt64(&sm.commands)
}

func (sm *ServerMetrics) GetPendingCommands() int64 {
	return atomic.LoadInt64(&sm.pendingcommands)
}

func (sm *ServerMetrics) GetConnections() int64 {
	return atomic.LoadInt64(&sm.connections)
}

func (sm *ServerMetrics) GetPendingConnections() int64 {
	return atomic.LoadInt64(&sm.pendingconnections)
}

func (sm *ServerMetrics) addConnections(n int64) {
	atomic.AddInt64(&sm.connections, n)
}

func (sm *ServerMetrics) addPendingConnections(n int64) {
	atomic.AddInt64(&sm.pendingconnections, n)
}

func (sm *ServerMetrics) addBytesRead(n int64) {
	atomic.AddInt64(&sm.numread, n)
}

func (sm *ServerMetrics) addBytesWrite(n int64) {
	atomic.AddInt64(&sm.numwrite, n)
}

func (sm *ServerMetrics) addCommands(n int64) {
	atomic.AddInt64(&sm.commands, n)
}

func (sm *ServerMetrics) addPendingCommands(n int64) {
	atomic.AddInt64(&sm.pendingcommands, n)
}
