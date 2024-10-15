package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
)

type RequestWriter struct {
	bytesChan chan []byte
	writer    io.Writer
	launched  bool
	errorChan chan error
}

func NewRequestWriter(writer io.Writer) *RequestWriter {
	return &RequestWriter{writer: writer, launched: false, errorChan: make(chan error)}
}

func (rw *RequestWriter) GetErrorChan() <-chan error {
	return rw.errorChan
}

func (rw *RequestWriter) SetBytesChan(bytesChan chan []byte) {
	rw.bytesChan = bytesChan
}

func (rw *RequestWriter) WriteFromChan(ctx context.Context) {
	if rw.launched {
		return
	}
	if rw.bytesChan == nil {
		panic("bytes chan must not be nil")
	}

	rw.launched = true

	go func() {

		for bytes := range rw.bytesChan {
			_, err := rw.writer.Write(bytes)
			if err != nil {
				rw.errorChan <- err
			}
			logrus.Infof("wrote bytes: %d", len(bytes))
			select {
			case <-ctx.Done():
				rw.launched = false
				return
			default:
			}
		}
	}()
}
