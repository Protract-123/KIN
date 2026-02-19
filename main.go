package main

import (
	"KIN/app"
	"KIN/icon"
	"KIN/info/active_app"
	"KIN/info/volume"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"fyne.io/systray"
)

type InfoFunction func(config app.PayloadConfig, deviceNameToDevice map[string]*app.DeviceConfig)

var payloadIDToInfoFunction = map[string]InfoFunction{
	"volume":     volume.SendVolumeData,
	"active_app": active_app.SendActiveWindowData,
}

const configDirectory = "KIN"

var applicationConfig = app.ApplicationConfig{}

func main() {
	// Config Initialization
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Unable to get user config directory: %v", err)
		shutdown()
	}

	configDir := filepath.Join(userConfigDir, configDirectory, "config.toml")

	err = app.InitializeConfigFile(configDir)
	if err != nil {
		log.Printf("Unable to initialize config: %v", err)
		shutdown()
	}

	err = app.LoadConfigFromFile(configDir, &applicationConfig)
	if err != nil {
		log.Printf("Unable to load config: %v", err)
		shutdown()
	}

	// HID Device Initialization
	for name, device := range applicationConfig.Devices {
		hidDevice, err := app.CreateHIDDevice(device)

		if err != nil {
			log.Printf("Unable to create HID device for %s: %v", name, err)
			continue
		}

		device.HIDDevice = hidDevice
		applicationConfig.Devices[name] = device
	}

	// Info Function Loops
	for payloadId, infoFunction := range payloadIDToInfoFunction {
		go func() {
			devices := map[string]*app.DeviceConfig{}

			for name, device := range applicationConfig.Devices {
				for _, payload := range device.AuthorizedPayloads {
					if payloadId == payload {
						devices[name] = &device
					}
				}
			}

			for {
				infoFunction(applicationConfig.Payloads[payloadId], devices)
			}
		}()
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitSignal
		systray.Quit()
	}()

	log.Printf("Application started")
	systray.Run(createTray, func() {})
	shutdown()
}

func createTray() {
	icon.SetTrayIcon()
	systray.SetTooltip("Keyboard Information Negotiator")

	quitMenuItem := systray.AddMenuItem("Quit", "Close KIN")
	quitMenuItem.SetIcon(icon.TrayIcon)
	go func() {
		<-quitMenuItem.ClickedCh
		systray.Quit()
	}()
}

func shutdown() {
	for name := range applicationConfig.Devices {
		if applicationConfig.Devices[name].HIDDevice == nil {
			continue
		}

		err := applicationConfig.Devices[name].HIDDevice.Close()
		if err != nil {
			log.Printf("Failed to close device %s: %v", name, err)
		}
	}

	log.Print("Exited gracefully")
	os.Exit(0)
}
