package view

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"sync"
)

type View struct {
	Window *widgets.QMainWindow
	Button *widgets.QPushButton
	Input  *widgets.QLineEdit
	Output *widgets.QLabel
	Mu     sync.Mutex
}

func NewView(title, buttonText string) *View {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle(title)
	window.SetFixedSize2(300, 100)

	centralWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(centralWidget)

	layout := widgets.NewQVBoxLayout()

	input := widgets.NewQLineEdit(nil)
	layout.AddWidget(input, 0, 0)

	button := widgets.NewQPushButton2(buttonText, nil)
	layout.AddWidget(button, 0, 0)

	output := widgets.NewQLabel2("", nil, 0)
	layout.AddWidget(output, 0, 0)
	output.Hide()
	output.SetAlignment(core.Qt__AlignCenter)

	font := gui.NewQFont()
	font.SetPointSize(14)
	output.SetFont(font)

	centralWidget.SetLayout(layout)
	return &View{Window: window, Button: button, Input: input, Output: output}
}

func (v *View) DisableInputs() {
	v.switchInputs(false)
}

func (v *View) EnableInputs() {
	v.switchInputs(true)
}

func (v *View) switchInputs(enable bool) {
	v.Button.SetEnabled(enable)
	v.Input.SetEnabled(enable)
}

func (v *View) ShowError(title, description string) {
	widgets.QMessageBox_Critical(nil, title, description, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}
