package tbasego

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	time_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/time.helper"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/model"
	error_code "gitlab.alipay-inc.com/octopus/tbase-go-client/model/ecode"
	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
	"mosn.io/pkg/log"
)

var commandExecutorCounter int32 = 0

type CommandExecutor struct {
	clientHandler   *ClientHandler
	connectionInfo  *model.ConnectionInfo
	instanceName    string
	failureCount    int32
	lastFailureTime int64
	closed          bool
	socketTimeout   time.Duration
}

func newCommandExecutor(connectioInfo *model.ConnectionInfo, handler *ClientHandler) *CommandExecutor {
	instanceName := fmt.Sprintf("CommandExecutor-(%v.%v#%v)", connectioInfo.Cluster, connectioInfo.Tenant, atomic.AddInt32(&commandExecutorCounter, 1))
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[COMMAND_EXECUTOR] init command executor instance '%v' with connection string: '%v'", instanceName, connectioInfo.ToString())
	}
	return &CommandExecutor{
		clientHandler:   handler,
		connectionInfo:  connectioInfo,
		failureCount:    0,
		instanceName:    instanceName,
		lastFailureTime: math.MaxInt64,
		closed:          false,
		socketTimeout:   10 * time.Second,
	}
}

func (commandExecutor *CommandExecutor) do(ta *model.TBaseAction) (string, error) {
	key, err := commandExecutor.checkAction(ta)
	if err != nil {
		return "", err
	}

	var tbaseError error = nil
	var movedHints = make(map[int]string)
	for ta.Attempts > 0 {
		clientResource, err := commandExecutor.clientHandler.takeClient(ta, []byte(key), time.Duration(ta.RemainTime()),
			commandExecutor.connectionInfo.MaxQueueSize, movedHints,
			commandExecutor.socketTimeout)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[COMMAND_EXECUTOR] take client error, error: %v", err)
			tbaseError = commandExecutor.processError(err, ta, movedHints, key)
			if isNeedReturn := needReturnBack(tbaseError, ta.Attempts); isNeedReturn {
				return ta.EndpointString, tbaseError
			}
			continue
		}
		err = clientResource.asyncClient.Do(ta)
		if err != nil {
			tbaseError = commandExecutor.processError(err, ta, movedHints, key)
			if isNeedReturn := needReturnBack(tbaseError, ta.Attempts); isNeedReturn {
				return ta.EndpointString, tbaseError
			}
			continue
		}
		return ta.EndpointString, nil
	}

	return ta.EndpointString, tbaseError
}

func (commandExecutor *CommandExecutor) checkAction(action *model.TBaseAction) (string, error) {
	if !model.SingKeyCmds[strings.ToUpper(action.Cmd)] {
		return "", error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("unsupported command '%v'", action.Cmd))
	}

	if action.Flat && (len(action.FlatKey[0]) > commandExecutor.connectionInfo.MaxKeySize || len(action.FlatKey[0]) <= 0) {
		return "", error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("expect 0 < key size <= %v, but actual key size is %v",
			commandExecutor.connectionInfo.MaxKeySize, len(action.FlatKey[0])))
	}

	if !action.Flat && (len(action.Args[0]) > commandExecutor.connectionInfo.MaxKeySize || len(action.Args[0]) <= 0) {
		return "", error2.NewTBaseClientIllegalArgumentsError(fmt.Sprintf("expect 0 < key size <= %v, but actual key size is %v",
			commandExecutor.connectionInfo.MaxKeySize, len(action.Args[0])))
	}

	// todo value 大小的check
	if action.Flat {
		return action.FlatKey[0], nil
	} else {
		return action.Args[0], nil
	}
}

func (commandExecutor *CommandExecutor) processError(err error, ta *model.TBaseAction,
	movedHints map[int]string, key string) error {

	var tbaseError error = nil
	if _, ok := err.(*error2.TBaseClientInternalError); ok {
		tbaseError = err
		commandExecutor.processFailure(false, ta, tbaseError, key)
	} else if _, ok := err.(*error2.TBaseClientTimeoutError); ok {
		tbaseError = err
		commandExecutor.processFailure(false, ta, tbaseError, key)
	} else if _, ok := err.(*error2.TBaseClientConnectionError); ok {
		tbaseError = err
		commandExecutor.processFailure(false, ta, tbaseError, key)
	} else if strings.Contains(err.Error(), error_code.ERROR_MOVED) {
		tbaseError = error2.NewTBaseClientMovedDataError(err.Error())
		commandExecutor.processFailure(true, ta, tbaseError, key)
		// this type cast will not panic
		err := extractMovedHints(tbaseError.(*error2.TBaseClientMovedDataError), movedHints)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[COMMAND_EXECUTOR] extract moved hints error, error: %v", err)
			return err
		}
	} else if strings.Contains(err.Error(), error_code.ERROR_READONLY) {
		tbaseError = error2.NewTBaseClientReadOnlyError(err.Error())
		commandExecutor.processFailure(true, ta, tbaseError, key)
	} else if strings.Contains(err.Error(), error_code.ERROR_HOTKEY) {
		tbaseError = error2.NewTBaseClientHotKeyDataError(err.Error())
		commandExecutor.processFailure(false, ta, tbaseError, key)
	} else if strings.Contains(err.Error(), error_code.ERROR_LOADING) {
		tbaseError = error2.NewTBaseClientLoadingError(err.Error())
		commandExecutor.processFailure(true, ta, tbaseError, key)
	} else if strings.Contains(err.Error(), error_code.ERROR_HANDSHAKE_ERR) {
		tbaseError = error2.NewTBaseClientHandShakeError(err.Error())
		commandExecutor.processFailure(true, ta, tbaseError, key)
	} else if strings.Contains(err.Error(), error_code.ERROR_OVERLOAD) {
		tbaseError = error2.NewTBaseClientOverloadError(err.Error())
	} else if strings.Contains(err.Error(), error_code.ERROR_LAG) {
		tbaseError = error2.NewTBaseClientLagError(err.Error())
		commandExecutor.processFailure(false, ta, tbaseError, key)
	} else if _, ok := err.(*error2.TBaseClientDataError); ok {
		tbaseError = err
	} else if error2.IsConnectionErr(err) {
		if netError, _ := err.(net.Error); netError != nil {
			if netError.Timeout() {
				tbaseError = error2.NewTBaseClientTimeoutError(err.Error())
				commandExecutor.processFailure(false, ta, tbaseError, key)
			} else {
				tbaseError = error2.NewTBaseClientConnectionError(err.Error())
				commandExecutor.processFailure(false, ta, tbaseError, key)
			}
		} else {
			tbaseError = error2.NewTBaseClientConnectionError(err.Error())
			commandExecutor.processFailure(false, ta, tbaseError, key)
		}
	} else if strings.HasPrefix(err.Error(), error_code.ERROR_DATA_ERROR) {
		tbaseError = error2.NewTBaseClientDataError(err.Error())
	} else {
		tbaseError = error2.NewTBaseClientInternalError(err.Error())
		commandExecutor.processFailure(false, ta, tbaseError, key)
	}

	return tbaseError
}

func (commandExecutor *CommandExecutor) processFailure(forceRefresh bool, ta *model.TBaseAction,
	err error, key string) {
	if forceRefresh {
		commandExecutor.clientHandler.refresh()
		atomic.StoreInt32(&commandExecutor.failureCount, 0)
	} else {
		commandExecutor.countFailure()
	}

	ta.DecrementAttempts()
	tbase_log.TBaseLogger.Errorf("[COMMAND_EXECUTOR] do action error. command: %v, key: %v, dests: %v, "+
		"submit time: %v, nanos: %v, elapsed: %v(ms), attempts: %v, error: %v",
		ta.Cmd, key, ta.EndpointString, time.Unix(0, ta.SubmitTime),
		ta.SubmitTime, time_helper.NsToMs(ta.Elapsed()), ta.Attempts, err)
}

func (commandExecutor *CommandExecutor) countFailure() {
	atomic.AddInt32(&commandExecutor.failureCount, 1)
	if commandExecutor.failureCount > commandExecutor.connectionInfo.FailuresToRefresh {
		commandExecutor.clientHandler.refresh()
		atomic.StoreInt32(&commandExecutor.failureCount, 0)
	} else if time.Now().UnixNano()-commandExecutor.lastFailureTime > commandExecutor.connectionInfo.FailureDetectInterval {
		atomic.StoreInt32(&commandExecutor.failureCount, 0)
	}
	atomic.StoreInt64(&commandExecutor.lastFailureTime, time.Now().UnixNano())
}

func (commandExecutor *CommandExecutor) Close() {
	if !commandExecutor.closed {
		commandExecutor.closed = true
		if commandExecutor.clientHandler != nil {
			commandExecutor.clientHandler.close()
			if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
				tbase_log.TBaseLogger.Infof("[COMMAND_EXECUTOR] command executor %v is closed", commandExecutor.instanceName)
			}
		}
	}
}

func needReturnBack(err error, attempt int) bool {
	if attempt <= 0 {
		return true
	}
	if _, ok := err.(*error2.TBaseClientTimeoutError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientReadOnlyError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientLoadingError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientHandShakeError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientLagError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientDataError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientInternalError); ok {
		return true
	}

	if _, ok := err.(*error2.TBaseClientOverLoadError); ok {
		return true
	}

	return false
}

func extractMovedHints(movedError *error2.TBaseClientMovedDataError, movedHints map[int]string) error {
	movedHintIndex := strings.LastIndex(movedError.Message, error_code.ERROR_MOVED)
	if movedHintIndex == -1 {
		tbase_log.TBaseLogger.Errorf("[COMMAND_EXECUTOR] extract shardId error. can't find MOVED hint")
		return movedError
	}

	movedHintStr := movedError.Message[movedHintIndex:len(movedError.Message)]
	messageItems := strings.Split(movedHintStr, " ")
	shardId, err := strconv.Atoi(messageItems[1])
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[COMMAND_EXECUTOR] extract shardId error. moved error message: %v, convert error: %v", messageItems, err)
		return err
	}

	movedHints[shardId] = messageItems[2]
	return nil
}
