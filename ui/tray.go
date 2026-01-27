package ui

import (
	"os"

	"github.com/mappu/miqt/qt6"
)

func CreateTray(win *qt6.QMainWindow) {
	tray := qt6.NewQSystemTrayIcon()

	icon := qt6.NewQIcon4("icon.png")
	tray.SetIcon(icon)
	tray.SetToolTip("KIN")

	menu := qt6.NewQMenu(nil)

	showAction := qt6.NewQAction2("Show Config")
	quitAction := qt6.NewQAction2("Quit")

	menu.AddAction(showAction)
	menu.AddSeparator()
	menu.AddAction(quitAction)

	tray.SetContextMenu(menu)
	tray.Show()

	showAction.OnTriggered(func() {
		win.Show()
		win.Raise()
		win.ActivateWindow()
	})

	quitAction.OnTriggered(func() {
		tray.Hide()
		os.Exit(0)
	})
}
