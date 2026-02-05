#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

const char* active_app_name(void) {
    CFArrayRef windows = CGWindowListCopyWindowInfo(
        kCGWindowListOptionOnScreenOnly |
        kCGWindowListExcludeDesktopElements,
        kCGNullWindowID
    );

    if (!windows) return NULL;

    const char* result = NULL;
    CFIndex count = CFArrayGetCount(windows);

    for (CFIndex i = 0; i < count; i++) {
        CFDictionaryRef win =
            (CFDictionaryRef)CFArrayGetValueAtIndex(windows, i);

        // Only layer 0 windows (real app windows)
        int layer = 0;
        CFNumberRef layerNum =
            CFDictionaryGetValue(win, kCGWindowLayer);
        if (!layerNum ||
            !CFNumberGetValue(layerNum, kCFNumberIntType, &layer) ||
            layer != 0) {
            continue;
        }

        CFStringRef owner =
            CFDictionaryGetValue(win, kCGWindowOwnerName);
        if (!owner) continue;

        char buffer[256];
        if (CFStringGetCString(owner, buffer, sizeof(buffer),
                                kCFStringEncodingUTF8)) {
            result = strdup(buffer);
            break;
        }
    }

    CFRelease(windows);
    return result;
}
