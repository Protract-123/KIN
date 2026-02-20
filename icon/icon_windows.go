//go:build windows

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "export/tray_icon.ico"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}
