package app

import (
	"errors"

	"rafaelmartins.com/p/usbhid"
)

func CreateHIDDevice(cfg DeviceConfig) (*usbhid.Device, error) {
	deviceFilter := func(device *usbhid.Device) bool {
		if device.VendorId() != cfg.VendorID.GetUint16() {
			return false
		}
		if device.ProductId() != cfg.ProductID.GetUint16() {
			return false
		}
		if device.Usage() != cfg.Usage.GetUint16() {
			return false
		}
		if device.UsagePage() != cfg.UsagePage.GetUint16() {
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

func SendPayload(dev *usbhid.Device, payloadType PayloadType, data []byte, reportLength int) error {
	if !dev.IsOpen() {
		return errors.New("USB device not open")
	}

	report := make([]byte, reportLength)

	// First byte is payload type
	report[0] = byte(payloadType)

	maxDataLen := reportLength - 1
	if len(data) > maxDataLen {
		data = data[:maxDataLen]
	}

	copy(report[1:], data)

	err := dev.SetOutputReport(0x00, report)
	return err
}

func StringToCString(s string, maxLen int) []byte {
	data := []byte(s + "\x00")

	if len(data) > maxLen {
		data = data[:maxLen]
		data[maxLen-1] = 0
	}

	return data
}
