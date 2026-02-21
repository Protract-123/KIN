//go:build windows

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "ico/tray_icon.ico"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}

var CrossIcon []byte

var TickIcon []byte

var QuitIcon []byte

var ConfigIcon []byte
