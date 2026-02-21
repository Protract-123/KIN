# KIN
Release Version: v1.0.0

KIN is a portable lightweight tray utility which sends realtime data from your system to USB HID devices.
It can be integrated with keyboards, macro pads and other HID to make them more intelligent.

An example of an HID device which uses KIN is my [hackpad](https://github.com/Protract-123/Hackpad/tree/main/firmware/qmk).
It uses KIN in order to have per application macros, and to display volume when changing volume
for finer adjustment

Currently supported payloads are:

| Payload    | Description                               | Default refresh rate |
|------------|-------------------------------------------|----------------------|
| volume     | System audio volume (0-100)               | 200ms                |
| active_app | Name of the currently focused application | 1000ms               |

## Installation

### Platform support

🟢 **Strong Support** - Payload supported through native APIs  
🟡 **Weak Support** - Payload supported through CLIs  
🔴 **No Support** - Payload not supported on platform

| OS      | Volume                | Active app      |
|---------|-----------------------|-----------------|
| macOS   | 🟢 CoreAudio          | 🟢 CoreGraphics |
| Windows | 🟢 Windows Core Audio | 🟢 Win32 API    |
| Linux   | 🟡 wpctl & pactl      | 🟡 hyprctl      |

### Install from release

Download the appropriate binary for your operating system from the [Releases](https://github.com/Protract-123/KIN/releases)
page, unzip the archive, and place the binary wherever you'd like.

### Install from source

The only prerequisite to install from source is Go. To install Go follow the directions found at
<https://go.dev/doc/install> or install through your preferred package manager.

```bash
git clone https://github.com/Protract-123/KIN.git
```

#### Windows

```bash
go build -ldflags="-s -w -H=windowsgui" -trimpath .
```

#### macOS/Linux

```bash
go build -ldflags="-s -w" -trimpath .
```

Running the commands above will output a `KIN` (or `KIN.exe`) binary in the repository directory, which
can then be moved anywhere.

## Usage

After running the binary, KIN will appear in the system tray and begin sending data to any configured devices
connected to the system.

From the tray menu you can :-

- See the connection status of all configured devices
- Open the config folder
- Quit KIN

Note: On Linux you may need to configure your system in order to allow information to be sent to the
configured HID devices.

## Configuration

KIN exposes a TOML configuration file which allows you to change application behavior. Its location
depends on your operating system, and can be directly opened through the tray item.

| OS      | PATH                                             |
|---------|--------------------------------------------------|
| Windows | `%AppData%\KIN\config.toml`                      |
| macOS   | `~/Library/Application\ Support/KIN/config.toml` |
| Linux   | `$XDG_CONFIG_HOME/KIN/config.toml`               |

### Sample Configuration

```toml
[devices]
  [devices.default]
    vendor_id = "0xFEED"
    product_id = "0x4020"
    usage_page = "0xFF60"
    usage = "0x61"
    report_length = 32
    authorized_payloads = ["volume", "active_app"]

[payloads]
  [payloads.active_app]
    refresh_rate = "1s"
    enabled = true
  [payloads.volume]
    refresh_rate = "200ms"
    enabled = true
```

### Devices

Each entry under `[devices]` represents a USB HID device that KIN will try to connect to. `[devices.default]`
represents a device which is named "default".

To find an HID device's information you can use "System Information" on macOS, "Device Manager" on Windows, and
"lsusb" on linux.

| Field                 | Description                               |
|-----------------------|-------------------------------------------|
| `vendor_id`           | USB vendor ID (hex)                       |
| `product_id`          | USB product ID (hex)                      |
| `usage_page`          | HID usage page (hex)                      |
| `usage`               | HID usage (hex)                           |
| `report_length`       | HID report length in bytes                |
| `authorized_payloads` | Which payloads this device should receive |

### Payloads

Each entry under `[payloads]` controls the behavior of a data source. `[payloads.active_app]` represents
the configuration for the payload `active_app`. A list of payloads can be found at the top of this README.

| Field          | Description                                     |
|----------------|-------------------------------------------------|
| `refresh_rate` | How often to send data (requires unit)          |
| `enabled`      | Set to `false` to disable this payload entirely |

## HID Report Format

KIN sends data using the following report structure:

```
Byte 0-2 : KIN identifier (0x10, 0xFF, 0x5B)
Byte 3   : Payload type (1 = active_app, 2 = volume)
Byte 4+  : Null-terminated UTF-8 string payload
```

## License

KIN is released under the [MIT License](LICENSE)  
Third party licenses can be found in the [Licenses](LICENSES) directory.
