package gateway

type ApiService struct {
	serviceKey string
	status     ApiStatus
	timeout    uint64
	rules      map[string]interface{}
	upstream   Upstream
}

func (a *ApiService) ServiceKey() string {
	return a.serviceKey
}

func (a *ApiService) SetServiceKey(serviceKey string) {
	a.serviceKey = serviceKey
}

func (a *ApiService) Status() ApiStatus {
	return a.status
}

func (a *ApiService) SetStatus(status ApiStatus) {
	a.status = status
}

func (a *ApiService) Timeout() uint64 {
	return a.timeout
}

func (a *ApiService) SetTimeout(timeout uint64) {
	a.timeout = timeout
}

func (a *ApiService) GetAttribute(key string) interface{} {
	return a.rules[key]
}

func (a *ApiService) SetAttribute(key string, value interface{}) {
	if a.rules == nil {
		a.rules = make(map[string]interface{})
	}
	a.rules[key] = value
}

func (a *ApiService) Upstream() Upstream {
	return a.upstream
}

func (a *ApiService) SetUpstream(upstream Upstream) {
	a.upstream = upstream
}

