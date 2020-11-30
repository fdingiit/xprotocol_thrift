package gateway

import (
	"reflect"
	"sync"

	"mosn.io/mosn/pkg/log"
)

var handlerFactory map[string]Handler
var configListenerFactory sync.Map

var upstreamParserFactory map[string]upstreamParserOrder
var downstreamCodecFactory map[string]downstreamCodecOrder
var upstreamCodecFactory map[string]upstreamCodecOrder
var downstreamStatusMappingFactory map[string]downstreamStatusMappingOrder
var upstreamStatusMappingFactory map[string]upstreamStatusMappingOrder

var serviceRuleParserFactory map[string]ServiceRuleParser
var gatewayRuleTypeFactory sync.Map
var appSpecTypeFactory sync.Map
var clusterSpecTypeFactory sync.Map

var configFilterFactory []ConfigFilter

var configListenerMutex = new(sync.Mutex)

func init() {
	handlerFactory = make(map[string]Handler)
	configListenerFactory = sync.Map{}
	gatewayRuleTypeFactory = sync.Map{}
	appSpecTypeFactory = sync.Map{}
	clusterSpecTypeFactory = sync.Map{}
	upstreamParserFactory = make(map[string]upstreamParserOrder)
	serviceRuleParserFactory = make(map[string]ServiceRuleParser)
	downstreamCodecFactory = make(map[string]downstreamCodecOrder)
	downstreamStatusMappingFactory = make(map[string]downstreamStatusMappingOrder)
	upstreamStatusMappingFactory = make(map[string]upstreamStatusMappingOrder)
	upstreamCodecFactory = make(map[string]upstreamCodecOrder)
}

func AddGatewayConfigFilter(filter ConfigFilter) {
	configFilterFactory = append(configFilterFactory, filter)
}

func RegisterHandler(handler Handler) {
	handlerFactory[handler.Name()] = handler
}

func GetHandler(name string) Handler {
	return handlerFactory[name]
}

func RegisterGatewayRuleType(configName string, control interface{}) {
	gatewayRuleTypeFactory.Store(configName, reflect.TypeOf(control).Elem())

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][factory] register gateway rule type [%s]", configName)
	}
}

func RegisterAppSpecType(configName string, control interface{}) {
	appSpecTypeFactory.Store(configName, reflect.TypeOf(control).Elem())

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][factory] register app spec type [%s]", configName)
	}
}

func RegisterClusterSpecType(configName string, control interface{}) {
	clusterSpecTypeFactory.Store(configName, reflect.TypeOf(control).Elem())

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][factory] register cluster spec type [%s]", configName)
	}
}

func RegisterConfigListener(name string, listener ConfigListener) {
	configListenerMutex.Lock()
	defer func() {
		configListenerMutex.Unlock()
	}()

	var newListeners []ConfigListener
	if value, ok := configListenerFactory.Load(name); ok {
		if listeners, ok := value.([]ConfigListener); ok {
			newListeners = listeners
		}
	}

	newListeners = append(newListeners, listener)
	configListenerFactory.Store(name, newListeners)

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][factory] register config listener [%s]", name)
	}
}

type upstreamParserOrder struct {
	parser UpstreamParser
	order  uint32
}

func RegisterUpstreamParser(protocol string, parser UpstreamParser, order uint32) {
	parserOrder := upstreamParserOrder{
		parser: parser,
		order:  order,
	}

	if exist, ok := upstreamParserFactory[protocol]; ok {
		if exist.order < order {
			upstreamParserFactory[protocol] = parserOrder

			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][factory] register upstream[%s] parser, order %d", protocol, order)
			}
		} else {
			if log.DefaultLogger.GetLogLevel() >= log.WARN {
				log.DefaultLogger.Warnf("[gateway][factory] current UpstreamParser[%s] register order is %d, order %d register failed", protocol, exist.order, order)
			}
		}
	} else {
		upstreamParserFactory[protocol] = parserOrder
	}
}

func GetUpstreamParser(protocol string) UpstreamParser {
	if order, ok := upstreamParserFactory[protocol]; ok {
		return order.parser
	}
	return nil
}

func RegisterServiceRuleParser(ruleName string, parser ServiceRuleParser) {
	serviceRuleParserFactory[ruleName] = parser

	if log.DefaultLogger.GetLogLevel() >= log.INFO {
		log.DefaultLogger.Infof("[gateway][factory] register service rule parser [%s]", ruleName)
	}
}

type downstreamCodecOrder struct {
	codec DownstreamCodec
	order uint32
}

func RegisterDownstreamCodec(protocol string, codec DownstreamCodec, order uint32) {
	codecOrder := downstreamCodecOrder{
		codec: codec,
		order: order,
	}

	if exist, ok := downstreamCodecFactory[protocol]; ok {
		if exist.order < order {
			downstreamCodecFactory[protocol] = codecOrder

			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][factory] register downstream[%s] codec, order %d", protocol, order)
			}
		} else {
			if log.DefaultLogger.GetLogLevel() >= log.WARN {
				log.DefaultLogger.Warnf("[gateway][factory] current DownstreamCodec[%s] register order is %d, order %d register failed", protocol, exist.order, order)
			}
		}
	} else {
		downstreamCodecFactory[protocol] = codecOrder
	}
}

func GetDownstreamCodec(protocol string) DownstreamCodec {
	if order, ok := downstreamCodecFactory[protocol]; ok {
		return order.codec
	}
	return nil
}

type downstreamStatusMappingOrder struct {
	mapping DownstreamStatusMapping
	order   uint32
}

func RegisterDownstreamStatusMapping(protocol string, mapping DownstreamStatusMapping, order uint32) {
	mappingOrder := downstreamStatusMappingOrder{
		mapping: mapping,
		order:   order,
	}

	if exist, ok := downstreamStatusMappingFactory[protocol]; ok {
		if exist.order < order {
			downstreamStatusMappingFactory[protocol] = mappingOrder

			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][factory] register downstream[%s] status mapping, order %d", protocol, order)
			}
		} else {
			if log.DefaultLogger.GetLogLevel() >= log.WARN {
				log.DefaultLogger.Warnf("[gateway][factory] current downstream[%s] status mapping register order is %d, order %d register failed", protocol, exist.order, order)
			}
		}
	} else {
		downstreamStatusMappingFactory[protocol] = mappingOrder
	}
}

func GetDownstreamStatusMapping(protocol string) DownstreamStatusMapping {
	if order, ok := downstreamStatusMappingFactory[protocol]; ok {
		return order.mapping
	}
	return nil
}

type upstreamStatusMappingOrder struct {
	mapping UpstreamStatusMapping
	order   uint32
}

func RegisterUpstreamStatusMapping(protocol string, mapping UpstreamStatusMapping, order uint32) {
	mappingOrder := upstreamStatusMappingOrder{
		mapping: mapping,
		order:   order,
	}

	if exist, ok := upstreamStatusMappingFactory[protocol]; ok {
		if exist.order < order {
			upstreamStatusMappingFactory[protocol] = mappingOrder

			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][factory] register upstream[%s] status mapping, order %d", protocol, order)
			}
		} else {
			if log.DefaultLogger.GetLogLevel() >= log.WARN {
				log.DefaultLogger.Warnf("[gateway][factory] current upstream[%s] status mapping register order is %d, order %d register failed", protocol, exist.order, order)
			}
		}
	} else {
		upstreamStatusMappingFactory[protocol] = mappingOrder
	}
}

func GetUpstreamStatusMapping(protocol string) UpstreamStatusMapping {
	if order, ok := upstreamStatusMappingFactory[protocol]; ok {
		return order.mapping
	}
	return nil
}

type upstreamCodecOrder struct {
	codec UpstreamCodec
	order uint32
}

func RegisterUpstreamCodec(protocol string, codec UpstreamCodec, order uint32) {
	codecOrder := upstreamCodecOrder{
		codec: codec,
		order: order,
	}

	if exist, ok := upstreamCodecFactory[protocol]; ok {
		if exist.order < order {
			upstreamCodecFactory[protocol] = codecOrder

			if log.DefaultLogger.GetLogLevel() >= log.INFO {
				log.DefaultLogger.Infof("[gateway][factory] register upstream[%s] codec, order %d", protocol, order)
			}
		} else {
			if log.DefaultLogger.GetLogLevel() >= log.WARN {
				log.DefaultLogger.Warnf("[gateway][factory] current UpstreamCodec[%s] register order is %d, order %d register failed", protocol, exist.order, order)
			}
		}
	} else {
		upstreamCodecFactory[protocol] = codecOrder
	}
}

func GetUpstreamCodec(protocol string) UpstreamCodec {
	if order, ok := upstreamCodecFactory[protocol]; ok {
		return order.codec
	}
	return nil
}
