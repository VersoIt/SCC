package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

type RequestWriter struct {
	bytesChan chan []byte
	writer    io.Writer
	errorChan chan error
	started   bool
	wg        sync.WaitGroup
}

func NewRequestWriter(writer io.Writer) *RequestWriter {
	return &RequestWriter{writer: writer, errorChan: make(chan error), started: false}
}

func (rw *RequestWriter) GetErrorChan() <-chan error {
	return rw.errorChan
}

func (rw *RequestWriter) SetBytesChan(bytesChan chan []byte) {
	rw.bytesChan = bytesChan
}

func (rw *RequestWriter) StartWritingFromChan(ctx context.Context) {
	if rw.bytesChan == nil {
		panic("bytes chan must not be nil")
	}

	if rw.started {
		panic("writer is already started")
	}
	rw.started = true

	rw.wg.Add(1)
	go func() {
		defer rw.wg.Done()
		for bytes := range rw.bytesChan {
			_, err := rw.writer.Write(bytes)
			if err != nil {
				rw.errorChan <- err
			}
			logrus.Infof("wrote bytes: %d", len(bytes))
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

func (rw *RequestWriter) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		rw.wg.Wait()
		close(done)
		close(rw.bytesChan)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
