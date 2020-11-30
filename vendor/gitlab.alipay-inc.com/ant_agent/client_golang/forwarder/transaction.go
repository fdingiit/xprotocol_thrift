package forwarder

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Transaction represents the task to process for a worker.
type Transaction interface {
	Process(ctx context.Context, cli *http.Client) error
	GetCreatedAt() time.Time
}

// HTTPTransaction represents one Payload for one Endpoint and one Domain.
type HTTPTransaction struct {
	Domain   string      // Domain represents the domain target by the HTTPTransaction
	Endpoint string      // Endpoint is the API Endpoint used by the HTTPTransaction
	Header   http.Header // Header is the HTTP Header used by the HTTPTransaction
	Payload  []byte      // Payload is the content delivered to the backend.

	logger    Logger
	createdAt time.Time
}

func NewHTTPTransaction(domain, endpoint string, header http.Header, payload []byte, logger Logger) *HTTPTransaction {
	return &HTTPTransaction{
		Domain:    domain,
		Endpoint:  endpoint,
		Header:    header,
		Payload:   payload,
		createdAt: time.Now(),
		logger:    logger,
	}
}

func (t *HTTPTransaction) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *HTTPTransaction) Process(ctx context.Context, cli *http.Client) error {
	url := t.Domain + t.Endpoint
	reader := bytes.NewReader(t.Payload)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		t.logger.Errorf("failed to new request for transaction to invalid URL %q, dropping transaction: %s", url, err)
		return nil
	}
	req = req.WithContext(ctx)
	req.Header = t.Header

	resp, err := cli.Do(req)
	if err != nil {
		// do not requeue transaction if that one was canceled
		if ctx.Err() == context.Canceled {
			t.logger.Warnf("the transaction was canceled, dropping")
			return nil
		}
		return fmt.Errorf("failed to send transaction, rescheduling it: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.logger.Errorf("failed to read the response Body: %s", err)
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send transaction to %q, rescheduling it: %q", url, resp.Status)
	}
	t.logger.Debugf("successfully post payload to %q: %s", url, string(body))
	return nil
}
