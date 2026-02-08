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

	volume := fetchVolume()

	if volume != "" {
		deviceNames := app.PayloadToDeviceNames[payloadKey]

		for _, name := range deviceNames {
			device := app.ActiveConfig.Devices[name]

			if device.HIDDevice == nil {
				continue
			}

			data := app.StringToCString(volume, device.ReportLength-1) // First byte reserved for Payload Type

			if err := app.SendPayload(device.HIDDevice, app.PayloadVolume, data, device.ReportLength); err != nil {
				log.Printf("Write to device %s failed: %v", name, err)
			}
		}
	}

	time.Sleep(time.Duration(payloadInfo.RefreshRate) * time.Millisecond)
}
