package server

import (
	"context"
	"net/http"
	"time"
)

type HttpServer struct {
	srv *http.Server
}

const (
	ReadTimeout    = 5 * time.Second
	WriteTimeout   = 5 * time.Second
	MaxHeaderBytes = 1 << 20 // 1 MB
)

func (s *HttpServer) Run(port string, handler http.Handler) error {
	s.srv = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: MaxHeaderBytes,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
	}

	return s.srv.ListenAndServe()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
