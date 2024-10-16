package main

import (
	"client_receiver/internal/service"
	"client_receiver/internal/view"
	"client_receiver/internal/vm"
	"client_receiver/pkg/config"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	app := widgets.NewQApplication(len(os.Args), os.Args)

	output, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(output)

	requestReader := service.NewRequestReader(fmt.Sprintf("%s/api/employee/", cfg.Host))
	requestWriter := service.NewRequestWriter(output)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	viewModel := vm.NewViewModel(service.NewService(requestReader, requestWriter))
	window, button := view.BuildWindowWithButton("Тестовое задание", "Получить данные с сервера", func(bool) {
		viewModel.TransmitChunks(ctx)
	})

	go func() {
		for bytesRead := range viewModel.GetBytesReadChan() {
			button.SetText(fmt.Sprintf("%d bytes", bytesRead))
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		for err := range viewModel.GetErrorChan() {
			logrus.Error(err)
			button.SetText(err.Error())
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-exit
		window.Close()

		if err = viewModel.ShutdownStreams(ctx); err != nil {
			logrus.Error(err)
		}
		cancel()
	}()

	window.Show()
	app.Exec()
}
