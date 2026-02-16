package main

import (
	"KIN/app"
	"KIN/info/active_app"
	"KIN/info/volume"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

var InfoFunctions = []func(){
	volume.SendVolumeData,
	active_app.SendActiveWindowData,
}

func main() {
	err := app.InitializeConfigFile()
	if err != nil {
		log.Printf("Unable to initialize config: %v", err)
		shutdown()
	}

	err = app.LoadConfigFromFile()
	if err != nil {
		log.Printf("Unable to load config: %v", err)
		shutdown()
	}

	app.InitializePayloadToDeviceNames()
	app.InitializeHIDDevices()

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

	log.Printf("Application started")
	systray.Run(createTray, func() {})
	shutdown()
}

func createTray() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Keyboard Information Negotiator")

	mQuit := systray.AddMenuItem("Quit", "Close KIN")
	mQuit.SetIcon(icon.Data)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func shutdown() {
	for name := range app.ActiveConfig.Devices {
		if app.ActiveConfig.Devices[name].HIDDevice == nil {
			continue
		}

		err := app.ActiveConfig.Devices[name].HIDDevice.Close()
		if err != nil {
			log.Printf("Failed to close device %s: %v", name, err)
		}
	}

	log.Print("Exited gracefully")
	os.Exit(0)
}
