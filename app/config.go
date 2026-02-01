package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ApplicationName = "KIN"

var DefaultConfig = ApplicationConfig{
	Keyboards: map[string]KeyboardConfig{
		"default": {
			VendorID:       0xFEED,
			ProductID:      0x4020,
			UsagePage:      0xFF60,
			Usage:          0x61,
			ReportLength:   32,
			ActivePayloads: []string{"volume", "active_window"},

			HIDDevice: nil,
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

var ActiveConfig = ApplicationConfig{}
var PayloadToKeyboardNames = map[string][]string{}

func InitializeConfigFile() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(base, ApplicationName)
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

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Error closing config file: %v", err)
		}
	}(configFile)

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

	configPath := filepath.Join(base, ApplicationName, "config.toml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(data, &ActiveConfig); err != nil {
		return err
	}

	return nil
}

func InitializePayloadToKeyboardNames(cfg ApplicationConfig) {
	result := make(map[string][]string)

	for name, kb := range cfg.Keyboards {
		for _, payload := range kb.ActivePayloads {
			result[payload] = append(result[payload], name)
		}
	}

	PayloadToKeyboardNames = result
}

func InitializeHIDDevices(cfg *ApplicationConfig) {
	for name, kb := range cfg.Keyboards {
		device, err := CreateHIDDevice(kb)
		if err != nil {
			log.Printf("Unable to find HID device for keyboard %s\n: %v", name, err)
			continue
		}
		kb.HIDDevice = device
		cfg.Keyboards[name] = kb
	}
}
