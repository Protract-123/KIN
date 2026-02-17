package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"rafaelmartins.com/p/usbhid"
)

type PayloadType uint8

const (
	PayloadActiveApp PayloadType = iota + 1
	PayloadVolume
)

type PayloadConfig struct {
	RefreshRate time.Duration `toml:"refresh_rate"`
	Enabled     bool          `toml:"enabled"`
}

type DeviceConfig struct {
	VendorID     HexUint16 `toml:"vendor_id"`
	ProductID    HexUint16 `toml:"product_id"`
	UsagePage    HexUint16 `toml:"usage_page"`
	Usage        HexUint16 `toml:"usage"`
	ReportLength int       `toml:"report_length"`

	AuthorizedPayloads []string       `toml:"authorized_payloads"`
	HIDDevice          *usbhid.Device `toml:"-"`
}

type ApplicationConfig struct {
	Devices  map[string]DeviceConfig  `toml:"devices"`
	Payloads map[string]PayloadConfig `toml:"payloads"`
}

type HexUint16 uint16

func (h HexUint16) Value() uint16 {
	return uint16(h)
}

func (h HexUint16) String() string {
	return fmt.Sprintf("0x%X", uint16(h))
}

func (h HexUint16) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

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
