//go:build windows

package active_window

import (
	"log"
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

var user32 = windows.NewLazySystemDLL("user32.dll").Handle()
var psapi = windows.NewLazySystemDLL("psapi.dll").Handle()

var getForegroundWindow func() uintptr = nil
var getWindowThreadProcessId func(uintptr, *uint32) uint32 = nil
var getModuleBaseNameW func(uintptr, uintptr, *uint16, uint32) uint32 = nil

var initOnce sync.Once

func initFunctions() {
	purego.RegisterLibFunc(&getForegroundWindow, user32, "GetForegroundWindow")
	purego.RegisterLibFunc(&getWindowThreadProcessId, user32, "GetWindowThreadProcessId")
	purego.RegisterLibFunc(&getModuleBaseNameW, psapi, "GetModuleBaseNameW")
}

func FetchActiveWindowName() string {
	initOnce.Do(initFunctions)
	if getForegroundWindow == nil || getWindowThreadProcessId == nil || getModuleBaseNameW == nil {
		return ""
	}

	foregroundWindowHandle := getForegroundWindow()
	if foregroundWindowHandle == 0 {
		return ""
	}

	var processId uint32
	threadId := getWindowThreadProcessId(
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
	stringLength := getModuleBaseNameW(
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

func stripBitnessSuffix(s string) string {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)

	for _, suf := range bitnessSuffixes {
		if strings.HasSuffix(lower, suf) {
			// ensure it's actually a suffix boundary
			cut := len(s) - len(suf)
			if cut > 0 {
				prev := s[cut-1]
				if prev == ' ' || prev == '-' || prev == '_' {
					return strings.TrimSpace(s[:cut-1])
				}
			}
			// exact match (e.g. "app64")
			return strings.TrimSpace(s[:cut])
		}
	}
	return s
}
