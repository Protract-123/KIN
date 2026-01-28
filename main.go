package main

import (
	"KIN/active_window"
	"KIN/app"
	"KIN/ui"
	"KIN/volume"
	"log"
	"os"

	"github.com/mappu/miqt/qt6"
	"github.com/sstallion/go-hid"
)

var InfoFunctions = []func(){
	volume.SendVolumeData,
	active_window.SendActiveWindowData,
}

func main() {
	if err := hid.Init(); err != nil {
		log.Fatalf("hid init failed: %v", err)
	}
	defer hid.Exit()

	qt6.NewQApplication(os.Args)
	configWindow := ui.NewConfigWindow()
	ui.CreateTray(configWindow)

	err := app.InitializeConfig()
	if err != nil {
		return
	}

	err = app.LoadConfig()
	if err != nil {
		return
	}

	app.BuildPayloadToKeyboards(app.ActiveConfig)
	app.BuildKeyboardToDevices(&app.ActiveConfig)

	for _, function := range InfoFunctions {
		go func() {
			for {
				function()
			}
		}()
	}

	qt6.QApplication_Exec()
}
