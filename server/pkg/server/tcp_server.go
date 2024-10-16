package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

type Handler interface {
	Handle(ReadWriteConn)
}

type TcpServer struct {
	listener net.Listener
	wg       sync.WaitGroup
	quit     chan struct{}
	handlers []func(ReadWriteConn)
}

type ReadWriteConn interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

func NewTcpServer(port string, handlers ...func(ReadWriteConn)) (*TcpServer, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	return &TcpServer{
		listener: listener,
		quit:     make(chan struct{}),
		handlers: handlers,
	}, nil
}

func (s *TcpServer) Run() error {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return net.ErrClosed
			default:
				logrus.Errorf("error accepting connection: %v", err)
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

func (s *TcpServer) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	for _, handler := range s.handlers {
		handler(conn)
	}
}

func (s *TcpServer) Shutdown(ctx context.Context) error {
	close(s.quit)
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(s.listener)

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
