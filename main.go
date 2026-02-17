package main

import (
	"KIN/app"
	"KIN/info/active_app"
	"KIN/info/volume"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

type InfoFunction func(config app.PayloadConfig, deviceNameToDevice map[string]*app.DeviceConfig)

var PayloadIDToInfoFunction = map[string]InfoFunction{
	"volume":     volume.SendVolumeData,
	"active_app": active_app.SendActiveWindowData,
}

const ConfigDirectory = "KIN"

var ApplicationConfig = app.ApplicationConfig{}

func main() {
	// Config Initialization
	base, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Unable to get user config directory: %v", err)
		shutdown()
	}

	configDir := filepath.Join(base, ConfigDirectory, "config.toml")

	err = app.InitializeConfigFile(configDir)
	if err != nil {
		log.Printf("Unable to initialize config: %v", err)
		shutdown()
	}

	err = app.LoadConfigFromFile(configDir, &ApplicationConfig)
	if err != nil {
		log.Printf("Unable to load config: %v", err)
		shutdown()
	}

	// HID Device Initialization
	for name, device := range ApplicationConfig.Devices {
		hidDevice, err := app.CreateHIDDevice(device)

		if err != nil {
			log.Printf("Unable to create HID device for %s: %v", name, err)
			continue
		}

		device.HIDDevice = hidDevice
		ApplicationConfig.Devices[name] = device
	}

	// Info Function Loops
	for payloadId, infoFunction := range PayloadIDToInfoFunction {
		go func() {
			devices := map[string]*app.DeviceConfig{}

			for name, device := range ApplicationConfig.Devices {
				for _, payload := range device.AuthorizedPayloads {
					if payloadId == payload {
						devices[name] = &device
					}
				}
			}

			for {
				infoFunction(ApplicationConfig.Payloads[payloadId], devices)
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
	for name := range ApplicationConfig.Devices {
		if ApplicationConfig.Devices[name].HIDDevice == nil {
			continue
		}

		err := ApplicationConfig.Devices[name].HIDDevice.Close()
		if err != nil {
			log.Printf("Failed to close device %s: %v", name, err)
		}
	}

	log.Print("Exited gracefully")
	os.Exit(0)
}
