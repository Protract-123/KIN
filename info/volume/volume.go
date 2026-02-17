package volume

import (
	"KIN/app"
	"log"
	"time"
)

func SendVolumeData(config app.PayloadConfig, deviceNameToDevice map[string]*app.DeviceConfig) {
	if !config.Enabled {
		return
	}

	volume := fetchVolume()

	if volume != "" {
		for deviceName, device := range deviceNameToDevice {
			if device.HIDDevice == nil {
				continue
			}

			if !device.HIDDevice.IsOpen() {
				continue
			}

			data := app.StringToCString(volume, device.ReportLength-1) // First byte reserved for Payload Type
			err := app.SendPayload(device.HIDDevice, app.PayloadVolume, data, device.ReportLength)

			if err != nil {
				log.Printf("Write to device %s failed: %v", deviceName, err)
			}
		}
	}

	time.Sleep(config.RefreshRate)
}
