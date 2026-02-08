package app

import (
	"errors"

	"github.com/sstallion/go-hid"
)

// CreateHIDDevice - Creates a reference to a hid.Device using a DeviceConfig to find an exact match
func CreateHIDDevice(cfg DeviceConfig) (*hid.Device, error) {
	var devicePath string

	err := hid.Enumerate(cfg.VendorID.GetUint16(), cfg.ProductID.GetUint16(),
		func(info *hid.DeviceInfo) error {
			if info.UsagePage == cfg.UsagePage.GetUint16() &&
				info.Usage == cfg.Usage.GetUint16() {
				devicePath = info.Path
				// Stop enumeration early
				return errors.New("device found")
			}
			return nil
		},
	)

	if devicePath == "" {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("HID device not found")
	}

	return hid.OpenPath(devicePath)
}

// SendPayload - Sends data to a hid.Device given a PayloadType, and the reportLength of the device
func SendPayload(dev *hid.Device, payloadType PayloadType, data []byte, reportLength int) error {
	// Report = ReportID (1 byte) + payload
	report := make([]byte, reportLength+1)

	// Payload starts after report ID
	payload := report[1:]

	// First byte is payload type
	payload[0] = byte(payloadType)

	maxDataLen := reportLength - 1
	if len(data) > maxDataLen {
		data = data[:maxDataLen]
	}

	copy(payload[1:], data)

	_, err := dev.Write(report)
	return err
}

// StringToCString - Converts a Go string to a C string which fits in maxLen bytes
func StringToCString(s string, maxLen int) []byte {
	data := []byte(s + "\x00")

	if len(data) > maxLen {
		data = data[:maxLen]
		data[maxLen-1] = 0
	}

	return data
}
