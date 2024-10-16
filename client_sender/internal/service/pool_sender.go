package service

import (
	"context"
	"encoding/base64"
	"net"
	"sync"
	"sync/atomic"
)

type PoolSender struct {
	addr               string
	wg                 sync.WaitGroup
	bytesCountSent     atomic.Int64
	bytesCountSentChan chan int64
	exit               chan struct{}
	errorChan          chan error
}

func NewPoolSender(addr string) *PoolSender {
	return &PoolSender{addr: addr, bytesCountSentChan: make(chan int64), exit: make(chan struct{}), errorChan: make(chan error)}
}

func (s *PoolSender) Send(poolCount int) {
	for i := 0; i < poolCount; i++ {
		select {
		case <-s.exit:
			return
		default:
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			for {
				conn, err := net.Dial("tcp", s.addr)
				if err != nil {
					s.errorChan <- err
					return
				}

				dataToSend := base64.StdEncoding.EncodeToString([]byte(randomDataToSend()))

				bytesSent, err := conn.Write([]byte(dataToSend))
				if err != nil {
					s.errorChan <- err
					return
				}

				s.bytesCountSent.Add(int64(bytesSent))
				s.bytesCountSentChan <- s.bytesCountSent.Load()

				err = conn.Close()
				if err != nil {
					s.errorChan <- err
					return
				}

				select {
				case <-s.exit:
					return
				default:
				}
			}
		}()
	}
}

func (s *PoolSender) GetBytesCountSentChan() <-chan int64 {
	return s.bytesCountSentChan
}

func randomDataToSend() string {
	return "Ruslan"
}

func (s *PoolSender) GetErrorChan() <-chan error {
	return s.errorChan
}

func (s *PoolSender) Shutdown(ctx context.Context) error {
	close(s.exit)
	done := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
