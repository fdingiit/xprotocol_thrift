package sofabolt

import (
	"encoding/json"
	"sync"
)

//go:generate syncmap -pkg sofabolt -o pool_generated.go -name PoolMap map[string]*Pool

type Pool struct {
	sync.RWMutex
	clients []*Client
	next    int
}

func NewPool() *Pool {
	return &Pool{
		clients: make([]*Client, 0, 8),
	}
}

func (p *Pool) Size() int {
	var n int
	p.RLock()
	n = len(p.clients)
	p.RUnlock()
	return n
}

func (p *Pool) Iterate(fn func(client *Client)) {
	p.Lock()
	for i := range p.clients {
		fn(p.clients[i])
	}
	p.Unlock()
}

func (p *Pool) Delete(client *Client) {
	var i int
	p.Lock()
	for i = 0; i < len(p.clients); i++ {
		if client == p.clients[i] {
			p.clients = append(p.clients[:i], p.clients[i+1:]...)
			break
		}
	}
	p.Unlock()
}

func (p *Pool) Push(client *Client) {
	p.Lock()
	p.clients = append(p.clients, client)
	p.Unlock()
}

func (p *Pool) Get() (*Client, bool) {
	var (
		client *Client
		n      int
	)

	p.Lock()
	n = len(p.clients)
	if n == 0 {
		p.Unlock()
		return nil, false
	}
	p.next = (p.next + 1) % n
	client = p.clients[p.next]
	p.Unlock()

	return client, true
}

func (p *Pool) MarshalJSON() ([]byte, error) {
	type clientStatus struct {
		Closed          bool  `json:"closed"`
		Ref             int64 `json:"ref"`
		Lasted          int64 `json:"lasted"`
		Used            int64 `json:"used"`
		Created         int64 `json:"created"`
		PendingRequests int64 `json:"pending_requests"`
	}

	type status struct {
		Next    int            `json:"next"`
		Clients []clientStatus `json:"clients"`
	}

	s := status{
		Clients: make([]clientStatus, 0, 8),
	}

	p.RLock()
	s.Next = p.next
	for i := 0; i < len(p.clients); i++ {
		s.Clients = append(s.Clients, clientStatus{
			Lasted:          p.clients[i].GetMetrics().GetLasted(),
			Closed:          p.clients[i].Closed(),
			Ref:             p.clients[i].GetMetrics().GetReferences(),
			Used:            p.clients[i].GetMetrics().GetUsed(),
			Created:         p.clients[i].GetMetrics().GetCreated(),
			PendingRequests: p.clients[i].GetMetrics().GetPendingCommands(),
		})
	}
	p.RUnlock()

	return json.Marshal(s)
}
