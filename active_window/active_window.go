package active_window

import (
	"KIN/app"
	"log"
	"strings"
	"time"
	"unicode"
)

const payloadKey = "active_window"

func SendActiveWindowData() {
	payloadInfo := app.ActiveConfig.Payloads[payloadKey]
	if !payloadInfo.Active {
		return
	}

	window := FetchActiveWindowName()

	if window != "" {
		keyboardNames := app.PayloadToKeyboardNames[payloadKey]

		for _, name := range keyboardNames {
			keyboard := app.ActiveConfig.Keyboards[name]

			if keyboard.HIDDevice == nil {
				continue
			}

			data := app.PrepareCString(window, keyboard.ReportLength-1) // First byte reserved for Payload Type
			payload := app.BuildPayload(app.PayloadActiveWindow, data, keyboard.ReportLength)

			if err := app.SendHIDReport(keyboard.HIDDevice, keyboard, payload); err != nil {
				log.Printf("Write to keyboard %s failed: %v", name, err)
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
