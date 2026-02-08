//go:build darwin

package activeapp

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	_CGWindowListCopyWindowInfo func(uint32, uint32) uintptr

	_CFArrayGetCount        func(uintptr) int64
	_CFArrayGetValueAtIndex func(uintptr, int64) unsafe.Pointer

	_CFDictionaryGetValue func(uintptr, unsafe.Pointer) unsafe.Pointer

	_CFNumberGetValue   func(uintptr, int32, unsafe.Pointer) bool
	_CFStringGetCString func(uintptr, unsafe.Pointer, int64, uint32) bool

	_CFRelease func(uintptr)

	_kCGWindowLayer     unsafe.Pointer
	_kCGWindowOwnerName unsafe.Pointer
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

	purego.RegisterLibFunc(&_CGWindowListCopyWindowInfo, coreGraphics, "CGWindowListCopyWindowInfo")

	purego.RegisterLibFunc(&_CFArrayGetCount, coreFoundation, "CFArrayGetCount")
	purego.RegisterLibFunc(&_CFArrayGetValueAtIndex, coreFoundation, "CFArrayGetValueAtIndex")
	purego.RegisterLibFunc(&_CFDictionaryGetValue, coreFoundation, "CFDictionaryGetValue")
	purego.RegisterLibFunc(&_CFNumberGetValue, coreFoundation, "CFNumberGetValue")
	purego.RegisterLibFunc(&_CFStringGetCString, coreFoundation, "CFStringGetCString")
	purego.RegisterLibFunc(&_CFRelease, coreFoundation, "CFRelease")

	sym, err := purego.Dlsym(coreGraphics, "kCGWindowLayer")
	if err != nil {
		panic(err)
	}
	_kCGWindowLayer = *(*unsafe.Pointer)(unsafe.Pointer(sym))

	sym, err = purego.Dlsym(coreGraphics, "kCGWindowOwnerName")
	if err != nil {
		panic(err)
	}
	_kCGWindowOwnerName = *(*unsafe.Pointer)(unsafe.Pointer(sym))
}

func fetchActiveAppName() string {
	initOnce.Do(initFunctions)

	windows := _CGWindowListCopyWindowInfo(
		_kCGWindowListOptionOnScreenOnly|
			_kCGWindowListExcludeDesktopElements,
		_kCGNullWindowID,
	)
	if windows == 0 {
		return ""
	}
	defer _CFRelease(windows)

	count := _CFArrayGetCount(windows)

	for i := int64(0); i < count; i++ {
		win := uintptr(
			_CFArrayGetValueAtIndex(windows, i),
		)

		// Layer check
		layerNum := uintptr(
			_CFDictionaryGetValue(win, _kCGWindowLayer),
		)
		if layerNum == 0 {
			continue
		}

		var layer int32
		if !_CFNumberGetValue(layerNum, _kCFNumberIntType, unsafe.Pointer(&layer)) {
			continue
		}
		if layer != 0 {
			continue
		}

		// Owner name
		owner := uintptr(
			_CFDictionaryGetValue(win, _kCGWindowOwnerName),
		)
		if owner == 0 {
			continue
		}

		buf := make([]byte, 256)
		if _CFStringGetCString(
			owner,
			unsafe.Pointer(&buf[0]),
			int64(len(buf)),
			_kCFStringEncodingUTF8,
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

/*
 The following constants were created based on CGWindow.h
	_kCGWindowListOptionOnScreenOnly     const
	_kCGWindowListExcludeDesktopElements const
	_kCGNullWindowID                     const
*/

const _kCGWindowListOptionOnScreenOnly = 1 << 0
const _kCGWindowListExcludeDesktopElements = 1 << 4
const _kCGNullWindowID = 0

/*
 The following constants were created based on CFNumber.h
	_kCFNumberIntType const
*/

const _kCFNumberIntType = 9

/*
 The following constants were created based on CFString.h
	_kCFStringEncodingUTF8 const
*/

const _kCFStringEncodingUTF8 = 0x08000100
