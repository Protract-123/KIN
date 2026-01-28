package volume

import (
	"KIN/app"
	"fmt"
	"log"
	"time"
)

const payloadKey = "volume"

func SendVolumeData() {
	payloadInfo := app.ActiveConfig.Payloads[payloadKey]
	if !payloadInfo.Active {
		return
	}

	volume, _ := FetchVolume()

	if volume != "" {
		keyboards := app.PayloadToKeyboards[payloadKey]

		for _, name := range keyboards {
			keyboard := app.ActiveConfig.Keyboards[name]

			data := app.PrepareCString(volume, keyboard.ReportLength-1) // First byte reserved for Payload Type
			payload := app.BuildPayload(app.PayloadVolume, data, keyboard.ReportLength)

			if keyboard.HIDDevice == nil {
				fmt.Printf("Keyboard %s not initialized\n", name)
				continue
			}

			if err := app.SendRawReport(keyboard.HIDDevice, keyboard, payload); err != nil {
				log.Printf("Write to keyboard %s failed: %v", name, err)
			}
		}
	}

	time.Sleep(time.Duration(payloadInfo.RefreshRate) * time.Millisecond)
}
