package app

import (
	"errors"

	"rafaelmartins.com/p/usbhid"
)

func CreateHIDDevice(deviceConfig DeviceConfig) (*usbhid.Device, error) {
	deviceFilter := func(device *usbhid.Device) bool {
		if device.VendorId() != deviceConfig.VendorID.Value() {
			return false
		}
		if device.ProductId() != deviceConfig.ProductID.Value() {
			return false
		}
		if device.Usage() != deviceConfig.Usage.Value() {
			return false
		}
		if device.UsagePage() != deviceConfig.UsagePage.Value() {
			return false
		}

		return true
	}

	device, err := usbhid.Get(deviceFilter, true, false)

	if err != nil {
		return nil, err
	}

	return device, nil
}

func SendPayload(device *usbhid.Device, payloadType PayloadType, payload []byte, reportSize int) error {
	if !device.IsOpen() {
		return errors.New("USB device not open")
	}

	report := make([]byte, reportSize)

	report[0] = byte(payloadType)

	dataLength := reportSize - 1
	if len(payload) > dataLength {
		payload = payload[:dataLength]
	}

	copy(report[1:], payload)

	err := device.SetOutputReport(0x00, report)
	return err
}

func StringToCString(text string, maxLength int) []byte {
	data := []byte(text + "\x00")

	if len(data) > maxLength {
		data = data[:maxLength]
		data[maxLength-1] = 0
	}

	return data
}
