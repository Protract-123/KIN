package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sstallion/go-hid"
)

var AppName = "KIN"

var DefaultConfig = AppConfig{
	Keyboards: map[string]KeyboardConfig{
		"default": {
			VendorID:       0xFEED,
			ProductID:      0x4020,
			UsagePage:      0xFF60,
			Usage:          0x61,
			ReportLength:   32,
			ActivePayloads: []string{"volume"},
		},
	},
	Payloads: map[string]PayloadConfig{
		"active_window": {
			RefreshRate: 1000,
			Active:      true,
		},
		"volume": {
			RefreshRate: 200,
			Active:      true,
		},
	},
}

var ActiveConfig = AppConfig{}

var PayloadToKeyboards = map[string][]KeyboardConfig{}
var KeyboardToDevice = map[string]hid.Device{}

func InitializeConfig() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(base, AppName)
	configPath := filepath.Join(configDir, "config.toml")

	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)

	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer configFile.Close()

	marshal, err := toml.Marshal(DefaultConfig)
	if err != nil {
		return err
	}

	_, err = configFile.WriteString(string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(base, AppName, "config.toml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(data, &ActiveConfig); err != nil {
		return err
	}

	return nil
}

func BuildPayloadToKeyboards(cfg AppConfig) {
	result := make(map[string][]KeyboardConfig)

	for _, kb := range cfg.Keyboards {
		for _, payload := range kb.ActivePayloads {
			result[payload] = append(result[payload], kb)
		}
	}

	PayloadToKeyboards = result
}

func BuildKeyboardToDevices() {

}
