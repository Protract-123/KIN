//go:build linux

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "export/tray_icon.png"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}
