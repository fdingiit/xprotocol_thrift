package sls

import (
	"fmt"
	"os"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
)

const (
	PRODUCER_SUCCESS        = 0
	PRODUCER_DATA_IS_EMPTY  = 1
	PRODUCER_BUFFER_IS_FULL = 2
)

func (s *LogClient) consumer() {
	go s.writeLoop()
}

func (s *LogClient) writeLoop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "log client consumer error: %v", err)
		}
	}()
	tick := time.NewTicker(time.Duration(s.pushPeriod) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			s.pushLogs()
		case item, ok := <-s.logChan:
			if !ok {
				return
			}
			timestamp := time.Now()
			contents := make([]*sls.LogContent, 0)
			for key, value := range item {
				contents = append(contents, &sls.LogContent{
					Key:   proto.String(key),
					Value: proto.String(value),
				})
			}

			l := &sls.Log{
				Time:     proto.Uint32(uint32(timestamp.Unix())),
				Contents: contents,
			}
			s.appendLog(l)
		}
	}
}

func (s *LogClient) Producer(data map[string]string) int {
	if data == nil || len(data) <= 0 {
		return PRODUCER_DATA_IS_EMPTY
	}

	select {
	case s.logChan <- data:
		return PRODUCER_SUCCESS
	default:
		return PRODUCER_BUFFER_IS_FULL
	}
}

func (s *LogClient) Stop() {
	close(s.logChan)
}

func (s *LogClient) appendLog(l *sls.Log) {
	s.logBuffer = append(s.logBuffer, l)
	if len(s.logBuffer) < int(s.bufferSize) {
		return
	}
	// if log buffer is full, launch logs in buffer
	logGroup := &sls.LogGroup{
		Topic:  proto.String(s.topic),
		Source: proto.String(s.source),
		Logs:   s.logBuffer,
	}

	s.logStore.PutLogs(logGroup)
	s.logBuffer = s.logBuffer[0:0]
}

func (s *LogClient) pushLogs() {
	if len(s.logBuffer) <= 0 {
		return
	}
	logGroup := &sls.LogGroup{
		Topic:  proto.String(s.topic),
		Source: proto.String(s.source),
		Logs:   s.logBuffer,
	}
	s.logStore.PutLogs(logGroup)
	s.logBuffer = s.logBuffer[0:0]
}
