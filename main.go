package main

import (
	"KIN/active_window"
	"KIN/app"
	"KIN/volume"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
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
		systray.Quit()
	}()

	systray.Run(CreateTray, func() {})
	shutdown()
}

func CreateTray() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("KIN")
	systray.SetTooltip("Keyboard Information Negotiator")

	mQuit := systray.AddMenuItem("Quit", "Close KIN")
	mQuit.SetIcon(icon.Data)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
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
