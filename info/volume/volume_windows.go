//go:build windows

package volume

import "C"
import (
	"strconv"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

var ole32 = windows.NewLazySystemDLL("ole32.dll").Handle()

var _CoCreateInstance func(*windows.GUID, uintptr, uint32, *windows.GUID, *unsafe.Pointer) uintptr = nil
var _CoInitialize func(uintptr) uintptr = nil
var _CoUninitialize func() = nil

var initOnce sync.Once

func initFunctions() {
	purego.RegisterLibFunc(&_CoCreateInstance, ole32, "CoCreateInstance")
	purego.RegisterLibFunc(&_CoInitialize, ole32, "CoInitialize")
	purego.RegisterLibFunc(&_CoUninitialize, ole32, "CoUninitialize")
}

func fetchVolume() string {
	initOnce.Do(initFunctions)

	var (
		deviceEnumerator *_IMMDeviceEnumerator
		defaultDevice    *_IMMDevice
		endpointVolume   *_IAudioEndpointVolume
		currentVolume    float32
	)

	exit := func() string {
		if endpointVolume != nil {
			purego.SyscallN(endpointVolume.lpVtbl.Release,
				uintptr(unsafe.Pointer(endpointVolume)))
		}
		if defaultDevice != nil {
			purego.SyscallN(defaultDevice.lpVtbl.Release,
				uintptr(unsafe.Pointer(defaultDevice)))
		}
		if deviceEnumerator != nil {
			purego.SyscallN(deviceEnumerator.lpVtbl.Release,
				uintptr(unsafe.Pointer(deviceEnumerator)))
		}

		_CoUninitialize()
		return strconv.Itoa(int(currentVolume * 100))
	}

	_CoInitialize(0)

	result := _CoCreateInstance(
		&_CLSID_MMDeviceEnumerator,
		0,
		windows.CLSCTX_INPROC_SERVER,
		&_IID_IMMDeviceEnumerator,
		(*unsafe.Pointer)(unsafe.Pointer(&deviceEnumerator)),
	)
	if int32(result) < 0 {
		return exit()
	}

	result, _, _ = purego.SyscallN(
		deviceEnumerator.lpVtbl.GetDefaultAudioEndpoint,
		uintptr(unsafe.Pointer(deviceEnumerator)),
		uintptr(0), // eRender
		uintptr(0), // eConsole
		uintptr(unsafe.Pointer(&defaultDevice)),
	)
	if int32(result) < 0 {
		return exit()
	}

	result, _, _ = purego.SyscallN(
		defaultDevice.lpVtbl.Activate,
		uintptr(unsafe.Pointer(defaultDevice)),
		uintptr(unsafe.Pointer(&_IID_IAudioEndpointVolume)),
		uintptr(windows.CLSCTX_INPROC_SERVER),
		0,
		uintptr(unsafe.Pointer(&endpointVolume)),
	)
	if int32(result) < 0 {
		return exit()
	}

	result, _, _ = purego.SyscallN(
		endpointVolume.lpVtbl.GetMasterVolumeLevelScalar,
		uintptr(unsafe.Pointer(endpointVolume)),
		uintptr(unsafe.Pointer(&currentVolume)),
	)

	return exit()
}

/*
 The following structs/GUIDS were created based on endpointvolume.h
	_CLSID_MMDeviceEnumerator GUID
	_IID_IMMDeviceEnumerator  GUID
	_IMMDeviceEnumerator      struct
	_IMMDeviceEnumeratorVtbl  struct
	_IMMDevice                struct
	_IMMDeviceVtbl            struct

*/

var _CLSID_MMDeviceEnumerator = windows.GUID{
	0xBCDE0395,
	0xE52F,
	0x467C,
	[8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E},
}

var _IID_IMMDeviceEnumerator = windows.GUID{
	0xA95664D2,
	0x9614,
	0x4F35,
	[8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6},
}

type _IMMDeviceEnumerator struct {
	lpVtbl *_IMMDeviceEnumeratorVtbl
}

type _IMMDeviceEnumeratorVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	EnumAudioEndpoints                     uintptr
	GetDefaultAudioEndpoint                uintptr
	GetDevice                              uintptr
	RegisterEndpointNotificationCallback   uintptr
	UnregisterEndpointNotificationCallback uintptr
}

type _IMMDevice struct {
	lpVtbl *_IMMDeviceVtbl
}

type _IMMDeviceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	Activate          uintptr
	OpenPropertyStore uintptr
	GetId             uintptr
	GetState          uintptr
}

/*
 The following structs/GUIDS were created based on endpointvolume.h
	_IID_IAudioEndpointVolume GUID
	_IAudioEndpointVolume     struct
	_IAudioEndpointVolumeVtbl struct
*/

var _IID_IAudioEndpointVolume = windows.GUID{
	0x5CDF2C82,
	0x841E,
	0x4546,
	[8]byte{0x97, 0x22, 0x0C, 0xF7, 0x40, 0x78, 0x22, 0x9A},
}

type _IAudioEndpointVolume struct {
	lpVtbl *_IAudioEndpointVolumeVtbl
}

type _IAudioEndpointVolumeVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	RegisterControlChangeNotify   uintptr
	UnregisterControlChangeNotify uintptr
	GetChannelCount               uintptr
	SetMasterVolumeLevel          uintptr
	SetMasterVolumeLevelScalar    uintptr
	GetMasterVolumeLevel          uintptr
	GetMasterVolumeLevelScalar    uintptr
	SetChannelVolumeLevel         uintptr
	SetChannelVolumeLevelScalar   uintptr
	GetChannelVolumeLevel         uintptr
	GetChannelVolumeLevelScalar   uintptr
	SetMute                       uintptr
	GetMute                       uintptr
	GetVolumeStepInfo             uintptr
	VolumeStepUp                  uintptr
	VolumeStepDown                uintptr
	QueryHardwareSupport          uintptr
	GetVolumeRange                uintptr
}
