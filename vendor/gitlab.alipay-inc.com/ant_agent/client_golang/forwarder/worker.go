package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Worker consumes Transaction from the Forwarder and process them.
// If the transaction fail to be processed, the Worker will send it back to the Forward
// to be retried later.
type Worker struct {
	Client       *http.Client       // Client the http client used to process transactions.
	HighPrioChan <-chan Transaction // HighPrioChan is the channel used to receive high priority transaction from the Forwarder.
	LowPrioChan  <-chan Transaction // LowPrioChan is the channel used to receive low priority transaction from the Forwarder.
	RequeueChan  chan<- Transaction // RequeueChan is the channel used to send failed transaction back to the Forwarder.

	stopChan chan bool
	stopped  chan struct{}
	logger   Logger
}

// NewWorker returns a new worker to consume Transaction from inputChan and push erroneous ones into requeueChan.
func NewWorker(highPrioChan, lowPrioChan <-chan Transaction, requeueChan chan<- Transaction, logger Logger) *Worker {
	cli := &http.Client{
		Timeout:   30 * time.Second,
		Transport: HTTPTransport(),
	}

	w := &Worker{
		Client:       cli,
		HighPrioChan: highPrioChan,
		LowPrioChan:  lowPrioChan,
		RequeueChan:  requeueChan,
		stopChan:     make(chan bool),
		stopped:      make(chan struct{}),
		logger:       logger,
	}
	return w
}

// Stop stops the worker.
func (w *Worker) Stop() {
	w.stopChan <- true
	<-w.stopped
}

// Start starts a worker.
func (w *Worker) Start() {
	go func() {
		// notify that the worker did stop
		defer close(w.stopped)

		for {
			// handling high priority transactions first
			select {
			case t := <-w.HighPrioChan:
				if w.callProcess(t) == nil {
					continue
				}
				return
			case <-w.stopChan:
				return
			default:
			}

			select {
			case t := <-w.HighPrioChan:
				if w.callProcess(t) != nil {
					return
				}
			case t := <-w.LowPrioChan:
				if w.callProcess(t) != nil {
					return
				}
			case <-w.stopChan:
				return
			}
		}
	}()
}

// callProcess will process a transaction and cancel it if we need to stop the worker.
func (w *Worker) callProcess(t Transaction) error {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.process(ctx, t)
		done <- struct{}{}
	}()

	select {
	case <-done:
	// wait for the Transaction process to be over
	case <-w.stopChan:
		cancel() // cancel current Transaction if we need to stop the worker.
		<-done   // still need to wait for the process func to return
		return fmt.Errorf("the Worker was requested to stop")
	}
	cancel()
	return nil
}

func (w *Worker) process(ctx context.Context, t Transaction) {
	requeue := func() {
		select {
		case w.RequeueChan <- t:
		default:
			w.logger.Errorf("dropping transaction because the retry goroutines is too busy to handle another one")
		}
	}

	if err := t.Process(ctx, w.Client); err != nil {
		requeue()
		w.logger.Errorf("failed to process the transaction: %s", err)
	}
}
