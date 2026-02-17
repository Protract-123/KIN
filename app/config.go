package app

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

var DefaultConfig = ApplicationConfig{
	Devices: map[string]DeviceConfig{
		"default": {
			VendorID:     0xFEED,
			ProductID:    0x4020,
			UsagePage:    0xFF60,
			Usage:        0x61,
			ReportLength: 32,

			AuthorizedPayloads: []string{"volume", "active_app"},
			HIDDevice:          nil,
		},
	},
	Payloads: map[string]PayloadConfig{
		"active_app": {
			RefreshRate: 1000 * time.Millisecond,
			Enabled:     true,
		},
		"volume": {
			RefreshRate: 200 * time.Millisecond,
			Enabled:     true,
		},
	},
}

func InitializeConfigFile(configFilePath string) error {
	configDir := filepath.Dir(configFilePath)

	err := os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)

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

func LoadConfigFromFile(configFilePath string, config *ApplicationConfig) error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
