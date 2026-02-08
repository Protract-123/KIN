//go:build darwin

package volume

import (
	"log"
	"math"
	"strconv"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var _AudioObjectGetPropertyData func(uint32, *_AudioObjectPropertyAddress, uint32, unsafe.Pointer, *uint32, unsafe.Pointer) int32 = nil
var _AudioObjectHasPropertyData func(uint32, *_AudioObjectPropertyAddress) bool = nil

var initOnce sync.Once

func initFunctions() {
	coreAudio, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreAudio.framework/CoreAudio",
		purego.RTLD_NOW|purego.RTLD_GLOBAL,
	)

	if err != nil {
		log.Fatalf("Unable to load core audio: %v", err)
	}

	purego.RegisterLibFunc(&_AudioObjectGetPropertyData, coreAudio, "AudioObjectGetPropertyData")
	purego.RegisterLibFunc(&_AudioObjectHasPropertyData, coreAudio, "AudioObjectHasProperty")
}

func fetchVolume() string {
	initOnce.Do(initFunctions)

	var deviceID uint32
	dataSize := uint32(unsafe.Sizeof(deviceID))

	getOutputDeviceAddress := _AudioObjectPropertyAddress{
		Selector: _kAudioHardwarePropertyDefaultOutputDevice,
		Scope:    _kAudioObjectPropertyScopeGlobal,
		Element:  _kAudioObjectPropertyElementMain,
	}

	status := _AudioObjectGetPropertyData(
		_kAudioObjectSystemObject,
		&getOutputDeviceAddress,
		0,
		nil,
		&dataSize,
		unsafe.Pointer(&deviceID),
	)

	if status != _kAudioHardwareNoError {
		return ""
	}

	getVolumeAddress := _AudioObjectPropertyAddress{
		Selector: _kAudioHardwareServiceDeviceProperty_VirtualMainVolume,
		Scope:    _kAudioObjectPropertyScopeOutput,
		Element:  _kAudioObjectPropertyElementMain,
	}

	if !_AudioObjectHasPropertyData(deviceID, &getVolumeAddress) {
		return ""
	}

	var volume float32
	dataSize = uint32(unsafe.Sizeof(volume))

	status = _AudioObjectGetPropertyData(
		deviceID,
		&getVolumeAddress,
		0,
		nil,
		&dataSize,
		unsafe.Pointer(&volume),
	)

	if status != _kAudioHardwareNoError {
		return ""
	}

	return strconv.Itoa(int(math.Round(float64(volume * 100))))
}

/*
 The following structs/consts were created based on AudioHardwareBase.h
	_AudioObjectPropertyAddress      struct

	_kAudioObjectPropertyScopeGlobal const
	_kAudioObjectPropertyScopeOutput const
	_kAudioObjectPropertyElementMain const
	_kAudioHardwareNoError           const
*/

type _AudioObjectPropertyAddress struct {
	Selector uint32
	Scope    uint32
	Element  uint32
}

const _kAudioObjectPropertyScopeGlobal uint32 = 0x676C6F62 // 'glob'
const _kAudioObjectPropertyScopeOutput uint32 = 0x6F757470 // 'outp'
const _kAudioObjectPropertyElementMain uint32 = 0
const _kAudioHardwareNoError int32 = 0

/*
 The following constants were created based on AudioHardware.h
	_kAudioObjectSystemObject                  const
	_kAudioHardwarePropertyDefaultOutputDevice const
*/

const _kAudioObjectSystemObject uint32 = 1
const _kAudioHardwarePropertyDefaultOutputDevice uint32 = 0x644F7574

/*
 The following constants were created based on AudioHardwareService.h
	_kAudioHardwareServiceDeviceProperty_VirtualMainVolume const
*/

const _kAudioHardwareServiceDeviceProperty_VirtualMainVolume uint32 = 0x766D7663
