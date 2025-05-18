package server

import "time"

type Option func(srv *Server)

func SetReadTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.internalServer.ReadTimeout = duration
	}
}

func SetWriteTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.internalServer.WriteTimeout = duration
	}
}

func SetAddr(addr string) Option {
	return func(srv *Server) {
		srv.internalServer.Addr = addr
	}
}

func SetShutdownTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.shutdownTimeout = duration
	}
}
