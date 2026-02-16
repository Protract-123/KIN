package app

import (
	"fmt"
	"strconv"
	"strings"

	"rafaelmartins.com/p/usbhid"
)

// PayloadType - A unique 8-bit integer which represents a type of payload
type PayloadType uint8

const (
	PayloadUnknown PayloadType = iota
	PayloadActiveApp
	PayloadVolume
)

type PayloadConfig struct {
	RefreshRate int  `toml:"refresh_rate"` // RefreshRate - How often data is refreshed in milliseconds
	Active      bool `toml:"active"`
}

type DeviceConfig struct {
	VendorID     HexUint16 `toml:"vendor_id"`
	ProductID    HexUint16 `toml:"product_id"`
	UsagePage    HexUint16 `toml:"usage_page"`
	Usage        HexUint16 `toml:"usage"`
	ReportLength int       `toml:"report_length"`

	ActivePayloads []string       `toml:"active_payloads"`
	HIDDevice      *usbhid.Device `toml:"-"`
}

type ApplicationConfig struct {
	Devices  map[string]DeviceConfig  `toml:"devices"`
	Payloads map[string]PayloadConfig `toml:"payloads"`
}

type HexUint16 uint16

func (h *HexUint16) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))

	base := 10
	if strings.HasPrefix(s, "0x") {
		base = 16
		s = s[2:]
	}

	v, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return err
	}

	*h = HexUint16(v)
	return nil
}

func (h HexUint16) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("0x%X", uint16(h))), nil
}

func (h HexUint16) GetUint16() uint16 {
	return uint16(h)
}
