//go:build darwin

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "export/tray_icon.svg"
var TrayIcon []byte

//go:embed "export/tray_icon.png"
var BackupTrayIcon []byte

func SetTrayIcon() {
	systray.SetTemplateIcon(TrayIcon, BackupTrayIcon)
}
