//go:build darwin

package icon

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed "icon.svg"
var TrayIcon []byte

//go:embed "icon.png"
var BackupTrayIcon []byte

func SetTrayIcon() {
	systray.SetTemplateIcon(TrayIcon, BackupTrayIcon)
}
