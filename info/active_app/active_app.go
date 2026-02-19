package active_app

import (
	"KIN/app"
	"log"
	"strings"
	"time"
	"unicode"
)

func SendActiveWindowData(config app.PayloadConfig, deviceNameToDevice map[string]*app.DeviceConfig) {
	if !config.Enabled {
		return
	}

	applicationName := fetchActiveAppName()

	if applicationName != "" {
		for deviceName, device := range deviceNameToDevice {
			if device.HIDDevice == nil {
				continue
			}

			if !device.HIDDevice.IsOpen() {
				continue
			}

			data := app.StringToCString(applicationName, device.ReportLength-1) // First byte reserved for Payload Type
			err := app.SendPayload(device.HIDDevice, app.PayloadActiveApp, data, device.ReportLength)

			if err != nil {
				log.Printf("Write to device %s failed: %v", deviceName, err)
			}
		}
	}

	time.Sleep(config.RefreshRate)
}

func formatAppString(appName string) string {
	appName = strings.TrimSpace(appName)
	appName = strings.TrimSuffix(strings.ToLower(appName), ".exe")

	replacer := strings.NewReplacer(
		"-", " ",
		"_", " ",
	)
	appName = replacer.Replace(appName)

	words := strings.Fields(appName)

	for i, word := range words {
		runes := []rune(word)
		if len(runes) == 0 {
			continue
		}

		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}

		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}
