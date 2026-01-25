package kbhid

import (
	"errors"

	"github.com/sstallion/go-hid"
)

type KeyboardHIDInfo struct {
	VendorID     uint16
	ProductID    uint16
	UsagePage    uint16
	Usage        uint16
	ReportLength int
}

// FindRawHIDDevice finds and opens the first matching RAW HID interface
func FindRawHIDDevice(cfg KeyboardHIDInfo) (*hid.Device, error) {
	var devicePath string

	err := hid.Enumerate(cfg.VendorID, cfg.ProductID,
		func(info *hid.DeviceInfo) error {
			if info.UsagePage == cfg.UsagePage &&
				info.Usage == cfg.Usage {
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
func SendRawReport(dev *hid.Device, cfg KeyboardHIDInfo, payload []byte) error {
	report := make([]byte, cfg.ReportLength+1) // Report ID + payload
	copy(report[1:], payload)

	_, err := dev.Write(report)
	return err
}
