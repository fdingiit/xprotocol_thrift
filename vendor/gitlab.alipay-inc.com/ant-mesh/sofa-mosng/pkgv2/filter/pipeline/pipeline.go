package pipeline

import (
	"sort"
	"sync"

	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"mosn.io/mosn/pkg/protocol"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/filter"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

type pipeline struct {
	fMux    sync.RWMutex
	filters sortableFilters
	eh      types.ErrorHandler
	flag    int32
}

// 根据 Priority 进行排序，数值小在前面
type sortableFilters []types.GatewayFilter

func (f sortableFilters) Len() int           { return len(f) }
func (f sortableFilters) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f sortableFilters) Less(i, j int) bool { return f[i].Priority() < f[j].Priority() }
func (f sortableFilters) copy() sortableFilters {
	var gfs []types.GatewayFilter
	for _, gf := range f {
		copyf := new(types.GatewayFilter)
		*copyf = gf
		gfs = append(gfs, gf)
	}

	return gfs
}

const (
	FLAG types.AttributeKey = "x-mosng-pipeline-flag"
)

func (p *pipeline) DoInBound(ctx types.Context) error {
	p.fMux.RLock()
	defer p.fMux.RUnlock()
	for _, f := range p.filters {
		p.incIndex(ctx)
		if status, err := f.InBound(ctx); status == types.Error {
			return err
		}
	}
	return nil
}

func (p *pipeline) DoOutBound(ctx types.Context) error {
	p.fMux.RLock()
	defer p.fMux.RUnlock()

	for i := p.getIndex(ctx); i >= 0; i-- {
		p.subIndex(ctx)
		if status, err := p.filters[i].OutBound(ctx); status == types.Error {
			return err
		}
	}
	return nil
}

func (p *pipeline) AddNewFilter(filters []types.GatewayFilter) error {
	if filters == nil {
		//todo err
		return nil
	}

	p.fMux.Lock()
	defer p.fMux.Unlock()

	p.filters = append(p.filters, filters...)

	// todo err
	return nil
}

func (p *pipeline) AddOrUpdateFilter(filter types.GatewayFilter) error {
	p.fMux.Lock()
	defer p.fMux.Unlock()

	p.filters = append(p.filters, filter)

	// todo err
	return nil
}

func (p *pipeline) DelFilter(filter types.GatewayFilter) error {
	p.fMux.Lock()
	defer p.fMux.Unlock()

	// todo del
	return nil
}

func (p *pipeline) GetFilters() []types.GatewayFilter {
	p.fMux.RLock()
	defer p.fMux.RUnlock()
	return p.filters
}

func (p *pipeline) Sort() {
	p.fMux.Lock()
	defer p.fMux.Unlock()

	sort.Stable(p.filters)
}

func (p *pipeline) Copy() types.Pipeline {
	cp := &pipeline{fMux: sync.RWMutex{}, filters: p.filters.copy(), eh: p.eh}
	return cp
}

func (p *pipeline) HandleErr(ctx types.Context, err types.GatewayError) (httpCode int, headers protocol.CommonHeader, res []byte) {
	return p.eh.Handle(ctx, err)
}

func (p *pipeline) SetErrHandler(eh types.ErrorHandler) {
	p.eh = eh
}

func (p *pipeline) incIndex(ctx types.Context) {
	var i int
	index := ctx.GetAttribute(FLAG)
	if index == nil {
		i = 0
	} else {
		i = index.(int)
		i++
	}
	ctx.SetAttribute(FLAG, i)
}

func (p *pipeline) subIndex(ctx types.Context) {
	var i int
	index := ctx.GetAttribute(FLAG)
	if index == nil {
		i = 0
	} else {
		i = index.(int)
		i--
	}
	ctx.SetAttribute(FLAG, i)
}

func (p *pipeline) getIndex(ctx types.Context) int {
	index := ctx.GetAttribute(FLAG)
	if index == nil {
		index = 0
	}
	return index.(int)
}

func NewPipeline(rf []*v1.Filter, rc []types.GatewayFilter, sf []*v1.Filter, sc []types.GatewayFilter, gp []types.GatewayFilter) types.Pipeline {
	var filters []types.GatewayFilter

	if rf != nil {
		filters = appendFilter(filters, rf)
	}

	if rc != nil {
		filters = append(filters, rc...)
	}

	if sf != nil {
		filters = appendFilter(filters, rf)
	}

	if sc != nil {
		filters = append(filters, sc...)
	}

	if gp != nil {
		filters = append(filters, gp...)
	}

	return &pipeline{fMux: sync.RWMutex{}, filters: filters, eh: errors.GetDefaultErrorHandler()}
}

func appendFilter(filters []types.GatewayFilter, rf []*v1.Filter) []types.GatewayFilter {
	for _, conf := range rf {
		if gatewayFilter, e := filter.New(conf); e == nil {
			filters = append(filters, gatewayFilter)
		} else {
			panic(errors.Errorf("create filter error => filter name: %s; filter config: %v", conf.Name, conf.Metadata))
		}
	}
	return filters
}
