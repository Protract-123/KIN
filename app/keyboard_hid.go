package app

import (
	"errors"

	"github.com/sstallion/go-hid"
)

// FindRawHIDDevice finds and opens the first matching RAW HID interface
func FindRawHIDDevice(cfg KeyboardConfig) (*hid.Device, error) {
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
		return nil, errors.New("RAW HID device not found")
	}

	return hid.OpenPath(devicePath)
}

// SendRawReport sends a RAW HID report using the provided config
func SendRawReport(dev *hid.Device, cfg KeyboardConfig, payload []byte) error {
	report := make([]byte, cfg.ReportLength+1) // Report ID + payload
	copy(report[1:], payload)

	_, err := dev.Write(report)
	return err
}
