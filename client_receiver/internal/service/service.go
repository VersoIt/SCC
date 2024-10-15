package service

import (
	"context"
)

type stream interface {
	GetErrorChan() <-chan error
}

type byteFlowReader interface {
	ReadToChan(context.Context)
	GetBytesChan() chan []byte
	GetReadBytesCountChan() <-chan int64
	stream
}

type byteFlowWriter interface {
	WriteFromChan(context.Context)
	SetBytesChan(chan []byte)
	stream
}

type Service struct {
	byteFlowReader byteFlowReader
	byteFlowWriter byteFlowWriter
	launched       bool
	errorChan      chan error
}

func NewService(reader byteFlowReader, writer byteFlowWriter) *Service {
	return &Service{byteFlowReader: reader, byteFlowWriter: writer, launched: false, errorChan: make(chan error)}
}

func (s *Service) TransmitChunks(ctx context.Context) {
	s.byteFlowReader.ReadToChan(ctx)
	if s.launched {
		return
	}
	s.launched = true

	s.byteFlowWriter.SetBytesChan(s.byteFlowReader.GetBytesChan())
	s.byteFlowWriter.WriteFromChan(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-s.byteFlowReader.GetErrorChan():
				s.errorChan <- err
			case err := <-s.byteFlowWriter.GetErrorChan():
				s.errorChan <- err
			}
		}
	}()
}

func (s *Service) GetStreamErrorChan() <-chan error {
	return s.errorChan
}

func (s *Service) GetBytesReadChan() <-chan int64 {
	return s.byteFlowReader.GetReadBytesCountChan()
}
