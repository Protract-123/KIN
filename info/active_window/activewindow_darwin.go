//go:build darwin

package active_window

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	CGWindowListCopyWindowInfo func(uint32, uint32) uintptr

	CFArrayGetCount        func(uintptr) int64
	CFArrayGetValueAtIndex func(uintptr, int64) unsafe.Pointer

	CFDictionaryGetValue func(uintptr, unsafe.Pointer) unsafe.Pointer

	CFNumberGetValue   func(uintptr, int32, unsafe.Pointer) bool
	CFStringGetCString func(uintptr, unsafe.Pointer, int64, uint32) bool

	CFRelease func(uintptr)

	kCGWindowLayer     unsafe.Pointer
	kCGWindowOwnerName unsafe.Pointer
)

var initOnce sync.Once

func initFunctions() {
	coreGraphics, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics",
		purego.RTLD_NOW|purego.RTLD_GLOBAL,
	)
	if err != nil {
		panic(err)
	}

	coreFoundation, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation",
		purego.RTLD_NOW|purego.RTLD_GLOBAL,
	)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&CGWindowListCopyWindowInfo, coreGraphics, "CGWindowListCopyWindowInfo")

	purego.RegisterLibFunc(&CFArrayGetCount, coreFoundation, "CFArrayGetCount")
	purego.RegisterLibFunc(&CFArrayGetValueAtIndex, coreFoundation, "CFArrayGetValueAtIndex")
	purego.RegisterLibFunc(&CFDictionaryGetValue, coreFoundation, "CFDictionaryGetValue")
	purego.RegisterLibFunc(&CFNumberGetValue, coreFoundation, "CFNumberGetValue")
	purego.RegisterLibFunc(&CFStringGetCString, coreFoundation, "CFStringGetCString")
	purego.RegisterLibFunc(&CFRelease, coreFoundation, "CFRelease")

	sym, err := purego.Dlsym(coreGraphics, "kCGWindowLayer")
	if err != nil {
		panic(err)
	}
	kCGWindowLayer = *(*unsafe.Pointer)(unsafe.Pointer(sym))

	sym, err = purego.Dlsym(coreGraphics, "kCGWindowOwnerName")
	if err != nil {
		panic(err)
	}
	kCGWindowOwnerName = *(*unsafe.Pointer)(unsafe.Pointer(sym))
}

func FetchActiveWindowName() string {
	initOnce.Do(initFunctions)

	windows := CGWindowListCopyWindowInfo(
		kCGWindowListOptionOnScreenOnly|
			kCGWindowListExcludeDesktopElements,
		kCGNullWindowID,
	)
	if windows == 0 {
		return ""
	}
	defer CFRelease(windows)

	count := CFArrayGetCount(windows)

	for i := int64(0); i < count; i++ {
		win := uintptr(
			CFArrayGetValueAtIndex(windows, i),
		)

		// Layer check
		layerNum := uintptr(
			CFDictionaryGetValue(win, kCGWindowLayer),
		)
		if layerNum == 0 {
			continue
		}

		var layer int32
		if !CFNumberGetValue(layerNum, kCFNumberIntType, unsafe.Pointer(&layer)) {
			continue
		}
		if layer != 0 {
			continue
		}

		// Owner name
		owner := uintptr(
			CFDictionaryGetValue(win, kCGWindowOwnerName),
		)
		if owner == 0 {
			continue
		}

		buf := make([]byte, 256)
		if CFStringGetCString(
			owner,
			unsafe.Pointer(&buf[0]),
			int64(len(buf)),
			kCFStringEncodingUTF8,
		) {
			return string(buf[:clen(buf)])
		}
	}

	return ""
}

func clen(b []byte) int {
	for i, c := range b {
		if c == 0 {
			return i
		}
	}
	return len(b)
}

const kCGWindowListOptionOnScreenOnly = 1 << 0
const kCGWindowListExcludeDesktopElements = 1 << 4
const kCGNullWindowID = 0
const kCFNumberIntType = 9
const kCFStringEncodingUTF8 = 0x08000100
