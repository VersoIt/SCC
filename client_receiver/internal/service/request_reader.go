package service

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
)

const (
	requestBuffer       = 1024
	requestsCountAtTime = 1 << 20
)

type RequestReader struct {
	bytesChan          chan []byte
	endpoint           string
	errorChan          chan error
	bytesReadCountChan chan int64
	readBytes          atomic.Int64
	requests           chan func()
	wg                 sync.WaitGroup
}

func NewRequestReader(endpoint string) *RequestReader {
	rr := &RequestReader{bytesChan: make(chan []byte), endpoint: endpoint, bytesReadCountChan: make(chan int64), errorChan: make(chan error), requests: make(chan func(), requestsCountAtTime)}
	rr.launchRequestListener()
	return rr
}

func (rr *RequestReader) GetErrorChan() <-chan error {
	return rr.errorChan
}

func (rr *RequestReader) GetBytesChan() chan []byte {
	return rr.bytesChan
}

func (rr *RequestReader) launchRequestListener() {
	rr.wg.Add(1)
	go func() {
		defer rr.wg.Done()
		for request := range rr.requests {
			request()
		}
	}()
}

func (rr *RequestReader) ReadToChan() {
	if rr.bytesChan == nil {
		panic("bytes chan must not be nil")
	}

	rr.requests <- func() {
		response, err := http.Get(rr.endpoint)
		if err != nil {
			rr.errorChan <- err
			return
		}

		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				rr.errorChan <- err
			}
		}(response.Body)

		reader := bufio.NewReader(response.Body)
		for {
			bytes := make([]byte, requestBuffer)
			bytesRead, err := reader.Read(bytes)

			if err == io.EOF {
				return
			}

			if err != nil {
				rr.errorChan <- err
				return
			}

			logrus.Infof("read bytes: %d", bytesRead)

			rr.readBytes.Add(int64(bytesRead))
			rr.bytesChan <- bytes[:bytesRead]
			rr.bytesReadCountChan <- rr.readBytes.Load()
		}
	}
}

func (rr *RequestReader) GetReadBytesCountChan() <-chan int64 {
	return rr.bytesReadCountChan
}

func (rr *RequestReader) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		rr.wg.Wait()
		close(rr.requests)
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
