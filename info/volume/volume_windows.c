#include <windows.h>
#include <mmdeviceapi.h>
#include <endpointvolume.h>

const GUID CLSID_MMDeviceEnumerator_local = {0xBCDE0395, 0xE52F, 0x467C, {0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E}};
const GUID IID_IMMDeviceEnumerator_local = {0xA95664D2, 0x9614, 0x4F35, {0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6}};
const GUID IID_IAudioEndpointVolume_local = {0x5CDF2C82, 0x841E, 0x4546, {0x97, 0x22, 0x0C, 0xF7, 0x40, 0x78, 0x22, 0x9A}};

float getMasterVolume() {
    HRESULT hr;
    IMMDeviceEnumerator *deviceEnumerator = NULL;
    IMMDevice *defaultDevice = NULL;
    IAudioEndpointVolume *endpointVolume = NULL;
    float currentVolume = 0;

    CoInitialize(NULL);

    hr = CoCreateInstance(
        &CLSID_MMDeviceEnumerator_local,
        NULL,
        CLSCTX_INPROC_SERVER,
        &IID_IMMDeviceEnumerator_local,
        (LPVOID *)&deviceEnumerator);

    if (FAILED(hr)) goto cleanup;

    hr = deviceEnumerator->lpVtbl->GetDefaultAudioEndpoint(
        deviceEnumerator,
        eRender,
        eConsole,
        &defaultDevice);

    if (FAILED(hr)) goto cleanup;

    hr = defaultDevice->lpVtbl->Activate(
        defaultDevice,
        &IID_IAudioEndpointVolume_local,
        CLSCTX_INPROC_SERVER,
        NULL,
        (LPVOID *)&endpointVolume);

    if (FAILED(hr)) goto cleanup;

    hr = endpointVolume->lpVtbl->GetMasterVolumeLevelScalar(
        endpointVolume,
        &currentVolume);

cleanup:
    if (endpointVolume) endpointVolume->lpVtbl->Release(endpointVolume);
    if (defaultDevice) defaultDevice->lpVtbl->Release(defaultDevice);
    if (deviceEnumerator) deviceEnumerator->lpVtbl->Release(deviceEnumerator);
    CoUninitialize();

    return currentVolume;
}