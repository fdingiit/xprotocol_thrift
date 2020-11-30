package async

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/model"

	"gitlab.alipay-inc.com/octopus/radix"
)

type asyncCmd struct {
	*model.TBaseAction
	resCh chan error
}

type asyncCmdWithConn struct {
	*asyncCmd
	conn *radix.IoErrConn
}

type AsyncClient struct {
	endpoint  string
	asyncPool *AsyncPool
	limit     int

	reqCh         chan *asyncCmd         //accept cmd for request
	respCh        chan *asyncCmdWithConn //accept cmd for response
	closeReqChan  chan bool              //accept close signal for request channel
	closeRespChan chan bool              //accept close signal for response channel
	reqAndRespWg  *sync.WaitGroup

	closed     bool
	ConnUsable uint32 // 1: usable; 0: not usable

	lock sync.Mutex
}

func NewAsyncClient(endpoint string, connectionTimeout time.Duration, maxQueueSize int, socketTimeout time.Duration) (*AsyncClient, error) {
	asyncPool, err := NewAsyncPool(endpoint, connectionTimeout, socketTimeout)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[ASYNC_CMD_EXECUTOR] new async pool to endpoint %v error, error: %v", endpoint, err)
		return nil, err
	}
	a := &AsyncClient{
		endpoint:      endpoint,
		asyncPool:     asyncPool,
		limit:         maxQueueSize,
		reqCh:         make(chan *asyncCmd, maxQueueSize*2),
		respCh:        make(chan *asyncCmdWithConn, maxQueueSize*2),
		closeReqChan:  make(chan bool),
		closeRespChan: make(chan bool),
		reqAndRespWg:  new(sync.WaitGroup),
		ConnUsable:    1,
	}

	a.reqAndRespWg.Add(2)
	go func() {
		defer a.reqAndRespWg.Done()
		a.startRequestLoop()
	}()

	go func() {
		defer a.reqAndRespWg.Done()
		a.startResponseLoop()
	}()

	return a, nil
}

func (ac *AsyncClient) Do(ta *model.TBaseAction) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			tbase_log.TBaseLogger.Errorf("unexpected panic, panic is %v\ntrace%s\n", panicError, string(debug.Stack()))
			err = error2.NewTBaseClientInternalError(
				fmt.Sprintf("unexpected panic, panic is %v", panicError))
		}
	}()

	req := getAsyncCmd(ta)

	if len(ac.reqCh)+len(ac.respCh) > ac.limit {
		return error2.QUEUE_FULL_ERROR
	}

	if ac.closed {
		return error2.NewTBaseClientInternalError(
			fmt.Sprintf("instance %v is closed", ac.endpoint))
	}

	if atomic.LoadUint32(&ac.ConnUsable) == 0 {
		return error2.CONNECTION_NOT_USABLE_ERROR
	}

	ac.reqCh <- req

	select {
	case result := <-req.resCh:
		return result
	case <-time.After(time.Duration(ta.RemainTime())):
		return error2.COMMAND_TIMEOUT_ERROR
	}
}

// Close closes the AsyncClient and make sure that
// 1. all background goroutines are stopped before returning;
// 2. close connection pool
func (ac *AsyncClient) Close() {
	if !ac.closed {
		ac.closed = true
		// close the connection first to let 'Close()' finish quickly
		ac.asyncPool.Close()
		ac.closeReqChan <- true
		ac.reqAndRespWg.Wait()
		tbase_log.TBaseLogger.Infof("[ASYNC_CLIENT] async client is closed")
	}
}

func (ac *AsyncClient) startRequestLoop() {
	for {
		select {
		case <-ac.closeReqChan:
			close(ac.reqCh)
			ac.closeRespChan <- true
			tbase_log.TBaseLogger.Infof("[ASYNC_CLIENT] request channel closed")
			return

		case req := <-ac.reqCh:
			ac.processRequest(req)
		}
	}
}

func (ac *AsyncClient) startResponseLoop() {
	for {
		select {
		case <-ac.closeRespChan:
			close(ac.respCh)
			tbase_log.TBaseLogger.Infof("[ASYNC_CLIENT] response channel closed")
			return

		case asyncCmdWithConn := <-ac.respCh:
			ac.processResponse(asyncCmdWithConn)
		}
	}
}

func (ac *AsyncClient) processResponse(cmd *asyncCmdWithConn) {
	if cmd == nil {
		tbase_log.TBaseLogger.Errorf("[ASYNC_CLIENT] unexpected nil cmd, won't do anything")
		return
	}

	if ac.closed {
		cmd.asyncCmd.resCh <- error2.NewTBaseClientInternalError(
			fmt.Sprintf("instance %v is closed", ac.endpoint))
		return
	}

	if atomic.LoadUint32(&ac.ConnUsable) == 0 {
		cmd.asyncCmd.resCh <- error2.CONNECTION_NOT_USABLE_ERROR
		return
	}
	err, isNetErr := ac.readResponse(cmd)
	if isNetErr {
		atomic.StoreUint32(&ac.ConnUsable, 0)
	}
	cmd.asyncCmd.resCh <- err
}

func (ac *AsyncClient) closeConn(conn *radix.IoErrConn) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	conn.Close()
}

func (ac *AsyncClient) processRequest(req *asyncCmd) {
	if req == nil {
		tbase_log.TBaseLogger.Errorf("[ASYNC_CLIENT] unexpected nil cmd, won't do anything")
		return
	}

	if atomic.LoadUint32(&ac.ConnUsable) == 0 {
		req.resCh <- error2.CONNECTION_NOT_USABLE_ERROR
		return
	}

	if ac.closed {
		req.resCh <- error2.NewTBaseClientInternalError(
			fmt.Sprintf("instance %v is closed", ac.endpoint))
		return
	}

	err, isNetErr := ac.writeRequest(req)
	if err != nil {
		if isNetErr {
			atomic.StoreUint32(&ac.ConnUsable, 0)
		}
		req.resCh <- err
	}
}

//
// encode and flush cmd and return:
// error: actual error if an error happens, otherwise return nil
// bool: if an actual network error happens
//
func (ac *AsyncClient) writeRequest(cmd *asyncCmd) (error, bool) {
	if cmd.IsExpired() {
		return error2.NewTBaseClientTimeoutError("submit timeout"), false
	} else {
		conn := ac.asyncPool.Get()
		if err := ac.doEncode(conn, cmd.GetInnerAction()); err != nil {
			if error2.IsConnectionErr(err) {
				return err, true
			} else {
				return err, false
			}
		} else {
			ac.respCh <- getAsyncCmdWithConn(cmd, conn)
			return nil, false
		}
	}
}

//
// read and decode response. return:
// error: actual error if an error happens, otherwise return nil
// bool: if an actual network error happens
//
func (ac *AsyncClient) readResponse(cmd *asyncCmdWithConn) (error, bool) {
	if err := ac.doDecode(cmd.conn, cmd.asyncCmd.GetInnerAction()); err != nil {
		if error2.IsConnectionErr(err) {
			return err, true
		} else {
			return err, false
		}
	} else {
		return nil, false
	}
}

// radix encode may panic, so we need to handle panic
// defer can't have return value, but it can change it
func (ac *AsyncClient) doEncode(conn *radix.IoErrConn, action radix.CmdAction) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = error2.NewTBaseClientInternalError(fmt.Sprintf("catch panic, panic is %v", panicError))
		}
	}()

	return conn.Encode(action)
}

// radix encode may panic, so we need to handle panic
// defer can't have return value, but it can change it
func (ac *AsyncClient) doDecode(conn *radix.IoErrConn, action radix.CmdAction) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = error2.NewTBaseClientInternalError(fmt.Sprintf("catch panic, panic is %v", panicError))
		}
	}()

	return conn.Decode(action)
}

func getAsyncCmd(action *model.TBaseAction) *asyncCmd {
	return &asyncCmd{
		action,
		make(chan error, 1),
	}
}

func getAsyncCmdWithConn(cmd *asyncCmd, conn *radix.IoErrConn) *asyncCmdWithConn {
	return &asyncCmdWithConn{cmd, conn}
}
