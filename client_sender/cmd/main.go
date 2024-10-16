package main

import (
	"client_sender/internal/service"
	"client_sender/internal/view"
	"client_sender/pkg"
	"context"
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	v := view.NewView("Тестовое задание", "Запустить")
	poolSender := service.NewPoolSender("localhost:8081")

	v.Button.ConnectClicked(func(bool) {
		poolCount, err := strconv.Atoi(v.Input.Text())
		if err != nil {
			v.ShowError("Ошибка", "Ошибка ввода количества потоков")
			return
		}

		go func() {
			poolSender.Send(poolCount)
		}()
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for bytesCount := range poolSender.GetBytesCountSentChan() {
			v.Mu.Lock()
			v.DisableInputs()
			v.Output.Show()
			v.Output.SetText(strconv.FormatInt(bytesCount, 10) + " bytes sent")
			v.Mu.Unlock()

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		throttler := pkg.NewThrottler(time.Millisecond * 200)
		for err := range poolSender.GetErrorChan() {
			throttler.Call(func() { v.ShowError("Ошибка", err.Error()) })
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

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(ctx, time.Second*5)
		defer cancelShutdownCtx()

		if err := poolSender.Shutdown(shutdownCtx); err != nil {
			log.Printf("error: %v", err)
		}

		cancel()
		v.Mu.Lock()
		v.Window.Close()
		v.Mu.Unlock()
	}()

	v.Window.Show()
	app.Exec()
}
