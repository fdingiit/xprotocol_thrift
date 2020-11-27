package v1

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/utils"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
)

func (r *Router) Diff(new *Router) []event.DifferEvent {

	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	// diff timeout
	if r.Timeout != new.Timeout {
		return events
	}

	// diff filter chain
	if !md5Diff(r.FilterChains, new.FilterChains) {
		return events
	}

	// diff match
	if !md5Diff(r.Matches, new.Matches) {
		return events
	}

	// diff proxy
	if !md5Diff(r.Proxy, new.Proxy) {
		return events
	}

	return []event.DifferEvent{}
}

func (rg RouterGroup) Diff(new RouterGroup) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if !md5Diff(rg.Gateways, new.Gateways) {
		return events
	}

	return []event.DifferEvent{}
}

func (s GatewayService) Diff(new *GatewayService) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if s.Timeout != new.Timeout {
		return events
	}

	if !md5Diff(s.FilterChains, new.FilterChains) {
		return events
	}

	if s.Protocol != new.Protocol {
		return events
	}

	if s.LbType != new.LbType {
		return events
	}

	if !md5Diff(s.Hosts, new.Hosts) {
		return events
	}

	if !md5Diff(s.Filters, new.Filters) {
		return events
	}

	return []event.DifferEvent{}
}

func (gf GlobalFilter) Diff(new GlobalFilter) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if len(gf.Filters) != len(new.Filters) {
		return events
	}

	if !md5Diff(gf.Filters, new.Filters) {
		return events
	}

	return []event.DifferEvent{}
}

func (fc FilterChain) Diff(new *FilterChain) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if len(fc.Filters) != len(new.Filters) {
		return events
	}

	if !md5Diff(fc.Filters, new.Filters) {
		return events
	}

	return []event.DifferEvent{}
}

func (s Gateway) Diff(new *Gateway) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if !md5Diff(s, new) {
		return events
	}
	return []event.DifferEvent{}
}

func (c Metadata) Diff(new *Metadata) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if !md5Diff(c.Value, new.Value) {
		return events
	}

	return []event.DifferEvent{}
}

func (g GwConfig) Diff(new *GwConfig) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    new,
		EventType: event.Update,
	}}

	if !md5Diff(g, new) {
		return events
	}

	return []event.DifferEvent{}
}

func md5Diff(old, new interface{}) bool {
	return utils.Md5Check(old, new)
}
