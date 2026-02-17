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

func formatAppString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(strings.ToLower(s), ".exe")

	// Replace common separators with spaces
	replacer := strings.NewReplacer(
		"-", " ",
		"_", " ",
	)
	s = replacer.Replace(s)

	// Split on whitespace
	words := strings.Fields(s)

	// Capitalize each word
	for i, w := range words {
		runes := []rune(w)
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
