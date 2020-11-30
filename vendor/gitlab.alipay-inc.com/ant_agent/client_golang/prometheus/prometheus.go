package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"gitlab.alipay-inc.com/ant_agent/client_golang/forwarder"
)

const (
	TenantNameKey         = "_tenant_name"
	WorkspaceNameKey      = "_workspace_name"
	ResourceTypeKey       = "_resource_type"
	AppNameKey            = "app_name"
	AppInstanceNameKey    = "app_instance_name"
	AppInstanceVersionKey = "app_instance_version"
	ResourceType          = "CUSTOM"
)

const (
	ContentType = "application/json"
)

type EncodeFunc func() ([]byte, error)

type Forwarder struct {
	gather prometheus.Gatherer
	forwarder.Forwarder
	contentType expfmt.Format

	tenantName    string
	workspaceName string
	token         string
	interval      time.Duration
	labels        []*dto.LabelPair
	logger        forwarder.Logger
	stopChan      chan bool
	stopped       chan struct{}
	encodeFunc    EncodeFunc
}

func NewForwarder(domain, tenantName, workspaceName string, interval time.Duration, workerNum, chanSize int, flushInterval time.Duration,
	logger forwarder.Logger) *Forwarder {
	f := &Forwarder{
		Forwarder:     forwarder.NewDefaultForwarder(domain, workerNum, chanSize, flushInterval, logger),
		tenantName:    tenantName,
		workspaceName: workspaceName,
		contentType:   ContentType,
		interval:      interval,
		logger:        logger,
	}
	f.init()
	return f
}

func NewDefaultForwarder(domain, tenantName, workspaceName string, interval time.Duration, logger forwarder.Logger) *Forwarder {
	return NewForwarder(domain, tenantName, workspaceName, interval, forwarder.WorkerNumber, forwarder.ChanSize, forwarder.FlushInterval, logger)
}

func (f *Forwarder) init() {
	f.stopChan = make(chan bool)
	f.stopped = make(chan struct{})
	f.labels = []*dto.LabelPair{
		&dto.LabelPair{
			Name:  proto.String(TenantNameKey),
			Value: proto.String(f.tenantName),
		},
		&dto.LabelPair{
			Name:  proto.String(WorkspaceNameKey),
			Value: proto.String(f.workspaceName),
		},
		{
			Name:  proto.String(ResourceTypeKey),
			Value: proto.String(ResourceType),
		},
	}
	if f.contentType == ContentType {
		f.encodeFunc = f.CustomEncode
	} else {
		f.encodeFunc = f.Encode
	}
}

func (f *Forwarder) Start() error {
	_ = f.Forwarder.Start()
	go f.run()
	return nil
}

func (f *Forwarder) Stop() {
	f.stopChan <- true
	<-f.stopped
	f.Forwarder.Stop()
}

func (f *Forwarder) run() {
	ticker := time.NewTicker(f.interval)
	for {
		select {
		case <-ticker.C:
			if err := f.Submit(); err != nil {
				f.logger.Errorf(err.Error())
			}
		case <-f.stopChan:
			ticker.Stop()
			f.stopped <- struct{}{}
			return
		}
	}
}

func (f *Forwarder) AddLabelPairs(pairs ...string) error {
	if len(pairs)%2 != 0 {
		return fmt.Errorf("labels must be key value pairs")
	}
	var labels []*dto.LabelPair
	for i := 0; i < len(pairs); i += 2 {
		labels = append(labels, &dto.LabelPair{
			Name:  proto.String(pairs[i]),
			Value: proto.String(pairs[i+1]),
		})
	}
	f.labels = append(f.labels, labels...)
	return nil
}

func (f *Forwarder) SetGather(gather prometheus.Gatherer) {
	f.gather = gather
}

func (f *Forwarder) SetContentType(contentType expfmt.Format) {
	f.contentType = contentType
	if f.contentType == ContentType {
		f.encodeFunc = f.CustomEncode
	} else {
		f.encodeFunc = f.Encode
	}
}

func (f *Forwarder) SetToken(token string) {
	f.token = token
}

func (f *Forwarder) Encode() ([]byte, error) {
	if f.gather == nil {
		return nil, fmt.Errorf("must set gather first")
	}
	mfs, err := f.gather.Gather()
	if err != nil {
		return nil, fmt.Errorf("faile to gather prometheus metrics: %s", err)
	}

	var b bytes.Buffer
	enc := expfmt.NewEncoder(&b, f.contentType)
	for _, mf := range mfs {
		f.addLabels(mf)
		if err := enc.Encode(mf); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

func (f *Forwarder) CustomEncode() ([]byte, error) {
	if f.gather == nil {
		return nil, fmt.Errorf("must set gather first")
	}
	mfs, err := f.gather.Gather()
	if err != nil {
		return nil, fmt.Errorf("faile to gather prometheus metrics: %s", err)
	}
	for _, mf := range mfs {
		f.addLabels(mf)
	}

	metrics := flatMetricFamilies(mfs)
	if len(metrics) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return b, nil
}

func (f *Forwarder) Submit() error {
	payload, err := f.encodeFunc()
	if err != nil {
		return err
	}
	header := make(http.Header)
	header.Set("Content-Type", string(f.contentType))
	header.Set("Authorization", f.token)
	if err := f.Forward(payload, header); err != nil {
		return err
	}
	return nil
}

func (f *Forwarder) addLabels(mf *dto.MetricFamily) {
	for _, m := range mf.Metric {
		m.Label = append(m.Label, f.labels...)
	}
}
