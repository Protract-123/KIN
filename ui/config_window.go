package ui

import (
	"github.com/mappu/miqt/qt6"
)

func NewConfigWindow() *qt6.QMainWindow {
	win := qt6.NewQMainWindow2()
	win.SetWindowTitle("KIN Configuration")
	win.Resize(500, 400)

	// Central widget
	central := qt6.NewQWidget2()
	layout := qt6.NewQVBoxLayout2()

	// Example controls
	label := qt6.NewQLabel3("Information")
	checkbox := qt6.NewQCheckBox3("Send volume data")
	saveBtn := qt6.NewQPushButton3("Send active application data")

	layout.AddWidget(label.QWidget)
	layout.AddWidget(checkbox.QWidget)
	layout.AddStretchWithStretch(1)
	layout.AddWidget(saveBtn.QWidget)

	central.SetLayout(layout.QLayout)
	win.SetCentralWidget(central)

	saveBtn.OnClicked(func() {
		win.Hide()
	})

	win.SetAttribute2(qt6.WA_DeleteOnClose, false)
	win.OnCloseEvent(func(super func(e *qt6.QCloseEvent), event *qt6.QCloseEvent) {
		event.Ignore()
		win.Hide()
	})

	return win
}
