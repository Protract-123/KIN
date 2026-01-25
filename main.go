package main

import (
	"KIN/active_window"
	"KIN/kbhid"
	"KIN/volume"
	"log"
	"os"

	"github.com/mappu/miqt/qt6"
	"github.com/sstallion/go-hid"
)

var DefaultKeyboard = kbhid.KeyboardHIDInfo{
	VendorID:     0xFEED,
	ProductID:    0x4020,
	UsagePage:    0xFF60,
	Usage:        0x61,
	ReportLength: 32,
}

func main() {
	if err := hid.Init(); err != nil {
		log.Fatalf("hid init failed: %v", err)
	}
	defer hid.Exit()

	dev, err := kbhid.FindRawHIDDevice(DefaultKeyboard)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	qt6.NewQApplication(os.Args)
	tray := qt6.NewQSystemTrayIcon()

	icon := qt6.NewQIcon4("icon.png")
	tray.SetIcon(icon)
	tray.SetToolTip("KIN")

	menu := qt6.NewQMenu(nil)

	showAction := qt6.NewQAction2("Show Status")
	quitAction := qt6.NewQAction2("Quit")

	menu.AddAction(showAction)
	menu.AddSeparator()
	menu.AddAction(quitAction)

	tray.SetContextMenu(menu)
	tray.Show()

	showAction.OnTriggered(func() {
		tray.ShowMessage5(
			"Running",
			"Go application is active",
			qt6.QSystemTrayIcon__Information,
			2000,
		)
	})

	quitAction.OnTriggered(func() {
		tray.Hide()
		os.Exit(0)
	})

	go func() {
		for {
			active_window.SendActiveWindowData(dev, DefaultKeyboard)
		}
	}()
	go func() {
		for {
			volume.SendVolumeData(dev, DefaultKeyboard)
		}
	}()

	qt6.QApplication_Exec()

}
