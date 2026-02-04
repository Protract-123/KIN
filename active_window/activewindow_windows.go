//go:build windows

package active_window

import (
	"log"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	psapi                      = windows.NewLazySystemDLL("psapi.dll")
	procGetForegroundWindow    = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcess = user32.NewProc("GetWindowThreadProcessId")
	procGetModuleBaseNameW     = psapi.NewProc("GetModuleBaseNameW")
)

func FetchActiveWindowName() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}

	var pid uint32
	procGetWindowThreadProcess.Call(
		hwnd,
		uintptr(unsafe.Pointer(&pid)),
	)

	hProcess, err := windows.OpenProcess(
		windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ,
		false,
		pid,
	)
	if err != nil {
		log.Println("OpenProcess error:", err)
		return ""
	}
	defer windows.CloseHandle(hProcess)

	exe := make([]uint16, windows.MAX_PATH)
	procGetModuleBaseNameW.Call(
		uintptr(hProcess),
		0,
		uintptr(unsafe.Pointer(&exe[0])),
		uintptr(len(exe)),
	)
	processName := windows.UTF16ToString(exe)

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
