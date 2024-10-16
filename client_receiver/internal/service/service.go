package service

import (
	"context"
)

type stream interface {
	GetErrorChan() <-chan error
	Shutdown(ctx context.Context) error
}

type byteFlowReader interface {
	ReadToChan()
	GetBytesChan() chan []byte
	GetReadBytesCountChan() <-chan int64
	stream
}

type byteFlowWriter interface {
	StartWritingFromChan(context.Context)
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
	s.byteFlowReader.ReadToChan()
	if s.launched {
		return
	}
	s.launched = true

	s.byteFlowWriter.SetBytesChan(s.byteFlowReader.GetBytesChan())
	s.byteFlowWriter.StartWritingFromChan(ctx)

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

func (s *Service) ShutdownStreams(ctx context.Context) error {
	if err := s.byteFlowReader.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.byteFlowWriter.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
