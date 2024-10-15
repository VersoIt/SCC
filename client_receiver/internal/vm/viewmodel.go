package vm

import (
	"client_receiver/internal/service"
	"context"
)

type ViewModel struct {
	service *service.Service
}

func NewViewModel(service *service.Service) *ViewModel {
	return &ViewModel{service: service}
}

func (v *ViewModel) TransmitChunks(ctx context.Context) {
	v.service.TransmitChunks(ctx)
}

func (v *ViewModel) GetErrorChan() <-chan error {
	return v.service.GetStreamErrorChan()
}

func (v *ViewModel) GetBytesReadChan() <-chan int64 {
	return v.service.GetBytesReadChan()
}
