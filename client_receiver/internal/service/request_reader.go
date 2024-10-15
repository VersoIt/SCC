package service

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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
	launched           bool
	requests           chan func()
}

func NewRequestReader(endpoint string) *RequestReader {
	return &RequestReader{bytesChan: make(chan []byte), endpoint: endpoint, bytesReadCountChan: make(chan int64), errorChan: make(chan error), launched: false, requests: make(chan func(), requestsCountAtTime)}
}

func (rr *RequestReader) GetErrorChan() <-chan error {
	return rr.errorChan
}

func (rr *RequestReader) GetBytesChan() chan []byte {
	return rr.bytesChan
}

func (rr *RequestReader) ReadToChan(ctx context.Context) {
	if rr.bytesChan == nil {
		panic("bytes chan must not be nil")
	}

	if !rr.launched {
		rr.launched = true
		go func() {
			for request := range rr.requests {
				request()
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}()
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
