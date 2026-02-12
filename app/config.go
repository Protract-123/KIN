package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ConfigDirectory - The sub-folder in os.UserConfigDir which stores config.toml
const ConfigDirectory = "KIN"

// DefaultConfig - The default ApplicationConfig used to generate config.toml if it doesn't exist
var DefaultConfig = ApplicationConfig{
	Devices: map[string]DeviceConfig{
		"default": {
			VendorID:       0xFEED,
			ProductID:      0x4020,
			UsagePage:      0xFF60,
			Usage:          0x61,
			ReportLength:   32,
			ActivePayloads: []string{"volume", "active_app"},

			HIDDevice: nil,
		},
	},
	Payloads: map[string]PayloadConfig{
		"active_app": {
			RefreshRate: 1000,
			Active:      true,
		},
		"volume": {
			RefreshRate: 200,
			Active:      true,
		},
	},
}

// ActiveConfig - The loaded ApplicationConfig used by the application
var ActiveConfig = ApplicationConfig{}

// PayloadToDeviceNames - A map of payload IDs to device names
var PayloadToDeviceNames = map[string][]string{}

// InitializeConfigFile - Initializes config.toml in os.UserConfigDir / ConfigDirectory if it doesn't exist using DefaultConfig
func InitializeConfigFile() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(base, ConfigDirectory)
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

// LoadConfig - Initializes ActiveConfig from os.UserConfigDir / ConfigDirectory / config.toml
func LoadConfig() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(base, ConfigDirectory, "config.toml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(data, &ActiveConfig); err != nil {
		return err
	}

	return nil
}

// InitializePayloadToDeviceNames - Initializes PayloadToDeviceNames using ActiveConfig
func InitializePayloadToDeviceNames() {
	result := make(map[string][]string)

	for name, device := range ActiveConfig.Devices {
		for _, payload := range device.ActivePayloads {
			result[payload] = append(result[payload], name)
		}
	}

	PayloadToDeviceNames = result
}

// InitializeHIDDevices - Initializes all HID devices in ActiveConfig
func InitializeHIDDevices() {
	cfg := &ActiveConfig
	for name, device := range cfg.Devices {
		hidDevice, err := CreateHIDDevice(device)
		if err != nil {
			log.Printf("Unable to find HID device for device %s\n: %v", name, err)
			continue
		}
		device.HIDDevice = hidDevice
		cfg.Devices[name] = device
	}
}
