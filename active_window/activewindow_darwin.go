//go:build darwin

package active_window

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
#include <stdlib.h>

const char* active_app_name(void);
*/
import "C"
import "unsafe"

func FetchActiveWindowName() string {
	name := C.active_app_name()
	if name == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(name))
	return C.GoString(name)
}
