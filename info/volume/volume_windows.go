//go:build windows

package volume

/*
#cgo LDFLAGS: -lole32 -luuid

float getMasterVolume(void);
*/
import "C"
import "strconv"

func FetchVolume() string {
	return strconv.Itoa(int(C.getMasterVolume() * 100))
}
