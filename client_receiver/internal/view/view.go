package view

import (
	"github.com/therecipe/qt/widgets"
)

func BuildWindowWithButton(title string, text string, clickHandler func(bool)) (*widgets.QMainWindow, *widgets.QPushButton) {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle(title)
	window.SetFixedSize2(300, 100)

	button := widgets.NewQPushButton2(text, nil)
	button.ConnectClicked(clickHandler)

	window.SetCentralWidget(button)

	return window, button
}
