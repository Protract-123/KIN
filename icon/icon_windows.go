//go:build windows

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "icon.ico"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}
