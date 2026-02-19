//go:build windows

package active_app

import (
	"log"
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

var user32 = windows.NewLazySystemDLL("user32.dll").Handle()
var psapi = windows.NewLazySystemDLL("psapi.dll").Handle()

var _GetForegroundWindow func() uintptr = nil
var _GetWindowThreadProcessId func(uintptr, *uint32) uint32 = nil
var _GetModuleBaseNameW func(uintptr, uintptr, *uint16, uint32) uint32 = nil

var initOnce sync.Once

func initFunctions() {
	purego.RegisterLibFunc(&_GetForegroundWindow, user32, "GetForegroundWindow")
	purego.RegisterLibFunc(&_GetWindowThreadProcessId, user32, "GetWindowThreadProcessId")
	purego.RegisterLibFunc(&_GetModuleBaseNameW, psapi, "GetModuleBaseNameW")
}

func fetchActiveAppName() string {
	initOnce.Do(initFunctions)
	if _GetForegroundWindow == nil || _GetWindowThreadProcessId == nil || _GetModuleBaseNameW == nil {
		return ""
	}

	foregroundWindowHandle := _GetForegroundWindow()
	if foregroundWindowHandle == 0 {
		return ""
	}

	var processId uint32
	threadId := _GetWindowThreadProcessId(
		foregroundWindowHandle,
		&processId,
	)
	if threadId == 0 || processId == 0 {
		return ""
	}

	processHandle, err := windows.OpenProcess(
		windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ,
		false,
		processId,
	)
	if err != nil {
		log.Println("OpenProcess error:", err)
		return ""
	}
	defer windows.CloseHandle(processHandle)

	exe := make([]uint16, windows.MAX_PATH)
	stringLength := _GetModuleBaseNameW(
		uintptr(processHandle),
		0,
		&exe[0],
		uint32(len(exe)),
	)
	if stringLength == 0 {
		return ""
	}

	processName := windows.UTF16ToString(exe[:stringLength])

	return stripBitnessSuffix(formatAppString(processName))
}

var bitnessSuffixes = []string{
	"32", "64",
	"x86", "x64",
	"amd64",
	"win32", "win64",
}

func stripBitnessSuffix(appName string) string {
	appName = strings.TrimSpace(appName)
	lower := strings.ToLower(appName)

	for _, suf := range bitnessSuffixes {
		if strings.HasSuffix(lower, suf) {
			cut := len(appName) - len(suf)
			if cut > 0 {
				prev := appName[cut-1]
				if prev == ' ' || prev == '-' || prev == '_' {
					return strings.TrimSpace(appName[:cut-1])
				}
			}
			return strings.TrimSpace(appName[:cut])
		}
	}
	return appName
}
