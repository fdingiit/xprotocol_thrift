package sofabolt

import "time"

// serverOptionSetter configures a Server.
type serverOptionSetter interface {
	set(*Server)
}

type serverOptionSetterFunc func(*Server)

func (f serverOptionSetterFunc) set(srv *Server) {
	f(srv)
}

func WithServerMetrics(sm *ServerMetrics) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.metrics = sm
	})
}

func WithServerHandler(fn Handler) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.handler = fn
	})
}

func WithServerTimeout(readTimeout, writeTimeout, idleTimeout, flushInterval time.Duration) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.options.readTimeout = readTimeout
		srv.options.writeTimeout = writeTimeout
		srv.options.idleTimeout = idleTimeout
		srv.options.flushInterval = flushInterval
	})
}

func WithServerMaxConnctions(m int) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.options.maxConnections = m
	})
}

func WithServerMaxPendingCommands(m int) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.options.maxPendingCommand = m
	})
}

func WithServerOnEventHandler(e ServerOnEventHandler) serverOptionSetter {
	return serverOptionSetterFunc(func(srv *Server) {
		srv.onhandler = e
	})
}
