package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sstallion/go-hid"
)

type PayloadConfig struct {
	RefreshRate int  `toml:"refresh_rate"`
	Active      bool `toml:"active"`
}

type KeyboardConfig struct {
	VendorID     HexUint16 `toml:"vendor_id"`
	ProductID    HexUint16 `toml:"product_id"`
	UsagePage    HexUint16 `toml:"usage_page"`
	Usage        HexUint16 `toml:"usage"`
	ReportLength int       `toml:"report_length"`

	ActivePayloads []string    `toml:"active_payloads"`
	HIDDevice      *hid.Device `toml:"-"`
}

type ApplicationConfig struct {
	Keyboards map[string]KeyboardConfig `toml:"keyboards"`
	Payloads  map[string]PayloadConfig  `toml:"payloads"`
}

type HexUint16 uint16

func (u HexUint16) GetUint16() uint16 {
	return uint16(u)
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

func (h HexUint16) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("0x%X", uint16(h))), nil
}
