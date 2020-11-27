package sofabolt

import "sync/atomic"

type ClientMetrics struct {
	nread           int64
	nwrite          int64
	commands        int64
	pendingCommands int64
	references      int64
	used            int64
	lasted          int64
	created         int64
}

func (cm *ClientMetrics) GetBytesRead() int64       { return atomic.LoadInt64(&cm.nread) }
func (cm *ClientMetrics) GetBytesWrite() int64      { return atomic.LoadInt64(&cm.nwrite) }
func (cm *ClientMetrics) GetCommands() int64        { return atomic.LoadInt64(&cm.commands) }
func (cm *ClientMetrics) GetPendingCommands() int64 { return atomic.LoadInt64(&cm.pendingCommands) }
func (cm *ClientMetrics) ResetPendingCommands()     { atomic.StoreInt64(&cm.pendingCommands, 0) }
func (cm *ClientMetrics) GetReferences() int64      { return atomic.LoadInt64(&cm.references) }
func (cm *ClientMetrics) GetUsed() int64            { return atomic.LoadInt64(&cm.used) }
func (cm *ClientMetrics) GetLasted() int64          { return atomic.LoadInt64(&cm.lasted) }
func (cm *ClientMetrics) GetCreated() int64         { return atomic.LoadInt64(&cm.created) }
