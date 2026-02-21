package main

import (
	"KIN/app"
	"KIN/icon"
	"KIN/info/active_app"
	"KIN/info/volume"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

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

	configFilePath := filepath.Join(userConfigDir, configDirectory, "config.toml")

	err = app.InitializeConfigFile(configFilePath)
	if err != nil {
		log.Printf("Unable to initialize config: %v", err)
		shutdown()
	}

	err = app.LoadConfigFromFile(configFilePath, &applicationConfig)
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

	var deviceMenuItems = map[string]*systray.MenuItem{}
	for deviceName := range applicationConfig.Devices {
		menuItem := systray.AddMenuItem(deviceName, fmt.Sprintf("%s not connected", deviceName))
		deviceMenuItems[deviceName] = menuItem
	}

	go func() {
		for deviceName, device := range applicationConfig.Devices {
			if menuItem, exists := deviceMenuItems[deviceName]; exists {
				tooltip := "not connected"
				statusText := "Disconnected"
				statusIcon := icon.CrossIcon

				if device.HIDDevice != nil && device.HIDDevice.IsOpen() {
					tooltip = "is connected"
					statusText = "Connected"
					statusIcon = icon.TickIcon
				}

				menuItem.SetTooltip(fmt.Sprintf("%s %s", deviceName, tooltip))
				if statusIcon != nil {
					menuItem.SetIcon(statusIcon)
				} else {
					menuItem.SetTitle(fmt.Sprintf("%s - %s", deviceName, statusText))
				}
			}
		}

		time.Sleep(1 * time.Second)
	}()

	systray.AddSeparator()

	openConfigMenuItem := systray.AddMenuItem("Open Config", "Opens the configuration file")
	if icon.ConfigIcon != nil {
		openConfigMenuItem.SetIcon(icon.ConfigIcon)
	}

	go func() {
		<-openConfigMenuItem.ClickedCh

		userConfigDir, err := os.UserConfigDir()

		if err != nil {
			log.Printf("Unable to get user config directory: %v", err)
		}

		configDir := filepath.Join(userConfigDir, configDirectory)

		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", configDir)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", configDir)
		case "linux":
			cmd = exec.Command("xdg-open", configDir)
		default:
			log.Printf("Unable to open config file on OS: %s", runtime.GOOS)
		}

		err = cmd.Start()
		if err != nil {
			log.Printf("Unable to open config file: %v", err)
		}
	}()

	quitMenuItem := systray.AddMenuItem("Quit", "Close KIN")
	if icon.QuitIcon != nil {
		quitMenuItem.SetIcon(icon.QuitIcon)
	}
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
