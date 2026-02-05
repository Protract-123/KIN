package volume

import (
	"KIN/app"
	"log"
	"time"
)

const payloadKey = "volume"

func SendVolumeData() {
	payloadInfo := app.ActiveConfig.Payloads[payloadKey]
	if !payloadInfo.Active {
		return
	}

	volume := FetchVolume()

	if volume != "" {
		keyboardNames := app.PayloadToKeyboardNames[payloadKey]

		for _, name := range keyboardNames {
			keyboard := app.ActiveConfig.Keyboards[name]

			if keyboard.HIDDevice == nil {
				continue
			}

			data := app.PrepareCString(volume, keyboard.ReportLength-1) // First byte reserved for Payload Type
			payload := app.BuildPayload(app.PayloadVolume, data, keyboard.ReportLength)

			if err := app.SendHIDReport(keyboard.HIDDevice, keyboard, payload); err != nil {
				log.Printf("Write to keyboard %s failed: %v", name, err)
			}
		}
	}

	time.Sleep(time.Duration(payloadInfo.RefreshRate) * time.Millisecond)
}
