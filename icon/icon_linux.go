//go:build linux

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "png/tray_icon.png"
var TrayIcon []byte

func SetTrayIcon() {
	systray.SetIcon(TrayIcon)
}

//go:embed "png/cross_icon.png"
var CrossIcon []byte

//go:embed "png/tick_icon.png"
var TickIcon []byte

//go:embed "png/quit_icon.png"
var QuitIcon []byte

//go:embed "png/config_icon.png"
var ConfigIcon []byte
