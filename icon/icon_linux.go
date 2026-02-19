//go:build linux

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "icon.png"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}
