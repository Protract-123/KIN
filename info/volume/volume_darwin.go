//go:build darwin

package volume

import (
	"log"
	"strconv"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var audioObjectGetPropertyData func(uint32, *AudioObjectPropertyAddress, uint32, unsafe.Pointer, *uint32, unsafe.Pointer) int32 = nil
var audioObjectHasPropertyData func(uint32, *AudioObjectPropertyAddress) bool = nil

var initOnce sync.Once

func initFunctions() {
	coreAudio, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreAudio.framework/CoreAudio",
		purego.RTLD_NOW|purego.RTLD_GLOBAL,
	)

	if err != nil {
		log.Fatalf("Unable to load core audio: %v", err)
	}

	purego.RegisterLibFunc(&audioObjectGetPropertyData, coreAudio, "AudioObjectGetPropertyData")
	purego.RegisterLibFunc(&audioObjectHasPropertyData, coreAudio, "AudioObjectHasProperty")
	purego.RegisterLibFunc(&audioObjectGetPropertyData, coreAudio, "AudioObjectGetPropertyData")
}

func FetchVolume() string {
	initOnce.Do(initFunctions)

	var deviceID uint32
	dataSize := uint32(unsafe.Sizeof(deviceID))

	getOutputDeviceAddress := AudioObjectPropertyAddress{
		Selector: kAudioHardwarePropertyDefaultOutputDevice,
		Scope:    kAudioObjectPropertyScopeGlobal,
		Element:  kAudioObjectPropertyElementMain,
	}

	status := audioObjectGetPropertyData(
		kAudioObjectSystemObject,
		&getOutputDeviceAddress,
		0,
		nil,
		&dataSize,
		unsafe.Pointer(&deviceID),
	)

	if status != kAudioHardwareNoError {
		return ""
	}

	getVolumeAddress := AudioObjectPropertyAddress{
		Selector: kAudioHardwareServiceDeviceProperty_VirtualMainVolume,
		Scope:    kAudioObjectPropertyScopeOutput,
		Element:  kAudioObjectPropertyElementMain,
	}

	if !audioObjectHasPropertyData(deviceID, &getVolumeAddress) {
		return ""
	}

	var volume float32
	dataSize = uint32(unsafe.Sizeof(volume))

	status = audioObjectGetPropertyData(
		deviceID,
		&getVolumeAddress,
		0,
		nil,
		&dataSize,
		unsafe.Pointer(&volume),
	)

	if status != kAudioHardwareNoError {
		return ""
	}

	return strconv.Itoa(int(volume * 100))
}

/*
 The following structs/consts were found in AudioHardwareBase.h
	AudioObjectPropertyAddress      struct

	kAudioObjectPropertyScopeGlobal const
	kAudioObjectPropertyScopeOutput const
	kAudioObjectPropertyElementMain const
	kAudioHardwareNoError           const
*/

type AudioObjectPropertyAddress struct {
	Selector uint32
	Scope    uint32
	Element  uint32
}

const kAudioObjectPropertyScopeGlobal uint32 = 0x676C6F62 // 'glob'
const kAudioObjectPropertyScopeOutput uint32 = 0x6F757470 // 'outp'
const kAudioObjectPropertyElementMain uint32 = 0
const kAudioHardwareNoError int32 = 0

/*
 The following constants were found in AudioHardware.h
	kAudioObjectSystemObject                  const
	kAudioHardwarePropertyDefaultOutputDevice const
*/

const kAudioObjectSystemObject uint32 = 1
const kAudioHardwarePropertyDefaultOutputDevice uint32 = 0x644F7574

/*
 The following constants were found in AudioHardwareService.h
	kAudioHardwareServiceDeviceProperty_VirtualMainVolume const
*/

const kAudioHardwareServiceDeviceProperty_VirtualMainVolume uint32 = 0x766D7663
