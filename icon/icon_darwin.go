//go:build darwin

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "svg/tray_icon.svg"
var TrayIcon []byte

//go:embed "png/tray_icon.png"
var BackupTrayIcon []byte

func SetTrayIcon() {
	systray.SetTemplateIcon(TrayIcon, BackupTrayIcon)
}

//go:embed "png/cross_icon.png"
var CrossIcon []byte

//go:embed "png/tick_icon.png"
var TickIcon []byte

//go:embed "png/quit_icon.png"
var QuitIcon []byte

//go:embed "png/config_icon.png"
var ConfigIcon []byte
