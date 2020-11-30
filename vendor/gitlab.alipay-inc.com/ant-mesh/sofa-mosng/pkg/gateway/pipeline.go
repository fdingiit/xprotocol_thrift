package gateway

import (
	"context"
	"runtime/debug"
	"sync"

	"mosn.io/mosn/pkg/log"
)

type DefaultPipeline struct {
	Handlers     []Handler
	HandlerNames []string
}

var registerHandlerNames []string
var registerHandlers []Handler

var once sync.Once

func NewPipeline(handlerNames []string) *DefaultPipeline {
	once.Do(func() {
		len := len(handlerNames)
		registerHandlers = make([]Handler, 0, len)
		registerHandlerNames = make([]string, 0, len)
		for _, name := range handlerNames {
			if handler := GetHandler(name); handler != nil {
				registerHandlers = append(registerHandlers, handler)
				registerHandlerNames = append(registerHandlerNames, name)
			}
		}
	})

	p := &DefaultPipeline{
		Handlers:     registerHandlers,
		HandlerNames: registerHandlerNames,
	}

	return p
}

func (p *DefaultPipeline) RunInHandlers(ctx context.Context) (bool, error) {
	for i := 0; i < len(p.Handlers); i++ {
		GetGatewayContext(ctx).SetIndex(i)
		status, err := p.Handlers[i].HandleIn(ctx)
		if status == HandleStatusStopAndReturn {
			return true, err
		}
	}

	return false, nil
}

func (p *DefaultPipeline) RunOutHandlers(ctx context.Context) (bool, error) {
	gwCtx := GetGatewayContext(ctx)

	defer func() {
		if err := recover(); err != nil {
			log.Proxy.Errorf(ctx, "[gateway][%s][pipeline] handle out occurred error , %+v\n%s", gwCtx.UniqueId(), err, string(debug.Stack()))
		}
	}()

	for i := gwCtx.CurrentIndex(); i >= 0; i-- {
		gwCtx.SetIndex(i)
		status, err := p.Handlers[i].HandleOut(ctx)
		if status == HandleStatusStopAndReturn {
			if gwErr, ok := err.(*GatewayError); ok {
				response := NewGatewayResponse(gwErr.Headers, gwErr.DataBuf, gwErr.Trailers)
				response.SetResultStatus(gwErr.ResultStatus)
				gwCtx.SetResponse(response)
			}
			if log.Proxy.GetLogLevel() >= log.WARN {
				log.Proxy.Warnf(ctx, "[gateway][%s][pipeline] Handler[%s] handle out occurred error, %+v", gwCtx.UniqueId(), p.Handlers[i].Name(), err)
			}
		}
	}
	return false, nil
}

//func (p *DefaultPipeline) AddFirst(handler Handler) Pipeline {
//	p.Handlers = insertHandlerAtIndex(p.Handlers, []Handler{handler}, 0)
//	p.HandlerNames = insertNameAtIndex(p.HandlerNames, []string{handler.Name()}, 0)
//	return p
//}
//
//func (p *DefaultPipeline) AddLast(handler Handler) Pipeline {
//	p.HandlerNames = append(p.HandlerNames, handler.Name())
//	p.Handlers = append(p.Handlers, handler)
//	return p
//}
//
//func (p *DefaultPipeline) AddBefore(baseName string, handler Handler) Pipeline {
//	index := indexOf(baseName, p.HandlerNames)
//	p.Handlers = insertHandlerAtIndex(p.Handlers, []Handler{handler}, index)
//	p.HandlerNames = insertNameAtIndex(p.HandlerNames, []string{handler.Name()}, index)
//	return p
//}
//
//func (p *DefaultPipeline) AddAfter(baseName string, handler Handler) Pipeline {
//	index := indexOf(baseName, p.HandlerNames) + 1
//	p.Handlers = insertHandlerAtIndex(p.Handlers, []Handler{handler}, index)
//	p.HandlerNames = insertNameAtIndex(p.HandlerNames, []string{handler.Name()}, index)
//	return p
//}
//
//func (p *DefaultPipeline) AddListFirst(handlers ...Handler) Pipeline {
//	for _, handler := range handlers {
//		p.AddFirst(handler)
//	}
//	return p
//}
//
//func (p *DefaultPipeline) AddListLast(handlers ...Handler) Pipeline {
//	for _, handler := range handlers {
//		p.AddLast(handler)
//	}
//	return p
//}
//
//func (p *DefaultPipeline) Copy() Pipeline {
//	return &DefaultPipeline{
//		Handlers:     p.Handlers,
//		HandlerNames: p.HandlerNames,
//	}
//}
//
//func insertHandlerAtIndex(slice, insertion []Handler, index int) []Handler {
//	result := make([]Handler, len(slice)+len(insertion))
//	at := copy(result, slice[:index])
//	at += copy(result[at:], insertion)
//	copy(result[at:], slice[index:])
//	return result
//}
//
//func insertNameAtIndex(slice, insertion []string, index int) []string {
//	result := make([]string, len(slice)+len(insertion))
//	at := copy(result, slice[:index])
//	at += copy(result[at:], insertion)
//	copy(result[at:], slice[index:])
//	return result
//}
//
//func indexOf(name string, names []string) int {
//	for i, v := range names {
//		if v == name {
//			return i
//		}
//	}
//	return -1
//}
