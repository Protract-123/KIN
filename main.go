package main

import (
	"KIN/active_window"
	"KIN/app"
	"KIN/ui"
	"KIN/volume"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mappu/miqt/qt6"
	"github.com/sstallion/go-hid"
)

var InfoFunctions = []func(){
	volume.SendVolumeData,
	active_window.SendActiveWindowData,
}

func main() {
	if err := hid.Init(); err != nil {
		log.Fatalf("HID init failed: %v", err)
	}

	qt6.NewQApplication(os.Args)
	configWindow := ui.NewConfigWindow()
	ui.CreateTray(configWindow)

	err := app.InitializeConfigFile()
	if err != nil {
		log.Printf("Unable to initialize config: %v", err)
		shutdown()
	}

	err = app.LoadConfig()
	if err != nil {
		log.Printf("Unable to load config: %v", err)
		shutdown()
	}

	app.InitializePayloadToKeyboardNames(app.ActiveConfig)
	app.InitializeHIDDevices(&app.ActiveConfig)

	for _, function := range InfoFunctions {
		go func() {
			for {
				function()
			}
		}()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		qt6.QCoreApplication_Exit()
	}()

	qt6.QApplication_Exec()

	shutdown()
}

func shutdown() {
	for name := range app.ActiveConfig.Keyboards {
		if app.ActiveConfig.Keyboards[name].HIDDevice == nil {
			continue
		}

		err := app.ActiveConfig.Keyboards[name].HIDDevice.Close()
		if err != nil {
			log.Printf("Failed to close keyboard %s: %v", name, err)
		}
	}

	err := hid.Exit()
	if err != nil {
		log.Printf("HID exit failed: %v", err)
	}

	log.Print("Exited gracefully")
	os.Exit(0)
}
