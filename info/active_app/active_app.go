package active_app

import (
	"KIN/app"
	"log"
	"strings"
	"time"
	"unicode"
)

const payloadKey = "active_app"

func SendActiveWindowData() {
	payloadInfo := app.ActiveConfig.Payloads[payloadKey]
	if !payloadInfo.Active {
		return
	}

	applicationName := fetchActiveAppName()

	if applicationName != "" {
		deviceNames := app.PayloadToDeviceNames[payloadKey]

		for _, name := range deviceNames {
			device := app.ActiveConfig.Devices[name]

			if device.HIDDevice == nil {
				continue
			}

			data := app.StringToCString(applicationName, device.ReportLength-1) // First byte reserved for Payload Type

			if err := app.SendPayload(device.HIDDevice, app.PayloadActiveApp, data, device.ReportLength); err != nil {
				log.Printf("Write to device %s failed: %v", name, err)
			}
		}
	}

	time.Sleep(time.Duration(payloadInfo.RefreshRate) * time.Millisecond)
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
